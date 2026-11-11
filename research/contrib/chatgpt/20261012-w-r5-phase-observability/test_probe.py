"""Mechanical-only controls; no R4 primary stream or scientific outcome is evaluated."""
import unittest
try:
    from .probe import phase_signature, induce_classes
except ImportError:
    from probe import phase_signature, induce_classes

class MechanismTests(unittest.TestCase):
    def setUp(self):
        self.rows = ((1,2,3,4),(4,3,2,1))
    def test_support_order_invariance(self):
        self.assertEqual(phase_signature(self.rows,"A0"),phase_signature(tuple(reversed(self.rows)),"A0"))
    def test_phase_preserves_order_observation(self):
        self.assertNotEqual(phase_signature(self.rows,"A1"),phase_signature(tuple(reversed(self.rows)),"A1"))
    def test_no_entity_label_in_induction(self):
        self.assertEqual(induce_classes(["alpha","beta","alpha"]),((0,2),(1,)))
    def test_reject_bool_as_byte(self):
        with self.assertRaises(ValueError): phase_signature(((True,2,3,4),),"A0")
    def test_reject_duplicate_rows(self):
        with self.assertRaises(ValueError): phase_signature((self.rows[0],self.rows[0]),"A0")
    def test_reject_missing_rotation(self):
        with self.assertRaises(ValueError): phase_signature(self.rows,"A2")
    def test_reject_unknown_arm(self):
        with self.assertRaises(ValueError): phase_signature(self.rows,"A3")
    def test_rotation_index_is_deterministic(self):
        rows=tuple((i,i+16,i+32,i+48) for i in range(16))
        self.assertEqual(phase_signature(rows,"A2",family_control_index=0)[1],rows[5])
        self.assertEqual(phase_signature(rows,"A2",family_control_index=2)[1],rows[15])

if __name__ == "__main__": unittest.main()
