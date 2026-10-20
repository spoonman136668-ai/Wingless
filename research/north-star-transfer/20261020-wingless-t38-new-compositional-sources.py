"""T38 new SOURCE-ONLY unseen compositional causal gate family.

Unknown shared rule: hidden G(a,b) OP H(b,c) with two separately reusable
two-input Boolean gate atoms and hidden XOR/AND composition. Source TRAIN
pairs satisfy a frozen residue inequality; original prospective DEV pairs
satisfy its complement. Every individual atomic gate can appear in TRAIN
and DEV, but pairs in prospective DEV are new compositions. Graph families
C and D also change unseen temporal entity relationships between splits.

Never pass private G,H,OP or source sampling priors to training/inference.
Do not call build() for any DEV identity before TRAIN native admission.
"""
from __future__ import annotations
import random
NODES={"C":8,"D":10}
GATES=(8,14,6,7,1,9,2,11)
TRAIN_META=tuple(f"t38-{f}-train-{i:03d}" for f in ("C","D") for i in range(20))
TRAIN_VALID=tuple(f"t38-{f}-train-{i:03d}" for f in ("C","D") for i in range(20,28))
DEV=tuple(f"t38-{f}-dev-{i:03d}" for f in ("C","D") for i in range(32))
ALL=TRAIN_META+TRAIN_VALID+DEV
assert len(ALL)==120 and len(set(ALL))==120
def gate_output(g,i):
    return (GATES[g]>>i)&1
def composed_table(g,h,op):
    table=0
    for p in range(8):
        a=p&1;b=(p>>1)&1;c=(p>>2)&1
        x=gate_output(g,a|(b<<1))
        y=gate_output(h,b|(c<<1))
        v=x^y if op==0 else x&y
        table|=v<<p
    return table
HYPOTHESES=tuple((g,h,op,composed_table(g,h,op))
                 for g in range(8) for h in range(8) for op in range(2))
assert len(HYPOTHESES)==128
TRAIN_PAIRS=tuple((g,h) for g in range(8) for h in range(8)
                  if (g+3*h)%4!=0)
DEV_PAIRS=tuple((g,h) for g in range(8) for h in range(8)
                if (g+3*h)%4==0)
assert len(TRAIN_PAIRS)==48 and len(DEV_PAIRS)==16
assert set(TRAIN_PAIRS).isdisjoint(DEV_PAIRS)
def graph(f,split,r):
    n=NODES[f];g=[]
    for node in range(n):
        if split=="train" and f=="C":
            pool=list(range(max(0,node-4),node))
            a=r.choice(pool) if pool else (node+1)%n
            b=r.choice(pool) if pool else (node+2)%n
            c=r.choice(pool) if pool else (node+3)%n
        elif split=="train":
            a=(node-1)%n;b=(node+2)%n;c=(node-3)%n
        elif f=="C":
            a=(node+3)%n;b=(node-2)%n;c=(node+1)%n
            if r.randrange(3)==0:c=(node+4)%n
        else:
            a=r.choice((0,1,n//2))
            b=(node+4)%n
            c=r.choice((2,n-1,n//2))
        g.append((a,b,c))
    return tuple(g)
def pattern(history,parent):
    s0,s1,s2=history
    a,b,c=parent
    return ((s0>>a)&1)|(((s1>>b)&1)<<1)|(((s2>>c)&1)<<2)
def output(history,topology,truth_table,node):
    return ((history[0]>>node)&1)^((truth_table>>pattern(history,topology[node]))&1)
def public_clause(history,topology,node,value):
    return pattern(history,topology[node]),value^((history[0]>>node)&1)
def build(identifier):
    if identifier not in ALL:raise ValueError("T38_UNFROZEN_SOURCE_ID")
    _,family,split,index=identifier.split("-")
    n=NODES[family]
    r=random.Random(38002011+(0 if family=="C" else 870011)
            +(0 if split=="train" else 601291)+int(index)*9209)
    possible=TRAIN_PAIRS if split=="train" else DEV_PAIRS
    gweights=(5,4,3,2,3,2,1,1)
    hweights=(1,3,2,4,3,2,1,4)
    pairweights=[gweights[g]*hweights[h] for g,h in possible]
    g,h=r.choices(possible,weights=pairweights,k=1)[0]
    op=r.choices((0,1),weights=(3,1),k=1)[0]
    assert ((g+3*h)%4==0)==(split=="dev")
    rule=composed_table(g,h,op)
    topo=graph(family,split,r)
    probes=[]
    for k in range(5):
        history=tuple(r.randrange(1<<n) for _ in range(3))
        action=(r.randrange(n),r.randrange(2)) if k==2 else None
        observed=[]
        for node in r.sample(range(n),2):
            bit=output(history,topo,rule,node)
            if action is not None and action[0]==node:bit=action[1]
            observed.append((node,bit))
        probes.append({"history":history,"do":action,
                       "observed":tuple(observed)})
    queries=[];answers=[]
    for k in range(4):
        history=tuple(r.randrange(1<<n) for _ in range(3))
        node=r.randrange(n)
        action=(r.randrange(n),r.randrange(2)) if k==2 else None
        if action is not None and node==action[0]:node=(node+1)%n
        answer=output(history,topo,rule,node)
        queries.append({"history":history,"node":node,"do":action})
        answers.append(answer)
    return {"id":identifier,"public":{"family":family,"topology":topo,
           "probes":tuple(probes),"queries":tuple(queries)},
           "_private_g":g,"_private_h":h,"_private_op":op,
           "_private_rule":rule,"_answers":tuple(answers),
           "_source_provenance":"ORIGINAL_INTERNAL_NEW_COMPOSITION_NOT_EXTERNAL"}
def META(case):
    if case["id"] not in TRAIN_META:raise PermissionError("T38_META_TRAIN_ONLY")
    return case["public"]
def VALID(case):
    if case["id"] not in TRAIN_VALID:raise PermissionError("T38_TRAIN_VALIDATION_ONLY")
    return case["public"]
def FUTURE_DEV(case):
    if case["id"] not in DEV:raise PermissionError("T38_PROSPECTIVE_AFTER_ADMISSION_ONLY")
    return case["public"]
