package dk11

import (
 "bytes"
 "testing"
)

func observation(dir int,held bool) Public {
 p:=Public{Image:make([][][]int,7),Direction:dir,Mission:"use the key to open the door and then get to the goal"}
 for x:=0;x<7;x++{
  p.Image[x]=make([][]int,7)
  for y:=0;y<7;y++{p.Image[x][y]=[]int{1,0,0}}
 }
 p.Image[3][5]=[]int{1,0,0}
 if held{p.Image[3][6]=[]int{5,4,0}}
 return p
}
func TestOriginalSelectedObjectContinuityPreferenceAndColdRestore(t *testing.T){
 p:=New();held:=observation(0,true);dropped:=observation(0,false);moved:=observation(1,true)
 if e:=p.BeginEpisode();e!=nil{t.Fatal(e)}
 first,status,e:=p.Choose(held)
 if e!=nil||first==nil||*first<0||*first>6||status!="SELECT_NATIVE_UNKNOWN_PUBLIC_ACTION"{
  t.Fatalf("missing native no-oracle action: %v %s %v",first,status,e)
 }
 if e:=p.Commit(held,0);e!=nil{t.Fatal(e)}
 if e:=p.Feedback(held,1,held,0,false);e==nil{t.Fatal("wrong selected action accepted")}
 if p.Samples!=0||p.Policy.Samples!=0{t.Fatal("wrong feedback polluted original learner")}
 if e:=p.Feedback(held,0,dropped,0,false);e!=nil{t.Fatal(e)}
 if e:=p.Feedback(held,0,dropped,0,false);e==nil{t.Fatal("duplicate source receipt accepted")}
 if e:=p.Commit(held,1);e!=nil{t.Fatal(e)}
 if e:=p.Feedback(held,1,moved,0,false);e!=nil{t.Fatal(e)}
 if p.Samples!=2||p.Policy.Samples!=2{t.Fatal("selected sample conservation lost")}
 chosen,reason,e:=p.Choose(held)
 if e!=nil||chosen==nil||*chosen!=1||reason!="SELECT_NATIVE_LEARNED_PUBLIC_OBJECT_CONTINUITY"{
  t.Fatalf("not learned held-object continuity: %v %s %v",chosen,reason,e)
 }
 checkpoint,e:=p.Snapshot();if e!=nil{t.Fatal(e)}
 cold,e:=Restore(checkpoint);if e!=nil{t.Fatal(e)}
 second,e:=cold.Snapshot()
 if e!=nil||!bytes.Equal(checkpoint,second){t.Fatal("cold source checkpoint drift")}
 next,_,e:=cold.Choose(held)
 if e!=nil||next==nil||*next!=1{t.Fatal("cold native action changed")}
 if _,e=Restore(append(append([]byte{},checkpoint...),0));e==nil{t.Fatal("trailing checkpoint accepted")}
 if _,e=Restore(checkpoint[:len(checkpoint)-1]);e==nil{t.Fatal("truncated checkpoint accepted")}
}
func TestOriginalPublicMissionAndWrongFeedbackFailClosed(t *testing.T){
 m:=New();held:=observation(0,true)
 bad:=observation(0,true);bad.Mission="secret simulator waypoint"
 if a,_,e:=m.Choose(bad);e==nil||a!=nil{t.Fatal("private target accepted")}
 if e:=m.ObserveFit(bad,2,held,0,false);e==nil{t.Fatal("untrusted source FIT accepted")}
 if e:=m.Commit(held,2);e!=nil{t.Fatal(e)}
 if a,_,e:=m.Choose(held);e==nil||a!=nil{t.Fatal("reproposal while source pending")}
 if _,e:=m.Snapshot();e==nil{t.Fatal("checkpoint while source pending")}
 if e:=m.Feedback(held,2,held,0,false);e!=nil{t.Fatal(e)}
 if m.Samples!=1{t.Fatal("original source count drift")}
}
