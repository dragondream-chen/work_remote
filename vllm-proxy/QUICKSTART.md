# vLLM Proxy 快速开始指南（离线部署）

## 适用场景

- 公司内网环境，无法直接访问 Go 模块代理
- 可以通过 Git 访问 GitHub 仓库
- 需要离线部署 Go 项目

## 一键部署

### Linux/macOS

```bash
# 进入项目目录
cd d:\code\vscode\vllm-proxy

# 添加执行权限
chmod +x scripts/deploy.sh

# 执行一键部署
./scripts/deploy.sh
```

### Windows PowerShell

```powershell
# 进入项目目录
cd d:\code\vscode\vllm-proxy

# 执行一键部署
.\scripts\deploy.ps1
```

## 分步部署

如果一键部署失败，可以按以下步骤手动操作：

### 步骤 1：克隆依赖包

```bash
# Linux/macOS
./scripts/clone-deps.sh

# Windows PowerShell
.\scripts\clone-deps.ps1
```

### 步骤 2：生成 modules.txt

```bash
# Linux/macOS
./scripts/create-modules-txt.sh

# Windows PowerShell
.\scripts\create-modules-txt.ps1
```

### 步骤 3：构建项目

```bash
# Linux/macOS
go build -mod=vendor -o bin/vllm-proxy ./cmd

# Windows PowerShell
go build -mod=vendor -o bin\vllm-proxy.exe .\cmd
```

## 使用 Makefile

```bash
# 克隆依赖并生成 vendor 目录
make vendor-deps

# 使用 vendor 模式构建
make vendor-build
```

## 验证部署

```bash
# 查看帮助信息
./bin/vllm-proxy --help

# 使用配置文件运行
./bin/vllm-proxy --config configs/config.yaml

# 使用命令行参数运行
./bin/vllm-proxy \
  --host 0.0.0.0 \
  --port 8000 \
  --prefiller-hosts 10.0.0.1 10.0.0.2 \
  --prefiller-ports 8100 8101 \
  --decoder-hosts 10.0.0.3 10.0.0.4 \
  --decoder-ports 8200 8201
```

## 目录结构

```
vllm-proxy/
├── bin/                    # 构建输出目录
│   └── vllm-proxy         # 可执行文件
├── cmd/                    # 主程序入口
│   └── main.go
├── config/                 # 配置管理
├── internal/               # 内部包
│   ├── instance/          # 实例管理
│   ├── kvtransfer/        # KV 传输处理
│   ├── loadbalancer/      # 负载均衡
│   ├── metrics/           # 监控指标
│   └── server/            # HTTP 服务器
├── pkg/                    # 公共工具包
├── scripts/                # 部署脚本
│   ├── clone-deps.sh      # Linux/macOS 依赖克隆
│   ├── clone-deps.ps1     # Windows 依赖克隆
│   ├── create-modules-txt.sh   # Linux/macOS modules.txt 生成
│   ├── create-modules-txt.ps1  # Windows modules.txt 生成
│   ├── deploy.sh          # Linux/macOS 一键部署
│   └── deploy.ps1         # Windows 一键部署
├── vendor/                 # 依赖包目录（脚本自动生成）
│   ├── github.com/        # GitHub 依赖
│   ├── go.uber.org/       # Uber 依赖
│   ├── golang.org/x/      # Golang 官方扩展
│   ├── google.golang.org/ # Google 依赖
│   ├── gopkg.in/          # 第三方包
│   └── modules.txt        # 依赖清单
├── Makefile               # 构建脚本
├── go.mod                 # Go 模块定义
└── OFFLINE_DEPLOYMENT.md  # 详细部署文档
```

## 常见问题

### 1. Git 克隆失败

**原因**：网络问题或 Git 配置问题

**解决方案**：
```bash
# 检查 Git 配置
git config --global --list

# 配置 Git 代理（如果需要）
git config --global http.proxy http://proxy-server:port

# 或使用 SSH 协议
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

### 2. 构建失败

**原因**：缺少依赖或 modules.txt 不正确

**解决方案**：
```bash
# 检查 vendor 目录
ls -la vendor/

# 检查 modules.txt
cat vendor/modules.txt

# 重新生成 modules.txt
./scripts/create-modules-txt.sh
```

### 3. 运行时错误

**原因**：配置错误或依赖服务未启动

**解决方案**：
```bash
# 检查配置文件
cat configs/config.yaml

# 检查端口占用
netstat -tlnp | grep 8000

# 查看日志
./bin/vllm-proxy --config configs/config.yaml
```

## 技术支持

详细文档请参考：
- [离线部署指南](OFFLINE_DEPLOYMENT.md)
- [项目设计文档](high_performance_proxy_design.md)
- [项目 README](README.md)
