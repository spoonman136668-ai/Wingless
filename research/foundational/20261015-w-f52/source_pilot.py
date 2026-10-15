"""W-F52: quantized seven-state compositional inference with learned sign.

True relations are noisy-observed affine maps over Z7. Exact public Bayes
maintains TWO sign-conditional belief states; compressed approximation learns
one sign from FIT source only and propagates a quantized seven-value belief.
All policies see the same public observations/calibration. Operation accounting
counts identical posterior transitions, likelihood comparisons and scoring.
Pure-Python implementation is NOT physically qualified native 576B.
"""
from __future__ import annotations
import itertools,json,math,random
from collections import Counter
FIT=range(541000,541016)
DIAG=range(541100,541116)
FIT_FAM=(940,941,942,943)
DIAG_FAM=(944,945,946,947)
RESERVED=range(542000,542016)
N=7
ARMS=("trained_quantized_state","budget_matched_symbolic_approx",
      "full_exact_public_posterior","public_three_bit_rule",
      "public_support_hamming","no_observation",
      "private_truth_oracle_unfair")
LAWS=("parity","majority","xnor","exactly_one",
      "xor_and","or_not","threshold_two","impl")
PATTERNS=(
 ((0,1),(1,-1),(2,1),(3,1),(0,-1),(1,1),(3,-1),(2,-1)),
 ((2,-1),(3,1),(0,1),(2,1),(3,-1),(0,1),(1,-1),(3,1)),
 ((1,1),(3,-1),(2,1),(0,-1),(1,-1),(3,1),(2,-1),(0,1)),
 ((3,1),(0,1),(3,-1),(1,1),(2,1),(0,-1),(3,1),(1,-1)),
 ((0,-1),(2,1),(1,-1),(3,1),(2,-1),(1,1),(0,1),(3,-1)),
 ((3,-1),(1,-1),(0,1),(2,1),(3,-1),(2,1),(1,1),(0,-1)),
 ((2,1),(0,-1),(3,1),(1,-1),(2,1),(0,1),(3,-1),(1,1)),
 ((1,-1),(3,1),(2,-1),(0,1),(3,1),(1,-1),(2,1),(0,-1)),
)
CONTEXTS=tuple(itertools.product((0,1),repeat=3))

def law(label,c):
    a,b,d=c
    if label=="parity":return a^b^d
    if label=="majority":return int(a+b+d>=2)
    if label=="xnor":return 1^(a^b^d)
    if label=="exactly_one":return int(a+b+d==1)
    if label=="xor_and":return a^(b&d)
    if label=="or_not":return a|(b&(1^d))
    if label=="threshold_two":return int(a+2*b+d>=3)
    if label=="impl":return int((not a) or (b and d))
    raise ValueError("unknown law")

def source(seed,fam,length):
    if not ((seed in FIT and fam in FIT_FAM) or
            (seed in DIAG and fam in DIAG_FAM)) or length not in (4,6,8):
        raise ValueError("unregistered source or future original")
    if seed in FIT and length not in (4,6):raise ValueError("fit length mismatch")
    if seed in DIAG and length not in (6,8):raise ValueError("diagnostic length mismatch")
    rng=random.Random(seed*21169+fam*8093+length*487+52009)
    ctx=CONTEXTS[rng.randrange(8)]
    sign=-1 if law(LAWS[fam-940],ctx)^(rng.random()<.17) else 1
    shifts=[rng.randrange(N) for _ in range(4)]
    observations=[]
    for r in range(4):
        x=rng.randrange(N)
        y=(sign*x+shifts[r])%N
        if rng.random()<.17:y=(y+rng.randrange(1,N))%N
        observations.append((r,x,y))
    calib=[]
    for cal_ctx in rng.sample(CONTEXTS,4):
        cal_sign=-1 if law(LAWS[fam-940],cal_ctx)^(rng.random()<.13) else 1
        shift=rng.randrange(N)
        x=rng.randrange(N)
        y=(x+rng.randrange(1,N))%N
        calib.append((cal_ctx,(x,(cal_sign*x+shift)%N),
                      (y,(cal_sign*y+shift)%N)))
    pattern=PATTERNS[fam-940]
    # Heldout role topology differs substantially from the fit orientation.
    chain=pattern[:length]
    start=rng.randrange(N)
    at=start
    for r,d in chain:
        at=(sign*at+shifts[r])%N if d==1 else (sign*(at-shifts[r]))%N
    return dict(seed=seed,family=fam,length=length,context=ctx,sign=sign,
                shifts=tuple(shifts),observations=tuple(observations),
                calibration=tuple(calib),chain=chain,start=start,truth=at)

def sign_in_anchor(anchor):
    ctx,(x,y),(a,b)=anchor
    slope=(b-y)*pow((a-x)%N,-1,N)%N
    assert slope in (1,6)
    return -1 if slope==6 else 1

def public_support_pminus(d):
    votes=[]
    for a in d["calibration"]:
        distance=sum(x!=y for x,y in zip(a[0],d["context"]))
        weight=math.exp(-distance)
        votes.append((weight,sign_in_anchor(a)))
    return sum(w for w,s in votes if s==-1)/sum(w for w,s in votes)

def feat(d):
    c=d["context"]
    p=public_support_pminus(d)
    return (1.,2*c[0]-1,2*c[1]-1,2*c[2]-1,2*p-1,
            1. if p>.5 else -1.,sum(c)/3)

def train():
    corpus=[]
    for seed in FIT:
      for fam in FIT_FAM:
       for length in (4,6):
        d=source(seed,fam,length)
        corpus.append((feat(d),int(d["sign"]==-1)))
    w=[0.]*7
    for epoch in range(50):
        for f,y in corpus:
            z=sum(a*b for a,b in zip(w,f))
            p=1/(1+math.exp(-max(-12,min(12,z))))
            for k in range(7):
                w[k]-=.035*((p-y)*f[k]+.0003*w[k])
    # Quantize signed weights to 8-bit to freeze persistent model size.
    q=[max(-127,min(127,round(val*28))) for val in w]
    return tuple(q),len(corpus)

def sign_for(arm,d,weights):
    if arm=="trained_quantized_state":
        z=sum(w*f for w,f in zip(weights,feat(d)))
        return -1 if z>=0 else 1
    if arm=="budget_matched_symbolic_approx":
        return -1 if public_support_pminus(d)>=.5 else 1
    if arm=="public_three_bit_rule":
        return -1 if d["context"][0]^d["context"][1]^d["context"][2] else 1
    if arm=="public_support_hamming":
        a=min(d["calibration"],key=lambda x:(sum(k!=z for k,z in zip(x[0],d["context"])),x[0]))
        return sign_in_anchor(a)
    return 1

def role_likelihoods(d,sign):
    result=[]
    for r,x,y in d["observations"]:
        values=[]
        for shift in range(N):
            expected=(sign*x+shift)%N
            values.append(.83 if expected==y else .17/6)
        scale=sum(values)
        result.append(tuple(v/scale for v in values))
    return result

def propagate(d,sign,quantized):
    likelihood=role_likelihoods(d,sign)
    state=[0.]*N;state[d["start"]]=1.
    operations=4*7+4*7  # observation likelihood comparisons and normalizations
    for role,direction in d["chain"]:
        nextstate=[0.]*N
        # Compressed policy retains the two highest-weight active answer
        # hypotheses and explicitly counts each of its 2x7 transitions.
        # Full exact control propagates all seven positions for both signs.
        active=(sorted(range(N),key=lambda i:(-state[i],i))[:2]
                if quantized else range(N))
        for at in active:
            for shift in range(N):
                if direction==1:nxt=(sign*at+shift)%N
                else:nxt=(sign*(at-shift))%N
                nextstate[nxt]+=state[at]*likelihood[role][shift]
                operations+=1
        if quantized:
            # 7-byte recurrent belief, stable quantized normalization.
            ints=[min(255,max(0,int(round(v*255)))) for v in nextstate]
            if sum(ints)==0:ints[max(range(N),key=lambda i:nextstate[i])]=1
            state=[v/sum(ints) for v in ints]
            operations+=7
        else:state=nextstate
    return state,operations

def decide(arm,d,weights):
    if arm=="no_observation":
        return d["start"],0,1
    if arm=="full_exact_public_posterior":
        p=public_support_pminus(d)
        pos,cpos=propagate(d,1,False)
        neg,cneg=propagate(d,-1,False)
        posterior=[(1-p)*a+p*b for a,b in zip(pos,neg)]
        # Actual two-sign posterior computation and public source calibration.
        work=cpos+cneg+4*9+7
        return max(range(N),key=lambda i:(posterior[i],-i)),work,14
    sign=sign_for(arm,d,weights)
    if arm=="budget_matched_symbolic_approx":
        # Nonlearned fixed-support sign and 7-step mode approximation, but
        # same observation likelihood evidence as the learned candidate.
        prob,cost=propagate(d,sign,True)
        return max(range(N),key=lambda i:(prob[i],-i)),cost+7,7
    state,cost=propagate(d,sign,True)
    return max(range(N),key=lambda i:(state[i],-i)),cost+7,7

def run():
    assert set(FIT).isdisjoint(DIAG) and set(FIT_FAM).isdisjoint(DIAG_FAM)
    assert not set(DIAG)&set(RESERVED)
    w,nfit=train()
    rows=[];work={a:0 for a in ARMS[:-1]}
    for seed in DIAG:
      for fam in DIAG_FAM:
       for length in (6,8):
        d=source(seed,fam,length)
        result={"private_truth_oracle_unfair":d["truth"]}
        costs={};bytes_used={}
        for arm in ARMS[:-1]:
            guess,cost,b=decide(arm,d,w)
            result[arm]=guess;costs[arm]=cost;bytes_used[arm]=b
            work[arm]+=cost
        rows.append(dict(seed=seed,family=fam,length=length,
         public_context=list(d["context"]),public_observations=[list(x) for x in d["observations"]],
         public_calibration=[[list(a[0]),list(a[1]),list(a[2])] for a in d["calibration"]],
         chain=[list(x) for x in d["chain"]],start=d["start"],
         evaluator_truth=d["truth"],answers=result,
         abstract_operation_counts=costs,belief_state_bytes=bytes_used))
    scores={arm:dict(correct=sum(row["answers"][arm]==row["evaluator_truth"] for row in rows),
                       total=len(rows)) for arm in ARMS}
    for s in scores.values():s["accuracy"]=s["correct"]/s["total"]
    strata={str(f):{a:sum(row["answers"][a]==row["evaluator_truth"] for row in rows
                   if row["family"]==f) for a in ARMS} for f in DIAG_FAM}
    candidate="trained_quantized_state"
    budget_fair=("budget_matched_symbolic_approx","public_three_bit_rule","public_support_hamming")
    best=max(budget_fair,key=lambda a:scores[a]["accuracy"])
    gain=scores[candidate]["accuracy"]-scores[best]["accuracy"]
    exactloss=scores["full_exact_public_posterior"]["accuracy"]-scores[candidate]["accuracy"]
    ratio=work[candidate]/max(1,work["full_exact_public_posterior"])
    no_reg=all(strata[str(f)][candidate]>=max(strata[str(f)][a] for a in budget_fair)
               for f in DIAG_FAM)
    passed=gain>=.05 and exactloss<=.02 and ratio<=.5 and no_reg
    return dict(schema="wingless.f52.quantized-compositional-belief.trainonly.v1",
      stage="DISJOINT_TRAIN_ONLY_NOVEL_VARIABLE_BINDING_NOT_ORIGINAL",
      fit_seeds=[FIT.start,FIT.stop],fit_families=list(FIT_FAM),
      diagnostic_seeds=[DIAG.start,DIAG.stop],diagnostic_families=list(DIAG_FAM),
      future_original_seeds_NOT_RESERVED=[RESERVED.start,RESERVED.stop],
      trained_quantized_weights=list(w),persistent_weight_bytes=7,
      train_examples=nfit,hidden_belief_state_bytes=7,
      physical_total_576B_unverified=True,
      cases=len(rows),rows=rows,scores=scores,by_family=strata,
      posterior_abstract_operations=work,strongest_budget_matched_fair=best,
      candidate_gain=gain,absolute_loss_vs_exact=exactloss,
      operation_ratio_vs_exact=ratio,family_nondegradation=no_reg,
      frozen_gate_met=passed,
      verdict=("TRAIN_ONLY_QUANTIZED_COMPOSITION_SUPPORTED" if passed else
               "TRAIN_ONLY_QUANTIZED_COMPOSITION_NEGATIVE"),
      independent_F3=False,original_executed=False,accepted_ref_mutation=False,
      paid_api_calls=0)
if __name__=="__main__":print(json.dumps(run(),indent=2,sort_keys=True))
