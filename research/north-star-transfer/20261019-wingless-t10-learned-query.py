"""T10: TRAIN-learned public query utility, independent synthetic data."""
import importlib.util,json,random,time,tracemalloc,sys
from pathlib import Path
D=Path(__file__).parent
TRAIN=range(-22031,-21999);OPEN=range(-22163,-22099)
def load(path,name):
 s=importlib.util.spec_from_file_location(name,path)
 m=importlib.util.module_from_spec(s);sys.modules[name]=m;s.loader.exec_module(m);return m
def word(w):
 a=[]
 for t in w:
  if a and a[-1]==(t[0],not t[1]):a.pop()
  else:a.append(t)
 return tuple(a)
def source(seed,k):
 assert seed in TRAIN or seed in OPEN
 rng=random.Random(221017+abs(seed)*7229)
 dom=abs(seed)%2;n=11 if dom==0 else 3
 roles=[];truth=[]
 for _ in range(4):
  p=rng.choice(k.candidates(dom));x,y=rng.sample(range(n),2)
  roles.append(((x,k.operator_apply(dom,p,x)),(y,k.operator_apply(dom,p,y))))
  truth.append(p)
 w=[(0,False),(1,False)] if abs(seed)%4<2 else [(0,False),(1,True),(0,False),(2,False)]
 while len(w)<(10 if abs(seed)%4<2 else 12):
  j=rng.randrange(4);v=bool(rng.randrange(2));i=rng.randrange(len(w)+1);w[i:i]=[(j,v),(j,not v)]
 start=rng.randrange(n);state=k.start(dom)
 for j,v in w:state=k.compose_state(state,truth[j],dom,v)
 return (dom,tuple(roles),tuple(w),start),k.apply(state,start)
def possible(k,dom,rows):
 return [p for p in k.candidates(dom) if all(k.operator_apply(dom,p,x)==y for x,y in rows)]
def features(k,c,j):
 dom,roles,w,start=c
 active=word(w);f=sum(r==j for r,_ in active);n=len(possible(k,dom,roles[j][:1]))
 return (1,min(3,f),int(f>0),dom,min(5,n//2),int(len(active)>3))
def train(k,cases):
 weights=[0]*6;changes=0
 for _ in range(2):
  for c,answer in cases:
   dom,roles,w,start=c
   for j,rows in enumerate(roles):
    a=possible(k,dom,rows[:1])[0];b=possible(k,dom,rows)[0]
    # TRAIN-only causal intervention label using PUBLIC pairs:
    # replace only this role's first-pair estimate and determine whether
    # the final composed prediction changes, without hidden task answer.
    decoded=[possible(k,dom,role_rows)[0] for role_rows in roles]
    prior_ops=list(decoded);prior_ops[j]=a
    after_ops=list(decoded);after_ops[j]=b
    def run(ops):
     state=k.start(dom)
     for r,inv in word(w):state=k.compose_state(state,ops[r],dom,inv)
     return k.apply(state,start)
    y=1 if run(prior_ops)!=run(after_ops) else -1
    f=features(k,c,j)
    if y*sum(a*b for a,b in zip(weights,f))<=0:
     weights=[max(-127,min(127,a+y*b)) for a,b in zip(weights,f)]
     changes+=1
 return weights,changes
def predict(k,c,weights,mode,quota=None):
 dom,roles,w,start=c;counts={j:sum(r==j for r,v in word(w)) for j in range(4)}
 active={j for j in counts if counts[j]}
 if mode=="learn":selected={j for j in active if sum(a*b for a,b in zip(weights,features(k,c,j)))>0}
 elif mode=="full":selected=active
 elif mode=="fixed":selected=set(sorted(active,key=lambda j:(-counts[j],j))[:quota])
 else:selected=set()
 ops=[];service=0
 for j,rows in enumerate(roles):
  used=rows[:1]
  if j in selected:
   tick=time.process_time_ns();used=rows;service+=time.process_time_ns()-tick
  ops.append(possible(k,dom,used)[0])
 state=k.start(dom)
 for j,v in word(w):state=k.compose_state(state,ops[j],dom,v)
 return k.apply(state,start),len(selected),service
def measure(f):
 tracemalloc.start();t=time.process_time_ns();x=f()
 c={"cpu_ns":time.process_time_ns()-t,"peak_bytes":tracemalloc.get_traced_memory()[1]}
 tracemalloc.stop();return x,c
def main():
 frozen=json.loads((D/"20261019-wingless-t10-preregister.json").read_text())
 assert frozen["parent_head"]=="462ca47d62df18284db6a8fd7f07c6ff6cce2147" and not frozen["authority"]["protected_original_first_use"]
 parent=load(D/"20261019-wingless-t9-causal-query.py","w_t10_parent")
 k=parent.load()
 traincases=[source(i,k) for i in TRAIN];cases=[source(i,k) for i in OPEN]
 (weights,updates),tc=measure(lambda:train(k,traincases))
 assert updates>0 and len(weights)==6
 learned,lc=measure(lambda:[predict(k,c,weights,"learn") for c,v in cases])
 full,fc=measure(lambda:[predict(k,c,weights,"full") for c,v in cases])
 fixed,bc=measure(lambda:[predict(k,c,weights,"fixed",learned[i][1]) for i,(c,v) in enumerate(cases)])
 one,oc=measure(lambda:[predict(k,c,weights,"none") for c,v in cases])
 correct=lambda x:sum(a[0]==v for a,(c,v) in zip(x,cases))
 regression=lambda base:sum(a[0]==v and b[0]!=v for a,b,(c,v) in zip(base,learned,cases))
 result={"schema":"wingless.t10.learned-public-query.v1","TRAIN":len(traincases),"OPEN":len(cases),
 "trained_signed_weights":weights,"train_updates":updates,
 "learned_correct":correct(learned),"same_budget_fixed_correct":correct(fixed),
 "full_evidence_correct":correct(full),"one_pair_correct":correct(one),
 "regressions_vs_full":regression(full),"regressions_vs_one":regression(one),
 "learned_queries":sum(a[1] for a in learned),"equal_budget_control_queries":sum(a[1] for a in fixed),
 "full_queries":sum(a[1] for a in full),
 "learned_public_query_CPU_ns":sum(a[2] for a in learned),
 "training_cost":tc,"candidate_cost":lc,"fixed_cost":bc,"full_cost":fc,
 "model_bytes":6,"retained_state_bytes":3,
 "full_query_network_and_source_cost_proven":False,
 "scientific_acceptance":False,"protected_or_consumed_cases_used":0}
 assert result["learned_queries"]==result["equal_budget_control_queries"]
 print(json.dumps(result,sort_keys=True))
if __name__=="__main__":main()
