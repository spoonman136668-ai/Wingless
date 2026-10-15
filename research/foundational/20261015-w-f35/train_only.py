"""Wingless W-F35: train-only learned demonstration-efficient operator prior.

A meta-learner trains on public transition demonstrations from fresh meta-train
episodes. On disjoint evaluation episodes it chooses between a group-constrained
affine/XOR hypothesis and a permutation MAP fit, composing six-to-ten ordered
transitions. No reserved scientific original or previous pilot is generated.
"""
import json, math, random
from functools import lru_cache

META_SEEDS=range(231000,231024)
META_FAMILIES=(200,201,202)
EVAL_SEEDS=range(232000,232012)
EVAL_FAMILIES=(204,205,206)
RESERVED_SEEDS=range(202000,202016)
RESERVED_FAMILIES=(184,185,186,187)
ROLES=8
STATES=8
DEPTHS=(6,8,10)
ARMS=("learned_meta","uniform_group_map","permutation_map","last_demo","role_blind","privileged_oracle")
def guard(seed,family):
    if not ((seed in META_SEEDS and family in META_FAMILIES) or
            (seed in EVAL_SEEDS and family in EVAL_FAMILIES)):
        raise ValueError("source identities denied, including reserved original")

def generator(seed,family,depth):
    guard(seed,family)
    g=random.Random(seed*1000703+family*17011+depth*1229+991)
    hidden=[]
    for role in range(ROLES):
        if family%3==0:
            a=g.choice([1,3,5,7]);b=g.randrange(8)
            op=tuple((a*x+b)%8 for x in range(8))
        elif family%3==1:
            arr=list(range(8));g.shuffle(arr);op=tuple(arr)
        else:
            a=g.choice([1,3,5,7]);b=g.randrange(8);mask=g.randrange(1,8)
            op=tuple(((a*x+b)%8)^mask for x in range(8))
        hidden.append(op)
    demos=[]
    for role in range(ROLES):
        for k in range(4):
            x=g.randrange(8)
            y=hidden[role][x]
            if g.random()<.12:y=(y+g.randrange(1,8))%8
            demos.append((role,x,y))
    g.shuffle(demos)
    initial=g.randrange(8)
    ordered_roles=tuple(g.randrange(ROLES) for _ in range(depth))
    return tuple(hidden),tuple(demos),initial,ordered_roles

def compose(maps,start,sequence):
    value=start
    for r in sequence:value=maps[r][value]
    return value

def compose_independent(maps,start,sequence):
    values=list(range(8))
    for r in sequence:
        lookup=maps[r]
        values=[lookup[v] for v in values]
    return values[start]

def group_candidates():
    result=set()
    for a in (1,3,5,7):
        for b in range(8):
            result.add(tuple((a*x+b)%8 for x in range(8)))
            for mask in range(1,8):result.add(tuple(((a*x+b)%8)^mask for x in range(8)))
    return tuple(sorted(result))
GROUP=group_candidates()

def assign(scores):
    @lru_cache(None)
    def solve(row, mask):
        if row==8:return 0.0,()
        outcomes=[]
        for target in range(8):
            if (mask>>target)&1:continue
            s,remaining=solve(row+1,mask|(1<<target))
            outcomes.append((s+scores[row][target],(target,)+remaining))
        return max(outcomes,key=lambda p:(p[0],tuple(-x for x in p[1])))
    return solve(0,0)[1]

def infer_role(examples):
    counts=[[0. for _ in range(8)] for _ in range(8)]
    for x,y in examples:counts[x][y]+=1.
    best_perm=assign(counts)
    best_group=max(GROUP,key=lambda op:(sum(op[x]==y for x,y in examples),tuple(-i for i in op)))
    score_group=sum(best_group[x]==y for x,y in examples)
    score_perm=sum(best_perm[x]==y for x,y in examples)
    repeated=len({x for x,y in examples})
    return best_group,best_perm,(1.,float(score_group-score_perm),float(repeated)/4)

def per_role_examples(demos,role):
    return [(x,y) for r,x,y in demos if r==role]

def sigmoid(x):
    return 1/(1+math.exp(-max(-25,min(25,x))))

def train_meta():
    training=[]
    # Train ONLY on public demonstrations: last demo is label for validation
    # of a role model fit to the first three demonstrations.
    for seed in META_SEEDS:
        for family in META_FAMILIES:
            truth,demos,initial,query=generator(seed,family,6)
            for role in range(ROLES):
                records=per_role_examples(demos,role)
                group,perm,features=infer_role(records[:3])
                valx,valy=records[3]
                group_correct=int(group[valx]==valy)
                perm_correct=int(perm[valx]==valy)
                if group_correct==perm_correct:continue
                target=float(group_correct>perm_correct)
                training.append((features,target))
    w=[0.,0.,0.]
    for _ in range(120):
        grad=[0.,0.,0.]
        for x,y in training:
            p=sigmoid(sum(a*b for a,b in zip(w,x)))
            for j in range(3):grad[j]+=(p-y)*x[j]
        for j in range(3):w[j]-=.18*(grad[j]/max(1,len(training))+.002*w[j])
    return w,len(training)

def infer_all(demos,w):
    learned=[];group=[];permutation=[];last=[]
    for role in range(ROLES):
        examples=per_role_examples(demos,role)
        g,p,x=infer_role(examples)
        weight=sigmoid(sum(a*b for a,b in zip(w,x)))
        learned.append(g if weight>=.5 else p)
        group.append(g);permutation.append(p)
        mapping={i:i for i in range(8)}
        for i,j in examples:mapping[i]=j
        last.append(tuple(mapping[i] for i in range(8)))
    return tuple(learned),tuple(group),tuple(permutation),tuple(last)

def score(rows,name):
    n=len(rows);good=sum(r["answers"][name]==r["truth"] for r in rows)
    return dict(correct=good,total=n,accuracy=good/n)

def run():
    assert set(EVAL_SEEDS).isdisjoint(META_SEEDS)
    assert set(EVAL_FAMILIES).isdisjoint(META_FAMILIES)
    assert set(EVAL_SEEDS).isdisjoint(RESERVED_SEEDS)
    assert set(EVAL_FAMILIES).isdisjoint(RESERVED_FAMILIES)
    weights,ntrain=train_meta()
    rows=[]
    for seed in EVAL_SEEDS:
      for family in EVAL_FAMILIES:
       for depth in DEPTHS:
        truth_ops,demos,start,roles=generator(seed,family,depth)
        truth=compose(truth_ops,start,roles)
        assert compose_independent(truth_ops,start,roles)==truth
        learned,group,permutation,last=infer_all(demos,weights)
        answers={
           "learned_meta":compose(learned,start,roles),
           "uniform_group_map":compose(group,start,roles),
           "permutation_map":compose(permutation,start,roles),
           "last_demo":compose(last,start,roles),
           "role_blind":start,
           "privileged_oracle":truth,
        }
        rows.append(dict(seed=seed,family=family,depth=depth,roles=roles,
            truth=truth,answers=answers,public_demo_count=len(demos),
            learner_inference_coefficients=3,external_demo_bytes=len(demos)*3*8))
    scores={arm:score(rows,arm) for arm in ARMS}
    family_scores={str(f):{a:score([r for r in rows if r["family"]==f],a) for a in ARMS} for f in EVAL_FAMILIES}
    fair=ARMS[1:-1]
    gain=scores["learned_meta"]["accuracy"]-max(scores[a]["accuracy"] for a in fair)
    advance=gain>=.03 and all(
        family_scores[str(f)]["learned_meta"]["accuracy"]>=max(family_scores[str(f)][a]["accuracy"] for a in fair)
        for f in EVAL_FAMILIES)
    return dict(schema="wingless.f35.meta-learned-role-grammar.trainonly.v1",stage="TRAIN_ONLY_NO_PRIMARY_ORIGINAL",
        meta_train=[META_SEEDS.start,META_SEEDS.stop],meta_train_families=META_FAMILIES,
        eval_seeds=[EVAL_SEEDS.start,EVAL_SEEDS.stop],eval_families=EVAL_FAMILIES,
        reserved_original_seeds=[RESERVED_SEEDS.start,RESERVED_SEEDS.stop],
        reserved_original_families=RESERVED_FAMILIES,training_validation_examples=ntrain,
        learned_prior_weights=weights,questions=len(rows),scores=scores,family_scores=family_scores,rows=rows,
        inferred_model_bytes=ROLES*STATES+8*3,external_demo_bytes_per_episode=ROLES*4*3*8,
        native_core_cap=576,independent_F3_source_custody=False,primary_original_not_executed=True,
        verdict="TRAIN_ONLY_META_LEARNED_ROLE_GAP_QUALIFIED" if advance else "TRAIN_ONLY_META_LEARNED_ROLE_NO_GAIN",
        accepted_ref_mutation=False,paid_api_calls=0)
if __name__=="__main__":print(json.dumps(run(),sort_keys=True,indent=2))
