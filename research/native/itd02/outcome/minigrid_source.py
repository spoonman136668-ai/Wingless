"""Fail-closed MiniGrid public observation/transition adapter for ITD-02.

No learner IPC, training, private env attributes, or scientific acceptance here.
Install explicitly: python -m pip install 'minigrid==2.3.1' 'gymnasium==0.29.1'
The import is lazy so the decoder can be unit tested without those dependencies.
"""
from __future__ import annotations

from dataclasses import dataclass
import math
import operator
from typing import Any

# MiniGrid's published OBJECT_TO_IDX encoding. Refuse incompatible versions
# in a real runner instead of silently reinterpreting object identifiers.
EMPTY = 1
WALL = 2
FLOOR = 3
DOOR = 4
GOAL = 8
LAVA = 9
AGENT = 10
PASSABLE = {EMPTY, FLOOR, GOAL}
BLOCKED = {WALL, LAVA}

class SourceObservationError(ValueError):
    pass

@dataclass(frozen=True)
class Public:
    x: int
    y: int
    goal_x: int
    goal_y: int
    direction: int
    front_blocked: bool

    def wire(self) -> dict[str, Any]:
        return {
            "x": self.x, "y": self.y,
            "goal_x": self.goal_x, "goal_y": self.goal_y,
            "direction": self.direction,
            "front_blocked": self.front_blocked,
        }

def _cell(image: Any, x: int, y: int) -> tuple[int, int, int]:
    try:
        values = image[x, y]
        fields = tuple(int(v) for v in values)
    except (TypeError, ValueError, IndexError, KeyError) as exc:
        raise SourceObservationError("invalid encoded grid cell") from exc
    if len(fields) != 3 or any(v < 0 or v > 255 for v in fields):
        raise SourceObservationError("invalid encoded grid channels")
    return fields

def decode_public(observation: dict[str, Any], prior_goal: tuple[int, int] | None = None) -> Public:
    """Decode only FullyObsWrapper's public image/direction; never env.unwrapped."""
    if not isinstance(observation, dict) or "image" not in observation:
        raise SourceObservationError("missing public image")
    image = observation["image"]
    shape = getattr(image, "shape", None)
    if shape is None or len(shape) != 3 or int(shape[2]) != 3:
        raise SourceObservationError("not a full encoded grid image")
    width, height = int(shape[0]), int(shape[1])
    if width < 3 or height < 3 or width > 64 or height > 64:
        raise SourceObservationError("unsupported public grid dimensions")
    agents, goals = [], []
    for x in range(width):
        for y in range(height):
            obj, _color, state = _cell(image, x, y)
            if obj == AGENT:
                agents.append((x, y, state))
            elif obj == GOAL:
                goals.append((x, y))
    if len(agents) != 1:
        raise SourceObservationError("expected exactly one public agent")
    x, y, direction = agents[0]
    if direction not in (0, 1, 2, 3):
        raise SourceObservationError("invalid encoded direction")
    published_direction = observation.get("direction")
    # Random-start MiniGrid emits NumPy integer scalars for direction. Accept
    # exact indexable integers, never float or bool, and cross-check against
    # the independently encoded agent direction in the public image.
    if isinstance(published_direction, bool):
        raise SourceObservationError("invalid public direction type")
    try:
        numeric_direction = operator.index(published_direction)
    except TypeError as exc:
        raise SourceObservationError("invalid public direction type") from exc
    if numeric_direction != direction:
        raise SourceObservationError("public direction mismatch")
    if len(goals) == 1:
        goal = goals[0]
        if prior_goal is not None and prior_goal != goal:
            raise SourceObservationError("public goal changed during episode")
    elif len(goals) == 0 and prior_goal == (x, y):
        # FullyObsWrapper overlays the agent on a goal cell at termination.
        goal = prior_goal
    else:
        raise SourceObservationError("missing or ambiguous public goal")
    dx, dy = ((1, 0), (0, 1), (-1, 0), (0, -1))[direction]
    fx, fy = x + dx, y + dy
    if not (0 <= fx < width and 0 <= fy < height):
        raise SourceObservationError("front cell outside public image")
    front, _color, door_state = _cell(image, fx, fy)
    if front in PASSABLE:
        blocked = False
    elif front in BLOCKED:
        blocked = True
    elif front == DOOR and door_state in (0, 1, 2):
        blocked = door_state != 0
    else:
        raise SourceObservationError("unknown front traversability")
    return Public(x, y, goal[0], goal[1], direction, blocked)

def source_action(action: int) -> int:
    """MiniGrid: left=0, right=1, forward=2. Reject bool and other actions."""
    if type(action) is not int or action not in (0, 1, 2):
        raise SourceObservationError("non-native selected action")
    return action

class SourceEpisode:
    """Opt-in authentic environment holder; no access to hidden env state."""
    def __init__(self, environment: Any):
        self.environment = environment
        self.current: Public | None = None
        self.closed = True

    def reset(self, seed: int) -> Public:
        if type(seed) is not int or not 0 <= seed <= 2**32 - 1:
            raise SourceObservationError("invalid preregistered seed")
        raw, _info = self.environment.reset(seed=seed)
        candidate = decode_public(raw)
        self.current, self.closed = candidate, False
        return candidate

    def step(self, selected_action: int) -> tuple[Public, float, bool, bool]:
        if self.closed or self.current is None:
            raise SourceObservationError("no live selected episode")
        legal = source_action(selected_action)
        # A world step cannot be undone: close the episode on any bad return.
        self.closed = True
        raw, reward, terminated, truncated, _info = self.environment.step(legal)
        if type(terminated) is not bool or type(truncated) is not bool or (terminated and truncated):
            raise SourceObservationError("invalid source termination flags")
        value = float(reward)
        if not math.isfinite(value) or abs(value) > 100:
            raise SourceObservationError("source reward out of native range")
        new = decode_public(raw, (self.current.goal_x, self.current.goal_y))
        self.current = new
        self.closed = terminated or truncated
        return new, value, terminated, truncated

def open_source_environment(environment_id: str) -> SourceEpisode:
    """Fixed full-observation source; environment ID must be predeclared."""
    if environment_id not in ("MiniGrid-Empty-5x5-v0", "MiniGrid-Empty-Random-5x5-v0", "MiniGrid-Empty-8x8-v0"):
        raise SourceObservationError("environment not on bounded allowlist")
    import gymnasium as gym
    import minigrid  # noqa: F401 - registers source environments
    from minigrid.wrappers import FullyObsWrapper
    from minigrid.core.constants import OBJECT_TO_IDX
    expected = {
        "empty": EMPTY, "wall": WALL, "floor": FLOOR, "door": DOOR,
        "goal": GOAL, "lava": LAVA, "agent": AGENT,
    }
    if any(OBJECT_TO_IDX.get(k) != v for k, v in expected.items()):
        raise SourceObservationError("MiniGrid source encoding drift")
    if getattr(minigrid, "__version__", None) != "2.3.1":
        raise SourceObservationError("unpinned MiniGrid version")
    if getattr(gym, "__version__", None) != "0.29.1":
        raise SourceObservationError("unpinned Gymnasium version")
    return SourceEpisode(FullyObsWrapper(gym.make(environment_id)))
