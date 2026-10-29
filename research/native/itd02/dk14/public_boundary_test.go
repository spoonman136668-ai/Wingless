package dk14

import (
 "bytes"
 "encoding/json"
 "math"
 "testing"
 prior "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk13"
)

func fixture(mission string,front int)Public{
 p:=Public{Image:make([][][]int,7),Direction:0,Mission:mission}
 for i:=0;i<7;i++{
  p.Image[i]=make([][]int,7)
  for j:=0;j<7;j++{p.Image[i][j]=[]int{1,0,0}}
 }
 p.Image[3][5]=[]int{front,0,0}
 return p
}
func TestGenuineSourceMissionVerbatimAndLegacyDoorKeyParity(t *testing.T){
 old:=fixture(LegacyMission,1)
 if e:=old.Valid();e!=nil{t.Fatal(e)}
 legacy:=prior.Public{Image:old.Image,Direction:old.Direction,Mission:old.Mission}
 if e:=legacy.Valid();e!=nil{t.Fatal(e)}
 for _,mission:=range []string{"open the door","get to the green goal square","Bring the YELLOW Key Near the Door "}{
  p:=fixture(mission,4)
  raw,e:=json.Marshal(p);if e!=nil{t.Fatal(e)}
  q,e:=ReadPublic(raw);if e!=nil||q.Mission!=mission{t.Fatalf("source mission changed: %v %v",q.Mission,e)}
  if e:=q.Valid();e!=nil{t.Fatal(e)}
 }
 modified:=fixture("open the door",1)
 if e:= (prior.Public{Image:modified.Image,Direction:modified.Direction,Mission:modified.Mission}).Valid();e==nil{
  t.Fatal("original DK13 legacy mission guard unexpectedly changed")
 }
 one,e:=Digest(old);if e!=nil{t.Fatal(e)}
 two,e:=Digest(old);if e!=nil||one!=two{t.Fatal("original public SHA mismatch")}
}
func TestNoHiddenSourceStateOrInvalidOriginalFields(t *testing.T){
 p:=fixture("open the door",1)
 raw,e:=json.Marshal(p);if e!=nil{t.Fatal(e)}
 var object map[string]any
 if e=json.Unmarshal(raw,&object);e!=nil{t.Fatal(e)}
 for _,field:=range []string{"agent_pos","hidden_goal","oracle_door_stage","counterfactual"}{
  object[field]=42
  poisoned,_:=json.Marshal(object)
  if _,e:=ReadPublic(poisoned);e==nil{t.Fatal("private source accepted: "+field)}
  delete(object,field)
 }
 invalid:=[]Public{
  fixture("",1),fixture("bad\nmission",1),fixture(string([]byte{0xff}),1),
  fixture("open the door",11),
 }
 badDir:=fixture("open the door",1);badDir.Direction=4
 invalid=append(invalid,badDir)
 for _,x:=range invalid{if e=x.Valid();e==nil{t.Fatal("invalid input accepted")}}
}
func TestSelectedActionFeedbackColdCheckpointAndNoSourceOracle(t *testing.T){
 a:=fixture("open the door",1);b:=fixture("open the door",4)
 session:=New()
 if _,e:=session.Commit(5);e==nil{t.Fatal("before beginning accepted")}
 if e:=session.Begin(a);e!=nil{t.Fatal(e)}
 if _,e:=session.Commit(-1);e==nil{t.Fatal("negative source action accepted")}
 if _,e:=session.Commit(5);e!=nil{t.Fatal(e)}
 if _,e:=session.Feedback(1,b,0,false,false);e==nil{t.Fatal("wrong original action accepted")}
 bad:=fixture("get to the green goal square",4)
 if _,e:=session.Feedback(5,bad,0,false,false);e==nil{t.Fatal("in-episode mission forged")}
 if _,e:=session.Feedback(5,b,math.NaN(),false,false);e==nil{t.Fatal("NaN source reward accepted")}
 if session.Pending==nil||len(session.Events)!=0{t.Fatal("failed source action mutated")}
 rec,e:=session.Feedback(5,b,0,false,false);if e!=nil{t.Fatal(e)}
 if rec.Index!=1||rec.Action!=5{t.Fatal("selected source record mismatch")}
 start,_:=Digest(a);after,_:=Digest(b)
 if rec.Before!=start||rec.After!=after{t.Fatal("original public frame hash drift")}
 buf,e:=session.Snapshot();if e!=nil{t.Fatal(e)}
 cold,e:=Restore(buf);if e!=nil{t.Fatal(e)}
 again,e:=cold.Snapshot();if e!=nil||!bytes.Equal(buf,again){t.Fatal("cold checkpoint mismatch")}
 if _,e:=cold.Commit(2);e!=nil{t.Fatal(e)}
 pending,e:=cold.Snapshot();if e!=nil{t.Fatal(e)}
 next,e:=Restore(pending);if e!=nil{t.Fatal(e)}
 if _,e:=next.Feedback(2,b,1,true,false);e!=nil{t.Fatal(e)}
 if _,e:=next.Commit(2);e==nil{t.Fatal("source actions after terminal accepted")}
 if _,e:=Restore(append(append([]byte{},buf...),0));e==nil{t.Fatal("corrupted checkpoint accepted")}
}
