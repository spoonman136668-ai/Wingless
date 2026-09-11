$ErrorActionPreference = "Stop"

$Repo = Split-Path $PSScriptRoot -Parent
Set-Location $Repo

$ExpectedPlaneSource = "539d56fec273c8851aea131cf4d31425a62f7250"
$EvidenceDir = Join-Path $Repo "evidence"
$EvidencePath = Join-Path $EvidenceDir "stage3-5-isolated-plane-integration.json"
New-Item -ItemType Directory -Path $EvidenceDir -Force | Out-Null

function Invoke-Checked {
    param(
        [Parameter(Mandatory=$true)][string]$Label,
        [Parameter(Mandatory=$true)][scriptblock]$Command
    )
    Write-Host "[STAGE3.5] $Label"
    & $Command
    if ($LASTEXITCODE -ne 0) {
        throw "$Label failed"
    }
}

Write-Host "[STAGE3.5] verify Windows toolchain"
$GoVersion = (& go version).Trim()
if ($LASTEXITCODE -ne 0) { throw "Go unavailable" }
if ($GoVersion -ne "go version go1.25.5 windows/amd64") {
    throw "WINDOWS_TOOLCHAIN_MISMATCH: $GoVersion"
}

Write-Host "[STAGE3.5] verify frozen plane identity is pinned in inactive adapter"
$AdapterPath = Join-Path $Repo "integration\ckbplane\adapter.go"
if (-not (Test-Path -LiteralPath $AdapterPath)) { throw "adapter missing" }
$AdapterText = Get-Content -LiteralPath $AdapterPath -Raw
if ($AdapterText -notmatch [regex]::Escape($ExpectedPlaneSource)) {
    throw "FROZEN_PLANE_IDENTITY_MISMATCH"
}

Write-Host "[STAGE3.5] verify proof source is present"
$ProofPath = Join-Path $Repo "integration\ckbplane\proof_test.go"
if (-not (Test-Path -LiteralPath $ProofPath)) { throw "isolated proof missing" }
$ProofText = Get-Content -LiteralPath $ProofPath -Raw
foreach ($Name in @(
    "TestIsolatedProofCandidateFlow",
    "TestIsolatedProofAuthorizedDeepFallback",
    "TestIsolatedProofBackendFailureNoDuplicateRetry",
    "TestIsolatedProofCancellationPropagation"
)) {
    if ($ProofText -notmatch [regex]::Escape($Name)) {
        throw "ISOLATED_PROOF_CASE_MISSING: $Name"
    }
}

$FocusedLog = Join-Path $EvidenceDir "stage3-5-isolated-focused.txt"
$StackLog = Join-Path $EvidenceDir "stage3-5-isolated-stack.txt"
$FullLog = Join-Path $EvidenceDir "stage3-5-isolated-full.txt"
$BuildLog = Join-Path $EvidenceDir "stage3-5-isolated-build.txt"

Invoke-Checked "run isolated adapter/proof tests" {
    go test -buildvcs=false ./integration/ckbplane -count=1 -timeout 30s 2>&1 |
        Tee-Object -FilePath $FocusedLog
}

Invoke-Checked "run reconciled Wingless stack packages" {
    go test -buildvcs=false ./broker ./inference ./resources ./worker ./contextbroker ./integration/ckbplane -count=1 -timeout 60s 2>&1 |
        Tee-Object -FilePath $StackLog
}

Invoke-Checked "run full Wingless repository suite" {
    go test -buildvcs=false ./... -count=1 -timeout 120s 2>&1 |
        Tee-Object -FilePath $FullLog
}

Invoke-Checked "build Wingless command" {
    go build -buildvcs=false ./cmd/wingless 2>&1 |
        Tee-Object -FilePath $BuildLog
}

$Head = ""
if (Get-Command git -ErrorAction SilentlyContinue) {
    $Head = (& git rev-parse HEAD 2>$null).Trim()
    if ($LASTEXITCODE -ne 0) { $Head = "" }
}

$Evidence = [ordered]@{
    schema = "wingless.stage3_5.isolated-plane-integration.v1"
    observed_at_utc = (Get-Date).ToUniversalTime().ToString("o")
    wingless_commit = $Head
    frozen_plane_source_commit = $ExpectedPlaneSource
    go_version = $GoVersion
    focused_adapter_proof = "PASS"
    reconciled_stack_packages = "PASS"
    full_wingless_suite = "PASS"
    wingless_build = "PASS"
    candidate_acceptance = "external_required"
    wingless_work_order_retry_budget = 0
    live_plane_transport_used = $false
    live_plane_queue_mutated = $false
    accepted_refs_mutated = $false
    ckb_runtime_launched = $false
    coinbase_or_broker_used = $false
    credentials_read = $false
    production_state_touched = $false
}

$Evidence | ConvertTo-Json -Depth 6 | Set-Content -LiteralPath $EvidencePath -Encoding UTF8

Write-Host ""
Write-Host "WINGLESS STAGE 3.5 ISOLATED INTEGRATION PROOF PASS"
Write-Host "wingless_commit=$Head"
Write-Host "frozen_plane_source_commit=$ExpectedPlaneSource"
Write-Host "focused_adapter_proof=PASS"
Write-Host "reconciled_stack_packages=PASS"
Write-Host "full_wingless_suite=PASS"
Write-Host "wingless_build=PASS"
Write-Host "candidate_acceptance=external_required"
Write-Host "wingless_work_order_retry_budget=0"
Write-Host "live_plane_transport_used=False"
Write-Host "evidence=$EvidencePath"
