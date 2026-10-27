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
