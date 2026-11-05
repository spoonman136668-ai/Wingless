"""CNC-G9P witness-based storage rebase cannot grant GPU runtime admission."""
import unittest
from g9p_gpu_capacity_rebase import report, MIB, HARDWARE_EVIDENCE


class G9PTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls): cls.result=report()
    def test_hardware_identity_and_time_boundary(self):
        r=self.result
        self.assertEqual(r["hardware_receipt"]["original_run_id"],41030389354)
        self.assertEqual(r["hardware_receipt"]["original_artifact_id"],12966370600)
        self.assertEqual(r["gpu_total_bytes_from_archived_driver"],12227*MIB)
        self.assertEqual(r["gpu_free_bytes_at_archived_probe"],11812*MIB)
        self.assertEqual(r["historical_total_below_nominal_12gib_bytes"],61*MIB)
        self.assertEqual(r["historical_free_below_nominal_12gib_bytes"],476*MIB)
        self.assertFalse(r["hardware_measurement_is_current"])
        self.assertEqual(r["runner_online_now"],"UNVERIFIED")
        self.assertEqual(r["runner_usage_authorized"],"UNVERIFIED")
        self.assertEqual(r["scientific_acceptance"],"NOT_GRANTED")
    def test_all_historical_scenarios_reconcile(self):
        r=self.result
        self.assertEqual(r["scenario_count"],13)
        cases={x["name"]:x for x in r["scenarios"]}
        self.assertEqual(len(cases),13)
        base=cases["int2_block256_fp16"]
        self.assertEqual(base["remaining_against_historical_physical_total_bytes"],
                         448974848-61*MIB)
        self.assertEqual(base["remaining_against_historical_free_snapshot_bytes"],
                         448974848-476*MIB)
        self.assertTrue(base["would_fit_physical_total_analytically"])
        self.assertFalse(base["would_fit_historical_free_snapshot_analytically"])
        self.assertFalse(cases["sparse_k3"]["would_fit_historical_free_snapshot_analytically"])
        self.assertTrue(cases["lowrank_rank1"]["would_fit_physical_total_analytically"])
        self.assertFalse(cases["mixed_both_boundary_layers_int4"]["would_fit_historical_free_snapshot_analytically"])
        for c in cases.values():
            self.assertIsNone(c["current_available_vram_bytes"])
            self.assertFalse(c["cuda_context_peak_measured"])
            self.assertFalse(c["full_original_model_forward_executed"])
            self.assertEqual(c["total_planned_bytes_including_hypothetical_workspace"]+
                             c["remaining_against_historical_physical_total_bytes"],12227*MIB)
            self.assertEqual(c["total_planned_bytes_including_hypothetical_workspace"]+
                             c["remaining_against_historical_free_snapshot_bytes"],11812*MIB)


if __name__=="__main__": unittest.main()
