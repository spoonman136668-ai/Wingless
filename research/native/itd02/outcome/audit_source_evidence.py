"""Read-only source-run evidence auditor. Never grants scientific acceptance."""
from __future__ import annotations

import argparse
import hashlib
import json
import math
from pathlib import Path


class EvidenceError(ValueError):
    pass


def digest(value: object) -> str:
    encoded = json.dumps(value, sort_keys=True, separators=(",", ":"), allow_nan=False).encode()
    return hashlib.sha256(encoded).hexdigest()


def assess(record: dict) -> dict:
    if record.get("schema") != "wingless.itd02.source-training.v1" or record.get("science") != "NOT_ACCEPTED":
        raise EvidenceError("unrecognized or falsely accepted source record")
    if record.get("status") != "ENGINEERING_SOURCE_TRAINING_ONLY":
        raise EvidenceError("record exceeds source engineering authority")
    events, episodes = record.get("events"), record.get("episodes")
    if not isinstance(events, list) or not events or not isinstance(episodes, list) or not episodes:
        raise EvidenceError("missing source events or episode summaries")
    if record.get("event_digest_sha256") != digest(events):
        raise EvidenceError("source event digest mismatch")
    checkpoint = record.get("checkpoint_base64")
    if not isinstance(checkpoint, str) or not checkpoint or hashlib.sha256(checkpoint.encode("ascii")).hexdigest() != record.get("checkpoint_sha256"):
        raise EvidenceError("checkpoint digest mismatch")
    seeds = record.get("seeds")
    limit = record.get("max_steps")
    if not isinstance(seeds, list) or not all(type(s) is int for s in seeds) or not isinstance(limit, int) or not 1 <= limit <= 1000:
        raise EvidenceError("invalid declared source budget")
    if len(events) > len(seeds) * limit or len(episodes) > len(seeds):
        raise EvidenceError("source budget exceeded")
    counts = {}
    totals = {}
    done = {}
    for i, event in enumerate(events, 1):
        if type(event.get("event_id")) is not int or event["event_id"] != i:
            raise EvidenceError("non-sequential selected feedback")
        seed, action = event.get("seed"), event.get("selected_action")
        if type(seed) is not int or seed not in seeds or type(action) is not int or action not in (0, 1, 2):
            raise EvidenceError("invalid source seed or action")
        reward = event.get("reward")
        if type(reward) not in (int, float) or not math.isfinite(reward):
            raise EvidenceError("nonfinite reward")
        if type(event.get("terminal")) is not bool or type(event.get("truncated")) is not bool or event["terminal"] and event["truncated"]:
            raise EvidenceError("invalid source termination")
        for field in ("before", "after"):
            state = event.get(field)
            if not isinstance(state, dict) or set(state) != {"x", "y", "goal_x", "goal_y", "direction", "front_blocked"}:
                raise EvidenceError("unrecognized public state")
        if done.get(seed):
            raise EvidenceError("event after source termination")
        counts[seed] = counts.get(seed, 0) + 1
        totals[seed] = totals.get(seed, 0.0) + reward
        done[seed] = event["terminal"] or event["truncated"]
    for episode in episodes:
        seed = episode.get("seed")
        if seed not in counts or counts[seed] != episode.get("steps") or counts[seed] > limit:
            raise EvidenceError("source step accounting mismatch")
        if abs(totals[seed] - episode.get("return", float("inf"))) > 1e-9:
            raise EvidenceError("source reward accounting mismatch")
        if type(episode.get("source_done")) is not bool or done[seed] != episode["source_done"]:
            raise EvidenceError("source completion accounting mismatch")
    completed = sum(1 for e in episodes if e["source_done"])
    nonzero = sum(1 for e in episodes if e["return"] != 0)
    return {
        "evidence": "SOURCE_TRANSCRIPTS_INTEGRITY_PASS",
        "episodes": len(episodes),
        "steps": len(events),
        "completed": completed,
        "nonzero_reward_episodes": nonzero,
        "source_outcome": (
            "BOUNDED_INCOMPLETE_ZERO_REWARD" if completed == 0 and nonzero == 0
            else "SOURCE_OUTCOME_OBSERVED_NOT_TRANSFER_PROOF"
        ),
        "scientific_acceptance": "NOT_ACCEPTED",
    }


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("source_evidence", type=Path)
    args = parser.parse_args()
    result = assess(json.loads(args.source_evidence.read_text(encoding="utf-8")))
    print(json.dumps(result, sort_keys=True))


if __name__ == "__main__":
    main()
