"""W-R18 frozen local singlet temporal transfer on wholly unexposed cohorts."""
import importlib.util
import json
from pathlib import Path

SRC=Path(__file__).resolve().parent.parent/"20261013-w-r16b-noise-matrix"/"execute.py"
spec=importlib.util.spec_from_file_location("w16b_frozen",SRC)
assert spec and spec.loader
w16=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w16)

DEVELOPMENT=tuple(range(128,144))
HELDOUT=tuple(range(144,160))
TRAINING={
 "coarse_only":("reverse_tail","rotate3","weave_tail"),
 "single_adjacent":("swap2_3","swap5_6","swap8_9")
}
EVAL_FAMILIES=("swap10_11","swap12_13")
ARMS=("pooled","max_transition","top2_transition")


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
    if kind=="pooled":
        return w16.representation(a,b,"mean_mismatch")
    transitions=w16.representation(a,b,"transition")
    if kind=="max_transition":
        value=max(transitions)
    else:
        assert kind=="top2_transition"
        value=sum(sorted(transitions,reverse=True)[:2])/2.0
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
    scores={(x["curriculum"],x["representation"]):x for x in rows}
    gains={k:scores[("single_adjacent",k)]["accuracy"]-
             scores[("coarse_only",k)]["accuracy"] for k in ARMS}
    if any(
        scores[("single_adjacent",k)]["accuracy"]-
        scores[("coarse_only",k)]["accuracy"]>=.15
        and all(c["accuracy"]>=.80 for c in scores[("single_adjacent",k)]["matrix"])
        for k in ARMS):
        classification="LOCAL_SINGLE_SWAP_TRANSFER"
    elif any(value>=.05 for value in gains.values()):
        classification="PARTIAL"
    else:
        classification="NEGATIVE"
    return {"schema":"wingless.w-r18-frozen-singlet-transfer.v1",
            "parent_head":"6b5024543afe21509b6e7f1acdd6992064963c8b",
            "classification":classification,"gains":gains,
            "development_seeds":DEVELOPMENT,"heldout_seeds":HELDOUT,
            "arms":rows,"model_calls":0,
            "claim_limit":"Synthetic temporal local interchange discrimination only"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
