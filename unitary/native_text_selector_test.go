package unitary

import("bytes";"encoding/json";"testing")

func nativeTextExamples()[]NativeTextExample{return []NativeTextExample{{"aaaa","amber"},{"zzzz","violet"}}}
func TestNativeTextLearningAndCompleteState(t *testing.T){
 m,err:=TrainNativeTextSelector(nativeTextExamples());if err!=nil{t.Fatal(err)}
 before,err:=m.Checkpoint();if err!=nil{t.Fatal(err)}
 prediction,err:=m.Predict("aaaa");if err!=nil{t.Fatal(err)}
 if prediction.LanguageProof||prediction.NovelStringGeneration||prediction.OutputMode!="closed-response-selection"{t.Fatal("overstated output capability")}
 changed:=false
 for _,s:=range m.models{if s.model!=nci11WInit(false){changed=true}}
 if !changed{t.Fatal("native parameters did not learn")}
 untrained:=*m;untrained.models=nil
 for range m.models{untrained.models=append(untrained.models,&NativeSession{model:nci11WInit(false)})}
 control,err:=untrained.Predict("aaaa");if err!=nil{t.Fatal(err)}
 same:=true;for i,s:=range prediction.NativeScores{if s!=control.NativeScores[i]{same=false}}
 if same{t.Fatal("native learning has no score influence")}
 after,_:=m.Checkpoint();if !bytes.Equal(before,after){t.Fatal("prediction mutated learned state")}
 restored,err:=RestoreNativeTextSelector(before);if err!=nil{t.Fatal(err)}
 again,_:=restored.Checkpoint();if !bytes.Equal(before,again){t.Fatal("complete state changed after restore")}
 a,_:=json.Marshal(prediction);p,_:=restored.Predict("aaaa");b,_:=json.Marshal(p)
 if !bytes.Equal(a,b){t.Fatal("restored native scores changed")}
}
func TestNativeTextInvalidTrainingAndCheckpoint(t *testing.T){
 for _,xs:=range [][]NativeTextExample{nil,{{"x","a"},{"y","a"}},{{"","a"},{"y","b"}}}{
  if _,err:=TrainNativeTextSelector(xs);err==nil{t.Fatal("accepted invalid training")}
 }
 m,_:=TrainNativeTextSelector(nativeTextExamples());before,_:=m.Checkpoint()
 if _,err:=m.Predict("");err==nil{t.Fatal("accepted empty query")}
 after,_:=m.Checkpoint();if !bytes.Equal(before,after){t.Fatal("invalid query mutated state")}
 var e map[string]any;json.Unmarshal(before,&e);e["sha256"]="bad";bad,_:=json.Marshal(e)
 if _,err:=RestoreNativeTextSelector(bad);err==nil{t.Fatal("accepted corrupt checkpoint")}
 if _,err:=RestoreNativeTextSelector(append(before,[]byte("{}")...));err==nil{t.Fatal("accepted trailing JSON")}
}
