$ErrorActionPreference = "Stop"

$TestDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$WindowsScriptsDir = Split-Path -Parent $TestDir
$ProjectRoot = Split-Path -Parent (Split-Path -Parent $WindowsScriptsDir)
$DeployScript = Join-Path $WindowsScriptsDir "deploy-windows.ps1"
$FrontendAssetsPath = Join-Path $ProjectRoot "frontend\dist\assets"
$ExpectedApiBaseUrl = "http://localhost:8080/api"

& $DeployScript -FrontendMode Standalone -BuildOnly -SkipBackend

$ApiReference = Get-ChildItem -Path $FrontendAssetsPath -Filter *.js |
    Select-String -SimpleMatch $ExpectedApiBaseUrl |
    Select-Object -First 1

if (-not $ApiReference) {
    throw "Standalone frontend bundle does not contain API base URL '$ExpectedApiBaseUrl'."
}

Write-Host "Standalone frontend API base URL test passed." -ForegroundColor Green
