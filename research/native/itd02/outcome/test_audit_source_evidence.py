"""Synthetic evidence auditor regression, not a real environment claim."""
import copy
import importlib.util
import json
from pathlib import Path
import sys
import unittest

spec = importlib.util.spec_from_file_location("audit_source_evidence", Path(__file__).with_name("audit_source_evidence.py"))
audit = importlib.util.module_from_spec(spec)
sys.modules[spec.name] = audit
spec.loader.exec_module(audit)


def fixture():
    state = {"x": 1, "y": 1, "goal_x": 3, "goal_y": 3, "direction": 0, "front_blocked": False}
    events = [{
        "seed": 7, "event_id": 1, "before": state, "after": state,
        "selected_action": 0, "reward": 0.0, "terminal": False, "truncated": False,
    }]
    state_bytes = "c3ludGhldGlj"
    import hashlib
    return {
        "schema": "wingless.itd02.source-training.v1",
        "status": "ENGINEERING_SOURCE_TRAINING_ONLY",
        "science": "NOT_ACCEPTED",
        "seeds": [7], "max_steps": 2, "events": events,
        "episodes": [{"seed": 7, "steps": 1, "return": 0.0, "source_done": False}],
        "event_digest_sha256": audit.digest(events),
        "checkpoint_base64": state_bytes,
        "checkpoint_sha256": hashlib.sha256(state_bytes.encode()).hexdigest(),
    }


class EvidenceTests(unittest.TestCase):
    def test_incomplete_zero_reward_is_not_acceptance(self):
        got = audit.assess(fixture())
        self.assertEqual(got["source_outcome"], "BOUNDED_INCOMPLETE_ZERO_REWARD")
        self.assertEqual(got["scientific_acceptance"], "NOT_ACCEPTED")

    def test_tampered_reward_is_rejected(self):
        record = fixture()
        record["events"][0]["reward"] = 9
        with self.assertRaises(audit.EvidenceError):
            audit.assess(record)

    def test_false_completion_rejected_with_recomputed_digest(self):
        record = fixture()
        record["episodes"][0]["source_done"] = True
        with self.assertRaises(audit.EvidenceError):
            audit.assess(record)

    def test_changed_action_with_recomputed_digest_rejected(self):
        record = fixture()
        record["events"][0]["selected_action"] = 3
        record["event_digest_sha256"] = audit.digest(record["events"])
        with self.assertRaises(audit.EvidenceError):
            audit.assess(record)

    def test_false_scientific_promotion_rejected(self):
        record = fixture()
        record["science"] = "ACCEPTED"
        with self.assertRaises(audit.EvidenceError):
            audit.assess(record)


if __name__ == "__main__":
    unittest.main()
