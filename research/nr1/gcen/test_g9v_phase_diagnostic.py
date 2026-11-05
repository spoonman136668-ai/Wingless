"""Offline fail-closed G9V screening tests: no weights, network or CUDA."""
import unittest
from pathlib import Path
from g9v_phase_diagnostic import decide, SEEDS


def row(seed, packed=3.0, dense=2.0, traced=3.1, error=0.01, share=0.7):
    return {"seed":seed,"dense_median_ms":float(dense),"packed_median_ms":float(packed),
            "trace_median_ms":float(traced),"relative_squared_error":float(error),
            "unpack_decode_fraction":float(share)}


class G9VScienceControls(unittest.TestCase):
    def test_fidelity_stop_outranks_fusion(self):
        self.assertEqual(decide([row(SEEDS[0],error=0.12),row(SEEDS[1],error=0.11)]),
                         "STOP_INT2_FIDELITY_SCREEN")

    def test_unstable_seeds_stop(self):
        self.assertEqual(decide([row(SEEDS[0],packed=2),row(SEEDS[1],packed=3)]),
                         "NO_GO_UNSTABLE_GPU_LATENCY")

    def test_trace_overhead_stop(self):
        self.assertEqual(decide([row(SEEDS[0],traced=5.0),row(SEEDS[1])]),
                         "NO_GO_TRACE_OVERHEAD")

    def test_no_decode_hotspot_stop(self):
        self.assertEqual(decide([row(SEEDS[0],share=0.4),row(SEEDS[1])]),
                         "NO_GO_NO_DOMINANT_DECODE")

    def test_positive_screen_is_not_science_acceptance(self):
        self.assertEqual(decide([row(SEEDS[0]),row(SEEDS[1])]),
                         "POTENTIAL_FUSION_STUDY_NOT_SCIENTIFIC_PASS")

    def test_wrong_seeds_rejected(self):
        for rows in ([],[row(SEEDS[0])],[row(SEEDS[0]),row(SEEDS[0])]):
            with self.subTest(rows=rows), self.assertRaisesRegex(ValueError,"seed"):
                decide(rows)

    def test_nonfinite_and_zero_rejected(self):
        for rows in ([row(SEEDS[0],packed=float("nan")),row(SEEDS[1])],
                     [row(SEEDS[0],dense=0),row(SEEDS[1])]):
            with self.subTest(rows=rows), self.assertRaisesRegex(ValueError,"metric"):
                decide(rows)

    def test_no_checkpoint_download_or_full_model_launch(self):
        code=Path(__file__).with_name("g9v_phase_diagnostic.py").read_text(encoding="utf-8")
        for forbidden in ("hf_hub_download(", "snapshot_download(", "from_pretrained(",
                          "urlopen(", "requests.get(", ".generate("):
            self.assertNotIn(forbidden,code)
        self.assertIn("check_checkout(root, source_sha)",code)
        self.assertIn("load_verified_original_expert(shard, manifest)",code)
        self.assertIn('scientific_acceptance": "NOT_GRANTED"',code)


if __name__ == "__main__":
    unittest.main()
