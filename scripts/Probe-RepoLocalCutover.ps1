Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function Get-SidecarSnapshot {
    param(
        [Parameter(Mandatory=$true)][string]$Name,
        [Parameter(Mandatory=$true)][string]$Root,
        [Parameter(Mandatory=$true)][string]$Source,
        [Parameter(Mandatory=$true)][string]$TaskName
    )

    Write-Host "=== SIDECAR $Name ==="
    $Cfg=Join-Path $Root 'config\autonomy.yaml'
    $Exe=Join-Path $Root 'bin\research-sidecar.exe'
    $Py=Join-Path $Root 'python312\python.exe'
    if(-not(Test-Path -LiteralPath $Py -PathType Leaf)){
        $Py=(Get-Command python -ErrorAction Stop).Source
    }

    $Task=Get-ScheduledTask -TaskName $TaskName -ErrorAction SilentlyContinue
    Write-Host "TASK state=$(if($null -eq $Task){'MISSING'}else{[string]$Task.State})"

    if((Test-Path -LiteralPath $Exe -PathType Leaf) -and (Test-Path -LiteralPath $Cfg -PathType Leaf)){
        $Raw=& $Exe -mode status -config $Cfg 2>$null
        if($LASTEXITCODE -eq 0){
            $State=(($Raw|ForEach-Object{[string]$_}) -join [Environment]::NewLine)|ConvertFrom-Json
            Write-Host "STATUS status=$($State.status) cycle_count=$($State.cycle_count) blocked=$($State.blocked_cycle_count) last_cycle_id=$($State.last_cycle_id) last_cycle_status=$($State.last_cycle_status)"
        }
    }

    if(Test-Path -LiteralPath $Source -PathType Container){
        $Head=(& git -C $Source rev-parse HEAD).Trim()
        $Tree=(& git -C $Source rev-parse 'HEAD^{tree}').Trim()
        $BranchRaw=(& git -C $Source branch --show-current 2>$null | Select-Object -First 1)
        $Branch=if($null -eq $BranchRaw){''}else{([string]$BranchRaw).Trim()}
        Write-Host "SOURCE path=$Source head=$Head tree=$Tree branch=$Branch"
    } else {
        Write-Host "SOURCE missing=$Source"
    }

    $CfgRaw=[IO.File]::ReadAllText($Cfg)
    $M=[regex]::Match($CfgRaw,"(?m)^\s*sqlite_path:\s*['""]?([^\r\n'""]+)['""]?\s*$")
    if(-not$M.Success){throw "$Name sqlite_path missing"}
    $Db=$M.Groups[1].Value.Trim()
    if(-not[IO.Path]::IsPathRooted($Db)){$Db=Join-Path $Root ($Db -replace '/','\')}
    $Db=[IO.Path]::GetFullPath($Db)
    Write-Host "DB path=$Db"

    $Probe=Join-Path $env:RUNNER_TEMP ("cutover-db-"+$Name+".py")
    $Python=@'
import json, sqlite3, sys
p=sys.argv[1]
db=sqlite3.connect("file:"+p.replace("\\","/")+"?mode=ro",uri=True,timeout=5)
def dec(v):
    return v.decode("utf-8") if isinstance(v,(bytes,bytearray)) else v
def latest(table):
    cols=[r[1] for r in db.execute("pragma table_info("+table+")")]
    if not cols:
        return None
    row=db.execute("select * from "+table+" order by rowid desc limit 1").fetchone()
    if row is None:
        return None
    return {cols[i]:dec(row[i]) for i in range(len(cols))}
try:
    out={}
    for t in ("research_sidecar_results","research_implementation_freezes","research_cycles"):
        out[t]=latest(t)
    print(json.dumps(out,sort_keys=True))
finally:
    db.close()
'@
    [IO.File]::WriteAllText($Probe,$Python,(New-Object Text.UTF8Encoding($false)))
    $Obj=((& $Py $Probe $Db) -join [Environment]::NewLine)|ConvertFrom-Json

    $Targets=New-Object Collections.Generic.List[string]
    foreach($Name2 in @('research_sidecar_results','research_implementation_freezes','research_cycles')){
        $Row=$Obj.$Name2
        if($null -eq $Row){continue}
        $RawProp=$Row.PSObject.Properties['raw']
        if($null -eq $RawProp){continue}
        $RawText=[string]$RawProp.Value
        if([string]::IsNullOrWhiteSpace($RawText)){continue}
        $R=$RawText|ConvertFrom-Json
        Write-Host "$Name2 cycle_id=$($R.cycle_id) status=$($R.status) experiment_id=$($R.experiment_id) baseline_sha=$($R.baseline_sha) source_head=$($R.source_head) source_commit=$($R.source_commit)"
        foreach($Sha in @([string]$R.source_head,[string]$R.source_commit,[string]$R.baseline_sha)){
            if($Sha -match '^[0-9a-f]{40}$' -and -not$Targets.Contains($Sha)){$Targets.Add($Sha)}
        }
    }
    return @($Targets)
}

$WingTargets=@(Get-SidecarSnapshot -Name wing -Root 'C:\ProgramData\CKBR\research-sidecar' -Source 'C:\ProgramData\CKBR\codex\work\rs\source\Wingless' -TaskName 'CKBPlane Research Sidecar')
$YggTargets=@(Get-SidecarSnapshot -Name ygg -Root 'C:\ProgramData\CKBR\research-sidecar-yggdrasil' -Source 'C:\ProgramData\CKBR\research-sidecar-yggdrasil\source\Yggdrasil' -TaskName 'CKBPlane Research Sidecar Yggdrasil')

Write-Host '=== CODEX ==='
foreach($P in @(
    'C:\ProgramData\CKBR\codex\bin\codex.exe',
    'C:\ProgramData\CKBR\research-sidecar-yggdrasil\codex\bin\codex.exe'
)){
    if(Test-Path -LiteralPath $P -PathType Leaf){
        Write-Host "CODEX path=$P version=$((& $P --version 2>$null)-join' ')"
    }
}

function Find-CommitRepos {
    param([string]$Label,[string[]]$Targets,[string[]]$Roots)
    Write-Host "=== LOCATE $Label COMMITS ==="
    $Repos=New-Object Collections.Generic.HashSet[string] ([StringComparer]::OrdinalIgnoreCase)
    foreach($Root in $Roots){
        if(-not(Test-Path -LiteralPath $Root -PathType Container)){continue}
        foreach($Mark in @(Get-ChildItem -LiteralPath $Root -Force -Recurse -ErrorAction SilentlyContinue | Where-Object {$_.Name -eq '.git'})){
            $Repo=if($Mark.PSIsContainer){$Mark.Parent.FullName}else{$Mark.DirectoryName}
            [void]$Repos.Add($Repo)
        }
    }
    foreach($Target in $Targets){
        $Found=$false
        foreach($Repo in $Repos){
            & git -C $Repo cat-file -e "$Target^{commit}" 2>$null
            if($LASTEXITCODE -eq 0){
                $HeadRaw=(& git -C $Repo rev-parse HEAD 2>$null|Select-Object -First 1)
                $Head=if($null-eq$HeadRaw){''}else{([string]$HeadRaw).Trim()}
                $BranchRaw=(& git -C $Repo branch --show-current 2>$null|Select-Object -First 1)
                $Branch=if($null-eq$BranchRaw){''}else{([string]$BranchRaw).Trim()}
                $Dirty=((& git -C $Repo status --porcelain 2>$null)-join';')
                Write-Host "COMMIT_FOUND label=$Label sha=$Target repo=$Repo head=$Head branch=$Branch dirty=$([bool](-not[string]::IsNullOrWhiteSpace($Dirty)))"
                $Found=$true
            }
        }
        if(-not$Found){Write-Host "COMMIT_NOT_FOUND label=$Label sha=$Target"}
    }
}

Find-CommitRepos -Label wing -Targets $WingTargets -Roots @(
    'C:\ProgramData\CKBR\codex\work\rs',
    'C:\ProgramData\CKBR\research-sidecar'
)
Find-CommitRepos -Label ygg -Targets $YggTargets -Roots @(
    'C:\ProgramData\CKBR\research-sidecar-yggdrasil'
)

Write-Host 'REPO_LOCAL_CUTOVER_PROBE=PASS'
