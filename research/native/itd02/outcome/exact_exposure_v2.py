"""Prospective, source-free-reviewable successor to the frozen exact exposure gate.

The historical exact_exposure_gate.py, its failures and any output remain unchanged.
Only actual source transitions count. A reset after an incomplete episode is
forbidden unless the entire training budget has been exhausted.
"""
from __future__ import annotations

import argparse
from pathlib import Path
import exact_exposure_gate as frozen
from ipc_source_runner import BridgeError

TRAIN_STEPS = frozen.TRAIN_STEPS
TRAIN_SEEDS = frozen.TRAIN_SEEDS

def exact_train_v2(policy, run_episode=frozen.run_episode):
    remaining = TRAIN_STEPS
    records = []
    for seed in TRAIN_SEEDS:
        if remaining == 0:
            break
        # Unlike the frozen predecessor's 100-step intermediate cutoff,
        # a source world can continue up to the remaining authentic exposure.
        record = run_episode(frozen.TRAIN_ENV, policy, seed, remaining)
        steps = record.get("steps")
        if type(steps) is not int or not (1 <= steps <= remaining):
            raise BridgeError("invalid actual source transition count")
        records.append(record)
        remaining -= steps
        if not record.get("completed") and remaining:
            raise BridgeError("no source reset after incomplete episode")
    if remaining:
        raise BridgeError("predeclared seeds insufficient for authentic exposure")
    if sum(record["steps"] for record in records) != TRAIN_STEPS:
        raise BridgeError("exact exposure accounting mismatch")
    return records

def main():
    p = argparse.ArgumentParser()
    p.add_argument("--native-binary", required=True)
    p.add_argument("--output", required=True, type=Path)
    p.add_argument("--run-source", action="store_true")
    a = p.parse_args()
    if not a.run_source:
        p.error("SOURCE_NOT_RUN")
    # Explicitly opt into the successor: never rewrite the frozen module or
    # its historical results, controls, budget, seeds or frozen evaluation.
    frozen.exact_train = exact_train_v2
    frozen.execute(a.native_binary, a.output)

if __name__ == "__main__":
    main()
