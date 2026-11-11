$ErrorActionPreference = 'Stop'

$Repo     = 'https://github.com/spoonman136668-ai/Wingless.git'
$Branch   = 'research/wingless-up13-learned-phase-frame-r1'
$Baseline = '20ae7845042af43f60634926c96347d80e214265'
$Stamp    = Get-Date -Format 'yyyyMMdd_HHmmss'
$Proof    = Join-Path $env:USERPROFILE "Downloads\K_WINGLESS_UP13_$Stamp"

git clone --single-branch --branch $Branch $Repo $Proof
if ($LASTEXITCODE -ne 0) {
    throw 'WINGLESS_UP13_CLONE_FAILED'
}

$Head = (git -C $Proof rev-parse HEAD).Trim()

git -C $Proof merge-base --is-ancestor $Baseline $Head
if ($LASTEXITCODE -ne 0) {
    throw "WINGLESS_UP13_BASELINE_ANCESTRY_FAILED baseline=$Baseline head=$Head"
}

powershell -ExecutionPolicy Bypass -File (Join-Path $Proof 'scripts\Test-UP13.ps1')
if ($LASTEXITCODE -ne 0) {
    throw 'WINGLESS_UP13_ACCEPTANCE_FAILED'
}

Write-Host "`n=== POST-RUN WORKTREE ==="
git -C $Proof status --short

Write-Host "`nWINGLESS_UP13_WINDOWS_PROOF_COMPLETE"
Write-Host "Proof:    $Proof"
Write-Host "Baseline: $Baseline"
Write-Host "HEAD:     $Head"
