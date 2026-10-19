"""T11 Wingless: train *final-answer causal sensitivity*, not just operator change.

TRAIN labels are differences in composed PUBLIC predictions after replacing
a role's one-sample estimate with its second-public-sample estimate. The
OPEN acquisition policy does not see the second observation until requested.
Synthetic source namespace entirely new; registered/consumed cohorts denied.
"""
from __future__ import annotations
import types,json,sys,importlib.util,time,tracemalloc
from pathlib import Path
D=Path(__file__).parent
TRAIN=range(-26031,-25999);OPEN=range(-26163,-26099)
def load(p,n):
 s=importlib.util.spec_from_file_location(n,p);m=importlib.util.module_from_spec(s)
 sys.modules[n]=m;s.loader.exec_module(m);return m
def measured(f):
 tracemalloc.start();t=time.process_time_ns();x=f()
 rec={"cpu_ns":time.process_time_ns()-t,"Python_peak_bytes":tracemalloc.get_traced_memory()[1]}
 tracemalloc.stop();return x,rec
def encode(t10,k,c,role):
 d,roles,w,start=c
 reduce=t10.word(w);count=sum(j==role for j,v in reduce)
 ambiguous=len(t10.possible(k,d,roles[role][:1]))
 first=t10.possible(k,d,roles[role][:1])[0]
 possibilities=t10.possible(k,d,roles[role][:1])
 # Potential causal output variation with other roles held at
 # PUBLIC one-pair estimators; does NOT consult a hidden true operator.
 other=[t10.possible(k,d,rows[:1])[0] for rows in roles]
 answers=set()
 for alt in possibilities:
  ops=list(other);ops[role]=alt
  state=k.start(d)
  for j,inv in reduce:state=k.compose_state(state,ops[j],d,inv)
  answers.add(k.apply(state,start))
 return (1,min(count,3),int(ambiguous>1),min(len(answers),3),
         int(d==1),int(len(reduce)>=4))
def final(t10,k,c,ops):
 d,roles,w,start=c;s=k.start(d)
 for j,inv in t10.word(w):s=k.compose_state(s,ops[j],d,inv)
 return k.apply(s,start)
def fit(t10,k,cases):
 weights=[0]*6;updates=0;public_interventions=0
 for _ in range(2):
  for c,_ in cases:
   d,roles,w,start=c
   one=[t10.possible(k,d,rows[:1])[0] for rows in roles]
   initial=final(t10,k,c,one)
   for j,rows in enumerate(roles):
    twice=list(one);twice[j]=t10.possible(k,d,rows)[0]
    label=1 if final(t10,k,c,twice)!=initial else -1
    public_interventions+=1
    f=encode(t10,k,c,j);score=sum(a*b for a,b in zip(weights,f))
    if label*score<=0:
     weights=[max(-127,min(127,a+label*b)) for a,b in zip(weights,f)]
     updates+=1
 return tuple(weights),{"parameter_updates":updates,
     "TRAIN_public_causal_interventions":public_interventions,
     "oracle_labels":False,"retained_weight_bytes":6}
def choose(t10,k,c,weights,mode,budget=None):
 d,roles,w,start=c
 word=t10.word(w);active={j for j,v in word}
 if mode=="learned":selected={j for j in active if sum(a*b for a,b in zip(weights,encode(t10,k,c,j)))>0}
 elif mode=="full":selected=active
 elif mode=="none":selected=set()
 elif mode=="matched":
  counts={j:sum(r==j for r,v in word) for j in active}
  selected=set(sorted(active,key=lambda j:(-counts[j],j))[:budget])
 else:raise AssertionError("MODE")
 decoded=[];querycpu=0
 for j,rows in enumerate(roles):
  visible=rows[:1]
  if j in selected:
   before=time.process_time_ns();visible=rows
   querycpu+=time.process_time_ns()-before
  decoded.append(t10.possible(k,d,visible)[0])
 return final(t10,k,c,decoded),{"queries":len(selected),"service_cpu_ns":querycpu}
def main():
 frozen=json.loads((D/"20261019-wingless-t11-freeze.json").read_text())
 assert frozen["schema"]=="northstar.t11.wingless.preregister.v1" and frozen["parent_exact_head"]=="5c82d69a28717f2fa4b609f65db1ee0abd3c8d65"
 t10=load(D/"20261019-wingless-t10-learned-query.py","wing_t11_original_t10")
 k=load(D/"20261018-wingless-t3-typed-composition.py","wing_t11_typed")
 globs=dict(t10.source.__globals__,TRAIN=TRAIN,OPEN=OPEN)
 source=types.FunctionType(t10.source.__code__,globs,t10.source.__name__,t10.source.__defaults__)
 training=[source(i,k) for i in TRAIN];tests=[source(i,k) for i in OPEN]
 (weights,receipt),traincost=measured(lambda:fit(t10,k,training))
 assert len(weights)==6 and receipt["parameter_updates"]>0
 learned,learncost=measured(lambda:[choose(t10,k,c,weights,"learned") for c,y in tests])
 full,fullcost=measured(lambda:[choose(t10,k,c,weights,"full") for c,y in tests])
 none,nonecost=measured(lambda:[choose(t10,k,c,weights,"none") for c,y in tests])
 matched,matchedcost=measured(lambda:[
   choose(t10,k,c,weights,"matched",learned[i][1]["queries"]) for i,(c,y) in enumerate(tests)])
 def correct(a):return sum(x[0]==v for x,(c,v) in zip(a,tests))
 def lost(base):return sum(a[0]==v and b[0]!=v for a,b,(c,v) in zip(base,learned,tests))
 totalqueries=lambda a:sum(x[1]["queries"] for x in a)
 out={"schema":"wingless.t11.learned-final-causal-query.v1",
  "TRAIN":len(training),"OPEN":len(tests),
  "trained_int8_weights":weights,"training_updates":receipt["parameter_updates"],
  "TRAIN_public_intervention_examples":receipt["TRAIN_public_causal_interventions"],
  "learned_exact":correct(learned),"matched_queries_exact":correct(matched),
  "full_public_evidence_exact":correct(full),"one_pair_exact":correct(none),
  "regressed_full_correct":lost(full),"regressed_one_pair_correct":lost(none),
  "learned_queries":totalqueries(learned),"matched_queries":totalqueries(matched),
  "full_queries":totalqueries(full),
  "candidate_training_cost":traincost,"candidate_inference_cost":learncost,
  "strong_full_inference_cost":fullcost,"matched_inference_cost":matchedcost,
  "measured_extra_public_query_CPU_ns":sum(x[1]["service_cpu_ns"] for x in learned),
  "full_remote_service_cost_measured":False,
  "full_resource_parity_qualified":False,"novel_structure_replication":False,
  "original_protected_or_consumed_cohorts_used":0,"science_accepted":False}
 assert out["learned_queries"]==out["matched_queries"]
 print(json.dumps(out,sort_keys=True))
if __name__=="__main__":main()
