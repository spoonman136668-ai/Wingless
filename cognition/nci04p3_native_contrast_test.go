package cognition

// NCI-04P3: a single project-native and reversible source intervention.
// The witness is the actual existing ICE TestSourceWins Go assertion executed
// in a detached, disposable checkout. It is NOT unseen scientific evaluation,
// a general action policy, or proof of learned unfamiliar-task transfer.
import (
 "bytes"
 "context"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "os"
 "os/exec"
 "path/filepath"
 "strings"
 "testing"
 "time"
)

func TestNCI04P3ActualICECausalInterventionAndReversal(t *testing.T) {
 root,err:=filepath.Abs("..");if err!=nil{t.Fatal(err)}
 exact:=exec.Command("git","-C",root,"rev-parse","HEAD")
 ref,err:=exact.Output();if err!=nil{t.Fatal(err)}
 sha:=strings.TrimSpace(string(ref))
 if len(sha)!=40{t.Fatal("invalid original source identity")}
 checkout:=filepath.Join(t.TempDir(),"isolated-source")
 setup:=exec.Command("git","clone","--quiet","--no-hardlinks","--no-checkout",root,checkout)
 if out,e:=setup.CombinedOutput();e!=nil{t.Fatalf("source isolation unavailable: %v %s",e,out)}
 git:=exec.Command("git","-C",checkout,"checkout","--quiet","--detach",sha)
 if out,e:=git.CombinedOutput();e!=nil{t.Fatalf("pinned source unavailable: %v %s",e,out)}
 source:=filepath.Join(checkout,"ice","index.go")
 original,err:=os.ReadFile(source);if err!=nil{t.Fatal(err)}
 if bytes.Count(original,[]byte("if string(a) != string(b) {"))!=1{
  t.Fatal("native stale-source validator changed; no intervention authorized")
 }
 desc:=nativeReasoningTest{Name:"TestSourceWins",Package:"./ice",Source:"ice/index_test.go"}
 run:=func()(int,[]byte,time.Duration){
  t.Helper()
  ctx,cancel:=context.WithTimeout(context.Background(),35*time.Second)
  defer cancel()
  cmd:=exec.CommandContext(ctx,"go","test","-json","./ice","-run","^TestSourceWins$","-count=1")
  cmd.Dir=checkout
  cmd.Env=append(os.Environ(),"GOFLAGS=","GOWORK=off","GOPROXY=off","GOSUMDB=off","GOTOOLCHAIN=local")
  start:=time.Now()
  output,runErr:=cmd.CombinedOutput()
  elapsed:=time.Since(start)
  if ctx.Err()!=nil{t.Fatalf("native execution timed out: %v",ctx.Err())}
  label,e:=nativeGoTestVerdict(output,desc,runErr)
  if e!=nil{t.Fatalf("native result not qualified as assertion outcome: %v\n%s",e,output)}
  return label,output,elapsed
 }
 first,_,initialTime:=run()
 if first!=1{t.Fatal("original native ICE test not passing")}
 before:=sha256.Sum256(original)
 patched:=bytes.Replace(original,[]byte("if string(a) != string(b) {"),[]byte("if false {"),1)
 if err:=os.WriteFile(source,patched,0600);err!=nil{t.Fatal(err)}
 diff:=exec.Command("git","-C",checkout,"diff","--name-only")
 names,err:=diff.Output();if err!=nil||strings.TrimSpace(string(names))!="ice/index.go"{
  t.Fatalf("unbounded action source mutation: %q %v",names,err)
 }
 if bytes.Equal(original,patched){t.Fatal("native intervention had no effect on source bytes")}
 damaged,output,interventionTime:=run()
 if damaged!=0{t.Fatal("actual native validator intervention did not produce a distinct test failure")}
 if !bytes.Contains(output,[]byte("stale index permitted")){
  t.Fatal("the native assertion failed for an unrelated reason; no causal witness admitted")
 }
 if err:=os.WriteFile(source,original,0600);err!=nil{t.Fatal(err)}
 after,err:=os.ReadFile(source);if err!=nil{t.Fatal(err)}
 digest:=sha256.Sum256(after)
 if !bytes.Equal(original,after)||before!=digest{t.Fatal("isolated source restoration changed original bytes")}
 repaired,_,restoreTime:=run()
 if repaired!=1{t.Fatal("restored original native source did not recover its task outcome")}
 clean:=exec.Command("git","-C",checkout,"status","--porcelain")
 dirty,err:=clean.Output();if err!=nil||len(bytes.TrimSpace(dirty))!=0{t.Fatalf("intervention leaked beyond isolated source: %q %v",dirty,err)}
 evidence:=map[string]any{
  "gate":"NCI-04P3-W",
  "repository_commit":sha,
  "source_file":"ice/index.go",
  "original_source_sha256":hex.EncodeToString(before[:]),
  "observed_native_test_results":[]int{first,damaged,repaired},
  "independent_witness":"original Go ICE stale-source protection assertion, not caller label",
  "operational_failures_mislabeled":0,
  "action_isolation":"disposable detached checkout, one original source file",
  "native_test_execution_wall_ns":[]int64{initialTime.Nanoseconds(),interventionTime.Nanoseconds(),restoreTime.Nanoseconds()},
  "new_cognitive_parameters":0,
  "source_changes_reversed":true,
  "native_learner_consumed_contrast":false,
  "novel_task_transfer":"NOT_TESTED",
  "matched_cpu_rss":"NOT_QUALIFIED",
  "independent_scientific_first_use":false,
  "classification":"AUTHENTIC_NATIVE_CONTRAST_DEVELOPMENT_PASS / LEARNING_TRANSFER_NO_GO",
 }
 encoded,e:=json.Marshal(evidence);if e!=nil{t.Fatal(e)}
 t.Logf("NCI04P3_REAL_NATIVE_CONTRAST=%s",encoded)
}

func TestNCI04P3NativeAssertionRejectsOperationalFailure(t *testing.T) {
 test:=nativeReasoningTest{Name:"TestSourceWins",Package:"./ice",Source:"ice/index_test.go"}
 raw:=[]byte(`{"Action":"fail","Package":"github.com/spoonman136668-ai/Wingless/ice"}`+"\n")
 if _,err:=nativeGoTestVerdict(raw,test,errors.New("build failure"));err==nil{
  t.Fatal("unexecuted build failure admitted as causal label")
 }
}
