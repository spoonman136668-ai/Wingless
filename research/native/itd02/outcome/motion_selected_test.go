package outcome

import (
 "bytes"
 "testing"

 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
)
func publicConsequence(before r1.PublicGrid,action int)r1.PublicGrid{
 after:=before
 switch action {
 case 0: after.Direction=(before.Direction+3)%4
 case 1: after.Direction=(before.Direction+1)%4
 case 2:
  if !before.FrontBlocked {
   switch before.Direction {
   case 0:after.X++
   case 1:after.Y++
   case 2:after.X--
   case 3:after.Y--
   }
  }
 }
 if before.FrontBlocked && action==0{after.FrontBlocked=false}
 return after
}

func TestMotionLearnsForwardFromSelectedDisplacementAndRestarts(t *testing.T){
 learner:=NewMotionSelected()
 before:=selectedFixture()
 before.FrontBlocked=false
 var seen [3]int
 for i:=0;i<6;i++ {
  action,err:=learner.Choose(before,2);if err!=nil{t.Fatalf("choose %d: %v",i,err)}
  seen[action]++
  after:=publicConsequence(before,action)
  if err=learner.Feedback(before,action,after,Event{ID:uint64(i+1),Action:action,Terminal:true});err!=nil{
   t.Fatalf("feedback %d: %v",i,err)
  }
  if err=learner.BeginEpisode();err!=nil{t.Fatalf("reset %d: %v",i,err)}
 }
 if seen!=[3]int{2,2,2}{t.Fatalf("unbalanced selected coverage: %+v",seen)}
 if learner.Bins[0][2].Useful!=2||learner.Bins[0][0].Useful!=0||learner.Bins[0][1].Useful!=0 {
  t.Fatalf("claimed unobserved public movement: %+v",learner.Bins[0])
 }
 raw,err:=learner.Snapshot();if err!=nil{t.Fatal(err)}
 cold,err:=RestoreMotionSelected(raw);if err!=nil{t.Fatal(err)}
 novel:=before
 novel.X=7;novel.Y=12;novel.GoalX=39;novel.GoalY=41
 action,err:=cold.Choose(novel,2)
 if err!=nil||action!=2{t.Fatalf("cold public motion preference: %d %v",action,err)}
 if !cold.Base.Base.Policy.Armed||!cold.Base.Base.Policy.Core.Core.Route.Pending||
 cold.Base.Base.Policy.Core.Core.Route.Before!=novel||
 cold.Base.Base.Policy.Core.Core.Route.Action!=action{
  t.Fatal("motion action escaped native selected-action pending route")
 }
}

func TestMotionBlockedTurnClearsFrontUsingActualFeedback(t *testing.T){
 learner:=NewMotionSelected()
 before:=selectedFixture()
 for i:=0;i<6;i++{
  a,err:=learner.Choose(before,2);if err!=nil{t.Fatal(err)}
  after:=publicConsequence(before,a)
  if err=learner.Feedback(before,a,after,Event{ID:uint64(i+1),Action:a,Terminal:true});err!=nil{t.Fatal(err)}
  if err=learner.BeginEpisode();err!=nil{t.Fatal(err)}
 }
 raw,err:=learner.Snapshot();if err!=nil{t.Fatal(err)}
 cold,err:=RestoreMotionSelected(raw);if err!=nil{t.Fatal(err)}
 action,err:=cold.Choose(before,2)
 if err!=nil||action!=0{t.Fatalf("expected learned public unblock turn, got %d %v",action,err)}
 if cold.Bins[1][0].Useful!=2||cold.Bins[1][2].Useful!=0 {
  t.Fatal("unobserved blocked progress credited")
 }
}
func TestMotionRejectsForgedConsequenceAndUnselectedFeedbackAtomically(t *testing.T){
 m:=NewMotionSelected()
 before:=selectedFixture();before.FrontBlocked=false
 action,err:=m.Choose(before,2);if err!=nil{t.Fatal(err)}
 raw,err:=m.Snapshot();if err!=nil{t.Fatal(err)}
 forged:=publicConsequence(before,action)
 forged.X+=3
 if err=m.Feedback(before,action,forged,Event{ID:1,Action:action});err==nil{
  t.Fatal("forged source-public displacement was accepted")
 }
 stale,err:=m.Snapshot();if err!=nil||!bytes.Equal(raw,stale){
  t.Fatal("invalid public consequence consumed native state")
 }
 after:=publicConsequence(before,action)
 wrong:=(action+1)%3
 if err=m.Feedback(before,action,after,Event{ID:1,Action:wrong});err==nil{
  t.Fatal("unselected reward accepted")
 }
 stale,err=m.Snapshot();if err!=nil||!bytes.Equal(raw,stale){
  t.Fatal("rejected selected feedback changed checkpoint")
 }
 if err=m.Feedback(before,action,after,Event{ID:1,Action:action,Terminal:true});err!=nil{t.Fatal(err)}
 if err=m.BeginEpisode();err!=nil{t.Fatal(err)}
 stable,err:=m.Snapshot();if err!=nil{t.Fatal(err)}
 tampered:=append([]byte(nil),stable...)
 tampered[len(tampered)/2]^=1
 if _,err=RestoreMotionSelected(tampered);err==nil{t.Fatal("corrupt motion checkpoint accepted")}
}
