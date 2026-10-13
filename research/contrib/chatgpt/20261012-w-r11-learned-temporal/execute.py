"""W-R11: fixed-width supervised temporal transfer pilot; synthetic only."""
import json


def original(seed, entity):
    q = (5 * entity + seed) % 9
    return tuple(
        ((192 + seed + 7 * q + 3 * t) % 256,
         (208 + 3 * seed + 5 * q + 5 * t) % 256,
         (224 + 5 * seed + 11 * q + 7 * t) % 256,
         (240 + 7 * seed + 13 * q + 11 * t) % 256)
        for t in range(16)
    )


def negative(sequence, kind):
    tail = sequence[1:]
    if kind == "reversal":
        return sequence[:1] + tail[::-1]
    assert kind == "rotation"
    return sequence[:1] + tail[5:] + tail[:5]


def representation(sequence, kind):
    if kind == "histogram":
        bins = [0.0] * 16
        for row in sequence:
            bins[row[0] % 16] += 1.0 / 16.0
        return bins
    if kind == "position":
        return [row[0] / 255.0 for row in sequence]
    assert kind == "transition"
    return [
        ((sequence[(i + 1) % 16][0] - sequence[i][0]) % 256) / 256.0
        for i in range(16)
    ]


def difference(a, b, kind):
    va = representation(a, kind)
    vb = representation(b, kind)
    assert len(va) == len(vb) == 16
    return [abs(x - y) for x, y in zip(va, vb)]


def labeled_pair(seed, entity, kind, positive, negative_kind):
    seq = original(seed, entity)
    other = seq if positive else negative(seq, negative_kind)
    assert seq[0] == other[0] and sorted(seq) == sorted(other)
    assert (seq == other) == positive
    return difference(seq, other, kind), 1 if positive else -1


def train(kind):
    examples = [
        labeled_pair(seed, entity, kind, positive, "reversal")
        for seed in range(4)
        for entity in range(9)
        for positive in (True, False)
    ]
    weights = [0.0] * 16
    bias = 0.0
    updates = 0
    for _ in range(32):
        for features, label in examples:
            score = bias + sum(w * value for w, value in zip(weights, features))
            if label * score <= 0.0:
                weights = [
                    w + 0.1 * label * value
                    for w, value in zip(weights, features)
                ]
                bias += 0.1 * label
                updates += 1
    return weights, bias, updates


def test(kind, trained):
    weights, bias, updates = trained
    rows = []
    for seed in range(4, 8):
        tp = tn = fp = fn = 0
        for entity in range(9):
            for positive in (True, False):
                feat, true_label = labeled_pair(
                    seed, entity, kind, positive, "rotation"
                )
                predicted_positive = (
                    bias + sum(w * x for w, x in zip(weights, feat))
                ) >= 0.0
                if true_label == 1 and predicted_positive:
                    tp += 1
                elif true_label == -1 and not predicted_positive:
                    tn += 1
                elif true_label == -1:
                    fp += 1
                else:
                    fn += 1
        assert tp + tn + fp + fn == 18
        rows.append({
            "seed": seed,
            "cohort": "heldout",
            "kind": kind,
            "true_positive": tp,
            "true_negative": tn,
            "false_positive": fp,
            "false_negative": fn,
            "accuracy": (tp + tn) / 18.0,
        })
    return {
        "kind": kind,
        "width": len(weights),
        "training_updates": updates,
        "heldout_per_seed": rows,
        "heldout_accuracy": sum(r["true_positive"] + r["true_negative"] for r in rows) / 72.0,
    }


def run():
    arms = [test(kind, train(kind)) for kind in ("histogram", "position", "transition")]
    metrics = {a["kind"]: a["heldout_accuracy"] for a in arms}
    if metrics["position"] >= 0.85 and metrics["transition"] >= 0.85 and metrics["histogram"] <= 0.55:
        classification = "ORDER_MODEL_TRANSFER"
    elif max(metrics["position"], metrics["transition"]) > metrics["histogram"]:
        classification = "PARTIAL_TRANSFER"
    else:
        classification = "NEGATIVE_TRANSFER"
    return {
        "schema": "wingless.w-r11-fixed-width-learned-temporal.v1",
        "predecessor_head": "6c106dff38e90e95aebd23d3ea592f26937b2f37",
        "models": arms,
        "classification": classification,
        "model_calls": 0,
        "claim_limit": "Supervised synthetic temporal pair discriminators, not natural-language or generalized intelligence",
    }


if __name__ == "__main__":
    result = run()
    assert result == run()
    print(json.dumps(result, sort_keys=True, indent=2))
