// Package dk10 adds opt-in episodic public transition evidence to unchanged
// DK8U/DK7 native Go learning. No hidden MiniGrid state, route oracle or
// hand-coded DoorKey action sequence. Fit and live updates are selected-only.
package dk10

import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "math"

 dk7 "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk7"
 dk8 "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk8"
)

const Schema="wingless.itd02.dk10.episodic-public-transitions.v1"
const MaxStates=512
const MaxObservations=4096
const MaxSnapshotBytes=550000

type Actions [7]map[string]uint32

type Agent struct {
 Base *dk8.Explorer
 History map[string]Actions
 Visits map[string]uint32
 Samples uint64
}

func New()*Agent{
 return &Agent{Base:dk8.New(),History:map[string]Actions{},Visits:map[string]uint32{}}
}
func publicHash(p dk7.Public)(string,error){
 if e:=p.Valid();e!=nil{return "",e}
 b,e:=json.Marshal(p);if e!=nil{return "",e}
 s:=sha256.Sum256(b);return hex.EncodeToString(s[:]),nil
}
func hash(raw []byte)string{s:=sha256.Sum256(raw);return hex.EncodeToString(s[:])}
func validDigest(value string)bool{
 if len(value)!=64{return false}
 for _,b:=range value{if !(b>='0'&&b<='9'||b>='a'&&b<='f'){return false}}
 return true
}
func (m *Agent) valid()bool{
 if m==nil||m.Base==nil||m.Base.Native==nil||m.History==nil||m.Visits==nil||
   m.Samples>MaxObservations||m.Base.Native.Samples!=m.Samples||
   len(m.History)>MaxStates||len(m.Visits)>MaxObservations+1{return false}
 total:=uint64(0)
 for state,arms:=range m.History{
  if !validDigest(state){return false}
  for _,dests:=range arms{
   for dest,count:=range dests{
    if !validDigest(dest)||count==0||count>MaxObservations{return false}
    total+=uint64(count)
   }
  }
 }
 if total!=m.Samples{return false}
 for state,n:=range m.Visits{if !validDigest(state)||n==0||n>MaxObservations+1{return false}}
 return true
}

func (m *Agent) BeginEpisode()error{
 if !m.valid(){return errors.New("DK10_INVALID_EPISODIC_STATE")}
 if e:=m.Base.BeginEpisode();e!=nil{return e}
 m.Visits=map[string]uint32{}
 return nil
}

func (m *Agent) Choose(before dk7.Public)(*int,string,error){
 if !m.valid()||m.Base.Native.Pending!=nil{return nil,"",errors.New("DK10_INVALID_OR_PENDING")}
 current,e:=publicHash(before);if e!=nil{return nil,"",e}
 baseline,status,e:=m.Base.Choose(before);if e!=nil{return nil,"",e}
 if status=="SELECT_NATIVE_OBSERVED_POSITIVE_REWARD" && baseline!=nil{
  return baseline,"SELECT_NATIVE_OBSERVED_POSITIVE_REWARD",nil
 }
 row:=m.History[current]
 bestAction:=-1
 bestScore:=math.Inf(-1)
 for action,dests:=range row{
  total:=uint32(0)
  changed:=float64(0)
  visit:=float64(0)
  for dest,n:=range dests{
   total+=n
   if dest!=current{changed+=float64(n)}
   visit+=float64(n)*float64(m.Visits[dest])
  }
  if total==0||changed==0{continue}
  denom:=float64(total)
  score:=changed/denom-(visit/denom)/float64(1+len(m.Visits))
  if score>bestScore || score==bestScore && (bestAction<0||action<bestAction){
   bestScore=score;bestAction=action
  }
 }
 if bestAction>=0{
  return &bestAction,"SELECT_NATIVE_LEARNED_PUBLIC_TRANSITION_NOVELTY",nil
 }
 if baseline==nil{return nil,"",errors.New("DK10_BASE_NATIVE_ACTION_MISSING")}
 return baseline,"SELECT_NATIVE_BASE_EXPLORATION",nil
}

func (m *Agent) Commit(before dk7.Public,action int)error{
 if !m.valid(){return errors.New("DK10_INVALID_EPISODIC_STATE")}
 return m.Base.Commit(before,action)
}
func (m *Agent) check(before dk7.Public,action int,after dk7.Public)(string,string,error){
 if !m.valid()||m.Samples>=MaxObservations||action<0||action>6{
  return "","",errors.New("DK10_SELECTED_SAMPLE_LIMIT")
 }
 b,e:=publicHash(before);if e!=nil{return "","",e}
 a,e:=publicHash(after);if e!=nil{return "","",e}
 if _,ok:=m.History[b];!ok&&len(m.History)>=MaxStates{
  return "","",errors.New("DK10_SOURCE_PUBLIC_CONTEXT_CAP")
 }
 return b,a,nil
}
func (m *Agent) record(before,after string,action int){
 row:=m.History[before]
 if row[action]==nil{row[action]=map[string]uint32{}}
 row[action][after]++
 m.History[before]=row
 m.Samples++
 m.Visits[after]++
}
func (m *Agent) Feedback(before dk7.Public,action int,after dk7.Public,reward float64,terminal bool)error{
 b,a,e:=m.check(before,action,after);if e!=nil{return e}
 if e=m.Base.Feedback(before,action,after,reward,terminal);e!=nil{return e}
 m.record(b,a,action)
 if !m.valid(){return errors.New("DK10_POST_SELECTED_INTEGRITY")}
 return nil
}
func (m *Agent) ObserveFit(before dk7.Public,action int,after dk7.Public,reward float64,terminal bool)error{
 b,a,e:=m.check(before,action,after);if e!=nil{return e}
 if e=m.Base.ObserveFit(before,action,after,reward,terminal);e!=nil{return e}
 m.record(b,a,action)
 if !m.valid(){return errors.New("DK10_POST_FIT_INTEGRITY")}
 return nil
}

type state struct {
 Schema string `json:"schema"`
 Base json.RawMessage `json:"base"`
 History map[string]Actions `json:"history"`
 Visits map[string]uint32 `json:"visits"`
 Samples uint64 `json:"samples"`
}
type envelope struct{State state `json:"state"`;SHA256 string `json:"sha256"`}
func (m *Agent) Snapshot()([]byte,error){
 if !m.valid(){return nil,errors.New("DK10_UNHEALTHY_SNAPSHOT")}
 base,e:=m.Base.Snapshot();if e!=nil{return nil,e}
 s:=state{Schema:Schema,Base:base,History:m.History,Visits:m.Visits,Samples:m.Samples}
 raw,e:=json.Marshal(s);if e!=nil{return nil,e}
 data,e:=json.Marshal(envelope{State:s,SHA256:hash(raw)})
 if e!=nil||len(data)>MaxSnapshotBytes{return nil,errors.New("DK10_SNAPSHOT_CAP")}
 return data,nil
}
func Restore(data []byte)(*Agent,error){
 if len(data)==0||len(data)>MaxSnapshotBytes{return nil,errors.New("DK10_RESTORE_SIZE")}
 var wire envelope
 d:=json.NewDecoder(bytes.NewReader(data));d.DisallowUnknownFields()
 if e:=d.Decode(&wire);e!=nil{return nil,errors.New("DK10_RESTORE_DECODE")}
 if d.Decode(new(any))==nil{return nil,errors.New("DK10_RESTORE_TRAILING")}
 if wire.State.Schema!=Schema{return nil,errors.New("DK10_RESTORE_SCHEMA")}
 raw,e:=json.Marshal(wire.State);if e!=nil||hash(raw)!=wire.SHA256{
  return nil,errors.New("DK10_RESTORE_AUTH")
 }
 base,e:=dk8.Restore(wire.State.Base);if e!=nil{return nil,e}
 m:=&Agent{Base:base,History:wire.State.History,Visits:wire.State.Visits,Samples:wire.State.Samples}
 if !m.valid(){return nil,errors.New("DK10_RESTORE_SELECTED_COUNTS")}
 canonical,e:=m.Snapshot()
 if e!=nil||!bytes.Equal(canonical,data){return nil,errors.New("DK10_NONCANONICAL")}
 return m,nil
}
