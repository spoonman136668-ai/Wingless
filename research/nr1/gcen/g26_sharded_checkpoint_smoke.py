"""G26 tiny two-shard Qwen3-MoE checkpoint load-and-dispatch smoke.

The two shards contain ONLY synthetic weights from a 16-by-8 fake expert.
Pins G25 upstream runtime and verifies actual Safetensors CPU/disk dispatch.
No original Qwen checkpoint, task data or genuine model activation is accessed.
"""
import hashlib
import json
from pathlib import Path
import tempfile
from types import SimpleNamespace

import torch
from safetensors.torch import save_file
from accelerate import init_empty_weights, load_checkpoint_and_dispatch
from transformers.models.qwen3_moe.modeling_qwen3_moe import Qwen3MoeMLP

from g25_disk_offload_smoke import (
    MODEL, REVISION, TORCH_VERSION, TRANSFORMERS_VERSION,
    ACCELERATE_VERSION, TRANSFORMERS_QWEN3_MOE_GIT_BLOB,
    ACCELERATE_BIG_MODELING_GIT_BLOB, git_blob,
)


def test_runtime_identity():
    import transformers
    import accelerate
    import transformers.models.qwen3_moe.modeling_qwen3_moe as original_source
    import accelerate.big_modeling as original_offload
    if (torch.__version__.split("+")[0] != TORCH_VERSION
        or transformers.__version__ != TRANSFORMERS_VERSION
        or accelerate.__version__ != ACCELERATE_VERSION):
        raise RuntimeError("G26 runtime version identity mismatch")
    if (git_blob(original_source.__file__) != TRANSFORMERS_QWEN3_MOE_GIT_BLOB
        or git_blob(original_offload.__file__) != ACCELERATE_BIG_MODELING_GIT_BLOB):
        raise RuntimeError("G26 original runtime source mismatch")
    if torch.cuda.is_available():
        raise RuntimeError("G26 CPU-only fixture required")


def run(output_path):
    test_runtime_identity()
    torch.manual_seed(53)
    config = SimpleNamespace(hidden_size=16, intermediate_size=8, hidden_act="silu")
    reference = Qwen3MoeMLP(config).eval()
    data = torch.randn((3, 16), dtype=torch.float32)
    expected_hooks = []
    observed_hooks = []
    def hook(target):
        def capture(_module, inputs):
            if len(inputs) != 1 or inputs[0].shape != (3, 8):
                raise RuntimeError("unexpected down-projection activation shape")
            target.append(inputs[0].detach().float().cpu().clone())
        return capture
    pre = reference.down_proj.register_forward_pre_hook(hook(expected_hooks))
    with torch.inference_mode():
        expected = reference(data)
    pre.remove()
    if len(expected_hooks) != 1:
        raise RuntimeError("original synthetic expert hook missing")

    with tempfile.TemporaryDirectory(prefix="gcen-g26-tiny-sharded-") as tmp:
        root = Path(tmp)
        source = root / "synthetic_checkpoint"
        offload = root / "offload"
        source.mkdir()
        offload.mkdir()
        state = {key: value.detach().cpu().contiguous().clone()
                 for key, value in reference.state_dict().items()}
        if set(state) != {"gate_proj.weight", "up_proj.weight", "down_proj.weight"}:
            raise RuntimeError("unexpected pinned Qwen3-MoE expert parameter layout")
        filenames = ("model-00001-of-00002.safetensors",
                     "model-00002-of-00002.safetensors")
        parts = (
            {key: state[key] for key in ("gate_proj.weight", "up_proj.weight")},
            {"down_proj.weight": state["down_proj.weight"]},
        )
        for name, shard in zip(filenames, parts):
            save_file(shard, str(source / name))
        mapping = {
            "gate_proj.weight": filenames[0],
            "up_proj.weight": filenames[0],
            "down_proj.weight": filenames[1],
        }
        bytes_total = sum(t.numel() * t.element_size() for t in state.values())
        index = {"metadata": {"total_size": bytes_total},
                 "weight_map": mapping}
        with (source / "model.safetensors.index.json").open("x", encoding="utf-8") as stream:
            json.dump(index, stream, sort_keys=True)
        digest = {name: hashlib.sha256((source / name).read_bytes()).hexdigest()
                  for name in filenames}
        with init_empty_weights():
            empty = Qwen3MoeMLP(config).eval()
        if not all(p.is_meta for p in empty.parameters()):
            raise RuntimeError("empty-model precondition violated")
        loaded = load_checkpoint_and_dispatch(
            empty,
            checkpoint=str(source),
            device_map={
                "gate_proj": "cpu",
                "up_proj": "cpu",
                "down_proj": "disk",
            },
            offload_folder=str(offload),
            offload_state_dict=True,
            force_hooks=True,
            dtype=torch.float32,
            strict=True,
        )
        if not any(p.is_meta for p in loaded.down_proj.parameters()):
            raise RuntimeError("disk-dispatched down projection was not offloaded")
        if not (offload / "index.json").is_file():
            raise RuntimeError("missing disk-offload weight index")
        if not any(offload.glob("*.dat")):
            raise RuntimeError("missing disk-offload weight data")
        recorder = loaded.down_proj.register_forward_pre_hook(hook(observed_hooks))
        with torch.inference_mode():
            actual = loaded(data)
        recorder.remove()
        if len(observed_hooks) != 1:
            raise RuntimeError("disk-dispatched down-projection hook missing")
        if not torch.allclose(actual, expected, atol=1e-6, rtol=1e-5):
            raise RuntimeError("G26 sharded load changed output")
        if not torch.allclose(observed_hooks[0], expected_hooks[0], atol=1e-6, rtol=1e-5):
            raise RuntimeError("G26 sharded load changed expert input activation")
        max_output_error = float((actual - expected).abs().max().item())
        max_hook_error = float((observed_hooks[0] - expected_hooks[0]).abs().max().item())

    receipt = {
        "schema": "wingless.gcen.g26.tiny-sharded-checkpoint-hook.v1",
        "result": "PINNED_TINY_TWO_SHARD_DISK_DISPATCH_PASS",
        "source_model": MODEL,
        "source_revision": REVISION,
        "tiny_model_shape": [16, 8],
        "synthetic_checkpoint_shard_count": 2,
        "synthetic_original_parameter_bytes": bytes_total,
        "synthetic_shard_sha256": digest,
        "disk_device_map": {"gate_proj": "cpu", "up_proj": "cpu", "down_proj": "disk"},
        "max_output_abs_error": max_output_error,
        "max_expert_input_abs_error": max_hook_error,
        "real_checkpoint_shards_downloaded": 0,
        "original_model_executed": False,
        "authentic_original_activation_trace": False,
        "full_model_disk_dispatch": "UNVERIFIED",
        "source_authenticity": "PINNED_UPSTREAM_CODE_WITH_SYNTHETIC_TENSORS_ONLY",
        "scientific_acceptance": "NOT_GRANTED",
    }
    out = Path(output_path)
    out.parent.mkdir(parents=True, exist_ok=True)
    with out.open("x", encoding="utf-8") as stream:
        json.dump(receipt, stream, indent=2, sort_keys=True)
        stream.write("\n")
    print("GCEN_G26_PINNED_TINY_SHARDED_CHECKPOINT_PASS")


if __name__ == "__main__":
    import argparse
    p = argparse.ArgumentParser()
    p.add_argument("--output", required=True)
    args = p.parse_args()
    run(args.output)
