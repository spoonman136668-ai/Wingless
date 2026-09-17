param(
    [Parameter(Mandatory = $true)] [string]$LlamaRoot,
    [Parameter(Mandatory = $true)] [string]$TracerExe,
    [Parameter(Mandatory = $true)] [string]$ExpertMeasureExe,
    [Parameter(Mandatory = $true)] [string]$Model,
    [Parameter(Mandatory = $true)] [string]$OutputDir,
    [string]$Session = '',
    [int]$NPredict = 128
)

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

if (-not $Session) {
    $Session = 'nr1a-' + (Get-Date).ToUniversalTime().ToString('yyyyMMddTHHmmssZ')
}
if ($Session -notmatch '^[A-Za-z0-9][A-Za-z0-9_.:-]{0,127}$') {
    throw "Invalid NR-1A corpus id: $Session"
}
if (Test-Path -LiteralPath $OutputDir) {
    throw "Output directory already exists; refuse to mix corpus runs: $OutputDir"
}
New-Item -ItemType Directory -Path $OutputDir | Out-Null

$SingleRunner = Join-Path $PSScriptRoot 'Run-NR1ATrace.ps1'
$WorkloadRoot = Join-Path $PSScriptRoot 'workloads'
$SystemFile = Join-Path $WorkloadRoot 'system.txt'
if (-not (Test-Path -LiteralPath $SingleRunner -PathType Leaf)) { throw "Missing runner: $SingleRunner" }
if (-not (Test-Path -LiteralPath $SystemFile -PathType Leaf)) { throw "Missing system prompt: $SystemFile" }

$Workloads = @(
    [ordered]@{ id = 'code-generation';     task = 'code_generation';       file = 'code_generation.txt' },
    [ordered]@{ id = 'debugging';           task = 'debugging';             file = 'debugging.txt' },
    [ordered]@{ id = 'repair';              task = 'repair';                file = 'repair.txt' },
    [ordered]@{ id = 'repo-reasoning';      task = 'repo_reasoning';        file = 'repo_reasoning.txt' },
    [ordered]@{ id = 'planning';            task = 'planning';              file = 'planning.txt' },
    [ordered]@{ id = 'review';              task = 'review';                file = 'review.txt' },
    [ordered]@{ id = 'instruction-follow';  task = 'instruction_following'; file = 'instruction_following.txt' },
    [ordered]@{ id = 'long-coding';         task = 'long_coding';           file = 'long_coding.txt' }
)

$CorpusTrace = Join-Path $OutputDir 'corpus.jsonl'
$UTF8NoBOM = [Text.UTF8Encoding]::new($false)
$Runs = @()
$FirstStorage = $null

foreach ($Spec in $Workloads) {
    $UserFile = Join-Path $WorkloadRoot $Spec.file
    if (-not (Test-Path -LiteralPath $UserFile -PathType Leaf)) {
        throw "Missing required workload: $UserFile"
    }

    # Each tracer invocation creates a fresh model context, so it must have a
    # distinct session id. The outer $Session value is the corpus/run id only.
    $RunSession = "$Session-$($Spec.task)"
    if ($RunSession -notmatch '^[A-Za-z0-9][A-Za-z0-9_.:-]{0,127}$') {
        throw "Derived NR-1A session id is invalid or too long: $RunSession"
    }

    $Trace = Join-Path $OutputDir ($Spec.task + '.jsonl')
    & $SingleRunner `
        -LlamaRoot $LlamaRoot `
        -TracerExe $TracerExe `
        -ExpertMeasureExe $ExpertMeasureExe `
        -Model $Model `
        -UserFile $UserFile `
        -SystemFile $SystemFile `
        -Trace $Trace `
        -Session $RunSession `
        -Workload $Spec.id `
        -TaskFamily $Spec.task `
        -NPredict $NPredict
    if ($LASTEXITCODE -ne 0) { throw "NR-1A workload failed: $($Spec.task)" }

    $MetaPath = "$Trace.meta.json"
    $StoragePath = "$Trace.expert-storage.json"
    if (-not (Test-Path -LiteralPath $MetaPath -PathType Leaf)) { throw "Missing trace metadata: $MetaPath" }
    if (-not (Test-Path -LiteralPath $StoragePath -PathType Leaf)) { throw "Missing expert storage report: $StoragePath" }

    $Meta = Get-Content -LiteralPath $MetaPath -Raw | ConvertFrom-Json
    if ([string]$Meta.schema -ne 'wingless.nr1.trace-run.v2' -or
        -not [bool]$Meta.research_only -or [bool]$Meta.live_wingless_activation -or [bool]$Meta.ckb_plane_used) {
        throw "Trace metadata violated NR-1A research boundary: $MetaPath"
    }

    if ($null -eq $FirstStorage) {
        $FirstStorage = Get-Content -LiteralPath $StoragePath -Raw | ConvertFrom-Json
    } else {
        $Storage = Get-Content -LiteralPath $StoragePath -Raw | ConvertFrom-Json
        if ([uint64]$Storage.uniform_expert_bytes -ne [uint64]$FirstStorage.uniform_expert_bytes -or
            [uint64]$Storage.non_expert_tensor_bytes -ne [uint64]$FirstStorage.non_expert_tensor_bytes -or
            [uint64]$Storage.expert_pool_bytes -ne [uint64]$FirstStorage.expert_pool_bytes) {
            throw "Expert storage measurement drifted across corpus run: $($Spec.task)"
        }
    }

    $TraceText = [IO.File]::ReadAllText($Trace, $UTF8NoBOM)
    if (-not $TraceText.EndsWith("`n")) { throw "Trace is not newline terminated: $Trace" }
    [IO.File]::AppendAllText($CorpusTrace, $TraceText, $UTF8NoBOM)

    $Runs += [ordered]@{
        session_id = $RunSession
        workload_id = $Spec.id
        task_family = $Spec.task
        user_file = $Spec.file
        user_sha256 = (Get-FileHash -LiteralPath $UserFile -Algorithm SHA256).Hash.ToLowerInvariant()
        trace = [IO.Path]::GetFileName($Trace)
        trace_sha256 = (Get-FileHash -LiteralPath $Trace -Algorithm SHA256).Hash.ToLowerInvariant()
        meta = [IO.Path]::GetFileName($MetaPath)
    }
}

if ($null -eq $FirstStorage) { throw 'NR-1A corpus produced no storage measurement' }
$ExpertBytes = [uint64]$FirstStorage.uniform_expert_bytes
$NonExpertTensorBytes = [uint64]$FirstStorage.non_expert_tensor_bytes
if ($ExpertBytes -le 0) { throw 'Invalid corpus expert byte measurement' }

$WinglessRoot = (Resolve-Path (Join-Path $PSScriptRoot '..\..')).Path
$KVPayloadBytes4096F16 = [uint64]402653184
$RuntimeScenarios = [ordered]@{
    optimistic = [uint64]536870912
    expected = [uint64]1073741824
    pessimistic = [uint64]1610612736
}
$AggregateLocalityAll = Join-Path $OutputDir 'corpus.locality-all.json'
$AggregateLocalityDecode = Join-Path $OutputDir 'corpus.locality-decode.json'
$AggregateResidency = [ordered]@{}

Push-Location $WinglessRoot
try {
    go run ./cmd/nr1a -trace $CorpusTrace -phase all -out $AggregateLocalityAll
    if ($LASTEXITCODE -ne 0) { throw 'NR-1A aggregate all-phase locality analysis failed' }

    go run ./cmd/nr1a -trace $CorpusTrace -phase decode -out $AggregateLocalityDecode
    if ($LASTEXITCODE -ne 0) { throw 'NR-1A aggregate decode locality analysis failed' }

    foreach ($Scenario in $RuntimeScenarios.GetEnumerator()) {
        $Reserved = [uint64]($NonExpertTensorBytes + $KVPayloadBytes4096F16 + [uint64]$Scenario.Value)
        $ReportPath = Join-Path $OutputDir ("corpus.residency-decode-{0}.json" -f $Scenario.Key)
        go run ./cmd/nr1a `
            -trace $CorpusTrace `
            -phase decode `
            -expert-bytes $ExpertBytes `
            -reserved-bytes $Reserved `
            -budgets-gib '4,6,8,12,16' `
            -out $ReportPath
        if ($LASTEXITCODE -ne 0) { throw "NR-1A aggregate decode residency simulation failed: $($Scenario.Key)" }
        $AggregateResidency[$Scenario.Key] = [ordered]@{
            path = [IO.Path]::GetFileName($ReportPath)
            phase = 'decode'
            runtime_buffer_assumption_bytes = [uint64]$Scenario.Value
            reserved_bytes = $Reserved
        }
    }
} finally {
    Pop-Location
}

$CorpusHash = (Get-FileHash -LiteralPath $CorpusTrace -Algorithm SHA256).Hash.ToLowerInvariant()
$CorpusMeta = [ordered]@{
    schema = 'wingless.nr1.corpus-run.v2'
    recorded_at = (Get-Date).ToUniversalTime().ToString('o')
    research_only = $true
    ckb_plane_used = $false
    live_wingless_activation = $false
    corpus_id = $Session
    n_predict_per_workload = $NPredict
    workload_count = $Workloads.Count
    workloads = $Runs
    corpus_trace = [ordered]@{
        schema = 'wingless.nr1.router-trace.v2'
        path = [IO.Path]::GetFileName($CorpusTrace)
        sha256 = $CorpusHash
        bytes = [int64](Get-Item -LiteralPath $CorpusTrace).Length
        phases = @('prefill','decode')
    }
    expert_storage = [ordered]@{
        expert_bytes = $ExpertBytes
        expert_pool_bytes = [uint64]$FirstStorage.expert_pool_bytes
        non_expert_tensor_bytes = $NonExpertTensorBytes
        classification = 'measured_from_gguf_tensor_metadata'
    }
    aggregate_locality = [ordered]@{
        all = [IO.Path]::GetFileName($AggregateLocalityAll)
        decode = [IO.Path]::GetFileName($AggregateLocalityDecode)
    }
    aggregate_residency = $AggregateResidency
}
$CorpusMetaPath = Join-Path $OutputDir 'corpus.meta.json'
[IO.File]::WriteAllText($CorpusMetaPath, (($CorpusMeta | ConvertTo-Json -Depth 12) + "`r`n"), $UTF8NoBOM)

Write-Host "NR1A_CORPUS=$CorpusTrace"
Write-Host "NR1A_CORPUS_SHA256=$CorpusHash"
Write-Host "NR1A_CORPUS_META=$CorpusMetaPath"
Write-Host "NR1A_LOCALITY_ALL=$AggregateLocalityAll"
Write-Host "NR1A_LOCALITY_DECODE=$AggregateLocalityDecode"
foreach ($Scenario in $AggregateResidency.GetEnumerator()) {
    Write-Host ("NR1A_RESIDENCY_DECODE_{0}={1}" -f $Scenario.Key.ToUpperInvariant(), (Join-Path $OutputDir $Scenario.Value.path))
}
