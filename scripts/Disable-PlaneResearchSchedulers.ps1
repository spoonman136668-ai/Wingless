Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'

$Tasks=@(
  [pscustomobject]@{Name='CKBPlane Research Sidecar';Label='wing';Db='C:\ProgramData\CKBR\research-sidecar\state\research-sidecar.db'},
  [pscustomobject]@{Name='CKBPlane Research Sidecar Yggdrasil';Label='ygg';Db='C:\ProgramData\CKBR\research-sidecar-yggdrasil\state\recovered-r2-20260930_190855\research-sidecar.db'}
)
$Py='C:\ProgramData\CKBR\research-sidecar-yggdrasil\python312\python.exe'
if(-not(Test-Path -LiteralPath $Py -PathType Leaf)){throw 'SHARED_PYTHON_MISSING'}

foreach($T in $Tasks){
  $Task=Get-ScheduledTask -TaskName $T.Name -ErrorAction Stop
  Write-Host "BEFORE label=$($T.Label) task=$($T.Name) state=$($Task.State)"
  try { Stop-ScheduledTask -TaskName $T.Name -ErrorAction Stop } catch { throw "TASK_STOP_FAILED label=$($T.Label) error=$($_.Exception.Message)" }
  try { Disable-ScheduledTask -TaskName $T.Name -ErrorAction Stop | Out-Null } catch { throw "TASK_DISABLE_FAILED label=$($T.Label) error=$($_.Exception.Message)" }
}

Start-Sleep -Seconds 2

foreach($T in $Tasks){
  $Task=Get-ScheduledTask -TaskName $T.Name -ErrorAction Stop
  $Info=Get-ScheduledTaskInfo -TaskName $T.Name -ErrorAction Stop
  Write-Host "AFTER label=$($T.Label) task=$($T.Name) state=$($Task.State) last_run=$($Info.LastRunTime.ToUniversalTime().ToString('o')) last_result=$($Info.LastTaskResult)"
  if([string]$Task.State -cne 'Disabled'){throw "TASK_NOT_DISABLED label=$($T.Label) state=$($Task.State)"}
}

$Procs=@(Get-CimInstance Win32_Process -ErrorAction SilentlyContinue | Where-Object {
  ([string]$_.ExecutablePath -like 'C:\ProgramData\CKBR\research-sidecar*\bin\research-sidecar.exe') -or
  ([string]$_.CommandLine -like '*C:\ProgramData\CKBR\research-sidecar*research-sidecar.exe*')
})
foreach($P in $Procs){Write-Host "SIDECAR_PROCESS_REMAINS pid=$($P.ProcessId) exe=$($P.ExecutablePath)"}
if($Procs.Count-gt0){throw "SIDECAR_PROCESS_REMAINS count=$($Procs.Count)"}

$Probe=Join-Path $env:RUNNER_TEMP 'cutover-final-results.py'
$Code=@'
import json,sqlite3,sys
for item in sys.argv[1:]:
    label,path=item.split("=",1)
    db=sqlite3.connect("file:"+path.replace("\\","/")+"?mode=ro",uri=True,timeout=5)
    try:
        cols=[r[1] for r in db.execute("pragma table_info(research_sidecar_results)")]
        row=db.execute("select * from research_sidecar_results order by rowid desc limit 1").fetchone()
        if row is None:
            print(json.dumps({"label":label,"missing":True},sort_keys=True))
            continue
        d={cols[i]:row[i] for i in range(len(cols))}
        raw=d.get("notice_raw")
        if isinstance(raw,(bytes,bytearray)): raw=raw.decode("utf-8")
        n=json.loads(raw)
        print(json.dumps({
            "label":label,
            "cycle_id":d.get("cycle_id"),
            "created":d.get("created"),
            "result_class":d.get("result_class"),
            "result_sha256":d.get("result_sha256"),
            "experiment_id":n.get("experiment_id"),
            "source_head":n.get("source_head"),
            "baseline_sha":n.get("baseline_sha")
        },sort_keys=True))
    finally:
        db.close()
'@
[IO.File]::WriteAllText($Probe,$Code,(New-Object Text.UTF8Encoding($false)))
$Args=@($Probe)
foreach($T in $Tasks){$Args+=($T.Label+'='+$T.Db)}
& $Py @Args
if($LASTEXITCODE-ne0){throw 'FINAL_RESULT_READ_FAILED'}

Write-Host 'PLANE_RESEARCH_SCHEDULERS_DISABLED=PASS'
