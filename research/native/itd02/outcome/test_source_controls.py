"""Synthetic custody tests for source controls (not scientific episodes)."""
import importlib.util
from pathlib import Path
import sys
import unittest

root = Path(__file__).parent
for name in ("minigrid_source", "ipc_source_runner", "source_controls"):
    if name not in sys.modules:
        spec = importlib.util.spec_from_file_location(name, root / (name + ".py"))
        module = importlib.util.module_from_spec(spec)
        sys.modules[name] = module
        spec.loader.exec_module(module)

from minigrid_source import Public, SourceObservationError
from source_controls import FixedPolicy, episode


class World:
    def __init__(self, at):
        self.at = at
        self.n = 0
        self.calls = []

    def reset(self, seed):
        self.n = 0
        self.seed = seed
        return Public(1, 1, 3, 3, 0, False)

    def step(self, selected):
        self.n += 1
        self.calls.append(selected)
        return Public(1, 1, 3, 3, 0, True), float(selected == 2), self.n == self.at, False


class ControlsTests(unittest.TestCase):
    def test_fixed_public_legal_actions_and_rewards(self):
        env = World(2)
        policy = FixedPolicy()
        record = episode(env, policy, 11, 2)
        self.assertEqual(record["steps"], 2)
        self.assertEqual(env.calls, [0, 2])
        self.assertEqual(record["return"], 1.0)
        self.assertTrue(record["completed"])
        self.assertEqual(len(record["events_sha256"]), 64)

    def test_bound_exhaustion_not_false_terminal(self):
        env = World(10)
        policy = FixedPolicy()
        record = episode(env, policy, 7, 3)
        self.assertEqual(record["steps"], 3)
        self.assertFalse(record["completed"])
        self.assertFalse(record["events"][-1]["terminated"])

    def test_bad_selected_action_refused(self):
        class Bad(FixedPolicy):
            def choose(self, before):
                return True
        env = World(1)
        with self.assertRaises(Exception):
            episode(env, Bad(), 1, 2)
        self.assertEqual(env.calls, [])

    def test_rejected_source_feedback_stops(self):
        class Broken(World):
            def step(self, selected):
                raise SourceObservationError("source rejected")
        with self.assertRaises(SourceObservationError):
            episode(Broken(1), FixedPolicy(), 1, 2)

if __name__ == "__main__":
    unittest.main()
