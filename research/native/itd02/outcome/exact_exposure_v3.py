"""Prospective v3 source-exposure design, NOT a rerun of failed frozen v1/v2.

Predeclared independent training starts: seeds 1009, 1011, ..., 1519.
Each counts only authentic transitions, never manufactured padding. A source
reset is allowed only after real terminal/truncated feedback; incomplete
episodes fail closed. Held-out test seeds and 256-step per-arm target unchanged.
"""
from __future__ import annotations
import argparse
from pathlib import Path
import exact_exposure_gate as frozen
import exact_exposure_v2 as v2

TRAIN_SEEDS_V3 = tuple(1009 + 2*i for i in range(256))
assert len(TRAIN_SEEDS_V3)==256 and len(set(TRAIN_SEEDS_V3))==256
assert not (set(TRAIN_SEEDS_V3) & set(frozen.TEST_SEEDS))

def main():
    parser=argparse.ArgumentParser()
    parser.add_argument("--native-binary",required=True)
    parser.add_argument("--output",required=True,type=Path)
    parser.add_argument("--run-source",action="store_true")
    args=parser.parse_args()
    if not args.run_source:
        parser.error("SOURCE_NOT_RUN")
    # Monkeypatch only process-local module values. Original checked-in
    # experiments and their failure evidence remain byte-identical.
    v2.TRAIN_SEEDS=TRAIN_SEEDS_V3
    frozen.TRAIN_SEEDS=TRAIN_SEEDS_V3
    frozen.exact_train=lambda policy:v2.exact_train_v2(policy)
    frozen.execute(args.native_binary,args.output)

if __name__=="__main__":
    main()
