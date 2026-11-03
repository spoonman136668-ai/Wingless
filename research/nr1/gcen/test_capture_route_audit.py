"""G22 capture-route admission is evidence-only and never authorizes inference."""
import json
from pathlib import Path
import unittest

from capture_route_audit import audit_route
from activation_replay import MODEL, REVISION


class RouteAuditTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.manifest = json.loads((Path(__file__).parent / "upstream_full_manifest_unverified_weights.json").read_text())

    def proposal(self):
        return {
            "model": MODEL,
            "revision": REVISION,
            "checkpoint_format": "original_safetensors",
            "runtime": "native_full_model_forward",
            "capture_kind": "native_expert_input_forward_hook",
            "tensor_name": "model.layers.0.mlp.experts.0.down_proj.weight",
            "source_commit": "a"*40,
            "checkpoint_files": self.manifest["files"].copy(),
            "cohort": "development",
            "calibration_record_ids": ["fit-a"],
            "evaluation_record_ids": ["eval-b"],
            "authority_reference": "proposal-only:operator-approval-not-verified",
            "qualified_runner_reference": "proposal-only:runner-qualification-not-verified",
        }

    def test_valid_declaration_still_does_not_admit_source(self):
        r = audit_route(self.proposal(), self.manifest)
        self.assertEqual(r["status"], "DECLARATION_COMPATIBLE_SOURCE_UNVERIFIED")
        self.assertEqual(r["source_authenticity"], "UNVERIFIED")
        self.assertFalse(r["execution_authorized_by_this_audit"])
        self.assertEqual(r["scientific_acceptance"], "NOT_GRANTED")

    def test_nr1a_gguf_router_trace_not_compatible(self):
        p = self.proposal()
        p.update(model="Qwen3-Coder-30B-A3B-Instruct",checkpoint_format="GGUF_Q3_K_M",runtime="llama_cpp_router_trace",capture_kind="router_telemetry")
        r = audit_route(p,self.manifest)
        self.assertEqual(r["status"],"BLOCKED")
        self.assertIn("WRONG_ORIGINAL_MODEL",r["blockers"])
        self.assertIn("INCOMPATIBLE_CHECKPOINT_FORMAT",r["blockers"])
        self.assertIn("ROUTER_TELEMETRY_IS_NOT_EXPERT_ACTIVATIONS",r["blockers"])

    def test_wrong_revision_is_blocked(self):
        p=self.proposal();p["revision"]="0"*40
        self.assertIn("WRONG_ORIGINAL_MODEL",audit_route(p,self.manifest)["blockers"])

    def test_missing_shard_is_blocked(self):
        p=self.proposal();p["checkpoint_files"].pop("model-00001-of-00016.safetensors")
        self.assertIn("ORIGINAL_CHECKPOINT_FILE_IDENTITY_MISSING",audit_route(p,self.manifest)["blockers"])

    def test_mutated_shard_digest_is_blocked(self):
        p=self.proposal();p["checkpoint_files"]["model-00001-of-00016.safetensors"]="0"*64
        self.assertIn("ORIGINAL_CHECKPOINT_FILE_IDENTITY_MISSING",audit_route(p,self.manifest)["blockers"])

    def test_protected_cohort_blocked(self):
        p=self.proposal();p["cohort"]="protected"
        self.assertIn("PROTECTED_OR_UNFROZEN_COHORT_FORBIDDEN",audit_route(p,self.manifest)["blockers"])

    def test_fit_eval_overlap_blocked(self):
        p=self.proposal();p["evaluation_record_ids"]=["fit-a"]
        self.assertIn("CALIBRATION_EVALUATION_SEPARATION_UNVERIFIED",audit_route(p,self.manifest)["blockers"])

    def test_missing_runner_and_authority_blocked(self):
        p=self.proposal();p["authority_reference"]="";p.pop("qualified_runner_reference")
        r=audit_route(p,self.manifest)
        self.assertIn("EXECUTION_AUTHORITY_NOT_IDENTIFIED",r["blockers"])
        self.assertIn("QUALIFIED_RUNTIME_NOT_IDENTIFIED",r["blockers"])

    def test_missing_runtime_commit_blocked(self):
        p=self.proposal();p["source_commit"]=""
        self.assertIn("UNPINNED_RUNTIME_SOURCE",audit_route(p,self.manifest)["blockers"])

    def test_non_expert_tensor_blocked(self):
        p=self.proposal();p["tensor_name"]="model.layers.0.mlp.router.weight"
        self.assertIn("INVALID_EXPERT_TENSOR",audit_route(p,self.manifest)["blockers"])

    def test_untrusted_manifest_does_not_pass(self):
        m=dict(self.manifest);m["files"]={"x":"0"*64}
        self.assertIn("INCOMPLETE_PINNED_CHECKPOINT_MANIFEST",audit_route(self.proposal(),m)["blockers"])

    def test_missing_descriptor_fails_closed(self):
        self.assertEqual(audit_route(None,self.manifest)["status"],"BLOCKED")


if __name__=="__main__":
    unittest.main()
