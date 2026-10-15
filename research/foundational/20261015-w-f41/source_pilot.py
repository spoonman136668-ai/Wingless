"""W-F41: separate answer-information gain from future answer correctness.

All deployable selectors receive public relationship triples and posterior
completions only. Hidden grammar is accessible exclusively through a two-read
source port and evaluator. The evaluator may report a private hindsight
ceiling but never passes it to policy. New source identities only, no F3.
"""
import json,random,math,hashlib
from collections import Counter

SEEDS=range(311000,311012);FAMILIES=(350,351,352)
RESERVED=range(312000,312016);RESERVED_FAMILIES=(354,355,356,357)
LENGTHS=(3,5);N=8;ROLES=8;COVERAGE=3;READS=2
SAMPLE_BUDGETS=(12,48)
ARMS=("expected_answer_information_gain","expected_final_answer_correctness",
      "query_frontier","fixed_acquisition","no_purchase","private_hindsight_oracle_unfair")
def source(seed,family,length):
    if seed not in SEEDS or family not in FAMILIES or length not in LENGTHS:
        raise ValueError("unregistered and reserved source identities forbidden")
    r=random.Random(seed*19997+family*10093+length*491+437)
    maps=[]
    for role in range(ROLES):
        if family==350:
            a=r.choice((1,3,5,7));b=r.randrange(8)
            mapping=tuple((a*x+b)%8 for x in range(N))
        elif family==351:
            arr=list(range(N));r.shuffle(arr);mapping=tuple(arr)
        else:
            a=r.choice((1,3,5,7));b=r.randrange(8);mask=r.randrange(1,8)
            mapping=tuple(((a*x+b)%8)^mask for x in range(N))
        maps.append(mapping)
    public=[]
    for role in range(ROLES):
        for x in r.sample(range(N),COVERAGE):
            public.append((role,x,maps[role][x]))
    r.shuffle(public)
    start=r.randrange(N)
    path=tuple((r.randrange(ROLES),r.choice((-1,1))) for _ in range(length))
    return tuple(maps),tuple(public),start,path
def follow(maps,start,path):
    v=start
    for role,orientation in path:
        v=maps[role][v] if orientation==1 else maps[role].index(v)
    return v
def independent_follow(maps,start,path):
    p=[int(i==start) for i in range(N)]
    for role,orientation in path:
        new=[0]*N
        for i,weight in enumerate(p):
            j=maps[role][i] if orientation==1 else maps[role].index(i)
            new[j]+=weight
        p=new
    assert sum(p)==1
    return p.index(1)
def completions(public,samples,seed):
    known=[{} for _ in range(ROLES)]
    for role,x,y in public:known[role][x]=y
    if any(len(set(d.values()))!=len(d) for d in known):raise ValueError("bad visible bijection")
    r=random.Random(seed)
    bank=[]
    for _ in range(samples):
        worlds=[]
        for d in known:
            x=[i for i in range(N) if i not in d]
            y=[i for i in range(N) if i not in d.values()]
            r.shuffle(y);new=dict(d);new.update(zip(x,y))
            worlds.append(tuple(new[i] for i in range(N)))
        bank.append(tuple(worlds))
    return bank
def entropy(v):
    counts=Counter(v);n=len(v)
    return -sum((c/n)*math.log2(c/n) for c in counts.values() if c)
def candidates(public,bank,start,path):
    observed={(r,x) for r,x,y in public}
    frontier=Counter()
    for world in bank:
        state=start
        for role,direction in path:
            subject=state if direction==1 else world[role].index(state)
            if (role,subject) not in observed:frontier[(role,subject)]+=1
            state=world[role][state] if direction==1 else subject
    other=[(r,x) for r in range(ROLES) for x in range(N) if (r,x) not in observed]
    return other,frontier
def score_actions(public,start,path,samples,seed):
    bank=completions(public,samples,seed)
    outcomes=[follow(w,start,path) for w in bank]
    allowed,frontier=candidates(public,bank,start,path)
    prior_entropy=entropy(outcomes)
    original_best=max(Counter(outcomes).values())/len(outcomes)
    stats={}
    for r,x in allowed:
        groups={}
        for w,a in zip(bank,outcomes):groups.setdefault(w[r][x],[]).append(a)
        post_entropy=sum(len(v)/len(outcomes)*entropy(v) for v in groups.values())
        post_bayes_accuracy=sum(max(Counter(v).values()) for v in groups.values())/len(outcomes)
        stats[(r,x)]=(prior_entropy-post_entropy,
                       post_bayes_accuracy-original_best,frontier[(r,x)])
    return stats,original_best,prior_entropy,bank,outcomes
class Port:
    def __init__(self,maps,public):
        self.__maps=maps
        self.known={(r,x) for r,x,y in public}
        self.receipts=[]
    def buy(self,r,x):
        if len(self.receipts)>=READS or (r,x) in self.known:raise ValueError("third/duplicate purchase")
        value=self.__maps[r][x]
        self.known.add((r,x));self.receipts.append((r,x,value))
        return value
def select(method,public,start,path,seed):
    if method=="fixed_acquisition":
        known={(r,x) for r,x,y in public}
        return next(((r,x) for r in range(ROLES) for x in range(N) if (r,x) not in known),None)
    stats,accuracy,entropy_before,bank,answers=score_actions(public,start,path,48,seed)
    if not stats:return None
    if method=="expected_answer_information_gain":
        return max(stats,key=lambda q:(stats[q][0],stats[q][2],-q[0],-q[1]))
    if method=="expected_final_answer_correctness":
        return max(stats,key=lambda q:(stats[q][1],stats[q][2],-q[0],-q[1]))
    if method=="query_frontier":
        return max(stats,key=lambda q:(stats[q][2],-q[0],-q[1]))
    raise ValueError("unknown selector")
def majority_prediction(public,start,path,seed):
    bank=completions(public,48,seed)
    counts=Counter(follow(w,start,path) for w in bank)
    return max(range(N),key=lambda x:(counts[x],-x))
def run():
    assert not set(SEEDS)&set(RESERVED)
    assert not set(FAMILIES)&set(RESERVED_FAMILIES)
    rows=[];stable=agreements=0
    for seed in SEEDS:
      for fam in FAMILIES:
       for length in LENGTHS:
        secret,initial,start,path=source(seed,fam,length)
        truth=follow(secret,start,path)
        assert truth==independent_follow(secret,start,path)
        small,_,_,_,_=score_actions(initial,start,path,12,seed*17+fam)
        large,_,_,_,_=score_actions(initial,start,path,48,seed*17+fam)
        entropy_first=max(small,key=lambda q:(small[q][0],small[q][2],-q[0],-q[1]))
        entropy_high=max(large,key=lambda q:(large[q][0],large[q][2],-q[0],-q[1]))
        decision_first=max(small,key=lambda q:(small[q][1],small[q][2],-q[0],-q[1]))
        decision_high=max(large,key=lambda q:(large[q][1],large[q][2],-q[0],-q[1]))
        stable+=entropy_first==entropy_high;agreements+=decision_high==entropy_high
        answers={"private_hindsight_oracle_unfair":truth}
        receipts={}
        for arm in ARMS[:-1]:
            observed=list(initial)
            port=Port(secret,initial)
            for step in range(READS if arm!="no_purchase" else 0):
                choice=select(arm,observed,start,path,seed*31+fam+step*7)
                if choice is None:break
                r,x=choice
                y=port.buy(r,x);observed.append((r,x,y))
            answers[arm]=majority_prediction(observed,start,path,seed*31+fam+607)
            receipts[arm]=list(port.receipts)
        rows.append(dict(seed=seed,family=fam,length=length,truth=truth,answers=answers,receipts=receipts,
            public_sha256=hashlib.sha256(json.dumps([initial,start,path]).encode()).hexdigest(),
            small_entropy_argmax=entropy_first,large_entropy_argmax=entropy_high,
            small_decision_argmax=decision_first,large_decision_argmax=decision_high,
            high_budget_estimated_decision_value=large[decision_high][1],
            high_budget_estimated_information_gain=large[entropy_high][0]))
    scores={arm:dict(correct=sum(r["answers"][arm]==r["truth"] for r in rows),total=len(rows)) for arm in ARMS}
    for a in scores:scores[a]["accuracy"]=scores[a]["correct"]/len(rows)
    by_family={str(f):{a:sum(r["answers"][a]==r["truth"] for r in rows if r["family"]==f) for a in ARMS} for f in FAMILIES}
    disagreement=1-agreements/len(rows)
    stability=stable/len(rows)
    fair=ARMS[:-1]
    gain=max(scores[a]["accuracy"] for a in fair[:-1])-scores["no_purchase"]["accuracy"]
    useful=gain>=.03 and scores["private_hindsight_oracle_unfair"]["correct"]==len(rows)
    return dict(schema="wingless.f41.posterior-action-identifiability.trainonly.v1",
        stage="FRESH_SOURCE_DIAGNOSTIC_NOT_LEARNED_ORIGINAL",
        seeds=[SEEDS.start,SEEDS.stop],families=FAMILIES,lengths=LENGTHS,
        reserved_original_seeds=[RESERVED.start,RESERVED.stop],
        reserved_original_families=RESERVED_FAMILIES,
        questions=len(rows),physical_purchase_cap=READS,
        posterior_sample_budgets=SAMPLE_BUDGETS,
        posterior_external_memory_bytes_per_48_sample_bank=48*ROLES*N,
        entropy_action_rank_stability_12_vs_48=stability,
        entropy_vs_decision_action_disagreement_rate_48=disagreement,
        scores=scores,by_family=by_family,rows=rows,
        verdict="TRAIN_ONLY_DECISION_IDENTIFIABILITY_HEADROOM" if useful else "TRAIN_ONLY_DECISION_IDENTIFIABILITY_INSUFFICIENT",
        evaluator_private_true_maps_never_given_to_selector=True,
        source_independent_F3_custody=False,primary_original_not_executed=True,
        accepted_ref_mutation=False,paid_api_calls=0)
if __name__=="__main__":print(json.dumps(run(),sort_keys=True,indent=2))
