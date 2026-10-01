Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'
$Path='C:\ProgramData\CKBR\codex\work\rs\source\Wingless\research\bootstrap\wingless-lm-north-star.txt'
$Expected='efd6172772d76f03ac04bdab0af40941232311def8202d6f0f953d33b3d05dd8'
if(-not(Test-Path -LiteralPath $Path -PathType Leaf)){throw "NORTH_STAR_MISSING path=$Path"}
$Raw=(Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
$Text=[IO.File]::ReadAllText($Path)
$Norm=$Text.Replace("`r`n","`n").Replace("`r","`n")
$Tmp=Join-Path $env:RUNNER_TEMP 'wingless-north-star-normalized.txt'
[IO.File]::WriteAllText($Tmp,$Norm,(New-Object Text.UTF8Encoding($false)))
$Normalized=(Get-FileHash -LiteralPath $Tmp -Algorithm SHA256).Hash.ToLowerInvariant()
Write-Host "raw_sha256=$Raw"
Write-Host "normalized_sha256=$Normalized"
if($Raw-cne$Expected -and $Normalized-cne$Expected){throw "NORTH_STAR_IDENTITY_MISMATCH expected=$Expected raw=$Raw normalized=$Normalized"}
Write-Host 'NORTH_STAR_CONTENT_BEGIN'
Write-Host $Norm
Write-Host 'NORTH_STAR_CONTENT_END'
Write-Host 'WINGLESS_NORTH_STAR_RECOVERY=PASS'
