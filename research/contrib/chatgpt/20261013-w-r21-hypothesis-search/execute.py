"""W-R21: frozen 8-hypothesis second-order rule selection from four labels."""
import hashlib
import json

CANDIDATES = tuple((a, b) for a in (1, 2) for b in range(4))
assert len(CANDIDATES) == 8


def sequence(seed, episode):
    assert 72 <= seed <= 75 and 0 <= episode < 16
    a = 2 if seed % 2 == 0 else 1
    b = (seed - 72) % 4
    values = [(episode+seed)%4, ((episode//4)+3*seed)%4]
    while len(values)<32:
        values.append((values[-1]+a*values[-2]+b)%4)
    return tuple(values)


def observe_four(seed):
    return [(sequence(seed,e)[0],sequence(seed,e)[1],sequence(seed,e)[2])
            for e in range(4)]


def select_model(rows):
    matching=[
        (a,b) for a,b in CANDIDATES
        if all((prev+bias*a+offset)%4==target
               for bias,prev,target in rows)
    ]
    # The list-comprehension intentionally calculates candidate
    # x_last + a*x_before + b in the exact preregistered order.
    return matching


def rote_table(rows):
    table=[[0]*4 for _ in range(16)]
    for old,last,label in rows:
        table[4*old+last][label]+=1
    assert sum(sum(r) for r in table)==4
    return table


def predict(arm,past,table,selected):
    if arm=="fixed_previous_rule":
        return (past[-1]+past[-2]+2)%4
    if arm=="four_example_rote":
        counts=table[4*past[-2]+past[-1]]
        return max(range(4),key=lambda c:(counts[c],-c))
    a,b=selected
    if arm=="temporal_shuffle_selected":
        past=past[::2]+past[1::2]
    else:
        assert arm=="hypothesis_selected"
    return (past[-1]+a*past[-2]+b)%4


def evaluate(seed,arm,table,selected):
    cm=[[0]*4 for _ in range(4)]
    digest=hashlib.sha256()
    for episode in range(4,16):
        seq=sequence(seed,episode)
        for tick in range(2,32):
            pred=predict(arm,seq[:tick],table,selected)
            true=seq[tick]
            cm[true][pred]+=1
            digest.update(bytes((episode,tick,true,pred)))
    total=sum(sum(r) for r in cm)
    correct=sum(cm[i][i] for i in range(4))
    assert total==360
    return {
        "seed":seed,"arm":arm,"accuracy":correct/total,
        "correct":correct,"examples":total,"confusion":cm,
        "trace_sha256":digest.hexdigest()
    }


def run():
    arms=("fixed_previous_rule","four_example_rote",
          "hypothesis_selected","temporal_shuffle_selected")
    selected={}
    ambiguous=[]
    evaluations=[]
    for seed in range(72,76):
        calibrate=observe_four(seed)
        hypotheses=select_model(calibrate)
        if len(hypotheses)!=1:
            ambiguous.append({"seed":seed,"matching_hypotheses":hypotheses})
            continue
        a,b=hypotheses[0]
        selected[str(seed)]={"a":a,"b":b}
        table=rote_table(calibrate)
        evaluations.extend(evaluate(seed,arm,table,hypotheses[0]) for arm in arms)
    if ambiguous:
        label="HYPOTHESIS_AMBIGUITY"
        comparisons=[]
    else:
        comparisons=[{
            "arm":arm,
            "heldout_accuracy":sum(x["correct"] for x in evaluations
                                  if x["arm"]==arm)/(4*360),
            "per_seed":[x for x in evaluations if x["arm"]==arm]
        } for arm in arms]
        scores={row["arm"]:row["heldout_accuracy"] for row in comparisons}
        if (scores["hypothesis_selected"]>=.95
                and all(scores[arm]<=.65 for arm in arms if arm!="hypothesis_selected")):
            label="BOUNDED_HYPOTHESIS_TRANSFER"
        elif scores["hypothesis_selected"]>max(
                scores[arm] for arm in arms if arm!="hypothesis_selected"):
            label="PARTIAL"
        else:
            label="NEGATIVE"
    return {
        "schema":"wingless.w-r21-frozen-hypothesis-search.v1",
        "predecessor_head":"baa8e1eadde3fd10bba761735a50ff4d6944a82b",
        "classification":label,
        "candidate_hypothesis_count":8,
        "calibration_labels_per_seed":4,
        "selected_hypotheses":selected,
        "ambiguous":ambiguous,
        "arms":comparisons,"model_calls":0,
        "claim_limit":"Small fixed hypothesis search on synthetic modular recurrences; not unconstrained rule induction."
    }


if __name__=="__main__":
    outcome=run()
    assert outcome==run()
    print(json.dumps(outcome,sort_keys=True,indent=2))
