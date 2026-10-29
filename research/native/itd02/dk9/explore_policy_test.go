package dk9
import("bytes";"testing";original "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk7")
func public()Public{
 p:=Public{Image:make([][][]int,7),Direction:2,Mission:original.Mission}
 for x:=0;x<7;x++ {p.Image[x]=make([][]int,7);for y:=0;y<7;y++{p.Image[x][y]=[]int{1,0,0}}}
 p.Image[3][5]=[]int{5,0,0}
 return p
}
func TestNativePublicLeastObservedAfterSelectedRealFeedback(t *testing.T){
 p:=New();before:=public()
 first,status,err:=p.Choose(before,"carrying_changed")
 if err!=nil||first==nil||status!="SELECT_NATIVE_LEAST_OBSERVED_ACTION"{t.Fatalf("unknown decision %v %s %v",first,status,err)}
 if err=p.Commit(before,*first);err!=nil{t.Fatal(err)}
 if err=p.Feedback(before,*first,before,0,false);err!=nil{t.Fatal(err)}
 second,status,err:=p.Choose(before,"carrying_changed")
 if err!=nil||second==nil||*second==*first||status!="SELECT_NATIVE_LEAST_OBSERVED_ACTION"{t.Fatalf("least-count not adapted %v %s %v",second,status,err)}
 if p.Samples!=1||p.Live!=1{t.Fatalf("wrong native update %d/%d",p.Samples,p.Live)}
 snap,err:=p.Snapshot();if err!=nil{t.Fatal(err)}
 cold,err:=Restore(snap);if err!=nil{t.Fatal(err)}
 cp,err:=cold.Snapshot();if err!=nil||!bytes.Equal(cp,snap){t.Fatal("not exact restored")}
 other,status,err:=cold.Choose(before,"carrying_changed")
 if err!=nil||other==nil||*other!=*second||status!="SELECT_NATIVE_LEAST_OBSERVED_ACTION"{t.Fatal("restored native selector changed")}
}
func TestPreservesOriginalObservedEffectPrecedenceAndRejectsPrivateState(t *testing.T){
 p:=New();b:=public();a:=public();a.Image[3][6]=[]int{5,0,0}
 if err:=p.ObserveFit(b,3,a,0,false);err!=nil{t.Fatal(err)}
 choice,status,err:=p.Choose(b,"carrying_changed")
 if err!=nil||choice==nil||*choice!=3||status!="SELECT_LEARNED_PUBLIC_ACTION"{t.Fatalf("observed cause lost %v %s %v",choice,status,err)}
 b.Mission="secret route"
 if _,_,err:=p.Choose(b,"carrying_changed");err==nil{t.Fatal("hidden mission accepted")}
}
