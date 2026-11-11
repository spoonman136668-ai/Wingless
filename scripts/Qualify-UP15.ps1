$ErrorActionPreference = 'Stop'

$Repo     = 'https://github.com/spoonman136668-ai/Wingless.git'
$Branch   = 'research/wingless-up15-learned-anchor-r1'
$Baseline = 'eb8affea150f4564502145bc50607ef863b2f425'
$Stamp    = Get-Date -Format 'yyyyMMdd_HHmmss'
$Proof    = Join-Path $env:USERPROFILE "Downloads\K_WINGLESS_UP15_$Stamp"

git clone --single-branch --branch $Branch $Repo $Proof
if ($LASTEXITCODE -ne 0) {
    throw 'WINGLESS_UP15_CLONE_FAILED'
}

$Head = (git -C $Proof rev-parse HEAD).Trim()

git -C $Proof merge-base --is-ancestor $Baseline $Head
if ($LASTEXITCODE -ne 0) {
    throw "WINGLESS_UP15_BASELINE_ANCESTRY_FAILED baseline=$Baseline head=$Head"
}

powershell -ExecutionPolicy Bypass -File (Join-Path $Proof 'scripts\Test-UP15.ps1')
if ($LASTEXITCODE -ne 0) {
    throw 'WINGLESS_UP15_ACCEPTANCE_FAILED'
}

Write-Host "`n=== POST-RUN WORKTREE ==="
git -C $Proof status --short

Write-Host "`nWINGLESS_UP15_WINDOWS_PROOF_COMPLETE"
Write-Host "Proof:    $Proof"
Write-Host "Baseline: $Baseline"
Write-Host "HEAD:     $Head"
