"""W-F40: train-only learned query-conditioned value of paid evidence.

True hidden role permutations are available only to TRAIN-label construction
and evaluator Port. Prospective EVAL action receives public triples, an
ordered typed query and a bank of observation-consistent model samples.
A seven-feature linear action-value head is fitted to expected probability
of the TRAIN true *final answer* after one paid observation.
External completion banks, optimizer and histories are metered separately.
No independent F3 custody, no primary original, and never reuse F39 cohort.
"""
from __future__ import annotations
import hashlib,json,math,random
from collections import Counter
N=8; ROLES=8;VISIBLE=3;BUDGET=2;SAMPLES=24;FEATURES=7
TRAIN=range(301000,301016);TRAIN_FAMILIES=(310,311,312)
EVAL=range(302000,302012);EVAL_FAMILIES=(314,315,316)
RESERVED=range(303000,303016);RESERVED_FAMILIES=(318,319,320,321)
LENGTHS=(3,5)
ARMS=("learned_query_value","query_frontier","posterior_information_gain",
      "random_eligible","fixed_acquisition","no_purchase","privileged_complete_map_oracle_unfair")
def guard(seed,family):
    if not ((seed in TRAIN and family in TRAIN_FAMILIES) or
            (seed in EVAL and family in EVAL_FAMILIES)):
        raise ValueError("unregistered source identity, including reserved primary")
def case(seed,family,length):
    guard(seed,family)
    r=random.Random(9491+seed*81173+family*29137+length*353)
    maps=[]
    for role in range(ROLES):
        t=family%3
        if t==0:
            a=r.choice((1,3,5,7));b=r.randrange(N)
            op=tuple((a*x+b)%N for x in range(N))
        elif t==1:
            v=list(range(N));r.shuffle(v);op=tuple(v)
        else:
            a=r.choice((1,3,5,7));b=r.randrange(N);mask=r.randrange(1,N)
            op=tuple(((a*x+b)%N)^mask for x in range(N))
        maps.append(op)
    public=[]
    for role in range(ROLES):
        for x in r.sample(range(N),VISIBLE):
            public.append((role,x,maps[role][x]))
    r.shuffle(public)
    start=r.randrange(N)
    path=tuple((r.randrange(ROLES),r.choice((-1,1))) for _ in range(length))
    return tuple(maps),tuple(public),start,path
def follow(maps,start,path):
    x=start
    for role,d in path:
        op=maps[role];x=op[x] if d==1 else op.index(x)
    return x
def independent_follow(maps,start,path):
    p=[int(x==start) for x in range(N)]
    for role,d in path:
        nxt=[0]*N
        for x,val in enumerate(p):
            y=maps[role][x] if d==1 else maps[role].index(x)
            nxt[y]+=val
        p=nxt
    assert sum(p)==1
    return p.index(1)
def uniform_worlds(public,seed):
    r=random.Random(seed)
    known=[{} for _ in range(ROLES)]
    for role,x,y in public:
        known[role][x]=y
    if any(len(set(v.values()))!=len(v) for v in known):
        raise ValueError("non-bijective public source")
    worlds=[]
    for _ in range(SAMPLES):
        maps=[]
        for k in known:
            unseen=[x for x in range(N) if x not in k]
            unknown_y=[y for y in range(N) if y not in k.values()]
            r.shuffle(unknown_y);m=dict(k);m.update(zip(unseen,unknown_y))
            maps.append(tuple(m[x] for x in range(N)))
        worlds.append(tuple(maps))
    return worlds
def entropy(values):
    c=Counter(values);n=len(values)
    return -sum((v/n)*math.log2(v/n) for v in c.values() if v)
def available(public):
    seen={(r,x) for r,x,y in public}
    return [(r,x) for r in range(ROLES) for x in range(N) if (r,x) not in seen]
def frontier(worlds,start,path,public):
    seen={(r,x) for r,x,y in public}
    freq=Counter();early=Counter()
    for maps in worlds:
        state=start
        for i,(role,d) in enumerate(path):
            x=state if d==1 else maps[role].index(state)
            if (role,x) not in seen:
                freq[(role,x)]+=1
                if i==0:early[(role,x)]+=1
            state=maps[role][x] if d==1 else x
    return freq,early
def observation_features(worlds,answers,start,path,public,choice,freq,early):
    role,x=choice
    before=entropy(answers)
    groups={}
    for m,ans in zip(worlds,answers):
        groups.setdefault(m[role][x],[]).append(ans)
    expected=sum(len(v)/len(answers)*entropy(v) for v in groups.values())
    information=max(0.,before-expected)
    n=len(worlds)
    role_fraction=sum(r==role for r,d in path)/len(path)
    distinct=len(groups)/N
    return (1.,freq[choice]/max(1,n*len(path)),early[choice]/n,
            information/3.,role_fraction,distinct,before/3.)
def proposals(public,start,path,seed):
    worlds=uniform_worlds(public,seed)
    answers=[follow(m,start,path) for m in worlds]
    freq,early=frontier(worlds,start,path,public)
    permitted=available(public)
    ranked=sorted(permitted,key=lambda q:(freq[q],early[q],-q[0],-q[1]),reverse=True)
    # A bounded same-information proposal space for learned, frontier and IG.
    pool=ranked[:14]
    features={q:observation_features(worlds,answers,start,path,public,q,freq,early) for q in pool}
    return pool,features,worlds,answers,freq,early
def fit():
    training=[]
    for seed in TRAIN:
      for fam in TRAIN_FAMILIES:
       for length in LENGTHS:
        private,public,start,path=case(seed,fam,length)
        answer=follow(private,start,path)
        pool,features,worlds,old,freq,early=proposals(public,start,path,seed*101+fam+length)
        for role,x in pool:
            # TRAIN-only evaluator label: chance assigned to the true final
            # task answer after actually purchasing this mapping from source.
            observed=private[role][x]
            filtered=[follow(m,start,path) for m in worlds if m[role][x]==observed]
            target=(filtered.count(answer)+.5)/(len(filtered)+N*.5) if filtered else 1./N
            training.append((features[(role,x)],target))
    w=[0.]*FEATURES
    for epoch in range(160):
        g=[0.]*FEATURES
        for x,y in training:
            pred=sum(a*b for a,b in zip(w,x))
            err=max(-3.,min(3.,pred-y))
            for i,z in enumerate(x):g[i]+=err*z
        for i in range(FEATURES):
            w[i]-=.11*(g[i]/len(training)+.01*w[i])
    return w,len(training)
class PaidPort:
    def __init__(self,private,public):
        self._source=private;self.public={(role,x) for role,x,y in public};self.traces=[]
    def buy(self,role,x):
        if len(self.traces)>=BUDGET or (role,x) in self.public:
            raise ValueError("exceeded 2 purchases or duplicated mapping")
        v=self._source[role][x]
        self.public.add((role,x));self.traces.append((role,x,v));return v
def select(method,public,start,path,seed,w):
    permitted=available(public)
    if not permitted:return None
    if method=="fixed_acquisition":return permitted[0]
    if method=="random_eligible":return random.Random(seed+999).choice(permitted)
    pool,features,worlds,answers,freq,early=proposals(public,start,path,seed)
    if method=="query_frontier":
        return max(pool,key=lambda q:(early[q],freq[q],-q[0],-q[1])) if pool else None
    if method=="posterior_information_gain":
        return max(pool,key=lambda q:(features[q][3],freq[q],-q[0],-q[1])) if pool else None
    if method=="learned_query_value":
        return max(pool,key=lambda q:(sum(a*b for a,b in zip(w,features[q])),freq[q],-q[0],-q[1])) if pool else None
    raise ValueError(method)
def summarize(rows,arm):
    n=len(rows);c=sum(r["answers"][arm]==r["truth"] for r in rows)
    return dict(correct=c,total=n,accuracy=c/n)
def run():
    assert set(TRAIN).isdisjoint(EVAL) and set(EVAL).isdisjoint(RESERVED)
    assert set(TRAIN_FAMILIES).isdisjoint(EVAL_FAMILIES)
    assert set(EVAL_FAMILIES).isdisjoint(RESERVED_FAMILIES)
    w,trained=fit()
    rows=[]
    for seed in EVAL:
      for family in EVAL_FAMILIES:
       for length in LENGTHS:
        private,public,start,path=case(seed,family,length)
        truth=follow(private,start,path)
        assert independent_follow(private,start,path)==truth
        answers={"privileged_complete_map_oracle_unfair":truth};receipts={};entropy_change={}
        for method in ARMS[:-1]:
            now=list(public);port=PaidPort(private,public)
            h0=entropy([follow(m,start,path) for m in uniform_worlds(now,seed*911+family)])
            for step in range(0 if method=="no_purchase" else BUDGET):
                choice=select(method,now,start,path,seed*911+family+step*17,w)
                if choice is None:break
                r,x=choice;y=port.buy(r,x);now.append((r,x,y))
            pred_worlds=uniform_worlds(now,seed*911+family+700)
            votes=Counter(follow(m,start,path) for m in pred_worlds)
            answers[method]=max(range(N),key=lambda x:(votes[x],-x))
            entropy_change[method]=h0-entropy([follow(m,start,path) for m in pred_worlds])
            receipts[method]=port.traces
        rows.append(dict(seed=seed,family=family,length=length,truth=truth,answers=answers,
            receipts=receipts,entropy_reduction=entropy_change,
            public_sha256=hashlib.sha256(json.dumps([public,start,path]).encode()).hexdigest()))
    scores={arm:summarize(rows,arm) for arm in ARMS}
    by_family={str(f):{a:summarize([r for r in rows if r["family"]==f],a) for a in ARMS} for f in EVAL_FAMILIES}
    fair=("query_frontier","posterior_information_gain","random_eligible","fixed_acquisition","no_purchase")
    gain=scores["learned_query_value"]["accuracy"]-max(scores[a]["accuracy"] for a in fair)
    consistent=all(by_family[str(f)]["learned_query_value"]["accuracy"]>=max(by_family[str(f)][a]["accuracy"] for a in fair) for f in EVAL_FAMILIES)
    return dict(schema="wingless.f40.learned-active-evidence.trainonly.v1",
        phase="FRESH_TRAIN_ONLY_EVALUATION_NOT_PRIMARY_ORIGINAL",
        train_seeds=[TRAIN.start,TRAIN.stop],train_families=TRAIN_FAMILIES,
        eval_seeds=[EVAL.start,EVAL.stop],eval_families=EVAL_FAMILIES,
        reserved_original_seeds=[RESERVED.start,RESERVED.stop],reserved_original_families=RESERVED_FAMILIES,
        questions=len(rows),training_action_labels=trained,model_weights=w,model_weight_bytes=FEATURES*8,
        paid_mapping_read_budget=BUDGET,posterior_worlds_per_update=SAMPLES,
        external_posterior_bytes_per_bank=SAMPLES*ROLES*N,
        native_runtime_memory_not_independently_measured=True,
        primary_original_not_executed=True,source_independently_born_F3=False,
        scores=scores,by_family=by_family,rows=rows,
        verdict="TRAIN_ONLY_LEARNED_ACTIVE_VALUE_GAIN" if gain>=.03 and consistent else "TRAIN_ONLY_LEARNED_ACTIVE_VALUE_NO_QUALIFIED_GAIN",
        accepted_ref_mutation=False,paid_api_calls=0)
if __name__=="__main__":print(json.dumps(run(),sort_keys=True,indent=2))
