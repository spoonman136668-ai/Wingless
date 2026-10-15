"""Wingless F47: genuinely learned bounded value-of-computation gate.

A linear regularized estimator learns from disjoint TRAIN cases when two-step
posterior search improves terminal answer correctness relative to one-step.
At diagnostic time only public evidence and task topology determine selection.
An inexpensive global top-4 gate enforces <=8 deep-case prereg ceiling.
This reduced exact-source Python pilot does NOT qualify native <=576B memory.
"""
from __future__ import annotations
import importlib.util,json,math,random,sys,time
from collections import Counter
from pathlib import Path
ROOT=Path(__file__).resolve().parent
src=ROOT.parent/"20261015-w-f46"/"source_pilot.py"
sp=importlib.util.spec_from_file_location("w47_f46_controls",src)
prev=importlib.util.module_from_spec(sp);sys.modules[sp.name]=prev;sp.loader.exec_module(prev)
FIT=range(421000,421012)
DIAG=range(421100,421112)
FIT_FAM=(560,561,562,563)
DIAG_FAM=(564,565,566,567)
RESERVED=range(422000,422016)
RESERVED_FAM=(568,569,570,571)
ARMS=("learned_deep_gate","one_step_value","exact_entropy","query_frontier",
      "static_F46_query_edge","fixed_public_hash_deep_gate",
      "full_two_step_costly_reference","no_purchase","private_oracle_unfair")
FEATURE_DIM=9
LEARNED_DEEP_BUDGET=4
PREREG_DEEP_CEILING=8

def source(seed,family,length):
    if (seed in FIT and family in FIT_FAM) or (seed in DIAG and family in DIAG_FAM):
        pass
    else:raise ValueError("reserved/original or cross-fold source prohibited")
    if length not in (3,5):raise ValueError("unapproved chain length")
    rng=random.Random(seed*17389+family*1607+length*617+47003)
    private=[];public=[]
    for role in range(4):
        mapping=list(range(6));rng.shuffle(mapping);private.append(tuple(mapping))
        for x in rng.sample(range(6),3):public.append((role,x,mapping[x]))
    rng.shuffle(public)
    templates=(
      ((0,1),(1,-1),(0,-1),(2,1),(3,-1)),
      ((3,-1),(2,1),(3,1),(0,-1),(2,1)),
      ((1,1),(0,-1),(2,1),(1,-1),(3,1)),
      ((0,-1),(2,1),(3,-1),(0,1),(1,-1)),
    )
    chain=templates[family%4][:length]
    return tuple(private),tuple(public),rng.randrange(6),chain

def public_features(observed,start,chain):
    bank=prev.f45.prior.prior.worlds(observed)
    answers=Counter(prev.f45.prior.prior.follow(w,start,chain) for w in bank)
    sorted_counts=sorted(answers.values(),reverse=True)
    counts=Counter(r for r,d in chain)
    unique=len(counts)
    inverse=sum(d<0 for r,d in chain)
    entropy=-sum(v/len(bank)*math.log2(v/len(bank)) for v in answers.values())
    # All features remain public, bounded and insensitive to source identity.
    return (
      1.,len(chain)/5,unique/4,inverse/5,
      max(counts.values())/5,
      sorted_counts[0]/len(bank),
      (sorted_counts[0]-sorted_counts[1])/len(bank) if len(sorted_counts)>1 else 1.,
      entropy/2.6,
      sum(1 for r in counts if counts[r]>1)/4,
    )

def run_arm(method,private,public,start,chain):
    observed=list(public)
    port=prev.f45.prior.prior.PurchasePort(private,public)
    effort=[0]
    for step in range(2):
        if method=="no_purchase":selected=None
        elif method=="one_step_value":
            selected,_=prev.f45.one_step(observed,start,chain,effort)
        elif method=="static_F46_query_edge":
            selected=prev.select(observed,start,chain,effort,True)
        elif method in ("exact_entropy","query_frontier","full_two_step_costly_reference"):
            parent={"full_two_step_costly_reference":"two_step_counterfactual",
                    "exact_entropy":"exact_entropy","query_frontier":"query_frontier"}[method]
            selected=prev.f45.prior.choose(parent,observed,start,chain,2-step,effort)
        else:raise ValueError("unknown research arm")
        if selected is None:break
        val=port.buy(*selected);observed.append((*selected,val))
    answer=prev.f45.prior.prior.predict(observed,start,chain)
    return answer,[list(x) for x in port.receipts],effort[0]

def trained_gate():
    training=[]
    t0=time.perf_counter()
    for seed in FIT:
      for fam in FIT_FAM:
       for length in (3,5):
        private,public,start,chain=source(seed,fam,length)
        feats=public_features(public,start,chain)
        truth=prev.f45.prior.prior.follow(private,start,chain)
        short=run_arm("one_step_value",private,public,start,chain)[0]
        deep=run_arm("full_two_step_costly_reference",private,public,start,chain)[0]
        target=int(deep==truth)-int(short==truth)
        training.append((feats,target))
    weights=[0.]*FEATURE_DIM
    # Bounded L2-regularized squared terminal-gain regression, deterministic.
    for epoch in range(75):
        for feat,target in training:
            pred=sum(a*b for a,b in zip(weights,feat))
            error=pred-target
            for j,x in enumerate(feat):
                weights[j]-=.025*(error*x+.02*weights[j])
    return weights,len(training),time.perf_counter()-t0

def run():
    assert not set(FIT)&set(DIAG)
    assert not set(DIAG)&set(RESERVED)
    assert not set(FIT_FAM)&set(DIAG_FAM)
    assert not set(DIAG_FAM)&set(RESERVED_FAM)
    weight,nfit,training_seconds=trained_gate()
    cases=[]
    for seed in DIAG:
      for fam in DIAG_FAM:
       for length in (3,5):
        private,public,start,chain=source(seed,fam,length)
        feat=public_features(public,start,chain)
        # Ranking deliberately ignores all diagnostic labels and hidden truth.
        predicted=sum(a*b for a,b in zip(weight,feat))
        random_gate=random.Random(seed*10007+fam*101+length).random()
        cases.append(dict(seed=seed,family=fam,length=length,private=private,
                          public=public,start=start,chain=chain,features=feat,
                          predicted_computation_value=predicted,hash_gate=random_gate))
    ranked=sorted(range(len(cases)),key=lambda i:(-cases[i]["predicted_computation_value"],
                 cases[i]["seed"],cases[i]["family"],cases[i]["length"]))
    hashed=sorted(range(len(cases)),key=lambda i:(-cases[i]["hash_gate"],
                  cases[i]["seed"],cases[i]["family"],cases[i]["length"]))
    chosen=set(ranked[:LEARNED_DEEP_BUDGET])
    chosen_hash=set(hashed[:LEARNED_DEEP_BUDGET])
    assert len(chosen)==len(chosen_hash)==LEARNED_DEEP_BUDGET
    rows=[];work={a:0 for a in ARMS[:-1]}
    for index,case in enumerate(cases):
        private,public,start,chain=(case[k] for k in ("private","public","start","chain"))
        truth=prev.f45.prior.prior.follow(private,start,chain)
        answers={"private_oracle_unfair":truth};receipts={};costs={}
        for arm in ARMS[:-1]:
            if arm=="learned_deep_gate":
                method="full_two_step_costly_reference" if index in chosen else "one_step_value"
            elif arm=="fixed_public_hash_deep_gate":
                method="full_two_step_costly_reference" if index in chosen_hash else "one_step_value"
            else:method=arm
            pred,receipt,price=run_arm(method,private,public,start,chain)
            if arm in ("learned_deep_gate","fixed_public_hash_deep_gate"):
                # Account separately for the public posterior/world features used
                # to rank this batch before any diagnostic truth is observed.
                price+=1296
            answers[arm]=pred;receipts[arm]=receipt;costs[arm]=price
            work[arm]+=price
        rows.append(dict(seed=case["seed"],family=case["family"],length=case["length"],
                 truth=truth,answers=answers,receipts=receipts,work=costs,
                 learned_deep_used=index in chosen,hash_deep_used=index in chosen_hash,
                 predicted_value=case["predicted_computation_value"]))
    scores={a:dict(correct=sum(r["answers"][a]==r["truth"] for r in rows),total=len(rows))
            for a in ARMS}
    for val in scores.values():val["accuracy"]=val["correct"]/val["total"]
    by_fam={str(f):{a:sum(r["answers"][a]==r["truth"] for r in rows if r["family"]==f)
                    for a in ARMS} for f in DIAG_FAM}
    fair=("one_step_value","exact_entropy","query_frontier",
          "static_F46_query_edge","fixed_public_hash_deep_gate")
    strongest=max(fair,key=lambda a:(scores[a]["accuracy"],-work[a]))
    cand="learned_deep_gate"
    gain=scores[cand]["accuracy"]-scores[strongest]["accuracy"]
    positive_hash=scores[cand]["accuracy"]>scores["fixed_public_hash_deep_gate"]["accuracy"]
    nonregress=all(by_fam[str(f)][cand]>=max(by_fam[str(f)][a] for a in fair)
                   for f in DIAG_FAM)
    ratio=work[cand]/max(1,work[strongest])
    qualified=gain>=.03 and positive_hash and nonregress and ratio<=1.2
    return dict(schema="wingless.f47.learned-value-of-compute.trainonly.v1",
      stage="TRAIN_ONLY_LEARNED_GATE_NO_PRIMARY_ORIGINAL",
      fit_seeds=[FIT.start,FIT.stop],diagnostic_seeds=[DIAG.start,DIAG.stop],
      fit_families=list(FIT_FAM),diagnostic_families=list(DIAG_FAM),
      proposed_original_seeds_NOT_RESERVED=[RESERVED.start,RESERVED.stop],
      questions=len(rows),fit_training_examples=nfit,
      training_cpu_seconds=training_seconds,
      learned_feature_weights=weight,feature_count=FEATURE_DIM,
      chosen_deep_count=len(chosen),hash_deep_count=len(chosen_hash),
      posterior_world_work=work,scores=scores,by_family=by_fam,rows=rows,
      strongest_fair_control=strongest,candidate_gain=gain,
      positive_gain_vs_hash_gate=positive_hash,family_nondegradation=nonregress,
      posterior_work_ratio=ratio,frozen_criterion_met=qualified,
      verdict="TRAIN_ONLY_LEARNED_COMPUTE_GATE_SUPPORTED" if qualified else
              "TRAIN_ONLY_LEARNED_COMPUTE_GATE_NEGATIVE",
      source_independent_F3_custody=False,original_executed=False,
      native_576_bytes_residency_qualified=False,accepted_ref_mutation=False,
      paid_api_calls=0)

if __name__=="__main__":print(json.dumps(run(),sort_keys=True,indent=2))
