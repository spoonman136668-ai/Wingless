"""Preflight is source-only: no evaluation of all train or reserved outcomes."""
import importlib.util, pathlib, unittest
p = pathlib.Path(__file__).with_name("source_pilot.py")
spec = importlib.util.spec_from_file_location("wf43", p)
m = importlib.util.module_from_spec(spec)
spec.loader.exec_module(m)

class Preflight(unittest.TestCase):
    def test_source_and_reserved_separation(self):
        self.assertFalse(set(m.SEEDS) & set(m.RESERVED_SEEDS))
        self.assertFalse(set(m.FAMILIES) & set(m.RESERVED_FAMILIES))
        with self.assertRaises(ValueError):
            m.source(332000, 394, 3)
    def test_exact_support_certified(self):
        private, public, start, chain = m.source(331000, 390, 3)
        self.assertEqual(len(m.worlds(public)), 1296)
        self.assertEqual(len(public), 12)
        self.assertIn(private, m.worlds(public))
    def test_exact_rank_independent_of_random_seed(self):
        _, public, start, chain = m.source(331000, 392, 3)
        a, cost = m.choose("exact_answer_entropy", public, start, chain, 11)
        b, _ = m.choose("exact_answer_entropy", public, start, chain, 987654)
        self.assertEqual(a, b)
        self.assertEqual(cost, 1296)
    def test_matched_world_evaluations(self):
        _, public, start, chain = m.source(331000, 391, 3)
        for arm in ("sampled_information_gain_equal_world_count",
                    "sampled_decision_value_equal_world_count",
                    "query_frontier_equal_world_count"):
            action, cost = m.choose(arm, public, start, chain, 271)
            self.assertIn(action, m.candidate_actions(public))
            self.assertEqual(cost, 1296)
    def test_purchase_limit_duplicate_and_zero_step(self):
        private, public, start, chain = m.source(331000, 390, 3)
        port = m.PurchasePort(private, public)
        first, _ = m.choose("exact_answer_entropy", public, start, chain, 11)
        port.buy(*first)
        with self.assertRaises(ValueError): port.buy(*first)
        with self.assertRaises(ValueError): port.buy(public[0][0], public[0][1])
        second = next(p for p in m.candidate_actions(public) if p != first)
        port.buy(*second)
        with self.assertRaises(ValueError): port.buy(*next(p for p in m.candidate_actions(public) if p not in (first,second)))
        self.assertIsNone(m.choose("no_purchase", public, start, chain, 11)[0])
    def test_malformed_public_bijection_rejected(self):
        with self.assertRaises(ValueError): m.worlds([(0, 1, 2), (0, 3, 2)])

if __name__ == "__main__":
    unittest.main()
