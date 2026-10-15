"""F38 source-only checks. Do not generate full TRAIN-only pilot or reserved IDs."""
import importlib.util,pathlib,random,unittest
P=pathlib.Path(__file__).with_name("source_pilot.py")
s=importlib.util.spec_from_file_location("wf38",P);m=importlib.util.module_from_spec(s);s.loader.exec_module(m)
class Contract(unittest.TestCase):
    def test_reserved_id_fails_closed(self):
        self.assertFalse(set(m.SEEDS)&set(m.RESERVED_SEEDS))
        self.assertFalse(set(m.FAMILIES)&set(m.RESERVED_FAMILIES))
        with self.assertRaises(ValueError):m.gen(282000,284,3,4)
    def test_dual_truth_and_provable_public_equivalence(self):
        hidden,public,start,query=m.gen(281000,280,3,4)
        self.assertEqual(len(public),24)
        self.assertEqual(m.truth_walk(hidden,start,query),m.dual_truth(hidden,start,query))
        alt=m.completions(public,random.Random(123))
        self.assertTrue(all(alt[role][x]==y for role,x,y in public))
        self.assertTrue(all(sorted(role)==list(range(8)) for role in alt))
    def test_dense_public_determines_bijection_uniquely(self):
        hidden,public,start,query=m.gen(281000,281,7,4)
        a=m.completions(public,random.Random(1))
        b=m.completions(public,random.Random(2))
        self.assertEqual(a,b)
        self.assertEqual(a,hidden)
        self.assertEqual(len(public),56)
    def test_strong_equal_evidence_controls(self):
        hidden,public,start,query=m.gen(281001,282,3,7)
        a,b=m.map_controls(public)
        self.assertEqual(len(a),8)
        self.assertEqual(len(b),8)
        self.assertTrue(all(sorted(role)==list(range(8)) for role in a))
        self.assertEqual(m.truth_walk(hidden,start,((0,1),(0,-1))),start)
if __name__=="__main__":unittest.main()
