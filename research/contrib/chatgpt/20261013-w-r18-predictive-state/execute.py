"""W-R18 frozen finite order-two next-symbol prediction, equal-context controls."""
import hashlib
import json


def sequence(seed, episode):
    assert 48 <= seed <= 55 and 0 <= episode < 16
    values = [(episode + seed) % 4, ((episode // 4) + 3 * seed) % 4]
    while len(values) < 32:
        values.append((values[-1] + values[-2]) % 4)
    assert len(values) == 32 and all(0 <= v <= 3 for v in values)
    return tuple(values)


def context(values, episode, tick, kind):
    if kind == "order_two":
        return 4 * values[tick - 2] + values[tick - 1]
    if kind == "last_plus_episode":
        return 4 * values[tick - 1] + episode % 4
    assert kind == "shuffled_order_two"
    permuted = values[::2] + values[1::2]
    assert sorted(permuted) == sorted(values)
    return 4 * permuted[tick - 2] + permuted[tick - 1]


def examples(seed, kind):
    for episode in range(16):
        values = sequence(seed, episode)
        for tick in range(2, len(values)):
            ctx = context(values, episode, tick, kind)
            assert 0 <= ctx < 16
            yield (episode, tick, ctx, values[tick])


def fit(kind):
    table = [[0] * 4 for _ in range(16)]
    count = 0
    for seed in range(48, 52):
        for _, _, ctx, label in examples(seed, kind):
            table[ctx][label] += 1
            count += 1
    assert count == 4 * 16 * 30
    assert sum(sum(row) for row in table) == count
    seen = sum(bool(sum(row)) for row in table)
    return table, count, seen


def predict(table, ctx):
    return max(range(4), key=lambda label: (table[ctx][label], -label))


def evaluate(kind, fitted):
    table, trained_examples, covered = fitted
    heldout = []
    for seed in range(52, 56):
        correct = 0
        matrix = [[0] * 4 for _ in range(4)]
        digest = hashlib.sha256()
        for episode, tick, ctx, label in examples(seed, kind):
            guess = predict(table, ctx)
            matrix[label][guess] += 1
            correct += int(guess == label)
            digest.update(bytes((episode, tick, ctx, label, guess)))
        assert sum(sum(row) for row in matrix) == 16 * 30
        heldout.append({
            "seed": seed,
            "correct": correct,
            "examples": 16 * 30,
            "accuracy": correct / (16 * 30),
            "confusion": matrix,
            "prediction_trace_sha256": digest.hexdigest(),
        })
    return {
        "arm": kind,
        "model_capacity_integer_counts": 16 * 4,
        "contexts_covered": covered,
        "trained_examples": trained_examples,
        "heldout_accuracy": sum(r["correct"] for r in heldout) / (4 * 16 * 30),
        "heldout": heldout,
    }


def run():
    names = ("order_two", "last_plus_episode", "shuffled_order_two")
    arms = [evaluate(name, fit(name)) for name in names]
    scores = {a["arm"]: a["heldout_accuracy"] for a in arms}
    if (scores["order_two"] >= 0.95
            and scores["last_plus_episode"] <= 0.65
            and scores["shuffled_order_two"] <= 0.65):
        outcome = "ORDER_TWO_PREDICTIVE_TRANSFER"
    elif scores["order_two"] > max(
            scores["last_plus_episode"], scores["shuffled_order_two"]):
        outcome = "PARTIAL_TRANSFER"
    else:
        outcome = "NEGATIVE_TRANSFER"
    return {
        "schema": "wingless.w-r18-order-two-learned-predictor.v1",
        "predecessor_head": "21b3c3fee08e5af3ed6c22218d67f59ed40d60a9",
        "classification": outcome,
        "arms": arms,
        "model_calls": 0,
        "claim_limit": "Finite deterministic 4-symbol recurrence prediction only; no language or generalized intelligence.",
    }


if __name__ == "__main__":
    result = run()
    assert result == run()
    print(json.dumps(result, sort_keys=True, indent=2))
