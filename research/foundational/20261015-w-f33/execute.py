"""W-F33: learn a compact ordered role operator from whole path outcomes.

This is a PRE-ORIGINAL source implementation only. No future test seeds are
read by training/preflight. Input is directed graph, source, target and permitted
role query. Neither exact answer nor generator path enters the learner at test.
"""
from __future__ import annotations
import json,math,random,hashlib
ROLES=6
TRAIN_SEEDS=range(191000,191064)
RESERVED_SEEDS=range(192000,192012)
TRAIN_FAMILIES=(160,161,162,163)
RESERVED_FAMILIES=(164,165,166,167)
TRAIN_DEPTHS=(3,4,5)
RESERVED_DEPTHS=(6,7,8)
STATE_BUDGET=576
WEIGHT_BYTES=ROLES*ROLES*8
VARIANTS=("genuine","joint-relabel","query-order-swap","edge-role-swap",
          "node-alias","reverse-directed-edge","extra-decoys","false-target")
def rng(seed,family,depth,variant):
    return random.Random(seed*1000003+family*10007+depth*719+variant*1777)
def make(seed,family,depth,variant):
    r=rng(seed,family,depth,0)
    roles=[r.randrange(ROLES) for _ in range(depth)]
    for i in range(1,depth):
        if roles[i]==roles[i-1]:roles[i]=(roles[i]+1)%ROLES
    edges=[(i,i+1,roles[i]) for i in range(depth)]
    n=depth+1
    for j in range(depth):
        if family%2==0:
            edges.append((j,n,r.randrange(ROLES)));n+=1
        else:
            edges.extend(((j,n,r.randrange(ROLES)),(n,n+1,r.randrange(ROLES))));n+=2
    query=list(roles);src=0;dst=depth
    e=rng(seed,family,depth,variant)
    if variant==1:
        perm=list(range(ROLES));e.shuffle(perm)
        query=[perm[q] for q in query]
        edges=[(a,b,perm[q]) for a,b,q in edges]
    if variant==2:query[0],query[1]=query[1],query[0]
    if variant==3:
        a,b,q=edges[0];edges[0]=(a,b,(q+1)%ROLES)
    if variant==4:
        perm=list(range(n));e.shuffle(perm)
        edges=[(perm[a],perm[b],q) for a,b,q in edges]
        src=perm[src];dst=perm[dst]
    if variant==5:
        a,b,q=edges[depth//2];edges[depth//2]=(b,a,q)
    if variant==6:
        edges.append((src,n,e.randrange(ROLES)));n+=1
    if variant==7:dst=n;n+=1
    e.shuffle(edges)
    return {"seed":seed,"family":family,"depth":depth,"variant":variant,
            "edges":tuple(edges),"source":src,"target":dst,"query":tuple(query),"nodes":n}

def coded_oracle(x):
    active={x["source"]}
    for role in x["query"]:
        active={b for a,b,r in x["edges"] if a in active and r==role}
    return x["target"] in active

def sigmoid(z):
    return 1/(1+math.exp(-max(-25,min(25,z))))

def role_prob(w,edge,query):
    return sigmoid(w[edge*ROLES+query])

def propagate(x,w,return_path=False):
    # Soft directed role compositional transition. Full graph edges are
    # enumerated in input; no true chain or candidate joins provided.
    states={x["source"]:(1.,())};reads=0
    for step,needed in enumerate(x["query"]):
        nxt={}
        for a,b,edge_role in x["edges"]:
            reads+=1
            if a not in states:continue
            prob,path=states[a]
            candidate=prob*role_prob(w,edge_role,needed)
            if candidate>nxt.get(b,(0.,()))[0]:
                nxt[b]=(candidate,path+((edge_role,needed),))
        states=nxt
    score,path=states.get(x["target"],(0.,()))
    if return_path:return score,path,reads
    return score,reads

def initial_weights():
    return [0.0]*(ROLES*ROLES)

def train():
    w=initial_weights()
    cases=[]
    for seed in TRAIN_SEEDS:
        for family in TRAIN_FAMILIES:
            for depth in TRAIN_DEPTHS:
                for variant in range(8):
                    x=make(seed,family,depth,variant)
                    label=coded_oracle(x)
                    expect=variant in (0,1,4,6)
                    if label!=expect:raise ValueError("invalid independent oracle TRAIN")
                    cases.append((x,label))
    # Full-outcome logistic training, not precomputed local equality labels.
    # The normalized path score is compared with a frozen train-only
    # threshold. Backpropagated role-transitions follow the actual selected
    # path, not hidden generator roles.
    threshold=0.35
    for epoch in range(35):
        gradient=[0.]*len(w)
        for x,label in cases:
            score,path,_=propagate(x,w,return_path=True)
            p=max(1e-6,min(1-1e-6,score))
            target=float(label)
            if not path:continue
            dloss=(p-target)/(1-p)
            for e,q in path:
                ix=e*ROLES+q
                gradient[ix]+=dloss*(1-role_prob(w,e,q))
        for i in range(len(w)):
            grad=max(-5,min(5,gradient[i]/len(cases)))
            w[i]-=.50*(grad+.001*w[i])
    return tuple(w),len(cases),threshold

def predict(x,w,threshold=.35):
    score,reads=propagate(x,w)
    return score>=threshold**len(x["query"]),reads

def fixed(x):
    # Strong full-access coded comparator whose role equality is programmed.
    # It is NOT excluded simply because it performs better than the learner.
    return coded_oracle(x),len(x["edges"])*len(x["query"])

def no_role(x):
    active={x["source"]}
    for _ in x["query"]:active={b for a,b,r in x["edges"] if a in active}
    return x["target"] in active

def preflight():
    assert set(TRAIN_SEEDS).isdisjoint(RESERVED_SEEDS)
    assert set(TRAIN_FAMILIES).isdisjoint(RESERVED_FAMILIES)
    x=make(191000,160,3,0)
    assert coded_oracle(x)
    for v in range(8):
        z=make(191000,160,3,v)
        assert coded_oracle(z)==(v in (0,1,4,6))
    w,cases,t=train()
    assert cases==64*4*3*8
    assert len(w)==36 and len(w)*8<=STATE_BUDGET
    probe=make(191000,160,3,3)
    _,reads=predict(probe,w,t)
    assert reads==len(probe["query"])*len(probe["edges"])
    return {"schema":"wingless.f33.train-only-role-composition-preflight.v1",
            "stage":"TRAIN_ONLY_NO_RESERVED_ORIGINAL",
            "train_examples":cases,"weights_count":len(w),
            "train_coefficients_bytes":len(w)*8,"state_limit_bytes":STATE_BUDGET,
            "frozen_train_threshold":t,"training_weights":w,
            "train_seed_range":[TRAIN_SEEDS.start,TRAIN_SEEDS.stop],
            "reserved_test_range":[RESERVED_SEEDS.start,RESERVED_SEEDS.stop],
            "reserved_test_families":RESERVED_FAMILIES,
            "reserved_original_not_executed":True,"independent_F3_source":False}
if __name__=="__main__":print(json.dumps(preflight(),indent=2,sort_keys=True))
