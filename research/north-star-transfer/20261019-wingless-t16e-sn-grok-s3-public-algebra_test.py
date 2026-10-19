"""T16E MIT-upstream S3 public algebra equivalence. No TRAIN/OPEN cohorts."""
import importlib.util,sys,unittest
from pathlib import Path
D=Path(__file__).parent
def load(name,path):
 s=importlib.util.spec_from_file_location(name,path)
 m=importlib.util.module_from_spec(s);sys.modules[s.name]=m;s.loader.exec_module(m)
 return m
m=load("t16e_upstream_s3",D/"20261019-wingless-t16e-sn-grok-s3-public-algebra.py")
k=load("t16e_k_frozen",D/"20261018-wingless-t3-typed-composition.py")
class S3SemanticsTests(unittest.TestCase):
 def test_exact_36_multiplication_pairs(self):
  result=m.check_wingless_public_algebra(k)
  self.assertEqual(result["exhaustive_products"],36)
  self.assertTrue(result["exact_semantics_preserved"])
  self.assertFalse(result["novel_family_discovered"])
  self.assertFalse(result["original_firstuse_performed"])
 def test_noncommutativity_preserved(self):
  p=(1,0,2);q=(0,2,1)
  self.assertNotEqual(m.upstream_s3_compose(p,q),m.upstream_s3_compose(q,p))
 def test_associativity_216_triples(self):
  for a in m.DOMAIN:
   for b in m.DOMAIN:
    for c in m.DOMAIN:
     self.assertEqual(m.upstream_s3_compose(m.upstream_s3_compose(a,b),c),
       m.upstream_s3_compose(a,m.upstream_s3_compose(b,c)))
 def test_group_identities_inverses(self):
  e=(0,1,2)
  for p in m.DOMAIN:
   q=m.inverse(p)
   self.assertEqual(m.upstream_s3_compose(p,q),e)
   self.assertEqual(m.upstream_s3_compose(q,p),e)
 def test_wrong_operator_orientation_detected(self):
  p=(1,2,0);q=(1,0,2)
  self.assertNotEqual(m.upstream_s3_compose(p,q),tuple(p[q[i]] for i in range(3)))
 def test_external_commit_pin_nonadmission(self):
  self.assertEqual(m.UPSTREAM_COMMIT,"cb0673dba621a92b1ae4566cb2baf449b7fe6e06")
  self.assertEqual(m.UPSTREAM_PERMUTATION_BLOB,"aa2489d8e528bad04ba60af36cf12fff11a52dd7")
  receipt=m.check_wingless_public_algebra(k)
  self.assertFalse(receipt["independent_witness_verified"])
  self.assertFalse(receipt["optimizer_admission"])
  self.assertFalse(receipt["full_resource_parity_measured"])
if __name__=="__main__":unittest.main()
