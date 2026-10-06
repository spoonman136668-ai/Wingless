[CmdletBinding()]
param(
  [Parameter(Mandatory=$true)][ValidateSet("Wingless","Yggdrasil")][string]$Lane,
  [Parameter(Mandatory=$true)][string]$Repository,
  [Parameter(Mandatory=$true)][string]$RepositoryPath,
  [string]$ProxyRoot = "C:\ProgramData\CKBR\research-write-proxy",
  [string]$ControllerIdentityPath = "C:\ProgramData\CKBR\research\controller-qualified-sha.txt",
  [switch]$SelfTest
)

Set-StrictMode -Version 3.0
$ErrorActionPreference="Stop"

function Get-Sha256Bytes([byte[]]$Bytes){
  $sha=[Security.Cryptography.SHA256]::Create()
  try{return ([BitConverter]::ToString($sha.ComputeHash($Bytes))).Replace("-","").ToLowerInvariant()}
  finally{$sha.Dispose()}
}
function Get-Sha256File([string]$Path){Get-Sha256Bytes ([IO.File]::ReadAllBytes($Path))}
function Write-AtomicUtf8([string]$Path,[string]$Text){
  $dir=Split-Path -Parent $Path
  [IO.Directory]::CreateDirectory($dir)|Out-Null
  $tmp=$Path+".tmp."+[Guid]::NewGuid().ToString("N")
  [IO.File]::WriteAllText($tmp,$Text,(New-Object Text.UTF8Encoding($false)))
  Move-Item -LiteralPath $tmp -Destination $Path -Force
}
function Assert-Sha([string]$Value,[string]$Code,[switch]$AllowEmpty){
  if($AllowEmpty -and [string]::IsNullOrWhiteSpace($Value)){return}
  if($Value -notmatch '^[0-9a-f]{40}$'){throw "${Code}:$Value"}
}
function Assert-Branch([string]$Branch){
  if([string]::IsNullOrWhiteSpace($Branch) -or $Branch.Contains("\") -or $Branch.Contains("..") -or $Branch.Contains("//")){
    throw "PROXY_BRANCH_INVALID:$Branch"
  }
  if($Branch -notmatch '^research/[A-Za-z0-9][A-Za-z0-9._/-]{0,180}$' -and
     $Branch -notmatch '^fanout/ckb-[a-z0-9][a-z0-9-]{0,100}$'){
    throw "PROXY_BRANCH_FORBIDDEN:$Branch"
  }
}
function Get-RemoteHead([string]$Branch){
  $safe=$RepositoryPath.Replace("\","/")
  $out=@(& git.exe -c "safe.directory=$safe" -C $RepositoryPath ls-remote --heads origin ("refs/heads/"+$Branch) 2>&1)
  if($LASTEXITCODE-ne0){throw "PROXY_REMOTE_READ_FAILED:"+($out -join " ")}
  $raw=($out -join [Environment]::NewLine).Trim()
  if(!$raw){return ""}
  $sha=($raw -split '\s+')[0]
  Assert-Sha $sha "PROXY_REMOTE_SHA_INVALID"
  return $sha
}
function Write-Response($Request,[string]$RequestSha,[string]$Status,[string]$Error="",[hashtable]$Extra=@{}){
  $response=[ordered]@{
    schema="ckb-plane.research-write-proxy-response.v1"
    request_id=[string]$Request.request_id
    request_sha256=$RequestSha
    controller_sha=[string]$Request.controller_sha
    lane=$Lane
    repository=$Repository
    operation=[string]$Request.operation
    status=$Status
    error=$Error
  }
  foreach($k in $Extra.Keys){$response[$k]=$Extra[$k]}
  $json=($response|ConvertTo-Json -Depth 20 -Compress)+[Environment]::NewLine
  $path=Join-Path (Join-Path $ProxyRoot ("responses\"+$Lane)) (([string]$Request.request_id)+".json")
  Write-AtomicUtf8 $path $json
}

if($SelfTest){
  Assert-Branch "research/proxy-selftest-r1"
  Assert-Branch "fanout/ckb-wingless-selftest"
  $rejected=$false
  try{Assert-Branch "main"}catch{$rejected=$true}
  if(!$rejected){throw "PROXY_SELFTEST_MAIN_NOT_REJECTED"}
  $rejected=$false
  try{Assert-Branch "research/../main"}catch{$rejected=$true}
  if(!$rejected){throw "PROXY_SELFTEST_TRAVERSAL_NOT_REJECTED"}
  Write-Host "CKB_RESEARCH_WRITE_PROXY_SELFTEST=PASS"
  return
}

if($Repository -notin @("spoonman136668-ai/Wingless","spoonman136668-ai/Yggdrasil")){throw "PROXY_REPOSITORY_FORBIDDEN:$Repository"}
$expectedRepository=if($Lane-eq"Wingless"){"spoonman136668-ai/Wingless"}else{"spoonman136668-ai/Yggdrasil"}
if($Repository -ne $expectedRepository){throw "PROXY_LANE_REPOSITORY_MISMATCH"}
$expectedPath="C:\ProgramData\CKBR\research\repos\"+$Lane
if([IO.Path]::GetFullPath($RepositoryPath).TrimEnd("\") -ne [IO.Path]::GetFullPath($expectedPath).TrimEnd("\")){throw "PROXY_REPOSITORY_PATH_FORBIDDEN"}
if(!(Test-Path -LiteralPath (Join-Path $RepositoryPath ".git"))){throw "PROXY_REPOSITORY_MISSING"}
if(!(Test-Path -LiteralPath $ControllerIdentityPath -PathType Leaf)){throw "PROXY_CONTROLLER_IDENTITY_MISSING"}
$controller=([IO.File]::ReadAllText($ControllerIdentityPath,[Text.Encoding]::UTF8)).Trim().ToLowerInvariant()
Assert-Sha $controller "PROXY_CONTROLLER_IDENTITY_INVALID"

if([string]::IsNullOrWhiteSpace($env:GH_TOKEN)){throw "PROXY_GH_TOKEN_MISSING"}
  & gh.exe auth status | Out-Null
  if($LASTEXITCODE-ne0){throw "PROXY_GH_AUTH_FAILED"}
  & gh.exe auth setup-git | Out-Null
  if($LASTEXITCODE-ne0){throw "PROXY_GH_GIT_SETUP_FAILED"}

  $requestDir=Join-Path $ProxyRoot ("requests\"+$Lane)
  $responseDir=Join-Path $ProxyRoot ("responses\"+$Lane)
  [IO.Directory]::CreateDirectory($requestDir)|Out-Null
  [IO.Directory]::CreateDirectory($responseDir)|Out-Null
  $requests=@(Get-ChildItem -LiteralPath $requestDir -File -Filter "*.json" -ErrorAction SilentlyContinue|Sort-Object Name)
  foreach($file in $requests){
    $requestSha=Get-Sha256File $file.FullName
    $responsePath=Join-Path $responseDir $file.Name
    if(Test-Path -LiteralPath $responsePath){
      try{
        $existingResponse=Get-Content -LiteralPath $responsePath -Raw|ConvertFrom-Json
        if($existingResponse.status -eq "PASS"){continue}
      }catch{}
    }
    $request=$null
    try{
      $request=Get-Content -LiteralPath $file.FullName -Raw|ConvertFrom-Json
      if($request.schema -ne "ckb-plane.research-write-proxy.v1"){throw "PROXY_SCHEMA_INVALID"}
      if(([string]$request.request_id+".json") -ne $file.Name){throw "PROXY_REQUEST_FILENAME_MISMATCH"}
      if([string]$request.controller_sha -ne $controller){throw "PROXY_CONTROLLER_IDENTITY_MISMATCH"}
      if([string]$request.lane -ne $Lane){throw "PROXY_LANE_MISMATCH"}
      if([string]$request.repository -ne $Repository){throw "PROXY_REPOSITORY_MISMATCH"}
      $op=[string]$request.operation
      switch($op){
        "push" {
          $source=[string]$request.source_sha
          $branch=[string]$request.branch
          $expected=[string]$request.expected_remote_sha
          Assert-Sha $source "PROXY_SOURCE_SHA_INVALID"
          Assert-Sha $expected "PROXY_EXPECTED_SHA_INVALID" -AllowEmpty
          Assert-Branch $branch
          $safe=$RepositoryPath.Replace("\","/")
          & git.exe -c "safe.directory=$safe" -C $RepositoryPath cat-file -e ($source+"^{commit}")
          if($LASTEXITCODE-ne0){throw "PROXY_SOURCE_COMMIT_MISSING:$source"}
          $before=Get-RemoteHead $branch
          if($before -eq $source){
            Write-Response $request $requestSha "PASS" "" @{remote_sha=$before}
            continue
          }
          if($before -ne $expected){throw "PROXY_EXPECTED_REMOTE_MISMATCH:${expected}:$before"}
          $out=@(& git.exe -c "safe.directory=$safe" -C $RepositoryPath push origin ($source+":refs/heads/"+$branch) 2>&1)
          if($LASTEXITCODE-ne0){throw "PROXY_PUSH_FAILED:"+($out -join " ")}
          $after=Get-RemoteHead $branch
          if($after -ne $source){throw "PROXY_PUSH_VERIFY_MISMATCH:${source}:$after"}
          Write-Response $request $requestSha "PASS" "" @{remote_sha=$after}
        }
        "rerun" {
          $runId=[long]$request.run_id
          if($runId-le0){throw "PROXY_RUN_ID_INVALID:$runId"}
          & gh.exe run rerun ([string]$runId) --repo $Repository | Out-Null
          if($LASTEXITCODE-ne0){throw "PROXY_RERUN_FAILED:$runId"}
          Write-Response $request $requestSha "PASS" "" @{run_id=$runId}
        }
        "fanout_dispatch" {
          $package=[string]$request.package_sha
          $manifestRef=[string]$request.manifest_ref
          $manifestPath=[string]$request.manifest_path
          Assert-Sha $package "PROXY_FANOUT_PACKAGE_INVALID"
          Assert-Sha $manifestRef "PROXY_FANOUT_MANIFEST_REF_INVALID"
          if($manifestPath -notmatch '^research/fanout/[A-Za-z0-9._-]+\.json$' -or $manifestPath.Contains("..")){throw "PROXY_FANOUT_MANIFEST_PATH_INVALID"}
          & gh.exe workflow run free-fanout-sidecar.yml --repo $Repository --ref main -f ("package_sha="+$package) -f ("manifest_ref="+$manifestRef) -f ("manifest_path="+$manifestPath) | Out-Null
          if($LASTEXITCODE-ne0){throw "PROXY_FANOUT_DISPATCH_FAILED"}
          Write-Response $request $requestSha "PASS"
        }
        default {throw "PROXY_OPERATION_FORBIDDEN:$op"}
      }
    }catch{
      if($null-ne$request -and $request.PSObject.Properties["request_id"] -and $request.PSObject.Properties["controller_sha"] -and $request.PSObject.Properties["operation"]){
        $msg=$_.Exception.Message
        if($msg.Length-gt512){$msg=$msg.Substring(0,512)}
        Write-Response $request $requestSha "FAILED" $msg
      }else{
        throw
      }
    }
  }
}
