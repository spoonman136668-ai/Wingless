"""W-F34 gradient-learned ordered role transition pilot: TRAIN ONLY.

New TRAIN-only episodes (203000..203011) are disjoint from the prior 201000
source headroom pilot and reserved 202000 primary originals. A role-conditioned
softmax matrix is optimized on noisy public demonstrations, then projected to
a 64-byte set of discrete bijections for inference. Strong permutation MAP
and affine/XOR operator MAP use exactly the same public demonstrations.

This is authored synthetic development; no independent F3 custody, no primary
original, no claim of actual Python heap fitting inside 576 bytes.
"""
from __future__ import annotations
import hashlib
import json
import math
import random
from functools import lru_cache

SEEDS=range(203000,203012)
FAMILIES=(188,189,190)
DEPTHS=(6,8,10)
RESERVED_SEEDS=range(202000,202016)
RESERVED_FAMILIES=(184,185,186,187)
PREVIOUS_TRAIN_SEEDS=range(201000,201024)
PREVIOUS_TRAIN_FAMILIES=(180,181,182)
ROLES=8
STATES=8
PUBLIC_DEMOS_PER_ROLE=4
EPOCHS=80
MODEL_COEFFICIENTS=ROLES*STATES*STATES
MODEL_DISTILLED_BYTES=ROLES*STATES
CORE_BUDGET=576
ARMS=("learned","permutation_map","affine_map","nearest_demo","no_role","privileged_oracle")

def make_rng(seed,family,episode):
    if seed not in SEEDS or family not in FAMILIES:
        raise ValueError("source identity unauthorized; reserved original and consumed pilot denied")
    return random.Random(839117+seed*27197+family*809+episode*1009)

def private_grammar(seed,family):
    r=make_rng(seed,family,19)
    operations=[]
    for role in range(ROLES):
        if family==188:
            a=r.choice((1,3,5,7));b=r.randrange(8)
            ops=tuple((a*x+b)%8 for x in range(STATES))
        elif family==189:
            vals=list(range(STATES));r.shuffle(vals);ops=tuple(vals)
        else:
            a=r.choice((1,3,5,7));b=r.randrange(8);mask=r.randrange(1,8)
            ops=tuple(((a*x+b)%8)^mask for x in range(STATES))
        if sorted(ops)!=list(range(STATES)):raise ValueError("invalid bijective grammar")
        operations.append(ops)
    return tuple(operations)

def public_episode(seed,family,depth):
    operations=private_grammar(seed,family)
    r=make_rng(seed,family,depth+300)
    demos=[]
    for role in range(ROLES):
        for _ in range(PUBLIC_DEMOS_PER_ROLE):
            first=r.randrange(STATES)
            second=operations[role][first]
            if r.random()<.12: second=(second+r.randint(1,7))%STATES
            demos.append((role,first,second))
    r.shuffle(demos)
    initial=r.randrange(STATES)
    roles=tuple(r.randrange(ROLES) for _ in range(depth))
    return operations,tuple(demos),initial,roles

def apply(ops,initial,roles):
    state=initial
    for role in roles:state=ops[role][state]
    return state

def independent_apply(ops,initial,roles):
    composed=list(range(STATES))
    for role in roles:
        lookup=ops[role]
        composed=[lookup[x] for x in composed]
    return composed[initial]

def assignment(scores):
    """Exact max-weight assignment (row x output), cached DP across 2**8 masks."""
    @lru_cache(None)
    def best(row,mask):
        if row==STATES:return 0.0,()
        options=[]
        for out in range(STATES):
            if mask&(1<<out):continue
            sub,path=best(row+1,mask|(1<<out))
            options.append((scores[row][out]+sub,(out,)+path))
        return max(options,key=lambda x:(x[0],tuple(-v for v in x[1])))
    return best(0,0)[1]

def fitted_gradient(demos):
    # SGD learns parameters from observations; the MAP baseline below instead
    # uses closed-form counts and a different permutation assignment.
    logits=[[[0.0]*STATES for _ in range(STATES)] for _ in range(ROLES)]
    for epoch in range(EPOCHS):
        lr=.18/(1+epoch/80)
        for role,x,y in demos:
            row=logits[role][x]
            maxima=max(row)
            probs=[math.exp(v-maxima) for v in row]
            z=sum(probs)
            for out in range(STATES):
                gradient=probs[out]/z-(out==y)
                row[out]-=lr*(gradient+.0005*row[out])
        # Smooth rows to share a weak orthogonality prior: suppress reuse of
        # already-popular output nodes in every role, without true grammar.
        if epoch%4==0:
            for role in range(ROLES):
                for out in range(STATES):
                    soft_column=sum(math.exp(max(-20,min(20,logits[role][x][out]))) for x in range(STATES))
                    correction=.002*(soft_column-STATES)
                    for x in range(STATES):
                        logits[role][x][out]-=correction
    mapped=tuple(tuple(assignment(logits[role])) for role in range(ROLES))
    if any(sorted(op)!=list(range(STATES)) for op in mapped):
        raise ValueError("learned map violates permutation requirement")
    return mapped

def fitted_permutation_map(demos):
    counts=[[[0.0]*STATES for _ in range(STATES)] for _ in range(ROLES)]
    for role,x,y in demos:counts[role][x][y]+=1.
    return tuple(tuple(assignment(counts[role])) for role in range(ROLES))

def affine_candidates():
    result=set()
    for a in (1,3,5,7):
        for b in range(STATES):
            result.add(tuple((a*x+b)%STATES for x in range(STATES)))
            for mask in range(1,STATES):
                result.add(tuple(((a*x+b)%STATES)^mask for x in range(STATES)))
    return tuple(sorted(result))
CANDIDATES=affine_candidates()

def fitted_affine_map(demos):
    roles=[]
    for role in range(ROLES):
        samples=[(x,y) for r,x,y in demos if r==role]
        best=max(CANDIDATES,key=lambda op:(sum(op[x]==y for x,y in samples),tuple(-x for x in op)))
        roles.append(best)
    return tuple(roles)

def fitted_nearest(demos):
    mapping=[]
    for role in range(ROLES):
        observations={(x,y) for r,x,y in demos if r==role}
        fit=[]
        for x in range(STATES):
            vals=sorted(y for xx,y in observations if xx==x)
            fit.append(vals[-1] if vals else x)
        mapping.append(tuple(fit))
    return tuple(mapping)

def scores(rows,name):
    correct=sum(r["predictions"][name]==r["truth_end"] for r in rows)
    return dict(correct=correct,total=len(rows),accuracy=correct/len(rows))

def train_only_pilot():
    if (set(SEEDS)&set(RESERVED_SEEDS) or set(FAMILIES)&set(RESERVED_FAMILIES)
        or set(SEEDS)&set(PREVIOUS_TRAIN_SEEDS) or set(FAMILIES)&set(PREVIOUS_TRAIN_FAMILIES)):
        raise ValueError("train/test/previous pilot identity collision")
    rows=[]
    for seed in SEEDS:
        for family in FAMILIES:
            for depth in DEPTHS:
                secret,demos,initial,roles=public_episode(seed,family,depth)
                truth=apply(secret,initial,roles)
                if truth!=independent_apply(secret,initial,roles):raise ValueError("oracle mismatch")
                fitted={
                    "learned":fitted_gradient(demos),
                    "permutation_map":fitted_permutation_map(demos),
                    "affine_map":fitted_affine_map(demos),
                    "nearest_demo":fitted_nearest(demos),
                }
                answers={name:apply(ops,initial,roles) for name,ops in fitted.items()}
                answers.update(no_role=initial,privileged_oracle=truth)
                evidence=hashlib.sha256(json.dumps([demos,initial,roles]).encode()).hexdigest()
                rows.append(dict(seed=seed,family=family,depth=depth,truth_end=truth,
                      predictions=answers,public_evidence_hash=evidence,
                      ordered_role_query=roles,reverse_order_changes_truth=truth!=independent_apply(secret,initial,roles[::-1]),
                      distilled_learned_bytes=MODEL_DISTILLED_BYTES,transient_gradient_parameters=MODEL_COEFFICIENTS,
                      outside_core_demo_bytes=len(demos)*3*8))
    summary={name:scores(rows,name) for name in ARMS}
    by_family={str(f):{name:scores([r for r in rows if r["family"]==f],name) for name in ARMS} for f in FAMILIES}
    highest=max(summary[a]["accuracy"] for a in ARMS[1:-1])
    gain=summary["learned"]["accuracy"]-highest
    condition=all(by_family[str(f)]["learned"]["accuracy"]>=max(by_family[str(f)][a]["accuracy"] for a in ARMS[1:-1]) for f in FAMILIES)
    verdict=("TRAIN_ONLY_LEARNED_COMPOSITION_ELIGIBLE_FOR_SEPARATE_ORIGINAL" if gain>=.03 and condition else
             "TRAIN_ONLY_LEARNED_COMPOSITION_NO_QUALIFIED_GAIN")
    return dict(schema="wingless.f34.gradient-transition-train-only.v1",
        stage="TRAIN_ONLY_PROSPECTIVE_POLICY_QUALIFICATION",questions=len(rows),
        train_seed_range=[SEEDS.start,SEEDS.stop],train_families=FAMILIES,depths=DEPTHS,
        primary_original_seed_range=[RESERVED_SEEDS.start,RESERVED_SEEDS.stop],
        primary_original_families=RESERVED_FAMILIES,
        output_mapping_bytes=MODEL_DISTILLED_BYTES,native_core_budget_bytes=CORE_BUDGET,
        transient_training_parameters=MODEL_COEFFICIENTS,
        optimizer_scratch_explicitly_outside_native_core=True,
        demonstration_history_bytes_per_episode=ROLES*PUBLIC_DEMOS_PER_ROLE*3*8,
        source_synthetic_from_same_repository=True,independent_F3_custody=False,
        primary_original_not_executed=True, accepted_ref_mutation=False,
        scores=summary,per_family=by_family,rows=rows,verdict=verdict)
if __name__=="__main__":print(json.dumps(train_only_pilot(),sort_keys=True,indent=2))
