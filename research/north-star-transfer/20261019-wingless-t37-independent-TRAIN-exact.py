"""T37 independent TRAIN-only exact verifier.

Independently reconstructs causal response equations from PUBLIC histories,
parent indices, interventions and observation bits. Does not import or
reuse candidate evidence() or all_possible(), and enumerates 256 possible
hidden Boolean laws directly. Candidate forecast/certification is checked
against these independently computed full truth-table posteriors.

Only original TRAIN_meta + TRAIN_validation IDs are built. No DEV or
consumed previous research source is generated or touched.
"""
from __future__ import annotations
import importlib.util,json,sys
from pathlib import Path
ROOT=Path(__file__).resolve().parent
def load(name,path):
    spec=importlib.util.spec_from_file_location(name,ROOT/path)
    m=importlib.util.module_from_spec(spec)
    sys.modules[spec.name]=m;spec.loader.exec_module(m)
    return m
def rootpattern(history,parents):
    return ((history[0]>>parents[0])&1) | (((history[1]>>parents[1])&1)<<1) | (((history[2]>>parents[2])&1)<<2)
def independent_worlds(public):
    observed=[]
    for probe in public["sequences"]:
        for node,output in probe["observed"]:
            action=probe["do"]
            if action is not None and node==action[0]:
                if output!=action[1]:
                    raise RuntimeError("T37_INDEPENDENT_KNOWN_INTERVENTION_WRONG")
                continue
            pattern=rootpattern(probe["history"],public["topology"][node])
            residual=output^((probe["history"][0]>>node)&1)
            observed.append((pattern,residual))
    result=[]
    for truth_table in range(256):
        if all(((truth_table>>pattern)&1)==value for pattern,value in observed):
            result.append(truth_table)
    if not result:raise RuntimeError("T37_INDEPENDENT_EXACT_WORLDSET_EMPTY")
    return result
def check():
    source=load("t37_source_for_independent_check","20261019-wingless-t37-sequential-rule-sources.py")
    candidate=load("t37_candidate_for_independent_check","20261019-wingless-t37-learned-dynamics.py")
    training=[source.public_META(source.build(x)) for x in source.TRAIN_META]
    fitted=candidate.learn(source,training)
    worlds_retained=certified=abstained=checked=cases_changed=0
    for key in source.TRAIN_FIT:
        item=source.build(key)
        public=source.public_FIT(item)
        possible=independent_worlds(public)
        if item["_hidden_law"] not in possible:
            raise RuntimeError("T37_INDEPENDENT_TRUE_CAUSAL_DYNAMICS_ELIMINATED")
        worlds_retained+=1
        actual={arm:candidate.infer(public,source,fitted[public["family"]],arm)
            for arm in candidate.ARMS}
        learned=actual[candidate.ARMS[0]]
        independent_beta=actual[candidate.ARMS[2]]
        if any(abs(a["probability"]-b["probability"])>1e-9
               for a,b in zip(learned,independent_beta)):
            cases_changed+=1
        for index,query in enumerate(public["queries"]):
            h=query["history"];j=query["node"];action=query["do"]
            predictions=[]
            for law in possible:
                if action is not None and action[0]==j:
                    val=action[1]
                else:
                    p=rootpattern(h,public["topology"][j])
                    val=((h[0]>>j)&1)^((law>>p)&1)
                predictions.append(val)
            permitted=set(predictions)
            expected=next(iter(permitted)) if len(permitted)==1 else None
            if item["_truth"][index] not in permitted:
                raise RuntimeError("T37_INDEPENDENT_PRIVATE_TRUTH_INVALID")
            for arm,series in actual.items():
                claimed=series[index]["certified"]
                if claimed!=expected:
                    raise RuntimeError("T37_INDEPENDENT_NOT_EXACT_"+arm)
                if not 0<=series[index]["probability"]<=1:
                    raise RuntimeError("T37_PROBABILISTIC_BELIEF_INVALID_"+arm)
            checked+=1
            if expected is None:abstained+=1
            else:
                if expected!=item["_truth"][index]:
                    raise RuntimeError("T37_INDEPENDENT_HARD_FALSE_CERTIFICATION")
                certified+=1
    return {"schema":"wingless.T37.independent-exact-train-causal-source-reconstruction.v1",
       "result":"INDEPENDENT_TRAIN_ONLY_EXACT_PASS",
       "TRAIN_meta_cases_used_for_fit":len(training),
       "TRAIN_validation_cases_independently_checked":worlds_retained,
       "causal_future_queries_independently_proved":checked,
       "correct_certified":certified,"sound_abstentions":abstained,
       "hard_false_certifications":0,
       "independent_full_256_rule_enumeration_not_candidate_code":True,
       "learned_parameters_changed_actual_causal_beliefs_cases":cases_changed,
       "scientific_DEV_generated":0,"no_prior_consumed_scientific_source_touched":True}
if __name__=="__main__":print(json.dumps(check(),sort_keys=True))
