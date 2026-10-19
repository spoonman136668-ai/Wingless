"""T23 brand-new internal source-only multi-question causal families I and J.

No candidate, scorer, T16B, or past source import. 40 TRAIN / 64 DEV.
Source words do not intentionally insert S3 group-order identities.
Private answers/hidden permutations never enter an arm's public view.
"""
import itertools,random
OPS=tuple(itertools.permutations(range(3)))
FAMILIES=("I","J")
TRAIN=tuple(f"t23-{f}-train-{i:03d}" for f in FAMILIES for i in range(20))
DEV=tuple(f"t23-{f}-dev-{i:03d}" for f in FAMILIES for i in range(32))
def flip(t):return (t[0],not t[1])
def word(rng,n):
    return tuple((rng.randrange(3),bool(rng.getrandbits(1))) for _ in range(n))
def queries(rng,family,i):
    """Unfamiliar branch/reconverge families, independently source authored."""
    shared=word(rng,5+(i%3))
    if family=="I":
        # Two branching causal pathways revisit the same latent roles.
        left=shared+word(rng,3)
        right=shared+word(rng,2)+tuple(flip(x) for x in shared[-2:])
        third=shared[:4]+word(rng,5)+(shared[-1],)
        fourth=shared+word(rng,4)
        result=(left,right,third,fourth)
    else:
        # Reconverging suffix and discontinuous common prefixes.
        mid=word(rng,3+(i%3))
        base=word(rng,3)
        result=(base+mid+shared[:3],
                base[:2]+word(rng,3)+mid+shared[:2],
                word(rng,3)+mid+shared[:4],
                base+word(rng,4)+shared[:2])
    return tuple(result)
def build(identifier):
    if identifier not in TRAIN and identifier not in DEV:
        raise ValueError("T23_UNKNOWN_SOURCE")
    _,family,partition,index=identifier.split("-")
    i=int(index)
    rng=random.Random(23010611+ (0 if family=="I" else 194801)+
                      (0 if partition=="train" else 476911)+i*1543)
    ops=tuple(OPS[rng.randrange(6)] for _ in range(3))
    role_pairs=[]
    for op in ops:
        x,y=rng.sample(range(3),2)
        role_pairs.append(((x,op[x]),(y,op[y])))
    task_words=queries(rng,family,i)
    starts=(rng.randrange(3),rng.randrange(3),rng.randrange(3),rng.randrange(3))
    truths=[]
    for w,start in zip(task_words,starts):
        v=start
        for j,inv in w:
            op=ops[j]
            v=op.index(v) if inv else op[v]
        truths.append(v)
    assert len(task_words)==4 and all(6<=len(w)<=18 for w in task_words)
    return {"id":identifier,"family":family,"partition":partition,
            "words":task_words,"starts":starts,
            "role_pairs":tuple(role_pairs),"private_truths":tuple(truths),
            "origin":"NEW_INTERNAL_NOT_EXTERNAL"}
def train_public(case):
    if case["id"] not in TRAIN:raise ValueError("T23_TRAIN_ONLY")
    return {"words":case["words"],"starts":case["starts"],"role_pairs":case["role_pairs"]}
def dev_public(case):
    if case["id"] not in DEV:raise ValueError("T23_DEV_ONLY")
    bought=set()
    def request(role):
        if role not in range(3) or role in bought:raise ValueError("T23_BAD_OR_REPEAT_REQUEST")
        bought.add(role)
        return case["role_pairs"][role][1]
    return {"words":case["words"],"starts":case["starts"],
            "first":tuple((pairs[0],) for pairs in case["role_pairs"]),
            "request":request,"charged":lambda:len(bought)}
