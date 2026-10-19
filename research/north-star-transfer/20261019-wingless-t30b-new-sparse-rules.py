"""T30B NEW sparse-evidence unknown causal law source, independent from any learner.

Hidden function is a previously unspecified Boolean truth table shared across
all nodes of a sampled directed graph. Public evidence consists of five NEW\nreset-style causal intervention experiments with full initial state and only\ntwo observed poststate entity bits. Source origin is internal, not external.

TRAIN Boolean mechanisms are affine; DEV mechanism functions are genuinely
nonaffine, never used in TRAIN. Both are drawn from the entire 2^2/2^3
truth-table universe. Model sees ONLY public experiments and query.
"""
from __future__ import annotations
import random
FAMILIES={"M":(7,2),"N":(9,3)}
TRAIN=tuple(f"t30b-{f}-train-{i:03d}" for f in ("M","N") for i in range(20))
DEV=tuple(f"t30b-{f}-dev-{i:03d}" for f in ("M","N") for i in range(32))
ALL=TRAIN+DEV
def templates(family):
    n,k=FAMILIES[family]
    gen=random.Random(309123 if family=="M" else 309277)
    found=[];used=set()
    while len(found)<8:
        permutation=list(range(n));gen.shuffle(permutation)
        parents=[]
        for node in range(n):
            pool=[v for v in range(n) if v!=node]
            if family=="M" and node>=3:
                pool=[v for v in range(node)]
            picks=gen.sample(pool,k)
            parents.append(tuple(permutation[p] for p in picks))
        identity=tuple(parents)
        if identity not in used:used.add(identity);found.append(identity)
    return tuple(found)
TEMPLATES={f:templates(f) for f in FAMILIES}
assert all(len(s)==8 and len(set(s))==8 for s in TEMPLATES.values())
def affine_rules(arity):
    rules=set()
    for params in range(1<<(arity+1)):
        tt=0
        for index in range(1<<arity):
            val=params&1
            for p in range(arity):
                if (params>>(p+1))&1:val^=(index>>p)&1
            tt|=val<<index
        rules.add(tt)
    assert len(rules)==1<<(arity+1)
    return tuple(sorted(rules))
AFFINE={k:affine_rules(k) for k in (2,3)}
def update(family,pre,topology,truth_table,intervention=None):
    n,k=FAMILIES[family]
    after=0
    for j,parent in enumerate(topology):
        pattern=sum(((pre>>p)&1)<<a for a,p in enumerate(parent))
        bit=(truth_table>>pattern)&1
        after|=(((pre>>j)&1)^bit)<<j
    if intervention is not None:
        node,val=intervention
        after=(after&~(1<<node))|(val<<node)
    return after
def build(identity):
    if identity not in ALL:raise ValueError("T30B_ID_NOT_FROZEN")
    _,family,split,suffix=identity.split("-")
    n,k=FAMILIES[family]
    rng=random.Random(30977701+(0 if family=="M" else 607113)+
                      (0 if split=="train" else 506311)+int(suffix)*8117)
    graph=rng.randrange(8)
    if split=="train":
        law=rng.choice(AFFINE[k])
    else:
        law=rng.randrange(1<<(1<<k))
        while law in AFFINE[k]:law=rng.randrange(1<<(1<<k))
    assert (law in AFFINE[k]) == (split=="train")
    topology=TEMPLATES[family][graph]
    probes=[]
    for i in range(5):
        pre=rng.randrange(1<<n)
        action=(rng.randrange(n),rng.randrange(2)) if i%3==1 else None
        post=update(family,pre,topology,law,action)
        observed=tuple(sorted(rng.sample(range(n),2)))
        probe={"pre":pre,"intervention":action,
               "observed":tuple((j,(post>>j)&1) for j in observed)}
        probes.append(probe)
    start=rng.randrange(1<<n)
    do=(rng.randrange(n),rng.randrange(2))
    targets=tuple(rng.sample(range(n),3))
    after=update(family,start,topology,law,do)
    final=update(family,after,topology,law)
    answers=tuple((final>>j)&1 for j in targets)
    return {"id":identity,"family":family,
        "public":{"family":family,"probes":tuple(probes),
                  "query":{"start":start,"do":do,"targets":targets}},
        "_private_template_index":graph,"_private_law":law,
        "_private_truth":answers,
        "_origin":"INTERNAL_NEW_UNKNOWN_CAUSAL_MECHANISM"}
def public_train(item):
    if item["id"] not in TRAIN:raise PermissionError("T30B_NOT_TRAIN")
    return item["public"]
def public_DEV(item):
    if item["id"] not in DEV:raise PermissionError("T30B_NOT_DEV")
    return item["public"]
