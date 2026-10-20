"""T21 source-free certificate-critical decision-DAG feasibility audit.

No TRAIN/DEV case generation, truth scoring, learned predictor, or original source.
Exhaust all public 3-role/3-position S3 templates in the frozen T21-A grid.
Compare exact adaptive decisions with the STRONGEST of six static orders;
validate adaptive dynamic programming against a separately enumerated
12-policy full binary conditional-tree reference.
"""
from __future__ import annotations
import importlib.util
from functools import lru_cache
from fractions import Fraction
from itertools import permutations, product
import json
from pathlib import Path

HERE=Path(__file__).resolve().parent
PERMUTATIONS=tuple(permutations(range(3)))
BINARY_WORLDS=tuple(product((0,1),repeat=3))
WORDS=tuple(tuple((order[i],inv[i]) for i in range(3))
    for order in PERMUTATIONS for inv in product((False,True),repeat=3))
FIRST_PAIRS=tuple(tuple((r,values[r]) for r in range(3))
    for values in product(range(3),repeat=3))

def exact_certifier():
    path=HERE/"20261019-wingless-t16b-persistent-bindings.py"
    spec=importlib.util.spec_from_file_location("t21_exact_certifier",path)
    m=importlib.util.module_from_spec(spec)
    spec.loader.exec_module(m)
    return m

def get_world_outputs(cert,word,first_pairs,start):
    choices=tuple(cert.consistent_ops(1,(first_pairs[r],)) for r in range(3))
    assert all(len(ops)==2 for ops in choices), "T21_INVALID_PUBLIC_S3_PAIR"
    outputs=[]
    for bits in BINARY_WORLDS:
        val=start
        for role,inverse in word:
            val=cert.apply_op(1,choices[role][bits[role]],val,inverse)
        outputs.append(val)
    return tuple(outputs)

def exact_eligible(obs):
    return tuple(i for i,world in enumerate(BINARY_WORLDS)
                 if all(obs[j]<0 or world[j]==obs[j] for j in range(3)))

def certified(outputs,obs):
    eligible=exact_eligible(obs)
    assert eligible
    return len({outputs[i] for i in eligible})==1

def expected_adaptive(outputs):
    """Exact public-observation decision DP; non-observed bits stay unknown."""
    @lru_cache(None)
    def rec(obs):
        if certified(outputs,obs):
            return Fraction(0),None
        possible=exact_eligible(obs)
        ranked=[]
        for role in range(3):
            if obs[role]>=0:
                continue
            score=Fraction(1)
            for b in (0,1):
                nxt=tuple(b if j==role else v for j,v in enumerate(obs))
                branch_size=sum(BINARY_WORLDS[i][role]==b for i in possible)
                assert branch_size>0
                score+=Fraction(branch_size,len(possible))*rec(nxt)[0]
            ranked.append((score,role))
        assert ranked
        return min(ranked)
    return rec((-1,-1,-1))[0]

def policy_cost(outputs,world,chooser):
    """Reference simulator: acquire a real role bit, then certify by every
    surviving globally consistent world. Never use this world's output as
    a certificate by itself.
    """
    obs=[-1,-1,-1];spent=0
    while not certified(outputs,tuple(obs)):
        role=chooser(tuple(obs))
        assert role in range(3) and obs[role]<0
        obs[role]=world[role]
        spent+=1
        assert spent<=3
    return spent

def best_static(outputs):
    outcomes=[]
    for ordering in PERMUTATIONS:
        def chooser(obs):
            return next(r for r in ordering if obs[r]<0)
        total=sum(policy_cost(outputs,world,chooser) for world in BINARY_WORLDS)
        outcomes.append((Fraction(total,len(BINARY_WORLDS)),ordering))
    return min(outcomes)

def independent_tree_reference(outputs):
    """12 independently enumerated adaptive trees, all one-bit branches:
    first of three roles; choose one of two remaining separately in the
    branch for first observation=0 and first observation=1.
    """
    candidates=[]
    for first in range(3):
        remaining=tuple(r for r in range(3) if r!=first)
        for on_zero in remaining:
            for on_one in remaining:
                def chooser(obs):
                    if obs[first]<0:
                        return first
                    second=on_zero if obs[first]==0 else on_one
                    if obs[second]<0:
                        return second
                    return next(r for r in remaining if obs[r]<0)
                total=sum(policy_cost(outputs,world,chooser)
                          for world in BINARY_WORLDS)
                candidates.append((Fraction(total,len(BINARY_WORLDS)),
                                   (first,on_zero,on_one)))
    assert len(candidates)==12
    return min(candidates)

def role_rename(word,first_pairs,renaming):
    new_pairs=[None]*3
    for old,new in enumerate(renaming):
        new_pairs[new]=first_pairs[old]
    return tuple((renaming[role],inverse) for role,inverse in word),tuple(new_pairs)

def run():
    cert=exact_certifier()
    assert len(WORDS)==48 and len(FIRST_PAIRS)==27
    compared=0;strict=0;ties=0;gap_sum=Fraction(0)
    max_gap=Fraction(0);best_witness=None
    equivariance_checks=0
    predeclared_example=((0,True),(1,False),(2,False))
    exemplar=None
    for word in WORDS:
        for first in FIRST_PAIRS:
            for start in range(3):
                outputs=get_world_outputs(cert,word,first,start)
                optimal=expected_adaptive(outputs)
                reference,ref_tree=independent_tree_reference(outputs)
                strongest_static,static_order=best_static(outputs)
                assert optimal==reference, "T21_EXACT_TREE_DP_DISAGREE"
                assert optimal<=strongest_static, "T21_ADAPTIVE_WORSE_THAN_STATIC"
                gap=strongest_static-optimal
                compared+=1
                gap_sum+=gap
                if gap>0:
                    strict+=1
                else:
                    ties+=1
                if gap>max_gap:
                    max_gap=gap
                    best_witness={"word":list(word),"start":start,
                      "first_public_pairs":list(first),
                      "adaptive_expected_queries":str(optimal),
                      "strongest_static_expected_queries":str(strongest_static),
                      "strongest_static_order":list(static_order),
                      "adaptive_example_tree":list(ref_tree)}
                if word==predeclared_example and first==((0,1),(1,0),(2,0)) and start==0:
                    exemplar={"adaptive_expected_queries":str(optimal),
                              "strongest_static_expected_queries":str(strongest_static)}
                    assert optimal==Fraction(9,4) and strongest_static==Fraction(5,2)
                # Fixed deterministic prefix, no source identities or DEV cases.
                if compared<=60:
                    for renaming in PERMUTATIONS:
                        renamed_word,renamed_first=role_rename(word,first,renaming)
                        shifted=get_world_outputs(cert,renamed_word,renamed_first,start)
                        assert expected_adaptive(shifted)==optimal, "T21_ROLE_RENAME_ADAPTIVE"
                        assert best_static(shifted)[0]==strongest_static, "T21_ROLE_RENAME_STATIC"
                        equivariance_checks+=1
    assert compared==3888 and strict>0 and exemplar is not None
    assert strict+ties==compared and equivariance_checks==360
    return {"schema":"wingless.t21.source-free-adaptive-query-feasibility.v1",
      "classification":"SOURCE_FREE_MATHEMATICAL_FEASIBILITY_ONLY",
      "public_templates":compared,"words":48,"public_first_pair_patterns":27,
      "starts":3,"worlds_per_template":8,
      "adaptive_vs_12_independent_conditional_trees":"EXACT_MATCH_ALL",
      "strict_adaptive_advantages":strict,"static_adaptive_ties":ties,
      "max_exact_expected_query_gain":str(max_gap),
      "aggregate_expected_query_gain":str(gap_sum),
      "public_witness":best_witness,
      "predeclared_example":exemplar,
      "role_renaming_invariance_checks":equivariance_checks,
      "TRAIN_used":0,"DEV_used":0,"source_consumed":False,
      "T20_DEV_replayed":False,"T16_originals_opened":0,
      "independent_science":False,"learned_cognition_proven":False,
      "new_T21_source_admitted":False,"optimizer_admission":False,
      "scientific_acceptance":False,"merge_authorized":False}

if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True))
