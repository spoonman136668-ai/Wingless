"""G25 pinned Qwen3-MoE expert offload/hook smoke test.

Runs ONLY on an explicitly selected one-shot research CI runner with fixed
CPU-only Torch, Transformers and Accelerate versions. Uses synthetic tiny
weights/inputs, NOT the Qwen checkpoint and NOT a scientifically admitted trace.
"""
import copy
import hashlib
import json
from pathlib import Path
import tempfile
from types import SimpleNamespace

import torch
import transformers
import accelerate
from accelerate import disk_offload
from transformers.models.qwen3_moe.modeling_qwen3_moe import Qwen3MoeMLP
from runtime_feasibility import MODEL, REVISION

TORCH_VERSION = "2.10.0"
TRANSFORMERS_VERSION = "4.57.3"
ACCELERATE_VERSION = "1.10.1"
TRANSFORMERS_COMMIT = "47b0e478f324b54f177ea7998a0791870fdd0324"
ACCELERATE_COMMIT = "623cc0ba588d192c77758939a6e340adf07edd4f"
TRANSFORMERS_QWEN3_MOE_GIT_BLOB = "25c738744faba82db3c451b6987b01db74fadcfd"
ACCELERATE_BIG_MODELING_GIT_BLOB = "2bace0f5615f37c84c1c3e77fe116081ed691bac"


def git_blob(file_path):
    data = Path(file_path).read_bytes()
    return hashlib.sha1(b"blob " + str(len(data)).encode("ascii") + b"\0" + data).hexdigest()


def run(output_path):
    versions = {"torch": torch.__version__.split("+")[0],
                "transformers": transformers.__version__,
                "accelerate": accelerate.__version__}
    if versions != {"torch": TORCH_VERSION, "transformers": TRANSFORMERS_VERSION,
                    "accelerate": ACCELERATE_VERSION}:
        raise RuntimeError("pinned CPU runtime version mismatch")
    import transformers.models.qwen3_moe.modeling_qwen3_moe as model_source
    import accelerate.big_modeling as offload_source
    blobs = {"transformers": git_blob(model_source.__file__),
             "accelerate": git_blob(offload_source.__file__)}
    if blobs != {"transformers": TRANSFORMERS_QWEN3_MOE_GIT_BLOB,
                 "accelerate": ACCELERATE_BIG_MODELING_GIT_BLOB}:
        raise RuntimeError("installed source differs from pinned upstream Git blobs")
    if torch.cuda.is_available():
        raise RuntimeError("CPU-only development qualification expected")

    torch.manual_seed(37)
    config = SimpleNamespace(hidden_size=16, intermediate_size=8, hidden_act="silu")
    reference = Qwen3MoeMLP(config).eval()
    candidate = copy.deepcopy(reference).eval()
    inputs = torch.randn(3, config.hidden_size, dtype=torch.float32)
    reference_captures = []
    candidate_captures = []

    def record(target):
        def hook(module, inputs):
            if len(inputs) != 1 or inputs[0].ndim != 2:
                raise RuntimeError("unexpected expert down_proj input shape")
            target.append(inputs[0].detach().float().cpu().clone())
        return hook

    baseline_hook = reference.down_proj.register_forward_pre_hook(record(reference_captures))
    with torch.inference_mode():
        expected = reference(inputs)
    baseline_hook.remove()

    with tempfile.TemporaryDirectory(prefix="gcen-g25-tiny-offload-") as dirname:
        folder = Path(dirname)
        offloaded = disk_offload(candidate, folder, execution_device=torch.device("cpu"))
        if not (folder / "index.json").is_file():
            raise RuntimeError("disk-backed weight index not materialized")
        if not any(p.is_meta for p in offloaded.parameters()):
            raise RuntimeError("no parameters offloaded to meta device")
        hook = offloaded.down_proj.register_forward_pre_hook(record(candidate_captures))
        with torch.inference_mode():
            actual = offloaded(inputs)
        hook.remove()
        if not any(folder.glob("*.dat")):
            raise RuntimeError("disk-backed weight data not present")
    if len(reference_captures) != 1 or len(candidate_captures) != 1:
        raise RuntimeError("expert down-projection hook missing")
    if reference_captures[0].shape != (3, 8) or candidate_captures[0].shape != (3, 8):
        raise RuntimeError("unexpected hooked activation dimensions")
    max_output_difference = float(torch.max(torch.abs(actual - expected)).item())
    max_activation_difference = float(
        torch.max(torch.abs(reference_captures[0] - candidate_captures[0])).item()
    )
    if not torch.allclose(actual, expected, rtol=1e-5, atol=1e-6):
        raise RuntimeError("disk-offload changed tiny Qwen3-MoE expert output")
    if not torch.allclose(candidate_captures[0], reference_captures[0], rtol=1e-5, atol=1e-6):
        raise RuntimeError("disk-offload changed expert input captured by forward pre-hook")

    result = {
        "schema": "wingless.gcen.g25.tiny-offload-hook-smoke.v1",
        "result": "PINNED_TINY_SYNTHETIC_CPU_OFFLOAD_HOOK_PASS",
        "package_versions": versions,
        "source_git_blobs": blobs,
        "source_commits": {"transformers": TRANSFORMERS_COMMIT, "accelerate": ACCELERATE_COMMIT},
        "original_model": MODEL,
        "original_checkpoint_revision": REVISION,
        "tiny_shape": {"hidden_size": 16, "intermediate_size": 8, "tokens": 3},
        "max_output_difference": max_output_difference,
        "max_activation_difference": max_activation_difference,
        "full_original_model_executed": False,
        "original_model_weights_loaded": False,
        "authentic_expert_activations": False,
        "production_inference_qualified": False,
        "scientific_acceptance": "NOT_GRANTED",
        "further_requirements": "full pinned checkpoint load, resource feasibility and authentic original forward-hook provenance",
    }
    path = Path(output_path)
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("x", encoding="utf-8") as stream:
        json.dump(result, stream, sort_keys=True, indent=2)
        stream.write("\n")
    print("GCEN_G25_PINNED_TINY_OFFLOAD_HOOK_PASS")


if __name__ == "__main__":
    import argparse
    parser = argparse.ArgumentParser()
    parser.add_argument("--output", required=True)
    args = parser.parse_args()
    run(args.output)
