"""Wingless F57: public-belief relational transfer method source, FIT only.

Frozen precode: source calls are *strictly limited* to F57 FIT seed/family
identities. Role maps are private evaluator labels; policies receive only
public noisy calibration and a variable-binding chain. This implements a
factorized learned public-belief transition, not another joint map top-k
beam. No diagnostic, original or policy admission is possible here.
"""
from __future__ import annotations

import itertools
import json
import math
import random
import sys

FIT=range(661000,661016)
FAMILY=(1240,1241,1242,1243)
DIAG=range(661100,661116)
ORIGINAL=range(662000,662016)
P=7
MAPS=tuple((a,b) for a in (1,2) for b in range(P))
TRAIN=range(661000,661008)
VALID=range(661008,661016)

def observation_noise(family):
    if family not in FAMILY:raise ValueError("unknown FIT source family")
    return .035 if family%2==0 else .09

def source(seed,family,length):
    if seed not in FIT or family not in FAMILY or length not in (8,10):
        raise ValueError("only preregistered F57 FIT source; DIAG and original prohibited")
    rng=random.Random(seed*57047+family*7127+length*671+57)
    private={role:rng.choice(MAPS) for role in range(4)}
    noise=observation_noise(family)
    public=[]
    for role in range(4):
        xs=rng.sample(range(P),2)
        observations=[]
        for x in xs:
            slope,shift=private[role]
            y=(slope*x+shift)%P
            if rng.random()<noise:
                y=(y+rng.randrange(1,P))%P
            observations.append((x,y))
        public.append((role,tuple(observations)))
    chain=[]
    for i in range(length):
        # Future relation semantics are resampled per FIT case.
        role=(i*3+family+rng.randrange(4))%4
        if chain and role==chain[-1][0]:
            role=(role+1)%4
        direction=-1 if (i+family+rng.randrange(3))%4==0 else 1
        chain.append((role,direction))
    start=rng.randrange(P);value=start
    for role,direction in chain:
        a,b=private[role]
        value=(a*value+b)%P if direction==1 else ((value-b)*pow(a,-1,P))%P
    return {"seed":seed,"family":family,"length":length,"start":start,
        "public":tuple(public),"chain":tuple(chain),
        "evaluator_truth":value,"evaluator_private_maps":private}

def posterior(case):
    likelihood={}
    for role,pairs in case["public"]:
        weights=[]
        for a,b in MAPS:
            probability=1.
            for x,y in pairs:
                predicted=(a*x+b)%P
                noise=observation_noise(case["family"])
                probability*=1-noise if predicted==y else noise/(P-1)
            weights.append(probability)
        total=sum(weights)
        likelihood[role]=tuple(weight/total for weight in weights)
    if set(likelihood)!={0,1,2,3}:raise ValueError("missing public role observation")
    return likelihood

def compose(maps,case):
    current=case["start"]
    for role,direction in case["chain"]:
        a,b=maps[role]
        current=(a*current+b)%P if direction==1 else ((current-b)*pow(a,-1,P))%P
    return current

def exhaustive_public_bayes(case):
    probs=posterior(case);term=[0.]*P
    for combination in itertools.product(range(len(MAPS)),repeat=4):
        weight=math.prod(probs[role][idx] for role,idx in enumerate(combination))
        where=compose({role:MAPS[idx] for role,idx in enumerate(combination)},case)
        term[where]+=weight
    return max(range(P),key=lambda n:(term[n],-n)),len(MAPS)**4*len(case["chain"])

def frozen_F56_budgeted_beam(case):
    probs=posterior(case);term=[0.]*P
    choices=[sorted(range(len(MAPS)),key=lambda k:(-probs[r][k],k))[:2] for r in range(4)]
    for combination in itertools.product(*choices):
        weight=math.prod(probs[role][idx] for role,idx in enumerate(combination))
        where=compose({role:MAPS[idx] for role,idx in enumerate(combination)},case)
        term[where]+=weight
    return max(range(P),key=lambda n:(term[n],-n)),16*len(case["chain"])

def fixed_known_operator_bank(case):
    probs=posterior(case)
    maps={role:MAPS[max(range(len(MAPS)),key=lambda k:(probs[role][k],-k))] for role in range(4)}
    return compose(maps,case),len(case["chain"])

def factorized_transition(case,temperatures):
    """A genuine distributional recurrent state (7 values), not a map beam.

    Each role posterior is calibrated with a FIT-learned temperature and
    applied to the evolving probability distribution over possible bindings.
    Reusing a role requires a future genuinely learned persistent state model
    to avoid unjustified independent resampling. This limitation is explicit.
    """
    probs=posterior(case)
    maps={}
    for role in range(4):
        w=[max(0.,probs[role][k])**temperatures[role] for k in range(len(MAPS))]
        total=sum(w)
        maps[role]=[x/total for x in w]
    belief=[0.]*P
    belief[case["start"]]=1.
    for role,direction in case["chain"]:
        next_belief=[0.]*P
        for state,mass in enumerate(belief):
            if mass<=0.:continue
            for (a,b),prob in zip(MAPS,maps[role]):
                target=(a*state+b)%P if direction==1 else ((state-b)*pow(a,-1,P))%P
                next_belief[target]+=mass*prob
        belief=next_belief
    return max(range(P),key=lambda n:(belief[n],-n)),len(case["chain"])*P*len(MAPS)

def train_fit_temperature():
    """FIT labels only; calibration cannot see new F57 diagnostic questions."""
    counts={role:[[0.,0.] for _ in range(4)] for role in range(4)}
    # Small predetermined temperature bank, trained separately by role.
    temperatures=(.75,1.,1.5,2.)
    for seed in TRAIN:
        for family in FAMILY:
            for length in (8,10):
                case=source(seed,family,length)
                probs=posterior(case)
                for role in range(4):
                    truth=MAPS.index(case["evaluator_private_maps"][role])
                    for i,temperature in enumerate(temperatures):
                        normalized=[x**temperature for x in probs[role]]
                        z=sum(normalized)
                        counts[role][i][0]+=math.log(max(normalized[truth]/z,1e-12))
                        counts[role][i][1]+=1
    fitted=[]
    for role in range(4):
        idx=max(range(len(temperatures)),key=lambda i:counts[role][i][0])
        fitted.append(temperatures[idx])
    return tuple(fitted),8*4*2

def fit_preview():
    weights,ntrain=train_fit_temperature()
    rows=[]
    for seed in VALID:
        for family in FAMILY:
            for length in (8,10):
                case=source(seed,family,length)
                a=factorized_transition(case,weights)
                b=frozen_F56_budgeted_beam(case)
                c=fixed_known_operator_bank(case)
                rows.append({"seed":seed,"family":family,"length":length,
                    "correct":{"learned_public_belief_state":int(a[0]==case["evaluator_truth"]),
                               "frozen_F56_posterior_beam":int(b[0]==case["evaluator_truth"]),
                               "known_operator_MAP":int(c[0]==case["evaluator_truth"])},
                    "abstract_work":{"learned_public_belief_state":a[1],
                                     "frozen_F56_posterior_beam":b[1],
                                     "known_operator_MAP":c[1]}})
    return {"schema":"wingless.f57.public-relational-belief.fit-only.v1",
      "stage":"WITHIN_FIT_CALIBRATION_AND_VALIDATION_NO_DIAG",
      "train_fit_seeds":[TRAIN.start,TRAIN.stop],
      "fit_validation_seeds":[VALID.start,VALID.stop],
      "fit_families":list(FAMILY),"fit_calibration_examples":ntrain,
      "learned_temperatures":list(weights),
      "fit_validation_cases":len(rows),"rows":rows,
      "diag_unconsumed":True,"original_unconsumed":True,
      "independent_F3":False,"physical_native_576B_unverified":True}

if __name__=="__main__":
    if sys.argv[1:]!=["--fit-only"]:
        raise SystemExit("FIT only: diagnostics and originals forbidden")
    print(json.dumps(fit_preview(),sort_keys=True))
