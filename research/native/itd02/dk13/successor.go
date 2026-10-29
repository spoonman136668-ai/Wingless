// Package dk13 learns episode-bounded successor persistence on top of
// unchanged native DK12. Only original public selected before/after is stored.
package dk13
import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "io"
 prior "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk12"
)
const Schema="wingless.itd02.dk13.public-selected-successor-persistence.v1"
const MaxEpisode=192
const MaxContexts=192
const MaxCheckpoint=650000
type Public=prior.Public
type Attempt struct{
 Count uint16 `json:"count"`
 Last string `json:"last"`
}
type Actions [7]Attempt
type Policy struct{
 *prior.Policy
 Attempts map[string]Actions
 Visits map[string]uint16
}
func New()*Policy{return &Policy{Policy:prior.New(),Attempts:map[string]Actions{},Visits:map[string]uint16{}}}
func hash(b []byte)string{s:=sha256.Sum256(b);return hex.EncodeToString(s[:])}
func publicHash(p Public)(string,error){
 if e:=p.Valid();e!=nil{return "",e}
 b,e:=json.Marshal(p);if e!=nil{return "",e}
 return hash(b),nil
}
func digestOK(s string)bool{
 if len(s)!=64{return false}
 for _,b:=range s {if !((b>='0'&&b<='9')||(b>='a'&&b<='f')){return false}}
 return true
}
func (m *Policy) valid()bool{
 if m==nil||m.Policy==nil||m.Attempts==nil||m.Visits==nil||
 len(m.Attempts)>MaxContexts||len(m.Visits)>MaxEpisode||
 m.Policy.Agent.Samples>prior.MaxSamples{return false}
 var attempts,visits uint64
 for key,row:=range m.Attempts{
  if !digestOK(key){return false}
  for _,cell:=range row{
   if cell.Count>0&&!digestOK(cell.Last){return false}
   if cell.Count==0&&cell.Last!=""{return false}
   attempts+=uint64(cell.Count)
  }
 }
 for key,n:=range m.Visits{
  if !digestOK(key)||n==0||n>MaxEpisode{return false}
  visits+=uint64(n)
 }
 var parentEpisode uint64
 for _,n:=range m.Policy.Episode{parentEpisode+=uint64(n)}
 return attempts==visits&&visits==parentEpisode&&visits<=MaxEpisode
}
func (m *Policy) BeginEpisode()error{
 if !m.valid(){return errors.New("DK13_INVALID_SOURCE_EPISODE")}
 if e:=m.Policy.BeginEpisode();e!=nil{return e}
 m.Attempts=map[string]Actions{}
 m.Visits=map[string]uint16{}
 return nil
}
func (m *Policy) Choose(before Public)(*int,string,error){
 if !m.valid(){return nil,"",errors.New("DK13_UNHEALTHY_NATIVE")}
 key,e:=publicHash(before);if e!=nil{return nil,"",e}
 if _,e=m.Policy.Snapshot();e!=nil{return nil,"",e} // pending source fail closed
 row:=m.Attempts[key]
 carrying:=before.Image[3][6][0]>=5&&before.Image[3][6][0]<=7
 best:=-1
 bestScore:=int64(-1<<60)
 for a,s:=range m.Policy.Motion{
  num:=int64(12)*int64(s.SameHeadingChanged)-
    int64(4)*int64(s.HeadingChanged)-int64(8)*int64(s.HeldLost)
  if carrying{num+=int64(4)*int64(s.HeldRetainedChanged)}
  score:=(num*1024)/(int64(s.N)+64)-
    int64(96)*int64(m.Policy.Episode[a])-
    int64(3072)*int64(row[a].Count)
  if row[a].Last!=""{
   score-=int64(1024)*int64(m.Visits[row[a].Last])
  }
  if score>bestScore{best=a;bestScore=score}
 }
 if bestScore>0{return &best,"SELECT_NATIVE_PERSISTENT_PUBLIC_SUCCESSOR",nil}
 // No expected improvement: try actions with fewer actual selected attempts
 // at current exact public state. Raw-public digest determines tie rotation.
 min:=uint16(MaxEpisode+1)
 for _,c:=range row{if c.Count<min{min=c.Count}}
 b,_:=json.Marshal(before);start:=int(sha256.Sum256(b)[0])%7
 for i:=0;i<7;i++{
  a:=(start+i)%7
  if row[a].Count==min{return &a,"SELECT_NATIVE_UNTRIED_PUBLIC_SUCCESSOR",nil}
 }
 return nil,"",errors.New("DK13_NO_LEGAL_NATIVE_ACTION")
}
func (m *Policy) Commit(before Public,action int)error{
 if !m.valid(){return errors.New("DK13_INVALID_COMMIT_SOURCE")}
 return m.Policy.Commit(before,action)
}
func (m *Policy) ObserveFit(before Public,action int,after Public,reward float64,terminal bool)error{
 if !m.valid(){return errors.New("DK13_INVALID_FIT_CUSTODY")}
 return m.Policy.ObserveFit(before,action,after,reward,terminal)
}
func (m *Policy) Feedback(before Public,action int,after Public,reward float64,terminal bool)error{
 if !m.valid()||action<0||action>6{return errors.New("DK13_INVALID_SELECTED_SOURCE")}
 from,e:=publicHash(before);if e!=nil{return e}
 to,e:=publicHash(after);if e!=nil{return e}
 var parentEvents uint64
 for _,n:=range m.Policy.Episode{parentEvents+=uint64(n)}
 if parentEvents>=MaxEpisode{return errors.New("DK13_EPISODE_BUDGET")}
 if _,ok:=m.Attempts[from];!ok&&len(m.Attempts)>=MaxContexts{
  return errors.New("DK13_CONTEXT_BUDGET")
 }
 if _,ok:=m.Visits[to];!ok&&len(m.Visits)>=MaxEpisode{
  return errors.New("DK13_VISIT_BUDGET")
 }
 if e=m.Policy.Feedback(before,action,after,reward,terminal);e!=nil{return e}
 row:=m.Attempts[from]
 cell:=row[action];cell.Count++;cell.Last=to;row[action]=cell
 m.Attempts[from]=row
 m.Visits[to]++
 if !m.valid(){return errors.New("DK13_SELECTED_COUNTER_DRIFT")}
 return nil
}
type state struct{
 Schema string `json:"schema"`
 Parent json.RawMessage `json:"parent"`
 Attempts map[string]Actions `json:"attempts"`
 Visits map[string]uint16 `json:"visits"`
}
type envelope struct{State state `json:"state"`;SHA256 string `json:"sha256"`}
func (m *Policy) Snapshot()([]byte,error){
 if !m.valid(){return nil,errors.New("DK13_INVALID_SNAPSHOT")}
 parent,e:=m.Policy.Snapshot();if e!=nil{return nil,e}
 s:=state{Schema:Schema,Parent:parent,Attempts:m.Attempts,Visits:m.Visits}
 b,e:=json.Marshal(s);if e!=nil{return nil,e}
 raw,e:=json.Marshal(envelope{State:s,SHA256:hash(b)})
 if e!=nil||len(raw)>MaxCheckpoint{return nil,errors.New("DK13_SNAPSHOT_BUDGET")}
 return raw,nil
}
func Restore(raw []byte)(*Policy,error){
 if len(raw)==0||len(raw)>MaxCheckpoint{return nil,errors.New("DK13_RESTORE_SIZE")}
 d:=json.NewDecoder(bytes.NewReader(raw));d.DisallowUnknownFields()
 var wire envelope
 if e:=d.Decode(&wire);e!=nil{return nil,errors.New("DK13_RESTORE_DECODE")}
 if e:=d.Decode(new(any));e!=io.EOF{return nil,errors.New("DK13_RESTORE_TRAILING")}
 b,e:=json.Marshal(wire.State)
 if e!=nil||wire.State.Schema!=Schema||hash(b)!=wire.SHA256{
  return nil,errors.New("DK13_RESTORE_SHA")
 }
 p,e:=prior.Restore(wire.State.Parent);if e!=nil{return nil,e}
 obj:=&Policy{Policy:p,Attempts:wire.State.Attempts,Visits:wire.State.Visits}
 if !obj.valid(){return nil,errors.New("DK13_SOURCE_RECEIPT_TOTALS")}
 again,e:=obj.Snapshot()
 if e!=nil||!bytes.Equal(again,raw){return nil,errors.New("DK13_NONCANONICAL")}
 return obj,nil
}
