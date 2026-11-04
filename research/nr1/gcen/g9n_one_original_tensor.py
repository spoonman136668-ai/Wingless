"""G9N isolated original BF16 expert -> verified G9N disk sink experiment.

One SHA256-verified original shard, one full original expert tensor, disk only.
Does not create/dispatch Qwen, evaluate prompts, use heldout, or access GPU.
Explicitly separate source-custody and sink-integrity evidence.
"""
import argparse
import json
import os
from pathlib import Path
import shutil
import tempfile
import time
from urllib.request import Request, urlopen

from collect_upstream_manifest import build_manifest
from g31_original_shard_stage import (
    MAX_SHARD_BYTES, RESERVE_BYTES, SHARD, URL, check_upstream,
    inspect_header, parse_unique_json, peak_rss_bytes, stage_verified,
)
from g32_real_expert_scaling import select_experts
from g9n_assignment_sinks import HostAssignmentSinks
from out_of_core_preflight import pinned_metadata
from runtime_feasibility import MODEL, REVISION, check_manifest

MAX_TENSOR_BYTES = 64 * 1024**2
MAX_RSS_DELTA = 1024**3
RESULT = "G9N_ONE_VERIFIED_ORIGINAL_EXPERT_DISK_ROUNDTRIP"


def choose_target(header, weight_map):
    """Use the existing frozen G32 first-layer/first-expert selector."""
    selected = select_experts(header, weight_map)
    key = selected["names"][0]
    entry = header.get(key)
    if not isinstance(entry, dict) or entry.get("dtype") != "BF16":
        raise ValueError("missing original BF16 expert")
    shape = entry.get("shape")
    if (not isinstance(shape, list) or len(shape) != 2
        or any(type(x) is not int or x <= 0 for x in shape)
        or 2 * shape[0] * shape[1] > MAX_TENSOR_BYTES
        or weight_map.get(key) != SHARD):
        raise ValueError("selected original expert outside safe bounds")
    return key, tuple(shape), selected["layer"]


def run(manifest_path, scratch_root, output, source_sha, *, g9r_kernel=False):
    if (os.name != "posix" or type(source_sha) is not str
        or len(source_sha) != 40
        or any(c not in "0123456789abcdef" for c in source_sha)):
        raise ValueError("exact SHA and Linux resource telemetry required")
    if not scratch_root.is_dir() or output.exists():
        raise ValueError("missing scratch or already-existing receipt")
    manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    check_manifest(manifest)
    if shutil.disk_usage(scratch_root).free < MAX_SHARD_BYTES + RESERVE_BYTES:
        raise RuntimeError("not enough scratch space for one original shard")
    check_upstream(manifest, build_manifest())
    index = pinned_metadata("model.safetensors.index.json", manifest)
    started, rss0 = time.monotonic(), peak_rss_bytes()
    with tempfile.TemporaryDirectory(prefix="g9n-one-original-", dir=scratch_root) as temp:
        root = Path(temp)
        original = root / SHARD
        request = Request(URL, headers={"Accept-Encoding":"identity",
                                       "User-Agent":"Wingless-G9N-OneOriginalTensor/1"})
        with urlopen(request, timeout=90) as response:
            custody = stage_verified(response, original, manifest["files"][SHARD])
        details = inspect_header(original, index)
        with original.open("rb") as stream:
            count = int.from_bytes(stream.read(8), "little")
            if count != details["header_bytes"]:
                raise ValueError("original header drift")
            header = parse_unique_json(stream.read(count))
        key, shape, layer = choose_target(header, index["weight_map"])
        from safetensors import safe_open
        import torch
        plan = {
            "schema":"wingless.gcen.g9n.repository-source-placement.v1",
            "shapes":{key:shape}, "target_devices":{key:"disk"},
            "indexed_tensors":1, "original_tensor_bytes":2*shape[0]*shape[1],
            "cpu_bytes":0, "disk_bytes":2*shape[0]*shape[1],
        }
        with HostAssignmentSinks(plan, root) as sinks:
            with safe_open(str(original), framework="pt", device="cpu") as reader:
                if key not in reader.keys():
                    raise RuntimeError("selected tensor missing from verified original shard")
                tensor = reader.get_tensor(key)
                if tuple(tensor.shape) != shape or tensor.dtype != torch.bfloat16:
                    raise ValueError("original tensor payload shape/dtype mismatch")
                sinks.assign(key, tensor)
                receipt = sinks.verified_receipt()
                with safe_open(str(sinks.disk_paths[key]), framework="pt", device="cpu") as persisted:
                    if not torch.equal(tensor, persisted.get_tensor(key)):
                        raise RuntimeError("original BF16 tensor disk roundtrip mismatch")
                # Compare the authenticated full expert once, after an exact
                # disk round-trip. No model forward or original activation trace.
                from g9p_original_expert_comparison import compare_original_expert
                comparison=compare_original_expert(tensor.to(dtype=torch.float32).numpy())
                kernel_receipt=None
                if g9r_kernel:
                    from g9r_numba_packed_int2 import benchmark_original
                    kernel_receipt=benchmark_original(tensor.to(dtype=torch.float32).numpy())
                del tensor
            owned = sinks.directory
        if owned.exists():
            raise RuntimeError("G9N disk sink cleanup failed")
        if list(root.iterdir()) != [original]:
            raise RuntimeError("unexpected original/sink retained files")
        peak = peak_rss_bytes()
        if peak - rss0 > MAX_RSS_DELTA:
            raise RuntimeError("G9N one-tensor RSS increase above guard")
        evidence = {
            "schema":"wingless.gcen.g9n.one-original-tensor-development.v1",
            "result":RESULT, "source_sha":source_sha,
            "model":MODEL, "revision":REVISION,
            "original_shard":SHARD, "original_shard_sha256":manifest["files"][SHARD],
            "source_custody":custody,
            "selected_original_layer":layer,
            "selected_tensor":key, "selected_shape":list(shape),
            "selected_tensor_bytes":plan["disk_bytes"],
            "disk_integrity_receipt":receipt,
            "original_full_expert_operator_comparison":comparison,
            "g9r_packed_cpu_benchmark":kernel_receipt,
            "process_rss_baseline_bytes":rss0, "process_peak_rss_bytes":peak,
            "elapsed_seconds":time.monotonic()-started,
            "one_original_shard_authenticated":True,
            "all_16_shards_authenticated":False,
            "model_constructed":False, "full_model_forward":False,
            "gpu_vram_measured":False, "authentic_model_activations":False,
            "scientific_acceptance":"NOT_GRANTED",
        }
    evidence["source_and_sink_files_after_cleanup"] = 0
    output.parent.mkdir(parents=True,exist_ok=True)
    with output.open("x",encoding="utf-8") as f:
        json.dump(evidence,f,indent=2,sort_keys=True)
        f.write("\n")
    print("G9N_ONE_ORIGINAL_TENSOR_BF16_ROUNDTRIP_PASS")


if __name__ == "__main__":
    p=argparse.ArgumentParser()
    p.add_argument("--manifest",type=Path,required=True)
    p.add_argument("--scratch-root",type=Path,required=True)
    p.add_argument("--output",type=Path,required=True)
    p.add_argument("--source-sha",required=True)
    p.add_argument("--g9r-kernel",action="store_true")
    a=p.parse_args()
    run(a.manifest,a.scratch_root,a.output,a.source_sha,g9r_kernel=a.g9r_kernel)
