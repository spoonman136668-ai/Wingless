// Package r7 implements Wingless's independent, bounded, integer-valued
// online causal-state pilot. This is NOT a conventional LLM, GRU, or broker.
//
// This is a public-feedback sensory-causal ASSOCIATION smoke, not ITD-01
// unfamiliar task-family transfer or a full action-conditioned world model.
package r7

import (
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "fmt"
)

type Context string
const (A Context="A"; B Context="B")
type PublicObservation struct {
 Context Context `json:"context"`
 Sensors [4]bool `json:"sensors"`
}
type PublicFeedback struct {
 ObservedOutcome bool `json:"observed_outcome"`
 DelaySteps int `json:"delay_steps"`
}
type State struct {
 Version int `json:"version"`
 Weights [2][4]int16 `json:"weights"`
 Bias [2]int16 `json:"bias"`
 FeedbackEvents uint32 `json:"feedback_events"`
 Changes uint32 `json:"weight_changes"`
}
type NativePredictor struct {State State}
func New() *NativePredictor {return &NativePredictor{State:State{Version:1}}}
func index(c Context)(int,error) {switch c {case A:return 0,nil;case B:return 1,nil};return 0,errors.New("R7_CONTEXT_NOT_PUBLIC_A_OR_B")}
func sign(x bool) int32 {if x{return 1};return -1}
func clip(n int32) int16 {if n>127{return 127};if n< -127{return -127};return int16(n)}
func (p *NativePredictor) Predict(o PublicObservation)(bool,error) {
 i,e:=index(o.Context);if e!=nil{return false,e}
 score:=int32(p.State.Bias[i])
 for j,v:=range o.Sensors {score+=int32(p.State.Weights[i][j])*sign(v)}
 return score>0,nil
}
// Learn is available ONLY after externally observed, delayed public feedback.
// No hidden teacher IDs, causal graph or correct outcome are in Predict().
func (p *NativePredictor) Learn(o PublicObservation,f PublicFeedback) error {
 i,e:=index(o.Context);if e!=nil{return e}
 if f.DelaySteps<2 {return errors.New("R7_PUBLIC_FEEDBACK_NOT_YET_AVAILABLE")}
 guess,_:=p.Predict(o)
 p.State.FeedbackEvents++
 if guess==f.ObservedOutcome{return nil}
 delta:=int32(1);if !f.ObservedOutcome{delta=-1}
 for j,v:=range o.Sensors {p.State.Weights[i][j]=clip(int32(p.State.Weights[i][j])+delta*sign(v))}
 p.State.Bias[i]=clip(int32(p.State.Bias[i])+delta)
 p.State.Changes++
 return nil
}
type persisted struct { Magic string `json:"magic"`; Payload State `json:"payload"`; Digest string `json:"sha256"` }
func (p *NativePredictor) Save()([]byte,error){
 raw,e:=json.Marshal(p.State);if e!=nil{return nil,e}
 dig:=sha256.Sum256(raw)
 return json.Marshal(persisted{Magic:"wingless.r7.native.v1",Payload:p.State,Digest:hex.EncodeToString(dig[:])})
}
func Restore(raw []byte)(*NativePredictor,error){
 if len(raw)==0||len(raw)>32768{return nil,errors.New("R7_PERSISTED_CAP_OR_EMPTY")}
 var o persisted
 if e:=json.Unmarshal(raw,&o);e!=nil{return nil,e}
 if o.Magic!="wingless.r7.native.v1"||o.Payload.Version!=1{return nil,errors.New("R7_SOURCE_STATE_IDENTITY_INVALID")}
 check,e:=json.Marshal(o.Payload);if e!=nil{return nil,e}
 dig:=sha256.Sum256(check)
 if o.Digest!=hex.EncodeToString(dig[:]){return nil,errors.New("R7_PERSISTED_SHA_INVALID")}
 for _,row:=range o.Payload.Weights{for _,w:=range row{if w < -127 || w>127{return nil,errors.New("R7_WEIGHT_BOUNDS_INVALID")}}}
 return &NativePredictor{State:o.Payload},nil
}
// The generator is PRIVATE DEV infrastructure. Its hidden mechanism must
// never be passed to a NativePredictor instance.
type devSource struct {x uint64}
func source(seed int) *devSource {return &devSource{x:uint64(seed)*0x9e3779b97f4a7c15+0xd1b54a32d192ed03}}
func (d *devSource) next() uint64 {d.x ^= d.x<<13;d.x ^= d.x>>7;d.x ^= d.x<<17;return d.x}
func (d *devSource) nextExample(c Context)(PublicObservation,bool){
 o:=PublicObservation{Context:c}
 for j:=range o.Sensors{o.Sensors[j]=(d.next()&1)==1}
 y:=o.Sensors[2];if c==B{y=!o.Sensors[1]}
 if d.next()%10==0 {y=!y}
 return o,y
}
type Scores struct {
 Seed int `json:"seed"`
 PreA float64 `json:"a_before_b"`
 RetainedA float64 `json:"a_after_b"`
 B float64 `json:"b_after_training"`
 NoStateA float64 `json:"sham_a"`
 NoStateB float64 `json:"sham_b"`
 MajorityA float64 `json:"constant_a"`
 MajorityB float64 `json:"constant_b"`
 Updates uint32 `json:"delayed_public_feedback_updates"`
 Changes uint32 `json:"nonzero_weight_changes"`
 SerializedBytes int `json:"serialized_state_bytes"`
 ColdRestoreExact bool `json:"cold_restore_exact"`
}
type example struct{o PublicObservation; y bool}
func take(d *devSource,c Context,n int)[]example{
 xs:=make([]example,n)
 for i:=range xs{xs[i].o,xs[i].y=d.nextExample(c)}
 return xs
}
func accuracy(p *NativePredictor,xs []example)(float64,error){
 correct:=0
 for _,x:=range xs{g,e:=p.Predict(x.o);if e!=nil{return 0,e};if g==x.y{correct++}}
 return float64(correct)/float64(len(xs)),nil
}
func baseline(xs []example)float64{n:=0;for _,x:=range xs{if !x.y{n++}};return float64(n)/float64(len(xs))}
func training(p *NativePredictor,xs []example)error{
 for _,x:=range xs{
  if _,e:=p.Predict(x.o);e!=nil{return e}
  if e:=p.Learn(x.o,PublicFeedback{ObservedOutcome:x.y,DelaySteps:2});e!=nil{return e}
 }
 return nil
}
// Study pins the public cohort order before first observations. All score
// probes are held for evaluation only; feedback is never used to fit on eval.
func Study(seed int)(Scores,error){
 d:=source(seed+71)
 trainA:=take(d,A,320)
 valA:=take(d,A,160)
 trainB:=take(d,B,320)
 postA:=take(d,A,160)
 postB:=take(d,B,160)
 p:=New()
 if e:=training(p,trainA);e!=nil{return Scores{},e}
 pre,e:=accuracy(p,valA);if e!=nil{return Scores{},e}
 if e:=training(p,trainB);e!=nil{return Scores{},e}
 bytes,e:=p.Save();if e!=nil{return Scores{},e}
 saved,e:=Restore(bytes);if e!=nil{return Scores{},e}
 a,e:=accuracy(saved,postA);if e!=nil{return Scores{},e}
 b,e:=accuracy(saved,postB);if e!=nil{return Scores{},e}
 sham:=New()
 a0,e:=accuracy(sham,postA);if e!=nil{return Scores{},e}
 b0,e:=accuracy(sham,postB);if e!=nil{return Scores{},e}
 stable:=true
 for _,x:=range append(append([]example{},postA...),postB...){
  q,_:=p.Predict(x.o);r,_:=saved.Predict(x.o)
  if q!=r{stable=false;break}
 }
 if len(bytes)>32768{return Scores{},fmt.Errorf("R7_MEMORY_CAP_EXCEEDED_%d",len(bytes))}
 return Scores{Seed:seed,PreA:pre,RetainedA:a,B:b,NoStateA:a0,NoStateB:b0,
  MajorityA:baseline(postA),MajorityB:baseline(postB),
  Updates:p.State.FeedbackEvents,Changes:p.State.Changes,
  SerializedBytes:len(bytes),ColdRestoreExact:stable},nil
}
