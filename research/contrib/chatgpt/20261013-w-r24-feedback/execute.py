"""W-R24: causal fixed-gain forecast-error feedback without future target access."""
import importlib.util
import json
from pathlib import Path

PARENT=Path(__file__).resolve().parent.parent/"20261013-w-r23-robust-curvature"/"execute.py"
spec=importlib.util.spec_from_file_location("w23_frozen",PARENT)
assert spec and spec.loader
w23=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w23)
w22=w23.w22

ARMS=("frozen","naive_curvature","last_error_quarter","median3_error_quarter")


def train():
    diffs=[]
    for seed in range(96,100):
        for entity in range(9):
            obs=w22.series(seed,entity,"development")
            for t in range(2,16):
                diffs.append(w22.circular_signed(obs[t]-2*obs[t-1]+obs[t-2]))
    assert len(diffs)==504
    diffs.sort()
    n=len(diffs)
    return (diffs[n//2-1]+diffs[n//2])/2


def median(xs):
    v=sorted(xs)
    assert 1<=len(v)<=3
    n=len(v)
    return v[n//2] if n%2 else (v[n//2-1]+v[n//2])/2


def residual(obs,t,curvature):
    assert 2<=t<16
    prior_prediction=2*obs[t-1]-obs[t-2]+curvature
    # A causal evaluator can calculate this only after observing tick t.
    return w22.circular_signed(obs[t]-prior_prediction)


def predict(obs,t,arm,curvature):
    baseline=2*obs[t-1]-obs[t-2]+curvature
    if arm=="frozen":return baseline
    if arm=="naive_curvature":
        local=curvature if t<3 else w22.circular_signed(
            obs[t-1]-2*obs[t-2]+obs[t-3])
        return 2*obs[t-1]-obs[t-2]+local
    previous=[residual(obs,j,curvature) for j in range(max(2,t-3),t)]
    if arm=="last_error_quarter":
        correction=previous[-1] if previous else 0.0
    else:
        assert arm=="median3_error_quarter"
        correction=median(previous) if previous else 0.0
    return baseline+0.25*correction


def evaluate(arm,curvature):
    results=[]
    all_errors=[]
    pre=[]
    post=[]
    for seed in range(100,104):
        errors=[]
        for entity in range(9):
            obs=w22.series(seed,entity,"heldout")
            for t in range(2,16):
                predicted=predict(obs,t,arm,curvature)
                err=w22.circular_abs(predicted,obs[t])
                errors.append((t,err))
        assert len(errors)==126
        early=[e for t,e in errors if t<=8]
        late=[e for t,e in errors if t>=9]
        assert len(early)==len(late)==63
        results.append({
            "seed":seed,
            "full_mae":sum(e for t,e in errors)/len(errors),
            "prechange_mae":sum(early)/len(early),
            "postchange_mae":sum(late)/len(late)
        })
        all_errors.extend(e for t,e in errors)
        pre.extend(early)
        post.extend(late)
    return {
        "arm":arm,
        "n":len(all_errors),
        "full_mae":sum(all_errors)/len(all_errors),
        "prechange_mae":sum(pre)/len(pre),
        "postchange_mae":sum(post)/len(post),
        "by_seed":results
    }


def run():
    trained=train()
    rows=[evaluate(arm,trained) for arm in ARMS]
    scores={row["arm"]:row for row in rows}
    static=scores["frozen"]
    feedback=scores["median3_error_quarter"]
    count=sum(a["postchange_mae"]<b["postchange_mae"]
              for a,b in zip(feedback["by_seed"],static["by_seed"]))
    if (feedback["postchange_mae"]<=.9*static["postchange_mae"]
            and feedback["prechange_mae"]<=1.1*static["prechange_mae"]
            and count>=3):
        classification="FEEDBACK_TRANSFER"
    elif feedback["postchange_mae"]<static["postchange_mae"]:
        classification="PARTIAL"
    else:
        classification="NEGATIVE"
    return {
        "schema":"wingless.w-r24-causal-feedback.v1",
        "predecessor_head":"4646dbe94ae96eec5a3b30def36285f65fd3376d",
        "classification":classification,
        "development_curvature":trained,
        "feedback_gain":0.25,
        "heldout_improved_seeds":count,
        "models":rows,"model_calls":0,
        "scope":"Finite causal autoregression with observed-only error feedback"
    }


if __name__=="__main__":
    result=run()
    assert result==run()
    print(json.dumps(result,sort_keys=True,indent=2))
