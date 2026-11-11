"""Prospective opt-in motion vs original contextual R8 vs linear-Q comparison.

Development/exploration only. No earlier v3 evidence is modified. Source-world
interaction requires explicit --run-source. Frozen clone-per-choice is a
preference diagnostic, NOT a full stateful navigation or scientific transfer test.
"""
from __future__ import annotations
import argparse
import hashlib
import json
from pathlib import Path
import resource
import time

import exact_exposure_gate as frozen
import exact_exposure_v2 as v2
from ipc_source_runner import BridgeError

TRAIN_SEEDS = tuple(1701+2*i for i in range(256))
TEST_SEEDS = (2609,2617,2621,2633,2647,2657,2663,2671)
TRAIN_STEPS=256
TEST_STEPS=100
PREVIOUS_V3_TRAIN = frozenset(range(1009,1520,2))
PREVIOUS_V3_TEST = frozenset((229,233,239,241,251,257,263,269))
assert len(set(TRAIN_SEEDS))==256 and len(set(TEST_SEEDS))==8
assert not(set(TRAIN_SEEDS)&(PREVIOUS_V3_TRAIN|PREVIOUS_V3_TEST|set(TEST_SEEDS)))
assert not(set(TEST_SEEDS)&(PREVIOUS_V3_TRAIN|PREVIOUS_V3_TEST))

def diagnostics(episodes):
    events=[e for episode in episodes for e in episode["events"]]
    counts={str(a):sum(e["action"]==a for e in events) for a in range(3)}
    advances=sum((e["before"]["x"],e["before"]["y"])!=(e["after"]["x"],e["after"]["y"]) for e in events)
    return {"selected_action_counts":counts, "position_advances":advances,
            "goal_count":sum(e["goal_reached"] for e in episodes),
            "returns":[e["return"] for e in episodes],
            "episodes":len(episodes), "authentic_steps":len(events)}

def execute(original_binary,motion_binary,output):
    if output.exists():raise FileExistsError("refuse to overwrite experimental evidence")
    start=time.perf_counter()
    # Only process-local values; the historical source files and outputs stay frozen.
    v2.TRAIN_SEEDS=TRAIN_SEEDS
    frozen.exact_train=lambda p:v2.exact_train_v2(p)
    frozen.TEST_SEEDS=TEST_SEEDS
    arms={}
    for name,binary in (("motion_native",motion_binary),("original_native",original_binary)):
        records,snapshot,digest=frozen.train_native(binary,False)
        if sum(e["steps"] for e in records)!=TRAIN_STEPS:
            raise BridgeError("native real exposure not exact: "+name)
        evaluated=frozen.evaluate_native(binary,snapshot)
        arms[name]={"training":records,"checkpoint_sha256":digest,
                    "evaluation":evaluated,
                    "training_diagnostic":diagnostics(records),
                    "frozen_evaluation_diagnostic":diagnostics(evaluated)}
    records,weights,digest=frozen.train_q()
    if sum(e["steps"] for e in records)!=TRAIN_STEPS:raise BridgeError("linear-Q source exposure mismatch")
    evaluated=frozen.evaluate_q(weights)
    arms["linear_q"]={"training":records,"checkpoint_sha256":digest,
                     "evaluation":evaluated,"training_diagnostic":diagnostics(records),
                     "frozen_evaluation_diagnostic":diagnostics(evaluated)}
    for i in range(len(TEST_SEEDS)):
        views=[arms[k]["evaluation"][i]["events"][0]["before"] for k in arms]
        if any(x!=views[0] for x in views):
            raise BridgeError("unmatched held-out public start")
    r={
      "schema":"wingless.itd02.optin-motion-exact-source.v1",
      "science":"NOT_ACCEPTED",
      "scope":"PROSPECTIVE_EXPLORATORY_FROZEN_PREFERENCES_NOT_STATEFUL_TRANSFER",
      "source_family":frozen.TRAIN_ENV,"evaluation_family":frozen.TEST_ENV,
      "train_seeds":list(TRAIN_SEEDS),"eval_seeds":list(TEST_SEEDS),
      "exposure_per_arm":TRAIN_STEPS,"eval_steps_per_seed":TEST_STEPS,
      "evaluation_mode":"FROZEN_NATIVE_CLONE_PER_CHOICE_VERSUS_FROZEN_LINEAR_Q",
      "arms":arms,
      "resource":{"host":"github-hosted-ubuntu-not-retail",
                  "wall_seconds":time.perf_counter()-start,
                  "parent_peak_kib":resource.getrusage(resource.RUSAGE_SELF).ru_maxrss,
                  "child_peak_kib":resource.getrusage(resource.RUSAGE_CHILDREN).ru_maxrss,
                  "limitations":"non-simultaneous memory peaks, energy not measured"},
      "limitations":"Distinct predeclared seed starts, trajectories may differ; full stateful feedback omitted at frozen evaluation; 256 training transitions is small; conventional linear-Q is not independently optimized; exploratory not independent scientific admission."
    }
    output.write_text(json.dumps(r,sort_keys=True,separators=(",",":"),allow_nan=False)+"\n",encoding="utf-8")
    print(json.dumps({"training":{k:v["training_diagnostic"] for k,v in arms.items()},
       "frozen_eval":{k:v["frozen_evaluation_diagnostic"] for k,v in arms.items()},
       "science":r["science"]},sort_keys=True))
    return r

if __name__=="__main__":
    parser=argparse.ArgumentParser()
    parser.add_argument("--original-binary",required=True)
    parser.add_argument("--motion-binary",required=True)
    parser.add_argument("--output",required=True,type=Path)
    parser.add_argument("--run-source",action="store_true")
    args=parser.parse_args()
    if not args.run_source:parser.error("SOURCE_NOT_RUN")
    execute(args.original_binary,args.motion_binary,args.output)
