"""R5 causal adaptive admission for intact synthetic expert residency.

No GPU operations, original weights, original routing traces or model execution.
Every cache decision uses only the previously observed expert requests.
"""
import argparse
from collections import Counter, deque
import json
import math
from pathlib import Path

from route_cache_r1 import (
    BANDWIDTH_BPS, CAPACITIES, COMPUTE_MS, EVAL, EXPERT_BYTES, FANOUT,
    MISS_SETUP_MS, N, PATTERNS, SEEDS, POLICIES, simulate, traces,
)

RECENCY_WINDOW = 128
MIN_PRIOR_REQUESTS = 1
NEW_POLICY = "causal_repeat_admission"
BASELINES = ("no_cache", "lru", "lfu", "training_hotset")
BYTES_PER_REQUEST = EXPERT_BYTES


def adaptive_simulate(training, evaluation, capacity, *, window=RECENCY_WINDOW):
    """Never inspect future requests or alter the number of resident experts."""
    if type(capacity) is not int or capacity < 1 or capacity >= N:
        raise ValueError("invalid cache capacity")
    if type(window) is not int or window < FANOUT or window > 4096:
        raise ValueError("invalid bounded history window")
    if not isinstance(training, (list, tuple)) or not isinstance(evaluation, (list, tuple)) or not evaluation:
        raise ValueError("invalid frozen traces")
    for row in tuple(training) + tuple(evaluation):
        if (not isinstance(row, (tuple, list)) or len(row) != FANOUT
                or len(set(row)) != FANOUT or any(type(e) is not int or not 0 <= e < N for e in row)):
            raise ValueError("invalid routed expert tuple")

    # Training trace deliberately ignored: candidate adapts to actual prior
    # evaluation requests only and does not know the later phase transition.
    resident = set()
    last_access = {}
    recent = deque()
    frequencies = Counter()
    misses = hits = evictions = admission_bypasses = peak = 0
    serial = 0
    token_latencies = []
    token_misses = []

    for row in evaluation:
        faults = 0
        for expert in row:
            if expert in resident:
                hits += 1
            else:
                misses += 1
                faults += 1
                if frequencies[expert] >= MIN_PRIOR_REQUESTS:
                    if len(resident) >= capacity:
                        victim = min(resident, key=lambda e: (last_access.get(e, -1), e))
                        resident.remove(victim)
                        evictions += 1
                    resident.add(expert)
                else:
                    admission_bypasses += 1
            # Update after the decision, preserving one-way causality.
            serial += 1
            last_access[expert] = serial
            recent.append(expert)
            frequencies[expert] += 1
            if len(recent) > window:
                expired = recent.popleft()
                frequencies[expired] -= 1
                if frequencies[expired] == 0:
                    del frequencies[expired]
            peak = max(peak, len(resident))
        token_misses.append(faults)
        token_latencies.append(COMPUTE_MS + faults * (1000 * EXPERT_BYTES / BANDWIDTH_BPS + MISS_SETUP_MS))

    total_requests = len(evaluation) * FANOUT
    result = {
        "policy": NEW_POLICY,
        "future_access_used": False,
        "training_or_eval_answer_labels_used": False,
        "history_window_requests": window,
        "admission_requires_prior_requests": MIN_PRIOR_REQUESTS,
        "capacity_experts": capacity,
        "requests": total_requests,
        "hits": hits, "misses": misses, "evictions": evictions,
        "admission_bypasses": admission_bypasses,
        "peak_resident_experts": peak,
        "peak_gpu_resident_bytes_assumed": peak * EXPERT_BYTES,
        "transferred_bytes": misses * EXPERT_BYTES,
        "hit_fraction": hits / total_requests,
        "token_misses": token_misses,
        "predicted_total_ms": sum(token_latencies),
        "predicted_p95_token_ms": sorted(token_latencies)[math.ceil(.95 * len(token_latencies)) - 1],
        "measured_hardware_latency": False,
    }
    if peak > capacity or hits + misses != total_requests or admission_bypasses > misses:
        raise RuntimeError("invalid cache accounting")
    return result


def evaluate():
    cases = []
    for pattern in PATTERNS:
        for seed in SEEDS:
            train, test = traces(pattern, seed)
            if len(test) != EVAL:
                raise RuntimeError("frozen evaluation length drift")
            for cap in CAPACITIES:
                candidate = adaptive_simulate(train, test, cap)
                controls = {policy: simulate(train, test, policy, cap) for policy in BASELINES}
                if candidate["misses"] > controls["no_cache"]["miss_count"]:
                    raise RuntimeError("candidate worse than uncached accounting")
                if candidate["peak_gpu_resident_bytes_assumed"] > cap * EXPERT_BYTES:
                    raise RuntimeError("synthetic memory capacity exceeded")
                cases.append({
                    "pattern": pattern, "seed": seed, "capacity_experts": cap,
                    "candidate": candidate, "controls": controls,
                    "candidate_hit_delta_vs_lru": candidate["hits"] - controls["lru"]["hit_count"],
                    "candidate_hit_delta_vs_training_hotset": candidate["hits"] - controls["training_hotset"]["hit_count"],
                })
    return {
        "schema": "wingless.nr1.residency.r5.causal-admission-development.v1",
        "status": "R5_CAUSAL_ADAPTIVE_RESIDENCY_SCREEN_COMPLETE",
        "mechanism": "BOUNDED_HISTORY_REPEATED_EXPERT_ADMISSION",
        "patterns": list(PATTERNS), "seeds": list(SEEDS), "capacities_experts": list(CAPACITIES),
        "baseline_policies": list(BASELINES),
        "candidate_policy": NEW_POLICY,
        "frozen_scenario_count": len(cases),
        "cases": cases,
        "original_weights_used": False, "actual_router_traces_used": False,
        "gpu_transfer_measured": False, "inference_executed": False,
        "protected_heldout_used": False, "scientific_acceptance": "NOT_GRANTED",
    }


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--output", required=True, type=Path)
    args = parser.parse_args()
    if args.output.exists():
        raise FileExistsError("refusing to overwrite frozen R5 receipt")
    receipt = evaluate()
    args.output.parent.mkdir(parents=True, exist_ok=True)
    with args.output.open("x", encoding="utf-8") as stream:
        json.dump(receipt, stream, indent=2, sort_keys=True)
        stream.write("\n")
    print("R5_CAUSAL_ADMISSION_SYNTHETIC_PASS", receipt["frozen_scenario_count"])


if __name__ == "__main__":
    main()
