"""W-F58: persistent map identity under long-role reuse, FIT-only source qualification.

Every public role has three distinct calibration inputs over F7 with
independent family-dependent observation corruption. Exact 14^4 joint-map
posterior retains a map consistently across all repeated occurrences.
The fixed factorized belief deliberately re-samples the same inferred map
distribution on every reuse, exposing task-state continuity headroom.
No learner is trained and no diagnostic/primary original is executed here.
"""
from __future__ import annotations
import itertools,json,math,random,sys
FIT=range(681000,681016);FAMS=(1300,1301,1302,1303)
DIAG=range(681100,681116);ORIGINAL=range(682000,682016)
P=7;MAPS=tuple((a,b) for a in (1,2) for b in range(P))
ARMS=("exact_joint_public_posterior_extra_work","factorized_public_belief_fixed",
      "frozen_top_two_role_map_beam","public_role_MAP","private_truth_oracle_unfair")
def noise(fam):
 if fam not in FAMS:raise ValueError("FIT only")
 return (.03,.07,.10,.14)[fam-1300]
def source(seed,fam,length):
 if seed not in FIT or fam not in FAMS or length not in (8,10):
  raise ValueError("no diagnostic, original or source identity drift")
 rng=random.Random(seed*27991+fam*5107+length*1607+58031)
 maps=tuple(rng.choice(MAPS) for _ in range(4));rate=noise(fam)
 public=[];corrupt=0
 for role in range(4):
  xs=rng.sample(range(P),3);pairs=[]
  for x in xs:
   a,b=maps[role];y=(a*x+b)%P
   if rng.random()<rate:
    y=(y+rng.randrange(1,P))%P;corrupt+=1
   pairs.append((x,y))
  public.append((role,tuple(pairs)))
 # Every role occurs at least twice; inverse and forward directions
 # have independent FIT noise and long composition order.
 chain=[(role,1 if ((i+fam)%3) else -1)
        for i,role in enumerate([0,1,2,3,3,2,1,0])]
 if length==10:chain.extend(((rng.randrange(4),-1),(rng.randrange(4),1)))
 rng.shuffle(chain)
 start=rng.randrange(P);at=start
 for role,d in chain:
  a,b=maps[role]
  at=(a*at+b)%P if d==1 else ((at-b)*pow(a,-1,P))%P
 return dict(seed=seed,family=fam,length=length,public=tuple(public),
   chain=tuple(chain),start=start,truth=at,private_maps=maps,noisy_pairs=corrupt)
def posterior(case):
 rate=noise(case["family"]);res=[]
 for role,pairs in case["public"]:
  weights=[]
  for a,b in MAPS:
   value=1.
   for x,y in pairs:
    value*=1-rate if (a*x+b)%P==y else rate/(P-1)
   weights.append(value)
  z=sum(weights);assert z>0
  res.append(tuple(w/z for w in weights))
 return tuple(res)
def compose(maps,case):
 value=case["start"]
 for role,d in case["chain"]:
  a,b=maps[role]
  value=(a*value+b)%P if d==1 else ((value-b)*pow(a,-1,P))%P
 return value
def exact(case,post):
 out=[0.]*P
 for ix in itertools.product(range(len(MAPS)),repeat=4):
  prob=post[0][ix[0]]*post[1][ix[1]]*post[2][ix[2]]*post[3][ix[3]]
  if prob<1e-18:continue
  at=compose(tuple(MAPS[k] for k in ix),case)
  out[at]+=prob
 return max(range(P),key=lambda j:(out[j],-j)),len(MAPS)**4*len(case["chain"])
def factorized(case,post):
 belief=[0.]*P;belief[case["start"]]=1.
 for role,d in case["chain"]:
  new=[0.]*P
  for at in range(P):
   if belief[at]<=0:continue
   for (a,b),prob in zip(MAPS,post[role]):
    goal=(a*at+b)%P if d==1 else ((at-b)*pow(a,-1,P))%P
    new[goal]+=belief[at]*prob
  belief=new
 return max(range(P),key=lambda j:(belief[j],-j)),len(case["chain"])*P*len(MAPS)
def beam(case,post):
 choices=[sorted(range(len(MAPS)),key=lambda k:(-post[r][k],k))[:2] for r in range(4)]
 out=[0.]*P
 for ix in itertools.product(*choices):
  prob=math.prod(post[role][ix[role]] for role in range(4))
  out[compose(tuple(MAPS[k] for k in ix),case)]+=prob
 return max(range(P),key=lambda j:(out[j],-j)),16*len(case["chain"])
def map_control(case,post):
 selected=tuple(MAPS[max(range(len(MAPS)),key=lambda k:(post[r][k],-k))]
  for r in range(4))
 return compose(selected,case),len(case["chain"])
def fit_report():
 rows=[];scores={a:0 for a in ARMS};work={a:0 for a in ARMS[:-1]}
 for seed in FIT:
  for fam in FAMS:
   for length in (8,10):
    c=source(seed,fam,length);post=posterior(c)
    observations=[exact(c,post),factorized(c,post),beam(c,post),map_control(c,post),
                  (c["truth"],0)]
    d={arm:int(v==c["truth"]) for arm,(v,ops) in zip(ARMS,observations)}
    for arm,value in d.items():scores[arm]+=value
    cost={arm:observations[i][1] for i,arm in enumerate(ARMS[:-1])}
    for arm,value in cost.items():work[arm]+=value
    rows.append(dict(seed=seed,family=fam,length=length,public_noisy_pairs=c["noisy_pairs"],
       correct=d,abstract_operations=cost))
 assert len(rows)==128
 advantage=(scores[ARMS[0]]-scores[ARMS[1]])/128
 eligible=scores[ARMS[0]]/128>=.8 and advantage>=.05
 return dict(schema="wingless.f58.persistent-role-FIT-source.v1",
  stage="F58_NEW_FIT_SOURCE_ONLY_NO_MODEL_OR_DIAGNOSTIC",
  fit_seeds=[FIT.start,FIT.stop],fit_families=list(FAMS),
  cases=128,scores=scores,abstract_operations=work,
  exact_advantage_vs_factorized=advantage,
  public_exact_accuracy=scores[ARMS[0]]/128,
  frozen_source_gate_met=eligible,rows=rows,
  DIAG_untouched=True,originals_untouched=True,independent_F3=False,
  native_total_576B_unverified=True,accepted_refs_unchanged=True)
if __name__=="__main__":
 if sys.argv[1:]!=["--fit-only"]:raise SystemExit("FIT-only public feasibility")
 print(json.dumps(fit_report(),sort_keys=True))
