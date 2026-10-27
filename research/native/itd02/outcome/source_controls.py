"""ITD-02 bounded authentic source controls: fixed and frozen native R8.

Controls use the same public MiniGrid observation decoder and legal three-action
interface as the selected-outcome learner. Results are engineering source
observations only: no held-out transfer or science acceptance is inferred.
"""
from __future__ import annotations

import argparse
import hashlib
import json
from pathlib import Path
from typing import Any

from ipc_source_runner import BridgeError, NativeBridge, canonical
from minigrid_source import SourceEpisode, open_source_environment


class FixedPolicy:
    """Deterministic reward-blind left/forward baseline."""
    def __init__(self):
        self.armed = False
        self.steps = 0

    def choose(self, before: dict[str, Any]) -> int:
        if self.armed:
            raise BridgeError("fixed policy has unresolved action")
        self.armed = True
        # Public observation is supplied, but fixed policy never reads reward.
        return 2 if self.steps % 4 else 0

    def feedback(self, before: dict[str, Any], action: int,
                 after: dict[str, Any], reward: float,
                 terminal: bool, truncated: bool) -> None:
        if not self.armed:
            raise BridgeError("unselected fixed feedback")
        self.armed = False
        self.steps += 1

    def reset(self) -> None:
        if self.armed:
            raise BridgeError("fixed reset during pending feedback")

    def close(self) -> None:
        if self.armed:
            raise BridgeError("fixed close during pending feedback")


class R8Policy:
    """Native original R8 bridge: no outcome-reward selector."""
    def __init__(self, binary: str):
        self.bridge = NativeBridge(binary)
        self.next_id = 1

    def choose(self, before: dict[str, Any]) -> int:
        result = self.bridge.request("choose", before=before)
        action = result.get("action")
        if type(action) is not int or action not in (0, 1, 2):
            raise BridgeError("R8 selected illegal action")
        return action

    def feedback(self, before: dict[str, Any], action: int,
                 after: dict[str, Any], reward: float,
                 terminal: bool, truncated: bool) -> None:
        # R8 bridge's selected-feedback protocol does not accept reward. This
        # baseline therefore cannot obtain extra reward information.
        self.bridge.request("feedback", before=before, action=action, after=after)
        self.next_id += 1

    def reset(self) -> None:
        self.bridge.request("episode_start")

    def close(self) -> None:
        self.bridge.close()


class OutcomePolicy:
    """Native selected contextual reward learner."""
    def __init__(self, binary: str):
        self.bridge = NativeBridge(binary)
        self.next_id = 1

    def choose(self, before: dict[str, Any]) -> int:
        result = self.bridge.request("choose", before=before, min_count=2)
        action = result.get("action")
        if type(action) is not int or action not in (0, 1, 2):
            raise BridgeError("outcome selected illegal action")
        return action

    def feedback(self, before: dict[str, Any], action: int,
                 after: dict[str, Any], reward: float,
                 terminal: bool, truncated: bool) -> None:
        reply = self.bridge.request(
            "feedback", before=before, action=action, after=after,
            reward=reward, terminal=terminal, truncated=truncated,
            event_id=self.next_id,
        )
        if reply.get("updates") != self.next_id:
            raise BridgeError("outcome feedback counter mismatch")
        self.next_id += 1

    def reset(self) -> None:
        self.bridge.request("episode_start")

    def close(self) -> None:
        self.bridge.close()


def episode(source: SourceEpisode, policy: Any, seed: int, max_steps: int) -> dict[str, Any]:
    if type(seed) is not int or not 0 <= seed <= 2**32 - 1:
        raise ValueError("invalid seed")
    if type(max_steps) is not int or not 1 <= max_steps <= 1000:
        raise ValueError("invalid step bound")
    before = source.reset(seed)
    events = []
    for _ in range(max_steps):
        action = policy.choose(before.wire())
        if type(action) is not int or action not in (0, 1, 2):
            raise BridgeError("policy selected illegal action")
        after, reward, terminated, truncated = source.step(action)
        policy.feedback(before.wire(), action, after.wire(), reward, terminated, truncated)
        events.append({
            "before": before.wire(), "action": action, "after": after.wire(),
            "reward": reward, "terminated": terminated, "truncated": truncated,
        })
        before = after
        if terminated or truncated:
            policy.reset()
            break
    return {
        "seed": seed, "steps": len(events),
        "return": sum(e["reward"] for e in events),
        "completed": bool(events[-1]["terminated"] or events[-1]["truncated"]),
        "events_sha256": hashlib.sha256(canonical(events)).hexdigest(),
        "events": events,
    }


def compare(environment: str, seed: int, max_steps: int,
            outcome_binary: str, r8_binary: str, path: Path) -> dict[str, Any]:
    if path.exists():
        raise FileExistsError("will not overwrite control evidence")
    results = {}
    for name, policy in (
        ("fixed", FixedPolicy()),
        ("native_r8", R8Policy(r8_binary)),
        ("contextual_outcome", OutcomePolicy(outcome_binary)),
    ):
        source = open_source_environment(environment)
        try:
            results[name] = episode(source, policy, seed, max_steps)
        finally:
            policy.close()
            source.environment.close()
    record = {
        "schema": "wingless.itd02.source-controls-smoke.v1",
        "authority": "engineering-only-no-scientific-acceptance",
        "environment_id": environment,
        "seed": seed,
        "max_steps": max_steps,
        "controls": results,
        "transfer": "NOT_TESTED",
        "scientific_acceptance": "NOT_ACCEPTED",
    }
    with path.open("x", encoding="utf-8") as f:
        f.write(json.dumps(record, sort_keys=True, separators=(",", ":"), allow_nan=False) + "\n")
    return record


def main() -> None:
    p = argparse.ArgumentParser()
    p.add_argument("--environment", choices=("MiniGrid-Empty-5x5-v0", "MiniGrid-Empty-Random-5x5-v0", "MiniGrid-Empty-8x8-v0"), required=True)
    p.add_argument("--seed", type=int, required=True)
    p.add_argument("--max-steps", type=int, default=32)
    p.add_argument("--outcome-binary", required=True)
    p.add_argument("--r8-binary", required=True)
    p.add_argument("--output", type=Path, required=True)
    p.add_argument("--run-source", action="store_true")
    a = p.parse_args()
    if not a.run_source:
        p.error("SOURCE_NOT_RUN: explicit --run-source required")
    result = compare(a.environment, a.seed, a.max_steps, a.outcome_binary, a.r8_binary, a.output)
    print(json.dumps({
        "science": result["scientific_acceptance"],
        "transfer": result["transfer"],
        "results": {k: {"steps": v["steps"], "return": v["return"], "completed": v["completed"]}
                    for k, v in result["controls"].items()},
    }, sort_keys=True))


if __name__ == "__main__":
    main()
