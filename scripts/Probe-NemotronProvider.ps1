Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'
$Root='C:\ProgramData\CKBR\codex\work\q'
$Files=@()
foreach($Name in @('research_sidecar_openrouter.go','research_planner_provider.go')){
  $Candidates=@(Get-ChildItem -LiteralPath $Root -File -Recurse -Filter $Name -ErrorAction SilentlyContinue | Sort-Object LastWriteTimeUtc -Descending)
  if($Candidates.Count-eq0){throw "PROVIDER_SOURCE_NOT_FOUND name=$Name"}
  $F=$Candidates[0]
  Write-Host "=== SOURCE name=$Name path=$($F.FullName) modified=$($F.LastWriteTimeUtc.ToString('o')) ==="
  Get-Content -LiteralPath $F.FullName -Raw
}
$Archive='C:\ProgramData\CKBR\research-sidecar-yggdrasil\state\archive\role-swap-b-r1-20261001_061127'
foreach($P in @(
  (Join-Path $Archive 'research-runtime.yaml'),
  (Join-Path $Archive 'research-sidecar.yaml'),
  (Join-Path $Archive 'Lifecycle.ps1')
)){
  if(-not(Test-Path -LiteralPath $P -PathType Leaf)){continue}
  Write-Host "=== CONFIG path=$P ==="
  foreach($Line in Get-Content -LiteralPath $P){
    if($Line -match '(?i)(secret|token|password|credential|api[_-]?key)'){
      Write-Host ([regex]::Replace($Line,'(?i)([:=]\s*)([^\s#]+)','$1<redacted>'))
    } else { Write-Host $Line }
  }
}
Write-Host 'PROVIDER_CONTRACT_PROBE=PASS'
