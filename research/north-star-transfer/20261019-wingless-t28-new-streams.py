"""T28 new Y/Z source ONLY: incrementally observed sparse-sensor nonlinear causal histories.

This frozen Y/Z source has distinct family identities and seeds; no T16-T27 scientific cohort is reopened.\nThe actual hidden topology/gate/initial world and future counterfactual truth
remain private to the scorer. All candidate arms only see public masked
chronological frames, interventions, and future experiment instruction.

No source family or identity from T16-T25 is opened or reused.
"""
from __future__ import annotations
import random
N=6
PARENTS=((1,2),(2,3),(3,5),(4,0),(5,1),(0,2))
TRAIN=tuple(f"t28-{f}-train-{i:03d}" for f in ("Y","Z") for i in range(20))
DEV=tuple(f"t28-{f}-dev-{i:03d}" for f in ("Y","Z") for i in range(32))
ALL=TRAIN+DEV
def template(t):
    if not 0<=t<12:raise ValueError("T26_BAD_TEMPLATE")
    shift=t%6
    sign=1 if t<6 else -1
    perm=tuple((sign*i+shift)%6 for i in range(N))
    parents=[None]*N
    for canonical,(a,b) in enumerate(PARENTS):
        parents[perm[canonical]]=(perm[a],perm[b])
    return tuple(parents)
TEMPLATES=tuple(template(t) for t in range(12))
assert len(set(TEMPLATES))==12
def next_state(bits,parents,gate,intervention=None):
    """Synchronous causality: all node outputs use the SAME previous state."""
    out=0
    for node,(a,b) in enumerate(parents):
        own=(bits>>node)&1;lhs=(bits>>a)&1;rhs=(bits>>b)&1
        change=(lhs&rhs if gate==0 else lhs|rhs)
        out|=(own^change)<<node
    if intervention is not None:
        j,x=intervention
        out=(out & ~(1<<j))|(x<<j)
    return out
def mask_value(rng,full,selected):
    mask=sum(1<<i for i in selected)
    return (mask,full&mask)
def sample(identity):
    if identity not in ALL:raise ValueError("T28_UNFROZEN_SOURCE_ID")
    _,family,partition,k=identity.split("-")
    i=int(k)
    rng=random.Random(280205117 + (0 if family=="Y" else 707919) +
       (0 if partition=="train" else 851971)+i*5813)
    topology=rng.randrange(12);gate=rng.randrange(2);initial=rng.randrange(64)
    par=TEMPLATES[topology]
    state=initial
    states=[initial]
    actions=[]
    for t in range(4):
        intervention=(rng.randrange(N),rng.randrange(2)) if t in (1,3) else None
        actions.append(intervention)
        state=next_state(state,par,gate,intervention)
        states.append(state)
    observations=[]
    if family=="Y":
        phase=rng.randrange(N)
        for t,hidden in enumerate(states):
            visible=tuple(sorted({(phase+t+j*2)%6 for j in range(2)}))
            observations.append(mask_value(rng,hidden,visible))
    else:
        for hidden in states:
            visible=tuple(sorted(rng.sample(range(N),2)))
            observations.append(mask_value(rng,hidden,visible))
    future_do=(rng.randrange(N),rng.randrange(2))
    target_nodes=tuple(rng.sample(range(N),3))
    one=next_state(state,par,gate,future_do)
    two=next_state(one,par,gate)
    answers=tuple((two>>node)&1 for node in target_nodes)
    return {"id":identity,"public":{
        "observations":tuple(observations),"history_actions":tuple(actions),
        "future_intervention":future_do,"targets":target_nodes},
        "_truth":answers,"_private_initial":initial,
        "_private_topology":topology,"_private_mode":gate,
        "_private_terminal_state":state,"_origin":"INTERNALLY_ORIGINATED_NOT_EXTERNALLY_REPLICATED"}
def train_public(item):
    if item["id"] not in TRAIN:raise PermissionError("T28_TRAIN_ONLY")
    return item["public"]
def dev_public(item):
    if item["id"] not in DEV:raise PermissionError("T28_DEV_ONLY")
    return item["public"]
