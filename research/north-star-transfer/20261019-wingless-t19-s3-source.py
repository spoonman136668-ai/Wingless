"""T19 self-authored internal S3 development source, *not* independent science.
T19 source is disjoint by namespace from prior T7-T18 and T16 originals.
The policy receives only the public first pairs and query() responses.
"""
import random

TRAIN=tuple("t19-train-%03d"%i for i in range(40))
DEV=tuple("t19-dev-%03d"%i for i in range(64))
WORDS=(
 ((0,False),(2,True),(1,False),(0,True),(2,False),(1,True),(0,False),(2,True)),
 ((2,False),(0,False),(1,True),(2,True),(1,False),(0,True),(2,False),(0,False)),
 ((1,True),(0,False),(2,False),(1,False),(0,True),(2,True),(1,False),(0,False)),
 ((0,False),(1,True),(2,False),(0,False),(1,False),(2,True),(0,True),(1,False),(2,False)),
 ((2,False),(1,False),(0,True),(1,True),(2,True),(0,False),(2,False),(1,True),(0,False)),
 ((1,False),(2,True),(0,False),(2,False),(1,True),(0,True),(1,False),(2,True))
)
def build(label,k):
    if label not in TRAIN and label not in DEV:
        raise ValueError("T19_UNREGISTERED_SOURCE_ID")
    i=int(label.rsplit("-",1)[-1])
    train=label in TRAIN
    rng=random.Random(190417011+919*i+(0 if train else 997331))
    roles=[]; private_ops=[]
    for _ in range(3):
        op=rng.choice(k.candidates(1))
        a,b=rng.sample(range(3),2)
        roles.append(((a,k.operator_apply(1,op,a)),(b,k.operator_apply(1,op,b))))
        private_ops.append(op)
    word=WORDS[(i*7+(0 if train else 3))%len(WORDS)]
    start=rng.randrange(3)
    state=k.start(1)
    for role,inv in word:
        state=k.compose_state(state,private_ops[role],1,inverse=inv)
    return {"id":label,"domain":1,"roles":tuple(roles),"word":word,
            "start":start,"private_truth":k.apply(state,start),
            "origin":"INTERNAL_DEV_SELF_AUTHORED_NOT_T16_FIRST_USE"}

def observed_interface(fixture):
    """The policy never receives fixture, all rows, or truth."""
    first=tuple((r[0],) for r in fixture["roles"])
    seen=set()
    def request(role):
        if role not in range(3) or role in seen:
            raise ValueError("T19_REPEATED_OR_ILLEGAL_PUBLIC_QUERY")
        seen.add(role)
        return fixture["roles"][role][1]
    return first,request,lambda: len(seen)
