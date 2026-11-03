"""Frozen-data and provenance negative controls for activation replay."""
import unittest
import numpy as np
from activation_replay import MODEL, REVISION, activation_digest, evaluate_activation_output


def fixture(x):
    return {
        "model": MODEL,
        "revision": REVISION,
        "capture_kind": "synthetic_test_fixture",
        "cohort": "development",
        "tensor_name": "model.layers.0.mlp.experts.0.down_proj.weight",
        "record_ids": [f"synthetic:{i}" for i in range(len(x))],
        "activation_sha256": activation_digest(x),
    }


class ReplayTests(unittest.TestCase):
    def setUp(self):
        self.x = np.array([[2., 0.], [1., 0.]], dtype=np.float32)
        self.w = np.eye(2, dtype=np.float32)
        self.p = fixture(self.x)

    def test_identical_weights_zero_error(self):
        output=evaluate_activation_output(self.w, {"exact":self.w.copy()}, self.x,self.p)
        self.assertEqual(output["candidate_metrics"]["exact"]["relative_output_squared_error"],0)
        self.assertEqual(output["source_authenticity"],"UNVERIFIED_REQUIRES_INDEPENDENT_ORIGINAL_ARTIFACT_INSPECTION")

    def test_activation_metric_separates_equal_weight_error(self):
        first=self.w.copy(); first[0,0]=0
        second=self.w.copy(); second[0,1]=1
        metrics=evaluate_activation_output(self.w,{"high_impact":first,"low_impact":second},self.x,self.p)["candidate_metrics"]
        self.assertEqual(metrics["high_impact"]["relative_weight_squared_error"],metrics["low_impact"]["relative_weight_squared_error"])
        self.assertGreater(metrics["high_impact"]["relative_output_squared_error"],metrics["low_impact"]["relative_output_squared_error"])

    def test_checksum_tampering_fails(self):
        x=self.x.copy(); x[0,0]+=1
        with self.assertRaises(ValueError):evaluate_activation_output(self.w,{"x":self.w},x,self.p)

    def test_protected_cohort_fails(self):
        p=dict(self.p,cohort="protected")
        with self.assertRaises(ValueError):evaluate_activation_output(self.w,{"x":self.w},self.x,p)

    def test_wrong_model_fails(self):
        p=dict(self.p,model="OTHER")
        with self.assertRaises(ValueError):evaluate_activation_output(self.w,{"x":self.w},self.x,p)

    def test_duplicate_ids_fail(self):
        p=dict(self.p,record_ids=["duplicate","duplicate"])
        with self.assertRaises(ValueError):evaluate_activation_output(self.w,{"x":self.w},self.x,p)

    def test_native_kind_without_receipt_fails(self):
        p=dict(self.p,capture_kind="native_model_forward_hook")
        with self.assertRaises(ValueError):evaluate_activation_output(self.w,{"x":self.w},self.x,p)

    def test_native_kind_with_claimed_identity_still_not_accepted(self):
        p=dict(self.p,capture_kind="native_model_forward_hook",source_commit="a"*40,source_artifact_sha256="b"*64,source_workflow_run_id=7,calibration_record_ids=["different:1"])
        output=evaluate_activation_output(self.w,{"x":self.w},self.x,p)
        self.assertEqual(output["source_authenticity"],"UNVERIFIED_REQUIRES_INDEPENDENT_ORIGINAL_ARTIFACT_INSPECTION")

    def test_calibration_leak_rejected(self):
        p=dict(self.p,capture_kind="native_model_forward_hook",source_commit="a"*40,source_artifact_sha256="b"*64,source_workflow_run_id=7,calibration_record_ids=["synthetic:1"])
        with self.assertRaises(ValueError):evaluate_activation_output(self.w,{"x":self.w},self.x,p)

    def test_nonfinite_candidate_fails(self):
        w=self.w.copy();w[0,0]=np.nan
        with self.assertRaises(ValueError):evaluate_activation_output(self.w,{"bad":w},self.x,self.p)

    def test_shape_mismatch_fails(self):
        with self.assertRaises(ValueError):evaluate_activation_output(self.w,{"bad":self.w[:1]},self.x,self.p)

    def test_zero_reference_output_fails(self):
        with self.assertRaises(ValueError):evaluate_activation_output(np.zeros_like(self.w),{"zero":np.zeros_like(self.w)},self.x,self.p)


if __name__=="__main__": unittest.main()
