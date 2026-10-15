"""Audit the existing W-F48 artifact WITHOUT re-running the diagnostic.

Source used only to reconstruct each already-reported row and receipts.
No new model selection, paid reads, or primary originals.
"""
import importlib.util,json,math,pathlib,sys
root=pathlib.Path(__file__).resolve().parent
spec=importlib.util.spec_from_file_location("w48_existing_provenance",root/"source_pilot.py")
m=importlib.util.module_from_spec(spec);sys.modules[spec.name]=m;spec.loader.exec_module(m)
def audit(path):
 d=json.loads(pathlib.Path(path).read_text())
 assert d["schema"]=="wingless.f48.disjoint-trained-query-value.trainonly.v1"
 assert d["stage"]=="ONE_FRESH_FIT_AND_DIAGNOSTIC_NO_PRIMARY_ORIGINAL"
 assert d["questions"]==32 and d["training_questions"]==32
 assert d["fit_seeds"]==[411000,411004] and d["diagnostic_seeds"]==[411100,411104]
 assert d["families"]==list(m.FAMILIES)
 assert d["originals_reserved_not_consumed"]==[412000,412016]
 assert d["original_families_reserved_not_consumed"]==list(m.ORIGINAL_FAMILIES)
 assert d["diagnostic_truth_not_used_for_fit"] and d["public_only_action_inputs"]
 assert not d["original_executed"] and not d["independent_F3_custody"]
 assert not d["native_576_byte_resident_qualified"] and not d["accepted_ref_mutation"]
 assert d["paid_api_calls"]==0
 assert len(d["rows"])==32 and set(d["scores"])==set(m.ARMS)
 assert set(d["posterior_world_work"])==set(m.ARMS[:-1])
 assert d["selected_beta_train_only"] in m.BETAS
 assert set(d["training_beta_scores"])=={str(b) for b in m.BETAS}
 assert all(isinstance(n,int) and 0<=n<=32 for n in d["training_beta_scores"].values())
 assert d["selected_beta_train_only"]==min(m.BETAS,
     key=lambda b:(-d["training_beta_scores"][str(b)],abs(b),b))
 counts={a:0 for a in m.ARMS}
 costs={a:0 for a in m.ARMS[:-1]}
 keys=set()
 for r in d["rows"]:
   key=r["seed"],r["family"],r["length"]
   assert key not in keys;keys.add(key)
   seed,fam,length=key
   assert seed in m.DIAG and fam in m.FAMILIES and length in (3,5)
   private,public,start,chain=m.source(seed,fam,length)
   truth=m.old.f45.prior.prior.follow(private,start,chain)
   assert truth==r["truth"]
   assert set(r["answers"])==set(m.ARMS)
   assert set(r["receipts"])==set(m.ARMS[:-1])
   assert set(r["work"])==set(m.ARMS[:-1])
   assert r["answers"]["evaluator_oracle_unfair"]==truth
   already={(role,x) for role,x,y in public}
   for arm in m.ARMS[:-1]:
     receipts=r["receipts"][arm]
     assert len(receipts)<=2
     assert not receipts if arm=="no_purchase" else True
     observed=list(public)
     acquired=set()
     for pair in receipts:
       assert len(pair)==3 and all(isinstance(x,int) for x in pair)
       role,x,y=pair
       assert 0<=role<4 and 0<=x<6 and y==private[role][x]
       assert (role,x) not in already and (role,x) not in acquired
       acquired.add((role,x));observed.append((role,x,y))
     assert r["answers"][arm]==m.old.f45.prior.prior.predict(observed,start,chain)
     assert isinstance(r["work"][arm],int) and r["work"][arm]>=0
     costs[arm]+=r["work"][arm]
   for arm in m.ARMS:counts[arm]+=int(r["answers"][arm]==truth)
 assert len(keys)==32 and counts["evaluator_oracle_unfair"]==32
 for arm in m.ARMS:
   s=d["scores"][arm]
   assert (s["correct"],s["total"])==(counts[arm],32)
   assert math.isclose(s["accuracy"],counts[arm]/32,abs_tol=1e-12)
 for arm in m.ARMS[:-1]:assert costs[arm]==d["posterior_world_work"][arm]
 for fam in m.FAMILIES:
   sub=[r for r in d["rows"] if r["family"]==fam]
   assert len(sub)==8
   for arm in m.ARMS:
     assert d["by_family"][str(fam)][arm]==sum(r["answers"][arm]==r["truth"] for r in sub)
 fair=("one_step_value","exact_entropy","query_frontier","budget_matched_rank")
 best=max(fair,key=lambda a:(counts[a],-costs[a]))
 gain=(counts["learned_query_cost_value"]-counts[best])/32
 ratio=costs["learned_query_cost_value"]/max(1,costs[best])
 nonreg=all(d["by_family"][str(f)]["learned_query_cost_value"]>=max(
   d["by_family"][str(f)][a] for a in fair) for f in m.FAMILIES)
 passed=gain>=.05 and ratio<=1.2 and nonreg
 assert d["strongest_budget_fair_control"]==best
 assert math.isclose(d["candidate_gain"],gain,abs_tol=1e-12)
 assert math.isclose(d["posterior_work_ratio"],ratio,abs_tol=1e-12)
 assert d["family_nondegradation"]==nonreg and d["frozen_criterion_met"]==passed
 assert d["verdict"]==("TRAIN_ONLY_SUPPORTED_NOT_F3" if passed else "TRAIN_ONLY_NEGATIVE")
 out=dict(schema="wingless.f48.independent-existing-train-audit.v1",source_run=38805703700,
   cases=32,train_only=True,selected_beta=d["selected_beta_train_only"],scores_correct=counts,
   strongest_fair=best,gain=gain,work_ratio=ratio,per_family_nonregression=nonreg,
   frozen_gate_pass=passed,primary_original=False,independent_F3=False)
 print(json.dumps(out,indent=2,sort_keys=True))
 return out
if __name__=="__main__":
 if len(sys.argv)!=2:raise SystemExit("usage: audit_existing.py RESULT_JSON")
 audit(sys.argv[1])
