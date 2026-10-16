"""W-F61 FIT-only counterfactual triple-instance typed variable-binding source.

Typed roles act on a public four-element lattice using the eight symmetries
of a square (the dihedral D4 group). The latent role maps obey the PUBLIC
closure constraint R3 o R2 o R1 o R0 = identity, rather than independent
affine shifts. This forces joint role identity across long repeated-role and
inverse task chains. All inference sees exactly the same noisy PUBLIC
observations, role sequence and start; no action sees latent truth.

This module only permits preregistered FIT identities. The explicit exact
posterior is an extra-work reference, never an equal-cost comparator.
"""
from __future__ import annotations
import itertools,json,math,random,sys
FIT=range(751000,751016); FIT_FAMILIES=(1440,1441,1442,1443)
DIAG=range(751100,751116);ORIGINAL=range(752000,752016)
P=4
# Coordinate indices 0..3 represent (x,y) corners of a square.
COORD=((0,0),(1,0),(1,1),(0,1))
def transform(point,rotate,reflection):
 x,y=point
 if reflection:x=1-x
 for _ in range(rotate):x,y=1-y,x
 return COORD.index((x,y))
PERMS=tuple(tuple(transform(v,r,f) for v in COORD)
    for f in (0,1) for r in range(4))
assert len(set(PERMS))==8
IDENTITY=tuple(range(P))
def compose(a,b):
 """a o b (apply b first)."""
 return tuple(a[b[x]] for x in range(P))
def invert(a):
 return tuple(a.index(x) for x in range(P))
assert all(compose(p,invert(p))==IDENTITY for p in PERMS)
MAPS=tuple((a,b,c,invert(compose(c,compose(b,a))))
 for a,b,c in itertools.product(PERMS,repeat=3))
assert len(MAPS)==512
ARMS=("public_full_exact_joint_posterior_extra_work",
      "strong_structured_MAP","budget_matched_public_dp",
      "frozen_F60_public_coupled_MAP","public_local_rolewise",
      "private_oracle_unfair")
def noise(family):
 if family not in FIT_FAMILIES:raise ValueError("not FIT")
 return (.045,.085,.12,.16)[family-1440]
def source(seed,family,length=8):
 if seed not in FIT or family not in FIT_FAMILIES or length not in (8,10):
  raise ValueError("only fresh preregistered F59 FIT source is authorized")
 rng=random.Random(seed*32609+family*5617+length*613+61031)
 hidden=MAPS[rng.randrange(len(MAPS))]
 public=[];corrupt=0
 for role in range(4):
  pairs=[]
  for x in rng.sample(range(P),2):
   y=hidden[role][x]
   if rng.random()<noise(family):
    y=rng.choice([z for z in range(P) if z!=y]);corrupt+=1
   pairs.append((x,y))
  public.append(tuple(pairs))
 # Two dependent variable-binding problems from one persistent hidden
 # operator assignment; the second starts at the first public query's
 # latent endpoint. Neither answer is an inference-time observation.
 chain=[(r,1 if j%3 else -1) for j,r in enumerate((0,1,2,3,0,2,1,3))]
 if length==10:chain.extend(((2,-1),(1,1)))
 rng.shuffle(chain)
 start=rng.randrange(P)
 answer=apply_chain(hidden,start,chain)
 second_chain=tuple((r,1 if j%2 else -1) for j,r in enumerate((1,3,0,2,3,1)))
 # A separate counterfactual query shares the same hidden role identity
 # but starts from a newly sampled public position; no answer is exposed.
 counterfactual_start=rng.randrange(P)
 counterfactual_chain=tuple((r,1 if j%3 else -1) for j,r in enumerate((2,0,3,1,2,1)))
 second_answer=apply_chain(hidden,answer,second_chain)
 counterfactual_answer=apply_chain(hidden,counterfactual_start,counterfactual_chain)
 # Public pair-composition probes are sampled independently of the
 # private queries, and are made available to EVERY registered solver.
 bridge=[]
 for a,b in ((0,2),(1,3),(2,0)):
  x=rng.randrange(P);y=hidden[b][hidden[a][x]]
  if rng.random()<noise(family):
   y=rng.choice([z for z in range(P) if z!=y]);corrupt+=1
  bridge.append((a,b,x,y))
 return dict(seed=seed,family=family,length=length,start=start,
    public=tuple(public),chain=tuple(chain),truth=answer,
    second_chain=second_chain,second_truth=second_answer,
    counterfactual_start=counterfactual_start,
    counterfactual_chain=counterfactual_chain,
    counterfactual_truth=counterfactual_answer,
    public_bridges=tuple(bridge),private_maps=hidden,corrupted=corrupt)
def apply_chain(maps,start,chain):
 at=start
 for role,d in chain:
  p=maps[role]
  at=p[at] if d==1 else invert(p)[at]
 return at
def likelihood(maps,observations,rate,bridges=()):
 w=1.
 for role,pairs in enumerate(observations):
  for x,y in pairs:
   w*=1-rate if maps[role][x]==y else rate/(P-1)
 for a,b,x,y in bridges:
  guess=maps[b][maps[a][x]]
  w*=1-rate if guess==y else rate/(P-1)
 return w
def evaluate(case):
 public=case["public"];rate=noise(case["family"])
 weights=[likelihood(m,public,rate,case["public_bridges"]) for m in MAPS]
 total=sum(weights);assert total>0
 weights=[v/total for v in weights]
 def task_pair(m):
  first=apply_chain(m,case["start"],case["chain"])
  second=apply_chain(m,first,case["second_chain"])
  counter=apply_chain(m,case["counterfactual_start"],case["counterfactual_chain"])
  return (first,second,counter)
 joint={}
 for m,w in zip(MAPS,weights):
  pair=task_pair(m)
  joint[pair]=joint.get(pair,0.)+w
 exact=max(sorted(joint),key=lambda pair:joint[pair])
 best=max(range(len(MAPS)),key=lambda k:(weights[k],-k))
 structured=task_pair(MAPS[best])
 # Strong public independent-role MAP; ignores cross-role closure but is
 # not arbitrarily weakened or denied publicly available calibration.
 ranks=[]
 for role,pairs in enumerate(public[:3]):
  local=[math.prod((1-rate if perm[x]==y else rate/3) for x,y in pairs)
       for perm in PERMS]
  ranks.append(sorted(range(8),key=lambda k:(-local[k],k)))
 rolewise=[]
 for pairs in public:
  local=[math.prod((1-rate if perm[x]==y else rate/3) for x,y in pairs) for perm in PERMS]
  rolewise.append(PERMS[max(range(8),key=lambda k:(local[k],-k))])
 local=task_pair(tuple(rolewise))
 beam_votes={}
 for indices in itertools.product(*[r[:2] for r in ranks]):
  a,b,c=(PERMS[k] for k in indices)
  maps=(a,b,c,invert(compose(c,compose(b,a))))
  pair=task_pair(maps)
  weight=likelihood(maps,public,rate,case["public_bridges"])
  beam_votes[pair]=beam_votes.get(pair,0.)+weight
 beam=max(sorted(beam_votes),key=lambda pair:beam_votes[pair])
 result=dict(zip(ARMS,(exact,structured,beam,structured,local,
                       (case["truth"],case["second_truth"],case["counterfactual_truth"]))))
 length=case["length"]+len(case["second_chain"])+len(case["counterfactual_chain"])
 costs=dict(zip(ARMS,(512*(10+2*len(case["public_bridges"])+length),
  512*(10+2*len(case["public_bridges"]))+length,
  3*8*2+8*(10+2*len(case["public_bridges"])+length),
  512*(10+2*len(case["public_bridges"]))+length,
  4*8*2+length,0)))
 return result,costs,len(joint)
def fit_report():
 rows=[];correct={a:0 for a in ARMS};operations={a:0 for a in ARMS[:-1]}
 ablation=0
 for seed in FIT:
  for family in FIT_FAMILIES:
   for length in (8,10):
    c=source(seed,family,length);ans,cost,joint_count=evaluate(c)
    d={a:int(tuple(ans[a])==(c["truth"],c["second_truth"],c["counterfactual_truth"])) for a in ARMS}
    for a in ARMS:correct[a]+=d[a]
    for a in ARMS[:-1]:operations[a]+=cost[a]
    ablation+=int(ans[ARMS[1]]==ans[ARMS[2]])
    rows.append(dict(seed=seed,family=family,length=c["length"],
      public=[[[x,y] for x,y in pair] for pair in c["public"]],
      public_bridges=[list(q) for q in c["public_bridges"]],
      secondary_chain=[list(v) for v in c["second_chain"]],
      secondary_truth=c["second_truth"],
      counterfactual_start=c["counterfactual_start"],
      counterfactual_chain=[list(v) for v in c["counterfactual_chain"]],
      counterfactual_truth=c["counterfactual_truth"],
      exact_joint_candidate_count=joint_count,
      chain=[list(x) for x in c["chain"]],start=c["start"],
      truth=c["truth"],public_corruptions=c["corrupted"],
      predictions=ans,correct=d,abstract_operations=cost))
 assert len(rows)==128
 exact=correct[ARMS[0]]/len(rows)
 strongest=max(ARMS[1:-1],key=lambda a:correct[a])
 gain=(correct[ARMS[0]]-correct[strongest])/len(rows)
 eligible=exact>=.8 and gain>=.05
 return dict(schema="wingless.f61.counterfactual-binding-FIT-source.v1",
   stage="F61_DISJOINT_FIT_ONLY_SOURCE_NOT_LEARNED",
   fit_seeds=[FIT.start,FIT.stop],FIT_families=list(FIT_FAMILIES),
   cases=len(rows),required_prereg_cases=128,case_count_matches_prereg=True,
   scores=correct,operations=operations,rows=rows,
   strongest_affordable_structural_control=strongest,
   exact_gain_vs_strongest=gain,exact_public_accuracy=exact,
   structured_vs_rolewise_prediction_agreement=ablation,
   source_gate_qualified=eligible,
   dependent_query_instances_per_case=3,
   frozen_counterfactual_binding_source=True,
   diagnostic_unconsumed=True,originals_unconsumed=True,
   independent_F3=False,native_total_576B_unverified=True,
   accepted_refs_unchanged=True)
if __name__=="__main__":
 if sys.argv[1:]!=["--fit-only"]:raise SystemExit("FIT source only")
 print(json.dumps(fit_report(),sort_keys=True))
