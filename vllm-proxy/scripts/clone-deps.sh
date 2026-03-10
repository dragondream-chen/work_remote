#!/bin/bash
# 通过 Git 克隆所有依赖包到本地 vendor 目录
# 适用于内网环境，只需要 git 访问权限

set -e

echo "=== 开始通过 Git 克隆依赖包 ==="

# 创建 vendor 目录结构
mkdir -p vendor/github.com/gin-gonic
mkdir -p vendor/github.com/prometheus
mkdir -p vendor/github.com/spf13
mkdir -p vendor/github.com/go-playground
mkdir -p vendor/github.com/go-yaml
mkdir -p vendor/github.com/ugorji
mkdir -p vendor/github.com/modern-go
mkdir -p vendor/github.com/mattn
mkdir -p vendor/github.com/mitchellh
mkdir -p vendor/github.com/pelletier
mkdir -p vendor/github.com/hashicorp
mkdir -p vendor/github.com/inconshreveable
mkdir -p vendor/github.com/fsnotify
mkdir -p vendor/github.com/magiconair
mkdir -p vendor/github.com/subosito
mkdir -p vendor/github.com/sagikazarmark
mkdir -p vendor/github.com/sourcegraph
mkdir -p vendor/github.com/bytedance
mkdir -p vendor/github.com/chenzhuoyu
mkdir -p vendor/github.com/gabriel-vasile
mkdir -p vendor/github.com/goccy
mkdir -p vendor/github.com/klauspost
mkdir -p vendor/github.com/leodido
mkdir -p vendor/github.com/twitchyliquid64
mkdir -p vendor/github.com/cespare
mkdir -p vendor/github.com/json-iterator
mkdir -p vendor/go.uber.org
mkdir -p vendor/golang.org/x
mkdir -p vendor/google.golang.org
mkdir -p vendor/gopkg.in

echo "目录结构创建完成"

# 主要依赖包
echo "克隆主要依赖包..."

# 1. Gin Web Framework
if [ ! -d "vendor/github.com/gin-gonic/gin" ]; then
    echo "克隆 gin..."
    git clone --depth 1 --branch v1.9.1 https://github.com/gin-gonic/gin.git vendor/github.com/gin-gonic/gin
fi

# 2. Prometheus Client
if [ ! -d "vendor/github.com/prometheus/client_golang" ]; then
    echo "克隆 prometheus client_golang..."
    git clone --depth 1 --branch v1.17.0 https://github.com/prometheus/client_golang.git vendor/github.com/prometheus/client_golang
fi

# 3. Cobra CLI
if [ ! -d "vendor/github.com/spf13/cobra" ]; then
    echo "克隆 cobra..."
    git clone --depth 1 --branch v1.8.0 https://github.com/spf13/cobra.git vendor/github.com/spf13/cobra
fi

# 4. Viper Config
if [ ! -d "vendor/github.com/spf13/viper" ]; then
    echo "克隆 viper..."
    git clone --depth 1 --branch v1.18.2 https://github.com/spf13/viper.git vendor/github.com/spf13/viper
fi

# 5. Zap Logger
if [ ! -d "vendor/go.uber.org/zap" ]; then
    echo "克隆 zap..."
    git clone --depth 1 --branch v1.26.0 https://github.com/uber-go/zap.git vendor/go.uber.org/zap
fi

# 6. YAML
if [ ! -d "vendor/github.com/go-yaml/yaml" ]; then
    echo "克隆 yaml..."
    git clone --depth 1 --branch v3.0.1 https://github.com/go-yaml/yaml.git vendor/github.com/go-yaml/yaml
fi

echo "主要依赖包克隆完成"

# 间接依赖包
echo "克隆间接依赖包..."

# Gin dependencies
if [ ! -d "vendor/github.com/gin-contrib/sse" ]; then
    git clone --depth 1 --branch v0.1.0 https://github.com/gin-contrib/sse.git vendor/github.com/gin-contrib/sse
fi

if [ ! -d "vendor/github.com/go-playground/validator" ]; then
    git clone --depth 1 --branch v10.14.0 https://github.com/go-playground/validator.git vendor/github.com/go-playground/validator
fi

if [ ! -d "vendor/github.com/go-playground/locales" ]; then
    git clone --depth 1 --branch v0.14.1 https://github.com/go-playground/locales.git vendor/github.com/go-playground/locales
fi

if [ ! -d "vendor/github.com/go-playground/universal-translator" ]; then
    git clone --depth 1 --branch v0.18.1 https://github.com/go-playground/universal-translator.git vendor/github.com/go-playground/universal-translator
fi

if [ ! -d "vendor/github.com/ugorji/go" ]; then
    git clone --depth 1 --branch v1.2.11 https://github.com/ugorji/go.git vendor/github.com/ugorji/go
fi

if [ ! -d "vendor/github.com/gabriel-vasile/mimetype" ]; then
    git clone --depth 1 --branch v1.4.2 https://github.com/gabriel-vasile/mimetype.git vendor/github.com/gabriel-vasile/mimetype
fi

if [ ! -d "vendor/github.com/bytedance/sonic" ]; then
    git clone --depth 1 --branch v1.9.1 https://github.com/bytedance/sonic.git vendor/github.com/bytedance/sonic
fi

if [ ! -d "vendor/github.com/chenzhuoyu/base64x" ]; then
    git clone --depth 1 https://github.com/chenzhuoyu/base64x.git vendor/github.com/chenzhuoyu/base64x
fi

if [ ! -d "vendor/github.com/goccy/go-json" ]; then
    git clone --depth 1 --branch v0.10.2 https://github.com/goccy/go-json.git vendor/github.com/goccy/go-json
fi

if [ ! -d "vendor/github.com/klauspost/cpuid" ]; then
    git clone --depth 1 --branch v2.2.4 https://github.com/klauspost/cpuid.git vendor/github.com/klauspost/cpuid
fi

if [ ! -d "vendor/github.com/twitchyliquid64/golang-asm" ]; then
    git clone --depth 1 --branch v0.15.1 https://github.com/twitchyliquid64/golang-asm.git vendor/github.com/twitchyliquid64/golang-asm
fi

if [ ! -d "vendor/github.com/leodido/go-urn" ]; then
    git clone --depth 1 --branch v1.2.4 https://github.com/leodido/go-urn.git vendor/github.com/leodido/go-urn
fi

# Viper dependencies
if [ ! -d "vendor/github.com/fsnotify/fsnotify" ]; then
    git clone --depth 1 --branch v1.7.0 https://github.com/fsnotify/fsnotify.git vendor/github.com/fsnotify/fsnotify
fi

if [ ! -d "vendor/github.com/hashicorp/hcl" ]; then
    git clone --depth 1 --branch v1.0.0 https://github.com/hashicorp/hcl.git vendor/github.com/hashicorp/hcl
fi

if [ ! -d "vendor/github.com/magiconair/properties" ]; then
    git clone --depth 1 --branch v1.8.7 https://github.com/magiconair/properties.git vendor/github.com/magiconair/properties
fi

if [ ! -d "vendor/github.com/mitchellh/mapstructure" ]; then
    git clone --depth 1 --branch v1.5.0 https://github.com/mitchellh/mapstructure.git vendor/github.com/mitchellh/mapstructure
fi

if [ ! -d "vendor/github.com/pelletier/go-toml" ]; then
    git clone --depth 1 --branch v2.1.0 https://github.com/pelletier/go-toml.git vendor/github.com/pelletier/go-toml
fi

if [ ! -d "vendor/github.com/spf13/afero" ]; then
    git clone --depth 1 --branch v1.11.0 https://github.com/spf13/afero.git vendor/github.com/spf13/afero
fi

if [ ! -d "vendor/github.com/spf13/cast" ]; then
    git clone --depth 1 --branch v1.6.0 https://github.com/spf13/cast.git vendor/github.com/spf13/cast
fi

if [ ! -d "vendor/github.com/spf13/pflag" ]; then
    git clone --depth 1 --branch v1.0.5 https://github.com/spf13/pflag.git vendor/github.com/spf13/pflag
fi

if [ ! -d "vendor/github.com/subosito/gotenv" ]; then
    git clone --depth 1 --branch v1.6.0 https://github.com/subosito/gotenv.git vendor/github.com/subosito/gotenv
fi

if [ ! -d "vendor/github.com/sagikazarmark/locafero" ]; then
    git clone --depth 1 --branch v0.4.0 https://github.com/sagikazarmark/locafero.git vendor/github.com/sagikazarmark/locafero
fi

if [ ! -d "vendor/github.com/sagikazarmark/slog-shim" ]; then
    git clone --depth 1 --branch v0.1.0 https://github.com/sagikazarmark/slog-shim.git vendor/github.com/sagikazarmark/slog-shim
fi

if [ ! -d "vendor/github.com/sourcegraph/conc" ]; then
    git clone --depth 1 --branch v0.3.0 https://github.com/sourcegraph/conc.git vendor/github.com/sourcegraph/conc
fi

# Other dependencies
if [ ! -d "vendor/github.com/inconshreveable/mousetrap" ]; then
    git clone --depth 1 --branch v1.1.0 https://github.com/inconshreveable/mousetrap.git vendor/github.com/inconshreveable/mousetrap
fi

if [ ! -d "vendor/github.com/mattn/go-isatty" ]; then
    git clone --depth 1 --branch v0.0.19 https://github.com/mattn/go-isatty.git vendor/github.com/mattn/go-isatty
fi

if [ ! -d "vendor/github.com/modern-go/concurrent" ]; then
    git clone --depth 1 https://github.com/modern-go/concurrent.git vendor/github.com/modern-go/concurrent
fi

if [ ! -d "vendor/github.com/modern-go/reflect2" ]; then
    git clone --depth 1 --branch v1.0.2 https://github.com/modern-go/reflect2.git vendor/github.com/modern-go/reflect2
fi

if [ ! -d "vendor/github.com/json-iterator/go" ]; then
    git clone --depth 1 --branch v1.1.12 https://github.com/json-iterator/go.git vendor/github.com/json-iterator/go
fi

if [ ! -d "vendor/github.com/cespare/xxhash" ]; then
    git clone --depth 1 --branch v2.2.0 https://github.com/cespare/xxhash.git vendor/github.com/cespare/xxhash
fi

# Prometheus dependencies
if [ ! -d "vendor/github.com/prometheus/client_model" ]; then
    git clone --depth 1 --branch v0.5.0 https://github.com/prometheus/client_model.git vendor/github.com/prometheus/client_model
fi

if [ ! -d "vendor/github.com/prometheus/common" ]; then
    git clone --depth 1 --branch v0.45.0 https://github.com/prometheus/common.git vendor/github.com/prometheus/common
fi

if [ ! -d "vendor/github.com/prometheus/procfs" ]; then
    git clone --depth 1 --branch v0.12.0 https://github.com/prometheus/procfs.git vendor/github.com/prometheus/procfs
fi

# Uber dependencies
if [ ! -d "vendor/go.uber.org/multierr" ]; then
    git clone --depth 1 --branch v1.10.0 https://github.com/uber-go/multierr.git vendor/go.uber.org/multierr
fi

# Golang.org/x dependencies
echo "克隆 golang.org/x 依赖..."

if [ ! -d "vendor/golang.org/x/arch" ]; then
    git clone --depth 1 https://github.com/golang/arch.git vendor/golang.org/x/arch
fi

if [ ! -d "vendor/golang.org/x/crypto" ]; then
    git clone --depth 1 https://github.com/golang/crypto.git vendor/golang.org/x/crypto
fi

if [ ! -d "vendor/golang.org/x/net" ]; then
    git clone --depth 1 https://github.com/golang/net.git vendor/golang.org/x/net
fi

if [ ! -d "vendor/golang.org/x/sys" ]; then
    git clone --depth 1 https://github.com/golang/sys.git vendor/golang.org/x/sys
fi

if [ ! -d "vendor/golang.org/x/text" ]; then
    git clone --depth 1 https://github.com/golang/text.git vendor/golang.org/x/text
fi

if [ ! -d "vendor/golang.org/x/exp" ]; then
    git clone --depth 1 https://github.com/golang/exp.git vendor/golang.org/x/exp
fi

# Google protobuf
if [ ! -d "vendor/google.golang.org/protobuf" ]; then
    git clone --depth 1 --branch v1.31.0 https://github.com/protocolbuffers/protobuf-go.git vendor/google.golang.org/protobuf
fi

# gopkg.in dependencies
if [ ! -d "vendor/gopkg.in/ini.v1" ]; then
    git clone --depth 1 --branch v1.67.0 https://github.com/go-ini/ini.git vendor/gopkg.in/ini.v1
fi

echo "=== 所有依赖包克隆完成 ==="
echo ""
echo "下一步：运行 ./create-modules-txt.sh 生成 modules.txt 文件"
