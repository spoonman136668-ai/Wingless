Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'
$Expected='efd6172772d76f03ac04bdab0af40941232311def8202d6f0f953d33b3d05dd8'
$Needle='wingless-lm-north-star'
$Db='C:\ProgramData\CKBR\research-sidecar\state\research-sidecar.db'
$Py='C:\ProgramData\CKBR\research-sidecar-yggdrasil\python312\python.exe'
foreach($P in @($Db,$Py)){if(-not(Test-Path -LiteralPath $P -PathType Leaf)){throw "MISSING path=$P"}}
$Tmp=Join-Path $env:RUNNER_TEMP 'recover-wingless-north-star.py'
$Code=@'
import ast,json,sqlite3,sys
db_path,hash_needle,name_needle=sys.argv[1:4]
db=sqlite3.connect("file:"+db_path.replace("\\","/")+"?mode=ro",uri=True,timeout=5)
hits=[]
try:
  tables=[r[0] for r in db.execute("select name from sqlite_master where type='table' order by name")]
  for table in tables:
    cols=[r[1] for r in db.execute("pragma table_info("+table+")")]
    if not cols: continue
    try:
      rows=db.execute("select rowid,* from "+table+" order by rowid desc limit 200").fetchall()
    except Exception:
      continue
    for row in rows:
      rid=row[0]
      for i,v in enumerate(row[1:]):
        if isinstance(v,(bytes,bytearray)):
          try:s=v.decode("utf-8")
          except Exception:continue
        elif isinstance(v,str):
          s=v
        else:
          continue
        if hash_needle in s or name_needle.lower() in s.lower():
          hits.append((table,rid,cols[i],s))
  print("MATCH_COUNT",len(hits))
  for table,rid,col,s in hits[:20]:
    print("MATCH",table,rid,col)
    print("RAW_BEGIN")
    print(s[:100000])
    print("RAW_END")
finally:
  db.close()
'@
[IO.File]::WriteAllText($Tmp,$Code,(New-Object Text.UTF8Encoding($false)))
& $Py $Tmp $Db $Expected $Needle
if($LASTEXITCODE-ne0){throw "DB_SEARCH_FAILED exit=$LASTEXITCODE"}
Write-Host 'WINGLESS_NORTH_STAR_DB_SEARCH=PASS'
