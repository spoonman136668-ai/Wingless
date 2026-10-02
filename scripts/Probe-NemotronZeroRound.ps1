Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'

$Repo=(Resolve-Path '.').Path
$Branch=(& git branch --show-current).Trim()
$Head=(& git rev-parse HEAD).Trim()
$Prompt=Join-Path $env:RUNNER_TEMP 'nemotron-zero-round-prompt.txt'
$Last=Join-Path $env:RUNNER_TEMP 'nemotron-zero-round-last.txt'
[IO.File]::WriteAllText($Prompt,'zero-round initialization probe',(New-Object Text.UTF8Encoding($false)))
try{
  & .\scripts\Invoke-NemotronRepoAgent.ps1 -RepoPath $Repo -Program Wingless -ActiveBranch $Branch -StartSha $Head -PromptPath $Prompt -LastMessagePath $Last -MaxRounds 0
  throw 'ZERO_ROUND_UNEXPECTED_SUCCESS'
}catch{
  Write-Host "TYPE=$($_.Exception.GetType().FullName)"
  Write-Host "MESSAGE=$($_.Exception.Message)"
  Write-Host "SCRIPT_STACK=$($_.ScriptStackTrace)"
  Write-Host "POSITION=$($_.InvocationInfo.PositionMessage)"
  if($_.Exception.Message -eq 'NEMOTRON_AGENT_MAX_ROUNDS_EXCEEDED max_rounds=0'){
    Write-Host 'NEMOTRON_ZERO_ROUND_INIT=PASS'
    exit 0
  }
  throw
}
