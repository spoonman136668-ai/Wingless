"""T38 independent post-first-TRAIN exact source and scorer audit.

Read only frozen source/candidate. Independently derive public causal
truth-table clauses and query outcomes from original history arrays.
Explicitly enumerate all 256 Boolean rules; independently compute held-out
validation Brier losses and compare to candidate-reported TRAIN values.
Inspect TRAIN private generator identities for source provenance only;
never give private gates/rules or future labels to the fitted model.
Original T38 DEV identities are NOT built, tested, or consumed.
"""
from __future__ import annotations
import importlib.util,json,sys
from collections import Counter
from pathlib import Path
ROOT=Path(__file__).resolve().parent
def load(name,file):
    spec=importlib.util.spec_from_file_location(name,ROOT/file)
    x=importlib.util.module_from_spec(spec)
    sys.modules[spec.name]=x;spec.loader.exec_module(x)
    return x
def p(history,parents):
    a,b,c=parents
    return ((history[0]>>a)&1)|(((history[1]>>b)&1)<<1)|(((history[2]>>c)&1)<<2)
def public_allowed(pub):
    rows=[]
    for case in pub["probes"]:
        for j,y in case["observed"]:
            if case["do"] is not None and case["do"][0]==j:
                if case["do"][1]!=y:raise RuntimeError("T38_INDEPENDENT_ACTION_WRONG")
                continue
            idx=p(case["history"],pub["topology"][j])
            rows.append((idx,y^((case["history"][0]>>j)&1)))
    possible=[rule for rule in range(256) if all(((rule>>i)&1)==b for i,b in rows)]
    if not possible:raise RuntimeError("T38_INDEPENDENT_EXACT_POSTERIOR_EMPTY")
    return possible
def check():
    s=load("t38_new_orig_source_independent","20261020-wingless-t38-new-compositional-sources.py")
    learner=load("t38_frozen_original_candidate_readonly","20261020-wingless-t38-learned-composition.py")
    assert len(s.TRAIN_PAIRS)==48 and len(s.DEV_PAIRS)==16
    assert set(s.TRAIN_PAIRS).isdisjoint(s.DEV_PAIRS)
    # Verify the 128 2-atom truth tables independently, not with the source
    # truth-table composition helper.
    for g,h,op,table in s.HYPOTHESES:
        reconstruct=0
        for index in range(8):
            a=index&1;b=(index>>1)&1;c=(index>>2)&1
            x=(s.GATES[g]>>(a|(b<<1)))&1
            y=(s.GATES[h]>>(b|(c<<1)))&1
            v=x^y if op==0 else x&y
            reconstruct|=v<<index
        if table!=reconstruct:raise RuntimeError("T38_INDEPENDENT_ATOM_COMPOSITION_WRONG")
    trainings=[]
    meta_pair=set();meta_g=set();meta_h=set()
    for identity in s.TRAIN_META:
        world=s.build(identity)
        meta_pair.add((world["_private_g"],world["_private_h"]))
        meta_g.add(world["_private_g"]);meta_h.add(world["_private_h"])
        if (world["_private_g"],world["_private_h"]) not in s.TRAIN_PAIRS:
            raise RuntimeError("T38_TRAIN_SOURCE_NOT_IN_TRAIN_PARTITION")
        trainings.append(s.META(world))
    fitted={arm:learner.train(s,trainings,arm) for arm in learner.ARMS}
    result={arm:{"sum_brier_unresolved":0.0,"unresolved":0,
         "correct_certified":0} for arm in learner.ARMS}
    numeric_ablation_worlds=0
    allowed_worlds=0;scored=0
    for identity in s.TRAIN_VALID:
        world=s.build(identity);pub=s.VALID(world)
        if (world["_private_g"],world["_private_h"]) not in s.TRAIN_PAIRS:
            raise RuntimeError("T38_TRAIN_VALIDATION_NOT_FROZEN_TRAIN_PAIR")
        laws=public_allowed(pub)
        if world["_private_rule"] not in laws:
            raise RuntimeError("T38_INDEPENDENT_TRUE_RULE_NOT_RETAINED")
        allowed_worlds+=1
        preds={arm:learner.future(s,pub,learner.observed_rule_bits(pub,s),fitted[arm],arm)
               for arm in learner.ARMS}
        off={"G":(1/8,)*8,"H":(1/8,)*8,"OP":(0.5,0.5)}
        disabled=learner.future(s,pub,learner.observed_rule_bits(pub,s),off,learner.ARMS[0])
        if any(abs(a["P1"]-b["P1"])>1e-9 for a,b in zip(preds[learner.ARMS[0]],disabled)):
            numeric_ablation_worlds+=1
        for qi,q in enumerate(pub["queries"]):
            j=q["node"]
            action=q["do"]
            predicted_from_source=set()
            for law in laws:
                if action is not None and action[0]==j:
                    value=action[1]
                else:
                    idx=p(q["history"],pub["topology"][j])
                    value=((q["history"][0]>>j)&1)^((law>>idx)&1)
                predicted_from_source.add(value)
            proof=next(iter(predicted_from_source)) if len(predicted_from_source)==1 else None
            true=world["_answers"][qi]
            if true not in predicted_from_source:
                raise RuntimeError("T38_INDEPENDENT_TRUE_COUNTERFACTUAL_EXCLUDED")
            for arm,outputs in preds.items():
                pred=outputs[qi]
                if pred["certified"]!=proof:
                    raise RuntimeError("T38_INDEPENDENT_INCORRECT_CERTIFICATE:"+arm)
                if not 0<=pred["P1"]<=1:raise RuntimeError("T38_INVALID_PROBABILITY")
                if proof is None:
                    result[arm]["unresolved"]+=1
                    result[arm]["sum_brier_unresolved"]+=(pred["P1"]-true)**2
                else:
                    if proof!=true:raise RuntimeError("T38_INDEPENDENT_FALSE_CERTIFIED")
                    result[arm]["correct_certified"]+=1
            scored+=1
    for x in result.values():
        x["Brier_unresolved"]=x["sum_brier_unresolved"]/max(1,x["unresolved"])
    orig=learner.proof()
    for arm,x in result.items():
        expected=orig["learned_Brier_unresolved"] if arm==learner.ARMS[0] else orig["competitor_Brier_unresolved"][arm]
        if abs(x["Brier_unresolved"]-expected)>1e-12:
            raise RuntimeError("T38_INDEPENDENT_BRIER_RECALCULATION_MISMATCH:"+arm)
    if numeric_ablation_worlds!=orig["learned_parameters_changed_actual_belief_worlds"]:
        raise RuntimeError("T38_INDEPENDENT_PARAMETER_ABLATION_MISMATCH")
    return {"schema":"wingless.T38.independent-frozen-TRAIN-compositional-source-and-scorer.v1",
        "result":"PASS_READONLY_INDEPENDENT_EXACT_TRAIN_VALIDATION",
        "full_composed_tables_verified":128,
        "TRAIN_source_atom_G_coverage":len(meta_g),
        "TRAIN_source_atom_H_coverage":len(meta_h),
        "TRAIN_meta_distinct_gate_pairs":len(meta_pair),
        "heldout_DEV_pair_class_disjoint_and_not_generated":True,
        "TRAIN_meta_public_cases":40,"TRAIN_validation_independently_checked":allowed_worlds,
        "counterfactuals_exactly_verified":scored,
        "parameter_ablation_changed_validation_cases":numeric_ablation_worlds,
        "all_arms_recomputed_Brier":{arm:{"unresolved":x["unresolved"],
           "Brier_unresolved":x["Brier_unresolved"],
           "correct_certified":x["correct_certified"]} for arm,x in result.items()},
        "no_original_DEV_generation":True,"no_old_DEV_replay":True,
        "no_optimizer_admission":True}
if __name__=="__main__":print(json.dumps(check(),sort_keys=True))
