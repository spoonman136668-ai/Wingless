"""W-R18: train-only scalar discrimination of unseen adjacent swaps with true jitter."""
import json

SEEDS_DEV = range(48, 52)
SEEDS_HELD = range(52, 56)
ARMS = ("unordered_histogram", "mean_position", "max_position", "transition")


def original(seed, entity):
    q = (5 * entity + seed) % 9
    return tuple(
        (
            (192 + seed + 7 * q + 3 * t + t*t) % 256,
            (208 + 3*seed + 5*q + 5*t + t*t) % 256,
            (224 + 5*seed + 11*q + 7*t + 2*t*t) % 256,
            (240 + 7*seed + 13*q + 11*t + 3*t*t) % 256,
        )
        for t in range(16)
    )


def make_pair(seed, entity, positive, phase):
    base = original(seed, entity)
    dev = phase == "development"
    offset = 17 if dev else 61
    modulus = 3 if dev else 5
    coefficient = 5 if dev else 7
    jitters = [((coefficient*t + entity + seed) % modulus) - modulus//2
               for t in range(16)]
    assert len(set(jitters)) == modulus
    candidate = tuple(tuple((x + offset + jitters[t]) % 256 for x in row)
                      for t, row in enumerate(base))
    if positive:
        observed = candidate
    elif dev:
        tail = candidate[1:]
        observed = candidate[:1] + (
            tail[::-1] if entity % 2 == 0 else tail[3:] + tail[:3])
    else:
        mutable = list(candidate)
        mutable[7], mutable[8] = mutable[8], mutable[7]
        observed = tuple(mutable)
    assert candidate[0] == observed[0]
    assert sorted(candidate) == sorted(observed)
    assert (candidate == observed) == positive
    assert observed != base
    return base, observed


def circular_distance(a, b):
    delta = abs(a-b) % 256
    return min(delta, 256-delta) / 128.0


def anchor_relative(seq):
    return [(row[0] - seq[0][0]) % 256 for row in seq]


def score(a, b, kind):
    x = anchor_relative(a)
    y = anchor_relative(b)
    if kind == "unordered_histogram":
        h1 = [0]*16
        h2 = [0]*16
        for u,v in zip(x,y):
            h1[u%16] += 1
            h2[v%16] += 1
        return sum(abs(u-v) for u,v in zip(h1,h2)) / 32.0
    positional = [circular_distance(u,v) for u,v in zip(x,y)]
    if kind == "mean_position":
        return sum(positional)/16
    if kind == "max_position":
        return max(positional)
    assert kind == "transition"
    edges_x = [(x[i+1] - x[i])%256 for i in range(15)]
    edges_y = [(y[i+1] - y[i])%256 for i in range(15)]
    return sum(circular_distance(u,v) for u,v in zip(edges_x,edges_y)) / 15


def train(kind):
    positives = []
    negatives = []
    for seed in SEEDS_DEV:
        for entity in range(9):
            for positive in (True, False):
                a,b = make_pair(seed,entity,positive,"development")
                observed = score(a,b,kind)
                (positives if positive else negatives).append(observed)
    assert len(positives)==len(negatives)==36
    maxp=max(positives)
    minn=min(negatives)
    threshold=(maxp+minn)/2
    return {"threshold":threshold,"training_max_positive":maxp,
            "training_min_negative":minn,"train_separable":maxp<minn}


def test(kind, trained):
    rows=[]
    for seed in SEEDS_HELD:
        tp=tn=fp=fn=0
        for entity in range(9):
            for positive in (True,False):
                a,b=make_pair(seed,entity,positive,"heldout")
                pred=score(a,b,kind)<=trained["threshold"]
                if positive and pred:tp+=1
                elif not positive and not pred:tn+=1
                elif not positive:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        rows.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,
                     "accuracy":(tp+tn)/18})
    return {"arm":kind,"learned_scalars":1,"train":trained,
            "heldout_accuracy":sum(x["tp"]+x["tn"] for x in rows)/72,
            "heldout":rows}


def run():
    arms=[test(kind,train(kind)) for kind in ARMS]
    scores={a["arm"]:a["heldout_accuracy"] for a in arms}
    if scores["transition"]>=.85 and scores["unordered_histogram"]<=.55:
        label="LOCAL_ORDER_TRANSFER"
    elif scores["transition"]>max(scores["mean_position"],
                                  scores["unordered_histogram"]):
        label="PARTIAL"
    else:
        label="NEGATIVE"
    return {
        "schema":"wingless.w-r18-local-order-transfer.v1",
        "predecessor_head":"21b3c3fee08e5af3ed6c22218d67f59ed40d60a9",
        "classification":label,
        "model_calls":0,"arms":arms,
        "scope":"New nonlinear toy sensor family, truly varying jitter; frozen unseen local swap"
    }


if __name__=="__main__":
    outcome=run()
    assert outcome==run()
    print(json.dumps(outcome,sort_keys=True,indent=2))
