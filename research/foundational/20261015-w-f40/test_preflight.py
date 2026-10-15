"""Pre-original tests only: never generate EVAL or reserved cohort."""
import importlib.util,pathlib,unittest
P=pathlib.Path(__file__).with_name("train_only.py")
s=importlib.util.spec_from_file_location("wf40",P)
m=importlib.util.module_from_spec(s);s.loader.exec_module(m)
class Gate(unittest.TestCase):
    def test_disjoint_and_reserved(self):
        self.assertFalse(set(m.TRAIN)&set(m.EVAL))
        self.assertFalse(set(m.EVAL)&set(m.RESERVED))
        self.assertFalse(set(m.TRAIN_FAMILIES)&set(m.EVAL_FAMILIES))
        with self.assertRaises(ValueError):m.case(303000,318,3)
    def test_truth_oracles_public_only_proposals(self):
        private,public,start,path=m.case(301000,310,3)
        self.assertEqual(len(public),24)
        self.assertEqual(m.follow(private,start,path),m.independent_follow(private,start,path))
        pool,features,worlds,answers,freq,early=m.proposals(public,start,path,19301)
        self.assertEqual(len(worlds),24)
        self.assertEqual(len(answers),24)
        self.assertTrue(all(len(features[q])==m.FEATURES for q in pool))
    def test_paid_mapping_receipts_at_most_two(self):
        hidden,public,start,path=m.case(301001,311,5)
        port=m.PaidPort(hidden,public)
        first=m.available(public)[0]
        self.assertEqual(port.buy(*first),hidden[first[0]][first[1]])
        with self.assertRaises(ValueError):port.buy(*first)
        second=m.available(public)[1];port.buy(*second)
        with self.assertRaises(ValueError):port.buy(*m.available(public)[2])
    def test_learned_train_only_fit_has_real_examples(self):
        weight,n=m.fit()
        self.assertEqual(len(weight),7)
        self.assertGreater(n,0)
        self.assertLessEqual(len(weight)*8,576)
if __name__=="__main__":unittest.main()
