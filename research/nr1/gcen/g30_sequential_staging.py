"""G30 bounded synthetic sequential shard loader; no real Qwen weights."""
import argparse
import hashlib
import json
import os
import resource
import tempfile
from pathlib import Path
from types import SimpleNamespace
import torch
from accelerate import init_empty_weights, dispatch_model
from accelerate.utils import set_module_tensor_to_device
from accelerate.utils.offload import offload_weight, save_offload_index
from safetensors import safe_open
from safetensors.torch import save_file
from transformers.models.qwen3_moe.modeling_qwen3_moe import Qwen3MoeMLP
from g26_sharded_checkpoint_smoke import test_runtime_identity

FILES = ("model-00001-of-00002.safetensors", "model-00002-of-00002.safetensors")
LAYOUT = {"gate_proj.weight": FILES[0], "up_proj.weight": FILES[0],
          "down_proj.weight": FILES[1]}
PLACEMENT = {"gate_proj": "cpu", "up_proj": "cpu", "down_proj": "disk"}
MAX_SYNTHETIC_SHARD_BYTES = 64 * 1024 * 1024


def verified_stage(src, dst, expected):
    if dst.exists() or not 0 < src.stat().st_size <= MAX_SYNTHETIC_SHARD_BYTES:
        raise ValueError("invalid staging precondition")
    h, count = hashlib.sha256(), 0
    try:
        with src.open("rb") as reader, dst.open("xb") as writer:
            while chunk := reader.read(65536):
                count += len(chunk)
                if count > MAX_SYNTHETIC_SHARD_BYTES:
                    raise ValueError("staging size bound")
                h.update(chunk)
                writer.write(chunk)
        if count != src.stat().st_size or h.hexdigest() != expected:
            raise ValueError("shard SHA256 mismatch")
        return count
    except BaseException:
        dst.unlink(missing_ok=True)
        raise


def allocation(root):
    return sum(p.stat().st_blocks * 512 for p in root.rglob("*") if p.is_file())


def run(output):
    if os.name != "posix":
        raise RuntimeError("Linux RSS measurement required")
    test_runtime_identity()
    torch.manual_seed(53)
    cfg = SimpleNamespace(hidden_size=16, intermediate_size=8, hidden_act="silu")
    reference = Qwen3MoeMLP(cfg).eval()
    x = torch.randn((3, 16), dtype=torch.float32)
    expected_acts, actual_acts = [], []

    def capture(target):
        def hook(_module, args):
            if len(args) != 1 or tuple(args[0].shape) != (3, 8):
                raise ValueError("unexpected expert activation")
            target.append(args[0].detach().float().cpu().clone())
        return hook

    h = reference.down_proj.register_forward_pre_hook(capture(expected_acts))
    with torch.inference_mode():
        expected = reference(x)
    h.remove()
    state = {k: v.detach().contiguous().cpu().clone()
             for k, v in reference.state_dict().items()}
    if set(state) != set(LAYOUT):
        raise RuntimeError("unexpected synthetic key inventory")
    baseline_rss = resource.getrusage(resource.RUSAGE_SELF).ru_maxrss * 1024
    with tempfile.TemporaryDirectory(prefix="gcen-g30-synthetic-") as temporary:
        root = Path(temporary)
        src, staged, offload = (root / name for name in ("source", "staging", "offload"))
        for p in (src, staged, offload):
            p.mkdir()
        save_file({k: state[k] for k in ("gate_proj.weight", "up_proj.weight")},
                  str(src / FILES[0]))
        save_file({"down_proj.weight": state["down_proj.weight"]}, str(src / FILES[1]))
        sha = {name: hashlib.sha256((src / name).read_bytes()).hexdigest()
               for name in FILES}
        # Negative control: an altered source file must be rejected and removed.
        broken = root / "corrupt.safetensors"
        original = (src / FILES[0]).read_bytes()
        broken.write_bytes(original[:-1] + bytes([original[-1] ^ 1]))
        rejected = False
        try:
            verified_stage(broken, staged / "rejected.safetensors", sha[FILES[0]])
        except ValueError as exc:
            rejected = "SHA256 mismatch" in str(exc)
        if not rejected or any(staged.iterdir()):
            raise RuntimeError("negative digest-control failure")
        broken.unlink()
        with init_empty_weights():
            model = Qwen3MoeMLP(cfg).eval()
        if not all(p.is_meta for p in model.parameters()):
            raise RuntimeError("meta-parameter precondition failed")
        offload_index, keys_loaded = {}, set()
        peak_staged, peak_scratch = 0, 0
        peak_staging_shard_count = 0
        staging_count_samples = []
        order = []
        for filename in FILES:
            before_count = len(tuple(staged.iterdir()))
            staging_count_samples.append(before_count)
            if before_count != 0:
                raise RuntimeError("staging overlap")
            stage = staged / filename
            try:
                peak_staged = max(peak_staged, verified_stage(src / filename, stage, sha[filename]))
                present = tuple(staged.iterdir())
                staged_count = len(present)
                staging_count_samples.append(staged_count)
                peak_staging_shard_count = max(peak_staging_shard_count, staged_count)
                if staged_count != 1 or present[0] != stage or not stage.is_file():
                    raise RuntimeError("unexpected staging inventory")
                peak_scratch = max(peak_scratch, allocation(staged) + allocation(offload))
                with safe_open(str(stage), framework="pt", device="cpu") as handle:
                    keys = list(handle.keys())
                    if set(keys) != {k for k, f in LAYOUT.items() if f == filename}:
                        raise RuntimeError("shard key mismatch")
                    for key in keys:
                        if key in keys_loaded:
                            raise RuntimeError("repeated tensor key")
                        tensor = handle.get_tensor(key)
                        if key == "down_proj.weight":
                            offload_index = offload_weight(tensor, key, str(offload), index=offload_index)
                        else:
                            set_module_tensor_to_device(model, key, "cpu", value=tensor)
                        keys_loaded.add(key)
                        del tensor
                peak_scratch = max(peak_scratch, allocation(staged) + allocation(offload))
                order.append(filename)
            finally:
                stage.unlink(missing_ok=True)
            cleaned_count = len(tuple(staged.iterdir()))
            staging_count_samples.append(cleaned_count)
            if cleaned_count != 0:
                raise RuntimeError("staged shard not cleaned")
        if keys_loaded != set(LAYOUT) or set(offload_index) != {"down_proj.weight"}:
            raise RuntimeError("incomplete weight admission")
        save_offload_index(offload_index, str(offload))
        if not (offload / "index.json").is_file():
            raise RuntimeError("missing offload index")
        model = dispatch_model(model, device_map=PLACEMENT, offload_dir=str(offload),
                               offload_index=offload_index, force_hooks=True)
        if not all(p.is_meta for p in model.down_proj.parameters()):
            raise RuntimeError("disk offload not active")
        h = model.down_proj.register_forward_pre_hook(capture(actual_acts))
        with torch.inference_mode():
            actual = model(x)
        h.remove()
        if len(expected_acts) != 1 or len(actual_acts) != 1:
            raise RuntimeError("expert hook not observed")
        if not torch.allclose(actual, expected, atol=1e-6, rtol=1e-5):
            raise RuntimeError("output parity failure")
        if not torch.allclose(actual_acts[0], expected_acts[0], atol=1e-6, rtol=1e-5):
            raise RuntimeError("expert-input parity failure")
        result = {
            "schema": "wingless.gcen.g30.synthetic-single-stage.v1",
            "status": "SYNTHETIC_SEQUENTIAL_STAGE_PARITY_PASS",
            "shard_order": order,
            "verified_shard_sha256": sha,
            "negative_corruption_control_pass": rejected,
            "peak_staging_shard_count": peak_staging_shard_count,
            "staging_count_samples": staging_count_samples,
            "staging_files_remaining": len(list(staged.iterdir())),
            "largest_staged_shard_bytes": peak_staged,
            "measured_peak_scratch_allocated_bytes": peak_scratch,
            "baseline_peak_process_rss_bytes": baseline_rss,
            "measured_peak_process_rss_bytes": resource.getrusage(resource.RUSAGE_SELF).ru_maxrss * 1024,
            "output_max_abs_error": float((actual - expected).abs().max()),
            "expert_hook_max_abs_error": float((actual_acts[0] - expected_acts[0]).abs().max()),
            "original_checkpoint_weights_loaded": False,
            "full_original_model_streaming": "UNVERIFIED",
            "scientific_acceptance": "NOT_GRANTED",
            "production_activation": False
        }
    out = Path(output)
    out.parent.mkdir(parents=True, exist_ok=True)
    with out.open("x", encoding="utf-8") as f:
        json.dump(result, f, sort_keys=True, indent=2)
        f.write("\n")
    print("GCEN_G30_SYNTHETIC_STAGE_PASS")


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--output", required=True)
    run(parser.parse_args().output)
