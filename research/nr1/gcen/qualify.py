#!/usr/bin/env python3
"""GCEN development-only qualification: focused tests plus unchanged native suite.

Requires qualified Windows/Linux runner with Python, numpy, safetensors, Go,
and PowerShell (pwsh). Does not download or measure checkpoint weights.
Run from the repository root: python research/nr1/gcen/qualify.py
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
    # Reuse the existing native contract rather than a weaker Go-only substitute:
    # race-enabled tests, ICE validate/build/revalidate, host evidence, Go build.
    run("native-regression", ["pwsh", "-NoProfile", "-File", str(ROOT / "scripts/Test-Native.ps1")])
    print("GCEN_QUAL_ALL_PASS", flush=True)


if __name__ == "__main__":
    main()
