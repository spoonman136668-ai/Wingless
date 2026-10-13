"""W-R32: frozen candidate operation tables, hidden within-episode regime change."""
import hashlib
import importlib.util
import json
from pathlib import Path

SOURCE = Path(__file__).resolve().parent.parent / "20261013-w-r26-latent-alphabet" / "execute.py"
_spec = importlib.util.spec_from_file_location("w26_frozen", SOURCE)
assert _spec and _spec.loader
w26 = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(w26)
assert len(w26.PERMUTATIONS)==24 and w26.fixture()=="PASS"
ARMS=("FROZEN_ADD","ERROR_TRIGGERED_XOR","ABSTAIN_AFTER_ERROR")


def order(seed):
    assert 192 <= seed <= 199
    return tuple(sorted(range(16),key=lambda e:(
        hashlib.sha256(("W32-order/"+str(seed)+"/"+str(e)).encode()).digest(),e)))


def sequence(seed,episode):
    assert 192 <= seed <= 199 and 0<=episode<16
    e=order(seed)[episode]
    xs=[(e+seed)%4,((e//4)+3*seed)%4]
    while len(xs)<32:
        t=len(xs)
        a,b=xs[-2],xs[-1]
        xs.append((a^b) if seed%2==0 and t>=16 else (a+b)%4)
    p=w26.PERMUTATIONS[hashlib.sha256(
        ("W32/"+str(seed)).encode()).digest()[0]%24]
    return tuple(p[x] for x in xs)


def xor_table(p):
    inverse=[0]*4
    for latent,observed in enumerate(p):
        inverse[observed]=latent
    return tuple(p[inverse[a]^inverse[b]] for a,b in w26.ALL_CONTEXTS)


def calibration(seed):
    observed=[sequence(seed,e) for e in range(16)]
    cases=[(xs[0],xs[1],xs[2]) for xs in observed]
    if len(set((a,b) for a,b,_ in cases))!=16:
        return {"invalid":"MISSING_CALIBRATION_CONTEXT"}
    hypotheses=[p for p in w26.PERMUTATIONS if w26.fits(p,cases)]
    add_tables={w26.table_for(p) for p in hypotheses}
    xor_tables={xor_table(p) for p in hypotheses}
    if len(add_tables)!=1 or len(xor_tables)!=1:
        return {"invalid":"AMBIGUOUS_HYPOTHESIS","candidates":len(hypotheses),
                "add_tables":len(add_tables),"xor_tables":len(xor_tables)}
    add=next(iter(add_tables))
    xor=next(iter(xor_tables))
    return {"add":add,"xor":xor,"candidates":len(hypotheses),
            "add_sha256":hashlib.sha256(bytes(add)).hexdigest(),
            "xor_sha256":hashlib.sha256(bytes(xor)).hexdigest(),
            "labels":16}


def evaluate(seed,arm,add,xor):
    confusion=[[0]*4 for _ in range(4)]
    trace=hashlib.sha256()
    pre={"correct":0,"wrong":0,"abstained":0}
    post={"correct":0,"wrong":0,"abstained":0}
    triggers=[]
    for ep in range(16):
        observations=sequence(seed,ep)
        switched=False
        first=None
        for t in range(3,32):
            phase=pre if t<16 else post
            if switched and arm=="ABSTAIN_AFTER_ERROR":
                prediction=None
            else:
                model=xor if switched and arm=="ERROR_TRIGGERED_XOR" else add
                prediction=model[4*observations[t-2]+observations[t-1]]
            truth=observations[t]  # strictly after forecast; used as feedback next t
            if prediction is None:
                phase["abstained"]+=1
                trace.update(bytes((ep,t,truth,4)))
            else:
                confusion[truth][prediction]+=1
                if prediction==truth:
                    phase["correct"]+=1
                else:
                    phase["wrong"]+=1
                trace.update(bytes((ep,t,truth,prediction)))
                if prediction!=truth and not switched and arm!="FROZEN_ADD":
                    switched=True
                    first=t
        triggers.append(first)
    assert sum(sum(row) for row in confusion)+pre["abstained"]+post["abstained"]==16*29
    return {"seed":seed,"cohort":"development" if seed<196 else "heldout",
            "condition":"SHIFT" if seed%2==0 else "STABLE",
            "arm":arm,"before_t16":pre,"from_t16":post,
            "first_errors":triggers,"trigger_count":sum(x is not None for x in triggers),
            "matrix":confusion,"trace_sha256":trace.hexdigest()}


def run():
    fingerprints={str(s):hashlib.sha256(bytes(
        x for ep in range(16) for x in sequence(s,ep))).hexdigest()
        for s in range(192,200)}
    if len(set(fingerprints.values()))!=8:
        return {"schema":"wingless.w-r32-regime-shift.v1",
                "classification":"INVALID_IDENTITIES","fingerprints":fingerprints,
                "rows":[],"model_calls":0}
    rows=[]
    calib={}
    for seed in range(192,200):
        fit=calibration(seed)
        if "invalid" in fit:
            return {"schema":"wingless.w-r32-regime-shift.v1",
                    "classification":"AMBIGUOUS_HYPOTHESIS",
                    "fingerprints":fingerprints,"failing_seed":seed,
                    "failure":fit,"rows":[],"model_calls":0}
        add,xor=fit.pop("add"),fit.pop("xor")
        calib[str(seed)]=fit
        for arm in ARMS:
            rows.append(evaluate(seed,arm,add,xor))
    held=[r for r in rows if r["seed"]>=196]
    score={}
    for condition in ("STABLE","SHIFT"):
        score[condition]={}
        for arm in ARMS:
            subset=[r for r in held if r["condition"]==condition and r["arm"]==arm]
            score[condition][arm]={
                "pre":{k:sum(r["before_t16"][k] for r in subset)
                       for k in ("correct","wrong","abstained")},
                "post":{k:sum(r["from_t16"][k] for r in subset)
                        for k in ("correct","wrong","abstained")},
                "triggers":sum(r["trigger_count"] for r in subset)}
    stable=score["STABLE"]["ERROR_TRIGGERED_XOR"]
    shift=score["SHIFT"]
    prior_same=all(row["before_t16"]==ref["before_t16"]
                   for row in rows for ref in rows
                   if row["seed"]==ref["seed"])
    gain=shift["ERROR_TRIGGERED_XOR"]["post"]["correct"]>shift["FROZEN_ADD"]["post"]["correct"]
    if (prior_same and gain and stable["post"]["wrong"]==0
        and stable["pre"]["wrong"]==0 and stable["triggers"]==0):
        classification="CAUSAL_SWITCH_RECOVERY"
    else:
        classification="MIXED"
    return {"schema":"wingless.w-r32-regime-shift.v1",
            "classification":classification,
            "parent":"d44de8f85364b6a4b336ffddb1fbc9e1c125089e",
            "fingerprints":fingerprints,"calibration":calib,
            "heldout":score,"rows":rows,"model_calls":0,
            "claim_limit":"Synthetic known ADD and XOR table candidates with observed causal feedback; not unrestricted change discovery."}


if __name__=="__main__":
    result=run()
    assert result==run()
    print(json.dumps(result,sort_keys=True,indent=2))
