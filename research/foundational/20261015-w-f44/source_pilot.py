"""W-F44 preregistered synthetic TRAIN-only counterfactual acquisition source.

Public-only posterior counterfactuals use the exact 6-value/4-role enumerator
from F43. This explicitly is NOT an unrestricted 8-value exact inference,
native resident 576B proof, independent F3, or an original experiment.
"""
from __future__ import annotations
from collections import Counter, defaultdict
from pathlib import Path
import importlib.util, json, random, sys, time
HERE=Path(__file__).resolve().parent
UPSTREAM=HERE.parent/"20261015-w-f43"/"source_pilot.py"
spec=importlib.util.spec_from_file_location("w_f44_exact_posterior_base",UPSTREAM)
prior=importlib.util.module_from_spec(spec)
sys.modules[spec.name]=prior
spec.loader.exec_module(prior)

SEEDS=range(351000,351008)
FAMILIES=(430,431,432)
LENGTHS=(3,5)
RESERVED_SEEDS=range(352000,352016)
RESERVED_FAMILIES=(434,435,436,437)
N=6
ROLES=4
BEAM=4
ARMS=("two_step_counterfactual","one_step_value","exact_entropy",
      "query_frontier","no_purchase","private_oracle_unfair")

def source(seed,family,length):
    if seed not in SEEDS or family not in FAMILIES or length not in LENGTHS:
        raise ValueError("unregistered and reserved source prohibited")
    rng=random.Random(seed*19879+family*1913+length*659+44013)
    private=[];public=[]
    for role in range(ROLES):
        perm=list(range(N));rng.shuffle(perm);private.append(tuple(perm))
        for x in rng.sample(range(N),3):
            public.append((role,x,perm[x]))
    rng.shuffle(public)
    templates={
      430:((0,1),(1,-1),(0,1),(2,-1),(3,1)),
      431:((3,-1),(2,1),(3,1),(1,-1),(0,1)),
      432:((0,1),(1,1),(2,-1),(1,-1),(0,1))}
    return tuple(private),tuple(public),rng.randrange(N),templates[family][:length]

def best_success(bank,start,chain):
    if not bank:return 0
    return max(Counter(prior.follow(w,start,chain) for w in bank).values())

def split_success(bank,action,start,chain):
    # Sum of Bayes-optimal correct counts following one hypothetical paid mapping.
    groups=defaultdict(Counter)
    role,x=action
    for world in bank:
        groups[world[role][x]][prior.follow(world,start,chain)]+=1
    return sum(max(counts.values()) for counts in groups.values())

def two_step_value(bank,first,remaining,start,chain,effort):
    # An explicitly bounded 4-action first-stage beam; exhaustive second stage.
    # Conditional future purchase decisions use only public hypothetical results.
    first_groups=defaultdict(list)
    for world in bank:first_groups[world[first[0]][first[1]]].append(world)
    total=0
    for conditional_bank in first_groups.values():
        observed=set(remaining)
        observed.add(first)
        second=[a for a in remaining if a!=first]
        options=[best_success(conditional_bank,start,chain)]
        for a in second:
            options.append(split_success(conditional_bank,a,start,chain))
            effort[0]+=len(conditional_bank)
        total+=max(options)
    return total

def choose(method,observed,start,chain,steps_remaining,effort):
    if method=="no_purchase" or steps_remaining<=0:return None
    bank=prior.worlds(observed)
    choices=prior.candidate_actions(observed)
    if not choices:return None
    initial=best_success(bank,start,chain)
    if method=="two_step_counterfactual" and steps_remaining==2:
        scores=prior.score_actions(bank,observed,start,chain)
        # First-stage candidate admission uses observable path frontier only.
        shortlist=sorted(choices,key=lambda a:(-scores[a][2],a[0],a[1]))[:BEAM]
        ranked=[]
        for a in shortlist:
            value=two_step_value(bank,a,choices,start,chain,effort)
            ranked.append((value,-a[0],-a[1],a))
        best=max(ranked)
        return best[3] if best[0]>initial else None
    if method in ("two_step_counterfactual","one_step_value"):
        vals={a:split_success(bank,a,start,chain) for a in choices}
        effort[0]+=len(bank)*len(choices)
        action=max(choices,key=lambda a:(vals[a],-a[0],-a[1]))
        return action if vals[action]>initial else None
    scores=prior.score_actions(bank,observed,start,chain)
    effort[0]+=len(bank)*len(choices)
    slot={"exact_entropy":0,"query_frontier":2}.get(method)
    if slot is None:raise ValueError("unknown method")
    action=max(choices,key=lambda a:(scores[a][slot],-a[0],-a[1]))
    return action if scores[action][slot]>0 else None

def run():
    assert set(SEEDS).isdisjoint(RESERVED_SEEDS)
    assert set(FAMILIES).isdisjoint(RESERVED_FAMILIES)
    rows=[];elapsed={arm:0.0 for arm in ARMS[:-1]}
    work={arm:0 for arm in ARMS[:-1]}
    for seed in SEEDS:
      for fam in FAMILIES:
       for length in LENGTHS:
        private,public,start,chain=source(seed,fam,length)
        truth=prior.follow(private,start,chain)
        assert len(prior.worlds(public))==1296
        answers={"private_oracle_unfair":truth}
        receipts={};efforts={}
        for method in ARMS[:-1]:
            t0=time.perf_counter()
            observed=list(public)
            port=prior.PurchasePort(private,public)
            effort=[0]
            for step in range(2):
                selected=choose(method,observed,start,chain,2-step,effort)
                if selected is None:break
                bought=port.buy(*selected)
                observed.append((*selected,bought))
            answers[method]=prior.predict(observed,start,chain)
            receipts[method]=[list(x) for x in port.receipts]
            efforts[method]=effort[0]
            work[method]+=effort[0]
            elapsed[method]+=time.perf_counter()-t0
        rows.append(dict(seed=seed,family=fam,length=length,truth=truth,
                         answers=answers,receipts=receipts,work=efforts,
                         initial_posterior_worlds=1296))
    scores={arm:{"correct":sum(r["answers"][arm]==r["truth"] for r in rows),
                 "total":len(rows)} for arm in ARMS}
    for v in scores.values():v["accuracy"]=v["correct"]/v["total"]
    by_family={str(f):{a:sum(r["answers"][a]==r["truth"] for r in rows
                  if r["family"]==f) for a in ARMS} for f in FAMILIES}
    controls=ARMS[1:-1]
    baseline=max(scores[a]["accuracy"] for a in controls)
    gain=scores["two_step_counterfactual"]["accuracy"]-baseline
    no_regression=all(by_family[str(f)]["two_step_counterfactual"]>=max(
         by_family[str(f)][a] for a in controls) for f in FAMILIES)
    # Work is measured as conditional posterior-world contributions; runtime
    # measured separately. Do not treat an unequal-compute comparison as a win.
    comparator=max(controls,key=lambda a:scores[a]["accuracy"])
    reference_work=max(work[comparator],1)
    parity=work["two_step_counterfactual"]<=1.2*reference_work
    passed=gain>=.03 and no_regression and parity
    return dict(schema="wingless.f44.two-read-path-counterfactual.trainonly.v1",
      stage="FRESH_SYNTHETIC_TRAIN_ONLY_NOT_ORIGINAL",
      seeds=[SEEDS.start,SEEDS.stop],families=list(FAMILIES),
      reserved_seeds=[RESERVED_SEEDS.start,RESERVED_SEEDS.stop],
      reserved_families=list(RESERVED_FAMILIES),questions=len(rows),
      first_stage_beam=BEAM,max_paid_mappings=2,
      initial_posterior_worlds=1296,scores=scores,by_family=by_family,
      posterior_world_work=work,elapsed_policy_seconds=elapsed,
      strongest_fair_control=comparator,candidate_gain=gain,
      family_nondegradation=no_regression,
      measured_work_budget_parity=parity,rows=rows,
      verdict="TRAIN_ONLY_TWO_STEP_QUALIFIED" if passed else "TRAIN_ONLY_TWO_STEP_NOT_QUALIFIED",
      source_independent_F3_custody=False,primary_original_not_executed=True,
      native_576_byte_residency_qualified=False,accepted_ref_mutation=False,paid_api_calls=0)

if __name__=="__main__":
    print(json.dumps(run(),indent=2,sort_keys=True))
