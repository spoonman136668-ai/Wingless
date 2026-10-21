package cognition

import (
 "bytes"
 "encoding/json"
 "testing"
 "unsafe"
)

func TestNCI04G3FocusedCausalRevisionAndRestart(t *testing.T) {
 var learner CausalPartition
 sensor:=[FeatureWidth]float32{.25,.1,0,0}
 quiet:=[FeatureWidth]float32{0,1,0,0}
 diagnostic:=[FeatureWidth]float32{1,0,0,0}
 offered:=[][FeatureWidth]float32{quiet,diagnostic}
 if _,known,err:=learner.SelectIntervention(sensor,offered);err!=nil||known{t.Fatalf("untrained intervention asserted: known=%v err=%v",known,err)}
 for i:=0;i<24;i++ {
  if err:=learner.ObserveCausal(sensor,quiet,0,1);err!=nil{t.Fatal(err)}
  if err:=learner.ObserveCausal(sensor,diagnostic,i%2,1-i%2);err!=nil{t.Fatal(err)}
 }
 chosen,known,err:=learner.SelectIntervention(sensor,offered)
 if err!=nil||!known||!equalVec(chosen,diagnostic){t.Fatalf("failed to select consequence-informative descriptor: %v %v %v",chosen,known,err)}
 old,ok,err:=learner.Reidentify(sensor,diagnostic,0)
 if err!=nil||!ok||old!=1{t.Fatalf("prior mode overwritten: %d %v %v",old,ok,err)}
 revised,ok,err:=learner.Reidentify(sensor,diagnostic,1)
 if err!=nil||!ok||revised!=0{t.Fatalf("new distinction missing: %d %v %v",revised,ok,err)}
 unknownSensor:=[FeatureWidth]float32{-.9,0,0,0}
 if _,known,_:=learner.SelectIntervention(unknownSensor,offered);known{t.Fatal("invented unknown intervention")}
 if _,ok,_:=learner.Reidentify(unknownSensor,diagnostic,1);ok{t.Fatal("invented unknown state")}
 snapshot,err:=learner.CausalSnapshot();if err!=nil{t.Fatal(err)}
 var again CausalPartition
 if err:=again.CausalRestore(snapshot);err!=nil{t.Fatal(err)}
 replay,err:=again.CausalSnapshot();if err!=nil||!bytes.Equal(snapshot,replay){t.Fatal("snapshot/restart changed learned state",err)}
 if v,ok,_:=again.Reidentify(sensor,diagnostic,0);!ok||v!=1{t.Fatal("old relation did not survive revision and restart")}
 corrupted:=append([]byte(nil),snapshot...)
 var wire map[string]json.RawMessage
 if err:=json.Unmarshal(corrupted,&wire);err!=nil{t.Fatal(err)}
 wire["count"]=[]byte("17")
 corrupted,_=json.Marshal(wire)
 before,_:=again.CausalSnapshot()
 if err:=again.CausalRestore(corrupted);err==nil{t.Fatal("accepted excess modes")}
 after,_:=again.CausalSnapshot()
 if !bytes.Equal(before,after){t.Fatal("failed restore mutated accepted state")}
 if unsafe.Sizeof(learner)>unsafe.Sizeof(LatentState{}){t.Fatal("candidate resident struct exceeds exact-parent bound")}
 t.Logf("NCI04G3_FOCUSED_PASS modes=%d trained_probe_descriptors=%d struct_bytes=%d native_struct_bytes=%d",learner.Count,learner.ProbeCount,unsafe.Sizeof(learner),unsafe.Sizeof(LatentState{}))
}

func TestNCI04G3NoHiddenIdentityOrConflictingUpdate(t *testing.T) {
 var s CausalPartition
 sensor:=[FeatureWidth]float32{0,.3,0,0}
 probe:=[FeatureWidth]float32{-.5,.2,0,0}
 if err:=s.ObserveCausal(sensor,probe,0,1);err!=nil{t.Fatal(err)}
 before,_:=s.CausalSnapshot()
 if err:=s.ObserveCausal(sensor,probe,0,0);err==nil{t.Fatal("conflicting historical knowledge overwritten")}
 if err:=s.ObserveCausal(sensor,probe,7,0);err==nil{t.Fatal("invalid consequence accepted")}
 after,_:=s.CausalSnapshot();if !bytes.Equal(before,after){t.Fatal("failed observations mutated state")}
}
