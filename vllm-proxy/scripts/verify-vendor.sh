#!/bin/bash
# 验证和修复 vendor 目录
# 确保 vendor 目录结构和 modules.txt 格式正确

set -e

echo "=== 验证 vendor 目录 ==="

# 检查 vendor 目录是否存在
if [ ! -d "vendor" ]; then
    echo "错误: vendor 目录不存在"
    echo "请先运行 ./scripts/clone-deps.sh"
    exit 1
fi

# 检查 modules.txt 是否存在
if [ ! -f "vendor/modules.txt" ]; then
    echo "错误: vendor/modules.txt 不存在"
    echo "请先运行 ./scripts/create-modules-txt.sh"
    exit 1
fi

echo "✓ vendor 目录存在"
echo "✓ modules.txt 存在"

# 检查关键依赖目录
echo ""
echo "检查关键依赖目录..."

required_deps=(
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "go.uber.org/zap"
    "github.com/go-yaml/yaml"
)

missing_deps=()
for dep in "${required_deps[@]}"; do
    if [ -d "vendor/$dep" ]; then
        echo "✓ $dep"
    else
        echo "✗ $dep (缺失)"
        missing_deps+=("$dep")
    fi
done

if [ ${#missing_deps[@]} -gt 0 ]; then
    echo ""
    echo "警告: 以下依赖缺失:"
    for dep in "${missing_deps[@]}"; do
        echo "  - $dep"
    done
    echo ""
    echo "请运行 ./scripts/clone-deps.sh 克隆缺失的依赖"
fi

# 检查每个依赖是否有 go.mod 文件
echo ""
echo "检查依赖的 go.mod 文件..."

missing_gomod=()
for dep in "${required_deps[@]}"; do
    if [ -d "vendor/$dep" ]; then
        if [ -f "vendor/$dep/go.mod" ]; then
            echo "✓ $dep/go.mod"
        else
            echo "⚠ $dep/go.mod (缺失)"
            missing_gomod+=("$dep")
        fi
    fi
done

if [ ${#missing_gomod[@]} -gt 0 ]; then
    echo ""
    echo "警告: 以下依赖缺少 go.mod 文件:"
    for dep in "${missing_gomod[@]}"; do
        echo "  - $dep"
    done
    echo ""
    echo "这可能导致构建失败"
fi

# 验证 modules.txt 格式
echo ""
echo "验证 modules.txt 格式..."

# 检查 modules.txt 是否包含必要的依赖
missing_in_txt=()
for dep in "${required_deps[@]}"; do
    if grep -q "^$dep " vendor/modules.txt; then
        echo "✓ $dep 在 modules.txt 中"
    else
        echo "✗ $dep 不在 modules.txt 中"
        missing_in_txt+=("$dep")
    fi
done

if [ ${#missing_in_txt[@]} -gt 0 ]; then
    echo ""
    echo "警告: 以下依赖不在 modules.txt 中:"
    for dep in "${missing_in_txt[@]}"; do
        echo "  - $dep"
    done
    echo ""
    echo "请重新运行 ./scripts/create-modules-txt.sh"
fi

# 尝试构建
echo ""
echo "尝试构建..."
if go build -mod=vendor -o bin/vllm-proxy ./cmd 2>&1; then
    echo ""
    echo "========================================="
    echo "✓ 构建成功！"
    echo "========================================="
else
    echo ""
    echo "========================================="
    echo "✗ 构建失败"
    echo "========================================="
    echo ""
    echo "请检查错误信息并修复问题"
    exit 1
fi
