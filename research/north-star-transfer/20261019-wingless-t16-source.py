"""T16 Wingless inert prospective source grammar: never births on import.
New development identifiers, independent of all T7-T15 used sources.
TRAIN short compositions; OPEN order-shifted repeats and inverses.
Evaluator-only operator truth never enters candidate or public ranker.
"""
import random
TRAIN=tuple(range(-61031,-60999))
OPEN=tuple(range(-61163,-61099))
def birth(seed,k):
    assert seed in TRAIN or seed in OPEN,"T16_NOT_PREREGISTERED"
    rng=random.Random(16149091+abs(seed)*9377)
    domain=abs(seed)%2
    n=11 if domain==0 else 3
    roles=[];hidden=[]
    for i in range(4):
        operator=rng.choice(k.candidates(domain))
        a,b=rng.sample(range(n),2)
        roles.append(((a,k.operator_apply(domain,operator,a)),
                      (b,k.operator_apply(domain,operator,b))))
        hidden.append(operator)
    word=([(0,False),(1,True)] if seed in TRAIN and seed%2
       else [(0,False),(1,False),(2,True)] if seed in TRAIN
       else [(0,False),(1,True),(2,False),(0,True),(3,False),
             (2,True),(1,False),(3,True)])
    start=rng.randrange(n)
    value=k.start(domain)
    for j,inv in word:value=k.compose_state(value,hidden[j],domain,inv)
    return (domain,tuple(roles),tuple(word),start),k.apply(value,start)
