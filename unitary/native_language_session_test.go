package unitary

import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "testing"
)
// Invented engineering inputs, not an original scientific cohort.
func nativeSessionTestCase()NativeSessionCase {
 edges:=make([]NativeSessionEdge,12)
 for i:=range edges{edges[i]=NativeSessionEdge{i+1,i+2}}
 return NativeSessionCase{edges,[][]int{{1,5},{5,9}},[]bool{true,true}}
}
func TestNativeSessionPreservesModelAndMemory(t *testing.T){
 s,err:=TrainNativeSession([]NativeSessionCase{nativeSessionTestCase()},false)
 if err!=nil{t.Fatal(err)}
 if s.model==nci11WInit(false){t.Fatal("no learned parameter changed")}
 if err=s.ObserveGraph(nativeSessionTestCase().Edges);err!=nil{t.Fatal(err)}
 before,err:=s.Checkpoint();if err!=nil{t.Fatal(err)}
 restored,err:=RestoreNativeSession(before);if err!=nil{t.Fatal(err)}
 after,_:=restored.Checkpoint()
 if !bytes.Equal(before,after)||restored.model!=s.model||restored.memory!=s.memory{
  t.Fatal("complete learned state changed")
 }
 a,la,_:=s.Predict(1,5);b,lb,_:=restored.Predict(1,5)
 if a!=b||la!=lb{t.Fatal("query changed after restore")}
 zero:=*s;zero.memory=nci11WMemory{}
 _,lz,_:=zero.Predict(1,5)
 if la==lz{t.Fatal("native matrix has no influence on this engineering query")}
}
func TestNativeSessionInvalidInputsPreserveState(t *testing.T){
 s,err:=TrainNativeSession([]NativeSessionCase{nativeSessionTestCase()},true)
 if err!=nil{t.Fatal(err)}
 before,_:=s.Checkpoint()
 if s.ObserveGraph(nil)==nil{t.Fatal("accepted incomplete graph")}
 if _,_,err=s.Predict(0,5);err==nil{t.Fatal("accepted invalid query")}
 after,_:=s.Checkpoint();if !bytes.Equal(before,after){t.Fatal("invalid request mutated state")}
 for _,xs:=range [][]NativeSessionCase{nil,{{}}}{
  if _,err=TrainNativeSession(xs,false);err==nil{t.Fatal("accepted invalid training")}
 }
}
func TestNativeSessionCorruptionAndForeignSourceDenied(t *testing.T){
 s,_:=TrainNativeSession([]NativeSessionCase{nativeSessionTestCase()},false)
 raw,_:=s.Checkpoint()
 var e map[string]any;json.Unmarshal(raw,&e)
 e["sha256"]="bad";bad,_:=json.Marshal(e)
 if _,err:=RestoreNativeSession(bad);err==nil{t.Fatal("accepted changed digest")}
 if _,err:=RestoreNativeSession(append(raw,[]byte("{}")...));err==nil{t.Fatal("accepted trailing JSON")}
 if _,err:=RestoreNativeSession([]byte(`{"payload":{},"sha256":""}`));err==nil{t.Fatal("accepted incomplete state")}
 var full nativeSessionEnvelope;json.Unmarshal(raw,&full)
 full.Payload.SourceBlob="foreign"
 body,_:=json.Marshal(full.Payload);hash:=sha256.Sum256(body)
 full.SHA256=hex.EncodeToString(hash[:]);foreign,_:=json.Marshal(full)
 if _,err:=RestoreNativeSession(foreign);err==nil{t.Fatal("accepted foreign source with valid digest")}
 if ValidateNativeSessionJSON([]byte(`{"x":1,"x":2}`))==nil{t.Fatal("accepted duplicate JSON key")}
 if ValidateNativeSessionJSON([]byte(`{"x":null}`))==nil{t.Fatal("accepted null state")}
}
