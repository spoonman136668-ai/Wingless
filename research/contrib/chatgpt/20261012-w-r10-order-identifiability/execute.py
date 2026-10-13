"""W-R10: frozen temporal-order identifiability falsification.

This is synthetic mechanism evaluation, not a model-training claim.
"""
from collections import defaultdict
import json


def ordered_rows(seed, entity):
    q = (5 * entity + seed) % 9
    return tuple(
        (
            (192 + seed + 7 * q + 3 * tick) % 256,
            (208 + 3 * seed + 5 * q + 5 * tick) % 256,
            (224 + 5 * seed + 11 * q + 7 * tick) % 256,
            (240 + 7 * seed + 13 * q + 11 * tick) % 256,
        )
        for tick in range(16)
    )


def family_sequences(seed, scramble):
    items = []
    for entity in range(9):
        original = ordered_rows(seed, entity)
        tail = original[1:]
        for family in range(3):
            if scramble and family == 1:
                sequence = original[:1] + tail[::-1]
            elif scramble and family == 2:
                sequence = original[:1] + tail[5:] + tail[:5]
            else:
                sequence = original
            items.append((entity, family, sequence))
    return items


def group_cardinality(items, temporal_order_sensitive):
    groups = defaultdict(list)
    for entity, family, sequence in items:
        # The W-R9 scorer is order-insensitive except for the first row.
        key = tuple(sequence) if temporal_order_sensitive else (
            tuple(sorted(sequence)), sequence[0]
        )
        groups[key].append((entity, family))
    return sum(
        len(members) == 3
        and {f for _, f in members} == {0, 1, 2}
        and len({e for e, _ in members}) == 1
        for members in groups.values()
    )


def score_seed(seed):
    control = family_sequences(seed, False)
    scrambled = family_sequences(seed, True)
    for entity in range(9):
        original = ordered_rows(seed, entity)
        triples = [scrambled[entity * 3 + family][2] for family in range(3)]
        assert all(seq[0] == original[0] for seq in triples)
        assert all(tuple(sorted(seq)) == tuple(sorted(original)) for seq in triples)
        assert triples[0] == original
        assert triples[1] != original and triples[2] != original
        assert len(set(original)) == 16
    return {
        "seed": seed,
        "cohort": "development" if seed < 4 else "heldout",
        "control_sorted_anchor": group_cardinality(control, False),
        "scrambled_sorted_anchor": group_cardinality(scrambled, False),
        "scrambled_exact_order": group_cardinality(scrambled, True),
        "anchor_and_multiset_preserved": True,
        "temporal_order_changed": True,
    }


def run():
    rows = [score_seed(seed) for seed in range(8)]
    heldout = [row for row in rows if row["cohort"] == "heldout"]
    blind = all(
        row["control_sorted_anchor"] == 9
        and row["scrambled_sorted_anchor"] == 9
        and row["scrambled_exact_order"] == 0
        for row in heldout
    )
    return {
        "schema": "wingless.w-r10-temporal-order-identifiability.v1",
        "predecessor_head": "964149b9ac729ae567ce2ff7603d6ca634c8d6ca",
        "rows": rows,
        "classification": "TEMPORAL_ORDER_BLIND" if blind else "FALSIFICATION_NOT_CONFIRMED",
        "model_calls": 0,
        "claim_limit": "Falsification of one synthetic scorer; no general-language transfer or accepted-ref promotion",
    }


if __name__ == "__main__":
    outcome = run()
    assert outcome == run()
    print(json.dumps(outcome, sort_keys=True, indent=2))
