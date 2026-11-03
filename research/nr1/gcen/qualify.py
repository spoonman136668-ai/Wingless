#!/usr/bin/env python3
"""GCEN cross-platform qualification entrypoint; never runs pretrained measurements.
Runs explicit GCEN focused tests followed by mandatory Wingless Go regression/build.
Invocation from repository root: python research/nr1/gcen/qualify.py
"""
import pathlib
import subprocess
import sys

ROOT = pathlib.Path(__file__).resolve().parents[3]

def run(label, command):
    print(f"GCEN_QUAL_START {label}", flush=True)
    subprocess.run(command, cwd=ROOT, check=True)
    print(f"GCEN_QUAL_PASS {label}", flush=True)

def main():
    run("focused", [sys.executable, str(ROOT / "research/nr1/gcen/run_focused.py")])
    run("full-go-regression", ["go", "test", "./...", "-count=1"])
    run("full-go-build", ["go", "build", "./..."])
    print("GCEN_QUAL_ALL_PASS", flush=True)

if __name__ == "__main__":
    main()
