"""Prospective native goal-progress versus motion, original and linear-Q stateful evaluation.

Development-only. Preregisters a new 256-transition training cohort and
holds out a NEW preregistered evaluation cohort (never v3/v4 test seeds).
Native feedback remains enabled DURING each test episode. Every seed starts
from its own byte-identical cold restoration of the same TRAINED checkpoint.
Never silently promote online evaluation feedback to training or accepted
generalization; stateful test-time adaptation is disclosed separately.
"""
from __future__ import annotations
import argparse
import hashlib
import json
import resource
import time
from pathlib import Path

import exact_exposure_gate as frozen
import exact_exposure_v2 as v2
from ipc_source_runner import BridgeError
from linear_q_crossfamily import LinearQ
from motion_source_compare import diagnostics
from reward_ablation_restart import MaskedOutcome

TRAIN_SEEDS=tuple(5101+2*i for i in range(256))
TEST_SEEDS=(6101,6113,6121,6133,6151,6163,6173,6197)
TRAIN_STEPS=256
TEST_STEPS=100
PREVIOUS_HELDOUT=frozenset((229,233,239,241,251,257,263,269,
                             2609,2617,2621,2633,2647,2657,2663,2671,
                             2801,2803,2819,2833,2843,2851,2861,2879,
                             3701,3719,3733,3743,3761,3779,3793,3803,
                             4703,4721,4733,4747,4759,4783,4793,4801))
assert len(set(TEST_SEEDS))==8 and len(set(TRAIN_SEEDS))==256
assert not(set(TEST_SEEDS)&(set(TRAIN_SEEDS)|PREVIOUS_HELDOUT))

def configure_training():
    # Private process monkeypatch only: all v1/v2/v3/v4 source files immutable.
    v2.TRAIN_SEEDS=TRAIN_SEEDS
    frozen.exact_train=lambda policy:v2.exact_train_v2(policy)

def assert_update_delta(prior,after,steps,arm,seed):
    if type(prior) is not int or type(after) is not int or type(steps) is not int:
        raise BridgeError("noninteger evaluation accounting: "+arm)
    if prior!=TRAIN_STEPS or steps<1 or steps>TEST_STEPS or after-prior!=steps:
        raise BridgeError("unequal or missing evaluation feedback: "+arm+" "+str(seed))
    return True

def evaluate_native_stateful(binary,snapshot,arm):
    records=[]
    for seed in TEST_SEEDS:
        policy=MaskedOutcome(binary,False)
        try:
            reply=policy.bridge.request("load",state=snapshot)
            prior=reply.get("updates")
            if prior!=TRAIN_STEPS:
                raise BridgeError("cold native checkpoint update count: "+arm)
            saved=policy.bridge.request("snapshot")["state"]
            if saved!=snapshot:
                raise BridgeError("cold restoration changed checkpoint: "+arm)
            policy.next_id=TRAIN_STEPS+1
            episode=frozen.run_episode(frozen.TEST_ENV,policy,seed,TEST_STEPS)
            current=policy.bridge.request("snapshot")
            after=current.get("updates")
            assert_update_delta(prior,after,episode["steps"],arm,seed)
            if policy.next_id!=after+1:
                raise BridgeError("native feedback event-ID conservation: "+arm)
            changed=current.get("state")!=saved
            if not changed:
                raise BridgeError("native evaluation feedback did not persist: "+arm)
            records.append({"seed":seed,"episode":episode,"updates_before":prior,
                            "updates_after":after,"online_updates":after-prior,
                            "restored_checkpoint_sha256":hashlib.sha256(saved.encode()).hexdigest(),
                            "post_episode_checkpoint_sha256":hashlib.sha256(current["state"].encode()).hexdigest(),
                            "checkpoint_changed_during_online_eval":changed})
        finally:
            policy.close()
    return records

def evaluate_linear_stateful(weights):
    records=[]
    for seed in TEST_SEEDS:
        q=LinearQ(seed=20261027+seed)
        q.weights=[row[:] for row in weights]
        q.training=True
        before=[row[:] for row in q.weights]
        prior=q.updates
        episode=frozen.run_episode(frozen.TEST_ENV,q,seed,TEST_STEPS)
        after=q.updates
        if prior!=0 or after!=episode["steps"] or q.armed is not None:
            raise BridgeError("linear-Q online feedback accounting failed")
        records.append({"seed":seed,"episode":episode,"updates_before":prior,
                        "updates_after":after,"online_updates":after,
                        "weights_changed_during_online_eval":q.weights!=before})
    return records

def verify_matched_eval(arms):
    for index,seed in enumerate(TEST_SEEDS):
        starts=[]
        for name,arm in arms.items():
            row=arm["evaluation"][index]
            if row["seed"]!=seed or row["episode"]["seed"]!=seed:
                raise BridgeError("evaluation seed drift: "+name)
            if row["online_updates"]!=row["episode"]["steps"]:
                raise BridgeError("unmatched evaluation updates: "+name)
            starts.append(row["episode"]["events"][0]["before"])
        if any(x!=starts[0] for x in starts[1:]):
            raise BridgeError("unmatched public initial observation at seed "+str(seed))
    return True

def execute(original_binary,motion_binary,progress_binary,output):
    if output.exists():raise FileExistsError("refuse to overwrite experiment evidence")
    start=time.perf_counter()
    configure_training()
    arms={}
    for name,binary in (("progress_native",progress_binary),("motion_native",motion_binary),("original_native",original_binary)):
        training,snapshot,digest=frozen.train_native(binary,False)
        if sum(e["steps"] for e in training)!=TRAIN_STEPS:
            raise BridgeError("native training exposure failed: "+name)
        online=evaluate_native_stateful(binary,snapshot,name)
        arms[name]={"training":training,
                    "trained_checkpoint_sha256":digest,
                    "evaluation":online,
                    "training_diagnostic":diagnostics(training),
                    "online_eval_diagnostic":diagnostics([x["episode"] for x in online])}
    q_train,weights,digest=frozen.train_q()
    if sum(e["steps"] for e in q_train)!=TRAIN_STEPS:
        raise BridgeError("linear-Q training exposure failed")
    q_online=evaluate_linear_stateful(weights)
    arms["linear_q"]={"training":q_train,
                      "trained_checkpoint_sha256":digest,
                      "evaluation":q_online,
                      "training_diagnostic":diagnostics(q_train),
                      "online_eval_diagnostic":diagnostics([x["episode"] for x in q_online])}
    verify_matched_eval(arms)
    result={
        "schema":"wingless.itd02.progress-online-eval.v3",
        "science":"NOT_ACCEPTED",
        "scope":"PROSPECTIVE_GOAL_PROGRESS_STATEFUL_EVAL_WITH_TEST_TIME_UPDATES",
        "train_family":frozen.TRAIN_ENV,"eval_family":frozen.TEST_ENV,
        "train_seeds":list(TRAIN_SEEDS),"heldout_eval_seeds":list(TEST_SEEDS),
        "train_source_interactions_per_arm":TRAIN_STEPS,
        "max_eval_steps_per_seed":TEST_STEPS,
        "evaluation_protocol":"COLD_RESTORE_PER_SEED_WITH_ONLINE_SELECTED_FEEDBACK",
        "test_time_adaptation":True,
        "arms":arms,
        "resource":{"host":"github-hosted-ubuntu-not-retail-PC",
                    "wall_seconds":time.perf_counter()-start,
                    "parent_peak_kib":resource.getrusage(resource.RUSAGE_SELF).ru_maxrss,
                    "child_peak_kib":resource.getrusage(resource.RUSAGE_CHILDREN).ru_maxrss,
                    "energy_measured":False},
        "limitations":"Exploratory within-episode adaptation is not frozen zero-shot transfer; 256 training transitions per learner is small; per-seed cold resets prevent cross-seed leakage; linear-Q not independently tuned; no independent science acceptance or retail-device performance claim."
    }
    output.write_text(json.dumps(result,sort_keys=True,separators=(",",":"),allow_nan=False)+"\n",encoding="utf-8")
    print(json.dumps({"science":result["science"],
      "training":{k:v["training_diagnostic"] for k,v in arms.items()},
      "online_eval":{k:v["online_eval_diagnostic"] for k,v in arms.items()},
      "test_time_updates_per_seed":{k:[x["online_updates"] for x in v["evaluation"]] for k,v in arms.items()}},sort_keys=True))
    return result

if __name__=="__main__":
    parser=argparse.ArgumentParser()
    parser.add_argument("--original-binary",required=True)
    parser.add_argument("--motion-binary",required=True)
    parser.add_argument("--progress-binary",required=True)
    parser.add_argument("--output",required=True,type=Path)
    parser.add_argument("--run-source",action="store_true")
    a=parser.parse_args()
    if not a.run_source:parser.error("SOURCE_NOT_RUN")
    execute(a.original_binary,a.motion_binary,a.progress_binary,a.output)
