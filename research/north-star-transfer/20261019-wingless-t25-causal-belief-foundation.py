"""T25 exact partial-observation causal-belief foundation, no learned claims.

A node's latent rule is one of 2*choose(7,2)=42 AND/OR parent
configurations. Exact belief state retains EVERY rule consistent with
observed one-step interventions and mask. Compositional counterfactual
propagation uses sound interval abstraction; losing correlations may lead
to abstention, never a guessed hard-positive.

This module never reads old T16-T24 source or DEV. The TRAIN-only test
uses source's hidden graph solely as an independent truth witness.
"""
from __future__ import annotations
import itertools,json,statistics,sys
from pathlib import Path
N=8
RULES=tuple((a,b,g)
            for a,b in itertools.combinations(range(N),2)
            for g in (0,1))
# candidate rule excludes the target itself. 21*2=42 each node.
CANDIDATES=tuple(tuple(rule for rule in RULES if rule[0]!=target and rule[1]!=target)
                 for target in range(N))
assert all(len(x)==42 for x in CANDIDATES)
def step_bit(current,parents,gate):
    a,b=parents
    return current[0]^(current[a]&current[b] if gate==0 else current[a]|current[b])

def consistent(target,rule,probe):
    if probe["do"] is not None and probe["do"][0]==target:return True
    observed=dict(probe["observed"])
    if target not in observed:return True
    a,b,gate=rule
    before=probe["pre"]
    predicted=before[target]^(before[a]&before[b] if gate==0 else before[a]|before[b])
    return predicted==observed[target]

def posterior(public):
    observed=public["probes"]
    assert len(observed)==18
    beliefs=[]
    for target in range(N):
        eligible=tuple(rule for rule in CANDIDATES[target]
                       if all(consistent(target,rule,probe) for probe in observed))
        if not eligible:raise RuntimeError("T25_EMPTY_CAUSAL_POSTERIOR")
        beliefs.append(eligible)
    return tuple(beliefs)

def advance(bits,belief,intervention):
    nxt=[]
    for target in range(N):
        if intervention is not None and intervention[0]==target:
            nxt.append(frozenset((intervention[1],)))
            continue
        values=set()
        for a,b,g in belief[target]:
            for x in bits[target]:
                for lhs in bits[a]:
                    for rhs in bits[b]:
                        values.add(x^(lhs&rhs if g==0 else lhs|rhs))
        if not values:raise RuntimeError("T25_EMPTY_ABSTRACT_REACHABILITY")
        nxt.append(frozenset(values))
    return tuple(nxt)

def infer(beliefs,query):
    state=tuple(frozenset((x,)) for x in query["pre"])
    assert len(state)==N and len(query["do"])==2
    for intervention in query["do"]:
        state=advance(state,beliefs,intervention)
    possible=state[query["target"]]
    return next(iter(possible)) if len(possible)==1 else None,possible

def train_only(source):
    assert len(source.TRAIN)==40 and len(source.DEV)==64
    assert set(source.TRAIN).isdisjoint(source.DEV)
    posterior_checks=0;ambiguous_nodes=0;worlds_with_uncertainty=0
    question_count=0;certified=0;abstained=0;errors=0
    posterior_sizes=[]
    for identity in source.TRAIN:
        world=source.build(identity)
        public=source.train_public(world)
        belief=posterior(public)
        uncertain=0
        for j,hypotheses in enumerate(belief):
            true=(min(world["_private_rules"][j][0],world["_private_rules"][j][1]),
                  max(world["_private_rules"][j][0],world["_private_rules"][j][1]),
                  world["_private_rules"][j][2])
            if true not in hypotheses:
                raise RuntimeError("T25_TRUE_WORLD_ELIMINATED")
            if len(hypotheses)>1:uncertain+=1;ambiguous_nodes+=1
            posterior_checks+=len(CANDIDATES[j])*len(public["probes"])
            posterior_sizes.append(len(hypotheses))
        if uncertain:worlds_with_uncertainty+=1
        for query,truth in zip(public["queries"],world["_private_truth"]):
            val,possible=infer(belief,query)
            if truth not in possible:
                raise RuntimeError("T25_FALSE_NEGATIVE_EXCLUDED_GROUND_TRUTH")
            question_count+=1
            if val is None:abstained+=1
            else:
                certified+=1
                if val!=truth:errors+=1
    if errors:raise RuntimeError("T25_HARD_FALSE_POSITIVE")
    return {
        "schema":"wingless.T25.TRAIN-only-causal-world-headroom.v1",
        "status":"PASS_TRAIN_ONLY_DEV_UNOPENED",
        "source_families":2,
        "TRAIN_cases":40,"DEV_generated":0,
        "entities":8,"nonlinear_rule_candidates_each_node":42,
        "public_probes_per_world":18,"masked_observed_outputs_per_probe":6,
        "hidden_rule_contained_every_node":True,
        "candidate_probe_consistency_checks":posterior_checks,
        "ambiguous_target_nodes":ambiguous_nodes,
        "worlds_with_unresolved_causal_edges":worlds_with_uncertainty,
        "posterior_mean_size":statistics.mean(posterior_sizes),
        "posterior_max_size":max(posterior_sizes),
        "counterfactual_questions":question_count,
        "soundly_certified":certified,
        "abstained_not_guessed":abstained,
        "hard_false_positives":errors,
        "information_headroom_fraction":abstained/question_count,
        "source_origin":"INTERNAL_NOT_OUTSIDE_REPLICATED",
        "learned_candidate_trained":False,
        "scientific_acceptance":False,"optimizer_admission":False,
        "prior_consumed_sources_opened":0
    }
