# Windows PowerShell 一键部署脚本
# 适用于内网环境，自动克隆依赖并构建项目

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "  vLLM Proxy 离线部署脚本" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""

# 检查 Go 环境
Write-Host "检查 Go 环境..." -ForegroundColor Yellow
try {
    $goVersion = go version
    Write-Host "✓ $goVersion" -ForegroundColor Green
} catch {
    Write-Host "错误: 未找到 Go 环境" -ForegroundColor Red
    Write-Host "请先安装 Go 1.21 或更高版本" -ForegroundColor Yellow
    exit 1
}

# 检查 Git
Write-Host "检查 Git 环境..." -ForegroundColor Yellow
try {
    $gitVersion = git --version
    Write-Host "✓ $gitVersion" -ForegroundColor Green
} catch {
    Write-Host "错误: 未找到 Git" -ForegroundColor Red
    Write-Host "请先安装 Git" -ForegroundColor Yellow
    exit 1
}

Write-Host ""

# 步骤 1: 克隆依赖
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "步骤 1/3: 克隆依赖包" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""

if ((Test-Path "vendor") -and (Test-Path "vendor\modules.txt")) {
    Write-Host "vendor 目录已存在，跳过克隆步骤" -ForegroundColor Gray
    Write-Host "如需重新克隆，请先删除 vendor 目录" -ForegroundColor Gray
} else {
    if (Test-Path "scripts\clone-deps.ps1") {
        & ".\scripts\clone-deps.ps1"
    } else {
        Write-Host "错误: 找不到 scripts\clone-deps.ps1" -ForegroundColor Red
        exit 1
    }
}

Write-Host ""

# 步骤 2: 生成 modules.txt
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "步骤 2/3: 生成 modules.txt" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""

if (Test-Path "vendor\modules.txt") {
    $lineCount = (Get-Content "vendor\modules.txt").Count
    Write-Host "modules.txt 已存在" -ForegroundColor Gray
    Write-Host "文件大小: $lineCount 行" -ForegroundColor Gray
} else {
    if (Test-Path "scripts\create-modules-txt.ps1") {
        & ".\scripts\create-modules-txt.ps1"
    } else {
        Write-Host "错误: 找不到 scripts\create-modules-txt.ps1" -ForegroundColor Red
        exit 1
    }
}

Write-Host ""

# 步骤 3: 构建项目
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "步骤 3/3: 构建项目" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""

# 创建 bin 目录
if (-not (Test-Path "bin")) {
    New-Item -ItemType Directory -Path "bin" -Force | Out-Null
}

# 构建参数
$VERSION = if ($env:VERSION) { $env:VERSION } else { "1.0.0" }
$BUILD_TIME = Get-Date -Format "yyyy-MM-dd_HH:mm:ss"
$LDFLAGS = "-ldflags `"-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME`""

Write-Host "构建版本: $VERSION" -ForegroundColor Yellow
Write-Host "构建时间: $BUILD_TIME" -ForegroundColor Yellow
Write-Host ""

# 执行构建
Write-Host "开始构建..." -ForegroundColor Yellow

$buildCommand = "go build -mod=vendor $LDFLAGS -o bin\vllm-proxy.exe .\cmd"

try {
    Invoke-Expression $buildCommand
    
    if (Test-Path "bin\vllm-proxy.exe") {
        Write-Host ""
        Write-Host "=========================================" -ForegroundColor Green
        Write-Host "  ✓ 构建成功！" -ForegroundColor Green
        Write-Host "=========================================" -ForegroundColor Green
        Write-Host ""
        
        # 显示二进制文件信息
        $fileInfo = Get-Item "bin\vllm-proxy.exe"
        $fileSize = [math]::Round($fileInfo.Length / 1MB, 2)
        
        Write-Host "二进制文件: bin\vllm-proxy.exe" -ForegroundColor Cyan
        Write-Host "文件大小: $fileSize MB" -ForegroundColor Cyan
        Write-Host ""
        
        # 显示使用帮助
        Write-Host "使用方法:" -ForegroundColor Yellow
        Write-Host "  .\bin\vllm-proxy.exe --help" -ForegroundColor White
        Write-Host ""
        Write-Host "配置文件:" -ForegroundColor Yellow
        Write-Host "  .\bin\vllm-proxy.exe --config configs\config.yaml" -ForegroundColor White
        Write-Host ""
        Write-Host "命令行参数:" -ForegroundColor Yellow
        Write-Host "  .\bin\vllm-proxy.exe ``" -ForegroundColor White
        Write-Host "    --host 0.0.0.0 ``" -ForegroundColor White
        Write-Host "    --port 8000 ``" -ForegroundColor White
        Write-Host "    --prefiller-hosts 10.0.0.1 10.0.0.2 ``" -ForegroundColor White
        Write-Host "    --prefiller-ports 8100 8101 ``" -ForegroundColor White
        Write-Host "    --decoder-hosts 10.0.0.3 10.0.0.4 ``" -ForegroundColor White
        Write-Host "    --decoder-ports 8200 8201" -ForegroundColor White
    } else {
        throw "构建失败：二进制文件未生成"
    }
} catch {
    Write-Host ""
    Write-Host "=========================================" -ForegroundColor Red
    Write-Host "  ✗ 构建失败" -ForegroundColor Red
    Write-Host "=========================================" -ForegroundColor Red
    Write-Host ""
    Write-Host "错误信息: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "请检查错误信息并尝试手动构建" -ForegroundColor Yellow
    exit 1
}
