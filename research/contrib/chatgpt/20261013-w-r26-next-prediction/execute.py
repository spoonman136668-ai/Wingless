"""W-R26: first direct future raw-observation delta prediction, causal features."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r20-generator-shift"/"execute.py"
spec=importlib.util.spec_from_file_location("w20_frozen",SOURCE)
assert spec and spec.loader
w20=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w20)
w16=w20.w19.w16

DEV=tuple(range(384,400))
HELDOUT=tuple(range(400,416))
GENS_TRAIN=("affine_prior","quadratic_irregular")
GENS_TEST=("affine_prior","quadratic_irregular","xorshift_independent")
ARMS=("ordered_lag_deltas","unordered_histogram","current_four_channels_repeated")


def signature(seed,entity,generator,amp,phase):
    base=w20.base_sequence(seed,entity,generator)
    observed=w16.observed(base,seed,entity,amp,
                          17 if phase=="development" else 79)
    assert len(base)==len(observed)==16
    return base,observed


def label_at(base,t):
    assert 4<=t<=14
    future=(base[t+1][0]-base[t][0])%256
    return 1 if future<128 else -1


def featurize(observed,t,arm):
    assert 4<=t<=14
    window=observed[t-4:t+1]
    assert len(window)==5
    if arm=="ordered_lag_deltas":
        features=[]
        for channel in range(4):
            for index in range(4):
                difference=(window[index+1][channel]-
                            window[index][channel])%256
                signed=((difference+128)%256)-128
                features.append(signed/128.0)
    elif arm=="unordered_histogram":
        features=[0.0]*16
        for row in window:
            features[row[0]%16]+=1.0/5.0
    else:
        assert arm=="current_four_channels_repeated"
        features=[window[-1][channel]/255.0 for channel in range(4)]*4
    assert len(features)==16
    return features


def samples(seed,generator,amp,phase,arm):
    result=[]
    for entity in range(9):
        base,observed=signature(seed,entity,generator,amp,phase)
        for t in range(4,15):
            vec=featurize(observed,t,arm)
            label=label_at(base,t)
            result.append((vec,label))
    assert len(result)==9*11
    return result


def training_rows(kind):
    all_rows=[row
              for seed in DEV
              for generator in GENS_TRAIN
              for amp in (1,2)
              for row in samples(seed,generator,amp,"development",kind)]
    assert len(all_rows)==6336
    return all_rows


def train(kind):
    rows=training_rows(kind)
    weights=[0.0]*16
    bias=0.0
    updates=0
    for _ in range(32):
        for vector,label in rows:
            signed=label*(bias+sum(w*x for w,x in zip(weights,vector)))
            if signed<=0:
                weights=[w+.1*label*x for w,x in zip(weights,vector)]
                bias+=.1*label
                updates+=1
    return weights,bias,updates


def evaluate(kind,model,generator):
    weights,bias,_=model
    result=[]
    for seed in HELDOUT:
        tp=tn=fp=fn=0
        for vector,actual in samples(seed,generator,3,"heldout",kind):
            predicted=1 if bias+sum(w*x for w,x in zip(weights,vector))>=0 else -1
            if actual==1 and predicted==1:tp+=1
            elif actual==-1 and predicted==-1:tn+=1
            elif actual==-1:fp+=1
            else:fn+=1
        assert tp+tn+fp+fn==99
        result.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn})
    correct=sum(r["tp"]+r["tn"] for r in result)
    total=sum(sum(r[k] for k in ("tp","tn","fp","fn")) for r in result)
    positive=sum(r["tp"]+r["fn"] for r in result)
    return {"generator":generator,"correct":correct,"total":total,
            "accuracy":correct/total,"true_positive_rate":positive/total,
            "by_seed":result}


def run():
    train_labels=[actual for _,actual in training_rows("ordered_lag_deltas")]
    npositive=sum(label==1 for label in train_labels)
    majority=1 if npositive>=len(train_labels)/2 else -1
    arms=[]
    for kind in ARMS:
        model=train(kind)
        tests=[evaluate(kind,model,gen) for gen in GENS_TEST]
        arms.append({"arm":kind,"width":16,"training_examples":6336,
                     "epochs":32,"updates":model[2],"tests":tests})
    scores={a["arm"]:{r["generator"]:r["accuracy"] for r in a["tests"]}
            for a in arms}
    ordered=scores["ordered_lag_deltas"]
    hist=scores["unordered_histogram"]
    xorshift_label_rate=next(r["true_positive_rate"] for r in arms[0]["tests"]
                             if r["generator"]=="xorshift_independent")
    baseline=xorshift_label_rate if majority==1 else 1-xorshift_label_rate
    if (ordered["quadratic_irregular"]>=.65 and
        ordered["xorshift_independent"]>=.65 and
        ordered["xorshift_independent"]-hist["xorshift_independent"]>=.05 and
        ordered["xorshift_independent"]>baseline):
        label="RAW_NEXT_STEP_TRANSFER"
    elif ordered["xorshift_independent"]-hist["xorshift_independent"]>=.03:
        label="PARTIAL"
    else:
        label="NEGATIVE"
    return {"schema":"wingless.w-r26-causal-raw-next-delta.v1",
            "parent_head":"15cb5de904ebeac196c5411ebdc6a29048662b3f",
            "classification":label,"training_majority_label":majority,
            "xorshift_majority_baseline_accuracy":baseline,
            "development_seeds":DEV,"heldout_seeds":HELDOUT,
            "models":arms,"model_calls":0,
            "claim_limit":"Supervised finite synthetic binary next-step forecasting only"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
