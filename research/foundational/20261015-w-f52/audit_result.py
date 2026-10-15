"""Independent verification of W-F52 exact public information, state and cost.

This audit consumes the one archived result only after source execution.
It never calls the experiment's run() and never touches protected originals.
"""
import importlib.util,json,math,pathlib,sys
ROOT=pathlib.Path(__file__).resolve().parent
spec=importlib.util.spec_from_file_location("w52_audit_source",ROOT/"source_pilot.py")
m=importlib.util.module_from_spec(spec);sys.modules[spec.name]=m
spec.loader.exec_module(m)
def check(path,prereg):
 d=json.load(open(path));p=json.load(open(prereg))
 t=p["train"];v=p["diagnostic"];arms=p["controls"];g=p["thresholds"]
 assert d["schema"]=="wingless.f52.quantized-compositional-belief.trainonly.v1"
 assert d["stage"]=="DISJOINT_TRAIN_ONLY_NOVEL_VARIABLE_BINDING_NOT_ORIGINAL"
 assert d["fit_seeds"]==t["seeds"] and d["fit_families"]==t["families"]
 assert d["diagnostic_seeds"]==v["seeds"] and d["diagnostic_families"]==v["families"]
 assert set(t["families"]).isdisjoint(v["families"])
 assert d["cases"]==len(d["rows"])==v["expected_cases"]==128
 assert d["train_examples"]==t["expected_train_cases"]==128
 assert d["future_original_seeds_NOT_RESERVED"]==p["future_original_NOT_RESERVED"]["seeds"]
 assert len(d["trained_quantized_weights"])==7 and d["persistent_weight_bytes"]==7
 assert d["hidden_belief_state_bytes"]==7 and d["physical_total_576B_unverified"]
 assert not d["independent_F3"] and not d["original_executed"]
 assert not d["accepted_ref_mutation"] and d["paid_api_calls"]==0
 assert d["trained_quantized_weights"]==list(m.train()[0])
 seen=set();correct={a:0 for a in arms};work={a:0 for a in arms[:-1]}
 for row in d["rows"]:
  key=row["seed"],row["family"],row["length"]
  assert key not in seen;seen.add(key)
  assert 541100<=row["seed"]<541116 and row["family"] in v["families"]
  assert row["length"] in v["chain_lengths"]
  assert len(row["public_context"])==3
  assert len(row["public_observations"])==4
  assert len(row["public_calibration"])==4
  assert len(row["chain"])==row["length"]
  assert set(row["answers"])==set(arms)
  assert set(row["abstract_operation_counts"])==set(arms[:-1])
  assert set(row["belief_state_bytes"])==set(arms[:-1])
  identity=m.source(*key)
  assert list(identity["context"])==row["public_context"]
  assert [list(x) for x in identity["observations"]]==row["public_observations"]
  assert [[list(x[0]),list(x[1]),list(x[2])] for x in identity["calibration"]]==row["public_calibration"]
  assert [list(x) for x in identity["chain"]]==row["chain"]
  assert identity["start"]==row["start"]
  assert identity["truth"]==row["evaluator_truth"]
  assert row["answers"]["private_truth_oracle_unfair"]==identity["truth"]
  for arm in arms[:-1]:
   answer,cost,bytes_count=m.decide(arm,identity,d["trained_quantized_weights"])
   assert (answer,cost,bytes_count)==(row["answers"][arm],row["abstract_operation_counts"][arm],row["belief_state_bytes"][arm])
   if arm=="trained_quantized_state":assert bytes_count<=g["compressed_state_bytes_target"]
   correct[arm]+=int(answer==identity["truth"])
   work[arm]+=cost
  correct["private_truth_oracle_unfair"]+=1
 assert len(seen)==128
 for arm in arms:
  assert d["scores"][arm]["correct"]==correct[arm]
  assert d["scores"][arm]["total"]==128
  assert math.isclose(d["scores"][arm]["accuracy"],correct[arm]/128,abs_tol=1e-12)
 for arm in arms[:-1]:assert d["posterior_abstract_operations"][arm]==work[arm]
 for f in v["families"]:
  subset=[r for r in d["rows"] if r["family"]==f]
  assert len(subset)==32
  for arm in arms:
   assert d["by_family"][str(f)][arm]==sum(
       r["answers"][arm]==r["evaluator_truth"] for r in subset)
 candidate="trained_quantized_state"
 fair=("budget_matched_symbolic_approx","public_three_bit_rule","public_support_hamming")
 best=max(fair,key=lambda a:correct[a])
 gain=(correct[candidate]-correct[best])/128
 loss=(correct["full_exact_public_posterior"]-correct[candidate])/128
 ratio=work[candidate]/max(1,work["full_exact_public_posterior"])
 nonreg=all(d["by_family"][str(f)][candidate]>=max(d["by_family"][str(f)][a]
            for a in fair) for f in v["families"])
 passed=gain>=g["gain_vs_budget_matched_min"] and loss<=g["max_absolute_accuracy_loss_vs_exact"] and ratio<=g["max_measured_abstract_work_ratio_vs_exact"] and nonreg
 assert d["strongest_budget_matched_fair"]==best
 assert math.isclose(d["candidate_gain"],gain,abs_tol=1e-12)
 assert math.isclose(d["absolute_loss_vs_exact"],loss,abs_tol=1e-12)
 assert math.isclose(d["operation_ratio_vs_exact"],ratio,abs_tol=1e-12)
 assert d["family_nondegradation"]==nonreg and d["frozen_gate_met"]==passed
 assert d["verdict"]==("TRAIN_ONLY_QUANTIZED_COMPOSITION_SUPPORTED" if passed else "TRAIN_ONLY_QUANTIZED_COMPOSITION_NEGATIVE")
 print("W_F52",d["verdict"],"correct",correct,"gain",gain,"loss_vs_exact",loss,"op_ratio",ratio,"family_nondegrade",nonreg)
if __name__=="__main__":
 if len(sys.argv)!=3:raise SystemExit("expected RESULT_JSON PREREG_JSON")
 check(*sys.argv[1:])
