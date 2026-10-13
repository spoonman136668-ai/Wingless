"""W-R22: new prospective cohort, fixed 8-candidate hypothesis search."""
import hashlib
import json

CANDIDATES = tuple((a, b) for a in (1, 2) for b in (0, 1, 2, 3))
assert len(CANDIDATES) == 8


def transition(before, previous, a, b):
    return (previous + a * before + b) % 4


def fixture_preflight():
    # Artificial fixture, not from a development or heldout seed.
    pairs = [(0, 1), (1, 1), (2, 1), (3, 1)]
    examples = [(x, y, transition(x, y, 2, 3)) for x, y in pairs]
    assert matching_hypotheses(examples) == [(2, 3)]
    assert matching_hypotheses([(0, 1, 0), (1, 1, 0)]) == []
    return "PASS"


def sequence(seed, episode):
    assert 80 <= seed <= 83 and 0 <= episode < 16
    a = 2 if seed % 2 == 0 else 1
    b = (seed - 80) % 4
    values = [(episode + seed) % 4, ((episode // 4) + 3 * seed) % 4]
    while len(values) < 32:
        values.append(transition(values[-2], values[-1], a, b))
    return tuple(values)


def matching_hypotheses(observations):
    return [
        (a, b) for a, b in CANDIDATES
        if all(transition(before, previous, a, b) == true
               for before, previous, true in observations)
    ]


def four_examples(seed):
    examples = []
    for episode in range(4):
        seq = sequence(seed, episode)
        examples.append((seq[0], seq[1], seq[2]))
    assert len({(x, y) for x, y, label in examples}) == 4
    return examples


def train_rote(observations):
    rows = [[0] * 4 for _ in range(16)]
    for before, previous, true in observations:
        rows[4 * before + previous][true] += 1
    assert sum(map(sum, rows)) == 4
    return rows


def guess(kind, history, selected, table):
    if kind == "fixed_previous":
        return transition(history[-2], history[-1], 1, 2)
    if kind == "four_demo_rote":
        row = table[4 * history[-2] + history[-1]]
        return max(range(4), key=lambda value: (row[value], -value))
    a, b = selected
    if kind == "prefix_shuffled":
        history = history[::2] + history[1::2]
    else:
        assert kind == "eight_hypothesis"
    return transition(history[-2], history[-1], a, b)


def evaluate(seed, kind, selected, table):
    confusion = [[0] * 4 for _ in range(4)]
    digest = hashlib.sha256()
    for episode in range(4, 16):
        observations = sequence(seed, episode)
        for tick in range(2, 32):
            history = observations[:tick]   # never includes current/future label
            predicted = guess(kind, history, selected, table)
            truth = observations[tick]
            confusion[truth][predicted] += 1
            digest.update(bytes((episode, tick, truth, predicted)))
    total = sum(map(sum, confusion))
    correct = sum(confusion[i][i] for i in range(4))
    assert total == 360
    return {
        "seed": seed,
        "correct": correct,
        "total": total,
        "accuracy": correct / total,
        "confusion": confusion,
        "trace_sha256": digest.hexdigest(),
    }


def run():
    assert fixture_preflight() == "PASS"
    names = ("fixed_previous", "four_demo_rote", "eight_hypothesis", "prefix_shuffled")
    results = {kind: [] for kind in names}
    selected_by_seed = {}
    ambiguous = []
    for seed in range(80, 84):
        examples = four_examples(seed)
        matches = matching_hypotheses(examples)
        if len(matches) != 1:
            ambiguous.append({"seed": seed, "matches": matches})
            continue
        selected = matches[0]
        selected_by_seed[str(seed)] = {"a": selected[0], "b": selected[1]}
        table = train_rote(examples)
        for kind in names:
            results[kind].append(evaluate(seed, kind, selected, table))
    if ambiguous:
        classification = "HYPOTHESIS_AMBIGUITY"
        scores = {}
    else:
        scores = {
            kind: sum(x["correct"] for x in records) / 1440
            for kind, records in results.items()
        }
        if (scores["eight_hypothesis"] >= .95
                and all(scores[k] <= .65 for k in names if k != "eight_hypothesis")):
            classification = "BOUNDED_HYPOTHESIS_TRANSFER"
        elif scores["eight_hypothesis"] > max(
                scores[k] for k in names if k != "eight_hypothesis"):
            classification = "PARTIAL"
        else:
            classification = "NEGATIVE"
    return {
        "schema": "wingless.w-r22-frozen-candidate-search.v1",
        "predecessor_head": "baa8e1eadde3fd10bba761735a50ff4d6944a82b",
        "classification": classification,
        "fixture_preflight": "PASS",
        "candidates": len(CANDIDATES),
        "allowed_calibration_labels_per_seed": 4,
        "test_predictions_per_seed": 360,
        "selected": selected_by_seed,
        "ambiguous": ambiguous,
        "accuracy": scores,
        "arms": results,
        "model_calls": 0,
        "claim_limit": "Fixed eight-member modular hypothesis class, not unrestricted scientific discovery.",
    }


if __name__ == "__main__":
    outcome = run()
    assert outcome == run()
    print(json.dumps(outcome, sort_keys=True, indent=2))
