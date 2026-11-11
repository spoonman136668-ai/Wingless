// Package mg40 is an opt-in original-public selected task-effect memory.
// It wraps the unchanged native MG37/MG28 learner; it never steps MiniGrid.
package mg40

import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "io"
 "math"
 "regexp"
 "strings"
 original "github.com/spoonman136668-ai/Wingless/research/native/itd02/mg28"
 previous "github.com/spoonman136668-ai/Wingless/research/native/itd02/mg37"
)

const Schema="wingless.itd02.nr-m1q.selected-task-effect.v1"
const MaxContexts=4096
const MaxCount=4096
const MaxCheckpoint=12500000

type Public=original.Public
type Stat struct {
 Count uint32 `json:"count"`
 Effect uint32 `json:"effect"`
}
type Row [7]Stat
type Agent struct {
 Parent *previous.SourceAgent
 Task map[string]Row
 Enabled bool
 Pending bool
 PendingAction int
 PendingKey string
}
func New(enabled bool)*Agent{
 return &Agent{Parent:previous.NewSourceAgent(true),Task:map[string]Row{},Enabled:enabled}
}
func hash(b []byte)string{d:=sha256.Sum256(b);return hex.EncodeToString(d[:])}
func publicKey(p Public)(string,error){
 if len(p.Image)!=7||p.Direction<0||p.Direction>3||p.Mission=="" {return "",errors.New("NRM1Q_BAD_PUBLIC")}
 for _,r:=range p.Image {
  if len(r)!=7{return "",errors.New("NRM1Q_BAD_PUBLIC_ROW")}
  for _,pixel:=range r{if len(pixel)!=3{return "",errors.New("NRM1Q_BAD_PUBLIC_CELL")}}
 }
 f:=p.Image[3][5]
 shape:=struct{
  Mission string `json:"mission"`
  Direction int `json:"direction"`
  Front [3]int `json:"front"`
 }{p.Mission,p.Direction,[3]int{f[0],f[1],f[2]}}
 raw,e:=json.Marshal(shape);if e!=nil{return "",e}
 return hash(raw),nil
}
var literalPickup=regexp.MustCompile(`^pick up the (red|green|blue|purple|yellow|grey) (key|ball|box)$`)
var colors=map[string]int{"red":0,"green":1,"blue":2,"purple":3,"yellow":4,"grey":5}
var objects=map[string]int{"key":5,"ball":6,"box":7}

// TaskEffect refuses to credit image change, turns, door closing or mere sighting.
// This is explicitly engineered PUBLIC semantic labeling, not discovered physics.
func TaskEffect(before Public, selected int, after Public)(bool,error){
 if selected<0||selected>6{return false,errors.New("NRM1Q_INVALID_SELECTED_ACTION")}
 if _,e:=publicKey(before);e!=nil{return false,e}
 if _,e:=publicKey(after);e!=nil{return false,e}
 if before.Mission!=after.Mission{return false,errors.New("NRM1Q_MISSION_DRIFT")}
 if before.Direction!=after.Direction{return false,nil}
 b:=before.Image[3][5];a:=after.Image[3][5]
 if selected==5 {
  return b[0]==4&&a[0]==4&&b[1]==a[1]&&(b[2]==1||b[2]==2)&&a[2]==0,nil
 }
 if selected==3 {
  terms:=literalPickup.FindStringSubmatch(before.Mission)
  if terms==nil{return false,nil}
  return b[0]==objects[terms[2]]&&b[1]==colors[terms[1]]&&a[0]==1,nil
 }
 return false,nil
}
func score(s Stat)uint64{return 10000*(uint64(s.Effect)+1)/(uint64(s.Count)+2)}

// Propose is diagnostic and can alter an opt-in action ONLY after same-context,
// prior selected real task consequences. Never proposes turn/forward/view churn.
func Propose(parent int, row Row)(int,bool,error){
 if parent<0||parent>6{return 0,false,errors.New("NRM1Q_INVALID_PARENT")}
 p:=row[parent]
 if p.Count<2{return parent,false,nil}
 prior:=score(p)
 best:=parent
 high:=uint64(0)
 for _,a:=range []int{3,5}{
  if a==parent{continue}
  t:=row[a]
  if t.Count<2||t.Effect==0 {continue}
  s:=score(t)
  if s>=5000&&s>=prior+1000&&s>high{best=a;high=s}
 }
 return best,best!=parent,nil
}
func(a *Agent) BeginEpisode()error{
 if a==nil||a.Parent==nil||a.Task==nil||a.Pending{return errors.New("NRM1Q_UNFINISHED_SOURCE")}
 return a.Parent.BeginEpisode()
}
func(a *Agent) Choose(before Public)(*int,string,error){
 if a==nil||a.Parent==nil||a.Task==nil||a.Pending{return nil,"",errors.New("NRM1Q_PENDING")}
 k,e:=publicKey(before);if e!=nil{return nil,"",e}
 base,status,e:=a.Parent.Choose(before)
 if e!=nil||base==nil{return nil,"",errors.New("NRM1Q_PARENT_DECISION_MISSING")}
 if !a.Enabled{return base,status,nil}
 chosen,override,e:=Propose(*base,a.Task[k]);if e!=nil{return nil,"",e}
 if !override{return base,status,nil}
 return &chosen,"NRM1Q_PUBLIC_TASK_EFFECT_EVIDENCE_OVERRIDE",nil
}
func(a *Agent) Commit(before Public,selected int)error{
 if a==nil||a.Parent==nil||a.Task==nil||a.Pending||selected<0||selected>6{return errors.New("NRM1Q_BAD_COMMIT")}
 k,e:=publicKey(before);if e!=nil{return e}
 if e=a.Parent.Commit(before,selected);e!=nil{return e}
 a.Pending=true;a.PendingAction=selected;a.PendingKey=k
 return nil
}
func(a *Agent) Feedback(before Public,selected int,after Public,reward float64,terminal bool)error{
 if a==nil||a.Parent==nil||a.Task==nil||!a.Pending||a.PendingAction!=selected{return errors.New("NRM1Q_UNSELECTED_SOURCE_FEEDBACK")}
 k,e:=publicKey(before);if e!=nil||k!=a.PendingKey{return errors.New("NRM1Q_SOURCE_CONTEXT_MISMATCH")}
 if math.IsNaN(reward)||math.IsInf(reward,0)||reward< -1||reward>1{return errors.New("NRM1Q_UNQUALIFIED_REWARD")}
 effect,e:=TaskEffect(before,selected,after);if e!=nil{return e}
 if _,found:=a.Task[k];!found&&len(a.Task)>=MaxContexts{return errors.New("NRM1Q_CONTEXT_CAP")}
 row:=a.Task[k];if row[selected].Count>=MaxCount{return errors.New("NRM1Q_SAMPLE_CAP")}
 // Only authentic native source feedback can update the opt-in task cells.
 if e=a.Parent.Feedback(before,selected,after,reward,terminal);e!=nil{return e}
 stat:=row[selected];stat.Count++;if effect{stat.Effect++}
 row[selected]=stat;a.Task[k]=row
 a.Pending=false;a.PendingKey=""
 return nil
}
func(a *Agent) valid()bool{
 if a==nil||a.Parent==nil||a.Task==nil||len(a.Task)>MaxContexts||a.Pending{return false}
 if len(a.Task)!=len(a.Parent.Rewards){return false}
 for k,row:=range a.Task {
  if len(k)!=64||strings.Trim(k,"0123456789abcdef")!=""{return false}
  parent,ok:=a.Parent.Rewards[k];if !ok{return false}
  for i,s:=range row {
   if s.Count>MaxCount||s.Effect>s.Count||s.Count!=parent[i].Count{return false}
  }
 }
 return true
}
type state struct{
 Schema string `json:"schema"`
 Parent json.RawMessage `json:"parent"`
 Task map[string]Row `json:"task"`
 Enabled bool `json:"enabled"`
}
type envelope struct {
 State state `json:"state"`
 SHA string `json:"sha256"`
}
func(a *Agent) Snapshot()([]byte,error){
 if !a.valid(){return nil,errors.New("NRM1Q_UNTRUSTED_NATIVE_STATE")}
 parent,e:=a.Parent.Snapshot();if e!=nil{return nil,e}
 s:=state{Schema,parent,a.Task,a.Enabled}
 body,e:=json.Marshal(s);if e!=nil{return nil,e}
 wire,e:=json.Marshal(envelope{s,hash(body)})
 if e!=nil||len(wire)>MaxCheckpoint{return nil,errors.New("NRM1Q_CHECKPOINT_CAP")}
 return wire,nil
}
func Restore(raw []byte)(*Agent,error){
 if len(raw)==0||len(raw)>MaxCheckpoint{return nil,errors.New("NRM1Q_BAD_CHECKPOINT_BYTES")}
 decoder:=json.NewDecoder(bytes.NewReader(raw));decoder.DisallowUnknownFields()
 var env envelope
 if e:=decoder.Decode(&env);e!=nil{return nil,e}
 if decoder.Decode(new(any))!=io.EOF{return nil,errors.New("NRM1Q_TRAILING_SOURCE")}
 body,e:=json.Marshal(env.State)
 if e!=nil||env.State.Schema!=Schema||hash(body)!=env.SHA{return nil,errors.New("NRM1Q_CHECKPOINT_DIGEST")}
 parent,e:=previous.RestoreSourceAgent(env.State.Parent);if e!=nil{return nil,e}
 agent:=&Agent{Parent:parent,Task:env.State.Task,Enabled:env.State.Enabled}
 if !agent.valid(){return nil,errors.New("NRM1Q_CHECKPOINT_EVIDENCE_MISMATCH")}
 rebuilt,e:=agent.Snapshot()
 if e!=nil||!bytes.Equal(raw,rebuilt){return nil,errors.New("NRM1Q_NONCANONICAL")}
 return agent,nil
}
