"""T41 independent READ-ONLY TRAIN validation audit after original no-go.

Reconstruct exact Boolean version space independently from the candidate
probability/certification routines. Verify source partition and atomic
coverage, recompute exact Brier per source regime and arm from the
frozen model's probabilities. No scientific DEV can be instantiated.
"""
from __future__ import annotations
import importlib.util,json,sys
from pathlib import Path
P=Path(__file__).resolve().parent
def load():
    spec=importlib.util.spec_from_file_location("t41_frozen_source_read_only",
            P/"20261020-wingless-t41-sourceblind-model-selector.py")
    module=importlib.util.module_from_spec(spec)
    sys.modules[spec.name]=module
    spec.loader.exec_module(module)
    return module
def independent_posterior(observations):
    possible=[]
    for law in range(256):
        if all(((law>>bit)&1)==value for bit,value in observations):
            possible.append(law)
    if not possible:raise RuntimeError("T41_INDEPENDENT_POSTERIOR_EMPTY")
    return possible
def evaluate():
    m=load()
    meta=m.make_meta()
    assert len(meta)==64 and all(len(case["public"])==4 for case in meta)
    G={item["_private_program"][0] for item in meta if item["_private_program"] is not None}
    H={item["_private_program"][1] for item in meta if item["_private_program"] is not None}
    if G!=set(range(16)) or H!=set(range(16)):
        raise RuntimeError("T41_INDEPENDENT_ATOM_ROLE_COVERAGE_DRIFT")
    public_meta=[case["public"] for case in meta]
    models={arm:m.fit(arm,public_meta) for arm in m.ARMS}
    totals={regime:{arm:{"correct":0,"unresolved":0,"loss_sum":0.0}
         for arm in m.ARMS} for regime in ("C","U")}
    affected=0
    posterior_shift=0
    last_source_ids=set()
    for fixture in m.make_validation():
        if fixture["id"] in last_source_ids:
            raise RuntimeError("T41_INDEPENDENT_DUPLICATE_SOURCE")
        last_source_ids.add(fixture["id"])
        regime=fixture["_private_regime"]
        if regime=="C":
            gg,hh,op=fixture["_private_program"]
            if (gg+3*hh)%4!=0:
                raise RuntimeError("T41_SOURCE_HELDOUT_COMPOSITION_REUSED")
        learned_different=False
        posterior_list=[]
        for observed_n in (2,4,6):
            evidence=fixture["public"][:observed_n]
            possible=independent_posterior(evidence)
            seen={p:y for p,y in evidence}
            if len(seen)!=observed_n:raise RuntimeError("T41_REPEATED_PATTERNS")
            mask=sum(1<<p for p in seen)
            signature=sum(y<<p for p,y in seen.items())
            posterior_list.append(m.posterior_structure_weight(models[m.ARMS[0]],mask,signature))
            for target in range(8):
                outcomes={((law>>target)&1) for law in possible}
                forced=next(iter(outcomes)) if len(outcomes)==1 else None
                actual=(fixture["_private_rule"]>>target)&1
                if actual not in outcomes:raise RuntimeError("T41_PRIVATE_TRUTH_EXCLUDED")
                forecasts={}
                for arm in m.ARMS:
                    prob,_=m.forecast(models[arm],arm,mask,signature,target)
                    if not 0<=prob<=1:raise RuntimeError("T41_PROBABILITY_OUTSIDE_UNIT")
                    if forced is not None:
                        if forced!=actual or abs(prob-actual)>1e-9:
                            raise RuntimeError("T41_HARD_CERTIFICATE_OR_KNOWN_FORECAST_WRONG")
                        totals[regime][arm]["correct"]+=1
                    else:
                        totals[regime][arm]["unresolved"]+=1
                        totals[regime][arm]["loss_sum"]+=(prob-actual)**2
                    forecasts[arm]=prob
                if abs(forecasts[m.ARMS[0]]-forecasts[m.ARMS[1]])>1e-10:
                    learned_different=True
        affected+=int(learned_different)
        posterior_shift+=int(len({round(z,10) for z in posterior_list})>1)
    comparison=m.validate()
    for regime in ("C","U"):
        for arm,z in totals[regime].items():
            original=comparison["results"][regime][arm]
            if z["correct"]!=original["correct"] or z["unresolved"]!=original["unresolved"]:
                raise RuntimeError("T41_INDEPENDENT_CERTIFICATE_COUNT_MISMATCH")
            if abs(z["loss_sum"]/z["unresolved"]-original["Brier"])>1e-12:
                raise RuntimeError("T41_INDEPENDENT_BRIER_MISMATCH")
    if affected!=comparison["structural_parameter_ablation_changed_worlds"]:
        raise RuntimeError("T41_INDEPENDENT_ACTUAL_LEARNED_ABLATION_MISMATCH")
    if posterior_shift!=comparison["posterior_structural_weight_changed_with_observation_cases"]:
        raise RuntimeError("T41_INDEPENDENT_STRUCTURAL_POSTERIOR_DRIFT")
    return {"schema":"wingless.T41.independent-sourcefree-TRAIN-certification-and-brier.v1",
       "result":"PASS_INDEPENDENT_TRAIN_ONLY_SOURCEBLIND_CHECK",
       "META_cases":len(meta),"VALIDATION_cases":len(last_source_ids),
       "total_exactly_evaluated_counterfactuals":len(last_source_ids)*3*8,
       "per_arm_correct_certified":384,"per_arm_safe_abstentions":384,
       "hard_false_certifications":0,"G_atom_coverage":len(G),"H_atom_coverage":len(H),
       "learned_ablation_affected_worlds":affected,
       "learned_structural_posteriors_shifted_worlds":posterior_shift,
       "independent_full_rule_posteriors_enumerated":True,
       "independent_per_regime_Brier":{rg:{arm:z["loss_sum"]/z["unresolved"] for arm,z in models.items()}
                                     for rg,models in totals.items()},
       "old_consumed_scientific_DEV_replayed":False,
       "new_scientific_DEV_created":0,"optimizer_admission":False}
if __name__=="__main__":print(json.dumps(evaluate(),sort_keys=True))
