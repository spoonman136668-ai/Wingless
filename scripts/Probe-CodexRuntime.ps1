Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'
$Bins=@(
 'C:\ProgramData\CKBR\codex\bin\codex.exe',
 'C:\ProgramData\CKBR\research-sidecar-yggdrasil\codex\bin\codex.exe'
)
foreach($P in $Bins){
  if(-not(Test-Path -LiteralPath $P -PathType Leaf)){continue}
  Write-Host "=== CODEX $P ==="
  & $P --version
  Write-Host "--- top help ---"
  & $P --help
  Write-Host "--- features help ---"
  & $P features --help 2>&1
  Write-Host "--- features list ---"
  & $P features list 2>&1
  Write-Host "--- exec help ---"
  & $P exec --help 2>&1
  $Dir=Split-Path -Parent $P
  Write-Host "--- bin listing $Dir ---"
  Get-ChildItem -LiteralPath $Dir -Force | ForEach-Object {
    Write-Host ("BIN name={0} bytes={1}" -f $_.Name,$_.Length)
  }
}
foreach($Home in @(
 'C:\ProgramData\CKBR\codex\home',
 'C:\ProgramData\CKBR\research-sidecar-yggdrasil\codex\home'
)){
  if(-not(Test-Path -LiteralPath $Home -PathType Container)){continue}
  Write-Host "=== HOME $Home ==="
  foreach($Cfg in @('config.toml','config.json')){
    $Path=Join-Path $Home $Cfg
    if(-not(Test-Path -LiteralPath $Path -PathType Leaf)){continue}
    Write-Host "CONFIG path=$Path"
    foreach($Line in Get-Content -LiteralPath $Path){
      if($Line -match '(?i)(token|secret|password|api[_-]?key|credential)'){continue}
      Write-Host $Line
    }
  }
}
Write-Host 'CODEX_RUNTIME_DIAGNOSTIC=PASS'
