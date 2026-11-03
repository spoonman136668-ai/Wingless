"""Invented engineering text pairs; no language-accuracy/transfer claim."""
import json
import os
from pathlib import Path
import subprocess
import sys

binary = str(Path(sys.argv[1]).resolve())
env = {k: v for k, v in os.environ.items() if k in (
    "PATH", "SystemRoot", "WINDIR", "TEMP", "TMP", "HOME", "USERPROFILE")}


def run(requests):
    p = subprocess.run([binary], input="\n".join(json.dumps(x) for x in requests) + "\n",
        text=True, capture_output=True, check=True, timeout=30, env=env)
    rows = [json.loads(x) for x in p.stdout.splitlines()]
    assert len(rows) == len(requests), (rows, p.stderr)
    assert all(x["language_proof"] is False for x in rows)
    return rows


examples = [{"input": "aaaa", "response": "amber"}, {"input": "zzzz", "response": "violet"}]
rows = run([{"op": "train", "examples": examples}, {"op": "predict", "input": "aaaa"},
            {"op": "predict", "input": "zzzz"}, {"op": "checkpoint"},
            {"op": "predict", "input": ""}, {"op": "checkpoint"}])
assert all(rows[i]["ok"] for i in (0, 1, 2, 3, 5)), rows
assert not rows[4]["ok"]
prediction = rows[1]["result"]
assert prediction["response"] in ("amber", "violet")
assert prediction["novel_string_generation"] is False
checkpoint = rows[3]["result"]["checkpoint"]
assert checkpoint == rows[5]["result"]["checkpoint"]
fresh = run([{"op": "restore", "checkpoint": checkpoint}, {"op": "predict", "input": "aaaa"},
             {"op": "checkpoint"}, {"op": "train", "examples": examples}])
assert all(fresh[i]["ok"] for i in (0, 1, 2)), fresh
assert not fresh[3]["ok"]
assert fresh[1]["result"] == prediction
assert fresh[2]["result"]["checkpoint"] == checkpoint
print("WINGLESS_NATIVE_TEXT_ENGINEERING_OBSERVATION=" + json.dumps({
    "training_pair_predictions": [rows[1]["result"], rows[2]["result"]],
    "training_pair_correct": sum(rows[i+1]["result"]["response"] == x["response"] for i, x in enumerate(examples)),
    "training_pair_count": len(examples), "scientific_language_evidence": False}, sort_keys=True))
print("WINGLESS_NATIVE_TEXT_FRESH_PROCESS=PASS language_proof=false output=closed-response-selection")
