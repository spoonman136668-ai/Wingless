"""W-R17: actual varying train tick jitter with frozen fresh cohort."""
import importlib.util
import json
from pathlib import Path

PRIOR_SOURCE=Path(__file__).resolve().parent.parent/"20261012-w-r13-general-permutation"/"execute.py"
_spec=importlib.util.spec_from_file_location("w13_frozen",PRIOR_SOURCE)
assert _spec and _spec.loader
w13=importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(w13)
w12=w13._prior


def generate_pair(seed, entity, positive, phase):
    base=w12.original(seed,entity)
    development=phase=="development"
    modulus=3 if development else 5
    scale=17 if development else 61
    jitter=tuple(((5 if development else 7)*t+entity+seed)%modulus-(modulus//2)
                 for t in range(16))
    assert len(set(jitter))==modulus, "variable per-tick jitter invariant"
    shifted=w12.nuisance_shift(base,scale)
    observed=tuple(
        tuple((value+jitter[t])%256 for value in row)
        for t,row in enumerate(shifted)
    )
    if positive:
        candidate=observed
    else:
        tail=observed[1:]
        if development:
            candidate=observed[:1]+(tail[::-1] if entity%2==0 else tail[3:]+tail[:3])
        else:
            candidate=observed[:1]+tail[::2]+tail[1::2]
    assert candidate[0]==observed[0]
    assert sorted(candidate)==sorted(observed)
    assert (candidate==observed)==positive
    return base,candidate


def feature(a,b,kind):
    mode="unordered_histogram" if kind=="unordered" else "pooled"
    values=w13.pair_features(a,b,mode)
    assert len(values)==16
    return values


def example(seed,entity,positive,phase,kind):
    a,b=generate_pair(seed,entity,positive,phase)
    return feature(a,b,kind), 1 if positive else -1


def dev_threshold():
    positives,negatives=[],[]
    for seed in range(40,44):
        for entity in range(9):
            for pos in (True,False):
                vec,label=example(seed,entity,pos,"development","pooled")
                (positives if pos else negatives).append(vec[0])
    assert len(positives)==len(negatives)==36
    maxp=max(positives)
    minn=min(negatives)
    return (maxp+minn)/2,maxp,minn


def train(kind):
    cases=[
        example(seed,entity,pos,"development",kind)
        for seed in range(40,44)
        for entity in range(9)
        for pos in (True,False)
    ]
    assert len(cases)==72
    weights=[0.0]*16
    bias=0.0
    updates=0
    for _ in range(32):
        for vec,label in cases:
            score=bias+sum(w*x for w,x in zip(weights,vec))
            if label*score<=0:
                weights=[w+.1*label*x for w,x in zip(weights,vec)]
                bias+=.1*label
                updates+=1
    return weights,bias,updates


def assess(kind,model=None,threshold=None):
    heldout=[]
    for seed in range(44,48):
        tp=tn=fp=fn=0
        for entity in range(9):
            for positive in (True,False):
                vec,label=example(seed,entity,positive,"heldout",kind)
                if kind=="scalar_threshold":
                    pred=1 if vec[0]<=threshold else -1
                else:
                    w,b,_=model
                    pred=1 if b+sum(x*y for x,y in zip(w,vec))>=0 else -1
                if pred==1 and label==1:tp+=1
                elif pred==-1 and label==-1:tn+=1
                elif pred==1:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        heldout.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,
                        "accuracy":(tp+tn)/18})
    return {
        "kind":kind,"learned_parameters":1 if kind=="scalar_threshold" else 17,
        "training_updates":0 if kind=="scalar_threshold" else model[2],
        "heldout_accuracy":sum(a["tp"]+a["tn"] for a in heldout)/72,
        "heldout":heldout
    }


def run():
    theta,maxp,minn=dev_threshold()
    arms=[
        assess("unordered",train("unordered")),
        assess("pooled",train("pooled")),
        assess("scalar_threshold",threshold=theta)
    ]
    scores={a["kind"]:a["heldout_accuracy"] for a in arms}
    if scores["scalar_threshold"]>=.85 and scores["unordered"]<=.55:
        classification="ROBUST_TRANSFER"
    elif scores["scalar_threshold"]>max(scores["unordered"],scores["pooled"]):
        classification="PARTIAL"
    else:
        classification="NEGATIVE"
    return {
        "schema":"wingless.w-r17-genuine-tick-jitter.v1",
        "predecessor_head":"8e8d402ae3e1134adcfd3fcf948eb94a46dfe135",
        "classification":classification,
        "train_only_threshold":theta,
        "train_max_positive":maxp,
        "train_min_negative":minn,
        "development_jitter_level_count":3,
        "heldout_jitter_level_count":5,
        "arms":arms,"model_calls":0,
        "claim_limit":"Controlled synthetic jitter and temporal matching only"
    }


if __name__=="__main__":
    data=run()
    assert data==run()
    print(json.dumps(data,sort_keys=True,indent=2))
