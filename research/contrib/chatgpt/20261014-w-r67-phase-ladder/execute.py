"""W-R67: prospective circular phase-ladder falsification of global-shift fragility."""
import importlib.util
import json
from pathlib import Path

p = Path(__file__).resolve().parent.parent / "20261014-w-r66-cross-channel-alignment" / "execute.py"
spec = importlib.util.spec_from_file_location("w66_frozen", p)
assert spec and spec.loader
w66 = importlib.util.module_from_spec(spec)
spec.loader.exec_module(w66)

HELDOUT = tuple(range(10112, 10136))
PHASES = {"intact":0, "plus1":1, "minus1":-1, "plus2":2,
          "minus2":-2, "plus4":4, "minus4":-4, "plus8":8}
AMPS = (1,2,3)

def acquire(rows, order, phase):
    assert len(rows)==16 and sorted(order)==[0,1,2,3]
    assert phase in PHASES
    shift=PHASES[phase]
    return tuple(tuple(rows[(t+shift)%16][ch] for ch in order) for t in range(16))

def rank(seed,entity,site,amp,weights,order,phase):
    m=w66.w65.w64.w63.w61.m.m
    a,b,c=m.w37.recoded_example(seed,entity,"xorshift_independent","identity",site,amp,False)
    a,b,c=(acquire(x,order,phase) for x in (a,b,c))
    p,n=m.w42.feature(a,b),m.w42.feature(a,c)
    assert len(p)==len(n)==len(weights)==16
    score=sum(w*(v-u) for w,v,u in zip(weights,p,n))
    return m.w45.credit(score,True)+m.w45.credit(-score,False)

def run():
    models=w66.w65.w64.w63.w61.fit()
    records=[]
    for i in (1,2,3):
        name=f"{i}:withhold_xorshift_independent"
        fit=models[name]
        assert fit["training_counts"]["xorshift_independent"]==0
        weights=fit["weights"]
        assert len(weights)==16
        for amplitude in AMPS:
            for perm,order in w66.w65.w64.PERMS.items():
                for phase in PHASES:
                    by_seed=[]
                    for seed in HELDOUT:
                        credit=sum(rank(seed,entity,site,amplitude,weights,order,phase)
                                   for entity in w66.w65.w64.ENTITIES
                                   for site in w66.w65.w64.SITES)
                        by_seed.append({"seed":seed,"correct_credit":credit,
                                        "presentations":len(w66.w65.w64.ENTITIES)*len(w66.w65.w64.SITES)*2})
                    numerator=sum(x["correct_credit"] for x in by_seed)
                    denominator=sum(x["presentations"] for x in by_seed)
                    records.append({"model":name,"amplitude":amplitude,"permutation":perm,
                                    "phase":phase,"credit":numerator,
                                    "presentations":denominator,"accuracy":numerator/denominator,
                                    "by_seed":by_seed})
    assert len(records)==3*3*4*8
    scores={(r["model"],r["amplitude"],r["permutation"],r["phase"]):r["accuracy"]
            for r in records}
    comparisons=[]
    for i in (1,2,3):
        model=f"{i}:withhold_xorshift_independent"
        for amp in AMPS:
            for perm in w66.w65.w64.PERMS:
                prefix=(model,amp,perm)
                intact=scores[prefix+("intact",)]
                delta={name:scores[prefix+(name,)]-intact for name in PHASES if name!="intact"}
                comparisons.append({"model":model,"amplitude":amp,"permutation":perm,
                                    "intact_accuracy":intact,"delta":delta})
    assert len(comparisons)==36
    mean={name:sum(c["delta"][name] for c in comparisons)/36
          for name in PHASES if name!="intact"}
    broad=sum(value<=-.06 for value in mean.values())>=6
    directional=any(abs(mean[f"plus{k}"]-mean[f"minus{k}"])>=.04 for k in (1,2,4))
    if broad: classification="BROAD_GLOBAL_PHASE_FRAGILITY"
    elif directional:classification="DIRECTIONAL_PHASE_ASYMMETRY"
    else: classification="MIXED_PHASE_DEPENDENCE"
    return {"schema":"wingless.w-r67-global-phase-ladder.v1",
            "parent_head":"1181a9d408cb4ca1ba5388105be45a26d74b0705",
            "classification":classification,"sealed_heldout_seeds":HELDOUT,
            "development_blocks":w66.w65.w64.w63.w61.DEV,
            "generator":"xorshift_independent","phases":PHASES,
            "mean_accuracy_change":mean,"comparisons":comparisons,
            "all_sealed_records":records,"model_calls":0,
            "accepted_ref_mutation":False,
            "claim_limit":"Synthetic source-withheld fixed ridge phase shifts; not proof of learned invariance or RSI."}

if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
