"""G27 offline full-map negative controls, independent of torch or network."""
import ast
from pathlib import Path
import unittest

from g27_original_meta_mapping import (
    check_full_key_mapping, EXPECTED_INDEXED_TENSORS,
    EXPECTED_EXPERT_PROJECTIONS, SELECTED_TENSORS,
)
from out_of_core_preflight import SHARDS


def synthetic_reference_keys():
    # Exactly 48*(128*3 + 9) + 3 = 18,867 names, matching the observed
    # pinned original index count; these are dummy keys, never real weights.
    names = ["model.embed_tokens.weight", "model.norm.weight", "lm_head.weight"]
    for layer in range(48):
        prefix = f"model.layers.{layer}"
        for name in ("q_proj", "k_proj", "v_proj", "o_proj", "q_norm", "k_norm"):
            names.append(f"{prefix}.self_attn.{name}.weight")
        names.extend((
            f"{prefix}.input_layernorm.weight",
            f"{prefix}.post_attention_layernorm.weight",
            f"{prefix}.mlp.gate.weight",
        ))
        for expert in range(128):
            for projection in ("gate_proj", "up_proj", "down_proj"):
                names.append(f"{prefix}.mlp.experts.{expert}.{projection}.weight")
    assert len(names) == EXPECTED_INDEXED_TENSORS
    return names


class G27OriginalMetaMappingTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.names = synthetic_reference_keys()
        cls.mapping = {name: SHARDS[i % len(SHARDS)] for i, name in enumerate(cls.names)}

    def test_exact_layout_frozen(self):
        result = check_full_key_mapping(self.names, self.mapping)
        self.assertEqual(result["indexed_tensors"], 18867)
        self.assertEqual(result["expert_modules"], 6144)
        self.assertEqual(result["expert_projection_tensors"], EXPECTED_EXPERT_PROJECTIONS)
        self.assertEqual(result["missing_tensors"], 0)

    def test_missing_original_parameter_rejected(self):
        with self.assertRaisesRegex(ValueError, "meta parameter count mismatch"):
            check_full_key_mapping(self.names[:-1], self.mapping)

    def test_duplicate_original_parameter_rejected(self):
        names = list(self.names)
        names[-1] = names[-2]
        with self.assertRaisesRegex(ValueError, "duplicate registered"):
            check_full_key_mapping(names, self.mapping)

    def test_wrong_original_name_rejected_even_if_counts_match(self):
        names = list(self.names)
        names[-1] = "invented.weight"
        with self.assertRaisesRegex(ValueError, "checkpoint/module name mismatch"):
            check_full_key_mapping(names, self.mapping)

    def test_missing_shard_reference_rejected(self):
        mapping = dict(self.mapping)
        for key in mapping:
            if mapping[key] == SHARDS[-1]:
                mapping[key] = SHARDS[0]
        with self.assertRaisesRegex(ValueError, "incomplete 16-shard"):
            check_full_key_mapping(self.names, mapping)

    def test_corrupt_source_index_count_rejected(self):
        mapping = dict(self.mapping)
        mapping.pop(SELECTED_TENSORS[0])
        with self.assertRaisesRegex(ValueError, "unexpected original indexed tensor count"):
            check_full_key_mapping(self.names, mapping)

    def test_missing_expert_projection_rejected(self):
        bad = "model.layers.47.mlp.experts.127.down_proj.weight"
        fake = "model.layers.47.nonexistent.weight"
        names = [fake if name == bad else name for name in self.names]
        mapping = {fake if name == bad else name: shard for name, shard in self.mapping.items()}
        with self.assertRaisesRegex(ValueError, "required expert module tensor missing"):
            check_full_key_mapping(names, mapping)

    def test_no_original_weight_loading_in_probe(self):
        code = Path(__file__).with_name("g27_original_meta_mapping.py").read_text(encoding="utf-8")
        ast.parse(code)
        self.assertIn("init_empty_weights(include_buffers=True)", code)
        self.assertIn("tensor.device.type != \"meta\"", code)
        for forbidden in ("load_checkpoint_and_dispatch(", "disk_offload(", "from_pretrained(",
                          "hf_hub_download(", "snapshot_download(", ".generate("):
            self.assertNotIn(forbidden, code)
        self.assertIn('"original_model_forward_executed": False', code)
        self.assertIn('"scientific_acceptance": "NOT_GRANTED"', code)


if __name__ == "__main__":
    unittest.main()
