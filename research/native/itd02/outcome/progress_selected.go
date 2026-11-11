package outcome

// ProgressSelected is an opt-in native, source-public successor to MotionSelected.
// Evidence is accumulated exclusively from authenticated selected-action
// before/after public coordinates and orientations, never from source hidden state.
// The goal is PUBLIC via FullObsWrapper. This is an engineered learning feature,
// not an oracle, shortest path, or scientific proof of transfer.
import (
 "bytes"
 "crypto/sha256"
 "encoding/json"
 "errors"
 "io"

 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
)

const ProgressSchema="wingless.itd02.selected-goal-progress.v1"
type ProgressBin struct{
 Count uint32 `json:"count"`
 Sum int64 `json:"sum"`
}
type ProgressSelected struct{
 Base *MotionSelected
 Bins [6][3]ProgressBin
}
func NewProgressSelected()*ProgressSelected{return &ProgressSelected{Base:NewMotionSelected()}}
func pAbs(x int)int{if x<0{return -x};return x}
func goalDot(g r1.PublicGrid)int{
 dx:=g.GoalX-g.X;dy:=g.GoalY-g.Y
 switch g.Direction{case 0:return dx;case 1:return dy;case 2:return -dx;case 3:return -dy}
 return 0
}
func goalCross(g r1.PublicGrid)int{
 dx:=g.GoalX-g.X;dy:=g.GoalY-g.Y
 switch g.Direction{case 0:return dy;case 1:return -dx;case 2:return -dy;case 3:return dx}
 return 0
}
func pSign(x int)int{if x<0{return -1};if x>0{return 1};return 0}
func progressContext(g r1.PublicGrid)int{
 v:=pSign(goalCross(g))+1
 if g.FrontBlocked{v+=3}
 return v
}
// Progress quantifies ONLY the selected authentic public effect. Translation
// toward the publicly observed goal dominates direction-change shaping.
func publicProgress(before,after r1.PublicGrid)int64{
 d0:=pAbs(before.GoalX-before.X)+pAbs(before.GoalY-before.Y)
 d1:=pAbs(after.GoalX-after.X)+pAbs(after.GoalY-after.Y)
 return int64(8*(d0-d1)+pSign(goalDot(after))-pSign(goalDot(before)))
}
func(p *ProgressSelected)valid()bool{
 if p==nil||p.Base==nil||!p.Base.valid(){return false}
 for blocked:=0;blocked<2;blocked++ {
  for a:=0;a<3;a++{
   total:=uint64(0)
   for bearing:=0;bearing<3;bearing++{
    bin:=p.Bins[blocked*3+bearing][a]
    if bin.Sum>int64(bin.Count)*10||bin.Sum < -int64(bin.Count)*10{return false}
    total+=uint64(bin.Count)
   }
   if total!=uint64(p.Base.Bins[blocked][a].Count){return false}
  }
 }
 return true
}
// Native route remains the sole pending-action custodian. Every preference
// candidate must arm original R8/R5 choose first; no synthetic observations.
func(p *ProgressSelected)Choose(before r1.PublicGrid,minCount uint32)(int,error){
 if !p.valid()||p.Base.Base.Base.Policy.Armed||p.Base.Base.Base.Rewards.Closed||minCount<2{
  return 0,errors.New("PROGRESS_CHOOSE_CUSTODY")
 }
 context:=progressContext(before)
 row:=p.Bins[context]
 target:=-1
 for a,b:=range row {
  if b.Count<minCount&&(target<0||b.Count<row[target].Count){target=a}
 }
 if target<0{
  best:=-1
  unique:=true
  for a,b:=range row {
   if best<0||b.Sum*int64(row[best].Count)>row[best].Sum*int64(b.Count){
    best=a;unique=true
   }else if b.Sum*int64(row[best].Count)==row[best].Sum*int64(b.Count){
    unique=false
   }
  }
  if unique&&best>=0 {target=best}
 }
 // Preserve an unambiguous experienced external reward preference. Do not
 // replace learned reward evidence with engineered movement heuristics.
 prior:=p.Base.Base.Bins[contextIndex(before)]
 if prior[0].Count>=minCount&&prior[1].Count>=minCount&&prior[2].Count>=minCount {
  top:=-1;unique:=true
  for a,b:=range prior {
   if top<0||b.Sum*float64(prior[top].Count)>prior[top].Sum*float64(b.Count){
    top=a;unique=true
   }else if b.Sum*float64(prior[top].Count)==prior[top].Sum*float64(b.Count){
    unique=false
   }
  }
  if unique&&top>=0&&prior[top].Sum>0{target=top}
 }
 chosen,err:=p.Base.Choose(before,minCount)
 if err!=nil{return 0,err}
 if target<0||target==chosen{return chosen,nil}
 native:=p.Base.Base.Base.Policy
 if !native.Armed||!native.Core.Core.Route.Pending||
 native.Core.Core.Route.Before!=before||native.Core.Core.Route.Action!=chosen{
  return 0,errors.New("PROGRESS_ROUTE_CUSTODY")
 }
 native.Core.Core.Route.Action=target
 if !p.valid(){return 0,errors.New("PROGRESS_INVALID_CHOICE")}
 return target,nil
}
func(p *ProgressSelected)Feedback(before r1.PublicGrid,action int,after r1.PublicGrid,e Event)error{
 if !p.valid()||!p.Base.Base.Base.Policy.Armed{return errors.New("PROGRESS_FEEDBACK_CUSTODY")}
 // MotionSelected independently validates physical public consequences.
 snap,err:=p.Base.Snapshot();if err!=nil{return err}
 clone,err:=RestoreMotionSelected(snap);if err!=nil{return err}
 if err=clone.Feedback(before,action,after,e);err!=nil{return err}
 bins:=p.Bins
 bin:=&bins[progressContext(before)][action]
 if bin.Count>=MaxEvents{return errors.New("PROGRESS_EVENT_CAP")}
 bin.Count++
 bin.Sum+=publicProgress(before,after)
 candidate:=&ProgressSelected{Base:clone,Bins:bins}
 if !candidate.valid(){return errors.New("PROGRESS_INVALID_FEEDBACK")}
 p.Base=clone;p.Bins=bins
 return nil
}
func(p *ProgressSelected)BeginEpisode()error{
 if !p.valid(){return errors.New("PROGRESS_RESET_CUSTODY")}
 raw,err:=p.Base.Snapshot();if err!=nil{return err}
 candidate,err:=RestoreMotionSelected(raw);if err!=nil{return err}
 if err=candidate.BeginEpisode();err!=nil{return err}
 p.Base=candidate;return nil
}
type progressState struct{
 Base []byte `json:"base"`
 Bins [6][3]ProgressBin `json:"progress_bins"`
}
type progressEnvelope struct{
 Schema string `json:"schema"`
 State progressState `json:"state"`
 Digest [32]byte `json:"digest"`
}
func(p *ProgressSelected)Snapshot()([]byte,error){
 if !p.valid(){return nil,errors.New("PROGRESS_SNAPSHOT_CUSTODY")}
 raw,err:=p.Base.Snapshot();if err!=nil{return nil,err}
 state:=progressState{Base:raw,Bins:p.Bins}
 body,err:=json.Marshal(state);if err!=nil{return nil,err}
 return json.Marshal(progressEnvelope{Schema:ProgressSchema,State:state,Digest:sha256.Sum256(body)})
}
func RestoreProgressSelected(raw []byte)(*ProgressSelected,error){
 if len(raw)==0||len(raw)>6<<20{return nil,errors.New("PROGRESS_CHECKPOINT_SIZE")}
 d:=json.NewDecoder(bytes.NewReader(raw));d.DisallowUnknownFields()
 var env progressEnvelope
 if err:=d.Decode(&env);err!=nil{return nil,err}
 if err:=d.Decode(new(struct{}));err!=io.EOF{return nil,errors.New("PROGRESS_CHECKPOINT_TRAILING")}
 b,err:=json.Marshal(env)
 if err!=nil||!bytes.Equal(b,raw)||env.Schema!=ProgressSchema{return nil,errors.New("PROGRESS_CHECKPOINT_CANONICAL")}
 body,err:=json.Marshal(env.State);if err!=nil{return nil,err}
 if sha256.Sum256(body)!=env.Digest{return nil,errors.New("PROGRESS_CHECKPOINT_DIGEST")}
 base,err:=RestoreMotionSelected(env.State.Base);if err!=nil{return nil,err}
 p:=&ProgressSelected{Base:base,Bins:env.State.Bins}
 if !p.valid(){return nil,errors.New("PROGRESS_CHECKPOINT_COUNTS")}
 return p,nil
}
