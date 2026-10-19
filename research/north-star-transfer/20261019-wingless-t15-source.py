"""T15 Wingless source-only prospective developer cohort definitions.
This module only DEFINES a new source distribution. Importing never births data.
TRAIN and OPEN seed namespaces disjoint from T7-T14; no protected WNS use.
Novel test: OPEN has longer, repeated noncommutative role structures.
"""
import random
TRAIN=tuple(range(-51031,-50999))
OPEN=tuple(range(-51163,-51099))
def birth(seed,k):
    assert seed in TRAIN or seed in OPEN, "UNREGISTERED_T15_SEED"
    rng=random.Random(15129029+abs(seed)*8339)
    dom=abs(seed)%2;n=11 if dom==0 else 3
    roles=[];true=[]
    for _ in range(4):
        p=rng.choice(k.candidates(dom))
        x,y=rng.sample(range(n),2)
        roles.append(((x,k.operator_apply(dom,p,x)),(y,k.operator_apply(dom,p,y))))
        true.append(p)
    # Train on 2/3-step compositions; OPEN on previously unseen 5/7-step
    # repeated-role orderings. Neither true operators nor hidden answer
    # are passed to source-blind inference.
    if seed in TRAIN:
        word=[(0,False),(1,False)] if seed%2 else [(2,False),(1,True),(0,False)]
    else:
        word=[(0,False),(1,False),(2,True),(0,True),(3,False)]
        if seed%2:word += [(1,True),(3,True)]
    start=rng.randrange(n);state=k.start(dom)
    for j,inv in word:state=k.compose_state(state,true[j],dom,inv)
    return (dom,tuple(roles),tuple(word),start),k.apply(state,start)
