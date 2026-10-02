Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'
Write-Host ("NOW_UTC="+[DateTime]::UtcNow.ToString('o'))
$Rows=@()
foreach($P in @(Get-CimInstance Win32_Process -ErrorAction SilentlyContinue)){
  $Cmd=[string]$P.CommandLine
  $Exe=[string]$P.ExecutablePath
  if($Cmd -match '(?i)Invoke-NemotronRepoAgent|Invoke-RepoLocalResearchCycle|openrouter|codex exec|research-continuation' -or
     $Exe -match '(?i)codex|powershell|pwsh'){
    $Start=$null;$Cpu=$null;$WorkingSet=$null
    try{
      $GP=Get-Process -Id $P.ProcessId -ErrorAction Stop
      $Start=$GP.StartTime.ToUniversalTime().ToString('o')
      $Cpu=[Math]::Round([double]$GP.CPU,3)
      $WorkingSet=[Math]::Round([double]$GP.WorkingSet64/1MB,1)
    }catch{}
    $Safe=$Cmd
    if($Safe -match '(?i)(bearer|authorization|token|secret|password|api[_-]?key)'){$Safe='<redacted-sensitive-command-line>'}
    if($Safe.Length-gt1000){$Safe=$Safe.Substring(0,1000)}
    $Rows+=[pscustomobject]@{pid=$P.ProcessId;ppid=$P.ParentProcessId;start_utc=$Start;cpu_seconds=$Cpu;working_set_mb=$WorkingSet;exe=$Exe;cmd=$Safe}
  }
}
$Rows|Sort-Object start_utc,pid|ForEach-Object{
  Write-Host ("PROC pid={0} ppid={1} start={2} cpu_s={3} ws_mb={4} exe={5} cmd={6}" -f $_.pid,$_.ppid,$_.start_utc,$_.cpu_seconds,$_.working_set_mb,$_.exe,$_.cmd)
}
Write-Host ("MATCH_COUNT="+$Rows.Count)
Write-Host 'RESEARCH_PROCESS_PROBE=PASS'
