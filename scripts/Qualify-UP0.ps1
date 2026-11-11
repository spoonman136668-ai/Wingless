$ErrorActionPreference = 'Stop'

$Repo      = 'https://github.com/spoonman136668-ai/Wingless.git'
$Branch    = 'research/wingless-up0-unitary-routing-r1'
$Baseline  = '3196880e54f89f2f55509ccf06ce627409c29bd4'
$Stamp     = Get-Date -Format 'yyyyMMdd_HHmmss'
$Proof     = Join-Path $env:USERPROFILE "Downloads\K_WINGLESS_UP0_$Stamp"

git clone --single-branch --branch $Branch $Repo $Proof
if ($LASTEXITCODE -ne 0) {
    throw 'WINGLESS_UP0_CLONE_FAILED'
}

$Head = (git -C $Proof rev-parse HEAD).Trim()

git -C $Proof merge-base --is-ancestor $Baseline $Head
if ($LASTEXITCODE -ne 0) {
    throw "WINGLESS_UP0_BASELINE_ANCESTRY_FAILED baseline=$Baseline head=$Head"
}

powershell -ExecutionPolicy Bypass -File (Join-Path $Proof 'scripts\Test-UP0.ps1')
if ($LASTEXITCODE -ne 0) {
    throw 'WINGLESS_UP0_ACCEPTANCE_FAILED'
}

Write-Host "`n=== POST-RUN WORKTREE ==="
git -C $Proof status --short

Write-Host "`nWINGLESS_UP0_WINDOWS_PROOF_COMPLETE"
Write-Host "Proof:    $Proof"
Write-Host "Baseline: $Baseline"
Write-Host "HEAD:     $Head"
