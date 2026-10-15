"""W-F51: learned latent role-operator posterior and unseen relation composition.

TRAIN-only synthetic algebraic substrate; no independent F3 original or
physically qualified native 576B residency. Diagnostic policy sees one public
mapping per role and three public context bits; private sign stays evaluator-only.
"""
import json,random
FIT=range(521000,521016)
DIAG=range(521100,521116)
FIT_FAM=(900,901,902,903)
DIAG_FAM=(904,905,906,907)
RESERVED=range(522000,522016)
N=7
ARMS=("learned_operator_posterior","always_shift","always_reflect",
      "full_five_bit_public_nonlinear_rule","three_bit_public_parity","public_hash",
      "private_operator_oracle_unfair")
TEMPLATES=(
 ((0,1),(1,-1),(2,1),(0,-1),(3,1)),
 ((1,-1),(3,1),(1,1),(2,-1),(0,1)),
 ((3,1),(0,1),(2,-1),(3,-1),(1,1)),
 ((2,1),(1,1),(0,-1),(2,-1),(3,1)),
 ((0,-1),(2,1),(3,1),(1,-1),(0,1)),
 ((3,-1),(1,-1),(3,1),(0,1),(2,-1)),
 ((1,1),(0,-1),(1,-1),(3,1),(2,1)),
 ((2,-1),(3,1),(0,1),(2,1),(1,-1)),
)
def source(seed,family,length):
    if not ((seed in FIT and family in FIT_FAM) or
            (seed in DIAG and family in DIAG_FAM)) or length not in (3,5):
        raise ValueError("source identity outside frozen synthetic TRAIN split")
    rng=random.Random(seed*51787+family*1097+length*607+50031)
    ctx=tuple(rng.randrange(2) for _ in range(5))
    correct_parity=((ctx[0]&ctx[1])^(ctx[2]&ctx[3])^ctx[4])
    sign=-1 if correct_parity ^ (rng.random()<.14) else 1
    private_shifts=tuple(rng.randrange(N) for _ in range(4))
    visible=[]
    for role in range(4):
        x=rng.randrange(N)
        visible.append((role,x,(sign*x+private_shifts[role])%N))
    chain=TEMPLATES[family-900][:length]
    start=rng.randrange(N)
    return dict(seed=seed,family=family,length=length,context=ctx,
        sign=sign,shifts=private_shifts,public=tuple(visible),chain=chain,start=start)

def compose(ctx,public,chain,start,sign):
    # Only bought/public triples can infer the role-specific latent shift.
    shifts={}
    for role,x,y in public:
        shifts[role]=(y-sign*x)%N
    value=start
    for role,direction in chain:
        offset=shifts[role]
        if direction==1:value=(sign*value+offset)%N
        else:value=(sign*(value-offset))%N
    return value

def fit():
    counts={format(i,"05b"):[1,1] for i in range(32)}
    n=0
    for seed in FIT:
      for family in FIT_FAM:
       for length in (3,):
        d=source(seed,family,length)
        code="".join(map(str,d["context"]))
        counts[code][int(d["sign"]==-1)]+=1
        n+=1
    # Conjugate beta-Bernoulli posterior over role operator sign.
    posterior={c:v[1]/sum(v) for c,v in counts.items()}
    return posterior,n

def guess(arm,ctx,p_signminus,seed,family,length):
    if arm=="learned_operator_posterior":return -1 if p_signminus>=.5 else 1
    if arm=="always_shift":return 1
    if arm=="always_reflect":return -1
    if arm=="cue0_rule":return -1 if ctx[0] else 1
    if arm=="two_bit_public_parity":return -1 if ctx[0]^ctx[1] else 1
    if arm=="fixed_public_hash":
        return -1 if random.Random(seed*61+family*23+length).randrange(2) else 1
    raise ValueError("non-public policy")
def run():
    assert not set(FIT)&set(DIAG)
    assert not set(FIT_FAM)&set(DIAG_FAM)
    assert not set(DIAG)&set(RESERVED)
    posterior,ntrain=fit()
    rows=[]
    for seed in DIAG:
      for fam in DIAG_FAM:
       for length in (3,5):
        d=source(seed,fam,length)
        true=compose(d["context"],d["public"],d["chain"],d["start"],d["sign"])
        answers={"private_operator_oracle_unfair":true};signs={}
        for arm in ARMS[:-1]:
            sign=guess(arm,d["context"],posterior["".join(map(str,d["context"]))],
                       seed,fam,length)
            answers[arm]=compose(d["context"],d["public"],d["chain"],d["start"],sign)
            signs[arm]=sign
        rows.append(dict(seed=seed,family=fam,length=length,
           context=d["context"],public=[list(z) for z in d["public"]],
           chain=[list(z) for z in d["chain"]],start=d["start"],truth=true,
           answers=answers,chosen_signs=signs))
    scores={arm:dict(correct=sum(row["answers"][arm]==row["truth"] for row in rows),
                      total=len(rows)) for arm in ARMS}
    for s in scores.values():s["accuracy"]=s["correct"]/s["total"]
    by_fam={str(f):{arm:sum(row["answers"][arm]==row["truth"] for row in rows
                  if row["family"]==f) for arm in ARMS} for f in DIAG_FAM}
    fair=ARMS[1:-1]
    strongest=max(fair,key=lambda arm:scores[arm]["accuracy"])
    gain=scores["learned_operator_posterior"]["accuracy"]-scores[strongest]["accuracy"]
    nonreg=all(by_fam[str(f)]["learned_operator_posterior"]>=max(
         by_fam[str(f)][arm] for arm in fair) for f in DIAG_FAM)
    qualified=gain>=.05 and nonreg
    return dict(schema="wingless.f51.five-bit-public-operator.trainonly.v1",
       stage="DISJOINT_FIT_AND_NOVEL_COMPOSITION_TRAIN_ONLY",
       fit_seeds=[FIT.start,FIT.stop],fit_families=list(FIT_FAM),
       diagnostic_seeds=[DIAG.start,DIAG.stop],diagnostic_families=list(DIAG_FAM),
       proposed_original_seeds_NOT_RESERVED=[RESERVED.start,RESERVED.stop],
       train_examples=ntrain,posterior_minus_by_public_context=posterior,
       cases=len(rows),scores=scores,by_family=by_fam,
       strongest_fair_control=strongest,candidate_gain=gain,
       per_family_nondegradation=nonreg,frozen_criterion_met=qualified,
       verdict="TRAIN_ONLY_OPERATOR_BINDING_SUPPORTED" if qualified else
               "TRAIN_ONLY_OPERATOR_BINDING_NEGATIVE",
       rows=rows,source_independent_F3=False,original_executed=False,
       native_576_bytes_residency_qualified=False,accepted_ref_mutation=False,paid_api_calls=0)

if __name__=="__main__":print(json.dumps(run(),indent=2,sort_keys=True))
