"""W-F48 actual TRAIN-only learned public edge-value/cost rank study.

The beta parameter is selected exclusively by FIT-only downstream answer
success; diagnostic truth, hidden role maps and F3 originals are never inputs
to a policy. Extra computational work is counted but the native 576B physical
resident constraint is not claimed.
"""
import importlib.util,json,math,random,sys
from pathlib import Path
P=Path(__file__).resolve().parent.parent/"20261015-w-f46"/"source_pilot.py"
spec=importlib.util.spec_from_file_location("w_f48_parent_f46",P)
old=importlib.util.module_from_spec(spec);sys.modules[spec.name]=old;spec.loader.exec_module(old)
FIT=range(411000,411004)
DIAG=range(411100,411104)
ORIGINAL=range(412000,412016)
FAMILIES=(600,601,602,603)
ORIGINAL_FAMILIES=(604,605,606,607)
ARMS=("learned_query_cost_value","one_step_value","exact_entropy","query_frontier",
      "budget_matched_rank","no_purchase","evaluator_oracle_unfair")
BETAS=(-1.,-.5,0.,.5,1.,2.)
def source(seed,family,length):
 if seed not in FIT and seed not in DIAG:raise ValueError("unregistered or primary original source")
 if family not in FAMILIES or length not in (3,5):raise ValueError("unregistered family")
 rng=random.Random(seed*17713+family*3037+length*521+48001)
 private=[];public=[]
 for role in range(4):
  perm=list(range(6));rng.shuffle(perm);private.append(tuple(perm))
  for x in rng.sample(range(6),3):public.append((role,x,perm[x]))
 rng.shuffle(public)
 templates={
  600:((2,1),(0,-1),(1,1),(3,-1),(2,-1)),
  601:((3,-1),(0,1),(2,1),(2,-1),(1,1)),
  602:((1,1),(3,-1),(0,1),(1,-1),(2,1)),
  603:((0,-1),(2,1),(1,-1),(3,1),(0,1))}
 return tuple(private),tuple(public),rng.randrange(6),templates[family][:length]
def learned_action(observed,start,chain,beta,work):
 port=old.f45.prior.prior
 bank=port.worlds(observed)
 actions=port.candidate_actions(observed)
 if not actions:return None
 baseline=old.f45.prior.best_success(bank,start,chain)
 answers={a:old.f45.prior.split_success(bank,a,start,chain) for a in actions}
 work[0]+=len(bank)*len(actions)
 roles=[r for r,sign in chain]
 relevance={r:roles.count(r) for r in range(4)}
 def score(action):
  r,_=action
  gain=max(0,answers[action]-baseline)/len(bank)
  # Beta learned on disjoint FIT labels. At diagnostic time this uses
  # only public query roles, posterior work proxy and posterior gain.
  sparse_ratio=(4-relevance[r])/4
  estimated_cost=1+sparse_ratio
  return gain*(1+beta*relevance[r]/5)/estimated_cost
 best=min(actions,key=lambda a:(-score(a),a))
 return best if answers[best]>baseline else None
def simulate(seed,family,length,method,beta=0):
 truth_map,public,start,chain=source(seed,family,length)
 truth=old.f45.prior.prior.follow(truth_map,start,chain)
 observed=list(public)
 bought=old.f45.prior.prior.PurchasePort(truth_map,public)
 work=[0]
 if method=="evaluator_oracle_unfair":return truth,[],0,truth
 for step in range(2):
  if method=="no_purchase":selected=None
  elif method=="learned_query_cost_value":
   selected=learned_action(observed,start,chain,beta,work)
  elif method=="budget_matched_rank":
   selected=old.select(observed,start,chain,work,False)
  else:
   selected=old.f45.prior.choose(method,observed,start,chain,2-step,work)
  if selected is None:break
  value=bought.buy(*selected)
  observed.append((*selected,value))
 answer=old.f45.prior.prior.predict(observed,start,chain)
 return answer,[list(v) for v in bought.receipts],work[0],truth
def cases(seeds):
 for seed in seeds:
  for family in FAMILIES:
   for length in (3,5):yield seed,family,length
def fit():
 scores={}
 for beta in BETAS:
  correct=0
  for seed,family,length in cases(FIT):
   answer,_,_,truth=simulate(seed,family,length,"learned_query_cost_value",beta)
   correct+=answer==truth
  scores[str(beta)]=correct
 chosen=min(BETAS,key=lambda b:(-scores[str(b)],abs(b),b))
 return chosen,scores
def run():
 assert not (set(FIT)&set(DIAG) or set(FIT)&set(ORIGINAL) or set(DIAG)&set(ORIGINAL))
 beta,fit_scores=fit()
 rows=[];work={a:0 for a in ARMS[:-1]}
 for seed,family,length in cases(DIAG):
  answers={};receipts={};meter={}
  for arm in ARMS:
   answer,trace,cost,truth=simulate(seed,family,length,arm,beta)
   answers[arm]=answer
   if arm!="evaluator_oracle_unfair":
    receipts[arm]=trace;meter[arm]=cost;work[arm]+=cost
  rows.append(dict(seed=seed,family=family,length=length,truth=truth,
                   answers=answers,receipts=receipts,work=meter))
 scores={a:dict(correct=sum(r["answers"][a]==r["truth"] for r in rows),
                  total=len(rows)) for a in ARMS}
 for v in scores.values():v["accuracy"]=v["correct"]/v["total"]
 fair=("one_step_value","exact_entropy","query_frontier","budget_matched_rank")
 strongest=max(fair,key=lambda a:(scores[a]["accuracy"],-work[a]))
 candidate="learned_query_cost_value"
 gain=scores[candidate]["accuracy"]-scores[strongest]["accuracy"]
 ratio=work[candidate]/max(1,work[strongest])
 by_family={str(f):{a:sum(r["answers"][a]==r["truth"] for r in rows if r["family"]==f)
                    for a in ARMS} for f in FAMILIES}
 no_family_regression=all(by_family[str(f)][candidate]>=max(by_family[str(f)][a] for a in fair)
                          for f in FAMILIES)
 qualified=gain>=.05 and ratio<=1.2 and no_family_regression
 return dict(schema="wingless.f48.disjoint-trained-query-value.trainonly.v1",
  stage="ONE_FRESH_FIT_AND_DIAGNOSTIC_NO_PRIMARY_ORIGINAL",
  fit_seeds=[FIT.start,FIT.stop],diagnostic_seeds=[DIAG.start,DIAG.stop],
  families=list(FAMILIES),originals_reserved_not_consumed=[ORIGINAL.start,ORIGINAL.stop],
  original_families_reserved_not_consumed=list(ORIGINAL_FAMILIES),
  questions=len(rows),training_questions=len(list(cases(FIT))),
  training_beta_scores=fit_scores,selected_beta_train_only=beta,
  diagnostic_truth_not_used_for_fit=True,public_only_action_inputs=True,
  scores=scores,rows=rows,posterior_world_work=work,by_family=by_family,
  strongest_budget_fair_control=strongest,candidate_gain=gain,
  posterior_work_ratio=ratio,family_nondegradation=no_family_regression,
  frozen_criterion_met=qualified,
  verdict="TRAIN_ONLY_SUPPORTED_NOT_F3" if qualified else "TRAIN_ONLY_NEGATIVE",
  original_executed=False,independent_F3_custody=False,
  accepted_ref_mutation=False,native_576_byte_resident_qualified=False,paid_api_calls=0)
if __name__=="__main__":print(json.dumps(run(),sort_keys=True,indent=2))
