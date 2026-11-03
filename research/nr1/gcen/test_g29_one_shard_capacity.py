"""G29 shard-by-shard staging planning safety controls; no checkpoint data access."""
import ast
from pathlib import Path
import unittest
from types import SimpleNamespace

from g29_one_shard_capacity import shard_payload_footprints, sequential_staging_proposal
from g28_original_device_plan import GIB
from test_g27_original_meta_mapping import synthetic_reference_keys
from out_of_core_preflight import SHARDS


class FakeMeta:
    def __init__(self, count=1, device="meta"):
        self.count = count
        self.device = SimpleNamespace(type=device)
    def numel(self):
        return self.count


class G29OneShardPlanTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.names = synthetic_reference_keys()
        cls.parameters = [(name, FakeMeta()) for name in cls.names]
        cls.index = {name: SHARDS[i % 16] for i, name in enumerate(cls.names)}
        cls.original_size = len(cls.names) * 2 + 2048

    def test_reconciliation_and_conservative_header_upper(self):
        r = shard_payload_footprints(self.parameters, self.index, self.original_size)
        self.assertEqual(r["total_tensor_payload_bytes"], len(self.names) * 2)
        self.assertEqual(r["aggregate_file_overhead_bytes"], 2048)
        self.assertEqual(len(r["per_shard_bf16_payload_bytes"]), 16)
        self.assertEqual(sum(r["per_shard_bf16_payload_bytes"].values()), len(self.names) * 2)
        self.assertEqual(r["single_shard_physical_upper_bound_bytes"],
                         r["largest_shard_tensor_payload_bytes"] + 2048)
        self.assertEqual(r["actual_per_shard_file_bytes"], "UNVERIFIED_NOT_IN_INDEX")

    def test_missing_source_tensor_fails(self):
        bad = dict(self.index)
        bad.pop(self.names[2])
        with self.assertRaises(ValueError):
            shard_payload_footprints(self.parameters, bad, self.original_size)

    def test_missing_shard_coverage_fails(self):
        bad = dict(self.index)
        for key in bad:
            if bad[key] == SHARDS[-1]:
                bad[key] = SHARDS[0]
        with self.assertRaises(ValueError):
            shard_payload_footprints(self.parameters, bad, self.original_size)

    def test_inconsistent_original_file_total_fails(self):
        with self.assertRaisesRegex(ValueError, "reconciliation"):
            shard_payload_footprints(self.parameters, self.index, len(self.names) * 2 - 1)

    def test_oversized_unexplained_file_overhead_fails(self):
        with self.assertRaisesRegex(ValueError, "reconciliation"):
            shard_payload_footprints(self.parameters, self.index,
                                     len(self.names) * 2 + 33 * 1024**2)

    def test_real_cpu_parameter_fails(self):
        wrong = list(self.parameters)
        wrong[3] = (wrong[3][0], FakeMeta(device="cpu"))
        with self.assertRaisesRegex(ValueError, "must be meta"):
            shard_payload_footprints(wrong, self.index, self.original_size)

    def test_single_shard_requirement_is_conservative_addition(self):
        r = sequential_staging_proposal(50 * GIB, 5 * GIB, 80 * GIB)
        self.assertEqual(r["optimistic_sequential_scratch_plan_bytes"], 63 * GIB)
        self.assertTrue(r["fits_this_planning_estimate"])
        self.assertFalse(r["loader_implements_this_streaming_plan"])
        self.assertFalse(r["original_model_execution_authorized"])

    def test_insufficient_disk_does_not_pass(self):
        r = sequential_staging_proposal(50 * GIB, 5 * GIB, 60 * GIB)
        self.assertFalse(r["fits_this_planning_estimate"])
        self.assertEqual(r["optimistic_sequential_scratch_plan_bytes"], 63 * GIB)

    def test_unknown_or_negative_values_fail(self):
        for args in [(0, 5, 10), (5, 0, 10), (5, 5, -1),
                     (5, 5, True), (5, 5, 100, 1),
                     (1.0, 1, 100)]:
            with self.subTest(args=args), self.assertRaises(ValueError):
                sequential_staging_proposal(*args)

    def test_static_probe_never_executes_real_model(self):
        source = Path(__file__).with_name("g29_one_shard_capacity.py").read_text(encoding="utf-8")
        ast.parse(source)
        for prohibited in ("load_checkpoint_and_dispatch(", "from_pretrained(",
                           "hf_hub_download(", "snapshot_download(", "disk_offload(", ".generate("):
            self.assertNotIn(prohibited, source)
        self.assertIn('"full_original_checkpoint_bytes_accessed": False', source)
        self.assertIn('"scientific_acceptance": "NOT_GRANTED"', source)
        self.assertIn('"shard_staging_custody_tested": False', source)


if __name__ == "__main__":
    unittest.main()
