"""T25 frozen prereg TRAIN-only headroom: count distinct latent parent pairs.

Reconciles the predeclared *causal parent ambiguity* admission metric with
the earlier non-equivalent output-certifier abstention metric. No source or
candidate modification. DEV identities are never built/opened.
"""
from __future__ import annotations
import importlib.util,json,sys
from pathlib import Path
ROOT=Path(__file__).parent
def load(n,p):
    spec=importlib.util.spec_from_file_location(n,ROOT/p)
    obj=importlib.util.module_from_spec(spec);sys.modules[n]=obj
    spec.loader.exec_module(obj);return obj
def main():
    source=load("t25_metric_source","20261019-wingless-t25-new-world-source.py")
    belief=load("t25_metric_exact","20261019-wingless-t25-causal-belief-foundation.py")
    unresolved=0;gate_only=0;counterfactual_abstain=0;queries=0;distinct_counts=[]
    per_family={"P":{"total":0,"unresolved":0},"Q":{"total":0,"unresolved":0}}
    for identity in source.TRAIN:
        world=source.build(identity)
        b=belief.posterior(source.train_public(world))
        for query,truth in zip(world["public"]["queries"],world["_private_truth"]):
            queries+=1
            target=query["target"]
            parents={(a,c) for a,c,g in b[target]}
            distinct_counts.append(len(parents))
            per_family[world["family"]]["total"]+=1
            if len(parents)>1:
                unresolved+=1;per_family[world["family"]]["unresolved"]+=1
            elif len(b[target])>1:gate_only+=1
            value,possible=belief.infer(b,query)
            if value is None:counterfactual_abstain+=1
            elif value!=truth:raise RuntimeError("T25_CAUSAL_TRUTH_MISMATCH")
    assert queries==120
    return {"schema":"wingless.T25.frozen-parent-ambiguity.verified-TRAIN.v1",
      "status":"PASS_TRAIN_ONLY_MEASUREMENT",
      "parent_ambiguous_target_questions":unresolved,
      "parent_ambiguity_fraction":unresolved/queries,
      "gate_ambiguity_only_target_questions":gate_only,
      "output_abstention_questions":counterfactual_abstain,
      "queries":queries,"minimum_frozen_parent_ambiguity_fraction":0.20,
      "admission_passed":unresolved/queries>=0.20,
      "source_families":per_family,
      "mean_distinct_parent_pairs":sum(distinct_counts)/len(distinct_counts),
      "DEV_cases_generated":0,"source_OR_model_blob_edited":False,
      "learned_candidate_trained":False,"optimizer_admission":False,
      "no_change_to_hypothesis":True}
if __name__=="__main__":print(json.dumps(main(),sort_keys=True))
