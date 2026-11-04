"""CNC-02 frozen synthetic proof of learning without a weight master."""
import json
import unittest
import numpy as np
from compressed_learning import PackedCodeLearner, _code_at, _set_code, fit_control
from packed_int4 import apply_packed_int4, pack_int4


class PackedLearningTests(unittest.TestCase):
    def test_single_code_native_learning_generalizes_to_unseen_inputs(self):
        rng = np.random.default_rng(117)
        w = rng.standard_normal((4, 128)).astype(np.float32)
        student = PackedCodeLearner.from_weight(w)
        original = student.packed.copy()
        teacher_codes = original.copy()
        row, col = 1, 13
        code = _code_at(teacher_codes, row, col)
        new_code = code + 1 if code < 7 else code - 1
        _set_code(teacher_codes, row, col, new_code)
        train_x = rng.standard_normal((64, 128)).astype(np.float32)
        heldout_x = rng.standard_normal((64, 128)).astype(np.float32)
        train_y = apply_packed_int4(train_x, teacher_codes, student.scales)
        heldout_y = apply_packed_int4(heldout_x, teacher_codes, student.scales)
        def mse(a, b):
            return float(np.mean((a.astype(np.float64) - b.astype(np.float64))**2))
        before_train = mse(student.predict(train_x), train_y)
        before_heldout = mse(student.predict(heldout_x), heldout_y)
        report = student.train_step(train_x, train_y)
        after_train = mse(student.predict(train_x), train_y)
        after_heldout = mse(student.predict(heldout_x), heldout_y)
        self.assertTrue(report["changed"])
        self.assertEqual((report["updated_row"], report["updated_column"]), (row, col))
        self.assertTrue(np.array_equal(student.packed, teacher_codes))
        self.assertTrue(np.array_equal(student.scales, pack_int4(w)[1]))
        self.assertLess(after_train, before_train * 1e-8 + 1e-12)
        self.assertLess(after_heldout, before_heldout * 1e-8 + 1e-12)
        self.assertEqual(student.resident_weight_bytes, 4*64 + 4*2*4 + 24)
        self.assertEqual(set(vars(student)), {"packed", "scales", "steps"})
        self.assertTrue(np.array_equal(original, pack_int4(w)[0]))
        print("CNC02_COMPRESSED_LEARNING_SYNTHETIC " + json.dumps({
            "status": "SYNTHETIC_ONE_CODE_ADAPTATION_ONLY",
            "before_train_mse": before_train, "after_train_mse": after_train,
            "before_heldout_mse": before_heldout, "after_heldout_mse": after_heldout,
            "compressed_persistent_weight_bytes": student.resident_weight_bytes,
            "no_master_weight_or_optimizer_state": True, "real_qwen": False,
            "gpu_measured": False, "model_forward": False,
        }, sort_keys=True), flush=True)

    def test_saturated_code_and_no_update(self):
        w = np.zeros((2, 64), dtype=np.float32)
        learner = PackedCodeLearner.from_weight(w)
        x = np.ones((8, 64), dtype=np.float32)
        y = learner.predict(x)
        before = learner.packed.copy()
        result = learner.train_step(x, y)
        self.assertFalse(result["changed"])
        self.assertTrue(np.array_equal(before, learner.packed))

    def test_fail_closed_bad_inputs_and_fp32_reference_control(self):
        learner = PackedCodeLearner.from_weight(np.ones((2, 64), dtype=np.float32))
        with self.assertRaises(ValueError):
            learner.train_step(np.ones((2, 63)), np.ones((2, 2)))
        with self.assertRaises(ValueError):
            learner.train_step(np.ones((2, 64)), np.ones((2, 3)))
        with self.assertRaises(ValueError):
            learner.train_step(np.full((2, 64), np.nan), np.ones((2, 2)))
        with self.assertRaises(ValueError):
            fit_control(np.ones((2, 64)), np.ones((3, 64)), np.ones((3, 2)), steps=65)
        self.assertEqual(fit_control(np.ones((2, 64)), np.ones((3, 64)),
                                     np.ones((3, 2)), steps=0).shape, (2,64))


if __name__ == "__main__":
    unittest.main()
