import copy
import importlib.util
import json
import pathlib
import tempfile
import unittest

ROOT = pathlib.Path(__file__).resolve().parents[1]
SPEC = importlib.util.spec_from_file_location("free_fanout_sidecar", ROOT / "scripts" / "free_fanout_sidecar.py")
ff = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(ff)

PROGRAM = "Wingless"
PACKAGE_SHA = "24d88a08468a0df8124b6fe4dabb2eaa5d74fb7c"

def base_manifest():
    if PROGRAM == "Wingless":
        interface = {
            "signature": "state_former(raw_input, history) -> fixed_size_state",
            "fixed_state_dimension": 8,
            "fixed_readout_capacity": 1,
            "training_data_class": "historical-replay",
            "sealed_outcome_exposure": False,
            "post_result_tuning": False,
            "capacity_growth": False,
        }
    else:
        interface = {
            "signature": "cognition_consumer(retained_state, local_state) -> decision_state",
            "total_slots": 16,
            "active_slots": 7,
            "retained_slots": 9,
            "training_data_class": "historical-replay",
            "sealed_outcome_exposure": False,
            "post_result_tuning": False,
            "hidden_persistent_memory_growth": False,
            "addressing_change": False,
            "capacity_growth": False,
        }
    return {
        "schema": "research.free-fanout.v1",
        "program": PROGRAM,
        "package_sha": PACKAGE_SHA,
        "authority": "external-evidence-only",
        "data_class": "historical-replay",
        "sealed_inputs": False,
        "accepted_state_mutation": False,
        "successor_dispatch": False,
        "queue_write": False,
        "production_access": False,
        "promotion_authority": False,
        "candidate_timeout_seconds": 60,
        "entrypoint": "research/fanout/smoke_candidate.py",
        "candidate_interface": interface,
        "resource_envelope": {
            "comparison": "equal-or-lower-than-baseline",
            "baseline_id": "synthetic-smoke-baseline-v1",
            "max_parameters": 0,
            "max_context_bytes": 8192,
            "max_model_calls": 0,
            "max_peak_rss_kib": 1048576,
            "max_elapsed_seconds": 30,
        },
        "candidates": [
            {"id": "smoke-a", "config": {"variant": 0}},
            {"id": "smoke-b", "config": {"variant": 1}},
        ],
    }

class ManifestValidationTests(unittest.TestCase):
    def assertRejected(self, manifest, code, package_sha=PACKAGE_SHA):
        with self.assertRaises(ff.FanoutError) as ctx:
            ff.validate_manifest(manifest, PROGRAM, package_sha)
        self.assertIn(code, str(ctx.exception))

    def test_valid_manifest_accepted(self):
        ff.validate_manifest(base_manifest(), PROGRAM, PACKAGE_SHA)

    def test_invalid_sha_rejected(self):
        self.assertRejected(base_manifest(), "FREE_FANOUT_PACKAGE_SHA_FORMAT", package_sha="not-a-sha")

    def test_duplicate_candidate_ids_rejected(self):
        m = base_manifest()
        m["candidates"][1]["id"] = m["candidates"][0]["id"]
        self.assertRejected(m, "FREE_FANOUT_DUPLICATE_CANDIDATE_ID")

    def test_sealed_inputs_rejected(self):
        m = base_manifest()
        m["sealed_inputs"] = True
        self.assertRejected(m, "FREE_FANOUT_FORBIDDEN_SEALED_INPUTS")

    def test_queue_write_rejected(self):
        m = base_manifest()
        m["queue_write"] = True
        self.assertRejected(m, "FREE_FANOUT_FORBIDDEN_QUEUE_WRITE")

    def test_accepted_state_mutation_rejected(self):
        m = base_manifest()
        m["accepted_state_mutation"] = True
        self.assertRejected(m, "FREE_FANOUT_FORBIDDEN_ACCEPTED_STATE_MUTATION")

    def test_successor_authority_rejected(self):
        m = base_manifest()
        m["successor_dispatch"] = True
        self.assertRejected(m, "FREE_FANOUT_FORBIDDEN_SUCCESSOR_DISPATCH")

    def test_path_traversal_rejected(self):
        m = base_manifest()
        m["entrypoint"] = "research/fanout/../sealed.py"
        self.assertRejected(m, "FREE_FANOUT_PATH_INVALID")

    def test_more_than_32_candidates_rejected(self):
        m = base_manifest()
        m["candidates"] = [{"id": f"candidate-{i}"} for i in range(33)]
        self.assertRejected(m, "FREE_FANOUT_CANDIDATE_COUNT")

    def test_missing_entrypoint_rejected(self):
        m = base_manifest()
        del m["entrypoint"]
        self.assertRejected(m, "FREE_FANOUT_PATH_INVALID")

    def test_unsupported_data_class_rejected(self):
        m = base_manifest()
        m["data_class"] = "sealed-unseen"
        self.assertRejected(m, "FREE_FANOUT_DATA_CLASS")

    def test_excessive_timeout_rejected(self):
        m = base_manifest()
        m["candidate_timeout_seconds"] = 1801
        self.assertRejected(m, "FREE_FANOUT_TIMEOUT")

    def test_wrong_project_rejected(self):
        m = base_manifest()
        m["program"] = "WrongProject"
        self.assertRejected(m, "FREE_FANOUT_MANIFEST_PROGRAM")

    def test_wrong_schema_rejected(self):
        m = base_manifest()
        m["schema"] = "wrong.schema"
        self.assertRejected(m, "FREE_FANOUT_MANIFEST_SCHEMA")

    def test_manifest_path_escape_rejected(self):
        with self.assertRaises(ff.FanoutError):
            ff.validate_manifest_path("research/fanout/../../manifest.json")

class RuntimeValidationTests(unittest.TestCase):
    def make_package(self, source):
        td = tempfile.TemporaryDirectory()
        root = pathlib.Path(td.name)
        target = root / "research" / "fanout" / "smoke_candidate.py"
        target.parent.mkdir(parents=True)
        target.write_text(source, encoding="utf-8")
        return td, root

    def test_missing_result_rejected(self):
        m = base_manifest()
        source = "import sys\nsys.exit(0)\n"
        td, root = self.make_package(source)
        self.addCleanup(td.cleanup)
        with tempfile.TemporaryDirectory() as out:
            with self.assertRaises(ff.FanoutError) as ctx:
                ff.run_candidate(m, root, "smoke-a", out)
        self.assertIn("FREE_FANOUT_RESULT_MISSING", str(ctx.exception))

    def test_mismatched_candidate_id_rejected(self):
        m = base_manifest()
        source = """import json, pathlib, sys
out = pathlib.Path(sys.argv[sys.argv.index("--output") + 1])
out.write_text(json.dumps({"candidate_id": "wrong"}) + "\\n", encoding="utf-8")
"""
        td, root = self.make_package(source)
        self.addCleanup(td.cleanup)
        with tempfile.TemporaryDirectory() as out:
            with self.assertRaises(ff.FanoutError) as ctx:
                ff.run_candidate(m, root, "smoke-a", out)
        self.assertIn("FREE_FANOUT_RESULT_ID_MISMATCH", str(ctx.exception))

    def test_aggregate_rejects_missing_candidate_envelope(self):
        m = base_manifest()
        with tempfile.TemporaryDirectory() as td:
            with self.assertRaises(ff.FanoutError) as ctx:
                ff.aggregate(m, td)
        self.assertIn("FREE_FANOUT_AGGREGATE_CANDIDATE_SET_MISMATCH", str(ctx.exception))

class WorkflowStaticTests(unittest.TestCase):
    def test_workflow_is_manual_hosted_and_read_only(self):
        text = (ROOT / ".github" / "workflows" / "free-fanout-sidecar.yml").read_text(encoding="utf-8")
        self.assertIn("workflow_dispatch:", text)
        self.assertNotIn("schedule:", text)
        self.assertNotIn("\n  push:", text)
        self.assertIn("manifest_ref:", text)
        self.assertIn("path: manifest-package", text)
        self.assertIn("ref: ${{ inputs.manifest_ref }}", text)
        self.assertIn("controller.manifest.json", text)
        self.assertIn("run-name: CKB fanout Wingless", text)
        self.assertNotIn("self-hosted", text)
        self.assertIn("contents: read", text)
        self.assertIn("persist-credentials: false", text)
        self.assertIn("runs-on: ubuntu-latest", text)
        self.assertIn("max-parallel: 8", text)

if __name__ == "__main__":
    unittest.main()
