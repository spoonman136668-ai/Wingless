"""Source-free W-F57 FIT result/threshold audit."""
import json,sys
def verify(result,pre):
 d=json.load(open(result));p=json.load(open(pre))
 assert d["schema"]=="wingless.f57.public-coupled-fitonly.v1"
 assert d["stage"]=="FIT_ONLY_NO_DIAG"
 assert d["FIT_seeds"]==p["fit"]["seeds"]
 assert d["FIT_families"]==p["fit"]["families"]
 assert d["cases"]==len(d["rows"])==p["fit"]["expected_cases"]==128
 assert d["hypotheses"]==1372
 assert d["diagnostic_unused"] and d["originals_unused"]
 assert d["native_576B_total_unverified"] and not d["source_independent_F3"]
 arms=("full_exact_public_coupled_posterior","fixed_equal_work_cross_role_beam",
 "fixed_structured_symbolic_pruning","public_algebraic_consistency")
 scores={arm:0 for arm in arms};work={arm:0 for arm in arms};seen=set()
 for row in d["rows"]:
  k=(row["seed"],row["family"],row["length"])
  assert k not in seen;seen.add(k)
  assert p["fit"]["seeds"][0]<=k[0]<p["fit"]["seeds"][1]
  assert k[1] in p["fit"]["families"] and k[2] in p["fit"]["chain_lengths"]
  assert set(row["correctness"])==set(arms)
  assert set(row["operation_counts"])==set(arms)
  assert row["operation_counts"][arms[0]]==1372*(8+k[2])
  assert row["operation_counts"][arms[1]]==112+16*k[2]
  assert row["operation_counts"][arms[2]]==1372*8+k[2]
  assert 0<=row["noisy_pairs"]<=8
  for a in arms:
   assert row["correctness"][a] in (0,1)
   scores[a]+=row["correctness"][a];work[a]+=row["operation_counts"][a]
 assert len(seen)==128 and scores==d["correct"] and work==d["work"]
 best=max(arms[1:],key=lambda a:scores[a])
 gain=(scores[arms[0]]-scores[best])/128
 passed=(scores[arms[0]]/128>=p["FIT_feasibility_before_learning"]["exact_public_accuracy_min"]
     and gain>=p["FIT_feasibility_before_learning"]["gain_vs_strongest_equal_work_fair_min"])
 assert d["strongest_fair"]==best and abs(d["exact_headroom"]-gain)<1e-12
 assert d["source_feasible"]==passed
 print("W_F57_FIT_FEASIBILITY",json.dumps(dict(correct=scores,exact_headroom=gain,qualified=passed),sort_keys=True))
if __name__=="__main__":
 assert len(sys.argv)==3
 verify(sys.argv[1],sys.argv[2])
