"""W-F39: two-purchase active relationship-evidence source feasibility.

Different F38 sources and no reuse of consumed answers. Paid relation mapping
queries reveal ONLY (role, subject -> object). Evaluator-private full grammar
is behind the port. Uniform public-consistent posterior samples approximate
final-answer uncertainty. TRAIN-only diagnostic; not a learned F39 model or F3.
"""
import hashlib,json,random,math
from collections import Counter

SEEDS=range(291000,291012)
FAMILIES=(290,291,292)
RESERVED_SEEDS=range(292000,292016)
RESERVED_FAMILIES=(294,295,296,297)
DEPTHS=(4,7)
N=8
ROLES=8
VISIBLE=3
READ_BUDGET=2
POSTERIOR_SAMPLES=24
METHODS=("posterior_information_gain","query_frontier","random_eligible","fixed_acquisition","no_purchase","privileged_complete_map_oracle_unfair")

def make_case(seed,family,depth):
    if seed not in SEEDS or family not in FAMILIES or depth not in DEPTHS:
        raise ValueError("unauthorized science source identity")
    r=random.Random(seed*77311+family*9007+depth*443+793)
    maps=[]
    for role in range(ROLES):
        if family==290:
            a=r.choice((1,3,5,7));b=r.randrange(N)
            perm=tuple((a*x+b)%N for x in range(N))
        elif family==291:
            vals=list(range(N));r.shuffle(vals);perm=tuple(vals)
        else:
            a=r.choice((1,3,5,7));b=r.randrange(N);mask=r.randrange(1,N)
            perm=tuple(((a*x+b)%N)^mask for x in range(N))
        maps.append(perm)
    public=[]
    for role in range(ROLES):
        for x in r.sample(range(N),VISIBLE):
            public.append((role,x,maps[role][x]))
    r.shuffle(public)
    start=r.randrange(N)
    ordered_roles=tuple((r.randrange(ROLES),r.choice((-1,1))) for _ in range(depth))
    return tuple(maps),tuple(public),start,ordered_roles

def walk(maps,start,query):
    x=start
    for role,direction in query:
        p=maps[role]
        x=p[x] if direction==1 else p.index(x)
    return x

def independently_compose(maps,start,query):
    states=[int(x==start) for x in range(N)]
    for role,direction in query:
        next_states=[0]*N
        for x,count in enumerate(states):
            dst=maps[role][x] if direction==1 else maps[role].index(x)
            next_states[dst]+=count
        states=next_states
    assert sum(states)==1
    return states.index(1)

def sample_completions(public,seed):
    rng=random.Random(seed)
    known=[{} for _ in range(ROLES)]
    for role,x,y in public:
        if x in known[role] and known[role][x]!=y:raise ValueError("inconsistent mapping source")
        known[role][x]=y
    generated=[]
    for _ in range(POSTERIOR_SAMPLES):
        maps=[]
        for role in range(ROLES):
            src=[x for x in range(N) if x not in known[role]]
            dst=[x for x in range(N) if x not in known[role].values()]
            rng.shuffle(dst)
            m=dict(known[role]);m.update(zip(src,dst))
            maps.append(tuple(m[x] for x in range(N)))
        generated.append(tuple(maps))
    return generated

def entropy(v):
    c=Counter(v);n=len(v)
    return -sum((k/n)*math.log2(k/n) for k in c.values() if k)

def posterior_answer(public,start,query,seed):
    sampled=sample_completions(public,seed)
    answers=[walk(m,start,query) for m in sampled]
    votes=Counter(answers)
    prediction=max(range(N),key=lambda x:(votes[x],-x))
    return prediction,entropy(answers),sampled,answers

def unknown(role,x,public):
    return not any(r==role and sx==x for r,sx,y in public)

def candidate_queries(sampled,public,start,query):
    counter=Counter()
    first=Counter()
    for m in sampled:
        x=start
        for step,(role,direction) in enumerate(query):
            if direction==1:
                qx=x
                x=m[role][x]
            else:
                qx=m[role].index(x)
                x=qx
            if unknown(role,qx,public):
                counter[(role,qx)]+=1
                if step==0:first[(role,qx)]+=1
    # Grounded ONLY in posterior sample pathways; no secret source oracle.
    return counter,first

def propose(method,public,start,query,seed):
    pred,H,sampled,answers=posterior_answer(public,start,query,seed)
    candidates,frontier=candidate_queries(sampled,public,start,query)
    if method=="query_frontier":
        return max(candidates,key=lambda q:(frontier[q],candidates[q],-q[0],-q[1])) if candidates else None
    if method=="posterior_information_gain":
        prior=entropy(answers)
        scored=[]
        for role,x in candidates:
            by_output={}
            for m,answer in zip(sampled,answers):
                by_output.setdefault(m[role][x],[]).append(answer)
            expected=sum(len(v)/len(answers)*entropy(v) for v in by_output.values())
            gain=prior-expected
            scored.append((gain,candidates[(role,x)],-role,-x,(role,x)))
        return max(scored)[-1] if scored else None
    all_unknown=[(role,x) for role in range(ROLES) for x in range(N) if unknown(role,x,public)]
    if method=="random_eligible":
        return random.Random(seed+2009).choice(all_unknown) if all_unknown else None
    if method=="fixed_acquisition":
        return all_unknown[0] if all_unknown else None
    raise ValueError(method)

class PaidMappingPort:
    def __init__(self,secret,public):
        self._private=secret
        self._known=set((r,x) for r,x,y in public)
        self.receipts=[]
    def buy(self,role,subject):
        if len(self.receipts)>=READ_BUDGET:raise ValueError("paid query ceiling")
        if (role,subject) in self._known:raise ValueError("duplicate paid/public purchase")
        if role not in range(ROLES) or subject not in range(N):raise ValueError("invalid query")
        answer=self._private[role][subject]
        self._known.add((role,subject))
        self.receipts.append((role,subject,answer))
        return answer

def run():
    if set(SEEDS)&set(RESERVED_SEEDS) or set(FAMILIES)&set(RESERVED_FAMILIES):
        raise ValueError("reserved source overlap")
    rows=[]
    for seed in SEEDS:
      for family in FAMILIES:
       for depth in DEPTHS:
        secret,public,start,query=make_case(seed,family,depth)
        truth=walk(secret,start,query)
        if truth!=independently_compose(secret,start,query):raise AssertionError("dual truth drift")
        answers={"privileged_complete_map_oracle_unfair":truth}
        receipts={}
        before={}
        after={}
        for method in METHODS[:-1]:
            current=list(public)
            port=PaidMappingPort(secret,current)
            base_rng=seed*11003+family*307+depth*17
            initial,H,_,_=posterior_answer(current,start,query,base_rng)
            before[method]=H
            for i in range(READ_BUDGET if method!="no_purchase" else 0):
                choice=propose(method,current,start,query,base_rng+i*101)
                if choice is None:break
                role,x=choice
                y=port.buy(role,x)
                current.append((role,x,y))
            prediction,remaining,_,_=posterior_answer(current,start,query,base_rng+909)
            answers[method]=prediction
            after[method]=remaining
            receipts[method]=list(port.receipts)
        rows.append(dict(seed=seed,family=family,depth=depth,truth=truth,start=start,
          query=query,answers=answers,receipts=receipts,
          entropy_before=before,entropy_after=after,
          private_map_not_exposed_to_policies=True,public_count=len(public),
          public_evidence_sha256=hashlib.sha256(json.dumps([public,start,query]).encode()).hexdigest()))
    scores={a:{"correct":sum(r["answers"][a]==r["truth"] for r in rows),"total":len(rows)} for a in METHODS}
    for s in scores.values():s["accuracy"]=s["correct"]/s["total"]
    per_family={str(f):{a:sum(r["answers"][a]==r["truth"] for r in rows if r["family"]==f) for a in METHODS} for f in FAMILIES}
    entropy_delta={a:sum(r["entropy_before"][a]-r["entropy_after"][a] for r in rows)/len(rows) for a in METHODS[:-1]}
    best=max(scores[a]["accuracy"] for a in METHODS[:4])
    source_ready=(best-scores["no_purchase"]["accuracy"]>=.04 and
                 max(entropy_delta[a] for a in METHODS[:4])>0 and
                 scores["privileged_complete_map_oracle_unfair"]["correct"]==len(rows))
    return dict(schema="wingless.f39.paid-relationship-source.trainonly.v1",
       stage="ACTIVE_SOURCE_HEADROOM_NO_LEARNED_ORIGINAL",
       seeds=[SEEDS.start,SEEDS.stop],families=FAMILIES,depths=DEPTHS,
       reserved_seeds=[RESERVED_SEEDS.start,RESERVED_SEEDS.stop],
       reserved_families=RESERVED_FAMILIES,
       questions=len(rows),physical_query_budget=READ_BUDGET,
       posterior_samples=POSTERIOR_SAMPLES,public_mappings_per_role=VISIBLE,
       query_history_bytes_per_case=READ_BUDGET*3,posterior_sample_bytes_external=POSTERIOR_SAMPLES*ROLES*N,
       native_core_memory_not_physically_measured=True,
       per_family_correct=per_family,scores=scores,mean_answer_entropy_reduction_bits=entropy_delta,rows=rows,
       no_learned_policy_deployed=True,source_same_repo_not_F3=True,primary_original_not_executed=True,
       verdict="TRAIN_ONLY_ACTIVE_SOURCE_FEASIBLE" if source_ready else "TRAIN_ONLY_ACTIVE_SOURCE_INSUFFICIENT_HEADROOM",
       accepted_ref_mutation=False,paid_api_calls=0)
if __name__=="__main__":print(json.dumps(run(),indent=2,sort_keys=True))
