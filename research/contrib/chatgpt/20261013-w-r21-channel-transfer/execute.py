"""W-R21: heldout channel-permutation recognition with frozen one-scalar learning."""
import importlib.util
import json
from pathlib import Path

PARENT=Path(__file__).resolve().parent.parent/"20261013-w-r20-temporal-warp"/"execute.py"
spec=importlib.util.spec_from_file_location("wr20_frozen",PARENT)
assert spec and spec.loader
w20=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w20)
w18=w20.prior

ARMS=("unordered_ch0","transition_ch0","transition_max4","transition_mean4")


def example(seed,entity,positive,phase):
    dev=phase=="development"
    base=w18.original(seed,entity)
    modulus=3 if dev else 5
    mul=5 if dev else 7
    global_offset=17 if dev else 61
    noise=[(mul*t+entity+seed)%modulus-modulus//2 for t in range(16)]
    assert len(set(noise))==modulus
    shifted=tuple(tuple((v+global_offset+noise[t])%256 for v in row)
                  for t,row in enumerate(base))
    if positive:
        observed=shifted
    elif dev:
        tail=shifted[1:]
        observed=shifted[:1]+(tail[::-1] if entity%2==0 else tail[3:]+tail[:3])
    else:
        rows=list(shifted)
        rows[5],rows[6]=rows[6],rows[5]
        rows[10],rows[11]=rows[11],rows[10]
        observed=tuple(rows)
    assert observed[0]==shifted[0]
    assert sorted(observed)==sorted(shifted)
    assert (observed==shifted)==positive
    if not dev:
        permutation=(1,2,3,0)
        base=tuple(tuple(row[c] for c in permutation) for row in base)
        observed=tuple(tuple(row[c] for c in permutation) for row in observed)
    return base,observed


def edges(a,b,channel):
    x=[(row[channel]-a[0][channel])%256 for row in a]
    y=[(row[channel]-b[0][channel])%256 for row in b]
    return sum(
        w18.circular_distance((x[t+1]-x[t])%256,(y[t+1]-y[t])%256)
        for t in range(15)
    )/15


def score(a,b,arm):
    if arm=="unordered_ch0":
        return w18.score(a,b,"unordered_histogram")
    if arm=="transition_ch0":
        return edges(a,b,0)
    values=[edges(a,b,c) for c in range(4)]
    if arm=="transition_max4":return max(values)
    assert arm=="transition_mean4"
    return sum(values)/4


def train(arm):
    pos=[]
    neg=[]
    for seed in range(72,76):
        for entity in range(9):
            for positive in (True,False):
                a,b=example(seed,entity,positive,"development")
                (pos if positive else neg).append(score(a,b,arm))
    assert len(pos)==len(neg)==36
    maxp=max(pos)
    minn=min(neg)
    return {"threshold":(maxp+minn)/2,"max_positive":maxp,"min_negative":minn,
            "separable":maxp<minn}


def assess(arm,fit):
    rows=[]
    for seed in range(76,80):
        tp=tn=fp=fn=0
        for entity in range(9):
            for positive in (True,False):
                a,b=example(seed,entity,positive,"heldout")
                prediction=score(a,b,arm)<=fit["threshold"]
                if positive and prediction:tp+=1
                elif not positive and not prediction:tn+=1
                elif not positive:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        rows.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,
                     "accuracy":(tp+tn)/18})
    return {"arm":arm,"trained_scalars":1,"fit":fit,
            "accuracy":sum(r["tp"]+r["tn"] for r in rows)/72,
            "heldout":rows}


def run():
    results=[assess(arm,train(arm)) for arm in ARMS]
    scores={r["arm"]:r["accuracy"] for r in results}
    symmetric=max(scores["transition_max4"],scores["transition_mean4"])
    ch0=scores["transition_ch0"]
    if symmetric>=.85 and symmetric>ch0 and scores["unordered_ch0"]<=.55:
        label="ROLE_TRANSFER"
    elif symmetric>ch0:
        label="PARTIAL"
    else:
        label="NEGATIVE"
    return {
        "schema":"wingless.w-r21-channel-permutation.v1",
        "predecessor_head":"9fdd600c66b3dce4fefe31d3e64076cb436e823b",
        "classification":label,
        "heldout_channel_permutation":[1,2,3,0],
        "models":results,"model_calls":0,
        "scope":"Synthetic channel permutation invariance, not semantic role generalization"
    }


if __name__=="__main__":
    result=run()
    assert result==run()
    print(json.dumps(result,sort_keys=True,indent=2))
