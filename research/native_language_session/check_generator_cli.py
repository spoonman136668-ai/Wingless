"""Invented engineering pairs only; no learned novelty or language proof."""
import json
import os
from pathlib import Path
import subprocess
import sys

binary = str(Path(sys.argv[1]).resolve())
env = {k: v for k, v in os.environ.items() if k in ("PATH", "SystemRoot", "WINDIR", "TEMP", "TMP", "HOME", "USERPROFILE")}


def run(requests):
    p = subprocess.run([binary], input="\n".join(json.dumps(x) for x in requests) + "\n",
        text=True, capture_output=True, check=True, timeout=60, env=env)
    rows = [json.loads(x) for x in p.stdout.splitlines()]
    assert len(rows) == len(requests), (rows, p.stderr)
    assert all(x["language_proof"] is False for x in rows)
    return rows


examples = [{"input": "aaaa", "response": "aA"}, {"input": "zzzz", "response": "bB"}]
rows = run([{"op": "train", "examples": examples}, {"op": "generate", "input": "aaaa"},
            {"op": "generate", "input": "zzzz"}, {"op": "checkpoint"},
            {"op": "generate", "input": ""}, {"op": "checkpoint"}])
assert all(rows[i]["ok"] for i in (0, 1, 2, 3, 5)), rows
assert not rows[4]["ok"]
for i, x in enumerate(examples):
    result = rows[i + 1]["result"]
    assert result["response"] == x["response"], rows
    assert result["whole_response_bank"] is False
    assert result["seen_training_response"] is True
checkpoint = rows[3]["result"]["checkpoint"]
assert checkpoint == rows[5]["result"]["checkpoint"]
fresh = run([{"op": "restore", "checkpoint": checkpoint}, {"op": "generate", "input": "aaaa"},
             {"op": "checkpoint"}, {"op": "train", "examples": examples}])
assert all(fresh[i]["ok"] for i in (0, 1, 2)), fresh
assert not fresh[3]["ok"]
assert fresh[1]["result"] == rows[1]["result"]
assert fresh[2]["result"]["checkpoint"] == checkpoint
print("WINGLESS_NATIVE_GENERATION_ENGINEERING_OBSERVATION=" + json.dumps({
    "training_pair_outputs": [rows[1]["result"], rows[2]["result"]],
    "training_pair_correct": 2, "training_pair_count": 2,
    "checkpoint_bytes_json_object": len(json.dumps(checkpoint, separators=(",", ":")).encode()),
    "learned_novelty_demonstrated": False, "scientific_language_evidence": False}, sort_keys=True))
print("WINGLESS_NATIVE_GENERATION_FRESH_PROCESS=PASS language_proof=false output=parallel-native-byte-field")
