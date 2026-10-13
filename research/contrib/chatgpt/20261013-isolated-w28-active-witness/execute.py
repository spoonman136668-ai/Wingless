"""W-R28: prospective active label choice in a fixed 24-table model family.

Never expose a t2 calibration label to the learner before the corresponding
query. The adaptive acquisition score uses only unlabelled initial pairs
and model predictions.
"""
import hashlib
import importlib.util
import json
from pathlib import Path

PREVIOUS = Path(__file__).resolve().parent.parent / "20261013-w-r26-latent-alphabet" / "execute.py"
_spec = importlib.util.spec_from_file_location("w26_frozen", PREVIOUS)
assert _spec and _spec.loader
w26 = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(w26)

MODELS = tuple((p, w26.table_for(p)) for p in w26.PERMUTATIONS)
assert len(MODELS) == 24 and w26.fixture() == "PASS"


def permutation(seed):
    assert 128 <= seed <= 135
    digest = hashlib.sha256(("W28/" + str(seed)).encode("utf-8")).digest()
    return w26.PERMUTATIONS[digest[0] % 24]


def initial_pair(seed, episode):
    assert 0 <= episode < 16
    hidden = permutation(seed)
    return (hidden[(episode + seed) % 4],
            hidden[((episode // 4) + 3 * seed) % 4])


def observed_sequence(seed, episode):
    assert 0 <= episode < 16
    initial = [(episode + seed) % 4, ((episode // 4) + 3 * seed) % 4]
    while len(initial) < 32:
        before, last = initial[-2:]
        nxt = ((before + last) % 4 if seed % 2 == 0 else before ^ last)
        initial.append(nxt)
    perm = permutation(seed)
    return tuple(perm[x] for x in initial)


def unlabeled_pair_pool(seed):
    pairs = [initial_pair(seed, episode) for episode in range(16)]
    assert len(pairs) == len(set(pairs)) == 16
    return pairs


def acquisition(seed, arm):
    assert arm in ("passive", "active", "budget16")
    pool = unlabeled_pair_pool(seed)
    surviving = list(MODELS)
    queried = []
    trace = hashlib.sha256()
    decision = None
    while len(queried) < 16:
        if arm != "budget16":
            unique_tables = {table for _, table in surviving}
            if not surviving:
                decision = "REJECT"
                break
            if len(unique_tables) == 1:
                decision = "ACCEPT"
                break
        options = [episode for episode in range(16) if episode not in queried]
        if arm in ("passive", "budget16"):
            selected = options[0]
        else:
            def score(episode):
                before, last = pool[episode]
                ctx = 4 * before + last
                sizes = [sum(table[ctx] == value for _, table in surviving)
                         for value in range(4)]
                return (max(sizes), sum(n*n for n in sizes), episode)
            selected = min(options, key=score)
        before, last = pool[selected]
        rows = observed_sequence(seed, selected)
        assert (rows[0], rows[1]) == (before, last)
        target = rows[2]  # first access to the requested label, AFTER acquisition
        surviving = [(p, table) for p, table in surviving
                     if table[4 * before + last] == target]
        queried.append(selected)
        trace.update(bytes((selected, before, last, target, len(surviving))))
    tables = {table for _, table in surviving}
    if decision is None:
        decision = ("REJECT" if not surviving else
                    "ACCEPT" if len(tables) == 1 else "ABSTAIN")
    model = next(iter(tables)) if decision == "ACCEPT" else None
    return {
        "decision": decision,
        "queries": len(queried),
        "query_order": queried,
        "consistent_permutations": len(surviving),
        "distinct_predictive_tables": len(tables),
        "model": model,
        "query_sha256": trace.hexdigest(),
        "pool_sha256": hashlib.sha256(bytes(
            value for pair in pool for value in pair)).hexdigest(),
    }


def evaluate_accepted(seed, table):
    matrix = [[0]*4 for _ in range(4)]
    trace = hashlib.sha256()
    for episode in range(16):
        observations = observed_sequence(seed, episode)
        for tick in range(3, 32):
            context = 4*observations[tick-2] + observations[tick-1]
            guess = table[context]
            truth = observations[tick]
            matrix[truth][guess] += 1
            trace.update(bytes((episode, tick, truth, guess)))
    total = sum(map(sum, matrix))
    assert total == 464
    return {"correct": sum(matrix[i][i] for i in range(4)),
            "total": total, "confusion": matrix,
            "trace_sha256": trace.hexdigest()}


def run():
    entries = []
    for seed in range(128, 136):
        true_family = "ADD4" if seed % 2 == 0 else "XOR4"
        for arm in ("passive", "active", "budget16"):
            result = acquisition(seed, arm)
            fitted = result.pop("model")
            outcome = evaluate_accepted(seed, fitted) if fitted is not None else None
            entries.append({"seed": seed, "cohort": "development" if seed < 132
                            else "heldout", "family": true_family,
                            "arm": arm, **result, "test": outcome})
    heldout = [r for r in entries if r["seed"] >= 132]
    stats = {}
    for arm in ("passive", "active", "budget16"):
        rs = [r for r in heldout if r["arm"] == arm]
        add = [r for r in rs if r["family"] == "ADD4"]
        xor = [r for r in rs if r["family"] == "XOR4"]
        stats[arm] = {
            "add_accepted": sum(r["decision"] == "ACCEPT" for r in add),
            "add_perfect": sum(r["decision"] == "ACCEPT" and
                               r["test"]["correct"] == 464 for r in add),
            "xor_rejected": sum(r["decision"] == "REJECT" for r in xor),
            "xor_false_accepted": sum(r["decision"] == "ACCEPT" for r in xor),
            "add_queries": [r["queries"] for r in add],
            "xor_queries": [r["queries"] for r in xor],
        }
    false_accept = any(stats[x]["xor_false_accepted"] for x in ("passive", "active"))
    valid = all(stats[x]["add_perfect"] == 2 and stats[x]["xor_rejected"] == 2
                for x in ("passive", "active"))
    passive_mean = sum(stats["passive"]["add_queries"])/2
    active_mean = sum(stats["active"]["add_queries"])/2
    if false_accept:
        label = "FALSE_ACCEPT"
    elif valid and active_mean < passive_mean:
        label = "EFFICIENT_VALID_BOUNDARY"
    elif valid and active_mean <= passive_mean:
        label = "VALID_BOUNDARY_NO_EFFICIENCY"
    else:
        label = "MIXED_VALID"
    return {
        "schema": "wingless.w-r28-adaptive-witness.v1",
        "parent": "02a379212a8b28e434ad87131539abcfa89ef509",
        "classification": label,
        "heldout": stats,
        "rows": entries,
        "model_calls": 0,
        "claim_limit": "Fixed 24 finite conjugacy hypotheses with unlabeled initial-context acquisition; not unbounded discovery.",
    }


if __name__ == "__main__":
    original = run()
    assert original == run()
    print(json.dumps(original, sort_keys=True, indent=2))
