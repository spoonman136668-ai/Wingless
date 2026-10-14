"""W-R72: three-phase fixed-budget curriculum learned phase augmentation versus frozen ridge baseline.

Three preregistered 16-weight ridge learners: original training, +1-only training,
and balanced-phase augmentation. Training sample count/source balance is exactly
unchanged; only development-time temporal representation differs. All scores
use fresh unseen-source and unseen-entity heldout cases.
"""
import hashlib
import importlib.util
import json
from pathlib import Path

PARENT = Path(__file__).resolve().parent.parent / "20261014-w-r68-minus1-replication" / "execute.py"
spec = importlib.util.spec_from_file_location("w68_frozen", PARENT)
assert spec and spec.loader
w68 = importlib.util.module_from_spec(spec)
spec.loader.exec_module(w68)

w67 = w68
w61 = w67.w66.w65.w64.w63.w61
m = w61.m.m

DEV = w61.DEV
HELDOUT = tuple(range(15360, 15376))
TRAIN_ARMS = ("original_ridge", "plus1_curriculum", "three_phase_curriculum")
TEST_PHASES = ("intact", "plus1", "minus1", "plus2")
AMPS = (1, 2, 3)
PERMS = w67.w66.w65.w64.PERMS
ENTITIES = w67.w66.w65.w64.ENTITIES
SITES = w67.w66.w65.w64.SITES
EXCLUDED = "xorshift_independent"

def augmented_development(seeds,arm):
    assert arm in TRAIN_ARMS
    pairs = []
    counts = {g:0 for g in m.GENERATORS}
    pair_index = 0
    for seed in seeds:
        for gen_index,generator in enumerate(m.GENERATORS):
            if generator == EXCLUDED:
                continue
            for entity in range(9):
                if entity % 2 != (seed + gen_index) % 2:
                    continue
                for site in m.TRAIN_SITES:
                    for amplitude in (1,2):
                        if arm == "plus1_curriculum":
                            phase = "plus1"
                        elif arm == "three_phase_curriculum":
                            phase = ("intact","plus1","minus1")[pair_index % 3]
                        else:
                            phase = "intact"
                        base,pos,neg = m.w37.recoded_example(
                            seed,entity,generator,"identity",site,amplitude,True)
                        a,b,c = (w67.acquire(rows,(0,1,2,3),phase)
                                 for rows in (base,pos,neg))
                        positive = m.w42.feature(a,b)
                        negative = m.w42.feature(a,c)
                        delta = tuple(x-y for x,y in zip(positive,negative))
                        assert len(delta)==16
                        pairs.extend(((delta,1),(tuple(-x for x in delta),-1)))
                        counts[generator] += 2
                        pair_index += 1
    assert len(pairs)==1728 and pair_index==864
    assert counts[EXCLUDED]==0 and all(
        counts[g]==864 for g in m.GENERATORS if g!=EXCLUDED)
    digest = hashlib.sha256(json.dumps(
        [[list(v),int(y)] for v,y in pairs],separators=(",",":")).encode()).hexdigest()
    return pairs,counts,digest

def fit():
    results = {}
    for idx,seeds in enumerate(DEV,1):
        expected,base_counts,base_digest = m.build_examples(seeds,EXCLUDED)
        base_entries,c,sha = augmented_development(seeds,"original_ridge")
        assert base_entries == expected and c == base_counts and sha == base_digest
        for arm in TRAIN_ARMS:
            data,counts,digest = (base_entries,c,sha) if arm=="original_ridge" \
                else augmented_development(seeds,arm)
            weights=m.w53.ridge_fit(data)
            assert len(weights)==16
            key=(idx,arm)
            results[key]={"weights":weights,"source_counts":counts,
                          "development_fingerprint":digest,
                          "weight_sha256":hashlib.sha256(json.dumps(
                              weights,separators=(",",":")).encode()).hexdigest(),
                          "examples":len(data)}
    assert len(results)==9
    return results

def score(seed,entity,site,amp,weights,order,phase):
    base,pos,neg=m.w37.recoded_example(
        seed,entity,EXCLUDED,"identity",site,amp,False)
    a,b,c=(w67.acquire(rows,order,phase) for rows in (base,pos,neg))
    p,n=m.w42.feature(a,b),m.w42.feature(a,c)
    assert len(p)==len(n)==len(weights)==16
    margin=sum(w*(x-y) for w,x,y in zip(weights,p,n))
    return m.w45.credit(margin,True)+m.w45.credit(-margin,False)

def run():
    trained=fit()
    records=[]
    totals={arm:{phase:{"credit":0.0,"presentations":0} for phase in TEST_PHASES}
            for arm in TRAIN_ARMS}
    cohorts={i:{arm:{phase:{"credit":0.0,"presentations":0} for phase in TEST_PHASES}
                for arm in TRAIN_ARMS} for i in (1,2,3)}
    for i in (1,2,3):
        for arm in TRAIN_ARMS:
            weights=trained[(i,arm)]["weights"]
            for amp in AMPS:
                for perm,order in PERMS.items():
                    for phase in TEST_PHASES:
                        for seed in HELDOUT:
                            credit=sum(score(seed,entity,site,amp,weights,order,phase)
                                       for entity in ENTITIES for site in SITES)
                            presentation_count=2*len(ENTITIES)*len(SITES)
                            records.append({"cohort":i,"arm":arm,"amplitude":amp,
                                            "permutation":perm,"phase":phase,"seed":seed,
                                            "credit":credit,"presentations":presentation_count})
                            for bucket in (totals[arm][phase],cohorts[i][arm][phase]):
                                bucket["credit"]+=credit
                                bucket["presentations"]+=presentation_count
    assert len(records)==3*3*3*len(PERMS)*len(TEST_PHASES)*len(HELDOUT)
    acc={arm:{p:totals[arm][p]["credit"]/totals[arm][p]["presentations"]
              for p in TEST_PHASES} for arm in TRAIN_ARMS}
    delta={p:acc["three_phase_curriculum"][p]-acc["original_ridge"][p]
           for p in TEST_PHASES}
    specialist_delta={p:acc["three_phase_curriculum"][p]-acc["plus1_curriculum"][p]
                      for p in TEST_PHASES}
    plus_per_cohort=[cohorts[i]["three_phase_curriculum"]["plus1"]["credit"]/
                     cohorts[i]["three_phase_curriculum"]["plus1"]["presentations"]-
                     cohorts[i]["original_ridge"]["plus1"]["credit"]/
                     cohorts[i]["original_ridge"]["plus1"]["presentations"]
                     for i in (1,2,3)]
    if (acc["three_phase_curriculum"]["plus1"]>=.65
        and acc["three_phase_curriculum"]["minus1"]>=.65
        and acc["three_phase_curriculum"]["intact"]>=.65
        and sum(x>=.02 for x in plus_per_cohort)>=2):
        classification="DUAL_PHASE_ROBUST_LEARNED_READOUT"
    elif (acc["three_phase_curriculum"]["plus1"]>=.60
          and acc["three_phase_curriculum"]["minus1"]>=.60
          and acc["three_phase_curriculum"]["intact"]>=.60):
        classification="PARTIAL_DUAL_PHASE_READOUT"
    else:
        classification="NO_DUAL_PHASE_ROBUSTNESS"
    public_models={f"{i}:{arm}":{k:v for k,v in trained[(i,arm)].items() if k!="weights"}
                   for i in (1,2,3) for arm in TRAIN_ARMS}
    return {"schema":"wingless.w-r72-three-phase-curriculum.v1",
            "parent_head":"2a5d5605fb5bd638723c93950ea2379d9ed804b3",
            "classification":classification,"development_blocks":DEV,
            "sealed_heldout_seeds":HELDOUT,"test_phases":TEST_PHASES,
            "training_arms":TRAIN_ARMS,"source_excluded":EXCLUDED,
            "model_width":16,"source_counts_per_fit":{"xorshift_independent":0},
            "training_examples_per_fit":1728,"model_data":public_models,
            "aggregate_accuracy":acc,"balanced_minus_baseline":delta,
            "plus1_cohort_deltas":plus_per_cohort,"dual_minus_plus1_specialist":specialist_delta,"all_heldout_records":records,
            "model_calls":0,"accepted_ref_mutation":False,
            "claim_limit":"Synthetic ridge readout training under equal example/width budget; neither learned raw state formation nor RSI."}

if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
