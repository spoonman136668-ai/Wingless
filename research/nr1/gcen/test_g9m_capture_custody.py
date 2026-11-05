"""CNC-G9M offline proof boundary; no original model, network or trace upload."""
import hashlib
from pathlib import Path
import tempfile
import unittest
from g9m_capture_custody import SCHEMA, MODEL, REVISION, verify_capture


class G9MCustodyTests(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.root = Path(self.tmp.name)
        for n, content in (("cal.bin",b"fake-local-development-cal"),("eval.bin",b"fake-local-development-eval")):
            (self.root/n).write_bytes(content)
        self.data={
            "schema":SCHEMA,"model":MODEL,"revision":REVISION,"cohort":"development",
            "protected_heldout_used":False,"capture_kind":"native_model_forward_hook",
            "source_commit":"a"*40,"workflow_run_id":12345,
            "tensor_name":"model.layers.0.mlp.experts.0.down_proj.weight",
            "calibration_record_ids":["c1","c2"],"evaluation_record_ids":["e1","e2"],
            "trace_files":[
                {"cohort":c,"path":p,"sha256":hashlib.sha256((self.root/p).read_bytes()).hexdigest()}
                for c,p in (("calibration","cal.bin"),("evaluation","eval.bin"))
            ]
        }

    def test_checksum_and_disjoint_records_do_not_prove_original_model(self):
        result=verify_capture(self.data,self.root)
        self.assertTrue(result["local_sha256_verified"])
        self.assertEqual(result["original_forward_authenticity"],
                         "UNVERIFIED_REQUIRES_INDEPENDENT_ATTESTATION")
        self.assertFalse(result["source_model_activation_evaluation_admitted"])
        self.assertEqual(result["scientific_acceptance"],"NOT_GRANTED")
        self.assertEqual(result["gpu_resource_qualification"],"UNVERIFIED")

    def test_forbid_disallowed_claims(self):
        for key,value in (("revision","main"),("cohort","heldout"),
                          ("protected_heldout_used",True),
                          ("capture_kind","synthetic_test_fixture"),
                          ("source_commit","main"),("workflow_run_id",True),
                          ("tensor_name","wrong")):
            with self.subTest(key=key):
                m={**self.data,key:value}
                with self.assertRaises(ValueError):verify_capture(m,self.root)

    def test_forbid_tamper_overlap_symlink_and_traversal(self):
        m={**self.data,"evaluation_record_ids":["c1"]}
        with self.assertRaises(ValueError):verify_capture(m,self.root)
        m={**self.data,"trace_files":[dict(x) for x in self.data["trace_files"]]}
        m["trace_files"][0]["sha256"]="0"*64
        with self.assertRaises(ValueError):verify_capture(m,self.root)
        m["trace_files"][0]["sha256"]=self.data["trace_files"][0]["sha256"]
        m["trace_files"][0]["path"]="../cal.bin"
        with self.assertRaises(ValueError):verify_capture(m,self.root)
        (self.root/"cal.bin").write_bytes(b"tampered")
        with self.assertRaises(ValueError):verify_capture(self.data,self.root)

    def test_no_capture_files_or_heldout_claims_created(self):
        result=verify_capture(self.data,self.root)
        self.assertEqual(set(x.name for x in self.root.iterdir()),{"cal.bin","eval.bin"})
        self.assertFalse(result["published_raw_activations"])


if __name__=="__main__":
    unittest.main()
