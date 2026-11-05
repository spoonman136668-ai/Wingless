"""G9T current-GPU fail-closed controls, no GPU dependency or heavyweight CI."""
import unittest
from types import SimpleNamespace
from g9t_gpu_admission import inspect_current_cuda, MIB, SCHEMA

SHA="f"*40
TOTAL=12227*MIB

class FakeCuda:
    def __init__(self,*,available=True,devices=1,free=11812*MIB,
                 total=TOTAL,index=0,name="NVIDIA GeForce RTX 5070 Ti Laptop GPU"):
        self.available=available
        self.devices=devices
        self.free=free
        self.total=total
        self.index=index
        self.name=name
    def is_available(self):return self.available
    def device_count(self):return self.devices
    def current_device(self):return self.index
    def get_device_properties(self,index):
        return SimpleNamespace(name=self.name,total_memory=self.total)
    def mem_get_info(self,index):return self.free,self.total
    def get_device_capability(self,index):return (12,0)

def fake_torch(**kwargs):
    return SimpleNamespace(cuda=FakeCuda(**kwargs),version=SimpleNamespace(cuda="12.8"),
                           __version__="2.10.0+cu128")

class GPUAdmissionTests(unittest.TestCase):
    def admit(self,torch=None,sha=SHA,minimum_free_mib=1024):
        return inspect_current_cuda(torch or fake_torch(),
               expected_source_sha=sha,minimum_free_mib=minimum_free_mib)
    def test_positive_observation_is_not_authority_or_execution(self):
        r=self.admit()
        self.assertEqual(r["schema"],SCHEMA)
        self.assertEqual(r["physical_total_vram_bytes"],TOTAL)
        self.assertEqual(r["currently_free_vram_bytes"],11812*MIB)
        self.assertTrue(r["probe_is_current_in_this_process"])
        self.assertFalse(r["runner_use_authorized_by_probe"])
        self.assertFalse(r["cuda_kernel_executed"])
        self.assertFalse(r["original_shards_loaded"])
        self.assertEqual(r["scientific_acceptance"],"NOT_GRANTED")
    def test_reject_cpu_fallback(self):
        with self.assertRaisesRegex(RuntimeError,"CPU fallback"):
            self.admit(fake_torch(available=False))
    def test_reject_multiple_devices(self):
        with self.assertRaisesRegex(RuntimeError,"single-GPU"):
            self.admit(fake_torch(devices=2))
    def test_reject_nonzero_device(self):
        with self.assertRaisesRegex(RuntimeError,"active GPU"):
            self.admit(fake_torch(index=1))
    def test_reject_insufficient_current_vram(self):
        with self.assertRaisesRegex(RuntimeError,"floor"):
            self.admit(fake_torch(free=512*MIB))
    def test_reject_wrong_gpu_family(self):
        with self.assertRaisesRegex(RuntimeError,"non-NVIDIA"):
            self.admit(fake_torch(name="Example GPU"))
    def test_reject_invalid_exact_sha(self):
        with self.assertRaisesRegex(ValueError,"SHA"):
            self.admit(sha="unknown")
    def test_reject_floor_weakening(self):
        with self.assertRaisesRegex(ValueError,"weakened"):
            self.admit(minimum_free_mib=256)
    def test_reject_mismatched_physical_memory(self):
        gpu=fake_torch()
        gpu.cuda.mem_get_info=lambda index:(512*MIB,TOTAL+1)
        with self.assertRaisesRegex(RuntimeError,"disagree"):
            self.admit(gpu)
    def test_reject_cpu_only_torch_distribution(self):
        torch=fake_torch()
        torch.version.cuda=None
        with self.assertRaisesRegex(RuntimeError,"CPU-only"):
            self.admit(torch)

if __name__=="__main__":unittest.main()
