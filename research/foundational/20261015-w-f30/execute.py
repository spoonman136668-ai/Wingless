"""W-F30: train-only structural-trajectory scoring across anonymized graph encounters.

Development synthetic study, NOT independent F3 source custody. No node identifiers
or correct candidate location are given to the scoring function. The generator,
candidate enumeration, motif statistics, and training procedure are programmed.
"""
import hashlib
import json
import math
import random
from collections import Counter
from dataclasses import dataclass

TRAIN_SEEDS=range(161000,161080)
TEST_SEEDS=range(162500,162516)
TRAIN_FAMILIES=(120,121,122,123)
TEST_FAMILIES=(124,125,126,127)
CONDITIONS=("role-permutation","edge-dropout","bridge-rewiring")
ARMS=("learned","canonical","jaccard","two_of_three","fixed_trajectory","no_link","privileged_oracle")
NUM_CANDIDATES=8
CORE_BUDGET=576
FEATURE_COUNT=12
MASK=0xffffffff
def stream(seed,family,condition,serial):
    return random.Random(seed*1000003+family*10007+condition*431+serial*53)

def graph(seed,family,serial,test):
    rng=stream(seed,family,0,serial)
    n=9 if test else 8
    # Different structural families are different low-density graph processes.
    edges=set()
    for i in range(n):
        edges.add((min(i,(i+1)%n),max(i,(i+1)%n)))
    p=.12+.025*(family%4)
    for a in range(n):
        for b in range(a+2,n):
            if (a,b) in edges:continue
            if rng.random()<p:edges.add((a,b))
    if family%2:
        for b in range(2,n,3):edges.add((0,b))
    # Edge role labels are source-generated (not identities).
    return n,tuple(sorted((a,b,(a*7+b*5+family+seed+serial)%(6 if test else 4)) for a,b in edges))

def observed(source,seed,family,condition,encounter,truth,test):
    n,edges=source
    rng=stream(seed,family,condition+7,encounter+int(not truth)*1000)
    perm=list(range(n));rng.shuffle(perm)
    rolemap=list(range(6 if test else 4));rng.shuffle(rolemap)
    # Source and destination use different anon node IDs in each episode.
    result=[]
    for a,b,r in edges:
        if condition==1 and rng.random()<.16:continue
        if condition==2 and rng.random()<.12:b=(b+2)%n
        if a==b:continue
        pa,pb=sorted((perm[a],perm[b]))
        result.append((pa,pb,rolemap[r]))
    return n,tuple(sorted(set(result)))

def motif(g):
    n,edges=g
    deg=[0]*n;roles=Counter();pairs=set()
    for a,b,r in edges:
        deg[a]+=1;deg[b]+=1;roles[r]+=1;pairs.add((a,b))
    degree=sorted(deg)
    # No raw node IDs, no canonical semantic labels.
    triangles=sum((min(a,b),max(a,b)) in pairs and
                  (min(a,c),max(a,c)) in pairs and
                  (min(b,c),max(b,c)) in pairs
                  for a in range(n) for b in range(a+1,n) for c in range(b+1,n))
    wedges=sum(d*(d-1)//2 for d in deg)
    hist=tuple(sum(d==k for d in degree) for k in range(5))
    return (len(edges),triangles,wedges,*hist, max(deg),sum(d*d for d in deg))

def features(ref,candidate):
    # Compare two independent graph encounters using isomorphism-invariant
    # graph counts rather than static witness similarity.
    a=motif(ref);b=motif(candidate)
    d=[abs(x-y) for x,y in zip(a,b)]
    den=[1+max(abs(x),abs(y)) for x,y in zip(a,b)]
    x=[float(z)/v for z,v in zip(d,den)]
    x.append(sum(z>0 for z in d)/len(d))
    return (1.0,)+tuple(x)

def make_case(seed,family,condition,test):
    rng=stream(seed,family,condition,500)
    base=graph(seed,family,10,test)
    query=observed(base,seed,family,condition,0,True,test)
    correct=rng.randrange(NUM_CANDIDATES)
    candidates=[]
    for j in range(NUM_CANDIDATES):
        if j==correct:node=base
        else:
            node=graph(seed+3000+j*71,family,20+j,test)
            if j%3==0:
                # Hard negative: topology derived from true graph, but
                # structurally rewired rather than a semantically true join.
                n,edges=base;ed=list(edges)
                if ed:
                    k=j%len(ed);a,b,r=ed[k]
                    ed[k]=((a+1)%n,(b+3)%n,r)
                    fixed={(min(x,y),max(x,y),role) for x,y,role in ed if x!=y}
                    node=(n,tuple(sorted(fixed)))
        candidates.append(observed(node,seed,family,condition,1+j,j==correct,test))
    return query,candidates,correct

def train():
    samples=[]
    for seed in TRAIN_SEEDS:
        for fam in TRAIN_FAMILIES:
            for cond in range(3):
                ref,candidates,truth=make_case(seed,fam,cond,False)
                for idx,cand in enumerate(candidates):
                    samples.append((features(ref,cand),float(idx==truth)))
    weights=[0.0]*FEATURE_COUNT
    for epoch in range(140):
        grad=[0.0]*FEATURE_COUNT
        for x,y in samples:
            z=sum(w*v for w,v in zip(weights,x))
            z=max(-20,min(20,z))
            error=1/(1+math.exp(-z))-y
            for j,v in enumerate(x):grad[j]+=error*v
        for j in range(FEATURE_COUNT):weights[j]-=.24*(grad[j]/len(samples)+.0003*weights[j])
    return weights,len(samples)

def classify(ref,candidates,weights,mode):
    if mode=="no_link":return -1
    score=[]
    ref_m=motif(ref)
    for cand in candidates:
        x=features(ref,cand)
        if mode=="learned":s=sum(a*b for a,b in zip(x,weights))
        elif mode=="canonical":s=-sum(abs(a-b) for a,b in zip(ref_m,motif(cand)))
        elif mode=="jaccard":s=-sum(x[1:])  # equal-evidence normalized motif distance
        elif mode=="two_of_three":s=-sum(int(z>0.15) for z in x[1:4])
        elif mode=="fixed_trajectory":s=-sum(abs(a-b) for a,b in zip(ref_m[:3],motif(cand)[:3]))
        else:raise ValueError(mode)
        score.append(s)
    return max(range(len(score)),key=lambda i:(score[i],-i))

def accuracy(rows,arm):
    n=len(rows);correct=sum(r["predictions"][arm]==r["truth_index"] for r in rows)
    return dict(correct=correct,total=n,accuracy=correct/n)

def original():
    weights,ntrain=train()
    rows=[]
    for seed in TEST_SEEDS:
        for fam in TEST_FAMILIES:
            for condition in range(3):
                ref,cands,truth=make_case(seed,fam,condition,True)
                pred={arm:(truth if arm=="privileged_oracle" else classify(ref,cands,weights,arm)) for arm in ARMS}
                rows.append(dict(seed=seed,family=fam,condition=CONDITIONS[condition],
                    truth_index=truth,predictions=pred,
                    candidate_count=NUM_CANDIDATES,query_motif=motif(ref),
                    candidate_motifs=[motif(x) for x in cands],
                    pair_feature_sha256=hashlib.sha256(json.dumps([features(ref,x) for x in cands]).encode()).hexdigest()))
    scores={arm:accuracy(rows,arm) for arm in ARMS}
    by_family={str(f):{a:accuracy([r for r in rows if r["family"]==f],a) for a in ARMS} for f in TEST_FAMILIES}
    best=max(scores[a]["accuracy"] for a in ARMS[1:-1])
    family_ok=all(by_family[str(f)]["learned"]["accuracy"] >= max(by_family[str(f)][a]["accuracy"] for a in ARMS[1:-1]) for f in TEST_FAMILIES)
    advanced=(scores["learned"]["accuracy"]-best>=.03 and family_ok)
    return dict(schema="wingless.f30.structural-trajectory-development.v1",stage="DEVELOPMENT_ONLY",
      parent="f9f322dfc4760868491be67a5ba050ebcd4a690a",native_core_budget_bytes=CORE_BUDGET,
      actual_coefficient_bytes=len(weights)*8,external_encounter_storage_metered_separately=True,
      external_encounter_values_per_query_upper_bound=2*NUM_CANDIDATES*9*9,training_pairs=ntrain,training_epoch_count=140,
      weights=weights,train_seeds=[161000,161080],test_seeds=[162500,162516],
      train_families=TRAIN_FAMILIES,test_families=TEST_FAMILIES,conditions=CONDITIONS,
      all_candidates_enumerated_without_truth=True,truth_not_given_to_scorer=True,
      questions=len(rows),scores=scores,by_family=by_family,rows=rows,
      verdict="DEVELOPMENT_STRUCTURAL_TRANSFER_LEARNED_GAIN" if advanced else "DEVELOPMENT_STRUCTURAL_TRANSFER_NO_QUALIFIED_GAIN",
      independently_born_F3_source=False,sealed_first_use=False,accepted_ref_mutation=False,paid_api_calls=0)

if __name__=="__main__":print(json.dumps(original(),sort_keys=True,indent=2))
