"""Wingless F57 independent FIT-only public-belief and frozen-beam recount.

No source-generator import; candidate and fixed-control predictions are
reconstructed from recorded public observations after FIT validation.
"""
import itertools
import json
import math
import pathlib
import sys

P=7
MAPS=tuple((a,b) for a in (1,2) for b in range(P))
ARMS=("learned_public_belief_state","frozen_F56_posterior_beam","known_operator_MAP")
def posterior(row):
    noise=.035 if row["family"]%2==0 else .09
    result={}
    for entry in row["public"]:
        role=entry["role"]
        assert role in (0,1,2,3) and role not in result
        assert len(entry["pairs"])==2
        raw=[]
        for a,b in MAPS:
            z=1.
            for x,y in entry["pairs"]:
                assert 0<=x<P and 0<=y<P
                z*=1-noise if (a*x+b)%P==y else noise/(P-1)
            raw.append(z)
        total=sum(raw)
        result[role]=[x/total for x in raw]
    assert set(result)=={0,1,2,3}
    return result
def compose(case,maps):
    at=case["start"]
    for role,d in case["chain"]:
        a,b=maps[role]
        at=(a*at+b)%P if d==1 else ((at-b)*pow(a,-1,P))%P
    return at
def answer(row,temps):
    post=posterior(row)
    assert len(row["chain"])==row["length"]
    b=[0.]*P;b[row["start"]]=1.
    tempered={}
    for r in range(4):
        weights=[p**temps[r] for p in post[r]]
        z=sum(weights)
        tempered[r]=[x/z for x in weights]
    for role,d in row["chain"]:
        new=[0.]*P
        for old,p in enumerate(b):
            for (a,offset),v in zip(MAPS,tempered[role]):
                dst=(a*old+offset)%P if d==1 else ((old-offset)*pow(a,-1,P))%P
                new[dst]+=p*v
        b=new
    learned=max(range(P),key=lambda x:(b[x],-x))
    choices=[sorted(range(len(MAPS)),key=lambda k:(-post[r][k],k))[:2] for r in range(4)]
    scores=[0.]*P
    for choice in itertools.product(*choices):
        prob=math.prod(post[r][k] for r,k in enumerate(choice))
        dst=compose(row,{r:MAPS[k] for r,k in enumerate(choice)})
        scores[dst]+=prob
    beam=max(range(P),key=lambda x:(scores[x],-x))
    maps={r:MAPS[max(range(len(MAPS)),key=lambda k:(post[r][k],-k))] for r in range(4)}
    known=compose(row,maps)
    n=row["length"]
    return {"learned_public_belief_state":learned,
            "frozen_F56_posterior_beam":beam,"known_operator_MAP":known}, {
            "learned_public_belief_state":n*P*len(MAPS),
            "frozen_F56_posterior_beam":16*n,
            "known_operator_MAP":n}
def audit(data,pre):
    assert pre["schema"]=="wingless.f57.uncertainty-aware-compositional-belief.prereg.v1"
    assert data["schema"]=="wingless.f57.public-relational-belief.fit-only.v1"
    assert data["stage"]=="WITHIN_FIT_CALIBRATION_AND_VALIDATION_NO_DIAG"
    assert data["train_fit_seeds"]==[661000,661008]
    assert data["fit_validation_seeds"]==[661008,661016]
    assert data["fit_families"]==pre["fit"]["families"]
    assert data["fit_calibration_examples"]==64
    assert data["fit_validation_cases"]==len(data["rows"])==64
    assert data["diag_unconsumed"] and data["original_unconsumed"]
    assert data["independent_F3"] is False and data["physical_native_576B_unverified"] is True
    temps=data["learned_temperatures"]
    assert len(temps)==4 and all(x in (.75,1.,1.5,2.) for x in temps)
    seen=set();counts={a:0 for a in ARMS}
    for row in data["rows"]:
        key=row["seed"],row["family"],row["length"]
        assert key not in seen
        seen.add(key)
        assert 661008<=row["seed"]<661016
        assert row["family"] in pre["fit"]["families"]
        assert row["length"] in pre["fit"]["lengths"]
        assert len(row["public"])==4
        assert set(row["answers"])==set(row["correct"])==set(row["abstract_work"])==set(ARMS)
        preds,work=answer(row,temps)
        assert preds==row["answers"] and work==row["abstract_work"]
        assert row["truth"] in range(P)
        for arm,guess in preds.items():
            correct=int(guess==row["truth"])
            assert row["correct"][arm]==correct
            counts[arm]+=correct
    assert len(seen)==64
    return {"schema":"wingless.f57.sourcefree-FIT-predictor-audit.v1",
        "validation_count":64,"correct":counts,
        "source_fit_only":True,"frozen_controls_reconstructed":True,
        "prospective_DIAG_not_executed":True,
        "accepted_original":False,"policy_improvement_proven":False}
if __name__=="__main__":
    assert len(sys.argv)==3
    d=json.loads(pathlib.Path(sys.argv[1]).read_text())
    p=json.loads(pathlib.Path(sys.argv[2]).read_text())
    print(json.dumps(audit(d,p),sort_keys=True))
