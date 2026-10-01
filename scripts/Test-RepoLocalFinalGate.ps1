Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'

function Get-CanonicalTextSha256([string]$Path) {
  $Text=[IO.File]::ReadAllText($Path)
  $Normalized=$Text.Replace("`r`n","`n").Replace("`r","`n")
  $Utf8NoBom=New-Object Text.UTF8Encoding($false)
  $Bytes=$Utf8NoBom.GetBytes($Normalized)
  $Sha=[Security.Cryptography.SHA256]::Create()
  try{return ([BitConverter]::ToString($Sha.ComputeHash($Bytes))).Replace('-','').ToLowerInvariant()}finally{$Sha.Dispose()}
}

$Controller=(Resolve-Path 'controller').Path
$Research=(Resolve-Path 'research').Path
$Expected='584a292dc50fc7dbef24805b29e426f2d650f426'
$ExpectedNorth='efd6172772d76f03ac04bdab0af40941232311def8202d6f0f953d33b3d05dd8'

$C=Get-Content -Raw -LiteralPath (Join-Path $Controller '.research-autonomy\config.json')|ConvertFrom-Json
if($C.schema-cne'research.autonomy.v2'){throw 'CONFIG_SCHEMA_INVALID'}
if([bool]$C.enabled){throw 'CONFIG_MUST_BE_DISABLED_DURING_FINAL_GATE'}
if($C.ckb_plane_role-cne'governance-only'){throw 'CKB_PLANE_ROLE_INVALID'}
if($C.active_branch-cne'research/autonomy-wingless-active'){throw 'ACTIVE_BRANCH_INVALID'}
if($C.north_star_sha256-cne$ExpectedNorth){throw 'NORTH_STAR_CONFIG_IDENTITY_INVALID'}

$Path=Join-Path $Controller 'scripts\Invoke-RepoLocalResearchCycle.ps1'
$Tokens=$null;$Errors=$null
[void][Management.Automation.Language.Parser]::ParseFile($Path,[ref]$Tokens,[ref]$Errors)
if(@($Errors).Count){throw "CONTROLLER_PARSE_FAILED count=$(@($Errors).Count)"}

$Head=(& git -C $Research rev-parse HEAD).Trim()
if($Head-cne$Expected){throw "ACTIVE_HEAD_DRIFT expected=$Expected actual=$Head"}
if(((& git -C $Research status --porcelain)-join'').Trim()){throw 'ACTIVE_WORKTREE_DIRTY'}

$NorthPath=Join-Path $Research 'research\bootstrap\wingless-lm-north-star.txt'
$North=Get-CanonicalTextSha256 $NorthPath
if($North-cne$ExpectedNorth){throw "NORTH_STAR_DRIFT expected=$ExpectedNorth actual=$North"}

foreach($TaskName in @('CKBPlane Research Sidecar','CKBPlane Research Sidecar Yggdrasil')){
  $T=Get-ScheduledTask -TaskName $TaskName -ErrorAction Stop
  Write-Host "TASK name=$TaskName state=$($T.State)"
  if([string]$T.State-cne'Disabled'){throw "RESEARCH_SIDECAR_NOT_DISABLED name=$TaskName state=$($T.State)"}
}
$Codex='C:\ProgramData\CKBR\codex\bin\codex.exe'
$Auth='C:\ProgramData\CKBR\codex\home\auth.json'
if(-not(Test-Path -LiteralPath $Codex -PathType Leaf)){throw 'CODEX_MISSING'}
if(-not(Test-Path -LiteralPath $Auth -PathType Leaf)){throw 'CODEX_AUTH_MISSING'}
& $Codex --version
if($LASTEXITCODE-ne0){throw 'CODEX_VERSION_FAILED'}

Push-Location $Research
try{
  go test ./...
  if($LASTEXITCODE-ne0){throw "WINGLESS_FULL_REGRESSION_FAILED exit=$LASTEXITCODE"}
}finally{Pop-Location}

Write-Host 'WINGLESS_REPO_LOCAL_FINAL_GATE=PASS'
Write-Host "active_head=$Head"
Write-Host "north_star_sha256=$North"
