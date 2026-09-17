param(
    [Parameter(Mandatory = $true)] [string]$LlamaRoot,
    [Parameter(Mandatory = $true)] [string]$TracerExe,
    [Parameter(Mandatory = $true)] [string]$Model,
    [Parameter(Mandatory = $true)] [string]$PromptFile,
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

foreach ($Required in @($TracerExe, $Model, $PromptFile)) {
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
$TraceDir = Split-Path -Parent $Trace
if ($TraceDir) { New-Item -ItemType Directory -Force -Path $TraceDir | Out-Null }
if (Test-Path -LiteralPath $Trace) {
    throw "Trace output already exists: $Trace"
}

& $TracerExe `
    -m $Model `
    --trace $Trace `
    --session $Session `
    --workload $Workload `
    --task $TaskFamily `
    -ngl 32 `
    -c 4096 `
    -b 256 `
    -n $NPredict `
    --prompt-file $PromptFile
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
$PromptHash = (Get-FileHash -LiteralPath $PromptFile -Algorithm SHA256).Hash.ToLowerInvariant()

$MetaPath = "$Trace.meta.json"
$Meta = [ordered]@{
    schema = 'wingless.nr1.trace-run.v1'
    recorded_at = (Get-Date).ToUniversalTime().ToString('o')
    research_only = $true
    live_wingless_activation = $false
    ckb_plane_used = $false
    llama_source_commit = $ExpectedLlamaCommit
    tracer_sha256 = $TracerHash
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
        session_id = $Session
        workload_id = $Workload
        task_family = $TaskFamily
        prompt_sha256 = $PromptHash
        n_predict = $NPredict
    }
    trace = [ordered]@{
        path = $Trace
        bytes = [int64]$TraceInfo.Length
        sha256 = $TraceHash
    }
}
$Meta | ConvertTo-Json -Depth 8 | Set-Content -LiteralPath $MetaPath -Encoding UTF8

$WinglessRoot = (Resolve-Path (Join-Path $PSScriptRoot '..\..')).Path
$Report = "$Trace.locality.json"
Push-Location $WinglessRoot
try {
    go run ./cmd/nr1a -trace $Trace -out $Report
    if ($LASTEXITCODE -ne 0) { throw 'NR-1A locality analysis failed' }
} finally {
    Pop-Location
}

Write-Host "TRACE=$Trace"
Write-Host "TRACE_SHA256=$TraceHash"
Write-Host "META=$MetaPath"
Write-Host "LOCALITY_REPORT=$Report"
