"""W-R29: frozen quadratic term ablation for causal modular raw forecasting."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r28b-numeric"/"execute.py"
spec=importlib.util.spec_from_file_location("w28b_frozen",SOURCE)
assert spec and spec.loader
w28=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w28)
w20=w28.w26.w20

DEVELOPMENT=tuple(range(496,504))
HELDOUT=tuple(range(504,520))
GENERATORS=("affine_prior","pure_quadratic","quadratic_irregular","xorshift_independent")
AMPS=(1,3,5)
MODELS=("persistence","velocity","acceleration")


def sequence(seed,entity,generator):
    if generator=="pure_quadratic":
        return tuple(
            tuple((43*seed+29*entity+17*c+(5+2*c)*t+(11+3*c)*t*t)%256
                  for c in range(4))
            for t in range(16))
    return w20.base_sequence(seed,entity,generator)


def observed(seed,entity,generator,amp):
    base=sequence(seed,entity,generator)
    rows=tuple(
        tuple((val+79+
               ((11*t+7*c+3*entity+seed)%(2*amp+1))-amp)%256
              for c,val in enumerate(row))
        for t,row in enumerate(base))
    if generator!="pure_quadratic":
        assert rows==w28.observed_for(seed,entity,generator,amp)
    if generator=="quadratic_irregular":
        pure=sequence(seed,entity,"pure_quadratic")
        for t in range(16):
            for c in range(4):
                assert base[t][c]==(pure[t][c]+7*(((seed+c)*t)%19))%256
    return rows


def seed_score(seed,generator,amp,model):
    error=0
    within8=0
    exact=0
    n=0
    for entity in range(9):
        o=observed(seed,entity,generator,amp)
        for c in range(4):
            x=[row[c] for row in o]
            for t in range(2,15):
                prediction=w28.predict(x,t,model)
                hidden_target=x[t+1]
                diff=w28.distance(prediction,hidden_target)
                error+=diff
                within8+=(diff<=8)
                exact+=(diff==0)
                n+=1
    assert n==9*4*13
    return {"seed":seed,"absolute_error":error,"n":n,
            "within8":within8,"exact":exact}


def cell(generator,amp,model):
    rows=[seed_score(seed,generator,amp,model) for seed in HELDOUT]
    count=sum(r["n"] for r in rows)
    return {"generator":generator,"amplitude":amp,"model":model,
            "circular_mae":sum(r["absolute_error"] for r in rows)/count,
            "within8_accuracy":sum(r["within8"] for r in rows)/count,
            "exact_accuracy":sum(r["exact"] for r in rows)/count,
            "heldout_rows":rows}


def run():
    cells=[cell(g,amp,m) for g in GENERATORS for amp in AMPS for m in MODELS]
    assert len(cells)==36
    d={(c["generator"],c["amplitude"],c["model"]):c["circular_mae"] for c in cells}
    pure=d[("pure_quadratic",3,"acceleration")]
    irregular=d[("quadratic_irregular",3,"acceleration")]
    affine=d[("affine_prior",3,"velocity")]
    delta=irregular-pure
    classification=("IRREGULARITY_BOUNDARY" if pure<=20 and delta>=15 and affine<=10 else
                    "PARTIAL" if delta>=10 else "NEGATIVE")
    return {"schema":"wingless.w-r29-irregularity-ablated.v1",
            "parent_head":"785ec7782e955225c3c140d181b261e9310eb042",
            "classification":classification,
            "pure_acceleration_error":pure,
            "irregular_acceleration_error":irregular,
            "irregularity_error_increase":delta,
            "affine_velocity_error":affine,
            "development_seeds":DEVELOPMENT,
            "sealed_heldout_seeds":HELDOUT,
            "cells":cells,"model_calls":0,
            "claim_limit":"Fixed deterministic synthetic generator-factor ablation"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
