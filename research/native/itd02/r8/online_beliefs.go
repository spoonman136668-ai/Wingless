// Package r8 is a native Go learning layer in the original Wingless R5 lineage.
// It only learns source-public selected-action outcomes accepted by the existing
// R5 FIT and Choose/Feedback gates. No simulator, hidden grid or protected role.
package r8

import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "io"

 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
 r5 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r5"
)

const Schema = "wingless.itd02.r8.online-public-next-front-risk.v1"
const MaxEvidence = 65535
const MinTurnSamples = 4

// Bin is a Beta(1,1) posterior's observed binary sufficient statistics.
type Bin struct {
 Clear uint32 `json:"clear"`
 Blocked uint32 `json:"blocked"`
}
func (v Bin) count() uint64 {return uint64(v.Clear)+uint64(v.Blocked)}
func (v Bin) risk() float64 {
 return float64(v.Blocked+1)/float64(v.count()+2)
}

// Policy keeps ALL original R5 trained weights and episode routing in Core.
// Learned future-front uncertainty is a distinct opt-in native state.
type Policy struct {
 Core *r5.Policy
 Bins [8][2]Bin
 FitSeen uint64
 LiveSeen uint64
 EpisodeSteps uint32
 TurnSeen uint64
 FitPrevious uint8 // 0 none; 1 LEFT, 2 RIGHT, 3 FORWARD.
 Previous uint8
 Armed bool
 ArmedContext uint8
}
func New() *Policy {return &Policy{Core:r5.New()}}
func validPublic(s r1.PublicGrid)bool{
 return s.X>=0&&s.X<=64&&s.Y>=0&&s.Y<=64&&
   s.GoalX>=0&&s.GoalX<=64&&s.GoalY>=0&&s.GoalY<=64&&
   s.Direction>=0&&s.Direction<=3
}
func ctx(blocked bool,previous uint8)uint8{
 c:=previous
 if blocked{c+=4}
 return c
}
func(p *Policy) stateValid()bool{
 if p==nil||p.Core==nil||p.Core.Core==nil||p.Core.Core.Core==nil||
   p.FitPrevious>3||p.Previous>3||p.ArmedContext>7||
   p.FitSeen+p.LiveSeen>MaxEvidence||
   p.TurnSeen>p.FitSeen+p.LiveSeen||
   p.FitSeen!=p.Core.Core.Core.Events||
   p.EpisodeSteps!=p.Core.Core.Route.Steps||
   p.Armed!=p.Core.Core.Route.Pending{return false}
 var turns uint64
 for _,pair:=range p.Bins {
  for _,v:=range pair{turns+=v.count()}
 }
 return turns==p.TurnSeen
}
func(p *Policy) Healthy()bool{
 if !p.stateValid(){return false}
 _,e:=p.Core.Snapshot()
 return e==nil
}
func(p *Policy) canAdd(context uint8,action int)bool{
 if p.FitSeen+p.LiveSeen>=MaxEvidence{return false}
 if action==0||action==1{
  return p.Bins[context][action].count()<MaxEvidence
 }
 return action==2
}
func(p *Policy) record(context uint8,action int,returned r1.PublicGrid){
 if action==0||action==1{
  v:=&p.Bins[context][action]
  if returned.FrontBlocked{v.Blocked++}else{v.Clear++}
  p.TurnSeen++
 }
}

// ObserveFit has no separate admission path. The ORIGINAL Wingless R5/R3
// source-FIT policy must first accept the actual selected public transition.
func(p *Policy) ObserveFit(before r1.PublicGrid,action int,after r1.PublicGrid,
 reward float64,terminal bool)error{
 if !p.stateValid()||!validPublic(before)||!validPublic(after)||
  action<0||action>2||p.Armed||p.EpisodeSteps!=0{
  return errors.New("R8_FIT_SOURCE_CUSTODY")
 }
 context:=ctx(before.FrontBlocked,p.FitPrevious)
 if !p.canAdd(context,action){return errors.New("R8_FIT_EVIDENCE_BUDGET")}
 if e:=p.Core.ObserveFit(before,action,after,reward,terminal);e!=nil{return e}
 p.record(context,action,after)
 p.FitSeen++
 if terminal{p.FitPrevious=0}else{p.FitPrevious=uint8(action+1)}
 return nil
}

func(p *Policy) BeginEpisode()error{
 if !p.stateValid()||p.Armed{return errors.New("R8_INVALID_EPISODE_RESET")}
 if e:=p.Core.BeginEpisode();e!=nil{return e}
 p.Previous=0
 p.EpisodeSteps=0
 return nil
}

// EligiblePrediction is an actually learned next-front risk from selected
// public turn effects, never an imagined unobserved source consequence.
func(p *Policy) EligiblePrediction(s r1.PublicGrid)([2]float64,bool,error){
 if !p.Healthy()||!validPublic(s){return [2]float64{},false,errors.New("R8_INVALID_PUBLIC_PREDICTION")}
 k:=ctx(s.FrontBlocked,p.Previous)
 left,right:=p.Bins[k][0],p.Bins[k][1]
 risks:=[2]float64{left.risk(),right.risk()}
 return risks,left.count()>=MinTurnSamples&&right.count()>=MinTurnSamples,nil
}
func(p *Policy) Choose(s r1.PublicGrid)(int,error){
 if !p.Healthy()||!validPublic(s)||p.Armed{
  return 0,errors.New("R8_INVALID_PUBLIC_ACTION_REQUEST")
 }
 context:=ctx(s.FrontBlocked,p.Previous)
 chosen:=-1
 if s.FrontBlocked {
  left,right:=p.Bins[context][0],p.Bins[context][1]
  if left.count()>=MinTurnSamples&&right.count()>=MinTurnSamples{
   diff:=left.risk()-right.risk()
   if diff>=.25{chosen=1}
   if diff<=-.25{chosen=0}
  }
 }
 if chosen<0{
  a,e:=p.Core.Choose(s)
  if e!=nil{return 0,e}
  chosen=a
 }else{
  // Preserve the original R5 checkpoint and pending selected-action custodian.
  if _,e:=p.Core.Snapshot();e!=nil{return 0,e}
  p.Core.Core.Route.Pending=true
  p.Core.Core.Route.Before=s
  p.Core.Core.Route.Action=chosen
 }
 p.Armed=true
 p.ArmedContext=context
 return chosen,nil
}
func(p *Policy) Feedback(before r1.PublicGrid,action int,after r1.PublicGrid)error{
 if !p.stateValid()||!p.Armed||!validPublic(before)||!validPublic(after)||
  action<0||action>2||p.EpisodeSteps>=4096||
  p.Core.Core.Route.Before!=before||
  p.Core.Core.Route.Action!=action{
  return errors.New("R8_REJECT_UNSELECTED_SOURCE_FEEDBACK")
 }
 if !p.canAdd(p.ArmedContext,action){return errors.New("R8_SOURCE_EVIDENCE_CAPACITY")}
 if e:=p.Core.Feedback(before,action,after);e!=nil{return e}
 p.record(p.ArmedContext,action,after)
 p.LiveSeen++
 p.EpisodeSteps++
 p.Previous=uint8(action+1)
 p.Armed=false
 p.ArmedContext=0
 return nil
}

type state struct{
 OriginalR5 []byte `json:"original_r5_checkpoint"`
 Bins [8][2]Bin `json:"source_public_risk_counts"`
 FitSeen uint64 `json:"original_fit_accepted"`
 LiveSeen uint64 `json:"live_selected_feedback_accepted"`
 EpisodeSteps uint32 `json:"episode_selected_steps"`
 TurnSeen uint64 `json:"turn_consequence_samples"`
 FitPrevious uint8 `json:"fit_previous_action"`
 Previous uint8 `json:"episode_previous_action"`
 Armed bool `json:"armed_selected_action"`
 ArmedContext uint8 `json:"armed_public_context"`
}
type envelope struct{
 Schema string `json:"schema"`
 State state `json:"state"`
 SHA256 string `json:"sha256"`
}
func(p *Policy) Snapshot()([]byte,error){
 if !p.stateValid(){return nil,errors.New("R8_INVALID_NATIVE_CHECKPOINT")}
 raw,e:=p.Core.Snapshot();if e!=nil{return nil,e}
 v:=state{OriginalR5:raw,Bins:p.Bins,FitSeen:p.FitSeen,
  LiveSeen:p.LiveSeen,EpisodeSteps:p.EpisodeSteps,TurnSeen:p.TurnSeen,
  FitPrevious:p.FitPrevious,Previous:p.Previous,
  Armed:p.Armed,ArmedContext:p.ArmedContext}
 body,e:=json.Marshal(v);if e!=nil{return nil,e}
 h:=sha256.Sum256(body)
 return json.Marshal(envelope{Schema:Schema,State:v,SHA256:hex.EncodeToString(h[:])})
}
func Restore(raw []byte)(*Policy,error){
 if len(raw)==0||len(raw)>3<<20{return nil,errors.New("R8_CHECKPOINT_SIZE")}
 dec:=json.NewDecoder(bytes.NewReader(raw));dec.DisallowUnknownFields()
 var wire envelope
 if e:=dec.Decode(&wire);e!=nil{return nil,e}
 if e:=dec.Decode(new(struct{}));e!=io.EOF{return nil,errors.New("R8_CHECKPOINT_TRAILING_BYTES")}
 canonical,e:=json.Marshal(wire);if e!=nil||!bytes.Equal(canonical,raw){
  return nil,errors.New("R8_CHECKPOINT_NONCANONICAL")
 }
 if wire.Schema!=Schema{return nil,errors.New("R8_CHECKPOINT_SCHEMA")}
 body,e:=json.Marshal(wire.State);if e!=nil{return nil,e}
 h:=sha256.Sum256(body)
 if wire.SHA256!=hex.EncodeToString(h[:]){return nil,errors.New("R8_CHECKPOINT_DIGEST")}
 core,e:=r5.Restore(wire.State.OriginalR5);if e!=nil{return nil,e}
 p:=&Policy{Core:core,Bins:wire.State.Bins,FitSeen:wire.State.FitSeen,
  LiveSeen:wire.State.LiveSeen,EpisodeSteps:wire.State.EpisodeSteps,
  TurnSeen:wire.State.TurnSeen,FitPrevious:wire.State.FitPrevious,
  Previous:wire.State.Previous,Armed:wire.State.Armed,
  ArmedContext:wire.State.ArmedContext}
 if !p.Healthy(){return nil,errors.New("R8_CHECKPOINT_SOURCE_CONSERVATION")}
 return p,nil
}
