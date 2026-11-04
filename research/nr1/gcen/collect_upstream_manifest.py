#!/usr/bin/env python3
"""Read-only, pinned-revision upstream digest inventory for GCEN development.

Run explicitly on an authorized network-enabled research host. This tool does
not download weight shards, run inference, or admit a scientific experiment.
"""
import argparse
import hashlib
import json
from pathlib import Path
from urllib.parse import quote
from urllib.request import urlopen

from custody import REQUIRED
from source_manifest import MODEL, REVISION, admit_source_manifest

BASE = "https://huggingface.co"
API = f"{BASE}/api/models/{MODEL}/revision/{REVISION}?blobs=true"
TREE = f"{BASE}/{MODEL}/tree/{REVISION}"
SMALL_FILES = {"config.json", "model.safetensors.index.json"}


def get_bytes(url, limit):
    with urlopen(url, timeout=40) as response:
        data = response.read(limit + 1)
    if len(data) > limit:
        raise ValueError("upstream response exceeds size limit")
    return data


def build_manifest(fetch=get_bytes):
    raw = fetch(API, 10 * 1024 * 1024)
    model_info = json.loads(raw)
    if not isinstance(model_info, dict) or model_info.get("sha") != REVISION:
        raise ValueError("upstream model revision mismatch")
    siblings = model_info.get("siblings")
    if not isinstance(siblings, list):
        raise ValueError("missing upstream sibling inventory")
    mapped = {}
    for item in siblings:
        if not isinstance(item, dict):
            raise ValueError("invalid upstream sibling entry")
        name = item.get("rfilename")
        if name in REQUIRED:
            if name in mapped:
                raise ValueError("duplicate upstream file: " + name)
            mapped[name] = item
    if set(mapped) != set(REQUIRED):
        raise ValueError("upstream inventory incomplete")
    digests = {}
    for name in REQUIRED:
        if name in SMALL_FILES:
            source = f"{BASE}/{MODEL}/resolve/{REVISION}/{quote(name)}"
            digests[name] = hashlib.sha256(fetch(source, 20 * 1024 * 1024)).hexdigest()
        else:
            lfs = mapped[name].get("lfs")
            if not isinstance(lfs, dict):
                raise ValueError("missing LFS provenance: " + name)
            digest = lfs.get("sha256")
            if not isinstance(digest, str) or len(digest) != 64 or any(ch not in "0123456789abcdef" for ch in digest):
                raise ValueError("invalid upstream LFS SHA256: " + name)
            digests[name] = digest
    document = {
        "model": MODEL,
        "revision": REVISION,
        "source": {"uri": TREE, "revision": REVISION},
        "files": digests,
    }
    admit_source_manifest(document)
    return document


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--output", required=True, help="New JSON path, never overwritten")
    args = parser.parse_args()
    doc = build_manifest()
    path = Path(args.output)
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("x", encoding="utf-8") as handle:
        json.dump(doc, handle, sort_keys=True, indent=2)
        handle.write("\n")


if __name__ == "__main__":
    main()
