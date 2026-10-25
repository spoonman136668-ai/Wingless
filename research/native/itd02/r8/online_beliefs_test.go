package r8

import(
 "bytes"
 "os"
 "os/exec"
 "testing"
 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
 r5 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r5"
)

func public(blocked bool)r1.PublicGrid{
 return r1.PublicGrid{X:20,Y:20,Direction:0,GoalX:21,GoalY:22,FrontBlocked:blocked}
}
func actualTurn(s r1.PublicGrid,arm int,returnedBlocked bool)r1.PublicGrid{
 n:=s
 if arm==0{n.Direction=(s.Direction+3)%4}
 if arm==1{n.Direction=(s.Direction+1)%4}
 n.FrontBlocked=returnedBlocked
 return n
}
func fitRealPublicAPI(t *testing.T,p *Policy,leftBlocked bool,n int){
 t.Helper()
 s:=public(true)
 for i:=0;i<n;i++{
  for a:=0;a<2;a++{
   blocked:=leftBlocked
   if a==1{blocked=!blocked}
   // These are authored source-PUBLIC-SHAPED FIT engineering fixtures.
   // They are NOT authenticated original MiniGrid episodes or real reward.
   if e:=p.ObserveFit(s,a,actualTurn(s,a,blocked),0,true);e!=nil{
    t.Fatalf("original R5 source FIT gate failed: %v",e)
   }
  }
 }
}
func TestNativeLearnerActuallyChangesSelectedAction(t *testing.T){
 p:=New()
 fitRealPublicAPI(t,p,true,12)
 if !p.Healthy(){t.Fatal("native FIT feedback not conserved")}
 s:=public(true)
 risks,eligible,e:=p.EligiblePrediction(s)
 if e!=nil||!eligible||risks[1]>=risks[0]{t.Fatalf("no fitted outcome differences %#v %v %v",risks,eligible,e)}
 raw,e:=p.Core.Snapshot();if e!=nil{t.Fatal(e)}
 baseline,e:=r5.Restore(raw);if e!=nil{t.Fatal(e)}
 before,e:=baseline.Choose(s);if e!=nil{t.Fatal(e)}
 after,e:=p.Choose(s);if e!=nil{t.Fatal(e)}
 if after!=1{t.Fatalf("real native outcome choice not driven by learned public consequences: %d",after)}
 if before==after{t.Fatalf("must demonstrate actual change relative to same original R5 FIT model, got %d/%d",before,after)}
 if _,e:=p.Choose(s);e==nil{t.Fatal("second action without public source consequence accepted")}
 if e:=p.Feedback(s,after,actualTurn(s,after,false));e!=nil{t.Fatal(e)}
 if p.Core.Core.Core.Events!=24{t.Fatal("live selected feedback overwrote original R3 source FIT model")}
 if p.LiveSeen!=1||p.EpisodeSteps!=1||!p.Healthy(){t.Fatal("real native selected action not charged")}
}
func TestLearningContradictoryReturnsRevisesNativeChoice(t *testing.T){
 s:=public(true)
 p:=New()
 fitRealPublicAPI(t,p,true,8)
 first,eligible,e:=p.EligiblePrediction(s)
 if e!=nil||!eligible||first[1]>=first[0]{t.Fatal("initial learned risk state missing")}
 // Later legally selected FIT reports contradict prior experience.
 fitRealPublicAPI(t,p,false,20)
 revised,eligible,e:=p.EligiblePrediction(s)
 if e!=nil||!eligible||revised[0]>=revised[1]{t.Fatal("contradiction was not incorporated")}
 chosen,e:=p.Choose(s)
 if e!=nil||chosen!=0{t.Fatalf("native contradiction did not revise action %d %v",chosen,e)}
}
func TestUntrainedContextDelegatesExactOriginalR5(t *testing.T){
 p:=New()
 s:=public(true)
 orig,e:=p.Core.Snapshot();if e!=nil{t.Fatal(e)}
 baseline,e:=r5.Restore(orig);if e!=nil{t.Fatal(e)}
 expected,e:=baseline.Choose(s);if e!=nil{t.Fatal(e)}
 got,e:=p.Choose(s);if e!=nil||got!=expected{t.Fatalf("untrained context changed original R5 %d/%d %v",got,expected,e)}
 // Force feedback rejection and demand zero bookkeeping changes.
 before,e:=p.Snapshot();if e!=nil{t.Fatal(e)}
 if e:=p.Feedback(s,(got+1)%3,s);e==nil{t.Fatal("mismatched selected source action was admitted")}
 after,e:=p.Snapshot();if e!=nil||!bytes.Equal(before,after){t.Fatal("invalid source event changed native state")}
 if e:=p.Feedback(s,got,actualTurn(s,got,false));e!=nil{t.Fatal(e)}
}
func TestRejectedFitCannotCorruptNativeR5OrBelief(t *testing.T){
 p:=New()
 s:=public(true)
 before,e:=p.Snapshot();if e!=nil{t.Fatal(e)}
 invalid:=actualTurn(s,0,false);invalid.X=99
 if e:=p.ObserveFit(s,0,invalid,0,true);e==nil{t.Fatal("illegally shaped FIT accepted")}
 after,e:=p.Snapshot();if e!=nil||!bytes.Equal(before,after){t.Fatal("rejected FIT mutated original native state")}
}
func TestConventionalEqualPublicPosteriorTies(t *testing.T){
 p:=New()
 fitRealPublicAPI(t,p,true,8)
 s:=public(true)
 risk,eligible,e:=p.EligiblePrediction(s)
 if e!=nil||!eligible{t.Fatal("expected two source FIT-conditioned arms")}
 // Entirely independent equally FIT-trained conventional Beta(1,1)
 // calculation; superiority against equal-info controller is NOT claimed.
 leftBlocked,leftClear:=8,0
 rightBlocked,rightClear:=0,8
 conventionalLeft:=float64(leftBlocked+1)/float64(leftBlocked+leftClear+2)
 conventionalRight:=float64(rightBlocked+1)/float64(rightBlocked+rightClear+2)
 if risk[0]!=conventionalLeft||risk[1]!=conventionalRight{
  t.Fatal("native predicted probabilities differ from strong fully fitted count control")
 }
}
func TestRealFreshProcessColdRestoreAndNativeSourceCompatibility(t *testing.T){
 if filename:=os.Getenv("R8_NATIVE_COLD_CHECKPOINT");filename!=""{
  raw,e:=os.ReadFile(filename);if e!=nil{t.Fatal(e)}
  p,e:=Restore(raw);if e!=nil{t.Fatal(e)}
  if !p.Healthy()||p.FitSeen!=24{t.Fatal("cold accepted FIT events missing")}
  a,e:=p.Choose(public(true))
  if e!=nil||a!=1{t.Fatal("cold process forgot learned decision")}
  return
 }
 p:=New();fitRealPublicAPI(t,p,true,12)
 original,e:=p.Core.Snapshot();if e!=nil{t.Fatal(e)}
 restoredR5,e:=r5.Restore(original);if e!=nil{t.Fatal(e)}
 again,e:=restoredR5.Snapshot()
 if e!=nil||!bytes.Equal(original,again){t.Fatal("original R5 source/checkpoint compatibility lost")}
 raw,e:=p.Snapshot();if e!=nil{t.Fatal(e)}
 q,e:=Restore(raw);if e!=nil{t.Fatal(e)}
 restored,e:=q.Snapshot()
 if e!=nil||!bytes.Equal(restored,raw){t.Fatal("native R8 exact roundtrip drift")}
 if _,e:=Restore(append(append([]byte{},raw...),byte('x')));e==nil{t.Fatal("tampered native checkpoint accepted")}
 file:=t.TempDir()+"/r8-native.json"
 if e:=os.WriteFile(file,raw,0600);e!=nil{t.Fatal(e)}
 cmd:=exec.Command(os.Args[0],"-test.run=^TestRealFreshProcessColdRestoreAndNativeSourceCompatibility$")
 cmd.Env=append(os.Environ(),"R8_NATIVE_COLD_CHECKPOINT="+file)
 if out,e:=cmd.CombinedOutput();e!=nil{t.Fatalf("cold native process loss: %v %s",e,out)}
}
func TestExplicitEpisodeResetRetainsCausalLearning(t *testing.T){
 p:=New();fitRealPublicAPI(t,p,true,8)
 s:=public(true)
 before,e:=p.Core.Core.Core.Snapshot();if e!=nil{t.Fatal(e)}
 got,e:=p.Choose(s);if e!=nil{t.Fatal(e)}
 if e:=p.Feedback(s,got,actualTurn(s,got,false));e!=nil{t.Fatal(e)}
 if e:=p.BeginEpisode();e!=nil{t.Fatal(e)}
 after,e:=p.Core.Core.Core.Snapshot();if e!=nil||!bytes.Equal(before,after){
  t.Fatal("begin episode mutated learned original R3 source FIT model")
 }
 if p.Previous!=0||p.EpisodeSteps!=0||p.LiveSeen!=1||!p.Healthy(){t.Fatal("ephemeral source episode leaked or all-time learning lost")}
 a,e:=p.Choose(s)
 if e!=nil||a!=1{t.Fatal("retained causal source learning no longer influences new episode")}
}
