// Package dk12: learn the difference between *public camera turns* and
// *same-heading changes to source-visible scenes* from original selected actions.
// No coded mapping of action index to locomotion, and no external stage oracle.
package dk12

import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "io"
 prior "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk11"
)
const Schema="wingless.itd02.dk12.public-nonrotational-selected-effects.v1"
const MaxSamples=4096
const MaxEpisode=192
const MaxCheckpoint=600000
type Public=prior.Public
type Stat struct{
 N uint32 `json:"n"`
 SameHeadingChanged uint32 `json:"same_heading_changed"`
 HeadingChanged uint32 `json:"heading_changed"`
 HeldLost uint32 `json:"held_lost"`
 HeldRetainedChanged uint32 `json:"held_retained_changed"`
}
type Policy struct{
 *prior.Agent
 Motion [7]Stat
 Episode [7]uint16
 Live bool
}
func New()*Policy{return &Policy{Agent:prior.New()}}
func hash(raw []byte)string{h:=sha256.Sum256(raw);return hex.EncodeToString(h[:])}
func valid(m *Policy)bool{
 if m==nil||m.Agent==nil||m.Samples>MaxSamples{return false}
 total:=uint64(0)
 for _,v:=range m.Motion{
  if v.SameHeadingChanged>v.N||v.HeadingChanged>v.N||
    v.HeldLost>v.N||v.HeldRetainedChanged>v.SameHeadingChanged{return false}
  total+=uint64(v.N)
 }
 if total!=m.Samples{return false}
 count:=uint64(0);for _,n:=range m.Episode{count+=uint64(n)}
 return count<=MaxEpisode
}
func (m *Policy) BeginEpisode()error{
 if !valid(m){return errors.New("DK12_INVALID_NATIVE_SAMPLE_STATE")}
 if e:=m.Agent.BeginEpisode();e!=nil{return e}
 m.Live=true
 m.Episode=[7]uint16{}
 return nil
}
func (m *Policy) Choose(before Public)(*int,string,error){
 if !valid(m){return nil,"",errors.New("DK12_UNHEALTHY")}
 if e:=before.Valid();e!=nil{return nil,"",e}
 if _,e:=m.Agent.Snapshot();e!=nil{return nil,"",e}
 best:=-1;highest:=int64(-1<<60)
 for action,s:=range m.Motion{
  n:=int64(s.N)+64
  numerator:=int64(12)*int64(s.SameHeadingChanged)-
             int64(4)*int64(s.HeadingChanged)-
             int64(8)*int64(s.HeldLost)
  if before.Image[3][6][0]==5||before.Image[3][6][0]==6||before.Image[3][6][0]==7{
   numerator+=int64(4)*int64(s.HeldRetainedChanged)
  }
  score:=(numerator*1024)/n-int64(96)*int64(m.Episode[action])
  if score>highest{highest=score;best=action}
 }
 if highest>0{
  return &best,"SELECT_NATIVE_LEARNED_SAME_HEADING_PUBLIC_EFFECT",nil
 }
 // No source progress evidence: explore least actual selected action count,
 // with a raw public-data digest tie, not a hardwired "forward" action.
 least:=m.Episode[0]
 for _,n:=range m.Episode{if n<least{least=n}}
 raw,_:=json.Marshal(before)
 key:=sha256.Sum256(raw)
 start:=int(key[0])%7
 for i:=0;i<7;i++{
  a:=(start+i)%7
  if m.Episode[a]==least{return &a,"SELECT_NATIVE_ANTI_CYCLE_UNKNOWN_ACTION",nil}
 }
 return nil,"",errors.New("DK12_UNREACHABLE_NATIVE_SELECTION")
}
func (m *Policy) Commit(before Public,action int)error{
 if !valid(m){return errors.New("DK12_UNHEALTHY_BEFORE_COMMIT")}
 return m.Agent.Commit(before,action)
}
func (m *Policy) selected(before Public,action int,after Public,reward float64,terminal,fit bool)error{
 if !valid(m)||action<0||action>6||m.Samples>=MaxSamples{
  return errors.New("DK12_INVALID_SELECTED_SOURCE")
 }
 if e:=before.Valid();e!=nil{return e}
 if e:=after.Valid();e!=nil{return e}
 if !fit&&(!m.Live||m.Episode[action]>=MaxEpisode){
  return errors.New("DK12_LIVE_EPISODE_BOUND")
 }
 var e error
 if fit{e=m.Agent.ObserveFit(before,action,after,reward,terminal)
 }else{e=m.Agent.Feedback(before,action,after,reward,terminal)}
 if e!=nil{return e}
 stat:=m.Motion[action];stat.N++
 if before.Direction!=after.Direction{stat.HeadingChanged++}
 if before.Direction==after.Direction&&!bytes.Equal(mustJSON(before.Image),mustJSON(after.Image)){
  stat.SameHeadingChanged++
  bh,ah:=before.Image[3][6],after.Image[3][6]
  held:=bh[0]==5||bh[0]==6||bh[0]==7
  kept:=held&&bytes.Equal(mustJSON(bh),mustJSON(ah))
  if kept{stat.HeldRetainedChanged++}
 }
 b:=before.Image[3][6][0];a:=after.Image[3][6][0]
 heldBefore:=b==5||b==6||b==7
 heldAfter:=a==5||a==6||a==7
 if heldBefore&&!heldAfter{stat.HeldLost++}
 m.Motion[action]=stat
 if !fit{m.Episode[action]++}
 if !valid(m){return errors.New("DK12_SOURCE_COUNTS_DRIFT")}
 return nil
}
func mustJSON(v any)[]byte{b,_:=json.Marshal(v);return b}
func (m *Policy) ObserveFit(before Public,action int,after Public,reward float64,terminal bool)error{
 return m.selected(before,action,after,reward,terminal,true)
}
func (m *Policy) Feedback(before Public,action int,after Public,reward float64,terminal bool)error{
 return m.selected(before,action,after,reward,terminal,false)
}
type state struct{
 Schema string `json:"schema"`
 Parent json.RawMessage `json:"parent"`
 Motion [7]Stat `json:"motion"`
 Episode [7]uint16 `json:"episode"`
 Live bool `json:"live"`
}
type envelope struct{State state `json:"state"`;SHA256 string `json:"sha256"`}
func (m *Policy) Snapshot()([]byte,error){
 if !valid(m){return nil,errors.New("DK12_UNHEALTHY_CHECKPOINT")}
 old,e:=m.Agent.Snapshot();if e!=nil{return nil,e}
 state:=state{Schema:Schema,Parent:old,Motion:m.Motion,Episode:m.Episode,Live:m.Live}
 b,e:=json.Marshal(state);if e!=nil{return nil,e}
 encoded,e:=json.Marshal(envelope{State:state,SHA256:hash(b)})
 if e!=nil||len(encoded)>MaxCheckpoint{return nil,errors.New("DK12_SNAPSHOT_BUDGET")}
 return encoded,nil
}
func Restore(raw []byte)(*Policy,error){
 if len(raw)==0||len(raw)>MaxCheckpoint{return nil,errors.New("DK12_CHECKPOINT_SIZE")}
 d:=json.NewDecoder(bytes.NewReader(raw));d.DisallowUnknownFields()
 var e envelope
 if err:=d.Decode(&e);err!=nil{return nil,err}
 if err:=d.Decode(new(any));err!=io.EOF{return nil,errors.New("DK12_CHECKPOINT_TRAILING")}
 b,err:=json.Marshal(e.State)
 if err!=nil||e.State.Schema!=Schema||hash(b)!=e.SHA256{return nil,errors.New("DK12_UNTRUSTED_SHA")}
 orig,err:=prior.Restore(e.State.Parent);if err!=nil{return nil,err}
 m:=&Policy{Agent:orig,Motion:e.State.Motion,Episode:e.State.Episode,Live:e.State.Live}
 if !valid(m){return nil,errors.New("DK12_SOURCE_SELECTED_TOTALS")}
 again,err:=m.Snapshot()
 if err!=nil||!bytes.Equal(again,raw){return nil,errors.New("DK12_NONCANONICAL")}
 return m,nil
}
