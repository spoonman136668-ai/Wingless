"""G9N original disk headroom: resource negative controls, no model weights."""
from pathlib import Path
import tempfile
from types import SimpleNamespace
from unittest.mock import patch
import unittest
from g9n_resource_admission import full_checkpoint_disk_preflight, GIB


def fake_plan():
    return {"model":"Qwen/Qwen3-30B-A3B","revision":"a"*40,
            "shapes":{"original.weight":(2,2)}, "disk_bytes":10*GIB,
            "indexed_tensors":18867}


class AdmissionDiskTests(unittest.TestCase):
    def test_sufficient_scratch_only_proves_disk_floor(self):
        with tempfile.TemporaryDirectory() as d:
            with patch("g9n_resource_admission.freeze_original_assignment",return_value=fake_plan()), \
                 patch("g9n_resource_admission.preflight",return_value=(Path(d)/"source",{"shard":4*GIB})):
                rec=full_checkpoint_disk_preflight(manifest={},config={},index={},
                    device_map={},cpu_budget_bytes=4*GIB,disk_budget_bytes=16*GIB,
                    source_dir=Path(d)/"source",scratch_dir=d,
                    disk_usage=lambda _:SimpleNamespace(free=23*GIB))
                self.assertEqual(rec["minimum_scratch_free_bytes"],22*GIB)
                self.assertFalse(rec["full_original_admission_authorized"])
                self.assertFalse(rec["source_shard_hashes_verified"])
                self.assertFalse(rec["cpu_available_ram_measured"])

    def test_exact_shortfall_rejected(self):
        with tempfile.TemporaryDirectory() as d:
            with patch("g9n_resource_admission.freeze_original_assignment",return_value=fake_plan()), \
                 patch("g9n_resource_admission.preflight",return_value=(Path(d)/"source",{"shard":4*GIB})):
                with self.assertRaisesRegex(RuntimeError,"insufficient"):
                    full_checkpoint_disk_preflight(manifest={},config={},index={},
                        device_map={},cpu_budget_bytes=4*GIB,disk_budget_bytes=16*GIB,
                        source_dir=Path(d)/"source",scratch_dir=d,
                        disk_usage=lambda _:SimpleNamespace(free=22*GIB-1))

    def test_symlinked_scratch_rejected(self):
        with tempfile.TemporaryDirectory() as d:
            root=Path(d)
            actual=root/"scratch"
            actual.mkdir()
            linked=root/"link"
            try:
                linked.symlink_to(actual,target_is_directory=True)
            except (NotImplementedError,OSError):
                self.skipTest("symlink creation not supported")
            with patch("g9n_resource_admission.freeze_original_assignment",return_value=fake_plan()), \
                 patch("g9n_resource_admission.preflight",return_value=(root/"source",{"shard":4*GIB})):
                with self.assertRaisesRegex(ValueError,"symlinked"):
                    full_checkpoint_disk_preflight(manifest={},config={},index={},
                        device_map={},cpu_budget_bytes=4*GIB,disk_budget_bytes=16*GIB,
                        source_dir=root/"source",scratch_dir=linked)

    def test_source_and_scratch_same_rejected(self):
        with tempfile.TemporaryDirectory() as d:
            with patch("g9n_resource_admission.freeze_original_assignment",return_value=fake_plan()), \
                 patch("g9n_resource_admission.preflight",return_value=(Path(d),{"shard":4*GIB})):
                with self.assertRaisesRegex(ValueError,"cannot be"):
                    full_checkpoint_disk_preflight(manifest={},config={},index={},
                        device_map={},cpu_budget_bytes=4*GIB,disk_budget_bytes=16*GIB,
                        source_dir=d,scratch_dir=d)

    def test_original_inventory_must_be_present(self):
        with tempfile.TemporaryDirectory() as d:
            with patch("g9n_resource_admission.freeze_original_assignment",return_value=fake_plan()), \
                 patch("g9n_resource_admission.preflight",side_effect=ValueError("source shard missing")):
                with self.assertRaisesRegex(ValueError,"missing"):
                    full_checkpoint_disk_preflight(manifest={},config={},index={},
                        device_map={},cpu_budget_bytes=4*GIB,disk_budget_bytes=16*GIB,
                        source_dir=d,scratch_dir=d)
if __name__=="__main__":
    unittest.main()
