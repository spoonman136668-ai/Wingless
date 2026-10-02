Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'

Write-Host '=== ENV NAMES ==='
Get-ChildItem Env: | Where-Object {$_.Name -match '(?i)openrouter|nemotron|nvidia'} | ForEach-Object {
  Write-Host ("ENV_PRESENT name={0} nonempty={1}" -f $_.Name,(-not[string]::IsNullOrWhiteSpace([string]$_.Value)))
}

Write-Host '=== MACHINE/USER ENV NAMES ==='
foreach($Scope in @('Machine','User')){
  try{
    $Vars=[Environment]::GetEnvironmentVariables($Scope)
    foreach($K in $Vars.Keys){
      if([string]$K -match '(?i)openrouter|nemotron|nvidia'){
        $V=[string]$Vars[$K]
        Write-Host ("ENV_PRESENT scope={0} name={1} nonempty={2}" -f $Scope,$K,(-not[string]::IsNullOrWhiteSpace($V)))
      }
    }
  }catch{}
}

Write-Host '=== CONFIG REFERENCES ==='
$Roots=@('C:\ProgramData\CKBR')
$Exts=@('.json','.yaml','.yml','.toml','.ps1','.cmd','.bat','.txt','.md','.ini','.conf','.config','.env')
$Count=0
foreach($Root in $Roots){
  if(-not(Test-Path -LiteralPath $Root -PathType Container)){continue}
  foreach($F in @(Get-ChildItem -LiteralPath $Root -File -Recurse -ErrorAction SilentlyContinue)){
    if($F.Length -gt 2097152){continue}
    if($Exts -notcontains $F.Extension.ToLowerInvariant() -and $F.Name -notmatch '(?i)env|config'){continue}
    try{
      $Lines=Get-Content -LiteralPath $F.FullName -ErrorAction Stop
      for($I=0;$I-lt$Lines.Count;$I++){
        $Line=[string]$Lines[$I]
        if($Line -notmatch '(?i)openrouter|nemotron|nvidia'){continue}
        $Safe=$Line
        if($Safe -match '(?i)(key|token|secret|password|credential)'){
          $Safe=[regex]::Replace($Safe,'(?i)([:=]\s*)([^\s#]+)','$1<redacted>')
          $Safe=[regex]::Replace($Safe,'(?i)("(?:api_)?key"\s*:\s*")[^"]+','$1<redacted>')
        }
        if($Safe.Length-gt400){$Safe=$Safe.Substring(0,400)}
        Write-Host ("MATCH path={0} line={1} text={2}" -f $F.FullName,($I+1),$Safe)
        $Count++
      }
    }catch{}
  }
}
Write-Host "MATCH_COUNT=$Count"

Write-Host '=== CANDIDATE EXECUTABLES ==='
foreach($Root in @('C:\ProgramData\CKBR')){
  if(-not(Test-Path -LiteralPath $Root)){continue}
  Get-ChildItem -LiteralPath $Root -File -Recurse -ErrorAction SilentlyContinue |
    Where-Object {$_.Name -match '(?i)nemotron|openrouter|router|planner|agent'} |
    Select-Object -First 100 |
    ForEach-Object { Write-Host ("FILE path={0} bytes={1}" -f $_.FullName,$_.Length) }
}
Write-Host 'NEMOTRON_CONFIG_PROBE=PASS'
