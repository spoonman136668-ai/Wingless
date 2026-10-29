// Package dk12 adds a bounded own-selected source-public state-action novelty
// memory to the unchanged Go DK11 learner. It is an explicit exploration
// hypothesis, not a map, route, task-stage oracle, or Wingless original R8.
package dk12

import (
 "bytes"
 "crypto/sha256"
 "encoding/base64"
 "encoding/hex"
 "encoding/json"
 "errors"
 "io"
 prior "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk11"
)
const Schema="wingless.itd02.dk12.public-selected-episodic-novelty.v1"
const MaxSteps=192
const MaxStateBytes=800000
type Public=prior.Public
type Cell struct{
 Count [7]uint16 `json:"count"`
 Novel [7]uint16 `json:"novel"`
 Loop [7]uint16 `json:"loop"`
 Next [7]string `json:"next"`
}
type Policy struct{
 *prior.Agent
 State map[string]Cell
 Seen map[string]uint16
 Steps uint16
}
func New()*Policy{return &Policy{Agent:prior.New(),State:map[string]Cell{},Seen:map[string]uint16{}}}
func hash(raw []byte)string{x:=sha256.Sum256(raw);return hex.EncodeToString(x[:])}
func publicHash(p Public)(string,error){
 if e:=p.Valid();e!=nil{return "",e}
 data,e:=json.Marshal(p);if e!=nil{return "",e}
 return hash(data),nil
}
func(p *Policy) BeginEpisode()error{
 if e:=p.Agent.BeginEpisode();e!=nil{return e}
 p.State=map[string]Cell{}
 p.Seen=map[string]uint16{}
 p.Steps=0
 return nil
}
func(p *Policy) Choose(b Public)(*int,string,error){
 if p.Steps>=MaxSteps{return nil,"",errors.New("DK12_ORIGINAL_EPISODE_CAP")}
 state,e:=publicHash(b);if e!=nil{return nil,"",e}
 if _,e=p.Agent.Snapshot();e!=nil{return nil,"",e}
 row:=p.State[state]
 start,_:=hex.DecodeString(state[:2]);shift:=int(start[0])%7
 for i:=0;i<7;i++{
  a:=(shift+i)%7
  if row.Count[a]==0{
   out:=a
   return &out,"SELECT_NATIVE_UNTRIED_PUBLIC_STATE_ACTION",nil
  }
 }
 best:=-1000000
 var scores [7]int
 for a:=0;a<7;a++{
  succ:=0
  if row.Next[a]!=""{v:=p.Seen[row.Next[a]];if v>=2{succ=2}else{succ=int(v)}}
  score:=4*int(row.Novel[a])-3*int(row.Count[a])-2*int(row.Loop[a])-succ
  scores[a]=score
  if score>best{best=score}
 }
 for i:=0;i<7;i++{a:=(shift+i)%7;if scores[a]==best{out:=a;return &out,"SELECT_NATIVE_OBSERVED_PUBLIC_NOVELTY",nil}}
 return nil,"",errors.New("DK12_NO_NATIVE_SOURCE_ACTION")
}
func(p *Policy) Feedback(before Public,action int,after Public,reward float64,terminal bool)error{
 if p.Steps>=MaxSteps||action<0||action>6{return errors.New("DK12_SOURCE_ACTION_CAP")}
 key,e:=publicHash(before);if e!=nil{return e}
 next,e:=publicHash(after);if e!=nil{return e}
 if _,exists:=p.State[key];!exists&&len(p.State)>=MaxSteps{return errors.New("DK12_STATE_CAP")}
 if p.Agent.Feedback(before,action,after,reward,terminal)!=nil{return errors.New("DK12_ORIGINAL_SELECTED_FEEDBACK_REJECTED")}
 row:=p.State[key]
 row.Count[action]++
 if p.Seen[next]==0{row.Novel[action]++}
 if next==key{row.Loop[action]++}
 row.Next[action]=next;p.State[key]=row
 p.Seen[next]++
 p.Steps++
 return nil
}
type state struct{
 Schema string `json:"schema"`
 Inner string `json:"inner"`
 State map[string]Cell `json:"state"`
 Seen map[string]uint16 `json:"seen"`
 Steps uint16 `json:"steps"`
}
type envelope struct{State state `json:"state"`;SHA256 string `json:"sha256"`}
func(p *Policy) healthy()bool{
 if p==nil||p.Agent==nil||p.State==nil||p.Seen==nil||
 p.Steps>MaxSteps||len(p.State)>MaxSteps||len(p.Seen)>MaxSteps{return false}
 count:=uint32(0);seen:=uint32(0)
 for k,v:=range p.State{
  if len(k)!=64{return false}
  for a:=0;a<7;a++{
   if v.Novel[a]>v.Count[a]||v.Loop[a]>v.Count[a]||
    (v.Count[a]>0&&len(v.Next[a])!=64)|| (v.Count[a]==0&&v.Next[a]!=""){return false}
   count+=uint32(v.Count[a])
  }
 }
 for k,n:=range p.Seen{if len(k)!=64||n==0{return false};seen+=uint32(n)}
 return count==uint32(p.Steps)&&seen==uint32(p.Steps)
}
func(p *Policy) Snapshot()([]byte,error){
 if !p.healthy(){return nil,errors.New("DK12_BAD_PUBLIC_SOURCE_MEMORY")}
 base,e:=p.Agent.Snapshot();if e!=nil{return nil,e}
 st:=state{Schema:Schema,Inner:base64.StdEncoding.EncodeToString(base),State:p.State,Seen:p.Seen,Steps:p.Steps}
 body,e:=json.Marshal(st);if e!=nil{return nil,e}
 out,e:=json.Marshal(envelope{State:st,SHA256:hash(body)})
 if e!=nil||len(out)>MaxStateBytes{return nil,errors.New("DK12_CHECKPOINT_LIMIT")}
 return out,nil
}
func Restore(raw []byte)(*Policy,error){
 if len(raw)==0||len(raw)>MaxStateBytes{return nil,errors.New("DK12_UNTRUSTED_CHECKPOINT_BYTES")}
 dec:=json.NewDecoder(bytes.NewReader(raw));dec.DisallowUnknownFields()
 var env envelope
 if e:=dec.Decode(&env);e!=nil{return nil,e}
 if e:=dec.Decode(new(any));e!=io.EOF{return nil,errors.New("DK12_TRAILING_UNTRUSTED_SOURCE")}
 st:=env.State
 body,e:=json.Marshal(st)
 if e!=nil||st.Schema!=Schema||hash(body)!=env.SHA256{return nil,errors.New("DK12_CHECKPOINT_SHA")}
 original,e:=base64.StdEncoding.DecodeString(st.Inner);if e!=nil{return nil,e}
 parent,e:=prior.Restore(original);if e!=nil{return nil,e}
 p:=&Policy{Agent:parent,State:st.State,Seen:st.Seen,Steps:st.Steps}
 if !p.healthy(){return nil,errors.New("DK12_UNVERIFIED_SELECTED_HISTORY")}
 c,e:=p.Snapshot();if e!=nil||!bytes.Equal(c,raw){return nil,errors.New("DK12_NONCANONICAL_SOURCE")}
 return p,nil
}
