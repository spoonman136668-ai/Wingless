"""Source-free fault injection of exact frozen exposure control flow."""
import ast
import json
from pathlib import Path
import tempfile
import unittest

SOURCE = Path("research/native/itd02/outcome/exact_exposure_gate.py")


def harness(run):
    tree = ast.parse(SOURCE.read_text(encoding="utf8"))
    names = {"ExposureIncomplete", "exact_train", "execute"}
    nodes = [n for n in tree.body if getattr(n, "name", None) in names]
    ns = dict(BridgeError=RuntimeError, TRAIN_STEPS=256,
              TRAIN_SEEDS=tuple(range(12)), run_episode=run,
              TRAIN_ENV="source-free-fixture", json=json)
    exec(compile(ast.Module(body=nodes, type_ignores=[]), str(SOURCE), "exec"), ns)
    return ns


class ExposureEvidenceTests(unittest.TestCase):
    def test_short_completed_episodes_preserve_original_negative(self):
        calls = []
        def run(env, policy, seed, budget):
            calls.append(seed)
            return {"steps": 1, "completed": True}
        ns = harness(run)
        with self.assertRaises(ns["ExposureIncomplete"]) as failure:
            ns["exact_train"](object())
        self.assertEqual(calls, list(range(12)))
        self.assertEqual(failure.exception.remaining, 244)
        self.assertEqual(len(failure.exception.records), 12)

    def test_incomplete_episode_never_resets_again(self):
        calls = []
        def run(*args):
            calls.append(args)
            return {"steps": 100, "completed": False}
        ns = harness(run)
        with self.assertRaises(ns["ExposureIncomplete"]):
            ns["exact_train"](object())
        self.assertEqual(len(calls), 1)

    def test_exact_budget_and_final_cut_are_unchanged(self):
        budgets = []
        def run(env, policy, seed, budget):
            budgets.append(budget)
            return {"steps": budget, "completed": budget == 100}
        ns = harness(run)
        records = ns["exact_train"](object())
        self.assertEqual(budgets, [100, 100, 56])
        self.assertEqual(sum(r["steps"] for r in records), 256)

    def test_negative_receipt_preserved_without_success_or_overwrite(self):
        ns = harness(None)
        def incomplete(*args):
            raise ns["ExposureIncomplete"]("fixture incomplete", [{"steps": 1}], 255)
        ns["_execute"] = incomplete
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "receipt.json"
            with self.assertRaises(ns["ExposureIncomplete"]):
                ns["execute"]("unused", output)
            record = json.loads(output.read_text())
            self.assertEqual(record["observed_interactions"], 1)
            self.assertFalse(record["comparison_complete"])
            self.assertFalse(record["evaluation_authorized_by_receipt"])
            saved = output.read_bytes()
            with self.assertRaises(FileExistsError):
                ns["execute"]("unused", output)
            self.assertEqual(output.read_bytes(), saved)


if __name__ == "__main__":
    unittest.main()
