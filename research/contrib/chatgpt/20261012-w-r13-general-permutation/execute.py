"""W-R13: fresh-cohort order generalization with constant-capacity learned features."""
import importlib.util
import json
from pathlib import Path

PRIOR_SOURCE = Path(__file__).resolve().parent.parent / "20261012-w-r12-offset-invariance" / "execute.py"
_spec = importlib.util.spec_from_file_location("w12_frozen_source", PRIOR_SOURCE)
assert _spec and _spec.loader
_prior = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(_prior)


def transform(sequence, phase, entity):
    tail = sequence[1:]
    if phase == "development":
        return sequence[:1] + (tail[::-1] if entity % 2 == 0 else tail[3:] + tail[:3])
    assert phase == "heldout"
    # Deliberately unseen permutation of the 15-element tail.
    woven = tail[0::2] + tail[1::2]
    assert len(woven) == 15 and sorted(woven) == sorted(tail) and woven != tail
    return sequence[:1] + woven


def pair_features(original, candidate, kind):
    a = _prior.anchor_relative(original)
    b = _prior.anchor_relative(candidate)
    if kind == "unordered_histogram":
        ha = [0.0] * 16
        hb = [0.0] * 16
        for x, y in zip(a, b):
            ha[x % 16] += 1.0 / 16
            hb[y % 16] += 1.0 / 16
        return [abs(x - y) for x, y in zip(ha, hb)]
    mismatch = [_prior.circular_distance(x, y) for x, y in zip(a, b)]
    if kind == "position":
        return mismatch
    assert kind == "pooled"
    pooled = sum(mismatch) / 16
    return [pooled] * 16


def example(seed, entity, kind, positive, phase):
    base = _prior.original(seed, entity)
    shifted = _prior.nuisance_shift(base, 17 if phase == "development" else 47)
    candidate = shifted if positive else transform(shifted, phase, entity)
    assert base != shifted and shifted[0] == candidate[0]
    assert sorted(shifted) == sorted(candidate)
    assert (shifted == candidate) == positive
    f = pair_features(base, candidate, kind)
    assert len(f) == 16
    if kind == "unordered_histogram":
        assert not any(f)
    else:
        assert (not any(f)) == positive
    return f, 1 if positive else -1


def train(kind):
    data = [
        example(seed, entity, kind, positive, "development")
        for seed in range(8, 12)
        for entity in range(9)
        for positive in (True, False)
    ]
    assert len(data) == 72
    weights = [0.0] * 16
    bias = 0.0
    updates = 0
    for _ in range(32):
        for vec, target in data:
            val = bias + sum(x*y for x, y in zip(weights, vec))
            if target * val <= 0:
                weights = [x + 0.1 * target * y for x, y in zip(weights, vec)]
                bias += 0.1 * target
                updates += 1
    return weights, bias, updates


def evaluate(kind, model):
    weights, bias, updates = model
    rows = []
    for seed in range(12, 16):
        tp = tn = fp = fn = 0
        for entity in range(9):
            for positive in (True, False):
                vec, truth = example(seed, entity, kind, positive, "heldout")
                pred = 1 if bias + sum(x*y for x, y in zip(weights, vec)) >= 0 else -1
                if pred == 1 and truth == 1:
                    tp += 1
                elif pred == -1 and truth == -1:
                    tn += 1
                elif pred == 1:
                    fp += 1
                else:
                    fn += 1
        assert tp + tn + fp + fn == 18
        rows.append({"seed": seed, "tp": tp, "tn": tn, "fp": fp, "fn": fn,
                     "accuracy": (tp + tn) / 18})
    return {"arm": kind, "feature_width": len(weights), "updates": updates,
            "heldout_accuracy": sum(x["tp"]+x["tn"] for x in rows) / 72,
            "by_seed": rows}


def run():
    arms = [evaluate(kind, train(kind)) for kind in
            ("unordered_histogram", "position", "pooled")]
    scores = {a["arm"]: a["heldout_accuracy"] for a in arms}
    if scores["pooled"] >= 0.85 and scores["unordered_histogram"] <= 0.55:
        label = "TRANSFER"
    elif scores["pooled"] > scores["unordered_histogram"]:
        label = "PARTIAL"
    else:
        label = "NEGATIVE"
    return {
        "schema": "wingless.w-r13-unseen-weave-transfer.v1",
        "predecessor_head": "5e27cf50f70bfb8fb28364461e2cac0f0695a055",
        "classification": label,
        "arms": arms,
        "model_calls": 0,
        "claim_limit": "Supervised synthetic temporal order mismatch, not general reasoning"
    }


if __name__ == "__main__":
    data = run()
    assert data == run()
    print(json.dumps(data, sort_keys=True, indent=2))
