"""ITD-02 predeclared cross-family native-versus-linear-Q source evaluation.

A reproducible conventional linear Q learner receives only the exact same six
source-public fields and real rewards as native Wingless. All arms train on
the same 8 source seeds and evaluate on the same 8 *new* FourRooms seeds.
Neither a strong-baseline superiority claim nor scientific admission follows:
the conventional hyperparameters are fixed, not independently optimized.
"""
from __future__ import annotations

import argparse
import copy
import hashlib
import json
import math
from pathlib import Path
import random
import resource
import time

from ipc_source_runner import BridgeError
from minigrid_source import open_source_environment
from reward_ablation_restart import MaskedOutcome
from source_controls import episode

TRAIN_ENV = "MiniGrid-Empty-Random-5x5-v0"
TEST_ENV = "MiniGrid-FourRooms-v0"
TRAIN_SEEDS = (101, 103, 107, 109, 113, 127, 131, 137)
TEST_SEEDS = (139, 149, 151, 157, 163, 167, 173, 179)
TRAIN_BUDGET = 100
TEST_BUDGET = 100

def features(o):
    dx=o["goal_x"]-o["x"]
    dy=o["goal_y"]-o["y"]
    z=[1., max(-1.,min(1.,dx/16.)),max(-1.,min(1.,dy/16.)),
       max(-1.,min(1.,(abs(dx)+abs(dy))/32.))]
    z.extend(1. if o["direction"]==d else 0. for d in range(4))
    z.append(float(o["front_blocked"]))
    return z

class LinearQ:
    """Small conventional Q-learning agent, no oracle; deterministic replay."""
    def __init__(self, seed=20261027, alpha=.10, gamma=.90, epsilon=.20):
        self.weights=[[0.]*9 for _ in range(3)]
        self.rng=random.Random(seed)
        self.alpha=alpha
        self.gamma=gamma
        self.epsilon=epsilon
        self.armed=None
        self.updates=0
        self.training=True

    def scores(self,o):
        z=features(o)
        return [sum(w*f for w,f in zip(row,z)) for row in self.weights]

    def choose(self,before):
        if self.armed is not None:
            raise BridgeError("linear Q has unresolved chosen action")
        q=self.scores(before)
        action=max(range(3),key=lambda a:q[a])
        if self.training and self.rng.random()<self.epsilon:
            action=self.rng.randrange(3)
        self.armed=(dict(before),action)
        return action

    def feedback(self,before,action,after,reward,terminal,truncated):
        if self.armed!=(before,action):
            raise BridgeError("linear Q accepted mismatched source transition")
        if not math.isfinite(reward):
            raise BridgeError("linear Q rejected nonfinite reward")
        if self.training:
            z=features(before)
            current=self.scores(before)[action]
            target=reward if terminal or truncated else reward+self.gamma*max(self.scores(after))
            err=max(-1.,min(1.,target-current))
            self.weights[action]=[max(-32.,min(32.,w+self.alpha*err*v))
                                  for w,v in zip(self.weights[action],z)]
            self.updates+=1
        self.armed=None

    def reset(self):
        if self.armed is not None:
            raise BridgeError("linear Q episode reset with pending selection")

    def close(self):
        if self.armed is not None:
            raise BridgeError("linear Q closed with pending selection")

def run_episode(environment, policy, seed, budget):
    world=open_source_environment(environment)
    try:
        result=episode(world,policy,seed,budget)
        return result
    finally:
        world.environment.close()

def train_native(binary,blind):
    policy=MaskedOutcome(binary,blind)
    try:
        episodes=[run_episode(TRAIN_ENV,policy,seed,TRAIN_BUDGET) for seed in TRAIN_SEEDS]
        if sum(x["steps"] for x in episodes)!=policy.next_id-1:
            raise BridgeError("source step accounting failed")
        saved=policy.bridge.request("snapshot")["state"]
        digest=hashlib.sha256(saved.encode("ascii")).hexdigest()
        return episodes,saved,digest
    finally:
        policy.close()

def evaluate_native(binary,saved,trained_steps,blind):
    results=[]
    for seed in TEST_SEEDS:
        policy=MaskedOutcome(binary,blind)
        try:
            reply=policy.bridge.request("load",state=saved)
            if reply.get("updates")!=trained_steps:
                raise BridgeError("cold native restore lost update count")
            if policy.bridge.request("snapshot")["state"]!=saved:
                raise BridgeError("native checkpoint changed on restore")
            policy.next_id=trained_steps+1
            results.append(run_episode(TEST_ENV,policy,seed,TEST_BUDGET))
        finally:
            policy.close()
    return results

def train_linear():
    policy=LinearQ()
    episodes=[run_episode(TRAIN_ENV,policy,seed,TRAIN_BUDGET) for seed in TRAIN_SEEDS]
    if policy.updates!=sum(e["steps"] for e in episodes):
        raise BridgeError("linear Q training step count mismatch")
    original=[row[:] for row in policy.weights]
    return episodes,policy,original

def evaluate_linear(policy,weights):
    results=[]
    for seed in TEST_SEEDS:
        # Unlike native online probes, this evaluates *frozen* learned weights
        # with no reward updates, while still processing selected feedback custody.
        clone=LinearQ()
        clone.weights=[r[:] for r in weights]
        clone.training=False
        results.append(run_episode(TEST_ENV,clone,seed,TEST_BUDGET))
        if clone.weights!=weights or clone.armed is not None:
            raise BridgeError("conventional frozen Q weights changed")
    return results

def run(binary, output):
    if output.exists():
        raise FileExistsError("no overwriting scientific source artifacts")
    start=time.perf_counter()
    trained={}
    evaluated={}
    source_train=None
    for name,blind in (("reward_aware",False),("reward_blind",True)):
        episodes,saved,sha=train_native(binary,blind)
        trace=[[(step["before"],step["action"]) for step in e["events"]] for e in episodes]
        if source_train is None:
            source_train=trace
        else:
            # Trajectories may legitimately diverge after reward-based updates,
            # but they must start from the same public observation per seed.
            for i,record in enumerate(episodes):
                if record["events"][0]["before"]!=trained["reward_aware"]["episodes"][i]["events"][0]["before"]:
                    raise BridgeError("unequal training public starts")
        trained[name]={"episodes":episodes,"checkpoint_sha256":sha}
        evaluated[name]=evaluate_native(binary,saved,sum(e["steps"] for e in episodes),blind)
    conventional_training,policy,weights=train_linear()
    trained["linear_q"]={"episodes":conventional_training,
                         "weight_sha256":hashlib.sha256(json.dumps(weights).encode()).hexdigest(),
                         "updates":policy.updates}
    evaluated["linear_q"]=evaluate_linear(policy,weights)
    for i in range(len(TEST_SEEDS)):
        observations=[evaluated[k][i]["events"][0]["before"] for k in evaluated]
        if not all(x==observations[0] for x in observations):
            raise BridgeError("unequal unfamiliar public observations")
    if len({json.dumps(x["events"][0]["before"],sort_keys=True)
            for x in evaluated["linear_q"]})<3:
        raise BridgeError("insufficient unfamiliar public initial diversity")
    record={
        "schema":"wingless.itd02.linear-q-cross-family.v1",
        "science":"NOT_ACCEPTED",
        "scope":"EXPLORATORY_NEW_SEEDS_ONLINE_NATIVE_VS_FROZEN_CONVENTIONAL",
        "training_environment":TRAIN_ENV,"evaluation_environment":TEST_ENV,
        "training_seeds":list(TRAIN_SEEDS),"evaluation_seeds":list(TEST_SEEDS),
        "per_episode_step_limit":100,
        "native_binary_sha256":hashlib.sha256(Path(binary).read_bytes()).hexdigest(),
        "linear_q_hyperparameters":{"alpha":.10,"gamma":.90,"epsilon_train":.20,
                                    "seed":20261027,"weights":27},
        "trained":trained,"evaluated":evaluated,
        "resource":{"host":"GitHub-hosted Ubuntu, not consumer Windows",
                    "wall_seconds":time.perf_counter()-start,
                    "python_parent_maxrss_kib_linux":resource.getrusage(resource.RUSAGE_SELF).ru_maxrss,
                    "child_maxrss_kib_linux":resource.getrusage(resource.RUSAGE_CHILDREN).ru_maxrss,
                    "limitations":"Linux maxima not combined peak; CPU/energy and retail hardware not fully measured"},
        "comparison_limitations":"Linear Q is a legitimate unoptimized conventional learning control, not a strong optimized benchmark. Native probes keep learning during evaluation while linear Q remains frozen; no claims of matched update policy or scientific superiority.",
    }
    output.write_text(json.dumps(record,sort_keys=True,separators=(",",":"),allow_nan=False)+"\n",encoding="utf-8")
    print(json.dumps({k:{"goals":sum(e["goal_reached"] for e in v),
                          "mean_return":sum(e["return"] for e in v)/len(v)}
                      for k,v in evaluated.items()},sort_keys=True))
    return record

if __name__=="__main__":
    p=argparse.ArgumentParser()
    p.add_argument("--native-binary",required=True)
    p.add_argument("--output",type=Path,required=True)
    p.add_argument("--run-source",action="store_true")
    args=p.parse_args()
    if not args.run_source:
        p.error("SOURCE_NOT_RUN")
    run(args.native_binary,args.output)
