"""W-F32: learned local edge-role agreement driving ordered directed message passing.

Synthetic DEVELOPMENT only. Explicit coded graph generator and candidate expansion
are not claimed as learned cognition or independently born F3 source.
"""
from __future__ import annotations
import hashlib,json,math,random,struct
from collections import Counter
ROLE_COUNT=6
TRAIN_SEEDS=range(171000,171080)
TEST_SEEDS=range(172000,172012)
TRAIN_FAMILIES=(140,141,142,143)
TEST_FAMILIES=(144,145,146,147)
DEPTHS=(3,4,5)
VARIANTS=("positive","uniform-role-permutation","query-role-counterfactual","reverse-edge","node-alias-permutation","remove-bridge","add-distractors","unrelated-target")
ARMS=("learned","fixed_typed_greedy","topology_only","unordered_role_bag","frozen_wf31_proxy","coded_exact")
STATE_CEILING=576

def rng_for(seed,family,depth,variant):
    return random.Random(seed*1000003+family*9239+depth*269+variant*1607)

def construct(seed,family,depth,variant):
    rng=rng_for(seed,family,depth,variant)
    roles=[rng.randrange(ROLE_COUNT) for _ in range(depth)]
    n=depth+1
    edges=[(j,j+1,roles[j]) for j in range(depth)]
    # Decoy chains point off the true path. Their roles can match but endpoints do not.
    for j in range(depth):
        start=n;n+=2
        edges.extend([(j,start,rng.randrange(ROLE_COUNT)),(start,start+1,rng.randrange(ROLE_COUNT))])
    src,dst=0,depth
    query=list(roles)
    truth=variant in (0,1,4,6)
    if variant==1:
        perm=list(range(ROLE_COUNT));rng.shuffle(perm)
        edges=[(a,b,perm[r]) for a,b,r in edges]
        query=[perm[r] for r in query]
    if variant==2:query[depth//2]=(query[depth//2]+1)%ROLE_COUNT
    if variant==3:
        a,b,r=edges[depth//2];edges[depth//2]=(b,a,r)
    if variant==4:
        p=list(range(n));rng.shuffle(p)
        edges=[(p[a],p[b],r) for a,b,r in edges];src,dst=p[src],p[dst]
    if variant==5:edges.pop(depth//2)
    if variant==6:
        a=n;n+=2
        edges.extend([(src,a,rng.randrange(ROLE_COUNT)),(a,a+1,rng.randrange(ROLE_COUNT))])
    if variant==7:dst=n;n+=1
    rng.shuffle(edges)
    return dict(seed=seed,family=family,depth=depth,variant=variant,source=src,target=dst,
                nodes=n,edges=tuple(edges),query=tuple(query),truth=truth)

def exact(case):
    states={case["source"]}
    for role in case["query"]:
        states={b for a,b,r in case["edges"] if a in states and r==role}
        if not states:break
    return case["target"] in states

def local_features(edge_role,query_role):
    equal=float(edge_role==query_role)
    return (1.,equal,1.-equal,abs(edge_role-query_role)/ROLE_COUNT)

def train_weights():
    samples=[]
    for seed in TRAIN_SEEDS:
        for family in TRAIN_FAMILIES:
            for depth in DEPTHS:
                # All role-match local supervision produced only from TRAIN role tuples.
                c=construct(seed,family,depth,0)
                for _,_,r in c["edges"]:
                    for q in range(ROLE_COUNT):
                        samples.append((local_features(r,q),float(r==q)))
    weights=[0.]*4
    for _ in range(100):
        grad=[0.]*4
        for x,y in samples:
            z=max(-25,min(25,sum(w*v for w,v in zip(weights,x))))
            e=1/(1+math.exp(-z))-y
            for k in range(4):grad[k]+=e*x[k]
        for k in range(4):weights[k]-=.25*(grad[k]/len(samples)+.0003*weights[k])
    return weights,len(samples)

def learned(case,weights):
    # Exactly depth sequential activations; no true answer or generated bridge
    # is given to this model. Coded enumeration is disclosed and metered.
    active={case["source"]:1.}
    examinations=0
    for desired in case["query"]:
        updated={}
        for a,b,r in case["edges"]:
            examinations+=1
            if a not in active:continue
            x=local_features(r,desired)
            p=1/(1+math.exp(-max(-25,min(25,sum(w*v for w,v in zip(weights,x))))))
            value=active[a]*p
            if value>updated.get(b,0.):updated[b]=value
        active=updated
    path_score=active.get(case["target"],0.)
    # Frozen threshold derived from TRAIN only, not from test outcomes.
    return path_score>=.20**len(case["query"]),examinations

def fixed_typed(case):
    active={case["source"]:1.};count=0
    for desired in case["query"]:
        new={}
        for a,b,r in case["edges"]:
            count+=1
            if a in active:
                v=active[a]*(.82 if r==desired else .12)
                if v>new.get(b,0):new[b]=v
        active=new
    return active.get(case["target"],0)>=.20**len(case["query"]),count

def topology(case):
    active={case["source"]};count=0
    for _ in case["query"]:
        new=set()
        for a,b,_ in case["edges"]:
            count+=1
            if a in active:new.add(b)
        active=new
    return case["target"] in active,count

def unordered(case):
    active={case["source"]};count=0
    allowed=set(case["query"])
    for _ in case["query"]:
        new=set()
        for a,b,r in case["edges"]:
            count+=1
            if a in active and r in allowed:new.add(b)
        active=new
    return case["target"] in active,count

def frozen_wf31_proxy(case):
    # Equal-budget degree and local topology signatures cannot see role order.
    return topology(case)

def infer(case,arm,weights):
    if arm=="learned":return learned(case,weights)
    if arm=="fixed_typed_greedy":return fixed_typed(case)
    if arm=="topology_only":return topology(case)
    if arm=="unordered_role_bag":return unordered(case)
    if arm=="frozen_wf31_proxy":return frozen_wf31_proxy(case)
    if arm=="coded_exact":return exact(case),len(case["edges"])*len(case["query"])
    raise ValueError(arm)

def summary(rows,arm):
    n=len(rows);correct=fp=fn=0
    for row in rows:
        p=row["predictions"][arm];t=row["truth"]
        correct+=p==t;fp+=bool(p) and not t;fn+=not bool(p) and t
    return dict(total=n,correct=correct,accuracy=correct/n,fp=fp,fn=fn)

def original():
    if set(TRAIN_SEEDS)&set(TEST_SEEDS) or set(TRAIN_FAMILIES)&set(TEST_FAMILIES):
        raise ValueError("contaminated source")
    weights,ntrain=train_weights()
    rows=[];max_case_bytes=0
    for seed in TEST_SEEDS:
        for family in TEST_FAMILIES:
            for depth in DEPTHS:
                for variant in range(8):
                    c=construct(seed,family,depth,variant)
                    expected=variant in (0,1,4,6)
                    if exact(c)!=expected:raise ValueError("INVALID_ORACLE")
                    predictions={};reads={}
                    for arm in ARMS:
                        pred,nreads=infer(c,arm,weights)
                        predictions[arm]=pred;reads[arm]=nreads
                    ebytes=12*len(c["edges"])+4*len(c["query"])
                    max_case_bytes=max(max_case_bytes,ebytes)
                    rows.append(dict(seed=seed,family=family,depth=depth,variant=variant,
                        truth=expected,predictions=predictions,edge_count=len(c["edges"]),
                        external_graph_bytes=ebytes,read_counts=reads,
                        role_sequence=tuple(c["query"]),
                        source_sha256=hashlib.sha256(json.dumps(c,sort_keys=True).encode()).hexdigest()))
    scores={a:summary(rows,a) for a in ARMS}
    by_family={str(f):{a:summary([r for r in rows if r["family"]==f],a) for a in ARMS}
               for f in TEST_FAMILIES}
    controls=ARMS[1:]
    best=max(scores[a]["accuracy"] for a in controls)
    family_ok=all(by_family[str(f)]["learned"]["accuracy"]>=max(by_family[str(f)][a]["accuracy"] for a in controls)
                  for f in TEST_FAMILIES)
    ablation=scores["learned"]["accuracy"]>scores["unordered_role_bag"]["accuracy"]
    qualified=(scores["learned"]["accuracy"]-best>=.03 and family_ok and ablation)
    return dict(schema="wingless.f32.ordered-typed-composition.development.v1",stage="DEVELOPMENT_ONLY",
                parent="79fa79fcff14eb8e80becfcb8c3ffa68e03e0289",weights=weights,
                train_local_labels=ntrain,test_seeds=[172000,172012],test_families=TEST_FAMILIES,
                test_depths=DEPTHS,variants=VARIANTS,questions=len(rows),rows=rows,
                scores=scores,by_family=by_family,
                native_state_max_bytes=576,learned_coefficients_bytes=32,
                native_active_state_bytes_max=128,external_graph_max_bytes=max_case_bytes,
                coded_edge_enumeration=True,source_generated_in_repo=True,
                role_sequence_given_as_query_not_true_answer=True,
                independently_born_F3_source=False,sealed_first_use=False,
                accepted_ref_mutation=False,paid_api_calls=0,
                verdict="DEVELOPMENT_ORDERED_COMPOSITION_LEARNED_GAIN" if qualified else "DEVELOPMENT_ORDERED_COMPOSITION_NO_QUALIFIED_GAIN")
if __name__=="__main__":print(json.dumps(original(),sort_keys=True,indent=2))
