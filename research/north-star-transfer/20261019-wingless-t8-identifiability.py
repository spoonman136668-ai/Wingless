"""T8 Wingless public ambiguity impossibility and certified-abstention proof.

A known single calibration pair is compatible with two true causal
operators and two different outputs after exactly ten composition steps.
Thus any predictor overriding a fixed baseline can harm a baseline-correct
world using identical public information. A confidence prior alone cannot
guarantee zero regression. This does not consume or tune T6/T7 OPEN.
"""
from __future__ import annotations
import importlib.util,sys,json
from pathlib import Path
HERE=Path(__file__).resolve().parent
FROZEN=HERE/"20261019-wingless-t8-source-freeze.json"
KERNEL=HERE/"20261018-wingless-t3-typed-composition.py"
def need(ok,label):
    if not ok:raise AssertionError("WING_T8_"+label)
def load():
    spec=importlib.util.spec_from_file_location("w8_typed",KERNEL)
    mod=importlib.util.module_from_spec(spec);sys.modules[spec.name]=mod
    spec.loader.exec_module(mod);return mod
def compose(k,domain,roles,chain,start):
    s=k.start(domain)
    for idx,inverse in chain:
        s=k.compose_state(s,roles[idx],domain,inverse)
    return k.apply(s,start)
def ambiguity_witness(k,domain):
    if domain==k.AFFINE:
        role_public=tuple(((0,0),) for _ in range(4))
        left=((1,0),(1,0),(1,0),(1,0))
        right=((2,0),(1,0),(1,0),(1,0))
        start=1
        compatible=[r for r in k.candidates(domain) if k.operator_apply(domain,r,0)==0]
        need(len(compatible)==10,"TEN_INDISCERNIBLE_AFFINE_SLOPES")
    else:
        role_public=tuple(((2,2),) for _ in range(4))
        identity=(0,1,2);swap=(1,0,2)
        left=(identity,)*4;right=(swap,identity,identity,identity)
        start=0
        compatible=[r for r in k.candidates(domain) if k.operator_apply(domain,r,2)==2]
        need(len(compatible)==2,"TWO_INDISCERNIBLE_S3_PERMUTATIONS")
    # Role zero once, with nine identity-role steps: frozen chain length ten.
    chain=((0,False),)+tuple((1,False) for _ in range(9))
    need(len(chain)==10,"REAL_TEN_HOP_COMPOSITION")
    for j,pair in enumerate(role_public):
        need(all(k.operator_apply(domain,left[j],x)==y
                 and k.operator_apply(domain,right[j],x)==y for x,y in pair),
             "IDENTICAL_PUBLIC_DATA_FOR_DIFFERENT_TRUTH")
    a=compose(k,domain,left,chain,start)
    b=compose(k,domain,right,chain,start)
    need(a!=b,"DISTINCT_CAUSAL_OUTPUTS_FROM_IDENTICAL_PUBLIC")
    return {"domain":"F11 affine" if domain==0 else "S3 noncommuting",
       "public_information_identical":True,
       "calibration_pairs_per_role":1,
       "causal_chain_length":10,
       "number_of_compatible_role_zero_operators":len(compatible),
       "baseline_correct_world_answer":a,
       "alternative_correct_world_answer":b,
       "some_other_world_regresses_if_alternative_selected":True,
       "no_private_oracle_available":True}
def main():
    d=json.loads(FROZEN.read_text())
    need(d["schema"]=="northstar.t8.qualified-causal-transfer-source.v1"
         and d["parent_exact_head"]=="5e3b08a9a34724d5bcc902d01778c6d2e84110bf"
         and d["authorization"]["prospective_source_first_use"] is False,
         "SOURCE_PROOF_PREREG")
    k=load()
    cases=[ambiguity_witness(k,k.AFFINE),
           ambiguity_witness(k,k.NONCOMMUTING_PERM)]
    return {"schema":"wingless.t8.ambiguity-impossibility-nonconsuming.v1",
      "constructive_counterexamples":cases,
      "universally_zero_baseline_regressions_with_one_pair_and_nonzero_override":
          False,
      "required_next_scientific_change":"new publicly distinguishing calibration information or evidence-backed selective abstention with different coverage metric",
      "original_T6_T7_prospective_cases_born":0,
      "T8_new_prospective_first_used":False,
      "trained_public_prior_retuned":False,
      "scientific_acceptance":False}
if __name__=="__main__":print(json.dumps(main(),sort_keys=True))
