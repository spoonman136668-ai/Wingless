// Package dk14 provides an opt-in seven-action original public mission boundary.
// It is NOT DK13 native policy training or an actor. Existing DK13 is unchanged.
package dk14

import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "io"
 "math"
 "unicode"
 "unicode/utf8"
)

const Schema="wingless.itd02.dk14.original-public-task-neutral-boundary.v1"
const MaxSteps=192
const MaxCheckpoint=2000000
const LegacyMission="use the key to open the door and then get to the goal"

type Public struct {
 Image [][][]int `json:"image"`
 Direction int `json:"direction"`
 Mission string `json:"mission"`
}
func (p Public) Valid()error{
 if p.Direction<0||p.Direction>3{return errors.New("DK14_PUBLIC_DIRECTION")}
 if !utf8.ValidString(p.Mission)||utf8.RuneCountInString(p.Mission)<1||utf8.RuneCountInString(p.Mission)>256{
  return errors.New("DK14_PUBLIC_MISSION_UNTRUSTED")
 }
 for _,r:=range p.Mission{if unicode.IsControl(r)||unicode.Is(unicode.Cs,r){
  return errors.New("DK14_PUBLIC_MISSION_UNTRUSTED")
 }}
 if len(p.Image)!=7{return errors.New("DK14_PUBLIC_IMAGE_SHAPE")}
 for _,row:=range p.Image{
  if len(row)!=7{return errors.New("DK14_PUBLIC_IMAGE_ROW")}
  for _,pixel:=range row{
   if len(pixel)!=3{return errors.New("DK14_PUBLIC_IMAGE_PIXEL")}
   k,c,state:=pixel[0],pixel[1],pixel[2]
   if k<0||k>10||c<0||c>5||state<0||state>255||(k==4&&state>2){
    return errors.New("DK14_PUBLIC_IMAGE_RANGE")
   }
  }
 }
 return nil
}
// Strictly reject source-private fields when reading JSON transport.
// Direct typed Public inputs cannot contain any other fields.
func ReadPublic(raw []byte)(Public,error){
 var p Public
 if len(raw)==0||len(raw)>10000{return p,errors.New("DK14_PUBLIC_FRAME_SIZE")}
 d:=json.NewDecoder(bytes.NewReader(raw))
 d.DisallowUnknownFields()
 if err:=d.Decode(&p);err!=nil{return Public{},errors.New("DK14_PRIVILEGED_OR_INVALID_PUBLIC_FRAME")}
 var extra any
 if e:=d.Decode(&extra);e!=io.EOF{return Public{},errors.New("DK14_EXTRA_PUBLIC_FRAME")}
 if err:=p.Valid();err!=nil{return Public{},err}
 return p,nil
}
func sha(raw []byte)string{h:=sha256.Sum256(raw);return hex.EncodeToString(h[:])}
func canonical(v any)([]byte,error){return json.Marshal(v)}
func Digest(p Public)(string,error){
 if err:=p.Valid();err!=nil{return "",err}
 // Keys fixed in lexical order to match Python canonical original-public hash.
 raw,e:=canonical(map[string]any{"image":p.Image,"direction":p.Direction,"mission":p.Mission})
 if e!=nil{return "",e}
 return sha(raw),nil
}
type Receipt struct {
 Index int `json:"index"`
 Before string `json:"before"`
 Action int `json:"action"`
 After string `json:"after"`
 Reward float64 `json:"reward"`
 Terminated bool `json:"terminated"`
 Truncated bool `json:"truncated"`
}
type Event struct {
 Receipt Receipt `json:"receipt"`
 PublicAfter Public `json:"public_after"`
}
type Boundary struct{
 Start *Public `json:"start"`
 Current *Public `json:"current"`
 Events []Event `json:"events"`
 Pending *int `json:"pending"`
 Finished bool `json:"finished"`
}
func New()*Boundary{return &Boundary{Events:[]Event{}}}
func clone(p Public)Public{
 x:=Public{Direction:p.Direction,Mission:p.Mission,Image:make([][][]int,len(p.Image))}
 for i,row:=range p.Image{
  x.Image[i]=make([][]int,len(row))
  for j,pixel:=range row {x.Image[i][j]=append([]int(nil),pixel...)}
 }
 return x
}
func (b *Boundary) Begin(p Public)error{
 if b.Current!=nil{return errors.New("DK14_EPISODE_ALREADY_ACTIVE")}
 if e:=p.Valid();e!=nil{return e}
 cp:=clone(p);start:=clone(p)
 b.Start=&start;b.Current=&cp
 return nil
}
func (b *Boundary) Commit(a int)(int,error){
 if b.Current==nil||b.Pending!=nil||b.Finished{return 0,errors.New("DK14_NO_PUBLIC_ACTION_CAPABILITY")}
 if a<0||a>6{return 0,errors.New("DK14_BAD_SELECTED_SOURCE_ACTION")}
 if len(b.Events)>=MaxSteps{return 0,errors.New("DK14_EPISODE_BUDGET")}
 b.Pending=&a
 return a,nil
}
func (b *Boundary) Feedback(a int,after Public,reward float64,terminated,truncated bool)(Receipt,error){
 if b.Current==nil||b.Pending==nil||b.Finished{return Receipt{},errors.New("DK14_FEEDBACK_NOT_AUTHORIZED")}
 if a!=*b.Pending{return Receipt{},errors.New("DK14_WRONG_SELECTED_ACTION")}
 if e:=after.Valid();e!=nil{return Receipt{},e}
 if after.Mission!=b.Current.Mission{return Receipt{},errors.New("DK14_ORIGINAL_PUBLIC_MISSION_CHANGED")}
 if math.IsNaN(reward)||math.IsInf(reward,0)||reward<0||reward>1{
  return Receipt{},errors.New("DK14_UNTRUSTED_SOURCE_REWARD")
 }
 old,e:=Digest(*b.Current);if e!=nil{return Receipt{},e}
 next,e:=Digest(after);if e!=nil{return Receipt{},e}
 receipt:=Receipt{len(b.Events)+1,old,a,next,reward,terminated,truncated}
 b.Events=append(b.Events,Event{Receipt:receipt,PublicAfter:clone(after)})
 current:=clone(after)
 b.Current=&current
 b.Pending=nil
 b.Finished=terminated||truncated
 return receipt,nil
}
type State struct {
 Schema string `json:"schema"`
 Start *Public `json:"start"`
 Current *Public `json:"current"`
 Events []Event `json:"events"`
 Pending *int `json:"pending"`
 Finished bool `json:"finished"`
}
type Envelope struct{
 State State `json:"state"`
 Sha256 string `json:"sha256"`
}
func (b *Boundary) state()State{
 return State{Schema, b.Start,b.Current,b.Events,b.Pending,b.Finished}
}
func (b *Boundary) Snapshot()([]byte,error){
 state:=b.state()
 raw,e:=canonical(state);if e!=nil{return nil,e}
 buf,e:=canonical(Envelope{state,sha(raw)});if e!=nil{return nil,e}
 if len(buf)>MaxCheckpoint{return nil,errors.New("DK14_CHECKPOINT_TOO_LARGE")}
 return buf,nil
}
func Restore(blob []byte)(*Boundary,error){
 if len(blob)==0||len(blob)>MaxCheckpoint{return nil,errors.New("DK14_CHECKPOINT_UNTRUSTED_BYTES")}
 var env Envelope
 d:=json.NewDecoder(bytes.NewReader(blob));d.DisallowUnknownFields()
 if e:=d.Decode(&env);e!=nil{return nil,errors.New("DK14_CHECKPOINT_INVALID")}
 var extra any
 if e:=d.Decode(&extra);e!=io.EOF{return nil,errors.New("DK14_CHECKPOINT_EXTRA_JSON")}
 raw,e:=canonical(env.State)
 if e!=nil||sha(raw)!=env.Sha256||env.State.Schema!=Schema{
  return nil,errors.New("DK14_CHECKPOINT_HASH_OR_SCHEMA")
 }
 original:=env.State
 b:=New()
 if original.Start!=nil{
  if e:=b.Begin(*original.Start);e!=nil{return nil,e}
  if len(original.Events)>MaxSteps{return nil,errors.New("DK14_CHECKPOINT_TOO_MANY_EVENTS")}
  for _,event:=range original.Events{
   r:=event.Receipt
   if _,e:=b.Commit(r.Action);e!=nil{return nil,e}
   replay,e:=b.Feedback(r.Action,event.PublicAfter,r.Reward,r.Terminated,r.Truncated)
   if e!=nil||replay!=r{return nil,errors.New("DK14_CHECKPOINT_REPLAY_DRIFT")}
  }
  if original.Pending!=nil{
   if _,e:=b.Commit(*original.Pending);e!=nil{return nil,e}
  }
 } else if original.Current!=nil||len(original.Events)>0||original.Pending!=nil||original.Finished{
  return nil,errors.New("DK14_CHECKPOINT_INVALID_EMPTY_STATE")
 }
 if (b.Current==nil)!=(original.Current==nil)||b.Finished!=original.Finished{
  return nil,errors.New("DK14_CHECKPOINT_CURRENT_DRIFT")
 }
 if b.Current!=nil{
  expected,e:=Digest(*original.Current);if e!=nil{return nil,e}
  actual,e:=Digest(*b.Current);if e!=nil||actual!=expected{return nil,errors.New("DK14_CHECKPOINT_PUBLIC_DRIFT")}
 }
 if (b.Pending==nil)!=(original.Pending==nil){return nil,errors.New("DK14_CHECKPOINT_PENDING_DRIFT")}
 if b.Pending!=nil&&*b.Pending!=*original.Pending{return nil,errors.New("DK14_CHECKPOINT_ACTION_DRIFT")}
 again,e:=b.Snapshot()
 if e!=nil||!bytes.Equal(blob,again){return nil,errors.New("DK14_CHECKPOINT_NONCANONICAL")}
 return b,nil
}
