"""Engineering regression of reward-mask boundary; not a source experiment."""
import importlib.util
from pathlib import Path
import sys
import unittest
from unittest.mock import patch

root=Path(__file__).parent
for name in ("minigrid_source","ipc_source_runner","source_controls","reward_ablation_restart"):
    if name not in sys.modules:
        spec=importlib.util.spec_from_file_location(name,root/(name+".py"))
        m=importlib.util.module_from_spec(spec)
        sys.modules[name]=m
        spec.loader.exec_module(m)

from reward_ablation_restart import MaskedOutcome

class FakeBridge:
    def __init__(self):
        self.messages=[]
    def request(self,op,**payload):
        self.messages.append((op,payload))
        return {"ok":True,"updates":payload.get("event_id",0)}

class MaskTests(unittest.TestCase):
    def test_reward_blind_receives_zero_but_preserves_selected_action(self):
        for blind,expected in ((True,0.0),(False,0.875)):
            instance=MaskedOutcome.__new__(MaskedOutcome)
            instance.bridge=FakeBridge()
            instance.next_id=1
            instance.blind=blind
            before={"x":1}
            after={"x":2}
            instance.feedback(before,2,after,.875,True,False)
            op,payload=instance.bridge.messages[0]
            self.assertEqual(op,"feedback")
            self.assertEqual(payload["reward"],expected)
            self.assertEqual(payload["action"],2)
            self.assertEqual(payload["before"],before)
            self.assertEqual(payload["after"],after)
            self.assertTrue(payload["terminal"])
            self.assertEqual(instance.next_id,2)

if __name__=="__main__":
    unittest.main()
