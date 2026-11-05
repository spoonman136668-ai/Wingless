"""G9N host sinks: offline negative controls, no original weights or Torch needed."""
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch
from g9n_assignment_sinks import HostAssignmentSinks, source_bound_sinks


class FakeDevice:
    type = "cpu"


class FakeTensor:
    dtype = "torch.bfloat16"
    device = FakeDevice()

    def __init__(self, shape=(2, 2)):
        self.shape = shape


def plan():
    return {"schema": "wingless.gcen.g9n.repository-source-placement.v1",
            "shapes": {"cpu.weight": (2, 2), "disk.weight": (2, 2)},
            "target_devices": {"cpu.weight": "cpu", "disk.weight": "disk"},
            "indexed_tensors": 2, "original_tensor_bytes": 16,
            "cpu_bytes": 8, "disk_bytes": 8}


def writer(path, key, tensor):
    Path(path).write_bytes(b"tensor-test-only")


class HostSinkTests(unittest.TestCase):
    def test_host_routing_and_complete_receipt(self):
        with tempfile.TemporaryDirectory() as d:
            with HostAssignmentSinks(plan(), d, writer=writer) as session:
                session.assign("cpu.weight", FakeTensor())
                session.assign("disk.weight", FakeTensor())
                self.assertEqual(session.receipt()["disk_tensors"], 1)
                self.assertEqual(session.receipt()["cpu_tensors"], 1)
                self.assertTrue(next(iter(session.disk_paths.values())).is_file())
                owned = session.directory
            self.assertFalse(owned.exists())

    def test_malformed_plan_fails_before_disk_mutation(self):
        p = plan()
        p["disk_bytes"] = 9
        with tempfile.TemporaryDirectory() as d:
            with self.assertRaisesRegex(ValueError, "accounting"):
                HostAssignmentSinks(p, d, writer=writer)
            self.assertEqual(list(Path(d).iterdir()), [])

    def test_duplicate_and_unknown_rejected(self):
        with tempfile.TemporaryDirectory() as d:
            with HostAssignmentSinks(plan(), d, writer=writer) as s:
                s.assign("cpu.weight", FakeTensor())
                with self.assertRaisesRegex(ValueError, "duplicate"):
                    s.assign("cpu.weight", FakeTensor())
                with self.assertRaisesRegex(ValueError, "unknown"):
                    s.assign("unknown", FakeTensor())

    def test_wrong_shape_and_dtype_rejected(self):
        with tempfile.TemporaryDirectory() as d:
            with HostAssignmentSinks(plan(), d, writer=writer) as s:
                with self.assertRaisesRegex(ValueError, "shape"):
                    s.assign("disk.weight", FakeTensor((2, 3)))
                t = FakeTensor()
                t.dtype = "torch.float32"
                with self.assertRaisesRegex(ValueError, "BF16"):
                    s.assign("disk.weight", t)
                self.assertEqual(list(s.directory.iterdir()), [])

    def test_failed_disk_writer_cleans_temporary_files(self):
        def broken(path, *_):
            Path(path).write_bytes(b"partial")
            raise OSError("injected sink failure")
        with tempfile.TemporaryDirectory() as d:
            with HostAssignmentSinks(plan(), d, writer=broken) as s:
                with self.assertRaisesRegex(OSError, "injected"):
                    s.assign("disk.weight", FakeTensor())
                self.assertEqual(list(s.directory.iterdir()), [])
                s.discard_partial()
                self.assertFalse(s.directory.exists())

    def test_unverified_receipt_must_not_claim_disk_content_validation(self):
        with tempfile.TemporaryDirectory() as d:
            with HostAssignmentSinks(plan(), d, writer=writer) as session:
                session.assign("cpu.weight", FakeTensor())
                session.assign("disk.weight", FakeTensor())
                self.assertFalse(session.receipt()["disk_content_verified"])

    def test_incomplete_assignment_refuses_receipt(self):
        with tempfile.TemporaryDirectory() as d:
            with HostAssignmentSinks(plan(), d, writer=writer) as s:
                s.assign("cpu.weight", FakeTensor())
                with self.assertRaisesRegex(RuntimeError, "incomplete"):
                    s.receipt()

    def test_closed_session_refuses_assignments(self):
        with tempfile.TemporaryDirectory() as d:
            s = HostAssignmentSinks(plan(), d, writer=writer)
            s.close()
            s.close()
            with self.assertRaisesRegex(RuntimeError, "closed"):
                s.assign("disk.weight", FakeTensor())

    def test_source_bound_factory_calls_original_contract(self):
        with tempfile.TemporaryDirectory() as d:
            with patch("g9n_assignment_sinks.freeze_original_assignment", return_value=plan()) as frozen:
                with source_bound_sinks(manifest={"identity": "mock"}, config={}, index={},
                        device_map={"": "disk"}, cpu_budget_bytes=16,
                        disk_budget_bytes=16, scratch_dir=d, writer=writer):
                    pass
                frozen.assert_called_once()

    def test_symlinked_scratch_rejected(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root/"actual").mkdir()
            try:
                (root/"linked").symlink_to(root/"actual", target_is_directory=True)
            except (NotImplementedError, OSError):
                self.skipTest("symlink creation unavailable")
            with self.assertRaisesRegex(ValueError, "symlink"):
                HostAssignmentSinks(plan(), root/"linked", writer=writer)


if __name__ == "__main__":
    unittest.main()
