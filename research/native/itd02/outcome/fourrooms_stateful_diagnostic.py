"""FourRooms source-public navigation failure diagnosis (engineering only).

Compare cold-restored full *stateful* native actor against reward-blind full
actor, original R8, fixed baseline, and a goal-direction reactive heuristic.
No private map, privileged source state, or scientific transfer admission.
"""
from __future__ import annotations
import argparse
import base64
import hashlib
import json
from pathlib import Path
import resource
import time

from frozen_source_probe import arm
from reward_ablation_restart import MaskedOutcome
from source_controls import FixedPolicy, R8Policy, episode
from minigrid_source import open_source_environment
from ipc_source_runner import BridgeError

SOURCE = "MiniGrid-FourRooms-v0"
SEEDS = (79, 83, 89, 97)
BUDGET = 128

class PublicReactive:
    """Conventional nonlearning compass-and-obstacle baseline, not an oracle."""
    def __init__(self):
        self.armed = False
        self.turns = 0

    def choose(self, public):
        if self.armed:
            raise BridgeError("reactive selection pending")
        dx = public["goal_x"] - public["x"]
        dy = public["goal_y"] - public["y"]
        # Encoded direction 0 E, 1 S, 2 W, 3 N.
        preferred = (0 if dx > 0 else 2) if abs(dx) >= abs(dy) and dx != 0 else (1 if dy > 0 else 3)
        d = public["direction"]
        if d == preferred:
            action = 0 if public["front_blocked"] else 2
        elif (preferred - d) % 4 == 1:
            action = 1
        else:
            action = 0
        self.armed = True
        return action

    def feedback(self, before, action, after, reward, terminal, truncated):
        if not self.armed:
            raise BridgeError("reactive unselected feedback")
        self.armed = False

    def reset(self):
        if self.armed:
            raise BridgeError("reactive pending reset")

    def close(self):
        if self.armed:
            raise BridgeError("reactive pending close")


def diagnostic(result):
    events=result["events"]
    unique_positions={(e["before"]["x"],e["before"]["y"]) for e in events}
    counts={a:sum(e["action"]==a for e in events) for a in (0,1,2)}
    advances=sum((e["before"]["x"],e["before"]["y"])!=(e["after"]["x"],e["after"]["y"]) for e in events)
    blocked=sum(bool(e["before"]["front_blocked"]) for e in events)
    result["diagnostic"]={
        "unique_positions":len(unique_positions),
        "selected_action_counts":counts,
        "position_advances":advances,
        "blocked_front_observations":blocked,
        "source_goal_reached":result["goal_reached"],
        "source_truncated":result["source_truncated"],
    }


def test_arm(native_binary, r8_binary, name, seed, training_snapshot, training_steps):
    if name=="reward_aware" or name=="reward_blind":
        policy=MaskedOutcome(native_binary, name=="reward_blind")
        reply=policy.bridge.request("load",state=training_snapshot)
        if reply.get("updates")!=training_steps:
            raise BridgeError("trained native state not restored")
        policy.next_id=training_steps+1
        if policy.bridge.request("snapshot")["state"]!=training_snapshot:
            raise BridgeError("cold snapshot lost fidelity")
    elif name=="native_r8":
        policy=R8Policy(r8_binary)
    elif name=="fixed":
        policy=FixedPolicy()
    else:
        policy=PublicReactive()
    world=open_source_environment(SOURCE)
    try:
        outcome=episode(world,policy,seed,BUDGET)
        diagnostic(outcome)
        return outcome
    finally:
        policy.close()
        world.environment.close()


def run(native_binary, r8_binary, output):
    if output.exists():
        raise FileExistsError("refuse overwrite source evidence")
    start=time.perf_counter()
    ru_start=resource.getrusage(resource.RUSAGE_CHILDREN)
    arms={}
    for name,blind in (("reward_aware",False),("reward_blind",True)):
        training,snapshot,digest=arm(native_binary,blind)
        arms[name]={"training":training,"checkpoint_sha256":digest,"snapshot":snapshot}
    aware=arms["reward_aware"]["training"]["events"]
    blind=arms["reward_blind"]["training"]["events"]
    if len(aware)!=len(blind) or any(
        a["before"]!=b["before"] or a["action"]!=b["action"]
        for a,b in zip(aware,blind)
    ):
        raise BridgeError("training exposure mismatch")
    results={name:[] for name in ("reward_aware","reward_blind","native_r8","fixed","reactive")}
    for seed in SEEDS:
        for name in results:
            state=arms[name]["snapshot"] if name in arms else ""
            count=arms[name]["training"]["steps"] if name in arms else 0
            results[name].append(test_arm(native_binary,r8_binary,name,seed,state,count))
        starting=[results[name][-1]["events"][0]["before"] for name in results]
        if any(v!=starting[0] for v in starting[1:]):
            raise BridgeError("unmatched source-public seed start")
    if len({json.dumps(r["events"][0]["before"],sort_keys=True) for r in results["fixed"]})<2:
        raise BridgeError("public start diversity insufficient")
    ru_end=resource.getrusage(resource.RUSAGE_CHILDREN)
    record={
        "schema":"wingless.itd02.fourrooms-stateful-diagnostic.v1",
        "science":"NOT_ACCEPTED",
        "scope":"EXPLORATORY_SOURCE_DIAGNOSIS_ONLINE_FEEDBACK",
        "train_family":"MiniGrid-Empty-Random-5x5-v0",
        "target_family":SOURCE,"train_seed":19,"target_seeds":list(SEEDS),
        "max_requested_steps":BUDGET,
        "checkpoint_digests":{name:arms[name]["checkpoint_sha256"] for name in arms},
        "results":results,
        "resource":{
            "wall_seconds":time.perf_counter()-start,
            "child_cpu_user_seconds":ru_end.ru_utime-ru_start.ru_utime,
            "child_cpu_system_seconds":ru_end.ru_stime-ru_start.ru_stime,
            "native_binary_sha256":hashlib.sha256(Path(native_binary).read_bytes()).hexdigest(),
            "r8_binary_sha256":hashlib.sha256(Path(r8_binary).read_bytes()).hexdigest(),
            "host":"github-hosted-ubuntu; not consumer Windows",
            "limitations":"child RSS is a process maximum not a precise aggregate; energy not measured",
            "child_maxrss_kib_linux":ru_end.ru_maxrss,
        },
        "limitations":"Heuristic is nonlearning, not a strong conventional RL baseline; FourRooms episode may truncate before requested budget; probe feedback enables online adaptation",
    }
    output.write_text(json.dumps(record,sort_keys=True,separators=(",",":"),allow_nan=False)+"\n",encoding="utf-8")
    print(json.dumps({name:[{"seed":e["seed"],"goal":e["goal_reached"],"steps":e["steps"],"positions":e["diagnostic"]["unique_positions"],"advances":e["diagnostic"]["position_advances"]} for e in episodes] for name,episodes in results.items()},sort_keys=True))
    return record

if __name__=="__main__":
    p=argparse.ArgumentParser()
    p.add_argument("--outcome-binary",required=True)
    p.add_argument("--r8-binary",required=True)
    p.add_argument("--output",type=Path,required=True)
    p.add_argument("--run-source",action="store_true")
    args=p.parse_args()
    if not args.run_source:p.error("SOURCE_NOT_RUN")
    run(args.outcome_binary,args.r8_binary,args.output)
