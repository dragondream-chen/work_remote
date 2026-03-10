# Windows PowerShell 版本的依赖克隆脚本
# 适用于内网环境，只需要 git 访问权限

Write-Host "=== 开始通过 Git 克隆依赖包 ===" -ForegroundColor Green

# 创建 vendor 目录结构
$vendorDirs = @(
    "vendor\github.com\gin-gonic",
    "vendor\github.com\prometheus",
    "vendor\github.com\spf13",
    "vendor\github.com\go-playground",
    "vendor\github.com\go-yaml",
    "vendor\github.com\ugorji",
    "vendor\github.com\modern-go",
    "vendor\github.com\mattn",
    "vendor\github.com\mitchellh",
    "vendor\github.com\pelletier",
    "vendor\github.com\hashicorp",
    "vendor\github.com\inconshreveable",
    "vendor\github.com\fsnotify",
    "vendor\github.com\magiconair",
    "vendor\github.com\subosito",
    "vendor\github.com\sagikazarmark",
    "vendor\github.com\sourcegraph",
    "vendor\github.com\bytedance",
    "vendor\github.com\chenzhuoyu",
    "vendor\github.com\gabriel-vasile",
    "vendor\github.com\goccy",
    "vendor\github.com\klauspost",
    "vendor\github.com\leodido",
    "vendor\github.com\twitchyliquid64",
    "vendor\github.com\cespare",
    "vendor\github.com\json-iterator",
    "vendor\github.com\gin-contrib",
    "vendor\go.uber.org",
    "vendor\golang.org\x",
    "vendor\google.golang.org",
    "vendor\gopkg.in"
)

foreach ($dir in $vendorDirs) {
    if (-not (Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir -Force | Out-Null
    }
}

Write-Host "目录结构创建完成" -ForegroundColor Cyan

# 定义依赖包列表
$dependencies = @(
    @{Name="gin"; Path="vendor\github.com\gin-gonic\gin"; Url="https://github.com/gin-gonic/gin.git"; Tag="v1.9.1"},
    @{Name="prometheus client_golang"; Path="vendor\github.com\prometheus\client_golang"; Url="https://github.com/prometheus/client_golang.git"; Tag="v1.17.0"},
    @{Name="cobra"; Path="vendor\github.com\spf13\cobra"; Url="https://github.com/spf13/cobra.git"; Tag="v1.8.0"},
    @{Name="viper"; Path="vendor\github.com\spf13\viper"; Url="https://github.com/spf13/viper.git"; Tag="v1.18.2"},
    @{Name="zap"; Path="vendor\go.uber.org\zap"; Url="https://github.com/uber-go/zap.git"; Tag="v1.26.0"},
    @{Name="yaml"; Path="vendor\github.com\go-yaml\yaml"; Url="https://github.com/go-yaml/yaml.git"; Tag="v3.0.1"},
    @{Name="gin-contrib/sse"; Path="vendor\github.com\gin-contrib\sse"; Url="https://github.com/gin-contrib/sse.git"; Tag="v0.1.0"},
    @{Name="validator"; Path="vendor\github.com\go-playground\validator"; Url="https://github.com/go-playground/validator.git"; Tag="v10.14.0"},
    @{Name="locales"; Path="vendor\github.com\go-playground\locales"; Url="https://github.com/go-playground/locales.git"; Tag="v0.14.1"},
    @{Name="universal-translator"; Path="vendor\github.com\go-playground\universal-translator"; Url="https://github.com/go-playground/universal-translator.git"; Tag="v0.18.1"},
    @{Name="ugorji/go"; Path="vendor\github.com\ugorji\go"; Url="https://github.com/ugorji/go.git"; Tag="v1.2.11"},
    @{Name="mimetype"; Path="vendor\github.com\gabriel-vasile\mimetype"; Url="https://github.com/gabriel-vasile/mimetype.git"; Tag="v1.4.2"},
    @{Name="sonic"; Path="vendor\github.com\bytedance\sonic"; Url="https://github.com/bytedance/sonic.git"; Tag="v1.9.1"},
    @{Name="base64x"; Path="vendor\github.com\chenzhuoyu\base64x"; Url="https://github.com/chenzhuoyu/base64x.git"; Tag=""},
    @{Name="go-json"; Path="vendor\github.com\goccy\go-json"; Url="https://github.com/goccy/go-json.git"; Tag="v0.10.2"},
    @{Name="cpuid"; Path="vendor\github.com\klauspost\cpuid"; Url="https://github.com/klauspost/cpuid.git"; Tag="v2.2.4"},
    @{Name="golang-asm"; Path="vendor\github.com\twitchyliquid64\golang-asm"; Url="https://github.com/twitchyliquid64/golang-asm.git"; Tag="v0.15.1"},
    @{Name="go-urn"; Path="vendor\github.com\leodido\go-urn"; Url="https://github.com/leodido/go-urn.git"; Tag="v1.2.4"},
    @{Name="fsnotify"; Path="vendor\github.com\fsnotify\fsnotify"; Url="https://github.com/fsnotify/fsnotify.git"; Tag="v1.7.0"},
    @{Name="hcl"; Path="vendor\github.com\hashicorp\hcl"; Url="https://github.com/hashicorp/hcl.git"; Tag="v1.0.0"},
    @{Name="properties"; Path="vendor\github.com\magiconair\properties"; Url="https://github.com/magiconair/properties.git"; Tag="v1.8.7"},
    @{Name="mapstructure"; Path="vendor\github.com\mitchellh\mapstructure"; Url="https://github.com/mitchellh/mapstructure.git"; Tag="v1.5.0"},
    @{Name="go-toml"; Path="vendor\github.com\pelletier\go-toml"; Url="https://github.com/pelletier/go-toml.git"; Tag="v2.1.0"},
    @{Name="afero"; Path="vendor\github.com\spf13\afero"; Url="https://github.com/spf13/afero.git"; Tag="v1.11.0"},
    @{Name="cast"; Path="vendor\github.com\spf13\cast"; Url="https://github.com/spf13/cast.git"; Tag="v1.6.0"},
    @{Name="pflag"; Path="vendor\github.com\spf13\pflag"; Url="https://github.com/spf13/pflag.git"; Tag="v1.0.5"},
    @{Name="gotenv"; Path="vendor\github.com\subosito\gotenv"; Url="https://github.com/subosito/gotenv.git"; Tag="v1.6.0"},
    @{Name="locafero"; Path="vendor\github.com\sagikazarmark\locafero"; Url="https://github.com/sagikazarmark/locafero.git"; Tag="v0.4.0"},
    @{Name="slog-shim"; Path="vendor\github.com\sagikazarmark\slog-shim"; Url="https://github.com/sagikazarmark/slog-shim.git"; Tag="v0.1.0"},
    @{Name="conc"; Path="vendor\github.com\sourcegraph\conc"; Url="https://github.com/sourcegraph/conc.git"; Tag="v0.3.0"},
    @{Name="mousetrap"; Path="vendor\github.com\inconshreveable\mousetrap"; Url="https://github.com/inconshreveable/mousetrap.git"; Tag="v1.1.0"},
    @{Name="go-isatty"; Path="vendor\github.com\mattn\go-isatty"; Url="https://github.com/mattn/go-isatty.git"; Tag="v0.0.19"},
    @{Name="concurrent"; Path="vendor\github.com\modern-go\concurrent"; Url="https://github.com/modern-go/concurrent.git"; Tag=""},
    @{Name="reflect2"; Path="vendor\github.com\modern-go\reflect2"; Url="https://github.com/modern-go/reflect2.git"; Tag="v1.0.2"},
    @{Name="json-iterator"; Path="vendor\github.com\json-iterator\go"; Url="https://github.com/json-iterator/go.git"; Tag="v1.1.12"},
    @{Name="xxhash"; Path="vendor\github.com\cespare\xxhash"; Url="https://github.com/cespare/xxhash.git"; Tag="v2.2.0"},
    @{Name="client_model"; Path="vendor\github.com\prometheus\client_model"; Url="https://github.com/prometheus/client_model.git"; Tag="v0.5.0"},
    @{Name="common"; Path="vendor\github.com\prometheus\common"; Url="https://github.com/prometheus/common.git"; Tag="v0.45.0"},
    @{Name="procfs"; Path="vendor\github.com\prometheus\procfs"; Url="https://github.com/prometheus/procfs.git"; Tag="v0.12.0"},
    @{Name="multierr"; Path="vendor\go.uber.org\multierr"; Url="https://github.com/uber-go/multierr.git"; Tag="v1.10.0"},
    @{Name="arch"; Path="vendor\golang.org\x\arch"; Url="https://github.com/golang/arch.git"; Tag=""},
    @{Name="crypto"; Path="vendor\golang.org\x\crypto"; Url="https://github.com/golang/crypto.git"; Tag=""},
    @{Name="net"; Path="vendor\golang.org\x\net"; Url="https://github.com/golang/net.git"; Tag=""},
    @{Name="sys"; Path="vendor\golang.org\x\sys"; Url="https://github.com/golang/sys.git"; Tag=""},
    @{Name="text"; Path="vendor\golang.org\x\text"; Url="https://github.com/golang/text.git"; Tag=""},
    @{Name="exp"; Path="vendor\golang.org\x\exp"; Url="https://github.com/golang/exp.git"; Tag=""},
    @{Name="protobuf"; Path="vendor\google.golang.org\protobuf"; Url="https://github.com/protocolbuffers/protobuf-go.git"; Tag="v1.31.0"},
    @{Name="ini"; Path="vendor\gopkg.in\ini.v1"; Url="https://github.com/go-ini/ini.git"; Tag="v1.67.0"}
)

# 克隆依赖包
$successCount = 0
$failCount = 0

foreach ($dep in $dependencies) {
    $name = $dep.Name
    $path = $dep.Path
    $url = $dep.Url
    $tag = $dep.Tag
    
    Write-Host "克隆 $name ..." -ForegroundColor Yellow -NoNewline
    
    if (Test-Path $path) {
        Write-Host " [已存在]" -ForegroundColor Gray
        $successCount++
        continue
    }
    
    try {
        if ($tag -ne "") {
            git clone --depth 1 --branch $tag $url $path 2>&1 | Out-Null
        } else {
            git clone --depth 1 $url $path 2>&1 | Out-Null
        }
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host " [成功]" -ForegroundColor Green
            $successCount++
        } else {
            Write-Host " [失败]" -ForegroundColor Red
            $failCount++
        }
    } catch {
        Write-Host " [失败: $($_.Exception.Message)]" -ForegroundColor Red
        $failCount++
    }
}

Write-Host ""
Write-Host "=== 克隆完成 ===" -ForegroundColor Green
Write-Host "成功: $successCount" -ForegroundColor Green
Write-Host "失败: $failCount" -ForegroundColor $(if ($failCount -gt 0) { "Red" } else { "Green" })

if ($failCount -eq 0) {
    Write-Host ""
    Write-Host "下一步：运行 .\scripts\create-modules-txt.ps1 生成 modules.txt 文件" -ForegroundColor Cyan
} else {
    Write-Host ""
    Write-Host "部分依赖克隆失败，请检查网络连接或手动下载" -ForegroundColor Yellow
}
