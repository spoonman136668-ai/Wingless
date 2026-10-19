"""T13 Wingless: bounded safe public uncertainty + TRAIN learned acquisition rank.

The dynamic-programming output set overapproximates all shared-operator
assignments consistent with PUBLIC role samples, so singleton is a sound
proof of the final answer. We never use hidden true operators to decide.
Repeated roles are treated independently, which may cause extra public
queries but cannot create a false certificate.
"""
from __future__ import annotations
import types,json,importlib.util,sys,time,tracemalloc
from pathlib import Path
HERE=Path(__file__).resolve().parent
TRAIN=range(-34031,-33999)
OPEN=range(-34163,-34099)
def load(path,name):
 spec=importlib.util.spec_from_file_location(name,path)
 mod=importlib.util.module_from_spec(spec)
 sys.modules[spec.name]=mod;spec.loader.exec_module(mod);return mod
def measured(f):
 tracemalloc.start();before=time.process_time_ns()
 value=f();receipt={"cpu_ns":time.process_time_ns()-before,
  "python_peak_bytes":tracemalloc.get_traced_memory()[1]}
 tracemalloc.stop();return value,receipt
def image(k,domain,x,operator,inverse):
 if not inverse:return k.operator_apply(domain,operator,x)
 if domain==0:
  a,b=operator;return ((x-b)*pow(a,-1,11))%11
 return operator.index(x)
def conservative_possible_outputs(t10,k,case,public):
 dom,roles,word,start=case
 states={start};steps=0;maximum=1
 for role,inverse in t10.word(word):
  ops=t10.possible(k,dom,public[role])
  states={image(k,dom,x,op,inverse) for x in states for op in ops}
  steps+=len(ops)*maximum
  maximum=max(maximum,len(states))
 assert maximum<=(11 if dom==0 else 3)
 return states,steps,maximum
def acquire(t10,t11,k,case,weights):
 dom,roles,word,start=case
 active=sorted({j for j,inv in t10.word(word)})
 observed=[tuple(rows[:1]) for rows in roles]
 queries=0;ops_evaluated=0;maxfrontier=0;service_ns=0
 while True:
  possible,step,frontier=conservative_possible_outputs(t10,k,case,observed)
  ops_evaluated+=step;maxfrontier=max(maxfrontier,frontier)
  if len(possible)==1:
   return next(iter(possible)),{"queries":queries,
     "bounded_public_operator_evaluations":ops_evaluated,
     "max_public_reachable_output_states":maxfrontier,
     "public_query_service_cpu_ns":service_ns}
  choices=[j for j in active if len(observed[j])==1]
  if not choices:raise AssertionError("T13_CERTIFICATE_NOT_SINGLETON_DESPITE_COMPLETE_PAIR")
  role=max(choices,key=lambda j:(
      sum(a*b for a,b in zip(weights,t11.encode(t10,k,case,j))),-j))
  then=time.process_time_ns()
  second=roles[role][1]
  assert second[0]!=observed[role][0][0]
  observed[role]=observed[role]+(second,)
  service_ns+=time.process_time_ns()-then;queries+=1
def main():
 f=json.loads((HERE/"20261019-wingless-t13-freeze.json").read_text())
 assert f["schema"]=="northstar.t13.wingless.preimplementation-freeze.v1"
 assert f["science_head_parent"]=="adb17205c1f75b3a649646d5d63a7a2a6e68aea9"
 assert not f["safety"]["protected_science_first_use"]
 t10=load(HERE/"20261019-wingless-t10-learned-query.py","wing_t13_t10")
 t11=load(HERE/"20261019-wingless-t11-mechanism.py","wing_t13_t11")
 t12=load(HERE/"20261019-wingless-t12-mechanism.py","wing_t13_t12")
 k=load(HERE/"20261018-wingless-t3-typed-composition.py","wing_t13_kernel")
 source=types.FunctionType(t10.source.__code__,dict(t10.source.__globals__,TRAIN=TRAIN,OPEN=OPEN),t10.source.__name__,t10.source.__defaults__)
 train=[source(i,k) for i in TRAIN]
 open_cases=[source(i,k) for i in OPEN]
 (weights,training),traincost=measured(lambda:t11.fit(t10,k,train))
 assert len(weights)==6 and training["parameter_updates"]>0
 cand,ccost=measured(lambda:[acquire(t10,t11,k,c,weights) for c,y in open_cases])
 full,fcost=measured(lambda:[t11.choose(t10,k,c,weights,"full") for c,y in open_cases])
 prev,pcost=measured(lambda:[t12.audit(t10,t11,k,c,weights) for c,y in open_cases])
 score=lambda r:sum(v[0]==truth for v,(c,truth) in zip(r,open_cases))
 loss=lambda a,b:sum(x[0]==truth and y[0]!=truth
    for x,y,(c,truth) in zip(a,b,open_cases))
 out={"schema":"wingless.t13.sound-bounded-public-certificate.v1",
 "TRAIN":len(train),"OPEN":len(open_cases),"training_updates":training["parameter_updates"],
 "certificate_correct":score(cand),"strong_full_correct":score(full),
 "prior_exhaustive_cert_correct":score(prev),
 "baseline_correct_regressions":loss(full,cand),
 "candidate_public_queries":sum(v[1]["queries"] for v in cand),
 "full_public_queries":sum(v[1]["queries"] for v in full),
 "previous_exhaustive_queries":sum(v[1]["extra_public_queries"] for v in prev),
 "bounded_operator_transition_evaluations":sum(v[1]["bounded_public_operator_evaluations"] for v in cand),
 "max_reachable_states":max(v[1]["max_public_reachable_output_states"] for v in cand),
 "prior_exhaustive_assignments":sum(v[1]["sourceblind_full_composition_enumerations"] for v in prev),
 "candidate_cpu_peak":ccost,"prior_exhaustive_cpu_peak":pcost,
 "strong_full_cpu_peak":fcost,"train_cpu_peak":traincost,
 "training_amortized_candidate_cpu_ns":ccost["cpu_ns"]+traincost["cpu_ns"]/len(open_cases),
 "best_full_cpu_ns":fcost["cpu_ns"],
 "training_inclusive_peak_bytes":max(ccost["python_peak_bytes"],traincost["python_peak_bytes"]),
 "strong_full_peak_bytes":fcost["python_peak_bytes"],
 "remote_public_service_cost_measured":False,
 "protected_original_or_consumed_source_replayed":False,
 "learned_novel_cognition_accepted":False,"science_accepted":False}
 assert out["certificate_correct"]==len(OPEN) and out["baseline_correct_regressions"]==0
 print(json.dumps(out,sort_keys=True))
if __name__=="__main__":main()
