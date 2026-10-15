"""W-R8 original deterministic synthetic temporal alignment experiment."""
from collections import defaultdict
import json
def context(e,t):
    return (192+(7*e+3*t)%16,208+(5*e+5*t)%16,224+(11*e+7*t)%16,240+(13*e+11*t)%16)
def arm(t,shift):
    groups=defaultdict(list)
    for e in range(9):
        rows=tuple(context(e,i) for i in range(16))
        for f in range(3):
            cue=rows[(t+5*f)%16] if shift else rows[t]
            groups[(tuple(sorted(rows)),cue)].append((e,f))
    return {"classes":len(groups),"complete_triplets":sum(len(v)==3 and len({e for e,f in v})==1 for v in groups.values()),"collision_classes":sum(len({e for e,f in v})>1 for v in groups.values())}
def execute():
    aligned=[arm(t,False) for t in range(16)]
    shifted=[arm(t,True) for t in range(16)]
    valid=aligned[0]=={"classes":9,"complete_triplets":9,"collision_classes":0}
    all_aligned=all(x["complete_triplets"]==9 for x in aligned)
    all_shifted=all(x["complete_triplets"]<9 for x in shifted)
    return {"schema":"wingless.w-r8-temporal-alignment.v1","classification":"INVALID" if not valid else "SUPPORTED" if all_aligned and all_shifted else "MIXED" if all_aligned!=all_shifted else "NEGATIVE","aligned":aligned,"shifted":shifted,"ticks":16,"model_calls":0,"claim_limit":"Synthetic indexed formula; no transfer or primary Go qualification"}
if __name__=="__main__":
    x=execute()
    assert x==execute()
    print(json.dumps(x,sort_keys=True,indent=2))
