"""T14 source-only generator. Import does NOT originate any TRAIN/OPEN source.
Frozen T10 mathematical source grammar; unused, disjoint T14 negative IDs.
Only the scorer may receive the hidden final truth.
"""
import types
TRAIN=tuple(range(-41031,-40999))
OPEN=tuple(range(-41163,-41099))
def make_source(t10):
    scope=dict(t10.source.__globals__,TRAIN=range(-41031,-40999),
               OPEN=range(-41163,-41099))
    return types.FunctionType(t10.source.__code__,scope,t10.source.__name__,t10.source.__defaults__)
def birth(seed,t10,k):
    assert seed in TRAIN or seed in OPEN
    return make_source(t10)(seed,k)
