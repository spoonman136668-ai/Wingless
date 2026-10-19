"""T31 fresh INTERNAL unknown conditional causal programs. No old DEV reused.

Hidden causal world: 8 new graph templates, unknown global 2/3-input Boolean
truth table, unknown coupling mode. Full prestate and two poststate sensors
for each of five distinct public experiments. Future is 3-step, not T30B's 2.
TRAIN affine laws; unopened DEV strictly nonaffine. Private labels and truth
are scorer-only, not fitted or supplied to solvers.
"""
from __future__ import annotations
import random

FAMILIES={"P":(7,2),"Q":(9,3)}
TRAIN=tuple(f"t31-{f}-train-{i:03d}" for f in ("P","Q") for i in range(20))
DEV=tuple(f"t31-{f}-dev-{i:03d}" for f in ("P","Q") for i in range(32))
ALL=set(TRAIN+DEV)

def template_family(f):
    n,k=FAMILIES[f]
    r=random.Random(313175 if f=="P" else 313199)
    out=[];seen=set()
    while len(out)<8:
        shift=r.randrange(1,n)
        parents=[]
        for j in range(n):
            choices=[(j-shift)%n,(j+1)%n,(j+2)%n,(j+3)%n]
            candidates=list(dict.fromkeys(v for v in choices if v!=j))
            if len(candidates)<k: raise AssertionError("T31_BAD_GRAPH")
            r.shuffle(candidates)
            parents.append(tuple(candidates[:k]))
        graph=tuple(parents)
        if graph not in seen:
            seen.add(graph);out.append(graph)
    return tuple(out)
TEMPLATES={f:template_family(f) for f in FAMILIES}

def affine_rules(k):
    result=set()
    for coeff in range(1<<(k+1)):
        value=0
        for pattern in range(1<<k):
            bit=coeff&1
            for i in range(k):
                bit^=((coeff>>(i+1))&1)*((pattern>>i)&1)
            value|=bit<<pattern
        result.add(value)
    return tuple(sorted(result))
AFFINE={k:affine_rules(k) for k in (2,3)}

def pattern(pre,parents):
    return sum(((pre>>p)&1)<<i for i,p in enumerate(parents))
def update(f,pre,graph,law,mode,intervention=None):
    n,k=FAMILIES[f]
    next_state=0
    for j,parents in enumerate(graph):
        bit=((pre>>j)&1)^((law>>pattern(pre,parents))&1)
        if j%2:bit^=mode*((pre>>((j+1)%n))&1)
        next_state|=bit<<j
    if intervention is not None:
        node,val=intervention
        next_state=(next_state&~(1<<node))|(val<<node)
    return next_state

def build(identity):
    if identity not in ALL:raise ValueError("T31_UNREGISTERED_CASE")
    _,f,split,suffix=identity.split("-")
    n,k=FAMILIES[f]
    r=random.Random(31317701+(0 if f=="P" else 700031)+
                    (0 if split=="train" else 900071)+int(suffix)*9973)
    graph_id=r.randrange(8);mode=r.randrange(2)
    if split=="train":law=r.choice(AFFINE[k])
    else:
        while True:
            law=r.randrange(1<<(1<<k))
            if law not in AFFINE[k]:break
    graph=TEMPLATES[f][graph_id]
    probes=[]
    for i in range(5):
        pre=r.randrange(1<<n)
        action=(r.randrange(n),r.randrange(2)) if i%3==1 else None
        outcome=update(f,pre,graph,law,mode,action)
        observed=tuple(sorted(r.sample(range(n),2)))
        probes.append({"pre":pre,"do":action,
                      "observed":tuple((j,(outcome>>j)&1) for j in observed)})
    start=r.randrange(1<<n)
    action=(r.randrange(n),r.randrange(2))
    targets=tuple(r.sample(range(n),3))
    final=start
    for step in range(3):
        final=update(f,final,graph,law,mode,action if step==0 else None)
    return {"id":identity,"family":f,
            "public":{"family":f,"probes":tuple(probes),
                      "query":{"start":start,"do":action,"targets":targets}},
            "_truth":tuple((final>>j)&1 for j in targets),
            "_private_graph":graph_id,"_private_mode":mode,
            "_private_law":law,
            "_origin":"T31_NEW_INTERNAL_CAUSAL_FAMILIES_NOT_EXTERNAL"}
def public_train(item):
    if item["id"] not in TRAIN:raise PermissionError("T31_NOT_TRAIN")
    return item["public"]
def public_dev(item):
    if item["id"] not in DEV:raise PermissionError("T31_NOT_DEV")
    return item["public"]
