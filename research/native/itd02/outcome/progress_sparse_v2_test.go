package outcome

import "testing"

// This regression is based on the V1 authentic failure, but uses ONLY source-
// shaped engineering fixtures. It does not replay observed held-out seeds.
func TestProgressV2SparsePublicContextFallsBackToNativeMotion(t *testing.T){
 p:=NewProgressSelected()
 before:=selectedFixture()
 before.FrontBlocked=false
 for i:=0;i<6;i++{
  action,err:=p.Choose(before,2);if err!=nil{t.Fatalf("choose %d: %v",i,err)}
  after:=publicConsequence(before,action)
  if err=p.Feedback(before,action,after,Event{ID:uint64(i+1),Action:action,Terminal:true});err!=nil{t.Fatal(err)}
  if err=p.BeginEpisode();err!=nil{t.Fatal(err)}
 }
 novel:=before
 novel.GoalY=before.Y-3 // different public goal-relative bearing; no prior evidence
 sparse:=progressContext(novel)
 for a:=0;a<3;a++{
  if p.Bins[sparse][a].Count!=0{t.Fatalf("sparse context %d already populated",sparse)}
 }
 baselineBytes,err:=p.Base.Snapshot();if err!=nil{t.Fatal(err)}
 baseline,err:=RestoreMotionSelected(baselineBytes);if err!=nil{t.Fatal(err)}
 originalAction,err:=baseline.Choose(novel,2);if err!=nil{t.Fatal(err)}
 candidateAction,err:=p.Choose(novel,2)
 if err!=nil||candidateAction!=originalAction {
  t.Fatalf("sparse context chose %d instead of existing motion %d: %v",
   candidateAction,originalAction,err)
 }
 if candidateAction!=2{
  t.Fatalf("source-learned forward motion lost in novel sparse context: %d",candidateAction)
 }
 native:=p.Base.Base.Base.Policy
 if !native.Armed||!native.Core.Core.Route.Pending||native.Core.Core.Route.Action!=candidateAction {
  t.Fatal("sparse fallback did not preserve genuine selected native custody")
 }
}
func TestProgressV2CheckpointSchemaRejectsPredecessor(t *testing.T){
 p:=NewProgressSelected()
 snap,err:=p.Snapshot();if err!=nil{t.Fatal(err)}
 cold,err:=RestoreProgressSelected(snap);if err!=nil||!cold.valid(){t.Fatal(err)}
 // V1 has a different checkpoint authority and is never silently accepted
 // as v2, even if source-public features have the same shape.
 historical:=[]byte(`{"schema":"wingless.itd02.selected-goal-progress.v1"}`)
 if _,err=RestoreProgressSelected(historical);err==nil{
  t.Fatal("v1 checkpoint accepted in successor v2")
 }
}
