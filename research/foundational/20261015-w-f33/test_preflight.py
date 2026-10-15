"""W-F33 preflight: train-only source, demonstration privacy and model size."""
import importlib.util,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("pilot.py")
spec=importlib.util.spec_from_file_location("wf33",p)
m=importlib.util.module_from_spec(spec);sys.modules[spec.name]=m;spec.loader.exec_module(m)
class Contract(unittest.TestCase):
    def test_reserved_original_untouched(self):
        self.assertFalse(set(m.TRAIN_SEEDS)&set(m.RESERVED_TEST_SEEDS))
        self.assertFalse(set(m.TRAIN_FAMILIES)&set(m.RESERVED_TEST_FAMILIES))
    def test_train_public_demonstration_interface(self):
        secret,demos=m.source(191000,160)
        self.assertEqual(len(secret),6)
        self.assertEqual(len(demos),30)
        self.assertEqual((secret,demos),m.source(191000,160))
        self.assertEqual(len(m.learnt_logits(demos)),6)
    def test_temporal_order_composition(self):
        initial,roles,target=m.public_target(191000,160,5,m.source(191000,160)[0])
        self.assertEqual(m.predict(initial,roles,m.source(191000,160)[0]),target)
    def test_state_and_no_privileged_grammar_argument(self):
        self.assertLessEqual(m.ROLES*m.BITS*8+8,576)
        _,demos=m.source(191003,163)
        self.assertEqual(len(m.decoded_masks(m.learnt_logits(demos))),m.ROLES)
if __name__=="__main__":unittest.main()
