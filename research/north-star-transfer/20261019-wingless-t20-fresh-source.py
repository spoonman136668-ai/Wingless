"""T20 self-authored fresh S3 development family. Not independent first-use science.

This module does not import or replay T19 DEV or protected T16 sources.
Private operators and answers are supplied only to the scoring harness.
"""
import random

TRAIN=tuple("t20-train-%03d"%i for i in range(40))
DEV=tuple("t20-dev-%03d"%i for i in range(64))

def build(label, kernel):
    if label not in TRAIN and label not in DEV:
        raise ValueError("T20_UNREGISTERED_SOURCE_ID")
    i=int(label.rsplit("-",1)[-1])
    train=label in TRAIN
    rng=random.Random(202610190520+i*1693+(0 if train else 318791))
    operators=[rng.choice(kernel.candidates(1)) for _ in range(3)]
    # Two independent public calibration pairs for each shared role.
    roles=[]
    for op in operators:
        x,y=rng.sample(range(3),2)
        roles.append(((x,kernel.operator_apply(1,op,x)),
                      (y,kernel.operator_apply(1,op,y))))
    # Independently authored, longer composition grammar than T19.
    length=(6+i%7) if train else (12+i%5)
    for attempt in range(32):
        word=[(rng.randrange(3),bool(rng.getrandbits(1))) for _ in range(length)]
        reduced=[]
        for item in word:
            if reduced and reduced[-1]==(item[0],not item[1]):
                reduced.pop()
            else:
                reduced.append(item)
        if len(reduced)>=5 and len({r for r,_ in reduced})==3:
            break
    else:
        raise RuntimeError("T20_FRESH_GRAMMAR_EXHAUSTED")
    state=kernel.start(1)
    for role,inv in word:
        state=kernel.compose_state(state,operators[role],1,inverse=inv)
    start=rng.randrange(3)
    return {"id":label,"domain":1,"roles":tuple(roles),
            "word":tuple(word),"start":start,
            "private_truth":kernel.apply(state,start),
            "origin":"T20_FRESH_SELF_AUTHORED_INTERNAL_NOT_INDEPENDENT"}

def observed_interface(fixture):
    """Expose only first role observations; charge every second pair."""
    first=tuple((r[0],) for r in fixture["roles"])
    purchased=set()
    def request(role):
        if role not in range(3) or role in purchased:
            raise ValueError("T20_DUPLICATE_OR_ILLEGAL_QUERY")
        purchased.add(role)
        return fixture["roles"][role][1]
    return first,request,lambda:len(purchased)
