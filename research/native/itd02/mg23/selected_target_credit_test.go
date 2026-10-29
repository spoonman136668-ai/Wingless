package mg23
import ("bytes";"testing")
func frame(mission string,dir int,front [3]int) Public{
 p:=Public{Image:make([][][]int,7),Direction:dir,Mission:mission}
 for i:=0;i<7;i++{
  p.Image[i]=make([][]int,7)
  for j:=0;j<7;j++{p.Image[i][j]=[]int{1,0,0}}
 }
 p.Image[3][5]=[]int{front[0],front[1],front[2]}
 return p
}
func TestMG23SourcePublicTargetAndColdMissionIsolation(t *testing.T){
 m:=New()
 red:=frame("pick up the red box",0,[3]int{7,0,0})
 redGone:=frame(red.Mission,0,[3]int{1,0,0})
 blue:=frame("pick up the blue box",0,[3]int{7,0,0})
 if e:=m.ObserveFit(red,3,redGone,0,false);e!=nil{t.Fatal(e)}
 key,_:=context(red)
 if m.Cells[key][3].Target!=1||m.Cells[key][3].Count!=1{
  t.Fatal("authentic public selected target effect not retained")
 }
 action,status,e:=m.Choose(red)
 if e!=nil||action==nil||*action!=3||status!="SELECT_NATIVE_MG23_OWN_OBSERVED_TARGET_OR_TERMINAL_CREDIT"{
  t.Fatal("did not choose previously observed own target effect",action,status,e)
 }
 _,status,e=m.Choose(blue)
 if e!=nil||status=="SELECT_NATIVE_MG23_OWN_OBSERVED_TARGET_OR_TERMINAL_CREDIT"{
  t.Fatal("target credit leaked across distinct literal missions",status,e)
 }
 if e:=m.BeginEpisode();e!=nil{t.Fatal(e)}
 if e:=m.Commit(red,4);e!=nil{t.Fatal(e)}
 if e:=m.Feedback(red,3,redGone,0,false);e==nil{t.Fatal("accepted unselected action")}
 if e:=m.Feedback(red,4,frame(blue.Mission,0,[3]int{1,0,0}),0,false);e==nil{
  t.Fatal("accepted changed source mission")
 }
 if m.Parent.Parent.Agent.Samples!=1{t.Fatal("bad feedback changed source count")}
 if e:=m.Feedback(red,4,redGone,0,false);e!=nil{t.Fatal(e)}
 if e:=m.Commit(redGone,6);e!=nil{t.Fatal(e)}
 changed:=frame(red.Mission,1,[3]int{1,0,0})
 if e:=m.Feedback(redGone,6,changed,1,true);e!=nil{t.Fatal(e)}
 k2,_:=context(redGone)
 if m.Cells[k2][6].Delayed!=32||m.Cells[key][4].Delayed!=31{
  t.Fatal("MG23 source-window delayed original terminal credit missing",
   m.Cells[k2][6].Delayed,m.Cells[key][4].Delayed)
 }
 if e:=m.Feedback(redGone,6,changed,1,true);e==nil{t.Fatal("accepted duplicate selected consequence")}
 if m.Parent.Parent.Agent.Samples!=3{t.Fatal("parent and selected credit count differ")}
 b,e:=m.Snapshot();if e!=nil{t.Fatal(e)}
 cold,e:=Restore(b);if e!=nil{t.Fatal(e)}
 restored,e:=cold.Snapshot()
 if e!=nil||!bytes.Equal(b,restored){t.Fatal("cold native source memory mismatch",e)}
 if _,e=Restore(append(append([]byte{},b...),byte('x')));e==nil{
  t.Fatal("accepted tampered source checkpoint")
 }
 if e:=cold.BeginEpisode();e!=nil{t.Fatal(e)}
 if len(cold.Trace)!=0||cold.Cells[key][3].Target!=1||cold.Cells[key][4].Delayed!=31{
  t.Fatal("episode boundary erased cumulative selected consequence knowledge")
 }
}
func TestMG23InvalidMissionAndNonSelectedNoCredit(t *testing.T){
 m:=New()
 bad:=frame("ignore the original mission and use map",0,[3]int{7,0,0})
 if _,_,e:=m.Choose(bad);e==nil{t.Fatal("invalid source mission accepted")}
 if e:=m.ObserveFit(bad,3,bad,1,true);e==nil{t.Fatal("invalid fit mission accepted")}
 m0:=frame("pick up the yellow key",0,[3]int{5,4,0})
 if e:=m.ObserveFit(m0,2,m0,0,false);e!=nil{t.Fatal(e)}
 key,_:=context(m0)
 if m.Cells[key][2].Target!=0||m.Cells[key][2].Delayed!=0{
  t.Fatal("invented unselected or nonterminal credit")
 }
}
