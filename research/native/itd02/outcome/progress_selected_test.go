package outcome

import (
 "bytes"
 "testing"
)

func TestGoalProgressSourcePublicEvidenceLearnsAndColdRestores(t *testing.T){
 p:=NewProgressSelected()
 before:=selectedFixture()
 before.FrontBlocked=false
 context:=progressContext(before)
 var counts [3]int
 for i:=0;i<6;i++{
  a,err:=p.Choose(before,2);if err!=nil{t.Fatalf("choose %d: %v",i,err)}
  counts[a]++
  after:=publicConsequence(before,a)
  if err=p.Feedback(before,a,after,Event{ID:uint64(i+1),Action:a,Terminal:true});err!=nil{
   t.Fatalf("source-public consequence %d: %v",i,err)
  }
  if err=p.BeginEpisode();err!=nil{t.Fatalf("reset %d: %v",i,err)}
 }
 if counts!=[3]int{2,2,2}{t.Fatalf("missing equal selected-action evidence: %v",counts)}
 if p.Bins[context][2].Sum<=p.Bins[context][0].Sum ||
 p.Bins[context][2].Sum<=p.Bins[context][1].Sum{
  t.Fatalf("real forward goal progress was not retained: %+v",p.Bins[context])
 }
 snap,err:=p.Snapshot();if err!=nil{t.Fatal(err)}
 cold,err:=RestoreProgressSelected(snap);if err!=nil{t.Fatal(err)}
 fresh:=before
 fresh.X=13;fresh.Y=16;fresh.GoalX=14;fresh.GoalY=18
 a,err:=cold.Choose(fresh,2)
 if err!=nil||a!=2{t.Fatalf("goal-relative novel public pose action=%d error=%v",a,err)}
 native:=cold.Base.Base.Base.Policy
 if !native.Armed||!native.Core.Core.Route.Pending||
 native.Core.Core.Route.Before!=fresh||native.Core.Core.Route.Action!=a {
  t.Fatal("selected progress bypassed original native R8/R5 route custody")
 }
}
func TestGoalProgressRejectsForgedSelectedConsequencesAtomically(t *testing.T){
 p:=NewProgressSelected()
 before:=selectedFixture();before.FrontBlocked=false
 a,err:=p.Choose(before,2);if err!=nil{t.Fatal(err)}
 snap,err:=p.Snapshot();if err!=nil{t.Fatal(err)}
 after:=publicConsequence(before,a)
 tampered:=after;tampered.GoalX++
 if err=p.Feedback(before,a,tampered,Event{ID:1,Action:a});err==nil{
  t.Fatal("accepted changed public goal in selected transition")
 }
 unchanged,err:=p.Snapshot()
 if err!=nil||!bytes.Equal(unchanged,snap){t.Fatal("invalid feedback changed persistent state")}
 if err=p.Feedback(before,a,after,Event{ID:1,Action:(a+1)%3});err==nil{
  t.Fatal("accepted unselected event as progress")
 }
 unchanged,err=p.Snapshot()
 if err!=nil||!bytes.Equal(unchanged,snap){t.Fatal("unselected feedback consumed progress")}
 if err=p.Feedback(before,a,after,Event{ID:1,Action:a,Terminal:true});err!=nil{t.Fatal(err)}
 if err=p.BeginEpisode();err!=nil{t.Fatal(err)}
 valid,err:=p.Snapshot();if err!=nil{t.Fatal(err)}
 broken:=append([]byte(nil),valid...);broken[len(broken)/2]^=1
 if _,err=RestoreProgressSelected(broken);err==nil{t.Fatal("accepted corrupted cold checkpoint")}
}
func TestGoalProgressPublicFeatureNoHiddenWorldOrOracle(t *testing.T){
 before:=selectedFixture();before.FrontBlocked=false
 before.X=20;before.Y=20;before.GoalX=24;before.GoalY=22
 forward:=before;forward.X++
 turn:=before;turn.Direction=1
 away:=before;away.X--
 if !(publicProgress(before,forward)>publicProgress(before,turn) &&
       publicProgress(before,turn)>publicProgress(before,away)){
   t.Fatalf("public observed progress ordering: forward=%d turn=%d away=%d",
    publicProgress(before,forward),publicProgress(before,turn),publicProgress(before,away))
 }
 for _,f:=range []int{0,1}{
  before.FrontBlocked=f==1
  for d:=0;d<4;d++{
   before.Direction=d
   idx:=progressContext(before)
   if idx<0||idx>5{t.Fatalf("invalid bounded source context %d",idx)}
  }
 }
}
