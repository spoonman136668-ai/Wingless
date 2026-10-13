"""W-R19 frozen four-channel single-adjacent robustness on wholly unexposed cohorts."""
import importlib.util
import json
from pathlib import Path

SRC=Path(__file__).resolve().parent.parent/"20261013-w-r16b-noise-matrix"/"execute.py"
spec=importlib.util.spec_from_file_location("w16b_frozen",SRC)
assert spec and spec.loader
w16=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w16)

DEVELOPMENT=tuple(range(160,176))
HELDOUT=tuple(range(176,192))
TRAINING={"single_adjacent":("swap2_3","swap5_6","swap8_9")}
EVAL_FAMILIES=("swap10_11","swap12_13")
ARMS=("first_channel_max","mean_four_max","median_four_max","mean_top_two_64")


def permute(seq,transform):
    if transform in ("reverse_tail","rotate3","weave_tail"):
        return w16.permute(seq,transform)
    index={"swap2_3":2,"swap5_6":5,"swap8_9":8,
           "swap10_11":10,"swap12_13":12}[transform]
    tail=list(seq[1:])
    tail[index],tail[index+1]=tail[index+1],tail[index]
    out=seq[:1]+tuple(tail)
    assert out[0]==seq[0] and sorted(out)==sorted(seq)
    assert out!=seq
    return out


def pair(seed,entity,amp,phase,transform,positive):
    base=w16.original(seed,entity)
    candidate=w16.observed(base,seed,entity,amp,17 if phase=="development" else 79)
    altered=candidate if positive else permute(candidate,transform)
    assert candidate[0]==altered[0] and sorted(candidate)==sorted(altered)
    assert (candidate==altered)==positive
    return base,altered


def feature(a,b,kind):
    channel_maxima=[]
    all_differences=[]
    for channel in range(4):
        va=[row[channel] for row in a]
        vb=[row[channel] for row in b]
        da=[(va[(i+1)%16]-va[i])%256 for i in range(16)]
        db=[(vb[(i+1)%16]-vb[i])%256 for i in range(16)]
        err=[w16.circular_distance(x,y) for x,y in zip(da,db)]
        channel_maxima.append(max(err))
        all_differences.extend(err)
    if kind=="first_channel_max":
        value=channel_maxima[0]
    elif kind=="mean_four_max":
        value=sum(channel_maxima)/4
    elif kind=="median_four_max":
        top=sorted(channel_maxima)
        value=(top[1]+top[2])/2
    else:
        assert kind=="mean_top_two_64"
        value=sum(sorted(all_differences,reverse=True)[:2])/2
    return [value]*16


def train(curriculum,kind):
    examples=[]
    for seed in DEVELOPMENT:
        for entity in range(9):
            for transform in TRAINING[curriculum]:
                for amp in (1,2):
                    for positive in (True,False):
                        a,b=pair(seed,entity,amp,"development",transform,positive)
                        examples.append((feature(a,b,kind),1 if positive else -1))
    assert len(examples)==1728
    weights=[0.0]*16
    bias=0.0
    updates=0
    for _ in range(32):
        for vector,label in examples:
            signed=label*(bias+sum(w*v for w,v in zip(weights,vector)))
            if signed<=0:
                weights=[w+.1*label*v for w,v in zip(weights,vector)]
                bias+=.1*label
                updates+=1
    return weights,bias,updates


def evaluate(kind,model,amp,transform):
    weights,bias,updates=model
    by_seed=[]
    for seed in HELDOUT:
        tp=tn=fp=fn=0
        for entity in range(9):
            for positive in (True,False):
                a,b=pair(seed,entity,amp,"heldout",transform,positive)
                vector=feature(a,b,kind)
                predicted=bias+sum(w*v for w,v in zip(weights,vector))>=0
                if positive and predicted: tp+=1
                elif not positive and not predicted:tn+=1
                elif not positive:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        by_seed.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn})
    correct=sum(r["tp"]+r["tn"] for r in by_seed)
    return {"amp":amp,"transform":transform,"correct":correct,
            "accuracy":correct/288,"by_seed":by_seed}


def run():
    rows=[]
    for curriculum in TRAINING:
        for kind in ARMS:
            model=train(curriculum,kind)
            matrix=[evaluate(kind,model,amp,transform)
                    for amp in (1,2,3) for transform in EVAL_FAMILIES]
            total=sum(c["correct"] for c in matrix)
            rows.append({"curriculum":curriculum,"representation":kind,
                         "accuracy":total/(6*288),"updates":model[2],
                         "width":16,"training_examples":1728,"epochs":32,
                         "matrix":matrix})
    by={x["representation"]:x for x in rows}
    first=by["first_channel_max"]["accuracy"]
    gains={k:by[k]["accuracy"]-first for k in ARMS[1:]}
    robust=any(
        gain>=.05 and all(c["accuracy"]>=.80
                           for c in by[k]["matrix"] if c["amp"]==3)
        for k,gain in gains.items())
    classification=("MULTICHANNEL_ROBUSTNESS" if robust else
                    "PARTIAL" if any(g>=.03 for g in gains.values())
                    else "NO_GAIN")
    return {"schema":"wingless.w-r19-multichannel-local-robustness.v1",
            "parent_head":"bf4aa817f6ea00e06e3327549811b5c7e7255133",
            "classification":classification,"gains":gains,
            "development_seeds":DEVELOPMENT,"heldout_seeds":HELDOUT,
            "arms":rows,"model_calls":0,
            "claim_limit":"Fixed-width classifier on four already observed synthetic channels"}



if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
