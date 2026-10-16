#!/usr/bin/env python3
"""Offline mutation tests for the read-only R35 original-job-log parser."""
from __future__ import annotations
import importlib.util
import json
from pathlib import Path
import unittest

HERE = Path(__file__).resolve().parent
spec = importlib.util.spec_from_file_location("r35_original_log_verifier",
                                               HERE / "verify_r35_original_log.py")
verify = importlib.util.module_from_spec(spec)
spec.loader.exec_module(verify)
ADVISORY = json.loads((HERE / "r35-original-scientific-observation-advisory.json").read_text())


def original_shape():
    authentic = dict(source_sample_index=43, authentic=True,
                     full5_alias_class_count=23,
                     full5_cross_entity_collision_pair_count=0,
                     pair_preserving_all4_match_count=0,
                     pair_preserving_partial_match_count=0,
                     pair_preserving_zero_match_count=127)
    controls = [
        dict(source_sample_index=i, authentic=False, full5_alias_class_count=23,
             full5_cross_entity_collision_pair_count=0,
             pair_preserving_all4_match_count=120,
             pair_preserving_partial_match_count=0,
             pair_preserving_zero_match_count=7)
        for i in range(43)
    ]
    return dict(schema=verify.SCHEMA, experiment=verify.EXPERIMENT,
                classification="record_content_necessary",
                metrics=dict(ADVISORY["recovered_from_original_job_log"]["metrics"]),
                arms=[authentic] + controls)


def action_log(result):
    return ("\n".join((
        "2026-10-16T17:11:06.1259828Z " +
            json.dumps(result, separators=(",", ":")),
        "2026-10-16T17:11:06.2230506Z WINGLESS_AUTOMATED_QUALIFICATION_CLASSIFICATION: qualified-scientific-result",
        "2026-10-16T17:11:06.2313050Z WINGLESS_AUTOMATED_QUALIFICATION_PASS",
        "2026-10-16T17:11:08.1626415Z REPO_LOCAL_RESEARCH_CYCLE_PASS",
        "2026-10-16T17:11:08.1634808Z experiment=" + verify.EXPERIMENT,
    ))) + "\n"


class FrozenR35LogTests(unittest.TestCase):
    def test_matching_original_result_advisory_only(self):
        verdict = verify.verify_log(action_log(original_shape()), ADVISORY)
        self.assertEqual(verdict["arms"], 44)
        self.assertEqual(verdict["scientific_classification"],
                         "record_content_necessary")
        self.assertFalse(verdict["transport_authenticated_by_this_script"])
        self.assertFalse(verdict["independent_scientific_acceptance"])
        self.assertFalse(verdict["experiment_replayed"])

    def test_reject_duplicated_sample_identity(self):
        sample = original_shape()
        sample["arms"][12]["source_sample_index"] = 44
        with self.assertRaises(verify.EvidenceMismatch):
            verify.verify_log(action_log(sample), ADVISORY)

    def test_reject_mismatch_in_duplicate_controls(self):
        sample = original_shape()
        sample["arms"][-1]["pair_preserving_zero_match_count"] = 8
        with self.assertRaises(verify.EvidenceMismatch):
            verify.verify_log(action_log(sample), ADVISORY)

    def test_reject_changed_metric_even_if_markers_pass(self):
        sample = original_shape()
        sample["metrics"]["duplicate_unchanged_arm_count"] = 42
        with self.assertRaises(verify.EvidenceMismatch):
            verify.verify_log(action_log(sample), ADVISORY)

    def test_reject_missing_wrapper_pass(self):
        with self.assertRaises(verify.EvidenceMismatch):
            verify.verify_log(action_log(original_shape()).replace(
                "WINGLESS_AUTOMATED_QUALIFICATION_PASS", "MISSING"), ADVISORY)

    def test_reject_conflicting_result_copies(self):
        a = original_shape()
        b = original_shape()
        b["arms"][0]["full5_alias_class_count"] = 24
        combined = action_log(a) + action_log(b)
        with self.assertRaises(verify.EvidenceMismatch):
            verify.verify_log(combined, ADVISORY)

    def test_reject_truncated_json_line(self):
        log = action_log(original_shape())
        cut = log.splitlines()
        cut[0] = cut[0].split('"arms"')[0] + '"arms"'
        with self.assertRaises(verify.EvidenceMismatch):
            verify.verify_log("\n".join(cut), ADVISORY)

    def test_reject_no_scientific_result_even_with_wrapper_pass(self):
        log = action_log(original_shape())
        cut = "\n".join(log.splitlines()[1:])
        with self.assertRaises(verify.EvidenceMismatch):
            verify.verify_log(cut, ADVISORY)

    def test_reject_noninteger_metric(self):
        sample = original_shape()
        sample["metrics"]["duplicate_sufficient_arm_count"] = False
        with self.assertRaises(verify.EvidenceMismatch):
            verify.verify_log(action_log(sample), ADVISORY)

    def test_git_blob_hash_stable(self):
        self.assertEqual(verify.git_blob_sha(b"hello"),
                         "b6fc4c620b67d95f953a5c1c1230aaab5db5a1b0")


if __name__ == "__main__":
    unittest.main()
