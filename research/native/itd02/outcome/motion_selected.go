package outcome

// Isolated source-public motion-evidence successor. No environment runner,
// privileged map, scripted action answer or scientific acceptance.
import (
 "bytes"
 "crypto/sha256"
 "encoding/json"
 "errors"
 "io"

 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
)

const MotionSchema="wingless.itd02.source-public-motion-evidence.v1"
type MotionBin struct {
 Count uint32 `json:"count"`
 Useful uint32 `json:"useful"`
}

// MotionSelected augments original R8 contextual reward custody. This branch
// never changes existing R8, ContextSelected, or their checkpoint formats.
// Useful means observed forward displacement in clear-front context, or a
// blocked-front turn followed by source-public clearance. Neither is a goal.
type MotionSelected struct {
 Base *ContextSelected
 Bins [2][3]MotionBin
}
func NewMotionSelected()*MotionSelected{return &MotionSelected{Base:NewContextSelected()}}
func(m *MotionSelected)valid()bool{
 if m==nil||m.Base==nil||!m.Base.valid(){return false}
 for c:=0;c<2;c++{
  for a,b:=range m.Bins[c] {
   if b.Useful>b.Count||b.Count!=m.Base.Bins[c][a].Count{return false}
  }
 }
 return true
}

// Motion preference is learned only from actually selected public feedback.
// Zero-motion evidence is not automatically positive. When genuinely observed
// external rewards are nonzero, defer to existing contextual reward selection.
func(m *MotionSelected)Choose(before r1.PublicGrid,minCount uint32)(int,error){
 if !m.valid()||m.Base.Base.Policy.Armed||m.Base.Base.Rewards.Closed||minCount<2 {
  return 0,errors.New("MOTION_CHOOSE_CUSTODY")
 }
 ctx:=contextIndex(before)
 row:=m.Bins[ctx]
 target:=-1
 for a,b:=range row {
  if b.Count<minCount&&(target<0||b.Count<row[target].Count){target=a}
 }
 if target<0 {
  hasReward:=false
  for _,b:=range m.Base.Bins[ctx]{if b.Sum!=0{hasReward=true}}
  if !hasReward {
   best:=-1
   unique:=true
   for a,b:=range row{
    if b.Useful==0{continue}
    if best<0||uint64(b.Useful)*uint64(row[best].Count)>uint64(row[best].Useful)*uint64(b.Count) {
     best=a;unique=true
    } else if uint64(b.Useful)*uint64(row[best].Count)==uint64(row[best].Useful)*uint64(b.Count){
     unique=false
    }
   }
   if best>=0&&unique{target=best}
  }
 }
 chosen,err:=m.Base.Choose(before,minCount)
 if err!=nil{return 0,err}
 if target<0||target==chosen{return chosen,nil}
 p:=m.Base.Base.Policy
 if !p.Armed||!p.Core.Core.Route.Pending||p.Core.Core.Route.Before!=before||
    p.Core.Core.Route.Action!=chosen{
  return 0,errors.New("MOTION_ROUTE_CUSTODY")
 }
 p.Core.Core.Route.Action=target
 if !m.valid(){return 0,errors.New("MOTION_INVALID_CHOICE")}
 return target,nil
}

func observedMotion(before r1.PublicGrid,action int,after r1.PublicGrid)(bool,error){
 if action<0||action>2||before.GoalX!=after.GoalX||before.GoalY!=after.GoalY {
  return false,errors.New("MOTION_INVALID_PUBLIC_CONSEQUENCE")
 }
 if action==0||action==1{
  delta:=3
  if action==1{delta=1}
  if before.X!=after.X||before.Y!=after.Y||after.Direction!=(before.Direction+delta)%4{
   return false,errors.New("MOTION_INVALID_PUBLIC_TURN")
  }
  return before.FrontBlocked&&!after.FrontBlocked,nil
 }
 if before.Direction!=after.Direction{
  return false,errors.New("MOTION_INVALID_PUBLIC_FORWARD_DIRECTION")
 }
 dx,dy:=after.X-before.X,after.Y-before.Y
 moved:=dx!=0||dy!=0
 if moved {
  if before.FrontBlocked{return false,errors.New("MOTION_BLOCKED_ADVANCE")}
  switch before.Direction {
  case 0: if dx!=1||dy!=0{return false,errors.New("MOTION_FORWARD_EAST")}
  case 1: if dx!=0||dy!=1{return false,errors.New("MOTION_FORWARD_SOUTH")}
  case 2: if dx!= -1||dy!=0{return false,errors.New("MOTION_FORWARD_WEST")}
  case 3: if dx!=0||dy!= -1{return false,errors.New("MOTION_FORWARD_NORTH")}
  default:return false,errors.New("MOTION_INVALID_DIRECTION")
  }
 }
 return moved,nil
}
func(m *MotionSelected)Feedback(before r1.PublicGrid,action int,after r1.PublicGrid,e Event)error{
 if !m.valid()||!m.Base.Base.Policy.Armed{
  return errors.New("MOTION_FEEDBACK_CUSTODY")
 }
 useful,err:=observedMotion(before,action,after)
 if err!=nil{return err}
 snap,err:=m.Base.Snapshot()
 if err!=nil{return err}
 clone,err:=RestoreContextSelected(snap)
 if err!=nil{return err}
 if err=clone.Feedback(before,action,after,e);err!=nil{return err}
 bins:=m.Bins
 b:=&bins[contextIndex(before)][action]
 if b.Count>=MaxEvents{return errors.New("MOTION_EVENT_CAP")}
 b.Count++
 if useful{b.Useful++}
 next:=&MotionSelected{Base:clone,Bins:bins}
 if !next.valid(){return errors.New("MOTION_FEEDBACK_INVALID")}
 m.Base=clone;m.Bins=bins
 return nil
}
func(m *MotionSelected)BeginEpisode()error{
 if !m.valid(){return errors.New("MOTION_RESET_CUSTODY")}
 snap,err:=m.Base.Snapshot();if err!=nil{return err}
 clone,err:=RestoreContextSelected(snap);if err!=nil{return err}
 if err=clone.BeginEpisode();err!=nil{return err}
 m.Base=clone
 return nil
}
type motionState struct{
 Base []byte `json:"base"`
 Bins [2][3]MotionBin `json:"motion_bins"`
}
type motionEnvelope struct{
 Schema string `json:"schema"`
 State motionState `json:"state"`
 Digest [32]byte `json:"digest"`
}
func(m *MotionSelected)Snapshot()([]byte,error){
 if !m.valid(){return nil,errors.New("MOTION_SNAPSHOT_CUSTODY")}
 raw,err:=m.Base.Snapshot();if err!=nil{return nil,err}
 state:=motionState{Base:raw,Bins:m.Bins}
 body,err:=json.Marshal(state);if err!=nil{return nil,err}
 return json.Marshal(motionEnvelope{Schema:MotionSchema,State:state,Digest:sha256.Sum256(body)})
}
func RestoreMotionSelected(raw []byte)(*MotionSelected,error){
 if len(raw)==0||len(raw)>5<<20{return nil,errors.New("MOTION_CHECKPOINT_SIZE")}
 d:=json.NewDecoder(bytes.NewReader(raw));d.DisallowUnknownFields()
 var env motionEnvelope
 if err:=d.Decode(&env);err!=nil{return nil,err}
 if err:=d.Decode(new(struct{}));err!=io.EOF{return nil,errors.New("MOTION_CHECKPOINT_TRAILING")}
 b,err:=json.Marshal(env)
 if err!=nil||!bytes.Equal(b,raw)||env.Schema!=MotionSchema{return nil,errors.New("MOTION_CHECKPOINT_CANONICAL")}
 body,err:=json.Marshal(env.State);if err!=nil{return nil,err}
 if sha256.Sum256(body)!=env.Digest{return nil,errors.New("MOTION_CHECKPOINT_DIGEST")}
 base,err:=RestoreContextSelected(env.State.Base);if err!=nil{return nil,err}
 v:=&MotionSelected{Base:base,Bins:env.State.Bins}
 if !v.valid(){return nil,errors.New("MOTION_CHECKPOINT_COUNTS")}
 return v,nil
}
