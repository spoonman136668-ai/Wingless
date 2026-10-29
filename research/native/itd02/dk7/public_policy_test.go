package dk7

import (
 "bytes"
 "testing"
)
func pub()Public{
 p:=Public{Image:make([][][]int,7),Direction:1,Mission:Mission}
 for x:=0;x<7;x++{
  p.Image[x]=make([][]int,7)
  for y:=0;y<7;y++{p.Image[x][y]=[]int{1,0,0}}
 }
 p.Image[3][5]=[]int{5,0,0}
 p.Image[3][6]=[]int{1,0,0}
 return p
}
func TestSevenActionActualSelectedRewardAndNoUnselectedFeedback(t *testing.T){
 p:=New();s:=pub()
 if e:=p.ObserveFit(s,3,s,1.,false);e!=nil{t.Fatal(e)}
 if e:=p.ObserveFit(s,5,s,0.,false);e!=nil{t.Fatal(e)}
 a,status,e:=p.Choose(s)
 if e!=nil||a==nil||*a!=3||status!="SELECT_LEARNED_POSITIVE_ACTION"{t.Fatalf("choice %v %v %v",a,status,e)}
 if e=p.Commit(s,*a);e!=nil{t.Fatal(e)}
 if e=p.Feedback(s,5,s,1.,false);e==nil{t.Fatal("wrong selected action accepted")}
 if p.Samples!=2 {t.Fatal("unselected mutation")}
 if _,e=p.Snapshot();e==nil{t.Fatal("pending snapshot")}
 if e=p.Feedback(s,*a,s,1.,false);e!=nil{t.Fatal(e)}
 if p.Samples!=3||p.Live!=1{t.Fatal("selected count mismatch")}
 b,e:=p.Snapshot();if e!=nil{t.Fatal(e)}
 restored,e:=Restore(b);if e!=nil{t.Fatal(e)}
 cold,e:=restored.Snapshot();if e!=nil||!bytes.Equal(cold,b){t.Fatal("cold restore differs")}
 if _,e=Restore(append(append([]byte{},b...),1));e==nil{t.Fatal("tamper accepted")}
}
func TestAbstainAndPublicSchema(t *testing.T){
 p:=New();s:=pub()
 a,status,e:=p.Choose(s)
 if e!=nil||a!=nil||status!="ABSTAIN_NO_OBSERVED_POSITIVE_ACTION"{t.Fatal("must abstain on no public evidence")}
 s2:=pub();s2.Mission="secret world"
 if e=p.ObserveFit(s2,0,s,1.,false);e==nil{t.Fatal("private mission accepted")}
 if p.Samples!=0{t.Fatal("invalid selected reward mutated state")}
 s2=pub();s2.Image[3][5]=[]int{4,0,4}
 if e=p.ObserveFit(s2,0,s,1.,false);e==nil{t.Fatal("invalid original door state accepted")}
 if p.Samples!=0{t.Fatal("invalid source mutate")}
 if e=p.ObserveFit(s,7,s,1.,false);e==nil{t.Fatal("action 7 accepted")}
 if e=p.ObserveFit(s,0,s,-1.,false);e==nil{t.Fatal("negative reward")}
}
