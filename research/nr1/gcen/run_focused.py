#!/usr/bin/env python3
"""Explicit GCEN focused suite for qualified Windows/Linux research runners.
Run from repository root: python research/nr1/gcen/run_focused.py
Requires numpy and safetensors installed in the qualified environment.
"""
import subprocess
import sys
from pathlib import Path
here = Path(__file__).resolve().parent
cmd = [sys.executable, '-m', 'unittest', 'discover', '-s', str(here), '-p', 'test_*.py', '-v']
raise SystemExit(subprocess.call(cmd, cwd=str(here)))
