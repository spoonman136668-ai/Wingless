"""W-R69: frozen feature/margin/credit attribution of the minus-one exception."""
import importlib.util
import json
from pathlib import Path

PARENT=Path(__file__).resolve().parent.parent/"20261014-w-r68-minus1-replication"/"execute.py"
spec=importlib.util.spec_from_file_location("w68_frozen",PARENT)
assert spec and spec.loader
w68=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w68)

HELDOUT=tuple(range(12672,12688))
PHASES=("intact","minus1","plus1","minus2")
AMPS=(1,2,3)

def score(seed,entity,site,amp,weights,order,phase):
    m=w68.w67.w66.w65.w64.w63.w61.m.m
    a,b,c=m.w37.recoded_example(seed,entity,"xorshift_independent",
                                  "identity",site,amp,False)
    a,b,c=(w68.w67.acquire(x,order,phase) for x in (a,b,c))
    p,n=m.w42.feature(a,b),m.w42.feature(a,c)
    assert len(p)==len(n)==len(weights)==16
    margin=sum(w*(v-u) for w,v,u in zip(weights,p,n))
    credit=m.w45.credit(margin,True)+m.w45.credit(-margin,False)
    return {"positive":p,"negative":n,"margin":margin,"credit":credit}

def run():
    trained=w68.w67.w66.w65.w64.w63.w61.fit()
    permutations=w68.w67.w66.w65.w64.PERMS
    entities=w68.w67.w66.w65.w64.ENTITIES
    sites=w68.w67.w66.w65.w64.SITES
    rows=[]
    global_comparisons={p:{"examples":0,"feature_equal":0,
                           "margin_equal":0,"credit_equal":0,
                           "intact_credit":0.0,"phase_credit":0.0}
                        for p in PHASES if p!="intact"}
    for i in (1,2,3):
        model=f"{i}:withhold_xorshift_independent"
        fit=trained[model]
        assert fit["training_counts"]["xorshift_independent"]==0
        weights=fit["weights"]
        for amp in AMPS:
            for name,order in permutations.items():
                for seed in HELDOUT:
                    changes={p:{"examples":0,"feature_equal":0,"margin_equal":0,
                                "credit_equal":0,"intact_credit":0.0,
                                "phase_credit":0.0}
                             for p in global_comparisons}
                    for entity in entities:
                        for site in sites:
                            original=score(seed,entity,site,amp,weights,order,"intact")
                            for p in global_comparisons:
                                changed=score(seed,entity,site,amp,weights,order,p)
                                d=changes[p]
                                d["examples"]+=1
                                d["feature_equal"]+=int(changed["positive"]==original["positive"]
                                                        and changed["negative"]==original["negative"])
                                d["margin_equal"]+=int(changed["margin"]==original["margin"])
                                d["credit_equal"]+=int(changed["credit"]==original["credit"])
                                d["intact_credit"]+=original["credit"]
                                d["phase_credit"]+=changed["credit"]
                    rows.append({"model":model,"amplitude":amp,"permutation":name,
                                 "seed":seed,"comparisons":changes})
                    for p,d in changes.items():
                        for k,v in d.items():global_comparisons[p][k]+=v
    assert len(rows)==3*3*4*16
    for p,d in global_comparisons.items():
        assert d["examples"]==len(rows)*len(entities)*len(sites)
    minus=global_comparisons["minus1"]
    if minus["feature_equal"]==minus["examples"]:
        label="EXACT_FEATURE_LEVEL_MINUS1_INVARIANCE"
    elif minus["margin_equal"]==minus["examples"]:
        label="MARGIN_LEVEL_MINUS1_INVARIANCE"
    elif minus["credit_equal"]==minus["examples"]:
        label="CREDIT_LEVEL_MINUS1_INVARIANCE"
    elif all(r["comparisons"]["minus1"]["intact_credit"]==
             r["comparisons"]["minus1"]["phase_credit"] for r in rows):
        label="SEED_AGGREGATE_ONLY_MINUS1_EQUIVALENCE"
    else:
        label="MINUS1_EXCEPTION_BREAKS_AT_SEED_LEVEL"
    return {"schema":"wingless.w-r69-attributed-minusone-equivalence.v1",
            "parent_head":"9b1eac075006442c7e426ac0cb627d58102958b8",
            "classification":label,"sealed_heldout_seeds":HELDOUT,
            "source_family":"xorshift_independent","training_blocks":w68.w67.w66.w65.w64.w63.w61.DEV,
            "model_count":3,"amplitudes":AMPS,"permutations":permutations,
            "conditions":PHASES,"per_seed_diagnostics":rows,
            "aggregate":global_comparisons,"model_calls":0,
            "accepted_ref_mutation":False,
            "claim_limit":"Fixed synthetic ridge feature versus margin/credit attribution, not generalized temporal cognition or RSI."}

if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
