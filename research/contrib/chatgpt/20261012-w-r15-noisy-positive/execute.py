"""W-R15: bounded jitter-positive temporal learning on fresh sealed cohort."""
import importlib.util
import json
from pathlib import Path

PRIOR_SOURCE = Path(__file__).resolve().parent.parent / "20261012-w-r14-train-calibration" / "execute.py"
_spec = importlib.util.spec_from_file_location("w14_frozen", PRIOR_SOURCE)
assert _spec and _spec.loader
w14 = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(w14)
w13 = w14.w13
w12 = w14.w12


def candidate_pair(seed, entity, positive, cohort):
    base = w12.original(seed, entity)
    delta = 17 if cohort == "development" else 61
    modulus = 3 if cohort == "development" else 5
    shifted = w12.nuisance_shift(base, delta)
    observed = tuple(
        tuple((value + ((3*t + entity + seed) % modulus) - (modulus//2)) % 256
              for value in row)
        for t, row in enumerate(shifted)
    )
    assert len(observed) == 16
    if positive:
        answer = observed
    elif cohort == "development":
        tail = observed[1:]
        answer = observed[:1] + (
            tail[::-1] if entity % 2 == 0 else tail[3:] + tail[:3]
        )
    else:
        tail = observed[1:]
        answer = observed[:1] + tail[::2] + tail[1::2]
    assert answer[0] == observed[0]
    assert sorted(answer) == sorted(observed)
    assert (answer == observed) == positive
    return base, answer


def calibration_factor():
    negative = []
    for seed in range(24, 28):
        for entity in range(9):
            a, b = candidate_pair(seed, entity, False, "development")
            negative.append(w13.pair_features(a, b, "pooled")[0])
    minimum = min(negative)
    assert minimum > 0
    return 1.0 / minimum, minimum


def features(a, b, kind, factor):
    if kind == "unordered":
        result = w13.pair_features(a, b, "unordered_histogram")
    elif kind == "pooled":
        result = w13.pair_features(a, b, "pooled")
    else:
        assert kind == "calibrated_pooled"
        result = [v * factor for v in w13.pair_features(a, b, "pooled")]
    assert len(result) == 16
    return result


def train(kind, factor):
    cases = []
    for seed in range(24, 28):
        for entity in range(9):
            for positive in (True, False):
                a, b = candidate_pair(seed, entity, positive, "development")
                cases.append((features(a, b, kind, factor),
                              1 if positive else -1))
    assert len(cases) == 72
    weights = [0.0] * 16
    bias = 0.0
    updates = 0
    for _ in range(32):
        for vec, label in cases:
            score = bias + sum(w*v for w, v in zip(weights, vec))
            if label * score <= 0:
                weights = [w + 0.1 * label * v for w, v in zip(weights, vec)]
                bias += 0.1 * label
                updates += 1
    return weights, bias, updates


def evaluate(kind, trained, factor):
    weights, bias, updates = trained
    rows = []
    for seed in range(28, 32):
        tp = tn = fp = fn = 0
        for entity in range(9):
            for positive in (True, False):
                a, b = candidate_pair(seed, entity, positive, "heldout")
                vec = features(a, b, kind, factor)
                pred = 1 if bias + sum(w*v for w,v in zip(weights,vec)) >= 0 else -1
                label = 1 if positive else -1
                if pred == 1 and label == 1: tp += 1
                elif pred == -1 and label == -1: tn += 1
                elif pred == 1: fp += 1
                else: fn += 1
        assert tp + tn + fp + fn == 18
        rows.append({"seed": seed, "tp": tp, "tn": tn, "fp": fp, "fn": fn,
                     "accuracy": (tp+tn)/18})
    return {"kind": kind, "width": len(weights), "updates": updates,
            "heldout_accuracy": sum(x["tp"]+x["tn"] for x in rows)/72,
            "by_seed": rows}


def run():
    factor, min_negative = calibration_factor()
    rows = [evaluate(kind, train(kind,factor),factor)
            for kind in ("unordered", "pooled", "calibrated_pooled")]
    scores = {row["kind"]: row["heldout_accuracy"] for row in rows}
    best = max(scores["pooled"], scores["calibrated_pooled"])
    if best >= 0.85 and scores["unordered"] <= 0.55:
        classification = "ROBUST_TEMPORAL_TRANSFER"
    elif best > scores["unordered"]:
        classification = "PARTIAL"
    else:
        classification = "NEGATIVE"
    return {
        "schema": "wingless.w-r15-jitter-positive.v1",
        "predecessor_head": "690b08efaffb921835525029c3b5504c30272279",
        "classification": classification,
        "train_calibration_factor": factor,
        "train_min_negative_mismatch": min_negative,
        "arms": rows,
        "model_calls": 0,
        "claim_limit": "Frozen synthetic jittered pairs, no natural language or general intelligence"
    }


if __name__ == "__main__":
    result = run()
    assert result == run()
    print(json.dumps(result, sort_keys=True, indent=2))
