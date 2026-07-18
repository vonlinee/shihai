# Shihai Poetry Platform - Windows Terminal launcher

param(
    [ValidateSet("Development", "Deployment")]
    [string]$Mode = "Deployment",
    [ValidateSet("Standalone", "Embedded")]
    [string]$FrontendMode = "Embedded",
    [bool]$Redeploy = $true,
    [string]$DeployPath = "C:\shihai-deploy",
    [string]$StandaloneApiBaseUrl = "http://localhost:8080/api",
    [ValidateNotNullOrEmpty()]
    [string]$FrontendTabTitle = "Shihai Frontend",
    [ValidateNotNullOrEmpty()]
    [string]$BackendTabTitle,
    [switch]$DryRun
)

$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ScriptsDir = Split-Path -Parent $ScriptDir
$ProjectRoot = Split-Path -Parent $ScriptsDir
$FrontendPath = Join-Path $ProjectRoot "frontend"
$BackendPath = Join-Path $ProjectRoot "backend"
$DeployScript = Join-Path $ScriptDir "deploy-windows.ps1"
$StartFrontendScript = Join-Path $ScriptDir "start-frontend.ps1"
$StartBackendScript = Join-Path $ScriptDir "start-backend.ps1"
$IsWindowsTerminalSession = -not [string]::IsNullOrWhiteSpace($env:WT_SESSION)
$TerminalWindowSelector = if ($IsWindowsTerminalSession) { "0" } else { "new" }
$TerminalWindowLabel = if ($IsWindowsTerminalSession) { "Current" } else { "New" }
$ResolvedBackendTabTitle = if ($PSBoundParameters.ContainsKey("BackendTabTitle")) {
    $BackendTabTitle
} elseif ($Mode -eq "Deployment" -and $FrontendMode -eq "Embedded") {
    "Backend (with embedded frontend)"
} else {
    "Shihai Backend"
}

function Test-Command {
    param([string]$Command)

    return $null -ne (Get-Command $Command -ErrorAction SilentlyContinue)
}

function New-TabPlan {
    param(
        [string]$Title,
        [string]$WorkingDirectory,
        [string[]]$Command
    )

    return [PSCustomObject]@{
        Title = $Title
        WorkingDirectory = $WorkingDirectory
        Command = $Command
    }
}

function Write-LaunchPlan {
    param([object[]]$Tabs)

    Write-Host "Mode: $Mode"
    Write-Host "Terminal window: $TerminalWindowLabel ($TerminalWindowSelector)"
    if ($Mode -eq "Deployment") {
        Write-Host "Frontend mode: $FrontendMode"
        Write-Host "Redeploy: $Redeploy"
        Write-Host "Deploy path: $DeployPath"
        if ($FrontendMode -eq "Standalone") {
            Write-Host "Standalone API base URL: $StandaloneApiBaseUrl"
        }
    }
    Write-Host "Tab count: $($Tabs.Count)"

    foreach ($Tab in $Tabs) {
        Write-Host "Tab: $($Tab.Title)"
        Write-Host "  Directory: $($Tab.WorkingDirectory)"
        Write-Host "  Command: $($Tab.Command -join ' ')"
    }
}

function Open-WindowsTerminal {
    param(
        [object[]]$Tabs,
        [string]$WindowSelector
    )

    $TerminalArguments = @("-w", $WindowSelector)
    for ($Index = 0; $Index -lt $Tabs.Count; $Index++) {
        if ($Index -gt 0) {
            $TerminalArguments += ";"
        }

        $Tab = $Tabs[$Index]
        $TerminalArguments += @(
            "new-tab",
            "--title", $Tab.Title,
            "-d", $Tab.WorkingDirectory
        )
        $TerminalArguments += $Tab.Command
    }

    & wt.exe @TerminalArguments
    if ($LASTEXITCODE -ne 0) {
        throw "Windows Terminal failed to open. Exit code: $LASTEXITCODE"
    }
}

function Assert-DevelopmentPrerequisites {
    foreach ($Command in @("wt.exe", "node", "npm", "go")) {
        if (-not (Test-Command $Command)) {
            throw "$Command must be installed and available in PATH."
        }
    }
}

function Assert-DeploymentArtifacts {
    $BackendExecutable = Join-Path $DeployPath "shihai-server.exe"
    $ModeFile = Join-Path $DeployPath "frontend-mode.txt"

    if (-not (Test-Path $BackendExecutable)) {
        throw "Backend deployment not found at $BackendExecutable. Enable redeployment or run deploy-windows.ps1 first."
    }
    if (-not (Test-Path $ModeFile)) {
        throw "Deployment mode marker not found at $ModeFile. Redeploy the application first."
    }

    $DeployedMode = (Get-Content -LiteralPath $ModeFile -Raw).Trim()
    if ($DeployedMode -ne $FrontendMode) {
        throw "Existing deployment uses frontend mode '$DeployedMode', but '$FrontendMode' was requested."
    }

    if ($FrontendMode -eq "Standalone") {
        $FrontendIndex = Join-Path $DeployPath "frontend\index.html"
        if (-not (Test-Path $FrontendIndex)) {
            throw "Standalone frontend deployment not found at $FrontendIndex."
        }
    }
}

if ($Mode -eq "Development") {
    $Tabs = @(
        (New-TabPlan -Title $FrontendTabTitle -WorkingDirectory $FrontendPath -Command @(
            "powershell.exe", "-NoExit", "-Command", "npm run dev"
        )),
        (New-TabPlan -Title $ResolvedBackendTabTitle -WorkingDirectory $BackendPath -Command @(
            "powershell.exe", "-NoExit", "-Command", "go run ./cmd/server"
        ))
    )
} elseif ($FrontendMode -eq "Embedded") {
    $Tabs = @(
        (New-TabPlan -Title $ResolvedBackendTabTitle -WorkingDirectory $DeployPath -Command @(
            "powershell.exe", "-NoExit", "-ExecutionPolicy", "Bypass", "-File", $StartBackendScript,
            "-DeployPath", $DeployPath
        ))
    )
} else {
    $Tabs = @(
        (New-TabPlan -Title $FrontendTabTitle -WorkingDirectory $DeployPath -Command @(
            "powershell.exe", "-NoExit", "-ExecutionPolicy", "Bypass", "-File", $StartFrontendScript,
            "-DeployPath", $DeployPath
        )),
        (New-TabPlan -Title $ResolvedBackendTabTitle -WorkingDirectory $DeployPath -Command @(
            "powershell.exe", "-NoExit", "-ExecutionPolicy", "Bypass", "-File", $StartBackendScript,
            "-DeployPath", $DeployPath
        ))
    )
}

Write-LaunchPlan -Tabs $Tabs

if ($DryRun) {
    return
}

if ($Mode -eq "Development") {
    Assert-DevelopmentPrerequisites
} else {
    if ($Redeploy) {
        & $DeployScript -FrontendMode $FrontendMode -DeployPath $DeployPath -StandaloneApiBaseUrl $StandaloneApiBaseUrl
    } else {
        Assert-DeploymentArtifacts
    }

    if (-not (Test-Command "wt.exe")) {
        throw "wt.exe must be installed and available in PATH."
    }
}

Open-WindowsTerminal -Tabs $Tabs -WindowSelector $TerminalWindowSelector
