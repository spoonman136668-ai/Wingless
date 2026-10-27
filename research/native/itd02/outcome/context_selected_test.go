package outcome

import (
 "bytes"
 "testing"
)

// These are public-shaped engineering observations, not authenticated
// source-world episodes or scientific transfer evidence.
func TestContextSelectedLearnsOppositePublicStatesAndRestores(t *testing.T){
 c:=NewContextSelected()
 before:=selectedFixture()
 var counts [2][3]int
 for i:=0;i<12;i++{
  before.FrontBlocked=i%2==1
  ctx:=contextIndex(before)
  a,err:=c.Choose(before,2);if err!=nil{t.Fatalf("choose %d: %v",i,err)}
  counts[ctx][a]++
  after:=before
  if a==0{after.Direction=(before.Direction+3)%4
  }else if a==1{after.Direction=(before.Direction+1)%4}
  reward:=float64(0)
  if (ctx==0&&a==0)||(ctx==1&&a==1){reward=10}
  if err:=c.Feedback(before,a,after,Event{ID:uint64(i+1),Action:a,Reward:reward,Terminal:true});err!=nil{
   t.Fatalf("selected context feedback %d: %v",i,err)
  }
  if err:=c.BeginEpisode();err!=nil{t.Fatalf("reset %d: %v",i,err)}
 }
 for ctx,row:=range counts{
  for action,n:=range row{
   if n!=2{t.Fatalf("context %d action %d has %d samples",ctx,action,n)}
  }
 }
 raw,err:=c.Snapshot();if err!=nil{t.Fatal(err)}
 cold,err:=RestoreContextSelected(raw);if err!=nil{t.Fatal(err)}
 for ctx:=0;ctx<2;ctx++{
  fixture:=selectedFixture()
  fixture.FrontBlocked=ctx==1
  // Isolate pending-action state for each counterfactual context.
  again,err:=RestoreContextSelected(raw);if err!=nil{t.Fatal(err)}
  chosen,err:=again.Choose(fixture,2);if err!=nil{t.Fatal(err)}
  if chosen!=ctx{t.Fatalf("context %d selected %d, expected %d",ctx,chosen,ctx)}
  if !again.Base.Policy.Armed||again.Base.Policy.Core.Core.Route.Action!=chosen{
   t.Fatal("learned context action failed to arm native custody")
  }
  // A distinct legal public pose and goal were never part of training.
  // The only shared feature is the current source-public blocked flag.
  fresh:=selectedFixture()
  fresh.X=7+ctx
  fresh.Y=13+ctx
  fresh.GoalX=39
  fresh.GoalY=41
  fresh.Direction=2
  fresh.FrontBlocked=ctx==1
  transferred,err:=RestoreContextSelected(raw);if err!=nil{t.Fatal(err)}
  selected,err:=transferred.Choose(fresh,2)
  if err!=nil||selected!=ctx{
   t.Fatalf("public-context generalization %d got %d: %v",ctx,selected,err)
  }
  if !transferred.Base.Policy.Armed||
   transferred.Base.Policy.Core.Core.Route.Before!=fresh||
   transferred.Base.Policy.Core.Core.Route.Action!=selected{
   t.Fatal("generalized action bypassed original selected-action custody")
  }
 }
 if !cold.valid(){t.Fatal("cold context learner invalid")}
}

func TestContextSelectedRejectedFeedbackAndCorruptionAreAtomic(t *testing.T){
 c:=NewContextSelected()
 before:=selectedFixture()
 action,err:=c.Choose(before,2);if err!=nil{t.Fatal(err)}
 original,err:=c.Snapshot();if err!=nil{t.Fatal(err)}
 after:=before
 if action==0{after.Direction=3}else if action==1{after.Direction=1}
 bad:=(action+1)%3
 if err:=c.Feedback(before,bad,after,Event{ID:1,Action:bad,Reward:5});err==nil{
  t.Fatal("unselected context reward accepted")
 }
 unchanged,err:=c.Snapshot();if err!=nil{t.Fatal(err)}
 if !bytes.Equal(original,unchanged){t.Fatal("rejection mutated contextual checkpoint")}
 tampered:=append([]byte(nil),original...)
 tampered[len(tampered)/2]^=1
 if _,err:=RestoreContextSelected(tampered);err==nil{
  t.Fatal("corrupted context checkpoint accepted")
 }
}

func TestContextSelectedRevisesBlockedFrontRewardAfterContradictions(t *testing.T){
 c:=NewContextSelected()
 before:=selectedFixture()
 for i:=0;i<81;i++{
  action,err:=c.Choose(before,2)
  if err!=nil{t.Fatalf("choose event %d: %v",i+1,err)}
  after:=before
  if action==0{after.Direction=(before.Direction+3)%4
  }else if action==1{after.Direction=(before.Direction+1)%4}
  reward:=float64(0)
  if i<6{
   if action==0{reward=9}
  }else{
   if action==0{reward= -9}
   if action==1{reward=20}
  }
  if err:=c.Feedback(before,action,after,Event{ID:uint64(i+1),Action:action,Reward:reward,Terminal:true});err!=nil{
   t.Fatalf("feedback event %d: %v",i+1,err)
  }
  if err:=c.BeginEpisode();err!=nil{t.Fatalf("reset %d: %v",i+1,err)}
  if i==5 && c.Bins[1][0].Sum/float64(c.Bins[1][0].Count)!=9 {
   t.Fatal("initial reward preference not learned")
  }
 }
 if c.Bins[1][1].Count<=2{t.Fatal("context exploration stopped after initial coverage")}
 if c.Bins[1][1].Sum/float64(c.Bins[1][1].Count)<=
  c.Bins[1][0].Sum/float64(c.Bins[1][0].Count){
  t.Fatal("contradictory selected feedback did not reverse reward preference")
 }
 raw,err:=c.Snapshot();if err!=nil{t.Fatal(err)}
 cold,err:=RestoreContextSelected(raw);if err!=nil{t.Fatal(err)}
 a,err:=cold.Choose(before,2)
 if err!=nil||a!=1{t.Fatalf("cold restart failed revised blocked-front action: %d %v",a,err)}
}

func TestContextSelectedMatchedPublicInformationOppositeRewards(t *testing.T){
 // Both learners encounter the same public observations and selected-action
 // coverage. Only the reward assignments differ; this is an engineering
 // discrimination control, not an authenticated source-world experiment.
 learners:=[2]*ContextSelected{NewContextSelected(),NewContextSelected()}
 base:=selectedFixture()
 for i:=0;i<12;i++ {
  observation:=base
  observation.FrontBlocked=i%2==1
  var observedAction int
  for k,c:=range learners {
   action,err:=c.Choose(observation,2)
   if err!=nil{t.Fatalf("learner %d choose %d: %v",k,i,err)}
   if k==0{observedAction=action}else if action!=observedAction{
    t.Fatalf("matched public observations selected different exploration actions at %d",i)
   }
   after:=observation
   if action==0{after.Direction=(observation.Direction+3)%4
   }else if action==1{after.Direction=(observation.Direction+1)%4}
   preferred:=contextIndex(observation)
   if k==1{preferred=1-preferred}
   reward:=float64(0)
   if action==preferred{reward=10}
   if err:=c.Feedback(observation,action,after,Event{ID:uint64(i+1),Action:action,Reward:reward,Terminal:true});err!=nil{
    t.Fatalf("learner %d feedback %d: %v",k,i,err)
   }
   if err:=c.BeginEpisode();err!=nil{t.Fatal(err)}
  }
 }
 for ctx:=0;ctx<2;ctx++{
  fresh:=base
  fresh.X=3+ctx
  fresh.Y=17+ctx
  fresh.GoalX=48
  fresh.GoalY=50
  fresh.Direction=2
  fresh.FrontBlocked=ctx==1
  for k,c:=range learners{
   raw,err:=c.Snapshot();if err!=nil{t.Fatal(err)}
   cold,err:=RestoreContextSelected(raw);if err!=nil{t.Fatal(err)}
   action,err:=cold.Choose(fresh,2);if err!=nil{t.Fatal(err)}
   expected:=ctx
   if k==1{expected=1-ctx}
   if action!=expected{
    t.Fatalf("control learner %d context %d chose %d, want %d",k,ctx,action,expected)
   }
   if !cold.Base.Policy.Armed||cold.Base.Policy.Core.Core.Route.Before!=fresh{
    t.Fatal("matched-control selected action bypassed native pending custody")
   }
  }
 }
}
