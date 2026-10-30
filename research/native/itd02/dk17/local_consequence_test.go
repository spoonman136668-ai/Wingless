package dk17
import("bytes";"testing")
func scene(m string,dir,far,near int)Public{
 p:=Public{Image:make([][][]int,7),Direction:dir,Mission:m}
 for i:=0;i<7;i++{
  p.Image[i]=make([][]int,7)
  for j:=0;j<7;j++{p.Image[i][j]=[]int{1,0,0}}
 }
 p.Image[0][0]=[]int{far,0,0}
 p.Image[3][5]=[]int{near,0,0}
 return p
}
func TestDistinctFullViewsCanShareLearnedSelectedConsequence(t *testing.T){
 m:=New();mission:="use the key to open the door and then get to the goal"
 a:=scene(mission,0,1,1)
 after:=scene(mission,0,2,1)
 distant:=scene(mission,3,7,1)
 for i:=0;i<2;i++{
  if e:=m.ObserveFit(a,5,after,0,false);e!=nil{t.Fatal(e)}
 }
 key,e:=context(a);if e!=nil{t.Fatal(e)}
 k,e:=context(distant);if e!=nil||key!=k{t.Fatal("distinct full source views did not share local context")}
 if m.Local[key][5].Changed!=2||m.Agent.Samples!=2{
  t.Fatal("no genuine selected source consequence aggregation")
 }
 _,_,e=m.Choose(distant);if e!=nil{t.Fatal(e)}
 snapshot,e:=m.Snapshot();if e!=nil{t.Fatal(e)}
 restored,e:=Restore(snapshot);if e!=nil{t.Fatal(e)}
 copy,e:=restored.Snapshot()
 if e!=nil||!bytes.Equal(snapshot,copy){t.Fatal("cold source custody not identical")}
 if _,e=Restore(append(append([]byte{},snapshot...),0));e==nil{t.Fatal("tampered trailing source accepted")}
}
func TestDifferentLocalSceneOrMissionDoesNotGeneralize(t *testing.T){
 m:=New();goal:="get to the green goal square"
 a:=scene(goal,0,1,1);after:=scene(goal,0,2,1)
 for i:=0;i<2;i++{
  if e:=m.ObserveFit(a,4,after,0,false);e!=nil{t.Fatal(e)}
 }
 for _,b:=range []Public{scene(goal,2,1,3),scene("use the key to open the door and then get to the goal",2,1,1)}{
  action,reason,e:=m.Choose(b)
  fallback,why,fe:=m.Policy.Choose(b)
  if e!=nil||fe!=nil||action==nil||fallback==nil||*action!=*fallback||reason!=why{
   t.Fatalf("learned effect leaked across local cell or mission: %v %v %v %v",action,fallback,reason,why)
  }
 }
}
func TestRotationAndNoEffectNeverProduceProgressCredit(t *testing.T){
 m:=New();goal:="get to the green goal square"
 a:=scene(goal,0,1,1);rotate:=scene(goal,1,1,1)
 for i:=0;i<2;i++{
  if e:=m.ObserveFit(a,1,rotate,0,false);e!=nil{t.Fatal(e)}
  if e:=m.ObserveFit(a,2,a,0,false);e!=nil{t.Fatal(e)}
 }
 key,_:=context(a)
 if m.Local[key][1].Changed!=0||m.Local[key][1].NoEffect!=0||
    m.Local[key][2].NoEffect!=2{t.Fatal("heading rotation or no effect gained movement credit")}
 selected,status,e:=m.Choose(a);base,reason,pe:=m.Policy.Choose(a)
 if e!=nil||pe!=nil||*selected!=*base||status!=reason{t.Fatal("no-positive leaked into choices")}
 if e=m.BeginEpisode();e!=nil{t.Fatal(e)}
 if e=m.Commit(a,6);e!=nil{t.Fatal(e)}
 changedMission:=scene("use the key to open the door and then get to the goal",0,1,1)
 if e=m.Feedback(a,6,changedMission,1,true);e==nil{t.Fatal("changed source mission accepted")}
 if e=m.Feedback(a,5,a,0,false);e==nil{t.Fatal("unselected action accepted")}
 if m.Agent.Samples!=4{t.Fatal("failed source event contaminated counter")}
 if e=m.Feedback(a,6,a,0,false);e!=nil{t.Fatal(e)}
 if e=m.Feedback(a,6,a,0,false);e==nil{t.Fatal("duplicate source accepted")}
 if m.Agent.Samples!=5{t.Fatal("selected source counted incorrectly")}
}


func TestLocalSourceTransferCanChangeChoiceOnUnseenFullView(t *testing.T){
 mission:="use the key to open the door and then get to the goal"
 // A different local configuration teaches the parent a global source-effect
 // tendency. Real selected local evidence must overcome that in an unseen full
 // observation, not simply count training samples.
 model:=New()
 other:=scene(mission,0,1,4);otherNext:=scene(mission,0,2,4)
 for i:=0;i<12;i++{
  if e:=model.ObserveFit(other,0,otherNext,0,false);e!=nil{t.Fatal(e)}
 }
 a:=scene(mission,0,1,1);after:=scene(mission,0,2,1)
 for i:=0;i<2;i++{
  if e:=model.ObserveFit(a,5,after,0,false);e!=nil{t.Fatal(e)}
 }
 unseen:=scene(mission,3,7,1)
 original,_,e:=model.Policy.Choose(unseen)
 if e!=nil||original==nil{t.Fatal(e)}
 selected,reason,e:=model.Choose(unseen)
 if e!=nil||selected==nil{t.Fatal(e)}
 if *original!=5 {
  if *selected!=5||reason!="SELECT_NATIVE_LOCAL_SELECTED_CAUSAL_EFFECT"{
   t.Fatalf("local evidence did not change policy on unseen full view: original %d new %d reason %s",*original,*selected,reason)
  }
 }else{
  t.Log("parent happens to select local supported action already; local evidence remains non-global")
 }
 if model.Agent.Samples!=14{t.Fatal("selected original source count drift")}
}
