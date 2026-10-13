"""W-R29: independent unused-context witness challenge after finite-model selection."""
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
CHALLENGE_BUDGETS = (0, 2, 4, 8, 16)


def observed(seed, episode):
    assert 136 <= seed <= 143 and 0 <= episode < 16
    latent = [(episode + seed) % 4, ((episode // 4) + 3 * seed) % 4]
    while len(latent) < 32:
        before, last = latent[-2:]
        latent.append((before + last) % 4 if seed % 2 == 0 else before ^ last)
    permutation = w26.PERMUTATIONS[hashlib.sha256(
        ("W29/" + str(seed)).encode()).digest()[0] % 24]
    return tuple(permutation[x] for x in latent)


def pool(seed):
    # Unlabelled contexts only; labels are not used during query selection.
    pairs = [(observed(seed, ep)[0], observed(seed, ep)[1]) for ep in range(16)]
    assert len(set(pairs)) == 16
    return pairs


def choose_next(candidates, remaining, pairs):
    def score(ep):
        a, b = pairs[ep]
        ctx = 4*a+b
        sizes = [sum(table[ctx] == val for _, table in candidates)
                 for val in range(4)]
        return (max(sizes), sum(size*size for size in sizes), ep)
    return min(remaining, key=score)


def selected_model(seed, pairs):
    compatible = list(MODELS)
    used = []
    trace = hashlib.sha256()
    while len(used) < 16 and compatible:
        tables = {table for _, table in compatible}
        if len(tables) == 1:
            return "ACCEPT", next(iter(tables)), used, trace.hexdigest(), len(compatible)
        ep = choose_next(compatible, [i for i in range(16) if i not in used], pairs)
        a, b = pairs[ep]
        target = observed(seed, ep)[2]  # target seen only AFTER selected query
        compatible = [(p, table) for p, table in compatible if table[4*a+b] == target]
        used.append(ep)
        trace.update(bytes((ep,a,b,target,len(compatible))))
    if not compatible:
        return "REJECT", None, used, trace.hexdigest(), 0
    tables = {table for _, table in compatible}
    if len(tables) == 1:
        return "ACCEPT", next(iter(tables)), used, trace.hexdigest(), len(compatible)
    return "ABSTAIN", None, used, trace.hexdigest(), len(compatible)


def challenge(seed, selected, pairs, max_extra):
    decision, frozen_table, used, first_trace, compatible = selected
    tested = []
    second_trace = hashlib.sha256()
    if decision == "ACCEPT":
        remaining = [ep for ep in range(16) if ep not in used]
        remaining.sort(key=lambda ep: (hashlib.sha256(bytes(pairs[ep])).digest(), ep))
        for ep in remaining[:max_extra]:
            a, b = pairs[ep]
            true_label = observed(seed, ep)[2]
            predicted = frozen_table[4*a+b]
            tested.append(ep)
            second_trace.update(bytes((ep, a, b, true_label, predicted)))
            if true_label != predicted:
                decision = "REJECT"
                break
    return {
        "decision":decision,
        "selection_labels":len(used),
        "witness_labels":len(tested),
        "total_labels":len(used)+len(tested),
        "selection_order":used,
        "witness_order":tested,
        "selection_trace_sha256":first_trace,
        "witness_trace_sha256":second_trace.hexdigest(),
        "compatible_hypotheses_at_selection":compatible,
        "frozen_selected_table_sha256":(
            hashlib.sha256(bytes(frozen_table)).hexdigest() if frozen_table else None),
        "table":frozen_table if decision == "ACCEPT" else None,
    }


def accepted_score(seed, table):
    matrix = [[0]*4 for _ in range(4)]
    digest = hashlib.sha256()
    for ep in range(16):
        observations = observed(seed, ep)
        for tick in range(3,32):
            actual = observations[tick]
            predicted = table[4*observations[tick-2] + observations[tick-1]]
            matrix[actual][predicted] += 1
            digest.update(bytes((ep,tick,actual,predicted)))
    assert sum(map(sum,matrix)) == 464
    return {
        "correct":sum(matrix[i][i] for i in range(4)), "total":464,
        "confusion":matrix, "trace_sha256":digest.hexdigest(),
    }


def run():
    fingerprints = {
        str(seed):hashlib.sha256(bytes(
            value for ep in range(16) for value in observed(seed,ep))).hexdigest()
        for seed in range(136,144)
    }
    if len(set(fingerprints.values())) != 8:
        return {"schema":"wingless.w-r29-frozen-witness.v1",
                "classification":"NONUNIQUE_FINGERPRINT",
                "fingerprints":fingerprints,"rows":[],"model_calls":0}
    records = []
    for seed in range(136,144):
        pairs = pool(seed)
        first = selected_model(seed, pairs)
        for budget in CHALLENGE_BUDGETS:
            entry = challenge(seed, first, pairs, budget)
            fitted = entry.pop("table")
            records.append({
                "seed":seed,
                "cohort":"development" if seed < 140 else "heldout",
                "family":"ADD4" if seed%2==0 else "XOR4",
                "extra_budget":budget,
                **entry,
                "test":accepted_score(seed,fitted) if fitted is not None else None,
            })
    held=[r for r in records if r["seed"]>=140]
    stats={}
    for budget in CHALLENGE_BUDGETS:
        test=[r for r in held if r["extra_budget"]==budget]
        adds=[r for r in test if r["family"]=="ADD4"]
        xors=[r for r in test if r["family"]=="XOR4"]
        stats[str(budget)]={
            "add_accepted":sum(r["decision"]=="ACCEPT" for r in adds),
            "add_perfect":sum(r["decision"]=="ACCEPT" and
                              r["test"]["correct"]==464 for r in adds),
            "xor_rejected":sum(r["decision"]=="REJECT" for r in xors),
            "xor_false_accepted":sum(r["decision"]=="ACCEPT" for r in xors),
            "xor_abstained":sum(r["decision"]=="ABSTAIN" for r in xors),
            "label_counts":{str(r["seed"]):r["total_labels"] for r in test},
        }
    full = stats["16"]
    full_valid=(full["add_perfect"]==2 and full["xor_rejected"]==2)
    bounded=any(stats[str(k)]["add_perfect"]==2 and
                stats[str(k)]["xor_rejected"]==2 for k in (0,2,4))
    if full["xor_false_accepted"]:
        classification="FALSE_ACCEPT"
    elif full_valid and bounded:
        classification="SAFE_BOUNDED_CHALLENGE"
    elif full_valid:
        classification="CHALLENGE_INSUFFICIENT"
    else:
        classification="MIXED"
    return {"schema":"wingless.w-r29-frozen-witness.v1",
            "parent":"141e25c2a19f3ccfca5e11ae066bd8015b168a84",
            "classification":classification,
            "candidate_permutations":len(MODELS),
            "fingerprints":fingerprints,
            "heldout":stats,"rows":records,"model_calls":0,
            "claim_limit":"Finite algebraic model family; explicit independent target-label queries and incomplete-context false acceptance."}


if __name__=="__main__":
    output=run()
    assert output==run()
    print(json.dumps(output,sort_keys=True,indent=2))
