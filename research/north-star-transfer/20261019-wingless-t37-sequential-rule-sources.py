"""T37 new internally born delayed causal-rule source, independent of learner.

Each world has 7 or 11 entity bindings, publicly visible causal parents and
three earlier state snapshots. The 3-input Boolean transition law is hidden
and is shared across nodes within the world. TRAIN local rules have algebraic
degree at most two; original prospective DEV rules have degree exactly three.
DEV uses newly structured parent families, not a reroll of TRAIN graphs.
"""
from __future__ import annotations
import random
FAMILIES={"R":7,"S":11}
TRAIN_META=tuple(f"t37-{f}-train-{i:03d}" for f in ("R","S") for i in range(20))
TRAIN_FIT=tuple(f"t37-{f}-train-{i:03d}" for f in ("R","S") for i in range(20,28))
DEV=tuple(f"t37-{f}-dev-{i:03d}" for f in ("R","S") for i in range(32))
ALL=TRAIN_META+TRAIN_FIT+DEV
def degree(tt):
    coeff=[(tt>>i)&1 for i in range(8)]
    for k in range(3):
        for j in range(8):
            if j&(1<<k):coeff[j]^=coeff[j^(1<<k)]
    return max((j.bit_count() for j,v in enumerate(coeff) if v),default=0)
LOW=tuple(t for t in range(256) if degree(t)<=2)
HIGH=tuple(t for t in range(256) if degree(t)==3)
assert len(LOW)==128 and len(HIGH)==128
def graph(f,split,r):
    n=FAMILIES[f]
    result=[]
    for j in range(n):
        if f=="R" and split=="train":
            past=list(range(max(0,j-4),j))
            a=r.choice(past) if past else (j+1)%n
            b=r.choice(past) if past else (j+2)%n
            c=r.choice(past) if past else (j+3)%n
        elif f=="S" and split=="train":
            a=(j-1)%n;b=(j+2)%n;c=(j-3)%n
            if r.randrange(3)==0:a=(j-2)%n
        elif f=="R":
            a=(j+3)%n;b=(j-2)%n;c=(j+1)%n
            if r.randrange(3)==0:b=(j+4)%n
        else:
            a=r.choice((0,1,n//2));b=(j+4)%n;c=r.choice((2,n-1,n//2))
            if a==j:a=(a+1)%n
        result.append((a,b,c))
    return tuple(result)
def pattern(s0,s1,s2,parent):
    a,b,c=parent
    return ((s0>>a)&1)|(((s1>>b)&1)<<1)|(((s2>>c)&1)<<2)
def rule_output(f,topology,rule,history,node):
    s0,s1,s2=history
    return ((s0>>node)&1)^((rule>>pattern(s0,s1,s2,topology[node]))&1)
def public_answer(topology,history,node,value):
    # Each public output directly constrains one causal truth-table bit,
    # with entity binding, snapshots and known copy/residual term respected.
    p=pattern(*history,topology[node])
    return p,value^((history[0]>>node)&1)
def build(name):
    if name not in ALL:raise RuntimeError("T37_UNFROZEN_OR_CONSUMED_ID")
    _,family,split,index=name.split("-");n=FAMILIES[family]
    r=random.Random(37008991+(0 if family=="R" else 993331)+
                    (0 if split=="train" else 570099)+int(index)*53011)
    topology=graph(family,split,r)
    law=r.choice(LOW if split=="train" else HIGH)
    sequences=[]
    for timepoint in range(3):
        history=tuple(r.randrange(1<<n) for _ in range(3))
        action=(r.randrange(n),r.randrange(2)) if timepoint==1 else None
        nodes=r.sample(range(n),2)
        observed=[]
        for node in nodes:
            out=rule_output(family,topology,law,history,node)
            if action is not None and action[0]==node:out=action[1]
            observed.append((node,out))
        sequences.append({"history":history,"do":action,
                          "observed":tuple(observed)})
    targets=[]
    for q in range(4):
        history=tuple(r.randrange(1<<n) for _ in range(3))
        node=r.randrange(n)
        action=(r.randrange(n),r.randrange(2)) if q==1 else None
        # Avoid fully disclosed trivial copy on intervened target.
        if action is not None and action[0]==node:node=(node+1)%n
        output=rule_output(family,topology,law,history,node)
        if action is not None and action[0]==node:output=action[1]
        targets.append({"history":history,"node":node,"do":action,
                        "_private_truth":output})
    public={"family":family,"entity_count":n,
         "topology":topology,"sequences":tuple(sequences),
         "queries":tuple({k:q[k] for k in ("history","node","do")} for q in targets)}
    return {"id":name,"public":public,"_hidden_law":law,
       "_truth":tuple(q["_private_truth"] for q in targets),
       "_origin":"INTERNALLY_ORIGINATED_CAUSAL_LAW_NOT_EXTERNAL"}
def public_META(item):
    if item["id"] not in TRAIN_META:raise PermissionError("T37_META_ONLY")
    return item["public"]
def public_FIT(item):
    if item["id"] not in TRAIN_FIT:raise PermissionError("T37_TRAIN_VALIDATION_ONLY")
    return item["public"]
def public_DEV(item):
    if item["id"] not in DEV:raise PermissionError("T37_PROSPECTIVE_ONLY")
    return item["public"]
