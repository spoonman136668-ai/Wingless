"""W-F56 FIT-only feasibility for noisy partial public relation observations.

Each of four independent hidden affine maps over F7 has exactly two PUBLIC
input/output observations; observations can be corrupted. Exact public Bayes
enumerates all 14^4 operator combinations, scoring terminal composition under
a properly normalized independent per-role posterior. The fixed two-best
operator beam uses 2^4=16 joint hypotheses, with comparable public evidence.
Diagnostic/primary original source IDs are never accepted by this module.
"""
from __future__ import annotations
import itertools,json,random,sys
FIT=range(621000,621016)
DIAG=range(621100,621116)
FIT_FAMILIES=(1140,1141,1142,1143)
DIAG_FAMILIES=(1144,1145,1146,1147)
ORIGINAL=range(622000,622016)
PRIME=7
SLOPES=(1,2)
CANDIDATES=tuple((a,b) for a in SLOPES for b in range(PRIME))
NOISE=.065
def source(seed,family,length):
    if seed not in FIT or family not in FIT_FAMILIES:
        raise ValueError("FIT-only source; DIAG and originals forbidden")
    if length not in (6,8):raise ValueError("non-preregistered FIT length")
    rng=random.Random(seed*17713+family*4489+length*467+56021)
    hidden={r:rng.choice(CANDIDATES) for r in range(4)}
    public=[];corrupted=0
    for r in range(4):
        points=[]
        for j in range(2):
            x=rng.randrange(PRIME)
            a,b=hidden[r]
            y=(a*x+b)%PRIME
            if rng.random()<NOISE:
                y=(y+rng.randrange(1,PRIME))%PRIME
                corrupted+=1
            points.append((x,y))
        public.append((r,tuple(points)))
    chain=[]
    for step in range(length):
        role=(step*3+family+rng.randrange(4))%4
        if chain and role==chain[-1][0]:role=(role+1)%4
        direction=1 if (step+family+rng.randrange(3))%3 else -1
        chain.append((role,direction))
    start=rng.randrange(PRIME);at=start
    for role,direction in chain:
        a,b=hidden[role]
        if direction==1:at=(a*at+b)%PRIME
        else:at=((at-b)*pow(a,-1,PRIME))%PRIME
    return dict(seed=seed,family=family,length=length,start=start,
                public=tuple(public),chain=tuple(chain),
                private_truth=at,private_maps=hidden,corrupted=corrupted)
def role_posterior(public):
    out={}
    for role,pairs in public:
        raw=[]
        for a,b in CANDIDATES:
            like=1.
            for x,y in pairs:
                like*=1-NOISE if (a*x+b)%PRIME==y else NOISE/(PRIME-1)
            raw.append(like)
        total=sum(raw)
        out[role]=tuple(v/total for v in raw)
    return out
def compose(maps,case):
    at=case["start"]
    for role,d in case["chain"]:
        a,b=maps[role]
        at=(a*at+b)%PRIME if d==1 else ((at-b)*pow(a,-1,PRIME))%PRIME
    return at
def predictor(case,mode):
    posterior=role_posterior(case["public"])
    if mode=="full_exact_public_bayes":
        # Exact joint posterior under the preregistered independent role
        # prior and observation corruption model, not a sampler.
        choices=[tuple(range(len(CANDIDATES))) for _ in range(4)]
    elif mode=="budget_matched_fixed_public_posterior_beam":
        choices=[tuple(sorted(range(len(CANDIDATES)),key=lambda k:(-posterior[r][k],k))[:2])
                 for r in range(4)]
    elif mode=="fixed_known_operator_bank":
        choices=[(max(range(len(CANDIDATES)),key=lambda k:(posterior[r][k],-k)),)
                 for r in range(4)]
    else:raise ValueError("unknown public solver")
    belief=[0.]*PRIME;hypotheses=0
    for idx in itertools.product(*choices):
        probability=1.;maps={}
        for r,k in enumerate(idx):
            probability*=posterior[r][k]
            maps[r]=CANDIDATES[k]
        guess=compose(maps,case)
        belief[guess]+=probability
        hypotheses+=1
    return max(range(PRIME),key=lambda x:(belief[x],-x)),hypotheses*len(case["chain"])
def fit_report():
    cases=[]
    for seed in FIT:
        for family in FIT_FAMILIES:
            for length in (6,8):
                c=source(seed,family,length)
                correct={};work={}
                for arm in ("full_exact_public_bayes","budget_matched_fixed_public_posterior_beam","fixed_known_operator_bank"):
                    guess,operations=predictor(c,arm)
                    correct[arm]=int(guess==c["private_truth"])
                    work[arm]=operations
                cases.append(dict(seed=seed,family=family,length=length,correct=correct,
                                  operations=work,public_noisy_pairs=c["corrupted"]))
    assert len(cases)==128
    arms=("full_exact_public_bayes","budget_matched_fixed_public_posterior_beam","fixed_known_operator_bank")
    count={a:sum(row["correct"][a] for row in cases) for a in arms}
    strongest=max(arms[1:],key=lambda a:count[a])
    gap=(count[arms[0]]-count[strongest])/128
    noised=sum(row["public_noisy_pairs"]>0 for row in cases)
    feasible=count[arms[0]]/128>=.75 and gap>=.05
    return dict(schema="wingless.f56.noisy-public-exact-posterior-fitonly.v1",
      stage="FIT_ONLY_SOURCE_FEASIBILITY_NOT_TRAINED_NOT_DIAGNOSTIC",
      fit_seeds=[FIT.start,FIT.stop],fit_families=list(FIT_FAMILIES),
      fit_cases=128,correct=count,strongest_fixed_control=strongest,
      exact_minus_best_fixed_gap=gap,noisy_source_cases=noised,
      exact_solver_min_gate=.75,budget_headroom_min_gate=.05,
      source_feasibility_pass=feasible,one_exact_posterior_per_case=True,
      diagnostic_unconsumed=True,originals_unconsumed=True,
      native_576B_total_unverified=True,rows=cases)
if __name__=="__main__":
    if len(sys.argv)!=2 or sys.argv[1]!="--fit-only":
        raise SystemExit("FIT-only preflight; no diagnostic access")
    print(json.dumps(fit_report(),sort_keys=True))
