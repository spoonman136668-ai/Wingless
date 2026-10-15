"""W-F54: can one PUBLIC calibration query identify a hidden operator law?

FIT chooses calibrated noise parameter. New diagnostic family/source identities;
target truth/private law used only to generate/evaluate observations, NEVER
by the prospective public query policy. Exactly one matched paid observation.
"""
import itertools,json,math,random
from pathlib import Path
FIT=range(601000,601008);DIAG=range(601100,601116)
FIT_FAM=(1200,1201,1202,1203);DIAG_FAM=(1204,1205,1206,1207)
RESERVED=range(602000,602016)
LAWS=("parity","majority","xnor","exactly_one","xor_and","or_not","threshold_two","impl")
CONTEXTS=tuple(itertools.product((0,1),repeat=3))
ARMS=("public_information_gain_purchase","public_target_nearest_purchase",
 "public_fixed_context_purchase","public_no_purchase","private_sign_oracle_UNFAIR")
def law(name,c):
 a,b,d=c
 return {"parity":a^b^d,"majority":int(a+b+d>=2),
 "xnor":int(not (a^b)),"exactly_one":int(a+b+d==1),
 "xor_and":(a^b)&d,"or_not":int(a or not d),
 "threshold_two":int(a+b+d==2),"impl":int((not a) or b)}[name]
def source(seed,family,sample):
 if not ((seed in FIT and family in FIT_FAM) or
         (seed in DIAG and family in DIAG_FAM)) or sample not in (0,1):
  raise ValueError("only preregistered FIT/DIAG, not originals")
 rng=random.Random(seed*7331+family*113+sample*151+5401)
 hidden=LAWS[family-1200]  # strictly source generation, not query selection
 candidates=list(CONTEXTS);rng.shuffle(candidates)
 initial=tuple(candidates[:3]);target=candidates[3]
 def observe(context):
  # One observation per context: record all possible outcomes prospectively,
  # but only reveal one chosen observation to the decision policy.
  independent=random.Random(seed*319991+family*773+sample*2707+
                            sum(x*(101**i) for i,x in enumerate(context))+54019)
  flipped=independent.random() < (.12 if seed in FIT else .2)
  return int(bool(law(hidden,context))^flipped)
 public=tuple((ctx,observe(ctx)) for ctx in initial)
 pool=tuple(ctx for ctx in CONTEXTS if ctx not in initial and ctx!=target)
 return dict(seed=seed,family=family,sample=sample,
  target=target,initial=public,
  observation_bank={''.join(map(str,c)):observe(c) for c in pool},
  true_target=law(hidden,target),hidden_law_EVALUATOR_ONLY=hidden)
def likelihood(anchors,eps):
 weights=[]
 for name in LAWS:
  w=1.
  for c,b in anchors:
   w*=1-eps if law(name,c)==b else eps
  weights.append(w)
 return tuple(w/sum(weights) for w in weights)
def public_predict(target,anchors,eps):
 posterior=likelihood(anchors,eps)
 p=sum(w*law(name,target) for w,name in zip(posterior,LAWS))
 return int(p>=.5)
def select(arm,target,anchors,eps):
 pool=tuple(c for c in CONTEXTS if c not in {a[0] for a in anchors} and c!=target)
 if arm=="public_no_purchase":return None
 if arm=="public_fixed_context_purchase":return min(pool)
 if arm=="public_target_nearest_purchase":
  return min(pool,key=lambda c:(sum(a!=b for a,b in zip(c,target)),c))
 if arm=="public_information_gain_purchase":
  posterior=likelihood(anchors,eps)
  # Expected Bernoulli entropy information maximized at balanced prospective
  # public reading, using ONLY the public law posterior and candidate context.
  return min(pool,key=lambda c:(-sum(w*law(n,c) for w,n in zip(posterior,LAWS))*
                     (1-sum(w*law(n,c) for w,n in zip(posterior,LAWS))),c))
 raise ValueError("unregistered public observation policy")
def fit():
 examples=[source(seed,f,s) for seed in FIT for f in FIT_FAM for s in (0,1)]
 errors={}
 for eps in (.08,.12,.16,.22):
  loss=0.
  for e in examples:
   p=sum(w*law(n,e["target"]) for w,n in
         zip(likelihood(e["initial"],eps),LAWS))
   p=min(1-1e-9,max(1e-9,p))
   loss-=math.log(p if e["true_target"] else 1-p)
  errors[eps]=loss
 return min(errors,key=lambda v:(errors[v],v)),len(examples)
def run():
 eps,nfit=fit()
 rows=[]
 for seed in DIAG:
  for family in DIAG_FAM:
   for sample in (0,1):
    e=source(seed,family,sample);out={};choices={}
    for arm in ARMS[:-1]:
     choice=select(arm,e["target"],e["initial"],eps)
     extra=None if choice is None else e["observation_bank"][''.join(map(str,choice))]
     augmented=e["initial"] if choice is None else (*e["initial"],(choice,extra))
     out[arm]=public_predict(e["target"],augmented,eps)
     choices[arm]=None if choice is None else list(choice)
    out["private_sign_oracle_UNFAIR"]=e["true_target"]
    rows.append(dict(seed=seed,family=family,sample=sample,target=list(e["target"]),
      initial_anchors=[[list(c),v] for c,v in e["initial"]],
      PUBLIC_extra_observations=e["observation_bank"],
      observed_contexts=choices,predictions=out,truth=e["true_target"]))
 scores={a:dict(correct=sum(r["predictions"][a]==r["truth"] for r in rows),
                    total=len(rows)) for a in ARMS}
 for v in scores.values():v["accuracy"]=v["correct"]/v["total"]
 by={str(f):{a:sum(r["predictions"][a]==r["truth"] for r in rows
                   if r["family"]==f) for a in ARMS} for f in DIAG_FAM}
 candidate=ARMS[0];fair=ARMS[1:-1];best=max(fair,key=lambda a:scores[a]["correct"])
 gain=scores[candidate]["accuracy"]-scores[best]["accuracy"]
 no_reg=all(by[str(f)][candidate]>=max(by[str(f)][a] for a in fair)
            for f in DIAG_FAM)
 passed=gain>=.05 and no_reg
 return dict(schema="wingless.f54.active-calibration-identifiability.trainonly.v1",
  stage="PUBLIC_ONCE_PAID_CALIBRATION_TRAIN_ONLY_NOT_ORIGINAL",
  fit_seeds=[FIT.start,FIT.stop],fit_families=list(FIT_FAM),
  diagnostic_seeds=[DIAG.start,DIAG.stop],diagnostic_families=list(DIAG_FAM),
  original_seeds_NOT_RESERVED=[RESERVED.start,RESERVED.stop],
  fitted_public_noise=eps,train_examples=nfit,public_initial_anchors=3,
  extra_calibration_max=1,cases=len(rows),rows=rows,scores=scores,
  by_family=by,strongest_fair_control=best,candidate_gain=gain,
  per_family_nondegradation=no_reg,frozen_gate_met=passed,
  verdict="TRAIN_ONLY_IDENTIFIABILITY_SUPPORTED_NOT_F3" if passed
   else "TRAIN_ONLY_IDENTIFIABILITY_NOT_SUPPORTED",
  private_source_law_not_selector_input=True,source_independent_F3=False,
  original_executed=False,accepted_ref_mutation=False,paid_api_calls=0)
if __name__=="__main__":print(json.dumps(run(),sort_keys=True))
