"""T33 FRESH causal-family source: learned predictive dual-syndrome witness.

Unlike T32 pair/triple constraint closure, public observations are
independent reset interventions on a hidden unknown 4-/5-input causal law,
and future queries ask parity after one unknown-law transition followed by
an independent, known deterministic transport across entities.

Unknown law, hidden graph, and future truth are never returned in public.
TRAIN laws have bounded polynomial degree; DEV laws are strictly outside
that class. Distinct source IDs and seeds from prior Wingless tranches.
"""
from __future__ import annotations
import random
FAMILIES={"U":(8,4,14),"V":(10,5,18)}
TRAIN=tuple(f"t33-{f}-train-{i:03d}" for f in ("U","V") for i in range(20))
DEV=tuple(f"t33-{f}-dev-{i:03d}" for f in ("U","V") for i in range(32))
ALL=TRAIN+DEV
def graph_templates(f):
    n,k,_=FAMILIES[f];r=random.Random(3300441+(0 if f=="U" else 991))
    unique=set();out=[]
    while len(out)<3:
        parents=tuple(tuple(sorted(r.sample([i for i in range(n) if i!=j],k))) for j in range(n))
        if parents not in unique:unique.add(parents);out.append(parents)
    return tuple(out)
TEMPLATES={f:graph_templates(f) for f in FAMILIES}
def anf_degree(truth,k):
    a=[(truth>>j)&1 for j in range(1<<k)]
    for i in range(k):
        for j in range(1<<k):
            if j&(1<<i):a[j]^=a[j^(1<<i)]
    return max((j.bit_count() for j,v in enumerate(a) if v),default=0)
def law_from_coeff(coeff,k):
    law=0
    for p in range(1<<k):
        v=0
        for q in range(1<<k):
            if (q&p)==q:v^=(coeff>>q)&1
        law|=v<<p
    return law
def choose_law(rng,f,split):
    n,k,_=FAMILIES[f]
    # Restricted low-degree source-born TRAIN and withheld higher-degree
    # DEV functions. The model is never told which subset generated a case.
    degree=1 if f=="U" else 2
    eligible=[i for i in range(1<<k) if i.bit_count()<=degree]
    if split=="train":
        coeff=sum(rng.randrange(2)<<i for i in eligible)
    else:
        outside=[i for i in range(1<<k) if i.bit_count()>degree]
        coeff=sum(rng.randrange(2)<<i for i in eligible)
        coeff|=1<<rng.choice(outside)
        for i in outside:
            if rng.randrange(4)==0:coeff^=1<<i
        if not any(coeff&(1<<i) for i in outside):coeff|=1<<outside[0]
    law=law_from_coeff(coeff,k)
    assert (anf_degree(law,k)<=degree)==(split=="train")
    return law
def input_pattern(pre,parents):
    return sum(((pre>>node)&1)<<i for i,node in enumerate(parents))
def transition(f,pre,topo,truth):
    n,k,_=FAMILIES[f]
    after=0
    for j,parents in enumerate(topo):
        value=((pre>>j)&1)^((truth>>input_pattern(pre,parents))&1)
        after|=value<<j
    return after
def transport(f,bits):
    n,_,_=FAMILIES[f]
    # known bijective one-step entity relation: new node j <- old j-2
    return sum(((bits>>((j-2)%n))&1)<<j for j in range(n))
def parity(bits,nodes):
    return sum((bits>>j)&1 for j in nodes)%2
def sample(identity):
    if identity not in ALL:raise ValueError("T33_UNFROZEN_SOURCE_ID")
    _,f,split,num=identity.split("-");n,k,count=FAMILIES[f]
    r=random.Random(33095503+(0 if f=="U" else 991919)+(0 if split=="train" else 177777)+int(num)*9733)
    graph=r.randrange(3);law=choose_law(r,f,split);topo=TEMPLATES[f][graph]
    pool=[r.randrange(1<<n) for _ in range(10)]
    probe=[]
    for j in range(count):
        pre=r.choice(pool)
        nodes=tuple(sorted(r.sample(range(n),2 if j%3 else 3)))
        after=transition(f,pre,topo,law)
        probe.append({"pre":pre,"nodes":nodes,"parity":parity(after,nodes)})
    query=[]
    for i in range(3):
        pre=r.choice(pool)
        targets=tuple(sorted(r.sample(range(n),2 if i!=1 else 3)))
        value=transport(f,transition(f,pre,topo,law))
        query.append({"pre":pre,"targets":targets,"_truth":parity(value,targets)})
    public_query=tuple({"pre":x["pre"],"targets":x["targets"]} for x in query)
    return {"id":identity,"family":f,
      "public":{"family":f,"probes":tuple(probe),"queries":public_query},
      "_hidden_graph":graph,"_hidden_law":law,
      "_future_truth":tuple(q["_truth"] for q in query)}
def public_train(case):
    if case["id"] not in TRAIN:raise PermissionError("T33_TRAIN_ONLY")
    return case["public"]
def public_dev(case):
    if case["id"] not in DEV:raise PermissionError("T33_DEV_ONLY")
    return case["public"]
