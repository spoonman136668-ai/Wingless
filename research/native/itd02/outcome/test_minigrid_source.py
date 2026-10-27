"""Dependency-free decoder tests; synthetic engineering fixtures, not source evidence."""
import importlib.util
from pathlib import Path
import unittest

spec = importlib.util.spec_from_file_location("minigrid_source", Path(__file__).with_name("minigrid_source.py"))
source = importlib.util.module_from_spec(spec)
import sys
sys.modules[spec.name] = source
spec.loader.exec_module(source)

class Image:
    shape = (5, 5, 3)
    def __init__(self, agent=(1, 1), direction=0, goal=(3, 3)):
        self.cells = {(x, y): (source.EMPTY, 0, 0) for x in range(5) for y in range(5)}
        for x in range(5):
            for y in (0, 4):
                self.cells[x, y] = (source.WALL, 0, 0)
        for y in range(5):
            for x in (0, 4):
                self.cells[x, y] = (source.WALL, 0, 0)
        if goal is not None:
            self.cells[goal] = (source.GOAL, 0, 0)
        self.cells[agent] = (source.AGENT, 0, direction)
    def __getitem__(self, pos):
        return self.cells[pos]

def public(agent=(1, 1), direction=0, goal=(3, 3)):
    return {"image": Image(agent, direction, goal), "direction": direction}

class Environment:
    def __init__(self):
        self.steps = []
    def reset(self, seed):
        self.seed = seed
        return public(), {}
    def step(self, action):
        self.steps.append(action)
        return public(agent=(2, 1)), 0.25, False, False, {}

class AdapterTests(unittest.TestCase):
    def test_public_decoding_and_bounded_step(self):
        env = Environment()
        episode = source.SourceEpisode(env)
        self.assertEqual(episode.reset(7).wire(), {
            "x": 1, "y": 1, "goal_x": 3, "goal_y": 3,
            "direction": 0, "front_blocked": False})
        after, reward, terminal, truncated = episode.step(2)
        self.assertEqual((after.x, reward, terminal, truncated), (2, .25, False, False))
        self.assertEqual(env.steps, [2])
        self.assertEqual(env.seed, 7)

    def test_rejects_ambiguous_public_image(self):
        obs = public()
        obs["image"].cells[2, 2] = (source.AGENT, 0, 0)
        with self.assertRaises(source.SourceObservationError):
            source.decode_public(obs)
        obs = public()
        obs["image"].cells[2, 2] = (source.GOAL, 0, 0)
        with self.assertRaises(source.SourceObservationError):
            source.decode_public(obs)
        obs = public()
        obs["direction"] = 1
        with self.assertRaises(source.SourceObservationError):
            source.decode_public(obs)
        obs = public()
        obs["image"].cells[2, 1] = (99, 0, 0)
        with self.assertRaises(source.SourceObservationError):
            source.decode_public(obs)

    def test_indexable_nonbuiltin_integer_direction_only(self):
        # Mimics a public NumPy scalar without importing optional MiniGrid.
        class Scalar:
            def __index__(self):
                return 0
        obs = public(direction=0)
        obs["direction"] = Scalar()
        self.assertEqual(source.decode_public(obs).direction, 0)
        for value in (True, 0.0, "0"):
            obs["direction"] = value
            with self.assertRaises(source.SourceObservationError):
                source.decode_public(obs)

    def test_goal_overlay_requires_prior_public_goal(self):
        obs = public(agent=(3, 3), goal=None)
        with self.assertRaises(source.SourceObservationError):
            source.decode_public(obs)
        result = source.decode_public(obs, prior_goal=(3, 3))
        self.assertEqual((result.goal_x, result.goal_y), (3, 3))

    def test_no_hidden_state_or_illegal_action(self):
        env = Environment()
        episode = source.SourceEpisode(env)
        episode.reset(1)
        for illegal in (True, -1, 3, 1.0):
            with self.assertRaises(source.SourceObservationError):
                episode.step(illegal)
        self.assertEqual(env.steps, [])
        episode.step(0)
        self.assertEqual(env.steps, [0])

    def test_fail_closed_after_invalid_environment_step(self):
        class Broken(Environment):
            def step(self, action):
                return public(), float("nan"), False, False, {}
        env = Broken()
        episode = source.SourceEpisode(env)
        episode.reset(2)
        with self.assertRaises(source.SourceObservationError):
            episode.step(2)
        with self.assertRaises(source.SourceObservationError):
            episode.step(2)

    def test_blocked_front_and_closed_door(self):
        obs = public()
        obs["image"].cells[2, 1] = (source.WALL, 0, 0)
        self.assertTrue(source.decode_public(obs).front_blocked)
        obs["image"].cells[2, 1] = (source.DOOR, 0, 1)
        self.assertTrue(source.decode_public(obs).front_blocked)
        obs["image"].cells[2, 1] = (source.DOOR, 0, 0)
        self.assertFalse(source.decode_public(obs).front_blocked)

if __name__ == "__main__":
    unittest.main()
