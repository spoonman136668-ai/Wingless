package cognition

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "testing"
    "time"
    "unsafe"

    tb "github.com/spoonman136668-ai/Wingless/toolboundary"
)

type p6hCounts struct {
    failure [2]uint16
    success [2]uint16
}
func (c *p6hCounts) index(pre [FeatureWidth]float32)int {
    if pre[0]<0{return 0}
    return 1
}
func (c *p6hCounts) observe(pre [FeatureWidth]float32,success bool) {
    i:=c.index(pre)
    if success {c.success[i]++}else{c.failure[i]++}
}
func (c *p6hCounts) predict(pre [FeatureWidth]float32)float64 {
    i:=c.index(pre)
    return float64(c.success[i]+1)/float64(c.failure[i]+c.success[i]+2)
}

func p6hNativeEvidence(t *testing.T, path,contents string, present bool)([FeatureWidth]float32,[FeatureWidth]float32,bool,time.Duration) {
    t.Helper()
    root:=t.TempDir()
    if present {
        if err:=os.WriteFile(filepath.Join(root,path),[]byte(contents),0600);err!=nil {t.Fatal(err)}
    }
    p:=tb.Policy{Workspace:root,AllowedTools:[]string{"read_file"},AllowedPaths:[]string{path},Timeout:time.Second,MaxOutputBytes:128}
    a:=tb.Proposal{Tool:"read_file",Path:path}
    observedPre,err:=tb.ObservePreaction(context.Background(),p,a)
    if err!=nil {t.Fatal("unauthenticated preaction native environment:",err)}
    before,err:=NativePreactionFeatures(observedPre)
    if err!=nil {t.Fatal(err)}
    began:=time.Now()
    result,executedErr:=tb.Execute(context.Background(),p,a)
    wall:=time.Since(began)
    success:=executedErr==nil && result.ExitCode==0
    if success!=present {t.Fatal("original native observed outcome violates prepared source state",executedErr)}
    after:=NativeExecutedConsequenceFeatures(result,executedErr)
    return before,after,success,wall
}

func p6hNativeTrain(t *testing.T)(LatentState,p6hCounts,time.Duration,time.Duration) {
    t.Helper()
    var state LatentState
    var conventional p6hCounts
    var nativeWall,modelWall time.Duration
    events:=[]struct{path,contents string;present bool}{
        {"first.txt","alpha",true},
        {"absent-a.txt","",false},
        {"second.txt","bravo",true},
        {"absent-b.txt","",false},
    }
    for _,e:=range events{
        before,after,success,nativeDuration:=p6hNativeEvidence(t,e.path,e.contents,e.present)
        nativeWall+=nativeDuration
        began:=time.Now()
        if err:=state.ObserveNativeTransition(before,after,0);err!=nil {t.Fatal("bounded native transition learning failed:",err)}
        modelWall+=time.Since(began)
        conventional.observe(before,success)
    }
    if state.Observations!=uint32(len(events)){t.Fatal("unowned native learning events")}
    if unsafe.Sizeof(state)!=6472{t.Fatalf("source bound fixed state grew: %d",unsafe.Sizeof(state))}
    return state,conventional,nativeWall,modelWall
}

func TestNCI04P4P6HNativeTransitionsAndTrainedControls(t *testing.T) {
    p6fPinOriginalSource(t)
    var empty LatentState
    unknown:= [FeatureWidth]float32{1,5.0/64,-1,0}
    if _,known,err:=empty.PredictNativeTransition(unknown,0);err!=nil||known {
        t.Fatal("untrained learner made confident prediction",known,err)
    }
    var memBefore,memAfter runtime.MemStats
    runtime.ReadMemStats(&memBefore)
    state,conventional,nativeWall,modelWall:=p6hNativeTrain(t)
    var candidateLoss,conditionedLoss,sourceInspectorLoss float64
    var predictionWall,controlWall time.Duration
    for _,e:=range []struct{path,contents string;present bool}{
        {"novel-third.txt","fresh",true},
        {"novel-absent.txt","",false},
    } {
        // Fresh, developer-selected native paths are development checks, not
        // an independently held unfamiliar task family or scientific first use.
        root:=t.TempDir()
        if e.present {
            if err:=os.WriteFile(filepath.Join(root,e.path),[]byte(e.contents),0600);err!=nil {t.Fatal(err)}
        }
        p:=tb.Policy{Workspace:root,AllowedTools:[]string{"read_file"},AllowedPaths:[]string{e.path},Timeout:time.Second,MaxOutputBytes:128}
        a:=tb.Proposal{Tool:"read_file",Path:e.path}
        pre,err:=tb.ObservePreaction(context.Background(),p,a)
        if err!=nil {t.Fatal(err)}
        before,err:=NativePreactionFeatures(pre)
        if err!=nil {t.Fatal(err)}
        began:=time.Now()
        predicted,known,err:=state.PredictNativeTransition(before,0)
        predictionWall+=time.Since(began)
        if err!=nil || !known {t.Fatal("learner failed same-tool new-identity development readout",known,err)}
        began=time.Now()
        countProbability:=conventional.predict(before)
        controlWall+=time.Since(began)
        // Independent same-information native source-inspection control:
        // the original tool is read_file; the authorized preaction inspection
        // already establishes present/missing for this stable workspace.
        sourcePrediction:=0.0
        if pre.Resource=="regular"{sourcePrediction=1}
        result,originalErr:=tb.Execute(context.Background(),p,a)
        observedSuccess:=originalErr==nil && result.ExitCode==0
        actual:=0.0
        if observedSuccess {actual=1}
        if observedSuccess!=e.present {t.Fatal("native heldout operation mismatch",originalErr)}
        candidateProbability:=0.0
        if predicted[0]>0{candidateProbability=1}
        candidateLoss+=(candidateProbability-actual)*(candidateProbability-actual)
        conditionedLoss+=(countProbability-actual)*(countProbability-actual)
        sourceInspectorLoss+=(sourcePrediction-actual)*(sourcePrediction-actual)
    }
    runtime.ReadMemStats(&memAfter)
    snapshot,err:=state.Snapshot()
    if err!=nil {t.Fatal(err)}
    receipt:=map[string]any{
        "gate":"NCI-04P6H","source_native_git_blob":p6fToolBlob,
        "authentic_training_events":int(state.Observations),
        "developer_selected_different_file_identities":2,
        "unseen_native_source_families":0,
        "native_tool_action":"read_file",
        "learner_prior_abstains":true,
        "learner_fixed_state_bytes":unsafe.Sizeof(state),
        "learner_snapshot_bytes":len(snapshot),
        "parent_runtime_heap_before_bytes":memBefore.Alloc,
        "parent_runtime_heap_after_bytes":memAfter.Alloc,
        "native_execution_wall_ns":nativeWall.Nanoseconds(),
        "learner_update_wall_ns":modelWall.Nanoseconds(),
        "learner_predict_wall_ns":predictionWall.Nanoseconds(),
        "trained_count_predict_wall_ns":controlWall.Nanoseconds(),
        "candidate_dev_brier_sum":candidateLoss,
        "trained_conditional_count_dev_brier_sum":conditionedLoss,
        "conventional_source_inspector_dev_brier_sum":sourceInspectorLoss,
        "source_inspection_superiority":candidateLoss<sourceInspectorLoss,
        "all_in_parent_child_cpu_rss_io":"NOT_QUALIFIED",
        "scientific_first_use":false,
        "claim":"FIXED_STATE_AUTHENTIC_SAME_TOOL_LEARNING_ONLY__UNFAMILIAR_FAMILY_TRANSFER_NOT_PROVEN",
    }
    b,err:=json.Marshal(receipt);if err!=nil {t.Fatal(err)}
    t.Log("P6H_NATIVE_LEARNING_RECEIPT="+string(b))
    if candidateLoss<0 || sourceInspectorLoss!=0 || candidateLoss<sourceInspectorLoss {
        t.Fatal("unearned native source-inspection advantage")
    }
}

type p6hSealed struct{
    Schema string `json:"schema"`
    OriginalGitBlob string `json:"original_git_blob"`
    State json.RawMessage `json:"state"`
}
func TestNCI04P4P6HColdProcessRestart(t *testing.T){
    p6fPinOriginalSource(t)
    if os.Getenv("P6H_COLD_CHILD")=="1" {
        b,err:=os.ReadFile(os.Getenv("P6H_SNAPSHOT"))
        if err!=nil {t.Fatal(err)}
        var seal p6hSealed
        if err:=json.Unmarshal(b,&seal);err!=nil {t.Fatal(err)}
        if seal.Schema!="wingless.nci04p6h.native.v1" || seal.OriginalGitBlob!=p6fToolBlob {
            t.Fatal("cross-process source/model seal mismatch")
        }
        var state LatentState
        if err:=state.Restore(seal.State);err!=nil {t.Fatal(err)}
        for _,e:=range []struct{before [FeatureWidth]float32;positive bool}{
            {[FeatureWidth]float32{1,5.0/64,-1,0},true},
            {[FeatureWidth]float32{-1,0,-1,0},false},
        }{
            next,known,err:=state.PredictNativeTransition(e.before,0)
            if err!=nil||!known||(next[0]>0)!=e.positive {t.Fatal("source-bound cold restart changed learned transition",err)}
        }
        t.Log("P6H_FRESH_PROCESS_RESTART_PASS")
        return
    }
    state,_,_,_:=p6hNativeTrain(t)
    snapshot,err:=state.Snapshot();if err!=nil {t.Fatal(err)}
    sealed,err:=json.Marshal(p6hSealed{Schema:"wingless.nci04p6h.native.v1",OriginalGitBlob:p6fToolBlob,State:snapshot})
    if err!=nil {t.Fatal(err)}
    path:=filepath.Join(t.TempDir(),"p6h-sealed-native.json")
    if err:=os.WriteFile(path,sealed,0600);err!=nil {t.Fatal(err)}
    cmd:=exec.Command(os.Args[0],"-test.run=^TestNCI04P4P6HColdProcessRestart$","-test.v")
    cmd.Env=append(os.Environ(),"P6H_COLD_CHILD=1","P6H_SNAPSHOT="+path)
    out,err:=cmd.CombinedOutput()
    if err!=nil || !bytes.Contains(out,[]byte("P6H_FRESH_PROCESS_RESTART_PASS")){
        t.Fatalf("true cold process restoration failed: %v %s",err,out)
    }
    unchanged,err:=os.ReadFile(path)
    if err!=nil || !bytes.Equal(unchanged,sealed){t.Fatal("cold restart modified source-pinned model checkpoint",err)}
    t.Log(fmt.Sprintf("P6H_COLD_PROCESS_RESTART_PASS snapshot_bytes=%d",len(sealed)))
}

func TestNCI04P4P6HMalformedTransitionFailsClosed(t *testing.T){
    p6fPinOriginalSource(t)
    var s LatentState
    before:=[FeatureWidth]float32{1,5.0/64,-1,0}
    after:=[FeatureWidth]float32{1,5.0/64,1,0}
    baseline,_:=s.Snapshot()
    if err:=s.ObserveNativeTransition(after,before,0);err==nil {t.Fatal("reversed modality accepted")}
    if err:=s.ObserveNativeTransition(before,after,MaxActions);err==nil{t.Fatal("invalid action accepted")}
    if _,_,err:=s.PredictNativeTransition(after,0);err==nil {t.Fatal("post-outcome feature permitted in pre-action prediction")}
    current,_:=s.Snapshot()
    if !bytes.Equal(current,baseline){t.Fatal("invalid native event mutated fixed state")}
}
