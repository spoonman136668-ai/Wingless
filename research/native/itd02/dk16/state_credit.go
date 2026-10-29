// DK16: positive reward credit only for the exact true-source public state.
package dk16
import (
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 prior "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk14"
)
const Schema="wingless.itd02.dk16.exact-state-reward.v1"
const MaxSamples=4096
const MaxContexts=4096
type Public=prior.Public
type Stat struct{Count uint32 `json:"count"`;Positive uint32 `json:"positive"`}
type Actions [7]Stat
type Policy struct{*prior.Policy;States map[string]Actions}
func New()*Policy{return &Policy{Policy:prior.New(),States:map[string]Actions{}}}
func hash(b []byte)string{v:=sha256.Sum256(b);return hex.EncodeToString(v[:])}
func sourceHash(p Public)(string,error){
 if p.Mission!=prior.DoorKey&&p.Mission!=prior.GoalOnly{return "",errors.New("DK16_BAD_MISSION")}
 original:=p;p.Mission=prior.DoorKey
 if e:=p.Valid();e!=nil{return "",e}
 raw,e:=json.Marshal(original);if e!=nil{return "",e}
 return hash(raw),nil
}
func (m *Policy) valid()bool{
 if m==nil||m.Policy==nil||m.States==nil||len(m.States)>MaxContexts||
    m.Policy.Agent.Samples>MaxSamples{return false}
 var total uint64
 for k,rows:=range m.States{
  if len(k)!=64{return false}
  for _,ch:=range k{if !((ch>='0'&&ch<='9')||(ch>='a'&&ch<='f')){return false}}
  for _,v:=range rows{if v.Positive>v.Count{return false};total+=uint64(v.Count)}
 }
 return total==uint64(m.Policy.Agent.Samples)
}
func (m *Policy) BeginEpisode()error{
 if !m.valid()||m.Policy.Pending!=nil{return errors.New("DK16_PENDING")}
 return m.Policy.BeginEpisode()
}
func (m *Policy) Choose(before Public)(*int,string,error){
 if !m.valid()||m.Policy.Pending!=nil{return nil,"",errors.New("DK16_BAD_SOURCE")}
 key,e:=sourceHash(before);if e!=nil{return nil,"",e}
 legacy:=before;legacy.Mission=prior.DoorKey
 baseline,status,e:=m.Policy.Policy.Choose(legacy)
 if e!=nil{return nil,"",e}
 best:=-1;high:=int64(0)
 for a,v:=range m.States[key]{
  if v.Positive==0{continue}
  score:=int64(4096)*int64(v.Positive)/(int64(v.Count)+64)
  if score>high{high=score;best=a}
 }
 if best>=0{return &best,"SELECT_NATIVE_EXACT_PUBLIC_REWARDED_SOURCE",nil}
 return baseline,status,nil
}
func (m *Policy) Commit(p Public,a int)error{
 if !m.valid(){return errors.New("DK16_BAD_COMMIT")}
 if _,e:=sourceHash(p);e!=nil{return e}
 return m.Policy.Commit(p,a)
}
func (m *Policy) learn(b Public,a int,n Public,reward float64,terminal,live bool)error{
 if !m.valid()||a<0||a>6||b.Mission!=n.Mission{return errors.New("DK16_BAD_SELECTED_MISSION")}
 key,e:=sourceHash(b);if e!=nil{return e}
 next,e:=sourceHash(n);if e!=nil{return e}
 if len(m.States)>=MaxContexts&&m.States[key]==(Actions{}){
  return errors.New("DK16_CONTEXT_LIMIT")
 }
 if m.Policy.Agent.Samples>=MaxSamples{return errors.New("DK16_SAMPLE_LIMIT")}
 if live{e=m.Policy.Feedback(b,a,n,reward,terminal)
 }else{e=m.Policy.ObserveFit(b,a,n,reward,terminal)}
 if e!=nil{return e}
 row:=m.States[key];v:=row[a];v.Count++
 if reward>0&&terminal&&key!=next{v.Positive++}
 row[a]=v;m.States[key]=row
 if !m.valid(){return errors.New("DK16_PARENT_COUNT_DRIFT")}
 return nil
}
func (m *Policy) ObserveFit(b Public,a int,n Public,r float64,t bool)error{
 return m.learn(b,a,n,r,t,false)
}
func (m *Policy) Feedback(b Public,a int,n Public,r float64,t bool)error{
 return m.learn(b,a,n,r,t,true)
}
