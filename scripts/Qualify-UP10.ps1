$ErrorActionPreference = 'Stop'

$Repo     = 'https://github.com/spoonman136668-ai/Wingless.git'
$Branch   = 'research/wingless-up10-compressed-frame-r1'
$Baseline = '8c4b2b9f108169fed4b01897c14dadb492541c78'
$Stamp    = Get-Date -Format 'yyyyMMdd_HHmmss'
$Proof    = Join-Path $env:USERPROFILE "Downloads\K_WINGLESS_UP10_$Stamp"

git clone --single-branch --branch $Branch $Repo $Proof
if ($LASTEXITCODE -ne 0) {
    throw 'WINGLESS_UP10_CLONE_FAILED'
}

$Head = (git -C $Proof rev-parse HEAD).Trim()

git -C $Proof merge-base --is-ancestor $Baseline $Head
if ($LASTEXITCODE -ne 0) {
    throw "WINGLESS_UP10_BASELINE_ANCESTRY_FAILED baseline=$Baseline head=$Head"
}

powershell -ExecutionPolicy Bypass -File (Join-Path $Proof 'scripts\Test-UP10.ps1')
if ($LASTEXITCODE -ne 0) {
    throw 'WINGLESS_UP10_ACCEPTANCE_FAILED'
}

Write-Host "`n=== POST-RUN WORKTREE ==="
git -C $Proof status --short

Write-Host "`nWINGLESS_UP10_WINDOWS_PROOF_COMPLETE"
Write-Host "Proof:    $Proof"
Write-Host "Baseline: $Baseline"
Write-Host "HEAD:     $Head"
