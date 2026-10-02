Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'

$Repo=(Resolve-Path '.').Path
$Branch=(& git branch --show-current).Trim()
$Head=(& git rev-parse HEAD).Trim()
$Prompt=Join-Path $env:RUNNER_TEMP 'nemotron-one-round-prompt.txt'
$Last=Join-Path $env:RUNNER_TEMP 'nemotron-one-round-last.txt'
[IO.File]::WriteAllText($Prompt,'This is a runtime protocol probe. Do not modify the repository. Reply with exactly NEMOTRON_ONE_ROUND_OK.',(New-Object Text.UTF8Encoding($false)))
try{
  & .\scripts\Invoke-NemotronRepoAgent.ps1 -RepoPath $Repo -Program Wingless -ActiveBranch $Branch -StartSha $Head -PromptPath $Prompt -LastMessagePath $Last -MaxRounds 1
  if(-not(Test-Path -LiteralPath $Last -PathType Leaf)){throw 'ONE_ROUND_LAST_MESSAGE_MISSING'}
  $Text=[IO.File]::ReadAllText($Last).Trim()
  Write-Host "LAST=$Text"
  Write-Host 'NEMOTRON_ONE_ROUND=PASS'
}catch{
  Write-Host "TYPE=$($_.Exception.GetType().FullName)"
  Write-Host "MESSAGE=$($_.Exception.Message)"
  Write-Host "SCRIPT_STACK=$($_.ScriptStackTrace)"
  Write-Host "POSITION=$($_.InvocationInfo.PositionMessage)"
  throw
}
