package r7
import("bytes";"os";"os/exec";"testing";r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1";r5 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r5")
func state(x,y,dir int)r1.PublicGrid{return r1.PublicGrid{X:x,Y:y,Direction:dir,GoalX:21,GoalY:20,FrontBlocked:false}}
func recorded(t *testing.T,p *Policy,b r1.PublicGrid,a int,n r1.PublicGrid){
 t.Helper();session:=&p.Core.Core.Core.Route
 session.Pending=true;session.Before=b;session.Action=a
 if e:=p.Feedback(b,a,n);e!=nil{t.Fatal(e)}
}
func rotations(t *testing.T,p *Policy){
 for d:=0;d<4;d++{recorded(t,p,state(20,20,d),1,state(20,20,(d+1)%4))}
}
func TestOnlyHeadingChangesNeverTriggerPhysicalFrontier(t *testing.T){
 p:=New();rotations(t,p)
 s:=state(20,20,0)
 if _,ok:=p.physical(s);ok{t.Fatal("new heading mistaken for new physical place")}
 original,e:=p.Core.Core.Snapshot();if e!=nil{t.Fatal(e)}
 baseline,e:=r5.Restore(original);if e!=nil{t.Fatal(e)}
 expected,e:=baseline.Choose(s);if e!=nil{t.Fatal(e)}
 actual,e:=p.Choose(s);if e!=nil||actual!=expected{t.Fatal("no physical evidence must delegate exact R5 decision")}
 if _,e:=p.Choose(s);e==nil{t.Fatal("duplicate choice without action consequence")}
 if e:=p.Feedback(s,(actual+1)%3,s);e==nil{t.Fatal("unselected action accepted")}
}
func TestActuallyObservedPhysicalGoalBeatsOrientationOnlyRoute(t *testing.T){
 p:=New();rotations(t,p)
 s:=state(20,20,0)
 recorded(t,p,s,2,state(21,20,0))
 a,found:=p.physical(s)
 if !found||a!=2{t.Fatalf("genuine physical progress lost %d %v",a,found)}
 chosen,e:=p.Choose(s);if e!=nil||chosen!=2{t.Fatal("verified public physical goal route not selected")}
 if e:=p.Feedback(s,2,state(21,20,0));e!=nil{t.Fatal(e)}
 if p.Core.Seen!=6{t.Fatal("selected real step did not conserve graph events")}
}
func TestContradictionCannotPromoteUnreliablePhysicalEdge(t *testing.T){
 p:=New();rotations(t,p)
 s:=state(20,20,0)
 recorded(t,p,s,2,state(21,20,0))
 recorded(t,p,s,2,state(20,21,0))
 edge:=p.Core.Edges[pose(s)*3+2]
 if edge.Contradictions!=1{t.Fatal("conflicting true public effect not retained")}
 if _,ok:=p.physical(s);ok{t.Fatal("contradicted physical edge promoted")}
}
func TestColdOSAndResetRetainsBaseWithoutEpisodeLeak(t *testing.T){
 if path:=os.Getenv("R7_COLD_CHECKPOINT");path!=""{
  b,e:=os.ReadFile(path);if e!=nil{t.Fatal(e)}
  p,e:=Restore(b);if e!=nil{t.Fatal(e)}
  if p.Core.Seen!=5||p.Core.Core.Core.Route.Steps!=5{t.Fatal("cold episode source lost")}
  if action,ok:=p.physical(state(20,20,0));!ok||action!=2{t.Fatal("cold physical route missing")}
  return
 }
 p:=New();rotations(t,p);recorded(t,p,state(20,20,0),2,state(21,20,0))
 original,e:=p.Snapshot();if e!=nil{t.Fatal(e)}
 restored,e:=Restore(original);if e!=nil{t.Fatal(e)}
 same,e:=restored.Snapshot();if e!=nil||!bytes.Equal(original,same){t.Fatal("source cold checkpoint drift")}
 if _,e:=Restore(append(append([]byte{},original...),byte('x')));e==nil{t.Fatal("tampered bytes accepted")}
 file:=t.TempDir()+"/state.json"
 if e:=os.WriteFile(file,original,0600);e!=nil{t.Fatal(e)}
 cmd:=exec.Command(os.Args[0],"-test.run=^TestColdOSAndResetRetainsBaseWithoutEpisodeLeak$")
 cmd.Env=append(os.Environ(),"R7_COLD_CHECKPOINT="+file)
 if output,e:=cmd.CombinedOutput();e!=nil{t.Fatalf("cold process failed %v %s",e,output)}
 learned,e:=p.Core.Core.Core.Core.Snapshot();if e!=nil{t.Fatal(e)}
 if e:=p.BeginEpisode();e!=nil{t.Fatal(e)}
 if p.Core.Seen!=0||len(p.Core.Edges)!=0{t.Fatal("previous maze graph leaked into next world")}
 later,e:=p.Core.Core.Core.Core.Snapshot()
 if e!=nil||!bytes.Equal(learned,later){t.Fatal("episode reset changed FIT model")}
}
