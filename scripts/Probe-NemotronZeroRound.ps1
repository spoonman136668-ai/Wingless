Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'

$Repo=(Resolve-Path '.').Path
$Branch=(& git branch --show-current).Trim()
$Head=(& git rev-parse HEAD).Trim()
$Prompt=Join-Path $env:RUNNER_TEMP 'nemotron-selftest-prompt.txt'
$Last=Join-Path $env:RUNNER_TEMP 'nemotron-selftest-last.txt'
[IO.File]::WriteAllText($Prompt,'self-test',(New-Object Text.UTF8Encoding($false)))
try{
  & .\scripts\Invoke-NemotronRepoAgent.ps1 -RepoPath $Repo -Program Wingless -ActiveBranch $Branch -StartSha $Head -PromptPath $Prompt -LastMessagePath $Last -SelfTest
  if(((& git status --porcelain)-join'').Trim()){throw 'NEMOTRON_SELFTEST_MUTATED_REPO'}
  Write-Host 'NEMOTRON_LOCAL_TOOL_SELFTEST=PASS'
}catch{
  Write-Host "TYPE=$($_.Exception.GetType().FullName)"
  Write-Host "MESSAGE=$($_.Exception.Message)"
  Write-Host "SCRIPT_STACK=$($_.ScriptStackTrace)"
  Write-Host "POSITION=$($_.InvocationInfo.PositionMessage)"
  throw
}

# adapter-root-guard-selftest
