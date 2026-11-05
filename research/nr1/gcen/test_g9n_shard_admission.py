"""Pure G9N shard-admission negative controls; no Torch required."""
import tempfile
from pathlib import Path
import unittest
from g9b_tensor_budget import MODEL, REVISION
from g9n_shard_admission import SHARDS, admit_local_shards, preflight


def fixture():
    return ({"model": MODEL, "revision": REVISION,
             "files": {x: "a"*64 for x in (*SHARDS, "config.json",
                                             "model.safetensors.index.json")}},
            {"weight_map": {f"tensor.{i}": name for i, name in enumerate(SHARDS)}},
            {f"tensor.{i}": (2, 2) for i in range(16)})


class G9NShardAdmissionTests(unittest.TestCase):
    def test_missing_last_shard_before_any_sink_mutation(self):
        manifest, index, shapes = fixture()
        with tempfile.TemporaryDirectory() as root:
            src = Path(root)/"source"
            scratch = Path(root)/"scratch"
            src.mkdir()
            scratch.mkdir()
            called = []
            with self.assertRaisesRegex(ValueError, "missing"):
                admit_local_shards(manifest=manifest, index=index,
                    expected_shapes=shapes, source_dir=src, scratch_dir=scratch,
                    assign_tensor=lambda *_: called.append("assign"),
                    discard_partial=lambda: called.append("discard"))
            self.assertEqual(called, [])

    def test_missing_manifest_metadata_rejected(self):
        manifest, index, shapes = fixture()
        del manifest["files"]["config.json"]
        with tempfile.TemporaryDirectory() as root:
            with self.assertRaises(ValueError):
                preflight(manifest, index, shapes, Path(root))

    def test_missing_discard_callback_rejected(self):
        manifest, index, shapes = fixture()
        with self.assertRaisesRegex(ValueError, "discard"):
            admit_local_shards(manifest=manifest, index=index,
                expected_shapes=shapes, source_dir="missing",
                scratch_dir="missing", assign_tensor=lambda *_: None,
                discard_partial=None)

    def test_unknown_tensor_rejected(self):
        manifest, index, shapes = fixture()
        shapes["extra"] = (2, 2)
        with tempfile.TemporaryDirectory() as root:
            with self.assertRaises(ValueError):
                preflight(manifest, index, shapes, Path(root))

    def test_symlinked_root_rejected(self):
        manifest, index, shapes = fixture()
        with tempfile.TemporaryDirectory() as root:
            p = Path(root)
            (p/"src").mkdir()
            (p/"link").symlink_to(p/"src", target_is_directory=True)
            with self.assertRaisesRegex(ValueError, "symlinked"):
                preflight(manifest, index, shapes, p/"link")

    def test_wrong_source_identity_rejected(self):
        manifest, index, shapes = fixture()
        manifest["revision"] = "0"*40
        with tempfile.TemporaryDirectory() as root:
            with self.assertRaises(ValueError):
                preflight(manifest, index, shapes, Path(root))


if __name__ == "__main__":
    unittest.main()
