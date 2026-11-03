"""G30 offline safety/contract tests; pinned Torch runtime runs separately in CI."""
import ast
from pathlib import Path
import unittest
CODE = Path(__file__).with_name("g30_sequential_staging.py").read_text(encoding="utf-8")


class G30SourceContract(unittest.TestCase):
    def test_syntax(self):
        ast.parse(CODE)

    def test_digest_before_deserialize(self):
        self.assertIn("reader.read(65536)", CODE)
        self.assertIn("h.hexdigest() != expected", CODE)
        self.assertLess(CODE.index("verified_stage(src / filename"), CODE.index("with safe_open("))
        self.assertIn("stage.unlink(missing_ok=True)", CODE)

    def test_negative_control(self):
        self.assertIn("negative digest-control failure", CODE)
        self.assertIn("broken.write_bytes(", CODE)

    def test_parity_and_offload(self):
        for marker in ("test_runtime_identity()", "offload_weight(", "save_offload_index(",
                       "dispatch_model(", "output parity failure", "expert-input parity failure"):
            self.assertIn(marker, CODE)

    def test_resources_and_claim_boundaries(self):
        for marker in ("ru_maxrss", "st_blocks * 512",
                       '"peak_staging_shard_count": peak_staging_shard_count',
                       '"staging_count_samples": staging_count_samples',
                       '"original_checkpoint_weights_loaded": False',
                       '"full_original_model_streaming": "UNVERIFIED"',
                       '"scientific_acceptance": "NOT_GRANTED"'):
            self.assertIn(marker, CODE)

    def test_staging_peak_is_observed_not_assumed(self):
        self.assertIn('peak_staging_shard_count = 0', CODE)
        self.assertIn('staging_count_samples.append(before_count)', CODE)
        self.assertIn('staging_count_samples.append(staged_count)', CODE)
        self.assertIn('staging_count_samples.append(cleaned_count)', CODE)
        self.assertIn('peak_staging_shard_count = max(peak_staging_shard_count, staged_count)', CODE)
        self.assertIn('staged_count != 1 or present[0] != stage', CODE)
        self.assertNotIn('"peak_staging_shard_count": 1', CODE)

    def test_ci_is_mandatory_and_prerequisite_is_checked(self):
        workflow = (Path(__file__).resolve().parents[3] / ".github" / "workflows" / "native.yml").read_text(encoding="utf-8")
        section = workflow.split("      - name: GCEN G30 sequential synthetic shard staging and dispatch parity", 1)[1].split("      - name: GCEN pinned upstream metadata custody probe", 1)[0]
        self.assertNotIn("continue-on-error: true", section)
        self.assertIn("G30 requires G25 prerequisite PASS", section)
        self.assertLess(section.index("G30 requires G25 prerequisite PASS"), section.index("python research/nr1/gcen/g30_sequential_staging.py"))
        self.assertIn('if: always()', workflow.split('uses: actions/upload-artifact@v4', 1)[1])

    def test_no_network_or_large_model(self):
        for forbidden in ("from_pretrained(", "snapshot_download(", "hf_hub_download(",
                          "urlopen(", "requests.get("):
            self.assertNotIn(forbidden, CODE)


if __name__ == "__main__":
    unittest.main()
