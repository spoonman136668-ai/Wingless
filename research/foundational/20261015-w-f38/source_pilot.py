"""W-F38 TRAIN-only identifiability / source evidence pilot.

For every fully consistent public set of partial bijections, construct a
different complete hidden grammar that yields *exactly the same public triples*.
A differing query answer proves non-identifiability. Failure to find a witness
is never presented as proof of identification, except when 7/8 public keys
uniquely determine every 8-node bijection.
"""
from __future__ import annotations
import hashlib,itertools,json,math,random
from collections import Counter
SEEDS=range(281000,281012)
FAMILIES=(280,281,282)
COVERAGE=(3,7)
DEPTHS=(4,7)
RESERVED_SEEDS=range(282000,282016)
RESERVED_FAMILIES=(284,285,286,287)
N=8;NROLES=8; WITNESS_SAMPLES=48
ARMS=("group_MAP","unrestricted_MAP","frequency_walk","uniform_completion_posterior","privileged_oracle")

def source_rng(seed,family,coverage,depth):
    if seed not in SEEDS or family not in FAMILIES or coverage not in COVERAGE or depth not in DEPTHS:
        raise ValueError("source identity forbidden, including all reserved original IDs")
    return random.Random(41771+seed*31013+family*17713+coverage*283+depth*719)
def gen(seed,family,coverage,depth):
    r=source_rng(seed,family,coverage,depth)
    true=[]
    for role in range(NROLES):
        if family==280:
            a=r.choice((1,3,5,7));b=r.randrange(N)
            op=tuple((a*x+b)%N for x in range(N))
        elif family==281:
            v=list(range(N));r.shuffle(v);op=tuple(v)
        else:
            a=r.choice((1,3,5,7));b=r.randrange(N);mask=r.randrange(1,N)
            op=tuple(((a*x+b)%N)^mask for x in range(N))
        true.append(op)
    triples=[]
    for role in range(NROLES):
        keys=r.sample(range(N),coverage)
        triples.extend((role,x,true[role][x]) for x in keys)
    r.shuffle(triples)
    start=r.randrange(N)
    query=tuple((r.randrange(NROLES),r.choice((-1,1))) for _ in range(depth))
    return tuple(true),tuple(triples),start,query

def truth_walk(maps,start,query):
    x=start
    for role,direction in query:
        op=maps[role]
        x=op[x] if direction==1 else op.index(x)
    return x

def dual_truth(maps,start,query):
    states=[int(x==start) for x in range(N)]
    for role,sign in query:
        out=[0]*N
        for x,val in enumerate(states):
            y=maps[role][x] if sign==1 else maps[role].index(x)
            out[y]+=val
        states=out
    assert sum(states)==1
    return states.index(1)

def completions(triples,rng):
    role_maps=[]
    for role in range(NROLES):
        known={x:y for r,x,y in triples if r==role}
        if len(set(known.values()))!=len(known):
            raise ValueError("invalid public bijection examples")
        src=[x for x in range(N) if x not in known]
        dst=[y for y in range(N) if y not in known.values()]
        rng.shuffle(dst)
        mapping=dict(known)
        mapping.update(zip(src,dst))
        op=tuple(mapping[x] for x in range(N))
        if sorted(op)!=list(range(N)):raise ValueError("invalid completion")
        role_maps.append(op)
    return tuple(role_maps)

def group_candidates():
    group=set()
    for a in (1,3,5,7):
        for b in range(N):
            group.add(tuple((a*x+b)%N for x in range(N)))
            for mask in range(1,N):
                group.add(tuple(((a*x+b)%N)^mask for x in range(N)))
    return tuple(sorted(group))
GROUP=group_candidates()
def map_controls(triples):
    group=[];unrestricted=[]
    for role in range(NROLES):
        known={x:y for r,x,y in triples if r==role}
        # Group model is a MATCHED inductive bias, not given hidden family.
        g=max(GROUP,key=lambda op:(sum(op[x]==y for x,y in known.items()),tuple(-v for v in op)))
        group.append(g)
        src=[x for x in range(N) if x not in known]
        dst=sorted(y for y in range(N) if y not in known.values())
        m=dict(known);m.update(zip(src,dst))
        unrestricted.append(tuple(m[x] for x in range(N)))
    return tuple(group),tuple(unrestricted)

def frequency(triples,start,query):
    current=start
    for role,sign in query:
        matches=[y if sign==1 else x for r,x,y in triples if r==role and (x if sign==1 else y)==current]
        if matches:current=min(matches)
    return current

def entropy(counts):
    n=sum(counts.values())
    return -sum((v/n)*math.log2(v/n) for v in counts.values() if v)

def run():
    assert set(SEEDS).isdisjoint(RESERVED_SEEDS)
    assert set(FAMILIES).isdisjoint(RESERVED_FAMILIES)
    rows=[]
    for seed in SEEDS:
      for family in FAMILIES:
       for coverage in COVERAGE:
        for depth in DEPTHS:
            hidden,public,start,query=gen(seed,family,coverage,depth)
            truth=truth_walk(hidden,start,query)
            if truth!=dual_truth(hidden,start,query):raise ValueError("dual oracle drift")
            g,u=map_controls(public)
            predictions={"group_MAP":truth_walk(g,start,query),
                         "unrestricted_MAP":truth_walk(u,start,query),
                         "frequency_walk":frequency(public,start,query),
                         "privileged_oracle":truth}
            rr=random.Random(seed*8317+family*1009+depth*97+coverage)
            sampled=[]
            witness=None
            for sample in range(WITNESS_SAMPLES):
                alt=completions(public,rr)
                assert all(alt[role][x]==y for role,x,y in public)
                target=dual_truth(alt,start,query)
                sampled.append(target)
                if witness is None and target!=truth:
                    witness={"alternate_answer":target,"maps_sha256":hashlib.sha256(json.dumps(alt).encode()).hexdigest()}
            votes=Counter(sampled)
            predictions["uniform_completion_posterior"]=max(range(N),key=lambda x:(votes[x],-x))
            if coverage==7:
                if witness is not None or any(y!=truth for y in sampled):
                    raise ValueError("7/8 visible bijective mapping must be unique")
            rows.append(dict(seed=seed,family=family,coverage=coverage,depth=depth,
                start=start,ordered_query=query,truth=truth,predictions=predictions,
                public_triples_count=len(public),
                public_evidence_sha256=hashlib.sha256(json.dumps([public,start,query]).encode()).hexdigest(),
                nonidentifiability_witness=witness,posterior_answer_entropy_bits=entropy(votes),
                sampled_completion_count=WITNESS_SAMPLES))
    scores={a:{"correct":sum(r["predictions"][a]==r["truth"] for r in rows),"total":len(rows)} for a in ARMS}
    for s in scores.values():s["accuracy"]=s["correct"]/s["total"]
    witness_fraction={str(c):sum(r["nonidentifiability_witness"] is not None for r in rows if r["coverage"]==c)/
                      sum(r["coverage"]==c for r in rows) for c in COVERAGE}
    avg_entropy={str(c):sum(r["posterior_answer_entropy_bits"] for r in rows if r["coverage"]==c)/
                 sum(r["coverage"]==c for r in rows) for c in COVERAGE}
    qualified=(witness_fraction["3"]>=.45 and witness_fraction["7"]==0 and
               witness_fraction["3"]-witness_fraction["7"]>=.40 and
               scores["privileged_oracle"]["correct"]==len(rows))
    return dict(schema="wingless.f38.identifiability-trainonly.v1",
        stage="SOURCE_DIAGNOSTIC_NOT_LEARNED_ORIGINAL",seeds=[SEEDS.start,SEEDS.stop],families=FAMILIES,
        coverage=COVERAGE,depths=DEPTHS,questions=len(rows),
        reserved_original_seeds=[RESERVED_SEEDS.start,RESERVED_SEEDS.stop],
        reserved_original_families=RESERVED_FAMILIES,
        witness_fraction=witness_fraction,posterior_answer_entropy_bits=avg_entropy,
        scores=scores,rows=rows,source_same_repository_not_F3=True,
        no_learned_candidate_executed=True,primary_original_not_executed=True,
        resident_budget_not_qualified_by_python=True,external_evidence_triples_metered=True,
        verdict="TRAIN_ONLY_IDENTIFIABILITY_GAP_VERIFIED" if qualified else "TRAIN_ONLY_IDENTIFIABILITY_SOURCE_GATE_BLOCKED",
        accepted_ref_mutation=False,paid_api_calls=0)
if __name__=="__main__":print(json.dumps(run(),sort_keys=True,indent=2))
