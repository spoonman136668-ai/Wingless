Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$SharedPython='C:\ProgramData\CKBR\research-sidecar-yggdrasil\python312\python.exe'
if(-not(Test-Path -LiteralPath $SharedPython -PathType Leaf)){throw 'SHARED_PYTHON_MISSING'}

$Sidecars = @(
    [pscustomobject]@{
        Name='wing'
        Root='C:\ProgramData\CKBR\research-sidecar'
        Source='C:\ProgramData\CKBR\codex\work\rs\source\Wingless'
        Task='CKBPlane Research Sidecar'
        SearchRoots=@('C:\ProgramData\CKBR\codex\work')
    },
    [pscustomobject]@{
        Name='ygg'
        Root='C:\ProgramData\CKBR\research-sidecar-yggdrasil'
        Source='C:\ProgramData\CKBR\research-sidecar-yggdrasil\source\Yggdrasil'
        Task='CKBPlane Research Sidecar Yggdrasil'
        SearchRoots=@('C:\ProgramData\CKBR\research-sidecar-yggdrasil')
    }
)

function Read-LatestRows {
    param([string]$Name,[string]$Root)
    $Cfg=Join-Path $Root 'config\autonomy.yaml'
    $CfgRaw=[IO.File]::ReadAllText($Cfg)
    $M=[regex]::Match($CfgRaw,"(?m)^\s*sqlite_path:\s*['""]?([^\r\n'""]+)['""]?\s*$")
    if(-not$M.Success){throw "$Name sqlite_path missing"}
    $Db=$M.Groups[1].Value.Trim()
    if(-not[IO.Path]::IsPathRooted($Db)){$Db=Join-Path $Root ($Db -replace '/','\')}
    $Db=[IO.Path]::GetFullPath($Db)
    Write-Host "DB label=$Name path=$Db"
    $Probe=Join-Path $env:RUNNER_TEMP ("cutover-db-"+$Name+".py")
    $Py=@'
import json,sqlite3,sys
db=sqlite3.connect("file:"+sys.argv[1].replace("\\","/")+"?mode=ro",uri=True,timeout=5)
def dec(v):
    return v.decode("utf-8") if isinstance(v,(bytes,bytearray)) else v
def latest(t):
    cols=[r[1] for r in db.execute("pragma table_info("+t+")")]
    if not cols:return None
    row=db.execute("select * from "+t+" order by rowid desc limit 1").fetchone()
    if row is None:return None
    return {cols[i]:dec(row[i]) for i in range(len(cols))}
try:
    print(json.dumps({t:latest(t) for t in ("research_sidecar_results","research_implementation_freezes","research_cycles")},sort_keys=True))
finally:
    db.close()
'@
    [IO.File]::WriteAllText($Probe,$Py,(New-Object Text.UTF8Encoding($false)))
    $Obj=((& $SharedPython $Probe $Db)-join[Environment]::NewLine)|ConvertFrom-Json
    $Targets=New-Object Collections.Generic.List[string]
    foreach($Table in @('research_sidecar_results','research_implementation_freezes','research_cycles')){
        $Row=$Obj.$Table
        if($null-eq$Row){continue}
        $RawProp=$Row.PSObject.Properties['raw']
        if($null-eq$RawProp){continue}
        $R=([string]$RawProp.Value)|ConvertFrom-Json
        function Get-OptionalValue([object]$Object,[string]$Property){
            $P=$Object.PSObject.Properties[$Property]
            if($null-eq$P){return ''}
            return [string]$P.Value
        }
        $Cycle=Get-OptionalValue $R 'cycle_id'
        $Status=Get-OptionalValue $R 'status'
        $Experiment=Get-OptionalValue $R 'experiment_id'
        $Baseline=Get-OptionalValue $R 'baseline_sha'
        $SourceHead=Get-OptionalValue $R 'source_head'
        $SourceCommit=Get-OptionalValue $R 'source_commit'
        Write-Host "LATEST label=$Name table=$Table cycle_id=$Cycle status=$Status experiment_id=$Experiment baseline_sha=$Baseline source_head=$SourceHead source_commit=$SourceCommit"
        foreach($Sha in @($SourceHead,$SourceCommit,$Baseline)){
            if($Sha -match '^[0-9a-f]{40}
    }
    return @($Targets)
}

function Get-RepoMarkers {
    param([string[]]$Roots)
    $Repos=New-Object Collections.Generic.HashSet[string] ([StringComparer]::OrdinalIgnoreCase)
    foreach($Root in $Roots){
        if(-not(Test-Path -LiteralPath $Root -PathType Container)){continue}
        Write-Host "SCAN_ROOT path=$Root"
        foreach($Mark in @(Get-ChildItem -LiteralPath $Root -Force -Recurse -ErrorAction SilentlyContinue | Where-Object {$_.Name -eq '.git'})){
            $Repo=if($Mark.PSIsContainer){$Mark.Parent.FullName}else{$Mark.DirectoryName}
            if(-not[string]::IsNullOrWhiteSpace($Repo)){[void]$Repos.Add($Repo)}
        }
    }
    Write-Host "SCAN_REPO_COUNT count=$($Repos.Count)"
    return @($Repos)
}

function Find-Targets {
    param([string]$Label,[string[]]$Targets,[string[]]$Roots)
    $Repos=@(Get-RepoMarkers -Roots $Roots)
    foreach($Repo in $Repos){
        $HeadOut=@(& git -C $Repo rev-parse HEAD 2>$null)
        if($LASTEXITCODE -ne 0 -or $HeadOut.Count -eq 0){continue}
        $Head=([string]$HeadOut[0]).Trim()
        $LogOut=@(& git -C $Repo log -1 --pretty='format:%ct|%H|%s' 2>$null)
        $Latest=if($LogOut.Count){[string]$LogOut[0]}else{''}
        Write-Host "REPO label=$Label path=$Repo head=$Head latest=$Latest"
        foreach($Target in $Targets){
            & git -C $Repo cat-file -e "$Target^{commit}" 2>$null
            if($LASTEXITCODE -eq 0){Write-Host "TARGET_FOUND label=$Label sha=$Target repo=$Repo"}
        }
    }
}

foreach($S in $Sidecars){
    Write-Host "=== SIDECAR $($S.Name) ==="
    $Cfg=Join-Path $S.Root 'config\autonomy.yaml'
    $Exe=Join-Path $S.Root 'bin\research-sidecar.exe'
    $Task=Get-ScheduledTask -TaskName $S.Task -ErrorAction SilentlyContinue
    Write-Host "TASK state=$(if($null-eq$Task){'MISSING'}else{[string]$Task.State})"
    if((Test-Path -LiteralPath $Exe -PathType Leaf) -and (Test-Path -LiteralPath $Cfg -PathType Leaf)){
        $Raw=& $Exe -mode status -config $Cfg 2>$null
        if($LASTEXITCODE -eq 0){
            $State=(($Raw|ForEach-Object{[string]$_})-join[Environment]::NewLine)|ConvertFrom-Json
            Write-Host "STATUS status=$($State.status) cycle_count=$($State.cycle_count) blocked=$($State.blocked_cycle_count) last_cycle_id=$($State.last_cycle_id) last_cycle_status=$($State.last_cycle_status)"
        }
    }
    if(Test-Path -LiteralPath $S.Source -PathType Container){
        $Head=((& git -C $S.Source rev-parse HEAD 2>$null)|Select-Object -First 1)
        Write-Host "SOURCE label=$($S.Name) path=$($S.Source) head=$Head"
    }
    $Targets=@(Read-LatestRows -Name $S.Name -Root $S.Root)
    Find-Targets -Label $S.Name -Targets $Targets -Roots $S.SearchRoots
}

Write-Host 'REPO_LOCAL_CUTOVER_PROBE=PASS'
 -and -not$Targets.Contains($Sha)){$Targets.Add($Sha)}
        }
    }
    return @($Targets)
}

function Get-RepoMarkers {
    param([string[]]$Roots)
    $Repos=New-Object Collections.Generic.HashSet[string] ([StringComparer]::OrdinalIgnoreCase)
    foreach($Root in $Roots){
        if(-not(Test-Path -LiteralPath $Root -PathType Container)){continue}
        Write-Host "SCAN_ROOT path=$Root"
        foreach($Mark in @(Get-ChildItem -LiteralPath $Root -Force -Recurse -ErrorAction SilentlyContinue | Where-Object {$_.Name -eq '.git'})){
            $Repo=if($Mark.PSIsContainer){$Mark.Parent.FullName}else{$Mark.DirectoryName}
            if(-not[string]::IsNullOrWhiteSpace($Repo)){[void]$Repos.Add($Repo)}
        }
    }
    Write-Host "SCAN_REPO_COUNT count=$($Repos.Count)"
    return @($Repos)
}

function Find-Targets {
    param([string]$Label,[string[]]$Targets,[string[]]$Roots)
    $Repos=@(Get-RepoMarkers -Roots $Roots)
    foreach($Repo in $Repos){
        $HeadOut=@(& git -C $Repo rev-parse HEAD 2>$null)
        if($LASTEXITCODE -ne 0 -or $HeadOut.Count -eq 0){continue}
        $Head=([string]$HeadOut[0]).Trim()
        $LogOut=@(& git -C $Repo log -1 --pretty='format:%ct|%H|%s' 2>$null)
        $Latest=if($LogOut.Count){[string]$LogOut[0]}else{''}
        Write-Host "REPO label=$Label path=$Repo head=$Head latest=$Latest"
        foreach($Target in $Targets){
            & git -C $Repo cat-file -e "$Target^{commit}" 2>$null
            if($LASTEXITCODE -eq 0){Write-Host "TARGET_FOUND label=$Label sha=$Target repo=$Repo"}
        }
    }
}

foreach($S in $Sidecars){
    Write-Host "=== SIDECAR $($S.Name) ==="
    $Cfg=Join-Path $S.Root 'config\autonomy.yaml'
    $Exe=Join-Path $S.Root 'bin\research-sidecar.exe'
    $Task=Get-ScheduledTask -TaskName $S.Task -ErrorAction SilentlyContinue
    Write-Host "TASK state=$(if($null-eq$Task){'MISSING'}else{[string]$Task.State})"
    if((Test-Path -LiteralPath $Exe -PathType Leaf) -and (Test-Path -LiteralPath $Cfg -PathType Leaf)){
        $Raw=& $Exe -mode status -config $Cfg 2>$null
        if($LASTEXITCODE -eq 0){
            $State=(($Raw|ForEach-Object{[string]$_})-join[Environment]::NewLine)|ConvertFrom-Json
            Write-Host "STATUS status=$($State.status) cycle_count=$($State.cycle_count) blocked=$($State.blocked_cycle_count) last_cycle_id=$($State.last_cycle_id) last_cycle_status=$($State.last_cycle_status)"
        }
    }
    if(Test-Path -LiteralPath $S.Source -PathType Container){
        $Head=((& git -C $S.Source rev-parse HEAD 2>$null)|Select-Object -First 1)
        Write-Host "SOURCE label=$($S.Name) path=$($S.Source) head=$Head"
    }
    $Targets=@(Read-LatestRows -Name $S.Name -Root $S.Root)
    Find-Targets -Label $S.Name -Targets $Targets -Roots $S.SearchRoots
}

Write-Host 'REPO_LOCAL_CUTOVER_PROBE=PASS'
