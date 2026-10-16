"""W-F59 source-free all-row public constraint and frozen gate audit.

The reference D4 permutations are reconstructed locally from public square
geometry, independent of the registered source generator. The auditor
recomputes exact Bayesian posterior, strong public coupled MAP, independent
role-wise MAP and restricted public symbolic beam from archived public pairs.
"""
import itertools,json,math,sys
C=((0,0),(1,0),(1,1),(0,1))
def perm(turn,flip):
 out=[]
 for x,y in C:
  if flip:x=1-x
  for z in range(turn):x,y=1-y,x
  out.append(C.index((x,y)))
 return tuple(out)
P=tuple(perm(t,f) for f in (0,1) for t in range(4))
def comp(a,b):return tuple(a[b[j]] for j in range(4))
def inv(a):return tuple(a.index(j) for j in range(4))
M=tuple((a,b,c,inv(comp(c,comp(b,a)))) for a,b,c in itertools.product(P,repeat=3))
ARMS=("public_exact_structured_solver_extra_work","public_structured_MAP",
      "fixed_public_rolewise_composition","matched_work_symbolic_state",
      "unfair_private_truth_oracle")
def calc(case):
 public=case["public"];rate=(.025,.05,.075,.10)[case["family"]-1340]
 def lik(mm):
  w=1.
  for i,pairs in enumerate(public):
   for x,y in pairs:w*=1-rate if mm[i][x]==y else rate/3
  return w
 def answer(mm):
  now=case["start"]
  for role,d in case["chain"]:
   m=mm[role]
   now=m[now] if d==1 else inv(m)[now]
  return now
 wt=[lik(mm) for mm in M];votes=[0.]*4
 for mm,w in zip(M,wt):votes[answer(mm)]+=w
 exact=max(range(4),key=lambda z:(votes[z],-z))
 coupled=answer(M[max(range(len(M)),key=lambda j:(wt[j],-j))])
 ranked=[]
 for role,pairs in enumerate(public):
  weights=[]
  for p in P:
   a=math.prod((1-rate if p[x]==y else rate/3) for x,y in pairs)
   weights.append(a)
  ranked.append(sorted(range(8),key=lambda k:(-weights[k],k)))
 local=answer(tuple(P[ranked[j][0]] for j in range(4)))
 beam=[0.]*4
 for idx in itertools.product(*[r[:2] for r in ranked[:3]]):
  a,b,c=(P[k] for k in idx)
  mm=(a,b,c,inv(comp(c,comp(b,a))))
  beam[answer(mm)]+=lik(mm)
 affordable=max(range(4),key=lambda z:(beam[z],-z))
 return exact,coupled,local,affordable
def audit(result,pre):
 d=json.load(open(result));p=json.load(open(pre))
 assert d["schema"]=="wingless.f59.typed-binding-fit-source.v1"
 assert d["stage"]=="NEW_DISJOINT_FIT_ONLY_SOURCE_FEASIBILITY_NOT_LEARNED"
 assert d["fit_seeds"]==p["fit"]["seeds"] and d["FIT_families"]==p["fit"]["families"]
 assert d["cases"]==len(d["rows"])==p["fit"]["expected_cases"]==128
 assert d["case_count_matches_prereg"] and d["required_prereg_cases"]==128
 assert d["diagnostic_unconsumed"] and d["originals_unconsumed"]
 assert d["native_total_576B_unverified"] and d["accepted_refs_unchanged"] and not d["independent_F3"]
 assert set(d["scores"])==set(ARMS)
 seen=set();scores={a:0 for a in ARMS};work={a:0 for a in ARMS[:-1]};same=0
 for row in d["rows"]:
  key=row["seed"],row["family"],row["length"]
  assert key not in seen;seen.add(key)
  assert p["fit"]["seeds"][0]<=key[0]<p["fit"]["seeds"][1]
  assert key[1] in p["fit"]["families"] and key[2] in (8,10)
  assert len(row["public"])==4 and all(len(pairs)==3 for pairs in row["public"])
  assert len(row["chain"])==key[2]
  assert set(row["predictions"])==set(ARMS)
  assert set(row["correct"])==set(ARMS)
  assert set(row["abstract_operations"])==set(ARMS)
  pred=calc(row)
  assert tuple(row["predictions"][a] for a in ARMS[:-1])==pred
  assert row["predictions"][ARMS[-1]]==row["truth"]
  same+=int(pred[1]==pred[2])
  req=(512*(12+key[2]),512*12+key[2],4*8*3+key[2],
      3*8*3+8*(12+key[2]),0)
  assert tuple(row["abstract_operations"][a] for a in ARMS)==req
  for arm in ARMS:
   valid=int(row["predictions"][arm]==row["truth"])
   assert valid==row["correct"][arm]
   scores[arm]+=valid
   if arm!=ARMS[-1]:work[arm]+=row["abstract_operations"][arm]
 assert len(seen)==128 and d["scores"]==scores and d["operations"]==work
 assert d["structured_vs_rolewise_prediction_agreement"]==same
 strong=max(ARMS[1:-1],key=lambda a:scores[a])
 gain=(scores[ARMS[0]]-scores[strong])/128
 exact=scores[ARMS[0]]/128
 qualified=exact>=p["FIT_source_requirements"]["public_exact_source_accuracy_min"] and gain>=p["FIT_source_requirements"]["full_public_minus_strongest_affordable_structural_control_min"]
 assert d["strongest_affordable_structural_control"]==strong
 assert abs(d["exact_gain_vs_strongest"]-gain)<1e-12
 assert abs(d["exact_public_accuracy"]-exact)<1e-12
 assert d["source_gate_qualified"]==qualified
 print("W_F59_SOURCE_AUDIT",json.dumps(dict(scores=scores,strong=strong,gain=gain,exact_accuracy=exact,qualified=qualified),sort_keys=True))
if __name__=="__main__":
 assert len(sys.argv)==3
 audit(sys.argv[1],sys.argv[2])
