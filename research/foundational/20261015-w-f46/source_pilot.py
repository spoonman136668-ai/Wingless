"""W-F46 new synthetic TRAIN-only experiment; no original or F3."""
import importlib.util,json,random,sys
from pathlib import Path
HERE=Path(__file__).resolve().parent
P=HERE.parent/"20261015-w-f45"/"source_pilot.py"
spec=importlib.util.spec_from_file_location("w_f46_parent_controls",P)
f45=importlib.util.module_from_spec(spec);sys.modules[spec.name]=f45;spec.loader.exec_module(f45)
FIT=range(391000,391012)
DIAG=range(391100,391112)
ORIGINAL=range(392000,392016)
OLD_ORIGINAL=range(372000,372016)
FAMILIES=(500,501,502,503)
ARMS=("query_edge_voc_per_compute","one_step_value","exact_entropy","query_frontier","budget_matched_rank","no_purchase","evaluator_oracle_unfair")
def source(seed,family,length):
 if seed not in DIAG or family not in FAMILIES or length not in (3,5):
  raise ValueError("protected or unregistered source")
 rng=random.Random(seed*17713+family*3037+length*521+46001)
 private=[];public=[]
 for role in range(4):
  perm=list(range(6));rng.shuffle(perm);private.append(tuple(perm))
  for x in rng.sample(range(6),3):public.append((role,x,perm[x]))
 rng.shuffle(public)
 templates={
  500:((0,1),(2,-1),(0,-1),(3,1),(1,-1)),
  501:((2,-1),(0,1),(3,-1),(1,1),(2,1)),
  502:((1,1),(3,1),(2,-1),(1,-1),(0,1)),
  503:((3,-1),(1,1),(3,1),(0,-1),(2,1))}
 return tuple(private),tuple(public),rng.randrange(6),templates[family][:length]
def select(observed,start,chain,work,weighted):
 bank=f45.prior.prior.worlds(observed)
 choices=f45.prior.prior.candidate_actions(observed)
 if not choices:return None
 baseline=f45.prior.best_success(bank,start,chain)
 success={a:f45.prior.split_success(bank,a,start,chain) for a in choices}
 work[0]+=len(bank)*len(choices)
 counts={r:sum(1 for a,b in chain if a==r) for r in range(4)}
 def score(a):
  gain=max(0,success[a]-baseline)
  role=counts[a[0]]
  # Public-only query-edge relevance vs public posterior cost estimate.
  return gain*(1+role)/(len(bank)*(1+max(0,2-role))) if weighted else gain/len(bank)
 ranked=sorted(choices,key=lambda a:(-score(a),a))
 return ranked[0] if success[ranked[0]]>baseline else None
def run():
 assert not set(DIAG)&(set(FIT)|set(ORIGINAL)|set(OLD_ORIGINAL))
 rows=[];effort={a:0 for a in ARMS[:-1]}
 for seed in DIAG:
  for family in FAMILIES:
   for length in (3,5):
    private,public,start,chain=source(seed,family,length)
    truth=f45.prior.prior.follow(private,start,chain)
    answers={"evaluator_oracle_unfair":truth};receipts={};costs={}
    for arm in ARMS[:-1]:
     observed=list(public);port=f45.prior.prior.PurchasePort(private,public);work=[0]
     for step in range(2):
      if arm=="no_purchase":chosen=None
      elif arm in ("query_edge_voc_per_compute","budget_matched_rank"):
       chosen=select(observed,start,chain,work,arm=="query_edge_voc_per_compute")
      else:chosen=f45.prior.choose(arm,observed,start,chain,2-step,work)
      if chosen is None:break
      value=port.buy(*chosen);observed.append((*chosen,value))
     answers[arm]=f45.prior.prior.predict(observed,start,chain)
     receipts[arm]=[list(x) for x in port.receipts]
     assert len(receipts[arm])<=2
     costs[arm]=work[0];effort[arm]+=work[0]
    rows.append(dict(seed=seed,family=family,length=length,truth=truth,
             answers=answers,receipts=receipts,work=costs))
 scores={a:dict(correct=sum(r["answers"][a]==r["truth"] for r in rows),total=len(rows)) for a in ARMS}
 for v in scores.values():v["accuracy"]=v["correct"]/v["total"]
 fair=("one_step_value","exact_entropy","query_frontier","budget_matched_rank")
 strongest=max(fair,key=lambda a:(scores[a]["accuracy"],-effort[a]))
 cand="query_edge_voc_per_compute"
 gain=scores[cand]["accuracy"]-scores[strongest]["accuracy"]
 ratio=effort[cand]/max(1,effort[strongest])
 by_family={str(f):{a:sum(r["answers"][a]==r["truth"] for r in rows if r["family"]==f) for a in ARMS} for f in FAMILIES}
 no_regress=all(by_family[str(f)][cand]>=max(by_family[str(f)][a] for a in fair) for f in FAMILIES)
 passed=gain>=.03 and ratio<=1.2 and no_regress
 return dict(schema="wingless.f46.query-edge-voc.trainonly.v1",stage="TRAIN_ONLY_NOT_ORIGINAL",
   questions=len(rows),seeds=[DIAG.start,DIAG.stop],families=list(FAMILIES),
   fit_seeds_not_consumed=[FIT.start,FIT.stop],
   reserved_originals_not_consumed=[ORIGINAL.start,ORIGINAL.stop],
   old_originals_not_consumed=[OLD_ORIGINAL.start,OLD_ORIGINAL.stop],
   scores=scores,rows=rows,posterior_world_work=effort,by_family=by_family,
   strongest_budget_fair_control=strongest,candidate_gain=gain,
   posterior_work_ratio=ratio,family_nondegradation=no_regress,
   frozen_criterion_met=passed,
   verdict="TRAIN_ONLY_SUPPORTED" if passed else "TRAIN_ONLY_NEGATIVE",
   original_executed=False,independent_F3_custody=False,accepted_ref_mutation=False,
   native_576_byte_resident_qualified=False,paid_api_calls=0)
if __name__=="__main__":print(json.dumps(run(),sort_keys=True,indent=2))
