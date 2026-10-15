"""W-F44 non-science preflight only; never consume all prospective TRAIN cases."""
import importlib.util,pathlib,sys,unittest
path=pathlib.Path(__file__).with_name("source_pilot.py")
spec=importlib.util.spec_from_file_location("w_f44_preflight",path)
m=importlib.util.module_from_spec(spec);sys.modules[spec.name]=m;spec.loader.exec_module(m)

class Gate(unittest.TestCase):
    def test_disjoint_source(self):
        self.assertFalse(set(m.SEEDS)&set(m.RESERVED_SEEDS))
        self.assertFalse(set(m.FAMILIES)&set(m.RESERVED_FAMILIES))
        with self.assertRaises(ValueError):m.source(352000,434,5)
    def test_reused_inverse_typed_roles(self):
        _,p,start,q=m.source(351000,430,5)
        self.assertEqual(len(p),12)
        self.assertEqual(q[0],(0,1))
        self.assertEqual(q[2],(0,1))
        self.assertEqual(q[1],(1,-1))
    def test_exact_counterfactual_is_public_only(self):
        private,p,start,chain=m.source(351000,431,3)
        posterior=m.prior.worlds(p)
        self.assertEqual(len(posterior),1296)
        allowed=m.prior.candidate_actions(p)
        effort=[0]
        selected=m.choose("two_step_counterfactual",list(p),start,chain,2,effort)
        self.assertTrue(selected is None or selected in allowed)
        self.assertGreater(effort[0],0)
        # All score arguments came from public triples, not private mapping.
        self.assertIn(private,posterior)
    def test_two_purchase_guard(self):
        private,p,start,chain=m.source(351000,432,3)
        port=m.prior.PurchasePort(private,p)
        choices=m.prior.candidate_actions(p)
        port.buy(*choices[0]);port.buy(*choices[1])
        with self.assertRaises(ValueError):port.buy(*choices[2])
    def test_zero_one_termination_available(self):
        private,p,start,chain=m.source(351000,430,3)
        self.assertIsNone(m.choose("no_purchase",list(p),start,chain,2,[0]))
        complete=[(r,x,private[r][x]) for r in range(4) for x in range(6)]
        self.assertIsNone(m.choose("two_step_counterfactual",complete,start,chain,2,[0]))

if __name__=="__main__":unittest.main()
