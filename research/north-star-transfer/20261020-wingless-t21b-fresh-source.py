"""Fresh T21-B internal causal composition development source, originated post-freeze.

Only the harness holds complete second calibration rows and private truth.
The policy receives the first rows, word/start, and separately purchased rows.
Self-authored, not independently sourced or scientific first-use evidence.
"""
import random
TRAIN=tuple("t21b-train-%03d"%i for i in range(40))
DEV=tuple("t21b-dev-%03d"%i for i in range(64))

def build(label,cert):
    if label not in TRAIN and label not in DEV:
        raise ValueError("T21B_UNREGISTERED_SOURCE")
    train=label in TRAIN
    i=int(label.rsplit("-",1)[-1])
    rng=random.Random(202610200621+i*1889+(0 if train else 618719))
    operators=tuple(rng.choice(cert.PERMS) for _ in range(3))
    roles=[]
    for op in operators:
        x,y=rng.sample(range(3),2)
        roles.append(((x,op[x]),(y,op[y])))
    relabel=list(rng.sample(range(3),3))
    length=(8+i%4) if train else (13+i%4)
    if train:
        # Cyclic binding reuse with sporadic departures.
        base=[relabel[j%3] for j in range(length)]
        for j in range(3,length):
            if j%4==3:
                options=[r for r in range(3) if r!=base[j-1] and (j+1>=length or r!=base[j+1])]
                base[j]=rng.choice(options)
    else:
        # Disjoint DEV generator regime: longer interleaved mirrored
        # causal blocks with repeated bindings, not a T20 word replay.
        motif=relabel+list(reversed(relabel))
        base=[motif[(j+i)%6] for j in range(length)]
        for j in range(1,length):
            if base[j]==base[j-1]:
                base[j]=next(r for r in relabel if r!=base[j-1])
        for j in range(4,length,5):
            base[j]=next(r for r in relabel if r!=base[j-1] and
                        (j+1>=length or r!=base[j+1]))
    # A single operator role always has a persistent hidden identity.
    word=tuple((role,bool(rng.getrandbits(1))) for role in base)
    assert len(word)==length and len(set(r for r,_ in word))==3
    assert all(word[j][0]!=word[j-1][0] for j in range(1,length))
    start=rng.randrange(3)
    target=start
    for role,inv in word:
        target=cert.apply_op(1,operators[role],target,inv)
    return {"id":label,"word":word,"start":start,"full_rows":tuple(roles),
            "private_truth":target,"origin":"T21B_FRESH_INTERNAL_SELF_AUTHORED"}

def public_interface(case):
    seen=set()
    visible=tuple((role[0],) for role in case["full_rows"])
    def request(role):
        if role not in range(3) or role in seen:
            raise ValueError("T21B_DUPLICATE_OR_ILLEGAL_PURCHASE")
        seen.add(role)
        return case["full_rows"][role][1]
    return visible,request,lambda:len(seen)
