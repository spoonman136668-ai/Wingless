package dk14
import ("bytes";"testing")
func visible(mission string)Public{
 p:=Public{Image:make([][][]int,7),Direction:0,Mission:mission}
 for x:=0;x<7;x++{p.Image[x]=make([][]int,7);for y:=0;y<7;y++{p.Image[x][y]=[]int{1,0,0}}}
 return p
}
func TestOriginalMissionInputsAreNotEnvironmentMasked(t *testing.T){
 model:=New()
 for _,mission:=range []string{DoorKey,GoalOnly}{
  p:=visible(mission)
  if _,_,e:=model.Choose(p);e!=nil{t.Fatalf("original mission rejected: %v",e)}
 }
 invalid:=visible("hidden global task route")
 if _,_,e:=model.Choose(invalid);e==nil{t.Fatal("untrusted mission accepted")}
}
func TestDifferentOriginalMissionSelectedRewardChangesNativeChoice(t *testing.T){
 model:=New()
 door:=visible(DoorKey);goal:=visible(GoalOnly)
 if e:=model.ObserveFit(door,1,door,1,true);e!=nil{t.Fatal(e)}
 if e:=model.ObserveFit(goal,4,goal,1,true);e!=nil{t.Fatal(e)}
 a,reason,e:=model.Choose(door)
 if e!=nil||a==nil||*a!=1||reason!="SELECT_NATIVE_GENUINE_MISSION_SELECTED_REWARD"{
  t.Fatalf("DoorKey source reward not learned: %v %s %v",a,reason,e)
 }
 b,reason,e:=model.Choose(goal)
 if e!=nil||b==nil||*b!=4||reason!="SELECT_NATIVE_GENUINE_MISSION_SELECTED_REWARD"{
  t.Fatalf("goal-only source reward not learned: %v %s %v",b,reason,e)
 }
 if model.Policy.Agent.Samples!=2{t.Fatal("original parent sample loss")}
 snap,e:=model.Snapshot();if e!=nil{t.Fatal(e)}
 restored,e:=Restore(snap);if e!=nil{t.Fatal(e)}
 rec,e:=restored.Snapshot()
 if e!=nil||!bytes.Equal(snap,rec){t.Fatal("native cold original mission drift")}
 if _,e=Restore(append(append([]byte{},snap...),0));e==nil{t.Fatal("trailing corruption accepted")}
}
func TestWrongMissionAndDuplicateSelectedFeedbackCannotMutate(t *testing.T){
 m:=New();door:=visible(DoorKey);goal:=visible(GoalOnly)
 if e:=m.BeginEpisode();e!=nil{t.Fatal(e)}
 if e:=m.Commit(door,2);e!=nil{t.Fatal(e)}
 if e:=m.Feedback(door,2,goal,1,true);e==nil{t.Fatal("changed mission accepted")}
 if e:=m.Feedback(goal,2,goal,1,true);e==nil{t.Fatal("mismatched before mission accepted")}
 if m.Policy.Agent.Samples!=0{t.Fatal("failed source mutated parent")}
 if e:=m.Feedback(door,2,door,0,false);e!=nil{t.Fatal(e)}
 if e:=m.Feedback(door,2,door,1,true);e==nil{t.Fatal("duplicate source accepted")}
 if m.Policy.Agent.Samples!=1||m.Missions[DoorKey][2].Success!=0{t.Fatal("selected source custody drift")}
}
