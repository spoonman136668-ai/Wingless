"""R3 matched, counterbalanced real synthetic CPU file-read experiments.

Controls policy order and app-cache reset; OS page-cache and scheduling remain
uncontrolled. A measured disk file read is not a measured physical disk transfer.
"""
import argparse
from collections import Counter
import hashlib
import itertools
import json
from pathlib import Path
import random
import statistics
import tempfile
import time

from route_cache_r1 import traces
from route_io_r2 import (BYTES_PER_EXPERT, CAPACITY, CORPUS, PATTERNS, POLICIES,
                         SEEDS, create_corpus, frozen_requests, measure)

ROUNDS = 6
PERMUTATIONS = tuple(itertools.permutations(POLICIES))
assert len(PERMUTATIONS) == 6


def balanced_order(seed, rounds=ROUNDS):
    """Each of the 3! orders is used once in every full block."""
    if type(seed) is not int or type(rounds) is not int or rounds != ROUNDS:
        raise ValueError("frozen rounds/seed mismatch")
    rng = random.Random(seed + 20261103)
    orders = list(PERMUTATIONS)
    rng.shuffle(orders)
    if sorted(orders) != sorted(PERMUTATIONS):
        raise RuntimeError("counterbalancing lost an order")
    return tuple(orders)


def checked_warmup(directory, hashes):
    """Read + verify once before timing; does not claim cold disk state."""
    for e in range(CORPUS):
        content = (directory / f"expert-{e:03d}.bin").read_bytes()
        if len(content) != BYTES_PER_EXPERT or hashlib.sha256(content).hexdigest() != hashes[e]:
            raise ValueError("warmup source-custody failure")


def matched_trial(directory, hashes, training, requests, order):
    if tuple(sorted(order)) != tuple(sorted(POLICIES)):
        raise ValueError("trial must run all and only matching control policies")
    results = {}
    for p in order:
        before = time.perf_counter_ns()
        v = measure(directory, hashes, training, requests, p, capacity=CAPACITY)
        elapsed = time.perf_counter_ns() - before
        if v["requests"] != len(requests):
            raise RuntimeError("request accounting drift")
        if elapsed < v["observed_total_request_wall_ns"]:
            raise RuntimeError("nested timing contradiction")
        v["observed_trial_wall_ns"] = elapsed
        v["application_cache_initial_state"] = "EMPTY"
        v["host_page_cache_state"] = "UNCONTROLLED_AFTER_VERIFIED_WARMUP"
        results[p] = v
    return results


def experiment():
    records = []
    with tempfile.TemporaryDirectory(prefix="wingless-residency-r3-") as td:
        root = Path(td)
        hashes = create_corpus(root)
        checked_warmup(root, hashes)
        for pattern in PATTERNS:
            for seed in SEEDS:
                train, ev = traces(pattern, seed)
                training = [e for row in train for e in row]
                test = frozen_requests(ev)
                if len(test) != 48:
                    raise RuntimeError("frozen sample changed")
                for round_no, order in enumerate(balanced_order(seed)):
                    controls = matched_trial(root, hashes, training, test, order)
                    if any(controls[p]["requests"] != 48 for p in POLICIES):
                        raise RuntimeError("incomplete paired trial")
                    if controls["no_cache"]["hits"] != 0:
                        raise RuntimeError("no-cache negative control had hits")
                    if any(v["application_peak_capacity_bytes"] != CAPACITY * BYTES_PER_EXPERT
                           for v in controls.values()):
                        raise RuntimeError("unequal cache budget")
                    records.append({"pattern": pattern, "seed": seed, "round": round_no,
                                    "policy_order": list(order), "controls": controls})
    if len(records) != len(PATTERNS) * len(SEEDS) * ROUNDS:
        raise RuntimeError("missing balanced trials")
    summaries = []
    for pattern in PATTERNS:
        matches = [r for r in records if r["pattern"] == pattern]
        for policy in POLICIES:
            values = [r["controls"][policy] for r in matches]
            summaries.append({
                "pattern": pattern, "policy": policy,
                "trial_count": len(values),
                "median_observed_trial_wall_ns": statistics.median(
                    v["observed_trial_wall_ns"] for v in values),
                "median_file_payload_bytes_read": statistics.median(
                    v["read_syscall_payload_bytes"] for v in values),
                "median_misses": statistics.median(v["misses"] for v in values),
                "median_hits": statistics.median(v["hits"] for v in values),
            })
    return {
        "schema": "wingless.nr1.residency.r3.counterbalanced-synthetic-file-reads.v1",
        "status": "R3_COUNTERBALANCED_SYNTHETIC_FILE_READS_COMPLETE",
        "base_source": "DETERMINISTIC_SYNTHETIC_EXPERT_CORPUS",
        "expert_count": CORPUS, "expert_bytes": BYTES_PER_EXPERT,
        "capacity_experts": CAPACITY,
        "patterns": list(PATTERNS), "seeds": list(SEEDS),
        "policies": list(POLICIES), "frozen_rounds": ROUNDS,
        "order_control": "ALL_SIX_PERMUTATIONS_ONCE_PER_SEED_AND_PATTERN",
        "application_cache_reset_each_policy": True,
        "filesystem_warmed_before_trials": True,
        "os_page_cache_control": "UNCONTROLLED",
        "cpu_disk_io_payload_observed": True,
        "physical_disk_bytes_measured": False,
        "gpu_transfer_measured": False,
        "original_weights_used": False, "authentic_router_traces_used": False,
        "inference_executed": False,
        "protected_heldout_used": False,
        "scientific_acceptance": "NOT_GRANTED",
        "records": records, "summaries": summaries,
    }


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--output", required=True, type=Path)
    args = ap.parse_args()
    if args.output.exists():
        raise FileExistsError("refuse to overwrite R3 receipt")
    receipt = experiment()
    args.output.parent.mkdir(parents=True, exist_ok=True)
    with args.output.open("x", encoding="utf-8") as f:
        json.dump(receipt, f, indent=2, sort_keys=True)
        f.write("\n")
    print("NR1_RESIDENCY_R3_COUNTERBALANCED_PASS", len(receipt["records"]))


if __name__ == "__main__":
    main()
