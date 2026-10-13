"""W-R36: active evidence acquisition for finite operation-family identification."""
import hashlib, importlib.util, json
from pathlib import Path

p=Path(__file__).resolve().parent.parent/"20261013-w-r26-latent-alphabet"/"execute.py"
spec=importlib.util.spec_from_file_location("w26_frozen",p)
w26=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w26)
assert w26.fixture()=="PASS" and len(w26.PERMUTATIONS)==24
OPS=("ADD4","XOR4","MUL4","SUB4","AFFINE4")
ARMS=(("PASSIVE_4",4,False),("ACTIVE_4",4,True),
      ("ACTIVE_8",8,True),("ACTIVE_16",16,True))


def op(a,b,kind):
    if kind=="ADD4":return (a+b)%4
    if kind=="XOR4":return a^b
    if kind=="MUL4":return (a*b)%4
    if kind=="SUB4":return (a-b)%4
    if kind=="AFFINE4":return (2*a+b+1)%4
    raise ValueError(kind)


def candidate(perm,kind):
    inv=[0]*4
    for i,observed in enumerate(perm):inv[observed]=i
    return tuple(perm[op(inv[a],inv[b],kind)]
                 for a,b in w26.ALL_CONTEXTS)


CANDIDATE_TABLES=tuple(sorted({
    candidate(p,kind)
    for p in w26.PERMUTATIONS
    for kind in OPS[:4]}))
assert CANDIDATE_TABLES and all(len(t)==16 for t in CANDIDATE_TABLES)


def order(seed):
    assert 232<=seed<=241
    return sorted(range(16),key=lambda i:(
        hashlib.sha256(f"W36-order/{seed}/{i}".encode()).digest(),i))


def sequence(seed,ep):
    o=order(seed)[ep]
    vals=[(o+seed)%4,((o//4)+3*seed)%4]
    family=OPS[seed%5]
    while len(vals)<32:
        vals.append(op(vals[-2],vals[-1],family))
    perm=w26.PERMUTATIONS[hashlib.sha256(
        f"W36-perm/{seed}".encode()).digest()[0]%24]
    return tuple(perm[k] for k in vals)


def pairs(seed):
    x=[sequence(seed,e)[:2] for e in range(16)]
    assert len(set(x))==16
    return x


def selection(seed,name,budget,active):
    inputs=pairs(seed)
    hypotheses=list(CANDIDATE_TABLES)
    queried=[]
    trail=hashlib.sha256()
    while len(queried)<budget and hypotheses:
        remaining=[e for e in range(16) if e not in queried]
        if not active:
            ep=remaining[0]
        else:
            def score(e):
                a,b=inputs[e]
                bins=[sum(t[4*a+b]==v for t in hypotheses) for v in range(4)]
                return (max(bins),sum(n*n for n in bins),e)
            ep=min(remaining,key=score)
        a,b=inputs[ep]
        target=sequence(seed,ep)[2]  # reveal only AFTER selecting this query
        hypotheses=[t for t in hypotheses if t[4*a+b]==target]
        queried.append(ep)
        trail.update(bytes((ep,a,b,target,len(hypotheses))))
    if not hypotheses:
        decision="REJECT";table=None
    elif len(hypotheses)==1:
        decision="ACCEPT";table=hypotheses[0]
    else:
        decision="ABSTAIN";table=None
    return {"seed":seed,"cohort":"development" if seed<237 else "heldout",
            "family":OPS[seed%5],"arm":name,"budget":budget,
            "labels_used":len(queried),"queried_episodes":queried,
            "remaining_tables":len(hypotheses),"decision":decision,
            "table_sha256":hashlib.sha256(bytes(table)).hexdigest() if table else None,
            "query_trace_sha256":trail.hexdigest()},table


def forecast(seed,table):
    matrix=[[0]*4 for _ in range(4)]
    trace=hashlib.sha256()
    for ep in range(16):
        xs=sequence(seed,ep)
        for t in range(3,32):
            predicted=table[4*xs[t-2]+xs[t-1]]
            actual=xs[t]
            matrix[actual][predicted]+=1
            trace.update(bytes((ep,t,actual,predicted)))
    assert sum(map(sum,matrix))==464
    return {"correct":sum(matrix[i][i] for i in range(4)),
            "total":464,"confusion":matrix,"trace_sha256":trace.hexdigest()}


def run():
    fingerprints={str(seed):hashlib.sha256(bytes(
        x for ep in range(16) for x in sequence(seed,ep))).hexdigest()
        for seed in range(232,242)}
    if len(set(fingerprints.values()))!=10:
        return {"schema":"wingless.w-r36-family.v1",
                "classification":"INVALID_IDENTITIES",
                "fingerprints":fingerprints,"rows":[],"model_calls":0}
    rows=[]
    for seed in range(232,242):
        for name,budget,active in ARMS:
            record,table=selection(seed,name,budget,active)
            rows.append({**record,"test":forecast(seed,table) if table else None})
    held=[r for r in rows if r["seed"]>=237]
    stats={}
    for name,_,_ in ARMS:
        subset=[r for r in held if r["arm"]==name]
        stats[name]={family:{
            "decision":next(r["decision"] for r in subset if r["family"]==family),
            "labels":next(r["labels_used"] for r in subset if r["family"]==family),
            "correct":next((r["test"]["correct"] if r["test"] else None)
                           for r in subset if r["family"]==family)}
            for family in OPS}
    def fullcorrect(a,f):
        return (stats[a][f]["decision"]=="ACCEPT" and
                stats[a][f]["correct"]==464)
    unknown_accepted=any(stats[a]["AFFINE4"]["decision"]=="ACCEPT"
                         for a,_,_ in ARMS)
    exhaustive=(all(fullcorrect("ACTIVE_16",f) for f in OPS[:4]) and
                stats["ACTIVE_16"]["AFFINE4"]["decision"]=="REJECT")
    bounded=(sum(fullcorrect("ACTIVE_8",f) for f in OPS[:4])>=3 and
             stats["ACTIVE_8"]["AFFINE4"]["decision"]=="REJECT")
    if unknown_accepted:
        label="FALSE_NOVEL_ACCEPT"
    elif exhaustive and bounded:
        n_active=sum(fullcorrect("ACTIVE_4",f) for f in OPS[:4])
        n_passive=sum(fullcorrect("PASSIVE_4",f) for f in OPS[:4])
        label="ACTIVE_EFFICIENCY" if n_active>n_passive else "BOUNDED_OPENSET"
    else:label="MIXED"
    return {"schema":"wingless.w-r36-family.v1",
            "classification":label,
            "parent":"ec08b32ef7617a51fad46e2bbf1b098f9bf85d6b",
            "candidate_tables":len(CANDIDATE_TABLES),
            "fingerprints":fingerprints,"heldout":stats,
            "rows":rows,"model_calls":0,
            "claim_limit":"Four finite candidate operation families and a predesigned unseen fifth; explicitly supervised label acquisition."}


if __name__=="__main__":
    original=run()
    assert original==run()
    print(json.dumps(original,sort_keys=True,indent=2))
