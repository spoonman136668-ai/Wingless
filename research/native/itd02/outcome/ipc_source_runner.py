"""Opt-in ITD-02 source-world / native IPC training harness.

Engineering plumbing only; no claimed scientific transfer. Runs real MiniGrid
only when explicitly invoked. Never inspects source environment internals.
"""
from __future__ import annotations

import argparse
import hashlib
import json
from pathlib import Path
import subprocess
import sys
from typing import Any

from minigrid_source import SourceEpisode, SourceObservationError, open_source_environment


class BridgeError(RuntimeError):
    pass


def canonical(value: Any) -> bytes:
    return json.dumps(value, sort_keys=True, separators=(",", ":"), allow_nan=False).encode("utf-8")


class NativeBridge:
    def __init__(self, binary: str):
        self.process = subprocess.Popen(
            [binary], stdin=subprocess.PIPE, stdout=subprocess.PIPE,
            stderr=subprocess.DEVNULL, text=True, encoding="utf-8", bufsize=1,
        )
        if self.process.stdin is None or self.process.stdout is None:
            raise BridgeError("missing native pipes")

    def request(self, op: str, **payload: Any) -> dict[str, Any]:
        request = {"op": op, **payload}
        try:
            self.process.stdin.write(canonical(request).decode("utf-8") + "\n")
            self.process.stdin.flush()
            line = self.process.stdout.readline()
            if not line or len(line) > 8 << 20:
                raise BridgeError("missing or oversized native reply")
            reply = json.loads(line)
        except (OSError, UnicodeError, ValueError) as exc:
            raise BridgeError("native protocol failed") from exc
        if not isinstance(reply, dict) or reply.get("ok") is not True:
            raise BridgeError("native rejected " + op + ": " + str(reply.get("error") if isinstance(reply, dict) else reply))
        return reply

    def close(self) -> None:
        if self.process.poll() is None:
            self.process.stdin.close()
            try:
                self.process.wait(timeout=5)
            except subprocess.TimeoutExpired:
                self.process.kill()
                self.process.wait(timeout=5)


def run_training_episode(
    source: SourceEpisode, bridge: Any, seed: int, max_steps: int, start_id: int,
    min_count: int = 2,
) -> tuple[list[dict[str, Any]], int, bool]:
    """One episode; only authentic completed transitions are accepted.

    A budget-expired episode is explicitly incomplete, not a fabricated
    source truncation. No held-out evaluation happens in this function.
    """
    if not 1 <= max_steps <= 1000 or not 2 <= min_count <= 65535 or start_id < 1:
        raise ValueError("invalid predeclared budget or event id")
    before = source.reset(seed)
    events: list[dict[str, Any]] = []
    completed = False
    for _ in range(max_steps):
        reply = bridge.request("choose", before=before.wire(), min_count=min_count)
        action = reply.get("action")
        if type(action) is not int or action not in (0, 1, 2):
            raise BridgeError("native returned illegal selected action")
        after, reward, terminal, truncated = source.step(action)
        event_id = start_id + len(events)
        feedback = bridge.request(
            "feedback", before=before.wire(), action=action, after=after.wire(),
            reward=reward, terminal=terminal, truncated=truncated, event_id=event_id,
        )
        if feedback.get("updates") != event_id:
            raise BridgeError("native selected feedback count mismatch")
        events.append({
            "seed": seed, "event_id": event_id, "before": before.wire(),
            "selected_action": action, "after": after.wire(), "reward": reward,
            "terminal": terminal, "truncated": truncated,
        })
        before = after
        if terminal or truncated:
            completed = True
            bridge.request("episode_start")
            break
    return events, start_id + len(events), completed


def run_source_training(binary: str, environment_id: str, seeds: list[int],
                        max_steps: int, output: Path) -> dict[str, Any]:
    if not seeds or len(seeds) > 16 or len(set(seeds)) != len(seeds):
        raise ValueError("invalid or duplicated predeclared training seeds")
    if output.exists():
        raise FileExistsError("refusing to overwrite existing evidence")
    source = open_source_environment(environment_id)
    bridge = NativeBridge(binary)
    events: list[dict[str, Any]] = []
    records: list[dict[str, Any]] = []
    next_id = 1
    try:
        for seed in seeds:
            episode, next_id, completed = run_training_episode(
                source, bridge, seed, max_steps, next_id)
            events.extend(episode)
            records.append({"seed": seed, "steps": len(episode), "source_done": completed,
                            "return": sum(row["reward"] for row in episode)})
            # The learner remains valid, but an incomplete world episode cannot
            # be treated as an authentic source terminal. Stop without claiming
            # the next source reset constitutes native episode completion.
            if not completed:
                break
        snapshot = bridge.request("snapshot")
        state = snapshot.get("state")
        if not isinstance(state, str) or not state:
            raise BridgeError("missing final native checkpoint")
    finally:
        bridge.close()
        close = getattr(source.environment, "close", None)
        if callable(close):
            close()
    result = {
        "schema": "wingless.itd02.source-training.v1",
        "status": "ENGINEERING_SOURCE_TRAINING_ONLY",
        "environment_id": environment_id, "seeds": seeds, "max_steps": max_steps,
        "events": events, "episodes": records,
        "checkpoint_base64": state,
        "checkpoint_sha256": hashlib.sha256(state.encode("ascii")).hexdigest(),
        "event_digest_sha256": hashlib.sha256(canonical(events)).hexdigest(),
        "science": "NOT_ACCEPTED",
    }
    # Atomic, exclusive creation; a write failure must not overwrite evidence.
    with output.open("x", encoding="utf-8") as stream:
        stream.write(json.dumps(result, sort_keys=True, separators=(",", ":"), allow_nan=False) + "\n")
    return result


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--native-binary", required=True)
    parser.add_argument("--environment", choices=("MiniGrid-Empty-5x5-v0", "MiniGrid-Empty-Random-5x5-v0", "MiniGrid-Empty-8x8-v0"), required=True)
    parser.add_argument("--seeds", type=int, nargs="+", required=True)
    parser.add_argument("--max-steps", type=int, default=100)
    parser.add_argument("--output", type=Path, required=True)
    parser.add_argument("--run-source", action="store_true", help="explicit source execution opt-in")
    args = parser.parse_args()
    if not args.run_source:
        parser.error("SOURCE_NOT_RUN: --run-source required")
    result = run_source_training(args.native_binary, args.environment, args.seeds, args.max_steps, args.output)
    print(json.dumps({"status": result["status"], "episodes": len(result["episodes"]),
                      "events": len(result["events"]), "science": result["science"]}))


if __name__ == "__main__":
    main()
