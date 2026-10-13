"""W-R20: four-example compositional transition rule transfer."""
import hashlib
import json


def sequence(seed, episode):
    assert 64 <= seed <= 67 and 0 <= episode < 16
    x = [(episode + seed) % 4, ((episode//4)+3*seed)%4]
    while len(x)<32:
        x.append((x[-1]+x[-2]+2)%4)
    return tuple(x)


def calibration(seed):
    observations=[]
    for episode in range(4):
        x=sequence(seed,episode)
        observations.append((x[0],x[1],x[2]))
    assert len(observations)==4
    assert len({(a,b) for a,b,_ in observations})==4
    return observations


def infer_bias(exemplars):
    residuals=[(label-a-b)%4 for a,b,label in exemplars]
    if len(set(residuals))!=1:
        return None,residuals
    return residuals[0],residuals


def rote(exemplars):
    rows=[[0]*4 for _ in range(16)]
    for a,b,label in exemplars:
        rows[4*a+b][label]+=1
    assert sum(sum(row) for row in rows)==4
    return rows


def predict(arm, history, rote_table, bias):
    if arm=="old_fixed_b0":
        return (history[-1]+history[-2])%4
    if arm=="rote_16_context":
        row=rote_table[4*history[-2]+history[-1]]
        return max(range(4),key=lambda val:(row[val],-val))
    if arm=="compositional_bias":
        return (history[-1]+history[-2]+bias)%4
    assert arm=="shuffled_history_compositional"
    past=tuple(history)
    scrambled=past[::2]+past[1::2]
    assert sorted(scrambled)==sorted(past)
    return (scrambled[-1]+scrambled[-2]+bias)%4


def run_seed(seed,arm,table,bias):
    cm=[[0]*4 for _ in range(4)]
    digest=hashlib.sha256()
    for episode in range(4,16):
        seq=sequence(seed,episode)
        for tick in range(2,32):
            pred=predict(arm,seq[:tick],table,bias)
            true=seq[tick]
            cm[true][pred]+=1
            digest.update(bytes((episode,tick,true,pred)))
    correct=sum(cm[i][i] for i in range(4))
    total=sum(sum(row) for row in cm)
    assert total==12*30
    return {
        "seed":seed,"correct":correct,"examples":total,
        "accuracy":correct/total,"confusion":cm,
        "prediction_trace_sha256":digest.hexdigest()
    }


def run():
    arms=("old_fixed_b0","rote_16_context","compositional_bias",
          "shuffled_history_compositional")
    rows=[]
    inferred_offsets={}
    for seed in range(64,68):
        examples=calibration(seed)
        bias,residuals=infer_bias(examples)
        assert bias is not None and len(set(residuals))==1
        inferred_offsets[str(seed)]=bias
        table=rote(examples)
        for arm in arms:
            row=run_seed(seed,arm,table,bias)
            row["arm"]=arm
            rows.append(row)
    results=[{
        "arm":arm,"heldout_accuracy":sum(r["correct"] for r in rows if r["arm"]==arm)/(4*12*30),
        "by_seed":[r for r in rows if r["arm"]==arm]
    } for arm in arms]
    scores={r["arm"]:r["heldout_accuracy"] for r in results}
    if scores["compositional_bias"]>=.95 and all(
        scores[x]<=.65 for x in arms if x!="compositional_bias"):
        label="COMPOSITIONAL_GENERALIZATION"
    elif scores["compositional_bias"]>max(
        scores[x] for x in arms if x!="compositional_bias"):
        label="PARTIAL"
    else:
        label="NEGATIVE"
    return {
        "schema":"wingless.w-r20-four-demo-composition.v1",
        "predecessor_head":"fc6056a230156401db4e4cd48f7ca918f43019f8",
        "classification":label,
        "inferred_train_only_bias_by_seed":inferred_offsets,
        "calibration_labels_per_seed":4,
        "test_predictions_per_seed":360,
        "arms":results,"model_calls":0,
        "claim_limit":"A predefined modular arithmetic rule family explains successful transfer from 4 examples; not general symbolic discovery."
    }


if __name__=="__main__":
    result=run()
    assert result==run()
    print(json.dumps(result,sort_keys=True,indent=2))
