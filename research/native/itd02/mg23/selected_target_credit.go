// Package mg23 independently adds selected-public target/event association and
// delayed original terminal source credit without modifying the MG19 parent.
package mg23

import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "io"
 "strings"
 old "github.com/spoonman136668-ai/Wingless/research/native/itd02/mg19"
)

const Schema="wingless.itd02.mg23.original-public-target-delayed-selected-credit.v1"
const MaxSamples=4096
const MaxContext=4096
const Window=32
const MaxCheckpoint=5500000
type Public=old.Public
type Cell struct{
 Count uint32 `json:"count"`
 Target uint32 `json:"target"`
 Delayed uint32 `json:"delayed"`
}
type Row [7]Cell
type Selected struct{
 Context string `json:"context"`
 Action int `json:"action"`
}
type Agent struct {
 Parent *old.Model
 Cells map[string]Row
 Trace []Selected
}
func New()*Agent{return &Agent{Parent:old.New(),Cells:map[string]Row{},Trace:[]Selected{}}}
func hash(raw []byte)string {d:=sha256.Sum256(raw);return hex.EncodeToString(d[:])}
func allowedPublic(p Public)bool{
 if len(p.Image)!=7||p.Direction<0||p.Direction>3{return false}
 for _,row:=range p.Image{
  if len(row)!=7{return false}
  for _,pixel:=range row{
   if len(pixel)!=3{return false}
   for _,v:=range pixel{if v<0||v>255{return false}}
  }
 }
 if p.Mission=="use the key to open the door and then get to the goal"||p.Mission=="get to the green goal square"{return true}
 if len(p.Mission)>96||!strings.HasPrefix(p.Mission,"pick up the "){return false}
 words:=strings.Split(strings.TrimPrefix(p.Mission,"pick up the ")," ")
 if len(words)<2||len(words)>4{return false}
 for _,word:=range words{
  if len(word)==0{return false}
  for _,c:=range word{if c<'a'||c>'z'{return false}}
 }
 return true
}
func context(p Public)(string,error){
 if !allowedPublic(p){return "",errors.New("MG23_BAD_ORIGINAL_PUBLIC")}
 pixel:=p.Image[3][5]
 value:=struct{
  Mission string `json:"mission"`
  Direction int `json:"direction"`
  Front [3]int `json:"front"`
 }{p.Mission,p.Direction,[3]int{pixel[0],pixel[1],pixel[2]}}
 data,e:=json.Marshal(value)
 if e!=nil{return "",e}
 return hash(data),nil
}
func missionObject(m string)(int,int,bool){
 if !strings.HasPrefix(m,"pick up the "){return 0,0,false}
 words:=strings.Split(strings.TrimPrefix(m,"pick up the ")," ")
 if len(words)!=2{return 0,0,false}
 color,ok:=map[string]int{"red":0,"green":1,"blue":2,"purple":3,"yellow":4,"grey":5}[words[0]]
 if !ok{return 0,0,false}
 obj,ok:=map[string]int{"key":5,"ball":6,"box":7}[words[1]]
 return obj,color,ok
}
func targetEffect(before,after Public, action int)bool{
 object,color,valid:=missionObject(before.Mission)
 if !valid||before.Direction!=after.Direction||action!=3{return false}
 b,a:=before.Image[3][5],after.Image[3][5]
 return b[0]==object&&b[1]==color&&a[0]==1
}
func(a *Agent)valid()bool{
 if a==nil||a.Parent==nil||a.Cells==nil||a.Trace==nil||
 len(a.Cells)>MaxContext||a.Parent.Parent.Agent.Samples>MaxSamples||
 len(a.Trace)>Window{return false}
 var n uint64
 for k,row:=range a.Cells{
  if len(k)!=64{return false}
  for _,c:=range k{if !((c>='0'&&c<='9')||(c>='a'&&c<='f')){return false}}
  for _,v:=range row{
   if v.Target>v.Count{return false}
   n+=uint64(v.Count)
  }
 }
 if n!=a.Parent.Parent.Agent.Samples{return false}
 for _,v:=range a.Trace{
  if len(v.Context)!=64||v.Action<0||v.Action>6{return false}
  if _,ok:=a.Cells[v.Context];!ok{return false}
 }
 return true
}
func(a *Agent)BeginEpisode()error{
 if !a.valid(){return errors.New("MG23_INVALID_SOURCE_EPISODE")}
 if e:=a.Parent.BeginEpisode();e!=nil{return e}
 a.Trace=[]Selected{}
 return nil
}
func(a *Agent)Choose(before Public)(*int,string,error){
 if !a.valid(){return nil,"",errors.New("MG23_INVALID_SOURCE_CHOOSE")}
 key,e:=context(before)
 if e!=nil{return nil,"",e}
 base,status,e:=a.Parent.Choose(before)
 if e!=nil||base==nil{return nil,"",errors.New("MG23_PARENT_SOURCE_CHOOSE_FAILED")}
 if *base<0||*base>6{return nil,"",errors.New("MG23_INVALID_NATIVE_ACTION")}
 row:=a.Cells[key]
 best:=*base
 highest:=uint64(0)
 for action,v:=range row{
  if v.Count==0{continue}
  score:=uint64(8)*uint64(v.Target)+uint64(v.Delayed)
  if score>highest{highest=score;best=action}
 }
 if highest>0{
  selected:=best
  return &selected,"SELECT_NATIVE_MG23_OWN_OBSERVED_TARGET_OR_TERMINAL_CREDIT",nil
 }
 return base,status,nil
}
func(a *Agent)Commit(before Public,action int)error{
 if !a.valid()||action<0||action>6{return errors.New("MG23_COMMIT_BAD_SOURCE")}
 if _,e:=context(before);e!=nil{return e}
 return a.Parent.Commit(before,action)
}
func(a *Agent)learn(before Public,action int,after Public,reward float64,terminal,live bool)error{
 if !a.valid()||action<0||action>6||before.Mission!=after.Mission{
  return errors.New("MG23_MISSION_OR_ACTION_SOURCE_DRIFT")
 }
 key,e:=context(before);if e!=nil{return e}
 if _,e=context(after);e!=nil{return e}
 if a.Parent.Parent.Agent.Samples>=MaxSamples{return errors.New("MG23_SAMPLE_BOUND")}
 if _,has:=a.Cells[key];!has&&len(a.Cells)>=MaxContext{return errors.New("MG23_CONTEXT_BOUND")}
 if live{e=a.Parent.Feedback(before,action,after,reward,terminal)
 }else{e=a.Parent.ObserveFit(before,action,after,reward,terminal)}
 if e!=nil{return e}
 row:=a.Cells[key];v:=row[action];v.Count++
 if targetEffect(before,after,action){v.Target++}
 row[action]=v;a.Cells[key]=row
 if len(a.Trace)>=Window{a.Trace=append([]Selected{},a.Trace[1:]...)}
 a.Trace=append(a.Trace,Selected{key,action})
 if terminal&&reward>0{
  for i:=len(a.Trace)-1;i>=0;i--{
   observation:=a.Trace[i]
   row:=a.Cells[observation.Context]
   stat:=row[observation.Action]
   stat.Delayed+=uint32(Window-(len(a.Trace)-1-i))
   row[observation.Action]=stat
   a.Cells[observation.Context]=row
  }
 }
 if !a.valid(){return errors.New("MG23_SELECTED_CREDIT_COUNTER_DRIFT")}
 return nil
}
func(a *Agent)ObserveFit(b Public,selected int,n Public,reward float64,terminal bool)error{
 return a.learn(b,selected,n,reward,terminal,false)
}
func(a *Agent)Feedback(b Public,selected int,n Public,reward float64,terminal bool)error{
 return a.learn(b,selected,n,reward,terminal,true)
}
type state struct{
 Schema string `json:"schema"`
 Parent json.RawMessage `json:"parent"`
 Cells map[string]Row `json:"cells"`
 Trace []Selected `json:"trace"`
}
type envelope struct{
 State state `json:"state"`
 SHA string `json:"sha256"`
}
func(a *Agent)Snapshot()([]byte,error){
 if !a.valid(){return nil,errors.New("MG23_SOURCE_SNAPSHOT_INVALID")}
 parent,e:=a.Parent.Snapshot();if e!=nil{return nil,e}
 value:=state{Schema, parent,a.Cells,a.Trace}
 canonical,e:=json.Marshal(value);if e!=nil{return nil,e}
 raw,e:=json.Marshal(envelope{value,hash(canonical)})
 if e!=nil||len(raw)>MaxCheckpoint{return nil,errors.New("MG23_CHECKPOINT_CAP")}
 return raw,nil
}
func Restore(blob []byte)(*Agent,error){
 if len(blob)==0||len(blob)>MaxCheckpoint{return nil,errors.New("MG23_RESTORE_SIZE")}
 reader:=json.NewDecoder(bytes.NewReader(blob))
 reader.DisallowUnknownFields()
 var wire envelope
 if e:=reader.Decode(&wire);e!=nil{return nil,e}
 if reader.Decode(new(any))!=io.EOF{return nil,errors.New("MG23_TRAILING_CHECKPOINT")}
 bytes0,e:=json.Marshal(wire.State)
 if e!=nil||wire.State.Schema!=Schema||wire.SHA!=hash(bytes0){
  return nil,errors.New("MG23_INVALID_CHECKPOINT_SHA")
 }
 parent,e:=old.Restore(wire.State.Parent);if e!=nil{return nil,e}
 a:=&Agent{parent,wire.State.Cells,wire.State.Trace}
 if !a.valid(){return nil,errors.New("MG23_INVALID_NATIVE_SELECTED_CREDIT")}
 recoded,e:=a.Snapshot()
 if e!=nil||!bytes.Equal(blob,recoded){return nil,errors.New("MG23_NONCANONICAL_CHECKPOINT")}
 return a,nil
}
