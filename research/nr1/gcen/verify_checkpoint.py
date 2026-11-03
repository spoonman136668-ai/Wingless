#!/usr/bin/env python3
"""Offline GCEN checkpoint byte verification. Never downloads weights or runs inference.

The source manifest must first be independently authenticated by the research custodian.
A successful receipt proves local bytes match the supplied manifest, NOT upstream authority.
"""
import argparse
import hashlib
import json
from pathlib import Path

from custody import REQUIRED, verify_manifest
from source_manifest import admit_source_manifest


def verify_checkpoint(root, document):
    digests = admit_source_manifest(document)
    verify_manifest(root, digests)
    return {
        "status": "LOCAL_BYTES_VERIFIED_SOURCE_AUTHORITY_UNVERIFIED",
        "model": document["model"],
        "revision": document["revision"],
        "source": document["source"],
        "files": [{"name": name, "sha256": digests[name]} for name in REQUIRED],
        "manifest_sha256": hashlib.sha256(
            json.dumps(document, sort_keys=True, separators=(",", ":")).encode("utf-8")
        ).hexdigest(),
    }


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--checkpoint-root", required=True)
    parser.add_argument("--source-manifest", required=True)
    parser.add_argument("--receipt", required=True)
    args = parser.parse_args()
    manifest = json.loads(Path(args.source_manifest).read_text(encoding="utf-8"))
    receipt = verify_checkpoint(args.checkpoint_root, manifest)
    output = Path(args.receipt)
    if output.exists():
        raise ValueError("refusing to overwrite an existing custody receipt")
    output.parent.mkdir(parents=True, exist_ok=True)
    with output.open("x", encoding="utf-8") as stream:
        json.dump(receipt, stream, sort_keys=True, indent=2)
        stream.write("\n")


if __name__ == "__main__":
    main()
