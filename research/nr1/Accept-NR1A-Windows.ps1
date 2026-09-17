param(
    [string]$Model = '',
    [int]$NPredict = 128
)

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

$ExpectedBranch = 'nr-1a-router-telemetry-r1'
$ExpectedLlamaCommit = '5266f24da75dc449bd56cbed7addb9c8e4a6a73e'
$ExpectedModelFile = 'Qwen3-Coder-30B-A3B-Instruct-Q3_K_M.gguf'
$ExpectedModelBytes = [int64]14711848640
$ExpectedModelSHA256 = 'f67299d72124ed68b4cbf3a776079df289bd854a17dfe6febfc835d1e1fbdefb'
$ExpectedLayers = 48
$ExpectedExpertsPerLayer = 128
$ExpectedExpertsPerToken = 8

function Fail-NR1A([string]$Message) {
    throw "NR1A_WINDOWS_R1_FAILED: $Message"
}

function Get-SHA256([string]$Path) {
    return (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
}

function Require-Command([string]$Name) {
    if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
        Fail-NR1A "required command unavailable: $Name"
    }
}

function Find-BuiltExe([string]$BuildDir, [string]$Name) {
    $Candidates = @(
        (Join-Path $BuildDir "bin\Release\$Name.exe"),
        (Join-Path $BuildDir "bin\$Name.exe"),
        (Join-Path $BuildDir "examples\trace-moe\Release\$Name.exe"),
        (Join-Path $BuildDir "examples\trace-moe\$Name.exe")
    )
    return $Candidates | Where-Object { Test-Path -LiteralPath $_ -PathType Leaf } | Select-Object -First 1
}

function Resolve-QualifiedModel([string]$Explicit) {
    if ($Explicit) {
        if (-not (Test-Path -LiteralPath $Explicit -PathType Leaf)) {
            Fail-NR1A "explicit model path unavailable: $Explicit"
        }
        return (Resolve-Path -LiteralPath $Explicit).Path
    }

    $Downloads = Join-Path $env:USERPROFILE 'Downloads'
    $Candidates = @()
    foreach ($Dir in @(Get-ChildItem -LiteralPath $Downloads -Directory -Filter 'K_QWEN3CODER_30BA3B_Q3KM_PROBE_R1_*' -ErrorAction SilentlyContinue | Sort-Object LastWriteTime -Descending)) {
        $Candidate = Join-Path $Dir.FullName "model\$ExpectedModelFile"
        if (Test-Path -LiteralPath $Candidate -PathType Leaf) {
            $Info = Get-Item -LiteralPath $Candidate
            if ($Info.Length -eq $ExpectedModelBytes) {
                $Candidates += $Candidate
            }
        }
    }
    if ($Candidates.Count -eq 0) {
        Fail-NR1A "qualified model was not found in prior probe artifacts; rerun with -Model <path-to-$ExpectedModelFile>"
    }
    foreach ($Candidate in $Candidates) {
        Write-Host "[NR-1A] verify model candidate: $Candidate"
        if ((Get-SHA256 $Candidate) -eq $ExpectedModelSHA256) {
            return (Resolve-Path -LiteralPath $Candidate).Path
        }
    }
    Fail-NR1A 'no discovered model candidate matched the qualified SHA-256'
}

if ($env:OS -ne 'Windows_NT') {
    Fail-NR1A 'Windows PowerShell is authoritative for this run'
}
if ($NPredict -lt 1 -or $NPredict -gt 2048) {
    Fail-NR1A "NPredict out of range: $NPredict"
}
foreach ($Command in @('git','go','cmake','nvidia-smi')) {
    Require-Command $Command
}

$WinglessRoot = (Resolve-Path (Join-Path $PSScriptRoot '..\..')).Path
$Branch = (& git -C $WinglessRoot branch --show-current).Trim()
if ($LASTEXITCODE -ne 0 -or $Branch -ne $ExpectedBranch) {
    Fail-NR1A "run from $ExpectedBranch; current branch is $Branch"
}
if (@(& git -C $WinglessRoot status --porcelain --untracked-files=all).Count -ne 0) {
    Fail-NR1A 'Wingless research checkout is dirty'
}
& git -C $WinglessRoot fetch origin $ExpectedBranch --quiet
if ($LASTEXITCODE -ne 0) { Fail-NR1A 'could not refresh remote research branch identity' }
$Head = (& git -C $WinglessRoot rev-parse HEAD).Trim()
$RemoteHead = (& git -C $WinglessRoot rev-parse "origin/$ExpectedBranch").Trim()
if ($LASTEXITCODE -ne 0 -or $Head -ne $RemoteHead) {
    Fail-NR1A "local research head is not the current remote branch head: local=$Head remote=$RemoteHead"
}

$ModelPath = Resolve-QualifiedModel $Model
$ModelInfo = Get-Item -LiteralPath $ModelPath
if ($ModelInfo.Name -ne $ExpectedModelFile -or $ModelInfo.Length -ne $ExpectedModelBytes) {
    Fail-NR1A 'qualified model filename/size mismatch'
}
$ModelHash = Get-SHA256 $ModelPath
if ($ModelHash -ne $ExpectedModelSHA256) {
    Fail-NR1A "qualified model hash mismatch: $ModelHash"
}

$Stamp = Get-Date -Format 'yyyyMMdd_HHmmss'
$Root = Join-Path (Join-Path $env:USERPROFILE 'Downloads') ("K_WINGLESS_NR1A_WINDOWS_R1_" + $Stamp)
$LlamaRoot = Join-Path $Root 'llama.cpp'
$OutputDir = Join-Path $Root 'results'
$SummaryPath = Join-Path $Root 'nr1a-windows-summary.json'
if (Test-Path -LiteralPath $Root) {
    Fail-NR1A "run root already exists: $Root"
}
New-Item -ItemType Directory -Path $Root | Out-Null

Write-Host '=== WINGLESS NR-1A WINDOWS R1 ==='
Write-Host "Wingless head : $Head"
Write-Host "Model         : $ModelPath"
Write-Host "Run root      : $Root"
Write-Host 'Boundary      : research-only; no CKB-plane; no live Wingless activation'
Write-Host ''

Push-Location $WinglessRoot
try {
    Write-Host '[1/6] Wingless full test suite'
    & go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { Fail-NR1A 'Wingless full test suite failed' }

    Write-Host '[2/6] Wingless full build'
    & go build ./...
    if ($LASTEXITCODE -ne 0) { Fail-NR1A 'Wingless full build failed' }
} finally {
    Pop-Location
}
if (@(& git -C $WinglessRoot status --porcelain --untracked-files=all).Count -ne 0) {
    Fail-NR1A 'tests/build mutated the Wingless research checkout'
}

Write-Host '[3/6] clone exact llama.cpp source for isolated research overlay'
& git clone --filter=blob:none --no-checkout https://github.com/ggml-org/llama.cpp.git $LlamaRoot
if ($LASTEXITCODE -ne 0) { Fail-NR1A 'llama.cpp clone failed' }
& git -C $LlamaRoot checkout --detach $ExpectedLlamaCommit
if ($LASTEXITCODE -ne 0) { Fail-NR1A 'llama.cpp exact commit checkout failed' }
$LlamaHead = (& git -C $LlamaRoot rev-parse HEAD).Trim()
if ($LlamaHead -ne $ExpectedLlamaCommit) {
    Fail-NR1A "llama.cpp identity mismatch after checkout: $LlamaHead"
}
if (@(& git -C $LlamaRoot status --porcelain --untracked-files=all).Count -ne 0) {
    Fail-NR1A 'fresh llama.cpp checkout is dirty'
}

Write-Host '[4/6] install and build passive CUDA router tracer'
$Installer = Join-Path $PSScriptRoot 'Install-NR1ATracer.ps1'
& $Installer -LlamaRoot $LlamaRoot -Build -Backend CUDA
if ($LASTEXITCODE -ne 0) { Fail-NR1A 'NR-1A tracer installer/build failed' }
$BuildDir = Join-Path $LlamaRoot 'build-nr1a'
$TracerExe = Find-BuiltExe $BuildDir 'llama-trace-moe'
$MeasureExe = Find-BuiltExe $BuildDir 'llama-measure-experts'
if (-not $TracerExe -or -not $MeasureExe) {
    Fail-NR1A 'NR-1A executables were not found after successful build'
}
$TracerHash = Get-SHA256 $TracerExe
$MeasureHash = Get-SHA256 $MeasureExe

Write-Host '[5/6] run exact-model representative routing corpus'
$CorpusRunner = Join-Path $PSScriptRoot 'Run-NR1ACorpus.ps1'
& $CorpusRunner `
    -LlamaRoot $LlamaRoot `
    -TracerExe $TracerExe `
    -ExpertMeasureExe $MeasureExe `
    -Model $ModelPath `
    -OutputDir $OutputDir `
    -NPredict $NPredict
if ($LASTEXITCODE -ne 0) { Fail-NR1A 'NR-1A corpus run failed' }

Write-Host '[6/6] verify corpus evidence and emit decision summary'
$MetaPath = Join-Path $OutputDir 'corpus.meta.json'
$LocalityPath = Join-Path $OutputDir 'corpus.locality-decode.json'
$ExpectedResidencyPath = Join-Path $OutputDir 'corpus.residency-decode-expected.json'
foreach ($Required in @($MetaPath, $LocalityPath, $ExpectedResidencyPath)) {
    if (-not (Test-Path -LiteralPath $Required -PathType Leaf)) {
        Fail-NR1A "missing required corpus evidence: $Required"
    }
}

$Meta = Get-Content -LiteralPath $MetaPath -Raw | ConvertFrom-Json
if ([string]$Meta.schema -ne 'wingless.nr1.corpus-run.v2' -or
    -not [bool]$Meta.research_only -or [bool]$Meta.ckb_plane_used -or [bool]$Meta.live_wingless_activation -or
    [int]$Meta.workload_count -ne 8) {
    Fail-NR1A 'corpus metadata boundary/schema validation failed'
}
$LocalityEnvelope = Get-Content -LiteralPath $LocalityPath -Raw | ConvertFrom-Json
if ([string]$LocalityEnvelope.schema -ne 'wingless.nr1a.report.v2' -or [string]$LocalityEnvelope.phase -ne 'decode') {
    Fail-NR1A 'decode locality report envelope invalid'
}
$Locality = $LocalityEnvelope.locality
if ([int]$Locality.sessions -ne 8 -or [int]$Locality.workloads -ne 8 -or [int]$Locality.task_families -ne 8 -or [int]$Locality.layers.Count -ne $ExpectedLayers) {
    Fail-NR1A 'decode locality report topology/corpus cardinality invalid'
}
foreach ($Layer in @($Locality.layers)) {
    if ([int]$Layer.layer -lt 0 -or [int]$Layer.layer -ge $ExpectedLayers -or [int]$Layer.unique_experts -gt $ExpectedExpertsPerLayer) {
        Fail-NR1A 'decode locality layer report invalid'
    }
}

$ResidencyEnvelope = Get-Content -LiteralPath $ExpectedResidencyPath -Raw | ConvertFrom-Json
if ([string]$ResidencyEnvelope.schema -ne 'wingless.nr1a.report.v2' -or [string]$ResidencyEnvelope.phase -ne 'decode' -or $null -eq $ResidencyEnvelope.residency) {
    Fail-NR1A 'decode residency report envelope invalid'
}
$Budget8 = [uint64](8GB)
$At8 = @($ResidencyEnvelope.residency.results | Where-Object { [uint64]$_.budget_bytes -eq $Budget8 })
$Online8 = @($At8 | Where-Object { -not [bool]$_.uses_future_trace })
$Oracle8 = @($At8 | Where-Object { [bool]$_.uses_future_trace })
if ($Online8.Count -ne 3 -or $Oracle8.Count -ne 2) {
    Fail-NR1A '8 GiB residency policy classification invalid'
}

$Coverage = [ordered]@{}
foreach ($Point in @($Locality.global_coverage)) {
    $Coverage[[string]$Point.target] = [int]$Point.experts
}
$OnlineSummary = @()
foreach ($R in $Online8) {
    $OnlineSummary += [ordered]@{
        policy = [string]$R.policy
        slots = [int]$R.slots
        hit_rate = [double]$R.hit_rate
        warm_bytes_per_token = [double]$R.warm_bytes_per_token
        cache_lifecycle = [string]$R.cache_lifecycle
    }
}
$OracleSummary = @()
foreach ($R in $Oracle8) {
    $OracleSummary += [ordered]@{
        policy = [string]$R.policy
        slots = [int]$R.slots
        hit_rate = [double]$R.hit_rate
        warm_bytes_per_token = [double]$R.warm_bytes_per_token
        uses_future_trace = [bool]$R.uses_future_trace
        cache_lifecycle = [string]$R.cache_lifecycle
    }
}

$Summary = [ordered]@{
    schema = 'wingless.nr1a.windows-summary.v1'
    recorded_at = (Get-Date).ToUniversalTime().ToString('o')
    status = 'MEASURED'
    research_only = $true
    ckb_plane_used = $false
    live_wingless_activation = $false
    wingless = [ordered]@{ branch = $ExpectedBranch; head = $Head }
    llama = [ordered]@{ commit = $LlamaHead; tracer_sha256 = $TracerHash; expert_measure_sha256 = $MeasureHash }
    model = [ordered]@{ file = $ExpectedModelFile; bytes = $ExpectedModelBytes; sha256 = $ModelHash; ngl = 32 }
    corpus = [ordered]@{
        workload_count = 8
        decode_tokens = [int]$Locality.tokens
        routed_layers = $ExpectedLayers
        experts_per_layer = $ExpectedExpertsPerLayer
        experts_per_token = $ExpectedExpertsPerToken
        adjacent_reuse_rate = [double]$Locality.adjacent_reuse_rate
        unique_expert_keys = [int]$Locality.unique_expert_keys
        global_coverage_expert_keys = $Coverage
    }
    expected_8gib = [ordered]@{
        online = $OnlineSummary
        oracle_bounds = $OracleSummary
    }
    evidence = [ordered]@{
        corpus_meta = $MetaPath
        decode_locality = $LocalityPath
        decode_residency_expected = $ExpectedResidencyPath
    }
}
$UTF8NoBOM = [Text.UTF8Encoding]::new($false)
[IO.File]::WriteAllText($SummaryPath, (($Summary | ConvertTo-Json -Depth 12) + "`r`n"), $UTF8NoBOM)
$SummaryHash = Get-SHA256 $SummaryPath

Write-Host ''
Write-Host '=== NR-1A WINDOWS R1 MEASUREMENT COMPLETE ==='
Write-Host "SUMMARY=$SummaryPath"
Write-Host "SUMMARY_SHA256=$SummaryHash"
Write-Host "DECODE_TOKENS=$($Locality.tokens)"
Write-Host ("ADJACENT_REUSE_RATE={0:N6}" -f [double]$Locality.adjacent_reuse_rate)
foreach ($R in $Online8) {
    Write-Host ("8G_ONLINE {0}: hit={1:N6} warm_bytes_per_token={2:N0} slots={3}" -f $R.policy, [double]$R.hit_rate, [double]$R.warm_bytes_per_token, [int]$R.slots)
}
foreach ($R in $Oracle8) {
    Write-Host ("8G_ORACLE {0}: hit={1:N6} warm_bytes_per_token={2:N0} slots={3}" -f $R.policy, [double]$R.hit_rate, [double]$R.warm_bytes_per_token, [int]$R.slots)
}
Write-Host 'NR1A_WINDOWS_R1_PASS'
