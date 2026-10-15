"""W-F45 source-only selective, compute-capped two-read counterfactual research.

This is a new synthetic TRAIN-only source, not a first-use original or F3.
All policy choices use public mappings and posterior completions, never hidden
truth or real source maps. Native <=576B residency is not qualified.
"""
from __future__ import annotations
import hashlib, importlib.util, json, math, random, sys, time
from collections import Counter
from pathlib import Path

HERE = Path(__file__).resolve().parent
P = HERE.parent / "20261015-w-f44" / "source_pilot.py"
spec = importlib.util.spec_from_file_location("w_f45_f44_controls", P)
prior = importlib.util.module_from_spec(spec)
sys.modules[spec.name] = prior
spec.loader.exec_module(prior)

SEEDS = range(371000, 371012)
FAMILIES = (460, 461, 462)
LENGTHS = (3, 5)
RESERVED = range(372000, 372016)
RESERVED_FAMILIES = (464, 465, 466, 467)
ARMS = ("selective_budgeted_lookahead", "one_step_value", "exact_entropy",
        "query_frontier", "full_two_step_costly_reference", "no_purchase",
        "evaluator_oracle_unfair")
FIRST_BEAM = 4
MAX_DEEP_FRAC = .10
MAX_DEEP_ABSOLUTE = 3

def source(seed, family, length):
    if seed not in SEEDS or family not in FAMILIES or length not in LENGTHS:
        raise ValueError("reserved or unregistered source forbidden")
    rng = random.Random(seed*17713 + family*3037 + length*521 + 45001)
    private = []
    public = []
    for role in range(4):
        perm = list(range(6))
        rng.shuffle(perm)
        private.append(tuple(perm))
        for x in rng.sample(range(6), 3):
            public.append((role, x, perm[x]))
    rng.shuffle(public)
    templates = {
        460: ((0,1),(1,-1),(0,-1),(2,1),(3,-1)),
        461: ((3,-1),(2,1),(2,-1),(1,1),(0,-1)),
        462: ((1,1),(0,1),(1,-1),(0,-1),(3,1)),
    }
    return tuple(private), tuple(public), rng.randrange(6), templates[family][:length]

def one_step(observed, start, chain, effort):
    bank = prior.prior.worlds(observed)
    actions = prior.prior.candidate_actions(observed)
    if not actions: return None, 0.
    base = prior.best_success(bank, start, chain)
    scores = {a: prior.split_success(bank, a, start, chain) for a in actions}
    effort[0] += len(bank)*len(actions)
    order = sorted(actions, key=lambda a:(-scores[a],a[0],a[1]))
    best = order[0]
    margin = (scores[best]-scores[order[1]])/len(bank) if len(order)>1 else 1.
    return (best if scores[best]>base else None), margin

def run():
    assert set(SEEDS).isdisjoint(RESERVED)
    assert set(FAMILIES).isdisjoint(RESERVED_FAMILIES)
    rows = []
    work = {a:0 for a in ARMS[:-1]}
    elapsed = {a:0. for a in ARMS[:-1]}
    deep_count = 0
    deep_flags = []
    for seed in SEEDS:
      for family in FAMILIES:
       for length in LENGTHS:
        private, public, start, chain = source(seed, family, length)
        truth = prior.prior.follow(private, start, chain)
        assert len(prior.prior.worlds(public))==1296
        answers={"evaluator_oracle_unfair":truth}
        receipts={};per_row_work={}
        decision_flags={}
        for method in ARMS[:-1]:
            began=time.perf_counter()
            observed=list(public)
            port=prior.prior.PurchasePort(private,public)
            meter=[0]
            for step in range(2):
                if method=="no_purchase":
                    selected=None
                elif method=="selective_budgeted_lookahead":
                    selected,margin=one_step(observed,start,chain,meter)
                    if step==0:
                        eligible=length==5 and margin<.05
                        used_deep=eligible and deep_count<MAX_DEEP_ABSOLUTE
                        decision_flags=dict(eligible=eligible,used_deep=used_deep,
                                            one_step_rank_margin=margin)
                        if used_deep:
                            # Count all base and counterfactual work; a conditional
                            # first-stage branch never receives hidden values.
                            extra=[0]
                            deeper=prior.choose("two_step_counterfactual",observed,start,chain,2,extra)
                            meter[0]+=extra[0]
                            deep_count+=1
                            if deeper is not None: selected=deeper
                    # Second purchase uses the public posterior after the first.
                elif method=="full_two_step_costly_reference":
                    selected=prior.choose("two_step_counterfactual",observed,start,chain,2-step,meter)
                else:
                    selected=prior.choose(method,observed,start,chain,2-step,meter)
                if selected is None:break
                value=port.buy(*selected)
                observed.append((*selected,value))
            answers[method]=prior.prior.predict(observed,start,chain)
            receipts[method]=[list(x) for x in port.receipts]
            per_row_work[method]=meter[0]
            work[method]+=meter[0]
            elapsed[method]+=time.perf_counter()-began
        deep_flags.append(decision_flags.get("used_deep",False))
        rows.append(dict(seed=seed,family=family,length=length,truth=truth,
                  answers=answers,receipts=receipts,work=per_row_work,
                  deep_gate=decision_flags,initial_posterior_worlds=1296))
    scores={a:dict(correct=sum(r["answers"][a]==r["truth"] for r in rows),
                   total=len(rows)) for a in ARMS}
    for s in scores.values():s["accuracy"]=s["correct"]/s["total"]
    by_family={str(f):{a:sum(r["answers"][a]==r["truth"] for r in rows if r["family"]==f)
                       for a in ARMS} for f in FAMILIES}
    budget_fair=("one_step_value","exact_entropy","query_frontier")
    strongest=max(budget_fair,key=lambda a:scores[a]["accuracy"])
    gain=scores["selective_budgeted_lookahead"]["accuracy"]-scores[strongest]["accuracy"]
    family_no_regress=all(by_family[str(f)]["selective_budgeted_lookahead"]>=
          max(by_family[str(f)][a] for a in budget_fair) for f in FAMILIES)
    parity=work["selective_budgeted_lookahead"]<=1.2*max(work[strongest],1)
    deep_frac=deep_count/len(rows)
    qualified=gain>=.03 and family_no_regress and parity and deep_frac<=MAX_DEEP_FRAC
    return dict(schema="wingless.f45.cost-capped-selective-lookahead.trainonly.v1",
      stage="SYNTHETIC_FRESH_TRAIN_ONLY_NO_ORIGINAL",seeds=[SEEDS.start,SEEDS.stop],
      families=list(FAMILIES),reserved_seeds=[RESERVED.start,RESERVED.stop],
      reserved_families=list(RESERVED_FAMILIES),questions=len(rows),
      initial_posterior_worlds=1296,first_stage_beam=FIRST_BEAM,
      max_paid_mappings=2,deep_gate_count=deep_count,deep_gate_fraction=deep_frac,
      posterior_world_work=work,elapsed_policy_seconds=elapsed,
      strongest_budget_fair_control=strongest,candidate_gain=gain,
      family_nondegradation=family_no_regress,posterior_work_parity=parity,
      scores=scores,by_family=by_family,rows=rows,
      verdict="TRAIN_ONLY_COST_GATED_LOOKAHEAD_QUALIFIED" if qualified else
              "TRAIN_ONLY_COST_GATED_LOOKAHEAD_NOT_QUALIFIED",
      original_executed=False,independent_F3_custody=False,
      native_576_byte_resident_qualified=False,
      accepted_ref_mutation=False,paid_api_calls=0)

if __name__=="__main__":print(json.dumps(run(),sort_keys=True,indent=2))
