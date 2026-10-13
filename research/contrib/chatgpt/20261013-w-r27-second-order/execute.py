"""W-R27: preregistered causal next-delta learner using equal-width second differences."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r26-next-prediction"/"execute.py"
spec=importlib.util.spec_from_file_location("w26_frozen",SOURCE)
assert spec and spec.loader
w26=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w26)

DEV=tuple(range(416,432))
HELDOUT=tuple(range(432,448))
TRAIN_GENERATORS=("affine_prior","quadratic_irregular")
TEST_GENERATORS=("affine_prior","quadratic_irregular","xorshift_independent")
ARMS=("first_lag_16","first_plus_second_order_16","unordered_hist_16")


def signed_delta(a,b):
    return (((b-a)%256+128)%256-128)/128.0


def causal_feature(observed,t,arm):
    assert 4<=t<=14
    window=observed[t-4:t+1]
    assert len(window)==5
    if arm=="unordered_hist_16":
        f=[0.0]*16
        for row in window:f[row[0]%16]+=1.0/5
    else:
        f=[]
        for channel in range(4):
            past=[signed_delta(window[i][channel],window[i+1][channel])
                  for i in range(4)]
            if arm=="first_lag_16":
                f.extend(past)
            else:
                assert arm=="first_plus_second_order_16"
                second=[signed_delta(int(round(past[i]*128))%256,
                                     int(round(past[i+1]*128))%256)
                        for i in range(3)]
                f.extend([past[-1]]+second)
    assert len(f)==16
    return f


def examples(seed,generator,amp,phase,arm):
    rows=[]
    for entity in range(9):
        base,observed=w26.signature(seed,entity,generator,amp,phase)
        for t in range(4,15):
            f=causal_feature(observed,t,arm)
            target=w26.label_at(base,t)
            rows.append((f,target))
    assert len(rows)==99
    return rows


def train(arm):
    data=[r for seed in DEV for gen in TRAIN_GENERATORS
          for amp in (1,2) for r in examples(seed,gen,amp,"development",arm)]
    assert len(data)==6336
    weights=[0.0]*16
    bias=0.0
    updates=0
    for _ in range(32):
        for f,label in data:
            if label*(bias+sum(w*x for w,x in zip(weights,f)))<=0:
                weights=[w+.1*label*x for w,x in zip(weights,f)]
                bias+=.1*label
                updates+=1
    return weights,bias,updates


def evaluate(arm,model,generator):
    weights,bias,_=model
    seeds=[]
    for seed in HELDOUT:
        tp=tn=fp=fn=0
        for features,actual in examples(seed,generator,3,"heldout",arm):
            pred=1 if bias+sum(w*x for w,x in zip(weights,features))>=0 else -1
            if actual==1 and pred==1:tp+=1
            elif actual==-1 and pred==-1:tn+=1
            elif actual==-1:fp+=1
            else:fn+=1
        assert tp+tn+fp+fn==99
        seeds.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn})
    good=sum(r["tp"]+r["tn"] for r in seeds)
    return {"generator":generator,"accuracy":good/(16*99),
            "correct":good,"total":16*99,"heldout":seeds}


def run():
    out=[]
    for arm in ARMS:
        model=train(arm)
        out.append({"arm":arm,"training_examples":6336,"width":16,"epochs":32,
                    "updates":model[2],
                    "generators":[evaluate(arm,model,g) for g in TEST_GENERATORS]})
    a={m["arm"]:{r["generator"]:r["accuracy"] for r in m["generators"]} for m in out}
    gain=a["first_plus_second_order_16"]["quadratic_irregular"]-a["first_lag_16"]["quadratic_irregular"]
    if a["first_plus_second_order_16"]["quadratic_irregular"]>=.68 and gain>=.05 and a["first_plus_second_order_16"]["xorshift_independent"]<=.60:
        label="SECOND_ORDER_GAIN"
    elif gain>=.03:
        label="PARTIAL"
    else:
        label="NEGATIVE"
    return {"schema":"wingless.w-r27-second-order-future-delta.v1",
            "parent_head":"7b36200f66fc9ea766a6451190637a5a57c47ccb",
            "classification":label,"quadratic_delta_gain":gain,
            "development_seeds":DEV,"sealed_heldout_seeds":HELDOUT,
            "arms":out,"model_calls":0,
            "claim_limit":"Synthetic causal binary next-delta prediction only"}

if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
