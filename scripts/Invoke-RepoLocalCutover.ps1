param(
    [Parameter(Mandatory=$true)][ValidateSet('Wingless','Yggdrasil')][string]$Program,
    [Parameter(Mandatory=$true)][string]$RepoPath,
    [Parameter(Mandatory=$true)][string]$ActiveBranch,
    [Parameter(Mandatory=$true)][string]$TaskName,
    [Parameter(Mandatory=$true)][string]$DbPath,
    [Parameter(Mandatory=$true)][string[]]$SearchRoots
)

Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'
$HandoffPushed=$false
$TaskDisabled=$false

function Invoke-Git {
    param([Parameter(ValueFromRemainingArguments=$true)][string[]]$Args)
    & git -C $RepoPath @Args
    if($LASTEXITCODE-ne0){throw "GIT_FAILED exit=$LASTEXITCODE args=$($Args -join ' ')"}
}

function Resolve-Python {
    foreach($P in @(
        'C:\ProgramData\CKBR\research-sidecar-yggdrasil\python312\python.exe',
        'C:\Python312\python.exe'
    )){
        if(Test-Path -LiteralPath $P -PathType Leaf){return $P}
    }
    throw 'CUTOVER_PYTHON_MISSING'
}

if(-not(Test-Path -LiteralPath $RepoPath -PathType Container)){throw "REPO_PATH_MISSING path=$RepoPath"}
if(-not(Test-Path -LiteralPath $DbPath -PathType Leaf)){throw "SIDECAR_DB_MISSING path=$DbPath"}

try{
    $Task=Get-ScheduledTask -TaskName $TaskName -ErrorAction Stop
    Write-Host "CUTOVER_TASK_INITIAL name=$TaskName state=$($Task.State)"
    Disable-ScheduledTask -TaskName $TaskName -ErrorAction Stop|Out-Null
    $TaskDisabled=$true
    try{Stop-ScheduledTask -TaskName $TaskName -ErrorAction Stop}catch{
        $AfterDisable=Get-ScheduledTask -TaskName $TaskName -ErrorAction Stop
        if([string]$AfterDisable.State-ceq'Running'){throw}
    }
    Start-Sleep -Seconds 2
    $Stopped=Get-ScheduledTask -TaskName $TaskName -ErrorAction Stop
    if([string]$Stopped.State-ceq'Running'){throw "SIDECAR_TASK_STILL_RUNNING task=$TaskName"}
    Write-Host "CUTOVER_TASK_FROZEN name=$TaskName state=$($Stopped.State)"

    $Py=Resolve-Python
    $Temp=Join-Path $env:RUNNER_TEMP ("repo-local-cutover-"+$Program.ToLowerInvariant())
    if(Test-Path -LiteralPath $Temp){Remove-Item -LiteralPath $Temp -Recurse -Force}
    New-Item -ItemType Directory -Force -Path $Temp|Out-Null
    $SnapshotPy=Join-Path $Temp 'snapshot.py'
    $SnapshotJson=Join-Path $Temp 'snapshot.json'
    $RootsJson=($SearchRoots|ConvertTo-Json -Compress)
    $PyText=@'
import ast, datetime, hashlib, json, os, sqlite3, subprocess, sys

program, db_path, roots_json, out_path = sys.argv[1:5]
roots=json.loads(roots_json)

def decode(v):
    if isinstance(v,(bytes,bytearray)):
        return v.decode("utf-8")
    return v

def latest(db, table):
    cols=[r[1] for r in db.execute("pragma table_info("+table+")")]
    if not cols:
        return None
    row=db.execute("select * from "+table+" order by rowid desc limit 1").fetchone()
    if row is None:
        return None
    return {cols[i]:decode(row[i]) for i in range(len(cols))}

def parsed(raw):
    if not isinstance(raw,str) or not raw:
        return {}
    for candidate in (raw,):
        try:return json.loads(candidate)
        except Exception:pass
    if raw.startswith("b'") or raw.startswith('b"'):
        try:
            b=ast.literal_eval(raw)
            if isinstance(b,(bytes,bytearray)):
                return json.loads(b.decode("utf-8"))
        except Exception:
            pass
    return {"unparsed_raw":raw}

db=sqlite3.connect("file:"+db_path.replace("\\","/")+"?mode=ro",uri=True,timeout=5)
try:
    result_row=latest(db,"research_sidecar_results") or {}
    freeze_row=latest(db,"research_implementation_freezes") or {}
    cycle_row=latest(db,"research_cycles") or {}
finally:
    db.close()

notice=parsed(result_row.get("notice_raw"))
freeze=parsed(freeze_row.get("raw"))
cycle=parsed(cycle_row.get("raw"))
source_head=str(notice.get("source_head") or freeze.get("source_commit") or cycle.get("baseline_sha") or "")
if len(source_head)!=40 or any(c not in "0123456789abcdefABCDEF" for c in source_head):
    raise SystemExit("SOURCE_HEAD_INVALID "+source_head)

def git_has(repo,sha):
    p=subprocess.run(["git","-C",repo,"cat-file","-e",sha+"^{commit}"],stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL)
    return p.returncode==0

repos=[]
for root in roots:
    if not os.path.isdir(root):
        continue
    if os.path.exists(os.path.join(root,".git")):
        repos.append(root)
    for current,dirs,files in os.walk(root):
        if ".git" in dirs or ".git" in files:
            if current not in repos: repos.append(current)
            if ".git" in dirs: dirs.remove(".git")

source_repo=None
for repo in repos:
    if git_has(repo,source_head):
        source_repo=repo
        break
if source_repo is None:
    raise SystemExit("SOURCE_COMMIT_NOT_FOUND "+source_head)

scientific=None
wanted=str(result_row.get("result_sha256") or "").lower()
if wanted:
    for root in roots:
        if not os.path.isdir(root): continue
        for current,dirs,files in os.walk(root):
            if "result.json" not in files: continue
            p=os.path.join(current,"result.json")
            try:data=open(p,"rb").read()
            except OSError:continue
            if hashlib.sha256(data).hexdigest()==wanted:
                try:scientific=json.loads(data.decode("utf-8"))
                except Exception:scientific={"raw_utf8":data.decode("utf-8","replace")}
                break
        if scientific is not None: break

snapshot={
    "schema":"research.repo-local-handoff.v1",
    "program":program,
    "captured_at_utc":datetime.datetime.now(datetime.timezone.utc).isoformat(),
    "source_head":source_head.lower(),
    "source_repo":source_repo,
    "experiment_id":notice.get("experiment_id"),
    "result_class":result_row.get("result_class"),
    "result_sha256":result_row.get("result_sha256"),
    "observation_sha256":result_row.get("observation_sha256"),
    "north_star":notice.get("north_star"),
    "completion_notice":notice,
    "implementation_freeze":freeze,
    "pending_cycle":cycle,
    "scientific_result_available":scientific is not None,
    "scientific_result":scientific,
}
with open(out_path,"w",encoding="utf-8") as f:
    json.dump(snapshot,f,indent=2,sort_keys=True)
print(json.dumps({"source_head":snapshot["source_head"],"source_repo":source_repo,"experiment_id":snapshot["experiment_id"],"result_class":snapshot["result_class"],"scientific_result_available":snapshot["scientific_result_available"]},sort_keys=True))
'@
    [IO.File]::WriteAllText($SnapshotPy,$PyText,(New-Object Text.UTF8Encoding($false)))
    $SnapshotLine=& $Py $SnapshotPy $Program $DbPath $RootsJson $SnapshotJson
    if($LASTEXITCODE-ne0){throw "SIDECAR_SNAPSHOT_FAILED exit=$LASTEXITCODE"}
    $SnapshotInfo=$SnapshotLine|ConvertFrom-Json
    $SourceHead=[string]$SnapshotInfo.source_head
    $SourceRepo=[string]$SnapshotInfo.source_repo
    if($SourceHead-notmatch'^[0-9a-f]{40}$'){throw "SNAPSHOT_SOURCE_HEAD_INVALID sha=$SourceHead"}
    if(-not(Test-Path -LiteralPath $SourceRepo -PathType Container)){throw "SNAPSHOT_SOURCE_REPO_MISSING path=$SourceRepo"}
    Write-Host "CUTOVER_SNAPSHOT program=$Program experiment=$($SnapshotInfo.experiment_id) source_head=$SourceHead result_class=$($SnapshotInfo.result_class) scientific_result_available=$($SnapshotInfo.scientific_result_available)"

    Push-Location $RepoPath
    try{
        git ls-remote --exit-code --heads origin "refs/heads/$ActiveBranch" *> $null
        if($LASTEXITCODE-eq0){throw "ACTIVE_BRANCH_ALREADY_EXISTS branch=$ActiveBranch"}
        Invoke-Git fetch $SourceRepo $SourceHead
        Invoke-Git switch -C $ActiveBranch FETCH_HEAD
        $Exact=(& git rev-parse HEAD).Trim()
        if($Exact-cne$SourceHead){throw "HANDOFF_SOURCE_MISMATCH expected=$SourceHead actual=$Exact"}
        if(((& git status --porcelain)-join'').Trim()){throw 'HANDOFF_SOURCE_DIRTY'}

        $HandoffDir=Join-Path $RepoPath '.research-autonomy'
        New-Item -ItemType Directory -Force -Path $HandoffDir|Out-Null
        Copy-Item -LiteralPath $SnapshotJson -Destination (Join-Path $HandoffDir 'handoff.json') -Force

        $State=[ordered]@{
            schema='research.repo-local-state.v1'
            program=$Program
            active_branch=$ActiveBranch
            handoff_source_sha=$SourceHead
            experiment=[string]$SnapshotInfo.experiment_id
            result_class=[string]$SnapshotInfo.result_class
            last_completed_utc=[DateTime]::UtcNow.ToString('o')
            authority='research-only'
            ckb_plane_role='governance-only'
            sidecar_task=$TaskName
            sidecar_disabled=$true
        }
        $State|ConvertTo-Json -Depth 20|Set-Content -LiteralPath (Join-Path $HandoffDir 'state.json') -Encoding UTF8
        git config user.name 'repo-local-research-cutover'
        git config user.email 'repo-local-research-cutover@users.noreply.github.com'
        Invoke-Git add '.research-autonomy'
        Invoke-Git commit -m 'research: hand off sidecar state to repo-local autonomy'
        $HandoffHead=(& git rev-parse HEAD).Trim()
        Invoke-Git push origin ("HEAD:refs/heads/"+$ActiveBranch)
        $HandoffPushed=$true
        Write-Host 'REPO_LOCAL_CUTOVER_HANDOFF_PASS'
        Write-Host "program=$Program"
        Write-Host "source_head=$SourceHead"
        Write-Host "handoff_head=$HandoffHead"
        Write-Host "active_branch=$ActiveBranch"
    } finally {
        Pop-Location
    }
}
catch{
    if($TaskDisabled -and -not$HandoffPushed){
        Write-Host "CUTOVER_ROLLBACK_REENABLE task=$TaskName reason=$($_.Exception.Message)"
        try{
            Enable-ScheduledTask -TaskName $TaskName -ErrorAction Stop|Out-Null
            Start-ScheduledTask -TaskName $TaskName -ErrorAction Stop
        }catch{
            Write-Host "CUTOVER_ROLLBACK_FAILED task=$TaskName error=$($_.Exception.Message)"
        }
    }
    throw
}
