"""Independent frozen F51 Wingless row recount from existing source output.

Audits the public calibration sign derived from paired values and algebraic
composition. Does not call source_pilot.run or select a new experiment.
"""
import json,math,sys
def audit(path,prereg):
 d=json.load(open(path));p=json.load(open(prereg))
 v=p["diagnostic"];t=p["fit"];arms=p["controls"]
 assert d["schema"]=="wingless.f51.calibrated-unseen-operator.trainonly.v1"
 assert d["stage"]=="TRAIN_ONLY_UNSEEN_CONTEXT_OPERATOR_LAWS_NOT_ORIGINAL"
 assert d["fit_seeds"]==t["seeds"] and d["fit_families"]==t["families"]
 assert d["diagnostic_seeds"]==v["seeds"] and d["diagnostic_families"]==v["families"]
 assert set(t["families"]).isdisjoint(v["families"])
 assert d["cases"]==len(d["rows"])==v["expected_cases"]==128
 assert d["train_examples"]==t["fit_cases"]==64
 assert d["original_seeds_NOT_RESERVED"]==p["future_original_NOT_RESERVED"]["seeds"]
 assert len(d["learned_attention_weights"])==3
 assert d["diagnostic_private_sign_never_selector_input"]
 assert not d["source_independent_F3"] and not d["original_executed"]
 assert not d["native_576_bytes_residency_qualified"] and not d["accepted_ref_mutation"]
 assert d["paid_api_calls"]==0
 count={a:0 for a in arms};seen=set()
 for row in d["rows"]:
  key=row["seed"],row["family"],row["length"]
  assert key not in seen;seen.add(key)
  assert 521100<=row["seed"]<521116 and row["family"] in v["families"]
  assert row["length"] in v["task_lengths"]
  assert len(row["context"])==3 and all(c in (0,1) for c in row["context"])
  assert len(row["public"])==4 and len({r for r,x,y in row["public"]})==4
  assert len(row["calibration"])==4 and len(row["chain"])==row["length"]
  assert set(row["answers"])==set(arms)
  assert set(row["chosen_signs"])==set(arms)-{"private_sign_oracle_unfair"}
  calibrated=[]
  for ctx,left,right in row["calibration"]:
   assert len(ctx)==3 and all(b in (0,1) for b in ctx)
   x,y=left;z,t=right
   assert x!=z and 0<=x<7 and 0<=y<7 and 0<=z<7 and 0<=t<7
   sign=(t-y)*pow((z-x)%7,-1,7)%7
   assert sign in (1,6)
   calibrated.append((ctx,-1 if sign==6 else 1))
  for arm,sign in row["chosen_signs"].items():
   assert sign in (-1,1)
   if arm=="public_three_bit_parity":
    assert sign==(-1 if row["context"][0]^row["context"][1]^row["context"][2] else 1)
   if arm=="public_three_bit_majority":
    assert sign==(-1 if sum(row["context"])>=2 else 1)
   if arm=="always_shift":assert sign==1
   if arm=="always_reflect":assert sign==-1
   if arm=="public_support_majority":
    assert sign==(-1 if sum(s==-1 for c,s in calibrated)>=3 else 1)
   if arm=="public_hamming_nearest":
    near=min(calibrated,key=lambda item:(sum(a!=b for a,b in zip(row["context"],item[0])),item[0]))
    assert sign==near[1]
   if arm=="learned_support_attention":
    weights=d["learned_attention_weights"]
    votes=[]
    for ctx,s in calibrated:
     z=sum(w*(1 if a==b else -1) for a,b,w in zip(row["context"],ctx,weights))
     votes.append((math.exp(max(-15,min(15,z))),s))
    pminus=sum(w for w,s in votes if s==-1)/sum(w for w,s in votes)
    assert sign==(-1 if pminus>=.5 else 1)
   shifts={r:(y-sign*x)%7 for r,x,y in row["public"]}
   cur=row["start"]
   for role,direction in row["chain"]:
    if direction==1:cur=(sign*cur+shifts[role])%7
    else:cur=(sign*(cur-shifts[role]))%7
   assert row["answers"][arm]==cur
  true=row["evaluator_private_sign"]
  assert true in (-1,1)
  true_shifts={r:(y-true*x)%7 for r,x,y in row["public"]}
  value=row["start"]
  for role,direction in row["chain"]:
   if direction==1:value=(true*value+true_shifts[role])%7
   else:value=(true*(value-true_shifts[role]))%7
  assert row["truth"]==value
  assert row["answers"]["private_sign_oracle_unfair"]==value
  for a in arms:count[a]+=int(row["answers"][a]==value)
 assert len(seen)==128
 for a in arms:
  assert d["scores"][a]["correct"]==count[a] and d["scores"][a]["total"]==128
  assert math.isclose(d["scores"][a]["accuracy"],count[a]/128,abs_tol=1e-12)
 for fam in v["families"]:
  sub=[row for row in d["rows"] if row["family"]==fam]
  assert len(sub)==32
  for arm in arms:assert d["by_family"][str(fam)][arm]==sum(
   row["answers"][arm]==row["truth"] for row in sub)
 fair=arms[1:-1];best=max(fair,key=lambda a:count[a])
 gain=(count["learned_support_attention"]-count[best])/128
 no_reg=all(d["by_family"][str(f)]["learned_support_attention"]>=max(
   d["by_family"][str(f)][a] for a in fair) for f in v["families"])
 passed=gain>=p["thresholds"]["gain_vs_strongest_fair_min"] and no_reg
 assert d["strongest_fair_control"]==best
 assert math.isclose(d["candidate_gain"],gain,abs_tol=1e-12)
 assert d["per_family_nondegradation"]==no_reg
 assert d["frozen_gate_met"]==passed
 assert d["verdict"]==("TRAIN_ONLY_UNSEEN_OPERATOR_SUPPORTED" if passed else
                       "TRAIN_ONLY_UNSEEN_OPERATOR_NEGATIVE")
 print("W_F51",d["verdict"],"correct",count,"gain",gain,"family_guard",no_reg)
if __name__=="__main__":
 assert len(sys.argv)==3
 audit(sys.argv[1],sys.argv[2])
