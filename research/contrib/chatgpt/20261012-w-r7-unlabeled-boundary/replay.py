"""W-R7 deterministic unlabeled boundary ablation on frozen synthetic formula."""
from collections import defaultdict
import json

def context(e,t):
    return (192+(7*e+3*t)%16,208+(5*e+5*t)%16,224+(11*e+7*t)%16,240+(13*e+11*t)%16)

def execute():
    outcome={}
    for arm in ("A0","A1","A2","A3"):
        groups=defaultdict(list)
        for e in range(9):
            for f in range(3):
                rows=tuple(context(e,t) for t in range(16))
                support=tuple(sorted(rows))
                cue={"A0":None,"A1":rows[0],"A2":min(rows),"A3":max(rows)}[arm]
                groups[(support,cue)].append((e,f))
        outcome[arm]={"classes":len(groups),
                      "complete_triplets":sum(len(g)==3 and len({e for e,f in g})==1 for g in groups.values()),
                      "collision_classes":sum(len({e for e,f in g})>1 for g in groups.values())}
    valid=outcome["A0"]=={"classes":8,"complete_triplets":7,"collision_classes":1} and outcome["A1"]=={"classes":9,"complete_triplets":9,"collision_classes":0}
    classification="INVALID" if not valid else "SUPPORTED" if all(outcome[k]["complete_triplets"]<9 for k in ("A2","A3")) else "NEGATIVE"
    return {"schema":"wingless.w-r7-synthetic-replay.v1","arms":outcome,"classification":classification,"records":432,"model_calls":0,"interpretation":"Existing synthetic formula replay, not independent unseen language evidence"}

if __name__=="__main__":
    x=execute()
    assert x==execute()
    print(json.dumps(x,sort_keys=True,indent=2))
