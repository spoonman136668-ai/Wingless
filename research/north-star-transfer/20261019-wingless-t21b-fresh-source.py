"""T21-B fresh synthetic internal TRAIN/DEV; no T19/T20 source replay.
Generated after the separately committed T21-B hypothesis.
No scientific independence or externally authenticated first-use is claimed.
"""
from __future__ import annotations
import random

TRAIN=tuple(f"t21b-train-{i:03d}" for i in range(72))
DEV=tuple(f"t21b-dev-{i:03d}" for i in range(96))

def make_word(rng,train,i):
    a,b,c=rng.sample(range(3),3)
    if train:
        # Short repeated-binding commutator/conjugation curriculum.
        word=[(a,False),(b,False),(a,True),(b,True),
              (c,False),(a,False),(c,True)]
        target=8+i%7
    else:
        # Unseen longer layered composition and binding-retrieval grammar.
        word=[(a,False),(b,False),(a,True),(b,True),
              (c,False),(b,False),(c,True),(b,True),
              (a,False),(c,False),(a,True),(c,True),
              (b,False),(a,False)]
        target=15+i%9
    while len(word)<target:
        role=(a,b,c)[rng.randrange(3)]
        word.append((role,bool(rng.getrandbits(1))))
    return tuple(word)

def build(label,kernel):
    if label not in TRAIN and label not in DEV:
        raise ValueError("T21B_UNREGISTERED_SOURCE_ID")
    train=label in TRAIN
    index=int(label.rsplit("-",1)[-1])
    seed=(202610192140 if train else 202610192311)+index*104729
    rng=random.Random(seed)
    ops=tuple(rng.choice(kernel.candidates(1)) for _ in range(3))
    rows=[]
    for op in ops:
        x,y=rng.sample(range(3),2)
        rows.append(((x,kernel.operator_apply(1,op,x)),
                     (y,kernel.operator_apply(1,op,y))))
    word=make_word(rng,train,index)
    start=rng.randrange(3)
    state=kernel.start(1)
    for role,inverse in word:
        state=kernel.compose_state(state,ops[role],1,inverse=inverse)
    return {"id":label,"domain":1,"word":word,"start":start,
            "roles":tuple(rows),"private_truth":kernel.apply(state,start),
            "origin":"T21B_SELF_AUTHORED_INTERNAL_NOT_INDEPENDENT"}

def observed_interface(case):
    public=tuple((rows[0],) for rows in case["roles"])
    bought=set()
    def request(role):
        if role not in range(3) or role in bought:
            raise ValueError("T21B_DUPLICATE_OR_ILLEGAL_QUERY")
        bought.add(role)
        return case["roles"][role][1]
    return public,request,lambda:len(bought)
