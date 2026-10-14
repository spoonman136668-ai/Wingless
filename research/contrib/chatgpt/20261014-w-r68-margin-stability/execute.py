"""W-R68: paired fresh heldout feature-to-margin-to-decision causality."""
import importlib.util,json,math
from pathlib import Path
p=Path(__file__).resolve().parent.parent/"20261014-w-r67-phase-ladder"/"execute.py"
spec=importlib.util.spec_from_file_location("w67_frozen",p)
assert spec and spec.loader
w67=importlib.util.module_from_spec(spec);spec.loader.exec_module(w67)

HELDOUT=tuple(range(10208,10232))
PHASES=("intact","minus1","plus1")
AMPS=(1,2,3)

def evaluation(seed,entity,site,amp,weights,order,phase):
    m=w67.w66.w65.w64.w63.w61.m.m
    base,pos,neg=m.w37.recoded_example(
        seed,entity,"xorshift_independent","identity",site,amp,False)
    a,b,c=(w67.acquire(rows,order,phase) for rows in (base,pos,neg))
    positive,negative=m.w42.feature(a,b),m.w42.feature(a,c)
    assert len(positive)==len(negative)==len(weights)==16
    margin=sum(w*(x-y) for w,x,y in zip(weights,positive,negative))
    credit=m.w45.credit(margin,True)+m.w45.credit(-margin,False)
    return positive,negative,margin,credit

def run():
    trained=w67.w66.w65.w64.w63.w61.fit()
    records=[]
    for i in (1,2,3):
        model=f"{i}:withhold_xorshift_independent"
        fit=trained[model]
        assert fit["training_counts"]["xorshift_independent"]==0
        weights=fit["weights"]
        for amp in AMPS:
            for perm,order in w67.w66.w65.w64.PERMS.items():
                for seed in HELDOUT:
                    metrics={
                        "pairs":0,
                        "minus1_pos_feature_changed":0,
                        "minus1_neg_feature_changed":0,
                        "minus1_margin_changed":0,
                        "minus1_credit_changed":0,
                        "plus1_margin_changed":0,
                        "plus1_credit_changed":0,
                        "minus1_abs_margin_delta":0.0,
                        "plus1_abs_margin_delta":0.0,
                        "minus1_abs_credit_delta":0.0,
                        "plus1_abs_credit_delta":0.0
                    }
                    for entity in w67.w66.w65.w64.ENTITIES:
                        for site in w67.w66.w65.w64.SITES:
                            original=evaluation(seed,entity,site,amp,weights,order,"intact")
                            minus=evaluation(seed,entity,site,amp,weights,order,"minus1")
                            plus=evaluation(seed,entity,site,amp,weights,order,"plus1")
                            metrics["pairs"]+=1
                            metrics["minus1_pos_feature_changed"]+=int(original[0]!=minus[0])
                            metrics["minus1_neg_feature_changed"]+=int(original[1]!=minus[1])
                            metrics["minus1_margin_changed"]+=int(abs(original[2]-minus[2])>1e-9)
                            metrics["minus1_credit_changed"]+=int(abs(original[3]-minus[3])>1e-9)
                            metrics["plus1_margin_changed"]+=int(abs(original[2]-plus[2])>1e-9)
                            metrics["plus1_credit_changed"]+=int(abs(original[3]-plus[3])>1e-9)
                            metrics["minus1_abs_margin_delta"]+=abs(original[2]-minus[2])
                            metrics["plus1_abs_margin_delta"]+=abs(original[2]-plus[2])
                            metrics["minus1_abs_credit_delta"]+=abs(original[3]-minus[3])
                            metrics["plus1_abs_credit_delta"]+=abs(original[3]-plus[3])
                    assert metrics["pairs"]==len(w67.w66.w65.w64.ENTITIES)*len(w67.w66.w65.w64.SITES)
                    records.append({"model":model,"amplitude":amp,"permutation":perm,
                                    "seed":seed,"metrics":metrics})
    assert len(records)==3*3*4*24
    metrics={key:sum(r["metrics"][key] for r in records) for key in records[0]["metrics"]}
    n=metrics["pairs"]
    assert n==3*3*4*24*len(w67.w66.w65.w64.ENTITIES)*len(w67.w66.w65.w64.SITES)
    frac_margin=metrics["minus1_margin_changed"]/n
    frac_credit=metrics["minus1_credit_changed"]/n
    if frac_margin>=.5 and frac_credit<=.05:
        label="DECISION_STABLE_DESPITE_MARGIN_DRIFT"
    elif frac_margin<=.05:
        label="MARGIN_INVARIANCE"
    else:
        label="MIXED_MARGIN_DECISION_EFFECT"
    return {"schema":"wingless.w-r68-paired-margin-causality.v1",
            "parent_head":"5a73ddf04a9bad1c2b41c8d5154e472f976ea4aa",
            "classification":label,
            "sealed_heldout_seeds":HELDOUT,
            "development_blocks":w67.w66.w65.w64.w63.w61.DEV,
            "generator":"xorshift_independent",
            "aggregate":metrics,
            "records":records,"model_calls":0,
            "accepted_ref_mutation":False,
            "claim_limit":"Synthetic fixed ridge feature-margin-decision causality, not learned self-repair or RSI."}
if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
