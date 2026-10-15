"""W-F54: FIT-trained role/direction-conditioned relation transition kernel.

Unlike F52/F53 this fits state *transition concentrations* using intermediate
FIT ground truth. Diagnostic decisions consume only public context, noisy
observations, public paired sign calibrations and the current 7-state belief.
The comparator is an explicit two-sign full public posterior under identical
likelihood assumptions, with strong rule and equal-work beam controls.
Only synthetic TRAIN science: NOT source-independent F3 or native 576B.
"""
from __future__ import annotations
import importlib.util,itertools,json,math,random,sys
from pathlib import Path
sp=importlib.util.spec_from_file_location("w54_f53_frozen_control",
 Path(__file__).resolve().parent.parent/"20261015-w-f53"/"source_pilot.py")
prev=importlib.util.module_from_spec(sp);sys.modules[sp.name]=prev;sp.loader.exec_module(prev)
FIT=range(581000,581016);DIAG=range(581100,581116)
FIT_FAM=(1040,1041,1042,1043);DIAG_FAM=(1044,1045,1046,1047)
ORIGINAL=range(582000,582016)
ARMS=("trained_compositional_transition","F53_frozen_learned_two_beam",
"fixed_mass_two_beam","fixed_entropy_two_beam","budget_matched_symbolic_transition",
"full_public_parity","public_support_nearest","exact_full_public_bayes",
"private_truth_oracle_unfair")
N=7
CONTEXTS=tuple(itertools.product((0,1),repeat=3))
PATHS=prev.TEMPLATES
TEMPERATURES=(.4,.7,1.,1.4,2.,3.)
def law(fam,context):
 a,b,c=context
 return (a^b^c,int(a+b+c>=2),a^(b&c),int(not a or (b and c)),
         (a|b)^c,int(a+2*b+c>=3),int(a+b+c==1),b^(a&c))[(fam-1040)%8]
def source(seed,fam,length):
 if not ((seed in FIT and fam in FIT_FAM) or (seed in DIAG and fam in DIAG_FAM)):
  raise ValueError("wrong experiment source or protected original")
 if length not in ((6,8) if seed in FIT else (8,10)):raise ValueError("length outside frozen fold")
 rng=random.Random(seed*26189+fam*1451+length*911+54011)
 context=CONTEXTS[rng.randrange(8)]
 noise=.12 if seed in FIT else .23
 sign=-1 if law(fam,context)^(rng.random()<noise) else 1
 shifts=tuple(rng.randrange(N) for _ in range(4))
 obs=[]
 for role in range(4):
  x=rng.randrange(N);y=(sign*x+shifts[role])%N
  if rng.random()<noise:y=(y+rng.randrange(1,N))%N
  obs.append((role,x,y))
 anchors=[]
 for calibration_context in rng.sample(CONTEXTS,4):
  s=-1 if law(fam,calibration_context)^(rng.random()<noise) else 1
  b=rng.randrange(N);x=rng.randrange(N);z=(x+rng.randrange(1,N))%N
  anchors.append((calibration_context,(x,(s*x+b)%N),(z,(s*z+b)%N)))
 chain=PATHS[fam-1040][:length]
 at=rng.randrange(N);start=at;states=[]
 for r,d in chain:
  at=(sign*at+shifts[r])%N if d==1 else (sign*(at-shifts[r]))%N
  states.append(at)
 return dict(seed=seed,family=fam,length=length,context=context,noise=noise,
    sign=sign,shifts=shifts,observations=tuple(obs),calibration=tuple(anchors),
    chain=tuple(chain),start=start,truth=at,intermediates=tuple(states))

def calibrated_minus(case):
 scores=[]
 for item in case["calibration"]:
  c,(x,y),(z,t)=item
  slope=((t-y)*pow((z-x)%N,-1,N))%N
  assert slope in (1,6)
  w=math.exp(-sum(v!=w for v,w in zip(case["context"],c)))
  scores.append((w,slope==6))
 return sum(w for w,b in scores if b)/sum(w for w,b in scores)
def likelihood(case,sign):
 n=case["noise"];out={}
 for role,x,y in case["observations"]:
  out[role]=tuple((1-n if (sign*x+s)%N==y else n/6) for s in range(N))
 return out
def step(state,role,direction,sign,likelihoods,active,temperature=1.,quantized=False):
 raw=likelihoods[role]
 if temperature!=1.:
  tempered=[p**temperature for p in raw]
  total=sum(tempered)
  loc=[p/total for p in tempered]
 else:loc=raw
 nxt=[0.]*N;work=0
 for at in active:
  for shift in range(N):
   to=(sign*at+shift)%N if direction==1 else (sign*(at-shift))%N
   nxt[to]+=state[at]*loc[shift];work+=1
 if quantized:
  fixed=[min(255,max(0,round(v*255))) for v in nxt]
  if not any(fixed):fixed[max(range(N),key=lambda i:nxt[i])]=1
  total=sum(fixed);nxt=[v/total for v in fixed];work+=7
 else:
  total=sum(nxt)
  if total:nxt=[v/total for v in nxt]
 return nxt,work

def propagate(case,sign,mode,weights=None,temperature=None,intermediate=False):
 lik=likelihood(case,sign);state=[0.]*N;state[case["start"]]=1.
 work=56;track=[]
 for j,(role,direction) in enumerate(case["chain"]):
  if mode=="exact":active=range(N)
  elif j==0:active=(case["start"],)
  elif mode=="F53_frozen_learned_two_beam":
   active=prev.pick(state,j,case["chain"],"learned_task_value_two_beam",weights)
   work+=24
  elif mode=="fixed_entropy_two_beam":
   active=prev.pick(state,j,case["chain"],"fixed_entropy_two_beam",[0]*8)
   work+=24
  else:
   active=prev.pick(state,j,case["chain"],"fixed_mass_two_beam",[0]*8)
   work+=24
  temp=1. if temperature is None else temperature[role*2+(direction<0)]
  state,c=step(state,role,direction,sign,lik,active,temp,mode!="exact")
  work+=c
  track.append(tuple(state))
 return state,work,track

def learner_fit():
 # Explicit FIT-only task-state supervision, no score fitting on diagnosis.
 cases=[source(seed,fam,length) for seed in FIT for fam in FIT_FAM for length in (6,8)]
 # Train eight shared role-direction transition temperatures by maximizing
 # intermediate correct-state log probability, not source-sign labels.
 params=[]
 for role in range(4):
  for direction in (1,-1):
   best_temp=1.;best_score=-float("inf")
   for value in TEMPERATURES:
    score=0.
    for case in cases:
     # Training may use the FIT hidden sign as an evaluator for transition
     # likelihood. Runtime sign is inferred entirely from public data.
     weights=[1.]*8;weights[role*2+(direction<0)]=value
     states,work,path=propagate(case,case["sign"],"fixed_mass_two_beam",temperature=weights)
     for j,(r,d) in enumerate(case["chain"]):
      if r==role and d==direction:
       score+=math.log(max(path[j][case["intermediates"][j]],1e-9))
    if score>best_score+1e-10 or (abs(score-best_score)<1e-10 and abs(value-1)<abs(best_temp-1)):
     best_score=score;best_temp=value
   params.append(best_temp)
 # Each temperature code is 1 byte, model weights 8B, 7B state and 2B
 # active indices. Runtime working set (Python objects) is not 17 bytes.
 quant=[TEMPERATURES.index(v) for v in params]
 return tuple(TEMPERATURES[i] for i in quant),tuple(quant),len(cases)

def public_sign(case,method):
 if method=="full_public_parity":return -1 if sum(case["context"])%2 else 1
 if method=="public_support_nearest":
  a=min(case["calibration"],key=lambda q:(sum(x!=y for x,y in zip(case["context"],q[0])),q[0]))
  (x,y),(z,t)=a[1:]
  return -1 if ((t-y)*pow((z-x)%7,-1,7))%7==6 else 1
 return -1 if calibrated_minus(case)>=.5 else 1

def infer(case,arm,params,f53_weights):
 if arm=="private_truth_oracle_unfair":return case["truth"],0,0
 if arm=="exact_full_public_bayes":
  p=calibrated_minus(case)
  a,ca,_=propagate(case,1,"exact")
  b,cb,_=propagate(case,-1,"exact")
  combined=[(1-p)*x+p*y for x,y in zip(a,b)]
  return max(range(N),key=lambda i:(combined[i],-i)),ca+cb+36,14
 sign=public_sign(case,arm)
 mode=("F53_frozen_learned_two_beam" if arm=="F53_frozen_learned_two_beam" else
       "fixed_entropy_two_beam" if arm=="fixed_entropy_two_beam" else
       "fixed_mass_two_beam")
 temp=params if arm=="trained_compositional_transition" else None
 state,work,_=propagate(case,sign,mode,f53_weights,temp)
 # Every control is charged matched inference and selection overhead.
 return max(range(N),key=lambda i:(state[i],-i)),work,7

def run():
 assert not set(FIT)&set(DIAG) and not set(FIT_FAM)&set(DIAG_FAM)
 params,encoded,nfit=learner_fit()
 # Prior model is a frozen historical policy fitted strictly on its own
 # original F53 FIT identities, never on new F54 diagnostic data.
 f53_weights,_=prev.trained_model()
 rows=[];costs={a:0 for a in ARMS[:-1]}
 for seed in DIAG:
  for fam in DIAG_FAM:
   for length in (8,10):
    d=source(seed,fam,length);answers={};ops={};bytes_used={}
    for arm in ARMS:
     answer,work,payload=infer(d,arm,params,f53_weights)
     answers[arm]=answer
     if arm!=ARMS[-1]:ops[arm]=work;bytes_used[arm]=payload;costs[arm]+=work
    rows.append(dict(seed=seed,family=fam,length=length,
        public_context=list(d["context"]),public_observations=[list(x) for x in d["observations"]],
        public_calibration=[[list(a[0]),list(a[1]),list(a[2])] for a in d["calibration"]],
        chain=[list(x) for x in d["chain"]],start=d["start"],
        answers=answers,truth=d["truth"],abstract_operations=ops,state_payload_bytes=bytes_used))
 counts={a:sum(row["answers"][a]==row["truth"] for row in rows) for a in ARMS}
 scores={a:dict(correct=counts[a],total=len(rows),accuracy=counts[a]/len(rows)) for a in ARMS}
 by_family={str(f):{a:sum(r["answers"][a]==r["truth"] for r in rows if r["family"]==f)
     for a in ARMS} for f in DIAG_FAM}
 fair=ARMS[1:7];strongest=max(fair,key=lambda a:counts[a])
 gain=(counts[ARMS[0]]-counts[strongest])/len(rows)
 work_ratio=costs[ARMS[0]]/max(1,costs["exact_full_public_bayes"])
 safe=all(by_family[str(f)][ARMS[0]]>=max(by_family[str(f)][a] for a in fair) for f in DIAG_FAM)
 passed=gain>=.05 and work_ratio<=.5 and safe
 return dict(schema="wingless.f54.learned-role-transition.trainonly.v1",
   stage="TRAIN_ONLY_DISJOINT_FIT_DIAG_LEARNED_TRANSITIONS",
   fit_seeds=[FIT.start,FIT.stop],fit_families=list(FIT_FAM),
   diagnostic_seeds=[DIAG.start,DIAG.stop],diagnostic_families=list(DIAG_FAM),
   original_seeds_NOT_RESERVED=[ORIGINAL.start,ORIGINAL.stop],
   train_examples=nfit,learned_transition_temperature_codes=list(encoded),
   F53_frozen_public_fitted_weights=list(f53_weights),transition_parameter_bytes=8,
   recurrent_belief_bytes=7,selected_state_index_bytes=2,
   physical_native_576B_unverified=True,cases=len(rows),rows=rows,
   scores=scores,by_family=by_family,posterior_abstract_operations=costs,
   strongest_fair_control=strongest,candidate_gain=gain,
   posterior_work_ratio_vs_exact=work_ratio,per_family_nondegradation=safe,
   frozen_gate_met=passed,
   verdict="TRAIN_ONLY_TRANSITION_GENERALIZATION_SUPPORTED" if passed else "TRAIN_ONLY_TRANSITION_GENERALIZATION_NEGATIVE",
   independent_F3=False,original_executed=False,accepted_ref_mutation=False,paid_api_calls=0)
if __name__=="__main__":print(json.dumps(run(),indent=2,sort_keys=True))
