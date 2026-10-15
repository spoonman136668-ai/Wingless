"""W-F42 frozen-source diagnostic: posterior rank convergence and acquisition value.
Never a learned-policy original. Private role mappings only in evaluator/paid
Port, no synthetic F3 independent custody. Reuse F41 *analysis functions*
but generate entirely new source identities, structures, and provenance.
"""
from __future__ import annotations
import importlib.util,json,random,hashlib,math,sys
from pathlib import Path
from collections import Counter
P=Path(__file__).resolve().parent
f=P.parent/"20261015-w-f41"/"source_pilot.py"
s=importlib.util.spec_from_file_location("wf42_prior_analytics",f)
prior=importlib.util.module_from_spec(s);sys.modules[s.name]=prior;s.loader.exec_module(prior)
SEEDS=range(321000,321012);FAMILIES=(370,371,372)
RESERVED=range(322000,322016);RESERVED_FAMILIES=(374,375,376,377)
LENGTHS=(3,5);BUDGETS=(12,48,192)
ARMS=("query_frontier","decision_value_48","decision_value_192",
      "information_gain_192","no_purchase","privileged_hindsight_unfair")
def case(seed,family,length):
    if seed not in SEEDS or family not in FAMILIES or length not in LENGTHS:
        raise ValueError("reserved/consumed source not authorized")
    rng=random.Random(seed*77417+family*9337+length*659+8803)
    maps=[]
    for role in range(8):
        if family==370:
            a=rng.choice((1,3,5,7));b=rng.randrange(8)
            op=tuple((a*x+b)%8 for x in range(8))
        elif family==371:
            arr=list(range(8));rng.shuffle(arr);op=tuple(arr)
        else:
            a=rng.choice((1,3,5,7));b=rng.randrange(8);mask=rng.randrange(1,8)
            op=tuple(((a*x+b)%8)^mask for x in range(8))
        maps.append(op)
    public=[]
    for role in range(8):
        for x in rng.sample(range(8),3):public.append((role,x,maps[role][x]))
    rng.shuffle(public)
    start=rng.randrange(8)
    query=tuple((rng.randrange(8),rng.choice((-1,1))) for _ in range(length))
    return tuple(maps),tuple(public),start,query

def select(method,observed,start,query,seed):
    if method=="no_purchase":return None
    sample_count=192 if "192" in method else 48
    scores,_,_,_,_=prior.score_actions(observed,start,query,sample_count,seed)
    if method=="query_frontier":key=lambda q:(scores[q][2],-q[0],-q[1])
    elif method.startswith("decision_value"):key=lambda q:(scores[q][1],scores[q][2],-q[0],-q[1])
    elif method=="information_gain_192":key=lambda q:(scores[q][0],scores[q][2],-q[0],-q[1])
    else:raise ValueError("unknown arm")
    return max(scores,key=key) if scores else None

def pick_bayes(observed,start,query,seed):
    return prior.majority_prediction(observed,start,query,seed)

def run():
    assert not set(SEEDS)&set(RESERVED)
    assert not set(FAMILIES)&set(RESERVED_FAMILIES)
    rows=[];agreements=dict(b12_48=0,b48_192=0,b192_192_independent=0)
    for seed in SEEDS:
      for family in FAMILIES:
       for length in LENGTHS:
        private,public,start,query=case(seed,family,length)
        truth=prior.follow(private,start,query)
        assert truth==prior.independent_follow(private,start,query)
        initial_seed=seed*101+family+length
        ranks={}
        for n,token in ((12,"b12"),(48,"b48"),(192,"b192"),(192,"b192_independent")):
            salt=initial_seed+(997 if token.endswith("independent") else 0)
            scores,_,_,_,_=prior.score_actions(public,start,query,n,salt)
            ranks[token]=max(scores,key=lambda q:(scores[q][1],scores[q][2],-q[0],-q[1]))
        agreements["b12_48"]+=ranks["b12"]==ranks["b48"]
        agreements["b48_192"]+=ranks["b48"]==ranks["b192"]
        agreements["b192_192_independent"]+=ranks["b192"]==ranks["b192_independent"]
        guesses={"privileged_hindsight_unfair":truth};receipts={}
        for method in ARMS[:-1]:
            observed=list(public);port=prior.Port(private,public)
            for i in range(0 if method=="no_purchase" else 2):
                pick=select(method,observed,start,query,initial_seed+201*i)
                if pick is None:break
                y=port.buy(*pick)
                observed.append((*pick,y))
            guesses[method]=pick_bayes(observed,start,query,initial_seed+1103)
            receipts[method]=[list(v) for v in port.receipts]
        rows.append(dict(seed=seed,family=family,length=length,truth=truth,
            predictions=guesses,receipts=receipts,ranking=ranks,
            public_sha256=hashlib.sha256(json.dumps((public,start,query)).encode()).hexdigest(),
            paid_mapping_max=2))
    scores={arm:dict(correct=sum(r["predictions"][arm]==r["truth"] for r in rows),total=len(rows))
            for arm in ARMS}
    for v in scores.values():v["accuracy"]=v["correct"]/v["total"]
    family_scores={str(f):{a:sum(r["predictions"][a]==r["truth"] for r in rows if r["family"]==f) for a in ARMS}
                   for f in FAMILIES}
    agreement={k:v/len(rows) for k,v in agreements.items()}
    strongest=max(scores[a]["accuracy"] for a in ("query_frontier","decision_value_48","information_gain_192"))
    gain=scores["decision_value_192"]["accuracy"]-strongest
    no_degrade=all(family_scores[str(f)]["decision_value_192"]>=max(
       family_scores[str(f)][a] for a in ("query_frontier","decision_value_48","information_gain_192"))
       for f in FAMILIES)
    passed=(agreement["b48_192"]>=.55 and gain>=.03 and no_degrade)
    return dict(schema="wingless.f42.action-rank-sample-stability.trainonly.v1",
      stage="TRAIN_ONLY_SOURCE_DIAGNOSTIC_NOT_LEARNED_ORIGINAL",
      source_seeds=[SEEDS.start,SEEDS.stop],source_families=list(FAMILIES),
      reserved_original_seeds=[RESERVED.start,RESERVED.stop],
      reserved_original_families=list(RESERVED_FAMILIES),
      questions=len(rows),posterior_budgets=list(BUDGETS),paid_mapping_cap=2,
      action_rank_agreements=agreement,scores=scores,by_family=family_scores,rows=rows,
      high_sample_posterior_external_bytes=192*8*8,
      private_maps_only_in_port_and_evaluator=True,synthetic_same_repo_not_F3=True,
      primary_original_not_executed=True,accepted_ref_mutation=False,paid_api_calls=0,
      verdict="TRAIN_ONLY_RANK_STABILITY_QUALIFIED" if passed else "TRAIN_ONLY_RANK_STABILITY_NOT_QUALIFIED")
if __name__=="__main__":print(json.dumps(run(),indent=2,sort_keys=True))
