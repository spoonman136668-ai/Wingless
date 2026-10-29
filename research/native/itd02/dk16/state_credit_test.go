package dk16
import ("bytes";"testing")
func sample(mission string,dir int)Public{
 p:=Public{Image:make([][][]int,7),Direction:dir,Mission:mission}
 for i:=0;i<7;i++{p.Image[i]=make([][]int,7);for j:=0;j<7;j++{p.Image[i][j]=[]int{1,0,0}}}
 return p
}
func TestFullSourcePublicStateAndMissionCreditIsolated(t *testing.T){
 m:=New()
 a:=sample("use the key to open the door and then get to the goal",0)
 b:=sample(a.Mission,1)
 changed:=sample(a.Mission,2)
 other:=sample("get to the green goal square",0)
 // Exact observed success is supported only from state a; no action semantics.
 if e:=m.ObserveFit(a,4,b,1,true);e!=nil{t.Fatal(e)}
 chosen,reason,e:=m.Choose(a)
 if e!=nil||chosen==nil||*chosen!=4||
  reason!="SELECT_NATIVE_EXACT_PUBLIC_REWARDED_SOURCE"{t.Fatal("positive state omitted",chosen,reason,e)}
 for _,p:=range []Public{changed,other}{
  x,r,e:=m.Choose(p)
  legacy:=p;legacy.Mission="use the key to open the door and then get to the goal"
  baseline,why,eb:=m.Policy.Policy.Choose(legacy)
  if e!=nil||eb!=nil||x==nil||baseline==nil||*x!=*baseline||r!=why{
   t.Fatal("mission-global contamination",r,why,e,eb)
  }
 }
 if m.Agent.Samples!=1{t.Fatal("inherited selected total")}
 raw,e:=m.Snapshot();if e!=nil{t.Fatal(e)}
 cold,e:=Restore(raw);if e!=nil{t.Fatal(e)}
 again,e:=cold.Snapshot()
 if e!=nil||!bytes.Equal(raw,again){t.Fatal("cold exact context drift")}
 if _,e=Restore(append(append([]byte{},raw...),0));e==nil{t.Fatal("tampered source accepted")}
}
func TestUnchangedSceneAndWrongSelectedFeedbackFailClosed(t *testing.T){
 m:=New();a:=sample("get to the green goal square",0)
 if e:=m.ObserveFit(a,2,a,1,true);e!=nil{t.Fatal(e)}
 chosen,reason,e:=m.Choose(a)
 legacy:=a;legacy.Mission="use the key to open the door and then get to the goal"
 base,why,eb:=m.Policy.Policy.Choose(legacy)
 if e!=nil||eb!=nil||*chosen!=*base||reason!=why{t.Fatal("unchanged scene source reward leaked")}
 if e=m.BeginEpisode();e!=nil{t.Fatal(e)}
 if e=m.Commit(a,1);e!=nil{t.Fatal(e)}
 mismatch:=sample("use the key to open the door and then get to the goal",0)
 if e=m.Feedback(a,1,mismatch,1,true);e==nil{t.Fatal("changed original mission accepted")}
 if e=m.Feedback(a,2,a,0,false);e==nil{t.Fatal("wrong selected source accepted")}
 if m.Agent.Samples!=1{t.Fatal("failed event modified source")}
 if e=m.Feedback(a,1,a,0,false);e!=nil{t.Fatal(e)}
 if e=m.Feedback(a,1,a,0,false);e==nil{t.Fatal("duplicate source accepted")}
 if m.Agent.Samples!=2{t.Fatal("duplicate modified count")}
}
