"""W-F61 independent source-free D4 counterfactual triple-binding audit.

Reconstructs all possible public D4 relations from archived public pairs and
paired query descriptions; never imports science source or private maps.
The published evaluator truth is treated as a claim (no independent physical
custody); no primary-original conclusion can follow from this audit.
"""
import itertools,json,math,sys
C=((0,0),(1,0),(1,1),(0,1))
def makeperm(k):
 flip=k//4;rot=k%4;ret=[]
 for x,y in C:
  if flip:x=1-x
  for _ in range(rot):x,y=1-y,x
  ret.append(C.index((x,y)))
 return tuple(ret)
PERMS=tuple(makeperm(k) for k in range(8))
def compose(a,b):return tuple(a[b[i]] for i in range(4))
def invert(p):return tuple(p.index(i) for i in range(4))
MAPS=tuple((a,b,c,invert(compose(c,compose(b,a))))
 for a,b,c in itertools.product(PERMS,repeat=3))
ARMS=("public_full_exact_joint_posterior_extra_work","strong_structured_MAP",
      "budget_matched_public_dp","frozen_F60_public_coupled_MAP",
      "public_local_rolewise","private_oracle_unfair")
def recompute(row):
 rate=(.045,.085,.12,.16)[row["family"]-1440]
 observations=row["public"];bridge=row["public_bridges"]
 chain=row["chain"];second=row["secondary_chain"];counter=row["counterfactual_chain"]
 def likelihood(maps):
  weight=1.
  for r,cal in enumerate(observations):
   for x,y in cal:
    weight*=1-rate if maps[r][x]==y else rate/3
  for a,b,x,y in bridge:
   weight*=1-rate if maps[b][maps[a][x]]==y else rate/3
  return weight
 def run(m,start,path):
  value=start
  for role,d in path:
   p=m[role];value=p[value] if d==1 else invert(p)[value]
  return value
 def triple(m):
  a=run(m,row["start"],chain)
  return (a,run(m,a,second),run(m,row["counterfactual_start"],counter))
 # Normalize the same PUBLIC joint likelihood distribution before
 # summation. With D4 symmetries, nearly tied pair posteriors can flip the
 # floating-point tie choice if one verifier sums unnormalized weights.
 # This is a numeric identity correction, not source/data/threshold tuning.
 raw=[likelihood(m) for m in MAPS]
 total=sum(raw);assert total>0
 weights=[v/total for v in raw]
 joint={}
 for m,v in zip(MAPS,weights):
  k=triple(m);joint[k]=joint.get(k,0.)+v
 exact=max(sorted(joint),key=lambda p:joint[p])
 best=max(range(512),key=lambda k:(weights[k],-k))
 structured=triple(MAPS[best])
 ranks=[]
 for r,cal in enumerate(observations):
  likelihoods=[math.prod(1-rate if p[x]==y else rate/3 for x,y in cal)
     for p in PERMS]
  ranks.append(sorted(range(8),key=lambda k:(-likelihoods[k],k)))
 independent=triple(tuple(PERMS[rr[0]] for rr in ranks))
 beam={}
 for indices in itertools.product(*[r[:2] for r in ranks[:3]]):
  a,b,c=(PERMS[k] for k in indices)
  maps=(a,b,c,invert(compose(c,compose(b,a))))
  pair=triple(maps);beam[pair]=beam.get(pair,0.)+likelihood(maps)
 selected=max(sorted(beam),key=lambda p:beam[p])
 return (exact,structured,selected,structured,independent),len(joint)
def audit(result,prereg):
 d=json.load(open(result));p=json.load(open(prereg))
 assert d["schema"]=="wingless.f61.counterfactual-binding-FIT-source.v1"
 assert d["stage"]=="F61_DISJOINT_FIT_ONLY_SOURCE_NOT_LEARNED"
 assert d["fit_seeds"]==p["fit"]["seeds"] and d["FIT_families"]==p["fit"]["families"]
 assert d["cases"]==len(d["rows"])==p["fit"]["expected_joint_cases"]==128
 assert d["dependent_query_instances_per_case"]==3
 assert d["frozen_counterfactual_binding_source"]
 assert d["diagnostic_unconsumed"] and d["originals_unconsumed"]
 assert d["native_total_576B_unverified"] and d["accepted_refs_unchanged"] and not d["independent_F3"]
 assert set(d["scores"])==set(ARMS)
 assert p["controls"]==list(ARMS)
 seen=set();totals={a:0 for a in ARMS};work={a:0 for a in ARMS[:-1]};same=0
 for row in d["rows"]:
  seed,fam,n=row["seed"],row["family"],row["length"]
  assert (seed,fam,n) not in seen;seen.add((seed,fam,n))
  assert p["fit"]["seeds"][0]<=seed<p["fit"]["seeds"][1] and fam in p["fit"]["families"]
  assert n in (8,10) and len(row["chain"])==n
  assert len(row["secondary_chain"])==6 and len(row["counterfactual_chain"])==6 and len(row["public_bridges"])==3
  assert len(row["public"])==4 and all(len(c)==2 for c in row["public"])
  assert len(row["predictions"])==len(ARMS)
  independent,num=recompute(row)
  assert [tuple(row["predictions"][a]) for a in ARMS[:-1]]==list(independent)
  assert row["exact_joint_candidate_count"]==num
  assert tuple(row["predictions"][ARMS[-1]])==(row["truth"],row["secondary_truth"],row["counterfactual_truth"])
  truth=(row["truth"],row["secondary_truth"],row["counterfactual_truth"])
  full_len=n+6+6;rate=3
  expected=(512*(10+2*rate+full_len),512*(10+2*rate)+full_len,
      3*8*2+8*(10+2*rate+full_len),
      512*(10+2*rate)+full_len,
      4*8*2+full_len,0)
  assert tuple(row["abstract_operations"][a] for a in ARMS)==expected
  for a in ARMS:
   val=int(tuple(row["predictions"][a])==truth)
   assert row["correct"][a]==val
   totals[a]+=val
   if a!=ARMS[-1]:work[a]+=row["abstract_operations"][a]
  same+=int(independent[1]==independent[2])
 assert len(seen)==128 and totals==d["scores"] and work==d["operations"]
 assert d["structured_vs_rolewise_prediction_agreement"]==same
 strong=max(ARMS[1:-1],key=lambda a:totals[a])
 gain=(totals[ARMS[0]]-totals[strong])/128
 accuracy=totals[ARMS[0]]/128
 qualified=(accuracy>=p["FIT_source_gates"]["minimum_exact_accuracy"]
   and gain>=p["FIT_source_gates"]["minimum_full_exact_gain_vs_strongest_affordable"])
 assert d["strongest_affordable_structural_control"]==strong
 assert abs(d["exact_gain_vs_strongest"]-gain)<1e-12
 assert abs(d["exact_public_accuracy"]-accuracy)<1e-12
 assert d["source_gate_qualified"]==qualified
 print("W_F61_SOURCE_AUDIT",json.dumps(dict(correct=totals,exact_accuracy=accuracy,headroom=gain,qualified=qualified),sort_keys=True))
if __name__=="__main__":
 assert len(sys.argv)==3
 audit(sys.argv[1],sys.argv[2])
