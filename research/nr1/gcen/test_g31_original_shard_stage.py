"""G31 offline synthetic negative controls; never request original weights."""
import hashlib
import io
import json
from pathlib import Path
import tempfile
import unittest

from g31_original_shard_stage import (
    SHARD, SHARDS, MAX_SHARD_BYTES, check_upstream, inspect_header,
    parse_unique_json, stage_verified
)


class FakeResponse:
    def __init__(self, data, *, status=200, length=True, url="https://cdn.example/shard",
                 content_range=None, encoding=None):
        self.data = io.BytesIO(data)
        self.status = status
        self.headers = {"Content-Length": str(len(data)) if length else None,
                        "Content-Range": content_range, "Content-Encoding": encoding}
        self.url = url

    def read(self, n):
        return self.data.read(n)

    def geturl(self):
        return self.url


class G31Tests(unittest.TestCase):
    @staticmethod
    def index():
        return {"weight_map": {"the.weight": SHARD, **{
            f"other.{i}": name for i, name in enumerate(SHARDS[1:])}}}

    @staticmethod
    def fake_safetensors(offsets=(0, 8), shape=(2, 2), dtype="BF16"):
        header = json.dumps({"the.weight": {
            "dtype": dtype, "shape": list(shape), "data_offsets": list(offsets)}}).encode()
        return len(header).to_bytes(8, "little") + header + b"\x01" * 8

    def test_valid_header_only(self):
        with tempfile.TemporaryDirectory() as d:
            p = Path(d) / SHARD
            p.write_bytes(self.fake_safetensors())
            info = inspect_header(p, self.index())
            self.assertEqual(info["tensor_count"], 1)
            self.assertEqual(info["indexed_payload_bytes"], 8)

    def test_missing_index_key_fails(self):
        with tempfile.TemporaryDirectory() as d:
            p = Path(d) / SHARD
            p.write_bytes(self.fake_safetensors())
            with self.assertRaisesRegex(ValueError, "mapping"):
                inspect_header(p, {"weight_map": {"wrong": SHARD, **{
                    f"other.{i}": name for i, name in enumerate(SHARDS[1:])}}})

    def test_bad_offset_and_shape_fail(self):
        for offsets, shape in (((0, 9), (2, 2)), ((1, 8), (2, 2)), ((0, 8), (0, 2))):
            with self.subTest(offsets=offsets, shape=shape):
                with tempfile.TemporaryDirectory() as d:
                    p = Path(d) / SHARD
                    p.write_bytes(self.fake_safetensors(offsets, shape))
                    with self.assertRaises(ValueError):
                        inspect_header(p, self.index())

    def test_duplicate_json_key_fails(self):
        with self.assertRaisesRegex(ValueError, "duplicate"):
            parse_unique_json(b'{"a":1,"a":2}')

    def test_valid_stream_digest_and_allocation(self):
        data = b"abcdef"
        with tempfile.TemporaryDirectory() as d:
            p = Path(d) / SHARD
            result = stage_verified(FakeResponse(data), p, hashlib.sha256(data).hexdigest(),
                                    min_bytes=1, max_bytes=12)
            self.assertEqual(result["bytes_verified"], 6)
            self.assertEqual(result["observed_peak_staged_file_count"], 1)
            self.assertEqual(p.read_bytes(), data)

    def test_corrupted_digest_removes_stage(self):
        with tempfile.TemporaryDirectory() as d:
            p = Path(d) / SHARD
            with self.assertRaisesRegex(ValueError, "SHA256"):
                stage_verified(FakeResponse(b"bad"), p, hashlib.sha256(b"good").hexdigest(),
                               min_bytes=1, max_bytes=12)
            self.assertFalse(p.exists())

    def test_truncated_http_length_fails(self):
        with tempfile.TemporaryDirectory() as d:
            p = Path(d) / SHARD
            response = FakeResponse(b"a", length=False)
            response.headers["Content-Length"] = "3"
            with self.assertRaisesRegex(ValueError, "truncated"):
                stage_verified(response, p, hashlib.sha256(b"a").hexdigest(),
                               min_bytes=1, max_bytes=10)
            self.assertFalse(p.exists())

    def test_partial_http_rejected_before_stage(self):
        with tempfile.TemporaryDirectory() as d:
            p = Path(d) / SHARD
            response = FakeResponse(b"a", status=206, content_range="bytes 0-0/3")
            with self.assertRaisesRegex(ValueError, "partial"):
                stage_verified(response, p, hashlib.sha256(b"a").hexdigest(),
                               min_bytes=1, max_bytes=10)
            self.assertFalse(p.exists())

    def test_encoding_and_insecure_redirect_rejected(self):
        for overrides in ({"encoding": "gzip"}, {"url": "http://untrusted/shard"}):
            with self.subTest(overrides=overrides):
                with tempfile.TemporaryDirectory() as d:
                    p = Path(d) / SHARD
                    with self.assertRaises(ValueError):
                        stage_verified(FakeResponse(b"a", **overrides), p,
                                       hashlib.sha256(b"a").hexdigest(),
                                       min_bytes=1, max_bytes=10)
                    self.assertFalse(p.exists())

    def test_exceeds_download_limit_and_cleans(self):
        with tempfile.TemporaryDirectory() as d:
            p = Path(d) / SHARD
            with self.assertRaisesRegex(ValueError, "cap"):
                stage_verified(FakeResponse(b"abcdef", length=False), p,
                               hashlib.sha256(b"abcdef").hexdigest(),
                               min_bytes=1, max_bytes=4)
            self.assertFalse(p.exists())

    def test_staging_overlap_rejected(self):
        with tempfile.TemporaryDirectory() as d:
            p = Path(d) / SHARD
            (Path(d) / "unrelated").write_bytes(b"x")
            with self.assertRaisesRegex(RuntimeError, "more than one"):
                stage_verified(FakeResponse(b"abc"), p, hashlib.sha256(b"abc").hexdigest(),
                               min_bytes=1, max_bytes=10)
            self.assertFalse(p.exists())

    def test_upstream_drift_rejected(self):
        from runtime_feasibility import MODEL, REVISION
        from out_of_core_preflight import SHARDS
        a = {"model": MODEL, "revision": REVISION,
             "files": {**{k: "a" * 64 for k in SHARDS},
                       "config.json": "b" * 64,
                       "model.safetensors.index.json": "c" * 64}}
        a["source"] = {"uri": "https://huggingface.co/Qwen/Qwen3-30B-A3B/tree/" + REVISION,
                       "revision": REVISION}
        b = json.loads(json.dumps(a))
        b["files"][SHARD] = "d" * 64
        with self.assertRaisesRegex(ValueError, "mismatch"):
            check_upstream(a, b)

    def test_never_conflates_staged_shard_with_model_inference(self):
        code = Path(__file__).with_name("g31_original_shard_stage.py").read_text(encoding="utf-8")
        for prohibited in ("load_checkpoint_and_dispatch(", "from_pretrained(", ".generate(",
                           "model(", "snapshot_download(", "hf_hub_download("):
            self.assertNotIn(prohibited, code)
        self.assertIn('"entire_model_loaded": False', code)
        self.assertIn('"scientific_acceptance": "NOT_GRANTED"', code)
        self.assertEqual(MAX_SHARD_BYTES, 6 * 1024**3)


if __name__ == "__main__":
    unittest.main()
