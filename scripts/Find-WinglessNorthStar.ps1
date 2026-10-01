Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'
$Expected='efd6172772d76f03ac04bdab0af40941232311def8202d6f0f953d33b3d05dd8'
$Roots=@(
  'C:\ProgramData\CKBR\research-sidecar',
  'C:\ProgramData\CKBR\codex\work\rs'
)
$Found=$false
foreach($Root in $Roots){
  if(-not(Test-Path -LiteralPath $Root -PathType Container)){continue}
  foreach($F in @(Get-ChildItem -LiteralPath $Root -File -Recurse -ErrorAction SilentlyContinue)){
    if($F.Length -gt 1048576){continue}
    try{
      $H=(Get-FileHash -LiteralPath $F.FullName -Algorithm SHA256).Hash.ToLowerInvariant()
      if($H -cne $Expected){continue}
      $Found=$true
      Write-Host "MATCH path=$($F.FullName) length=$($F.Length) sha256=$H"
      Write-Host 'CONTENT_BEGIN'
      Get-Content -LiteralPath $F.FullName -Raw
      Write-Host 'CONTENT_END'
    }catch{}
  }
}
if(-not$Found){throw 'NORTH_STAR_HASH_NOT_FOUND'}
Write-Host 'NORTH_STAR_HASH_RECOVERY=PASS'
