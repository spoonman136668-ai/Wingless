"""T21-C internally authored, source-blind prospective S3 relational families.

This is deliberately NOT outside-authored or independently replicated science.
Source code imports no model, certifier or candidate. The private truth is
never passed through observed_interface or public_train_case.
Source namespace is fresh and disjoint from all T16-T20 originals/DEV.
"""
import itertools
import random

OPS=tuple(itertools.permutations(range(3)))
FAMILIES=("E","F")
TRAIN=tuple(f"t21i-{family}-train-{i:03d}" for family in FAMILIES for i in range(20))
DEV=tuple(f"t21i-{family}-dev-{i:03d}" for family in FAMILIES for i in range(32))

def flip(item):
    return (item[0],not item[1])

def word_e(rng,i):
    """Intervention-conditioned binding, two composition routes and re-entry."""
    x,y,z=rng.sample(range(3),3)
    context=rng.getrandbits(1)
    switch=rng.randrange(3)
    a=(x,False); b=(y,True); c=(z,False)
    # Preserve cross-context binding of the same shared roles; vary order.
    if context:
        head=(a,c,flip(b),flip(a),b,c)
        tail=(flip(c),b,a,flip(b),flip(a),c)
    else:
        head=(b,a,flip(c),flip(b),c,a)
        tail=(flip(a),flip(c),b,c,flip(b),a)
    if switch==0:return head+tail
    if switch==1:return head+(flip(head[2]),head[1])+tail[:4]
    return head[:4]+(tail[0],tail[1])+tail

def word_f(rng,i):
    """New nested noncommutative commutator/braid grammar."""
    p,q,r=rng.sample(range(3),3)
    a=(p,False);b=(q,False);c=(r,False)
    comm=(a,b,flip(a),flip(b))
    other=(b,c,flip(b),flip(c))
    j=rng.randrange(4)
    if j==0:return comm+other+tuple(flip(x) for x in reversed(comm))
    if j==1:return (c,)+comm+(flip(c),)+other+(a,flip(b))
    if j==2:return other+comm+(flip(c),a,flip(b),c)
    return (a,c,flip(b),flip(a),b,flip(c))+comm+(c,flip(a))

def build(source_id):
    if source_id not in TRAIN and source_id not in DEV:
        raise ValueError("T21C_SOURCE_ID_NOT_FROZEN")
    parts=source_id.split("-")
    family=parts[1];partition=parts[2];i=int(parts[3])
    seed=21910371+(0 if family=="E" else 917411)+(0 if partition=="train" else 482039)+i*1471
    rng=random.Random(seed)
    op_idx=[rng.randrange(6) for _ in range(3)]
    roles=[]
    for j in range(3):
        x,y=rng.sample(range(3),2)
        op=OPS[op_idx[j]]
        roles.append(((x,op[x]),(y,op[y])))
    word=word_e(rng,i) if family=="E" else word_f(rng,i)
    assert 8<=len(word)<=16 and {k for k,_ in word}=={0,1,2}
    start=rng.randrange(3)
    value=start
    for j,inverse in word:
        op=OPS[op_idx[j]]
        value=op.index(value) if inverse else op[value]
    return {"source_id":source_id,"word":word,"start":start,
      "roles":tuple(roles),"private_truth":value,
      "family":family,"partition":partition,
      "source_origin":"INTERNALLY_AUTHORED_PROSPECTIVE_NOT_EXTERNALLY_REPLICATED"}

def public_train_case(case):
    if case["source_id"] not in TRAIN or case["partition"]!="train":
        raise PermissionError("T21C_TRAIN_ONLY")
    return {"word":case["word"],"start":case["start"],"roles":case["roles"]}

def observed_interface(case):
    if case["source_id"] not in DEV:
        raise PermissionError("T21C_DEV_ONLY")
    visible=tuple((r[0],) for r in case["roles"])
    queried=set()
    def request(role):
        if role not in range(3) or role in queried:
            raise ValueError("T21C_DUPLICATE_OR_ILLEGAL_PUBLIC_QUERY")
        queried.add(role)
        return case["roles"][role][1]
    return {"word":case["word"],"start":case["start"],
            "first_public_pairs":visible,
            "query_second_pair":request,
            "charged_queries":lambda:len(queried)}
