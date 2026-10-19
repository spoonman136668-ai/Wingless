"""T12 Wingless: exact PUBLIC uncertainty certificate + learned acquisition order.

Source-neutral ONLY. Enumerate compatible role operators using already
disclosed public pairs and their 10/12-hop final outcomes. If all compatible
assignments yield the same final output, inference needs no further public
query. Otherwise use TRAIN-only six-int8 causal-priority weights to choose
which unresolved role receives a second PUBLIC observation. No oracle labels
or true private operators at prediction time.
"""
from __future__ import annotations
import types,itertools,json,time,tracemalloc
from pathlib import Path
D=Path(__file__).parent
TRAIN=range(-30031,-29999);OPEN=range(-30163,-30099)
def measured(f):
 tracemalloc.start();t=time.process_time_ns();r=f()
 cost={"cpu_ns":time.process_time_ns()-t,"Python_peak_bytes":tracemalloc.get_traced_memory()[1]}
 tracemalloc.stop();return r,cost
def audit(t10,t11,k,c,w):
 dom,roles,word,start=c
 active=tuple(sorted({j for j,inv in t10.word(word)}))
 rows=[list(roles[j][:1]) for j in range(4)]
 guesses=t10.possible
 query_cpu=0;extra_queries=0;oracle_free=True;enumerations=0
 def outputs():
  nonlocal enumerations
  choices=[guesses(k,dom,r) if j in active else guesses(k,dom,r)[:1]
           for j,r in enumerate(rows)]
  all_answers=set()
  for option in itertools.product(*choices):
   enumerations+=1
   all_answers.add(t11.final(t10,k,c,option))
  return all_answers
 while True:
  out=outputs()
  if len(out)==1:
   return next(iter(out)),{"extra_public_queries":extra_queries,
       "public_service_CPU_ns":query_cpu,
       "sourceblind_full_composition_enumerations":enumerations,
       "proof_unique_public_final_output":True}
  remaining=[j for j in active if len(rows[j])==1]
  if not remaining:
   raise ValueError("WING_T12_TWO_PUBLIC_ROWS_STILL_AMBIGUOUS")
  # TRAIN-only score ranks which public uncertainty to resolve first.
  j=max(remaining,key=lambda j:(sum(a*b for a,b in zip(w,t11.encode(t10,k,c,j))),-j))
  ts=time.process_time_ns()
  newly_public=roles[j][1]
  assert newly_public[0]!=rows[j][0][0]
  rows[j].append(newly_public)
  query_cpu+=time.process_time_ns()-ts
  extra_queries+=1
def main():
 t11=__import__("importlib").util
 from importlib.util import spec_from_file_location,module_from_spec
 import sys
 def load(p,n):
  sp=spec_from_file_location(n,p);m=module_from_spec(sp);sys.modules[n]=m;sp.loader.exec_module(m);return m
 f=json.loads((D/"20261019-wingless-t12-freeze.json").read_text())
 assert f["schema"]=="northstar.t12.wingless.frozen-preimplementation.v1" and f["parent_exact_head"]=="2194103d1f9a404160d4932771365ba769e2b040"
 t11=load(D/"20261019-wingless-t11-mechanism.py","wing_t12_t11")
 t10=load(D/"20261019-wingless-t10-learned-query.py","wing_t12_t10")
 k=load(D/"20261018-wingless-t3-typed-composition.py","wing_t12_typed")
 source=types.FunctionType(t10.source.__code__,dict(t10.source.__globals__,TRAIN=TRAIN,OPEN=OPEN),t10.source.__name__,t10.source.__defaults__)
 train=[source(i,k) for i in TRAIN]
 OPEN_rows=[source(i,k) for i in OPEN]
 (w,fit),tcost=measured(lambda:t11.fit(t10,k,train))
 assert fit["parameter_updates"]>0 and len(w)==6
 candidate,cost=measured(lambda:[audit(t10,t11,k,c,w) for c,y in OPEN_rows])
 full,fullcost=measured(lambda:[t11.choose(t10,k,c,w,"full") for c,y in OPEN_rows])
 one,onecost=measured(lambda:[t11.choose(t10,k,c,w,"none") for c,y in OPEN_rows])
 advisory,advisorycost=measured(lambda:[t11.choose(t10,k,c,w,"learned") for c,y in OPEN_rows])
 exact=lambda a:sum(v[0]==y for v,(c,y) in zip(a,OPEN_rows))
 query_total=lambda a:sum(row[1]["extra_public_queries"] for row in a)
 actual_quota=sum(z[1]["extra_public_queries"] for z in candidate)
 baseline_positive_lost=sum(ba[0]==y and ca[0]!=y for ba,ca,(c,y) in zip(full,candidate,OPEN_rows))
 out={"schema":"wingless.t12.certified-public-final-output-acquisition.v1",
      "TRAIN":len(train),"OPEN":len(OPEN_rows),"train_parameter_updates":fit["parameter_updates"],
      "weights_int8":w,"learned_certificate_correct":exact(candidate),
      "full_second_pair_control_correct":exact(full),
      "untrained_one_pair_correct":exact(one),
      "uncertified_learned_T11_correct":exact(advisory),
      "lost_full_baseline_positives":baseline_positive_lost,
      "certified_public_second_pair_queries":actual_quota,
      "full_active_second_pair_queries":sum(r[1]["queries"] for r in full),
      "public_service_measured_CPU_ns":sum(c[1]["public_service_CPU_ns"] for c in candidate),
      "candidate_sourceblind_composition_evaluations":sum(c[1]["sourceblind_full_composition_enumerations"] for c in candidate),
      "training_CPU_peak":tcost,"candidate_inference_CPU_peak":cost,
      "full_control_inference_CPU_peak":fullcost,
      "train_amortized_total_CPU_ns":cost["cpu_ns"]+tcost["cpu_ns"]/len(OPEN_rows),
      "full_control_total_CPU_ns":fullcost["cpu_ns"],
      "full_CPU_nonincrease":cost["cpu_ns"]+tcost["cpu_ns"]/len(OPEN_rows)<=fullcost["cpu_ns"],
      "full_peak_nonincrease":max(cost["Python_peak_bytes"],tcost["Python_peak_bytes"])<=fullcost["Python_peak_bytes"],
      "remote_public_request_network_cost_accounted":False,
      "protected_original_source_used":0,"science_accepted":False}
 assert out["learned_certificate_correct"]==out["full_second_pair_control_correct"]==len(OPEN_rows)
 print(json.dumps(out,sort_keys=True))
if __name__=="__main__":main()
