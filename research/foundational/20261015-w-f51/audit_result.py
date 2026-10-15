"""Independent W-F51 per-row public role composition and frozen-control audit."""
import json,math,sys
def run(report,prereg):
 d=json.load(open(report));p=json.load(open(prereg))
 assert d["schema"]=="wingless.f51.five-bit-public-operator.trainonly.v1"
 assert d["stage"]=="DISJOINT_FIT_AND_NOVEL_COMPOSITION_TRAIN_ONLY"
 assert d["fit_seeds"]==p["fit_seeds"] and d["fit_families"]==p["fit_families"]
 assert d["diagnostic_seeds"]==p["diagnostic_seeds"]
 assert d["diagnostic_families"]==p["diagnostic_families"]
 assert d["proposed_original_seeds_NOT_RESERVED"]==p["future_original_seeds_NOT_RESERVED"]
 assert d["train_examples"]==64 and d["cases"]==128 and len(d["rows"])==128
 assert d["source_independent_F3"] is False and d["original_executed"] is False
 assert d["accepted_ref_mutation"] is False and d["paid_api_calls"]==0
 assert set(d["scores"])==set(p["controls"])
 assert len(d["posterior_minus_by_public_context"])==32
 controls=p["controls"];counts={a:0 for a in controls};seen=set()
 for r in d["rows"]:
  key=r["seed"],r["family"],r["length"]
  assert key not in seen and 541100<=r["seed"]<541116 and r["family"] in p["diagnostic_families"]
  seen.add(key)
  assert len(r["context"])==5 and all(x in (0,1) for x in r["context"])
  assert len(r["public"])==4 and len({x[0] for x in r["public"]})==4
  assert set(r["answers"])==set(controls)
  assert set(r["chosen_signs"])==set(controls)-{"private_operator_oracle_unfair"}
  bits=r["context"]
  sign=-1 if ((bits[0]&bits[1])^(bits[2]&bits[3])^bits[4]) else 1
  assert r["chosen_signs"]["full_five_bit_public_nonlinear_rule"]==sign
  for arm,g in r["chosen_signs"].items():
   assert g in (-1,1)
   offsets={role:(y-g*x)%7 for role,x,y in r["public"]}
   at=r["start"]
   for role,direction in r["chain"]:
    at=(g*at+offsets[role])%7 if direction==1 else (g*(at-offsets[role]))%7
   assert at==r["answers"][arm]
  assert r["answers"]["private_operator_oracle_unfair"]==r["truth"]
  for arm in controls:counts[arm]+=int(r["answers"][arm]==r["truth"])
 assert len(seen)==128
 for arm in controls:
  assert d["scores"][arm]["correct"]==counts[arm]
  assert math.isclose(d["scores"][arm]["accuracy"],counts[arm]/128,abs_tol=1e-12)
 for f in p["diagnostic_families"]:
  fam=[r for r in d["rows"] if r["family"]==f]
  assert len(fam)==32
  assert all(d["by_family"][str(f)][a]==sum(r["answers"][a]==r["truth"] for r in fam) for a in controls)
 candidate="learned_operator_posterior";fair=controls[1:-1]
 winner=max(fair,key=lambda a:counts[a])
 gain=(counts[candidate]-counts[winner])/128
 no_regression=all(d["by_family"][str(f)][candidate]>=max(d["by_family"][str(f)][a] for a in fair)
  for f in p["diagnostic_families"])
 qualified=gain>=p["min_gain_over_strongest_fair"] and no_regression
 assert d["strongest_fair_control"]==winner and math.isclose(d["candidate_gain"],gain,abs_tol=1e-12)
 assert d["per_family_nondegradation"]==no_regression and d["frozen_criterion_met"]==qualified
 print("W_F51_INDEPENDENT_RECOUNT=PASS",counts,"gain",gain,"qualified_train",qualified)
if __name__=="__main__":
 assert len(sys.argv)==3;run(sys.argv[1],sys.argv[2])
