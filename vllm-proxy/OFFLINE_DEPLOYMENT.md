# vLLM Proxy 离线部署指南

本指南适用于公司内网环境，通过 Git 克隆依赖包的方式进行离线部署。

## 前置要求

1. **Go 环境**：Go 1.21 或更高版本
2. **Git**：能够访问 GitHub（公司内网通常允许 Git 协议）
3. **项目代码**：vllm-proxy 项目完整代码

## 部署步骤

### 方案 1：使用自动化脚本（推荐）

#### Linux/macOS 环境

```bash
# 1. 进入项目目录
cd d:\code\vscode\vllm-proxy

# 2. 添加执行权限
chmod +x scripts/*.sh

# 3. 克隆所有依赖包
./scripts/clone-deps.sh

# 4. 生成 modules.txt 文件
./scripts/create-modules-txt.sh

# 5. 构建项目
go build -mod=vendor -o bin/vllm-proxy ./cmd
```

#### Windows PowerShell 环境

```powershell
# 1. 进入项目目录
cd d:\code\vscode\vllm-proxy

# 2. 克隆所有依赖包
.\scripts\clone-deps.ps1

# 3. 生成 modules.txt 文件
.\scripts\create-modules-txt.ps1

# 4. 构建项目
go build -mod=vendor -o bin\vllm-proxy .\cmd
```

### 方案 2：手动部署

#### 步骤 1：创建 vendor 目录结构

```bash
# Linux/macOS
mkdir -p vendor/github.com/{gin-gonic,prometheus,spf13,go-playground,go-yaml}
mkdir -p vendor/go.uber.org
mkdir -p vendor/golang.org/x
mkdir -p vendor/google.golang.org
mkdir -p vendor/gopkg.in

# Windows PowerShell
New-Item -ItemType Directory -Path "vendor\github.com\gin-gonic" -Force
New-Item -ItemType Directory -Path "vendor\github.com\prometheus" -Force
New-Item -ItemType Directory -Path "vendor\github.com\spf13" -Force
# ... 创建其他目录
```

#### 步骤 2：克隆主要依赖包

```bash
# Gin Web Framework
git clone --depth 1 --branch v1.9.1 https://github.com/gin-gonic/gin.git vendor/github.com/gin-gonic/gin

# Prometheus Client
git clone --depth 1 --branch v1.17.0 https://github.com/prometheus/client_golang.git vendor/github.com/prometheus/client_golang

# Cobra CLI
git clone --depth 1 --branch v1.8.0 https://github.com/spf13/cobra.git vendor/github.com/spf13/cobra

# Viper Config
git clone --depth 1 --branch v1.18.2 https://github.com/spf13/viper.git vendor/github.com/spf13/viper

# Zap Logger
git clone --depth 1 --branch v1.26.0 https://github.com/uber-go/zap.git vendor/go.uber.org/zap

# YAML
git clone --depth 1 --branch v3.0.1 https://github.com/go-yaml/yaml.git vendor/github.com/go-yaml/yaml
```

#### 步骤 3：克隆间接依赖包

```bash
# Gin 依赖
git clone --depth 1 --branch v0.1.0 https://github.com/gin-contrib/sse.git vendor/github.com/gin-contrib/sse
git clone --depth 1 --branch v10.14.0 https://github.com/go-playground/validator.git vendor/github.com/go-playground/validator
git clone --depth 1 --branch v0.14.1 https://github.com/go-playground/locales.git vendor/github.com/go-playground/locales
git clone --depth 1 --branch v0.18.1 https://github.com/go-playground/universal-translator.git vendor/github.com/go-playground/universal-translator

# Viper 依赖
git clone --depth 1 --branch v1.7.0 https://github.com/fsnotify/fsnotify.git vendor/github.com/fsnotify/fsnotify
git clone --depth 1 --branch v1.0.0 https://github.com/hashicorp/hcl.git vendor/github.com/hashicorp/hcl
git clone --depth 1 --branch v1.5.0 https://github.com/mitchellh/mapstructure.git vendor/github.com/mitchellh/mapstructure

# golang.org/x 依赖
git clone --depth 1 https://github.com/golang/crypto.git vendor/golang.org/x/crypto
git clone --depth 1 https://github.com/golang/net.git vendor/golang.org/x/net
git clone --depth 1 https://github.com/golang/sys.git vendor/golang.org/x/sys
git clone --depth 1 https://github.com/golang/text.git vendor/golang.org/x/text

# 其他依赖
# ... 参考 scripts/clone-deps.sh 中的完整列表
```

#### 步骤 4：创建 modules.txt

手动创建 `vendor/modules.txt` 文件，内容参考 `scripts/create-modules-txt.sh`。

#### 步骤 5：构建项目

```bash
go build -mod=vendor -o bin/vllm-proxy ./cmd
```

### 方案 3：从有网络的环境复制

#### 步骤 1：在有网络的环境中准备

```bash
# 1. 在有网络的环境中
cd d:\code\vscode\vllm-proxy

# 2. 下载所有依赖
go mod download

# 3. 生成 vendor 目录
go mod vendor

# 4. 打包整个项目
tar -czf vllm-proxy-vendor.tar.gz vendor/ go.mod go.sum *.go cmd/ config/ internal/ pkg/ configs/
```

#### 步骤 2：在公司服务器上部署

```bash
# 1. 解压项目
tar -xzf vllm-proxy-vendor.tar.gz

# 2. 构建
go build -mod=vendor -o bin/vllm-proxy ./cmd
```

## 验证部署

### 1. 检查二进制文件

```bash
# Linux/macOS
ls -lh bin/vllm-proxy
./bin/vllm-proxy --help

# Windows
dir bin\vllm-proxy.exe
.\bin\vllm-proxy.exe --help
```

### 2. 运行测试

```bash
# 使用 vendor 模式运行测试
go test -mod=vendor ./...
```

### 3. 启动服务

```bash
# 使用默认配置
./bin/vllm-proxy

# 使用配置文件
./bin/vllm-proxy --config configs/config.yaml

# 使用命令行参数
./bin/vllm-proxy \
  --host 0.0.0.0 \
  --port 8000 \
  --prefiller-hosts 10.0.0.1 10.0.0.2 \
  --prefiller-ports 8100 8101 \
  --decoder-hosts 10.0.0.3 10.0.0.4 \
  --decoder-ports 8200 8201
```

## 常见问题

### 1. Git 克隆失败

**问题**：某些仓库无法克隆

**解决方案**：
- 检查 Git 是否配置了正确的代理
- 尝试使用 SSH 协议代替 HTTPS
- 手动下载 ZIP 包并解压到 vendor 目录

### 2. modules.txt 不一致

**问题**：`go: inconsistent vendoring` 错误

**解决方案**：
```bash
# 重新生成 modules.txt
./scripts/create-modules-txt.sh

# 或者手动编辑 vendor/modules.txt
```

### 3. 缺少依赖包

**问题**：构建时提示缺少某个包

**解决方案**：
```bash
# 查看缺少的包
go list -m -mod=vendor all

# 手动克隆缺失的包
git clone --depth 1 https://github.com/xxx/xxx.git vendor/path/to/package
```

### 4. 版本不匹配

**问题**：依赖版本不匹配

**解决方案**：
- 检查 go.mod 中的版本要求
- 克隆指定版本的标签：`git clone --branch v1.2.3 ...`
- 更新 modules.txt 中的版本号

## 依赖包清单

### 主要依赖（直接依赖）

| 包名 | 版本 | GitHub 仓库 |
|------|------|-------------|
| gin | v1.9.1 | github.com/gin-gonic/gin |
| prometheus/client_golang | v1.17.0 | github.com/prometheus/client_golang |
| cobra | v1.8.0 | github.com/spf13/cobra |
| viper | v1.18.2 | github.com/spf13/viper |
| zap | v1.26.0 | go.uber.org/zap |
| yaml | v3.0.1 | gopkg.in/yaml.v3 |

### 间接依赖

完整列表请参考 `scripts/clone-deps.sh` 或 `go.mod` 文件。

## 构建选项

### 基本构建

```bash
go build -mod=vendor -o bin/vllm-proxy ./cmd
```

### 带版本信息的构建

```bash
VERSION=1.0.0
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

go build -mod=vendor \
  -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME" \
  -o bin/vllm-proxy ./cmd
```

### 优化构建

```bash
# 减小二进制文件大小
go build -mod=vendor -ldflags="-s -w" -o bin/vllm-proxy ./cmd

# 或者使用 UPX 压缩
upx --best bin/vllm-proxy
```

## 更新依赖

如果需要更新某个依赖包：

```bash
# 1. 删除旧的依赖
rm -rf vendor/github.com/gin-gonic/gin

# 2. 克隆新版本
git clone --depth 1 --branch v1.10.0 https://github.com/gin-gonic/gin.git vendor/github.com/gin-gonic/gin

# 3. 更新 modules.txt
# 手动编辑 vendor/modules.txt，更新版本号

# 4. 重新构建
go build -mod=vendor -o bin/vllm-proxy ./cmd
```

## 技术支持

如有问题，请参考：
- 项目文档：`README.md`
- 设计文档：`high_performance_proxy_design.md`
- 脚本目录：`scripts/`
