package dk15
import (
 "bytes"
 "testing"
 old "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk13"
)
func frame(mission string,front int)Public{
 p:=Public{Image:make([][][]int,7),Direction:0,Mission:mission}
 for i:=0;i<7;i++{
  p.Image[i]=make([][]int,7)
  for j:=0;j<7;j++{p.Image[i][j]=[]int{1,0,0}}
 }
 p.Image[3][5]=[]int{front,0,0}
 return p
}
func TestExactOriginalDoorKeyNativeActionAndCheckpointParity(t *testing.T){
 p:=frame("use the key to open the door and then get to the goal",1)
 actor:=New();baseline:=old.New()
 if e:=actor.Begin(p);e!=nil{t.Fatal(e)}
 if e:=baseline.BeginEpisode();e!=nil{t.Fatal(e)}
 candidate,reason,e:=actor.Choose();if e!=nil{t.Fatal(e)}
 original,oldReason,e:=baseline.Choose(old.Public{Image:p.Image,Direction:p.Direction,Mission:p.Mission})
 if e!=nil||candidate==nil||*candidate!=*original||reason!=oldReason{
  t.Fatalf("old action changed: %v %s / %v %s / %v",candidate,reason,original,oldReason,e)
 }
 c,e:=actor.Original.Snapshot();if e!=nil{t.Fatal(e)}
 b,e:=baseline.Snapshot();if e!=nil||!bytes.Equal(b,c){t.Fatal("original fit state changed")}
 if e:=actor.Commit(*candidate);e!=nil{t.Fatal(e)}
 if e:=baseline.Commit(old.Public{Image:p.Image,Direction:p.Direction,Mission:p.Mission},*original);e!=nil{t.Fatal(e)}
 after:=frame(p.Mission,4)
 if e:=actor.Feedback(*candidate,after,0,false,false);e!=nil{t.Fatal(e)}
 if e:=baseline.Feedback(old.Public{Image:p.Image,Direction:p.Direction,Mission:p.Mission},*original,
  old.Public{Image:after.Image,Direction:after.Direction,Mission:after.Mission},0,false);e!=nil{t.Fatal(e)}
 c,e=actor.Original.Snapshot();if e!=nil{t.Fatal(e)}
 b,e=baseline.Snapshot();if e!=nil||!bytes.Equal(b,c){t.Fatal("exact old source feedback drift")}
 if actor.Boundary.Current.Mission!=p.Mission{t.Fatal("original mission changed")}
}
func TestRealUnlockMissionNativeActionAndSelfSelectedLearning(t *testing.T){
 a:=New();before:=frame("open the door",1)
 if e:=a.Begin(before);e!=nil{t.Fatal(e)}
 if a.Mode!="NATIVE_NEW_MISSION_PUBLIC_SELECTED"{t.Fatal(a.Mode)}
 oldSnapshot,e:=a.Original.Snapshot();if e!=nil{t.Fatal(e)}
 action,_,e:=a.Choose();if e!=nil||action==nil||*action<0||*action>6{t.Fatal(e)}
 if e:=a.Commit(*action);e!=nil{t.Fatal(e)}
 if e:=a.Feedback(*action+1,before,0,false,false);e==nil{t.Fatal("wrong source action accepted")}
 if e:=a.Feedback(*action,frame("get to the green goal square",4),0,false,false);e==nil{
  t.Fatal("source mission substituted")}
 after:=frame("open the door",4)
 if e:=a.Feedback(*action,after,0,false,false);e!=nil{t.Fatal(e)}
 b,e:=a.Original.Snapshot();if e!=nil||!bytes.Equal(oldSnapshot,b){
  t.Fatal("new mission mutated historical old native model")}
 if len(a.Boundary.Events)!=1||a.Boundary.Current.Mission!="open the door"||
 a.AttemptsMustTotal()!=1{t.Fatal("genuine source-selected learning not recorded")}
 saved,e:=a.Snapshot();if e!=nil{t.Fatal(e)}
 cold,e:=Restore(saved);if e!=nil{t.Fatal(e)}
 again,e:=cold.Snapshot();if e!=nil||!bytes.Equal(again,saved){t.Fatal("opt-in cold restore drift")}
 next1,reason1,e:=a.Choose();if e!=nil{t.Fatal(e)}
 next2,reason2,e:=cold.Choose();if e!=nil||*next1!=*next2||reason1!=reason2{
  t.Fatal("restored native action changed")
 }
}
func (a *Actor) AttemptsMustTotal()int{
 n:=0;for _,r:=range a.Attempts{for _,v:=range r{n+=int(v.Count)}}
 return n
}
func TestGoalOnlyOriginalMissionAndTerminalNoOracle(t *testing.T){
 a:=New();p:=frame("get to the green goal square",1)
 if e:=a.Begin(p);e!=nil{t.Fatal(e)}
 act,_,e:=a.Choose();if e!=nil{t.Fatal(e)}
 if e:=a.Commit(*act);e!=nil{t.Fatal(e)}
 if e:=a.Feedback(*act,p,1,true,false);e!=nil{t.Fatal(e)}
 if _,_,e:=a.Choose();e==nil{t.Fatal("action after source-terminal accepted")}
 blob,e:=a.Snapshot();if e!=nil{t.Fatal(e)}
 restored,e:=Restore(blob);if e!=nil||!restored.Boundary.Finished{t.Fatal(e)}
 if _,e:=Restore(append(append([]byte{},blob...),0));e==nil{t.Fatal("corruption accepted")}
}
func TestPrivatePublicPayloadAndWrongPhase(t *testing.T){
 a:=New()
 if _,_,e:=a.Choose();e==nil{t.Fatal("not begun")}
 if e:=a.Commit(1);e==nil{t.Fatal("not begun")}
 p:=frame("open the door",1);p.Image[0][0]=[]int{4,2,9}
 if e:=a.Begin(p);e==nil{t.Fatal("invalid public image accepted")}
 b:=New();if e:=b.Begin(frame("open the door",1));e!=nil{t.Fatal(e)}
 if e:=b.Commit(7);e==nil{t.Fatal("invalid source action")}
 if e:=b.Commit(5);e!=nil{t.Fatal(e)}
 if e:=b.Commit(5);e==nil{t.Fatal("duplicate commit")}
}
