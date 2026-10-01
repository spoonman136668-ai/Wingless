Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'
$Codex='C:\ProgramData\CKBR\codex\bin\codex.exe'
if(-not(Test-Path -LiteralPath $Codex -PathType Leaf)){throw "CODEX_MISSING path=$Codex"}
Write-Host '=== CODEX VERSION ==='
& $Codex --version
if($LASTEXITCODE-ne0){throw "CODEX_VERSION_FAILED exit=$LASTEXITCODE"}
Write-Host '=== CODEX EXEC HELP ==='
& $Codex exec --help
if($LASTEXITCODE-ne0){throw "CODEX_EXEC_HELP_FAILED exit=$LASTEXITCODE"}
Write-Host 'CODEX_EXEC_HELP_PROBE=PASS'
