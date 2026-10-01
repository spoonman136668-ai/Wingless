Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'

$Codex='C:\ProgramData\CKBR\codex\bin\codex.exe'
$CodexHome='C:\ProgramData\CKBR\codex\home'
if(-not(Test-Path -LiteralPath $Codex -PathType Leaf)){throw 'CODEX_BINARY_MISSING'}
if(-not(Test-Path -LiteralPath (Join-Path $CodexHome 'auth.json') -PathType Leaf)){throw 'CODEX_AUTH_MISSING'}

$env:CODEX_HOME=$CodexHome
$Root=Join-Path $env:RUNNER_TEMP 'repo-local-codex-smoke'
if(Test-Path -LiteralPath $Root){Remove-Item -LiteralPath $Root -Recurse -Force}
New-Item -ItemType Directory -Force -Path $Root|Out-Null
$Prompt=Join-Path $Root 'prompt.txt'
$Out=Join-Path $Root 'last.txt'
[IO.File]::WriteAllText($Prompt,'Output exactly REPO_LOCAL_CODEX_OK. Do not inspect or modify any files.',(New-Object Text.UTF8Encoding($false)))

Push-Location $Root
try{
  Get-Content -LiteralPath $Prompt -Raw | & $Codex exec -m gpt-5.6-luna --sandbox read-only --output-last-message $Out -
  if($LASTEXITCODE-ne0){throw "CODEX_SMOKE_FAILED exit=$LASTEXITCODE"}
}finally{Pop-Location}

$Last=(Get-Content -LiteralPath $Out -Raw).Trim()
Write-Host "codex_last=$Last"
if($Last -cne 'REPO_LOCAL_CODEX_OK'){throw "CODEX_SMOKE_OUTPUT_MISMATCH actual=$Last"}
Write-Host 'REPO_LOCAL_CODEX_SMOKE=PASS'
