package r6

import(
 "bytes"
 "os"
 "os/exec"
 "testing"
 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
)
func source(x,y,dir int)r1.PublicGrid{
 return r1.PublicGrid{X:x,Y:y,Direction:dir,GoalX:21,GoalY:20}
}
func feedback(t *testing.T,p *Policy,b r1.PublicGrid,action int,n r1.PublicGrid){
 t.Helper()
 if !p.valid(){t.Fatal("unqualified initial native route")}
 p.Core.Core.Route.Pending=true
 p.Core.Core.Route.Before=b
 p.Core.Core.Route.Action=action
 if e:=p.Feedback(b,action,n);e!=nil{t.Fatal(e)}
}
// Synthetic source-shaped fixtures are NOT independent MiniGrid evidence.
func fixture(t *testing.T,p *Policy){
 a:=source(20,20,0);b:=source(20,21,0)
 d:=source(21,21,0);c:=source(21,20,0)
 feedback(t,p,a,0,b)
 feedback(t,p,b,1,d)
 feedback(t,p,d,2,c)
 feedback(t,p,a,2,c)
}
func TestObservedShortestPathNotR5ImaginedPhysics(t *testing.T){
 p:=New();a:=source(20,20,0)
 if _,ok:=p.observedShortest(a);ok{t.Fatal("invented observed edge from untrained R3")}
 fixture(t,p)
 next,ok:=p.observedShortest(a)
 if !ok||next!=2{t.Fatalf("empirically shorter action not chosen: %d %v",next,ok)}
 decided,e:=p.Choose(a);if e!=nil||decided!=2{t.Fatalf("real verified path not chosen: %d %v",decided,e)}
 if _,e=p.Choose(a);e==nil{t.Fatal("two actions without a selected public consequence")}
 old:=p.Seen
 if e=p.Feedback(a,0,source(21,20,0));e==nil{t.Fatal("unselected feedback accepted")}
 if old!=p.Seen{t.Fatal("rejected event updated trusted edge counter")}
 if e=p.Feedback(a,2,source(21,20,0));e!=nil{t.Fatal(e)}
 if p.Seen!=old+1||p.Core.Core.Route.Steps!=p.Seen{t.Fatal("count conservation failed")}
}
func TestContradictionRemovesUnreliableObservedRoute(t *testing.T){
 p:=New();fixture(t,p)
 a:=source(20,20,0)
 // A genuine alternative public consequence invalidates the direct edge.
 feedback(t,p,a,2,source(20,20,1))
 e:=p.Edges[edgeKey(a,2)]
 if e.Count!=2||e.Contradictions!=1{t.Fatal("conflicting observed cause not recorded")}
 action,ok:=p.observedShortest(a)
 if !ok||action!=0{t.Fatalf("contradicted edge remained preferred: %d %v",action,ok)}
}
func TestRejectedNonlocalObservationDoesNotAlterNativeCustody(t *testing.T){
 p:=New()
 a:=source(20,20,0)
 p.Core.Core.Route.Pending=true;p.Core.Core.Route.Before=a;p.Core.Core.Route.Action=2
 snap,e:=p.Snapshot();if e!=nil{t.Fatal(e)}
 if e=p.Feedback(a,2,source(25,25,0));e==nil{t.Fatal("nonlocal real consequence admitted")}
 again,e:=p.Snapshot()
 if e!=nil||!bytes.Equal(snap,again){t.Fatal("bad action changed source checkpoint")}
}
func TestFreshWorldClearsEpisodeGraphAndKeepsLearnedCore(t *testing.T){
 p:=New();fixture(t,p)
 previous,e:=p.Core.Core.Core.Snapshot();if e!=nil{t.Fatal(e)}
 if e=p.BeginEpisode();e!=nil{t.Fatal(e)}
 if p.Seen!=0||len(p.Edges)!=0||p.Core.Core.Route.Steps!=0{
  t.Fatal("observed local graph contaminated next world")
 }
 after,e:=p.Core.Core.Core.Snapshot()
 if e!=nil||!bytes.Equal(previous,after){t.Fatal("native R3 model changed on session reset")}
}
func TestAuthenticatedColdOSStateAndTamper(t *testing.T){
 if path:=os.Getenv("R6_TEST_COLD");path!=""{
  raw,e:=os.ReadFile(path);if e!=nil{t.Fatal(e)}
  p,e:=Restore(raw);if e!=nil{t.Fatal(e)}
  if p.Seen!=4||len(p.Edges)!=4||p.Core.Core.Route.Steps!=4{
   t.Fatal("cold process did not retain selected public edge history")
  }
  action,ok:=p.observedShortest(source(20,20,0))
  if !ok||action!=2{t.Fatal("cold process lost shortest actual source route")}
  return
 }
 p:=New();fixture(t,p)
 original,e:=p.Snapshot();if e!=nil{t.Fatal(e)}
 q,e:=Restore(original);if e!=nil{t.Fatal(e)}
 repeated,e:=q.Snapshot()
 if e!=nil||!bytes.Equal(original,repeated){t.Fatal("source checkpoint non-deterministic")}
 if _,e=Restore(append(append([]byte{},original...),byte('X')));e==nil{
  t.Fatal("trailing ciphertext accepted")
 }
 file:=t.TempDir()+"/real-source-shaped.json"
 if e=os.WriteFile(file,original,0600);e!=nil{t.Fatal(e)}
 run:=exec.Command(os.Args[0],"-test.run=^TestAuthenticatedColdOSStateAndTamper$")
 run.Env=append(os.Environ(),"R6_TEST_COLD="+file)
 if output,e:=run.CombinedOutput();e!=nil{t.Fatalf("real cold OS process failed %v %s",e,output)}
}
func TestNoFITUpdatesAfterPublicSession(t *testing.T){
 p:=New();fixture(t,p)
 old,e:=p.Core.Core.Core.Snapshot();if e!=nil{t.Fatal(e)}
 b:=source(20,20,0);n:=source(21,20,0)
 if e:=p.ObserveFit(b,2,n,0,true);e==nil{
  t.Fatal("public episode feedback became fitting source")
 }
 now,e:=p.Core.Core.Core.Snapshot()
 if e!=nil||!bytes.Equal(old,now){t.Fatal("evaluation rewrote FIT model")}
}
