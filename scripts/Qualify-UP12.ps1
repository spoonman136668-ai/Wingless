$ErrorActionPreference = 'Stop'

$Repo     = 'https://github.com/spoonman136668-ai/Wingless.git'
$Branch   = 'research/wingless-up12-fullrank-integration-r1'
$Baseline = '48db0c6fe4180a71a2f1a118eaea83691736c35f'
$Stamp    = Get-Date -Format 'yyyyMMdd_HHmmss'
$Proof    = Join-Path $env:USERPROFILE "Downloads\K_WINGLESS_UP12_$Stamp"

git clone --single-branch --branch $Branch $Repo $Proof
if ($LASTEXITCODE -ne 0) {
    throw 'WINGLESS_UP12_CLONE_FAILED'
}

$Head = (git -C $Proof rev-parse HEAD).Trim()

git -C $Proof merge-base --is-ancestor $Baseline $Head
if ($LASTEXITCODE -ne 0) {
    throw "WINGLESS_UP12_BASELINE_ANCESTRY_FAILED baseline=$Baseline head=$Head"
}

powershell -ExecutionPolicy Bypass -File (Join-Path $Proof 'scripts\Test-UP12.ps1')
if ($LASTEXITCODE -ne 0) {
    throw 'WINGLESS_UP12_ACCEPTANCE_FAILED'
}

Write-Host "`n=== POST-RUN WORKTREE ==="
git -C $Proof status --short

Write-Host "`nWINGLESS_UP12_WINDOWS_PROOF_COMPLETE"
Write-Host "Proof:    $Proof"
Write-Host "Baseline: $Baseline"
Write-Host "HEAD:     $Head"
