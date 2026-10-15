"""W-F37: train-only inverse-consistent learned typed variable bindings.

Unlike row-only W-F36, the same doubly-stochastic role operator serves
forward and reverse edges. Sinkhorn projection approximates a bijection from
noisy public triples. Fresh synthetic families, no independent F3 source.
"""
from __future__ import annotations
import hashlib,json,math,random
from functools import lru_cache
TRAIN_SEEDS=range(261000,261012)
TRAIN_FAMILIES=(240,241,242)
RESERVED_SEEDS=range(262000,262016)
RESERVED_FAMILIES=(244,245,246,247)
DEPTHS=(4,6,8);N=8;ROLES=8;DEMOS=5
ARMS=("learned_sinkhorn","permutation_MAP","affine_MAP","frequency_walk","role_scrambled","privileged_oracle")

def get_rng(seed,family,depth):
    if seed not in TRAIN_SEEDS or family not in TRAIN_FAMILIES:
        raise ValueError("reserved or unregistered source")
    return random.Random(7391+seed*70001+family*17713+depth*293)

def make_case(seed,family,depth):
    r=get_rng(seed,family,depth)
    hidden=[]
    for role in range(ROLES):
        if family==240:
            a=r.choice((1,3,5,7));b=r.randrange(8)
            op=tuple((a*x+b)%8 for x in range(N))
        elif family==241:
            v=list(range(N));r.shuffle(v);op=tuple(v)
        else:
            a=r.choice((1,3,5,7));b=r.randrange(8);mask=r.randrange(1,8)
            op=tuple(((a*x+b)%8)^mask for x in range(N))
        hidden.append(op)
    public=[]
    for role in range(ROLES):
        for _ in range(DEMOS):
            src=r.randrange(N);dst=hidden[role][src]
            if r.random()<.10:dst=(dst+r.randrange(1,N))%N
            public.append((role,src,dst))
    r.shuffle(public)
    start=r.randrange(N)
    query=tuple((r.randrange(ROLES),1 if r.random()<.65 else -1) for _ in range(depth))
    return tuple(hidden),tuple(public),start,query

def secret_compose(hidden,start,query):
    x=start
    for role,direction in query:
        op=hidden[role]
        x=op[x] if direction==1 else op.index(x)
    return x

def oracle_independent(hidden,start,query):
    states=[float(i==start) for i in range(N)]
    for role,direction in query:
        target=[0.]*N
        for x,p in enumerate(states):
            y=hidden[role][x] if direction==1 else hidden[role].index(x)
            target[y]+=p
        states=target
    assert abs(sum(states)-1)<1e-9
    return states.index(1.)

def softmax(logits):
    m=max(logits);p=[math.exp(z-m) for z in logits];s=sum(p)
    return [z/s for z in p]

def fit_sinkhorn(public):
    scores=[[[0.0]*N for _ in range(N)] for _ in range(ROLES)]
    for epoch in range(90):
        lr=.15/(1+epoch/45)
        for role,src,dst in public:
            row=scores[role][src];p=softmax(row)
            for j in range(N):row[j]-=lr*(p[j]-(j==dst)+.001*row[j])
    matrices=[]
    for role in range(ROLES):
        a=[softmax(row) for row in scores[role]]
        # Project to a doubly stochastic matrix. Preserve uncertainty; the
        # inverse uses the same q8 matrix's transpose.
        for _ in range(40):
            for j in range(N):
                den=sum(a[x][j] for x in range(N))
                for x in range(N):a[x][j]/=max(1e-10,den)
            for x in range(N):
                den=sum(a[x])
                a[x]=[z/max(1e-10,den) for z in a[x]]
        q=[]
        for x in range(N):
            values=[round(255*z) for z in a[x]]
            imax=max(range(N),key=lambda i:values[i])
            values[imax]+=255-sum(values)
            if min(values)<0 or max(values)>255 or sum(values)!=255:raise ValueError("bad q8 projection")
            q.append(tuple(values))
        matrices.append(tuple(q))
    return tuple(matrices)

def learned_compose(q,start,query):
    state=[0.]*N;state[start]=1.
    for role,direction in query:
        nxt=[0.]*N
        for src,p in enumerate(state):
            for dst in range(N):
                if direction==1:
                    nxt[dst]+=p*q[role][src][dst]/255
                else:
                    # Query asks for inverse edge from current source value:
                    # probability over possible preceding subjects.
                    nxt[dst]+=p*q[role][dst][src]/255
        mass=sum(nxt)
        if mass<=0:raise ValueError("zero inverse mass")
        state=[x/mass for x in nxt]
    return max(range(N),key=lambda j:(state[j],-j))

def solve_assignment(scores):
    @lru_cache(None)
    def recurse(x,mask):
        if x==N:return (0.,())
        options=[]
        for y in range(N):
            if mask>>y&1:continue
            val,path=recurse(x+1,mask|(1<<y))
            options.append((val+scores[x][y],(y,)+path))
        return max(options,key=lambda z:(z[0],tuple(-i for i in z[1])))
    return recurse(0,0)[1]

def candidates():
    opts=set()
    for a in (1,3,5,7):
        for b in range(8):
            opts.add(tuple((a*x+b)%8 for x in range(8)))
            for mask in range(1,8):opts.add(tuple(((a*x+b)%8)^mask for x in range(8)))
    return tuple(sorted(opts))
GROUP=candidates()

def controls(public):
    counts=[[[0]*N for _ in range(N)] for _ in range(ROLES)]
    for role,src,dst in public:counts[role][src][dst]+=1
    assignment=tuple(solve_assignment(counts[role]) for role in range(ROLES))
    group=[]
    for role in range(ROLES):
        obs=[(src,dst) for r,src,dst in public if r==role]
        group.append(max(GROUP,key=lambda g:(sum(g[x]==y for x,y in obs),tuple(-j for j in g))))
    return assignment,tuple(group),counts

def frequency_compose(counts,start,query):
    x=start
    for role,direction in query:
        if direction==1:x=max(range(N),key=lambda j:(counts[role][x][j],-j))
        else:x=max(range(N),key=lambda j:(counts[role][j][x],-j))
    return x

def summarize(rows,arm):
    n=len(rows);correct=sum(r["predictions"][arm]==r["truth"] for r in rows)
    return dict(total=n,correct=correct,accuracy=correct/n)

def run():
    assert set(TRAIN_SEEDS).isdisjoint(RESERVED_SEEDS)
    assert set(TRAIN_FAMILIES).isdisjoint(RESERVED_FAMILIES)
    rows=[]
    for seed in TRAIN_SEEDS:
      for family in TRAIN_FAMILIES:
       for depth in DEPTHS:
        private,public,start,query=make_case(seed,family,depth)
        truth=secret_compose(private,start,query)
        if oracle_independent(private,start,query)!=truth:raise ValueError("independent oracle mismatch")
        q=fit_sinkhorn(public)
        hard,affine,counts=controls(public)
        reversed_query=tuple(reversed(query))
        predictions=dict(
            learned_sinkhorn=learned_compose(q,start,query),
            permutation_MAP=secret_compose(hard,start,query),
            affine_MAP=secret_compose(affine,start,query),
            frequency_walk=frequency_compose(counts,start,query),
            role_scrambled=learned_compose(q,start,reversed_query),
            privileged_oracle=truth)
        rows.append(dict(seed=seed,family=family,depth=depth,truth=truth,query=query,
                         predictions=predictions,public_triples=len(public),
                         public_sha256=hashlib.sha256(json.dumps([public,start,query]).encode()).hexdigest(),
                         order_changes_truth=truth!=secret_compose(private,start,reversed_query)))
    scores={arm:summarize(rows,arm) for arm in ARMS}
    family={str(f):{arm:summarize([r for r in rows if r["family"]==f],arm) for arm in ARMS} for f in TRAIN_FAMILIES}
    fair=ARMS[1:-1]
    gain=scores["learned_sinkhorn"]["accuracy"]-max(scores[a]["accuracy"] for a in fair)
    robust=all(family[str(f)]["learned_sinkhorn"]["accuracy"]>=max(family[str(f)][a]["accuracy"] for a in fair) for f in TRAIN_FAMILIES)
    return dict(schema="wingless.f37.doubly-stochastic-trainonly.v1",
        stage="DEVELOPMENT_TRAIN_ONLY_NO_PRIMARY_ORIGINAL",
        seeds=[TRAIN_SEEDS.start,TRAIN_SEEDS.stop],families=TRAIN_FAMILIES,depths=DEPTHS,
        reserved_primary_seeds=[RESERVED_SEEDS.start,RESERVED_SEEDS.stop],
        reserved_primary_families=RESERVED_FAMILIES,
        questions=len(rows),q8_model_bytes=512,distribution_scratch_bytes=16,native_core_limit_bytes=576,
        source_same_repo_not_independent_F3=True,external_optimizer_memory_metered_separately=True,
        primary_original_not_executed=True,order_counterfactual_changes=sum(r["order_changes_truth"] for r in rows),
        scores=scores,by_family=family,rows=rows,
        verdict="TRAIN_ONLY_SINKHORN_LEARNED_GAIN" if gain>=.03 and robust else "TRAIN_ONLY_SINKHORN_NO_QUALIFIED_GAIN",
        accepted_ref_mutation=False,paid_API_calls=0)
if __name__=="__main__":print(json.dumps(run(),sort_keys=True,indent=2))
