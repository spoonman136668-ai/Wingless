"""R1: frozen synthetic MoE expert-cache screen. No weights, inference or router traces."""
import argparse
from collections import Counter, defaultdict, deque
import json
import math
from pathlib import Path
import random

N = 128
FANOUT = 8
TRAIN = 32
EVAL = 96
SEEDS = (11, 23, 37)
PATTERNS = ("uniform", "hotset", "phase_shift", "cyclic_adversarial")
CAPACITIES = (8, 16, 32)
POLICIES = ("no_cache", "lru", "lfu", "training_hotset", "hindsight_belady")
EXPERT_BYTES = 8 * 1024 * 1024
BANDWIDTH_BPS = 8 * 1024**3
MISS_SETUP_MS = 0.15
COMPUTE_MS = 1.0


def traces(pattern, seed):
    if pattern not in PATTERNS or type(seed) is not int:
        raise ValueError("unknown frozen trace")
    rng = random.Random(seed)
    rows = []
    for t in range(TRAIN + EVAL):
        if pattern == "uniform":
            row = rng.sample(range(N), FANOUT)
        elif pattern == "hotset":
            row = rng.sample(range(16) if rng.random() < .9 else range(16, N), FANOUT)
        elif pattern == "phase_shift":
            row = rng.sample(range(16) if t < TRAIN + EVAL//2 else range(64, 80), FANOUT)
        else:
            row = [(t * FANOUT + j) % 40 for j in range(FANOUT)]
        rows.append(tuple(row))
    return tuple(rows[:TRAIN]), tuple(rows[TRAIN:])


def simulate(training, evaluation, policy, capacity, bandwidth=BANDWIDTH_BPS):
    if policy not in POLICIES:
        raise ValueError("unsupported policy")
    if type(capacity) is not int or not 1 <= capacity < N:
        raise ValueError("cache cannot hold a complete expert or exceeds corpus")
    if not isinstance(bandwidth, (int, float)) or not math.isfinite(bandwidth) or bandwidth <= 0:
        raise ValueError("invalid assumed bandwidth")
    if not isinstance(training, (list, tuple)) or not isinstance(evaluation, (list, tuple)) or not evaluation:
        raise ValueError("missing frozen traces")
    for token in tuple(training) + tuple(evaluation):
        if not isinstance(token, (list, tuple)) or len(token) != FANOUT or len(set(token)) != FANOUT:
            raise ValueError("wrong routing fanout or duplicate expert")
        if any(type(e) is not int or not 0 <= e < N for e in token):
            raise ValueError("invalid expert id")
    counts = Counter(e for row in training for e in row)
    hot = set(sorted(range(N), key=lambda e: (-counts[e], e))[:capacity])
    future = defaultdict(deque)
    requests = [e for row in evaluation for e in row]
    if policy == "hindsight_belady":
        for i, e in enumerate(requests):
            future[e].append(i)
    resident, last, freq = set(), {}, Counter()
    hits = misses = evictions = max_resident = 0
    token_ms = []
    serial = 0
    for token in evaluation:
        faults = 0
        for e in token:
            if policy == "hindsight_belady":
                if not future[e] or future[e][0] != serial:
                    raise RuntimeError("invalid hindsight chronology")
                future[e].popleft()
            serial += 1
            if e in resident:
                hits += 1
            else:
                misses += 1
                faults += 1
                if policy != "no_cache" and (policy != "training_hotset" or e in hot):
                    if len(resident) == capacity:
                        if policy == "lru":
                            victim = min(resident, key=lambda v: (last.get(v, -1), v))
                        elif policy == "lfu":
                            victim = min(resident, key=lambda v: (freq[v], last.get(v, -1), v))
                        elif policy == "hindsight_belady":
                            victim = max(resident, key=lambda v: (future[v][0] if future[v] else math.inf, v))
                        else:
                            raise RuntimeError("training hotset exceeded capacity")
                        resident.remove(victim)
                        evictions += 1
                    resident.add(e)
            last[e] = serial
            freq[e] += 1
            max_resident = max(max_resident, len(resident))
        token_ms.append(COMPUTE_MS + faults * (1000 * EXPERT_BYTES / bandwidth + MISS_SETUP_MS))
    ordered = sorted(token_ms)
    return {
        "policy": policy,
        "future_access_used": policy == "hindsight_belady",
        "capacity_experts": capacity,
        "hit_count": hits, "miss_count": misses, "evictions": evictions,
        "blocking_page_faults": misses,
        "hit_fraction": hits / len(requests),
        "transferred_bytes": misses * EXPERT_BYTES,
        "peak_gpu_resident_bytes_assumed": max_resident * EXPERT_BYTES,
        "entire_disk_corpus_bytes_assumed": N * EXPERT_BYTES,
        "cpu_staging_bytes_per_miss_assumed": EXPERT_BYTES,
        "bandwidth_bytes_per_second_assumed": bandwidth,
        "predicted_total_ms": sum(token_ms),
        "predicted_p95_token_ms": ordered[math.ceil(.95 * len(ordered)) - 1],
        "runtime_measured": False,
    }


def evaluate():
    scenarios = []
    for pattern in PATTERNS:
        for seed in SEEDS:
            training, evaluation = traces(pattern, seed)
            for cap in CAPACITIES:
                results = {p: simulate(training, evaluation, p, cap) for p in POLICIES}
                oracle = results["hindsight_belady"]
                if oracle["miss_count"] > results["lru"]["miss_count"]:
                    raise RuntimeError("hindsight lower bound violated")
                if results["no_cache"]["hit_count"] != 0:
                    raise RuntimeError("uncached control unexpectedly hit")
                scenarios.append({"pattern": pattern, "seed": seed, "capacity": cap,
                                  "training_tokens": len(training), "eval_tokens": len(evaluation),
                                  "results": results})
    return {
        "schema": "wingless.nr1.residency.r1.synthetic.v1",
        "result": "FROZEN_ROUTER_CACHE_SYNTHETIC_SCREEN_PASS",
        "hypothesis": "intact_expert_routing_aware_paging",
        "frozen_patterns": list(PATTERNS), "frozen_seeds": list(SEEDS),
        "capacities_experts": list(CAPACITIES), "scenarios": scenarios,
        "oracle_role": "NON_DEPLOYABLE_HINDSIGHT_BOUND",
        "static_training": "DISJOINT_PRE_EVALUATION_PREFIX_ONLY",
        "original_weights_loaded": False, "authentic_router_traces": False,
        "gpu_measurement": "UNVERIFIED", "full_model_inference": False,
        "protected_heldout_used": False, "scientific_acceptance": "NOT_GRANTED",
    }


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    if args.output.exists():
        raise FileExistsError("refuse overwrite")
    result = evaluate()
    args.output.parent.mkdir(parents=True, exist_ok=True)
    with args.output.open("x", encoding="utf-8") as fp:
        json.dump(result, fp, indent=2, sort_keys=True)
        fp.write("\n")
    print("RESIDENCY_R1_SYNTHETIC_PASS", len(result["scenarios"]))


if __name__ == "__main__":
    main()
