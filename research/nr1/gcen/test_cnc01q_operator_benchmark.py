"""CNC-01Q focused benchmark contract, synthetic only."""
import unittest

from cnc01q_operator_benchmark import measure


class CNC01QBenchmarkTests(unittest.TestCase):
    def test_bounded_operator_comparison(self):
        receipt = measure(repeats=3, warmups=1)
        self.assertEqual(receipt["status"], "SYNTHETIC_DEVELOPMENT_DIAGNOSTIC_ONLY")
        self.assertEqual(receipt["runtime_peak_vram_bytes"], "UNMEASURED")
        self.assertLess(receipt["factor_fp16_storage_bytes"], receipt["dense_fp16_storage_bytes"] / 4)
        self.assertEqual(set(receipt["results"]), {"dense_original", "dense_reconstructed", "direct_factors"})
        self.assertLess(receipt["results"]["direct_factors"]["relative_output_squared_error"], 1.0)
        self.assertAlmostEqual(
            receipt["results"]["direct_factors"]["relative_output_squared_error"],
            receipt["results"]["dense_reconstructed"]["relative_output_squared_error"],
            delta=1e-5,
        )
        for item in receipt["results"].values():
            self.assertGreaterEqual(item["median_seconds"], 0)

    def test_bounded_sampling_fails_closed(self):
        for arguments in ({"repeats": 0}, {"repeats": True}, {"repeats": 10001},
                          {"warmups": -1}, {"warmups": 1001}):
            with self.subTest(arguments=arguments):
                with self.assertRaises(ValueError):
                    measure(**arguments)


if __name__ == "__main__":
    unittest.main()
