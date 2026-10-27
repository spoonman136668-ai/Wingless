"""Read-only frozen-choice source test, exploratory training seed 19.

Frozen-choice semantics: for each new public observation, native choose runs
on an isolated clone of the trained checkpoint. The active native checkpoint
never gets armed or updated. This deliberately excludes learned route dynamics
within a probe and must not be mislabeled as ordinary online evaluation.
"""
from __future__ import annotations
import argparse
import base64
import hashlib
import json
from pathlib import Path
from ipc_source_runner import BridgeError
from source_controls import episode
from reward_ablation_restart import MaskedOutcome
from minigrid_source import open_source_environment

TRAIN_SEED=19
PROBE_SEEDS=(41,43,47,53)
LIMIT=64
SOURCE="MiniGrid-Empty-Random-5x5-v0"

class FrozenPolicy:
    def __init__(self, bridge):
        self.bridge=bridge
    def choose(self,before):
        response=self.bridge.request("peek_frozen",before=before,min_count=2)
        action=response.get("action")
        if type(action) is not int or action not in (0,1,2):
            raise BridgeError("illegal frozen source action")
        return action
    def feedback(self,before,action,after,reward,terminal,truncated):
        # This is a read-only preference snapshot test: source consequences
        # are recorded but never submitted as learning feedback.
        return None
    def reset(self):
        return None

def arm(binary,blind):
    learner=MaskedOutcome(binary,blind)
    world=open_source_environment(SOURCE)
    try:
        training=episode(world,learner,TRAIN_SEED,LIMIT)
        snap=learner.bridge.request("snapshot")["state"]
        digest=hashlib.sha256(base64.b64decode(snap,validate=True)).hexdigest()
        return training,snap,digest
    finally:
        learner.close()
        world.environment.close()

def evaluate(binary,snapshot,expected_steps):
    learner=MaskedOutcome(binary,False)
    source_results=[]
    try:
        loaded=learner.bridge.request("load",state=snapshot)
        if loaded.get("updates")!=expected_steps:
            raise BridgeError("frozen restore update mismatch")
        for seed in PROBE_SEEDS:
            before=learner.bridge.request("snapshot")["state"]
            if before!=snapshot:
                raise BridgeError("frozen snapshot differs before probe")
            source=open_source_environment(SOURCE)
            try:
                result=episode(source,FrozenPolicy(learner.bridge),seed,LIMIT)
            finally:
                source.environment.close()
            after=learner.bridge.request("snapshot")["state"]
            if after!=before or after!=snapshot:
                raise BridgeError("frozen native state mutated during evaluation")
            source_results.append(result)
        return source_results
    finally:
        learner.close()

def run(binary,output):
    if output.exists():
        raise FileExistsError("refuse to overwrite frozen probe evidence")
    outcome={}
    for blind in (False,True):
        training,snapshot,digest=arm(binary,blind)
        probes=evaluate(binary,snapshot,training["steps"])
        outcome["reward_blind" if blind else "reward_aware"]={
            "training":training, "checkpoint_sha256":digest,"probes":probes}
    a,b=outcome["reward_aware"],outcome["reward_blind"]
    if len(a["training"]["events"])!=len(b["training"]["events"]) or any(
        x["action"]!=y["action"] or x["before"]!=y["before"]
        for x,y in zip(a["training"]["events"],b["training"]["events"])
    ):
        raise BridgeError("unmatched source training before reward ablation")
    if a["checkpoint_sha256"]==b["checkpoint_sha256"]:
        raise BridgeError("source reward ablation failed to alter checkpoint")
    record={
        "schema":"wingless.itd02.frozen-preference-source.v1",
        "authority":"EXPLORATORY_NO_SCIENTIFIC_ACCEPTANCE",
        "source":SOURCE, "training_seed":TRAIN_SEED, "probe_seeds":list(PROBE_SEEDS),
        "max_steps":LIMIT, "choice_semantics":"FROZEN_CHECKPOINT_CLONE_PER_STEP",
        "reward_aware":a,"reward_blind":b,
        "scientific_acceptance":"NOT_ACCEPTED",
        "checkpoint_immutable_during_probes":True,
    }
    output.write_text(json.dumps(record,sort_keys=True,separators=(",",":"),allow_nan=False)+"\n",encoding="utf-8")
    print(json.dumps({k:[{"seed":v["seed"],"steps":v["steps"],"return":v["return"]} for v in record[k]["probes"]] for k in ("reward_aware","reward_blind")},sort_keys=True))

if __name__=="__main__":
    p=argparse.ArgumentParser()
    p.add_argument("--native-binary",required=True)
    p.add_argument("--output",type=Path,required=True)
    p.add_argument("--run-source",action="store_true")
    args=p.parse_args()
    if not args.run_source:
        p.error("SOURCE_NOT_RUN")
    run(args.native_binary,args.output)
