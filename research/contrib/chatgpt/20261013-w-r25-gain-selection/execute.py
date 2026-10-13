"""W-R25: frozen development-only selection of causal forecast residual gain."""
import importlib.util
import json
from pathlib import Path

PARENT=Path(__file__).resolve().parent.parent/"20261013-w-r24-feedback"/"execute.py"
spec=importlib.util.spec_from_file_location("w24_frozen",PARENT)
assert spec and spec.loader
w24=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w24)
w22=w24.w22

GAINS=(0.0,0.125,0.25,0.5)


def series(seed,entity,stage):
    s=(5*entity+seed)%9
    if stage=="development_train":
        shift=0
        modulus=3
        coefficient=5
    elif stage=="development_validation":
        shift=2
        modulus=3
        coefficient=5
    else:
        assert stage=="heldout"
        shift=3
        modulus=5
        coefficient=7
    noise=[(coefficient*t+seed+entity)%modulus-modulus//2 for t in range(16)]
    assert len(set(noise))==modulus
    return tuple((
        192+seed+7*s+3*t+t*t+
        shift*max(0,t-8)**2+noise[t]
    )%256 for t in range(16))


def learn_curvature():
    diffs=[]
    for seed in range(104,106):
        for entity in range(9):
            readings=series(seed,entity,"development_train")
            for t in range(2,16):
                diffs.append(w22.circular_signed(readings[t]-2*readings[t-1]+readings[t-2]))
    assert len(diffs)==252
    diffs.sort()
    n=len(diffs)
    return (diffs[n//2-1]+diffs[n//2])/2


def median(values):
    ordered=sorted(values)
    n=len(ordered)
    assert 1<=n<=3
    return ordered[n//2] if n%2 else (ordered[n//2-1]+ordered[n//2])/2


def prediction(readings,t,curvature,gain):
    baseline=2*readings[t-1]-readings[t-2]+curvature
    residuals=[]
    for j in range(max(2,t-3),t):
        preceding_forecast=2*readings[j-1]-readings[j-2]+curvature
        residuals.append(w22.circular_signed(readings[j]-preceding_forecast))
    correction=median(residuals) if residuals else 0.0
    return baseline+gain*correction


def metric(seeds,stage,curvature,gain):
    out=[]
    early=[]
    late=[]
    for seed in seeds:
        errors=[]
        for entity in range(9):
            readings=series(seed,entity,stage)
            for t in range(2,16):
                forecast=prediction(readings,t,curvature,gain)
                errors.append((t,w22.circular_abs(forecast,readings[t])))
        assert len(errors)==126
        pre=[e for t,e in errors if t<=8]
        post=[e for t,e in errors if t>=9]
        assert len(pre)==len(post)==63
        out.append({"seed":seed,"prechange_mae":sum(pre)/len(pre),
                    "postchange_mae":sum(post)/len(post),
                    "full_mae":sum(e for t,e in errors)/len(errors)})
        early.extend(pre)
        late.extend(post)
    return {"gain":gain,"prechange_mae":sum(early)/len(early),
            "postchange_mae":sum(late)/len(late),
            "full_mae":(sum(early)+sum(late))/(len(early)+len(late)),
            "by_seed":out}


def run():
    curvature=learn_curvature()
    validations=[metric(range(106,108),"development_validation",curvature,gain)
                 for gain in GAINS]
    # Precommitted tie-breaker prefers smaller gain, no heldout information.
    choice=min(validations,key=lambda row:(row["postchange_mae"],row["gain"]))["gain"]
    controls=(0.0,0.25,choice)
    held=[metric(range(108,112),"heldout",curvature,gain)
          for gain in controls]
    frozen,fixed,selected=held
    improved=selected["postchange_mae"]<frozen["postchange_mae"]
    if (selected["postchange_mae"]<=.9*frozen["postchange_mae"]
        and selected["prechange_mae"]<=1.1*frozen["prechange_mae"]
        and selected["postchange_mae"]<fixed["postchange_mae"]):
        label="DEV_GAIN_TRANSFER"
    elif improved:
        label="PARTIAL"
    else:
        label="NEGATIVE"
    return {
        "schema":"wingless.w-r25-train-only-feedback-gain.v1",
        "predecessor_head":"3c02f2afc1868c97659e17a667c8e8d88277c9fa",
        "classification":label,
        "development_curvature":curvature,
        "development_gain_validation":validations,
        "chosen_development_gain":choice,
        "heldout_arms":[{"arm":"static","metrics":frozen},
                        {"arm":"fixed_quarter","metrics":fixed},
                        {"arm":"dev_selected","metrics":selected}],
        "model_calls":0,
        "scope":"Synthetic supervised development gain selection; sealed heldout shift amplitude3"
    }


if __name__=="__main__":
    outcome=run()
    assert outcome==run()
    print(json.dumps(outcome,sort_keys=True,indent=2))
