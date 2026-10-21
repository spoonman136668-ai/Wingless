package cognition

// P4 is an opt-in integration of authentic, project-owned native source
// intervention consequences with the existing fixed-capacity learned state.
// It is an engineering feasibility seam, NOT unfamiliar-task transfer or an
// independently blind scientific test.
import (
 "bytes"
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

type NativeInterventionLearningReceipt struct {
 SourceCommit string `json:"source_commit"`
 ICEManifestSHA string `json:"ice_manifest_sha"`
 IntervenedFile string `json:"intervened_file"`
 OriginalSHA256 string `json:"original_sha256"`
 NativeOutcomes [3]uint8 `json:"native_outcomes"`
 NativeWallNS [3]int64 `json:"native_wall_ns"`
 BaselinePrior float64 `json:"baseline_prior"`
 IntervenedPrior float64 `json:"intervened_prior"`
 BaselinePosterior float64 `json:"baseline_posterior"`
 IntervenedPosterior float64 `json:"intervened_posterior"`
 LearnedObservations uint32 `json:"learned_observations"`
 EvidenceBasis string `json:"evidence_basis"`
 Admission string `json:"admission"`
}

const iceNativeQueryGuard = "if e := Validate(root); e != nil {"
const iceNativeDisabledGuard = "if e := Validate(root); false {"

// executeIsolatedICETest is an actual project-native assertion measurement.
// Compilation, timeout and operational errors are never converted into labels.
func executeIsolatedICETest(ctx context.Context,checkout string)(uint8,int64,[]byte,error){
 task:=nativeReasoningTest{Name:"TestSourceWins",Package:"./ice",Source:"ice/index_test.go"}
 deadline,cancel:=context.WithTimeout(ctx,35*time.Second);defer cancel()
 cmd:=exec.CommandContext(deadline,"go","test","-json","./ice","-run","^TestSourceWins$","-count=1")
 cmd.Dir=checkout
 cmd.Env=append(os.Environ(),"GOFLAGS=","GOWORK=off","GOPROXY=off","GOSUMDB=off","GOTOOLCHAIN=local")
 start:=time.Now()
 output,err:=cmd.CombinedOutput()
 elapsed:=time.Since(start).Nanoseconds()
 if deadline.Err()!=nil{return 0,elapsed,nil,deadline.Err()}
 observed,check:=nativeGoTestVerdict(output,task,err)
 if check!=nil{return 0,elapsed,nil,fmt.Errorf("nonadmissible native test result: %w",check)}
 if len(output)>1<<20{return 0,elapsed,nil,errors.New("oversized native test result")}
 return uint8(observed),elapsed,output,nil
}

// ObserveNativeICEIntervention measures the original test's own outcomes
// before, during, and after ONE tightly bounded change to Query's stale-source
// guard inside a disposable detached Git checkout. It never learns from a
// caller-supplied result. The learner commits only after full restoration.
// The intervention policy is hard-pinned; no arbitrary code execution input.
func (r *SourceReasoner) ObserveNativeICEIntervention(ctx context.Context)(NativeInterventionLearningReceipt,error){
 var empty NativeInterventionLearningReceipt
 if r==nil{return empty,errors.New("nil native reasoning owner")}
 if err:=r.validate();err!=nil{return empty,err}
 if r.State.Observations!=0{return empty,errors.New("requires untrained intervention-specific native state")}
 source,err:=r.symbol("TestSourceWins");if err!=nil{return empty,err}
 if source.File!="ice/index_test.go"||source.Kind!="function"{return empty,errors.New("native test source identity drift")}
 target,err:=r.symbol("ObserveExistingCall");if err!=nil{return empty,err}
 src,dst:=vectorFromSource(source),vectorFromSource(target)
 baselinePrior,err:=r.State.Predict(src,dst,0);if err!=nil{return empty,err}
 actionPrior,err:=r.State.Predict(src,dst,1);if err!=nil{return empty,err}
 head:=exec.CommandContext(ctx,"git","rev-parse","HEAD");head.Dir=r.Root
 output,err:=head.Output();if err!=nil{return empty,err}
 sha:=strings.TrimSpace(string(output))
 if len(sha)!=40{return empty,errors.New("untrusted native source HEAD")}
 originalPath:=filepath.Join(r.Root,"ice","index.go")
 working,err:=os.ReadFile(originalPath);if err!=nil{return empty,err}
 committed:=exec.CommandContext(ctx,"git","show","HEAD:ice/index.go");committed.Dir=r.Root
 trusted,err:=committed.Output();if err!=nil||!bytes.Equal(working,trusted){return empty,errors.New("native source differs from pinned Git blob")}
 if bytes.Count(working,[]byte(iceNativeQueryGuard))!=1{return empty,errors.New("native ICE Query guard no longer uniquely identified")}
 scratch,err:=os.MkdirTemp("","wingless-p4-native-");if err!=nil{return empty,err}
 defer os.RemoveAll(scratch)
 checkout:=filepath.Join(scratch,"source")
 clone:=exec.CommandContext(ctx,"git","clone","--quiet","--no-hardlinks","--no-checkout",r.Root,checkout)
 if out,e:=clone.CombinedOutput();e!=nil{return empty,fmt.Errorf("native scratch clone unavailable: %w %s",e,out)}
 detach:=exec.CommandContext(ctx,"git","-C",checkout,"checkout","--quiet","--detach",sha)
 if out,e:=detach.CombinedOutput();e!=nil{return empty,fmt.Errorf("exact native checkout unavailable: %w %s",e,out)}
 local:=filepath.Join(checkout,"ice","index.go")
 before,err:=os.ReadFile(local);if err!=nil||!bytes.Equal(before,trusted){return empty,errors.New("native clone source-blob identity mismatch")}
 var labels [3]uint8
 var walls [3]int64
 labels[0],walls[0],_,err=executeIsolatedICETest(ctx,checkout);if err!=nil{return empty,err}
 if labels[0]!=1{return empty,errors.New("native control assertion not passing")}
 changed:=bytes.Replace(before,[]byte(iceNativeQueryGuard),[]byte(iceNativeDisabledGuard),1)
 if bytes.Equal(changed,before){return empty,errors.New("native source intervention missing")}
 if err:=os.WriteFile(local,changed,0600);err!=nil{return empty,err}
 diff:=exec.CommandContext(ctx,"git","-C",checkout,"diff","--name-only")
 paths,err:=diff.Output()
 if err!=nil||strings.TrimSpace(string(paths))!="ice/index.go"{return empty,errors.New("native intervention escaped authorized source file")}
 var failing []byte
 labels[1],walls[1],failing,err=executeIsolatedICETest(ctx,checkout);if err!=nil{return empty,err}
 if labels[1]!=0||!bytes.Contains(failing,[]byte("stale index permitted")){
  return empty,errors.New("native test failure not caused by expected original ICE guard")
 }
 if err:=os.WriteFile(local,before,0600);err!=nil{return empty,err}
 labels[2],walls[2],_,err=executeIsolatedICETest(ctx,checkout);if err!=nil{return empty,err}
 if labels[2]!=1{return empty,errors.New("native source reversal did not restore original test")}
 restored,err:=os.ReadFile(local);if err!=nil||!bytes.Equal(restored,before){return empty,errors.New("source reverse intervention bytes changed")}
 clean:=exec.CommandContext(ctx,"git","-C",checkout,"status","--porcelain")
 dirty,err:=clean.Output();if err!=nil||len(bytes.TrimSpace(dirty))!=0{return empty,errors.New("isolated native source not restored cleanly")}
 if err:=r.validate();err!=nil{return empty,fmt.Errorf("live ICE source identity changed: %w",err)}
 // Only after complete independent process evidence and cleanup validation
 // do we update native cognitive state. Learning is transactionally bounded.
 next:=r.State
 for i,action:=range []int{0,1,0}{
  if err:=next.Observe(src,dst,action,int(labels[i]));err!=nil{return empty,err}
 }
 p0,err:=next.Predict(src,dst,0);if err!=nil{return empty,err}
 p1,err:=next.Predict(src,dst,1);if err!=nil{return empty,err}
 r.State=next
 hash:=sha256.Sum256(before)
 return NativeInterventionLearningReceipt{
  SourceCommit:sha,ICEManifestSHA:r.IndexSHA,IntervenedFile:"ice/index.go",
  OriginalSHA256:hex.EncodeToString(hash[:]),NativeOutcomes:labels,NativeWallNS:walls,
  BaselinePrior:baselinePrior,IntervenedPrior:actionPrior,
  BaselinePosterior:p0,IntervenedPosterior:p1,
  LearnedObservations:next.Observations,
  EvidenceBasis:"exact-ICE-index + Git source blob + isolated native Go assertion + clean reversal",
  Admission:"NATIVE_LEARNING_INTEGRATION_ONLY__TRANSFER_AND_RESOURCE_NONREGRESSION_NOT_QUALIFIED",
 },nil
}
