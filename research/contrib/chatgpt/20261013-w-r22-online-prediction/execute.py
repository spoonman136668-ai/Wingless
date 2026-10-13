"""W-R22: frozen next-step raw sensor prediction with causal online curvature."""
import json

DEVELOPMENT=range(80,84)
HELDOUT=range(84,88)
PREDICTORS=("persistence","momentum","frozen_curvature","online_curvature")


def series(seed,entity,phase):
    development=phase=="development"
    q=(5*entity+seed)%9
    modulation=3 if development else 5
    multiplier=5 if development else 7
    jitter=[(multiplier*t+entity+seed)%modulation-modulation//2
            for t in range(16)]
    assert len(set(jitter))==modulation
    return tuple((
        192+seed+7*q+3*t+t*t+
        (0 if development else 2*max(0,t-8)**2)+jitter[t]
    )%256 for t in range(16))


def circular_signed(value):
    return (value+128)%256-128


def circular_abs(predicted,observed):
    return abs(circular_signed(predicted-observed))


def learn_dev_curvature():
    differences=[]
    for seed in DEVELOPMENT:
        for entity in range(9):
            seq=series(seed,entity,"development")
            for tick in range(2,16):
                diff=seq[tick]-2*seq[tick-1]+seq[tick-2]
                differences.append(circular_signed(diff))
    assert len(differences)==4*9*14
    differences.sort()
    n=len(differences)
    estimate=(differences[n//2-1]+differences[n//2])/2
    return estimate,{"count":n,"q0":differences[0],
                      "median_curvature":estimate,"q100":differences[-1]}


def predict(seq,tick,arm,curvature):
    last=seq[tick-1]
    before=seq[tick-2]
    if arm=="persistence":return last
    if arm=="momentum":return 2*last-before
    if arm=="frozen_curvature":return 2*last-before+curvature
    assert arm=="online_curvature"
    if tick==2:local=curvature
    else:
        local=circular_signed(seq[tick-1]-2*seq[tick-2]+seq[tick-3])
    return 2*last-before+local


def assess(arm,curvature):
    per_seed=[]
    total=[]
    pre=[]
    post=[]
    for seed in HELDOUT:
        errors=[]
        for entity in range(9):
            observed=series(seed,entity,"heldout")
            for t in range(2,16):
                error=circular_abs(predict(observed,t,arm,curvature),observed[t])
                assert 0<=error<=128
                record={"tick":t,"entity":entity,"error":error}
                errors.append(record)
        assert len(errors)==9*14
        aggregate=[r["error"] for r in errors]
        early=[r["error"] for r in errors if r["tick"]<=8]
        late=[r["error"] for r in errors if r["tick"]>=9]
        assert len(early)==len(late)==9*7
        per_seed.append({"seed":seed,
                         "mae":sum(aggregate)/len(aggregate),
                         "prechange_mae":sum(early)/len(early),
                         "postchange_mae":sum(late)/len(late),
                         "n":len(aggregate)})
        total.extend(aggregate)
        pre.extend(early)
        post.extend(late)
    return {"arm":arm,"total_mae":sum(total)/len(total),
            "prechange_mae":sum(pre)/len(pre),
            "postchange_mae":sum(post)/len(post),
            "n":len(total),"by_seed":per_seed}


def run():
    curvature,provenance=learn_dev_curvature()
    arms=[assess(arm,curvature) for arm in PREDICTORS]
    results={r["arm"]:r for r in arms}
    learned=results["frozen_curvature"]
    online=results["online_curvature"]
    momentum=results["momentum"]
    improvements=sum(a["postchange_mae"]<b["postchange_mae"]
                     for a,b in zip(online["by_seed"],learned["by_seed"]))
    if (online["postchange_mae"]<=.90*learned["postchange_mae"]
        and online["postchange_mae"]<=momentum["postchange_mae"]
        and improvements>=3):
        label="ADAPTIVE_PREDICTION_TRANSFER"
    elif online["postchange_mae"]<learned["postchange_mae"]:
        label="PARTIAL"
    else:
        label="NEGATIVE"
    return {
        "schema":"wingless.w-r22-raw-prediction.v1",
        "predecessor_head":"cf92ea14e2806b3a7162383da819107286fa2515",
        "classification":label,
        "learned_curvature":curvature,
        "training_provenance":provenance,
        "heldout_seeds_improved":improvements,
        "models":arms,"model_calls":0,
        "scope":"Toy raw tick prediction, causal past-only adaptation, not language reasoning"
    }


if __name__=="__main__":
    result=run()
    assert result==run()
    print(json.dumps(result,sort_keys=True,indent=2))
