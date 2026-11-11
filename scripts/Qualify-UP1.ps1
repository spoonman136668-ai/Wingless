$ErrorActionPreference = 'Stop'

$Repo     = 'https://github.com/spoonman136668-ai/Wingless.git'
$Branch   = 'research/wingless-up1-trained-propagator-r1'
$Baseline = '9cd058dd43bd042d0dc03df4a35b69ecbd5f36e8'
$Stamp    = Get-Date -Format 'yyyyMMdd_HHmmss'
$Proof    = Join-Path $env:USERPROFILE "Downloads\K_WINGLESS_UP1_$Stamp"

git clone --single-branch --branch $Branch $Repo $Proof
if ($LASTEXITCODE -ne 0) {
    throw 'WINGLESS_UP1_CLONE_FAILED'
}

$Head = (git -C $Proof rev-parse HEAD).Trim()

git -C $Proof merge-base --is-ancestor $Baseline $Head
if ($LASTEXITCODE -ne 0) {
    throw "WINGLESS_UP1_BASELINE_ANCESTRY_FAILED baseline=$Baseline head=$Head"
}

powershell -ExecutionPolicy Bypass -File (Join-Path $Proof 'scripts\Test-UP1.ps1')
if ($LASTEXITCODE -ne 0) {
    throw 'WINGLESS_UP1_ACCEPTANCE_FAILED'
}

Write-Host "`n=== POST-RUN WORKTREE ==="
git -C $Proof status --short

Write-Host "`nWINGLESS_UP1_WINDOWS_PROOF_COMPLETE"
Write-Host "Proof:    $Proof"
Write-Host "Baseline: $Baseline"
Write-Host "HEAD:     $Head"
