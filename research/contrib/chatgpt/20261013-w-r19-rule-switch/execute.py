"""W-R19: fixed-capacity supervised few-shot shift of second-order symbol rule."""
import hashlib
import importlib.util
import json
from pathlib import Path

PREVIOUS = Path(__file__).resolve().parent.parent / "20261013-w-r18-predictive-state" / "execute.py"
_spec = importlib.util.spec_from_file_location("w18_frozen", PREVIOUS)
assert _spec and _spec.loader
w18 = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(w18)


def shifted_sequence(seed, episode):
    assert 56 <= seed <= 59 and 0 <= episode < 16
    seq = [(episode + seed) % 4, ((episode // 4) + 3 * seed) % 4]
    while len(seq) < 32:
        seq.append((seq[-1] + seq[-2] + 1) % 4)
    return tuple(seq)


def original_context(sequence, tick):
    return 4 * sequence[tick-2] + sequence[tick-1]


def shuffled_prefix_context(sequence, tick):
    observed = sequence[:tick]
    scrambled = observed[::2] + observed[1::2]
    assert len(observed) == len(scrambled) == tick
    assert sorted(observed) == sorted(scrambled)
    return 4 * scrambled[-2] + scrambled[-1]


def fit_new_rule(seed):
    table = [[0]*4 for _ in range(16)]
    for episode in range(16):
        seq = shifted_sequence(seed, episode)
        ctx = original_context(seq, 2)
        table[ctx][seq[2]] += 1
    assert sum(sum(row) for row in table) == 16
    assert all(sum(row) == 1 for row in table)
    return table


def evaluate_seed(seed, arm, trained_old, trained_new):
    if arm == "frozen_old_rule":
        table = trained_old
    else:
        table = trained_new
    cm = [[0]*4 for _ in range(4)]
    digest = hashlib.sha256()
    for episode in range(16):
        seq = shifted_sequence(seed, episode)
        for tick in range(3, 32):
            context = (shuffled_prefix_context(seq,tick)
                       if arm == "fewshot_temporal_shuffle"
                       else original_context(seq,tick))
            prediction = w18.predict(table, context)
            actual = seq[tick]
            cm[actual][prediction] += 1
            digest.update(bytes((episode, tick, context, actual, prediction)))
    total = sum(sum(row) for row in cm)
    correct = sum(cm[i][i] for i in range(4))
    assert total == 16 * 29
    return {
        "seed": seed,
        "arm": arm,
        "calibration_labels": 0 if arm == "frozen_old_rule" else 16,
        "test_examples": total,
        "correct": correct,
        "accuracy": correct / total,
        "confusion": cm,
        "prediction_trace_sha256": digest.hexdigest(),
    }


def run():
    old_table, historical_labels, old_coverage = w18.fit("order_two")
    assert historical_labels == 1920 and old_coverage == 16
    names = ("frozen_old_rule","fewshot_order_two","fewshot_temporal_shuffle")
    rows = []
    for seed in range(56, 60):
        adapted = fit_new_rule(seed)
        for name in names:
            rows.append(evaluate_seed(seed, name, old_table, adapted))
    arms = [{
        "arm": name,
        "model_capacity_integer_counts": 64,
        "heldout_accuracy": sum(r["correct"] for r in rows if r["arm"]==name)
                            / (4*16*29),
        "per_seed": [r for r in rows if r["arm"] == name],
    } for name in names]
    scores = {a["arm"]: a["heldout_accuracy"] for a in arms}
    if (scores["fewshot_order_two"] >= .95
            and scores["frozen_old_rule"] <= .40
            and scores["fewshot_temporal_shuffle"] <= .65):
        label = "ADAPTIVE_RULE_RECOVERY"
    elif scores["fewshot_order_two"] > max(
            scores["frozen_old_rule"], scores["fewshot_temporal_shuffle"]):
        label = "PARTIAL"
    else:
        label = "NEGATIVE"
    return {
        "schema":"wingless.w-r19-fewshot-rule-switch.v1",
        "predecessor_head":"452f0e37a449a9e68a1f10371582cbe7013ccb61",
        "classification":label,
        "historical_old_rule_labels":1920,
        "new_rule_calibration_labels_per_seed":16,
        "heldout_test_examples_per_seed":464,
        "arms":arms,
        "model_calls":0,
        "claim_limit":"All 16 new-rule contexts explicitly demonstrated once; finite supervised rule adaptation, not zero-shot dynamics transfer.",
    }


if __name__ == "__main__":
    result = run()
    assert result == run()
    print(json.dumps(result, sort_keys=True, indent=2))
