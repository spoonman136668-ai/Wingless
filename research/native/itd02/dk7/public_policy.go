// Package dk7 is a new opt-in, source-only native Go public seven-action
// selected-consequence learner. This is NOT a port/replacement of original R8.
// Original Wingless Go R8/R5, Go source and accepted memory stay unchanged.
package dk7

import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "math"
 "strconv"
)

const Schema="wingless.itd02.dk7.public-seven-action-native-go-sourceonly.v1"
const Mission="use the key to open the door and then get to the goal"
const MaxSamples=4096
const MaxContexts=512
const MaxCheckpoint=150000

type Public struct {
 Image [][][]int `json:"image"`
 Direction int `json:"direction"`
 Mission string `json:"mission"`
}
func (p Public) Valid()error {
 if p.Direction<0 || p.Direction>3 || p.Mission!=Mission || len(p.Image)!=7 {return errors.New("DK7_PUBLIC_MISSION_OR_IMAGE_OR_HEADING")}
 for x:=0;x<7;x++ {
  if len(p.Image[x])!=7{return errors.New("DK7_PUBLIC_IMAGE_SHAPE")}
  for y:=0;y<7;y++ {
   pixel:=p.Image[x][y]
   if len(pixel)!=3{return errors.New("DK7_PUBLIC_PIXEL_SIZE")}
   if pixel[0]<0||pixel[0]>10||pixel[1]<0||pixel[1]>5||pixel[2]<0||pixel[2]>255{return errors.New("DK7_PUBLIC_PIXEL_RANGE")}
   if pixel[0]==4 && pixel[2]>2{return errors.New("DK7_PUBLIC_DOOR_STATE")}
  }
 }
 return nil
}
func Context(p Public)(string,error) {
 if e:=p.Valid();e!=nil{return "",e}
 f:=p.Image[3][5]
 held:=p.Image[3][6]
 kind,color:=0,0
 if held[0]==5||held[0]==6||held[0]==7{kind,color=held[0],held[1]}
 return fmt.Sprintf("%d:%d:%d:%d:%d:%d",f[0],f[1],f[2],kind,color,p.Direction),nil
}
type Bin struct {Count uint32 `json:"count"`; Positive uint32 `json:"positive"`}
type Cells [7]Bin
type Pending struct {BeforeSHA string `json:"before_sha"`; Action int `json:"action"`}
type Policy struct {
 Counts map[string]Cells
 Samples uint64
 Live uint64
 Pending *Pending
}
func New()*Policy{return &Policy{Counts:map[string]Cells{}}}
func canonical(p Public)([]byte,error){if e:=p.Valid();e!=nil{return nil,e};return json.Marshal(p)}
func hash(b []byte)string{s:=sha256.Sum256(b);return hex.EncodeToString(s[:])}
func (m *Policy) healthy()bool{
 if m==nil||m.Counts==nil||m.Samples>MaxSamples||m.Live>m.Samples||len(m.Counts)>MaxContexts{return false}
 var count uint64
 for _,cells:=range m.Counts{
  for _,bin:=range cells{
   if bin.Positive>bin.Count{return false}
   count+=uint64(bin.Count)
  }
 }
 if count!=m.Samples{return false}
 if m.Pending!=nil && (m.Pending.Action<0||m.Pending.Action>6||len(m.Pending.BeforeSHA)!=64){return false}
 return true
}
func (m *Policy) BeginEpisode()error{
 if !m.healthy()||m.Pending!=nil{return errors.New("DK7_EPISODE_WITH_PENDING")}
 return nil
}
func (m *Policy) Choose(before Public)(*int,string,error){
 if !m.healthy()||m.Pending!=nil{return nil,"",errors.New("DK7_UNHEALTHY_OR_PENDING")}
 key,e:=Context(before);if e!=nil{return nil,"",e}
 cells:=m.Counts[key]
 choice:=-1
 best:=-1.0
 for a:=0;a<7;a++{
  b:=cells[a]
  if b.Count>0&&b.Positive>0{
   confidence:=float64(b.Positive)/float64(b.Count)
   if confidence>best{choice=a;best=confidence}
  }
 }
 if choice<0{return nil,"ABSTAIN_NO_OBSERVED_POSITIVE_ACTION",nil}
 return &choice,"SELECT_LEARNED_POSITIVE_ACTION",nil
}
func (m *Policy) Commit(before Public,action int)error{
 if !m.healthy()||m.Pending!=nil||action<0||action>6{return errors.New("DK7_COMMIT_ILLEGAL_OR_PENDING")}
 raw,e:=canonical(before);if e!=nil{return e}
 m.Pending=&Pending{BeforeSHA:hash(raw),Action:action}
 return nil
}
func (m *Policy) ObserveFit(before Public,action int,after Public,reward float64,terminal bool)error{
 if !m.healthy()||m.Pending!=nil{return errors.New("DK7_FIT_PENDING_OR_CORRUPT")}
 return m.observe(before,action,after,reward,terminal,false)
}
func (m *Policy) Feedback(before Public,action int,after Public,reward float64,terminal bool)error{
 if !m.healthy()||m.Pending==nil{return errors.New("DK7_NO_SELECTED_ACTION")}
 raw,e:=canonical(before)
 if e!=nil||action!=m.Pending.Action||hash(raw)!=m.Pending.BeforeSHA{
  return errors.New("DK7_SELECTED_PUBLIC_FEEDBACK_MISMATCH")
 }
 if e=m.observe(before,action,after,reward,terminal,true);e!=nil{return e}
 m.Pending=nil
 return nil
}
func (m *Policy) observe(before Public,action int,after Public,reward float64,terminal bool,live bool)error{
 if !m.healthy()||action<0||action>6||math.IsNaN(reward)||math.IsInf(reward,0)||reward<0||reward>1||m.Samples>=MaxSamples{return errors.New("DK7_SELECTED_FEEDBACK_INVALID")}
 key,e:=Context(before);if e!=nil{return e}
 if e=after.Valid();e!=nil{return e}
 row,existing:=m.Counts[key]
 if !existing&&len(m.Counts)>=MaxContexts{return errors.New("DK7_CONTEXT_LIMIT")}
 if row[action].Count>=MaxSamples{return errors.New("DK7_SELECTED_ACTION_OVERFLOW")}
 b:=row[action];b.Count++;if reward>0{b.Positive++};row[action]=b
 m.Counts[key]=row
 m.Samples++
 if live{m.Live++}
 if !m.healthy(){return errors.New("DK7_POSTUPDATE_CONSERVATION")}
 return nil
}
type snapshot struct {
 Schema string `json:"schema"`
 Counts map[string]Cells `json:"counts"`
 Samples uint64 `json:"samples"`
 Live uint64 `json:"live"`
}
type envelope struct {Payload snapshot `json:"payload"`; SHA256 string `json:"sha256"`}
func (m *Policy) Snapshot()([]byte,error){
 if !m.healthy()||m.Pending!=nil{return nil,errors.New("DK7_CHECKPOINT_PENDING_OR_INVALID")}
 state:=snapshot{Schema:Schema,Counts:m.Counts,Samples:m.Samples,Live:m.Live}
 raw,e:=json.Marshal(state);if e!=nil{return nil,e}
 outer,e:=json.Marshal(envelope{Payload:state,SHA256:hash(raw)})
 if e!=nil||len(outer)>MaxCheckpoint{return nil,errors.New("DK7_CHECKPOINT_BUDGET")}
 return outer,nil
}
func Restore(raw []byte)(*Policy,error){
 if len(raw)==0||len(raw)>MaxCheckpoint{return nil,errors.New("DK7_CHECKPOINT_SIZE")}
 decoder:=json.NewDecoder(bytes.NewReader(raw))
 decoder.DisallowUnknownFields()
 var wire envelope
 if e:=decoder.Decode(&wire);e!=nil{return nil,e}
 if e:=decoder.Decode(new(any));e!=io.EOF{return nil,errors.New("DK7_TRAILING_JSON")}
 if wire.Payload.Schema!=Schema{return nil,errors.New("DK7_SCHEMA")}
 inner,e:=json.Marshal(wire.Payload);if e!=nil||hash(inner)!=wire.SHA256{return nil,errors.New("DK7_CHECKPOINT_DIGEST")}
 obj:=&Policy{Counts:wire.Payload.Counts,Samples:wire.Payload.Samples,Live:wire.Payload.Live}
 if !obj.healthy(){return nil,errors.New("DK7_CHECKPOINT_COUNTERS")}
 canon,e:=obj.Snapshot()
 if e!=nil||!bytes.Equal(raw,canon){return nil,errors.New("DK7_CHECKPOINT_NONCANONICAL")}
 return obj,nil
}
func Status(m *Policy)map[string]string{
 return map[string]string{"updates":strconv.FormatUint(m.Samples,10),"live":strconv.FormatUint(m.Live,10)}
}
