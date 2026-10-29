// Package mg28 isolates opt-in native selected original-public next-change signal.
// It wraps the unchanged MG23 learner and never synthesizes a source outcome.
package mg28
import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "io"
 old "github.com/spoonman136668-ai/Wingless/research/native/itd02/mg23"
)
const Schema="wingless.itd02.mg28.selected-public-next-change.v1"
const MaxSamples=4096
const MaxContexts=4096
const MaxCheckpoint=6500000
type Public=old.Public
type Stat struct{Count uint32 `json:"count"`;Changed uint32 `json:"changed"`}
type Row [7]Stat
type Agent struct{
 Parent *old.Agent
 Cells map[string]Row
 Enabled bool
}
func New(enabled bool)*Agent{return &Agent{Parent:old.New(),Cells:map[string]Row{},Enabled:enabled}}
func hash(b []byte)string{v:=sha256.Sum256(b);return hex.EncodeToString(v[:])}
func key(p Public)(string,error){
 if len(p.Image)!=7||p.Direction<0||p.Direction>3||len(p.Mission)==0{return "",errors.New("MG28_NOT_ORIGINAL_PUBLIC")}
 for _,r:=range p.Image{
  if len(r)!=7{return "",errors.New("MG28_BAD_PUBLIC_ROW")}
  for _,c:=range r{if len(c)!=3{return "",errors.New("MG28_BAD_PUBLIC_PIXEL")}}
 }
 point:=p.Image[3][5]
 record:=struct{
  Mission string `json:"mission"`
  Direction int `json:"direction"`
  Front [3]int `json:"front"`
 }{p.Mission,p.Direction,[3]int{point[0],point[1],point[2]}}
 raw,e:=json.Marshal(record)
 if e!=nil{return "",e}
 return hash(raw),nil
}
func changed(b,n Public)(bool,error){
 a,e:=json.Marshal(b);if e!=nil{return false,e}
 c,e:=json.Marshal(n);if e!=nil{return false,e}
 return !bytes.Equal(a,c),nil
}
func samples(a *Agent)uint64{return a.Parent.Parent.Parent.Agent.Samples}
func(a *Agent)valid()bool{
 if a==nil||a.Parent==nil||a.Cells==nil||len(a.Cells)>MaxContexts{return false}
 if samples(a)>MaxSamples{return false}
 var count uint64
 for k,row:=range a.Cells{
  if len(k)!=64{return false}
  for _,v:=range k{if !((v>='0'&&v<='9')||(v>='a'&&v<='f')){return false}}
  for _,s:=range row{
   if s.Changed>s.Count{return false}
   count+=uint64(s.Count)
  }
 }
 return count==samples(a)
}
func(a *Agent)BeginEpisode()error{
 if !a.valid(){return errors.New("MG28_BAD_EPISODE")}
 return a.Parent.BeginEpisode()
}
func(a *Agent)Choose(b Public)(*int,string,error){
 if !a.valid(){return nil,"",errors.New("MG28_BAD_SOURCE_CHOOSE")}
 k,e:=key(b);if e!=nil{return nil,"",e}
 base,status,e:=a.Parent.Choose(b)
 if e!=nil||base==nil||*base<0||*base>6{return nil,"",errors.New("MG28_MISSING_NATIVE_PARENT_SOURCE")}
 if !a.Enabled{return base,status,nil}
 row:=a.Cells[k]
 baseline:=row[*base]
 baselineScore:=uint64(10000)*(uint64(baseline.Changed)+1)/(uint64(baseline.Count)+2)
 best:=-1;highest:=uint64(0)
 for action,stat:=range row{
  if stat.Count<3{continue}
  score:=uint64(10000)*(uint64(stat.Changed)+1)/(uint64(stat.Count)+2)
  if score<7500||score<baselineScore+1000{continue}
  if score>highest{highest=score;best=action}
 }
 if best>=0&&best!=*base{
  result:=best
  return &result,"MG28_NATIVE_OWN_SELECTED_REAL_NEXT_PUBLIC_CHANGE",nil
 }
 return base,status,nil
}
func(a *Agent)Commit(b Public,action int)error{
 if !a.valid(){return errors.New("MG28_COMMIT_INVALID_SOURCE_STATE")}
 if _,e:=key(b);e!=nil{return e}
 return a.Parent.Commit(b,action)
}
func(a *Agent)learn(b Public,action int,n Public,reward float64,terminal,live bool)error{
 if !a.valid()||action<0||action>6||b.Mission!=n.Mission{
  return errors.New("MG28_UNTRUSTED_SELECTED_SOURCE")
 }
 k,e:=key(b);if e!=nil{return e}
 if _,e=key(n);e!=nil{return e}
 delta,e:=changed(b,n);if e!=nil{return e}
 if samples(a)>=MaxSamples{return errors.New("MG28_SOURCE_SAMPLE_CAP")}
 if _,ok:=a.Cells[k];!ok&&len(a.Cells)>=MaxContexts{return errors.New("MG28_CONTEXT_CAP")}
 if live{e=a.Parent.Feedback(b,action,n,reward,terminal)
 }else{e=a.Parent.ObserveFit(b,action,n,reward,terminal)}
 if e!=nil{return e}
 row:=a.Cells[k];stat:=row[action];stat.Count++
 if delta{stat.Changed++}
 row[action]=stat;a.Cells[k]=row
 if !a.valid(){return errors.New("MG28_SOURCE_COUNT_CONSERVATION")}
 return nil
}
func(a *Agent)Feedback(b Public,action int,n Public,r float64,t bool)error{
 return a.learn(b,action,n,r,t,true)
}
func(a *Agent)ObserveFit(b Public,action int,n Public,r float64,t bool)error{
 return a.learn(b,action,n,r,t,false)
}
type state struct{
 Schema string `json:"schema"`
 Parent json.RawMessage `json:"parent"`
 Cells map[string]Row `json:"cells"`
 Enabled bool `json:"enabled"`
}
type envelope struct{State state `json:"state"`;SHA string `json:"sha256"`}
func(a *Agent)Snapshot()([]byte,error){
 if !a.valid(){return nil,errors.New("MG28_INVALID_SNAPSHOT")}
 parent,e:=a.Parent.Snapshot();if e!=nil{return nil,e}
 value:=state{Schema,parent,a.Cells,a.Enabled}
 raw,e:=json.Marshal(value);if e!=nil{return nil,e}
 wire,e:=json.Marshal(envelope{value,hash(raw)})
 if e!=nil||len(wire)>MaxCheckpoint{return nil,errors.New("MG28_SNAPSHOT_CAP")}
 return wire,nil
}
func Restore(raw []byte)(*Agent,error){
 if len(raw)==0||len(raw)>MaxCheckpoint{return nil,errors.New("MG28_RESTORE_SIZE")}
 decoder:=json.NewDecoder(bytes.NewReader(raw));decoder.DisallowUnknownFields()
 var wire envelope
 if e:=decoder.Decode(&wire);e!=nil{return nil,e}
 if decoder.Decode(new(any))!=io.EOF{return nil,errors.New("MG28_RESTORE_TRAILING")}
 stateRaw,e:=json.Marshal(wire.State)
 if e!=nil||wire.State.Schema!=Schema||hash(stateRaw)!=wire.SHA{return nil,errors.New("MG28_RESTORE_DIGEST")}
 parent,e:=old.Restore(wire.State.Parent);if e!=nil{return nil,e}
 a:=&Agent{Parent:parent,Cells:wire.State.Cells,Enabled:wire.State.Enabled}
 if !a.valid(){return nil,errors.New("MG28_UNTRUSTED_SOURCE_STATE")}
 rebuild,e:=a.Snapshot()
 if e!=nil||!bytes.Equal(raw,rebuild){return nil,errors.New("MG28_NONCANONICAL_RESTORE")}
 return a,nil
}
