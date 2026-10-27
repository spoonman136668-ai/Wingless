"""Exploratory authentic-source discrimination: rewards vs zeroed rewards.

NOT blinded held-out scientific evaluation: seed 19 is an already-observed
successful seed. This gate measures whether real reward information changes
a native policy after cold restart, without making a transfer claim.
"""
from __future__ import annotations
import argparse
import base64
import hashlib
import json
from pathlib import Path
from source_controls import OutcomePolicy, episode
from ipc_source_runner import canonical, BridgeError
from minigrid_source import open_source_environment

SOURCE = "MiniGrid-Empty-Random-5x5-v0"
TRAIN_SEED = 19
PROBE_SEEDS = (7, 11)
LIMIT = 64

class MaskedOutcome(OutcomePolicy):
    def __init__(self, binary: str, blind: bool):
        super().__init__(binary)
        self.blind = blind
    def feedback(self, before, action, after, reward, terminal, truncated):
        return super().feedback(before, action, after, 0.0 if self.blind else reward,
                                terminal, truncated)

def trial(binary: str, blind: bool) -> dict:
    train = MaskedOutcome(binary, blind)
    source = open_source_environment(SOURCE)
    try:
        training = episode(source, train, TRAIN_SEED, LIMIT)
        snap = train.bridge.request("snapshot")["state"]
        base64.b64decode(snap, validate=True)
        digest = hashlib.sha256(base64.b64decode(snap)).hexdigest()
    finally:
        train.close()
        source.environment.close()
    probes = []
    # Each probe starts from an independent, cold-restored training checkpoint.
    # Feedback during a probe necessarily updates its disposable copy only.
    for seed in PROBE_SEEDS:
        cold = MaskedOutcome(binary, blind)
        world = open_source_environment(SOURCE)
        try:
            reply = cold.bridge.request("load", state=snap)
            if reply.get("updates") != training["steps"]:
                raise BridgeError("cold checkpoint did not preserve update count")
            cold.next_id = training["steps"] + 1
            loaded = cold.bridge.request("snapshot")["state"]
            if loaded != snap:
                raise BridgeError("cold checkpoint bytes changed during restore")
            result = episode(world, cold, seed, LIMIT)
            probes.append(result)
        finally:
            cold.close()
            world.environment.close()
    return {"training": training, "train_checkpoint_sha256": digest,
            "probes": probes, "reward_blind": blind}

def compare(binary: str, output: Path) -> dict:
    if output.exists():
        raise FileExistsError("cannot overwrite experimental receipt")
    informed = trial(binary, False)
    blind = trial(binary, True)
    ta, tb = informed["training"], blind["training"]
    # Until selected real reward diverges, both policies must have had matched
    # training observations and selected actions; refuse a mismatched control.
    first_difference = next((i for i, (a,b) in enumerate(zip(ta["events"],tb["events"]))
       if a["action"] != b["action"] or a["before"] != b["before"]), None)
    observations_matched = (first_difference is None and ta["steps"] == tb["steps"])
    distinct_initial = all(
        informed["probes"][i]["events"][0]["before"] ==
        blind["probes"][i]["events"][0]["before"] for i in range(len(PROBE_SEEDS))
    )
    if not distinct_initial:
        raise BridgeError("held-out probe initial observations differ")
    record = {
        "schema":"wingless.itd02.reward-ablation-cold-restart.v1",
        "science":"NOT_ACCEPTED",
        "scope":"EXPLORATORY_KNOWN_SEED_NO_TRANSFER",
        "environment_id":SOURCE,
        "train_seed":TRAIN_SEED, "probe_seeds":list(PROBE_SEEDS),
        "max_steps_per_episode":LIMIT,
        "reward_aware":informed,
        "reward_blind":blind,
        "matched_training_trajectory":observations_matched,
        "first_training_trajectory_difference":first_difference,
        "distinct_trained_checkpoints":informed["train_checkpoint_sha256"]!=blind["train_checkpoint_sha256"],
        "cold_restart_retained":True,
    }
    output.write_text(json.dumps(record,sort_keys=True,separators=(",",":"),allow_nan=False)+"\n",encoding="utf-8")
    return record

def main():
    p=argparse.ArgumentParser()
    p.add_argument("--native-binary",required=True)
    p.add_argument("--output",required=True,type=Path)
    p.add_argument("--run-source",action="store_true")
    args=p.parse_args()
    if not args.run_source:
        p.error("SOURCE_NOT_RUN")
    r=compare(args.native_binary,args.output)
    print(json.dumps({
      "scope":r["scope"],"science":r["science"],
      "matched_training":r["matched_training_trajectory"],
      "distinct_checkpoints":r["distinct_trained_checkpoints"],
      "reward_aware_training_return":r["reward_aware"]["training"]["return"],
      "reward_blind_training_return":r["reward_blind"]["training"]["return"],
      "probe_returns_aware":[x["return"] for x in r["reward_aware"]["probes"]],
      "probe_returns_blind":[x["return"] for x in r["reward_blind"]["probes"]],
    },sort_keys=True))

if __name__=="__main__":
    main()
