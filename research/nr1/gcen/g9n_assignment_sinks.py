"""G9N: bounded host-only tensor assignment sinks for verified shard admission.

No model construction, forward pass, GPU/offload hook, checkpoint download, or
runtime promotion. A session owns its disk files and must be explicitly closed.
Default disk serialization uses one Safetensors file per tensor, not Accelerate
dispatch. Source-bound factory reconstructs G9B/G28 placement independently.
"""
import hashlib
import os
from pathlib import Path
import shutil
import tempfile

from g9n_source_contract import freeze_original_assignment

SCHEMA = "wingless.gcen.g9n.repository-source-placement.v1"
MAX_TENSORS = 18867


def _write_safetensor(path, key, tensor):
    """Heavy dependencies are deferred until a disk assignment is actually made."""
    from safetensors.torch import save_file
    save_file({key: tensor.contiguous()}, str(path))


class HostAssignmentSinks:
    """Own a private disk staging tree and fail closed on any partial assignment."""

    def __init__(self, plan, scratch_dir, *, writer=None):
        if not isinstance(plan, dict) or plan.get("schema") != SCHEMA:
            raise ValueError("source placement plan required")
        shapes, targets = plan.get("shapes"), plan.get("target_devices")
        if (not isinstance(shapes, dict) or not isinstance(targets, dict)
            or not shapes or len(shapes) > MAX_TENSORS or set(shapes) != set(targets)
            or plan.get("indexed_tensors") != len(shapes)):
            raise ValueError("source plan inventory mismatch")
        if (any(v not in ("cpu", "disk") for v in targets.values())
            or not any(v == "disk" for v in targets.values())):
            raise ValueError("original offload plan must include disk")
        if any(not isinstance(k, str) or not k or
               not isinstance(s, (tuple, list)) or not s or
               any(type(dim) is not int or dim <= 0 for dim in s)
               for k, s in shapes.items()):
            raise ValueError("invalid tensor identities/shapes")
        sums = {"cpu": 0, "disk": 0}
        for key, shape in shapes.items():
            n = 2
            for dim in shape:
                n *= dim
            sums[targets[key]] += n
        if (sums["cpu"] != plan.get("cpu_bytes") or
            sums["disk"] != plan.get("disk_bytes") or
            sum(sums.values()) != plan.get("original_tensor_bytes")):
            raise ValueError("source plan byte accounting mismatch")
        if writer is not None and not callable(writer):
            raise ValueError("invalid disk writer")
        root = Path(scratch_dir)
        if root.is_symlink() or not root.is_dir():
            raise ValueError("scratch directory missing or symlinked")
        root = root.resolve(strict=True)
        self._directory = Path(tempfile.mkdtemp(prefix="g9n-sinks-", dir=root))
        self._targets = dict(targets)
        self._shapes = {k: tuple(s) for k, s in shapes.items()}
        self._budgets = sums
        self._used = {"cpu": 0, "disk": 0}
        self._writer = writer or _write_safetensor
        self.cpu_tensors = {}
        self.disk_paths = {}
        self._closed = False

    @property
    def directory(self):
        return self._directory

    def assign(self, key, tensor):
        if self._closed:
            raise RuntimeError("assignment session closed")
        if key not in self._targets or key in self.cpu_tensors or key in self.disk_paths:
            raise ValueError("unknown or duplicate tensor")
        if tuple(getattr(tensor, "shape", ())) != self._shapes[key]:
            raise ValueError("tensor shape mismatch")
        if str(getattr(tensor, "dtype", None)) != "torch.bfloat16":
            raise ValueError("original tensor must be BF16")
        if getattr(getattr(tensor, "device", None), "type", None) != "cpu":
            raise ValueError("tensor must be staged on CPU")
        size = 2
        for dim in self._shapes[key]:
            size *= dim
        target = self._targets[key]
        if self._used[target] + size > self._budgets[target]:
            raise ValueError("tensor allocation exceeds plan")
        if target == "cpu":
            self.cpu_tensors[key] = tensor
        else:
            name = hashlib.sha256(key.encode("utf-8")).hexdigest() + ".safetensors"
            dest = self._directory / name
            temp = self._directory / (name + ".partial")
            if dest.exists() or temp.exists():
                raise RuntimeError("disk tensor filename collision")
            try:
                self._writer(temp, key, tensor)
                if temp.is_symlink() or not temp.is_file() or temp.stat().st_size <= 0:
                    raise RuntimeError("disk tensor writer did not produce a regular file")
                os.replace(temp, dest)
            finally:
                temp.unlink(missing_ok=True)
            self.disk_paths[key] = dest
        self._used[target] += size

    def receipt(self):
        if self._closed:
            raise RuntimeError("session closed")
        if set(self.cpu_tensors) | set(self.disk_paths) != set(self._targets):
            raise RuntimeError("incomplete tensor assignment")
        if self._used != self._budgets:
            raise RuntimeError("source tensor byte totals not reconciled")
        return {"schema": "wingless.gcen.g9n.host-sinks.v1",
                "cpu_tensors": len(self.cpu_tensors),
                "disk_tensors": len(self.disk_paths),
                "cpu_bytes": self._used["cpu"],
                "disk_bytes": self._used["disk"],
                "original_model_forward_executed": False,
                "accelerate_dispatch_verified": False,
                "gpu_residency_verified": False,
                "scientific_acceptance": "NOT_GRANTED"}

    def discard_partial(self):
        if self._closed:
            return
        self._closed = True
        self.cpu_tensors.clear()
        self.disk_paths.clear()
        shutil.rmtree(self._directory)

    def close(self):
        self.discard_partial()

    def __enter__(self):
        return self

    def __exit__(self, *_):
        self.close()


def source_bound_sinks(*, manifest, config, index, device_map,
                       cpu_budget_bytes, disk_budget_bytes, scratch_dir, writer=None):
    """Build session from verified source metadata, not a caller-forged plan."""
    plan = freeze_original_assignment(
        manifest=manifest, config=config, index=index,
        device_map=device_map, cpu_budget_bytes=cpu_budget_bytes,
        disk_budget_bytes=disk_budget_bytes)
    return HostAssignmentSinks(plan, scratch_dir, writer=writer)
