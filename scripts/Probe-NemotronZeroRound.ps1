Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'

$Repo=(Resolve-Path '.').Path
$Branch=(& git branch --show-current).Trim()
$Head=(& git rev-parse HEAD).Trim()
$Prompt=Join-Path $env:RUNNER_TEMP 'nemotron-tool-round-prompt.txt'
$Last=Join-Path $env:RUNNER_TEMP 'nemotron-tool-round-last.txt'
[IO.File]::WriteAllText($Prompt,'This is a runtime protocol probe. Do not modify the repository. Call git_status, inspect the result, and when you have enough information reply with exactly NEMOTRON_TOOL_ROUND_OK. Do not modify the repository.',(New-Object Text.UTF8Encoding($false)))
try{
  & .\scripts\Invoke-NemotronRepoAgent.ps1 -RepoPath $Repo -Program Wingless -ActiveBranch $Branch -StartSha $Head -PromptPath $Prompt -LastMessagePath $Last -MaxRounds 12
  if(-not(Test-Path -LiteralPath $Last -PathType Leaf)){throw 'TOOL_ROUND_LAST_MESSAGE_MISSING'}
  $Text=[IO.File]::ReadAllText($Last).Trim()
  Write-Host "LAST=$Text"
  if($Text-cne'NEMOTRON_TOOL_ROUND_OK'){throw "TOOL_ROUND_OUTPUT_INVALID output=$Text"}
  if(((& git status --porcelain)-join'').Trim()){throw 'TOOL_ROUND_MUTATED_REPO'}
  Write-Host 'NEMOTRON_TOOL_ROUND=PASS'
}catch{
  Write-Host "TYPE=$($_.Exception.GetType().FullName)"
  Write-Host "MESSAGE=$($_.Exception.Message)"
  Write-Host "SCRIPT_STACK=$($_.ScriptStackTrace)"
  Write-Host "POSITION=$($_.InvocationInfo.PositionMessage)"
  throw
}
