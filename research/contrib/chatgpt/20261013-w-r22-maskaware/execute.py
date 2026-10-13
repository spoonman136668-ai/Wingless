"""W-R22: frozen fixed-width mask-aware learning on disjoint cohorts."""
import importlib.util
import json
from pathlib import Path

SOURCE = Path(__file__).resolve().parent.parent / "20261013-w-r21-occlusion" / "execute.py"
_spec = importlib.util.spec_from_file_location("w21_frozen", SOURCE)
assert _spec and _spec.loader
w21 = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(w21)
w20 = w21.w20
w19 = w20.w19
w16 = w19.w16

TRAIN_SEEDS = tuple(range(240,256))
HELDOUT_SEEDS = tuple(range(256,272))
TRAIN_SWAP_INDICES = (2,5,8)
TEST_SWAP_FAMILIES = ("swap10_11","swap12_13")
MASK_ORDER = (2,5,8,10,3,6)
ARMS = ("naive_clean","naive_augmented","maskaware_clean","maskaware_augmented")


def swap_tail(candidate, index):
    assert index in (2,5,8)
    tail=list(candidate[1:])
    tail[index],tail[index+1]=tail[index+1],tail[index]
    out=candidate[:1]+tuple(tail)
    assert out[0]==candidate[0] and sorted(out)==sorted(candidate)
    assert out!=candidate
    return out


def prepared(seed,entity,generator,amp,phase,mask_count,negative,swap):
    base=w20.base_sequence(seed,entity,generator)
    observed=w16.observed(base,seed,entity,amp,17 if phase=="development" else 79)
    assert mask_count in (0,2,3,4,5)
    mask=set(MASK_ORDER[:mask_count])
    assert not mask.intersection({0,11,12,13,14})
    candidate=tuple((0,0,0,0) if tick in mask else row
                    for tick,row in enumerate(observed))
    if negative:
        other=swap_tail(candidate,swap) if phase=="development" else w19.permute(candidate,swap)
    else:
        other=candidate
    assert other[0]==candidate[0] and sorted(other)==sorted(candidate)
    assert (other==candidate)==(not negative)
    return base,other


def maskaware_feature(a,b):
    differences=[]
    for channel in range(4):
        a_vals=[row[channel] for row in a]
        b_vals=[row[channel] for row in b]
        for i in range(16):
            nxt=(i+1)%16
            # Zeroed four-channel rows are observed, not a privileged oracle mask.
            if b[i]==(0,0,0,0) or b[nxt]==(0,0,0,0):
                continue
            da=(a_vals[nxt]-a_vals[i])%256
            db=(b_vals[nxt]-b_vals[i])%256
            differences.append(w16.circular_distance(da,db))
    if len(differences)>=2:
        value=sum(sorted(differences,reverse=True)[:2])/2
    elif differences:
        value=differences[0]
    else:
        value=0.0
    return [value]*16


def featurize(a,b,arm):
    if arm.startswith("maskaware"):
        vec=maskaware_feature(a,b)
    else:
        vec=w19.feature(a,b,"mean_top_two_64")
    assert len(vec)==16
    return vec


def train(arm):
    examples=[]
    augmented=arm.endswith("augmented")
    for seed in TRAIN_SEEDS:
        for entity in range(9):
            for index in TRAIN_SWAP_INDICES:
                for amp in (1,2):
                    count=(0,2,4)[(seed+entity+index+amp)%3] if augmented else 0
                    for negative in (False,True):
                        a,b=prepared(seed,entity,"affine_prior",amp,
                                     "development",count,negative,index)
                        examples.append((featurize(a,b,arm),-1 if negative else 1))
    assert len(examples)==1728
    weights=[0.0]*16
    bias=0.0
    updates=0
    for _ in range(32):
        for vec,label in examples:
            score=bias+sum(w*x for w,x in zip(weights,vec))
            if label*score<=0:
                weights=[w+0.1*label*x for w,x in zip(weights,vec)]
                bias+=0.1*label
                updates+=1
    return weights,bias,updates


def test(arm,model,generator,masked):
    weights,bias,updates=model
    seeds=[]
    for seed in HELDOUT_SEEDS:
        tp=tn=fp=fn=0
        for entity in range(9):
            for swap in TEST_SWAP_FAMILIES:
                pos_a,pos_b=prepared(seed,entity,generator,3,
                                     "heldout",masked,False,swap)
                neg_a,neg_b=prepared(seed,entity,generator,3,
                                     "heldout",masked,True,swap)
                assert pos_a==neg_a
                assert pos_b[0]==neg_b[0] and sorted(pos_b)==sorted(neg_b)
                for positive,a,b in ((True,pos_a,pos_b),(False,neg_a,neg_b)):
                    vec=featurize(a,b,arm)
                    prediction=bias+sum(w*x for w,x in zip(weights,vec))>=0
                    if positive and prediction:tp+=1
                    elif not positive and not prediction:tn+=1
                    elif not positive:fp+=1
                    else:fn+=1
        assert tp+tn+fp+fn==36
        seeds.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn})
    good=sum(s["tp"]+s["tn"] for s in seeds)
    return {"generator":generator,"masked_ticks":masked,
            "correct":good,"total":16*36,"accuracy":good/(16*36),
            "heldout_by_seed":seeds}


def run():
    arms=[]
    for arm in ARMS:
        model=train(arm)
        cells=[test(arm,model,generator,count)
               for generator in w20.GENERATORS for count in (3,5)]
        arms.append({"arm":arm,"width":16,"training_examples":1728,
                     "epochs":32,"training_updates":model[2],
                     "accuracy":sum(x["correct"] for x in cells)/
                                sum(x["total"] for x in cells),
                     "cells":cells})
    m={x["arm"]:x for x in arms}
    gain=m["maskaware_augmented"]["accuracy"]-m["naive_augmented"]["accuracy"]
    if all(c["accuracy"]>=.8 for c in m["maskaware_augmented"]["cells"]) and gain>=.10:
        classification="MASK_AWARE_GENERALIZATION"
    elif gain>=.05:
        classification="PARTIAL"
    else:
        classification="NEGATIVE"
    return {"schema":"wingless.w-r22-mask-aware-augmentation.v1",
            "parent_head":"fd54ad48b14aa493ee922b057f23c3b75c81fe76",
            "classification":classification,"gain_maskaware_aug_vs_naive_aug":gain,
            "development_seeds":TRAIN_SEEDS,"sealed_heldout_seeds":HELDOUT_SEEDS,
            "arms":arms,"model_calls":0,
            "claim_limit":"Synthetic paired temporal observation missingness only"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
