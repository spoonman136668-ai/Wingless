// Package mg19 wraps unchanged DK16 with a genuine original-public mission envelope.
// Source-only prospective engineering: this does not prove learning or transfer.
package mg19

import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "io"
 "strings"
 base "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk16"
)

const Schema="wingless.itd02.mg19.original-full-public-mission-envelope.v1"
const MaxSamples=4096
const MaxContexts=4096
const MaxCheckpoint=3500000
const Prefix="pick up the "
type Public=base.Public
type Stat struct {Count uint32 `json:"count"`; Positive uint32 `json:"positive"`}
type Actions [7]Stat
type awaiting struct {Before string; Mission string; Action int}
type Model struct {
 Parent *base.Policy
 Contexts map[string]Actions
 Missions map[string]uint64
 Pending *awaiting
}
func New()*Model {
 return &Model{Parent:base.New(),Contexts:map[string]Actions{},Missions:map[string]uint64{}}
}
func hash(raw []byte)string{h:=sha256.Sum256(raw);return hex.EncodeToString(h[:])}
func admitted(m string)bool{
 if m=="use the key to open the door and then get to the goal" ||
    m=="get to the green goal square"{return true}
 if len(m)>96||!strings.HasPrefix(m,Prefix){return false}
 rest:=strings.TrimPrefix(m,Prefix)
 words:=strings.Split(rest," ")
 if len(words)<2||len(words)>4{return false}
 for _,word:=range words{
  if len(word)==0{return false}
  for _,c:=range word {if c<'a'||c>'z'{return false}}
 }
 return true
}
// The source mission is kept unchanged in all MG19 keys. Only the inherited
// DoorKey-only native parent's private compatibility input is projected.
func legacy(p Public)(Public,error){
 if !admitted(p.Mission){return p,errors.New("MG19_UNADMITTED_SOURCE_MISSION")}
 projected:=p
 if p.Mission!="use the key to open the door and then get to the goal"&&
    p.Mission!="get to the green goal square"{
  projected.Mission="use the key to open the door and then get to the goal"
 }
 if err:=projected.Valid();err!=nil{return p,err}
 return projected,nil
}
func sourceHash(p Public)(string,error){
 if _,e:=legacy(p);e!=nil{return "",e}
 raw,e:=json.Marshal(p);if e!=nil{return "",e}
 return hash(raw),nil
}
func (m *Model) valid()bool{
 if m==nil||m.Parent==nil||m.Contexts==nil||m.Missions==nil||
 len(m.Contexts)>MaxContexts||m.Parent.Agent.Samples>MaxSamples{return false}
 var n uint64
 for k,row:=range m.Contexts{
  if len(k)!=64{return false}
  for _,v:=range k {if !((v>='0'&&v<='9')||(v>='a'&&v<='f')){return false}}
  for _,stat:=range row{
   if stat.Positive>stat.Count{return false}
   n+=uint64(stat.Count)
  }
 }
 var missions uint64
 for name,count:=range m.Missions{
  if !admitted(name)||count==0{return false}
  missions+=count
 }
 if m.Pending!=nil && (!admitted(m.Pending.Mission)||m.Pending.Action<0||
 m.Pending.Action>6||len(m.Pending.Before)!=64){return false}
 return n==m.Parent.Agent.Samples&&missions==n
}
func (m *Model) BeginEpisode()error{
 if !m.valid()||m.Pending!=nil{return errors.New("MG19_INVALID_EPISODE")}
 return m.Parent.BeginEpisode()
}
func (m *Model) Choose(before Public)(*int,string,error){
 if !m.valid()||m.Pending!=nil{return nil,"",errors.New("MG19_CHOOSE_PENDING")}
 projected,e:=legacy(before);if e!=nil{return nil,"",e}
 key,e:=sourceHash(before);if e!=nil{return nil,"",e}
 // On genuine unfamiliar pickup missions DK14's global mission reward is NOT
 // consulted; retain the original DK13 native contextual public successor.
 if before.Mission=="use the key to open the door and then get to the goal"||
    before.Mission=="get to the green goal square"{
  return m.Parent.Choose(before)
 }
 native,status,e:=m.Parent.Policy.Policy.Choose(projected)
 if e!=nil{return nil,"",e}
 best:=-1;highest:=int64(0)
 for i,st:=range m.Contexts[key]{
  if st.Positive==0{continue}
  score:=4096*int64(st.Positive)/(int64(st.Count)+64)
  if score>highest {highest=score;best=i}
 }
 if best>=0{return &best,"MG19_NATIVE_EXACT_ORIGINAL_MISSION_STATE_REWARD",nil}
 return native,status,nil
}
func (m *Model) Commit(before Public,action int)error{
 if !m.valid()||m.Pending!=nil||action<0||action>6{return errors.New("MG19_BAD_SOURCE_COMMIT")}
 projected,e:=legacy(before);if e!=nil{return e}
 beforeHash,e:=sourceHash(before);if e!=nil{return e}
 if e=m.Parent.Commit(projected,action);e!=nil{return e}
 m.Pending=&awaiting{Before:beforeHash,Mission:before.Mission,Action:action}
 return nil
}
func (m *Model) learn(before Public,action int,after Public,reward float64,terminal,live bool)error{
 if !m.valid()||action<0||action>6||before.Mission!=after.Mission{return errors.New("MG19_WRONG_SOURCE_FEEDBACK")}
 b,e:=legacy(before);if e!=nil{return e}
 n,e:=legacy(after);if e!=nil{return e}
 key,e:=sourceHash(before);if e!=nil{return e}
 next,e:=sourceHash(after);if e!=nil{return e}
 if live{
  if m.Pending==nil||m.Pending.Before!=key||m.Pending.Action!=action||
    m.Pending.Mission!=before.Mission{return errors.New("MG19_NOT_ORIGINAL_SELECTED_SOURCE")}
 }else if m.Pending!=nil{return errors.New("MG19_FIT_WITH_PENDING_LIVE_SOURCE")}
 if m.Parent.Agent.Samples>=MaxSamples{return errors.New("MG19_TOTAL_BUDGET")}
 if _,ok:=m.Contexts[key];!ok&&len(m.Contexts)>=MaxContexts{return errors.New("MG19_CONTEXT_BUDGET")}
 if live {e=m.Parent.Feedback(b,action,n,reward,terminal)
 }else {e=m.Parent.ObserveFit(b,action,n,reward,terminal)}
 if e!=nil{return e}
 stats:=m.Contexts[key];cell:=stats[action];cell.Count++
 if reward>0&&terminal&&key!=next {cell.Positive++}
 stats[action]=cell;m.Contexts[key]=stats
 m.Missions[before.Mission]++
 if live{m.Pending=nil}
 if !m.valid(){return errors.New("MG19_NATIVE_SOURCE_COUNTERS")}
 return nil
}
func (m *Model) ObserveFit(b Public,a int,n Public,r float64,terminal bool)error{
 return m.learn(b,a,n,r,terminal,false)
}
func (m *Model) Feedback(b Public,a int,n Public,r float64,terminal bool)error{
 return m.learn(b,a,n,r,terminal,true)
}
type state struct{
 Schema string `json:"schema"`
 Parent json.RawMessage `json:"parent"`
 Contexts map[string]Actions `json:"contexts"`
 Missions map[string]uint64 `json:"missions"`
}
type envelope struct{State state `json:"state"`;SHA256 string `json:"sha256"`}
func (m *Model) Snapshot()([]byte,error){
 if !m.valid()||m.Pending!=nil{return nil,errors.New("MG19_UNSAFE_SNAPSHOT")}
 parent,e:=m.Parent.Snapshot();if e!=nil{return nil,e}
 value:=state{Schema:Schema,Parent:parent,Contexts:m.Contexts,Missions:m.Missions}
 canonical,e:=json.Marshal(value);if e!=nil{return nil,e}
 raw,e:=json.Marshal(envelope{State:value,SHA256:hash(canonical)})
 if e!=nil||len(raw)>MaxCheckpoint{return nil,errors.New("MG19_SNAPSHOT_CAP")}
 return raw,nil
}
func Restore(raw []byte)(*Model,error){
 if len(raw)==0||len(raw)>MaxCheckpoint{return nil,errors.New("MG19_RESTORE_CAP")}
 decoder:=json.NewDecoder(bytes.NewReader(raw));decoder.DisallowUnknownFields()
 var wire envelope
 if e:=decoder.Decode(&wire);e!=nil{return nil,e}
 if decoder.Decode(new(any))!=io.EOF{return nil,errors.New("MG19_RESTORE_TRAILING")}
 canonical,e:=json.Marshal(wire.State)
 if e!=nil||wire.State.Schema!=Schema||hash(canonical)!=wire.SHA256{return nil,errors.New("MG19_RESTORE_SHA")}
 parent,e:=base.Restore(wire.State.Parent);if e!=nil{return nil,e}
 m:=&Model{Parent:parent,Contexts:wire.State.Contexts,Missions:wire.State.Missions}
 if !m.valid(){return nil,errors.New("MG19_RESTORE_NATIVE_SOURCE_INVALID")}
 reencoded,e:=m.Snapshot()
 if e!=nil||!bytes.Equal(raw,reencoded){return nil,errors.New("MG19_NONCANONICAL_SNAPSHOT")}
 return m,nil
}
