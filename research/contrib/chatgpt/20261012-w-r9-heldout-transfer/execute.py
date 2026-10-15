"""Frozen W-R9 new-seed synthetic transfer, no model or tuning."""
import json
from collections import defaultdict
def evaluate(seed,shift):
    groups=defaultdict(list)
    for e in range(9):
        q=(5*e+seed)%9
        rows=[((192+seed+7*q+3*t)%256,(208+3*seed+5*q+5*t)%256,(224+5*seed+11*q+7*t)%256,(240+7*seed+13*q+11*t)%256) for t in range(16)]
        for f in range(3):
            groups[(tuple(sorted(rows)),rows[(5*f)%16] if shift else rows[0])].append((e,f))
    return sum(len(v)==3 and len({e for e,f in v})==1 for v in groups.values())
def run():
    rows=[{"seed":s,"cohort":"development" if s<4 else "heldout","aligned":evaluate(s,False),"shifted":evaluate(s,True)} for s in range(8)]
    return {"schema":"wingless.w-r9-heldout-transfer.v1","rows":rows,"heldout_aligned_all_nine":all(x["aligned"]==9 for x in rows[4:]),"heldout_shifted_below_nine":all(x["shifted"]<9 for x in rows[4:]),"model_calls":0,"classification":"SUPPORTED" if all(x["aligned"]==9 and x["shifted"]<9 for x in rows[4:]) else "NEGATIVE","claim_limit":"new deterministic synthetic seeds only; no real-language generalization"}
if __name__=="__main__":
    o=run();assert o==run()
    print(json.dumps(o,sort_keys=True,indent=2))
