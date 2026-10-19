"""T25 original internal partially observable nonlinear causal worlds, families P/Q.

Eight entities carry binary state and a hidden two-parent AND/OR update rule.
Interventions override selected post-transition entity. Public observations
reveal full pre-state but only a masked subset of next-state values.
The model must not import or access hidden graph rules or query truth.
Source generation does not import any learned candidate, baseline or certifier.
"""
from __future__ import annotations
import itertools, random
N=8
TRAIN=tuple(f"t25-{f}-train-{i:03d}" for f in ("P","Q") for i in range(20))
DEV=tuple(f"t25-{f}-dev-{i:03d}" for f in ("P","Q") for i in range(32))
ALL=TRAIN+DEV
def parents_source(rng,family):
    perm=list(range(N));rng.shuffle(perm)
    rules=[None]*N
    for pos,node in enumerate(perm):
        if family=="P":
            # Causal fan-in: mostly directed, some distant mediation.
            pool=perm[:pos] if pos>=2 else [x for x in range(N) if x!=node]
            if len(pool)<2:pool=[x for x in range(N) if x!=node]
            parents=rng.sample(pool,2)
        else:
            # Recurrent and long-range causal links, breaking DAG structure.
            pool=[x for x in range(N) if x!=node]
            p=rng.choice(pool)
            q=rng.choice([x for x in pool if x!=p])
            parents=(p,q)
        rules[node]=(parents[0],parents[1],rng.randrange(2))
    return tuple(rules)
def step(pre,rules,intervention):
    out=[]
    for node,(a,b,gate) in enumerate(rules):
        term=pre[a]&pre[b] if gate==0 else pre[a]|pre[b]
        value=pre[node]^term
        if intervention is not None and intervention[0]==node:
            value=intervention[1]
        out.append(value)
    return tuple(out)
def stimulus(rng):
    state=tuple(rng.randrange(2) for _ in range(N))
    if rng.randrange(5)<2:return state,(rng.randrange(N),rng.randrange(2))
    return state,None
def build(identity):
    if identity not in ALL:raise ValueError("T25_UNFROZEN_SOURCE_IDENTITY")
    _,family,partition,idx=identity.split("-")
    i=int(idx)
    rng=random.Random(250104701+(0 if family=="P" else 1800271)+
      (0 if partition=="train" else 720077)+i*2029)
    hidden=parents_source(rng,family)
    # Public probes are individual source interventions, with sparse output
    # observation and random entity renaming already applied to hidden graph.
    probes=[]
    for k in range(18):
        pre,intervention=stimulus(rng)
        post=step(pre,hidden,intervention)
        observed=tuple(sorted(rng.sample(range(N),6)))
        probes.append({"pre":pre,"do":intervention,
                "observed":tuple((j,post[j]) for j in observed)})
    # Unfamiliar two-step counterfactual plans, 3 diverse target queries.
    queries=[];truth=[]
    chosen=rng.sample(range(N),3)
    for target in chosen:
        pre=tuple(rng.randrange(2) for _ in range(N))
        dos=((rng.randrange(N),rng.randrange(2)),
              (rng.randrange(N),rng.randrange(2)))
        first=step(pre,hidden,dos[0])
        second=step(first,hidden,dos[1])
        queries.append({"pre":pre,"do":dos,"target":target})
        truth.append(second[target])
    return {"id":identity,"family":family,"partition":partition,
      "public":{"probes":tuple(probes),"queries":tuple(queries)},
      "_private_rules":hidden,"_private_truth":tuple(truth),
      "_origin":"INTERNAL_NOT_EXTERNALLY_REPLICATED"}
def train_public(case):
    if case["id"] not in TRAIN:raise ValueError("T25_TRAIN_NOT_REGISTERED")
    return case["public"]
def dev_public(case):
    if case["id"] not in DEV:raise ValueError("T25_DEV_NOT_REGISTERED")
    return case["public"]
