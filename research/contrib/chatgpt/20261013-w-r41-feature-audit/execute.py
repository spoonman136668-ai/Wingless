"""W-R41: exact observational and feature-label collision bounds, no model fitting."""
import importlib.util
import hashlib
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r40-positional"/"execute.py"
spec=importlib.util.spec_from_file_location("w40_pinned",SOURCE)
assert spec and spec.loader
w40=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w40)
w39=w40.w39
w37=w40.w37

HELDOUT=tuple(range(1920,1936))
CONDITIONS=("identity","xor55")
ARMS=("raw_pair","scalar_repeated","ordered_lag1","ordered_lag2","top16_sorted")
GENERATORS=w40.GENERATORS
SITES=w40.TEST_SITES


def feature(base,candidate,arm):
    if arm=="raw_pair":
        return tuple(v for row in (base+candidate) for v in row)
    if arm=="top16_sorted":
        return tuple(w39.feature(base,candidate,"top16_sorted"))
    return tuple(w40.feature(base,candidate,arm))


def run():
    results=[]
    sha=hashlib.sha256()
    for condition in CONDITIONS:
        counts={arm:{} for arm in ARMS}
        by_arm={arm:{"identical_pairs":0,"pairs":0,"by_seed":[]}
                for arm in ARMS}
        for generator in GENERATORS:
            for site in SITES:
                cell_by_arm={arm:{"generator":generator,"site":site,
                                  "identical_pairs":0,"pairs":0,"by_seed":[]}
                             for arm in ARMS}
                for seed in HELDOUT:
                    local={arm:0 for arm in ARMS}
                    raw_equal=0
                    for entity in range(9):
                        base,pos,neg=w37.recoded_example(
                            seed,entity,generator,condition,site,3,False)
                        raw_base,raw_pos,raw_neg=w37.recoded_example(
                            seed,entity,generator,"identity",site,3,False)
                        raw_equal+=int(pos==neg)
                        sha.update(json.dumps([condition,generator,site,seed,entity,
                                               base,pos,neg],
                                              separators=(",",":")).encode())
                        for arm in ARMS:
                            p=feature(base,pos,arm)
                            n=feature(base,neg,arm)
                            if arm!="raw_pair":
                                assert p==feature(raw_base,raw_pos,arm)
                                assert n==feature(raw_base,raw_neg,arm)
                            is_same=p==n
                            local[arm]+=int(is_same)
                            distribution=counts[arm]
                            distribution.setdefault(p,[0,0])[0]+=1
                            distribution.setdefault(n,[0,0])[1]+=1
                    for arm in ARMS:
                        entry=cell_by_arm[arm]
                        entry["identical_pairs"]+=local[arm]
                        entry["pairs"]+=9
                        entry["by_seed"].append({
                            "seed":seed,"identical_feature_pairs":local[arm],
                            "identical_raw_observations":raw_equal,
                            "counterfactual_pairs":9})
                for arm in ARMS:
                    by_arm[arm]["identical_pairs"]+=cell_by_arm[arm]["identical_pairs"]
                    by_arm[arm]["pairs"]+=cell_by_arm[arm]["pairs"]
                    by_arm[arm]["by_seed"].append(cell_by_arm[arm])
        for arm in ARMS:
            r=by_arm[arm]
            assert r["pairs"]==16*3*2*9
            r["arm"]=arm
            r["condition"]=condition
            r["identity_fraction"]=r["identical_pairs"]/r["pairs"]
            r["matched_pair_optimistic_accuracy_upper_bound"]=1-.5*r["identity_fraction"]
            r["exact_feature_memorization_upper_bound_on_this_sample"]=(
                sum(max(v) for v in counts[arm].values())/(2*r["pairs"]))
            r["distinct_feature_values"]=len(counts[arm])
            results.append(r)
    by={(a["condition"],a["arm"]):a for a in results}
    raw_rate=by["identity","raw_pair"]["identity_fraction"]
    max_feature_rate=max(by["identity",arm]["identity_fraction"]
                         for arm in ARMS if arm!="raw_pair")
    if raw_rate>=.20:
        label="OBSERVATION_LIMIT"
    elif raw_rate<.02 and max_feature_rate>=.20:
        label="FEATURE_COLLAPSE"
    else:
        label="UNEXPLAINED_BY_EXACT_PAIR_COLLISION"
    for arm in ARMS:
        if arm!="raw_pair":
            assert by["identity",arm]["identical_pairs"]==by["xor55",arm]["identical_pairs"]
    return {"schema":"wingless.w-r41-exact-feature-label-collisions.v1",
            "parent_head":"7f538f0fbe0e67418a6a52544d086557fd7e85b7",
            "classification":label,"sealed_heldout_seeds":HELDOUT,
            "generators":GENERATORS,"sites":SITES,
            "raw_observation_identical_fraction":raw_rate,
            "worst_feature_identical_fraction":max_feature_rate,
            "original_fixture_sha256":sha.hexdigest(),
            "arms":results,"training_updates":0,"model_calls":0,
            "claim_limit":"Exact finite-cohort collision upper bounds are evaluator-only oracles and not predictive model accuracy, training, or independent generalization."}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
