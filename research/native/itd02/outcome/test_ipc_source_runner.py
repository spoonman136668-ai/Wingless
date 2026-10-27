"""Bounded source-to-IPC custody tests. Fixtures are not scientific episodes."""
import importlib.util
from pathlib import Path
import sys
import unittest

folder = Path(__file__).parent
for module in ("minigrid_source", "ipc_source_runner"):
    spec = importlib.util.spec_from_file_location(module, folder / (module + ".py"))
    obj = importlib.util.module_from_spec(spec)
    sys.modules[module] = obj
    spec.loader.exec_module(obj)

from minigrid_source import Public, SourceObservationError
from ipc_source_runner import BridgeError, run_training_episode


class ControlledSource:
    def __init__(self, done_at: int):
        self.calls = 0
        self.done_at = done_at
        self.current = None

    def reset(self, seed):
        self.seed = seed
        self.calls = 0
        self.current = Public(1, 1, 3, 3, 0, False)
        return self.current

    def step(self, action):
        self.calls += 1
        if action != 0:
            raise SourceObservationError("unexpected chosen action")
        self.current = Public(1, 1, 3, 3, 3, True)
        return self.current, 1.0, self.calls == self.done_at, False


class ControlledBridge:
    def __init__(self, reject_feedback=False):
        self.commands = []
        self.reject = reject_feedback

    def request(self, op, **payload):
        self.commands.append((op, payload))
        if op == "choose":
            return {"ok": True, "action": 0}
        if op == "feedback":
            if self.reject:
                raise BridgeError("custody rejected")
            return {"ok": True, "updates": payload["event_id"]}
        if op == "episode_start":
            return {"ok": True}
        raise AssertionError(op)


class TrainingBoundaryTests(unittest.TestCase):
    def test_exact_completed_feedback_and_episode_reset(self):
        world, native = ControlledSource(1), ControlledBridge()
        events, next_id, done = run_training_episode(world, native, seed=7, max_steps=3, start_id=19)
        self.assertTrue(done)
        self.assertEqual(next_id, 20)
        self.assertEqual(len(events), 1)
        self.assertEqual(events[0]["event_id"], 19)
        self.assertEqual(events[0]["selected_action"], 0)
        self.assertEqual(events[0]["reward"], 1.0)
        self.assertEqual([op for op, _ in native.commands], ["choose", "feedback", "episode_start"])
        feedback = native.commands[1][1]
        self.assertEqual(feedback["before"], Public(1, 1, 3, 3, 0, False).wire())
        self.assertEqual(feedback["after"], Public(1, 1, 3, 3, 3, True).wire())
        self.assertTrue(feedback["terminal"])
        self.assertFalse(feedback["truncated"])

    def test_step_budget_does_not_forge_terminal(self):
        world, native = ControlledSource(5), ControlledBridge()
        events, next_id, done = run_training_episode(world, native, 8, 1, 1)
        self.assertFalse(done)
        self.assertEqual(next_id, 2)
        self.assertFalse(events[0]["terminal"])
        self.assertFalse(events[0]["truncated"])
        self.assertNotIn("episode_start", [op for op, _ in native.commands])

    def test_rejected_feedback_aborts_without_reset(self):
        native = ControlledBridge(reject_feedback=True)
        with self.assertRaises(BridgeError):
            run_training_episode(ControlledSource(1), native, 9, 2, 1)
        self.assertNotIn("episode_start", [op for op, _ in native.commands])

    def test_rejects_non_native_and_malformed_counts(self):
        class Invalid(ControlledBridge):
            def request(self, op, **payload):
                if op == "choose":
                    return {"ok": True, "action": True}
                return super().request(op, **payload)
        with self.assertRaises(BridgeError):
            run_training_episode(ControlledSource(1), Invalid(), 1, 1, 1)
        class WrongCount(ControlledBridge):
            def request(self, op, **payload):
                result = super().request(op, **payload)
                if op == "feedback":
                    result["updates"] -= 1
                return result
        with self.assertRaises(BridgeError):
            run_training_episode(ControlledSource(1), WrongCount(), 1, 1, 1)

if __name__ == "__main__":
    unittest.main()
