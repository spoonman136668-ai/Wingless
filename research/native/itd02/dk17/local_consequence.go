// Package dk17 learns whether independently selected source actions actually
// change the scene from a local public configuration. It does not prescribe
// action codes or infer hidden objective/position.
package dk17
import(
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "io"
 prior "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk16"
)
const Schema="wingless.itd02.dk17.selected-local-public-consequence.v1"
const MaxSamples=4096
const MaxContexts=4096
const MaxCheckpoint=4000000
type Public=prior.Public
type Stat struct{N uint32 `json:"n"`;Changed uint32 `json:"changed"`;NoEffect uint32 `json:"no_effect"`}
type Actions [7]Stat
type Agent struct{*prior.Policy;Local map[string]Actions}
func New()*Agent{return &Agent{Policy:prior.New(),Local:make(map[string]Actions)}}
func checksum(raw []byte)string{s:=sha256.Sum256(raw);return hex.EncodeToString(s[:])}
func context(p Public)(string,error){
 if p.Mission!="use the key to open the door and then get to the goal"&&p.Mission!="get to the green goal square"{
  return "",errors.New("DK17_UNREGISTERED_ORIGINAL_MISSION")
 }
 copy:=p;copy.Mission="use the key to open the door and then get to the goal"
 if e:=copy.Valid();e!=nil{return "",e}
 key:=struct{
  Mission string `json:"mission"`
  Visible [3][3]int `json:"visible"`
 }{Mission:p.Mission}
 for i,col:=range [3]int{4,5,6}{
  for k:=0;k<3;k++{key.Visible[i][k]=p.Image[3][col][k]}
 }
 raw,e:=json.Marshal(key);if e!=nil{return "",e}
 return checksum(raw),nil
}
func sameImage(a,b Public)bool{
 if len(a.Image)!=len(b.Image){return false}
 for i:=range a.Image{
  if len(a.Image[i])!=len(b.Image[i]){return false}
  for j:=range a.Image[i]{
   if len(a.Image[i][j])!=len(b.Image[i][j]){return false}
   for k:=range a.Image[i][j]{
    if a.Image[i][j][k]!=b.Image[i][j][k]{return false}
   }
  }
 }
 return true
}
func (m *Agent) valid()bool{
 if m==nil||m.Policy==nil||m.Local==nil||len(m.Local)>MaxContexts||
    m.Policy.Agent.Samples>MaxSamples{return false}
 var total uint64
 for key,rows:=range m.Local{
  if len(key)!=64{return false}
  for _,c:=range key{if !((c>='0'&&c<='9')||(c>='a'&&c<='f')){return false}}
  for _,r:=range rows{
   if r.Changed>r.N||r.NoEffect>r.N||r.Changed+r.NoEffect>r.N{return false}
   total+=uint64(r.N)
  }
 }
 return total==uint64(m.Policy.Agent.Samples)
}
func (m *Agent) BeginEpisode()error{
 if !m.valid(){return errors.New("DK17_SOURCE_DRIFT")}
 return m.Policy.BeginEpisode()
}
func score(s Stat)int64{
 return int64(4096)*(int64(s.Changed)-int64(s.NoEffect))/(int64(s.N)+64)
}
func (m *Agent) Choose(p Public)(*int,string,error){
 if !m.valid(){return nil,"",errors.New("DK17_SOURCE_INVALID")}
 key,e:=context(p);if e!=nil{return nil,"",e}
 baseline,reason,e:=m.Policy.Choose(p)
 if e!=nil||baseline==nil{return baseline,reason,e}
 rows:=m.Local[key]
 baseScore:=score(rows[*baseline])
 best:=-1;bestScore:=baseScore
 for action,r:=range rows{
  if r.Changed<2{continue}
  v:=score(r)
  if v>bestScore{best=action;bestScore=v}
 }
 if best>=0{
  return &best,"SELECT_NATIVE_LOCAL_SELECTED_CAUSAL_EFFECT",nil
 }
 return baseline,reason,nil
}
func (m *Agent) Commit(p Public,a int)error{
 if !m.valid(){return errors.New("DK17_COMMIT_DRIFT")}
 if _,e:=context(p);e!=nil{return e}
 return m.Policy.Commit(p,a)
}
func (m *Agent) learn(b Public,a int,n Public,reward float64,terminal,live bool)error{
 if !m.valid()||a<0||a>6||b.Mission!=n.Mission{
  return errors.New("DK17_BAD_SELECTED_SOURCE")
 }
 key,e:=context(b);if e!=nil{return e}
 if _,e=context(n);e!=nil{return e}
 if m.Policy.Agent.Samples>=MaxSamples{return errors.New("DK17_SOURCE_BUDGET")}
 if _,exists:=m.Local[key];!exists&&len(m.Local)>=MaxContexts{
  return errors.New("DK17_LOCAL_CONTEXT_BUDGET")
 }
 if live{e=m.Policy.Feedback(b,a,n,reward,terminal)
 }else{e=m.Policy.ObserveFit(b,a,n,reward,terminal)}
 if e!=nil{return e}
 row:=m.Local[key]
 stat:=row[a];stat.N++
 identical:=sameImage(b,n)
 if b.Direction==n.Direction&&!identical{stat.Changed++}
 if b.Direction==n.Direction&&identical{stat.NoEffect++}
 row[a]=stat;m.Local[key]=row
 if !m.valid(){return errors.New("DK17_SELECTED_SOURCE_COUNTER_DRIFT")}
 return nil
}
func (m *Agent) ObserveFit(b Public,a int,n Public,r float64,t bool)error{
 return m.learn(b,a,n,r,t,false)
}
func (m *Agent) Feedback(b Public,a int,n Public,r float64,t bool)error{
 return m.learn(b,a,n,r,t,true)
}
type checkpoint struct{
 Schema string `json:"schema"`
 Parent json.RawMessage `json:"parent"`
 Local map[string]Actions `json:"local"`
}
type envelope struct{
 State checkpoint `json:"state"`
 SHA256 string `json:"sha256"`
}
func (m *Agent) Snapshot()([]byte,error){
 if !m.valid(){return nil,errors.New("DK17_SNAPSHOT_SOURCE_DRIFT")}
 parent,e:=m.Policy.Snapshot();if e!=nil{return nil,e}
 state:=checkpoint{Schema:Schema,Parent:parent,Local:m.Local}
 raw,e:=json.Marshal(state);if e!=nil{return nil,e}
 full,e:=json.Marshal(envelope{State:state,SHA256:checksum(raw)})
 if e!=nil||len(full)>MaxCheckpoint{return nil,errors.New("DK17_SNAPSHOT_SIZE")}
 return full,nil
}
func Restore(raw []byte)(*Agent,error){
 if len(raw)==0||len(raw)>MaxCheckpoint{return nil,errors.New("DK17_RESTORE_SIZE")}
 dec:=json.NewDecoder(bytes.NewReader(raw));dec.DisallowUnknownFields()
 var wire envelope
 if e:=dec.Decode(&wire);e!=nil{return nil,e}
 if e:=dec.Decode(new(any));e!=io.EOF{return nil,errors.New("DK17_TRAILING_BYTES")}
 nested,e:=json.Marshal(wire.State)
 if e!=nil||wire.State.Schema!=Schema||checksum(nested)!=wire.SHA256{
  return nil,errors.New("DK17_CHECKPOINT_HASH_MISMATCH")
 }
 parent,e:=prior.Restore(wire.State.Parent);if e!=nil{return nil,e}
 model:=&Agent{Policy:parent,Local:wire.State.Local}
 if !model.valid(){return nil,errors.New("DK17_CORRUPT_SOURCE_COUNTERS")}
 again,e:=model.Snapshot()
 if e!=nil||!bytes.Equal(again,raw){return nil,errors.New("DK17_NONCANONICAL_RESTORE")}
 return model,nil
}
