"""W-F58 independent FIT-only score, identity, posterior-work recount.

Runs no source generator or learner; rejects every identity outside the
precode-frozen FIT. This is a source feasibility audit, not an independent
original or a scientific learning-success claim.
"""
import json,sys
def audit(result,prereg):
 d=json.load(open(result));p=json.load(open(prereg))
 assert d["schema"]=="wingless.f58.persistent-role-FIT-source.v1"
 assert d["stage"]=="F58_NEW_FIT_SOURCE_ONLY_NO_MODEL_OR_DIAGNOSTIC"
 assert d["fit_seeds"]==p["fit"]["seed_range"]
 assert d["fit_families"]==p["fit"]["families"]
 assert len(d["rows"])==d["cases"]==p["fit"]["cases"]==128
 assert d["DIAG_untouched"] and d["originals_untouched"]
 assert d["native_total_576B_unverified"] and d["accepted_refs_unchanged"] and not d["independent_F3"]
 arms=p["FIT_source_methods"]
 assert set(d["scores"])==set(arms)
 assert set(d["abstract_operations"])==set(arms[:-1])
 total={a:0 for a in arms};cost={a:0 for a in arms[:-1]};seen=set()
 for row in d["rows"]:
  seed,fam,length=row["seed"],row["family"],row["length"]
  key=(seed,fam,length);assert key not in seen;seen.add(key)
  assert p["fit"]["seed_range"][0]<=seed<p["fit"]["seed_range"][1]
  assert fam in p["fit"]["families"] and length in p["fit"]["lengths"]
  assert 0<=row["public_noisy_pairs"]<=12
  assert set(row["correct"])==set(arms)
  assert set(row["abstract_operations"])==set(arms[:-1])
  assert row["correct"][arms[-1]]==1
  assert row["abstract_operations"][arms[0]]==14**4*length
  assert row["abstract_operations"][arms[1]]==length*7*14
  assert row["abstract_operations"][arms[2]]==16*length
  assert row["abstract_operations"][arms[3]]==length
  for arm in arms:
   assert row["correct"][arm] in (0,1)
   total[arm]+=row["correct"][arm]
   if arm!=arms[-1]:cost[arm]+=row["abstract_operations"][arm]
 assert len(seen)==128 and total==d["scores"] and cost==d["abstract_operations"]
 advantage=(total[arms[0]]-total[arms[1]])/128
 exact=total[arms[0]]/128
 passed=(exact>=p["FIT_source_gates"]["exact_public_accuracy_min"]
  and advantage>=p["FIT_source_gates"]["exact_advantage_vs_factorized_min_absolute"])
 assert abs(d["exact_advantage_vs_factorized"]-advantage)<1e-12
 assert abs(d["public_exact_accuracy"]-exact)<1e-12
 assert d["frozen_source_gate_met"]==passed
 print("W_F58_FIT_SOURCE",json.dumps(dict(correct=total,exact_accuracy=exact,headroom=advantage,qualified=passed),sort_keys=True))
if __name__=="__main__":
 assert len(sys.argv)==3
 audit(sys.argv[1],sys.argv[2])
