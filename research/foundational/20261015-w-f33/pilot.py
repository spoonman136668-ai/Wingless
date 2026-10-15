"""W-F33 latent role grammar feasibility — train-only pilot, no primary original.

Each episode has an independently permuted mapping from six public role tokens
to unknown 3-bit transition operators. Demonstrations disclose only
(start, token, observed end); true maps remain private to source generator.
The learned model optimizes 18 local weights by gradient, then composes
new sequences (length 3-5). The strong empirical posterior sees same demos.
"""
from __future__ import annotations
import json,math,random
from collections import defaultdict

TRAIN_SEEDS=range(191000,191040)
RESERVED_TEST_SEEDS=range(192000,192012)
TRAIN_FAMILIES=(160,161,162,163)
RESERVED_TEST_FAMILIES=(164,165,166,167)
ROLES=6;BITS=3;DEMOS=5
DEPTHS=(3,4,5)
ARMS=("learned","empirical_bayes","nearest_demo","no_role","oracle")
SEED_BASIS=9127341

def source(seed,family):
    g=random.Random(SEED_BASIS+seed*12317+family*937)
    # Hidden semantics of tokens. A role is a nontrivial bit-flip operator.
    masks=[g.randint(1,7) for _ in range(ROLES)]
    if family%2==1: masks=[m^((m<<1)&7) or 1 for m in masks]
    demos=[]
    for role in range(ROLES):
        for i in range(DEMOS):
            start=g.randrange(8)
            corrupt=(1<<g.randrange(BITS)) if g.random()<.16 else 0
            observed=start^masks[role]^corrupt
            demos.append((role,start,observed))
    g.shuffle(demos)
    return masks,demos

def public_target(seed,family,depth,masks):
    g=random.Random(SEED_BASIS+seed*713+family*1901+depth*17011)
    first=g.randrange(8)
    roles=[g.randrange(ROLES) for _ in range(depth)]
    true_end=first
    for role in roles:true_end^=masks[role]
    return first,roles,true_end

def learnt_logits(demos):
    # Local parameter learning is from public observed demonstration transitions
    # only. The true role grammar is NOT part of this function signature.
    params=[[0. for _ in range(BITS)] for _ in range(ROLES)]
    for step in range(110):
        grad=[[0. for _ in range(BITS)] for _ in range(ROLES)]
        count=[[0 for _ in range(BITS)] for _ in range(ROLES)]
        for role,start,end in demos:
            flips=start^end
            for b in range(BITS):
                y=float((flips>>b)&1)
                p=1/(1+math.exp(-params[role][b]))
                grad[role][b]+=p-y
                count[role][b]+=1
        for r in range(ROLES):
            for b in range(BITS):
                params[r][b]-=.35*(grad[r][b]/max(1,count[r][b])+.001*params[r][b])
    return params

def decoded_masks(params):
    return [sum((1<<b) for b in range(BITS) if w[b]>=0) for w in params]

def empirical_masks(demos):
    # Strong Beta(1,1) bit posterior/mode with exactly matched demonstrations.
    votes=[[1.,1.] for _ in range(ROLES*BITS)]
    for role,start,end in demos:
        for b in range(BITS):
            idx=role*BITS+b
            votes[idx][((start^end)>>b)&1]+=1
    return [sum((1<<b) for b in range(BITS)
                if votes[r*BITS+b][1]>=votes[r*BITS+b][0]) for r in range(ROLES)]

def nearest_masks(demos):
    observed={}
    for role,start,end in demos:
        observed[role]=start^end
    return [observed.get(role,0) for role in range(ROLES)]

def predict(initial,roles,masks):
    state=initial
    for role in roles:state^=masks[role]
    return state

def summarize(rows,name):
    correct=sum(r["predictions"][name]==r["truth_end"] for r in rows)
    return dict(correct=correct,total=len(rows),accuracy=correct/len(rows))

def run():
    if set(TRAIN_SEEDS)&set(RESERVED_TEST_SEEDS) or set(TRAIN_FAMILIES)&set(RESERVED_TEST_FAMILIES):
        raise ValueError("source identity collision")
    rows=[]
    for seed in TRAIN_SEEDS:
        for family in TRAIN_FAMILIES:
            secret,demos=source(seed,family)
            learned=decoded_masks(learnt_logits(demos))
            empirical=empirical_masks(demos)
            nearest=nearest_masks(demos)
            for depth in DEPTHS:
                initial,roles,target=public_target(seed,family,depth,secret)
                prediction={
                    "learned":predict(initial,roles,learned),
                    "empirical_bayes":predict(initial,roles,empirical),
                    "nearest_demo":predict(initial,roles,nearest),
                    "no_role":initial,
                    "oracle":target
                }
                rows.append(dict(seed=seed,family=family,depth=depth,query_roles=roles,
                                 query_initial=initial,truth_end=target,
                                 predictions=prediction,public_demonstration_count=len(demos),
                                 external_demonstration_bytes=len(demos)*3*8,
                                 learned_core_bytes=ROLES*BITS*8+8,
                                 first_use_f3=False))
    scores={arm:summarize(rows,arm) for arm in ARMS}
    best=max(scores[name]["accuracy"] for name in ("empirical_bayes","nearest_demo","no_role"))
    gain=scores["learned"]["accuracy"]-best
    role_gain=scores["learned"]["accuracy"]-scores["no_role"]["accuracy"]
    eligible=(role_gain>0 and gain>=0)
    return dict(schema="wingless.f33.train-only-latent-grammar-pilot.v1",
                status="PRE_ORIGINAL_TRAIN_ONLY",parent="6b46109cfda106696b8d68e23163c32cdc364420",
                seeds=[191000,191040],reserved_seeds=[192000,192012],
                families=TRAIN_FAMILIES,reserved_families=RESERVED_TEST_FAMILIES,
                questions=len(rows),scores=scores,rows=rows,
                learned_gradient_steps_per_grammar=110,
                learned_coefficients=18,native_state_bytes=152,
                external_demonstration_bytes_per_grammar=6*5*3*8,
                no_true_role_grammar_to_learner=True,source_generated_within_repo=True,
                privileged_oracle_correct=len(rows),
                verdict="PILOT_ROLE_GRAMMAR_QUALIFIES_FOR_SEPARATE_SOURCE" if eligible else
                        "PILOT_ROLE_GRAMMAR_NO_QUALIFIED_GAIN",
                primary_original_not_executed=True,F3_source_custody=False,
                accepted_ref_mutation=False,paid_api_calls=0)

if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
