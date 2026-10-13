"""W-R33: delayed feedback controls reversibility of a two-table forecasting model."""
import hashlib
import importlib.util
import json
from pathlib import Path

PREVIOUS=Path(__file__).resolve().parent.parent/"20261013-w-r26-latent-alphabet"/"execute.py"
spec=importlib.util.spec_from_file_location("w26_frozen",PREVIOUS)
assert spec and spec.loader
w26=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w26)
assert w26.fixture()=="PASS" and len(w26.PERMUTATIONS)==24
ARMS=("FROZEN_ADD","ONE_WAY_SWITCH","TWO_WAY_SWITCH")
PHASES=(("pre",3,12),("mid",12,22),("late",22,32))


def ordering(seed):
    assert 200<=seed<=207
    return sorted(range(16),key=lambda e:(
        hashlib.sha256(("W33-order/"+str(seed)+"/"+str(e)).encode()).digest(),e))


def observed(seed,ep):
    assert 200<=seed<=207 and 0<=ep<16
    e=ordering(seed)[ep]
    xs=[(e+seed)%4,((e//4)+3*seed)%4]
    while len(xs)<32:
        t=len(xs)
        a,b=xs[-2:]
        new=(a^b) if seed%2==0 and 12<=t<22 else (a+b)%4
        xs.append(new)
    p=w26.PERMUTATIONS[hashlib.sha256(
        ("W33/"+str(seed)).encode()).digest()[0]%24]
    return tuple(p[v] for v in xs)


def other_table(permutation):
    inverse=[0]*4
    for latent,label in enumerate(permutation):
        inverse[label]=latent
    return tuple(permutation[inverse[a]^inverse[b]]
                 for a,b in w26.ALL_CONTEXTS)


def fit(seed):
    calibration=[observed(seed,e) for e in range(16)]
    cases=[(seq[0],seq[1],seq[2]) for seq in calibration]
    if len({(a,b) for a,b,_ in cases})!=16:
        return {"invalid":"AMBIGUOUS_HYPOTHESIS","calibration_contexts":len(set(
            (a,b) for a,b,_ in cases))}
    fits=[p for p in w26.PERMUTATIONS if w26.fits(p,cases)]
    add={w26.table_for(p) for p in fits}
    xor={other_table(p) for p in fits}
    if len(add)!=1 or len(xor)!=1:
        return {"invalid":"AMBIGUOUS_HYPOTHESIS","fits":len(fits),
                "add":len(add),"xor":len(xor)}
    a=next(iter(add));x=next(iter(xor))
    return {"add":a,"xor":x,"sha_add":hashlib.sha256(bytes(a)).hexdigest(),
            "sha_xor":hashlib.sha256(bytes(x)).hexdigest(),
            "candidates":len(fits),"calibration_labels":16}


def eval_arm(seed,arm,add,xor):
    phases={k:{"correct":0,"wrong":0} for k,_,_ in PHASES}
    matrix=[[0]*4 for _ in range(4)]
    trace=hashlib.sha256()
    switches=[]
    for e in range(16):
        xs=observed(seed,e)
        mode=0
        events=[]
        for t in range(3,32):
            name=next(k for k,start,end in PHASES if start<=t<end)
            table=xor if mode==1 else add
            predicted=table[4*xs[t-2]+xs[t-1]]
            true=xs[t]  # feedback after prediction
            if predicted==true:
                phases[name]["correct"]+=1
            else:
                phases[name]["wrong"]+=1
            matrix[true][predicted]+=1
            trace.update(bytes((e,t,true,predicted,mode)))
            if predicted!=true and arm!="FROZEN_ADD":
                if arm=="TWO_WAY_SWITCH" or mode==0:
                    mode=1-mode
                    events.append(t)
        switches.append(events)
    assert sum(z["correct"]+z["wrong"] for z in phases.values())==16*29
    return {"seed":seed,"cohort":"development" if seed<204 else "heldout",
            "condition":"DOUBLE_SWITCH" if seed%2==0 else "STABLE",
            "arm":arm,"phases":phases,
            "switches":switches,"switch_count":sum(map(len,switches)),
            "confusion":matrix,"trace_sha256":trace.hexdigest()}


def run():
    fingers={str(seed):hashlib.sha256(bytes(
        t for e in range(16) for t in observed(seed,e))).hexdigest()
        for seed in range(200,208)}
    if len(set(fingers.values()))!=8:
        return {"schema":"wingless.w-r33-reversible-regime.v1",
                "classification":"INVALID_IDENTITIES","fingerprints":fingers,
                "rows":[],"model_calls":0}
    inputs={}
    rows=[]
    for seed in range(200,208):
        result=fit(seed)
        if "invalid" in result:
            return {"schema":"wingless.w-r33-reversible-regime.v1",
                    "classification":"AMBIGUOUS_HYPOTHESIS","seed":seed,
                    "diagnostic":result,"rows":[],"model_calls":0}
        add=result.pop("add");xor=result.pop("xor")
        inputs[str(seed)]=result
        rows.extend(eval_arm(seed,arm,add,xor) for arm in ARMS)
    held=[r for r in rows if r["seed"]>=204]
    agg={condition:{arm:{
        **{phase:{metric:sum(r["phases"][phase][metric] for r in held
                              if r["condition"]==condition and r["arm"]==arm)
                     for metric in ("correct","wrong")}
           for phase,_,_ in PHASES},
        "switch_count":sum(r["switch_count"] for r in held
                           if r["condition"]==condition and r["arm"]==arm)}
        for arm in ARMS}
        for condition in ("STABLE","DOUBLE_SWITCH")}
    stable=agg["STABLE"]["TWO_WAY_SWITCH"]
    switch=agg["DOUBLE_SWITCH"]
    if (switch["TWO_WAY_SWITCH"]["late"]["correct"]>
        switch["ONE_WAY_SWITCH"]["late"]["correct"]
        and stable["switch_count"]==0
        and all(stable[p]["wrong"]==0 for p,_,_ in PHASES)
        and all(agg[condition][arm]["pre"]["wrong"]==0 for condition in agg for arm in ARMS)):
        classification="REVERSAL_RECOVERY"
    else:
        classification="MIXED"
    return {"schema":"wingless.w-r33-reversible-regime.v1",
            "classification":classification,
            "parent":"9a83922dc60d32c194477982aa38445dabdb77df",
            "fingerprints":fingers,"calibration":inputs,
            "heldout":agg,"rows":rows,"model_calls":0,
            "claim_limit":"Two predefined operators and full sequential truthful feedback; no free-form operator inference."}


if __name__=="__main__":
    result=run()
    assert result==run()
    print(json.dumps(result,sort_keys=True,indent=2))
