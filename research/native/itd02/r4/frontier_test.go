package r4

import(
 "bytes"
 "encoding/json"
 "os"
 "os/exec"
 "strconv"
 "testing"
 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
)
func sample(x,y,dir,gx,gy int,blocked bool)r1.PublicGrid{
 return r1.PublicGrid{X:x,Y:y,GoalX:gx,GoalY:gy,Direction:dir,FrontBlocked:blocked}
}
func trainingEffect(s r1.PublicGrid,a int)r1.PublicGrid{
 out:=s
 switch a{
 case 0:out.Direction=(s.Direction+3)%4
 case 1:out.Direction=(s.Direction+1)%4
 case 2:if !s.FrontBlocked{
  switch s.Direction{case 0:out.X++;case 1:out.Y++;case 2:out.X--;case 3:out.Y--}
 }
 }
 return out
}
func teach(t *testing.T,p *Policy){
 t.Helper()
 for heading:=0;heading<4;heading++{
  for blocked:=0;blocked<2;blocked++{
   for action:=0;action<3;action++{
    s:=sample(12+heading*6+blocked,16+action,heading,40,40,blocked==1)
    for i:=0;i<8;i++{
     if e:=p.ObserveFit(s,action,trainingEffect(s,action),0,true);e!=nil{
      t.Fatalf("failed synthetic training of real-event-shaped transition: %v",e)
     }
    }
   }
  }
 }
}
func TestNoUnseenPhysicsOrEpisodeInterleaving(t *testing.T){
 p:=New()
 s:=sample(20,20,0,22,20,false)
 if _,ok,e:=p.Core.Predict(s,2);e!=nil||ok{t.Fatal("untrained physics available")}
 a,e:=p.Choose(s);if e!=nil{t.Fatal(e)}
 if _,e=p.Choose(s);e==nil{t.Fatal("second choice with no observed feedback")}
 if e=p.ObserveFit(s,a,s,0,true);e==nil{t.Fatal("evaluation observation trained model")}
 if e=p.Feedback(s,(a+1)%3,s);e==nil{t.Fatal("unselected action gained a label")}
 if p.Route.Steps!=0{t.Fatal("rejected actual action changed session")}
 if e=p.Feedback(s,a,s);e!=nil{t.Fatal(e)}
 if p.Core.Events!=0||p.Route.Steps!=1{t.Fatal("evaluation trained original model")}
 if e=p.BeginEpisode();e!=nil{t.Fatal(e)}
 if p.Route.Steps!=0||p.Core.Events!=0{t.Fatal("episode reset drift")}
}
func TestStalledPublicWorldCanLookForUnseenFrontier(t *testing.T){
 p:=New();teach(t,p)
 s:=sample(30,30,0,33,30,true)
 // Simulate a physically verified four-turn stall. Each chosen action is
 // followed by the true synthetic public action consequence.
 for i:=0;i<4;i++{
  a,e:=p.Choose(s);if e!=nil{t.Fatal(e)}
  after:=trainingEffect(s,a)
  if e=p.Feedback(s,a,after);e!=nil{t.Fatal(e)}
  s=after
  s.FrontBlocked=true
 }
 if p.Route.Steps!=4||p.Route.Stationary!=4{
  t.Fatalf("fixture did not produce a bounded turning loop: %+v",p.Route)
 }
 first,found,e:=p.frontier(s)
 if e!=nil||!found||first<0||first>2{
  t.Fatalf("no learned candidate frontier: first=%d found=%v err=%v",first,found,e)
 }
 a,e:=p.Choose(s)
 if e!=nil||a!=first{t.Fatalf("stalled episode failed to choose learned frontier: %d, %d, %v",a,first,e)}
 before:=p.Core.Events
 if e=p.Feedback(s,a,trainingEffect(s,a));e!=nil{t.Fatal(e)}
 if p.Core.Events!=before{t.Fatal("heldout action learned on future evidence")}
}
func TestSourceCorrectionAndResetNotOldTrainingWall(t *testing.T){
 p:=New();teach(t,p)
 s:=sample(32,32,0,35,35,true)
 for i:=0;i<3;i++{
  a,e:=p.Choose(s);if e!=nil{t.Fatal(e)}
  nxt:=trainingEffect(s,a)
  if e=p.Feedback(s,a,nxt);e!=nil{t.Fatal(e)}
  s=nxt;s.FrontBlocked=true
 }
 known,ok:=seenBlocked(p.Route,s)
 if ok&&known!=true{t.Fatal("incorrect episode-local obstacle")}
 raw,e:=p.Core.Snapshot();if e!=nil{t.Fatal(e)}
 if e=p.BeginEpisode();e!=nil{t.Fatal(e)}
 if len(p.Route.Poses)!=0||len(p.Route.Visits)!=0||p.Route.Steps!=0{
  t.Fatal("old-world occupancy contaminated reset")
 }
 again,e:=p.Core.Snapshot();if e!=nil||!bytes.Equal(raw,again){
  t.Fatal("new episode rewrote native learned effects")
 }
}
func TestRejectedSessionEventsAtomic(t *testing.T){
 p:=New();teach(t,p)
 s:=sample(20,20,0,25,25,false)
 a,e:=p.Choose(s);if e!=nil{t.Fatal(e)}
 stable,e:=p.Snapshot();if e!=nil{t.Fatal(e)}
 rejects:=[]r1.PublicGrid{
  sample(25,25,0,25,25,false),
  sample(66,20,0,25,25,false),
  sample(-1,20,0,25,25,false),
 }
 for _,after:=range rejects{
  if e=p.Feedback(s,a,after);e==nil{t.Fatal("invalid public episode accepted")}
  raw,e:=p.Snapshot();if e!=nil||!bytes.Equal(raw,stable){
   t.Fatal("invalid episode mutated native checkpoint")
  }
 }
 if e=p.BeginEpisode();e==nil{t.Fatal("reset while action pending accepted")}
 if e=p.Feedback(s,a,trainingEffect(s,a));e!=nil{t.Fatal(e)}
 if e=p.BeginEpisode();e!=nil{t.Fatal(e)}
}
func TestSourceLocalColdProcessAndCheckpointIntegrity(t *testing.T){
 if file:=os.Getenv("R4_CHECKPOINT_CHILD_FILE");file!=""{
  raw,e:=os.ReadFile(file);if e!=nil{t.Fatal(e)}
  p,e:=Restore(raw);if e!=nil{t.Fatal(e)}
  if p.Core.Events!=192||p.Route.Steps!=4||p.Route.Stationary!=4{
   t.Fatal("cold process lost native learned effects or session")
  }
  if _,e:=p.Choose(sample(30,30,0,33,30,true));e!=nil{t.Fatal(e)}
  return
 }
 p:=New();teach(t,p)
 s:=sample(30,30,0,33,30,true)
 for i:=0;i<4;i++{
  a,e:=p.Choose(s);if e!=nil{t.Fatal(e)}
  nxt:=trainingEffect(s,a)
  if e=p.Feedback(s,a,nxt);e!=nil{t.Fatal(e)}
  s=nxt;s.FrontBlocked=true
 }
 raw,e:=p.Snapshot();if e!=nil{t.Fatal(e)}
 restored,e:=Restore(raw);if e!=nil{t.Fatal(e)}
 same,e:=restored.Snapshot();if e!=nil||!bytes.Equal(raw,same){t.Fatal("roundtrip not deterministic")}
 var wire map[string]any
 if e=json.Unmarshal(raw,&wire);e!=nil{t.Fatal(e)}
 wire["sha256"]="0000000000000000000000000000000000000000000000000000000000000000"
 tampered,_:=json.Marshal(wire)
 if _,e:=Restore(tampered);e==nil{t.Fatal("modified checkpoint accepted")}
 if _,e:=Restore(append(append([]byte{},raw...),byte(88)));e==nil{t.Fatal("trailing data accepted")}
 path:=t.TempDir()+"/state.json"
 if e=os.WriteFile(path,raw,0600);e!=nil{t.Fatal(e)}
 cmd:=exec.Command(os.Args[0],"-test.run=^TestSourceLocalColdProcessAndCheckpointIntegrity$")
 cmd.Env=append(os.Environ(),"R4_CHECKPOINT_CHILD_FILE="+path,
                   "R4_CHECKPOINT_CHILD_EXPECT="+strconv.Itoa(len(raw)))
 if out,e:=cmd.CombinedOutput();e!=nil{t.Fatalf("cold real OS-process failed: %v: %s",e,out)}
}
