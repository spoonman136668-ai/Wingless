"""Read-only W-F53 scorer and cost audit of one source result (NO original replay)."""
import importlib.util,json,math,pathlib,sys
ROOT=pathlib.Path(__file__).resolve().parent
sp=importlib.util.spec_from_file_location("w53_pinned_source",ROOT/"source_pilot.py")
m=importlib.util.module_from_spec(sp);sys.modules[sp.name]=m;sp.loader.exec_module(m)
def audit(path,prereg):
 d=json.load(open(path));p=json.load(open(prereg))
 t=p["fit"];v=p["diagnostic"];arms=p["controls"];g=p["thresholds"]
 assert d["schema"]=="wingless.f53.future-value-beam.trainonly.v1"
 assert d["stage"]=="TRAIN_ONLY_LEARNED_BELIEF_PRUNING_NO_ORIGINAL"
 assert d["fit_seeds"]==t["seeds"] and d["fit_families"]==t["families"]
 assert d["diagnostic_seeds"]==v["seeds"] and d["diagnostic_families"]==v["families"]
 assert set(t["families"]).isdisjoint(v["families"])
 assert d["cases"]==len(d["rows"])==v["expected_cases"]==128
 assert d["proposed_original_NOT_RESERVED"]==p["future_original_NOT_RESERVED"]["seeds"]
 assert not d["independent_F3"] and not d["original_executed"]
 assert d["native_576B_total_unverified"] and not d["accepted_ref_mutation"]
 assert d["paid_api_calls"]==0 and d["train_examples"]>0
 assert d["model_parameter_bytes"]==len(d["quantized_fit_weights"])==8
 assert d["recurrent_belief_state_bytes"]==7 and d["selected_state_index_bytes"]==2
 assert d["model_parameter_bytes"]+d["recurrent_belief_state_bytes"]+d["selected_state_index_bytes"]<=p["resource_limits"]["learned_state_payload_max_bytes"]
 w,n=m.trained_model()
 assert list(w)==d["quantized_fit_weights"] and n==d["train_examples"]
 actual={a:0 for a in arms};work={a:0 for a in arms[:-1]};seen=set()
 for row in d["rows"]:
  key=row["seed"],row["family"],row["length"]
  assert key not in seen;seen.add(key)
  seed,fam,length=key
  assert seed in m.DIAG and fam in m.DIAG_FAM and length in v["path_lengths"]
  case=m.source(*key)
  assert row["public_context"]==list(case["context"])
  assert row["public_observations"]==[list(x) for x in case["observations"]]
  assert row["public_calibration"]==[[list(a[0]),list(a[1]),list(a[2])] for a in case["calibration"]]
  assert row["chain"]==[list(x) for x in case["chain"]]
  assert row["start"]==case["start"] and row["truth"]==case["truth"]
  assert set(row["answers"])==set(arms)
  assert set(row["costs"])==set(arms[:-1])
  assert set(row["model_state_bytes"])==set(arms[:-1])
  for arm in arms:
   guess,cost,size=m.run_arm(case,arm,w)
   assert row["answers"][arm]==guess
   actual[arm]+=int(guess==case["truth"])
   if arm!=arms[-1]:
    assert row["costs"][arm]==cost and row["model_state_bytes"][arm]==size
    assert cost>=0 and size<=64
    work[arm]+=cost
 assert len(seen)==128
 for arm in arms:
  assert d["scores"][arm]["correct"]==actual[arm]
  assert d["scores"][arm]["total"]==128
  assert math.isclose(d["scores"][arm]["accuracy"],actual[arm]/128,abs_tol=1e-12)
 assert actual[arms[-1]]==128
 for arm in arms[:-1]:assert d["posterior_abstract_operations"][arm]==work[arm]
 for fam in v["families"]:
  rows=[r for r in d["rows"] if r["family"]==fam]
  assert len(rows)==32
  for arm in arms:
   assert d["by_family"][str(fam)][arm]==sum(r["answers"][arm]==r["truth"] for r in rows)
 fair=arms[1:6];best=max(fair,key=lambda a:actual[a])
 gain=(actual[arms[0]]-actual[best])/128
 ratio=work[arms[0]]/work["full_exact_public_bayes"]
 guard=all(d["by_family"][str(f)][arms[0]]>=max(d["by_family"][str(f)][a] for a in fair) for f in v["families"])
 accepted=gain>=g["gain_vs_strongest_equal_work_fair_min"] and ratio<=g["posterior_operation_ratio_vs_exact_max"] and guard
 assert d["strongest_equal_work_fair"]==best
 assert math.isclose(d["candidate_gain"],gain,abs_tol=1e-12)
 assert math.isclose(d["posterior_operation_ratio_vs_exact"],ratio,abs_tol=1e-12)
 assert d["family_nondegradation"]==guard and d["frozen_gate_met"]==accepted
 assert d["verdict"]==("TRAIN_ONLY_LEARNED_TASK_VALUE_SUPPORTED" if accepted else "TRAIN_ONLY_LEARNED_TASK_VALUE_NEGATIVE")
 print("W_F53",d["verdict"],"correct",actual,"gain",gain,"work_ratio",ratio,"family_guard",guard)
if __name__=="__main__":
 if len(sys.argv)!=3:raise SystemExit("expected RESULT_JSON FROZEN_PREREG_JSON")
 audit(*sys.argv[1:])
