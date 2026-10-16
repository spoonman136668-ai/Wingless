"""Diagnostic of ONE preserved W-F60 FIT result, NEVER a source replay."""
import json,sys
from audit_fit import recompute,ARMS
d=json.load(open(sys.argv[1]))
drift={a:0 for a in ARMS[:-1]}
for row in d["rows"]:
 rec,num=recompute(row)
 bad=[]
 for arm,val in zip(ARMS[:-1],rec):
  if tuple(row["predictions"][arm])!=val:
   drift[arm]+=1;bad.append((arm,row["predictions"][arm],val))
 if bad and sum(drift.values())<=20:
  print("INDEPENDENT_PREDICTION_DRIFT",row["seed"],row["family"],row["length"],bad)
print("TOTAL_DIAGNOSTIC_DRIFT",json.dumps(drift,sort_keys=True))
print("ORIGINAL_SOURCE_NOT_EXECUTED=true")
