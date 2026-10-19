"""T29: two NEW structural causal grammars, F/7 entities and G/9 entities.

F: two-parent competing mediators with XOR/NAND nonlinear change.
G: three-parent feedback with majority/parity nonlinear change.
Hidden template, mode and initial state survive chronological masked
observations. Public source never exposes private labels or future answers
to candidates. Unlike T26-T28, node count, parent arity and transition
grammar genuinely differ. No candidate/baseline/model import.
"""
from __future__ import annotations
import random
NODES={"F":7,"G":9}
TRAIN=tuple(f"t29-{f}-train-{k:03d}" for f in ("F","G") for k in range(20))
DEV=tuple(f"t29-{f}-dev-{k:03d}" for f in ("F","G") for k in range(32))
ALL=TRAIN+DEV
def build_templates(family):
    n=NODES[family]
    rng=random.Random(290103+ (0 if family=="F" else 440123))
    templates=[];seen=set()
    while len(templates)<8:
        perm=list(range(n));rng.shuffle(perm)
        parents=[None]*n
        for i in range(n):
            if family=="F":
                # A causal feeder DAG with a seeded feedback root.
                pool=[p for p in range(n) if p!=i]
                if i>=3:pool=[p for p in range(i)]
                chosen=rng.sample(pool,2)
            else:
                # Three-parent recurrent feedback with nonlocal inputs.
                chosen=rng.sample([p for p in range(n) if p!=i],3)
            parents[perm[i]]=tuple(perm[p] for p in chosen)
        out=tuple(parents)
        if out not in seen:templates.append(out);seen.add(out)
    return tuple(templates)
TEMPLATES={f:build_templates(f) for f in ("F","G")}
assert all(len(v)==8 and len(set(v))==8 for v in TEMPLATES.values())
def next_state(family,bits,template,mode,intervention=None):
    n=NODES[family];result=0
    for node,parents in enumerate(template):
        own=(bits>>node)&1
        children=[(bits>>a)&1 for a in parents]
        if family=="F":
            a,b=children
            change=(a^b) if mode==0 else 1-(a&b)
        else:
            a,b,c=children
            change=int(a+b+c>=2) if mode==0 else a^b^c
        result|=(own^change)<<node
    if intervention is not None:
        node,value=intervention
        result=(result&~(1<<node)) | (value<<node)
    return result
def sample(identity):
    if identity not in ALL:raise ValueError("T29_SOURCE_ID_NOT_FROZEN")
    _,family,split,number=identity.split("-");k=int(number)
    seed=29010913+(0 if family=="F" else 717911)+(0 if split=="train" else 901037)+k*10937
    rng=random.Random(seed);n=NODES[family]
    chosen=rng.randrange(8);mode=rng.randrange(2);state=rng.randrange(1<<n)
    template=TEMPLATES[family][chosen]
    frames=[]
    def observe(value,t):
        count=2 if family=="F" else 3
        if family=="F":
            phase=(k+chosen+t)%n
            indices=tuple(sorted({(phase+2*j)%n for j in range(count)}))
        else:
            indices=tuple(sorted(rng.sample(range(n),count)))
        mask=sum(1<<p for p in indices)
        return(mask,value&mask)
    frames.append(observe(state,0))
    actions=[]
    for t in range(4):
        action=(rng.randrange(n),rng.randrange(2)) if t in (1,3) else None
        actions.append(action)
        state=next_state(family,state,template,mode,action)
        frames.append(observe(state,t+1))
    future=(rng.randrange(n),rng.randrange(2))
    targets=tuple(rng.sample(range(n),3))
    after=next_state(family,state,template,mode,future)
    after=next_state(family,after,template,mode)
    truth=tuple((after>>node)&1 for node in targets)
    return {"id":identity,"family":family,
       "public":{"family":family,"observations":tuple(frames),
                 "history_actions":tuple(actions),"future_intervention":future,
                 "targets":targets},
       "_hidden_template":chosen,"_hidden_mode":mode,
       "_hidden_current_state":state,"_private_future_answers":truth,
       "_source_origin":"INTERNAL_NEW_CAUSAL_GRAMMAR_NOT_EXTERNAL"}
def train_public(item):
    if item["id"] not in TRAIN:raise PermissionError("T29_TRAIN_ONLY")
    return item["public"]
def dev_public(item):
    if item["id"] not in DEV:raise PermissionError("T29_DEV_ONLY")
    return item["public"]
