"""Standalone F53 source-isolated raw TRAIN and public-relation evidence audit."""
import json,math,sys
def compose(pub,chain,start,sign):
 shifts={r:(y-sign*x)%7 for r,x,y in pub};at=start
 for r,d in chain:
  at=(sign*at+shifts[r])%7 if d==1 else (sign*(at-shifts[r]))%7
 return at
def check(x,p):
 assert x["schema"]=="wingless.f53.group-minimax-public-law.trainonly.v1"
 assert x["stage"]=="TRAIN_ONLY_UNSEEN_CONTEXT_OPERATOR_LAWS_NOT_ORIGINAL"
 assert x["fit_seeds"]==p["fit_seeds"] and x["diagnostic_seeds"]==p["diagnostic_seeds"]
 assert x["fit_families"]==p["fit_families"] and x["diagnostic_families"]==p["diagnostic_families"]
 assert x["original_seeds_NOT_RESERVED"]==p["proposed_original_seeds_NOT_RESERVED"]
 assert x["source_independent_F3"] is False and x["original_executed"] is False
 assert x["accepted_ref_mutation"] is False and x["paid_api_calls"]==0
 assert x["diagnostic_calibration_noise"]==p["diagnostic_calibration_noise"]
 arms=[p["candidate"],*p["controls"],p["unfair_control"]]
 assert len(set(arms))==len(arms) and set(x["scores"])==set(arms)
 assert x["train_examples"]==64 and x["cases"]==128 and len(x["rows"])==128
 assert x["fit_group_params"][0] in (0.,.08,.16,.24,.32,.40)
 assert x["fit_group_params"][1] in ("nearest","support","majority")
 assert set(x["fit_family_correct"])==set(map(str,p["fit_families"]))
 seen=set();counts={a:0 for a in arms}
 for row in x["rows"]:
  key=row["seed"],row["family"],row["length"]
  assert key not in seen and 581100<=key[0]<581116 and key[1] in p["diagnostic_families"] and key[2] in (3,5)
  seen.add(key)
  assert len(row["public"])==4 and len(row["calibration"])==5
  assert len(row["chain"])==row["length"] and len(row["context"])==3
  assert set(row["chosen_signs"])==set(arms[:-1]) and set(row["answers"])==set(arms)
  for arm,sgn in row["chosen_signs"].items():
   assert sgn in (-1,1) and compose(row["public"],row["chain"],row["start"],sgn)==row["answers"][arm]
  assert compose(row["public"],row["chain"],row["start"],row["evaluator_private_sign"])==row["truth"]
  assert row["answers"][p["unfair_control"]]==row["truth"]
  for a in arms:counts[a]+=int(row["answers"][a]==row["truth"])
 assert len(seen)==128
 for a in arms:
  assert x["scores"][a]["correct"]==counts[a] and x["scores"][a]["total"]==128
  assert math.isclose(x["scores"][a]["accuracy"],counts[a]/128,abs_tol=1e-12)
 by={}
 for f in p["diagnostic_families"]:
  sub=[r for r in x["rows"] if r["family"]==f];assert len(sub)==32
  by[str(f)]={a:sum(r["answers"][a]==r["truth"] for r in sub) for a in arms}
  assert x["by_family"][str(f)]==by[str(f)]
 fair=p["controls"];winner=max(fair,key=lambda a:counts[a])
 gain=(counts[p["candidate"]]-counts[winner])/128
 nonreg=all(by[str(f)][p["candidate"]]>=max(by[str(f)][a] for a in fair)
            for f in p["diagnostic_families"])
 qualifies=gain>=p["threshold_gain_pp"]/100 and nonreg
 assert x["strongest_fair_control"]==winner
 assert math.isclose(x["candidate_gain"],gain,abs_tol=1e-12)
 assert x["per_family_nondegradation"]==nonreg and x["frozen_gate_met"]==qualifies
 return {"candidate_correct":counts[p["candidate"]],"strongest_fair":winner,
         "strongest_correct":counts[winner],"gain_pp":round(100*gain,5),
         "no_family_regression":nonreg,"passed_frozen_train":qualifies}
if __name__=="__main__":
 assert len(sys.argv)==3
 print("W_F53_INDEPENDENT_RECOUNT_PASS",json.dumps(check(json.load(open(sys.argv[1])),json.load(open(sys.argv[2]))),sort_keys=True))
