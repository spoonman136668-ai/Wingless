"""Offline no-download source controls for the optional G25 tiny runtime smoke."""
import ast
from pathlib import Path
import unittest

SMOKE = Path(__file__).with_name("g25_disk_offload_smoke.py")


class G25RuntimeContractTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.text = SMOKE.read_text(encoding="utf-8")
        cls.tree = ast.parse(cls.text)

    def test_expected_source_commit_ids_present(self):
        self.assertIn("47b0e478f324b54f177ea7998a0791870fdd0324", self.text)
        self.assertIn("623cc0ba588d192c77758939a6e340adf07edd4f", self.text)

    def test_expected_upstream_blob_ids_present(self):
        self.assertIn("25c738744faba82db3c451b6987b01db74fadcfd", self.text)
        self.assertIn("2bace0f5615f37c84c1c3e77fe116081ed691bac", self.text)

    def test_no_model_checkpoint_download_calls(self):
        forbidden = ("from_pretrained(", "hf_hub_download(", "snapshot_download(",
                     "urlopen(", "requests.get(", "load_checkpoint_and_dispatch(")
        for name in forbidden:
            with self.subTest(name=name):
                self.assertNotIn(name, self.text)

    def test_source_and_version_check_fail_closed(self):
        self.assertIn("pinned CPU runtime version mismatch", self.text)
        self.assertIn("installed source differs from pinned upstream Git blobs", self.text)

    def test_capture_requires_matching_original_baseline(self):
        self.assertIn("disk-offload changed expert input", self.text)
        self.assertIn("disk-offload changed tiny Qwen3-MoE expert output", self.text)

    def test_no_scientific_acceptance(self):
        self.assertIn('"authentic_expert_activations": False', self.text)
        self.assertIn('"production_inference_qualified": False', self.text)
        self.assertIn('"scientific_acceptance": "NOT_GRANTED"', self.text)

    def test_only_entrypoint_runs_real_torch_smoke(self):
        functions = [node.name for node in self.tree.body if isinstance(node, ast.FunctionDef)]
        self.assertIn("run", functions)
        self.assertIn("git_blob", functions)


if __name__ == "__main__":
    unittest.main()
