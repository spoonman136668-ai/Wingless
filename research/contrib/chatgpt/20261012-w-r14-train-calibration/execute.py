"""W-R14: train-only calibration of circular temporal mismatch on fresh heldout."""
import importlib.util
import json
from pathlib import Path

PRIOR_SOURCE = Path(__file__).resolve().parent.parent / "20261012-w-r13-general-permutation" / "execute.py"
_spec = importlib.util.spec_from_file_location("w13_frozen", PRIOR_SOURCE)
assert _spec and _spec.loader
w13 = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(w13)
w12 = w13._prior


def pair(seed, entity, positive, phase):
    base = w12.original(seed, entity)
    shifted = w12.nuisance_shift(base, 17 if phase == "development" else 61)
    tail = shifted[1:]
    if positive:
        candidate = shifted
    elif phase == "development":
        candidate = shifted[:1] + (tail[::-1] if entity % 2 == 0 else tail[3:] + tail[:3])
    else:
        candidate = shifted[:1] + tail[7:] + tail[:7]
    assert base != shifted and shifted[0] == candidate[0]
    assert sorted(shifted) == sorted(candidate)
    assert (shifted == candidate) == positive
    return base, candidate


def train_only_scale():
    values = []
    for seed in range(16, 20):
        for entity in range(9):
            a, b = pair(seed, entity, False, "development")
            values.append(w13.pair_features(a, b, "pooled")[0])
    assert len(values) == 36
    minimum = min(values)
    assert minimum > 0.0
    return 1.0 / minimum, minimum


def vector(a, b, kind, scale):
    if kind == "histogram":
        f = w13.pair_features(a, b, "unordered_histogram")
        assert all(x == 0 for x in f)
    else:
        f = w13.pair_features(a, b, "pooled")
        if kind == "calibrated_pooled":
            f = [x * scale for x in f]
        else:
            assert kind == "unscaled_pooled"
    assert len(f) == 16
    return f


def train(kind, scale):
    records = []
    for seed in range(16, 20):
        for entity in range(9):
            for positive in (True, False):
                a, b = pair(seed, entity, positive, "development")
                records.append((vector(a, b, kind, scale), 1 if positive else -1))
    assert len(records) == 72
    weights = [0.0] * 16
    bias = 0.0
    updates = 0
    for _ in range(32):
        for values, target in records:
            score = bias + sum(w*x for w, x in zip(weights, values))
            if target * score <= 0:
                weights = [w + 0.1 * target * x for w, x in zip(weights, values)]
                bias += 0.1 * target
                updates += 1
    return weights, bias, updates


def evaluate(kind, trained, scale):
    weights, bias, updates = trained
    rows = []
    for seed in range(20, 24):
        tp = tn = fp = fn = 0
        for entity in range(9):
            for positive in (True, False):
                a, b = pair(seed, entity, positive, "heldout")
                values = vector(a, b, kind, scale)
                pred = 1 if bias + sum(w*x for w, x in zip(weights, values)) >= 0 else -1
                label = 1 if positive else -1
                if label == 1 and pred == 1: tp += 1
                elif label == -1 and pred == -1: tn += 1
                elif label == -1: fp += 1
                else: fn += 1
        assert tp+tn+fp+fn == 18
        rows.append({"seed": seed, "tp": tp, "tn": tn, "fp": fp, "fn": fn,
                     "accuracy": (tp+tn)/18})
    return {"kind": kind, "width": len(weights), "updates": updates,
            "heldout_accuracy": sum(r["tp"]+r["tn"] for r in rows)/72,
            "heldout_rows": rows}


def run():
    factor, development_min = train_only_scale()
    kinds = ("histogram", "unscaled_pooled", "calibrated_pooled")
    arms = [evaluate(kind, train(kind, factor), factor) for kind in kinds]
    score = {a["kind"]: a["heldout_accuracy"] for a in arms}
    if score["calibrated_pooled"] >= .85 and score["histogram"] <= .55:
        classification = "TRAIN_CALIBRATED_TRANSFER"
    elif score["calibrated_pooled"] > max(score["histogram"], score["unscaled_pooled"]):
        classification = "PARTIAL"
    else:
        classification = "NEGATIVE"
    return {
        "schema": "wingless.w-r14-train-only-calibrated.v1",
        "predecessor_head": "e9fa91c490316fcc486677e8c272d3cfeb9011a8",
        "classification": classification,
        "train_only_calibration_factor": factor,
        "train_only_min_negative_mismatch": development_min,
        "arms": arms,
        "model_calls": 0,
        "claim_limit": "Synthetic fixed-width train-only normalization; no general reasoning"
    }


if __name__ == "__main__":
    result = run()
    assert result == run()
    print(json.dumps(result, sort_keys=True, indent=2))
