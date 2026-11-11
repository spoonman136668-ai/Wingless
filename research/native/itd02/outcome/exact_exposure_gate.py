"""Equal *actual source transitions* MiniGrid cross-family gate.

Fresh seeds, 256 train steps per arm. Evaluations are frozen, with native
peek_frozen clone-per-choice and unchanging linear-Q weights; this is a
bounded preference comparison, not a full-stateful deployment comparison.
"""
from __future__ import annotations
import argparse
import hashlib
import json
import time
import resource
from pathlib import Path
from frozen_source_probe import FrozenPolicy
from ipc_source_runner import BridgeError
from linear_q_crossfamily import LinearQ, run_episode
from reward_ablation_restart import MaskedOutcome

TRAIN_ENV="MiniGrid-Empty-Random-5x5-v0"
TEST_ENV="MiniGrid-FourRooms-v0"
TRAIN_SEEDS=(181,191,193,197,199,211,223,227,271,277,281,283)
TEST_SEEDS=(229,233,239,241,251,257,263,269)
TRAIN_STEPS=256
TEST_STEPS=100

def exact_train(policy):
    remaining=TRAIN_STEPS
    records=[]
    for seed in TRAIN_SEEDS:
        if remaining<=0:
            break
        record=run_episode(TRAIN_ENV,policy,seed,min(remaining,100))
        records.append(record)
        remaining-=record["steps"]
        # A source episode whose budget expired without real termination is
        # never followed by another source reset. Only the final train cut
        # is permitted to leave the policy in an unfinished world episode.
        if not record["completed"] and remaining>0:
            raise BridgeError("incomplete source episode before train boundary")
    if remaining!=0:
        raise BridgeError("failed to reach exact authentic source training budget")
    return records

def train_native(binary,blind):
    policy=MaskedOutcome(binary,blind)
    try:
        records=exact_train(policy)
        if policy.next_id!=TRAIN_STEPS+1:
            raise BridgeError("native feedback did not match training exposure")
        saved=policy.bridge.request("snapshot")["state"]
        return records,saved,hashlib.sha256(saved.encode()).hexdigest()
    finally:
        policy.close()

def train_q():
    q=LinearQ(seed=20261027)
    records=exact_train(q)
    if q.updates!=TRAIN_STEPS:
        raise BridgeError("linear Q missing training updates")
    weights=[r[:] for r in q.weights]
    return records,weights,hashlib.sha256(json.dumps(weights).encode()).hexdigest()

def evaluate_native(binary,saved):
    results=[]
    for seed in TEST_SEEDS:
        learner=MaskedOutcome(binary,False)
        try:
            reply=learner.bridge.request("load",state=saved)
            if reply.get("updates")!=TRAIN_STEPS:
                raise BridgeError("native cold restoration update mismatch")
            before=learner.bridge.request("snapshot")["state"]
            if before!=saved:
                raise BridgeError("native cold restoration bytes mismatch")
            result=run_episode(TEST_ENV,FrozenPolicy(learner.bridge),seed,TEST_STEPS)
            if learner.bridge.request("snapshot")["state"]!=saved:
                raise BridgeError("frozen native state changed during evaluation")
            results.append(result)
        finally:
            learner.close()
    return results

def evaluate_q(weights):
    results=[]
    for seed in TEST_SEEDS:
        q=LinearQ()
        q.weights=[w[:] for w in weights]
        q.training=False
        result=run_episode(TEST_ENV,q,seed,TEST_STEPS)
        if q.weights!=weights or q.updates!=0:
            raise BridgeError("frozen conventional learning state changed")
        results.append(result)
    return results

def execute(binary,output):
    if output.exists():
        raise FileExistsError("cannot overwrite evidence")
    start=time.perf_counter()
    trained={}
    evaluated={}
    for name,blind in (("reward_aware",False),("reward_blind",True)):
        records,saved,digest=train_native(binary,blind)
        trained[name]={"episodes":records,"state_sha256":digest}
        evaluated[name]=evaluate_native(binary,saved)
    q_records,weights,digest=train_q()
    trained["linear_q"]={"episodes":q_records,"state_sha256":digest}
    evaluated["linear_q"]=evaluate_q(weights)
    for label,data in trained.items():
        if sum(x["steps"] for x in data["episodes"])!=TRAIN_STEPS:
            raise BridgeError("exact source exposure mismatch: "+label)
    for i in range(len(TEST_SEEDS)):
        start_views=[evaluated[k][i]["events"][0]["before"] for k in evaluated]
        if start_views.count(start_views[0])!=len(start_views):
            raise BridgeError("unequal test public start")
    if len({json.dumps(x["events"][0]["before"],sort_keys=True)
            for x in evaluated["linear_q"]})<3:
        raise BridgeError("insufficient test public start diversity")
    r={
      "schema":"wingless.itd02.exact-exposure-frozen.v1",
      "science":"NOT_ACCEPTED",
      "scope":"EXPLORATORY_EQUAL_SOURCE_INTERACTIONS_FROZEN_PREFERENCES",
      "training_environment":TRAIN_ENV,"evaluation_environment":TEST_ENV,
      "train_seeds_allowed":list(TRAIN_SEEDS),"test_seeds":list(TEST_SEEDS),
      "exact_train_interactions_per_arm":TRAIN_STEPS,
      "training_interactions":{k:sum(x["steps"] for x in v["episodes"]) for k,v in trained.items()},
      "eval_update_mode":{"reward_aware":"FROZEN_CLONE_PER_ACTION","reward_blind":"FROZEN_CLONE_PER_ACTION","linear_q":"FROZEN_WEIGHTS"},
      "train":trained,"evaluation":evaluated,
      "resource":{"host":"github-hosted-ubuntu-not-retail-PC",
                  "wall_seconds":time.perf_counter()-start,
                  "parent_peak_kib":resource.getrusage(resource.RUSAGE_SELF).ru_maxrss,
                  "child_peak_kib":resource.getrusage(resource.RUSAGE_CHILDREN).ru_maxrss,
                  "limitations":"peak fields are not simultaneous aggregate footprint; energy not measured"},
      "limitations":"Unoptimized conventional linear Q is not a strong RL reference; frozen native clone-per-action omits route-state; training seed schedule starts are same but policies visit different source trajectories; this is exploratory and not independent science acceptance."
    }
    output.write_text(json.dumps(r,sort_keys=True,separators=(",",":"),allow_nan=False)+"\n",encoding="utf-8")
    print(json.dumps({"exposure":r["training_interactions"],
      "results":{name:{"goals":sum(e["goal_reached"] for e in results),
                       "returns":[e["return"] for e in results]}
                 for name,results in evaluated.items()}},sort_keys=True))
    return r

if __name__=="__main__":
    p=argparse.ArgumentParser()
    p.add_argument("--native-binary",required=True)
    p.add_argument("--output",required=True,type=Path)
    p.add_argument("--run-source",action="store_true")
    a=p.parse_args()
    if not a.run_source:p.error("SOURCE_NOT_RUN")
    execute(a.native_binary,a.output)
