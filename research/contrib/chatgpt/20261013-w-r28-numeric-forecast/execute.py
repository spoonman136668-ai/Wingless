"""W-R28 frozen causal numerical forecast benchmark; no future-byte access."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r27-second-order"/"execute.py"
spec=importlib.util.spec_from_file_location("w27_frozen",SOURCE)
assert spec and spec.loader
w27=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w27)
w26=w27.w26

DEVELOPMENT=tuple(range(448,456))
HELDOUT=tuple(range(456,472))
GENERATORS=("affine_prior","quadratic_irregular","xorshift_independent")
MODELS=("persistence","velocity","acceleration")
AMPS=(1,3,5)


def predict(values,t,model):
    a=values[t]
    b=values[t-1]
    c=values[t-2]
    if model=="persistence":return a
    if model=="velocity":return (2*a-b)%256
    assert model=="acceleration"
    return (3*a-3*b+c)%256


def distance(a,b):
    diff=(a-b)%256
    return min(diff,256-diff)


def evaluate_seed(seed,generator,amp,model):
    error=0
    n=0
    exact=0
    within8=0
    for entity in range(9):
        _,observed=w26.signature(seed,entity,generator,amp,"heldout")
        for channel in range(4):
            values=[row[channel] for row in observed]
            for t in range(2,15):
                # The target at t+1 enters evaluation only.
                estimate=predict(values,t,model)
                target=values[t+1]
                e=distance(estimate,target)
                error+=e
                exact+=(e==0)
                within8+=(e<=8)
                n+=1
    assert n==9*4*13
    return {"seed":seed,"abs_circular_error":error,"samples":n,
            "exact":exact,"within8":within8}


def evaluate_group(seeds,generator,amp,model):
    rows=[evaluate_seed(seed,generator,amp,model) for seed in seeds]
    n=sum(r["samples"] for r in rows)
    return {"generator":generator,"amplitude":amp,"predictor":model,
            "samples":n,"circular_mae":sum(r["abs_circular_error"] for r in rows)/n,
            "exact_accuracy":sum(r["exact"] for r in rows)/n,
            "within8_accuracy":sum(r["within8"] for r in rows)/n,
            "by_seed":rows}


def select_from_development():
    # A global model choice determined only by development affine+quadratic streams.
    losses={}
    for model in MODELS:
        cells=[evaluate_group(DEVELOPMENT,gen,amp,model)
               for gen in ("affine_prior","quadratic_irregular") for amp in (1,3)]
        losses[model]=sum(c["circular_mae"]*c["samples"] for c in cells)/sum(c["samples"] for c in cells)
    # Frozen tie-break matches MODELS order.
    winner=min(MODELS,key=lambda x:(losses[x],MODELS.index(x)))
    return winner,losses


def run():
    winner,dev_losses=select_from_development()
    cells=[evaluate_group(HELDOUT,g,amp,model)
           for g in GENERATORS for amp in AMPS for model in MODELS]
    assert len(cells)==27
    def err(generator,model):
        row=next(c for c in cells if c["generator"]==generator and
                 c["amplitude"]==3 and c["predictor"]==model)
        return row["circular_mae"]
    affine_gain=err("affine_prior","persistence")-err("affine_prior","velocity")
    quadratic_gain=err("quadratic_irregular","velocity")-err("quadratic_irregular","acceleration")
    label=("STRUCTURED_FORECAST" if affine_gain>=8 and quadratic_gain>=8 else
           "PARTIAL" if max(affine_gain,quadratic_gain)>=8 else "NEGATIVE")
    return {"schema":"wingless.w-r28-numeric-next-byte.v1",
            "parent_head":"859d4af60957941fa3c28dd0c41841c42296d5df",
            "classification":label,"affine_velocity_gain":affine_gain,
            "quadratic_acceleration_gain":quadratic_gain,
            "development_only_selected_model":winner,
            "development_model_losses":dev_losses,
            "development_seeds":DEVELOPMENT,"sealed_heldout_seeds":HELDOUT,
            "heldout_cells":cells,"model_calls":0,
            "claim_limit":"Fixed deterministic modular recurrence on synthetic raw streams"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
