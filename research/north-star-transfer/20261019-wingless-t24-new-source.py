"""T24 fresh source-only 4-role S3 causal information experiments K and L.

No learner/certifier import. TRAIN and DEV identities were frozen earlier.
Generator does not optimize for group-order identities or query savings.
"""
import itertools,random
OPS=tuple(itertools.permutations(range(3)))
TRAIN=tuple(f"t24-{family}-train-{i:03d}" for family in ("K","L") for i in range(20))
DEV=tuple(f"t24-{family}-dev-{i:03d}" for family in ("K","L") for i in range(32))
def token(rng,roles):return (rng.choice(roles),bool(rng.getrandbits(1)))
def walk(rng,roles,n):return tuple(token(rng,roles) for _ in range(n))
def graph_words(rng,family,index):
    a,b,c,d=rng.sample(range(4),4)
    # Counterfactual shared mediator paths; no forced algebraic identities.
    if family=="K":
        root=walk(rng,(a,b,c),3+index%2)
        mediated=root+walk(rng,(b,d),2+index%2)
        altered=root[:-1]+walk(rng,(c,d),3)
        contrast=walk(rng,(a,c),2)+root[-2:]+walk(rng,(b,d),2)
        return (mediated,altered,contrast)
    shared=walk(rng,(a,c),3+index%3)
    branch=walk(rng,(b,d),2+index%2)
    # Discontinuous context transfer, new latent role dependencies.
    return (shared+branch+walk(rng,(d,c),2),
            branch+shared[:2]+walk(rng,(a,d),2),
            walk(rng,(c,b),3)+shared[-2:]+walk(rng,(a,d),2))
def build(identifier):
    if identifier not in TRAIN and identifier not in DEV:raise ValueError("T24_ID_UNFROZEN")
    parts=identifier.split("-")
    family,partition,i=parts[1],parts[2],int(parts[3])
    rng=random.Random(240106781+
        (0 if family=="K" else 319747)+
        (0 if partition=="train" else 700021)+i*2309)
    op_idx=tuple(rng.randrange(6) for _ in range(4))
    rows=[]
    for j in op_idx:
        x,y=rng.sample(range(3),2)
        op=OPS[j]
        rows.append(((x,op[x]),(y,op[y])))
    words=graph_words(rng,family,i)
    starts=tuple(rng.randrange(3) for _ in words)
    truth=[]
    for w,start in zip(words,starts):
        val=start
        for role,inv in w:
            op=OPS[op_idx[role]]
            val=op.index(val) if inv else op[val]
        truth.append(val)
    if len(words)!=3 or any(not 5<=len(w)<=16 for w in words):
        raise RuntimeError("T24_SOURCE_WORD_LENGTH_DRIFT")
    return {"id":identifier,"words":words,"starts":starts,
            "rows":tuple(rows),"truth":tuple(truth),
            "family":family,"partition":partition,
            "origin":"INTERNAL_NOT_EXTERNAL"}
def visible_training(item):
    if item["id"] not in TRAIN:raise PermissionError("T24_TRAIN_ONLY")
    return {"words":item["words"],"starts":item["starts"],
            "first":tuple((row[0],) for row in item["rows"])}
def visible_dev(item):
    if item["id"] not in DEV:raise PermissionError("T24_DEV_ONLY")
    seen=set()
    def purchase(role):
        if role not in range(4) or role in seen:raise ValueError("T24_BAD_DUPLICATE_QUERY")
        seen.add(role)
        return item["rows"][role][1]
    return {"words":item["words"],"starts":item["starts"],
            "first":tuple((row[0],) for row in item["rows"]),
            "purchase":purchase,"charged":lambda:len(seen)}
