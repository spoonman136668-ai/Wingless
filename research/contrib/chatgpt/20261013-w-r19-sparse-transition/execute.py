"""W-R19: preregistered sparse-edge statistics on unseen double local swaps."""
import importlib.util
import json
from pathlib import Path

PARENT=Path(__file__).resolve().parent.parent/"20261013-w-r18-local-order"/"execute.py"
spec=importlib.util.spec_from_file_location("wr18_frozen",PARENT)
assert spec and spec.loader
prior=importlib.util.module_from_spec(spec)
spec.loader.exec_module(prior)

ARMS=("unordered_histogram","mean_transition","max_transition","top3_transition")


def pair(seed,entity,positive,phase):
    if phase=="development":
        return prior.make_pair(seed,entity,positive,phase)
    base=prior.original(seed,entity)
    modulus=5
    noise=[((7*t+entity+seed)%modulus)-2 for t in range(16)]
    assert len(set(noise))==5
    candidate=tuple(tuple((value+61+noise[t])%256 for value in row)
                    for t,row in enumerate(base))
    observed=list(candidate)
    if not positive:
        observed[5],observed[6]=observed[6],observed[5]
        observed[10],observed[11]=observed[11],observed[10]
    observed=tuple(observed)
    assert observed[0]==candidate[0]
    assert sorted(candidate)==sorted(observed)
    assert (observed==candidate)==positive
    return base,observed


def score(a,b,arm):
    if arm=="unordered_histogram":
        return prior.score(a,b,"unordered_histogram")
    x=prior.anchor_relative(a)
    y=prior.anchor_relative(b)
    deltas=[
        prior.circular_distance((x[i+1]-x[i])%256,(y[i+1]-y[i])%256)
        for i in range(15)
    ]
    if arm=="mean_transition":return sum(deltas)/15
    if arm=="max_transition":return max(deltas)
    assert arm=="top3_transition"
    return sum(sorted(deltas,reverse=True)[:3])/3


def training(arm):
    yes,no=[],[]
    for seed in range(56,60):
        for entity in range(9):
            for positive in (True,False):
                a,b=pair(seed,entity,positive,"development")
                (yes if positive else no).append(score(a,b,arm))
    assert len(yes)==len(no)==36
    highest_positive=max(yes)
    lowest_negative=min(no)
    return {
        "threshold":(highest_positive+lowest_negative)/2,
        "highest_development_positive":highest_positive,
        "lowest_development_negative":lowest_negative,
        "development_separable":highest_positive<lowest_negative
    }


def test(arm,fit):
    rows=[]
    for seed in range(60,64):
        tp=tn=fp=fn=0
        for entity in range(9):
            for positive in (True,False):
                a,b=pair(seed,entity,positive,"heldout")
                prediction=score(a,b,arm)<=fit["threshold"]
                if positive and prediction:tp+=1
                elif not positive and not prediction:tn+=1
                elif not positive:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        rows.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,
                     "accuracy":(tp+tn)/18})
    return {
        "arm":arm,"trained_parameters":1,"fit":fit,
        "heldout_accuracy":sum(r["tp"]+r["tn"] for r in rows)/72,
        "heldout":rows
    }


def run():
    arms=[test(arm,training(arm)) for arm in ARMS]
    s={r["arm"]:r["heldout_accuracy"] for r in arms}
    best=max(s["max_transition"],s["top3_transition"])
    if best>=.85 and s["unordered_histogram"]<=.55 and best>s["mean_transition"]:
        label="SPARSE_EDGE_TRANSFER"
    elif best>s["mean_transition"]:
        label="PARTIAL"
    else:
        label="NEGATIVE"
    return {
        "schema":"wingless.w-r19-sparse-edge.v1",
        "predecessor_head":"cc98b0829b1f238fe6f950985b8f3420bd9ef8f5",
        "classification":label,"arms":arms,"model_calls":0,
        "scope":"New unseen double-adjacent swap, fixed nonlinear toy sequence"
    }


if __name__=="__main__":
    result=run()
    assert result==run()
    print(json.dumps(result,sort_keys=True,indent=2))
