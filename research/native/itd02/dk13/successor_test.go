package dk13
import ("bytes";"testing")
func public(dir int,front int)Public{
 p:=Public{Image:make([][][]int,7),Direction:dir,Mission:"use the key to open the door and then get to the goal"}
 for i:=0;i<7;i++{p.Image[i]=make([][]int,7);for j:=0;j<7;j++{p.Image[i][j]=[]int{1,0,0}}}
 if front!=1{p.Image[3][5]=[]int{front,0,0}}
 return p
}
func TestSourceReversalDisfavoredAndColdRestore(t *testing.T){
 agent:=New();a:=public(0,1);b:=public(0,4)
 if e:=agent.BeginEpisode();e!=nil{t.Fatal(e)}
 // Live observed A --action 5--> B --action 5--> A.
 if e:=agent.Commit(a,5);e!=nil{t.Fatal(e)}
 if e:=agent.Feedback(a,5,b,0,false);e!=nil{t.Fatal(e)}
 if e:=agent.Commit(b,5);e!=nil{t.Fatal(e)}
 if e:=agent.Feedback(b,5,a,0,false);e!=nil{t.Fatal(e)}
 // Once revisited, previously toggling action 5 should lose to an
 // action not yet tried from A; no action-code semantics are assumed.
 proposal,reason,e:=agent.Choose(a)
 if e!=nil||proposal==nil||*proposal==5{
  t.Fatalf("reversal not disfavored: %v %s %v",proposal,reason,e)
 }
 checkpoint,e:=agent.Snapshot();if e!=nil{t.Fatal(e)}
 cold,e:=Restore(checkpoint);if e!=nil{t.Fatal(e)}
 again,e:=cold.Snapshot()
 if e!=nil||!bytes.Equal(checkpoint,again){t.Fatal("source cold restore drift")}
 other,_,e:=cold.Choose(a)
 if e!=nil||other==nil||*other!=*proposal{t.Fatal("cold action changed")}
 if _,e=Restore(append(append([]byte{},checkpoint...),0));e==nil{t.Fatal("trailing corruption accepted")}
}
func TestNoEffectSourceEscapesRepeatedSelection(t *testing.T){
 p:=New();a:=public(0,1)
 if e:=p.BeginEpisode();e!=nil{t.Fatal(e)}
 seen:=map[int]bool{}
 for i:=0;i<9;i++{
  act,_,e:=p.Choose(a);if e!=nil||act==nil{t.Fatal(e)}
  seen[*act]=true
  if e=p.Commit(a,*act);e!=nil{t.Fatal(e)}
  if e=p.Feedback(a,*act,a,0,false);e!=nil{t.Fatal(e)}
 }
 if len(seen)<4{t.Fatalf("repeated action not escaped: %v",seen)}
 if p.Policy.Agent.Samples!=9||p.Policy.Motion[0].N>9{t.Fatal("source totals")}
 snap,e:=p.Snapshot();if e!=nil{t.Fatal(e)}
 if _,e=Restore(snap);e!=nil{t.Fatal(e)}
}
func TestWrongActionNoSourceMutationAndEpisodeReset(t *testing.T){
 p:=New();a:=public(0,1)
 if e:=p.BeginEpisode();e!=nil{t.Fatal(e)}
 if e:=p.Commit(a,2);e!=nil{t.Fatal(e)}
 if e:=p.Feedback(a,5,a,0,false);e==nil{t.Fatal("wrong source action accepted")}
 if p.Policy.Agent.Samples!=0||len(p.Attempts)!=0{t.Fatal("selected attempt contaminated")}
 if e:=p.Feedback(a,2,a,0,false);e!=nil{t.Fatal(e)}
 if e:=p.Feedback(a,2,a,0,false);e==nil{t.Fatal("duplicate accepted")}
 if e:=p.BeginEpisode();e!=nil{t.Fatal(e)}
 if len(p.Attempts)!=0||len(p.Visits)!=0||p.Policy.Agent.Samples!=1{t.Fatal("episode reset lost parent or leaked state")}
}
