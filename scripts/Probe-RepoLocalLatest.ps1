Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'

function Test-GitCommit {
  param([string]$Repo,[string]$Spec)
  $Old=$ErrorActionPreference
  try{
    $ErrorActionPreference='SilentlyContinue'
    & git -C $Repo cat-file -e $Spec 2>$null
    return ($LASTEXITCODE -eq 0)
  } finally {
    $ErrorActionPreference=$Old
  }
}

$Specs=@(
  [pscustomobject]@{
    Name='wing'
    Db='C:\ProgramData\CKBR\research-sidecar\state\research-sidecar.db'
    Python='C:\ProgramData\CKBR\research-sidecar-yggdrasil\python312\python.exe'
    Source='C:\ProgramData\CKBR\codex\work\rs\source\Wingless'
    Inspect=@('.wingless/qualification-request.json','.github/workflows/research-qualify-windows.yml','research/bootstrap/wingless-lm-north-star.txt')
  },
  [pscustomobject]@{
    Name='ygg'
    Db='C:\ProgramData\CKBR\research-sidecar-yggdrasil\state\recovered-r2-20260930_190855\research-sidecar.db'
    Python='C:\ProgramData\CKBR\research-sidecar-yggdrasil\python312\python.exe'
    Source='C:\ProgramData\CKBR\research-sidecar-yggdrasil\source\Yggdrasil'
    Inspect=@('.yggdrasil/qualification-request.json','.yggdrasil/isolated-run.json','.github/workflows/plane-isolated-research.yml','research/architecture/yggdrasil-north-star.ice')
  }
)

foreach($S in $Specs){
  Write-Host "=== $($S.Name) LATEST SEALED ==="
  if(-not(Test-Path -LiteralPath $S.Db -PathType Leaf)){throw "DB_MISSING $($S.Db)"}
  if(-not(Test-Path -LiteralPath $S.Python -PathType Leaf)){throw "PYTHON_MISSING $($S.Python)"}
  if(-not(Test-Path -LiteralPath $S.Source -PathType Container)){throw "SOURCE_MISSING $($S.Source)"}

  $PyFile=Join-Path $env:RUNNER_TEMP ("latest-"+$S.Name+".py")
  $Py=@'
import json,sqlite3,sys
db=sqlite3.connect("file:"+sys.argv[1].replace("\\","/")+"?mode=ro",uri=True,timeout=5)
try:
    cols=[r[1] for r in db.execute("pragma table_info(research_sidecar_results)")]
    rows=db.execute("select * from research_sidecar_results order by rowid desc limit 3").fetchall()
    for row in rows:
        d={cols[i]:row[i] for i in range(len(cols))}
        raw=d.get("notice_raw")
        if isinstance(raw,(bytes,bytearray)):
            raw=raw.decode("utf-8")
        try:
            n=json.loads(raw)
        except Exception:
            n={}
        print(json.dumps({
            "cycle_id":d.get("cycle_id"),
            "created":d.get("created"),
            "result_class":d.get("result_class"),
            "result_sha256":d.get("result_sha256"),
            "experiment_id":n.get("experiment_id"),
            "source_head":n.get("source_head"),
            "baseline_sha":n.get("baseline_sha"),
            "north_star":n.get("north_star")
        },sort_keys=True))
finally:
    db.close()
'@
  [IO.File]::WriteAllText($PyFile,$Py,(New-Object Text.UTF8Encoding($false)))
  $Rows=@(& $S.Python $PyFile $S.Db)
  if($LASTEXITCODE-ne0 -or $Rows.Count-eq0){throw "LATEST_RESULT_QUERY_FAILED $($S.Name)"}
  $Index=0
  foreach($Line in $Rows){
    $R=$Line|ConvertFrom-Json
    Write-Host "RESULT[$Index] cycle=$($R.cycle_id) experiment=$($R.experiment_id) class=$($R.result_class) source_head=$($R.source_head) baseline=$($R.baseline_sha) created=$($R.created)"
    if($Index-eq0){
      $Sha=[string]$R.source_head
      if($Sha-notmatch'^[0-9a-f]{40}$'){throw "LATEST_SOURCE_SHA_INVALID $($S.Name)"}
      $CommitSpec=$Sha+'^{commit}'
      $RepoForCommit=$null
      if(Test-GitCommit -Repo $S.Source -Spec $CommitSpec){
        $RepoForCommit=$S.Source
      } else {
        Write-Host "LATEST_SOURCE_CANONICAL_MISSING name=$($S.Name) sha=$Sha"
        $Roots=if($S.Name -eq 'wing'){
          @('C:\ProgramData\CKBR\codex\work','C:\ProgramData\CKBR\research-sidecar')
        }else{
          @('C:\ProgramData\CKBR\research-sidecar-yggdrasil')
        }
        $Seen=New-Object Collections.Generic.HashSet[string] ([StringComparer]::OrdinalIgnoreCase)
        foreach($Root in $Roots){
          if(-not(Test-Path -LiteralPath $Root -PathType Container)){continue}
          foreach($Mark in @(Get-ChildItem -LiteralPath $Root -Force -Recurse -ErrorAction SilentlyContinue | Where-Object {$_.Name -eq '.git'})){
            $Repo=if($Mark.PSIsContainer){$Mark.Parent.FullName}else{$Mark.DirectoryName}
            if([string]::IsNullOrWhiteSpace($Repo) -or $Seen.Contains($Repo)){continue}
            [void]$Seen.Add($Repo)
            if(Test-GitCommit -Repo $Repo -Spec $CommitSpec){
              Write-Host "LATEST_SOURCE_FOUND name=$($S.Name) sha=$Sha repo=$Repo"
              if($null-eq$RepoForCommit){$RepoForCommit=$Repo}
            }
          }
        }
      }
      if($null-eq$RepoForCommit){
        Write-Host "LATEST_SOURCE_COMMIT_UNLOCATED name=$($S.Name) sha=$Sha"
      } else {
        $TreeSpec=$Sha+'^{tree}'
        $Tree=(& git -C $RepoForCommit rev-parse $TreeSpec).Trim()
        Write-Host "LATEST_SOURCE name=$($S.Name) sha=$Sha tree=$Tree repo=$RepoForCommit"
        & git -C $RepoForCommit show -s --pretty='format:LATEST_LOG %H %ct %s' $Sha
        Write-Host 'LATEST_FILES_BEGIN'
        & git -C $RepoForCommit diff-tree --no-commit-id --name-only -r $Sha
        Write-Host 'LATEST_FILES_END'
        foreach($Path in $S.Inspect){
          $BlobSpec=$Sha+':'+$Path
          if(-not(Test-GitCommit -Repo $RepoForCommit -Spec $BlobSpec)){continue}
          Write-Host "INSPECT_BEGIN $Path"
          & git -C $RepoForCommit show $BlobSpec
          Write-Host "INSPECT_END $Path"
        }
      }
    }
    $Index++
  }
}

Write-Host '=== CODEX RUNTIME ==='
$CodexHome='C:\ProgramData\CKBR\codex\home'
Write-Host "CODEX_HOME_PATH=$CodexHome exists=$(Test-Path -LiteralPath $CodexHome -PathType Container)"
if(Test-Path -LiteralPath $CodexHome -PathType Container){
  foreach($Entry in @(Get-ChildItem -LiteralPath $CodexHome -Force -ErrorAction SilentlyContinue)){
    Write-Host "CODEX_HOME_ENTRY=$($Entry.Name)"
  }
}
foreach($P in @('C:\ProgramData\CKBR\codex\bin\codex.exe','C:\ProgramData\CKBR\research-sidecar-yggdrasil\codex\bin\codex.exe')){
  if(Test-Path -LiteralPath $P -PathType Leaf){Write-Host "CODEX_BIN=$P"}
}
Write-Host 'REPO_LOCAL_LATEST_PROBE=PASS'
