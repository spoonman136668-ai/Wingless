"""Fresh-process JSONL engineering check; no language/scientific evaluation."""
import json
import os
from pathlib import Path
import subprocess
import sys


binary = str(Path(sys.argv[1]).resolve())
env = {k: v for k, v in os.environ.items() if k in (
    "PATH", "SystemRoot", "WINDIR", "TEMP", "TMP", "HOME", "USERPROFILE")}


def run(requests):
    p = subprocess.run([binary], input="\n".join(
        x if isinstance(x, str) else json.dumps(x) for x in requests) + "\n",
        text=True, capture_output=True, check=True, timeout=30, env=env)
    rows = [json.loads(x) for x in p.stdout.splitlines()]
    assert len(rows) == len(requests), (rows, p.stderr)
    assert all(x["language_capability"] is False for x in rows)
    return rows


edges = [{"from": i, "to": i + 1} for i in range(1, 13)]
case = {"edges": edges, "questions": [[1, 5], [5, 9]], "truth": [True, True]}
rows = run([
    {"op": "train", "examples": [case], "leaky": False},
    {"op": "observe", "edges": edges},
    {"op": "predict", "from": 1, "to": 5}, {"op": "checkpoint"},
    {"op": "observe", "edges": []}, {"op": "checkpoint"},
    {"op": "train", "examples": [case], "leaky": False},
    {"op": "checkpoint"},
])
assert all(rows[i]["ok"] for i in (0, 1, 2, 3, 5, 7)), rows
assert not rows[4]["ok"] and not rows[6]["ok"], rows
checkpoint = rows[3]["result"]["checkpoint"]
assert checkpoint == rows[5]["result"]["checkpoint"] == rows[7]["result"]["checkpoint"]
restored = run([
    {"op": "restore", "checkpoint": checkpoint},
    {"op": "predict", "from": 1, "to": 5}, {"op": "checkpoint"},
    {"op": "restore", "checkpoint": checkpoint},
])
assert all(restored[i]["ok"] for i in (0, 1, 2)), restored
assert restored[1]["result"] == rows[2]["result"], restored
assert restored[2]["result"]["checkpoint"] == checkpoint
assert not restored[3]["ok"]
assert not run(['{"op":"checkpoint","op":"train"}'])[0]["ok"]
assert not run([{"op": "train", "examples": [case], "leaky": None}])[0]["ok"]
bad_case = dict(case, questions=[[1, 5, 9], [5, 9]])
assert not run([{"op": "train", "examples": [bad_case], "leaky": False}])[0]["ok"]
print("WINGLESS_NATIVE_SESSION_FRESH_PROCESS=PASS language_capability=false")
