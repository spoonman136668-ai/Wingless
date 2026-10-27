// Package outcome provides an opt-in, source-public selected-action reward ledger.
// It does not change the frozen R8 policy or authorize environment execution.
package outcome

import (
 "crypto/sha256"
 "encoding/json"
 "errors"
 "math"
)

const Schema = "wingless.itd02.selected-outcome.v1"
const MaxEvents = 65535
const RewardLimit = 1000000.0

type Event struct {
 ID uint64 `json:"id"`
 Action int `json:"action"`
 Reward float64 `json:"reward"`
 Terminal bool `json:"terminal"`
 Truncated bool `json:"truncated"`
}
type Bin struct { Count uint32 `json:"count"`; Sum float64 `json:"sum"` }
type Ledger struct { Schema string `json:"schema"`; NextID uint64 `json:"next_id"`; Bins [3]Bin `json:"bins"`; Closed bool `json:"closed"` }
func New() *Ledger { return &Ledger{Schema:Schema, NextID:1} }
func (l *Ledger) Valid() bool {
 if l==nil || l.Schema!=Schema || l.NextID<1 || l.NextID>MaxEvents+1{return false}
 var n uint64
 for _,b:=range l.Bins { n+=uint64(b.Count); if math.IsNaN(b.Sum)||math.IsInf(b.Sum,0)||math.Abs(b.Sum)>float64(b.Count)*RewardLimit{return false} }
 return n==l.NextID-1
}
// Accept must be called only after the original selected-action custodian
// successfully accepts the corresponding public transition.
func(l *Ledger) Accept(e Event)error{
 if !l.Valid()||l.Closed||e.ID!=l.NextID||e.Action<0||e.Action>2||e.Terminal&&e.Truncated||math.IsNaN(e.Reward)||math.IsInf(e.Reward,0)||math.Abs(e.Reward)>RewardLimit||l.NextID>MaxEvents {
  return errors.New("OUTCOME_REJECTED")
 }
 b:=&l.Bins[e.Action]
 next:=b.Sum+e.Reward
 if math.IsNaN(next)||math.IsInf(next,0){return errors.New("OUTCOME_OVERFLOW")}
 b.Sum=next;b.Count++;l.NextID++
 l.Closed=e.Terminal||e.Truncated
 return nil
}
func(l *Ledger) Reset()error {if !l.Valid()||!l.Closed{return errors.New("OUTCOME_RESET_REJECTED")};l.Closed=false;return nil}
func(l *Ledger) Value(action int)(float64,bool){
 if !l.Valid()||action<0||action>2||l.Bins[action].Count==0{return 0,false}
 b:=l.Bins[action];return b.Sum/float64(b.Count),true
}
type checkpoint struct { Ledger Ledger `json:"ledger"`; Digest [32]byte `json:"digest"` }
func(l *Ledger) Snapshot()([]byte,error){
 if !l.Valid(){return nil,errors.New("OUTCOME_INVALID_STATE")}
 raw,_:=json.Marshal(l);return json.Marshal(checkpoint{Ledger:*l,Digest:sha256.Sum256(raw)})
}
func Restore(raw []byte)(*Ledger,error){
 var c checkpoint
 if len(raw)==0||len(raw)>4096{return nil,errors.New("OUTCOME_INVALID_CHECKPOINT")}
 if err:=json.Unmarshal(raw,&c);err!=nil{return nil,err}
 b,_:=json.Marshal(c.Ledger)
 if c.Digest!=sha256.Sum256(b)||!c.Ledger.Valid(){return nil,errors.New("OUTCOME_CHECKPOINT_REJECTED")}
 return &c.Ledger,nil
}
