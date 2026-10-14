"""W-R74: test learnable global phase origin from raw temporal history.

Three fixed sixteen-component history encoders, each with three ridge heads
trained only on public development phase labels. Holdout phases are never fed
to the candidate classifier, and excluded source/identity remain sealed.
"""
import hashlib,importlib.util,json
from pathlib import Path
p=Path(__file__).resolve().parent.parent/"20261014-w-r73-canonical-temporal-state"/"execute.py"
spec=importlib.util.spec_from_file_location("w73_pinned",p)
assert spec and spec.loader
w73=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w73)
m=w73.m
PRE=w73.PRE
EXCLUDED="xorshift_independent"
MODES=("level16","delta16","absolute16")
PHASES=("intact","plus1","minus1")
TRAIN=tuple(range(9008,9024))
VALID=tuple(range(9040,9048))+tuple(range(9072,9080))
SEALED=tuple(range(16896,16912))
assert EXCLUDED in m.GENERATORS
def state(rows,mode):
    assert len(rows)==16 and all(len(r)==4 for r in rows)
    features=[]
    for t in (0,4,8,12):
        for ch in range(4):
            if mode=="level16": value=rows[t][ch]
            elif mode=="delta16": value=rows[t][ch]-rows[(t-1)%16][ch]
            elif mode=="absolute16":value=abs(rows[t][ch])
            else:raise ValueError("BAD_ENCODER")
            features.append(value)
    assert len(features)==16
    return tuple(features)
def build_train():
    samples=[]
    included=[g for g in m.GENERATORS if g!=EXCLUDED]
    for seed in TRAIN:
        for gi,generator in enumerate(m.GENERATORS):
            if generator==EXCLUDED:continue
            for entity in range(9):
                if entity%2!=(seed+gi)%2:continue
                for site in m.TRAIN_SITES:
                    for amplitude in (1,2):
                        base,_,_=m.w37.recoded_example(
                            seed,entity,generator,"identity",site,amplitude,True)
                        for phase in PHASES:
                            rows=PRE.acquire(base,(0,1,2,3),phase)
                            samples.append((rows,phase))
    assert len(samples)==2592
    return samples
def fit(samples):
    trained={}
    for mode in MODES:
        encoded=[(state(rows,mode),phase) for rows,phase in samples]
        heads={}
        for class_label in PHASES:
            labelled=[(features,1 if label==class_label else -1)
                      for features,label in encoded]
            weights=m.w53.ridge_fit(labelled)
            assert len(weights)==16
            heads[class_label]=weights
        trained[mode]=heads
    return trained
def guess(rows,heads,mode):
    x=state(rows,mode)
    scores={name:sum(w*v for w,v in zip(weights,x))
            for name,weights in heads.items()}
    return max(PHASES,key=lambda label:(scores[label],-PHASES.index(label)))
def eval_samples(seed_range,source_generators,entities,sites,amps,model,mode):
    stats={"total":0,"correct":0,"confusion":{p:{q:0 for q in PHASES} for p in PHASES}}
    for seed in seed_range:
        for generator in source_generators:
            for entity in entities:
                for site in sites:
                    for amp in amps:
                        base,_,_=m.w37.recoded_example(
                            seed,entity,generator,"identity",site,amp,False)
                        for phase in PHASES:
                            rows=PRE.acquire(base,(0,1,2,3),phase)
                            prediction=guess(rows,model,mode)
                            stats["total"]+=1
                            stats["correct"]+=int(prediction==phase)
                            stats["confusion"][phase][prediction]+=1
    return stats
def run():
    examples=build_train()
    fitted=fit(examples)
    included=tuple(g for g in m.GENERATORS if g!=EXCLUDED)
    dev={}
    for mode in MODES:
        r=eval_samples(VALID,included,range(9),m.TEST_SITES,(2,),fitted[mode],mode)
        dev[mode]={"accuracy":r["correct"]/r["total"],"n":r["total"]}
    selected=sorted(MODES,key=lambda mode:(-dev[mode]["accuracy"],mode))[0]
    trial_rows=[]
    matrix={phase:{p:0 for p in PHASES} for phase in PHASES}
    for seed in SEALED:
        for entity in range(9,18):
            for site in m.TEST_SITES:
                for amp in (1,2,3):
                    base,_,_=m.w37.recoded_example(
                        seed,entity,EXCLUDED,"identity",site,amp,False)
                    for phase in PHASES:
                        observed=PRE.acquire(base,(0,1,2,3),phase)
                        prediction=guess(observed,fitted[selected],selected)
                        trial_rows.append({"seed":seed,"entity":entity,"site":site,
                                           "amplitude":amp,"phase":phase,"predicted":prediction})
                        matrix[phase][prediction]+=1
    assert len(trial_rows)==16*9*len(m.TEST_SITES)*3*3
    n=len(trial_rows)
    correct=sum(matrix[p][p] for p in PHASES)
    acc=correct/n
    recalls={p:matrix[p][p]/sum(matrix[p].values()) for p in PHASES}
    if acc>=.45 and min(recalls.values())>=.40:
        label="UNSEEN_PHASE_ORIGIN_IDENTIFIABLE"
    elif acc>=.37:
        label="PARTIAL_PHASE_IDENTIFICATION"
    else:
        label="NO_PHASE_ORIGIN_IDENTIFIABILITY"
    return {"schema":"wingless.w-r74-learned-phase-origin.v1",
            "parent_head":"3237efbcd46c25e28f797e1588fc6d73233b59d7",
            "classification":label,"training_seeds":TRAIN,"validation_seeds":VALID,
            "sealed_heldout_seeds":SEALED,
            "source_excluded":EXCLUDED,"new_unseen_entity_ids":list(range(9,18)),
            "candidate_encoders":MODES,"selected_encoder":selected,
            "dev_accuracy":dev,"training_unique_items":len(examples),
            "inference_state_width":16,"inference_weight_count":48,
            "model_hashes":{mode:{head:hashlib.sha256(json.dumps(weights).encode()).hexdigest()
                      for head,weights in heads.items()} for mode,heads in fitted.items()},
            "baseline_uniform_chance_accuracy":1/3,
            "confusion":matrix,"accuracy":acc,"per_class_recalls":recalls,
            "heldout_records":trial_rows,"model_calls":0,"accepted_ref_mutation":False,
            "claim_limit":"Supervised three-class phase-origin prediction under synthetic hidden origin. State encoder itself hand-fixed; 48 learned readout weights and added compute. Not phase-robust reasoning or RSI."}
if __name__=="__main__":print(json.dumps(run(),sort_keys=True,indent=2))
