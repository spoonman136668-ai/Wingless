"""CNC-01Q-G7 packed INT4 tests, synthetic development only."""
import unittest

import numpy as np

from direct_expert import factorize, apply_direct
from packed_int4 import pack_int4, apply_packed_int4, stored_bytes


class PackedInt4Tests(unittest.TestCase):
    def test_exact_representable_codes_and_zero_group(self):
        vals = np.tile(np.arange(-7, 8, dtype=np.float32), 5)[:64]
        vals[-1] = 7
        w = np.stack([vals, np.zeros_like(vals)]).astype(np.float32)
        x = np.random.default_rng(17).standard_normal((3, 64)).astype(np.float32)
        packed, scales = pack_int4(w)
        np.testing.assert_array_equal(scales, np.ones((2, 1), dtype=np.float32))
        got = apply_packed_int4(x, packed, scales)
        np.testing.assert_allclose(got, x @ w.T, rtol=1e-5, atol=1e-5)

    def test_quantized_operator_accuracy_vs_svd_negative_control(self):
        rng = np.random.default_rng(23)
        w = rng.standard_normal((128, 256)).astype(np.float32)
        x = rng.standard_normal((32, 256)).astype(np.float32)
        packed, scales = pack_int4(w)
        actual = apply_packed_int4(x, packed, scales)
        reference = x @ w.T
        a, b = factorize(w, 14)
        lowrank = apply_direct(x, a, b)
        energy = float(np.sum(reference.astype(np.float64)**2))
        err_int4 = float(np.sum((actual.astype(np.float64) - reference.astype(np.float64))**2)/energy)
        err_svd = float(np.sum((lowrank.astype(np.float64) - reference.astype(np.float64))**2)/energy)
        # Development-only timed comparison; no speed threshold or GPU claim.
        import json
        from time import perf_counter
        def timed(operation):
            for _ in range(2):
                operation()
            samples = []
            for _ in range(5):
                start = perf_counter()
                operation()
                samples.append(perf_counter() - start)
            return float(np.median(samples))
        print("CNC01Q_PACKED_INT4_JSON " + json.dumps({
            "status": "SYNTHETIC_CPU_DEVELOPMENT_ONLY",
            "int4_relative_output_squared_error": err_int4,
            "svd_rank14_relative_output_squared_error": err_svd,
            "dense_fp16_bytes": int(w.size * 2),
            "packed_codes_bytes": int(packed.nbytes),
            "scales_float32_bytes": int(scales.nbytes),
            "total_packed_bytes_with_metadata": stored_bytes(packed, scales, 256),
            "dense_median_seconds": timed(lambda: x @ w.T),
            "packed_int4_median_seconds": timed(lambda: apply_packed_int4(x, packed, scales)),
            "peak_gpu_memory": "UNMEASURED", "real_qwen": False,
        }, sort_keys=True), flush=True)
        self.assertLess(err_int4, err_svd)
        self.assertLess(err_int4, .08)  # frozen synthetic diagnostic only
        self.assertLess(stored_bytes(packed, scales, 256), (w.size * 2) / 3)
        self.assertEqual(packed.shape, (128, 128))
        self.assertEqual(scales.shape, (128, 4))

    def test_reject_invalid_or_tampered_codes(self):
        w = np.ones((2, 64), dtype=np.float32)
        for block in (0, 30, True, 65):
            with self.subTest(block=block), self.assertRaises(ValueError):
                pack_int4(w, block=block)
        with self.assertRaises(ValueError):
            pack_int4(np.full((2, 64), np.nan))
        packed, scales = pack_int4(w)
        invalid_scales = scales.copy()
        invalid_scales[0, 0] = 0
        with self.assertRaises(ValueError):
            apply_packed_int4(np.ones((1, 64)), packed, invalid_scales)
        with self.assertRaises(ValueError):
            apply_packed_int4(np.ones((1, 32)), packed, scales)
        tampered = packed.copy()
        tampered[0, 0] = 0x88  # out-of-range -8 nibble
        with self.assertRaises(ValueError):
            apply_packed_int4(np.ones((1, 64)), tampered, scales)


if __name__ == "__main__":
    unittest.main()
