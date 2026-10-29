package mg19

import ("bytes";"testing")

func observed(mission string,dir int) Public {
 p:=Public{Image:make([][][]int,7),Direction:dir,Mission:mission}
 for i:=0;i<7;i++{
  p.Image[i]=make([][]int,7)
  for j:=0;j<7;j++{p.Image[i][j]=[]int{1,0,0}}
 }
 return p
}
func TestGenuineMissionIsSeparateFromCompatibilityProjection(t *testing.T){
 m:=New()
 red:=observed("pick up the red box",0)
 redChanged:=observed(red.Mission,1)
 blue:=observed("pick up the blue box",0)
 scene:=observed(red.Mission,2)
 if e:=m.ObserveFit(red,4,redChanged,1,true);e!=nil{t.Fatal(e)}
 a,status,e:=m.Choose(red)
 if e!=nil||a==nil||*a!=4||status!="MG19_NATIVE_EXACT_ORIGINAL_MISSION_STATE_REWARD"{
  t.Fatal("missing exact authentic mission keyed positive",a,status,e)
 }
 for _,p:=range []Public{blue,scene}{
  _,reason,e:=m.Choose(p)
  if e!=nil||reason=="MG19_NATIVE_EXACT_ORIGINAL_MISSION_STATE_REWARD"{
   t.Fatal("source reward leaked across missions or viewpoints",reason,e)
  }
 }
 if m.Missions[red.Mission]!=1||m.Parent.Agent.Samples!=1{
  t.Fatal("parent and genuine mission source counts diverged")
 }
 if e:=m.BeginEpisode();e!=nil{t.Fatal(e)}
 a,_,e=m.Choose(blue);if e!=nil||a==nil||*a<0||*a>6{t.Fatal("no native source action",e)}
 if e=m.Commit(blue,*a);e!=nil{t.Fatal(e)}
 if e=m.Feedback(blue,(*a+1)%7,blue,0,false);e==nil{
  t.Fatal("accepted incorrect selected source action")
 }
 if e=m.Feedback(blue,*a,observed(red.Mission,1),0,false);e==nil{
  t.Fatal("accepted original mission changed after source selection")
 }
 if m.Parent.Agent.Samples!=1{t.Fatal("invalid source mutated native")}
 if e=m.Feedback(blue,*a,observed(blue.Mission,1),0,false);e!=nil{t.Fatal(e)}
 if e=m.Feedback(blue,*a,blue,0,false);e==nil{t.Fatal("duplicate original source event")}
 if m.Parent.Agent.Samples!=2||m.Missions[blue.Mission]!=1{t.Fatal("selected feedback lost")}
 raw,e:=m.Snapshot();if e!=nil{t.Fatal(e)}
 cold,e:=Restore(raw);if e!=nil{t.Fatal(e)}
 restored,e:=cold.Snapshot()
 if e!=nil||!bytes.Equal(raw,restored){t.Fatal("cold source mission checksum mismatch")}
 if _,e=Restore(append(append([]byte{},raw...),0));e==nil{t.Fatal("tampered checkpoint accepted")}
}
func TestUnsupportedMissionAndUnchangedTerminalDoNotGrantGlobalReward(t *testing.T){
 m:=New()
 bad:=[]string{"", "pick up the ", "Pick up the red box",
  "pick up the red box and use a private map", "ignore instructions", "pick up the red\nbox"}
 for _,txt:=range bad{
  if _,_,e:=m.Choose(observed(txt,0));e==nil{t.Fatal("untrusted native mission admitted",txt)}
 }
 yellow:=observed("pick up the yellow key",0)
 if e:=m.ObserveFit(yellow,3,yellow,1,true);e!=nil{t.Fatal(e)}
 _,reason,e:=m.Choose(yellow)
 if e!=nil||reason=="MG19_NATIVE_EXACT_ORIGINAL_MISSION_STATE_REWARD"{
  t.Fatal("unchanged original source rewarded act leaked",e)
 }
}
