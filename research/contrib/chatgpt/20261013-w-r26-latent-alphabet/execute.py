"""W-R26: supervised latent alphabet conjugacy search, with abstention."""
import hashlib
import itertools
import json

PERMUTATIONS = tuple(itertools.permutations(range(4)))
assert len(PERMUTATIONS) == 24
ALL_CONTEXTS = tuple((a, b) for a in range(4) for b in range(4))
BUDGETS = (4, 8, 16)


def true_permutation(seed):
    assert 112 <= seed <= 119
    digest = hashlib.sha256(("W26/" + str(seed)).encode("utf-8")).digest()
    return PERMUTATIONS[digest[0] % 24]


def latent_sequence(seed, episode):
    values = [(episode + seed) % 4, ((episode // 4) + 3 * seed) % 4]
    while len(values) < 32:
        values.append((values[-2] + values[-1]) % 4)
    assert len(values) == 32
    return tuple(values)


def observed_sequence(seed, episode):
    perm = true_permutation(seed)
    return tuple(perm[value] for value in latent_sequence(seed, episode))


def table_for(perm):
    inverse = [0] * 4
    for latent, observed in enumerate(perm):
        inverse[observed] = latent
    return tuple(perm[(inverse[before] + inverse[previous]) % 4]
                 for before, previous in ALL_CONTEXTS)


def fits(perm, observed_triples):
    table = table_for(perm)
    return all(table[4 * before + previous] == target
               for before, previous, target in observed_triples)


def fixture():
    perm = (2, 0, 3, 1)
    table = table_for(perm)
    cases = [(a, b, table[4*a+b]) for a, b in ALL_CONTEXTS]
    matches = [p for p in PERMUTATIONS if fits(p, cases)]
    assert matches
    assert len({table_for(p) for p in matches}) == 1
    return "PASS"


def calibrate(seed, count):
    evidence = []
    for episode in range(count):
        rows = observed_sequence(seed, episode)
        evidence.append((rows[0], rows[1], rows[2]))
    assert len(evidence) == count
    assert len({(a,b) for a,b,_ in evidence}) == count
    return evidence


def evaluate(seed, table, kind):
    matrix = [[0] * 4 for _ in range(4)]
    digest = hashlib.sha256()
    for episode in range(16):
        rows = observed_sequence(seed, episode)
        for tick in range(3,32):
            previous = rows[tick-1]
            before = rows[tick-2]
            prediction = ((before+previous)%4 if kind=="baseline"
                          else table[4*before+previous])
            actual = rows[tick]
            matrix[actual][prediction] += 1
            digest.update(bytes((episode,tick,actual,prediction)))
    assert sum(map(sum,matrix)) == 16*29
    correct = sum(matrix[i][i] for i in range(4))
    return {
        "correct": correct,"total":16*29,
        "accuracy":correct/(16*29),
        "confusion":matrix,"trace_sha256":digest.hexdigest(),
    }


def run():
    assert fixture()=="PASS"
    rows=[]
    for seed in range(112,120):
        baseline=evaluate(seed,None,"baseline")
        rows.append({"seed":seed,"arm":"no_adaptation",
                     "calibration_labels":0,"accepted":True,
                     "matched_permutations":0,"distinct_predictive_tables":1,
                     "test":baseline})
        for budget in BUDGETS:
            cases=calibrate(seed,budget)
            consistent=[perm for perm in PERMUTATIONS if fits(perm,cases)]
            assert consistent, "true permutation must explain its own calibration"
            induced={table_for(perm) for perm in consistent}
            accepted=len(induced)==1
            model=next(iter(induced)) if accepted else None
            rows.append({
                "seed":seed,"arm":"enumeration_k"+str(budget),
                "calibration_labels":budget,
                "accepted":accepted,
                "matched_permutations":len(consistent),
                "distinct_predictive_tables":len(induced),
                "table_sha256":hashlib.sha256(bytes(model)).hexdigest()
                               if accepted else None,
                "test":evaluate(seed,model,"inferred") if accepted else None,
            })
    coverage={}
    accuracy={}
    for arm in ("no_adaptation", "enumeration_k4", "enumeration_k8","enumeration_k16"):
        selected=[r for r in rows if r["arm"]==arm]
        coverage[arm]=sum(bool(r["accepted"]) for r in selected)
        n=sum(r["test"]["total"] for r in selected if r["accepted"])
        accuracy[arm]=(sum(r["test"]["correct"] for r in selected if r["accepted"])/n
                       if n else None)
    if (coverage["enumeration_k16"]==8
        and accuracy["enumeration_k16"]==1
        and coverage["enumeration_k16"]>=coverage["enumeration_k4"]
        and accuracy["no_adaptation"]<.9):
        result="REPRESENTATION_RECOVERY"
    elif any(coverage["enumeration_k"+str(b)] for b in BUDGETS):
        result="PARTIAL"
    else:
        result="NEGATIVE"
    return {
        "schema":"wingless.w-r26-alphabet-permutation.v1",
        "predecessor_head":"ebead7f1a2d194553d28db3e352f6e9bb039cedb",
        "classification":result,
        "fixture":"PASS","candidate_permutations":24,
        "coverage_seeds_by_arm":coverage,
        "accepted_accuracy_by_arm":accuracy,
        "test_predictions_per_seed":464,
        "rows":rows,"model_calls":0,
        "claim_limit":"Fixed 24-symbol-permutation search and known modular recurrence, not unrestricted learned semantics."
    }


if __name__=="__main__":
    outcome=run()
    assert outcome==run()
    print(json.dumps(outcome,sort_keys=True,indent=2))
