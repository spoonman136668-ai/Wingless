"""W-F53: FIT-supervised value of intermediate relational belief states.

Latent state-vs-future-task supervision comes from FIT labels only. Diagnostic
action receives public role observations/calibration, recurrent 7-value belief
and its train-frozen 8-bit state-scoring weights, never any private truth.
All synthetic; operations are abstract, NOT measured native 576B residency.
"""
from __future__ import annotations
import itertools,json,math,random
from pathlib import Path
FIT=range(561000,561016)
DIAG=range(561100,561116)
FIT_FAM=(980,981,982,983)
DIAG_FAM=(984,985,986,987)
RESERVED=range(562000,562016)
ARMS=("learned_task_value_two_beam","fixed_mass_two_beam",
      "fixed_entropy_two_beam","budget_matched_symbolic",
      "public_context_parity","public_support_hamming",
      "full_exact_public_bayes","private_truth_oracle_unfair")
N=7
DIM=8
CONTEXTS=tuple(itertools.product((0,1),repeat=3))
TEMPLATES=(
 ((0,1),(2,-1),(1,1),(3,1),(0,-1),(3,-1),(2,1),(1,-1),(0,1),(2,-1)),
 ((3,-1),(1,1),(2,1),(1,-1),(0,1),(3,1),(1,1),(2,-1),(0,1),(3,-1)),
 ((2,1),(0,-1),(3,1),(2,-1),(1,1),(0,1),(3,-1),(2,1),(1,-1),(0,1)),
 ((1,-1),(3,1),(0,-1),(2,1),(3,-1),(1,1),(2,-1),(0,1),(3,1),(1,-1)),
 ((0,1),(1,-1),(3,-1),(2,1),(0,-1),(2,-1),(1,1),(3,1),(0,1),(1,-1)),
 ((2,-1),(0,1),(3,1),(1,1),(2,1),(3,-1),(0,-1),(1,-1),(2,1),(3,1)),
 ((3,1),(2,1),(0,-1),(3,-1),(1,1),(0,1),(2,-1),(1,-1),(3,1),(2,1)),
 ((1,1),(3,-1),(2,-1),(0,1),(1,-1),(3,1),(0,-1),(2,1),(3,-1),(1,1)))
def law(fam,c):
    a,b,d=c
    k=(fam-980)%8
    values=(a^b^d,int(a+b+d>=2),a^(b&d),int((not a) or b),
            (a|b)^d,int(a+b*2+d>=3),int(a+b+d==1),b^(a&d))
    return values[k]
def source(seed,fam,length):
    if not ((seed in FIT and fam in FIT_FAM) or
            (seed in DIAG and fam in DIAG_FAM)):raise ValueError("wrong source identity")
    if length not in ((6,8) if seed in FIT else (8,10)):
        raise ValueError("unapproved chain length")
    rng=random.Random(seed*25121+fam*4339+length*101+53003)
    context=CONTEXTS[rng.randrange(8)]
    noise=.12 if seed in FIT else .23
    sign=-1 if law(fam,context)^(rng.random()<noise) else 1
    shifts=tuple(rng.randrange(7) for i in range(4))
    obs=[]
    for role in range(4):
        x=rng.randrange(7);y=(sign*x+shifts[role])%7
        if rng.random()<noise:y=(y+rng.randrange(1,7))%7
        obs.append((role,x,y))
    calib=[]
    for c in rng.sample(CONTEXTS,4):
        cal_sign=-1 if law(fam,c)^(rng.random()<noise) else 1
        b=rng.randrange(7);x=rng.randrange(7)
        z=(x+rng.randrange(1,7))%7
        calib.append((c,(x,(cal_sign*x+b)%7),(z,(cal_sign*z+b)%7)))
    chain=TEMPLATES[fam-980][:length]
    # Target is evaluated strictly on hidden role maps; the learner cannot
    # see exact hidden shifts or true intermediate states at diagnostic time.
    at=rng.randrange(7);start=at;states=[]
    for r,d in chain:
        at=(sign*at+shifts[r])%7 if d==1 else (sign*(at-shifts[r]))%7
        states.append(at)
    return dict(seed=seed,family=fam,length=length,context=context,
      sign=sign,shifts=shifts,observations=tuple(obs),calibration=tuple(calib),
      chain=tuple(chain),start=start,truth=at,hidden_intermediates=tuple(states))

def anchor_sign(a):
    c,(x,y),(z,t)=a
    slope=((t-y)*pow((z-x)%7,-1,7))%7
    assert slope in (1,6)
    return -1 if slope==6 else 1

def minus_probability(d):
    scores=[]
    for item in d["calibration"]:
        wt=math.exp(-sum(x!=y for x,y in zip(d["context"],item[0])))
        scores.append((wt,anchor_sign(item)))
    return sum(w for w,s in scores if s==-1)/sum(w for w,s in scores)

def likelihoods(d,sign):
    noise=.12 if d["seed"] in FIT else .23
    # Diagnoser is allowed the preregistered noise rate per public dataset,
    # but not any hidden source sign, shift or label.
    out={}
    for role,x,y in d["observations"]:
        probs=tuple((1-noise if (sign*x+b)%7==y else noise/6)
                    for b in range(7))
        assert abs(sum(probs)-1)<1e-12
        out[role]=probs
    return out

def move(state,role,direction,sign,likelihood,active,quantized):
    nxt=[0.]*7;work=0
    for x in active:
        for shift in range(7):
            dest=(sign*x+shift)%7 if direction==1 else (sign*(x-shift))%7
            nxt[dest]+=state[x]*likelihood[role][shift]
            work+=1
    if quantized:
        ints=[min(255,max(0,round(v*255))) for v in nxt]
        if not any(ints):ints[max(range(7),key=lambda i:nxt[i])]=1
        total=sum(ints)
        nxt=[i/total for i in ints]
        work+=7
    else:
        total=sum(nxt)
        if total:nxt=[a/total for a in nxt]
    return nxt,work

def features(state,index,at,chain):
    val=state[at]
    hi=max(state)
    ent=-sum(x*math.log(x+1e-15) for x in state)
    unused=chain[index:]
    repeat=sum(r==chain[index][0] for r,d in unused)/max(len(unused),1)
    return (1.,val, val/hi if hi else 0.,
      (at-3)/3,ent/math.log(7),index/max(len(chain),1),
      repeat, (val-max(x for i,x in enumerate(state) if i!=at)))

def pick(state,step,chain,mode,weights):
    if mode=="learned_task_value_two_beam":
        # The public mass filter allocates exactly three task-value scores.
        # Every control is conservatively charged the same selection work.
        eligible=sorted(range(7),key=lambda i:(-state[i],i))[:3]
        ranked=sorted(eligible,key=lambda i:(-sum(x*y for x,y in zip(weights,features(state,step,i,chain))),-state[i],i))
    elif mode=="fixed_entropy_two_beam":
        top=max(range(7),key=lambda i:(state[i],-i))
        ranked=[top]+sorted((i for i in range(7) if i!=top),
                    key=lambda i:(-state[i]*min((i-top)%7,(top-i)%7),i))
    else:
        ranked=sorted(range(7),key=lambda i:(-state[i],i))
    return tuple(ranked[:2])

def task_belief(d,mode,sign,weights):
    probs=likelihoods(d,sign)
    state=[0.]*7;state[d["start"]]=1.
    work=7*4+7*4
    history=[]
    for step,(role,direction) in enumerate(d["chain"]):
        if step==0:active=(d["start"],)
        else:
            active=pick(state,step,d["chain"],mode,weights)
            work+=3*DIM
        state,spend=move(state,role,direction,sign,probs,active,True)
        work+=spend
        history.append(tuple(state))
    return state,work,history

def exact_belief(d):
    p=minus_probability(d)
    all_states=[];spent=7*4*2+7*4*2
    for sign in (1,-1):
        likelihood=likelihoods(d,sign)
        state=[0.]*7;state[d["start"]]=1.
        for role,direction in d["chain"]:
            state,c=move(state,role,direction,sign,likelihood,range(7),False)
            spent+=c
        all_states.append(state)
    return [(1-p)*a+p*b for a,b in zip(*all_states)],spent+9*4

def trained_model():
    examples=[]
    for seed in FIT:
      for fam in FIT_FAM:
       for length in (6,8):
        d=source(seed,fam,length)
        sign=-1 if minus_probability(d)>.5 else 1
        state=[0.]*7;state[d["start"]]=1.
        prior=likelihoods(d,sign)
        for step,(r,direction) in enumerate(d["chain"]):
            state,_=move(state,r,direction,sign,prior,range(7),False)
            truth=d["hidden_intermediates"][step]
            for target in range(7):
                examples.append((features(state,step,target,d["chain"]),float(target==truth)))
    # Real supervised gradient optimization from only FIT private outcomes,
    # frozen before any diagnostic case is generated.
    w=[0.]*DIM
    for epoch in range(16):
        for f,y in examples:
            z=sum(a*b for a,b in zip(w,f))
            prob=1/(1+math.exp(-max(-15,min(15,z))))
            factor=6 if y else 1
            for k in range(DIM):
                w[k]-=.009*((prob-y)*factor*f[k]+.0003*w[k])
    quant=tuple(min(127,max(-127,int(round(wk*24)))) for wk in w)
    return quant,len(examples)

def run_arm(d,arm,weights):
    if arm=="full_exact_public_bayes":
        prob,work=exact_belief(d)
        return max(range(7),key=lambda i:(prob[i],-i)),work,14
    if arm=="private_truth_oracle_unfair":return d["truth"],0,0
    if arm=="public_context_parity":sign=-1 if sum(d["context"])%2 else 1
    elif arm=="public_support_hamming":
        first=min(d["calibration"],key=lambda a:(sum(x!=y for x,y in zip(d["context"],a[0])),a[0]))
        sign=anchor_sign(first)
    else:sign=-1 if minus_probability(d)>=.5 else 1
    if arm=="budget_matched_symbolic":
        mode="fixed_mass_two_beam"
    elif arm in ("public_context_parity","public_support_hamming"):
        mode="fixed_mass_two_beam"
    else:mode=arm
    state,work,_=task_belief(d,mode,sign,weights)
    return max(range(7),key=lambda i:(state[i],-i)),work,7

def run():
    assert not set(FIT)&set(DIAG) and not set(FIT_FAM)&set(DIAG_FAM)
    assert not set(DIAG)&set(RESERVED)
    weights,ntrain=trained_model()
    rows=[];work={arm:0 for arm in ARMS[:-1]}
    for seed in DIAG:
      for fam in DIAG_FAM:
       for length in (8,10):
        d=source(seed,fam,length)
        answers={};costs={};memory={}
        for arm in ARMS:
            a,c,b=run_arm(d,arm,weights)
            answers[arm]=a
            if arm!=ARMS[-1]:
                costs[arm]=c;memory[arm]=b
                work[arm]+=c
        rows.append(dict(seed=seed,family=fam,length=length,
           public_context=list(d["context"]),
           public_observations=[list(x) for x in d["observations"]],
           public_calibration=[[list(a[0]),list(a[1]),list(a[2])] for a in d["calibration"]],
           chain=[list(a) for a in d["chain"]],start=d["start"],
           truth=d["truth"],answers=answers,costs=costs,model_state_bytes=memory))
    scores={arm:dict(correct=sum(r["answers"][arm]==r["truth"] for r in rows),total=len(rows))
            for arm in ARMS}
    for v in scores.values():v["accuracy"]=v["correct"]/v["total"]
    per_family={str(f):{arm:sum(r["answers"][arm]==r["truth"] for r in rows if r["family"]==f)
                 for arm in ARMS} for f in DIAG_FAM}
    fair=("fixed_mass_two_beam","fixed_entropy_two_beam","budget_matched_symbolic",
          "public_context_parity","public_support_hamming")
    strongest=max(fair,key=lambda a:scores[a]["accuracy"])
    gain=scores[ARMS[0]]["accuracy"]-scores[strongest]["accuracy"]
    ratio=work[ARMS[0]]/work["full_exact_public_bayes"]
    protected=all(per_family[str(f)][ARMS[0]]>=max(per_family[str(f)][a] for a in fair)
                  for f in DIAG_FAM)
    passed=gain>=.05 and ratio<=.5 and protected
    return dict(schema="wingless.f53.future-value-beam.trainonly.v1",
     stage="TRAIN_ONLY_LEARNED_BELIEF_PRUNING_NO_ORIGINAL",
     fit_seeds=[FIT.start,FIT.stop],fit_families=list(FIT_FAM),
     diagnostic_seeds=[DIAG.start,DIAG.stop],diagnostic_families=list(DIAG_FAM),
     proposed_original_NOT_RESERVED=[RESERVED.start,RESERVED.stop],
     train_examples=ntrain,quantized_fit_weights=list(weights),
     model_parameter_bytes=DIM,recurrent_belief_state_bytes=7,
     selected_state_index_bytes=2,
     native_576B_total_unverified=True,cases=len(rows),rows=rows,
     scores=scores,by_family=per_family,posterior_abstract_operations=work,
     strongest_equal_work_fair=strongest,candidate_gain=gain,
     posterior_operation_ratio_vs_exact=ratio,
     family_nondegradation=protected,frozen_gate_met=passed,
     verdict="TRAIN_ONLY_LEARNED_TASK_VALUE_SUPPORTED" if passed else
             "TRAIN_ONLY_LEARNED_TASK_VALUE_NEGATIVE",
     independent_F3=False,original_executed=False,accepted_ref_mutation=False,paid_api_calls=0)
if __name__=="__main__":print(json.dumps(run(),sort_keys=True,indent=2))
