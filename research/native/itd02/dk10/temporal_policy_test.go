package dk10
import ("bytes";"testing";parent "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk7")
func pub()Public{
 p:=Public{Image:make([][][]int,7),Direction:1,Mission:parent.Mission}
 for x:=0;x<7;x++{p.Image[x]=make([][]int,7);for y:=0;y<7;y++{p.Image[x][y]=[]int{1,0,0}}}
 p.Image[3][5]=[]int{5,0,0};return p
}
func TestOriginalActualPublicSelectedLoopMemoryAndCheckpoint(t *testing.T){
 m:=New();p:=pub()
 if e:=m.ObserveFit(p,3,p,0,false);e!=nil{t.Fatal(e)}
 if e:=m.BeginEpisode();e!=nil{t.Fatal(e)}
 first,s,e:=m.Choose(p,"carrying_changed")
 if e!=nil||first==nil||s!="SELECT_NATIVE_PUBLIC_TEMPORAL_ACTION"{t.Fatalf("choice %v %s %v",first,s,e)}
 if e=m.Commit(p,*first);e!=nil{t.Fatal(e)}
 if e=m.Feedback(p,*first,p,0,false);e!=nil{t.Fatal(e)}
 if m.Steps!=1||m.Live!=1||m.Samples!=2{t.Fatal("not real selected source credit")}
 if len(m.Visits)!=1||len(m.Seen)!=1{t.Fatal("not storing public state transitions")}
 second,s,e:=m.Choose(p,"carrying_changed")
 if e!=nil||second==nil||*second==*first{t.Fatalf("did not avoid observed action loop %v %s %v",second,s,e)}
 b,e:=m.Snapshot();if e!=nil{t.Fatal(e)}
 copy,e:=Restore(b);if e!=nil{t.Fatal(e)}
 cb,e:=copy.Snapshot();if e!=nil||!bytes.Equal(cb,b){t.Fatal("cold original native source drift")}
 next,_,e:=copy.Choose(p,"carrying_changed")
 if e!=nil||next==nil||*next!=*second{t.Fatal("cold selector changes source behavior")}
 if _,e=Restore(append(append([]byte{},b...),1));e==nil{t.Fatal("invalid original checkpoint passed")}
}
func TestSelectedActionMismatchAndInvalidPrivateWorld(t *testing.T){
 m:=New();p:=pub()
 if e:=m.BeginEpisode();e!=nil{t.Fatal(e)}
 if e:=m.Commit(p,4);e!=nil{t.Fatal(e)}
 if e:=m.Feedback(p,5,p,0,false);e==nil{t.Fatal("selected-action spoof accepted")}
 if m.Steps!=0||m.Live!=0{t.Fatal("invalid event mutated memory")}
 if _,e:=m.Snapshot();e==nil{t.Fatal("snapshot allowed outstanding selected action")}
 q:=pub();q.Mission="secret answer"
 if _,_,e:=m.Choose(q,"carrying_changed");e==nil{t.Fatal("invalid mission")}
}
