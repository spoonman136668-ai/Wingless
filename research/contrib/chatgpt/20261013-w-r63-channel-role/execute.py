"""W-R63: frozen source identity vs displayed channel-slot ablation."""
import importlib.util
import json
from pathlib import Path

PARENT = Path(__file__).resolve().parent.parent / "20261013-w-r62-channel-outage" / "execute.py"
spec = importlib.util.spec_from_file_location("w62_r2_pinned", PARENT)
assert spec and spec.loader
w62 = importlib.util.module_from_spec(spec)
spec.loader.exec_module(w62)
w61 = w62.w61
HELDOUT = tuple(range(9552, 9576))
PERMS = {"identity":(0,1,2,3), "rotate":(1,2,3,0),
         "swap01":(1,0,2,3), "reverse":(3,2,1,0)}
CONDITIONS = ("intact_permuted", "erase_presented_slot0",
              "erase_original_source_channel0")
AMPS = (1,2,3)
GENERATORS = w61.GENERATORS
GROUPS = w61.GROUPS
SITES = w61.m.m.TEST_SITES

def acquire(observations, order, condition):
    assert sorted(order)==[0,1,2,3] and len(observations)==16
    drop = None if condition=="intact_permuted" else (
        0 if condition=="erase_presented_slot0" else order.index(0))
    assert condition in CONDITIONS
    return tuple(tuple(0 if drop == slot else row[src]
                       for slot,src in enumerate(order)) for row in observations)

def rank(seed, entity, generator, site, amplitude, weights, order, condition):
    base,pos,neg=w61.m.m.w37.recoded_example(
        seed,entity,generator,"identity",site,amplitude,False)
    base=acquire(base,order,condition)
    pos=acquire(pos,order,condition)
    neg=acquire(neg,order,condition)
    p=w61.m.m.w42.feature(base,pos)
    n=w61.m.m.w42.feature(base,neg)
    assert len(p)==len(n)==len(weights)==16
    margin=sum(w*(a-b) for w,a,b in zip(weights,p,n))
    return w61.m.m.w45.credit(margin,True)+w61.m.m.w45.credit(-margin,False)

def run():
    trained=w61.fit()
    assert len(trained)==9
    results=[]
    for i in (1,2,3):
        name=f"{i}:withhold_xorshift_independent"
        weights=trained[name]["weights"]
        assert trained[name]["training_counts"]["xorshift_independent"]==0
        for amp in AMPS:
            for generator in GENERATORS:
                for group,entities in GROUPS.items():
                    for label,order in PERMS.items():
                        for condition in CONDITIONS:
                            seed_rows=[]
                            score=0.0
                            count=0
                            for seed in HELDOUT:
                                credit=0.0
                                presentations=0
                                for site in SITES:
                                    for entity in entities:
                                        credit+=rank(seed,entity,generator,site,amp,
                                                     weights,order,condition)
                                        presentations+=2
                                seed_rows.append({"seed":seed,"credit":credit,
                                                  "presentations":presentations})
                                score+=credit
                                count+=presentations
                            results.append({"model":name,"amplitude":amp,
                                            "generator":generator,"entity_group":group,
                                            "permutation":label,"condition":condition,
                                            "accuracy":score/count,"credit":score,
                                            "presentations":count,"by_seed":seed_rows})
    assert len(results)==3*3*3*2*4*3
    scores={(r["model"],r["amplitude"],r["generator"],r["entity_group"],
             r["permutation"],r["condition"]):r["accuracy"] for r in results}
    comparisons=[]
    for i in (1,2,3):
        name=f"{i}:withhold_xorshift_independent"
        for amp in AMPS:
            for perm in ("rotate","swap01","reverse"):
                p=(name,amp,"xorshift_independent","unseen_identity_transfer",perm)
                intact=scores[p+("intact_permuted",)]
                slot=scores[p+("erase_presented_slot0",)]
                source=scores[p+("erase_original_source_channel0",)]
                comparisons.append({"model":name,"amplitude":amp,"permutation":perm,
                                    "intact_accuracy":intact,
                                    "slot0_accuracy":slot,"source0_accuracy":source,
                                    "slot_minus_intact":slot-intact,
                                    "source_minus_intact":source-intact})
    assert len(comparisons)==27
    slot_mean=sum(x["slot_minus_intact"] for x in comparisons)/27
    source_mean=sum(x["source_minus_intact"] for x in comparisons)/27
    if source_mean<=-.07 and slot_mean>=-.05:
        verdict="SOURCE_CHANNEL_ROLE"
    elif slot_mean<=-.07 and source_mean>=-.05:
        verdict="FEATURE_SLOT_ROLE"
    else:
        verdict="MIXED_OR_INCONCLUSIVE"
    return {"schema":"wingless.w-r63-channel-role.v1",
            "parent_head":"6825ce8eab1e3ee4ca3434853466da82dd082237",
            "classification":verdict,"sealed_heldout_seeds":HELDOUT,
            "training_blocks":w61.DEV,"permutations":PERMS,
            "amplitudes":AMPS,"generators":GENERATORS,
            "entity_groups":GROUPS,"conditions":CONDITIONS,
            "source_channel_mean_effect":source_mean,
            "presented_slot_mean_effect":slot_mean,
            "primary_unseen_xorshift_comparisons":comparisons,
            "all_sealed_records":results,"model_calls":0,
            "accepted_ref_mutation":False,
            "claim_limit":"Synthetic source versus feature-slot channel erasure mechanism; no semantics or general intelligence."}

if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
