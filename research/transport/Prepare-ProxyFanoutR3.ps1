$ErrorActionPreference='Stop'
Set-StrictMode -Version 3.0
$p='research/transport/Invoke-CKBResearchWriteProxy.ps1'
$raw=Get-Content -LiteralPath $p -Raw
$oldSkip='        if($existingResponse.status -eq "PASS" -and [string]$existingResponse.request_sha256 -eq $requestSha){continue}'
$newSkip='        if($existingResponse.status -eq "PASS" -and [string]$existingResponse.request_sha256 -eq $requestSha -and [string]$existingResponse.operation -ne "fanout_dispatch"){continue}'
if(([regex]::Matches($raw,[regex]::Escape($oldSkip))).Count-ne1){throw 'PROXY_R3_SKIP_SITE'}
$raw=$raw.Replace($oldSkip,$newSkip)
$old=@'
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
'@
$new=@'
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
'@
if(([regex]::Matches($raw,[regex]::Escape($old))).Count-ne1){throw 'PROXY_R3_FANOUT_SITE'}
$raw=$raw.Replace($old,$new)
[IO.File]::WriteAllText((Resolve-Path $p),$raw,(New-Object Text.UTF8Encoding($false)))
