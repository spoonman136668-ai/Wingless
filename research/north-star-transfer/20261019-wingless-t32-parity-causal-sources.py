"""T32 FRESH source-born parity-observed unknown causal transition functions.

Public probes reveal TWO parities of *post-transition* outputs: one XOR
of 2 nodes, one XOR of 3 nodes, never individual node outputs. Exact
unknown local truth-table law is shared by all nodes. True graph/mode/law
and future answers remain source-private.
"""
from __future__ import annotations
import random
FAMILIES={"H":(7,2),"J":(9,3)}
TRAIN=tuple(f"t32-{f}-train-{i:03d}" for f in ("H","J") for i in range(20))
DEV=tuple(f"t32-{f}-dev-{i:03d}" for f in ("H","J") for i in range(32))
ALL=TRAIN+DEV
def graph_templates(f):
    n,k=FAMILIES[f];rng=random.Random(3203041+(0 if f=="H" else 101))
    out=[];visited=set()
    while len(out)<8:
        perm=list(range(n));rng.shuffle(perm)
        par=tuple(tuple(perm[v] for v in rng.sample([x for x in range(n) if x!=j],k))
                  for j in range(n))
        if par not in visited:
            visited.add(par);out.append(par)
    return tuple(out)
TEMPLATES={f:graph_templates(f) for f in FAMILIES}
def affine_set(k):
    rules=set()
    for coeff in range(1<<(k+1)):
        table=0
        for bits in range(1<<k):
            v=coeff&1
            for b in range(k):v^=((coeff>>(b+1))&1)*((bits>>b)&1)
            table|=v<<bits
        rules.add(table)
    return tuple(sorted(rules))
AFFINE={k:affine_set(k) for k in (2,3)}
def pattern(pre,parents):
    return sum(((pre>>i)&1)<<j for j,i in enumerate(parents))
def update(f,pre,graph,mode,law,action=None):
    n,_=FAMILIES[f]
    result=0
    for j,parents in enumerate(graph):
        bit=((pre>>j)&1)^((law>>pattern(pre,parents))&1)
        if (j%2)==1:bit^=mode*((pre>>((j+1)%n))&1)
        result|=bit<<j
    if action is not None:
        node,v=action
        result=(result&~(1<<node))|(v<<node)
    return result
def sample(name):
    if name not in ALL:raise ValueError("T32_UNKNOWN_SOURCE_ID")
    _,f,split,num=name.split("-")
    n,k=FAMILIES[f]
    rng=random.Random(320999111+(0 if f=="H" else 700113)
                      +(0 if split=="train" else 107177)+int(num)*7739)
    graph_id=rng.randrange(8);mode=rng.randrange(2)
    if split=="train":law=rng.choice(AFFINE[k])
    else:
        law=rng.randrange(1<<(1<<k))
        while law in AFFINE[k]:law=rng.randrange(1<<(1<<k))
    graph=TEMPLATES[f][graph_id]
    probes=[]
    for step in range(5):
        pre=rng.randrange(1<<n)
        action=(rng.randrange(n),rng.randrange(2)) if step in (1,3) else None
        out=update(f,pre,graph,mode,law,action)
        pair=tuple(sorted(rng.sample(range(n),2)))
        triple=tuple(sorted(rng.sample(range(n),3)))
        def parity(nodes):return sum((out>>node)&1 for node in nodes)%2
        probes.append({"pre":pre,"do":action,
                       "observed_parities":((pair,parity(pair)),(triple,parity(triple)))})
    start=rng.randrange(1<<n)
    action=(rng.randrange(n),rng.randrange(2))
    targets=tuple(rng.sample(range(n),3))
    state=start
    for step in range(3):
        state=update(f,state,graph,mode,law,action if step==0 else None)
    truth=tuple((state>>node)&1 for node in targets)
    return {"id":name,
      "public":{"family":f,"probes":tuple(probes),
                "query":{"start":start,"do":action,"targets":targets}},
      "_true_graph":graph_id,"_true_mode":mode,"_true_law":law,
      "_answers":truth,"_origin":"T32_INTERNAL_FRESH_TWO_GRAMMARS_NOT_EXTERNAL"}
def train_public(x):
    if x["id"] not in TRAIN:raise PermissionError("T32_TRAIN_ONLY")
    return x["public"]
def dev_public(x):
    if x["id"] not in DEV:raise PermissionError("T32_DEV_ONLY")
    return x["public"]
