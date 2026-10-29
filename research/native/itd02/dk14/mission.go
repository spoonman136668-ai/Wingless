// Package dk14 accepts real original missions at the native Go source boundary.
// Unchanged DK13 requires the DoorKey mission internally; that compatibility
// projection is explicit and never made by the environment/evaluator.
package dk14
import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "io"
 prior "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk13"
)
const Schema="wingless.itd02.dk14.actual-source-mission-conditioned.v1"
const DoorKey="use the key to open the door and then get to the goal"
const GoalOnly="get to the green goal square"
const MaxCheckpoint=850000
type Public=prior.Public
type Bin struct{Count uint32 `json:"count"`;Success uint32 `json:"success"`}
type Stats [7]Bin
type pending struct{Before string `json:"before"`;Mission string `json:"mission"`;Action int `json:"action"`}
type Policy struct{
 *prior.Policy
 Missions map[string]Stats
 Pending *pending
}
func New()*Policy{return &Policy{Policy:prior.New(),Missions:map[string]Stats{}}}
func hash(b []byte)string{v:=sha256.Sum256(b);return hex.EncodeToString(v[:])}
func missionOK(m string)bool{return m==DoorKey||m==GoalOnly}
func projected(p Public)(Public,error){
 if !missionOK(p.Mission){return p,errors.New("DK14_UNTRUSTED_OR_UNREGISTERED_ORIGINAL_MISSION")}
 p.Mission=DoorKey
 if e:=p.Valid();e!=nil{return p,e}
 return p,nil
}
func publicHash(p Public)(string,error){
 if _,e:=projected(p);e!=nil{return "",e}
 b,e:=json.Marshal(p);if e!=nil{return "",e}
 return hash(b),nil
}
func (m *Policy) valid()bool{
 if m==nil||m.Policy==nil||m.Missions==nil||len(m.Missions)>2{return false}
 var count uint64
 for mission,stats:=range m.Missions{
  if !missionOK(mission){return false}
  for _,b:=range stats{if b.Success>b.Count{return false};count+=uint64(b.Count)}
 }
 return count==m.Policy.Agent.Samples
}
func (m *Policy) BeginEpisode()error{
 if !m.valid()||m.Pending!=nil{return errors.New("DK14_EPISODE_PENDING_OR_DRIFT")}
 return m.Policy.BeginEpisode()
}
func (m *Policy) Choose(before Public)(*int,string,error){
 if !m.valid()||m.Pending!=nil{return nil,"",errors.New("DK14_SOURCE_PENDING_OR_INVALID")}
 adapted,e:=projected(before);if e!=nil{return nil,"",e}
 base,status,e:=m.Policy.Choose(adapted)
 if e!=nil{return nil,"",e}
 // Selection uses ONLY previously observed genuine source-positive outcomes
 // for this original mission, never a hard-coded task-action mapping.
 best:=-1;highest:=int64(0)
 for a,b:=range m.Missions[before.Mission]{
  if b.Success==0{continue}
  confidence:=int64(4096)*int64(b.Success)/(int64(b.Count)+64)
  if confidence>highest{highest=confidence;best=a}
 }
 if best>=0{return &best,"SELECT_NATIVE_GENUINE_MISSION_SELECTED_REWARD",nil}
 return base,status,nil
}
func (m *Policy) Commit(before Public,action int)error{
 if !m.valid()||m.Pending!=nil||action<0||action>6{return errors.New("DK14_INVALID_SOURCE_COMMIT")}
 adapted,e:=projected(before);if e!=nil{return e}
 h,e:=publicHash(before);if e!=nil{return e}
 if e=m.Policy.Commit(adapted,action);e!=nil{return e}
 m.Pending=&pending{Before:h,Mission:before.Mission,Action:action}
 return nil
}
func (m *Policy) learn(before Public,action int,after Public,reward float64,terminal,live bool)error{
 if !m.valid()||action<0||action>6||before.Mission!=after.Mission{
  return errors.New("DK14_MISSION_CHANGED_OR_SELECTED_SOURCE_INVALID")
 }
 b,e:=projected(before);if e!=nil{return e}
 a,e:=projected(after);if e!=nil{return e}
 if live{
  h,e:=publicHash(before);if e!=nil{return e}
  if m.Pending==nil||m.Pending.Before!=h||m.Pending.Action!=action||
     m.Pending.Mission!=before.Mission{return errors.New("DK14_SELECTED_FEEDBACK_MISMATCH")}
 }else if m.Pending!=nil{return errors.New("DK14_FIT_WITH_PENDING_SOURCE")}
 if live{e=m.Policy.Feedback(b,action,a,reward,terminal)
 }else{e=m.Policy.ObserveFit(b,action,a,reward,terminal)}
 if e!=nil{return e}
 stat:=m.Missions[before.Mission];one:=stat[action];one.Count++
 if reward>0{one.Success++};stat[action]=one;m.Missions[before.Mission]=stat
 if live{m.Pending=nil}
 if !m.valid(){return errors.New("DK14_NATIVE_MISSION_COUNTS_DRIFT")}
 return nil
}
func (m *Policy) ObserveFit(before Public,action int,after Public,reward float64,terminal bool)error{
 return m.learn(before,action,after,reward,terminal,false)
}
func (m *Policy) Feedback(before Public,action int,after Public,reward float64,terminal bool)error{
 return m.learn(before,action,after,reward,terminal,true)
}
type state struct{
 Schema string `json:"schema"`
 Parent json.RawMessage `json:"parent"`
 Missions map[string]Stats `json:"missions"`
}
type envelope struct{State state `json:"state"`;SHA256 string `json:"sha256"`}
func (m *Policy) Snapshot()([]byte,error){
 if !m.valid()||m.Pending!=nil{return nil,errors.New("DK14_PENDING_OR_DRIFT")}
 raw,e:=m.Policy.Snapshot();if e!=nil{return nil,e}
 state:=state{Schema:Schema,Parent:raw,Missions:m.Missions}
 b,e:=json.Marshal(state);if e!=nil{return nil,e}
 out,e:=json.Marshal(envelope{State:state,SHA256:hash(b)})
 if e!=nil||len(out)>MaxCheckpoint{return nil,errors.New("DK14_SNAPSHOT_BUDGET")}
 return out,nil
}
func Restore(raw []byte)(*Policy,error){
 if len(raw)==0||len(raw)>MaxCheckpoint{return nil,errors.New("DK14_RESTORE_BUDGET")}
 d:=json.NewDecoder(bytes.NewReader(raw));d.DisallowUnknownFields()
 var wire envelope
 if e:=d.Decode(&wire);e!=nil{return nil,e}
 if e:=d.Decode(new(any));e!=io.EOF{return nil,errors.New("DK14_RESTORE_TRAILING")}
 b,e:=json.Marshal(wire.State)
 if e!=nil||wire.State.Schema!=Schema||hash(b)!=wire.SHA256{return nil,errors.New("DK14_RESTORE_SHA")}
 parent,e:=prior.Restore(wire.State.Parent);if e!=nil{return nil,e}
 m:=&Policy{Policy:parent,Missions:wire.State.Missions}
 if !m.valid(){return nil,errors.New("DK14_MISSION_SOURCE_COUNTERS")}
 again,e:=m.Snapshot()
 if e!=nil||!bytes.Equal(again,raw){return nil,errors.New("DK14_NONCANONICAL_MISSION_CHECKPOINT")}
 return m,nil
}
