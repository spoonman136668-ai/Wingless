"""W-F60 FIT-only multi-instance dependent typed variable-binding source.

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
FIT=range(721000,721016); FIT_FAMILIES=(1380,1381,1382,1383)
DIAG=range(721100,721116);ORIGINAL=range(722000,722016)
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
ARMS=("public_exact_structured_solver_extra_work",
      "public_structured_MAP","fixed_public_rolewise_composition",
      "matched_work_symbolic_state","unfair_private_truth_oracle")
def noise(family):
 if family not in FIT_FAMILIES:raise ValueError("not FIT")
 return (.045,.07,.09,.11)[family-1380]
def source(seed,family,length=8):
 if seed not in FIT or family not in FIT_FAMILIES or length not in (8,10):
  raise ValueError("only fresh preregistered F59 FIT source is authorized")
 rng=random.Random(seed*31091+family*5279+length*571+60031)
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
 second_chain=tuple((r,-d) for r,d in reversed(chain))
 second_answer=apply_chain(hidden,answer,second_chain)
 # Public pair-composition probes are sampled independently of the
 # private queries, and are made available to EVERY registered solver.
 bridge=[]
 for a,b in ((0,2),(1,3)):
  x=rng.randrange(P);y=hidden[b][hidden[a][x]]
  if rng.random()<noise(family):
   y=rng.choice([z for z in range(P) if z!=y]);corrupt+=1
  bridge.append((a,b,x,y))
 return dict(seed=seed,family=family,length=length,start=start,
    public=tuple(public),chain=tuple(chain),truth=answer,
    second_chain=second_chain,second_truth=second_answer,
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
 public=case["public"];rate=noise(case["family"]);weights=[]
 for maps in MAPS:weights.append(likelihood(maps,public,rate,case["public_bridges"]))
 total=sum(weights);assert total>0
 weights=[v/total for v in weights]
 votes=[0.]*P
 for maps,w in zip(MAPS,weights):
  votes[apply_chain(maps,case["start"],case["chain"])]+=w
 exact=max(range(P),key=lambda z:(votes[z],-z))
 best_map=max(range(len(MAPS)),key=lambda k:(weights[k],-k))
 structured_map=apply_chain(MAPS[best_map],case["start"],case["chain"])
 # Independent public per-role MAP drops the global closure constraint.
 choices=[]
 for role,pairs in enumerate(public):
  ww=[math.prod((1-rate if p[x]==y else rate/(P-1)) for x,y in pairs)
       for p in PERMS]
  choices.append(max(range(8),key=lambda k:(ww[k],-k)))
 local=apply_chain(tuple(PERMS[x] for x in choices),case["start"],case["chain"])
 # Matched-work structurally valid restricted beam: choose 2 locally
 # preferred operators for each of roles 0,1,2, deriving role 3 via
 # PUBLIC closure. This sees the same evidence with 8 joint branches.
 rankings=[]
 for role,pairs in enumerate(public[:3]):
  ww=[math.prod((1-rate if p[x]==y else rate/(P-1)) for x,y in pairs)
      for p in PERMS]
  rankings.append(sorted(range(8),key=lambda k:(-ww[k],k))[:2])
 beam_votes=[0.]*P
 for indices in itertools.product(*rankings):
  maps=(PERMS[indices[0]],PERMS[indices[1]],PERMS[indices[2]])
  last=invert(compose(maps[2],compose(maps[1],maps[0])))
  allmaps=maps+(last,)
  w=likelihood(allmaps,public,rate,case["public_bridges"])
  beam_votes[apply_chain(allmaps,case["start"],case["chain"])]+=w
 beam=max(range(P),key=lambda z:(beam_votes[z],-z))
 answers=(exact,structured_map,local,beam,case["truth"])
 # Two dependent tasks share a fixed hidden map. The secondary chain
 # is a novel variable binding, but no inference method receives its
 # true starting value; a separate marginal posterior is evaluated.
 secondary={}
 for maps,w in zip(MAPS,weights):
  pair=(apply_chain(maps,case["start"],case["chain"]),
        apply_chain(maps,apply_chain(maps,case["start"],case["chain"]),case["second_chain"]))
  secondary[pair]=secondary.get(pair,0)+w
 return dict(zip(ARMS,answers)),dict(zip(ARMS,(512*(14+case["length"]),512*14+case["length"],
     4*8*2+case["length"],3*8*2+8*(14+case["length"]),0))),secondary
 return dict(zip(ARMS,answers)),dict(zip(ARMS,costs))
def fit_report():
 rows=[];correct={a:0 for a in ARMS};operations={a:0 for a in ARMS[:-1]}
 ablation=0
 for seed in FIT:
  for family in FIT_FAMILIES:
   for length in (8,10):
    c=source(seed,family,length);ans,cost,joint=evaluate(c)
    d={a:int(ans[a]==c["truth"]) for a in ARMS}
    for a in ARMS:correct[a]+=d[a]
    for a in ARMS[:-1]:operations[a]+=cost[a]
    ablation+=int(ans[ARMS[1]]==ans[ARMS[2]])
    rows.append(dict(seed=seed,family=family,length=c["length"],
      public=[[[x,y] for x,y in pair] for pair in c["public"]],
      public_bridges=[list(q) for q in c["public_bridges"]],
      secondary_chain=[list(v) for v in c["second_chain"]],
      secondary_truth=c["second_truth"],
      exact_joint_candidate_count=len(joint),
      chain=[list(x) for x in c["chain"]],start=c["start"],
      truth=c["truth"],public_corruptions=c["corrupted"],
      predictions=ans,correct=d,abstract_operations=cost))
 assert len(rows)==128
 exact=correct[ARMS[0]]/len(rows)
 strongest=max(ARMS[1:-1],key=lambda a:correct[a])
 gain=(correct[ARMS[0]]-correct[strongest])/len(rows)
 eligible=exact>=.8 and gain>=.05
 return dict(schema="wingless.f60.multi-instance-binding-FIT-source.v1",
   stage="F60_DISJOINT_FIT_ONLY_SOURCE_NOT_LEARNED",
   fit_seeds=[FIT.start,FIT.stop],FIT_families=list(FIT_FAMILIES),
   cases=len(rows),required_prereg_cases=128,case_count_matches_prereg=True,
   scores=correct,operations=operations,rows=rows,
   strongest_affordable_structural_control=strongest,
   exact_gain_vs_strongest=gain,exact_public_accuracy=exact,
   structured_vs_rolewise_prediction_agreement=ablation,
   source_gate_qualified=eligible,
   dependent_query_instances_per_case=2,
   frozen_counterfactual_binding_source=True,
   diagnostic_unconsumed=True,originals_unconsumed=True,
   independent_F3=False,native_total_576B_unverified=True,
   accepted_refs_unchanged=True)
if __name__=="__main__":
 if sys.argv[1:]!=["--fit-only"]:raise SystemExit("FIT source only")
 print(json.dumps(fit_report(),sort_keys=True))
