#!/bin/bash
# 生成 vendor/modules.txt 文件
# 这是 Go modules vendor 模式必需的清单文件

set -e

echo "=== 生成 vendor/modules.txt ==="

# 确保 vendor 目录存在
mkdir -p vendor

# 创建正确格式的 modules.txt
# Go modules vendor 格式要求：
# - 每个模块一行：module-path version
# - 所有模块都需要列出，无论是直接依赖还是间接依赖

cat > vendor/modules.txt << 'EOF'
github.com/bytedance/sonic v1.9.1
github.com/cespare/xxhash/v2 v2.2.0
github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311
github.com/fsnotify/fsnotify v1.7.0
github.com/gabriel-vasile/mimetype v1.4.2
github.com/gin-contrib/sse v0.1.0
github.com/gin-gonic/gin v1.9.1
github.com/go-playground/locales v0.14.1
github.com/go-playground/universal-translator v0.18.1
github.com/go-playground/validator/v10 v10.14.0
github.com/goccy/go-json v0.10.2
github.com/hashicorp/hcl v1.0.0
github.com/inconshreveable/mousetrap v1.1.0
github.com/json-iterator/go v1.1.12
github.com/klauspost/cpuid/v2 v2.2.4
github.com/leodido/go-urn v1.2.4
github.com/magiconair/properties v1.8.7
github.com/mattn/go-isatty v0.0.19
github.com/mitchellh/mapstructure v1.5.0
github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
github.com/modern-go/reflect2 v1.0.2
github.com/pelletier/go-toml/v2 v2.1.0
github.com/prometheus/client_golang v1.17.0
github.com/prometheus/client_model v0.5.0
github.com/prometheus/common v0.45.0
github.com/prometheus/procfs v0.12.0
github.com/sagikazarmark/locafero v0.4.0
github.com/sagikazarmark/slog-shim v0.1.0
github.com/sourcegraph/conc v0.3.0
github.com/spf13/afero v1.11.0
github.com/spf13/cast v1.6.0
github.com/spf13/cobra v1.8.0
github.com/spf13/pflag v1.0.5
github.com/spf13/viper v1.18.2
github.com/subosito/gotenv v1.6.0
github.com/twitchyliquid64/golang-asm v0.15.1
github.com/ugorji/go/codec v1.2.11
go.uber.org/multierr v1.10.0
go.uber.org/zap v1.26.0
golang.org/x/arch v0.3.0
golang.org/x/crypto v0.16.0
golang.org/x/exp v0.0.0-20230905200255-921286631fa9
golang.org/x/net v0.19.0
golang.org/x/sys v0.15.0
golang.org/x/text v0.14.0
google.golang.org/protobuf v1.31.0
gopkg.in/ini.v1 v1.67.0
gopkg.in/yaml.v3 v3.0.1
EOF

echo "modules.txt 已生成"
echo ""

# 检查 vendor 目录中的关键依赖
echo "检查关键依赖目录..."
required_dirs=(
    "vendor/github.com/gin-gonic/gin"
    "vendor/github.com/prometheus/client_golang"
    "vendor/github.com/spf13/cobra"
    "vendor/github.com/spf13/viper"
    "vendor/go.uber.org/zap"
    "vendor/github.com/go-yaml/yaml"
)

all_exist=true
for dir in "${required_dirs[@]}"; do
    if [ -d "$dir" ]; then
        echo "✓ $dir"
    else
        echo "✗ $dir (缺失)"
        all_exist=false
    fi
done

if [ "$all_exist" = false ]; then
    echo ""
    echo "警告：部分依赖目录缺失，请先运行 clone-deps.sh"
fi

echo ""
echo "文件内容预览："
head -20 vendor/modules.txt
