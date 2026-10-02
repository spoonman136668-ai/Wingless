Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'
Write-Host ("NOW_UTC="+[DateTime]::UtcNow.ToString('o'))
$Rows=@()
foreach($P in @(Get-CimInstance Win32_Process -ErrorAction SilentlyContinue)){
  $Cmd=[string]$P.CommandLine
  $Exe=[string]$P.ExecutablePath
  if($Cmd -match '(?i)Invoke-NemotronRepoAgent|Invoke-RepoLocalResearchCycle|openrouter|codex exec|research-continuation' -or
     $Exe -match '(?i)codex|powershell|pwsh'){
    $Start=$null
    try{$Start=(Get-Process -Id $P.ProcessId -ErrorAction Stop).StartTime.ToUniversalTime().ToString('o')}catch{}
    $Safe=$Cmd
    if($Safe -match '(?i)(bearer|authorization|token|secret|password|api[_-]?key)'){$Safe='<redacted-sensitive-command-line>'}
    if($Safe.Length-gt1000){$Safe=$Safe.Substring(0,1000)}
    $Rows+=[pscustomobject]@{pid=$P.ProcessId;ppid=$P.ParentProcessId;start_utc=$Start;exe=$Exe;cmd=$Safe}
  }
}
$Rows|Sort-Object start_utc,pid|ForEach-Object{
  Write-Host ("PROC pid={0} ppid={1} start={2} exe={3} cmd={4}" -f $_.pid,$_.ppid,$_.start_utc,$_.exe,$_.cmd)
}
Write-Host ("MATCH_COUNT="+$Rows.Count)
Write-Host 'RESEARCH_PROCESS_PROBE=PASS'
