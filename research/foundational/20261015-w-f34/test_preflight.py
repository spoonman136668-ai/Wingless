"""W-F34 source-only checks; no reserved test cohort or pilot evaluation."""
import importlib.util
import pathlib
import unittest
here=pathlib.Path(__file__).resolve().parent
s=importlib.util.spec_from_file_location("wf34_source",here/"source_pilot.py")
m=importlib.util.module_from_spec(s);s.loader.exec_module(m)

class Gate(unittest.TestCase):
    def test_no_reserved_source_ids(self):
        self.assertFalse(set(m.SEEDS)&set(m.RESERVED_SEEDS))
        self.assertFalse(set(m.FAMILIES)&set(m.RESERVED_FAMILIES))
        with self.assertRaises(ValueError):
            m.hidden_operators(202000,184)
    def test_three_distinct_grammar_processes(self):
        for family in m.FAMILIES:
            ops=m.hidden_operators(201000,family)
            self.assertEqual(len(ops),8)
            self.assertTrue(all(sorted(op)==list(range(8)) for op in ops))
        self.assertNotEqual(m.hidden_operators(201000,180),m.hidden_operators(201000,181))
    def test_long_ordered_role_sequences_and_dual_oracle(self):
        for depth in m.DEPTHS:
            private,demo,x,roles=m.demos_and_query(201000,180,depth)
            self.assertEqual(len(roles),depth)
            self.assertEqual(len(demo),32)
            self.assertEqual(m.eval_source(private,x,roles),
                             m.eval_independent(private,x,roles))
    def test_train_only_empirical_control(self):
        private,demo,x,roles=m.demos_and_query(201001,182,6)
        fitted=m.empirical_map(demo)
        self.assertEqual(len(fitted),8)
        self.assertEqual(len(m.nearest_map(demo)),8)
        self.assertGreater(len(m.CANDIDATES),30)
if __name__=="__main__":unittest.main()
