package loadbalancer

import (
	"testing"
)

func TestNewServerState(t *testing.T) {
	server := NewServerState("localhost", 8000, 1)

	if server.Host != "localhost" {
		t.Errorf("Expected host localhost, got %s", server.Host)
	}
	if server.Port != 8000 {
		t.Errorf("Expected port 8000, got %d", server.Port)
	}
	if server.Weight != 1 {
		t.Errorf("Expected weight 1, got %d", server.Weight)
	}
	if !server.Healthy {
		t.Error("Expected server to be healthy by default")
	}
}

func TestServerStateAddTokens(t *testing.T) {
	server := NewServerState("localhost", 8000, 1)

	server.AddTokens(100)
	if server.ActiveTokens != 100 {
		t.Errorf("Expected 100 active tokens, got %d", server.ActiveTokens)
	}

	server.AddTokens(50)
	if server.ActiveTokens != 150 {
		t.Errorf("Expected 150 active tokens, got %d", server.ActiveTokens)
	}
}

func TestServerStateReleaseTokens(t *testing.T) {
	server := NewServerState("localhost", 8000, 1)

	server.AddTokens(100)
	server.ReleaseTokens(50)

	if server.ActiveTokens != 50 {
		t.Errorf("Expected 50 active tokens, got %d", server.ActiveTokens)
	}
}

func TestServerStateCalculatePrefillPriority(t *testing.T) {
	server := NewServerState("localhost", 8000, 1)

	server.AddTokens(100)
	server.AddKVCache(100)

	priority := server.CalculatePrefillPriority()
	expected := float64(100) + float64(100)*0.3

	if priority != expected {
		t.Errorf("Expected priority %f, got %f", expected, priority)
	}
}

func TestServerStateCalculateDecodePriority(t *testing.T) {
	server := NewServerState("localhost", 8000, 1)

	server.AddTokens(100)

	priority := server.CalculateDecodePriority()
	expected := float64(100)

	if priority != expected {
		t.Errorf("Expected priority %f, got %f", expected, priority)
	}
}

func TestNewLoadBalancer(t *testing.T) {
	prefillers := []*ServerState{
		NewServerState("localhost", 8100, 1),
		NewServerState("localhost", 8101, 1),
	}
	decoders := []*ServerState{
		NewServerState("localhost", 8200, 1),
		NewServerState("localhost", 8201, 1),
	}

	lb := NewLoadBalancer(prefillers, decoders)

	if lb.PrefillerCount() != 2 {
		t.Errorf("Expected 2 prefillers, got %d", lb.PrefillerCount())
	}
	if lb.DecoderCount() != 2 {
		t.Errorf("Expected 2 decoders, got %d", lb.DecoderCount())
	}
}

func TestLoadBalancerSelectPrefiller(t *testing.T) {
	prefillers := []*ServerState{
		NewServerState("localhost", 8100, 1),
		NewServerState("localhost", 8101, 1),
	}
	decoders := []*ServerState{}

	lb := NewLoadBalancer(prefillers, decoders)

	server, idx := lb.SelectPrefiller(100)
	if server == nil {
		t.Error("Expected to select a prefiller")
	}
	if idx < 0 || idx >= len(prefillers) {
		t.Errorf("Invalid prefiller index: %d", idx)
	}
}

func TestLoadBalancerSelectDecoder(t *testing.T) {
	prefillers := []*ServerState{}
	decoders := []*ServerState{
		NewServerState("localhost", 8200, 1),
		NewServerState("localhost", 8201, 1),
	}

	lb := NewLoadBalancer(prefillers, decoders)

	server, idx := lb.SelectDecoder(100)
	if server == nil {
		t.Error("Expected to select a decoder")
	}
	if idx < 0 || idx >= len(decoders) {
		t.Errorf("Invalid decoder index: %d", idx)
	}
}

func TestCalculatePrefillScore(t *testing.T) {
	score := CalculatePrefillScore(1000)
	expected := float64(1000)/4.0*0.0345 + 120.0745

	if score != expected {
		t.Errorf("Expected score %f, got %f", expected, score)
	}
}

func TestCalculateDecodeScore(t *testing.T) {
	score := CalculateDecodeScore(1000)
	expected := float64(1000)

	if score != expected {
		t.Errorf("Expected score %f, got %f", expected, score)
	}
}

func TestLoadBalancerAddPrefiller(t *testing.T) {
	lb := NewLoadBalancer(nil, nil)

	server := NewServerState("localhost", 8100, 1)
	lb.AddPrefiller(server)

	if lb.PrefillerCount() != 1 {
		t.Errorf("Expected 1 prefiller, got %d", lb.PrefillerCount())
	}
}

func TestLoadBalancerAddDecoder(t *testing.T) {
	lb := NewLoadBalancer(nil, nil)

	server := NewServerState("localhost", 8200, 1)
	lb.AddDecoder(server)

	if lb.DecoderCount() != 1 {
		t.Errorf("Expected 1 decoder, got %d", lb.DecoderCount())
	}
}

func TestLoadBalancerRemovePrefiller(t *testing.T) {
	server := NewServerState("localhost", 8100, 1)
	lb := NewLoadBalancer([]*ServerState{server}, nil)

	lb.RemovePrefiller(server)

	if lb.PrefillerCount() != 0 {
		t.Errorf("Expected 0 prefillers, got %d", lb.PrefillerCount())
	}
}

func TestLoadBalancerRemoveDecoder(t *testing.T) {
	server := NewServerState("localhost", 8200, 1)
	lb := NewLoadBalancer(nil, []*ServerState{server})

	lb.RemoveDecoder(server)

	if lb.DecoderCount() != 0 {
		t.Errorf("Expected 0 decoders, got %d", lb.DecoderCount())
	}
}
