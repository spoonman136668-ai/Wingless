"""Frozen W-R6 replay verification; no additional primary trials."""
import importlib.util
import json
from pathlib import Path
import unittest

ROOT=Path(__file__).resolve().parent
spec=importlib.util.spec_from_file_location("w_r6_replay",ROOT/"replay.py")
m=importlib.util.module_from_spec(spec)
spec.loader.exec_module(m)

class TestW_R6(unittest.TestCase):
    def test_exact_frozen_result(self):
        expected=json.loads((ROOT/"result.json").read_text(encoding="utf-8"))
        self.assertEqual(m.outcome(),expected)
        self.assertEqual(m.outcome(),m.outcome())
    def test_controls(self):
        arms=m.outcome()["arms"]
        self.assertEqual(arms["A0"]["signature_classes"],8)
        self.assertEqual(arms["A1"]["complete_alias_triplets"],9)
        self.assertEqual(arms["A2"]["complete_alias_triplets"],0)
        self.assertEqual(arms["A3"]["complete_alias_triplets"],0)
        self.assertEqual(m.outcome()["classification"],"SUPPORTED")

if __name__=="__main__":
    unittest.main()
