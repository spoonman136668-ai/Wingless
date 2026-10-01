Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'

$ExpectedSha='93a3da439a2ec6f8b5eef4b6193b26a11351294b'
$ExpectedTree='d25bc25d3ba52e9ebb20b0b5552fa559974ee50c'
$Cache='C:\ProgramData\CKBR\codex\work\rs\workspaces\wingless-baseline-cache'
$TaskName='CKBPlane Research Sidecar'
$SidecarExe='C:\ProgramData\CKBR\research-sidecar\bin\research-sidecar.exe'
$Branch='research/wingless-wlm-si-dense-sparse-r2-r1'
$Remote='https://github.com/spoonman136668-ai/Wingless.git'

if([string]::IsNullOrWhiteSpace($env:GITHUB_TOKEN)){throw 'GITHUB_TOKEN_MISSING'}
if(-not(Test-Path -LiteralPath $Cache -PathType Container)){throw "CACHE_MISSING path=$Cache"}

$ActualSha=(& git -C $Cache rev-parse $ExpectedSha).Trim()
if($LASTEXITCODE-ne0 -or $ActualSha-cne$ExpectedSha){throw "CACHE_SHA_MISMATCH actual=$ActualSha"}
$ActualTree=(& git -C $Cache rev-parse "$ExpectedSha^{tree}").Trim()
if($LASTEXITCODE-ne0 -or $ActualTree-cne$ExpectedTree){throw "CACHE_TREE_MISMATCH expected=$ExpectedTree actual=$ActualTree"}

$Pair="x-access-token:$env:GITHUB_TOKEN"
$Basic=[Convert]::ToBase64String([Text.Encoding]::ASCII.GetBytes($Pair))
$Auth="http.https://github.com/.extraheader=AUTHORIZATION: basic $Basic"
$Existing=@(& git -c $Auth ls-remote --heads $Remote "refs/heads/$Branch")
if($LASTEXITCODE-ne0){throw 'REMOTE_BRANCH_PROBE_FAILED'}
if($Existing.Count-gt0){throw "REMOTE_BRANCH_ALREADY_EXISTS branch=$Branch"}

$Task=Get-ScheduledTask -TaskName $TaskName -ErrorAction SilentlyContinue
if($null-eq$Task){throw "SIDECAR_TASK_MISSING name=$TaskName"}
Disable-ScheduledTask -TaskName $TaskName -ErrorAction Stop | Out-Null
try{Stop-ScheduledTask -TaskName $TaskName -ErrorAction Stop}catch{}
Start-Sleep -Seconds 2
foreach($P in @(Get-CimInstance Win32_Process -ErrorAction SilentlyContinue | Where-Object {
    ([string]$_.ExecutablePath).Trim() -ieq $SidecarExe
})){
    Stop-Process -Id $P.ProcessId -Force -ErrorAction Stop
}
Start-Sleep -Seconds 1
$Task=Get-ScheduledTask -TaskName $TaskName -ErrorAction Stop
if([string]$Task.State -eq 'Running'){throw 'SIDECAR_TASK_STILL_RUNNING'}
$Procs=@(Get-CimInstance Win32_Process -ErrorAction SilentlyContinue | Where-Object {
    ([string]$_.ExecutablePath).Trim() -ieq $SidecarExe
})
if($Procs.Count-ne0){throw "SIDECAR_PROCESS_STILL_RUNNING count=$($Procs.Count)"}
Write-Host "WINGLESS_SIDECAR_DISABLED state=$($Task.State)"

$Work=Join-Path $env:RUNNER_TEMP 'wingless-repo-local-cutover-r1'
if(Test-Path -LiteralPath $Work){Remove-Item -LiteralPath $Work -Recurse -Force}
try{
    & git -C $Cache worktree add --detach $Work $ExpectedSha
    if($LASTEXITCODE-ne0){throw 'WORKTREE_ADD_FAILED'}
    $Head=(& git -C $Work rev-parse HEAD).Trim()
    $Tree=(& git -C $Work rev-parse 'HEAD^{tree}').Trim()
    if($Head-cne$ExpectedSha -or $Tree-cne$ExpectedTree){throw 'WORKTREE_IDENTITY_MISMATCH'}

    $ReqPath=Join-Path $Work '.wingless\qualification-request.json'
    $Req=Get-Content -LiteralPath $ReqPath -Raw|ConvertFrom-Json
    if([string]$Req.schema-cne'wingless.research-qualification-request.v1'){throw 'REQUEST_SCHEMA_INVALID'}
    if([string]$Req.branch-cne$Branch){throw "REQUEST_BRANCH_UNEXPECTED actual=$($Req.branch)"}
    if([string]$Req.experiment-cne'WLM-SI-DENSE-SPARSE-R2'){throw "REQUEST_EXPERIMENT_UNEXPECTED actual=$($Req.experiment)"}
    $Req.automation_revision='repo-local-cutover-r1'
    $Json=$Req|ConvertTo-Json -Depth 20 -Compress
    [IO.File]::WriteAllText($ReqPath,$Json,(New-Object Text.UTF8Encoding($false)))

    $WfPath=Join-Path $Work '.github\workflows\research-qualify-windows.yml'
    $Wf=[IO.File]::ReadAllText($WfPath)
    if($Wf-notmatch [regex]::Escape("'research/wingless-up**'")){throw 'WORKFLOW_OLD_FILTER_MISSING'}
    $Wf=$Wf.Replace("'research/wingless-up**'","'research/wingless-*'")
    [IO.File]::WriteAllText($WfPath,$Wf,(New-Object Text.UTF8Encoding($false)))

    & git -C $Work config user.name 'wingless-research-bot'
    & git -C $Work config user.email 'actions@users.noreply.github.com'
    & git -C $Work add -- '.wingless/qualification-request.json' '.github/workflows/research-qualify-windows.yml'
    $Changed=@(& git -C $Work diff --cached --name-only)
    if($Changed.Count-ne2 -or $Changed-notcontains'.wingless/qualification-request.json' -or $Changed-notcontains'.github/workflows/research-qualify-windows.yml'){
        throw "CUTOVER_SCOPE_INVALID files=$($Changed -join ',')"
    }
    & git -C $Work commit -m 'infra: cut over Wingless research to repo-local continuation'
    if($LASTEXITCODE-ne0){throw 'CUTOVER_COMMIT_FAILED'}
    $CutoverHead=(& git -C $Work rev-parse HEAD).Trim()

    & git -C $Work -c $Auth push $Remote "HEAD:refs/heads/$Branch"
    if($LASTEXITCODE-ne0){throw 'CUTOVER_PUSH_FAILED'}

    $RemoteLine=@(& git -c $Auth ls-remote --heads $Remote "refs/heads/$Branch")
    if($LASTEXITCODE-ne0 -or $RemoteLine.Count-ne1){throw 'CUTOVER_REMOTE_VERIFY_FAILED'}
    $RemoteSha=(([string]$RemoteLine[0])-split '\s+')[0]
    if($RemoteSha-cne$CutoverHead){throw "CUTOVER_REMOTE_SHA_MISMATCH expected=$CutoverHead actual=$RemoteSha"}

    Write-Host 'WINGLESS_REPO_LOCAL_CUTOVER=PASS'
    Write-Host "sealed_sha=$ExpectedSha"
    Write-Host "sealed_tree=$ExpectedTree"
    Write-Host "cutover_head=$CutoverHead"
    Write-Host "branch=$Branch"
}
finally{
    if(Test-Path -LiteralPath $Work){
        & git -C $Cache worktree remove --force $Work 2>$null
    }
}
