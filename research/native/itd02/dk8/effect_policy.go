// Package dk8 is an opt-in native Go learned visible-effect policy for
// original public DoorKey seven-action transitions. It is NOT Wingless R8.
// No external effect label, seed, hidden map or alternative reward is read.
package dk8

import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "math"
 "strings"
 prior "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk7"
)
const Schema="wingless.itd02.dk8.visible-public-consequence-sourceonly.v1"
const MaxSamples=4096
const MaxContexts=512
const MaxCheckpoint=150000
const Carrying="carrying_changed"
const Door="visible_door_state_changed"
type Public=prior.Public
type Cell struct{
 N uint32 `json:"n"`
 Carry uint32 `json:"carry"`
 Door uint32 `json:"door"`
}
type Row [7]Cell
type Pending struct{ BeforeSHA string `json:"before_sha"`; Action int `json:"action"` }
type Policy struct{
 Counts map[string]Row
 Samples uint64
 Live uint64
 Pending *Pending
}
func New()*Policy{return &Policy{Counts:map[string]Row{}}}
func digest(data []byte)string{s:=sha256.Sum256(data);return hex.EncodeToString(s[:])}
func canon(p Public)([]byte,error){if e:=p.Valid();e!=nil{return nil,e};return json.Marshal(p)}
func Context(p Public)(string,error){
 key,e:=prior.Context(p);if e!=nil{return "",e}
 i:=strings.LastIndexByte(key,':')
 if i<=0{return "",errors.New("DK8_INVALID_ORIGINAL_PUBLIC_CONTEXT")}
 return key[:i],nil // R11 context WITHOUT direction; action is separate row index.
}
func carried(p Public)[2]int{
 v:=p.Image[3][6]
 if v[0]==5||v[0]==6||v[0]==7{return [2]int{v[0],v[1]}}
 return [2]int{0,0}
}
func observedEffects(before,after Public)(bool,bool,error){
 if e:=before.Valid();e!=nil{return false,false,e}
 if e:=after.Valid();e!=nil{return false,false,e}
 b:=before.Image[3][5];a:=after.Image[3][5]
 return carried(before)!=carried(after),b[0]==4&&a[0]==4&&b[2]!=a[2],nil
}
func (p *Policy) healthy()bool{
 if p==nil||p.Counts==nil||p.Samples>MaxSamples||p.Live>p.Samples||len(p.Counts)>MaxContexts{return false}
 var n uint64
 for key,row:=range p.Counts{
  if len(strings.Split(key,":"))!=5{return false}
  for _,c:=range row{
   if c.Carry>c.N||c.Door>c.N{return false}
   n+=uint64(c.N)
  }
 }
 return n==p.Samples&&(p.Pending==nil||(len(p.Pending.BeforeSHA)==64&&p.Pending.Action>=0&&p.Pending.Action<7))
}
func (p *Policy) BeginEpisode()error{
 if !p.healthy()||p.Pending!=nil{return errors.New("DK8_EPISODE_PENDING_OR_CORRUPT")}
 return nil
}
func (p *Policy) Choose(before Public,objective string)(*int,string,error){
 if !p.healthy()||p.Pending!=nil{return nil,"",errors.New("DK8_PENDING")}
 if objective!=Carrying&&objective!=Door{return nil,"",errors.New("DK8_OBJECTIVE_NOT_DECLARED")}
 key,e:=Context(before);if e!=nil{return nil,"",e}
 row:=p.Counts[key]
 best:=-1.0;choice:=-1
 for a,c:=range row{
  if c.N==0{continue}
  successes:=c.Carry
  if objective==Door{successes=c.Door}
  if successes==0{continue}
  benefit:=float64(successes)/float64(c.N)
  if benefit>best{best=benefit;choice=a}
 }
 if choice<0{return nil,"ABSTAIN_NO_LEARNED_POSITIVE_ACTION",nil}
 return &choice,"SELECT_LEARNED_PUBLIC_ACTION",nil
}
func (p *Policy) Commit(before Public,action int)error{
 if !p.healthy()||p.Pending!=nil||action<0||action>=7{return errors.New("DK8_INVALID_COMMIT")}
 raw,e:=canon(before);if e!=nil{return e}
 p.Pending=&Pending{BeforeSHA:digest(raw),Action:action}
 return nil
}
func (p *Policy) ObserveFit(before Public,action int,after Public,reward float64,terminal bool)error{
 if !p.healthy()||p.Pending!=nil{return errors.New("DK8_FIT_INVALID_PENDING")}
 return p.observe(before,action,after,reward,terminal,false)
}
func (p *Policy) Feedback(before Public,action int,after Public,reward float64,terminal bool)error{
 if !p.healthy()||p.Pending==nil{return errors.New("DK8_NO_SELECTED_PUBLIC_COMMIT")}
 raw,e:=canon(before)
 if e!=nil||action!=p.Pending.Action||digest(raw)!=p.Pending.BeforeSHA{
  return errors.New("DK8_PUBLIC_SELECTED_ACTION_PROVENANCE")}
 if e=p.observe(before,action,after,reward,terminal,true);e!=nil{return e}
 p.Pending=nil
 return nil
}
func (p *Policy) observe(before Public,action int,after Public,reward float64,terminal bool,live bool)error{
 if !p.healthy()||action<0||action>=7||math.IsNaN(reward)||math.IsInf(reward,0)||reward<0||reward>1||p.Samples>=MaxSamples{return errors.New("DK8_INVALID_REAL_SELECTED_REWARD")}
 key,e:=Context(before);if e!=nil{return e}
 carry,door,e:=observedEffects(before,after);if e!=nil{return e}
 row,seen:=p.Counts[key]
 if !seen&&len(p.Counts)>=MaxContexts{return errors.New("DK8_CONTEXT_BUDGET")}
 if row[action].N>=MaxSamples{return errors.New("DK8_ACTION_COUNT_OVERFLOW")}
 b:=row[action];b.N++;if carry{b.Carry++};if door{b.Door++};row[action]=b
 p.Counts[key]=row;p.Samples++;if live{p.Live++}
 if !p.healthy(){return errors.New("DK8_POSTUPDATE_CONSERVATION")}
 return nil
}
type State struct{
 Schema string `json:"schema"`
 Counts map[string]Row `json:"counts"`
 Samples uint64 `json:"samples"`
 Live uint64 `json:"live"`
}
type Envelope struct{State State `json:"state"`; SHA256 string `json:"sha256"`}
func (p *Policy) Snapshot()([]byte,error){
 if !p.healthy()||p.Pending!=nil{return nil,errors.New("DK8_PENDING_SNAPSHOT")}
 s:=State{Schema:Schema,Counts:p.Counts,Samples:p.Samples,Live:p.Live}
 raw,e:=json.Marshal(s);if e!=nil{return nil,e}
 blob,e:=json.Marshal(Envelope{State:s,SHA256:digest(raw)})
 if e!=nil||len(blob)>MaxCheckpoint{return nil,errors.New("DK8_SNAPSHOT_TOO_LARGE")}
 return blob,nil
}
func Restore(blob []byte)(*Policy,error){
 if len(blob)==0||len(blob)>MaxCheckpoint{return nil,errors.New("DK8_UNTRUSTED_LENGTH")}
 dec:=json.NewDecoder(bytes.NewReader(blob));dec.DisallowUnknownFields()
 var e Envelope
 if err:=dec.Decode(&e);err!=nil{return nil,err}
 if err:=dec.Decode(new(any));err!=io.EOF{return nil,errors.New("DK8_TRAILING_DATA")}
 raw,err:=json.Marshal(e.State)
 if err!=nil||e.State.Schema!=Schema||digest(raw)!=e.SHA256{return nil,errors.New("DK8_BAD_AUTH_OR_SCHEMA")}
 p:=&Policy{Counts:e.State.Counts,Samples:e.State.Samples,Live:e.State.Live}
 if !p.healthy(){return nil,errors.New("DK8_INCONSISTENT_REAL_EVENT_COUNTS")}
 copy,err:=p.Snapshot();if err!=nil||!bytes.Equal(copy,blob){return nil,errors.New("DK8_NONCANONICAL")}
 return p,nil
}
func (p *Policy) Summary()string{return fmt.Sprintf("%d/%d",p.Samples,p.Live)}
