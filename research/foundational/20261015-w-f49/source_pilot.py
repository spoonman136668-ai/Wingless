"""W-F49 public posterior reveal-value identifiability and action rank stability.

No private mapping/truth is used by selectors; evaluator alone sees outcomes.
No owner original, independent F3, Codex, Cloudflare, or accepted ref mutation.
"""
import json
import random
from pathlib import Path
import importlib.util
import sys

P = Path(__file__).resolve().parent.parent / "20261015-w-f48" / "source_pilot.py"
spec = importlib.util.spec_from_file_location("w_f49_f48_controls", P)
prior = importlib.util.module_from_spec(spec)
sys.modules[spec.name] = prior
spec.loader.exec_module(prior)
portlib = prior.old.f45.prior.prior
FIT = range(431000,431004)
DIAG = range(431100,431108)
ORIGINAL = range(432000,432016)
FAMILIES=(640,641,642,643)
ORIGINAL_FAMILIES=(644,645,646,647)
ARMS=("public_stability_gate","one_step_value","exact_entropy","query_frontier",
      "budget_matched_rank","no_purchase","evaluator_oracle_unfair")
FAIR=("one_step_value","exact_entropy","query_frontier","budget_matched_rank")

def source(seed, family, length):
    if seed not in FIT and seed not in DIAG:
        raise ValueError("unregistered or reserved seed")
    if family not in FAMILIES or length not in (3,5):
        raise ValueError("unregistered relation family")
    rng=random.Random(seed*17713+family*3037+length*521+49001)
    private=[];public=[]
    for role in range(4):
        perm=list(range(6));rng.shuffle(perm);private.append(tuple(perm))
        for x in rng.sample(range(6),3):public.append((role,x,perm[x]))
    rng.shuffle(public)
    patterns={
        640:((3,-1),(2,1),(0,-1),(3,1),(1,1)),
        641:((0,1),(1,-1),(3,1),(2,-1),(0,1)),
        642:((2,-1),(3,1),(1,1),(0,-1),(2,1)),
        643:((1,-1),(0,1),(2,-1),(3,1),(1,-1))}
    return tuple(private),tuple(public),rng.randrange(6),patterns[family][:length]

def rank_stability(observed, start, chain, work):
    """Independent disjoint partitions of PUBLIC posterior worlds; no oracle access."""
    worlds=portlib.worlds(observed)
    actions=portlib.candidate_actions(observed)
    if not actions:return None,False,0,0
    parts=(worlds[::2],worlds[1::2])
    if not all(parts):return None,False,0,0
    def pick(bank):
        base=prior.old.f45.prior.best_success(bank,start,chain)
        values={a:prior.old.f45.prior.split_success(bank,a,start,chain) for a in actions}
        work[0]+=len(bank)*len(actions)
        winner=min(actions,key=lambda a:(-values[a],a))
        return winner, values[winner]>base, values[winner]-base
    a,apos,again=pick(parts[0])
    b,bpos,bgain=pick(parts[1])
    agree=bool(a==b and apos and bpos)
    # Action decision only from registered public graph and posterior worlds.
    return a if agree else None,agree,again,bgain

def simulate(seed,family,length,arm):
    maps,public,start,chain=source(seed,family,length)
    truth=portlib.follow(maps,start,chain)
    observed=list(public)
    purchased=portlib.PurchasePort(maps,public)
    work=[0]; first_agree=None; first_gains=(0,0)
    if arm=="evaluator_oracle_unfair":return truth,[],0,truth,first_agree,first_gains
    for step in range(2):
        if arm=="no_purchase":selected=None
        elif arm=="public_stability_gate":
            selected,agree,g1,g2=rank_stability(observed,start,chain,work)
            if step==0:first_agree,first_gains=agree,(g1,g2)
            if selected is None:
                # Deterministic public posterior one-step fallback.
                selected=prior.old.f45.prior.choose("one_step_value",
                       observed,start,chain,2-step,work)
        elif arm=="budget_matched_rank":
            selected=prior.old.select(observed,start,chain,work,False)
        else:
            selected=prior.old.f45.prior.choose(arm,observed,start,chain,2-step,work)
        if selected is None:break
        value=purchased.buy(*selected)
        observed.append((*selected,value))
    pred=portlib.predict(observed,start,chain)
    return pred,[list(x) for x in purchased.receipts],work[0],truth,first_agree,first_gains

def cases():
    for seed in DIAG:
        for family in FAMILIES:
            for length in (3,5):yield seed,family,length

def run():
    assert set(FIT).isdisjoint(DIAG) and set(DIAG).isdisjoint(ORIGINAL)
    assert set(FAMILIES).isdisjoint(ORIGINAL_FAMILIES)
    rows=[];total_work={a:0 for a in ARMS[:-1]}
    for seed,family,length in cases():
        answers={};receipts={};work={};first=None;gains=(0,0)
        for arm in ARMS:
            answer,trace,cost,truth,agreement,delta=simulate(seed,family,length,arm)
            answers[arm]=answer
            if arm!="evaluator_oracle_unfair":
                receipts[arm]=trace;work[arm]=cost;total_work[arm]+=cost
            if arm=="public_stability_gate":first=agreement;gains=delta
        rows.append(dict(seed=seed,family=family,length=length,truth=truth,
               answers=answers,receipts=receipts,work=work,
               first_purchase_public_rank_agreement=first,
               public_partition_gains=list(gains)))
    n=len(rows)
    scores={a:dict(correct=sum(r["answers"][a]==r["truth"] for r in rows),total=n)
            for a in ARMS}
    for v in scores.values():v["accuracy"]=v["correct"]/n
    best=max(FAIR,key=lambda a:(scores[a]["correct"],-total_work[a],a))
    candidate="public_stability_gate"
    gain=(scores[candidate]["correct"]-scores[best]["correct"])/n
    ratio=total_work[candidate]/max(1,total_work[best])
    by_family={str(f):{a:sum(r["answers"][a]==r["truth"] for r in rows
                 if r["family"]==f) for a in ARMS} for f in FAMILIES}
    family_ok=all(by_family[str(f)][candidate]>=
                 max(by_family[str(f)][a] for a in FAIR) for f in FAMILIES)
    stable=[r for r in rows if r["first_purchase_public_rank_agreement"]]
    unstable=[r for r in rows if not r["first_purchase_public_rank_agreement"]]
    def uplift(rs):
        if not rs:return 0.
        return sum(int(r["answers"][candidate]==r["truth"])-
                   int(r["answers"]["no_purchase"]==r["truth"])
                   for r in rs)/len(rs)
    stable_uplift=uplift(stable)
    unstable_uplift=uplift(unstable)
    stability_rate=len(stable)/n
    passed=(len(stable)>=4 and len(unstable)>=4 and stability_rate>=.60
            and stable_uplift-unstable_uplift>=.05 and gain>=.03
            and ratio<=3 and family_ok)
    return dict(schema="wingless.f49.public-rank-information.trainonly.v1",
        stage="NEW_PUBLIC_RANK_DIAGNOSTIC_NOT_PRIMARY_ORIGINAL",
        diagnostic_seeds=[DIAG.start,DIAG.stop],
        prereg_fit_seeds_held_not_consumed=[FIT.start,FIT.stop],
        families=list(FAMILIES),proposed_originals_not_consumed=[ORIGINAL.start,ORIGINAL.stop],
        original_families_not_consumed=list(ORIGINAL_FAMILIES),
        questions=n,rows=rows,scores=scores,posterior_world_work=total_work,
        by_family=by_family,strongest_fair_control=best,
        candidate_gain=gain,work_ratio=ratio,family_nondegradation=family_ok,
        stable_cases=len(stable),unstable_cases=len(unstable),
        public_rank_agreement_rate=stability_rate,
        stable_ex_post_purchase_uplift=stable_uplift,
        unstable_ex_post_purchase_uplift=unstable_uplift,
        frozen_criterion_met=passed,
        verdict="TRAIN_ONLY_SUPPORTED_NOT_F3" if passed else "TRAIN_ONLY_NEGATIVE",
        evaluator_private_truth_not_selector_input=True,
        primary_original_not_executed=True,independent_F3_custody=False,
        accepted_ref_mutation=False,native_576_byte_resident_qualified=False,
        paid_api_calls=0)
if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
