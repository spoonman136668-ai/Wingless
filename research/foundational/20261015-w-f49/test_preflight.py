import importlib.util
from pathlib import Path
import sys
import unittest
p=Path(__file__).parent/"source_pilot.py"
s=importlib.util.spec_from_file_location("w_f49_preflight",p)
m=importlib.util.module_from_spec(s);sys.modules[s.name]=m;s.loader.exec_module(m)

class PublicRankF49(unittest.TestCase):
    def test_reserved_source_blocked(self):
        for seed,family in ((432000,644),(431100,644),(411100,640),(372000,640),(392000,640)):
            with self.assertRaises(ValueError):
                m.source(seed,family,5)
    def test_registered_diagnostic_and_family(self):
        self.assertEqual(len(tuple(m.cases())),64)
        self.assertFalse(set(m.FIT)&set(m.DIAG))
        self.assertFalse(set(m.ORIGINAL)&set(m.DIAG))
        self.assertFalse(set(m.FAMILIES)&set(m.ORIGINAL_FAMILIES))
    def test_public_action_rank_repeatable(self):
        _,public,start,chain=m.source(431100,640,3)
        a=m.rank_stability(list(public),start,chain,[0])
        b=m.rank_stability(list(public),start,chain,[0])
        self.assertEqual(a,b)
        if a[0] is not None:
            self.assertIn(a[0],m.portlib.candidate_actions(public))
    def test_new_relationship_chains(self):
        a=m.source(431100,640,5)
        b=m.source(431100,643,5)
        self.assertNotEqual(a[3],b[3])
    def test_no_third_paid_read(self):
        maps,public,start,chain=m.source(431100,640,3)
        port=m.portlib.PurchasePort(maps,public)
        choices=m.portlib.candidate_actions(public)
        port.buy(*choices[0]);port.buy(*choices[1])
        with self.assertRaises(ValueError):
            port.buy(*choices[2])

if __name__=="__main__":unittest.main()
