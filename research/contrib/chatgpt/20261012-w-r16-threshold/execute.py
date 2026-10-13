"""W-R16: heldout-shielded train-only scalar separation of jittered pairs."""
import importlib.util
import json
from pathlib import Path

PRIOR_SOURCE = Path(__file__).resolve().parent.parent / "20261012-w-r15-noisy-positive" / "execute.py"
_spec = importlib.util.spec_from_file_location("w15_frozen", PRIOR_SOURCE)
assert _spec and _spec.loader
w15 = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(w15)


def example(seed, entity, positive, phase, kind):
    a, b = w15.candidate_pair(seed, entity, positive, phase)
    vec = w15.features(a,b,"unordered" if kind=="unordered" else "pooled",1.0)
    assert len(vec)==16
    return vec, 1 if positive else -1


def threshold_from_development():
    pos, neg = [],[]
    for seed in range(32,36):
        for entity in range(9):
            for positive in (True,False):
                vec,label=example(seed,entity,positive,"development","pooled")
                (pos if positive else neg).append(vec[0])
    assert len(pos)==len(neg)==36
    high_pos=max(pos)
    low_neg=min(neg)
    return (high_pos+low_neg)/2.0, high_pos, low_neg


def train_perceptron(kind):
    cases=[
        example(seed,entity,positive,"development",kind)
        for seed in range(32,36)
        for entity in range(9)
        for positive in (True,False)
    ]
    assert len(cases)==72
    weights=[0.0]*16
    bias=0.0
    updates=0
    for _ in range(32):
        for vec,label in cases:
            s=bias+sum(a*b for a,b in zip(weights,vec))
            if label*s<=0:
                weights=[a+0.1*label*b for a,b in zip(weights,vec)]
                bias+=0.1*label
                updates+=1
    return weights,bias,updates


def evaluate(kind,model=None,threshold=None):
    rows=[]
    for seed in range(36,40):
        tp=tn=fp=fn=0
        for entity in range(9):
            for positive in (True,False):
                vec,label=example(seed,entity,positive,"heldout",kind)
                if kind=="scalar_threshold":
                    pred=1 if vec[0]<=threshold else -1
                else:
                    weights,bias,_=model
                    pred=1 if bias+sum(a*b for a,b in zip(weights,vec))>=0 else -1
                if pred==1 and label==1:tp+=1
                elif pred==-1 and label==-1:tn+=1
                elif pred==1:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        rows.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,
                     "accuracy":(tp+tn)/18})
    return {
        "arm":kind,"model_capacity":1 if kind=="scalar_threshold" else 16,
        "training_updates":0 if kind=="scalar_threshold" else model[2],
        "heldout_accuracy":sum(r["tp"]+r["tn"] for r in rows)/72,
        "rows":rows
    }


def run():
    theta,maximum_pos,minimum_neg=threshold_from_development()
    arms=[
        evaluate("unordered",train_perceptron("unordered")),
        evaluate("pooled",train_perceptron("pooled")),
        evaluate("scalar_threshold",threshold=theta),
    ]
    scores={a["arm"]:a["heldout_accuracy"] for a in arms}
    if scores["scalar_threshold"]>=.85 and scores["unordered"]<=.55:
        label="THRESHOLD_TRANSFER"
    elif scores["scalar_threshold"]>max(scores["unordered"],scores["pooled"]):
        label="PARTIAL"
    else:
        label="NEGATIVE"
    return {
        "schema":"wingless.w-r16-train-only-threshold.v1",
        "predecessor_head":"cd9ed8c98e655e9ebf77204750f9ca733ced83e5",
        "classification":label,
        "train_only_threshold":theta,
        "train_max_positive":maximum_pos,
        "train_min_negative":minimum_neg,
        "train_separable":maximum_pos<minimum_neg,
        "arms":arms,"model_calls":0,
        "claim_limit":"Train-only threshold on synthetic jittered pairs; not general intelligence"
    }


if __name__=="__main__":
    x=run()
    assert x==run()
    print(json.dumps(x,sort_keys=True,indent=2))
