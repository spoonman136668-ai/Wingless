// Package dk11 tests generic visible-object continuity learned from selected
// source-public consequences. Original DK10/DK9/DK8 and R8 remain unchanged.
package dk11

import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "fmt"
 "io"

 prior "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk10"
)

const Schema="wingless.itd02.dk11.public-object-continuity.v1"
const MaxSamples=4096
const MaxContexts=128
const MaxCheckpoint=600000

type Public=prior.Public
type Stat struct{
 Count uint32 `json:"count"`
 Changed uint32 `json:"changed"`
 RetainedChanged uint32 `json:"retained_changed"`
 Lost uint32 `json:"lost"`
 Gained uint32 `json:"gained"`
}
type Contexts [7]Stat
type Agent struct {
 *prior.Policy
 Memory map[string]Contexts
 Samples uint64
}
func New()*Agent{return &Agent{Policy:prior.New(),Memory:map[string]Contexts{}}}
func hash(b []byte)string{d:=sha256.Sum256(b);return hex.EncodeToString(d[:])}
func holding(p Public)bool{
 kind:=p.Image[3][6][0]
 return kind==5||kind==6||kind==7
}
func feature(p Public)(string,error){
 if e:=p.Valid();e!=nil{return "",e}
 hold:=0;if holding(p){hold=1}
 return fmt.Sprintf("%d:%d",p.Image[3][5][0],hold),nil
}
func (m *Agent) healthy()bool{
 if m==nil||m.Policy==nil||m.Memory==nil||m.Samples>MaxSamples||
  m.Policy.Samples!=m.Samples||len(m.Memory)>MaxContexts{return false}
 total:=uint64(0)
 for _,bins:=range m.Memory{
  for _,bin:=range bins{
   if bin.Changed>bin.Count||bin.RetainedChanged>bin.Changed||
      bin.Lost>bin.Count||bin.Gained>bin.Count{return false}
   total+=uint64(bin.Count)
  }
 }
 return total==m.Samples
}
func (m *Agent) BeginEpisode()error{
 if !m.healthy(){return errors.New("DK11_INVALID_NATIVE_STATE")}
 return m.Policy.BeginEpisode()
}
// Choose never accepts a manual DoorKey stage objective. All inputs are
// authentic original source-public image/direction/mission and OWN evidence.
func (m *Agent) Choose(before Public)(*int,string,error){
 if !m.healthy(){return nil,"",errors.New("DK11_UNHEALTHY_NATIVE")}
 ctx,e:=feature(before);if e!=nil{return nil,"",e}
 // Parent pending-state guard, without invoking the old objective-taking chooser.
 if _,e=m.Policy.Snapshot();e!=nil{return nil,"",e}
 bins:=m.Memory[ctx]
 best:=-1;bestScore:=0
 for a,b:=range bins{
  if b.Count==0{continue}
  score:=int64(0)
  if holding(before){
   score=6*int64(b.RetainedChanged)+2*int64(b.Changed)-8*int64(b.Lost)
  }else{
   score=6*int64(b.Gained)+2*int64(b.Changed)
  }
  normalized:=int(score*1024/int64(b.Count))
  if normalized>bestScore{bestScore=normalized;best=a}
 }
 if best>=0{return &best,"SELECT_NATIVE_LEARNED_PUBLIC_OBJECT_CONTINUITY",nil}
 // Genuine native zero-evidence exploratory action, never an external collector.
 least:=uint32(MaxSamples+1)
 for _,b:=range bins{if b.Count<least{least=b.Count}}
 raw,_:=json.Marshal(before);offset:=int(sha256.Sum256(raw)[0])%7
 for i:=0;i<7;i++{
  a:=(offset+i)%7
  if bins[a].Count==least{return &a,"SELECT_NATIVE_UNKNOWN_PUBLIC_ACTION",nil}
 }
 return nil,"",errors.New("DK11_ACTION_NOT_SELECTED")
}
func (m *Agent) check(before Public,action int,after Public)(string,error){
 if !m.healthy()||action<0||action>6||m.Samples>=MaxSamples{return "",errors.New("DK11_SOURCE_SAMPLE_BOUND")}
 key,e:=feature(before);if e!=nil{return "",e}
 if e=after.Valid();e!=nil{return "",e}
 if _,exists:=m.Memory[key];!exists&&len(m.Memory)>=MaxContexts{
  return "",errors.New("DK11_CONTEXT_BOUND")
 }
 return key,nil
}
func (m *Agent) record(key string,before Public,action int,after Public){
 bins:=m.Memory[key]
 s:=bins[action];s.Count++
 raw0,_:=json.Marshal(before);raw1,_:=json.Marshal(after)
 changed:=!bytes.Equal(raw0,raw1)
 h0,h1:=holding(before),holding(after)
 if changed{s.Changed++}
 if changed&&h0&&h1&&bytes.Equal([]byte(fmt.Sprint(before.Image[3][6])),[]byte(fmt.Sprint(after.Image[3][6]))){s.RetainedChanged++}
 if h0&&!h1{s.Lost++}
 if !h0&&h1{s.Gained++}
 bins[action]=s;m.Memory[key]=bins;m.Samples++
}
func (m *Agent) ObserveFit(before Public,action int,after Public,reward float64,terminal bool)error{
 key,e:=m.check(before,action,after);if e!=nil{return e}
 if e=m.Policy.ObserveFit(before,action,after,reward,terminal);e!=nil{return e}
 m.record(key,before,action,after)
 if !m.healthy(){return errors.New("DK11_FIT_COUNTERS_DIFFER")}
 return nil
}
func (m *Agent) Feedback(before Public,action int,after Public,reward float64,terminal bool)error{
 key,e:=m.check(before,action,after);if e!=nil{return e}
 if e=m.Policy.Feedback(before,action,after,reward,terminal);e!=nil{return e}
 m.record(key,before,action,after)
 if !m.healthy(){return errors.New("DK11_SELECTED_COUNTERS_DIFFER")}
 return nil
}
type snapshot struct{
 Schema string `json:"schema"`
 Parent json.RawMessage `json:"parent"`
 Memory map[string]Contexts `json:"memory"`
 Samples uint64 `json:"samples"`
}
type envelope struct{State snapshot `json:"state"`;SHA256 string `json:"sha256"`}
func (m *Agent) Snapshot()([]byte,error){
 if !m.healthy(){return nil,errors.New("DK11_INVALID_CHECKPOINT_STATE")}
 old,e:=m.Policy.Snapshot();if e!=nil{return nil,e}
 state:=snapshot{Schema:Schema,Parent:old,Memory:m.Memory,Samples:m.Samples}
 b,e:=json.Marshal(state);if e!=nil{return nil,e}
 raw,e:=json.Marshal(envelope{State:state,SHA256:hash(b)})
 if e!=nil||len(raw)>MaxCheckpoint{return nil,errors.New("DK11_CHECKPOINT_BUDGET")}
 return raw,nil
}
func Restore(raw []byte)(*Agent,error){
 if len(raw)==0||len(raw)>MaxCheckpoint{return nil,errors.New("DK11_CHECKPOINT_SIZE")}
 dec:=json.NewDecoder(bytes.NewReader(raw));dec.DisallowUnknownFields()
 var env envelope
 if e:=dec.Decode(&env);e!=nil{return nil,errors.New("DK11_CHECKPOINT_DECODE")}
 if e:=dec.Decode(new(any));e!=io.EOF{return nil,errors.New("DK11_CHECKPOINT_TRAILING")}
 b,e:=json.Marshal(env.State)
 if e!=nil||env.State.Schema!=Schema||hash(b)!=env.SHA256{
  return nil,errors.New("DK11_CHECKPOINT_SHA")
 }
 old,e:=prior.Restore(env.State.Parent);if e!=nil{return nil,e}
 m:=&Agent{Policy:old,Memory:env.State.Memory,Samples:env.State.Samples}
 if !m.healthy(){return nil,errors.New("DK11_SOURCE_SELECTED_COUNTERS")}
 copy,e:=m.Snapshot()
 if e!=nil||!bytes.Equal(copy,raw){return nil,errors.New("DK11_NONCANONICAL")}
 return m,nil
}
