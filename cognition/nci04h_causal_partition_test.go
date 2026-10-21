package cognition
import (
 "bytes"
 "math"
 "path/filepath"
 "testing"
 "unsafe"
)

// Tests one necessary learning condition within the original model budget.
// It is not an independently selected unfamiliar task or scientific result.
func TestNCI04HContradictionDrivenNativeStateRevision(t *testing.T){
 var s LatentState
 src:=[4]float32{0,0,0,0}
 dst:=[4]float32{1,0,0,0}
 for n:=0;n<100;n++ {
  if err:=s.Observe(src,dst,0,1);err!=nil{t.Fatal(err)}
 }
 strong,err:=s.Predict(src,dst,0)
 if err!=nil||strong<.75{t.Fatalf("prior not established: %v %v",strong,err)}
 baseline,_:=s.Snapshot()
 same,err:=s.ObservePartitioned(src,dst,0,0)
 if err!=nil||!same.Abstained||same.Split{t.Fatalf("indistinguishable evidence wrongly split: %+v %v",same,err)}
 afterUnknown,_:=s.Snapshot()
 if !bytes.Equal(baseline,afterUnknown){t.Fatal("unidentifiable conflict modified existing learned capability")}
 nearby:=[4]float32{.10,0,0,0}
 changed,err:=s.ObservePartitioned(nearby,dst,0,0)
 if err!=nil||!changed.Split||changed.Count!=3{t.Fatalf("source-distinguishable conflict not partitioned %+v %v",changed,err)}
 preserved,_:=s.Predict(src,dst,0)
 newProb,_:=s.Predict(nearby,dst,0)
 if math.Abs(preserved-strong)>1e-12 || newProb>=.5{t.Fatalf("prior knowledge lost or counterexample overclaimed: old=%.5f new=%.5f",preserved,newProb)}
 for n:=0;n<45;n++ {
  u,e:=s.ObservePartitioned(nearby,dst,0,0)
  if e!=nil||u.Split{t.Fatalf("repeat novel evidence duplicated identity: %+v %v",u,e)}
 }
 learned,_:=s.Predict(nearby,dst,0)
 if learned>=.30{t.Fatalf("new partition did not learn from observed consequences %.5f",learned)}
 aliased:=[4]float32{.102,0,0,0}
 transferred,_:=s.Predict(aliased,dst,0)
 if math.Abs(transferred-learned)>1e-10{t.Fatal("failed sensory alias transfer inside newly learned causal state")}
 receipt,_:=s.Snapshot()
 var resumed LatentState
 if e:=resumed.Restore(receipt);e!=nil{t.Fatal(e)}
 other,_:=resumed.Snapshot()
 if !bytes.Equal(receipt,other){t.Fatal("persistent causal partitions not restored")}
 if unsafe.Sizeof(s)!=unsafe.Sizeof(LatentState{}){t.Fatal("native model budget grew")}
 t.Logf("NCI04H_WINGLESS_ENGINEERING_PASS fixed_state_bytes=%d learned_split=1 strong_prior=%.6f novel_probability=%.6f",unsafe.Sizeof(s),strong,learned)
}

func TestNCI04HICEPositiveEvidenceAndNoUnverifiedNegatives(t *testing.T){
 root,err:=filepath.Abs("..");if err!=nil{t.Fatal(err)}
 r,err:=OpenSourceReasoner(root);if err!=nil{t.Fatal(err)}
 before:=r.State.Observations
 result,err:=r.ObserveVerifiedCallPartition("runICESourceCognition","ObserveExistingCall")
 if err!=nil||result.Split||result.Abstained||r.State.Observations!=before+1{t.Fatalf("native ICE-positive update failed: %+v %v",result,err)}
 beforeSnap,_:=r.State.Snapshot()
 if _,err:=r.ObserveVerifiedCallPartition("runICESourceCognition","PredictCall");err==nil{
  t.Fatal("unsupported source negative incorrectly admitted")
 }
 afterSnap,_:=r.State.Snapshot()
 if !bytes.Equal(beforeSnap,afterSnap){t.Fatal("unverified source evidence changed model")}
}
