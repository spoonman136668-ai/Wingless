Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'
$Expected='efd6172772d76f03ac04bdab0af40941232311def8202d6f0f953d33b3d05dd8'
$Roots=@(
  'C:\ProgramData\CKBR',
  'C:\Users\camar\GolandProjects'
)
$Candidates=New-Object Collections.Generic.List[object]
foreach($Root in $Roots){
  if(-not(Test-Path -LiteralPath $Root -PathType Container)){continue}
  Write-Host "SCAN_ROOT=$Root"
  foreach($F in @(Get-ChildItem -LiteralPath $Root -File -Recurse -ErrorAction SilentlyContinue | Where-Object {
    $_.Length -le 1048576 -and (
      $_.Name -ieq 'wingless-lm-north-star.txt' -or
      $_.FullName -match '(?i)north.star|bootstrap|research'
    )
  })){
    $Candidates.Add($F)
  }
}
$Seen=@{}
$Found=$false
foreach($F in $Candidates){
  if($Seen.ContainsKey($F.FullName)){continue}
  $Seen[$F.FullName]=$true
  try{
    $H=(Get-FileHash -LiteralPath $F.FullName -Algorithm SHA256).Hash.ToLowerInvariant()
    if($H -ne $Expected){continue}
    $Found=$true
    Write-Host "MATCH path=$($F.FullName) bytes=$($F.Length) sha256=$H"
    Write-Host 'CONTENT_BEGIN'
    Get-Content -LiteralPath $F.FullName -Raw
    Write-Host 'CONTENT_END'
  }catch{}
}
if(-not$Found){
  Write-Host "CANDIDATES_SCANNED=$($Seen.Count)"
  throw 'NORTH_STAR_HASH_NOT_FOUND_GLOBAL'
}
Write-Host 'NORTH_STAR_HASH_GLOBAL_RECOVERY=PASS'
