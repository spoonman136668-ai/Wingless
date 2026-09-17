param(
    [Parameter(Mandatory = $true)] [string]$LlamaRoot,
    [Parameter(Mandatory = $true)] [string]$TracerExe,
    [Parameter(Mandatory = $true)] [string]$ExpertMeasureExe,
    [Parameter(Mandatory = $true)] [string]$Model,
    [Parameter(Mandatory = $true)] [string]$UserFile,
    [string]$SystemFile = '',
    [Parameter(Mandatory = $true)] [string]$Trace,
    [Parameter(Mandatory = $true)] [string]$Session,
    [Parameter(Mandatory = $true)] [string]$Workload,
    [Parameter(Mandatory = $true)] [string]$TaskFamily,
    [int]$NPredict = 128
)

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

$ExpectedLlamaCommit = '5266f24da75dc449bd56cbed7addb9c8e4a6a73e'
$ExpectedModelSHA256 = 'f67299d72124ed68b4cbf3a776079df289bd854a17dfe6febfc835d1e1fbdefb'
$QualifiedServerSHA256 = 'cb29f66008d4d73cce17cab2c2569ab318eb244b0b2ca2a15024b44d90dfcd3f'
$QualifiedWinglessCommit = 'af99cca38ca34b5d3c45111edebe09477e991ce1'
$KVPayloadBytes4096F16 = [uint64]402653184

$Head = (& git -C $LlamaRoot rev-parse HEAD).Trim()
if ($LASTEXITCODE -ne 0 -or $Head -ne $ExpectedLlamaCommit) {
    throw "llama.cpp identity mismatch. Expected $ExpectedLlamaCommit, got $Head"
}

$Changed = @(& git -C $LlamaRoot status --porcelain --untracked-files=all)
if ($LASTEXITCODE -ne 0) { throw 'git status failed' }
foreach ($Line in $Changed) {
    $Path = ($Line.Substring(3) -replace '\\','/').Trim('"')
    if ($Path -ne 'examples/CMakeLists.txt' -and -not $Path.StartsWith('examples/trace-moe/')) {
        throw "Unexpected llama.cpp research-checkout mutation: $Line"
    }
}

$RequiredFiles = @($TracerExe, $ExpertMeasureExe, $Model, $UserFile)
if ($SystemFile) { $RequiredFiles += $SystemFile }
foreach ($Required in $RequiredFiles) {
    if (-not (Test-Path -LiteralPath $Required -PathType Leaf)) {
        throw "Missing required file: $Required"
    }
}
if ($NPredict -lt 0 -or $NPredict -gt 2048) {
    throw "NPredict out of bounded research range: $NPredict"
}

$ModelInfo = Get-Item -LiteralPath $Model
$ModelHash = (Get-FileHash -LiteralPath $Model -Algorithm SHA256).Hash.ToLowerInvariant()
if ($ModelHash -ne $ExpectedModelSHA256) {
    throw "Model SHA256 mismatch. Expected $ExpectedModelSHA256, got $ModelHash"
}
if ($ModelInfo.Length -ne 14711848640) {
    throw "Model byte-size mismatch. Expected 14711848640, got $($ModelInfo.Length)"
}

$TracerHash = (Get-FileHash -LiteralPath $TracerExe -Algorithm SHA256).Hash.ToLowerInvariant()
$MeasureHash = (Get-FileHash -LiteralPath $ExpertMeasureExe -Algorithm SHA256).Hash.ToLowerInvariant()
$TraceDir = Split-Path -Parent $Trace
if ($TraceDir) { New-Item -ItemType Directory -Force -Path $TraceDir | Out-Null }
if (Test-Path -LiteralPath $Trace) {
    throw "Trace output already exists: $Trace"
}

$StoragePath = "$Trace.expert-storage.json"
if (Test-Path -LiteralPath $StoragePath) {
    throw "Expert storage output already exists: $StoragePath"
}
$StorageRaw = (& $ExpertMeasureExe $Model) -join "`n"
if ($LASTEXITCODE -ne 0) {
    throw "llama-measure-experts failed with exit code $LASTEXITCODE"
}
$Storage = $StorageRaw | ConvertFrom-Json
if ([string]$Storage.schema -ne 'wingless.nr1.expert-storage.v1' -or
    [string]$Storage.llama_source_commit -ne $ExpectedLlamaCommit -or
    [int]$Storage.layers -ne 48 -or [int]$Storage.experts_per_layer -ne 128) {
    throw 'Expert storage report identity/shape validation failed'
}
if ($null -eq $Storage.uniform_expert_bytes -or [uint64]$Storage.uniform_expert_bytes -le 0) {
    throw 'Checkpoint does not expose a uniform per-layer expert byte size; fixed-slot residency simulation is not valid'
}
[IO.File]::WriteAllText($StoragePath, ($StorageRaw.TrimEnd() + "`r`n"), [Text.UTF8Encoding]::new($false))
$ExpertBytes = [uint64]$Storage.uniform_expert_bytes
$NonExpertTensorBytes = [uint64]$Storage.non_expert_tensor_bytes

$TraceArgs = @(
    '-m', $Model,
    '--trace', $Trace,
    '--session', $Session,
    '--workload', $Workload,
    '--task', $TaskFamily,
    '-ngl', '32',
    '-c', '4096',
    '-b', '256',
    '-n', [string]$NPredict,
    '--user-file', $UserFile
)
if ($SystemFile) {
    $TraceArgs += @('--system-file', $SystemFile)
}

& $TracerExe @TraceArgs
if ($LASTEXITCODE -ne 0) {
    throw "llama-trace-moe failed with exit code $LASTEXITCODE"
}
if (-not (Test-Path -LiteralPath $Trace -PathType Leaf)) {
    throw "Tracer did not create output: $Trace"
}

$TraceInfo = Get-Item -LiteralPath $Trace
if ($TraceInfo.Length -le 0) {
    throw 'Trace is empty'
}
$TraceHash = (Get-FileHash -LiteralPath $Trace -Algorithm SHA256).Hash.ToLowerInvariant()
$StorageHash = (Get-FileHash -LiteralPath $StoragePath -Algorithm SHA256).Hash.ToLowerInvariant()
$UserHash = (Get-FileHash -LiteralPath $UserFile -Algorithm SHA256).Hash.ToLowerInvariant()
$SystemHash = $null
if ($SystemFile) {
    $SystemHash = (Get-FileHash -LiteralPath $SystemFile -Algorithm SHA256).Hash.ToLowerInvariant()
}

$WinglessRoot = (Resolve-Path (Join-Path $PSScriptRoot '..\..')).Path
$LocalityReport = "$Trace.locality.json"
$RuntimeScenarios = [ordered]@{
    optimistic = [uint64]536870912
    expected = [uint64]1073741824
    pessimistic = [uint64]1610612736
}
$ResidencyReports = [ordered]@{}

Push-Location $WinglessRoot
try {
    go run ./cmd/nr1a -trace $Trace -out $LocalityReport
    if ($LASTEXITCODE -ne 0) { throw 'NR-1A locality analysis failed' }

    foreach ($Scenario in $RuntimeScenarios.GetEnumerator()) {
        $Reserved = [uint64]($NonExpertTensorBytes + $KVPayloadBytes4096F16 + [uint64]$Scenario.Value)
        $ReportPath = "$Trace.residency-$($Scenario.Key).json"
        go run ./cmd/nr1a `
            -trace $Trace `
            -expert-bytes $ExpertBytes `
            -reserved-bytes $Reserved `
            -budgets-gib '4,6,8,12,16' `
            -out $ReportPath
        if ($LASTEXITCODE -ne 0) { throw "NR-1A residency simulation failed: $($Scenario.Key)" }
        $ResidencyReports[$Scenario.Key] = [ordered]@{
            path = $ReportPath
            runtime_buffer_assumption_bytes = [uint64]$Scenario.Value
            reserved_bytes = $Reserved
        }
    }
} finally {
    Pop-Location
}

$MetaPath = "$Trace.meta.json"
$Meta = [ordered]@{
    schema = 'wingless.nr1.trace-run.v1'
    recorded_at = (Get-Date).ToUniversalTime().ToString('o')
    research_only = $true
    live_wingless_activation = $false
    ckb_plane_used = $false
    llama_source_commit = $ExpectedLlamaCommit
    tracer_sha256 = $TracerHash
    expert_measure_sha256 = $MeasureHash
    model = [ordered]@{
        family = 'Qwen3-Coder-30B-A3B-Instruct'
        quantization = 'Q3_K_M'
        sha256 = $ModelHash
        bytes = [int64]$ModelInfo.Length
        ngl = 32
    }
    qualification_reference = [ordered]@{
        wingless_commit = $QualifiedWinglessCommit
        qualified_llama_server_sha256 = $QualifiedServerSHA256
        context_tokens = 4096
        note = 'Research tracer is a separate executable and does not inherit qualified server identity.'
    }
    workload = [ordered]@{
        prompt_mode = 'model_chat_template'
        session_id = $Session
        workload_id = $Workload
        task_family = $TaskFamily
        user_sha256 = $UserHash
        system_sha256 = $SystemHash
        n_predict = $NPredict
    }
    trace = [ordered]@{
        path = $Trace
        bytes = [int64]$TraceInfo.Length
        sha256 = $TraceHash
    }
    storage_measurement = [ordered]@{
        path = $StoragePath
        sha256 = $StorageHash
        expert_pool_bytes = [uint64]$Storage.expert_pool_bytes
        non_expert_tensor_bytes = $NonExpertTensorBytes
        expert_bytes = $ExpertBytes
        classification = 'measured_from_gguf_tensor_metadata'
    }
    residency_assumptions = [ordered]@{
        kv_payload_bytes = $KVPayloadBytes4096F16
        kv_classification = 'analytical_48_layers_4_kv_heads_128_head_dim_k_and_v_f16_4096_tokens_excludes_allocator_overhead'
        runtime_buffer_classification = 'hypothetical_scenarios_not_measured'
        reports = $ResidencyReports
    }
}
$Meta | ConvertTo-Json -Depth 10 | Set-Content -LiteralPath $MetaPath -Encoding UTF8

Write-Host "TRACE=$Trace"
Write-Host "TRACE_SHA256=$TraceHash"
Write-Host "EXPERT_STORAGE=$StoragePath"
Write-Host "EXPERT_BYTES=$ExpertBytes"
Write-Host "META=$MetaPath"
Write-Host "LOCALITY_REPORT=$LocalityReport"
foreach ($Scenario in $ResidencyReports.GetEnumerator()) {
    Write-Host ("RESIDENCY_{0}={1}" -f $Scenario.Key.ToUpperInvariant(), $Scenario.Value.path)
}
