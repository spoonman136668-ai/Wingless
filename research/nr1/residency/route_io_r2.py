"""R2 bounded measured CPU/disk synthetic expert residency pilot.

Measures real file reads, checks hashes and cold *application-level* cache, not
cold operating-system page cache or GPU / inference performance.
"""
import argparse
from collections import Counter
import hashlib
import json
from pathlib import Path
import statistics
import tempfile
import time

from route_cache_r1 import EVAL, traces

CORPUS = 128
BYTES_PER_EXPERT = 64 * 1024
CAPACITY = 8
POLICIES = ("no_cache", "lru", "training_hotset")
PATTERNS = ("uniform", "hotset", "phase_shift")
SEEDS = (11, 23, 37)
SAMPLES = 48


def payload(expert, size=BYTES_PER_EXPERT):
    if type(expert) is not int or not 0 <= expert < CORPUS:
        raise ValueError("expert identity invalid")
    if type(size) is not int or not 0 < size <= 1024 * 1024:
        raise ValueError("unbounded payload")
    digest = hashlib.sha256(f"r2-expert-{expert}".encode()).digest()
    return (digest * (size // len(digest) + 1))[:size]


def create_corpus(directory, size=BYTES_PER_EXPERT):
    if not directory.is_dir() or any(directory.iterdir()):
        raise ValueError("corpus directory must be empty")
    hashes = {}
    for expert in range(CORPUS):
        content = payload(expert, size)
        path = directory / f"expert-{expert:03d}.bin"
        with path.open("xb") as handle:
            handle.write(content)
        hashes[expert] = hashlib.sha256(content).hexdigest()
    return hashes


def measure(directory, hashes, training, requests, policy, capacity=CAPACITY):
    if policy not in POLICIES or type(capacity) is not int or not 1 <= capacity < CORPUS:
        raise ValueError("invalid residency policy/capacity")
    if not requests or any(type(e) is not int or not 0 <= e < CORPUS for e in requests):
        raise ValueError("invalid evaluation requests")
    if any(type(e) is not int or not 0 <= e < CORPUS for e in training):
        raise ValueError("invalid training requests")
    if not isinstance(hashes, dict) or len(hashes) != CORPUS:
        raise ValueError("incomplete source custody")
    hot_counts = Counter(training)
    hot = set(sorted(range(CORPUS), key=lambda i: (-hot_counts[i], i))[:capacity])
    cache = {}
    clock = {}
    hits = misses = evictions = bytes_read = 0
    latencies_ns = []
    for t, expert in enumerate(requests):
        start = time.perf_counter_ns()
        if expert in cache:
            data = cache[expert]
            hits += 1
        else:
            misses += 1
            path = directory / f"expert-{expert:03d}.bin"
            with path.open("rb") as stream:
                data = stream.read()
            bytes_read += len(data)
            if hashlib.sha256(data).hexdigest() != hashes[expert]:
                raise ValueError("tampered synthetic expert bytes")
            admit = policy == "lru" or (policy == "training_hotset" and expert in hot)
            if admit:
                if len(cache) >= capacity:
                    victim = min(cache, key=lambda k: (clock[k], k))
                    del cache[victim]
                    del clock[victim]
                    evictions += 1
                cache[expert] = data
        if len(data) != BYTES_PER_EXPERT:
            raise ValueError("unexpected expert size")
        clock[expert] = t if expert in cache else -1
        latencies_ns.append(time.perf_counter_ns() - start)
    if bytes_read != misses * BYTES_PER_EXPERT or hits + misses != len(requests):
        raise RuntimeError("resource accounting inconsistent")
    if len(cache) > capacity:
        raise RuntimeError("capacity exceeded")
    return {
        "policy": policy, "requests": len(requests),
        "hits": hits, "misses": misses, "evictions": evictions,
        "application_cached_experts_at_end": len(cache),
        "application_peak_capacity_bytes": capacity * BYTES_PER_EXPERT,
        "read_syscall_payload_bytes": bytes_read,
        "observed_total_request_wall_ns": sum(latencies_ns),
        "observed_p95_request_wall_ns": sorted(latencies_ns)[(95 * len(latencies_ns) + 99) // 100 - 1],
        "os_page_cache_cold": "UNVERIFIED",
        "observed_gpu_transfer": False,
        "original_model_executed": False,
    }


def frozen_requests(evaluation):
    """Keep both halves of a routing phase shift; no future-driven policy tuning."""
    if len(evaluation) != EVAL or SAMPLES % 2 or SAMPLES > len(evaluation):
        raise ValueError("invalid frozen evaluation split")
    half = len(evaluation) // 2
    a = [expert for row in evaluation[:half] for expert in row][:SAMPLES // 2]
    b = [expert for row in evaluation[half:] for expert in row][:SAMPLES // 2]
    if len(a) != SAMPLES // 2 or len(b) != SAMPLES // 2:
        raise ValueError("incomplete independent phase windows")
    return a + b


def experiment():
    observations = []
    with tempfile.TemporaryDirectory(prefix="wingless-residency-r2-") as root:
        path = Path(root)
        hashes = create_corpus(path)
        for pattern in PATTERNS:
            for seed in SEEDS:
                train, evaluation = traces(pattern, seed)
                # Frozen 32-expert subset; retain original chronological order.
                train = [i for row in train for i in row]
                test = frozen_requests(evaluation)
                if pattern == "phase_shift":
                    first, second = test[:SAMPLES // 2], test[SAMPLES // 2:]
                    if not (all(0 <= i < 16 for i in first)
                            and all(64 <= i < 80 for i in second)):
                        raise RuntimeError("frozen phase-shift transition missing")
                for policy in POLICIES:
                    result = measure(path, hashes, train, test, policy)
                    result.update(pattern=pattern, seed=seed)
                    observations.append(result)
    return {
        "schema": "wingless.nr1.residency.r2.measured-disk-synthetic.v1",
        "status": "R2_SYNTHETIC_CPU_DISK_EXPERT_READS_COMPLETE",
        "source": "DETERMINISTIC_SYNTHETIC_FILE_CORPUS_ONLY",
        "expert_count": CORPUS, "expert_bytes": BYTES_PER_EXPERT,
        "sampling": "FIRST_24_REQUESTS_PER_EVALUATION_HALF",
        "phase_shift_later_experts_included": True,
        "frozen_patterns": list(PATTERNS), "frozen_seeds": list(SEEDS),
        "policy_set": list(POLICIES), "observations": observations,
        "source_files_cleaned": True,
        "actual_file_reads_measured": True,
        "os_page_cache_control": "NOT_COLD_OR_CONTROLLED",
        "gpu_transfers_measured": False, "original_weights_used": False,
        "model_inference_executed": False, "scientific_acceptance": "NOT_GRANTED",
    }


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    if args.output.exists():
        raise FileExistsError("refuse to overwrite evidence")
    data = experiment()
    args.output.parent.mkdir(parents=True, exist_ok=True)
    with args.output.open("x", encoding="utf-8") as stream:
        json.dump(data, stream, indent=2, sort_keys=True)
        stream.write("\n")
    print("RESIDENCY_R2_REAL_FILE_READS_PASS", len(data["observations"]))


if __name__ == "__main__":
    main()
