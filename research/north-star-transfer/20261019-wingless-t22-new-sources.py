"""T22 source-only prospective internal source families G/H.

This code has NO imports of candidate, T16B or previous source modules.
Its private labels/true operators must only be used by evaluation scorer.
Unlike the consumed T21 E/F source, this has fresh source namespace and
distinct exact/nonidentity algebraic workloads.
"""
import itertools
import random
OPS=tuple(itertools.permutations(range(3)))
FAMILIES=("G","H")
TRAIN=tuple(f"t22-{family}-train-{i:03d}" for family in FAMILIES for i in range(20))
DEV=tuple(f"t22-{family}-dev-{i:03d}" for family in FAMILIES for i in range(32))

def invert(word):
    return tuple((r,not inverse) for r,inverse in reversed(tuple(word)))

def make_word(rng,family,index):
    a,b,c=rng.sample(range(3),3)
    p=(a,False);q=(b,False);t=(c,False)
    if family=="G":
        # Six powers universal identity in S3; five or seven are deliberately
        # NOT universal. Mixing context changes which live role matters.
        n=(6,5,7,6)[index%4]
        m=(6,7,6,5)[(index//4)%4]
        head=tuple((a,bool(index%3==2)) for _ in range(n))
        other=tuple((b,bool((index//3)%2)) for _ in range(m))
        if index%3==0: return head+other+(t,p,q)
        if index%3==1: return (t,)+head+(q,)+other+(p,)
        return other+(t,)+head+(q,p)
    # [a,b] always lies in A3 and cube is identity. The 2nd/4th powers
    # form nontrivial controls; no source labels are learner features.
    comm=(p,q,(a,True),(b,True))
    repetitions=(3,2,4,3)[index%4]
    inner=comm*repetitions
    if index%3==0:return inner+(t,p)
    if index%3==1:return (t,)+inner+(q,(c,True),p)
    return (q,)+inner+(t,(b,True),p)

def build(identifier):
    if identifier not in TRAIN and identifier not in DEV:
        raise ValueError("T22_SOURCE_ID_UNREGISTERED")
    parts=identifier.split("-")
    family,partition,index=parts[1],parts[2],int(parts[3])
    # Disjoint namespaces and source-provenance seeds, not external origin.
    seed=220620181+(0 if family=="G" else 910001)+(0 if partition=="train" else 341021)+index*1901
    rng=random.Random(seed)
    ops=tuple(OPS[rng.randrange(6)] for _ in range(3))
    pairs=[]
    for op in ops:
        x,y=rng.sample(range(3),2)
        pairs.append(((x,op[x]),(y,op[y])))
    word=tuple(make_word(rng,family,index))
    assert 4<=len(word)<=32 and set(role for role,inv in word)=={0,1,2}
    start=rng.randrange(3)
    value=start
    for role,inv in word:
        op=ops[role]
        value=op.index(value) if inv else op[value]
    return {"source_id":identifier,"word":word,"start":start,
            "roles":tuple(pairs),"private_truth":value,
            "internal_origin":True,"family":family}

def public_train_case(case):
    if case["source_id"] not in TRAIN:
        raise ValueError("T22_TRAIN_ONLY")
    return {"word":case["word"],"start":case["start"],"roles":case["roles"]}

def public_dev_interface(case):
    if case["source_id"] not in DEV:
        raise ValueError("T22_DEV_ONLY")
    first=tuple((pairs[0],) for pairs in case["roles"])
    seen=set()
    def request(role):
        if role not in range(3) or role in seen:
            raise ValueError("T22_INVALID_QUERY")
        seen.add(role)
        return case["roles"][role][1]
    return {"word":case["word"],"start":case["start"],
            "first":first,"request":request,"charged":lambda:len(seen)}
