"""W-R16 frozen factorial evaluation: genuine time-varying noise and novel order changes."""
import hashlib
import json


TRAIN_SEEDS = tuple(range(32, 48))
HELDOUT_SEEDS = tuple(range(48, 64))
TRAIN_FAMILIES = ("reverse_tail", "rotate3", "weave_tail")
TEST_FAMILIES = ("rotate7", "swap_tail_blocks", "pairwise_swap_tail")
TRAIN_AMPS = (1, 2)
TEST_AMPS = (1, 2, 3)
ARMS = ("unordered_histogram", "circular_position", "mean_mismatch", "transition")


def original(seed, entity):
    q = (5 * entity + seed) % 9
    return tuple(
        ((192 + seed + 7*q + 3*t) % 256,
         (208 + 3*seed + 5*q + 5*t) % 256,
         (224 + 5*seed + 11*q + 7*t) % 256,
         (240 + 7*seed + 13*q + 11*t) % 256)
        for t in range(16)
    )


def observed(base, seed, entity, amplitude, offset):
    assert amplitude in (1, 2, 3)
    # Multiplier 11 is coprime to 3/5/7: jitter genuinely varies by tick.
    return tuple(
        tuple((value + offset
               + (11*t + 7*channel + 3*entity + seed) % (2*amplitude+1)
               - amplitude) % 256
              for channel, value in enumerate(row))
        for t, row in enumerate(base)
    )


def permute(sequence, family):
    tail = sequence[1:]
    if family == "reverse_tail":
        transformed = tail[::-1]
    elif family == "rotate3":
        transformed = tail[3:] + tail[:3]
    elif family == "weave_tail":
        transformed = tail[::2] + tail[1::2]
    elif family == "rotate7":
        transformed = tail[7:] + tail[:7]
    elif family == "swap_tail_blocks":
        # Exchange the first seven with the following four; preserve the final four.
        transformed = tail[7:11] + tail[:7] + tail[11:]
    else:
        assert family == "pairwise_swap_tail"
        transformed = tuple(tail[i+1] if i%2==0 and i+1<len(tail)
                            else tail[i-1] if i%2==1 else tail[i]
                            for i in range(len(tail)))
    out = sequence[:1] + transformed
    assert out[0] == sequence[0]
    assert sorted(out) == sorted(sequence)
    assert out != sequence
    return out


def pair(seed, entity, amplitude, offset, negative_family, positive):
    base = original(seed, entity)
    candidate = observed(base, seed, entity, amplitude, offset)
    assert len(set(row[0] for row in candidate)) == 16
    # Verify jitter values across ticks are not constant as in W-R15.
    jitter = [((11*t + 3*entity + seed) % (2*amplitude+1)) - amplitude
              for t in range(16)]
    assert len(set(jitter)) == 2*amplitude+1
    other = candidate if positive else permute(candidate, negative_family)
    assert candidate[0] == other[0] and sorted(candidate) == sorted(other)
    assert (candidate == other) == positive
    return base, other


def circular_distance(a, b):
    delta = (a-b) % 256
    return min(delta, 256-delta) / 128.0


def representation(base, other, kind):
    a = tuple((row[0]-base[0][0]) % 256 for row in base)
    b = tuple((row[0]-other[0][0]) % 256 for row in other)
    if kind == "unordered_histogram":
        ha = [0.0]*16
        hb = [0.0]*16
        for x, y in zip(a, b):
            ha[x%16] += 1.0/16
            hb[y%16] += 1.0/16
        features = [abs(x-y) for x,y in zip(ha,hb)]
    elif kind == "circular_position":
        features = [circular_distance(x,y) for x,y in zip(a,b)]
    elif kind == "mean_mismatch":
        mean = sum(circular_distance(x,y) for x,y in zip(a,b))/16
        features = [mean]*16
    else:
        assert kind == "transition"
        da = [((a[(i+1)%16]-a[i])%256) for i in range(16)]
        db = [((b[(i+1)%16]-b[i])%256) for i in range(16)]
        features = [circular_distance(x,y) for x,y in zip(da,db)]
    assert len(features) == 16 and all(0<=x<=1 for x in features)
    return features


def train(kind):
    examples = []
    for seed in TRAIN_SEEDS:
        for entity in range(9):
            for family in TRAIN_FAMILIES:
                for amp in TRAIN_AMPS:
                    for positive in (True,False):
                        a,b = pair(seed,entity,amp,17,family,positive)
                        features = representation(a,b,kind)
                        examples.append((features,1 if positive else -1))
    assert len(examples)==1728
    weights = [0.0]*16
    bias = 0.0
    updates = 0
    for _ in range(32):
        for features,label in examples:
            margin = label*(bias+sum(w*v for w,v in zip(weights,features)))
            if margin<=0:
                weights = [w + .1*label*v for w,v in zip(weights,features)]
                bias += .1*label
                updates += 1
    model_digest = hashlib.sha256(json.dumps([weights,bias],sort_keys=True).encode()).hexdigest()
    return weights,bias,updates,model_digest


def evaluate(kind,model,amp,family):
    weights,bias,updates,digest = model
    seeds = []
    for seed in HELDOUT_SEEDS:
        tp=tn=fp=fn=0
        for entity in range(9):
            positive_base,positive_candidate = pair(seed,entity,amp,61,family,True)
            negative_base,negative_candidate = pair(seed,entity,amp,61,family,False)
            assert positive_base == negative_base
            assert sorted(positive_candidate)==sorted(negative_candidate)
            assert positive_candidate[0] == negative_candidate[0]
            if kind=="unordered_histogram":
                assert representation(positive_base,positive_candidate,kind)==representation(negative_base,negative_candidate,kind)
            for positive,(base,candidate) in ((True,(positive_base,positive_candidate)),
                                             (False,(negative_base,negative_candidate))):
                features = representation(base,candidate,kind)
                pred = (bias+sum(w*v for w,v in zip(weights,features)))>=0
                if positive and pred:tp+=1
                elif not positive and not pred:tn+=1
                elif not positive:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        seeds.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn})
    good = sum(s["tp"]+s["tn"] for s in seeds)
    return {"amplitude":amp,"negative_family":family,"correct":good,
            "total":len(seeds)*18,"accuracy":good/(len(seeds)*18),
            "by_seed":seeds}


def run():
    result = []
    for kind in ARMS:
        model = train(kind)
        cells = [evaluate(kind,model,amp,family)
                 for amp in TEST_AMPS for family in TEST_FAMILIES]
        total = sum(x["total"] for x in cells)
        correct = sum(x["correct"] for x in cells)
        result.append({"model":kind,"width":16,"training_examples":1728,
                       "epochs":32,"updates":model[2],"trained_model_sha256":model[3],
                       "aggregate_accuracy":correct/total,"matrix":cells})
    scores = {x["model"]:x for x in result}
    best_any = max(scores[k]["aggregate_accuracy"] for k in ARMS[1:])
    strict = any(
        all(c["accuracy"]>=.85 for c in scores[k]["matrix"])
        for k in ARMS[1:]
    ) and scores["unordered_histogram"]["aggregate_accuracy"]<=.60
    classification = ("STRONG_FAMILY_TRANSFER" if strict else
                      "PARTIAL" if best_any>scores["unordered_histogram"]["aggregate_accuracy"]
                      else "NEGATIVE")
    return {
        "schema":"wingless.w-r16-preregistered-noise-matrix.v1",
        "parent_head":"cd9ed8c98e655e9ebf77204750f9ca733ced83e5",
        "classification":classification,"model_calls":0,
        "development_seeds":TRAIN_SEEDS,"sealed_heldout_seeds":HELDOUT_SEEDS,
        "models":result,
        "claim_limit":"Finite supervised synthetic sequence difference discrimination only"
    }


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
