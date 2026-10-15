"""W-F53: learn public-support attention for unseen algebraic source laws.

Exactly four PUBLIC calibration contexts per diagnostic task, each backed by
two public role-mapping values which unambiguously identify that anchor's
operator sign. Target private sign is never supplied to a policy. Fit and
diagnostic source-law sets are disjoint, with fresh role composition orders.
"""
from __future__ import annotations
import itertools,json,math,random
FIT=range(581000,581008)
DIAG=range(581100,581116)
FIT_FAM=(1100,1101,1102,1103)
DIAG_FAM=(1104,1105,1106,1107)
RESERVED=range(582000,582016)
ARMS=("fit_minimax_public_mixture","public_bayesian_law_posterior","learned_support_attention","public_hamming_nearest",
      "public_support_majority","public_three_bit_parity",
      "public_three_bit_majority","always_shift","always_reflect",
      "private_sign_oracle_unfair")
LAWS=("parity","majority","xnor","exactly_one",
      "xor_and","or_not","threshold_two","impl")
CONTEXTS=tuple(itertools.product((0,1),repeat=3))
PATHS=(
 ((0,1),(2,-1),(1,1),(3,1),(0,-1)),
 ((3,-1),(1,1),(2,-1),(1,-1),(0,1)),
 ((1,1),(3,-1),(0,1),(2,1),(3,-1)),
 ((2,1),(0,-1),(3,1),(2,-1),(1,1)),
 ((1,-1),(2,1),(0,1),(3,-1),(2,-1)),
 ((3,1),(0,1),(1,-1),(3,1),(2,-1)),
 ((0,1),(1,-1),(3,-1),(0,-1),(2,1)),
 ((2,-1),(3,1),(2,1),(0,1),(3,-1)))
PATHS=tuple(tuple(reversed(path)) for path in PATHS)
N=7
def law(name,c):
    a,b,d=c
    if name=="parity":return a^b^d
    if name=="majority":return int(a+b+d>=2)
    if name=="xnor":return 1^(a^b^d)
    if name=="exactly_one":return int(a+b+d==1)
    if name=="xor_and":return a^(b&d)
    if name=="or_not":return a|(b&(1^d))
    if name=="threshold_two":return int(a+b*2+d>=3)
    if name=="impl":return int((not a) or (b and d))
    raise ValueError("private source law unknown")
def affine(sign,shift,x):
    return (sign*x+shift)%N
def infer_sign(anchor):
    ctx,(x,y),(z,t)=anchor
    assert x!=z
    # Multiplicative inverse exists for nonzero differences modulo prime 7.
    slope=((t-y)*pow((z-x)%N,-1,N))%N
    if slope not in (1,6):raise ValueError("invalid public calibration anchor")
    return -1 if slope==6 else 1
def source(seed,fam,length):
    if not ((seed in FIT and fam in FIT_FAM) or
            (seed in DIAG and fam in DIAG_FAM)) or length not in (3,5):
        raise ValueError("source outside frozen FIT/DIAG; originals denied")
    rng=random.Random(seed*19139+fam*3797+length*139+51019)
    law_name=LAWS[fam-1100]
    ctx=CONTEXTS[rng.randrange(8)]
    hidden_sign=-1 if law(law_name,ctx)^(rng.random()<.11) else 1
    shifts=[rng.randrange(N) for role in range(4)]
    public=[]
    for role in range(4):
        x=rng.randrange(N)
        public.append((role,x,affine(hidden_sign,shifts[role],x)))
    # Four calibration contexts are public pairs, not hidden operator labels.
    picked=rng.sample(CONTEXTS,5)
    anchors=[]
    for calibration_context in picked:
        calibration_sign=-1 if law(law_name,calibration_context)^(rng.random() < (.08 if seed in FIT else .14)) else 1
        b=rng.randrange(N)
        x=rng.randrange(N);y=(x+rng.randrange(1,N))%N
        anchors.append((calibration_context,(x,affine(calibration_sign,b,x)),
                        (y,affine(calibration_sign,b,y))))
    chain=PATHS[fam-1100][:length]
    start=rng.randrange(N)
    return dict(seed=seed,family=fam,length=length,context=ctx,
       private_sign=hidden_sign,public=public,calibration=anchors,
       chain=chain,start=start)
def compose(pub,chain,start,sign):
    shifts={role:(y-sign*x)%N for role,x,y in pub}
    at=start
    for role,d in chain:
        at=(sign*at+shifts[role])%N if d==1 else (sign*(at-shifts[role]))%N
    return at
def similarity(c,ref,w):
    return sum(weight*(1 if x==y else -1)
               for x,y,weight in zip(c,ref,w))
def support_probability(context,anchors,weights):
    total=neg=0.
    for a in anchors:
        score=math.exp(max(-15,min(15,similarity(context,a[0],weights))))
        total+=score
        if infer_sign(a)<0:neg+=score
    return neg/max(total,1e-12)
def public_law_probability(context,anchors):
    # 8 * PUBLIC registered hypothesis laws, 5 public sign calibrations;
    # neither hidden diagnostic law/family nor target true sign is inspected.
    weights=[]
    for law_name in LAWS:
        likelihood=1.
        for calibration in anchors:
            expected=-1 if law(law_name,calibration[0]) else 1
            likelihood*=.92 if infer_sign(calibration)==expected else .08
        guess_negative=int(law(law_name,context))
        weights.append((likelihood,guess_negative))
    denom=sum(w for w,_ in weights)
    return sum(w*(.89 if neg else .11) for w,neg in weights)/denom

def group_robust_sign(c,anchors,weights,params):
    # Source-limited public actions: params fit on known FIT families only.
    cutoff,fallback=params
    probability=public_law_probability(c,anchors)
    if abs(probability-.5)>=cutoff:return -1 if probability>=.5 else 1
    if fallback=="nearest":
        best=min(anchors,key=lambda a:(sum(x!=y for x,y in zip(c,a[0])),a[0]))
        return infer_sign(best)
    if fallback=="support":
        return -1 if support_probability(c,anchors,weights)>=.5 else 1
    return -1 if sum(infer_sign(a)<0 for a in anchors)>=3 else 1

def fit_group_robust(weights):
    # FIT-only minimax family score; no diagnostic sample or hidden heldout law.
    data={f:[] for f in FIT_FAM}
    for seed in FIT:
        for family in FIT_FAM:
            for length in (3,5):
                d=source(seed,family,length)
                data[family].append(d)
    grid=tuple((v,f) for v in (0.,.08,.16,.24,.32,.40)
               for f in ("nearest","support","majority"))
    counts={}
    for params in grid:
        per={fam:sum(group_robust_sign(d["context"],d["calibration"],weights,params)
                      ==d["private_sign"] for d in observations)
             for fam,observations in data.items()}
        counts[params]=per
    best=min(grid,key=lambda params:
        (-min(counts[params].values()),-sum(counts[params].values()),params))
    return best,counts[best]

def predict(arm,d,weights,params=(0.,'nearest')):
    c=d["context"];anchors=d["calibration"]
    if arm=="fit_minimax_public_mixture":
        return group_robust_sign(c,anchors,weights,params)
    if arm=="public_bayesian_law_posterior":
        return -1 if public_law_probability(c,anchors)>=.5 else 1
    if arm=="learned_support_attention":
        return -1 if support_probability(c,anchors,weights)>=.5 else 1
    if arm=="public_hamming_nearest":
        best=min(anchors,key=lambda a:(sum(x!=y for x,y in zip(c,a[0])),a[0]))
        return infer_sign(best)
    if arm=="public_support_majority":
        return -1 if sum(infer_sign(a)<0 for a in anchors)>=3 else 1
    if arm=="public_three_bit_parity":return -1 if (c[0]^c[1]^c[2]) else 1
    if arm=="public_three_bit_majority":return -1 if sum(c)>=2 else 1
    if arm=="always_shift":return 1
    if arm=="always_reflect":return -1
    raise ValueError("unregistered public control")
def fit():
    data=[]
    for seed in FIT:
      for fam in FIT_FAM:
       for length in (3,5):
        d=source(seed,fam,length)
        data.append((d["context"],d["calibration"],d["private_sign"]))
    # Discrete finite, predeclared attention-weight search. All labels from FIT.
    choices=(0.,0.5,2.,4.)
    settings=list(itertools.product(choices,repeat=3))
    accuracy={}
    for weights in settings:
        correct=sum(((-1 if support_probability(c,anc,weights)>=.5 else 1)==sign)
                    for c,anc,sign in data)
        accuracy[weights]=correct
    best=min(settings,key=lambda w:(-accuracy[w],sum(w),w))
    return best,len(data),accuracy[best]
def run():
    assert not set(FIT)&set(DIAG) and not set(FIT_FAM)&set(DIAG_FAM)
    assert not set(DIAG)&set(RESERVED)
    weights,nfit,fit_correct=fit()
    fit_params,fit_family_correct=fit_group_robust(weights)
    rows=[]
    for seed in DIAG:
      for fam in DIAG_FAM:
       for length in (3,5):
        d=source(seed,fam,length)
        truth=compose(d["public"],d["chain"],d["start"],d["private_sign"])
        pred={"private_sign_oracle_unfair":truth};signs={}
        for arm in ARMS[:-1]:
            sign=predict(arm,d,weights,fit_params)
            signs[arm]=sign
            pred[arm]=compose(d["public"],d["chain"],d["start"],sign)
        rows.append(dict(seed=seed,family=fam,length=length,
          context=d["context"],public=[list(p) for p in d["public"]],
          calibration=[[list(a[0]),list(a[1]),list(a[2])] for a in d["calibration"]],
          chain=[list(z) for z in d["chain"]],start=d["start"],
          evaluator_private_sign=d["private_sign"],truth=truth,
          answers=pred,chosen_signs=signs))
    scores={a:dict(correct=sum(r["answers"][a]==r["truth"] for r in rows),total=len(rows))
            for a in ARMS}
    for s in scores.values():s["accuracy"]=s["correct"]/s["total"]
    by_fam={str(f):{a:sum(r["answers"][a]==r["truth"] for r in rows if r["family"]==f)
                    for a in ARMS} for f in DIAG_FAM}
    fair=ARMS[1:-1]
    best=max(fair,key=lambda a:scores[a]["accuracy"])
    gain=scores["fit_minimax_public_mixture"]["accuracy"]-scores[best]["accuracy"]
    family_pass=all(by_fam[str(f)]["fit_minimax_public_mixture"]>=
                    max(by_fam[str(f)][a] for a in fair) for f in DIAG_FAM)
    qualified=gain>=.05 and family_pass
    return dict(schema="wingless.f53.group-minimax-public-law.trainonly.v1",
       stage="TRAIN_ONLY_UNSEEN_CONTEXT_OPERATOR_LAWS_NOT_ORIGINAL",
       fit_seeds=[FIT.start,FIT.stop],fit_families=list(FIT_FAM),
       diagnostic_seeds=[DIAG.start,DIAG.stop],diagnostic_families=list(DIAG_FAM),
       original_seeds_NOT_RESERVED=[RESERVED.start,RESERVED.stop],
       train_examples=nfit,train_best_correct=fit_correct,
       learned_attention_weights=list(weights),fit_group_params=list(fit_params),
       fit_family_correct={str(k):v for k,v in fit_family_correct.items()},
       diagnostic_calibration_noise=.14,cases=len(rows),rows=rows,
       scores=scores,by_family=by_fam,strongest_fair_control=best,
       candidate_gain=gain,per_family_nondegradation=family_pass,
       frozen_gate_met=qualified,verdict=("TRAIN_ONLY_UNSEEN_OPERATOR_SUPPORTED"
         if qualified else "TRAIN_ONLY_UNSEEN_OPERATOR_NEGATIVE"),
       diagnostic_private_sign_never_selector_input=True,
       source_independent_F3=False,original_executed=False,
       native_576_bytes_residency_qualified=False,accepted_ref_mutation=False,
       paid_api_calls=0)
if __name__=="__main__":print(json.dumps(run(),indent=2,sort_keys=True))
