"""Bounded cross-family frozen-preference comparison: Empty -> FourRooms.

Evidence is exploratory: fixed baseline is not a trained conventional learner,
and clone-per-action frozen selection is not the full native online actor.
"""
from __future__ import annotations
import argparse
import base64
import hashlib
import json
from pathlib import Path
import resource
import time
from frozen_source_probe import FrozenPolicy, arm
from ipc_source_runner import BridgeError
from source_controls import FixedPolicy, episode
from reward_ablation_restart import MaskedOutcome
from minigrid_source import open_source_environment

TARGET="MiniGrid-FourRooms-v0"
SEEDS=(61,67,71,73)
LIMIT=128

def run(binary: str, destination: Path) -> dict:
    if destination.exists():
        raise FileExistsError("refusing to overwrite source evidence")
    start=time.perf_counter()
    binary_sha256=hashlib.sha256(Path(binary).read_bytes()).hexdigest()
    outcomes={}
    trained={}
    for name,blind in (("reward_aware",False),("reward_blind",True)):
        training,snapshot,checkpoint_sha=arm(binary,blind)
        trained[name]={"training":training,"checkpoint_sha256":checkpoint_sha}
        learner=MaskedOutcome(binary,False)
        try:
            restored=learner.bridge.request("load",state=snapshot)
            if restored.get("updates")!=training["steps"]:
                raise BridgeError("cross-family checkpoint update mismatch")
            episodes=[]
            for seed in SEEDS:
                original=learner.bridge.request("snapshot")["state"]
                if original!=snapshot:
                    raise BridgeError("native checkpoint mutated before cross-family step")
                world=open_source_environment(TARGET)
                try:
                    result=episode(world,FrozenPolicy(learner.bridge),seed,LIMIT)
                finally:
                    world.environment.close()
                if learner.bridge.request("snapshot")["state"]!=snapshot:
                    raise BridgeError("cross-family frozen checkpoint mutated")
                episodes.append(result)
            outcomes[name]=episodes
        finally:
            learner.close()
    fixed=[]
    for seed in SEEDS:
        world=open_source_environment(TARGET)
        try:
            fixed.append(episode(world,FixedPolicy(),seed,LIMIT))
        finally:
            world.environment.close()
    for i in range(len(SEEDS)):
        public=[outcomes[k][i]["events"][0]["before"] for k in ("reward_aware","reward_blind")]
        public.append(fixed[i]["events"][0]["before"])
        if not (public[0]==public[1]==public[2]):
            raise BridgeError("cross-family public initial state mismatch")
    if len({json.dumps(e["events"][0]["before"],sort_keys=True) for e in fixed})<2:
        raise BridgeError("cross-family source public diversity insufficient")
    if len(trained["reward_aware"]["training"]["events"])!=len(trained["reward_blind"]["training"]["events"]) or any(
        a["action"]!=b["action"] or a["before"]!=b["before"]
        for a,b in zip(trained["reward_aware"]["training"]["events"],trained["reward_blind"]["training"]["events"])
    ):
        raise BridgeError("unmatched training public trajectories")
    usage=resource.getrusage(resource.RUSAGE_SELF)
    data={
      "schema":"wingless.itd02.unfamiliar-family-source.v1",
      "scientific_acceptance":"NOT_ACCEPTED",
      "scope":"EXPLORATORY_CROSS_FAMILY_FROZEN_PREFERENCE",
      "training_family":"MiniGrid-Empty-Random-5x5-v0",
      "target_family":TARGET,
      "training_seed":19,"target_seeds":list(SEEDS),"step_budget":LIMIT,
      "minigrid_version":"2.3.1","gymnasium_version":"0.29.1",
      "native_binary_sha256":binary_sha256,
      "trained":trained,"results":{**outcomes,"fixed":fixed},
      "frozen_checkpoint_verified_each_episode":True,
      "resource_accounting":{
        "process_wall_seconds":time.perf_counter()-start,
        "python_parent_maxrss_kib_linux":usage.ru_maxrss,
        "scope":"Python parent process only; native subprocess RSS not included",
        "hardware":"GitHub hosted ubuntu-24.04, not consumer PC"
      },
      "limitations":"Only fixed control, not strong conventional learning baseline; four seeds, previously examined training seed, clone-per-action selection",
    }
    destination.write_text(json.dumps(data,sort_keys=True,separators=(",",":"),allow_nan=False)+"\n",encoding="utf-8")
    print(json.dumps({k:[{"seed":e["seed"],"steps":e["steps"],"return":e["return"],"completed":e["completed"]} for e in v] for k,v in data["results"].items()},sort_keys=True))
    return data

if __name__=="__main__":
    p=argparse.ArgumentParser()
    p.add_argument("--native-binary",required=True)
    p.add_argument("--output",required=True,type=Path)
    p.add_argument("--run-source",action="store_true")
    a=p.parse_args()
    if not a.run_source:
        p.error("SOURCE_NOT_RUN")
    run(a.native_binary,a.output)
