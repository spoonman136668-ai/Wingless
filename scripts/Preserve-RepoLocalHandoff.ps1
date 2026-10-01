Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'

$Repo='spoonman136668-ai/Wingless'
$Cache='C:\ProgramData\CKBR\codex\work\rs\workspaces\wingless-baseline-cache'
$ExpectedSha='0bef2799dee4f78a6de2face17407f9ce2e66372'
$ExpectedExperiment='WLM-SI-QUEUE-CONTAINER-R1'

$Task=Get-ScheduledTask -TaskName 'CKBPlane Research Sidecar' -ErrorAction Stop
if([string]$Task.State -cne 'Disabled'){throw "SIDECAR_NOT_FROZEN state=$($Task.State)"}
if([string]::IsNullOrWhiteSpace($env:GITHUB_TOKEN)){throw 'GITHUB_TOKEN_MISSING'}
if(-not(Test-Path -LiteralPath $Cache -PathType Container)){throw "CACHE_MISSING path=$Cache"}

& git -C $Cache cat-file -e ($ExpectedSha+'^{commit}')
if($LASTEXITCODE-ne0){throw "EXPECTED_COMMIT_MISSING sha=$ExpectedSha"}

$RequestRaw=& git -C $Cache show ($ExpectedSha+':.wingless/qualification-request.json')
if($LASTEXITCODE-ne0){throw 'QUALIFICATION_REQUEST_READ_FAILED'}
$Q=($RequestRaw -join [Environment]::NewLine)|ConvertFrom-Json
if([string]$Q.schema -cne 'wingless.research-qualification-request.v1'){throw 'QUALIFICATION_SCHEMA_INVALID'}
if([string]$Q.experiment -cne $ExpectedExperiment){throw "EXPERIMENT_MISMATCH actual=$($Q.experiment)"}
$Branch=[string]$Q.branch
if($Branch -notmatch '^research/wingless-[a-z0-9][a-z0-9._-]{0,120}$'){throw "BRANCH_INVALID branch=$Branch"}

$Tree=(& git -C $Cache rev-parse ($ExpectedSha+'^{tree}')).Trim()
if($LASTEXITCODE-ne0){throw 'TREE_READ_FAILED'}
$Remote='https://github.com/'+$Repo+'.git'
$Pair="x-access-token:$env:GITHUB_TOKEN"
$Basic=[Convert]::ToBase64String([Text.Encoding]::ASCII.GetBytes($Pair))
$Header="AUTHORIZATION: basic $Basic"

$Before=@(& git -c "http.https://github.com/.extraheader=$Header" ls-remote --heads $Remote ("refs/heads/"+$Branch))
if($LASTEXITCODE-ne0){throw 'REMOTE_READ_FAILED'}
if($Before.Count-gt0){
  $Existing=(([string]$Before[0]) -split '\s+')[0]
  if($Existing -cne $ExpectedSha){throw "REMOTE_BRANCH_CONFLICT expected=$ExpectedSha actual=$Existing branch=$Branch"}
  Write-Host "already_preserved=true"
}else{
  & git -c "http.https://github.com/.extraheader=$Header" -C $Cache push $Remote ($ExpectedSha+":refs/heads/"+$Branch)
  if($LASTEXITCODE-ne0){throw "PUSH_FAILED branch=$Branch"}
}

$After=@(& git -c "http.https://github.com/.extraheader=$Header" ls-remote --heads $Remote ("refs/heads/"+$Branch))
if($LASTEXITCODE-ne0 -or $After.Count-ne1){throw 'REMOTE_VERIFY_FAILED'}
$Actual=(([string]$After[0]) -split '\s+')[0]
if($Actual -cne $ExpectedSha){throw "REMOTE_VERIFY_SHA_MISMATCH expected=$ExpectedSha actual=$Actual"}

Write-Host 'WINGLESS_FINAL_CUTOVER_PRESERVE=PASS'
Write-Host "experiment=$ExpectedExperiment"
Write-Host "branch=$Branch"
Write-Host "sha=$ExpectedSha"
Write-Host "tree=$Tree"
