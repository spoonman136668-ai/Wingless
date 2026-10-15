"""W-F55: source-feasibility pilot for public-identifiable variable binding.

Public observations include exactly two non-degenerate input/output pairs per
role over F_11.  An exact constructive solver obtains each affine role operator
from those pairs, then composes unseen chains and inverse relationships.
This source file deliberately does NOT train a model or execute DIAG in its
FIT-only qualification entry point. Any scientific learning claim requires
a different prospective frozen evaluation.
"""
from __future__ import annotations
import json, random, sys
FIT=range(601000,601016)
DIAG=range(601100,601116)
FIT_FAM=(1100,1101,1102,1103)
DIAG_FAM=(1104,1105,1106,1107)
ORIGINAL=range(602000,602016)
PRIME=11
SLOPES=(-3,-2,-1,1,2,3)
def source(seed,family,length):
    fit=seed in FIT and family in FIT_FAM
    diag=seed in DIAG and family in DIAG_FAM
    if not (fit or diag):raise ValueError("protected or mixed source identity")
    if length not in ((6,8) if fit else (8,10)):
        raise ValueError("unregistered composition length")
    rng=random.Random(seed*19069+family*4337+length*307+55021)
    # Novel relation semantics are newly sampled per task, not determined
    # by visible context parity or one shared sign for all roles.
    private={}
    public=[]
    for role in range(4):
        slope=rng.choice(SLOPES)
        shift=rng.randrange(PRIME)
        x=rng.randrange(PRIME)
        z=(x+rng.randrange(1,PRIME))%PRIME
        pairs=((x,(slope*x+shift)%PRIME),(z,(slope*z+shift)%PRIME))
        public.append(dict(role=role,pairs=pairs))
        private[role]=(slope,shift)
    roles=[0,1,2,3]
    chain=[]
    for j in range(length):
        # Enforce role reuse and inverse steps; use heldout role-order shift.
        candidate=roles[(j*3+(family-1100)+rng.randrange(4))%4]
        if chain and candidate==chain[-1][0]:
            candidate=(candidate+1)%4
        direction=1 if ((j+family+rng.randrange(3))%3)!=0 else -1
        chain.append((candidate,direction))
    at=rng.randrange(PRIME);start=at
    for role,direction in chain:
        slope,shift=private[role]
        if direction==1:at=(slope*at+shift)%PRIME
        else:at=((at-shift)*pow(slope%PRIME,-1,PRIME))%PRIME
    return dict(seed=seed,family=family,length=length,start=start,
      public_calibration=public,chain=tuple(chain),evaluator_truth=at,
      evaluator_private_maps=private)
def constructive_public_only(case):
    maps={}
    for observation in case["public_calibration"]:
        role=observation["role"]
        (x,y),(z,t)=observation["pairs"]
        if x==z or role in maps:raise ValueError("invalid public calibration")
        slope=((t-y)*pow((z-x)%PRIME,-1,PRIME))%PRIME
        shift=(y-slope*x)%PRIME
        if slope==0:raise ValueError("non-invertible public operator")
        maps[role]=(slope,shift)
    assert len(maps)==4
    at=case["start"]
    for role,direction in case["chain"]:
        slope,shift=maps[role]
        if direction==1:at=(slope*at+shift)%PRIME
        else:at=((at-shift)*pow(slope,-1,PRIME))%PRIME
    return at,maps
def fit_only_report():
    rows=[]
    for seed in FIT:
        for family in FIT_FAM:
            for length in (6,8):
                c=source(seed,family,length)
                answer,recovered=constructive_public_only(c)
                rows.append(dict(seed=seed,family=family,length=length,
                  correct=int(answer==c["evaluator_truth"]),maps=len(recovered)))
    assert len(rows)==128
    correct=sum(row["correct"] for row in rows)
    return dict(schema="wingless.f55.public-role-source-feasibility.fit-only.v1",
       frozen_fit_seeds=[FIT.start,FIT.stop],fit_families=list(FIT_FAM),
       train_fit_cases=len(rows),public_constructive_correct=correct,
       public_constructive_accuracy=correct/len(rows),
       min_required_public_constructive_accuracy=.85,
       eligible_for_future_learning_design=correct/len(rows)>=.85,
       disjoint_diagnostic_unused=True,originals_unused=True,
       native_576B_total_unverified=True,rows=rows)
if __name__=="__main__":
    if len(sys.argv)!=2 or sys.argv[1]!="--fit-only":
        raise SystemExit("Only --fit-only is authorized; no diagnostic from this module")
    print(json.dumps(fit_only_report(),sort_keys=True))
