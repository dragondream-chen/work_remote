package loadbalancer

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	TaintPriority = 1e15
)

type ServerState struct {
	Host          string
	Port          int
	URL           string
	ActiveTokens  int64
	ActiveKVCache int64
	Healthy       bool
	Tainted       bool
	Weight        int

	mu sync.RWMutex
}

func NewServerState(host string, port int, weight int) *ServerState {
	return &ServerState{
		Host:    host,
		Port:    port,
		URL:     fmt.Sprintf("http://%s:%d/v1", host, port),
		Healthy: true,
		Tainted: false,
		Weight:  weight,
	}
}

func (s *ServerState) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s *ServerState) CalculatePrefillPriority() float64 {
	return float64(atomic.LoadInt64(&s.ActiveTokens)) + float64(atomic.LoadInt64(&s.ActiveKVCache))*0.3
}

func (s *ServerState) CalculateDecodePriority() float64 {
	return float64(atomic.LoadInt64(&s.ActiveTokens))
}

func (s *ServerState) AddTokens(count int64) {
	atomic.AddInt64(&s.ActiveTokens, count)
}

func (s *ServerState) AddKVCache(count int64) {
	atomic.AddInt64(&s.ActiveKVCache, count)
}

func (s *ServerState) ReleaseTokens(count int64) {
	atomic.AddInt64(&s.ActiveTokens, -count)
}

func (s *ServerState) ReleaseKVCache(count int64) {
	for {
		current := atomic.LoadInt64(&s.ActiveKVCache)
		if current <= 0 {
			return
		}
		newVal := current - count
		if newVal < 0 {
			newVal = 0
		}
		if atomic.CompareAndSwapInt64(&s.ActiveKVCache, current, newVal) {
			return
		}
	}
}

type PriorityItem struct {
	Priority float64
	Index    int
	Server   *ServerState
}

type PriorityQueue []*PriorityItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*PriorityItem)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

type ServerPool struct {
	servers []*ServerState
	heap    PriorityQueue
	index   map[*ServerState]int
	mu      sync.RWMutex
}

func NewServerPool(servers []*ServerState) *ServerPool {
	pool := &ServerPool{
		servers: servers,
		heap:    make(PriorityQueue, len(servers)),
		index:   make(map[*ServerState]int),
	}

	for i, server := range servers {
		pool.heap[i] = &PriorityItem{
			Priority: 0,
			Index:    i,
			Server:   server,
		}
		pool.index[server] = i
	}

	return pool
}

func (p *ServerPool) Select(calculatePriority func(*ServerState) float64) (*ServerState, int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.heap) == 0 {
		return nil, -1
	}

	minIdx := 0
	for i := 1; i < len(p.heap); i++ {
		if p.heap[i].Priority < p.heap[minIdx].Priority {
			minIdx = i
		}
	}

	item := p.heap[minIdx]
	p.heap = append(p.heap[:minIdx], p.heap[minIdx+1:]...)

	return item.Server, item.Index
}

func (p *ServerPool) UpdatePriority(idx int, priority float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if idx < 0 || idx >= len(p.servers) {
		return
	}

	server := p.servers[idx]
	for i, item := range p.heap {
		if item.Index == idx {
			p.heap[i].Priority = priority
			break
		}
	}
}

func (p *ServerPool) AddServer(server *ServerState) {
	p.mu.Lock()
	defer p.mu.Unlock()

	idx := len(p.servers)
	p.servers = append(p.servers, server)
	p.index[server] = idx

	p.heap = append(p.heap, &PriorityItem{
		Priority: 0,
		Index:    idx,
		Server:   server,
	})
}

func (p *ServerPool) RemoveServer(server *ServerState) {
	p.mu.Lock()
	defer p.mu.Unlock()

	idx, exists := p.index[server]
	if !exists {
		return
	}

	delete(p.index, server)
	p.servers = append(p.servers[:idx], p.servers[idx+1:]...)

	for i := idx; i < len(p.servers); i++ {
		p.index[p.servers[i]] = i
	}

	p.rebuildHeap()
}

func (p *ServerPool) rebuildHeap() {
	p.heap = make(PriorityQueue, len(p.servers))
	for i, server := range p.servers {
		p.heap[i] = &PriorityItem{
			Priority: server.CalculatePrefillPriority(),
			Index:    i,
			Server:   server,
		}
	}
}

func (p *ServerPool) TaintServer(server *ServerState) {
	p.mu.Lock()
	defer p.mu.Unlock()

	server.Tainted = true
	p.rebuildHeap()
}

func (p *ServerPool) GetServers() []*ServerState {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make([]*ServerState, len(p.servers))
	copy(result, p.servers)
	return result
}

func (p *ServerPool) Len() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.servers)
}

type LoadBalancer struct {
	prefillerPool *ServerPool
	decoderPool   *ServerPool
}

func NewLoadBalancer(prefillers, decoders []*ServerState) *LoadBalancer {
	return &LoadBalancer{
		prefillerPool: NewServerPool(prefillers),
		decoderPool:   NewServerPool(decoders),
	}
}

func (lb *LoadBalancer) SelectPrefiller(score float64) (*ServerState, int) {
	server, idx := lb.prefillerPool.Select((*ServerState).CalculatePrefillPriority)
	if server != nil {
		server.AddTokens(int64(score))
		server.AddKVCache(int64(score))
		lb.prefillerPool.UpdatePriority(idx, server.CalculatePrefillPriority())
	}
	return server, idx
}

func (lb *LoadBalancer) ReleasePrefiller(idx int, score float64) {
	if idx < 0 || idx >= lb.prefillerPool.Len() {
		return
	}

	servers := lb.prefillerPool.GetServers()
	if idx >= len(servers) {
		return
	}

	server := servers[idx]
	server.ReleaseTokens(int64(score))
	lb.prefillerPool.UpdatePriority(idx, server.CalculatePrefillPriority())
}

func (lb *LoadBalancer) ReleasePrefillerKV(idx int, score float64) {
	if idx < 0 || idx >= lb.prefillerPool.Len() {
		return
	}

	servers := lb.prefillerPool.GetServers()
	if idx >= len(servers) {
		return
	}

	server := servers[idx]
	server.ReleaseKVCache(int64(score))
	lb.prefillerPool.UpdatePriority(idx, server.CalculatePrefillPriority())
}

func (lb *LoadBalancer) SelectDecoder(score float64) (*ServerState, int) {
	server, idx := lb.decoderPool.Select((*ServerState).CalculateDecodePriority)
	if server != nil {
		server.AddTokens(int64(score))
		lb.decoderPool.UpdatePriority(idx, server.CalculateDecodePriority())
	}
	return server, idx
}

func (lb *LoadBalancer) ReleaseDecoder(idx int, score float64) {
	if idx < 0 || idx >= lb.decoderPool.Len() {
		return
	}

	servers := lb.decoderPool.GetServers()
	if idx >= len(servers) {
		return
	}

	server := servers[idx]
	server.ReleaseTokens(int64(score))
	lb.decoderPool.UpdatePriority(idx, server.CalculateDecodePriority())
}

func (lb *LoadBalancer) AddPrefiller(server *ServerState) {
	lb.prefillerPool.AddServer(server)
}

func (lb *LoadBalancer) AddDecoder(server *ServerState) {
	lb.decoderPool.AddServer(server)
}

func (lb *LoadBalancer) RemovePrefiller(server *ServerState) {
	lb.prefillerPool.RemoveServer(server)
}

func (lb *LoadBalancer) RemoveDecoder(server *ServerState) {
	lb.decoderPool.RemoveServer(server)
}

func (lb *LoadBalancer) GetPrefillers() []*ServerState {
	return lb.prefillerPool.GetServers()
}

func (lb *LoadBalancer) GetDecoders() []*ServerState {
	return lb.decoderPool.GetServers()
}

func (lb *LoadBalancer) PrefillerCount() int {
	return lb.prefillerPool.Len()
}

func (lb *LoadBalancer) DecoderCount() int {
	return lb.decoderPool.Len()
}

func CalculatePrefillScore(requestLength int) float64 {
	lengthScore := float64(requestLength) / 4.0
	return lengthScore*0.0345 + 120.0745
}

func CalculateDecodeScore(requestLength int) float64 {
	return float64(requestLength)
}
