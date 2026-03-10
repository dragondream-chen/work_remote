# Windows PowerShell 版本的 modules.txt 生成脚本
# 按照 Go 官方规范格式生成

Write-Host "=== 生成 vendor\modules.txt ===" -ForegroundColor Green

# 确保 vendor 目录存在
if (-not (Test-Path "vendor")) {
    New-Item -ItemType Directory -Path "vendor" -Force | Out-Null
}

# 创建符合 Go 官方规范的 modules.txt
# 格式说明：
# # module-path version
# ## explicit; go 1.xx  (如果是显式依赖)
# module-path/package
# module-path/package/subpackage

$modulesContent = @"
# github.com/gin-gonic/gin v1.9.1
## explicit; go 1.20
github.com/gin-gonic/gin
github.com/gin-gonic/gin/binding
github.com/gin-gonic/gin/internal/bytesconv
github.com/gin-gonic/gin/internal/json
github.com/gin-gonic/gin/render
github.com/gin-gonic/gin/testdata/protoexample
# github.com/prometheus/client_golang v1.17.0
## explicit; go 1.19
github.com/prometheus/client_golang/prometheus
github.com/prometheus/client_golang/prometheus/collectors
github.com/prometheus/client_golang/prometheus/internal
github.com/prometheus/client_golang/prometheus/promauto
github.com/prometheus/client_golang/prometheus/promhttp
github.com/prometheus/client_golang/prometheus/testutil/promlint
github.com/prometheus/client_golang/prometheus/testutil/testutil
# github.com/spf13/cobra v1.8.0
## explicit; go 1.15
github.com/spf13/cobra
# github.com/spf13/viper v1.18.2
## explicit; go 1.20
github.com/spf13/viper
github.com/spf13/viper/internal/encoding
github.com/spf13/viper/internal/encoding/dotenv
github.com/spf13/viper/internal/encoding/json
github.com/spf13/viper/internal/encoding/toml
github.com/spf13/viper/internal/encoding/yaml
github.com/spf13/viper/internal/features
# go.uber.org/zap v1.26.0
## explicit; go 1.19
go.uber.org/zap
go.uber.org/zap/buffer
go.uber.org/zap/internal
go.uber.org/zap/internal/bufferpool
go.uber.org/zap/internal/color
go.uber.org/zap/internal/exit
go.uber.org/zap/internal/pool
go.uber.org/zap/internal/ztest
go.uber.org/zap/zapcore
go.uber.org/zap/zapgrpc
go.uber.org/zap/zaptest
go.uber.org/zap/zaptest/observer
# gopkg.in/yaml.v3 v3.0.1
## explicit
gopkg.in/yaml.v3
# github.com/bytedance/sonic v1.9.1
github.com/bytedance/sonic
github.com/bytedance/sonic/ast
github.com/bytedance/sonic/internal/base64
github.com/bytedance/sonic/internal/caching
github.com/bytedance/sonic/internal/decoder
github.com/bytedance/sonic/internal/encoder
github.com/bytedance/sonic/internal/encoder/alg
github.com/bytedance/sonic/internal/encoder/alg/x86
github.com/bytedance/sonic/internal/encoder/alg/x86/quote
github.com/bytedance/sonic/internal/encoder/stream
github.com/bytedance/sonic/internal/rt
github.com/bytedance/sonic/loader
# github.com/cespare/xxhash/v2 v2.2.0
github.com/cespare/xxhash/v2
# github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311
github.com/chenzhuoyu/base64x
# github.com/fsnotify/fsnotify v1.7.0
## explicit; go 1.17
github.com/fsnotify/fsnotify
# github.com/gabriel-vasile/mimetype v1.4.2
github.com/gabriel-vasile/mimetype
github.com/gabriel-vasile/mimetype/internal/charset
github.com/gabriel-vasile/mimetype/internal/json
github.com/gabriel-vasile/mimetype/internal/magic
# github.com/gin-contrib/sse v0.1.0
github.com/gin-contrib/sse
# github.com/go-playground/locales v0.14.1
github.com/go-playground/locales
github.com/go-playground/locales/currency
# github.com/go-playground/universal-translator v0.18.1
github.com/go-playground/universal-translator
# github.com/go-playground/validator/v10 v10.14.0
github.com/go-playground/validator/v10
# github.com/goccy/go-json v0.10.2
github.com/goccy/go-json
github.com/goccy/go-json/internal/decoder
github.com/goccy/go-json/internal/encoder
github.com/goccy/go-json/internal/errors
github.com/goccy/go-json/internal/runtime
# github.com/hashicorp/hcl v1.0.0
github.com/hashicorp/hcl
github.com/hashicorp/hcl/hcl/ast
github.com/hashicorp/hcl/hcl/parser
github.com/hashicorp/hcl/hcl/printer
github.com/hashicorp/hcl/hcl/scanner
github.com/hashicorp/hcl/hcl/strconv
github.com/hashicorp/hcl/hcl/token
github.com/hashicorp/hcl/json/parser
github.com/hashicorp/hcl/json/scanner
github.com/hashicorp/hcl/json/token
# github.com/inconshreveable/mousetrap v1.1.0
github.com/inconshreveable/mousetrap
# github.com/json-iterator/go v1.1.12
github.com/json-iterator/go
# github.com/klauspost/cpuid/v2 v2.2.4
github.com/klauspost/cpuid/v2
# github.com/leodido/go-urn v1.2.4
github.com/leodido/go-urn
# github.com/magiconair/properties v1.8.7
github.com/magiconair/properties
# github.com/mattn/go-isatty v0.0.19
github.com/mattn/go-isatty
# github.com/mitchellh/mapstructure v1.5.0
## explicit; go 1.14
github.com/mitchellh/mapstructure
# github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
github.com/modern-go/concurrent
# github.com/modern-go/reflect2 v1.0.2
github.com/modern-go/reflect2
# github.com/pelletier/go-toml/v2 v2.1.0
github.com/pelletier/go-toml/v2
github.com/pelletier/go-toml/v2/internal/characters
github.com/pelletier/go-toml/v2/internal/danger
github.com/pelletier/go-toml/v2/internal/tracker
# github.com/prometheus/client_model v0.5.0
github.com/prometheus/client_model/go
# github.com/prometheus/common v0.45.0
github.com/prometheus/common/expfmt
github.com/prometheus/common/internal/bitbucket.org/ww/goautoneg
github.com/prometheus/common/model
# github.com/prometheus/procfs v0.12.0
github.com/prometheus/procfs
github.com/prometheus/procfs/internal/fs
github.com/prometheus/procfs/internal/util
# github.com/sagikazarmark/locafero v0.4.0
github.com/sagikazarmark/locafero
# github.com/sagikazarmark/slog-shim v0.1.0
github.com/sagikazarmark/slog-shim
# github.com/sourcegraph/conc v0.3.0
github.com/sourcegraph/conc
github.com/sourcegraph/conc/internal/multierror
github.com/sourcegraph/conc/iter
github.com/sourcegraph/conc/mutex
github.com/sourcegraph/conc/panics
github.com/sourcegraph/conc/waitgroup
# github.com/spf13/afero v1.11.0
## explicit; go 1.19
github.com/spf13/afero
github.com/spf13/afero/internal/common
github.com/spf13/afero/mem
github.com/spf13/afero/sftp
# github.com/spf13/cast v1.6.0
github.com/spf13/cast
# github.com/spf13/pflag v1.0.5
github.com/spf13/pflag
# github.com/subosito/gotenv v1.6.0
github.com/subosito/gotenv
# github.com/twitchyliquid64/golang-asm v0.15.1
github.com/twitchyliquid64/golang-asm
github.com/twitchyliquid64/golang-asm/obj
github.com/twitchyliquid64/golang-asm/obj/arm
github.com/twitchyliquid64/golang-asm/obj/arm64
github.com/twitchyliquid64/golang-asm/obj/mips
github.com/twitchyliquid64/golang-asm/obj/ppc64
github.com/twitchyliquid64/golang-asm/obj/riscv
github.com/twitchyliquid64/golang-asm/obj/s390x
github.com/twitchyliquid64/golang-asm/obj/wasm
github.com/twitchyliquid64/golang-asm/obj/x86
github.com/twitchyliquid64/golang-asm/objabi
github.com/twitchyliquid64/golang-asm/src
github.com/twitchyliquid64/golang-asm/sys
github.com/twitchyliquid64/golang-asm/types
# github.com/ugorji/go/codec v1.2.11
github.com/ugorji/go/codec
# go.uber.org/multierr v1.10.0
go.uber.org/multierr
# golang.org/x/arch v0.3.0
golang.org/x/arch/arm/armasm
golang.org/x/arch/arm64/arm64asm
golang.org/x/arch/ppc64/ppc64asm
golang.org/x/arch/x86/x86asm
# golang.org/x/crypto v0.16.0
golang.org/x/crypto/acme
golang.org/x/crypto/acme/autocert
golang.org/x/crypto/bcrypt
golang.org/x/crypto/blake2b
golang.org/x/crypto/blake2s
golang.org/x/crypto/blowfish
golang.org/x/crypto/bn256
golang.org/x/crypto/cast5
golang.org/x/crypto/chacha20
golang.org/x/crypto/chacha20poly1305
golang.org/x/crypto/cryptobyte
golang.org/x/crypto/cryptobyte/asn1
golang.org/x/crypto/curve25519
golang.org/x/crypto/curve25519/internal/field
golang.org/x/crypto/ed25519
golang.org/x/crypto/hkdf
golang.org/x/crypto/internal/alias
golang.org/x/crypto/internal/poly1305
golang.org/x/crypto/md4
golang.org/x/crypto/nacl/box
golang.org/x/crypto/nacl/secretbox
golang.org/x/crypto/nacl/sign
golang.org/x/crypto/ocsp
golang.org/x/crypto/openpgp
golang.org/x/crypto/openpgp/armor
golang.org/x/crypto/openpgp/elgamal
golang.org/x/crypto/openpgp/errors
golang.org/x/crypto/openpgp/packet
golang.org/x/crypto/openpgp/s2k
golang.org/x/crypto/pbkdf2
golang.org/x/crypto/pkcs12
golang.org/x/crypto/pkcs12/internal/rc2
golang.org/x/crypto/poly1305
golang.org/x/crypto/ripemd160
golang.org/x/crypto/salsa20
golang.org/x/crypto/salsa20/salsa
golang.org/x/crypto/scrypt
golang.org/x/crypto/sha3
golang.org/x/crypto/ssh
golang.org/x/crypto/ssh/agent
golang.org/x/crypto/ssh/internal/bcrypt_pbkdf
golang.org/x/crypto/ssh/knownhosts
golang.org/x/crypto/tea
golang.org/x/crypto/twofish
golang.org/x/crypto/xtea
golang.org/x/crypto/xts
# golang.org/x/exp v0.0.0-20230905200255-921286631fa9
golang.org/x/exp/constraints
golang.org/x/exp/maps
golang.org/x/exp/slices
# golang.org/x/net v0.19.0
golang.org/x/net/bpf
golang.org/x/net/context
golang.org/x/net/context/ctxhttp
golang.org/x/net/html
golang.org/x/net/html/atom
golang.org/x/net/html/charset
golang.org/x/net/http/httpguts
golang.org/x/net/http2
golang.org/x/net/http2/h2c
golang.org/x/net/http2/hpack
golang.org/x/net/idna
golang.org/x/net/internal/iana
golang.org/x/net/internal/socket
golang.org/x/net/internal/timeseries
golang.org/x/net/ipv4
golang.org/x/net/ipv6
golang.org/x/net/nettest
golang.org/x/net/netutil
golang.org/x/net/publicsuffix
golang.org/x/net/trace
golang.org/x/net/webdav
golang.org/x/net/webdav/internal/xml
golang.org/x/net/websocket
# golang.org/x/sys v0.15.0
golang.org/x/sys/cpu
golang.org/x/sys/execabs
golang.org/x/sys/plan9
golang.org/x/sys/unix
golang.org/x/sys/windows
golang.org/x/sys/windows/registry
golang.org/x/sys/windows/svc
golang.org/x/sys/windows/svc/debug
golang.org/x/sys/windows/svc/eventlog
# golang.org/x/text v0.14.0
golang.org/x/text/cases
golang.org/x/text/collate
golang.org/x/text/collate/build
golang.org/x/text/currency
golang.org/x/text/encoding
golang.org/x/text/encoding/charmap
golang.org/x/text/encoding/htmlindex
golang.org/x/text/encoding/ianaindex
golang.org/x/text/encoding/internal
golang.org/x/text/encoding/internal/identifier
golang.org/x/text/encoding/japanese
golang.org/x/text/encoding/korean
golang.org/x/text/encoding/simplifiedchinese
golang.org/x/text/encoding/traditionalchinese
golang.org/x/text/encoding/unicode
golang.org/x/text/encoding/unicode/override
golang.org/x/text/feature/plural
golang.org/x/text/internal
golang.org/x/text/internal/catmsg
golang.org/x/text/internal/colltab
golang.org/x/text/internal/export/idna
golang.org/x/text/internal/format
golang.org/x/text/internal/internal
golang.org/x/text/internal/language
golang.org/x/text/internal/language/compact
golang.org/x/text/internal/number
golang.org/x/text/internal/stringset
golang.org/x/text/internal/tag
golang.org/x/text/internal/testtext
golang.org/x/text/internal/utf8internal
golang.org/x/text/language
golang.org/x/text/language/display
golang.org/x/text/message
golang.org/x/text/message/catalog
golang.org/x/text/number
golang.org/x/text/runes
golang.org/x/text/search
golang.org/x/text/search/internal
golang.org/x/text/secure
golang.org/x/text/secure/bidirule
golang.org/x/text/secure/precis
golang.org/x/text/transform
golang.org/x/text/unicode
golang.org/x/text/unicode/bidi
golang.org/x/text/unicode/cldr
golang.org/x/text/unicode/norm
golang.org/x/text/unicode/rangetable
golang.org/x/text/width
# google.golang.org/protobuf v1.31.0
google.golang.org/protobuf/encoding/protodelim
google.golang.org/protobuf/encoding/protojson
google.golang.org/protobuf/encoding/prototext
google.golang.org/protobuf/encoding/protowire
google.golang.org/protobuf/internal/descfmt
google.golang.org/protobuf/internal/descopts
google.golang.org/protobuf/internal/detrand
google.golang.org/protobuf/internal/encoding/defval
google.golang.org/protobuf/internal/encoding/json
google.golang.org/protobuf/internal/encoding/messageset
google.golang.org/protobuf/internal/encoding/tag
google.golang.org/protobuf/internal/encoding/text
google.golang.org/protobuf/internal/errors
google.golang.org/protobuf/internal/filedesc
google.golang.org/protobuf/internal/filetype
google.golang.org/protobuf/internal/flags
google.golang.org/protobuf/internal/genid
google.golang.org/protobuf/internal/impl
google.golang.org/protobuf/internal/msgfmt
google.golang.org/protobuf/internal/order
google.golang.org/protobuf/internal/pragma
google.golang.org/protobuf/internal/set
google.golang.org/protobuf/internal/strs
google.golang.org/protobuf/internal/version
google.golang.org/protobuf/proto
google.golang.org/protobuf/protoadapt
google.golang.org/protobuf/reflect/protodesc
google.golang.org/protobuf/reflect/protoreflect
google.golang.org/protobuf/reflect/protoregistry
google.golang.org/protobuf/runtime/protoiface
google.golang.org/protobuf/runtime/protoimpl
google.golang.org/protobuf/testing/protocmp
google.golang.org/protobuf/testing/protopack
google.golang.org/protobuf/types/dynamicpb
google.golang.org/protobuf/types/known/anypb
google.golang.org/protobuf/types/known/durationpb
google.golang.org/protobuf/types/known/emptypb
google.golang.org/protobuf/types/known/fieldmaskpb
google.golang.org/protobuf/types/known/sourcecontextpb
google.golang.org/protobuf/types/known/structpb
google.golang.org/protobuf/types/known/timestamppb
google.golang.org/protobuf/types/known/wrapperspb
# gopkg.in/ini.v1 v1.67.0
gopkg.in/ini.v1
"@

# 写入文件
$modulesContent | Out-File -FilePath "vendor\modules.txt" -Encoding UTF8 -NoNewline

Write-Host "modules.txt 已生成" -ForegroundColor Green
Write-Host ""

# 检查 vendor 目录中的关键依赖
Write-Host "检查关键依赖目录..." -ForegroundColor Yellow
$requiredDirs = @(
    "vendor\github.com\gin-gonic\gin",
    "vendor\github.com\prometheus\client_golang",
    "vendor\github.com\spf13\cobra",
    "vendor\github.com\spf13\viper",
    "vendor\go.uber.org\zap",
    "vendor\github.com\go-yaml\yaml"
)

$allExist = $true
foreach ($dir in $requiredDirs) {
    if (Test-Path $dir) {
        Write-Host "✓ $dir" -ForegroundColor Green
    } else {
        Write-Host "✗ $dir (缺失)" -ForegroundColor Red
        $allExist = $false
    }
}

if (-not $allExist) {
    Write-Host ""
    Write-Host "警告：部分依赖目录缺失，请先运行 clone-deps.ps1" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "文件内容预览：" -ForegroundColor Cyan
Get-Content "vendor\modules.txt" | Select-Object -First 30

Write-Host ""
Write-Host "下一步：运行 go build -mod=vendor -o bin\vllm-proxy .\cmd 进行构建" -ForegroundColor Cyan
