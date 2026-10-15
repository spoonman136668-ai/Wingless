"""W-F34 TRAIN-ONLY source feasibility, not a learned science original.

Three distinct hidden role-grammar processes; six-to-ten-step directed
queries. Model-facing evidence contains noisy demonstrations and the public
ordered query, never the private transition tables or oracle truth.
This project-authored generator is NOT independently born F3 evidence.
"""
from __future__ import annotations
import hashlib, json, random
from collections import Counter

SEEDS=range(201000,201024)
FAMILIES=(180,181,182)
DEPTHS=(6,8,10)
ROLES=8
STATES=8
DEMOS=4
RESERVED_SEEDS=range(202000,202016)
RESERVED_FAMILIES=(184,185,186,187)
ARMS=("empirical_posterior","nearest_demo","no_role","privileged_oracle")

def _rng(seed,fam,extra):
    if seed not in SEEDS or fam not in FAMILIES:
        raise ValueError("train-only source: RESERVED ORIGINAL IDENTITY DENIED")
    return random.Random(784901+seed*41047+fam*9719+extra*601)

def hidden_operators(seed,fam):
    g=_rng(seed,fam,10)
    ops=[]
    for r in range(ROLES):
        if fam==180:
            a=g.choice((1,3,5,7));b=g.randrange(8)
            op=tuple((a*x+b)%8 for x in range(8))
        elif fam==181:
            vals=list(range(8));g.shuffle(vals);op=tuple(vals)
        else:
            a=g.choice((1,3,5,7));b=g.randrange(8);mask=g.randrange(1,8)
            op=tuple(((a*x+b)%8)^mask for x in range(8))
        assert sorted(op)==list(range(8))
        ops.append(op)
    return tuple(ops)

def demos_and_query(seed,fam,depth):
    ops=hidden_operators(seed,fam)
    rng=_rng(seed,fam,depth+700)
    demos=[]
    for role in range(ROLES):
        for d in range(DEMOS):
            x=rng.randrange(8)
            correct=ops[role][x]
            y=correct if rng.random()>.12 else (correct+rng.randint(1,7))%8
            demos.append((role,x,y))
    rng.shuffle(demos)
    initial=rng.randrange(8)
    roles=tuple(rng.randrange(ROLES) for _ in range(depth))
    # Source-private tables are ONLY returned to the evaluator wrapper;
    # the model's public evidence is separately constructed below.
    return ops, tuple(demos), initial, roles

def eval_source(ops,initial,roles):
    value=initial
    for role in roles:
        value=ops[role][value]
    return value

def eval_independent(ops,initial,roles):
    # Independently composed complete-state transition rather than direct walk.
    composed=list(range(STATES))
    for token in roles:
        table=ops[token]
        composed=[table[prior] for prior in composed]
    return composed[initial]

def candidate_ops():
    # Fair programmed empirical posterior's hypothesis bank is public, not
    # conditioned on hidden generator family.
    opts=set()
    for a in (1,3,5,7):
        for b in range(8):
            opts.add(tuple((a*x+b)%8 for x in range(8)))
            for mask in range(1,8):
                opts.add(tuple(((a*x+b)%8)^mask for x in range(8)))
    return sorted(opts)
CANDIDATES=candidate_ops()

def empirical_map(demos):
    fitted=[]
    for role in range(ROLES):
        obs=[(x,y) for r,x,y in demos if r==role]
        # Equal-data empirical likelihood: 0.88 correct and 0.12/7 corruption.
        best=max(CANDIDATES,key=lambda op:(sum(op[x]==y for x,y in obs),-CANDIDATES.index(op)))
        fitted.append(best)
    return tuple(fitted)

def nearest_map(demos):
    roles=[]
    for role in range(ROLES):
        votes={}
        for r,x,y in demos:
            if r==role:votes.setdefault(x,[]).append(y)
        roles.append(tuple(Counter(votes[x]).most_common(1)[0][0] if x in votes else x for x in range(STATES)))
    return tuple(roles)

def run():
    assert not set(SEEDS)&set(RESERVED_SEEDS)
    assert not set(FAMILIES)&set(RESERVED_FAMILIES)
    rows=[]
    for seed in SEEDS:
        for fam in FAMILIES:
            for depth in DEPTHS:
                ops,demos,initial,roles=demos_and_query(seed,fam,depth)
                oracle=eval_source(ops,initial,roles)
                assert oracle==eval_independent(ops,initial,roles)
                empirical=empirical_map(demos)
                nearest=nearest_map(demos)
                predictions={
                    "empirical_posterior":eval_independent(empirical,initial,roles),
                    "nearest_demo":eval_independent(nearest,initial,roles),
                    "no_role":initial,
                    "privileged_oracle":oracle,
                }
                reversed_answer=eval_independent(ops,initial,roles[::-1])
                evidence_hash=hashlib.sha256(json.dumps((demos,initial,roles)).encode()).hexdigest()
                rows.append(dict(seed=seed,family=fam,depth=depth,query_roles=roles,
                    initial=initial,truth_end=oracle,predictions=predictions,
                    reversed_role_order_changes_truth=reversed_answer!=oracle,
                    public_evidence_sha256=evidence_hash,
                    public_demonstration_count=len(demos),
                    external_demo_bytes=len(demos)*3*8,
                    native_core_budget_bytes=576))
    scores={a:dict(correct=sum(r["predictions"][a]==r["truth_end"] for r in rows),total=len(rows))
            for a in ARMS}
    for a in scores:scores[a]["accuracy"]=scores[a]["correct"]/len(rows)
    order_changes=sum(r["reversed_role_order_changes_truth"] for r in rows)
    strongest=max(scores[a]["accuracy"] for a in ARMS[:-1])
    pass_gate=(scores["privileged_oracle"]["correct"]==len(rows)
               and strongest<=.96 and order_changes>=len(rows)//20)
    return dict(schema="wingless.f34.train-only-independent-grammar-source-pilot.v1",
        status="TRAIN_ONLY_SOURCE_QUALIFICATION_NOT_PRIMARY_ORIGINAL",
        train_seeds=[SEEDS.start,SEEDS.stop],reserved_original_seeds=[RESERVED_SEEDS.start,RESERVED_SEEDS.stop],
        train_families=FAMILIES,reserved_original_families=RESERVED_FAMILIES,
        depths=DEPTHS,questions=len(rows),scores=scores,
        role_order_changes=order_changes,role_order_change_rate=order_changes/len(rows),
        fixed_demo_count=ROLES*DEMOS,latent_generator_mechanisms=3,
        dual_oracle_implementation_agreement=True,
        source_generator_authored_in_same_repo=True,
        separately_controlled_F3_custody=False,primary_original_not_executed=True,
        resident_memory_measured=False,native_core_limit_bytes=576,
        external_demo_bytes_per_query=ROLES*DEMOS*3*8,
        verdict="TRAIN_ONLY_SOURCE_HEADROOM_ELIGIBLE" if pass_gate else "TRAIN_ONLY_SOURCE_HEADROOM_BLOCKED",
        accepted_ref_mutation=False,paid_api_calls=0,rows=rows)
if __name__=="__main__":print(json.dumps(run(),sort_keys=True,indent=2))
