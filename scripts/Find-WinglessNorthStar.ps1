Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'
$Db='C:\ProgramData\CKBR\research-sidecar\state\research-sidecar.db'
$Py='C:\ProgramData\CKBR\research-sidecar-yggdrasil\python312\python.exe'
$Expected='efd6172772d76f03ac04bdab0af40941232311def8202d6f0f953d33b3d05dd8'
$Tmp=Join-Path $env:RUNNER_TEMP 'recover-wingless-material.py'
$Code=@'
import json,sqlite3,sys
db=sqlite3.connect("file:"+sys.argv[1].replace("\\","/")+"?mode=ro",uri=True,timeout=5)
want=sys.argv[2].lower()
try:
    cols=[r[1] for r in db.execute("pragma table_info(research_materials)")]
    print("COLS",",".join(cols))
    rows=db.execute("select rowid,* from research_materials order by rowid desc").fetchall()
    hits=0
    for row in rows:
        d={cols[i]:row[i+1] for i in range(len(cols))}
        parts=[]
        for k,v in d.items():
            if isinstance(v,(bytes,bytearray)):
                try:s=v.decode("utf-8")
                except:s=v.hex()
            else:s=str(v)
            parts.append(s)
        joined="\n".join(parts).lower()
        if want in joined or "wingless-lm-north-star" in joined:
            hits+=1
            print("MATCH rowid",row[0])
            for k,v in d.items():
                if isinstance(v,(bytes,bytearray)):
                    try:v=v.decode("utf-8")
                    except:v=v.hex()
                print("FIELD",k)
                print(str(v)[:300000])
                print("END_FIELD")
    print("MATCH_COUNT",hits)
finally:
    db.close()
'@
[IO.File]::WriteAllText($Tmp,$Code,(New-Object Text.UTF8Encoding($false)))
& $Py $Tmp $Db $Expected
if($LASTEXITCODE-ne0){throw "MATERIAL_QUERY_FAILED exit=$LASTEXITCODE"}
Write-Host 'WINGLESS_MATERIAL_QUERY=PASS'
