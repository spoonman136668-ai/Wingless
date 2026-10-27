package outcome

import (
 "bytes"
 "testing"
)

// Prospective source-shaped fixtures, never an authentic world or replay of
// previously observed FourRooms test seeds.
func TestProgressV3BoundedLeftRightAcquisitionAcrossColdRestart(t *testing.T){
 p:=NewProgressSelected()
 before:=selectedFixture()
 before.FrontBlocked=false
 before.GoalX=before.X+5
 before.GoalY=before.Y+3 // public lateral goal; no prescribed turn answer
 var choices [3]int
 for i:=0;i<64;i++{
  action,err:=p.Choose(before,2)
  if err!=nil{t.Fatalf("choose %d: %v",i,err)}
  choices[action]++
  after:=publicConsequence(before,action)
  if err=p.Feedback(before,action,after,Event{ID:uint64(i+1),Action:action,Terminal:true});err!=nil{
   t.Fatalf("feedback %d: %v",i,err)
  }
  if err=p.BeginEpisode();err!=nil{t.Fatalf("terminal reset %d: %v",i,err)}
 }
 if choices[0]==0||choices[1]==0||choices[2]==0{
  t.Fatalf("missing sampled left/right/forward evidence: %+v",choices)
 }
 if choices[0]+choices[1]>=48{
  t.Fatalf("unbounded turning displaced native forward motion: %+v",choices)
 }
 if p.Base.Base.Base.Rewards.NextID!=65||p.Base.Base.Base.Policy.LiveSeen!=64 {
  t.Fatalf("evidence counter mismatch: next=%d live=%d",
   p.Base.Base.Base.Rewards.NextID,p.Base.Base.Base.Policy.LiveSeen)
 }
 raw,err:=p.Snapshot();if err!=nil{t.Fatal(err)}
 cold,err:=RestoreProgressSelected(raw)
 if err!=nil||!cold.valid(){t.Fatalf("native persistent progress not restored: %v",err)}
 identical,err:=cold.Snapshot()
 if err!=nil||!bytes.Equal(raw,identical){t.Fatal("cold persistence lost exact bytes")}
}
func TestProgressV3NoForcedTurnWithoutLateralGoalOrObstacle(t *testing.T){
 p:=NewProgressSelected()
 before:=selectedFixture()
 before.FrontBlocked=false
 before.GoalX=before.X+4
 before.GoalY=before.Y
 before.Direction=0 // goal straight forward
 for i:=0;i<8;i++{
  a,err:=p.Choose(before,2);if err!=nil{t.Fatal(err)}
  after:=publicConsequence(before,a)
  if err=p.Feedback(before,a,after,Event{ID:uint64(i+1),Action:a,Terminal:true});err!=nil{t.Fatal(err)}
  if err=p.BeginEpisode();err!=nil{t.Fatal(err)}
 }
 baselineBytes,err:=p.Base.Snapshot();if err!=nil{t.Fatal(err)}
 original,err:=RestoreMotionSelected(baselineBytes);if err!=nil{t.Fatal(err)}
 expected,err:=original.Choose(before,2);if err!=nil{t.Fatal(err)}
 chosen,err:=p.Choose(before,2)
 if err!=nil||chosen!=expected{
  t.Fatalf("public aligned goal must not force turning: %d motion %d %v",chosen,expected,err)
 }
}
func TestProgressV3OldCheckpointSchemaIsNotSilentUpgrade(t *testing.T){
 p:=NewProgressSelected()
 raw,err:=p.Snapshot();if err!=nil{t.Fatal(err)}
 cold,err:=RestoreProgressSelected(raw);if err!=nil||!cold.valid(){t.Fatal(err)}
 // Old artifact cannot silently acquire authority of the new learning rule.
 tampered:=bytes.Replace(raw,[]byte("selected-goal-progress.v3"),
                              []byte("selected-goal-progress.v2"),1)
 if bytes.Equal(tampered,raw){t.Fatal("schema substitution test missing source marker")}
 if _,err=RestoreProgressSelected(tampered);err==nil{t.Fatal("v2 checkpoint accepted as v3")}
}
