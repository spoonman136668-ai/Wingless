"""W-F57 frozen FIT method checks. Never uses diagnostic/original cohorts."""
import importlib.util,json,pathlib,sys,unittest

HERE=pathlib.Path(__file__).resolve().parent
spec=importlib.util.spec_from_file_location("w57",HERE/"source_fit.py")
m=importlib.util.module_from_spec(spec)
sys.modules[spec.name]=m
spec.loader.exec_module(m)
class FITBelief(unittest.TestCase):
    def test_sealed_folds_rejected(self):
        for seed,family,length in ((661100,1244,10),(662000,1248,8),
                                  (661000,1244,8),(661100,1240,8)):
            with self.assertRaises(ValueError):m.source(seed,family,length)
    def test_source_matches_frozen_prereg(self):
        p=json.loads((HERE/"prereg.json").read_text())
        self.assertEqual(p["fit"]["seed_range"],[m.FIT.start,m.FIT.stop])
        self.assertEqual(p["fit"]["families"],list(m.FAMILY))
        self.assertEqual(p["fit"]["lengths"],[8,10])
    def test_only_public_information_enters_policy(self):
        case=m.source(661000,1240,8)
        clone=dict(case)
        clone["evaluator_truth"]=(case["evaluator_truth"]+1)%m.P
        clone["evaluator_private_maps"]={j:(1,0) for j in range(4)}
        weights=(1.,1.,1.,1.)
        self.assertEqual(m.factorized_transition(case,weights),
                         m.factorized_transition(clone,weights))
        self.assertEqual(m.frozen_F56_budgeted_beam(case),
                         m.frozen_F56_budgeted_beam(clone))
        self.assertEqual(m.exhaustive_public_bayes(case),
                         m.exhaustive_public_bayes(clone))
    def test_sample_fit_likelihood_and_state_normalization(self):
        case=m.source(661001,1241,10)
        distributions=m.posterior(case)
        self.assertEqual(set(distributions),{0,1,2,3})
        for prob in distributions.values():
            self.assertAlmostEqual(sum(prob),1.,places=10)
            self.assertEqual(len(prob),len(m.MAPS))
        prediction,operations=m.factorized_transition(case,(1.,1.,1.,1.))
        self.assertIn(prediction,range(m.P))
        self.assertLess(operations, len(m.MAPS)**4*len(case["chain"])/2)
    def test_fit_temperature_learning_only(self):
        params,n=m.train_fit_temperature()
        self.assertEqual(len(params),4)
        self.assertEqual(n,64)
        self.assertTrue(all(x in (.75,1.,1.5,2.) for x in params))

if __name__=="__main__":unittest.main()
