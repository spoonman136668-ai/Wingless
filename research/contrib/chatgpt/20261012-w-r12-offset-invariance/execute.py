"""W-R12: frozen shifted-positive temporal recognition without identical-pair shortcut."""
import json


def original(seed, entity):
    q = (5 * entity + seed) % 9
    return tuple(
        ((192 + seed + 7 * q + 3 * tick) % 256,
         (208 + 3 * seed + 5 * q + 5 * tick) % 256,
         (224 + 5 * seed + 11 * q + 7 * tick) % 256,
         (240 + 7 * seed + 13 * q + 11 * tick) % 256)
        for tick in range(16)
    )


def nuisance_shift(rows, delta):
    return tuple(tuple((value + delta) % 256 for value in row) for row in rows)


def corrupt_tail(rows, transform):
    tail = rows[1:]
    if transform == "reverse":
        return rows[:1] + tail[::-1]
    assert transform == "rotate_five"
    return rows[:1] + tail[5:] + tail[:5]


def circular_distance(a, b):
    delta = abs(a - b) % 256
    return min(delta, 256 - delta) / 128.0


def anchor_relative(rows):
    origin = rows[0][0]
    return tuple((row[0] - origin) % 256 for row in rows)


def features(a, b, kind):
    if kind == "raw_position":
        return [abs(x[0] - y[0]) / 255.0 for x, y in zip(a, b)]
    aa = anchor_relative(a)
    bb = anchor_relative(b)
    if kind == "unordered_relative_histogram":
        h1 = [0.0] * 16
        h2 = [0.0] * 16
        for x, y in zip(aa, bb):
            h1[x % 16] += 1.0 / 16.0
            h2[y % 16] += 1.0 / 16.0
        return [abs(x - y) for x, y in zip(h1, h2)]
    assert kind == "anchor_relative_circular_position"
    return [circular_distance(x, y) for x, y in zip(aa, bb)]


def labeled_pair(seed, entity, kind, positive, phase):
    base = original(seed, entity)
    shifted = nuisance_shift(base, 17 if phase == "development" else 31)
    observed = shifted if positive else corrupt_tail(
        shifted, "reverse" if phase == "development" else "rotate_five"
    )
    # Positive and negative cases share the SAME first row and multiset.
    assert base != shifted
    assert observed[0] == shifted[0]
    assert sorted(observed) == sorted(shifted)
    assert (observed == shifted) == positive
    assert anchor_relative(base) == anchor_relative(shifted)
    f = features(base, observed, kind)
    assert len(f) == 16
    if kind == "unordered_relative_histogram":
        assert all(v == 0 for v in f)
    if kind == "anchor_relative_circular_position":
        assert (all(v == 0 for v in f)) == positive
    return f, 1 if positive else -1


def train(kind):
    rows = [
        labeled_pair(seed, entity, kind, positive, "development")
        for seed in range(4)
        for entity in range(9)
        for positive in (True, False)
    ]
    weights = [0.0] * 16
    bias = 0.0
    updates = 0
    for _ in range(32):
        for vector, label in rows:
            signed_margin = label * (
                bias + sum(weight * value for weight, value in zip(weights, vector))
            )
            if signed_margin <= 0.0:
                weights = [
                    weight + 0.1 * label * value
                    for weight, value in zip(weights, vector)
                ]
                bias += 0.1 * label
                updates += 1
    return weights, bias, updates


def evaluate(kind, trained):
    weights, bias, updates = trained
    outcome = []
    for seed in range(4, 8):
        tp = tn = fp = fn = 0
        for entity in range(9):
            for positive in (True, False):
                vector, label = labeled_pair(seed, entity, kind, positive, "heldout")
                predicted = 1 if bias + sum(
                    weight * value for weight, value in zip(weights, vector)
                ) >= 0.0 else -1
                if label == 1 and predicted == 1:
                    tp += 1
                elif label == -1 and predicted == -1:
                    tn += 1
                elif label == -1:
                    fp += 1
                else:
                    fn += 1
        assert tp + tn + fp + fn == 18
        outcome.append({"seed": seed, "tp": tp, "tn": tn, "fp": fp, "fn": fn,
                        "accuracy": (tp + tn) / 18.0})
    return {
        "kind": kind,
        "features": len(weights),
        "training_updates": updates,
        "heldout_by_seed": outcome,
        "heldout_accuracy": sum(x["tp"] + x["tn"] for x in outcome) / 72.0,
    }


def run():
    kinds = ("raw_position", "unordered_relative_histogram",
             "anchor_relative_circular_position")
    results = [evaluate(kind, train(kind)) for kind in kinds]
    score = {row["kind"]: row["heldout_accuracy"] for row in results}
    ordered = score["anchor_relative_circular_position"]
    if ordered >= 0.85 and score["unordered_relative_histogram"] <= 0.55:
        classification = "NUISANCE_INVARIANT_TRANSFER"
    elif ordered > max(score["raw_position"], score["unordered_relative_histogram"]):
        classification = "PARTIAL"
    else:
        classification = "NEGATIVE"
    return {
        "schema": "wingless.w-r12-nuisance-shift.v1",
        "predecessor_head": "4f7da44d03544559bbb7753052811b93f7d6d823",
        "classification": classification,
        "models": results,
        "model_calls": 0,
        "claim_limit": "Narrow synthetic temporal recognition with frozen nuisance; no natural-language claims",
    }


if __name__ == "__main__":
    result = run()
    assert result == run()
    print(json.dumps(result, sort_keys=True, indent=2))
