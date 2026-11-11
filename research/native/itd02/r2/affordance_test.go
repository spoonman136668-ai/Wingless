package r2
import(
 "bytes"
 "encoding/base64"
 "os"
 "os/exec"
 "strconv"
 "testing"
 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
)
func sample(blocked bool)r1.PublicGrid{return r1.PublicGrid{X:3,Y:4,GoalX:6,GoalY:7,Direction:0,FrontBlocked:blocked}}
func TestNoUnobservedHardcodedAvoidance(t *testing.T){
 p:=New()
 s:=sample(true)
 got,e:=p.Scores(s);if e!=nil{t.Fatal(e)}
 base,e:=p.Core.Scores(s);if e!=nil||got!=base{t.Fatal("untrained prior isn't original source")}
 if p.BlockedAttempts!=0{t.Fatal("not native first use")}
}
func TestActualEightBlockedCausalEventsLearnPenalty(t *testing.T){
 p:=New();s:=sample(true);u:=sample(false)
 for i:=0;i<8;i++{if e:=p.Observe(s,2,s,0,true);e!=nil{t.Fatal(e)}}
 scores,e:=p.Scores(s);if e!=nil{t.Fatal(e)}
 standard,e:=p.Core.Scores(s);if e!=nil{t.Fatal(e)}
 if standard[2]-scores[2]<6{t.Fatal("failed causal affordance learning")}
 a,e:=p.Choose(s);if e!=nil||a==2{t.Fatalf("still chooses blocked action: %d",a)}
 before,e:=p.Scores(u);if e!=nil{t.Fatal(e)}
 original,e:=p.Core.Scores(u);if e!=nil||before!=original{t.Fatal("penalty incorrectly generalized to unblocked movement")}
 if p.BlockedAttempts!=8||p.BlockedProgress!=0{t.Fatal("wrong source counts")}
}
func TestOnlySelectedActionAndTransactionalInvalid(t *testing.T){
 p:=New();s:=sample(true)
 for i:=0;i<8;i++{if e:=p.Observe(s,0,s,0,true);e!=nil{t.Fatal(e)}}
 if p.BlockedAttempts!=0{t.Fatal("turn updates blocked forward outcome")}
 old,_:=p.Snapshot()
 if e:=p.Observe(s,3,s,0,true);e==nil{t.Fatal("illegal 4th action accepted")}
 now,_:=p.Snapshot();if !bytes.Equal(old,now){t.Fatal("bad update mutated original native state")}
}
func TestAuthenticationAndColdRestart(t *testing.T){
 if b:=os.Getenv("K4Y_GO_CHILD");b!=""{
  raw,e:=base64.StdEncoding.DecodeString(b);if e!=nil{t.Fatal(e)}
  p,e:=Restore(raw);if e!=nil{t.Fatal(e)}
  a,e:=p.Choose(sample(true));if e!=nil{t.Fatal(e)}
  want,_:=strconv.Atoi(os.Getenv("K4Y_GO_WANT"))
  if a!=want||p.BlockedAttempts!=8{t.Fatal("restored true blocked state lost")}
  return
 }
 p:=New();s:=sample(true)
 for i:=0;i<8;i++{if e:=p.Observe(s,2,s,0,true);e!=nil{t.Fatal(e)}}
 raw,e:=p.Snapshot();if e!=nil{t.Fatal(e)}
 restored,e:=Restore(raw);if e!=nil{t.Fatal(e)}
 a,e:=p.Choose(s);if e!=nil{t.Fatal(e)}
 b,e:=restored.Choose(s);if e!=nil||a!=b{t.Fatal("restored prediction drift")}
 bad:=bytes.Replace(raw,[]byte("\"blocked_attempts\":8"),[]byte("\"blocked_attempts\":7"),1)
 if bytes.Equal(bad,raw){t.Fatal("tamper fixture missed")}
 if _,e=Restore(bad);e==nil{t.Fatal("tampered counts accepted")}
 cmd:=exec.Command(os.Args[0],"-test.run=^TestAuthenticationAndColdRestart$")
 cmd.Env=append(os.Environ(),"K4Y_GO_CHILD="+base64.StdEncoding.EncodeToString(raw),"K4Y_GO_WANT="+strconv.Itoa(a))
 result,e:=cmd.CombinedOutput();if e!=nil{t.Fatalf("cold OS restore: %s %v",result,e)}
}
func TestSourceIsNotNavigationSuccess(t *testing.T){if Schema==""{t.Fatal("missing schema")}}
