// Package dk15 is a separate, opt-in native source-selected actor continuation.
// Genuine source missions are NEVER replaced with a fake DoorKey mission.
// Original DoorKey delegates to unchanged DK13; other public missions use a
// developmental stateful public successor learner initialized from immutable
// DK13 FIT motion stats and then trained only on own selected returned effects.
// No real environment calls or original-source edits live in this package.
package dk15

import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "io"
 "math"
 old "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk13"
 motion "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk12"
 boundary "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk14"
)

const Schema="wingless.itd02.dk15.mission-preserved-native-public-actor.v1"
const MaxBytes=2800000
const MaxSteps=192
type Public=boundary.Public
type Cell struct{
 Count uint16 `json:"count"`
 Last string `json:"last"`
}
type Choices [7]Cell
type Actor struct{
 Original *old.Policy
 Boundary *boundary.Boundary
 Mode string
 Motion [7]motion.Stat
 Attempts map[string]Choices
 Visits map[string]uint16
 Counts [7]uint16
}
func New()*Actor{
 return &Actor{Original:old.New(),Boundary:boundary.New(),
 Attempts:map[string]Choices{},Visits:map[string]uint16{}}
}
func hash(p Public)(string,error){return boundary.Digest(p)}
func (a *Actor) Begin(p Public)error{
 if a.Mode!=""{return errors.New("DK15_EPISODE_ALREADY_STARTED")}
 if e:=a.Boundary.Begin(p);e!=nil{return e}
 if p.Mission==boundary.LegacyMission{
  if e:=a.Original.BeginEpisode();e!=nil{return e}
  a.Mode="ORIGINAL_DOORKEY_UNCHANGED"
 }else{
  a.Motion=a.Original.Motion // read-only FIT motion stats copied, not borrowed mutable state
  a.Mode="NATIVE_NEW_MISSION_PUBLIC_SELECTED"
 }
 return nil
}
func (a *Actor) Choose()(*int,string,error){
 if a.Mode==""||a.Boundary==nil||a.Boundary.Current==nil||
 a.Boundary.Pending!=nil||a.Boundary.Finished{
  return nil,"",errors.New("DK15_NO_SELECTED_SOURCE_ACTION_PHASE")
 }
 p:=*a.Boundary.Current
 if a.Mode=="ORIGINAL_DOORKEY_UNCHANGED"{return a.Original.Choose(
  old.Public{Image:p.Image,Direction:p.Direction,Mission:p.Mission})}
 if a.Mode!="NATIVE_NEW_MISSION_PUBLIC_SELECTED"{
  return nil,"",errors.New("DK15_UNAUTHORIZED_ACTOR_MODE")
 }
 key,e:=hash(p);if e!=nil{return nil,"",e}
 row:=a.Attempts[key]
 held:=p.Image[3][6][0]>=5&&p.Image[3][6][0]<=7
 best:=-1
 high:=int64(-1<<60)
 for action,stat:=range a.Motion{
  num:=int64(12)*int64(stat.SameHeadingChanged)-4*int64(stat.HeadingChanged)-
   8*int64(stat.HeldLost)
  if held{num+=4*int64(stat.HeldRetainedChanged)}
  score:=num*1024/(int64(stat.N)+64)-
     96*int64(a.Counts[action])-3072*int64(row[action].Count)
  if row[action].Last!=""{score-=1024*int64(a.Visits[row[action].Last])}
  if score>high{high=score;best=action}
 }
 if high>0{return &best,"SELECT_DK15_NATIVE_PUBLIC_SUCCESSOR",nil}
 minimum:=uint16(MaxSteps+1)
 for _,c:=range row{if c.Count<minimum{minimum=c.Count}}
 b,_:=json.Marshal(old.Public{Image:p.Image,Direction:p.Direction,Mission:p.Mission})
 start:=int(sha256.Sum256(b)[0])%7
 for j:=0;j<7;j++{
  action:=(start+j)%7
  if row[action].Count==minimum{return &action,"SELECT_DK15_NATIVE_UNTRIED_PUBLIC_SUCCESSOR",nil}
 }
 return nil,"",errors.New("DK15_NO_SELECTED_SOURCE_ACTION")
}
func (a *Actor) Commit(action int)error{
 if a.Mode==""{return errors.New("DK15_ACTION_BEFORE_PUBLIC_BEGIN")}
 if _,e:=a.Boundary.Commit(action);e!=nil{return e}
 if a.Mode=="ORIGINAL_DOORKEY_UNCHANGED"{
  p:=*a.Boundary.Current
  return a.Original.Commit(old.Public{Image:p.Image,Direction:p.Direction,Mission:p.Mission},action)
 }
 return nil
}
func (a *Actor) Feedback(action int,next Public,reward float64,term,trunc bool)error{
 if a.Mode==""||a.Boundary==nil||a.Boundary.Current==nil||
 a.Boundary.Pending==nil||*a.Boundary.Pending!=action||a.Boundary.Finished{
  return errors.New("DK15_UNAUTHENTICATED_SELECTED_FEEDBACK")
 }
 if e:=next.Valid();e!=nil{return e}
 before:=*a.Boundary.Current
 if next.Mission!=before.Mission{return errors.New("DK15_GENUINE_SOURCE_MISSION_DRIFT")}
 if math.IsNaN(reward)||math.IsInf(reward,0)||reward<0||reward>1{
  return errors.New("DK15_BAD_REAL_SOURCE_REWARD")
 }
 if a.Mode=="ORIGINAL_DOORKEY_UNCHANGED"{
  if trunc{return errors.New("DK15_LEGACY_TRUNCATION_UNSUPPORTED")}
  if e:=a.Original.Feedback(old.Public{Image:before.Image,Direction:before.Direction,Mission:before.Mission},action,
   old.Public{Image:next.Image,Direction:next.Direction,Mission:next.Mission},reward,term);e!=nil{return e}
  _,e:=a.Boundary.Feedback(action,next,reward,term,trunc)
  return e
 }
 if a.Mode!="NATIVE_NEW_MISSION_PUBLIC_SELECTED"{return errors.New("DK15_BAD_SOURCE_MODE")}
 from,e:=hash(before);if e!=nil{return e}
 to,e:=hash(next);if e!=nil{return e}
 if len(a.Boundary.Events)>=MaxSteps||(len(a.Attempts)>=MaxSteps&&a.Attempts[from]==(Choices{}))||
 (len(a.Visits)>=MaxSteps&&a.Visits[to]==0){
  return errors.New("DK15_NATIVE_SOURCE_BOUND")
 }
 if _,e=a.Boundary.Feedback(action,next,reward,term,trunc);e!=nil{return e}
 row:=a.Attempts[from];cell:=row[action];cell.Count++;cell.Last=to;row[action]=cell
 a.Attempts[from]=row;a.Visits[to]++;a.Counts[action]++
 stat:=a.Motion[action];stat.N++
 if before.Direction!=next.Direction{stat.HeadingChanged++}
 if before.Direction==next.Direction{
  oldb,_:=json.Marshal(before.Image);newb,_:=json.Marshal(next.Image)
  if !bytes.Equal(oldb,newb){
   stat.SameHeadingChanged++
   carried:=before.Image[3][6]
   if carried[0]>=5&&carried[0]<=7{
    newcarry,_:=json.Marshal(next.Image[3][6])
    oldcarry,_:=json.Marshal(carried)
    if bytes.Equal(newcarry,oldcarry){stat.HeldRetainedChanged++}
   }
  }
 }
 oldHeld:=before.Image[3][6][0]>=5&&before.Image[3][6][0]<=7
 newHeld:=next.Image[3][6][0]>=5&&next.Image[3][6][0]<=7
 if oldHeld&&!newHeld{stat.HeldLost++}
 a.Motion[action]=stat
 return nil
}
type state struct{
 Schema string `json:"schema"`
 Original json.RawMessage `json:"original"`
 Boundary json.RawMessage `json:"boundary"`
 Mode string `json:"mode"`
 Motion [7]motion.Stat `json:"motion"`
 Attempts map[string]Choices `json:"attempts"`
 Visits map[string]uint16 `json:"visits"`
 Counts [7]uint16 `json:"counts"`
}
type envelope struct{
 State state `json:"state"`
 SHA256 string `json:"sha256"`
}
func sha(b []byte)string{h:=sha256.Sum256(b);return hex.EncodeToString(h[:])}
func (a *Actor) Snapshot()([]byte,error){
 orig,e:=a.Original.Snapshot();if e!=nil{return nil,e}
 pub,e:=a.Boundary.Snapshot();if e!=nil{return nil,e}
 state:=state{Schema,orig,pub,a.Mode,a.Motion,a.Attempts,a.Visits,a.Counts}
 b,e:=json.Marshal(state);if e!=nil{return nil,e}
 wire,e:=json.Marshal(envelope{state,sha(b)})
 if e!=nil||len(wire)>MaxBytes{return nil,errors.New("DK15_SNAPSHOT_LIMIT")}
 return wire,nil
}
func Restore(raw []byte)(*Actor,error){
 if len(raw)==0||len(raw)>MaxBytes{return nil,errors.New("DK15_BAD_CHECKPOINT_SIZE")}
 var wire envelope
 dec:=json.NewDecoder(bytes.NewReader(raw));dec.DisallowUnknownFields()
 if e:=dec.Decode(&wire);e!=nil{return nil,errors.New("DK15_BAD_CHECKPOINT")}
 var extra any
 if e:=dec.Decode(&extra);e!=io.EOF{return nil,errors.New("DK15_TRAILING_CHECKPOINT")}
 b,e:=json.Marshal(wire.State)
 if e!=nil||wire.State.Schema!=Schema||sha(b)!=wire.SHA256{
  return nil,errors.New("DK15_SNAPSHOT_SHA_DRIFT")
 }
 oldactor,e:=old.Restore(wire.State.Original);if e!=nil{return nil,e}
 pub,e:=boundary.Restore(wire.State.Boundary);if e!=nil{return nil,e}
 s:=wire.State
 if len(s.Attempts)>MaxSteps||len(s.Visits)>MaxSteps||
 (s.Mode!="ORIGINAL_DOORKEY_UNCHANGED"&&s.Mode!="NATIVE_NEW_MISSION_PUBLIC_SELECTED"&&s.Mode!=""){
  return nil,errors.New("DK15_BAD_ACTOR_STATE")
 }
 actor:=&Actor{Original:oldactor,Boundary:pub,Mode:s.Mode,Motion:s.Motion,
 Attempts:s.Attempts,Visits:s.Visits,Counts:s.Counts}
 again,e:=actor.Snapshot()
 if e!=nil||!bytes.Equal(again,raw){return nil,errors.New("DK15_NONCANONICAL_CHECKPOINT")}
 return actor,nil
}
