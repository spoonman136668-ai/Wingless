package r3

import(
 "bytes"
 "encoding/base64"
 "encoding/json"
 "math"
 "os"
 "os/exec"
 "strconv"
 "testing"
 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
)

func state(x,y,dir,gx,gy int,blocked bool)r1.PublicGrid{
 return r1.PublicGrid{X:x,Y:y,GoalX:gx,GoalY:gy,Direction:dir,FrontBlocked:blocked}
}
// Entirely synthetic, source-only fixture for mechanism and integrity tests.
// Real navigation and conventional-control parity remain FORGE-owned science.
func afterTestFixture(s r1.PublicGrid,action int)r1.PublicGrid{
 a:=s
 switch action{
 case 0:a.Direction=(s.Direction+3)%4
 case 1:a.Direction=(s.Direction+1)%4
 case 2:
  if !s.FrontBlocked{
   switch s.Direction{case 0:a.X++;case 1:a.Y++;case 2:a.X--;case 3:a.Y--}
  }
 }
 return a
}
func teachAll(t *testing.T,p *Policy){
 t.Helper()
 for dir:=0;dir<4;dir++{
  for block:=0;block<2;block++{
   for arm:=0;arm<3;arm++{
    s:=state(12+dir*6+block,18+arm,dir,50,50,block==1)
    for i:=0;i<8;i++{
     a:=afterTestFixture(s,arm)
     if e:=p.Observe(s,arm,a,0,true);e!=nil{t.Fatalf("learn dir=%d blocked=%d arm=%d: %v",dir,block,arm,e)}
    }
   }
  }
 }
}
func TestZeroSourceActionsNoHardcodedWorld(t *testing.T){
 p:=New();s:=state(20,20,0,22,20,true)
 _,ok,e:=p.Plan(s)
 if e!=nil||ok{t.Fatalf("unseen physics leaked to planner: %v %v",ok,e)}
 a,e:=p.Choose(s);if e!=nil{t.Fatal(e)}
 base,e:=p.Core.Choose(s)
 if e!=nil||a!=base{t.Fatal("untrained R3 cannot replace original R2")}
 if p.Events!=0||len(p.Cells)!=0{t.Fatal("first-use violation")}
}
func TestLearnedActionsCanPlanWithoutHardcodedDynamics(t *testing.T){
 p:=New();teachAll(t,p)
 if p.Events!=192||p.Core.Core.Updates!=192{t.Fatal("selected consequence training counts diverged")}
 s:=state(40,40,0,42,40,false)
 predicted,ok,e:=p.Predict(s,2)
 if e!=nil||!ok||predicted.X!=41||predicted.Y!=40||predicted.Direction!=0{
  t.Fatalf("learned forward consequence invalid: %+v %v %v",predicted,ok,e)
 }
 a,found,e:=p.Plan(s)
 if e!=nil||!found||a!=2{t.Fatalf("learned route absent: arm=%d found=%v err=%v",a,found,e)}
 blocked:=state(40,40,0,42,40,true)
 unchanged,ok,e:=p.Predict(blocked,2)
 if e!=nil||!ok||unchanged.X!=40{t.Fatal("blocked action did not learn no motion")}
 a,found,e=p.Plan(blocked)
 if e!=nil||!found||a==2{t.Fatalf("blocked source should learn a turn first: %d %v %v",a,found,e)}
 // The learned map stores no goal label; an unfamiliar goal changes the route
 // through decision-time public goal observation, not a stored source answer.
 changedGoal:=state(40,40,0,38,40,false)
 turn,found,e:=p.Plan(changedGoal)
 if e!=nil||!found||turn==2{t.Fatalf("new goal did not change planned action: %d %v %v",turn,found,e)}
}
func TestConflictingNativeObservationsReviseLocalCell(t *testing.T){
 p:=New()
 a:=state(10,10,0,20,20,true)
 for i:=0;i<4;i++{if e:=p.Observe(a,2,a,0,true);e!=nil{t.Fatal(e)}}
 key:=poseKey(10,10,0)
 if !p.Cells[key].Blocked{t.Fatal("blocked state lost")}
 b:=a;b.FrontBlocked=false
 for i:=0;i<5;i++{if e:=p.Observe(b,2,afterTestFixture(b,2),0,true);e!=nil{t.Fatal(e)}}
 c:=p.Cells[key]
 if c.Blocked||c.Contradictions!=1||c.Seen!=9{t.Fatalf("failure to revise native local evidence: %+v",c)}
 blocked,ok,e:=p.Predict(a,2)
 if e!=nil||!ok||blocked.X!=a.X{t.Fatal("old causal outcome erased without evidence")}
 moved,ok,e:=p.Predict(b,2)
 if e!=nil||!ok||moved.X!=b.X+1{t.Fatal("new empirical consequence absent")}
}
func TestInvalidIsTransactionalAndBoundsAreReal(t *testing.T){
 p:=New();teachAll(t,p)
 original,e:=p.Snapshot();if e!=nil{t.Fatal(e)}
 s:=state(10,10,0,20,20,false)
 for _,bad:=range []struct{before r1.PublicGrid;arm int;after r1.PublicGrid;reward float64}{
  {s,3,s,0},
  {s,2,state(13,10,0,20,20,false),0},
  {s,2,state(65,10,0,20,20,false),0},
  {s,2,s,math.NaN()},
  {state(-1,10,0,20,20,false),2,s,0},
 }{
  if e:=p.Observe(bad.before,bad.arm,bad.after,bad.reward,true);e==nil{t.Fatal("invalid source transition accepted")}
  now,e:=p.Snapshot();if e!=nil||!bytes.Equal(original,now){t.Fatal("rejected source mutated learning or checkpoint")}
 }
 // Coupled-source custody regression: a model may not invent action physics,
 // world occupancy, or completed observations and then checkpoint them.
 {
s:=state(20,20,0,22,20,false)
 for _,alter:=range []struct{
  name string
  mutate func(*Policy)
 }{
  {"fabricated_action_effect",func(p *Policy){
   v:=&p.Learned[effectKey(s,2)]
   v.Samples=4
   v.Outcomes[0]=Effect{DX:1,DY:0,Heading:0,Count:4}
  }},
  {"fabricated_public_pose",func(p *Policy){
   p.Cells[poseKey(s.X,s.Y,s.Direction)]=CellEvidence{Seen:4,Blocked:false}
  }},
  {"fabricated_total_events",func(p *Policy){p.Events=4}},
  {"fabricated_original_core_events",func(p *Policy){p.Core.Core.Updates=4}},
  {"internally_balanced_but_without_original_core",func(p *Policy){
   v:=&p.Learned[effectKey(s,2)]
   v.Samples=4
   v.Outcomes[0]=Effect{DX:1,DY:0,Heading:0,Count:4}
   p.Cells[poseKey(s.X,s.Y,s.Direction)]=CellEvidence{Seen:4}
   p.Events=4
  }},
  {"original_core_events_do_not_replace_public_source",func(p *Policy){
   p.Core.Core.Updates=4
   p.Events=4
  }},
 }{
  t.Run(alter.name,func(t *testing.T){
   p:=New()
   alter.mutate(p)
   if _,_,e:=p.Predict(s,2);e==nil{t.Fatal("unobserved transition accepted")}
   if _,_,e:=p.Plan(s);e==nil{t.Fatal("unobserved plan accepted")}
   if _,e:=p.Snapshot();e==nil{t.Fatal("unobserved state serialized as signed model")}
  })
 }
 p:=New()
 teachAll(t,p)
 if p.Events!=p.Core.Core.Updates{t.Fatal("genuine source custody missing")}
 if _,e:=p.Snapshot();e!=nil{t.Fatalf("genuine observed model incorrectly rejected: %v",e)}

 }

}
func TestStrictModelCheckpointTamperAndColdOSRestart(t *testing.T){
 if input:=os.Getenv("WINGLESS_R3_CHILD_CHECKPOINT");input!=""{
  raw,e:=base64.StdEncoding.DecodeString(input);if e!=nil{t.Fatal(e)}
  p,e:=Restore(raw);if e!=nil{t.Fatal(e)}
  s:=state(40,40,0,42,40,true)
  a,e:=p.Choose(s);if e!=nil{t.Fatal(e)}
  want,e:=strconv.Atoi(os.Getenv("WINGLESS_R3_CHILD_ACTION"));if e!=nil{t.Fatal(e)}
  if a!=want||p.Events!=192{t.Fatal("cold OS restore lost independent learned state")}
  return
 }
 p:=New();teachAll(t,p)
 raw,e:=p.Snapshot();if e!=nil{t.Fatal(e)}
 restored,e:=Restore(raw);if e!=nil{t.Fatal(e)}
 again,e:=restored.Snapshot();if e!=nil||!bytes.Equal(raw,again){t.Fatal("snapshot not canonical or restorable")}
 probe:=state(40,40,0,42,40,true)
 want,e:=p.Choose(probe);if e!=nil{t.Fatal(e)}
 a,e:=restored.Choose(probe);if e!=nil||a!=want{t.Fatal("restored decision drift")}
 var wire map[string]interface{}
 if e=json.Unmarshal(raw,&wire);e!=nil{t.Fatal(e)}
 data:=wire["data"].(map[string]interface{})
 data["completed_events"]=float64(191)
 altered,_:=json.Marshal(wire)
 if _,e=Restore(altered);e==nil{t.Fatal("signed checkpoint accepted modified event identity")}
 if _,e=Restore(append(append([]byte{},raw...),byte('X')));e==nil{t.Fatal("trailing garbage accepted")}
 cmd:=exec.Command(os.Args[0],"-test.run=^TestStrictModelCheckpointTamperAndColdOSRestart$")
 cmd.Env=append(os.Environ(),"WINGLESS_R3_CHILD_CHECKPOINT="+base64.StdEncoding.EncodeToString(raw),
  "WINGLESS_R3_CHILD_ACTION="+strconv.Itoa(want))
 if output,e:=cmd.CombinedOutput();e!=nil{t.Fatalf("cold OS process model restore failed: %v %s",e,output)}
}
func TestIncompleteSourceNoPrematurePlanning(t *testing.T){
 p:=New()
 s:=state(20,20,0,22,20,false)
 for i:=0;i<3;i++{if e:=p.Observe(s,2,afterTestFixture(s,2),0,true);e!=nil{t.Fatal(e)}}
 _,ok,e:=p.Predict(s,2);if e!=nil||ok{t.Fatal("under-evidenced action given a fictitious prediction")}
 _,ok,e=p.Plan(s);if e!=nil||ok{t.Fatal("under-evidenced plan should fail closed")}
}
