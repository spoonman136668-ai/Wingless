Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'
$Db='C:\ProgramData\CKBR\research-sidecar\state\research-sidecar.db'
$Py='C:\ProgramData\CKBR\research-sidecar-yggdrasil\python312\python.exe'
foreach($P in @($Db,$Py)){if(-not(Test-Path -LiteralPath $P -PathType Leaf)){throw "MISSING path=$P"}}
$Tmp=Join-Path $env:RUNNER_TEMP 'recover-wingless-contexts.py'
$Code=@'
import ast,json,sqlite3,sys
db=sqlite3.connect("file:"+sys.argv[1].replace("\\","/")+"?mode=ro",uri=True,timeout=5)
def dec(v):
    if isinstance(v,(bytes,bytearray)):
        try:return v.decode("utf-8")
        except:return repr(v[:200])
    return v
try:
    tables=[r[0] for r in db.execute("select name from sqlite_master where type='table' order by name")]
    print("TABLES",",".join(tables))
    for table in ("research_implementation_contexts","research_implementation_bundles","research_cycle_artifacts","settings"):
        if table not in tables: continue
        cols=[r[1] for r in db.execute("pragma table_info("+table+")")]
        print("COLS",table,",".join(cols))
        rows=db.execute("select rowid,* from "+table+" order by rowid desc limit 12").fetchall()
        for row in rows:
            d={cols[i]:dec(row[i+1]) for i in range(len(cols))}
            print("ROW",table,row[0])
            raw=d.get("raw")
            if isinstance(raw,str):
                try:
                    obj=json.loads(raw)
                except Exception:
                    try:
                        b=ast.literal_eval(raw); obj=json.loads(b.decode("utf-8")) if isinstance(b,(bytes,bytearray)) else {}
                    except Exception: obj={}
                if isinstance(obj,dict):
                    files=obj.get("files")
                    if isinstance(files,list):
                        print("FILES",json.dumps([{"path":x.get("path"),"sha256":x.get("sha256"),"content":x.get("content")} for x in files],sort_keys=True)[:200000])
                    for k in ("context","north_star","prompt","planning_context","model_context"):
                        if k in obj:
                            print("FIELD",k,json.dumps(obj[k],sort_keys=True)[:100000])
            for k,v in d.items():
                if isinstance(v,str) and ("wingless-lm-north-star" in v.lower() or "efd6172772d76f03ac04bdab0af40941232311def8202d6f0f953d33b3d05dd8" in v.lower()):
                    print("MATCHCOL",k,v[:200000])
finally:
    db.close()
'@
[IO.File]::WriteAllText($Tmp,$Code,(New-Object Text.UTF8Encoding($false)))
& $Py $Tmp $Db
if($LASTEXITCODE-ne0){throw "CONTEXT_DB_DUMP_FAILED exit=$LASTEXITCODE"}
Write-Host 'WINGLESS_CONTEXT_DB_DUMP=PASS'
