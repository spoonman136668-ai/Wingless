"""W-R68: prospectively falsify the exact -1 phase exception on a fresh heldout cohort."""
import importlib.util
import json
from pathlib import Path

p = Path(__file__).resolve().parent.parent / "20261014-w-r66-cross-channel-alignment" / "execute.py"
spec = importlib.util.spec_from_file_location("w66_frozen", p)
assert spec and spec.loader
w66 = importlib.util.module_from_spec(spec)
spec.loader.exec_module(w66)

HELDOUT = tuple(range(10432, 10456))
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
    exact_minus1_cells=sum(c["delta"]["minus1"]==0 for c in comparisons)
    if exact_minus1_cells==36 and mean["plus1"]<=-.06 and mean["minus2"]<=-.06:
        classification="MINUS_ONE_EXCEPTION_REPLICATED"
    elif mean["minus1"]<=-.06:
        classification="MINUS_ONE_EXCEPTION_FALSIFIED"
    else:
        classification="MINUS_ONE_PARTIAL_OR_INCONCLUSIVE"
    return {"schema":"wingless.w-r68-minus1-fresh-cohort.v1",
            "parent_head":"5a73ddf04a9bad1c2b41c8d5154e472f976ea4aa",
            "classification":classification,"sealed_heldout_seeds":HELDOUT,
            "development_blocks":w66.w65.w64.w63.w61.DEV,
            "generator":"xorshift_independent","phases":PHASES,
            "mean_accuracy_change":mean,"exact_minus1_zero_cells":exact_minus1_cells,"comparisons":comparisons,
            "all_sealed_records":records,"model_calls":0,
            "accepted_ref_mutation":False,
            "claim_limit":"Fresh source-withheld fixed ridge -1 phase exception; no general transfer or RSI."}

if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
