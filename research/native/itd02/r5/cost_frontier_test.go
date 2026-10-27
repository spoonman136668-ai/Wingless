package r5

import(
 "bytes"
 "os"
 "os/exec"
 "testing"
 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
)
func observation(x,y,dir,gx,gy int,blocked bool)r1.PublicGrid{
 return r1.PublicGrid{X:x,Y:y,GoalX:gx,GoalY:gy,Direction:dir,FrontBlocked:blocked}
}
func observedConsequence(s r1.PublicGrid,arm int)r1.PublicGrid{
 a:=s
 switch arm{
 case 0:a.Direction=(s.Direction+3)%4
 case 1:a.Direction=(s.Direction+1)%4
 case 2:
  if !s.FrontBlocked{
   switch s.Direction{case 0:a.X++;case 1:a.Y++;case 2:a.X--;case 3:a.Y--}
  }
 }
 return a
}
func train(t *testing.T,p *Policy){
 t.Helper()
 for direction:=0;direction<4;direction++{
  for blocked:=0;blocked<2;blocked++{
   for arm:=0;arm<3;arm++{
    s:=observation(12+direction*6+blocked,18+arm,direction,45,45,blocked==1)
    for i:=0;i<8;i++{
     if e:=p.ObserveFit(s,arm,observedConsequence(s,arm),0,true);e!=nil{t.Fatal(e)}
    }
   }
  }
 }
}
func TestUntrainedWorldCannotInventActionPhysics(t *testing.T){
 p:=New()
 s:=observation(20,20,0,21,20,false)
 _,ok,e:=p.chooseFrontier(s)
 if e!=nil||ok{t.Fatal("unlearned native environment dynamics hallucinated")}
 a,e:=p.Choose(s);if e!=nil{t.Fatal(e)}
 if _,e=p.Choose(s);e==nil{t.Fatal("second action accepted without public source consequence")}
 if e=p.Feedback(s,(a+1)%3,s);e==nil{t.Fatal("unselected action treated as true source evidence")}
 if p.Core.Route.Steps!=0{t.Fatal("rejected evidence mutated learned session")}
 if e=p.Feedback(s,a,observedConsequence(s,a));e!=nil{t.Fatal(e)}
 if p.Core.Core.Events!=0{t.Fatal("evaluation trained original causal model")}
}
func TestSourceOnlyCostFrontierIsBoundedAndDistinct(t *testing.T){
 p:=New();train(t,p)
 s:=observation(30,30,0,34,30,true)
 for i:=0;i<4;i++{
  a,e:=p.Choose(s);if e!=nil{t.Fatal(e)}
  n:=observedConsequence(s,a)
  if e=p.Feedback(s,a,n);e!=nil{t.Fatal(e)}
  s=n;s.FrontBlocked=true
 }
 if p.Core.Route.Stationary<4{t.Fatal("synthetic fixture did not induce stall")}
 a,ok,e:=p.chooseFrontier(s)
 if e!=nil||!ok||a<0||a>2{t.Fatalf("bounded learned frontier absent: %v %v %v",a,ok,e)}
 chosen,e:=p.Choose(s)
 if e!=nil||chosen!=a{t.Fatalf("chosen action differs from frozen cost objective %d %d %v",chosen,a,e)}
 if p.Core.Core.Events!=192{t.Fatal("untrained heldout model mutation")}
 if e=p.Feedback(s,chosen,observedConsequence(s,chosen));e!=nil{t.Fatal(e)}
 before,e:=p.Core.Core.Snapshot();if e!=nil{t.Fatal(e)}
 if e=p.BeginEpisode();e!=nil{t.Fatal(e)}
 after,e:=p.Core.Core.Snapshot();if e!=nil||!bytes.Equal(before,after){t.Fatal("episode reset rewrote original R3 FIT model")}
 if p.Core.Route.Steps!=0||len(p.Core.Route.Visits)!=0{t.Fatal("new world contaminated by old exploration")}
}
func TestColdProcessPreservesTrainingAndActualEpisodicCustody(t *testing.T){
 if state:=os.Getenv("R5_NATIVE_COLD_STATE");state!=""{
  raw,e:=os.ReadFile(state);if e!=nil{t.Fatal(e)}
  p,e:=Restore(raw);if e!=nil{t.Fatal(e)}
  if p.Core.Core.Events!=192||p.Core.Route.Steps!=4||p.Core.Route.Stationary!=4{
   t.Fatal("true OS restore lost fitted model or episode state")
  }
  s:=observation(31,30,0,34,30,true)
  if _,e:=p.Choose(s);e!=nil{t.Fatal(e)}
  return
 }
 p:=New();train(t,p)
 s:=observation(30,30,0,34,30,true)
 for i:=0;i<4;i++{
  a,e:=p.Choose(s);if e!=nil{t.Fatal(e)}
  n:=observedConsequence(s,a)
  if e=p.Feedback(s,a,n);e!=nil{t.Fatal(e)}
  s=n;s.FrontBlocked=true
 }
 raw,e:=p.Snapshot();if e!=nil{t.Fatal(e)}
 q,e:=Restore(raw);if e!=nil{t.Fatal(e)}
 identical,e:=q.Snapshot();if e!=nil||!bytes.Equal(raw,identical){t.Fatal("roundtrip drift")}
 if _,e=Restore(append(append([]byte{},raw...),byte('X')));e==nil{t.Fatal("trailing checksum forgery accepted")}
 f:=t.TempDir()+"/original.json"
 if e=os.WriteFile(f,raw,0600);e!=nil{t.Fatal(e)}
 c:=exec.Command(os.Args[0],"-test.run=^TestColdProcessPreservesTrainingAndActualEpisodicCustody$")
 c.Env=append(os.Environ(),"R5_NATIVE_COLD_STATE="+f)
 if out,e:=c.CombinedOutput();e!=nil{t.Fatalf("cold process lost original source custody: %v %s",e,out)}
}
func TestInvalidPublicObservationAndOriginalModelNoGo(t *testing.T){
 p:=New();train(t,p)
 s:=observation(25,25,0,40,40,false)
 for _,x:=range []r1.PublicGrid{
  observation(-1,25,0,40,40,false),
  observation(65,25,0,40,40,false),
  observation(25,25,4,40,40,false),
 }{
  if _,e:=p.Choose(x);e==nil{t.Fatal("invalid public observation accepted")}
 }
 original,e:=p.Core.Core.Snapshot();if e!=nil{t.Fatal(e)}
 a,e:=p.Choose(s);if e!=nil{t.Fatal(e)}
 if e=p.ObserveFit(s,a,observedConsequence(s,a),0,true);e==nil{
  t.Fatal("heldout episode masqueraded as training")
 }
 if e=p.Feedback(s,a,observedConsequence(s,a));e!=nil{t.Fatal(e)}
 current,e:=p.Core.Core.Snapshot();if e!=nil||!bytes.Equal(original,current){
  t.Fatal("true original native model changed outside FIT")
 }
}
