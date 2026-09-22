$ErrorActionPreference = 'Stop'

$Repo     = 'https://github.com/spoonman136668-ai/Wingless.git'
$Branch   = 'research/wingless-up0-unitary-routing-r1'
$Expected = '18c1f58007dcb0e42a00e10ee4a804abf3e47c86'
$Stamp    = Get-Date -Format 'yyyyMMdd_HHmmss'
$Proof    = Join-Path $env:USERPROFILE "Downloads\K_WINGLESS_UP0_$Stamp"

git clone --single-branch --branch $Branch $Repo $Proof
if ($LASTEXITCODE -ne 0) {
    throw 'WINGLESS_UP0_CLONE_FAILED'
}

$Head = (git -C $Proof rev-parse HEAD).Trim()
if ($Head -cne $Expected) {
    throw "WINGLESS_UP0_HEAD_MISMATCH expected=$Expected actual=$Head"
}

powershell -ExecutionPolicy Bypass -File (Join-Path $Proof 'scripts\Test-UP0.ps1')
if ($LASTEXITCODE -ne 0) {
    throw 'WINGLESS_UP0_ACCEPTANCE_FAILED'
}

Write-Host "`n=== POST-RUN WORKTREE ==="
git -C $Proof status --short

Write-Host "`nWINGLESS_UP0_WINDOWS_PROOF_COMPLETE"
Write-Host "Proof: $Proof"
Write-Host "HEAD:  $Head"
