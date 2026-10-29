package dk12

import (
 "bytes"
 "testing"
)
func observed(dir int,held bool)Public{
 p:=Public{Image:make([][][]int,7),Direction:dir,Mission:"use the key to open the door and then get to the goal"}
 for x:=0;x<7;x++{
  p.Image[x]=make([][]int,7)
  for y:=0;y<7;y++{p.Image[x][y]=[]int{1,0,0}}
 }
 if held{p.Image[3][6]=[]int{5,4,0}}
 return p
}
func TestSourceLearnedSceneEffectBeatsCameraTurnNotActionOracle(t *testing.T){
 m:=New()
 before:=observed(0,false)
 turned:=observed(1,false)
 changed:=observed(0,false);changed.Image[2][3]=[]int{4,0,1}
 if e:=m.ObserveFit(before,0,turned,0,false);e!=nil{t.Fatal(e)}
 if e:=m.ObserveFit(before,3,changed,0,false);e!=nil{t.Fatal(e)}
 if m.Samples!=2||m.Agent.Samples!=2{t.Fatal("native parent selected sample parity")}
 a,status,e:=m.Choose(before)
 if e!=nil||a==nil||*a!=3||status!="SELECT_NATIVE_LEARNED_SAME_HEADING_PUBLIC_EFFECT"{
  t.Fatalf("empirical same-heading selected effect not valued: %v %s %v",a,status,e)
 }
 snap,e:=m.Snapshot();if e!=nil{t.Fatal(e)}
 cold,e:=Restore(snap);if e!=nil{t.Fatal(e)}
 again,e:=cold.Snapshot()
 if e!=nil||!bytes.Equal(again,snap){t.Fatal("cold source checkpoint mismatch")}
 b,mode,e:=cold.Choose(before)
 if e!=nil||b==nil||*b!=3||mode!=status{t.Fatal("cold native proposal drift")}
 if _,e:=Restore(append(append([]byte{},snap...),0));e==nil{t.Fatal("trailing bytes accepted")}
}
func TestNoObservedProgressBreaksSingleRepeatedAction(t *testing.T){
 m:=New()
 before:=observed(0,false)
 if e:=m.BeginEpisode();e!=nil{t.Fatal(e)}
 seen:=map[int]bool{}
 for i:=0;i<9;i++{
  choice,_,e:=m.Choose(before);if e!=nil||choice==nil{t.Fatal(e)}
  seen[*choice]=true
  if e:=m.Commit(before,*choice);e!=nil{t.Fatal(e)}
  if e:=m.Feedback(before,*choice,before,0,false);e!=nil{t.Fatal(e)}
 }
 if len(seen)<4{t.Fatalf("still trapped in a repeated single action: %v",seen)}
 if m.Samples!=9||m.Agent.Samples!=9{t.Fatal("selected parent conservation failure")}
}
func TestRejectSelectedFeedbackMismatchAndPrivatePublic(t *testing.T){
 m:=New()
 before:=observed(0,false)
 fake:=observed(0,false);fake.Mission="hidden global waypoint"
 if a,_,e:=m.Choose(fake);e==nil||a!=nil{t.Fatal("untrusted public mission accepted")}
 if e:=m.BeginEpisode();e!=nil{t.Fatal(e)}
 if e:=m.Commit(before,1);e!=nil{t.Fatal(e)}
 if e:=m.Feedback(before,2,before,0,false);e==nil{t.Fatal("wrong selected consequence accepted")}
 if m.Samples!=0||m.Agent.Samples!=0{t.Fatal("invalid selected mutated native")}
 if e:=m.Feedback(before,1,before,0,false);e!=nil{t.Fatal(e)}
 if e:=m.Feedback(before,1,before,0,false);e==nil{t.Fatal("duplicate consequence accepted")}
 if m.Samples!=1{t.Fatal("duplicate changed counts")}
}
