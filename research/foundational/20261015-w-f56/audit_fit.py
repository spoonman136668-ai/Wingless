"""W-F56 independent source-free aggregate audit for frozen FIT-only JSON.

This verifier never runs the scientific source or accesses a heldout cohort.
Scientific feasibility FAIL is a valid audited source result, not permission
to lower the preregistered thresholds after inspection.
"""
import json,sys
def audit(path,prereg):
 d=json.load(open(path));p=json.load(open(prereg))
 assert d["schema"]=="wingless.f56.noisy-public-exact-posterior-fitonly.v1"
 assert d["stage"]=="FIT_ONLY_SOURCE_FEASIBILITY_NOT_TRAINED_NOT_DIAGNOSTIC"
 assert d["fit_seeds"]==p["fit"]["seeds"]
 assert d["fit_families"]==p["fit"]["families"]
 assert d["fit_cases"]==len(d["rows"])==p["fit"]["expected_cases"]==128
 assert d["diagnostic_unconsumed"] and d["originals_unconsumed"]
 assert d["one_exact_posterior_per_case"] and d["native_576B_total_unverified"]
 assert d["exact_solver_min_gate"]==p["feasibility_first"]["FIT_only_exact_public_solver_min_accuracy"]
 assert d["budget_headroom_min_gate"]==p["feasibility_first"]["FIT_only_fair_budget_headroom_min_pp"]
 seen=set();arms=("full_exact_public_bayes","budget_matched_fixed_public_posterior_beam","fixed_known_operator_bank")
 count={a:0 for a in arms};noised=0
 for row in d["rows"]:
  key=row["seed"],row["family"],row["length"]
  assert key not in seen;seen.add(key)
  assert p["fit"]["seeds"][0]<=key[0]<p["fit"]["seeds"][1]
  assert key[1] in p["fit"]["families"] and key[2] in p["fit"]["chain_lengths"]
  assert set(row["correct"])==set(arms)
  assert set(row["operations"])==set(arms)
  assert row["operations"]["full_exact_public_bayes"]==14**4*key[2]
  assert row["operations"]["budget_matched_fixed_public_posterior_beam"]==16*key[2]
  assert row["operations"]["fixed_known_operator_bank"]==key[2]
  for arm in arms:
   assert row["correct"][arm] in (0,1)
   count[arm]+=row["correct"][arm]
  assert 0<=row["public_noisy_pairs"]<=8
  noised+=row["public_noisy_pairs"]>0
 assert len(seen)==128 and count==d["correct"] and noised==d["noisy_source_cases"]
 fair=arms[1:];best=max(fair,key=lambda a:count[a])
 gain=(count[arms[0]]-count[best])/128
 feasible=count[arms[0]]/128>=p["feasibility_first"]["FIT_only_exact_public_solver_min_accuracy"] and gain>=p["feasibility_first"]["FIT_only_fair_budget_headroom_min_pp"]
 assert d["strongest_fixed_control"]==best
 assert abs(gain-d["exact_minus_best_fixed_gap"])<1e-12
 assert feasible==d["source_feasibility_pass"]
 print("W_F56_FIT_ONLY_FEASIBILITY",json.dumps(dict(correct=count,headroom=gain,noise_cases=noised,source_feasible=feasible),sort_keys=True))
if __name__=="__main__":
 assert len(sys.argv)==3
 audit(*sys.argv[1:])
