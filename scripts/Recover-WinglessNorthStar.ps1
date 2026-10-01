Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'
$Expected='efd6172772d76f03ac04bdab0af40941232311def8202d6f0f953d33b3d05dd8'
$Roots=@('C:\ProgramData\CKBR\codex\work\rs','C:\ProgramData\CKBR\research-sidecar')
$Found=$null

function Get-CanonicalHash([string]$Path){
  $Text=[IO.File]::ReadAllText($Path)
  $Norm=$Text.Replace("`r`n","`n").Replace("`r","`n")
  $Tmp=Join-Path $env:RUNNER_TEMP ('northstar-'+[Guid]::NewGuid().ToString('N')+'.txt')
  [IO.File]::WriteAllText($Tmp,$Norm,(New-Object Text.UTF8Encoding($false)))
  try{return (Get-FileHash -LiteralPath $Tmp -Algorithm SHA256).Hash.ToLowerInvariant()}finally{Remove-Item -LiteralPath $Tmp -Force -ErrorAction SilentlyContinue}
}

foreach($Root in $Roots){
  if(-not(Test-Path -LiteralPath $Root -PathType Container)){continue}
  Write-Host "SEARCH_ROOT=$Root"
  foreach($F in @(Get-ChildItem -LiteralPath $Root -Recurse -File -Filter 'wingless-lm-north-star.txt' -ErrorAction SilentlyContinue)){
    $Raw=(Get-FileHash -LiteralPath $F.FullName -Algorithm SHA256).Hash.ToLowerInvariant()
    $Norm=Get-CanonicalHash $F.FullName
    Write-Host "CANDIDATE path=$($F.FullName) raw=$Raw normalized=$Norm"
    if($Raw-ceq$Expected -or $Norm-ceq$Expected){$Found=$F.FullName;break}
  }
  if($Found){break}
}

if(-not$Found){
  $Repo='C:\ProgramData\CKBR\codex\work\rs\source\Wingless'
  if(Test-Path -LiteralPath $Repo -PathType Container){
    Write-Host 'GIT_OBJECT_SEARCH_BEGIN'
    $Objs=@(& git -C $Repo rev-list --all --objects 2>$null | Select-String -SimpleMatch 'wingless-lm-north-star')
    foreach($O in $Objs){Write-Host ([string]$O)}
    Write-Host 'GIT_OBJECT_SEARCH_END'
    foreach($O in $Objs){
      $Line=[string]$O
      $Parts=$Line -split ' ',2
      if($Parts.Count-lt2){continue}
      $Blob=$Parts[0];$Name=$Parts[1]
      $Tmp=Join-Path $env:RUNNER_TEMP ('northstar-blob-'+$Blob+'.txt')
      & git -C $Repo cat-file blob $Blob | Set-Content -LiteralPath $Tmp -Encoding UTF8
      if($LASTEXITCODE-ne0){continue}
      $Norm=Get-CanonicalHash $Tmp
      Write-Host "GIT_CANDIDATE blob=$Blob name=$Name normalized=$Norm"
      if($Norm-ceq$Expected){$Found=$Tmp;break}
    }
  }
}

if(-not$Found){throw 'FROZEN_NORTH_STAR_NOT_FOUND_BY_IDENTITY'}
$Text=[IO.File]::ReadAllText($Found)
$Norm=$Text.Replace("`r`n","`n").Replace("`r","`n")
Write-Host "FOUND_PATH=$Found"
Write-Host 'NORTH_STAR_CONTENT_BEGIN'
Write-Host $Norm
Write-Host 'NORTH_STAR_CONTENT_END'
Write-Host 'WINGLESS_NORTH_STAR_RECOVERY=PASS'
