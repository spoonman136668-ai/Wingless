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
function Assert-DocumentPath([string]$Path){
  if([string]::IsNullOrWhiteSpace($Path) -or $Path.Contains("\\") -or $Path.Contains("..") -or $Path.Contains("//")){
    throw "PROXY_DOCUMENT_PATH_INVALID:$Path"
  }
  $prefix="docs/experiments/"
  if(!$Path.StartsWith($prefix,[StringComparison]::Ordinal)){throw "PROXY_DOCUMENT_PATH_FORBIDDEN:$Path"}
  $leaf=$Path.Substring($prefix.Length)
  if([string]::IsNullOrWhiteSpace($leaf) -or $leaf.Contains("/") -or !$leaf.EndsWith(".json",[StringComparison]::Ordinal) -or $leaf.Length-gt185 -or $leaf -match "[^A-Za-z0-9._-]"){
    throw "PROXY_DOCUMENT_PATH_FORBIDDEN:$Path"
  }
}
function Invoke-GhJson([string]$Method,[string]$Endpoint,$Body,[string]$Jq){
  $tmp=Join-Path $env:RUNNER_TEMP ("ckb-proxy-gh-"+[Guid]::NewGuid().ToString("N")+".json")
  try{
    [IO.File]::WriteAllText($tmp,($Body|ConvertTo-Json -Depth 50 -Compress),(New-Object Text.UTF8Encoding($false)))
    $args=@("api","--method",$Method,$Endpoint,"--input",$tmp)
    if($Jq){$args+=@("--jq",$Jq)}
    $out=@(& gh.exe @args 2>&1)
    if($LASTEXITCODE-ne0){throw ("PROXY_GH_API_FAILED:"+$Method+":"+$Endpoint+":"+($out -join " "))}
    return ($out -join [Environment]::NewLine).Trim()
  }finally{Remove-Item -LiteralPath $tmp -Force -ErrorAction SilentlyContinue}
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
  Assert-DocumentPath "docs/experiments/proxy-selftest.json"
  $rejected=$false
  try{Assert-Branch "main"}catch{$rejected=$true}
  if(!$rejected){throw "PROXY_SELFTEST_MAIN_NOT_REJECTED"}
  $rejected=$false
  try{Assert-Branch "research/../main"}catch{$rejected=$true}
  if(!$rejected){throw "PROXY_SELFTEST_TRAVERSAL_NOT_REJECTED"}
  foreach($bad in @(".github/workflows/bad.yml","docs/experiments/../bad.json","docs/experiments/bad.yml","docs/experiments/sub/bad.json")){
    $rejected=$false
    try{Assert-DocumentPath $bad}catch{$rejected=$true}
    if(!$rejected){throw "PROXY_SELFTEST_DOCUMENT_PATH_NOT_REJECTED:$bad"}
  }
  if("research-r49-static.yml" -ne "research-r49-static.yml"){throw "PROXY_SELFTEST_RESEARCH_WORKFLOW_ALLOWLIST"}
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
        if($existingResponse.status -eq "PASS" -and [string]$existingResponse.request_sha256 -eq $requestSha -and [string]$existingResponse.operation -ne "fanout_dispatch"){continue}
        if($existingResponse.status -eq "FAILED" -and [string]$existingResponse.request_sha256 -eq $requestSha -and [string]$existingResponse.controller_sha -ne $controller){continue}
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
          if($LASTEXITCODE-ne0){
            $oldNativeEap=$ErrorActionPreference
            try{
              $ErrorActionPreference="Continue"
              $fetchOut=@(& git.exe -c "safe.directory=$safe" -C $RepositoryPath fetch --no-tags origin $source 2>&1)
              $fetchExit=$LASTEXITCODE
            }finally{$ErrorActionPreference=$oldNativeEap}
            if($fetchExit-ne0){throw "PROXY_SOURCE_FETCH_FAILED:"+($fetchOut -join " ")}
            & git.exe -c "safe.directory=$safe" -C $RepositoryPath cat-file -e ($source+"^{commit}")
            if($LASTEXITCODE-ne0){throw "PROXY_SOURCE_COMMIT_MISSING_AFTER_FETCH:$source"}
          }
          $before=Get-RemoteHead $branch
          if($before -eq $source){
            Write-Response $request $requestSha "PASS" "" @{remote_sha=$before}
            continue
          }
          if($before -ne $expected){throw "PROXY_EXPECTED_REMOTE_MISMATCH:${expected}:$before"}
          $oldNativeEap=$ErrorActionPreference
          try{
            $ErrorActionPreference="Continue"
            $out=@(& git.exe -c "safe.directory=$safe" -C $RepositoryPath push origin ($source+":refs/heads/"+$branch) 2>&1)
            $pushExit=$LASTEXITCODE
          }finally{
            $ErrorActionPreference=$oldNativeEap
          }
          if($pushExit-ne0){throw "PROXY_PUSH_FAILED:"+($out -join " ")}
          $after=Get-RemoteHead $branch
          if($after -ne $source){throw "PROXY_PUSH_VERIFY_MISMATCH:${source}:$after"}
          Write-Response $request $requestSha "PASS" "" @{remote_sha=$after}
        }
        "create_documents_branch" {
          $parent=[string]$request.parent_sha
          $branch=[string]$request.branch
          $expected=[string]$request.expected_remote_sha
          $message=[string]$request.commit_message
          Assert-Sha $parent "PROXY_DOCUMENT_PARENT_SHA_INVALID"
          Assert-Sha $expected "PROXY_DOCUMENT_EXPECTED_SHA_INVALID" -AllowEmpty
          Assert-Branch $branch
          if(!$branch.StartsWith("research/")){throw "PROXY_DOCUMENT_BRANCH_FORBIDDEN:$branch"}
          if([string]::IsNullOrWhiteSpace($message) -or $message.Length-gt200 -or $message.Contains([Environment]::NewLine)){throw "PROXY_DOCUMENT_COMMIT_MESSAGE_INVALID"}
          $files=@($request.files)
          if($files.Count-lt1 -or $files.Count-gt4){throw "PROXY_DOCUMENT_FILE_COUNT_INVALID:$($files.Count)"}
          $paths=[System.Collections.Generic.HashSet[string]]::new([System.StringComparer]::Ordinal)
          $entries=[System.Collections.Generic.List[object]]::new()
          foreach($doc in $files){
            $path=[string]$doc.path
            $expectedSha=[string]$doc.sha256
            $contentB64=[string]$doc.content_b64
            Assert-DocumentPath $path
            if(!$paths.Add($path)){throw "PROXY_DOCUMENT_DUPLICATE_PATH:$path"}
            if($expectedSha.Length-ne64 -or $expectedSha -match "[^0-9a-f]"){throw "PROXY_DOCUMENT_SHA256_INVALID:$path"}
            try{$bytes=[Convert]::FromBase64String($contentB64)}catch{throw "PROXY_DOCUMENT_BASE64_INVALID:$path"}
            if($bytes.Length-le0 -or $bytes.Length-gt262144){throw ("PROXY_DOCUMENT_SIZE_INVALID:"+$path+":"+$bytes.Length)}
            if((Get-Sha256Bytes $bytes)-ne$expectedSha){throw "PROXY_DOCUMENT_SHA256_MISMATCH:$path"}
            $blobBody=[ordered]@{content=$contentB64;encoding="base64"}
            $blobSha=Invoke-GhJson "POST" ("repos/"+$Repository+"/git/blobs") $blobBody ".sha"
            Assert-Sha $blobSha "PROXY_DOCUMENT_BLOB_SHA_INVALID"
            [void]$entries.Add([ordered]@{path=$path;mode="100644";type="blob";sha=$blobSha})
          }
          $parentRaw=@(& gh.exe api ("repos/"+$Repository+"/git/commits/"+$parent) 2>&1)
          if($LASTEXITCODE-ne0){throw "PROXY_DOCUMENT_PARENT_REMOTE_MISSING:"+($parentRaw -join " ")}
          $parentObj=(($parentRaw -join [Environment]::NewLine)|ConvertFrom-Json)
          if([string]$parentObj.sha-ne$parent){throw "PROXY_DOCUMENT_PARENT_REMOTE_MISMATCH"}
          $baseTree=[string]$parentObj.tree.sha
          Assert-Sha $baseTree "PROXY_DOCUMENT_BASE_TREE_INVALID"
          $treeBody=[ordered]@{base_tree=$baseTree;tree=@($entries)}
          $treeSha=Invoke-GhJson "POST" ("repos/"+$Repository+"/git/trees") $treeBody ".sha"
          Assert-Sha $treeSha "PROXY_DOCUMENT_TREE_SHA_INVALID"
          $before=Get-RemoteHead $branch
          if($before){
            $existingRaw=@(& gh.exe api ("repos/"+$Repository+"/git/commits/"+$before) 2>&1)
            if($LASTEXITCODE-ne0){throw "PROXY_DOCUMENT_EXISTING_COMMIT_READ_FAILED"}
            $existing=(($existingRaw -join [Environment]::NewLine)|ConvertFrom-Json)
            $parents=@($existing.parents)
            if($parents.Count-eq1 -and [string]$parents[0].sha-eq$parent -and [string]$existing.tree.sha-eq$treeSha){
              Write-Response $request $requestSha "PASS" "" @{remote_sha=$before;tree_sha=$treeSha}
              continue
            }
            throw "PROXY_DOCUMENT_BRANCH_EXISTS_DIFFERENT:$before"
          }
          if($expected){throw ("PROXY_DOCUMENT_EXPECTED_REMOTE_MISMATCH:"+$expected+":")}
          $commitBody=[ordered]@{message=$message;tree=$treeSha;parents=@($parent)}
          $commitSha=Invoke-GhJson "POST" ("repos/"+$Repository+"/git/commits") $commitBody ".sha"
          Assert-Sha $commitSha "PROXY_DOCUMENT_COMMIT_SHA_INVALID"
          $refBody=[ordered]@{ref=("refs/heads/"+$branch);sha=$commitSha}
          $createdRef=Invoke-GhJson "POST" ("repos/"+$Repository+"/git/refs") $refBody ".object.sha"
          if($createdRef-ne$commitSha){throw ("PROXY_DOCUMENT_REF_CREATE_MISMATCH:"+$commitSha+":"+$createdRef)}
          $after=Get-RemoteHead $branch
          if($after-ne$commitSha){throw ("PROXY_DOCUMENT_REF_VERIFY_MISMATCH:"+$commitSha+":"+$after)}
          Write-Response $request $requestSha "PASS" "" @{remote_sha=$commitSha;tree_sha=$treeSha}
        }
        "rerun" {
          $runId=[long]$request.run_id
          if($runId-le0){throw "PROXY_RUN_ID_INVALID:$runId"}
          & gh.exe run rerun ([string]$runId) --repo $Repository | Out-Null
          if($LASTEXITCODE-ne0){throw "PROXY_RERUN_FAILED:$runId"}
          Write-Response $request $requestSha "PASS" "" @{run_id=$runId}
        }
        "research_dispatch" {
          $package=[string]$request.package_sha
          $workflow=[string]$request.workflow
          Assert-Sha $package "PROXY_RESEARCH_PACKAGE_INVALID"
          if($Lane-ne"Wingless"){throw "PROXY_RESEARCH_DISPATCH_LANE_FORBIDDEN:$Lane"}
          if($workflow-ne"research-r49-static.yml"){throw "PROXY_RESEARCH_WORKFLOW_FORBIDDEN:$workflow"}
          $commitRaw=@(& gh.exe api ("repos/"+$Repository+"/git/commits/"+$package) 2>&1)
          if($LASTEXITCODE-ne0){throw "PROXY_RESEARCH_PACKAGE_REMOTE_MISSING:"+($commitRaw -join " ")}
          $commit=(($commitRaw -join [Environment]::NewLine)|ConvertFrom-Json)
          if([string]$commit.sha-ne$package){throw "PROXY_RESEARCH_PACKAGE_REMOTE_MISMATCH"}
          $title="CKB research R49 $([string]$request.request_id) $package"
          $run=$null
          $deadline=(Get-Date).AddSeconds(90)
          $dispatched=$false
          do{
            $oldNativeEap=$ErrorActionPreference
            try{
              $ErrorActionPreference="Continue"
              $listed=@(& gh.exe run list --repo $Repository --workflow $workflow --event workflow_dispatch --limit 50 --json databaseId,displayTitle,status,conclusion,createdAt)
              $listExit=$LASTEXITCODE
            }finally{$ErrorActionPreference=$oldNativeEap}
            if($listExit-ne0){throw "PROXY_RESEARCH_LIST_FAILED"}
            $rows=@((($listed -join [Environment]::NewLine)|ConvertFrom-Json))
            $run=$rows|Where-Object{
              $d=$_.PSObject.Properties["displayTitle"]
              $d -and [string]$d.Value -ceq $title
            }|Sort-Object {
              $created=$_.PSObject.Properties["createdAt"]
              if($created){[datetime]$created.Value}else{[datetime]::MinValue}
            } -Descending|Select-Object -First 1
            if($run){break}
            if(!$dispatched){
              $oldNativeEap=$ErrorActionPreference
              try{
                $ErrorActionPreference="Continue"
                $dispatchOut=@(& gh.exe workflow run $workflow --repo $Repository --ref main -f ("package_sha="+$package) -f ("request_id="+[string]$request.request_id) 2>&1)
                $dispatchExit=$LASTEXITCODE
              }finally{$ErrorActionPreference=$oldNativeEap}
              if($dispatchExit-ne0){throw "PROXY_RESEARCH_DISPATCH_FAILED:"+($dispatchOut -join " ")}
              $dispatched=$true
            }
            Start-Sleep -Seconds 2
          }while((Get-Date)-lt$deadline)
          if(!$run){throw "PROXY_RESEARCH_RUN_DISCOVERY_TIMEOUT"}
          $runId=[long]$run.databaseId
          if($runId-le0){throw "PROXY_RESEARCH_RUN_ID_INVALID"}
          Write-Response $request $requestSha "PASS" "" @{
            run_id=$runId
            hosted_status=[string]$run.status
            conclusion=[string]$run.conclusion
            workflow=$workflow
            package_sha=$package
          }
        }
        "fanout_dispatch" {
          $package=[string]$request.package_sha
          $manifestRef=[string]$request.manifest_ref
          $manifestPath=[string]$request.manifest_path
          Assert-Sha $package "PROXY_FANOUT_PACKAGE_INVALID"
          Assert-Sha $manifestRef "PROXY_FANOUT_MANIFEST_REF_INVALID"
          if($manifestPath -notmatch '^research/fanout/[A-Za-z0-9._-]+\.json$' -or $manifestPath.Contains("..")){throw "PROXY_FANOUT_MANIFEST_PATH_INVALID"}
          $title="CKB fanout $Lane $package $manifestRef"
          $run=$null
          $deadline=(Get-Date).AddSeconds(60)
          $dispatched=$false
          do{
            $oldNativeEap=$ErrorActionPreference
            try{
              $ErrorActionPreference="Continue"
              $listed=@(& gh.exe run list --repo $Repository --workflow free-fanout-sidecar.yml --event workflow_dispatch --limit 50 --json databaseId,displayTitle,status,conclusion,createdAt)
              $listExit=$LASTEXITCODE
            }finally{$ErrorActionPreference=$oldNativeEap}
            if($listExit-ne0){throw "PROXY_FANOUT_LIST_FAILED"}
            $rows=@((($listed -join [Environment]::NewLine)|ConvertFrom-Json))
            $run=$rows|Where-Object{
              $d=$_.PSObject.Properties["displayTitle"]
              $d -and [string]$d.Value -ceq $title
            }|Sort-Object {
              $c=$_.PSObject.Properties["createdAt"]
              if($c){[datetime]$c.Value}else{[datetime]::MinValue}
            } -Descending|Select-Object -First 1
            if($run){break}
            if(!$dispatched){
              & gh.exe workflow run free-fanout-sidecar.yml --repo $Repository --ref main -f ("package_sha="+$package) -f ("manifest_ref="+$manifestRef) -f ("manifest_path="+$manifestPath) | Out-Null
              if($LASTEXITCODE-ne0){throw "PROXY_FANOUT_DISPATCH_FAILED"}
              $dispatched=$true
            }
            Start-Sleep -Seconds 2
          }while((Get-Date)-lt$deadline)
          if(!$run){throw "PROXY_FANOUT_RUN_DISCOVERY_TIMEOUT"}
          $runId=[long]$run.databaseId
          if($runId-le0){throw "PROXY_FANOUT_RUN_ID_INVALID"}
          Write-Response $request $requestSha "PASS" "" @{run_id=$runId;fanout_status=[string]$run.status;conclusion=[string]$run.conclusion}
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
  $global:LASTEXITCODE=0
  exit 0
