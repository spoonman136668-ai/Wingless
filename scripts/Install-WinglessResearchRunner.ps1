param(
    [string]$RepositoryUrl = 'https://github.com/spoonman136668-ai/Wingless',
    [string]$InstallRoot = 'C:\actions-runner-wingless',
    [string]$RunnerName = ('WINGLESS-' + $env:COMPUTERNAME),
    [string]$Labels = 'wingless-research-safe',
    [string]$RegistrationToken
)

$ErrorActionPreference = 'Stop'

function Test-IsAdministrator {
    $identity = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($identity)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

if (-not (Test-IsAdministrator)) {
    throw 'WINGLESS_RUNNER_INSTALL_REQUIRES_ADMINISTRATOR'
}

if (Test-Path (Join-Path $InstallRoot '.runner')) {
    throw "WINGLESS_RUNNER_ALREADY_CONFIGURED: $InstallRoot"
}

if ([string]::IsNullOrWhiteSpace($RegistrationToken)) {
    $gh = Get-Command gh -ErrorAction SilentlyContinue
    if ($null -eq $gh) {
        throw 'WINGLESS_RUNNER_TOKEN_REQUIRED: pass -RegistrationToken from Settings > Actions > Runners > New self-hosted runner, or install/authenticate GitHub CLI.'
    }

    & gh auth status
    if ($LASTEXITCODE -ne 0) {
        throw 'WINGLESS_RUNNER_GH_NOT_AUTHENTICATED'
    }

    $RegistrationToken = (
        & gh api -X POST `
            repos/spoonman136668-ai/Wingless/actions/runners/registration-token `
            --jq .token
    ).Trim()

    if ($LASTEXITCODE -ne 0 -or [string]::IsNullOrWhiteSpace($RegistrationToken)) {
        throw 'WINGLESS_RUNNER_TOKEN_ACQUISITION_FAILED'
    }
}

New-Item -ItemType Directory -Force -Path $InstallRoot | Out-Null

$release = Invoke-RestMethod `
    -Headers @{ 'User-Agent' = 'Wingless-Research-Runner-Installer' } `
    -Uri 'https://api.github.com/repos/actions/runner/releases/latest'

$asset = @(
    $release.assets |
    Where-Object {
        $_.name.StartsWith('actions-runner-win-x64-') -and
        $_.name.EndsWith('.zip')
    } |
    Select-Object -First 1
)

if ($asset.Count -ne 1) {
    $available = @($release.assets | ForEach-Object { $_.name }) -join ', '
    throw "WINGLESS_RUNNER_PACKAGE_NOT_FOUND available=$available"
}

$zip = Join-Path $env:TEMP $asset[0].name
Invoke-WebRequest -Uri $asset[0].browser_download_url -OutFile $zip

try {
    Expand-Archive -Path $zip -DestinationPath $InstallRoot -Force
}
finally {
    Remove-Item -Force $zip -ErrorAction SilentlyContinue
}

$Config = Join-Path $InstallRoot 'config.cmd'
if (-not (Test-Path $Config)) {
    throw "WINGLESS_RUNNER_CONFIG_CMD_MISSING: $Config"
}

Push-Location $InstallRoot
try {
    & $Config `
        --unattended `
        --url $RepositoryUrl `
        --token $RegistrationToken `
        --name $RunnerName `
        --labels $Labels `
        --work '_work' `
        --runasservice `
        --replace

    if ($LASTEXITCODE -ne 0) {
        throw "WINGLESS_RUNNER_CONFIG_FAILED exit=$LASTEXITCODE"
    }
}
finally {
    Pop-Location
}

$runnerConfig = Join-Path $InstallRoot '.runner'
if (-not (Test-Path $runnerConfig)) {
    throw "WINGLESS_RUNNER_REGISTRATION_MISSING: $runnerConfig"
}

$matchingServices = @(
    Get-Service |
    Where-Object {
        $_.Name -like 'actions.runner.*' -and
        (
            $_.Name -like '*WINGLESS*' -or
            $_.DisplayName -like '*WINGLESS*'
        )
    }
)

if ($matchingServices.Count -lt 1) {
    throw 'WINGLESS_RUNNER_SERVICE_NOT_FOUND_AFTER_CONFIG'
}

foreach ($service in $matchingServices) {
    if ($service.Status -ne 'Running') {
        Start-Service $service.Name
    }
}

Write-Host ''
Write-Host 'WINGLESS_RESEARCH_RUNNER_INSTALL_COMPLETE'
Write-Host "InstallRoot: $InstallRoot"
Write-Host "RunnerName:  $RunnerName"
Write-Host "Labels:      $Labels"
Write-Host 'Existing KTRADE/CKB runner paths were not modified.'
Write-Host ''
Write-Host 'Wingless runner services:'
Get-Service |
    Where-Object { $_.Name -like 'actions.runner.*' } |
    Select-Object Name, Status
