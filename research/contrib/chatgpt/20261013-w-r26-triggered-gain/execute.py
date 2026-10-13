"""W-R26: train-only residual-trigger gate and gain choice for causal raw prediction."""
import importlib.util
import json
import math
from pathlib import Path

PARENT=Path(__file__).resolve().parent.parent/"20261013-w-r25-gain-selection"/"execute.py"
spec=importlib.util.spec_from_file_location("w25_frozen",PARENT)
assert spec and spec.loader
w25=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w25)
w22=w25.w22

GAINS=(0.0,0.125,0.25,0.5)


def curvature():
    ds=[]
    for seed in range(112,114):
        for entity in range(9):
            seq=w25.series(seed,entity,"development_train")
            for t in range(2,16):
                ds.append(w22.circular_signed(seq[t]-2*seq[t-1]+seq[t-2]))
    ds.sort()
    assert len(ds)==252
    n=len(ds)
    return (ds[n//2-1]+ds[n//2])/2


def past_residual(seq,j,curv):
    assert j>=2
    prediction=2*seq[j-1]-seq[j-2]+curv
    return w22.circular_signed(seq[j]-prediction)


def median(values):
    seq=sorted(values)
    n=len(seq)
    assert 1<=n<=3
    return seq[n//2] if n%2 else (seq[n//2-1]+seq[n//2])/2


def learn_gate(curv):
    residual_abs=[]
    for seed in range(112,114):
        for entity in range(9):
            seq=w25.series(seed,entity,"development_train")
            for t in range(2,16):
                residual_abs.append(abs(past_residual(seq,t,curv)))
    residual_abs.sort()
    assert len(residual_abs)==252
    index=math.ceil(.95*len(residual_abs))-1
    return residual_abs[index],{"sample_count":len(residual_abs),
                                "percentile_index":index,
                                "train_max_residual":max(residual_abs)}


def step(seq,t,curv,gain,gate,mode):
    base=2*seq[t-1]-seq[t-2]+curv
    residuals=[past_residual(seq,j,curv) for j in range(max(2,t-3),t)]
    median_r=median(residuals) if residuals else 0.0
    active=bool(residuals) and (mode=="always" or abs(median_r)>gate)
    return base+(gain*median_r if active else 0.0),active


def evaluate(seeds,stage,curv,gain,gate,mode):
    per_seed=[]
    pre=[]
    post=[]
    active_pre=active_post=0
    for seed in seeds:
        errors=[]
        for entity in range(9):
            seq=w25.series(seed,entity,stage)
            for t in range(2,16):
                predicted,activated=step(seq,t,curv,gain,gate,mode)
                err=w22.circular_abs(predicted,seq[t])
                errors.append((t,err))
                if activated and t<=8:active_pre+=1
                if activated and t>=9:active_post+=1
        assert len(errors)==126
        early=[err for t,err in errors if t<=8]
        late=[err for t,err in errors if t>=9]
        assert len(early)==len(late)==63
        per_seed.append({"seed":seed,"prechange_mae":sum(early)/len(early),
                         "postchange_mae":sum(late)/len(late)})
        pre.extend(early)
        post.extend(late)
    return {"gain":gain,"mode":mode,
            "prechange_mae":sum(pre)/len(pre),
            "postchange_mae":sum(post)/len(post),
            "full_mae":(sum(pre)+sum(post))/(len(pre)+len(post)),
            "triggered_pre":active_pre,"triggered_post":active_post,
            "by_seed":per_seed}


def run():
    curv=curvature()
    gate,provenance=learn_gate(curv)
    dev_results=[evaluate(range(114,116),"development_validation",curv,gain,gate,"trigger")
                 for gain in GAINS]
    chosen=min(dev_results,key=lambda r:(r["postchange_mae"],r["gain"]))["gain"]
    data=[
        evaluate(range(116,120),"heldout",curv,0.0,gate,"always"),
        evaluate(range(116,120),"heldout",curv,.25,gate,"always"),
        evaluate(range(116,120),"heldout",curv,chosen,gate,"always"),
        evaluate(range(116,120),"heldout",curv,chosen,gate,"trigger"),
    ]
    static,quarter,always,trigger=data
    if (trigger["postchange_mae"]<=.9*static["postchange_mae"]
        and trigger["prechange_mae"]<=1.1*static["prechange_mae"]
        and trigger["postchange_mae"]<always["postchange_mae"]):
        classification="FEEDBACK_TRIGGER_TRANSFER"
    elif trigger["postchange_mae"]<static["postchange_mae"]:
        classification="PARTIAL"
    else:
        classification="NEGATIVE"
    return {
        "schema":"wingless.w-r26-causal-feedback-trigger.v1",
        "predecessor_head":"3eae59154cdc2dd88b3e8e82f54a052098c35358",
        "classification":classification,
        "train_curvature":curv,
        "dev_residual_gate":gate,
        "gate_provenance":provenance,
        "selected_dev_gain":chosen,
        "development_gain_validation":dev_results,
        "heldout_models":[{"arm":name,"metrics":out} for name,out in zip(
            ("static","fixed_quarter","selected_always","selected_triggered"),data)],
        "model_calls":0,
        "scope":"Synthetic development-trained causal change trigger, no future target access"
    }


if __name__=="__main__":
    result=run()
    assert result==run()
    print(json.dumps(result,sort_keys=True,indent=2))
