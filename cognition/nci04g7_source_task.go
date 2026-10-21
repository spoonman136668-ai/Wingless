package cognition

// NCI-04G7: an opt-in, actual project-native source-reasoning task outcome.
// This is not the prior generated single-file probe. The consequence is the
// exit status of an executed, specifically allowlisted Wingless test sourced
// from the exact validated ICE index and git commit. Failed executions are
// errors, NOT an invented negative supervision label.
import (
 "context"
 "bufio"
 "bytes"
 "encoding/json"
 "crypto/sha256"
 "encoding/hex"
 "errors"
 "fmt"
 "os"
 "os/exec"
 "path/filepath"
 "strings"
 "time"
)

type NativeReasoningConsequence struct{
 TaskName string `json:"task_name"`
 SourceCommit string `json:"source_commit"`
 SourceFileSHA256 string `json:"source_file_sha256"`
 IndexSHA string `json:"ice_manifest_sha256"`
 Action uint8 `json:"action"`
 ObservedPass uint8 `json:"observed_pass"`
 EvidenceBasis string `json:"evidence_basis"`
}

type nativeReasoningTest struct{
 Name,Package,Source string
}
var nativeReasoningTasks=[3]nativeReasoningTest{
 {"TestSourceWins","./ice","ice/index_test.go"},
 {"TestICEStateTaskRoundtripAndStrongExactControl","./cognition","cognition/ice_task_test.go"},
 {"TestICEStateAmbiguityAndNoInventedNegatives","./cognition","cognition/ice_task_test.go"},
}

// nativeGoTestVerdict admits only a completed Go test's own assertion result.
// A compile error, timeout, panic, skipped test, missing exact test event, or
// arbitrary process output is NOT a negative learning label. The verifier
// owns the observed consequence; the caller never supplies an answer.
func nativeGoTestVerdict(raw []byte, task nativeReasoningTest, executionErr error)(int,error){
 if len(raw)==0||len(raw)>1<<20{return 0,errors.New("unbounded or empty native test evidence")}
 expectedPackage:="github.com/spoonman136668-ai/Wingless/"+strings.TrimPrefix(task.Package,"./")
 scanner:=bufio.NewScanner(bytes.NewReader(raw))
 scanner.Buffer(make([]byte,4096),1<<20)
 var started,finished,packageFinished bool
 var result int
 for scanner.Scan(){
  var event struct{
   Action string `json:"Action"`
   Package string `json:"Package"`
   Test string `json:"Test"`
   Output string `json:"Output"`
  }
  if err:=json.Unmarshal(scanner.Bytes(),&event);err!=nil{
   return 0,errors.New("non-JSON native process output")
  }
  if event.Package!=expectedPackage{return 0,errors.New("native test package identity drift")}
  if event.Test!=""&&event.Test!=task.Name{return 0,errors.New("additional test outcome or nested subtest not admitted")}
  if strings.Contains(event.Output,"panic:")||
    strings.Contains(event.Output,"fatal error:")||
    strings.Contains(event.Output,"signal:")||
    strings.Contains(event.Output,"runtime error:"){
    return 0,errors.New("native test failure is operational, not an assertion label")
  }
  if event.Test==task.Name{
   switch event.Action{
   case "run":
    if started||finished{return 0,errors.New("duplicate native test start")}
    started=true
   case "pass","fail":
    if !started||finished{return 0,errors.New("native test result out of sequence")}
    finished=true
    if event.Action=="pass"{result=1}
   case "skip":
    return 0,errors.New("skipped native test cannot supervise cognition")
   }
  }else if event.Test==""&&(event.Action=="pass"||event.Action=="fail"){
   if !finished||packageFinished{return 0,errors.New("native package terminated without isolated test result")}
   packageFinished=true
   if (event.Action=="pass")!=(result==1){return 0,errors.New("native package and test outcomes conflict")}
  }
 }
 if err:=scanner.Err();err!=nil{return 0,err}
 if !started||!finished||!packageFinished{return 0,errors.New("missing completed native assertion receipt")}
 if (executionErr==nil)!=(result==1){return 0,errors.New("native process exit and assertion disagree")}
 return result,nil
}

// ObserveNativeReasoningTask learns from a real, source-checked Wingless
// reasoning/ICE test execution. No caller-supplied labels or arbitrary test
// names are permitted. All other SourceReasoner operations remain unchanged.
func (r *SourceReasoner) ObserveNativeReasoningTask(ctx context.Context,action int)(NativeReasoningConsequence,error){
 var empty NativeReasoningConsequence
 if r==nil||action<0||action>=len(nativeReasoningTasks){return empty,errors.New("invalid exact source task action")}
 if r.State.Observations>=3{return empty,errors.New("source-task observation budget exhausted")}
 if err:=r.validate();err!=nil{return empty,err}
 task:=nativeReasoningTasks[action]
 symbol,err:=r.symbol(task.Name);if err!=nil{return empty,err}
 if symbol.Kind!="function"||symbol.File!=task.Source{return empty,errors.New("native task symbol identity drift")}
 file:=filepath.Join(r.Root,task.Source)
 b,err:=os.ReadFile(file);if err!=nil{return empty,err}
 // Check that checked-out source is byte-identical with immutable Git HEAD.
 git:=exec.CommandContext(ctx,"git","show","HEAD:"+task.Source);git.Dir=r.Root
 committed,err:=git.Output();if err!=nil{return empty,fmt.Errorf("native source proof unavailable: %w",err)}
 if string(b)!=string(committed){return empty,errors.New("source not equal to committed native task")}
 commitCmd:=exec.CommandContext(ctx,"git","rev-parse","HEAD");commitCmd.Dir=r.Root
 sha,err:=commitCmd.Output();if err!=nil{return empty,err}
 commit:=strings.TrimSpace(string(sha))
 if len(commit)!=40{return empty,errors.New("invalid exact native source commit")}
 deadline,cancel:=context.WithTimeout(ctx,50*time.Second);defer cancel()
 cmd:=exec.CommandContext(deadline,"go","test","-json",task.Package,"-run","^"+task.Name+"$","-count=1")
 cmd.Dir=r.Root
 cmd.Env=append(os.Environ(),"GOFLAGS=","GOWORK=off","GOPROXY=off","GOSUMDB=off","GOTOOLCHAIN=local")
 output,runErr:=cmd.CombinedOutput()
 if deadline.Err()!=nil{return empty,deadline.Err()}
 // A Go test failure is evidence ONLY when the exact test ran to an
 // assertion result. Build failures and environment errors never train.
 observed,err:=nativeGoTestVerdict(output,task,runErr)
 if err!=nil{return empty,fmt.Errorf("unqualified native action consequence: %w",err)}
 dest,err:=r.symbol("ObserveExistingCall");if err!=nil{return empty,err}
 if err:=r.State.Observe(vectorFromSource(symbol),vectorFromSource(dest),action,observed);err!=nil{return empty,err}
 digest:=sha256.Sum256(b)
 return NativeReasoningConsequence{TaskName:task.Name,SourceCommit:commit,
  SourceFileSHA256:hex.EncodeToString(digest[:]),IndexSHA:r.IndexSHA,
  Action:uint8(action),ObservedPass:uint8(observed),
  EvidenceBasis:"exact-ICE-index+git-blob+completed-native-go-test-json-assertion"},nil
}
func (r *SourceReasoner) PredictNativeReasoningTask(action int)(float64,error){
 if r==nil||action<0||action>=len(nativeReasoningTasks){return 0,errors.New("invalid native reasoning action")}
 if err:=r.validate();err!=nil{return 0,err}
 source,err:=r.symbol(nativeReasoningTasks[action].Name);if err!=nil{return 0,err}
 dest,err:=r.symbol("ObserveExistingCall");if err!=nil{return 0,err}
 return r.State.Predict(vectorFromSource(source),vectorFromSource(dest),action)
}
