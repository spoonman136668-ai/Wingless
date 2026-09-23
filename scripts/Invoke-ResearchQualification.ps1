param(
    [Parameter(Mandatory = $true)]
    [string]$RequestPath,

    [string]$OutputDirectory = 'evidence/research-qualification',

    [int]$GuardWaitMinutes = 120
)

$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$RequestFull = Join-Path $Repo $RequestPath
$OutputFull = Join-Path $Repo $OutputDirectory

if (-not (Test-Path $RequestFull)) {
    throw "WINGLESS_QUAL_REQUEST_MISSING: $RequestFull"
}

$Request = Get-Content -Raw -Path $RequestFull | ConvertFrom-Json
if ($Request.schema -cne 'wingless.research-qualification-request.v1') {
    throw "WINGLESS_QUAL_REQUEST_SCHEMA_MISMATCH: $($Request.schema)"
}

foreach ($field in @('experiment', 'branch', 'baseline_sha', 'test_script')) {
    if ([string]::IsNullOrWhiteSpace([string]$Request.$field)) {
        throw "WINGLESS_QUAL_REQUEST_FIELD_MISSING: $field"
    }
}

$Head = (git -C $Repo rev-parse HEAD).Trim()
if ($LASTEXITCODE -ne 0) {
    throw 'WINGLESS_QUAL_HEAD_READ_FAILED'
}

if (-not [string]::IsNullOrWhiteSpace($env:GITHUB_SHA) -and $Head -cne $env:GITHUB_SHA) {
    throw "WINGLESS_QUAL_GITHUB_SHA_MISMATCH expected=$env:GITHUB_SHA actual=$Head"
}

$Branch = if (-not [string]::IsNullOrWhiteSpace($env:GITHUB_REF_NAME)) {
    $env:GITHUB_REF_NAME
}
else {
    (git -C $Repo branch --show-current).Trim()
}

if (-not [string]::IsNullOrWhiteSpace($Branch) -and $Branch -cne [string]$Request.branch) {
    throw "WINGLESS_QUAL_BRANCH_MISMATCH expected=$($Request.branch) actual=$Branch"
}

git -C $Repo merge-base --is-ancestor ([string]$Request.baseline_sha) $Head
if ($LASTEXITCODE -ne 0) {
    throw "WINGLESS_QUAL_BASELINE_ANCESTRY_FAILED baseline=$($Request.baseline_sha) head=$Head"
}

$TestScript = Join-Path $Repo ([string]$Request.test_script)
if (-not (Test-Path $TestScript)) {
    throw "WINGLESS_QUAL_TEST_SCRIPT_MISSING: $TestScript"
}

New-Item -ItemType Directory -Force -Path $OutputFull | Out-Null

$TranscriptPath = Join-Path $OutputFull 'transcript.txt'
$ProbePath = Join-Path $OutputFull 'probe.json'
$SummaryPath = Join-Path $OutputFull 'summary.json'
$StatusPath = Join-Path $OutputFull 'git-status.txt'

# This executes in the current PowerShell process. Its BelowNormal process
# priority and GOMAXPROCS settings therefore apply to the qualification and its
# child Go processes.
& (Join-Path $Repo 'scripts\Test-WinglessHostGuard.ps1') -WaitMinutes $GuardWaitMinutes

Write-Host "=== AUTOMATED WINGLESS RESEARCH QUALIFICATION ==="
Write-Host "Experiment: $($Request.experiment)"
Write-Host "Branch:     $($Request.branch)"
Write-Host "Baseline:   $($Request.baseline_sha)"
Write-Host "HEAD:       $Head"
Write-Host "Test:       $($Request.test_script)"

$TestExit = 0
& powershell -ExecutionPolicy Bypass -File $TestScript 2>&1 |
    Tee-Object -FilePath $TranscriptPath
$TestExit = $LASTEXITCODE

$Transcript = Get-Content -Raw -Path $TranscriptPath

$Probe = $null
$ProbeParseError = $null
$Match = [regex]::Match(
    $Transcript,
    '(?s)(\{\s*"schema".*?\r?\n\})\s*\r?\n\s*=== SCIENTIFIC DIAGNOSIS'
)

if ($Match.Success) {
    try {
        $Probe = $Match.Groups[1].Value | ConvertFrom-Json
        $Probe | ConvertTo-Json -Depth 100 | Set-Content -Encoding UTF8 -Path $ProbePath
    }
    catch {
        $ProbeParseError = $_.Exception.Message
    }
}
else {
    $ProbeParseError = 'No probe JSON block was found before the scientific-diagnosis marker.'
}

$StatusLines = @(
    git -C $Repo status --short
)
$StatusLines | Set-Content -Encoding UTF8 -Path $StatusPath

$AllowedDirty = @(
    if ($null -ne $Request.allowed_dirty_paths) {
        foreach ($path in $Request.allowed_dirty_paths) {
            [string]$path
        }
    }
)

$UnexpectedDirty = @(
    foreach ($line in $StatusLines) {
        if ([string]::IsNullOrWhiteSpace($line)) {
            continue
        }

        $path = if ($line.Length -gt 3) {
            $line.Substring(3).Trim()
        }
        else {
            $line.Trim()
        }

        if ($AllowedDirty -notcontains $path) {
            $line
        }
    }
)

$HarnessMarker = $Transcript -match 'WINGLESS_UP\d+_HARNESS_PASS'
$ScientificDiagnosis = if ($null -ne $Probe -and $null -ne $Probe.diagnosis) {
    $Probe.diagnosis
}
else {
    $null
}

$Classification = if ($TestExit -ne 0) {
    'harness-or-artifact-failure'
}
elseif ($UnexpectedDirty.Count -gt 0) {
    'unexpected-worktree-mutation'
}
elseif (-not $HarnessMarker) {
    'missing-harness-marker'
}
elseif ($null -eq $Probe) {
    'probe-parse-failure'
}
else {
    'qualified-scientific-result'
}

$Summary = [ordered]@{
    schema = 'wingless.research-qualification-summary.v1'
    experiment = [string]$Request.experiment
    branch = [string]$Request.branch
    baseline_sha = [string]$Request.baseline_sha
    head_sha = $Head
    test_script = [string]$Request.test_script
    platform = 'windows-self-hosted'
    runner_name = $env:RUNNER_NAME
    runner_os = $env:RUNNER_OS
    runner_arch = $env:RUNNER_ARCH
    github_run_id = $env:GITHUB_RUN_ID
    github_run_attempt = $env:GITHUB_RUN_ATTEMPT
    test_exit_code = $TestExit
    harness_marker = [bool]$HarnessMarker
    probe_parsed = [bool]($null -ne $Probe)
    probe_parse_error = $ProbeParseError
    classification = $Classification
    scientific_diagnosis = $ScientificDiagnosis
    allowed_dirty_paths = $AllowedDirty
    unexpected_dirty = $UnexpectedDirty
    production_priority_guard = 'passed'
    process_priority = [string][System.Diagnostics.Process]::GetCurrentProcess().PriorityClass
    gomaxprocs = $env:GOMAXPROCS
    generated_at_utc = [DateTime]::UtcNow.ToString('o')
}

$Summary | ConvertTo-Json -Depth 100 |
    Set-Content -Encoding UTF8 -Path $SummaryPath

Write-Host ""
Write-Host "WINGLESS_AUTOMATED_QUALIFICATION_CLASSIFICATION: $Classification"
Write-Host "WINGLESS_AUTOMATED_QUALIFICATION_SUMMARY: $SummaryPath"

if ($TestExit -ne 0) {
    throw "WINGLESS_AUTOMATED_QUALIFICATION_TEST_FAILED exit=$TestExit"
}
if ($UnexpectedDirty.Count -gt 0) {
    throw "WINGLESS_AUTOMATED_QUALIFICATION_UNEXPECTED_DIRT: $($UnexpectedDirty -join '; ')"
}
if (-not $HarnessMarker) {
    throw 'WINGLESS_AUTOMATED_QUALIFICATION_HARNESS_MARKER_MISSING'
}
if ($null -eq $Probe) {
    throw "WINGLESS_AUTOMATED_QUALIFICATION_PROBE_PARSE_FAILED: $ProbeParseError"
}

Write-Host 'WINGLESS_AUTOMATED_QUALIFICATION_PASS'
