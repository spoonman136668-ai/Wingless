"""W-R17: frozen local-order curricula transfer with unexposed cohorts."""
import importlib.util
import json
import hashlib
from pathlib import Path

BASE=Path(__file__).resolve().parent.parent/"20261013-w-r16b-noise-matrix"/"execute.py"
spec=importlib.util.spec_from_file_location("w16b_frozen",BASE)
assert spec and spec.loader
w16=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w16)

TRAIN_SEEDS=tuple(range(96,112))
TEST_SEEDS=tuple(range(112,128))
CURRICULA={
 "global_only":("reverse_tail","rotate3","weave_tail"),
 "local_mixed":("reverse_tail","pairwise_swap_tail","swap_tail_blocks"),
}
REPRESENTATIONS=("circular_position","mean_mismatch","transition")
TEST_FAMILIES=("offset_pair_swaps","single_middle_adjacent_swap")


def test_transform(sequence,family):
    tail=list(sequence[1:])
    if family=="offset_pair_swaps":
        for i in range(1,14,2):
            tail[i],tail[i+1]=tail[i+1],tail[i]
    else:
        assert family=="single_middle_adjacent_swap"
        tail[6],tail[7]=tail[7],tail[6]
    out=sequence[:1]+tuple(tail)
    assert out[0]==sequence[0] and sorted(out)==sorted(sequence)
    assert out!=sequence
    return out


def make_pair(seed,entity,amplitude,phase,family,positive):
    base=w16.original(seed,entity)
    offset=17 if phase=="development" else 79
    cand=w16.observed(base,seed,entity,amplitude,offset)
    if positive:
        other=cand
    elif phase=="development":
        other=w16.permute(cand,family)
    else:
        other=test_transform(cand,family)
    assert other[0]==cand[0] and sorted(other)==sorted(cand)
    assert (other==cand)==positive
    return base,other


def train(curriculum,kind):
    examples=[]
    for seed in TRAIN_SEEDS:
        for entity in range(9):
            for family in CURRICULA[curriculum]:
                for amp in (1,2):
                    for positive in (True,False):
                        base,cand=make_pair(seed,entity,amp,"development",family,positive)
                        examples.append((w16.representation(base,cand,kind),1 if positive else -1))
    assert len(examples)==1728
    weights=[0.0]*16
    bias=0.0
    updates=0
    for _ in range(32):
        for vec,label in examples:
            score=bias+sum(w*v for w,v in zip(weights,vec))
            if label*score<=0:
                weights=[w+0.1*label*v for w,v in zip(weights,vec)]
                bias+=0.1*label
                updates+=1
    digest=hashlib.sha256(json.dumps([weights,bias]).encode()).hexdigest()
    return weights,bias,updates,digest


def evaluate(kind,model,amp,family):
    weights,bias,updates,sha=model
    seeds=[]
    for seed in TEST_SEEDS:
        tp=tn=fp=fn=0
        for entity in range(9):
            for positive in (True,False):
                a,b=make_pair(seed,entity,amp,"heldout",family,positive)
                vector=w16.representation(a,b,kind)
                pred=(bias+sum(w*v for w,v in zip(weights,vector)))>=0
                if positive and pred:tp+=1
                elif not positive and not pred:tn+=1
                elif not positive:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        seeds.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn})
    correct=sum(r["tp"]+r["tn"] for r in seeds)
    return {"amp":amp,"family":family,"accuracy":correct/288,
            "correct":correct,"total":288,"by_seed":seeds}


def run():
    arms=[]
    for curriculum in CURRICULA:
        for kind in REPRESENTATIONS:
            model=train(curriculum,kind)
            cells=[evaluate(kind,model,amp,family)
                   for amp in (1,2,3) for family in TEST_FAMILIES]
            assert len(cells)==6
            acc=sum(c["correct"] for c in cells)/sum(c["total"] for c in cells)
            arms.append({"curriculum":curriculum,"representation":kind,"accuracy":acc,
                         "training_examples":1728,"feature_width":16,"epochs":32,
                         "updates":model[2],"model_sha256":model[3],"cells":cells})
    metrics={(r["curriculum"],r["representation"]):r["accuracy"] for r in arms}
    diffs={k:metrics[("local_mixed",k)]-metrics[("global_only",k)] for k in REPRESENTATIONS}
    if diffs["mean_mismatch"]>=0.05 and diffs["transition"]>=0.05:
        classification="MIXED_CURRICULUM_GAIN"
    elif any(g>=0.05 for g in diffs.values()):
        classification="TRADEOFF"
    else:
        classification="NO_CROSS_FEATURE_GAIN"
    return {"schema":"wingless.w-r17-local-order-curriculum.v1",
            "parent_head":"a8e9d82a8abbb4a423dfa0ffac2d99af026a1005",
            "classification":classification,"mean_model_differences":diffs,
            "development_seeds":TRAIN_SEEDS,"heldout_seeds":TEST_SEEDS,
            "arms":arms,"model_calls":0,
            "claim_limit":"Supervised synthetic local order experiment; not general intelligence"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
