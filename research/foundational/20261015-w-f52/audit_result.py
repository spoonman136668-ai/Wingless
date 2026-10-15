"""Independent W-F52 recount; never imports source or executes it."""
import json,math,sys
def compose(pub,chain,start,sign):
 shifts={r:(y-sign*x)%7 for r,x,y in pub}
 at=start
 for role,d in chain:
  at=(sign*at+shifts[role])%7 if d==1 else (sign*(at-shifts[role]))%7
 return at
def verify(d,p):
 assert d["schema"]=="wingless.f52.public-law-model-averaging.trainonly.v1"
 assert d["stage"]=="TRAIN_ONLY_UNSEEN_CONTEXT_OPERATOR_LAWS_NOT_ORIGINAL"
 assert d["fit_seeds"]==p["fit_seeds"] and d["diagnostic_seeds"]==p["diagnostic_seeds"]
 assert d["fit_families"]==p["fit_families"] and d["diagnostic_families"]==p["diagnostic_families"]
 assert d["original_seeds_NOT_RESERVED"]==p["future_originals_proposed_only"]
 assert d["source_independent_F3"] is False and d["original_executed"] is False
 assert d["accepted_ref_mutation"] is False and d["paid_api_calls"]==0
 arms=[p["candidate"],*p["baseline_controls"],p["unfair_control"]]
 assert len(arms)==len(set(arms))
 assert set(d["scores"])==set(arms)
 assert d["train_examples"]==64 and d["cases"]==128 and len(d["rows"])==128
 counts={arm:0 for arm in arms};seen=set()
 for row in d["rows"]:
  key=(row["seed"],row["family"],row["length"])
  assert key not in seen and 561100<=key[0]<561116
  assert key[1] in p["diagnostic_families"] and key[2] in (3,5)
  seen.add(key)
  assert len(row["context"])==3 and len(row["calibration"])==5
  assert len(row["public"])==4 and len(row["chain"])==row["length"]
  assert set(row["answers"])==set(arms)
  assert set(row["chosen_signs"])==set(arms)-{p["unfair_control"]}
  for anchor in row["calibration"]:
   c,(x,y),(z,t)=anchor
   assert len(c)==3 and x!=z and ((t-y)*pow((z-x)%7,-1,7))%7 in (1,6)
  for arm,sign in row["chosen_signs"].items():
   assert sign in (-1,1)
   assert compose(row["public"],row["chain"],row["start"],sign)==row["answers"][arm]
  assert compose(row["public"],row["chain"],row["start"],row["evaluator_private_sign"])==row["truth"]
  assert row["answers"][p["unfair_control"]]==row["truth"]
  for arm in arms:counts[arm]+=int(row["answers"][arm]==row["truth"])
 assert len(seen)==128
 for arm in arms:
  assert d["scores"][arm]["correct"]==counts[arm]
  assert d["scores"][arm]["total"]==128
  assert math.isclose(d["scores"][arm]["accuracy"],counts[arm]/128,abs_tol=1e-12)
 for f in p["diagnostic_families"]:
  sub=[r for r in d["rows"] if r["family"]==f];assert len(sub)==32
  for arm in arms:
   assert d["by_family"][str(f)][arm]==sum(r["answers"][arm]==r["truth"] for r in sub)
 fair=p["baseline_controls"];best=max(fair,key=lambda a:counts[a])
 gain=(counts[p["candidate"]]-counts[best])/128
 nonreg=all(d["by_family"][str(f)][p["candidate"]]>=max(d["by_family"][str(f)][a] for a in fair)
            for f in p["diagnostic_families"])
 passed=gain>=p["threshold_absolute_gain_pp"]/100 and nonreg
 assert d["strongest_fair_control"]==best
 assert math.isclose(d["candidate_gain"],gain,abs_tol=1e-12)
 assert d["per_family_nondegradation"]==nonreg and d["frozen_gate_met"]==passed
 return {"candidate":counts[p["candidate"]],"fair":counts[best],"strongest":best,
         "gain_pp":100*gain,"no_family_regression":nonreg,"threshold_met":passed}
if __name__=="__main__":
 assert len(sys.argv)==3
 d=json.load(open(sys.argv[1]));p=json.load(open(sys.argv[2]))
 print("W_F52_INDEPENDENT_RECOUNT_PASS",json.dumps(verify(d,p),sort_keys=True))
