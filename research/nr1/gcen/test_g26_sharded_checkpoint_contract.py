"""G26 offline controls; actual Torch/Accelerate test runs only on pinned one-shot CI."""
import ast
from pathlib import Path
import unittest

CODE = Path(__file__).with_name("g26_sharded_checkpoint_smoke.py").read_text(encoding="utf-8")


class G26ContractTests(unittest.TestCase):
    def test_valid_syntax(self):
        ast.parse(CODE)

    def test_empty_weights_before_checkpoint_dispatch(self):
        self.assertLess(CODE.index("with init_empty_weights():"), CODE.index("load_checkpoint_and_dispatch("))

    def test_exact_cpu_disk_device_map(self):
        for literal in ['"gate_proj": "cpu"', '"up_proj": "cpu"', '"down_proj": "disk"']:
            self.assertIn(literal, CODE)

    def test_requires_two_shards_and_source_hashes(self):
        self.assertIn('model-00001-of-00002.safetensors', CODE)
        self.assertIn('model-00002-of-00002.safetensors', CODE)
        self.assertIn('hashlib.sha256(', CODE)

    def test_no_large_external_checkpoint_load(self):
        for forbidden in ("snapshot_download(", "from_pretrained(", "hf_hub_download(",
                          "urlopen(", "requests.get(", "Qwen/Qwen3-30B-A3B/resolve"):
            with self.subTest(forbidden=forbidden):
                self.assertNotIn(forbidden, CODE)

    def test_fails_closed_if_hook_parity_fails(self):
        self.assertIn("disk-dispatched down-projection hook missing", CODE)
        self.assertIn("G26 sharded load changed expert input activation", CODE)
        self.assertIn("G26 sharded load changed output", CODE)

    def test_scientific_claims_explicitly_prohibited(self):
        self.assertIn('"authentic_original_activation_trace": False', CODE)
        self.assertIn('"original_model_executed": False', CODE)
        self.assertIn('"scientific_acceptance": "NOT_GRANTED"', CODE)


if __name__ == "__main__":
    unittest.main()
