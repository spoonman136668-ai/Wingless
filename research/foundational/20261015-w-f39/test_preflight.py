"""W-F39 independent source and paid relationship receipt guards; no cohort."""
import importlib.util,pathlib,unittest
P=pathlib.Path(__file__).with_name("source_pilot.py")
s=importlib.util.spec_from_file_location("wf39",P)
m=importlib.util.module_from_spec(s);s.loader.exec_module(m)
class SourceGate(unittest.TestCase):
    def test_reserved_original_source_is_inaccessible(self):
        self.assertFalse(set(m.SEEDS)&set(m.RESERVED_SEEDS))
        with self.assertRaises(ValueError):m.make_case(292000,294,4)
    def test_separate_hidden_source_and_dual_oracle(self):
        hidden,public,start,roles=m.make_case(291000,290,4)
        self.assertEqual(len(public),24)
        self.assertEqual(m.walk(hidden,start,roles),m.independently_compose(hidden,start,roles))
        self.assertEqual(len(roles),4)
    def test_only_paid_relation_is_revealed(self):
        secret,public,start,query=m.make_case(291000,291,7)
        port=m.PaidMappingPort(secret,public)
        known={(r,x) for r,x,y in public}
        candidate=next((r,x) for r in range(8) for x in range(8) if (r,x) not in known)
        result=port.buy(*candidate)
        self.assertEqual(result,secret[candidate[0]][candidate[1]])
        self.assertEqual(len(port.receipts),1)
        with self.assertRaises(ValueError):port.buy(*candidate)
        second=next((r,x) for r in range(8) for x in range(8) if (r,x) not in known and (r,x)!=candidate)
        port.buy(*second)
        with self.assertRaises(ValueError):port.buy(7,7)
    def test_conditional_posterior_uses_only_public(self):
        _,public,start,roles=m.make_case(291001,292,4)
        answer,uncertainty,worlds,answers=m.posterior_answer(public,start,roles,100)
        self.assertEqual(len(worlds),24)
        self.assertEqual(len(answers),24)
        self.assertIn(answer,range(8))
        self.assertGreaterEqual(uncertainty,0)
        self.assertIn(m.propose("posterior_information_gain",public,start,roles,100),(None,*[(r,x) for r in range(8) for x in range(8)]))
if __name__=="__main__":unittest.main()
