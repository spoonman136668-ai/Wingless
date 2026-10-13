"""W-R37: bounded modular expression bank vs fixed primitive operation bank."""
import hashlib, importlib.util, json
from pathlib import Path

path=Path(__file__).resolve().parent.parent/"20261013-w-r26-latent-alphabet"/"execute.py"
spec=importlib.util.spec_from_file_location("w26_pinned",path)
w26=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w26)
assert len(w26.PERMUTATIONS)==24 and w26.fixture()=="PASS"
FAMILIES=("ADD4","XOR4","AFFINE4","MUL4")
ARMS=(("BASE_8","BASE",8),("BASE_16","BASE",16),
      ("DSL_8","DSL",8),("DSL_16","DSL",16))


def source_op(a,b,name):
    if name=="ADD4":return (a+b)%4
    if name=="XOR4":return a^b
    if name=="AFFINE4":return (2*a+b+1)%4
    if name=="MUL4":return (a*b)%4
    raise ValueError(name)


def convert(p,fn):
    inv=[0]*4
    for index,label in enumerate(p):inv[label]=index
    return tuple(p[fn(inv[a],inv[b])] for a,b in w26.ALL_CONTEXTS)


BASE_TABLES=tuple(sorted({
    convert(p,fn)
    for p in w26.PERMUTATIONS
    for fn in (
        lambda a,b:(a+b)%4, lambda a,b:a^b,
        lambda a,b:(a*b)%4, lambda a,b:(a-b)%4
    )}))
DSL_TABLES=tuple(sorted({
    convert(p,lambda a,b,c0=c0,c1=c1,c2=c2:(c0+c1*a+c2*b)%4)
    for p in w26.PERMUTATIONS
    for c0 in range(4) for c1 in range(4) for c2 in range(4)
}.union({
    convert(p,lambda a,b:a^b) for p in w26.PERMUTATIONS
})))
assert BASE_TABLES and DSL_TABLES and all(len(t)==16 for t in DSL_TABLES)


def sequence(seed,ep):
    assert 248<=seed<=255 and 0<=ep<16
    ordered=sorted(range(16),key=lambda i:(
        hashlib.sha256(f"W37-order/{seed}/{i}".encode()).digest(),i))
    c=ordered[ep]
    xs=[(c+seed)%4,((c//4)+3*seed)%4]
    while len(xs)<32:
        xs.append(source_op(xs[-2],xs[-1],FAMILIES[seed%4]))
    p=w26.PERMUTATIONS[hashlib.sha256(
        f"W37-perm/{seed}".encode()).digest()[0]%24]
    return tuple(p[z] for z in xs)


def choose(seed,arm,which,budget):
    seqs=[sequence(seed,e) for e in range(16)]
    pairs=[(s[0],s[1]) for s in seqs]
    assert len(set(pairs))==16
    remaining=list(BASE_TABLES if which=="BASE" else DSL_TABLES)
    query=[]
    trace=hashlib.sha256()
    while len(query)<budget and remaining:
        options=[i for i in range(16) if i not in query]
        def score(e):
            a,b=pairs[e];idx=4*a+b
            counts=[0]*4
            for table in remaining:counts[table[idx]]+=1
            return (max(counts),sum(v*v for v in counts),e)
        ep=min(options,key=score)
        a,b=pairs[ep]
        label=seqs[ep][2]  # only AFTER query chosen
        remaining=[t for t in remaining if t[4*a+b]==label]
        query.append(ep)
        trace.update(bytes((ep,a,b,label))+
                     len(remaining).to_bytes(4,"big"))
    result="REJECT" if not remaining else ("ACCEPT" if len(remaining)==1
            else "ABSTAIN")
    table=remaining[0] if result=="ACCEPT" else None
    return {"seed":seed,"cohort":"development" if seed<252 else "heldout",
            "family":FAMILIES[seed%4],"arm":arm,"budget":budget,
            "bank":which,"initial_candidates":len(BASE_TABLES if which=="BASE" else DSL_TABLES),
            "remaining_candidates":len(remaining),
            "labels_used":len(query),"query_order":query,
            "decision":result,
            "table_sha256":hashlib.sha256(bytes(table)).hexdigest() if table else None,
            "query_sha256":trace.hexdigest()},table


def forecast(seed,table):
    confusion=[[0]*4 for _ in range(4)]
    trace=hashlib.sha256()
    for ep in range(16):
        seq=sequence(seed,ep)
        for t in range(3,32):
            guess=table[4*seq[t-2]+seq[t-1]]
            truth=seq[t]
            confusion[truth][guess]+=1
            trace.update(bytes((ep,t,truth,guess)))
    assert sum(map(sum,confusion))==464
    return {"correct":sum(confusion[i][i] for i in range(4)),
            "total":464,"confusion":confusion,"trace_sha256":trace.hexdigest()}


def run():
    fingerprints={str(seed):hashlib.sha256(bytes(
        value for ep in range(16) for value in sequence(seed,ep))).hexdigest()
        for seed in range(248,256)}
    if len(set(fingerprints.values()))!=8:
        return {"schema":"wingless.w-r37-bounded-dsl.v1",
                "classification":"INVALID_IDENTITIES","fingerprints":fingerprints,
                "rows":[],"model_calls":0}
    rows=[]
    for seed in range(248,256):
        for arm,which,budget in ARMS:
            result,table=choose(seed,arm,which,budget)
            rows.append({**result,"test":forecast(seed,table) if table else None})
    held=[r for r in rows if r["seed"]>=252]
    stats={}
    for arm,_,_ in ARMS:
        stats[arm]={family:{
            "decision":next(r["decision"] for r in held
                            if r["arm"]==arm and r["family"]==family),
            "correct":next((r["test"]["correct"] if r["test"] else None)
                           for r in held if r["arm"]==arm and r["family"]==family),
            "labels":next(r["labels_used"] for r in held
                          if r["arm"]==arm and r["family"]==family)}
            for family in FAMILIES}
    def perfect(arm,family):
        return stats[arm][family]["decision"]=="ACCEPT" and stats[arm][family]["correct"]==464
    if (stats["BASE_16"]["AFFINE4"]["decision"]=="ACCEPT" or
        stats["DSL_16"]["MUL4"]["decision"]=="ACCEPT"):
        outcome="FALSE_UNSEEN_ACCEPT"
    elif (perfect("BASE_16","MUL4") and
          stats["BASE_16"]["AFFINE4"]["decision"]=="REJECT" and
          perfect("DSL_16","AFFINE4") and
          stats["DSL_16"]["MUL4"]["decision"]=="REJECT"):
        outcome=("DSL_REACH_WITH_TRADEOFF" if perfect("DSL_8","AFFINE4")
                 else "PARTIAL")
    else:outcome="MIXED"
    return {"schema":"wingless.w-r37-bounded-dsl.v1",
            "classification":outcome,
            "parent":"fab90052e3850519f7182978876088f396d655b1",
            "bank_sizes":{"BASE":len(BASE_TABLES),"DSL":len(DSL_TABLES)},
            "fingerprints":fingerprints,"heldout":stats,"rows":rows,
            "model_calls":0,
            "claim_limit":"Fixed modular-affine expression grammar plus XOR, finite relabelings, and explicit supervised t2 queries."}


if __name__=="__main__":
    result=run()
    assert result==run()
    print(json.dumps(result,sort_keys=True,indent=2))
