"""W-F57 FIT-only public correlated-role feasibility test; no diagnostic use."""
import itertools,json,random,sys
FIT=range(641000,641016); DIAG=range(641100,641116)
FIT_FAM=(1180,1181,1182,1183); DIAG_FAM=(1184,1185,1186,1187)
ORIGINAL=range(642000,642016); P=7; NOISE=.11
MAPS=tuple(((s,b0),(t,b1),(s,b2),(t,(b0+b1-b2)%P))
 for s in (1,2) for t in (1,2)
 for b0 in range(P) for b1 in range(P) for b2 in range(P))
ARMS=("full_exact_public_coupled_posterior","fixed_equal_work_cross_role_beam",
      "fixed_structured_symbolic_pruning","public_algebraic_consistency")
def compose(maps,start,chain):
 at=start
 for role,d in chain:
  a,b=maps[role]
  at=(a*at+b)%P if d==1 else ((at-b)*pow(a,-1,P))%P
 return at
def source(seed,family,length):
 if seed not in FIT or family not in FIT_FAM:raise ValueError("FIT-only registered source")
 if length not in (6,8):raise ValueError("FIT chain length")
 rng=random.Random(seed*20269+family*4289+length*443+57013)
 hidden=MAPS[rng.randrange(len(MAPS))]
 public=[];corrupt=0
 for role in range(4):
  pairs=[]
  for j in range(2):
   x=rng.randrange(P);a,b=hidden[role];y=(a*x+b)%P
   if rng.random()<NOISE:y=(y+rng.randrange(1,P))%P;corrupt+=1
   pairs.append((x,y))
  public.append(tuple(pairs))
 chain=[]
 for j in range(length):
  role=(j*3+family+rng.randrange(4))%4
  if chain and role==chain[-1][0]:role=(role+1)%4
  chain.append((role,1 if (j+family+rng.randrange(3))%3 else -1))
 start=rng.randrange(P)
 return dict(seed=seed,family=family,length=length,public=public,chain=chain,
    start=start,truth=compose(hidden,start,chain),noisy_pairs=corrupt)
def likelihood(maps,observations):
 prob=1.
 for role in range(4):
  a,b=maps[role]
  for x,y in observations[role]:
   prob*=1-NOISE if (a*x+b)%P==y else NOISE/6
 return prob
def local_top(observations):
 result=[]
 for role in range(4):
  local=[]
  for a in (1,2):
   for b in range(P):
    score=1.
    for x,y in observations[role]:
     score*=1-NOISE if (a*x+b)%P==y else NOISE/6
    local.append(((a,b),score))
  result.append(sorted(local,key=lambda x:(-x[1],x[0])))
 return result
def predict(d,arm):
 votes=[0.]*P;length=len(d["chain"]);obs=d["public"]
 if arm==ARMS[0]:
  for maps in MAPS:votes[compose(maps,d["start"],d["chain"])]+=likelihood(maps,obs)
  work=len(MAPS)*(8+length)
 elif arm==ARMS[1]:
  local=local_top(obs)
  for idx in itertools.product((0,1),repeat=4):
   maps=tuple(local[r][idx[r]][0] for r in range(4))
   score=1.
   for r in range(4):score*=local[r][idx[r]][1]
   votes[compose(maps,d["start"],d["chain"])]+=score
  work=112+16*length
 elif arm==ARMS[2]:
  # Correlated constrained MAP control; scans all public hypothesis scores.
  winning=max(MAPS,key=lambda m:likelihood(m,obs))
  votes[compose(winning,d["start"],d["chain"])]=1.
  work=len(MAPS)*8+length
 elif arm==ARMS[3]:
  maps=tuple(local_top(obs)[r][0][0] for r in range(4))
  votes[compose(maps,d["start"],d["chain"])]=1.
  work=112+length
 else:raise ValueError("unknown control")
 return max(range(P),key=lambda x:(votes[x],-x)),work
def fit_report():
 rows=[];correct={a:0 for a in ARMS};work={a:0 for a in ARMS}
 for seed in FIT:
  for family in FIT_FAM:
   for length in (6,8):
    c=source(seed,family,length);score={};cost={}
    for a in ARMS:
     guess,ops=predict(c,a);score[a]=int(guess==c["truth"]);cost[a]=ops
     correct[a]+=score[a];work[a]+=ops
    rows.append(dict(seed=seed,family=family,length=length,
         correctness=score,operation_counts=cost,noisy_pairs=c["noisy_pairs"]))
 best=max(ARMS[1:],key=lambda a:correct[a])
 gain=(correct[ARMS[0]]-correct[best])/128
 good=correct[ARMS[0]]>=96 and gain>=.05
 return dict(schema="wingless.f57.public-coupled-fitonly.v1",
  stage="FIT_ONLY_NO_DIAG",FIT_seeds=[FIT.start,FIT.stop],FIT_families=list(FIT_FAM),
  cases=len(rows),hypotheses=len(MAPS),correct=correct,work=work,
  strongest_fair=best,exact_headroom=gain,source_feasible=good,
  diagnostic_unused=True,originals_unused=True,
  native_576B_total_unverified=True,source_independent_F3=False,rows=rows)
if __name__=="__main__":
 if len(sys.argv)!=2 or sys.argv[1]!="--fit-only":raise SystemExit("FIT only")
 print(json.dumps(fit_report(),sort_keys=True))
