package dk8
import ("bytes";"testing";original "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk7")
func view()Public{
 p:=Public{Image:make([][][]int,7),Direction:1,Mission:original.Mission}
 for x:=0;x<7;x++{p.Image[x]=make([][]int,7);for y:=0;y<7;y++{p.Image[x][y]=[]int{1,0,0}}}
 p.Image[3][5]=[]int{5,0,0}
 return p
}
func TestRealSelectedPublicCarryingAndDoorEffectCounts(t *testing.T){
 p:=New()
 before:=view();after:=view();after.Image[3][6]=[]int{5,0,0}
 if e:=p.ObserveFit(before,3,after,0,false);e!=nil{t.Fatal(e)}
 if e:=p.ObserveFit(before,0,before,0,false);e!=nil{t.Fatal(e)}
 choice,status,e:=p.Choose(before,Carrying)
 if e!=nil||choice==nil||*choice!=3||status!="SELECT_LEARNED_PUBLIC_ACTION"{t.Fatalf("carry action %v %s %v",choice,status,e)}
 if choice,status,e=p.Choose(before,Door);e!=nil||choice!=nil||status!="ABSTAIN_NO_LEARNED_POSITIVE_ACTION"{t.Fatalf("no door %v %s %v",choice,status,e)}
 doorBefore:=view();doorBefore.Image[3][5]=[]int{4,0,2}
 doorAfter:=view();doorAfter.Image[3][5]=[]int{4,0,0}
 if e=p.ObserveFit(doorBefore,5,doorAfter,0,false);e!=nil{t.Fatal(e)}
 choice,status,e=p.Choose(doorBefore,Door)
 if e!=nil||choice==nil||*choice!=5{t.Fatalf("door action %v %s %v",choice,status,e)}
 if e=p.Commit(doorBefore,5);e!=nil{t.Fatal(e)}
 if e=p.Feedback(doorBefore,3,doorAfter,0,false);e==nil{t.Fatal("unselected false action admitted")}
 if p.Samples!=3{t.Fatal("bad feedback modified original")}
 if _,e=p.Snapshot();e==nil{t.Fatal("pending episode snapshot accepted")}
 if e=p.Feedback(doorBefore,5,doorAfter,0,false);e!=nil{t.Fatal(e)}
 if p.Samples!=4||p.Live!=1{t.Fatal("incorrect exactly-once effect learning")}
 snap,e:=p.Snapshot();if e!=nil{t.Fatal(e)}
 cold,e:=Restore(snap);if e!=nil{t.Fatal(e)}
 again,e:=cold.Snapshot();if e!=nil||!bytes.Equal(snap,again){t.Fatal("cold restore different")}
 if _,e=Restore(append(append([]byte{},snap...),42));e==nil{t.Fatal("tampered canonical snapshot")}
}
func TestSourceOnlyNoHiddenStateAndTrulyUnknownAbstains(t *testing.T){
 m:=New();p:=view()
 result,status,e:=m.Choose(p,Carrying)
 if e!=nil||result!=nil||status!="ABSTAIN_NO_LEARNED_POSITIVE_ACTION"{t.Fatal("unknown action must abstain")}
 p2:=view();p2.Mission="world-answer"
 if e=m.ObserveFit(p2,3,p,0,false);e==nil{t.Fatal("hidden field")}
 if e=m.ObserveFit(p,7,p,0,false);e==nil{t.Fatal("illegal action")}
 if e=m.ObserveFit(p,3,p,-1,false);e==nil{t.Fatal("negative reward accepted")}
 if _,_,e=m.Choose(p,"positive_reward");e==nil{t.Fatal("unregistered objective accepted")}
 if m.Samples!=0{t.Fatal("invalid source mutated memory")}
}
