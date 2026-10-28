"""PG8 prospective Empty-8x8 geometry-shift development probe.

Pre-registered before this code in research/itd02-pg8-prospective-geometry-prereg-20261028.md.
Reuses frozen native V4 learner and equal-exposure source protocol without modifying them.
Not a test of causal entity binding, 24h retention, zero-shot transfer, or science acceptance.
"""
from __future__ import annotations

import argparse
import hashlib
import json
import re
from pathlib import Path

import exact_exposure_gate as frozen
import progress_stateful_eval_v4 as v4
from fourrooms_stateful_diagnostic import PublicReactive
from ipc_source_runner import BridgeError

SOURCE_HEAD = "4f547b79708c7f3249a6cf2afb96d0f6f5b39919"
V4_NATIVE_ANCESTOR = "21928da6e939b84cf71f70476215bfa348aecb85"
TARGET_ENV = "MiniGrid-Empty-8x8-v0"
TRAIN_ENV = "MiniGrid-Empty-Random-5x5-v0"
TARGET_SEEDS = (18041, 18047, 18061, 18077, 18089, 18101, 18119, 18133)
TEST_STEPS = 100
ARM_NAMES = ("progress_native", "progress_v3", "motion_native", "original_native", "linear_q")

def locked_contract() -> None:
    assert v4.TRAIN_SEEDS == tuple(7101 + 2*i for i in range(256))
    assert v4.TRAIN_STEPS == 256 and v4.TEST_STEPS == TEST_STEPS
    assert frozen.TRAIN_ENV == TRAIN_ENV and frozen.TEST_ENV == "MiniGrid-FourRooms-v0"
    assert v4.TEST_SEEDS == (8101, 8111, 8123, 8137, 8147, 8161, 8171, 8191)
    assert len(set(TARGET_SEEDS)) == len(TARGET_SEEDS) == 8
    assert not (set(TARGET_SEEDS) & (set(v4.TRAIN_SEEDS) | set(v4.TEST_SEEDS) | v4.PREVIOUS_HELDOUT))

def classify(goals: dict[str, int]) -> str:
    if set(goals) != set(ARM_NAMES) | {"public_reactive"}:
        raise BridgeError("PG8_ARM_SET_DRIFT")
    if any(type(n) is not int or not 0 <= n <= len(TARGET_SEEDS) for n in goals.values()):
        raise BridgeError("PG8_GOAL_COUNT_INVALID")
    best_control = max(n for k, n in goals.items() if k != "progress_native")
    if goals["progress_native"] >= 1 and goals["progress_native"] >= best_control:
        return "PG8_DEVELOPMENT_POSITIVE_NOT_SCIENCE"
    return "PG8_NATIVE_TRANSFER_NO_GO"

def execute(original_binary: str, motion_binary: str, progress_binary: str,
            predecessor_binary: str, raw_output: Path, report_output: Path,
            source_sha: str) -> dict:
    locked_contract()
    if not re.fullmatch(r"[0-9a-f]{40}", source_sha):
        raise BridgeError("PG8_SOURCE_SHA_INVALID")
    if raw_output == report_output or raw_output.exists() or report_output.exists():
        raise FileExistsError("PG8_REFUSE_EVIDENCE_OVERWRITE")
    if not raw_output.parent.is_dir() or not report_output.parent.is_dir():
        raise BridgeError("PG8_OUTPUT_DIRECTORY_MISSING")

    # Local module constants only. No native source, training protocol, or
    # original V4 evidence is mutated or post-result tuned.
    v4.TEST_SEEDS = TARGET_SEEDS
    frozen.TEST_ENV = TARGET_ENV
    data = v4.execute(original_binary, motion_binary, progress_binary,
                      predecessor_binary, raw_output)
    if data.get("science") != "NOT_ACCEPTED" or data.get("test_time_adaptation") is not True:
        raise BridgeError("PG8_SOURCE_SCIENCE_CONTRACT_DRIFT")
    if data.get("train_family") != TRAIN_ENV or data.get("eval_family") != TARGET_ENV:
        raise BridgeError("PG8_SOURCE_FAMILY_DRIFT")
    if data.get("train_seeds") != list(v4.TRAIN_SEEDS) or data.get("heldout_eval_seeds") != list(TARGET_SEEDS):
        raise BridgeError("PG8_SOURCE_SEED_DRIFT")
    if data.get("train_source_interactions_per_arm") != 256 or set(data["arms"]) != set(ARM_NAMES):
        raise BridgeError("PG8_EXPOSURE_OR_ARM_DRIFT")

    raw_goals = {}
    raw_per_seed = {}
    for name in ARM_NAMES:
        arm = data["arms"][name]
        if sum(x["steps"] for x in arm["training"]) != 256:
            raise BridgeError("PG8_TRAINING_ACCOUNTING_DRIFT " + name)
        if len(arm["evaluation"]) != len(TARGET_SEEDS):
            raise BridgeError("PG8_EVALUATION_COUNT_DRIFT " + name)
        rows = []
        for seed, evaluation in zip(TARGET_SEEDS, arm["evaluation"]):
            ep = evaluation["episode"]
            if evaluation["seed"] != seed or ep["seed"] != seed or not ep["events"]:
                raise BridgeError("PG8_EVALUATION_SEED_OR_EVENT_DRIFT " + name)
            if evaluation["online_updates"] != ep["steps"] or ep["steps"] < 1 or ep["steps"] > TEST_STEPS:
                raise BridgeError("PG8_ONLINE_UPDATE_MISMATCH " + name)
            if name != "linear_q" and (evaluation["updates_before"] != 256 or
                                       evaluation["updates_after"] != 256 + ep["steps"] or
                                       not evaluation["checkpoint_changed_during_online_eval"]):
                raise BridgeError("PG8_NATIVE_RESTORE_OR_FEEDBACK_DRIFT " + name)
            rows.append({"seed": seed, "goal": ep["goal_reached"], "return": ep["return"],
                         "steps": ep["steps"], "events_sha256": ep["events_sha256"]})
        raw_per_seed[name] = rows
        raw_goals[name] = sum(row["goal"] is True for row in rows)

    # Strong goal-directed but nonlearning reference uses the same public six
    # fields and starts. Its training budget is explicitly zero, not matched.
    reactive_rows = []
    for i, seed in enumerate(TARGET_SEEDS):
        result = frozen.run_episode(TARGET_ENV, PublicReactive(), seed, TEST_STEPS)
        start = result["events"][0]["before"]
        for name in ARM_NAMES:
            if data["arms"][name]["evaluation"][i]["episode"]["events"][0]["before"] != start:
                raise BridgeError("PG8_UNEQUAL_PUBLIC_START " + name)
        reactive_rows.append({"seed": seed, "goal": result["goal_reached"],
                              "return": result["return"], "steps": result["steps"],
                              "events_sha256": result["events_sha256"]})
    raw_per_seed["public_reactive"] = reactive_rows
    raw_goals["public_reactive"] = sum(row["goal"] is True for row in reactive_rows)
    classification = classify(raw_goals)

    report = {
        "schema": "wingless.itd02.pg8.prospective-geometry-source.v1",
        "scope": "FRESH_PREDECLARED_EMPTY8X8_SOURCE_PUBLIC_GEOMETRY_ONLY",
        "classification": classification,
        "science": "NOT_ACCEPTED",
        "scientific_first_use": False,
        "entity_role_binding_tested": False,
        "cold_24h_retention_tested": False,
        "retail_hardware_qualified": False,
        "source_sha": source_sha,
        "frozen_proposal_parent": SOURCE_HEAD,
        "native_source_ancestor": V4_NATIVE_ANCESTOR,
        "training_environment": TRAIN_ENV,
        "target_environment": TARGET_ENV,
        "training_source_steps_per_learning_arm": 256,
        "target_max_steps_per_seed": TEST_STEPS,
        "target_seeds": list(TARGET_SEEDS),
        "test_time_updates": True,
        "reactive_training_steps": 0,
        "goals_by_arm": raw_goals,
        "per_seed": raw_per_seed,
        "underlying_v4_evidence_sha256": hashlib.sha256(raw_output.read_bytes()).hexdigest(),
        "resource": data["resource"],
        "limitations": "Previously exposed training source; target seeds fresh only relative to listed V4 cohorts; no globally sealed unfamiliar source, no causal entity roles, no 24/48-hour delay; public goal-directed control is nonlearning and has zero training; conventional linear-Q is unoptimized; online feedback during evaluation; GitHub-hosted Linux is not an out-of-box retail PC."
    }
    with report_output.open("x", encoding="utf-8") as out:
        json.dump(report, out, sort_keys=True, separators=(",", ":"), allow_nan=False)
        out.write("\n")
    print(json.dumps({"classification": classification, "science": "NOT_ACCEPTED",
                      "goals_by_arm": raw_goals, "source_sha": source_sha}, sort_keys=True))
    return report

if __name__ == "__main__":
    p = argparse.ArgumentParser()
    p.add_argument("--original-binary", required=True)
    p.add_argument("--motion-binary", required=True)
    p.add_argument("--progress-binary", required=True)
    p.add_argument("--predecessor-binary", required=True)
    p.add_argument("--raw-output", type=Path, required=True)
    p.add_argument("--report-output", type=Path, required=True)
    p.add_argument("--source-sha", required=True)
    p.add_argument("--run-source", action="store_true")
    args = p.parse_args()
    if not args.run_source:
        p.error("PG8_SOURCE_NOT_RUN")
    execute(args.original_binary, args.motion_binary, args.progress_binary,
            args.predecessor_binary, args.raw_output, args.report_output, args.source_sha)
