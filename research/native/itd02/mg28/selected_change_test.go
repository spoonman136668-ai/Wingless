package mg28
import ("bytes";"testing")
func public(mission string,heading int)Public{
 p:=Public{Image:make([][][]int,7),Direction:heading,Mission:mission}
 for i:=0;i<7;i++{p.Image[i]=make([][]int,7);for j:=0;j<7;j++{p.Image[i][j]=[]int{1,0,0}}}
 return p
}
func TestMG28OriginalSelectedPredictiveSignalIsolation(t *testing.T){
 b:=public("pick up the red ball",0)
 n:=public(b.Mission,1)
 model:=New(true);negative:=New(false)
 proposal,_,e:=model.Choose(b)
 if e!=nil||proposal==nil{t.Fatal("missing native source proposal",e)}
 chosen:=(*proposal+1)%7
 for _,m:=range []*Agent{model,negative}{
  for i:=0;i<3;i++{
   if e=m.ObserveFit(b,chosen,n,0,false);e!=nil{t.Fatal("source fit failed",e)}
  }
  key,e:=key(b);if e!=nil{t.Fatal(e)}
  if m.Cells[key][chosen].Count!=3||m.Cells[key][chosen].Changed!=3{
   t.Fatal("genuine original selected public changes not retained")
  }
 }
 // The disabled signal must never emit the opt-in action selection receipt.
 _,disabled,e:=negative.Choose(b)
 if e!=nil||disabled=="MG28_NATIVE_OWN_SELECTED_REAL_NEXT_PUBLIC_CHANGE"{
  t.Fatal("no-signal ablation used proposed action effect",e)
 }
 blue:=public("pick up the blue ball",0)
 k1,_:=key(b);k2,_:=key(blue)
 if k1==k2{t.Fatal("source mission names aliased")}
 if _,found:=model.Cells[k2];found{t.Fatal("selected effect leaked into another mission")}
 if e=model.BeginEpisode();e!=nil{t.Fatal(e)}
 if e=model.Commit(b,4);e!=nil{t.Fatal(e)}
 if e=model.Feedback(b,3,n,0,false);e==nil{t.Fatal("wrong selected source action learned")}
 if e=model.Feedback(b,4,public(blue.Mission,1),0,false);e==nil{t.Fatal("wrong mission learned")}
 if e=model.Feedback(b,4,n,0,false);e!=nil{t.Fatal("correct selected source rejected",e)}
 if e=model.Feedback(b,4,n,0,false);e==nil{t.Fatal("duplicate selected source learned")}
 raw,e:=model.Snapshot();if e!=nil{t.Fatal(e)}
 restored,e:=Restore(raw);if e!=nil{t.Fatal(e)}
 reencoded,e:=restored.Snapshot()
 if e!=nil||!bytes.Equal(raw,reencoded){t.Fatal("source cold snapshot changed",e)}
 if _,e=Restore(append(append([]byte{},raw...),byte('x')));e==nil{
  t.Fatal("accepted trailing altered source checkpoint")
 }
 if _,e=Restore([]byte("not-source"));e==nil{t.Fatal("accepted forged source bytes")}
}
func TestMG28NegativeNoChangeAndTerminalNoFabricatedSignal(t *testing.T){
 b:=public("pick up the red key",0)
 a:=New(true)
 for i:=0;i<4;i++{
  if e:=a.ObserveFit(b,5,b,0,false);e!=nil{t.Fatal(e)}
 }
 k,_:=key(b)
 if a.Cells[k][5].Changed!=0{t.Fatal("counted source no-op as changed")}
 if _,status,e:=a.Choose(b);e!=nil||status=="MG28_NATIVE_OWN_SELECTED_REAL_NEXT_PUBLIC_CHANGE"{
  t.Fatal("created false positive next-change signal",status,e)
 }
}
