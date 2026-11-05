"""One-shot G9N CPU BF16 Safetensors round-trip. Synthetic data ONLY."""
import hashlib
import json
import tempfile
from pathlib import Path

import torch
from safetensors import safe_open
from safetensors.torch import save_file
from g9b_tensor_budget import MODEL, REVISION
from g9n_shard_admission import SHARDS, admit_local_shards
from g9n_assignment_sinks import HostAssignmentSinks


def fixture(root):
    src, scratch = root/"src", root/"scratch"
    src.mkdir()
    scratch.mkdir()
    hashes = {"config.json":"0"*64, "model.safetensors.index.json":"1"*64}
    shapes, mapping, expected = {}, {}, {}
    for i, shard in enumerate(SHARDS):
        key = f"synthetic.{i}.weight"
        value = torch.arange(4, dtype=torch.float32).reshape(2, 2).add(i).to(torch.bfloat16)
        save_file({key:value}, str(src/shard))
        hashes[shard] = hashlib.sha256((src/shard).read_bytes()).hexdigest()
        shapes[key], mapping[key], expected[key] = (2,2), shard, value
    targets = {k:("cpu" if i<2 else "disk") for i,k in enumerate(shapes)}
    plan = {"schema":"wingless.gcen.g9n.repository-source-placement.v1",
            "shapes":shapes, "target_devices":targets,
            "indexed_tensors":16, "original_tensor_bytes":128,
            "cpu_bytes":16, "disk_bytes":112}
    return src,scratch,{"model":MODEL,"revision":REVISION,"files":hashes}, {"weight_map":mapping},expected,plan


def ingest(src,scratch,manifest,index,expected,sinks):
    return admit_local_shards(manifest=manifest,index=index,
             expected_shapes={k:tuple(v.shape) for k,v in expected.items()},
             source_dir=src,scratch_dir=scratch,
             assign_tensor=sinks.assign,discard_partial=sinks.discard_partial)


def check_roundtrip():
    with tempfile.TemporaryDirectory() as d:
        src,scratch,manifest,index,expected,plan=fixture(Path(d))
        with HostAssignmentSinks(plan,scratch) as sinks:
            result=ingest(src,scratch,manifest,index,expected,sinks)
            assert (result["verified_shards"],result["verified_tensor_count"])==(16,16)
            rec=sinks.verified_receipt()
            assert rec["disk_content_verified"]
            assert (rec["cpu_tensors"],rec["disk_tensors"])==(2,14)
            for k,v in sinks.cpu_tensors.items():
                assert torch.equal(v,expected[k])
            for k,p in sinks.disk_paths.items():
                with safe_open(str(p),framework="pt",device="cpu") as f:
                    assert list(f.keys())==[k] and torch.equal(f.get_tensor(k),expected[k])
            owned=sinks.directory
        assert not owned.exists()
        assert not list(scratch.glob("g9n-*"))


def check_corrupt_stored_disk_tensor():
    with tempfile.TemporaryDirectory() as d:
        src,scratch,manifest,index,expected,plan=fixture(Path(d))
        with HostAssignmentSinks(plan,scratch) as sinks:
            ingest(src,scratch,manifest,index,expected,sinks)
            assert sinks.verified_receipt()["disk_content_verified"]
            path=next(iter(sinks.disk_paths.values()))
            mutated=bytearray(path.read_bytes())
            mutated[-1]^=1
            path.write_bytes(mutated)
            try:
                sinks.verified_receipt()
            except RuntimeError as e:
                assert "integrity mismatch" in str(e), str(e)
            else:
                raise AssertionError("corrupted sink disk tensor was accepted")
            sinks.discard_partial()
            assert not sinks.directory.exists()
        assert not list(scratch.glob("g9n-*"))


def check_corrupt_tail():
    with tempfile.TemporaryDirectory() as d:
        src,scratch,manifest,index,expected,plan=fixture(Path(d))
        p=src/SHARDS[-1]
        raw=bytearray(p.read_bytes())
        raw[-1]^=1
        p.write_bytes(raw)
        with HostAssignmentSinks(plan,scratch) as sinks:
            try:
                ingest(src,scratch,manifest,index,expected,sinks)
            except ValueError as e:
                assert "digest mismatch" in str(e), str(e)
            else:
                raise AssertionError("corrupt last shard accepted")
            assert not sinks.directory.exists() and not sinks.cpu_tensors and not sinks.disk_paths
            assert not list(scratch.glob("g9n-*"))


def check_writer_failure():
    with tempfile.TemporaryDirectory() as d:
        src,scratch,manifest,index,expected,plan=fixture(Path(d))
        counter=[0]
        def writer(path,key,value):
            counter[0]+=1
            save_file({key:value.contiguous()},str(path))
            if counter[0]==2:
                raise OSError("injected failure after actual disk write")
        with HostAssignmentSinks(plan,scratch,writer=writer) as sinks:
            try:
                ingest(src,scratch,manifest,index,expected,sinks)
            except OSError as e:
                assert "injected failure" in str(e), str(e)
            else:
                raise AssertionError("disk failure accepted")
            assert counter[0]==2 and not sinks.directory.exists()
            assert not sinks.cpu_tensors and not sinks.disk_paths
            assert not list(scratch.glob("g9n-*"))


if __name__=="__main__":
    check_roundtrip()
    check_corrupt_stored_disk_tensor()
    check_corrupt_tail()
    check_writer_failure()
    print("G9N_SYNTHETIC_BF16_REAL_IO_PASS "+json.dumps({
        "verified_shards":16,"verified_tensors":16,
        "negative_controls":["corrupt_last_shard","partial_disk_write","corrupt_stored_disk_tensor"],
        "source":"synthetic_only","disk_content_verified":True,"model_forward":False,
        "gpu_vram_measured":False,"scientific_acceptance":"NOT_GRANTED"},sort_keys=True))
