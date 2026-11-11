"""W-R23: frozen modular open-set hypothesis rejection on fresh inputs."""
import hashlib
import importlib.util
import json
from pathlib import Path

PREVIOUS = Path(__file__).resolve().parent.parent / "20261013-w-r22-new-hypothesis" / "execute.py"
spec = importlib.util.spec_from_file_location("w22_frozen", PREVIOUS)
assert spec and spec.loader
w22 = importlib.util.module_from_spec(spec)
spec.loader.exec_module(w22)
CANDIDATES = w22.CANDIDATES


def sequence(seed, episode):
    assert 88 <= seed <= 95 and 0 <= episode < 16
    a = 2 if seed % 2 == 0 else 3
    b = (seed - 88) % 4
    seq = [(episode + seed) % 4, ((episode // 4) + 3 * seed) % 4]
    while len(seq) < 32:
        seq.append(w22.transition(seq[-2], seq[-1], a, b))
    return tuple(seq)


def fixture_only():
    x = [(before, 1) for before in range(4)]
    inside = [(before, previous, w22.transition(before, previous, 2, 2))
              for before, previous in x]
    outside = [(before, previous, w22.transition(before, previous, 3, 2))
               for before, previous in x]
    assert w22.matching_hypotheses(inside) == [(2, 2)]
    assert w22.matching_hypotheses(outside) == []
    return "PASS"


def observations(seed):
    cases = []
    for episode in range(4):
        seq = sequence(seed, episode)
        cases.append((seq[0], seq[1], seq[2]))
    assert len(cases) == 4 and len({(a, b) for a, b, _ in cases}) == 4
    return cases


def score_candidate(candidate, triples):
    a, b = candidate
    return sum(w22.transition(before, previous, a, b) != actual
               for before, previous, actual in triples)


def evaluate(seed, model):
    cm = [[0] * 4 for _ in range(4)]
    digest = hashlib.sha256()
    for episode in range(4, 16):
        seq = sequence(seed, episode)
        for tick in range(2, 32):
            a, b = model
            predicted = w22.transition(seq[tick-2], seq[tick-1], a, b)
            actual = seq[tick]
            cm[actual][predicted] += 1
            digest.update(bytes((episode, tick, actual, predicted)))
    assert sum(map(sum, cm)) == 360
    correct = sum(cm[i][i] for i in range(4))
    return {"correct": correct, "total": 360,
            "accuracy": correct / 360, "confusion": cm,
            "trace_sha256": digest.hexdigest()}


def run():
    assert fixture_only() == "PASS"
    results = []
    for seed in range(88, 96):
        triples = observations(seed)
        errors = [(score_candidate(c, triples), c) for c in CANDIDATES]
        best_error = min(value for value, _ in errors)
        best_models = [c for value, c in errors if value == best_error]
        forced = min(best_models)
        perfectly_matching = [c for value, c in errors if value == 0]
        accepted = len(perfectly_matching) == 1
        gate_result = evaluate(seed, perfectly_matching[0]) if accepted else None
        forced_result = evaluate(seed, forced)
        results.append({
            "seed": seed,
            "truth_bucket": "in_class" if seed % 2 == 0 else "out_of_class",
            "best_calibration_errors": best_error,
            "matching_models": perfectly_matching,
            "accepted": accepted,
            "forced_model": {"a": forced[0], "b": forced[1]},
            "accepted_outcome": gate_result,
            "forced_outcome": forced_result,
        })
    known = [r for r in results if r["truth_bucket"] == "in_class"]
    unknown = [r for r in results if r["truth_bucket"] == "out_of_class"]
    wrong_actionable = sum(r["accepted_outcome"]["total"] - r["accepted_outcome"]["correct"]
                           for r in results if r["accepted"])
    if any(r["accepted"] for r in unknown):
        classification = "FALSE_ACCEPT"
    elif any(not r["accepted"] for r in known):
        classification = "FALSE_REJECT"
    elif wrong_actionable == 0 and all(r["accepted"] for r in known):
        classification = "OPEN_SET_REJECTION"
    else:
        classification = "NEGATIVE"
    return {
        "schema": "wingless.w-r23-open-set-hypothesis-gate.v1",
        "predecessor_head": "2ebe681a4ffa159251ef9e96a233ccc261f37e70",
        "classification": classification,
        "fixture_only": "PASS",
        "candidate_count": len(CANDIDATES),
        "total_heldout_positions": 8 * 360,
        "accepted_positions": sum(360 for r in results if r["accepted"]),
        "wrong_actionable_predictions": wrong_actionable,
        "in_class_accepted": sum(r["accepted"] for r in known),
        "out_of_class_rejected": sum(not r["accepted"] for r in unknown),
        "forced_in_class_accuracy": sum(r["forced_outcome"]["correct"] for r in known) / (4*360),
        "forced_out_class_accuracy": sum(r["forced_outcome"]["correct"] for r in unknown) / (4*360),
        "rows": results, "model_calls": 0,
        "claim_limit": "Exact rejection only for one excluded synthetic modular coefficient family.",
    }


if __name__ == "__main__":
    result = run()
    assert result == run()
    print(json.dumps(result, sort_keys=True, indent=2))
