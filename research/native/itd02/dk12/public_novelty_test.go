package dk12
import ("bytes";"testing")
func public()Public{
 p:=Public{Image:make([][][]int,7),Direction:1,Mission:"use the key to open the door and then get to the goal"}
 for x:=0;x<7;x++{p.Image[x]=make([][]int,7);for y:=0;y<7;y++{p.Image[x][y]=[]int{1,0,0}}}
 return p
}
func TestPublicOwnSelectedNoveltyAfterSevenSameSourceStates(t *testing.T){
 p:=New();b:=public()
 if e:=p.BeginEpisode();e!=nil{t.Fatal(e)}
 distinct:=map[int]bool{}
 for i:=0;i<7;i++{
  action,status,e:=p.Choose(b)
  if e!=nil||action==nil||status!="SELECT_NATIVE_UNTRIED_PUBLIC_STATE_ACTION"{t.Fatalf("selection %v %s %v",action,status,e)}
  if distinct[*action]{t.Fatal("repeated source action before trying all seven")}
  distinct[*action]=true
  if e=p.Commit(b,*action);e!=nil{t.Fatal(e)}
  if e=p.Feedback(b,*action,b,0,false);e!=nil{t.Fatal(e)}
 }
 if len(distinct)!=7||p.Steps!=7||p.Samples!=7||p.Live!=7{t.Fatal("source selected action budget not conserved")}
 chosen,status,e:=p.Choose(b)
 if e!=nil||chosen==nil||status!="SELECT_NATIVE_OBSERVED_PUBLIC_NOVELTY"{t.Fatalf("failed learned episode path %v %s %v",chosen,status,e)}
 snap,e:=p.Snapshot();if e!=nil{t.Fatal(e)}
 cold,e:=Restore(snap);if e!=nil{t.Fatal(e)}
 next,_,e:=cold.Choose(b);if e!=nil||next==nil||*next!=*chosen{t.Fatal("cold native source action drift")}
 cp,e:=cold.Snapshot();if e!=nil||!bytes.Equal(cp,snap){t.Fatal("cold checkpoint not exact")}
 if _,e=Restore(append(append([]byte{},snap...),1));e==nil{t.Fatal("accepted corrupt snapshot")}
}
func TestRejectPrivateInputAndWrongActualSelectedFeedback(t *testing.T){
 p:=New();b:=public();bad:=public();bad.Mission="private simulator global position"
 if a,_,e:=p.Choose(bad);e==nil||a!=nil{t.Fatal("accepted private mission")}
 if e:=p.Commit(b,0);e!=nil{t.Fatal(e)}
 if e:=p.Feedback(b,1,b,0,false);e==nil{t.Fatal("accepted wrong source event")}
 if p.Steps!=0||p.Live!=0{t.Fatal("polluted public source memory")}
 if e:=p.Feedback(b,0,b,0,false);e!=nil{t.Fatal(e)}
 if e:=p.BeginEpisode();e!=nil{t.Fatal(e)}
 if p.Steps!=0||p.Samples!=1{t.Fatal("reset wiped lifelong selected samples or not episode memory")}
}
