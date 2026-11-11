package r7

import (
 "bytes"
 "encoding/json"
 "fmt"
 "testing"
)
func TestPublicFeedbackMustBeDelayed(t *testing.T){
 p:=New()
 s:=PublicObservation{Context:A,Sensors:[4]bool{true,false,true,false}}
 if err:=p.Learn(s,PublicFeedback{ObservedOutcome:true,DelaySteps:0});err==nil{t.Fatal("undelayed private teacher admitted")}
 if p.State.FeedbackEvents!=0 || p.State.Changes!=0{t.Fatalf("early state update: %+v",p.State)}
 if err:=p.Learn(s,PublicFeedback{ObservedOutcome:true,DelaySteps:2});err!=nil{t.Fatal(err)}
 if p.State.FeedbackEvents!=1 || p.State.Changes!=1{t.Fatalf("public update missing: %+v",p.State)}
 if _,e:=p.Predict(PublicObservation{Context:"H1",Sensors:s.Sensors});e==nil{t.Fatal("unadmitted private task context")}
}
func TestColdRestoreIdentityAndTamper(t *testing.T){
 p:=New()
 for j:=0;j<100;j++{
  o:=PublicObservation{Context:A,Sensors:[4]bool{j%2==0,j%3==0,j%4==0,j%5==0}}
  if err:=p.Learn(o,PublicFeedback{ObservedOutcome:j%3>0,DelaySteps:2});err!=nil{t.Fatal(err)}
 }
 raw,e:=p.Save();if e!=nil{t.Fatal(e)}
 restored,e:=Restore(raw);if e!=nil{t.Fatal(e)}
 if !bytes.Equal(raw,mustSave(t,restored)){t.Fatal("serialized state not exactly preserved")}
 changed:=bytes.Replace(raw,[]byte("wingless.r7.native.v1"),[]byte("wingless.r7.native.v2"),1)
 if _,e:=Restore(changed);e==nil{t.Fatal("tampered state identity admitted")}
 var v map[string]any
 if e:=json.Unmarshal(raw,&v);e!=nil{t.Fatal(e)}
 v["sha256"]="0"
 altered,_:=json.Marshal(v)
 if _,e:=Restore(altered);e==nil{t.Fatal("tampered data without checksum detected")}
}
func mustSave(t *testing.T,p *NativePredictor)[]byte{
 t.Helper();v,e:=p.Save();if e!=nil{t.Fatal(e)};return v
}
func TestNativeLearnerActuallyChangesWithPublicHistory(t *testing.T){
 p:=New();xs:=take(source(335),A,320)
 before,_:=accuracy(p,xs)
 if err:=training(p,xs);err!=nil{t.Fatal(err)}
 after,_:=accuracy(p,xs)
 if after < before+0.15 {t.Fatalf("native public learning did not improve: %.3f -> %.3f",before,after)}
 if p.State.Changes==0{t.Fatal("native weight learning did not occur")}
 if p.State.FeedbackEvents!=320{t.Fatalf("wrong feedback count %d",p.State.FeedbackEvents)}
}
func TestFrozenStudyAllThreeSeeds(t *testing.T){
 for _,seed:=range []int{17,31,47}{
  x,e:=Study(seed);if e!=nil{t.Fatal(e)}
  raw,_:=json.Marshal(x)
  t.Logf("WINGLESS_R7_PUBLIC_DEV_RESULT=%s",raw)
  if x.Updates!=640||!x.ColdRestoreExact{t.Fatalf("learning/restart absent on %d: %+v",seed,x)}
  if x.SerializedBytes>32768||x.Changes==0{t.Fatal("capacity/update gate")}
  if x.RetainedA<.75||x.B<.75{t.Fatalf("learned success criterion failed: %+v",x)}
  if x.PreA-x.RetainedA>.05+1e-9{t.Fatalf("retention prereg gate failed: %+v",x)}
  if (x.RetainedA-x.NoStateA)<.15||(x.B-x.NoStateB)<.15{t.Fatalf("erased memory ablation gate failed: %+v",x)}
  if (x.RetainedA-x.MajorityA)<.15||(x.B-x.MajorityB)<.15{t.Fatalf("current-only control gate failed: %+v",x)}
 }
}
func TestPublicSeedReplayNoTrainingOnTest(t *testing.T){
 for _,seed:=range []int{17,31,47}{
  a,e:=Study(seed);if e!=nil{t.Fatal(e)}
  b,e:=Study(seed);if e!=nil{t.Fatal(e)}
  if fmt.Sprintf("%+v",a)!=fmt.Sprintf("%+v",b){t.Fatalf("reproducibility fails seed %d",seed)}
 }
}
