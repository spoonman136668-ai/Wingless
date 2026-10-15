"""W-F54 source-independent reproducible recalculation of all public choices."""
import itertools,json,math,sys
LAWS=("parity","majority","xnor","exactly_one","xor_and","or_not","threshold_two","impl")
CONTEXTS=tuple(itertools.product((0,1),repeat=3))
def law(name,c):
 a,b,d=c
 return {"parity":a^b^d,"majority":int(a+b+d>=2),"xnor":int(not(a^b)),
  "exactly_one":int(a+b+d==1),"xor_and":(a^b)&d,"or_not":int(a or not d),
  "threshold_two":int(a+b+d==2),"impl":int((not a) or b)}[name]
def posterior(anchors,eps):
 w=[math.prod((1-eps if law(n,c)==b else eps) for c,b in anchors) for n in LAWS]
 return [v/sum(w) for v in w]
def predict(target,anchors,eps):
 return int(sum(w*law(n,target) for w,n in zip(posterior(anchors,eps),LAWS))>=.5)
def choose(arm,target,anchors,eps):
 pool=[c for c in CONTEXTS if c not in [v[0] for v in anchors] and c!=target]
 if arm=="public_no_purchase":return None
 if arm=="public_fixed_context_purchase":return min(pool)
 if arm=="public_target_nearest_purchase":
  return min(pool,key=lambda c:(sum(a!=b for a,b in zip(c,target)),c))
 if arm=="public_information_gain_purchase":
  weights=posterior(anchors,eps)
  return min(pool,key=lambda c:(-(sum(w*law(n,c) for w,n in zip(weights,LAWS))*
  (1-sum(w*law(n,c) for w,n in zip(weights,LAWS)))),c))
 raise AssertionError("bad arm")
def check(x,p):
 assert x["schema"]=="wingless.f54.active-calibration-identifiability.trainonly.v1"
 assert x["stage"]=="PUBLIC_ONCE_PAID_CALIBRATION_TRAIN_ONLY_NOT_ORIGINAL"
 assert x["fit_seeds"]==p["fit_seeds"] and x["diagnostic_seeds"]==p["diagnostic_seeds"]
 assert x["diagnostic_families"]==p["diagnostic_families"]
 assert x["original_seeds_NOT_RESERVED"]==p["proposed_originals_NOT_RESERVED"]
 assert x["source_independent_F3"] is False and x["original_executed"] is False
 assert x["accepted_ref_mutation"] is False and x["paid_api_calls"]==0
 assert x["private_source_law_not_selector_input"] is True
 assert x["cases"]==128 and len(x["rows"])==128 and x["train_examples"]==64
 eps=x["fitted_public_noise"];assert eps in (.08,.12,.16,.22)
 arms=p["decision_arms"];assert set(x["scores"])==set(arms)
 seen=set();scores={a:0 for a in arms}
 for r in x["rows"]:
  k=r["seed"],r["family"],r["sample"]
  assert k not in seen and 601100<=k[0]<601116 and k[1] in p["diagnostic_families"] and k[2] in (0,1)
  seen.add(k)
  target=tuple(r["target"]);anchors=[(tuple(c),bit) for c,bit in r["initial_anchors"]]
  assert len(anchors)==3 and target not in [c for c,v in anchors] and len(set(c for c,v in anchors))==3
  assert set(r["observed_contexts"])==set(arms[:-1]) and set(r["predictions"])==set(arms)
  for a in arms[:-1]:
   c=choose(a,target,anchors,eps)
   assert r["observed_contexts"][a]==(None if c is None else list(c))
   more=anchors if c is None else anchors+[(c,r["PUBLIC_extra_observations"][''.join(map(str,c))])]
   assert predict(target,more,eps)==r["predictions"][a]
  assert r["predictions"][arms[-1]]==r["truth"]
  for a in arms:scores[a]+=int(r["predictions"][a]==r["truth"])
 for a in arms:assert x["scores"][a]["correct"]==scores[a] and x["scores"][a]["total"]==128
 by={}
 for f in p["diagnostic_families"]:
  sub=[r for r in x["rows"] if r["family"]==f];assert len(sub)==32
  by[str(f)]={a:sum(r["predictions"][a]==r["truth"] for r in sub) for a in arms}
  assert by[str(f)]==x["by_family"][str(f)]
 cand=arms[0];fair=p["fair_controls"];best=max(fair,key=lambda a:scores[a])
 gain=(scores[cand]-scores[best])/128
 no_reg=all(by[str(f)][cand]>=max(by[str(f)][a] for a in fair) for f in p["diagnostic_families"])
 passed=gain>=p["threshold_gain_pp"]/100 and no_reg
 assert x["strongest_fair_control"]==best and math.isclose(x["candidate_gain"],gain,abs_tol=1e-12)
 assert x["per_family_nondegradation"]==no_reg and x["frozen_gate_met"]==passed
 return dict(candidate=scores[cand],strongest_fair=best,fair_correct=scores[best],
  gain_pp=100*gain,no_family_regression=no_reg,criterion_met=passed)
if __name__=="__main__":
 assert len(sys.argv)==3
 print("W_F54_INDEPENDENT_RECOUNT_PASS",json.dumps(check(json.load(open(sys.argv[1])),json.load(open(sys.argv[2]))),sort_keys=True))
