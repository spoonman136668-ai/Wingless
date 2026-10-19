"""T18 INTERNAL source only. Does not instantiate or reuse T7-T17 originals.

Three persistent roles per novel algebraic word. Distinct TRAIN/DEV IDs.
Private operators/outputs only to fixture evaluator (not policy).
"""
import random

TRAIN=tuple(f"t18-train-{i:03d}" for i in range(24))
DEV=tuple(f"t18-dev-{i:03d}" for i in range(48))
WORDS=(
 ((0,False),(1,False),(2,False),(0,True),(1,True),(2,True)),
 ((0,False),(1,True),(0,True),(2,False),(1,False),(2,True)),
 ((0,False),(1,False),(2,False),(2,True),(1,True),(0,True)),
 ((2,False),(1,True),(0,False),(2,True),(1,False)),
 ((0,False),(1,False),(0,False),(2,False),(1,True),(0,True)),
 ((0,False),(1,False),(2,True),(0,False),(1,True),(2,False),(0,True)),
 ((1,False),(2,False),(0,True),(1,True),(2,True),(0,False)),
 ((0,False),(1,True),(2,False),(1,False),(0,True),(2,True))
)
def build(label,k):
    assert label in TRAIN or label in DEV, "T18_UNREGISTERED_INTERNAL_SOURCE"
    i=int(label.rsplit("-",1)[-1])
    istrain=label in TRAIN
    rng=random.Random(180384013+7373*i+(0 if istrain else 191111))
    domain=(i+(0 if istrain else 1))%2
    n=11 if domain==0 else 3
    rows=[]
    private_ops=[]
    for role in range(3):
        op=rng.choice(k.candidates(domain))
        a,b=rng.sample(range(n),2)
        rows.append(((a,k.operator_apply(domain,op,a)),
                     (b,k.operator_apply(domain,op,b))))
        private_ops.append(op)
    word=WORDS[(i*5+(0 if istrain else 3))%len(WORDS)]
    start=rng.randrange(n)
    state=k.start(domain)
    for role,inverse in word:
        state=k.compose_state(state,private_ops[role],domain,inverse=inverse)
    return {"source_id":label,"domain":domain,"roles":tuple(rows),
            "word":word,"start":start,"private_truth":k.apply(state,start),
            "oracle_never_shared_with_candidate":True}
