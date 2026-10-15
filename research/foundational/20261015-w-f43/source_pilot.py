"""W-F43 train-only deterministic posterior answer-information study.

Small exactly enumerated 6-value / 4-role source: three public triples per role
leave 3! compatible bijections per role (6**4 = 1,296 joint worlds).
No claim about exact inference in the unrestricted 8-value W-F42 substrate.
Private source maps stay in the bounded purchase port/evaluator.
"""
from __future__ import annotations
from collections import Counter, defaultdict
from itertools import permutations, product
from pathlib import Path
import json, math, random, time

SEEDS = range(331000, 331008)
FAMILIES = (390, 391, 392)
RESERVED_SEEDS = range(332000, 332016)
RESERVED_FAMILIES = (394, 395, 396, 397)
LENGTHS = (3, 5)
VALUES, ROLES, PUBLIC_PER_ROLE, PURCHASE_LIMIT = 6, 4, 3, 2
ARMS = ("exact_answer_entropy", "sampled_information_gain_equal_world_count",
        "sampled_decision_value_equal_world_count", "query_frontier_equal_world_count",
        "no_purchase", "private_truth_oracle_unfair")

def source(seed: int, family: int, length: int):
    if seed not in SEEDS or family not in FAMILIES or length not in LENGTHS:
        raise ValueError("source identity not admitted (including held-out)")
    rng = random.Random(seed * 12109 + family * 7919 + length * 601 + 43001)
    private = []
    public = []
    for role in range(ROLES):
        perm = list(range(VALUES))
        rng.shuffle(perm)
        private.append(tuple(perm))
        for x in rng.sample(range(VALUES), PUBLIC_PER_ROLE):
            public.append((role, x, perm[x]))
    rng.shuffle(public)
    start = rng.randrange(VALUES)
    if family == 390:
        chain = [(rng.randrange(ROLES), 1) for _ in range(length)]
    elif family == 391:
        cycle = [(2, -1), (1, 1), (2, 1), (0, -1), (3, 1)]
        chain = cycle[:length]
    else:
        cycle = [(0, 1), (1, -1), (0, -1), (2, 1), (3, -1)]
        chain = cycle[:length]
    return tuple(private), tuple(public), start, tuple(chain)

def follow(world, start, chain):
    x = start
    for role, direction in chain:
        m = world[role]
        x = m[x] if direction == 1 else m.index(x)
    return x

def worlds(public):
    constraints = [dict() for _ in range(ROLES)]
    for role, x, y in public:
        if x in constraints[role] and constraints[role][x] != y:
            raise ValueError("conflicting public mapping")
        constraints[role][x] = y
    per_role = []
    for known in constraints:
        if len(set(known.values())) != len(known):
            raise ValueError("non-bijective public input")
        unknown_x = [x for x in range(VALUES) if x not in known]
        unknown_y = [y for y in range(VALUES) if y not in known.values()]
        choices = []
        for fill in permutations(unknown_y):
            mapping = dict(known)
            mapping.update(zip(unknown_x, fill))
            choices.append(tuple(mapping[x] for x in range(VALUES)))
        per_role.append(choices)
    return list(product(*per_role))

def entropy(labels):
    n = len(labels)
    counts = Counter(labels)
    return -sum((c/n) * math.log2(c/n) for c in counts.values() if c)

def candidate_actions(public):
    seen = {(r, x) for r, x, _ in public}
    return [(r, x) for r in range(ROLES) for x in range(VALUES) if (r, x) not in seen]

def score_actions(bank, observed, start, chain):
    if not bank:
        raise ValueError("empty posterior")
    answers = [follow(w, start, chain) for w in bank]
    baseline_entropy = entropy(answers)
    base_correct = max(Counter(answers).values()) / len(answers)
    seen = {(r, x) for r, x, _ in observed}
    frontier = Counter()
    for w in bank:
        state = start
        for role, direction in chain:
            subject = state if direction == 1 else w[role].index(state)
            if (role, subject) not in seen:
                frontier[(role, subject)] += 1
            state = w[role][state] if direction == 1 else subject
    scored = {}
    for role, x in candidate_actions(observed):
        groups = defaultdict(list)
        for w, answer in zip(bank, answers):
            groups[w[role][x]].append(answer)
        post_entropy = sum(len(v)/len(bank)*entropy(v) for v in groups.values())
        post_accuracy = sum(max(Counter(v).values()) for v in groups.values()) / len(bank)
        scored[(role, x)] = (baseline_entropy - post_entropy,
                              post_accuracy - base_correct, frontier[(role, x)])
    return scored

def choose(method, observed, start, chain, seed):
    if method == "no_purchase":
        return None, 0
    bank = worlds(observed)
    if len(bank) > 1296:
        raise ValueError("certified source support bound exceeded")
    if method != "exact_answer_entropy":
        # Equal number of answer-world evaluations, fresh deterministic sampled bank.
        bank = random.Random(seed).choices(bank, k=len(bank))
    scores = score_actions(bank, observed, start, chain)
    if method in ("exact_answer_entropy", "sampled_information_gain_equal_world_count"):
        index = 0
    elif method == "sampled_decision_value_equal_world_count":
        index = 1
    elif method == "query_frontier_equal_world_count":
        index = 2
    else:
        raise ValueError("unsupported action method")
    action = max(scores, key=lambda a: (scores[a][index], -a[0], -a[1])) if scores else None
    return action, len(bank)

def predict(observed, start, chain):
    posterior = worlds(observed)
    counts = Counter(follow(w, start, chain) for w in posterior)
    return max(range(VALUES), key=lambda x: (counts[x], -x))

class PurchasePort:
    def __init__(self, private, public):
        self.__private = private
        self._known = {(r, x) for r, x, _ in public}
        self.receipts = []
    def buy(self, role, x):
        if len(self.receipts) >= PURCHASE_LIMIT:
            raise ValueError("purchase budget exhausted")
        if (role, x) in self._known:
            raise ValueError("public/duplicate mapping cannot be purchased")
        if not (0 <= role < ROLES and 0 <= x < VALUES):
            raise ValueError("invalid purchase")
        y = self.__private[role][x]
        self._known.add((role, x))
        self.receipts.append((role, x, y))
        return y

def run():
    assert set(SEEDS).isdisjoint(RESERVED_SEEDS)
    assert set(FAMILIES).isdisjoint(RESERVED_FAMILIES)
    t0 = time.perf_counter()
    rows = []
    for seed in SEEDS:
      for family in FAMILIES:
       for length in LENGTHS:
        private, public, start, chain = source(seed, family, length)
        truth = follow(private, start, chain)
        first_exact, initial_worlds = choose("exact_answer_entropy", public, start, chain, 0)
        second_exact, _ = choose("exact_answer_entropy", public, start, chain, 917)
        assert first_exact == second_exact
        answers = {"private_truth_oracle_unfair": truth}
        receipts = {}
        decision_world_counts = {}
        for arm in ARMS[:-1]:
            observed = list(public)
            port = PurchasePort(private, public)
            world_counts = []
            for step in range(0 if arm == "no_purchase" else PURCHASE_LIMIT):
                choice, worlds_evaluated = choose(arm, observed, start, chain,
                                                   seed * 37 + family * 7 + step * 10007)
                world_counts.append(worlds_evaluated)
                if choice is None:
                    break
                y = port.buy(*choice)
                observed.append((*choice, y))
            answers[arm] = predict(observed, start, chain)
            receipts[arm] = [list(p) for p in port.receipts]
            decision_world_counts[arm] = world_counts
        rows.append(dict(seed=seed, family=family, length=length, truth=truth,
                         answers=answers, receipts=receipts, first_exact=list(first_exact),
                         exact_rank_seed_invariant=first_exact == second_exact,
                         initial_worlds=initial_worlds,
                         world_evaluation_counts=decision_world_counts))
    scores = {arm: dict(correct=sum(r["answers"][arm] == r["truth"] for r in rows),
                        total=len(rows)) for arm in ARMS}
    for score in scores.values():
        score["accuracy"] = score["correct"] / score["total"]
    family_scores = {str(f): {arm: sum(r["answers"][arm] == r["truth"]
                          for r in rows if r["family"] == f) for arm in ARMS} for f in FAMILIES}
    fair = ARMS[1:-1]
    best = max(scores[a]["accuracy"] for a in fair)
    gain = scores["exact_answer_entropy"]["accuracy"] - best
    nondegrade = all(family_scores[str(f)]["exact_answer_entropy"] >=
                     max(family_scores[str(f)][a] for a in fair) for f in FAMILIES)
    qualified = gain >= .03 and nondegrade and all(r["exact_rank_seed_invariant"] for r in rows)
    return dict(schema="wingless.f43.certified-answer-information.trainonly.v1",
        stage="TRAIN_ONLY_EXACT_POSTERIOR_SOURCE_DIAGNOSTIC_NOT_ORIGINAL",
        seeds=[SEEDS.start, SEEDS.stop], families=list(FAMILIES),
        reserved_seeds=[RESERVED_SEEDS.start, RESERVED_SEEDS.stop],
        reserved_families=list(RESERVED_FAMILIES),
        questions=len(rows), exact_source_world_bound=1296,
        selector_world_evaluation_matching=True, posterior_exact_only_for_reduced_six_value_source=True,
        posterior_external_logical_payload_bytes=1296*ROLES*VALUES,
        wall_clock_seconds=time.perf_counter()-t0,
        scores=scores, by_family=family_scores, candidate_gain=gain,
        per_family_nondegrade=nondegrade, rows=rows,
        verdict="TRAIN_ONLY_EXACT_INFORMATION_HEADROOM" if qualified else "TRAIN_ONLY_EXACT_INFORMATION_NOT_QUALIFIED",
        private_truth_only_in_evaluator_and_paid_port=True,
        source_independent_F3_custody=False, primary_original_not_executed=True,
        native_576_byte_resident_qualified=False, accepted_ref_mutation=False,
        paid_api_calls=0)

if __name__ == "__main__":
    print(json.dumps(run(), sort_keys=True, indent=2))
