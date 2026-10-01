import json, os, sqlite3

SPECS=[
 ("wing", r"C:\ProgramData\CKBR\research-sidecar\state\research-sidecar.db"),
 ("ygg", r"C:\ProgramData\CKBR\research-sidecar-yggdrasil\state\recovered-r2-20260930_190855\research-sidecar.db"),
]
for name,path in SPECS:
    print("===DB",name,"===")
    db=sqlite3.connect("file:"+path.replace("\\","/")+"?mode=ro",uri=True,timeout=5)
    try:
        tables=[r[0] for r in db.execute("select name from sqlite_master where type='table' order by name")]
        for table in tables:
            cols=[r[1] for r in db.execute("pragma table_info("+table+")")]
            print("TABLE",table,"COLS",",".join(cols))
            if any(x in table.lower() for x in ("research","planner","draft","proposal","observation","result","cycle")):
                try:
                    rows=db.execute("select * from "+table+" order by rowid desc limit 2").fetchall()
                except Exception as exc:
                    print("ROW_ERROR",table,repr(exc)); continue
                for row in rows:
                    d={cols[i]:row[i] for i in range(len(cols))}
                    safe={}
                    for k,v in d.items():
                        if isinstance(v,(bytes,bytearray)):
                            try:v=v.decode("utf-8")
                            except Exception:v=f"<bytes:{len(v)}>"
                        if isinstance(v,str) and len(v)>12000:
                            v=v[:12000]+"<TRUNCATED>"
                        safe[k]=v
                    print("ROW",table,json.dumps(safe,sort_keys=True,default=str))
    finally:
        db.close()
print("SIDECAR_DB_PROBE=PASS")
