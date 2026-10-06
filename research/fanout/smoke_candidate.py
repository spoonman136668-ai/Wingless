#!/usr/bin/env python3
import argparse
import hashlib
import json
from pathlib import Path

STATE_DIMENSION = 8
READOUT_CAPACITY = 1
INTERFACE = "state_former(raw_input, history) -> fixed_size_state"

def state_former(raw_input, history, candidate_id, config):
    variant = int(config.get("variant", 0))
    salt = int(hashlib.sha256(candidate_id.encode("utf-8")).hexdigest()[:8], 16)
    history_total = sum(sum(int(v) for v in row) for row in history)
    raw_total = sum((i + 1) * int(v) for i, v in enumerate(raw_input))
    base = (salt + history_total + raw_total + 31 * variant) % 65521
    return [int((base + 17 * i + 3 * history_total + variant * i) % 257) for i in range(STATE_DIMENSION)]

def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--candidate", required=True)
    ap.add_argument("--candidate-config", required=True)
    ap.add_argument("--output", required=True)
    a = ap.parse_args()
    fixture = json.loads(Path(__file__).with_name("smoke_fixture.json").read_text(encoding="utf-8"))
    config = json.loads(Path(a.candidate_config).read_text(encoding="utf-8"))
    state = state_former(fixture["raw_input"], fixture["history"], a.candidate, config)
    result = {
        "candidate_id": a.candidate,
        "interface": INTERFACE,
        "fixed_state_dimension": STATE_DIMENSION,
        "fixed_readout_capacity": READOUT_CAPACITY,
        "fixture_class": "synthetic-replay",
        "metric": {"name": "synthetic_state_checksum", "value": round(sum(state) / (STATE_DIMENSION * 256.0), 12)},
        "resource_usage": {
            "parameters": 0,
            "context_bytes": len(json.dumps(fixture, sort_keys=True).encode("utf-8")),
            "model_calls": 0
        },
        "state_sha256": hashlib.sha256(json.dumps(state, separators=(",", ":")).encode("utf-8")).hexdigest()
    }
    Path(a.output).write_text(json.dumps(result, indent=2, sort_keys=True) + "\n", encoding="utf-8")

if __name__ == "__main__":
    main()
