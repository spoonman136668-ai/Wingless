package cognition

import (
 "bytes"
 "encoding/json"
 "fmt"
 "math"
 "testing"
 "unsafe"
)
func TestExperienceRawObservationAndPersistence(t *testing.T){
 m:=NewExperienceMemory()
 signatures:=[4][4]float64{{-.9,-.8,-.1,.1},{-.3,-.1,-.1,.1},{.4,.5,-.1,.1},{.9,.8,-.1,.1}}
 for repetition:=0;repetition<14;repetition++{
  for id,ob:=range signatures{
   after:=ob
   after[2],after[3]=.9,-.9
   if err:=m.ObserveExperience(ob,after,id%2,float64(id%2));err!=nil{t.Fatal(err)}
  }
 }
 used:=0;for _,x:=range m.Tracks{if x.Used{used++}}
 if used==0||used>ExperienceCapacity{t.Fatalf("wrong bounded allocation: %d",used)}
 before,err:=m.ExperienceSnapshot();if err!=nil{t.Fatal(err)}
 for _,q:=range signatures{_,_,e:=m.InferExperience(q,0);if e!=nil{t.Fatal(e)}}
 after,_:=m.ExperienceSnapshot()
 if !bytes.Equal(before,after){t.Fatal("query mutated observation memory")}
 n:=NewExperienceMemory()
 if err:=n.ExperienceRestore(before);err!=nil{t.Fatal(err)}
 newer,_:=n.ExperienceSnapshot()
 if !bytes.Equal(before,newer){t.Fatal("snapshot not exact")}
 invalid:=[][4]float64{{math.NaN(),0,0,0},{2,0,0,0}}
 for _,b:=range invalid{
  if e:=n.ObserveExperience(b,b,0,0);e==nil{t.Fatal("invalid observation accepted")}
 }
 orig,_:=n.ExperienceSnapshot()
 if !bytes.Equal(orig,before){t.Fatal("invalid observation changed memory")}
 if e:=n.ExperienceRestore([]byte("{}"));e==nil{t.Fatal("schema-less restore")}
}
func TestExperienceNoEntitySlotsAndAba(t *testing.T){
 m:=NewExperienceMemory()
 a:=[4]float64{-.8,-.9,0,0}; b:=[4]float64{.8,.9,0,0}
 for i:=0;i<30;i++{
  if e:=m.ObserveExperience(a,a,1,1);e!=nil{t.Fatal(e)}
 }
 p,known,e:=m.InferExperience(a,1);if e!=nil{t.Fatal(e)}
 if !known||p<.9{t.Fatalf("A not retained before distractor: %v %v",p,known)}
 for i:=0;i<30;i++{
  if e:=m.ObserveExperience(b,b,1,0);e!=nil{t.Fatal(e)}
 }
 p,known,e=m.InferExperience(a,1);if e!=nil{t.Fatal(e)}
 if !known||p<.85{t.Fatalf("A lost across B: %v %v",p,known)}
 fmt.Printf("NCI04_NATIVE_W_ABA_QUALITY=%.6f NCI04_NATIVE_W_RETAINED_BYTES=%d\n",p,unsafe.Sizeof(*m))
 if _,ok,e:=m.InferExperience([4]float64{0,0,0,0},1);e!=nil||ok{t.Fatalf("novel source should remain unknown: %v %v",ok,e)}
}
func TestExperienceComparatorEvidenceOnly(t *testing.T){
 // A fully supervised exact-signature count table is a strong upper bound.
 // It receives the same observed raw episodes; we do not claim model parity.
 m:=NewExperienceMemory()
 correct,covered:=0,0
 for i:=0;i<48;i++ {
   sig:=[4]float64{float64((i%3)-1)*.7,.3,0,0}
   label:=(i%3)%2
   if i>=24{
     p,ok,e:=m.InferExperience(sig,0);if e!=nil{t.Fatal(e)}
     if ok{covered++;if (p>=.5)==(label==1){correct++}}
   }
   if e:=m.ObserveExperience(sig,sig,0,float64(label));e!=nil{t.Fatal(e)}
 }
 strongOracle:=24
 fmt.Printf("NCI04_NATIVE_W_DEV_COVERED=%d NCI04_NATIVE_W_DEV_CORRECT=%d NCI04_NATIVE_W_STRONG_ORACLE=%d\n",covered,correct,strongOracle)
 report:=map[string]any{"correct":correct,"covered":covered,"oracle":strongOracle,"learned_entity_binding_accepted":false}
 if _,e:=json.Marshal(report);e!=nil{t.Fatal(e)}
}
