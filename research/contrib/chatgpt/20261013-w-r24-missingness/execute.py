"""W-R24: frozen observed-only fill-robust temporal difference learning."""
import importlib.util
import json
from pathlib import Path

PARENT=Path(__file__).resolve().parent.parent/"20261013-w-r23-fill-falsification"/"execute.py"
spec=importlib.util.spec_from_file_location("w23_frozen",PARENT)
assert spec and spec.loader
w23=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w23)
w22=w23.w22
w20=w23.w20
w19=w23.w19
w16=w23.w16

DEV=tuple(range(288,304))
HELDOUT=tuple(range(304,320))
TRAIN_SWAPS=(2,5,8)
TEST_SWAPS=("swap10_11","swap12_13")
TRAIN_MASK=(1,5,8,11)
HELDOUT_MASK=(1,4,7,9,15,6)
TRAIN_FILLS=("zero","carryforward","constant127")
TEST_FILLS=("zero","carryforward","constant127","constant73","linear_neighbour_interpolation")
GENERATORS=("affine_prior","quadratic_irregular","xorshift_independent")
ARMS=("naive_top2","zero_sentinel","structural_fill")


def apply_missing(observed,mask,fill):
    assert 0 not in mask
    if fill in ("zero","constant127","constant73"):
        replacement={"zero":(0,0,0,0),
                     "constant127":(127,127,127,127),
                     "constant73":(73,73,73,73)}[fill]
        return tuple(replacement if t in mask else row
                     for t,row in enumerate(observed))
    if fill=="carryforward":
        result=[]
        for t,row in enumerate(observed):
            result.append(result[-1] if t in mask else row)
        return tuple(result)
    assert fill=="linear_neighbour_interpolation"
    result=list(observed)
    for t in sorted(mask):
        left=(t-1)%16
        while left in mask:left=(left-1)%16
        right=(t+1)%16
        while right in mask:right=(right+1)%16
        length=(right-left)%16
        ratio=((t-left)%16)/length
        result[t]=tuple(int(round((1-ratio)*observed[left][channel]
                                  + ratio*observed[right][channel]))%256
                        for channel in range(4))
    return tuple(result)


def pair(seed,entity,amp,phase,swap,fill,count,positive,generator):
    base=w20.base_sequence(seed,entity,generator)
    candidate=w16.observed(base,seed,entity,amp,
                            17 if phase=="development" else 79)
    positions=TRAIN_MASK if phase=="development" else HELDOUT_MASK
    assert count in ((2,4) if phase=="development" else (2,4,6))
    missing=set(positions[:count])
    if phase=="development":
        assert not missing.intersection({3,4,6,7,9,10})
    else:
        assert not missing.intersection({0,11,12,13,14})
    edited=apply_missing(candidate,missing,fill)
    if positive:
        compared=edited
    elif phase=="development":
        compared=w22.swap_tail(edited,swap)
    else:
        compared=w19.permute(edited,swap)
    assert edited[0]==compared[0] and sorted(edited)==sorted(compared)
    assert (edited==compared)==positive
    return base,compared


def structural_feature(base,compared):
    invalid=set()
    for t,row in enumerate(compared):
        if len(set(row))==1 or row==compared[(t-1)%16]:
            invalid.add(t)
    errors=[]
    for channel in range(4):
        for t in range(16):
            nxt=(t+1)%16
            if t in invalid or nxt in invalid:
                continue
            left=(base[nxt][channel]-base[t][channel])%256
            right=(compared[nxt][channel]-compared[t][channel])%256
            errors.append(w16.circular_distance(left,right))
    top=sorted(errors,reverse=True)[:2]
    value=sum(top)/len(top) if top else 0.0
    return [value]*16


def feature(base,compared,arm):
    if arm=="naive_top2":
        return w19.feature(base,compared,"mean_top_two_64")
    if arm=="zero_sentinel":
        return w22.maskaware_feature(base,compared)
    assert arm=="structural_fill"
    return structural_feature(base,compared)


def train(arm):
    data=[]
    for seed in DEV:
        for entity in range(9):
            for swap in TRAIN_SWAPS:
                for amp in (1,2):
                    fill=TRAIN_FILLS[(seed+entity+swap+amp)%3]
                    count=(2,4)[(seed+entity+amp)%2]
                    for positive in (True,False):
                        a,b=pair(seed,entity,amp,"development",
                                 swap,fill,count,positive,"affine_prior")
                        data.append((feature(a,b,arm),1 if positive else -1))
    assert len(data)==1728
    weights=[0.0]*16
    bias=0.0
    updates=0
    for _ in range(32):
        for vec,label in data:
            score=bias+sum(w*v for w,v in zip(weights,vec))
            if label*score<=0:
                weights=[w+.1*label*v for w,v in zip(weights,vec)]
                bias+=.1*label
                updates+=1
    return weights,bias,updates


def test(arm,model,fill,generator,count,swap):
    weights,bias,_=model
    seeds=[]
    for seed in HELDOUT:
        tp=tn=fp=fn=0
        for entity in range(9):
            for positive in (True,False):
                a,b=pair(seed,entity,3,"heldout",swap,fill,count,positive,generator)
                vector=feature(a,b,arm)
                predicted=bias+sum(w*v for w,v in zip(weights,vector))>=0
                if positive and predicted:tp+=1
                elif not positive and not predicted:tn+=1
                elif not positive:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        seeds.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn})
    correct=sum(r["tp"]+r["tn"] for r in seeds)
    return {"fill":fill,"generator":generator,"count":count,"swap":swap,
            "correct":correct,"total":288,"accuracy":correct/288,
            "seed_rows":seeds}


def run():
    arms=[]
    for arm in ARMS:
        trained=train(arm)
        cells=[test(arm,trained,fill,generator,count,swap)
               for fill in TEST_FILLS for generator in GENERATORS
               for count in (2,4,6) for swap in TEST_SWAPS]
        assert len(cells)==90
        by_fill={fill:sum(c["correct"] for c in cells if c["fill"]==fill)/
                       sum(c["total"] for c in cells if c["fill"]==fill)
                 for fill in TEST_FILLS}
        arms.append({"arm":arm,"width":len(trained[0]),"epochs":32,
                     "training_examples":1728,"updates":trained[2],
                     "accuracy_by_fill":by_fill,"cells":cells})
    d={a["arm"]:a for a in arms}
    structural=d["structural_fill"]
    zero=d["zero_sentinel"]
    novel=("constant73","linear_neighbour_interpolation")
    structural_novel=sum(structural["accuracy_by_fill"][k] for k in novel)/2
    zero_novel=sum(zero["accuracy_by_fill"][k] for k in novel)/2
    all_novel=[c for c in structural["cells"] if c["fill"] in novel]
    if all(c["accuracy"]>=.8 for c in all_novel) and structural_novel-zero_novel>=.1:
        label="ENCODING_TRANSFER"
    elif structural_novel-zero_novel>=.05:
        label="PARTIAL"
    else:
        label="NEGATIVE"
    return {"schema":"wingless.w-r24-observed-only-fill-transfer.v1",
            "parent_head":"cf2da79ea3289aa82791448e103f86b89d97228e",
            "classification":label,"novel_gain_over_zero_detector":structural_novel-zero_novel,
            "development_seeds":DEV,"heldout_seeds":HELDOUT,
            "arms":arms,"model_calls":0,
            "claim_limit":"Synthetic observed-row structural features only, no semantic learning"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
