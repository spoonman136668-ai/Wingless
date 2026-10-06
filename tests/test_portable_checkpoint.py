import importlib.util
import json
import pathlib
import tempfile
import unittest

ROOT = pathlib.Path(__file__).resolve().parents[1]
SPEC = importlib.util.spec_from_file_location("portable_checkpoint", ROOT / "scripts" / "portable_checkpoint.py")
pc = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(pc)

PROJECT = "Wingless"
SOURCE_COMMIT = "72f9da765f7d358797e3cbbdd2bb356561b58e02"

def spec():
    return {
        "schema": "research.portable-checkpoint-export.v1",
        "project": PROJECT,
        "mechanism_id": "portable-smoke-mechanism",
        "source_commit_sha": SOURCE_COMMIT,
        "source_tree_sha": "1"*40,
        "source_schema_sha256": "2"*64,
        "qualification_contract_sha256": "3"*64,
        "toolchain_id": "portable-checkpoint-test-v1",
        "lineage": {"experiment_id":"EXP-PORTABLE-SMOKE", "parent_checkpoint_sha256":None},
        "resource_envelope": {
            "comparison":"equal-or-lower-than-baseline",
            "baseline_id":"portable-baseline-v1",
            "max_parameters":16,
            "max_context_bytes":4096,
            "max_model_calls":0,
        },
        "interface": {
            "signature": "state_former(raw_input, history) -> fixed_size_state",
            "fixed_state_dimension": 8,
            "fixed_readout_capacity": 1,
            "capacity_growth": False,
        },
        "files":[
            {"path":"model/weights.bin","role":"learned_parameters"},
            {"path":"state/resume.json","role":"resumable_state"},
            {"path":"replay/input.json","role":"replay_input"},
            {"path":"replay/output.json","role":"replay_output"},
        ],
        "replay":{"vectors":[{"id":"v1","input_path":"replay/input.json","output_path":"replay/output.json"}]},
        "authority":{
            "accepted_state_mutation":False,
            "execution_authority":False,
            "promotion_authority":False,
            "production_authority":False,
            "successor_authority":False,
            "retraining_required":False,
        },
    }

def write_source(root):
    files={
        "model/weights.bin":b"weights-v1",
        "state/resume.json":b'{"step":7}\n',
        "replay/input.json":b'{"x":[1,2,3]}\n',
        "replay/output.json":b'{"y":[4,5]}\n',
    }
    for rel, raw in files.items():
        p=root/rel;p.parent.mkdir(parents=True,exist_ok=True);p.write_bytes(raw)
    return files

class PortableCheckpointTests(unittest.TestCase):
    def test_valid_spec(self):
        pc.validate_spec(spec())

    def test_export_verify_restore_byte_exact_without_retraining(self):
        with tempfile.TemporaryDirectory() as td:
            base=pathlib.Path(td);src=base/"src";src.mkdir();expected=write_source(src)
            sp=base/"spec.json";sp.write_text(json.dumps(spec()),encoding="utf-8")
            cp=base/"checkpoint";pc.export_checkpoint(sp,src,cp)
            m=pc.verify_checkpoint(cp)
            restored=base/"restored";receipt=pc.restore_checkpoint(cp,restored)
            self.assertFalse(receipt["retraining_performed"])
            self.assertTrue(receipt["requalification_required"])
            self.assertFalse(receipt["promotion_authority"])
            for rel,raw in expected.items(): self.assertEqual((restored/rel).read_bytes(),raw)
            self.assertTrue(m["replay"]["runtime_replay_required_before_promotion"])

    def test_manifest_tamper_rejected(self):
        with tempfile.TemporaryDirectory() as td:
            base=pathlib.Path(td);src=base/"src";src.mkdir();write_source(src)
            sp=base/"spec.json";sp.write_text(json.dumps(spec()),encoding="utf-8")
            cp=base/"checkpoint";pc.export_checkpoint(sp,src,cp)
            p=cp/"manifest.json";m=json.loads(p.read_text());m["mechanism_id"]="tampered";p.write_text(json.dumps(m))
            with self.assertRaises(pc.PortabilityError): pc.verify_checkpoint(cp)

    def test_file_tamper_rejected(self):
        with tempfile.TemporaryDirectory() as td:
            base=pathlib.Path(td);src=base/"src";src.mkdir();write_source(src)
            sp=base/"spec.json";sp.write_text(json.dumps(spec()),encoding="utf-8")
            cp=base/"checkpoint";pc.export_checkpoint(sp,src,cp)
            (cp/"files/model/weights.bin").write_bytes(b"tamper")
            with self.assertRaises(pc.PortabilityError): pc.verify_checkpoint(cp)

    def test_traversal_rejected(self):
        s=spec();s["files"][0]["path"]="../weights.bin"
        with self.assertRaises(pc.PortabilityError): pc.validate_spec(s)

    def test_duplicate_path_rejected(self):
        s=spec();s["files"].append(dict(s["files"][0]))
        with self.assertRaises(pc.PortabilityError): pc.validate_spec(s)

    def test_missing_required_role_rejected(self):
        s=spec();s["files"]=[x for x in s["files"] if x["role"]!="resumable_state"]
        with self.assertRaises(pc.PortabilityError): pc.validate_spec(s)

    def test_bad_source_identity_rejected(self):
        s=spec();s["source_tree_sha"]="bad"
        with self.assertRaises(pc.PortabilityError): pc.validate_spec(s)

    def test_bad_qualification_identity_rejected(self):
        s=spec();s["qualification_contract_sha256"]="bad"
        with self.assertRaises(pc.PortabilityError): pc.validate_spec(s)

    def test_authority_widening_rejected(self):
        s=spec();s["authority"]["promotion_authority"]=True
        with self.assertRaises(pc.PortabilityError): pc.validate_spec(s)

    def test_resource_identity_required(self):
        s=spec();s["resource_envelope"]["baseline_id"]=""
        with self.assertRaises(pc.PortabilityError): pc.validate_spec(s)

    def test_replay_binding_required(self):
        s=spec();s["replay"]["vectors"][0]["output_path"]="model/weights.bin"
        with self.assertRaises(pc.PortabilityError): pc.validate_spec(s)

    def test_project_specific_contract(self):
        s=spec()
        if PROJECT=="Wingless": s["interface"]["fixed_state_dimension"]=0
        else: s["interface"]["retained_slots"]=8
        with self.assertRaises(pc.PortabilityError): pc.validate_spec(s)

if __name__=="__main__":
    unittest.main()
