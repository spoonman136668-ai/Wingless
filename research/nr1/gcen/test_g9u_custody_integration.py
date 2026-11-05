"""Dependency-isolated positive and failure-path tests; not authentic model evidence."""
import hashlib
import json
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch
import types
import numpy as np
import g9u_verified_loader as loader
from g31_original_shard_stage import SHARD
from runtime_feasibility import MODEL, REVISION

class FakeTensor:
    shape = (2048, 768)
    dtype = 'bf16'
    def to(self, *, dtype):
        return types.SimpleNamespace(numpy=lambda: np.ones(self.shape, dtype=np.float32))

class FakeReader:
    def __init__(self, tensor): self.tensor=tensor
    def __enter__(self): return self
    def __exit__(self, *a): return False
    def get_tensor(self, name): return self.tensor

class CustodyIntegration(unittest.TestCase):
    def setUp(self):
        self.root=tempfile.TemporaryDirectory()
        self.addCleanup(self.root.cleanup)
        self.shard=Path(self.root.name)/SHARD
        self.shard.write_bytes(b'\x02\x00\x00\x00\x00\x00\x00\x00{}payload')
        self.digest=hashlib.sha256(self.shard.read_bytes()).hexdigest()
        files={f'model-{i:05d}-of-00016.safetensors':'a'*64 for i in range(1,17)}
        files.update({'config.json':'a'*64, 'model.safetensors.index.json':'a'*64})
        files[SHARD]=self.digest
        self.manifest=Path(self.root.name)/'manifest.json'
        self.manifest.write_text(json.dumps({'model':MODEL,'revision':REVISION,'files':files}))
        self.name='model.layers.0.mlp.experts.0.down_proj.weight'
        self.index={'weight_map':{self.name:SHARD}}
        self.tensor=FakeTensor()
        self.patches=[
            patch.object(loader,'PINNED_SHARD_SHA256',self.digest),
            patch.object(loader,'pinned_metadata',return_value=self.index),
            patch.object(loader,'inspect_header',return_value={}),
            patch.object(loader,'parse_unique_json',return_value={self.name:{}}),
            patch.object(loader,'choose_target',return_value=(self.name,(2048,768),0)),
        ]
        for p in self.patches:
            p.start(); self.addCleanup(p.stop)
        import sys
        fake_st=types.SimpleNamespace(safe_open=lambda *a,**kw:FakeReader(self.tensor))
        fake_torch=types.SimpleNamespace(bfloat16='bf16',float32='fp32')
        p1=patch.dict(sys.modules,{'safetensors':fake_st,'torch':fake_torch})
        p1.start(); self.addCleanup(p1.stop)

    def test_positive_mocked_source(self):
        weight,receipt=loader.load_verified_original_expert(self.shard,self.manifest)
        self.assertEqual(weight.shape,(2048,768))
        self.assertEqual(receipt['shard_sha256'],self.digest)
        self.assertTrue(receipt['original_shard_authenticated'])

    def test_invalid_bf16_rejected(self):
        self.tensor.dtype='fp32'
        with self.assertRaisesRegex(ValueError,'dtype/shape'):
            loader.load_verified_original_expert(self.shard,self.manifest)

    def test_invalid_tensor_selection_rejected(self):
        with patch.object(loader,'choose_target',return_value=(self.name,(2048,768),1)):
            with self.assertRaisesRegex(ValueError,'selection'):
                loader.load_verified_original_expert(self.shard,self.manifest)

    def test_mutation_detected_after_read(self):
        class MutatingReader(FakeReader):
            def get_tensor(inner,name):
                self.shard.write_bytes(self.shard.read_bytes()+b'changed')
                return inner.tensor
        with patch.dict(__import__('sys').modules,{'safetensors':types.SimpleNamespace(safe_open=lambda *a,**kw:MutatingReader(self.tensor))}):
            with self.assertRaisesRegex(ValueError,'changed during'):
                loader.load_verified_original_expert(self.shard,self.manifest)

if __name__=='__main__': unittest.main()
