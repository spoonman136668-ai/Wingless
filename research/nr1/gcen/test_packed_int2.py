"""G9C tests: frozen synthetic INT2 numerical/packing/storage controls."""
import json
import unittest
import numpy as np
from packed_int2 import pack_int2, apply_packed_int2, stored_bytes
from packed_int4 import pack_int4, apply_packed_int4, stored_bytes as stored4


class PackedInt2Tests(unittest.TestCase):
    def test_codebook_roundtrip_and_block_size_controls(self):
        rng = np.random.default_rng(13)
        for block in (128, 256):
            code = np.tile(np.array([-3., -1., 1., 3.], dtype=np.float32), block//4)
            w = np.stack([code, np.zeros_like(code)]).astype(np.float32)
            x = rng.standard_normal((3, block)).astype(np.float32)
            packed, scales = pack_int2(w, block=block)
            self.assertEqual(packed.dtype, np.uint8)
            self.assertEqual(scales.dtype, np.float16)
            self.assertEqual(packed.shape, (2, block//4))
            self.assertEqual(stored_bytes(packed, scales, block, block=block),
                             packed.nbytes + scales.nbytes + 24)
            np.testing.assert_allclose(apply_packed_int2(x, packed, scales, block=block),
                                       x @ w.T, rtol=2e-5, atol=2e-5)

    def test_frozen_synthetic_quality_compared_with_int4(self):
        rng = np.random.default_rng(23)
        w = rng.standard_normal((128, 256)).astype(np.float32)
        x = rng.standard_normal((32, 256)).astype(np.float32)
        reference = x @ w.T
        energy = float(np.sum(reference.astype(np.float64)**2))
        p4, s4 = pack_int4(w)
        q4 = apply_packed_int4(x, p4, s4)
        err4 = float(np.sum((q4.astype(np.float64) - reference.astype(np.float64))**2) / energy)
        cases = []
        for block in (128, 256):
            packed, scales = pack_int2(w, block=block)
            y = apply_packed_int2(x, packed, scales, block=block)
            err = float(np.sum((y.astype(np.float64) - reference.astype(np.float64))**2) / energy)
            bytes2 = stored_bytes(packed, scales, 256, block=block)
            self.assertLess(bytes2, stored4(p4, s4, 256))
            self.assertTrue(np.isfinite(err))
            self.assertGreater(err, err4)
            self.assertLess(err, 2.)
            cases.append({"block": block, "error": err, "packed_bytes": bytes2})
        print("CNC_G9C_INT2_SYNTHETIC " + json.dumps({
            "status": "SYNTHETIC_DEVELOPMENT_ONLY_NOT_INFERENCE",
            "source": "SEEDED_FLOAT32_SYNTHETIC",
            "int4_error": err4,
            "int4_bytes": stored4(p4, s4, 256),
            "int2": cases, "full_model_quality": "UNVERIFIED",
            "gpu_vram": "UNMEASURED",
        }, sort_keys=True), flush=True)

    def test_fail_closed_scale_overflow_and_shape(self):
        for block in (0, 64, 129, True):
            with self.subTest(block=block), self.assertRaises(ValueError):
                pack_int2(np.ones((2, 256), dtype=np.float32), block=block)
        with self.assertRaises(ValueError):
            pack_int2(np.ones((2, 255), dtype=np.float32))
        with self.assertRaises(ValueError):
            pack_int2(np.full((2, 256), np.nan, dtype=np.float32))
        with self.assertRaises(ValueError):
            pack_int2(np.full((2, 256), 1e35, dtype=np.float32))
        packed, scales = pack_int2(np.ones((2, 256), dtype=np.float32))
        with self.assertRaises(ValueError):
            apply_packed_int2(np.ones((3, 255)), packed, scales)
        with self.assertRaises(ValueError):
            apply_packed_int2(np.ones((3, 256)), packed, -np.ones_like(scales))
        with self.assertRaises(ValueError):
            stored_bytes(packed, scales, 256, block=256)


if __name__ == "__main__":
    unittest.main()
