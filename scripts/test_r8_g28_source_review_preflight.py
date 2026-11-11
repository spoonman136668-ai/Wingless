"""Fail-closed exact-source R8 G28 preflight controls; not a review approval."""
from __future__ import annotations

import copy
from pathlib import Path
import re
import unittest

from scripts import r8_g28_source_review_preflight as m

ROOT = Path(__file__).resolve().parents[1]


class R8G28PreflightTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.result = m.run(ROOT)

    def test_full_exact_source_cohort_is_reconciled(self):
        x = self.result
        self.assertEqual(x["schema"], m.SCHEMA)
        self.assertEqual(x["accepted_source_sha"], m.ACCEPTED)
        self.assertEqual(x["structural_source_sha"], m.R8_INDEXED)
        self.assertEqual(x["historical_reviewed_identities"], 352)
        self.assertEqual(x["byte_identical_reviewed_declarations"], 351)
        self.assertEqual(x["new_symbol_count"], 303)
        self.assertEqual(len(x["new_unreviewed_provisional_declarations"]), 303)
        self.assertEqual(len({tuple(z["key"]) for z in
                              x["new_unreviewed_provisional_declarations"]}), 303)

    def test_changed_legacy_review_never_auto_promoted(self):
        x = self.result
        self.assertEqual(len(x["changed_historical_declarations"]), 1)
        old = x["changed_historical_declarations"][0]
        self.assertEqual(old["key"],
                         ["cmd/wingless/main.go", "", "main", "function"])
        self.assertEqual(old["old_span"], [19, 81])
        self.assertEqual(old["new_span"], [19, 88])
        self.assertIs(old["review_required"], True)
        self.assertIs(old["reviewed"], False)

    def test_new_declarations_are_only_provisional(self):
        records = self.result["new_unreviewed_provisional_declarations"]
        self.assertEqual(len({x["key"][0] for x in records}), 30)
        for row in records:
            self.assertEqual(row["target_source_commit"], m.R8_INDEXED)
            self.assertEqual(row["evidence_basis"], "lexical-provisional")
            self.assertIsNone(row["review_evidence"])
            self.assertIs(row["independent_approval"], False)
            self.assertRegex(row["file_sha256"], r"^[0-9a-f]{64}$")
            self.assertRegex(row["declaration_sha256"], r"^[0-9a-f]{64}$")
            self.assertNotIn("reviewer", row)
            self.assertNotIn("publication_authorized", row)

    def test_missing_authority_stays_explicit(self):
        x = self.result
        self.assertIs(x["source_review_completed"], False)
        self.assertIs(x["source_admission_authorized"], False)
        self.assertIs(x["scientific_first_use"], False)
        self.assertEqual(x["status"], "BLOCKED_INDEPENDENT_SOURCE_REVIEW")

    def test_readonly_blob_path_and_digest_negative(self):
        blob = b"package native\nfunc Hello() {}\n"
        good = m.digest(blob)
        self.assertEqual(m.verified_source(lambda sha, path: blob,
                                           "f" * 40, "native/x.go", good), blob)
        with self.assertRaisesRegex(m.R8G28Denied, "SOURCE_HASH_DRIFT"):
            m.verified_source(lambda sha, path: blob,
                              "f" * 40, "native/x.go", "0" * 64)
        for path in ("/etc/passwd", "../private.go", "src/../other.go",
                     "src\\shell.go"):
            with self.subTest(path=path):
                with self.assertRaisesRegex(m.R8G28Denied, "PATH_INVALID"):
                    m.verified_source(lambda sha, path: blob,
                                      "f" * 40, path, good)

    def test_invalid_declaration_span_fails_closed(self):
        raw = b"package native\nfunc Hello() {}\n"
        for args in ({"start": 0, "end": 1},
                     {"start": 1, "end": 4},
                     {"start": 2, "end": 1},
                     {"start": True, "end": 2}):
            with self.subTest(args=args):
                with self.assertRaisesRegex(m.R8G28Denied,
                                            "DECLARATION_SPAN_INVALID"):
                    m.declaration(raw, args)

    def test_sha_pins_are_not_soft_preferences(self):
        with self.assertRaisesRegex(m.R8G28Denied, "TARGET_SHA_NOT_PINNED"):
            m.inventory({}, {}, {}, lambda sha, path: b"",
                        target_sha="0" * 40)

    def test_no_fake_reviewer_receipt_is_generated(self):
        text = repr(self.result)
        self.assertNotIn("'reviewer_assertion'", text)
        self.assertNotIn("'source_review_completed': True", text)
        self.assertNotIn("'independent_approval': True", text)


if __name__ == "__main__":
    unittest.main()
