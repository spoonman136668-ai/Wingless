package cognition

// NCI-04G7: an opt-in, actual project-native source-reasoning task outcome.
// This is not the prior generated single-file probe. The consequence is the
// exit status of an executed, specifically allowlisted Wingless test sourced
// from the exact validated ICE index and git commit. Failed executions are
// errors, NOT an invented negative supervision label.
import (
 "context"
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
 cmd:=exec.CommandContext(deadline,"go","test",task.Package,"-run","^"+task.Name+"$","-count=1")
 cmd.Dir=r.Root
 cmd.Env=append(os.Environ(),"GOFLAGS=","GOWORK=off","GOPROXY=off","GOSUMDB=off","GOTOOLCHAIN=local")
 output,err:=cmd.CombinedOutput()
 if err!=nil{return empty,fmt.Errorf("native reasoning task failed or unavailable: %w",err)}
 if deadline.Err()!=nil{return empty,deadline.Err()}
 if !strings.Contains(string(output),"ok")||!strings.Contains(string(output),"Wingless/"){
   return empty,errors.New("native test lacked source-bound successful Go test receipt")
 }
 dest,err:=r.symbol("ObserveExistingCall");if err!=nil{return empty,err}
 if err:=r.State.Observe(vectorFromSource(symbol),vectorFromSource(dest),action,1);err!=nil{return empty,err}
 digest:=sha256.Sum256(b)
 return NativeReasoningConsequence{TaskName:task.Name,SourceCommit:commit,
  SourceFileSHA256:hex.EncodeToString(digest[:]),IndexSHA:r.IndexSHA,
  Action:uint8(action),ObservedPass:1,
  EvidenceBasis:"exact-ICE-index+git-blob+executed-native-Wingless-go-test"},nil
}
func (r *SourceReasoner) PredictNativeReasoningTask(action int)(float64,error){
 if r==nil||action<0||action>=len(nativeReasoningTasks){return 0,errors.New("invalid native reasoning action")}
 if err:=r.validate();err!=nil{return 0,err}
 source,err:=r.symbol(nativeReasoningTasks[action].Name);if err!=nil{return 0,err}
 dest,err:=r.symbol("ObserveExistingCall");if err!=nil{return 0,err}
 return r.State.Predict(vectorFromSource(source),vectorFromSource(dest),action)
}
