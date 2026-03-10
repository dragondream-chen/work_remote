#!/bin/bash
# 一键部署脚本 - 适用于内网环境
# 自动克隆依赖并构建项目

set -e

echo "========================================="
echo "  vLLM Proxy 离线部署脚本"
echo "========================================="
echo ""

# 检查 Go 环境
echo "检查 Go 环境..."
if ! command -v go &> /dev/null; then
    echo "错误: 未找到 Go 环境"
    echo "请先安装 Go 1.21 或更高版本"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✓ Go 版本: $GO_VERSION"

# 检查 Git
echo "检查 Git 环境..."
if ! command -v git &> /dev/null; then
    echo "错误: 未找到 Git"
    echo "请先安装 Git"
    exit 1
fi

GIT_VERSION=$(git --version | awk '{print $3}')
echo "✓ Git 版本: $GIT_VERSION"
echo ""

# 步骤 1: 克隆依赖
echo "========================================="
echo "步骤 1/3: 克隆依赖包"
echo "========================================="
echo ""

if [ -d "vendor" ] && [ -f "vendor/modules.txt" ]; then
    echo "vendor 目录已存在，跳过克隆步骤"
    echo "如需重新克隆，请先删除 vendor 目录"
else
    if [ -f "scripts/clone-deps.sh" ]; then
        chmod +x scripts/clone-deps.sh
        ./scripts/clone-deps.sh
    else
        echo "错误: 找不到 scripts/clone-deps.sh"
        exit 1
    fi
fi
echo ""

# 步骤 2: 生成 modules.txt
echo "========================================="
echo "步骤 2/3: 生成 modules.txt"
echo "========================================="
echo ""

if [ -f "vendor/modules.txt" ]; then
    echo "modules.txt 已存在"
    echo "文件大小: $(wc -l < vendor/modules.txt) 行"
else
    if [ -f "scripts/create-modules-txt.sh" ]; then
        chmod +x scripts/create-modules-txt.sh
        ./scripts/create-modules-txt.sh
    else
        echo "错误: 找不到 scripts/create-modules-txt.sh"
        exit 1
    fi
fi
echo ""

# 步骤 3: 构建项目
echo "========================================="
echo "步骤 3/3: 构建项目"
echo "========================================="
echo ""

# 创建 bin 目录
mkdir -p bin

# 构建参数
VERSION=${VERSION:-"1.0.0"}
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS="-ldflags -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME"

echo "构建版本: $VERSION"
echo "构建时间: $BUILD_TIME"
echo ""

# 执行构建
echo "开始构建..."
if go build -mod=vendor $LDFLAGS -o bin/vllm-proxy ./cmd; then
    echo ""
    echo "========================================="
    echo "  ✓ 构建成功！"
    echo "========================================="
    echo ""
    
    # 显示二进制文件信息
    if [ -f "bin/vllm-proxy" ]; then
        FILE_SIZE=$(du -h bin/vllm-proxy | cut -f1)
        echo "二进制文件: bin/vllm-proxy"
        echo "文件大小: $FILE_SIZE"
        echo ""
        
        # 显示使用帮助
        echo "使用方法:"
        echo "  ./bin/vllm-proxy --help"
        echo ""
        echo "配置文件:"
        echo "  ./bin/vllm-proxy --config configs/config.yaml"
        echo ""
        echo "命令行参数:"
        echo "  ./bin/vllm-proxy \\"
        echo "    --host 0.0.0.0 \\"
        echo "    --port 8000 \\"
        echo "    --prefiller-hosts 10.0.0.1 10.0.0.2 \\"
        echo "    --prefiller-ports 8100 8101 \\"
        echo "    --decoder-hosts 10.0.0.3 10.0.0.4 \\"
        echo "    --decoder-ports 8200 8201"
    fi
else
    echo ""
    echo "========================================="
    echo "  ✗ 构建失败"
    echo "========================================="
    echo ""
    echo "请检查错误信息并尝试手动构建"
    exit 1
fi
