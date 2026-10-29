// Package dk10 adds source-public within-episode transition/loop memory
// around the unchanged Go DK9/DK8 learner. NOT original Wingless R8.
package dk10

import (
  "bytes"
  "crypto/sha256"
  "encoding/base64"
  "encoding/hex"
  "encoding/json"
  "errors"
  "io"
  prior "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk9"
)
const Schema="wingless.itd02.dk10.source-public-temporal-loop-memory.v1"
const MaxSteps=192
const MaxCheckpoint=300000
type Public=prior.Public
type Counts [7]uint16
type Policy struct{
  *prior.Policy
  Visits map[string]Counts
  Edges map[string]map[string]string
  Seen map[string]uint16
  Steps uint16
}
func New()*Policy{return &Policy{Policy:prior.New(),Visits:map[string]Counts{},Edges:map[string]map[string]string{},Seen:map[string]uint16{}}}
func hash(b []byte)string{sum:=sha256.Sum256(b);return hex.EncodeToString(sum[:])}
func hashPublic(public Public)(string,error){
 if e:=public.Valid();e!=nil{return "",e}
 raw,e:=json.Marshal(public);if e!=nil{return "",e}
 return hash(raw),nil
}
func (p *Policy) BeginEpisode()error{
 if e:=p.Policy.BeginEpisode();e!=nil{return e}
 p.Visits=map[string]Counts{}
 p.Edges=map[string]map[string]string{}
 p.Seen=map[string]uint16{}
 p.Steps=0
 return nil
}
func (p *Policy) Choose(before Public,objective string)(*int,string,error){
 if p.Steps>MaxSteps{return nil,"",errors.New("DK10_EPISODE_LIMIT")}
 base,_,e:=p.Policy.Choose(before,objective)
 if e!=nil{return nil,"",e}
 if base==nil{return nil,"",errors.New("DK10_ORIGINAL_PARENT_ABSTAIN")}
 key,e:=hashPublic(before);if e!=nil{return nil,"",e}
 visits:=p.Visits[key]
 next:=p.Edges[key]
 highest:=-100000
 var score [7]int
 for action:=0;action<7;action++{
  s:= -3*int(visits[action])
  if action==*base{s+=4}
  if next!=nil {
   if successor,known:=next[string(rune('0'+action))];known{
    if successor==key{s-=2}
    n:=p.Seen[successor]
    if n>=2{s-=2}else{s-=int(n)}
   }
  }
  score[action]=s
  if s>highest{highest=s}
 }
 raw,_:=hex.DecodeString(key[:2]);start:=int(raw[0])%7
 for offset:=0;offset<7;offset++{
  a:=(start+offset)%7
  if score[a]==highest{
   result:=a
   return &result,"SELECT_NATIVE_PUBLIC_TEMPORAL_ACTION",nil
  }
 }
 return nil,"",errors.New("DK10_UNREACHABLE")
}
func(p *Policy) Feedback(before Public,action int,after Public,reward float64,terminal bool)error{
 if p.Steps>=MaxSteps{return errors.New("DK10_TOO_MANY_SOURCE_STEPS")}
 old,e:=hashPublic(before);if e!=nil{return e}
 newHash,e:=hashPublic(after);if e!=nil{return e}
 if action<0||action>6{return errors.New("DK10_ACTION_INVALID")}
 if len(p.Visits)>=MaxSteps&&p.Visits[old]==(Counts{}) {
  // No new context at the memory bound; no selected feedback may be lost.
  if _,exists:=p.Visits[old];!exists{return errors.New("DK10_MEMORY_FULL")}
 }
 if e=p.Policy.Feedback(before,action,after,reward,terminal);e!=nil{return e}
 row:=p.Visits[old];row[action]++;p.Visits[old]=row
 if p.Edges[old]==nil{p.Edges[old]=map[string]string{}}
 p.Edges[old][string(rune('0'+action))]=newHash
 p.Seen[newHash]++
 p.Steps++
 return nil
}
type state struct{
 Schema string `json:"schema"`
 Inner string `json:"inner"`
 Visits map[string]Counts `json:"visits"`
 Edges map[string]map[string]string `json:"edges"`
 Seen map[string]uint16 `json:"seen"`
 Steps uint16 `json:"steps"`
}
type envelope struct{State state `json:"state"`;SHA256 string `json:"sha256"`}
func(p *Policy) Snapshot()([]byte,error){
 inner,e:=p.Policy.Snapshot();if e!=nil{return nil,e}
 if !p.validMemory(){return nil,errors.New("DK10_BAD_MEMORY")}
 s:=state{Schema:Schema,Inner:base64.StdEncoding.EncodeToString(inner),
    Visits:p.Visits,Edges:p.Edges,Seen:p.Seen,Steps:p.Steps}
 raw,e:=json.Marshal(s);if e!=nil{return nil,e}
 out,e:=json.Marshal(envelope{State:s,SHA256:hash(raw)})
 if e!=nil||len(out)>MaxCheckpoint{return nil,errors.New("DK10_SNAPSHOT_LIMIT")}
 return out,nil
}
func(p *Policy) validMemory()bool{
 if p==nil||p.Policy==nil||p.Steps>MaxSteps||
    p.Visits==nil||p.Edges==nil||p.Seen==nil||
    len(p.Visits)>MaxSteps||len(p.Seen)>MaxSteps||len(p.Edges)!=len(p.Visits){return false}
 var total,after uint32
 for key,actions:=range p.Visits{
  if len(key)!=64{return false}
  for _,n:=range actions{total+=uint32(n)}
  if _,ok:=p.Edges[key];!ok{return false}
  for action,destination:=range p.Edges[key]{
   if len(action)!=1||action[0]<'0'||action[0]>'6'||len(destination)!=64{return false}
  }
 }
 for key,count:=range p.Seen{if len(key)!=64||count==0{return false};after+=uint32(count)}
 return total==uint32(p.Steps)&&after==uint32(p.Steps)
}
func Restore(blob []byte)(*Policy,error){
 if len(blob)==0||len(blob)>MaxCheckpoint{return nil,errors.New("DK10_CHECKPOINT_SIZE")}
 dec:=json.NewDecoder(bytes.NewReader(blob));dec.DisallowUnknownFields()
 var env envelope
 if e:=dec.Decode(&env);e!=nil{return nil,e}
 if e:=dec.Decode(new(any));e!=io.EOF{return nil,errors.New("DK10_TRAILING")}
 raw,e:=json.Marshal(env.State)
 if e!=nil||env.State.Schema!=Schema||hash(raw)!=env.SHA256{
  return nil,errors.New("DK10_SOURCE_SHA_OR_SCHEMA")}
 orig,e:=base64.StdEncoding.DecodeString(env.State.Inner);if e!=nil{return nil,e}
 original,e:=prior.Restore(orig);if e!=nil{return nil,e}
 p:=&Policy{Policy:original,Visits:env.State.Visits,
   Edges:env.State.Edges,Seen:env.State.Seen,Steps:env.State.Steps}
 if !p.validMemory(){return nil,errors.New("DK10_SOURCE_HISTORY_INVALID")}
 repeat,e:=p.Snapshot()
 if e!=nil||!bytes.Equal(repeat,blob){return nil,errors.New("DK10_NONCANONICAL")}
 return p,nil
}
