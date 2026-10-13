"""W-R22: separate new first-use eight-candidate rule identification, fixture checked."""
import hashlib
import json

CANDIDATES = tuple((a, b) for a in (1,2) for b in range(4))
assert len(CANDIDATES)==8


def select_hypotheses(observations):
    return [
        (a,b) for a,b in CANDIDATES
        if all((last+a*before+b)%4 == label
               for before,last,label in observations)
    ]


def smoke_fixture_only():
    # Manually authored sentinel tuples; NO heldout-seed data used here.
    sentinel=[(0,0,1),(1,0,2),(2,0,3),(3,0,0)]
    assert select_hypotheses(sentinel)==[(1,1)]


def episode(seed, index):
    assert 80<=seed<=83 and 0<=index<16
    coefficient=2 if seed%2==0 else 1
    offset=(seed-80)%4
    sequence=[(index+seed)%4,((index//4)+3*seed)%4]
    while len(sequence)<32:
        sequence.append((sequence[-1]+coefficient*sequence[-2]+offset)%4)
    return tuple(sequence)


def samples(seed):
    rows=[]
    for ep in range(4):
        seq=episode(seed,ep)
        rows.append((seq[0],seq[1],seq[2]))
    assert len({(old,last) for old,last,_ in rows})==4
    return rows


def counts(rows):
    model=[[0]*4 for _ in range(16)]
    for before,last,label in rows:
        model[4*before+last][label]+=1
    assert sum(sum(row) for row in model)==4
    return model


def predict(arm, observed, table, selected):
    if arm=="fixed_predecessor":
        return (observed[-1]+observed[-2]+2)%4
    if arm=="rote":
        row=table[4*observed[-2]+observed[-1]]
        return max(range(4),key=lambda c:(row[c],-c))
    coefficient,offset=selected
    if arm=="shuffled_history":
        original=observed
        observed=original[::2]+original[1::2]
        assert sorted(observed)==sorted(original)
    else:
        assert arm=="selected_hypothesis"
    return (observed[-1]+coefficient*observed[-2]+offset)%4


def evaluate(seed,arm,table,chosen):
    confusion=[[0]*4 for _ in range(4)]
    signature=hashlib.sha256()
    for ep in range(4,16):
        seq=episode(seed,ep)
        for tick in range(2,32):
            pred=predict(arm,seq[:tick],table,chosen)
            true=seq[tick]
            confusion[true][pred]+=1
            signature.update(bytes((ep,tick,true,pred)))
    total=sum(sum(row) for row in confusion)
    assert total==360
    correct=sum(confusion[i][i] for i in range(4))
    return {"seed":seed,"arm":arm,"correct":correct,"total":total,
            "accuracy":correct/total,"confusion":confusion,
            "trace_sha256":signature.hexdigest()}


def run():
    smoke_fixture_only()
    arms=("fixed_predecessor","rote","selected_hypothesis","shuffled_history")
    selected={}
    ambiguous=[]
    observations=[]
    for seed in range(80,84):
        calibration=samples(seed)
        matches=select_hypotheses(calibration)
        if len(matches)!=1:
            ambiguous.append({"seed":seed,"matches":matches})
            continue
        selected[str(seed)]={"coefficient":matches[0][0],"offset":matches[0][1]}
        table=counts(calibration)
        observations.extend(evaluate(seed,arm,table,matches[0]) for arm in arms)
    if ambiguous:
        label="HYPOTHESIS_AMBIGUITY"
        results=[]
    else:
        results=[{
            "arm":arm,
            "heldout_accuracy":sum(row["correct"] for row in observations
                                   if row["arm"]==arm)/(4*360),
            "by_seed":[row for row in observations if row["arm"]==arm]
        } for arm in arms]
        metrics={row["arm"]:row["heldout_accuracy"] for row in results}
        if (metrics["selected_hypothesis"]>=.95 and
                all(metrics[a]<=.65 for a in arms if a!="selected_hypothesis")):
            label="BOUNDED_HYPOTHESIS_TRANSFER"
        elif metrics["selected_hypothesis"]>max(
                metrics[a] for a in arms if a!="selected_hypothesis"):
            label="PARTIAL"
        else:
            label="NEGATIVE"
    return {
        "schema":"wingless.w-r22-new-cohort-hypothesis.v1",
        "parent_head":"baa8e1eadde3fd10bba761735a50ff4d6944a82b",
        "failed_predecessor":"W-R21 invalid NameError before scientific output; no replay",
        "classification":label,
        "candidate_count":len(CANDIDATES),
        "selected":selected,
        "ambiguity":ambiguous,
        "arms":results,"model_calls":0,
        "claim_limit":"Eight hardcoded synthetic dynamics candidates identified from four labels; not unrestricted discovery."
    }


if __name__=="__main__":
    result=run()
    assert result==run()
    print(json.dumps(result,sort_keys=True,indent=2))
