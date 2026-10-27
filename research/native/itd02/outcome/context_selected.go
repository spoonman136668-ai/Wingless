package outcome

import (
 "bytes"
 "crypto/sha256"
 "encoding/json"
 "errors"
 "io"
 "math"

 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
 r8 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r8"
)

// ContextSelected is an isolated optional successor to global-reward selection.
// Context is the current source-public FrontBlocked bit, never a private map.
// Original R8 and the existing Selected checkpoint format remain unchanged.
type ContextSelected struct {
 Base *Selected
 Bins [2][3]Bin
}
func NewContextSelected()*ContextSelected{return &ContextSelected{Base:NewSelected()}}
func contextIndex(s r1.PublicGrid)int {if s.FrontBlocked{return 1};return 0}
func(c *ContextSelected) valid()bool{
 if c==nil||c.Base==nil||c.Base.Policy==nil||c.Base.Rewards==nil||
  !c.Base.Policy.Healthy()||!c.Base.Rewards.Valid()||
  c.Base.Rewards.NextID!=c.Base.Policy.LiveSeen+1{return false}
 var total uint64
 var byAction [3]Bin
 for _,row:=range c.Bins {
  for a,b:=range row{
   if math.IsNaN(b.Sum)||math.IsInf(b.Sum,0)||
    math.Abs(b.Sum)>float64(b.Count)*RewardLimit{return false}
   total+=uint64(b.Count)
   byAction[a].Count+=b.Count
   byAction[a].Sum+=b.Sum
  }
 }
 if total!=c.Base.Policy.LiveSeen{return false}
 for a:=range byAction{
  global:=c.Base.Rewards.Bins[a]
  if byAction[a].Count!=global.Count||
   math.Abs(byAction[a].Sum-global.Sum)>1e-7{return false}
 }
 return true
}

// Choose first balances source-public coverage separately in each context.
// After coverage it exploits the uniquely highest empirical reward for that
// context; ties fall back to the original native R8 action.
func(c *ContextSelected)Choose(before r1.PublicGrid,minCount uint32)(int,error){
 if !c.valid()||c.Base.Policy.Armed||c.Base.Rewards.Closed||minCount<2{
  return 0,errors.New("OUTCOME_CONTEXT_CHOICE_CUSTODY")
 }
 raw,err:=c.Base.Policy.Snapshot();if err!=nil{return 0,err}
 candidate,err:=r8.Restore(raw);if err!=nil{return 0,err}
 row:=c.Bins[contextIndex(before)]
 target:=-1
 for a,b:=range row {
  if b.Count<minCount&&(target<0||b.Count<row[target].Count){target=a}
 }
 if target<0 && c.Base.Rewards.NextID>1 && (c.Base.Rewards.NextID-1)%5==0 {
  // Keep sampling each public context after initial coverage. Without
  // ongoing selected outcomes a context's early preference cannot recover
  // from a changing reward relationship.
  offset:=int((c.Base.Rewards.NextID-1)/5%3)
  target=offset
  for j:=1;j<3;j++{
   a:=(offset+j)%3
   if row[a].Count<row[target].Count{target=a}
  }
 }
 if target<0{
  best:=-1
  var top float64
  ties:=0
  for a,b:=range row{
   mean:=b.Sum/float64(b.Count)
   if best<0||mean>top{best=a;top=mean;ties=1
   }else if mean==top{ties++}
  }
  if ties==1{target=best}
 }
 native,err:=candidate.Choose(before);if err!=nil{return 0,err}
 if target<0{target=native}
 if target!=native{
  if !candidate.Armed||!candidate.Core.Core.Route.Pending||
   candidate.Core.Core.Route.Before!=before||
   candidate.Core.Core.Route.Action!=native{
   return 0,errors.New("OUTCOME_CONTEXT_ROUTE")
  }
  candidate.Core.Core.Route.Action=target
 }
 if !candidate.Healthy(){return 0,errors.New("OUTCOME_CONTEXT_INVALID")}
 c.Base.Policy=candidate
 return target,nil
}

// Feedback applies genuine R8 selected public feedback before committing the
// contextual evidence. Rejected feedback cannot mutate either component.
func(c *ContextSelected)Feedback(before r1.PublicGrid,action int,after r1.PublicGrid,e Event)error{
 if !c.valid()||!c.Base.Policy.Armed{return errors.New("OUTCOME_CONTEXT_FEEDBACK_CUSTODY")}
 raw,err:=c.Base.SnapshotSelected();if err!=nil{return err}
 candidate,err:=RestoreSelected(raw);if err!=nil{return err}
 if err=candidate.Feedback(before,action,after,e);err!=nil{return err}
 bins:=c.Bins
 b:=&bins[contextIndex(before)][action]
 b.Count++
 b.Sum+=e.Reward
 next:=&ContextSelected{Base:candidate,Bins:bins}
 if !next.valid(){return errors.New("OUTCOME_CONTEXT_INVALID_FEEDBACK")}
 c.Base=candidate;c.Bins=bins
 return nil
}
func(c *ContextSelected)BeginEpisode()error{
 if !c.valid(){return errors.New("OUTCOME_CONTEXT_RESET_CUSTODY")}
 raw,err:=c.Base.SnapshotSelected();if err!=nil{return err}
 candidate,err:=RestoreSelected(raw);if err!=nil{return err}
 if err=candidate.BeginEpisode();err!=nil{return err}
 c.Base=candidate
 return nil
}

const ContextSchema="wingless.itd02.selected-context-outcome.v1"
type contextState struct{
 Base []byte `json:"base"`
 Bins [2][3]Bin `json:"context_bins"`
}
type contextEnvelope struct{
 Schema string `json:"schema"`
 State contextState `json:"state"`
 Digest [32]byte `json:"digest"`
}
func(c *ContextSelected)Snapshot()([]byte,error){
 if !c.valid(){return nil,errors.New("OUTCOME_CONTEXT_SNAPSHOT_CUSTODY")}
 raw,err:=c.Base.SnapshotSelected();if err!=nil{return nil,err}
 state:=contextState{Base:raw,Bins:c.Bins}
 body,err:=json.Marshal(state);if err!=nil{return nil,err}
 return json.Marshal(contextEnvelope{Schema:ContextSchema,State:state,Digest:sha256.Sum256(body)})
}
func RestoreContextSelected(raw []byte)(*ContextSelected,error){
 if len(raw)==0||len(raw)>5<<20{return nil,errors.New("OUTCOME_CONTEXT_CHECKPOINT_SIZE")}
 dec:=json.NewDecoder(bytes.NewReader(raw));dec.DisallowUnknownFields()
 var wire contextEnvelope
 if err:=dec.Decode(&wire);err!=nil{return nil,err}
 if err:=dec.Decode(new(struct{}));err!=io.EOF{return nil,errors.New("OUTCOME_CONTEXT_TRAILING")}
 canonical,err:=json.Marshal(wire)
 if err!=nil||!bytes.Equal(raw,canonical)||wire.Schema!=ContextSchema{
  return nil,errors.New("OUTCOME_CONTEXT_NONCANONICAL")
 }
 body,err:=json.Marshal(wire.State);if err!=nil{return nil,err}
 if sha256.Sum256(body)!=wire.Digest{return nil,errors.New("OUTCOME_CONTEXT_DIGEST")}
 base,err:=RestoreSelected(wire.State.Base);if err!=nil{return nil,err}
 c:=&ContextSelected{Base:base,Bins:wire.State.Bins}
 if !c.valid(){return nil,errors.New("OUTCOME_CONTEXT_CUSTODY_MISMATCH")}
 return c,nil
}
