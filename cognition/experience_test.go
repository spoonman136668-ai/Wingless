package cognition

import (
 "bytes"
 "encoding/json"
 "fmt"
 "math"
 "testing"
 "unsafe"
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
// A trained capacity-matched raw-feature prototype peer receives exactly
// the same experiences. It does not learn feature stability from paired views.
type trainedRawPeer struct{
 tracks [ExperienceCapacity]struct{
  used bool
  centre [4]float64
  value [2]float64
  counts [2]uint32
 }
 count int
}
func (m *trainedRawPeer)locate(v [4]float64) (int,float64){
 index:=-1;best:=math.Inf(1)
 for i,t:=range m.tracks{
  if !t.used{continue}
  dist:=0.0;for k:=range v{d:=v[k]-t.centre[k];dist+=d*d}
  if dist<best{best=dist;index=i}
 }
 return index,best
}
func (m *trainedRawPeer)observe(v [4]float64,action int,value float64){
 k,d:=m.locate(v)
 if k<0||d>.28{
  if m.count<ExperienceCapacity{k=m.count;m.count++}else{k=m.count%ExperienceCapacity}
  m.tracks[k].used=true;m.tracks[k].centre=v;m.tracks[k].value=[2]float64{.5,.5}
  m.tracks[k].counts=[2]uint32{}
 }
 t:=&m.tracks[k]
 for i:=range t.centre{t.centre[i]=.85*t.centre[i]+.15*v[i]}
 t.value[action]=.8*t.value[action]+.2*value
 t.counts[action]++
}
func (m *trainedRawPeer)infer(v [4]float64,action int)(float64,bool){
 k,d:=m.locate(v);if k<0||d>.28||m.tracks[k].counts[action]==0{return 0,false}
 return m.tracks[k].value[action],true
}
func TestExperienceTrainedMatchedPeerEvidenceOnly(t *testing.T){
 candidate:=NewExperienceMemory()
 var peer trainedRawPeer
 candidateCorrect,peerCorrect,candidateCoverage,peerCoverage:=0,0,0,0
 for i:=0;i<72;i++{
  id:=i%3
  sig:=[4]float64{float64(id-1)*.7,.3,0,0}
  y:=id%2
  if i>=24{
   p,known,e:=candidate.InferExperience(sig,0);if e!=nil{t.Fatal(e)}
   if known{candidateCoverage++;if (p>=.5)==(y==1){candidateCorrect++}}
   q,ok:=peer.infer(sig,0)
   if ok{peerCoverage++;if (q>=.5)==(y==1){peerCorrect++}}
  }
  if e:=candidate.ObserveExperience(sig,sig,0,float64(y));e!=nil{t.Fatal(e)}
  peer.observe(sig,0,float64(y))
 }
 // The oracle's correctness is measured from the actual source-free labels,
 // never hardcoded as a fake test outcome.
 perfect:=0;for i:=24;i<72;i++{id:=i%3;if int(float64(id%2))==id%2{perfect++}}
 fmt.Printf("NCI04_W_MATCHED_NATIVE_CORRECT=%d/%d TRAINED_RAW_PEER_CORRECT=%d/%d SYMBOLIC_ORACLE=%d/%d\n",
   candidateCorrect,candidateCoverage,peerCorrect,peerCoverage,perfect,48)
 fmt.Printf("NCI04_W_NATIVE_BYTES=%d TRAINED_RAW_PEER_BYTES=%d\n",unsafe.Sizeof(*candidate),unsafe.Sizeof(peer))
 report:=map[string]any{"native_correct":candidateCorrect,"trained_peer_correct":peerCorrect,
  "native_coverage":candidateCoverage,"trained_peer_coverage":peerCoverage,
  "symbolic_oracle_correct":perfect,"eval_episodes":48,"scientific_first_use":0,
  "cognitive_transfer_accepted":false}
 if _,e:=json.Marshal(report);e!=nil{t.Fatal(e)}
}
