package cognition

// P6L: genuine native-effect conflict as an identifiability safety gate.
// Original LatentState.Observe mutates its own native state; original
// LatentState.Predict leaves that state unchanged. Both are called for real.
// The SOURCE-only preaction sensor is the same, so an action channel that
// omits the actual operation must abstain after observing both effects.
// This development experiment is not evidence of transferable cognition.
import (
    "bytes"
    "context"
    "encoding/json"
    "os"
    "os/exec"
    "path/filepath"
    "testing"
    "time"
    "unsafe"

    tb "github.com/spoonman136668-ai/Wingless/toolboundary"
)

func p6lOriginalCognitionPreaction(t *testing.T) [FeatureWidth]float32 {
    t.Helper()
    root,err:=filepath.Abs("..")
    if err!=nil {t.Fatal(err)}
    // Pin source without consuming semantic-atlas data or future outcomes.
    nativeSource,err:=os.ReadFile(filepath.Join(root,"cognition","latent_bind.go"))
    if err!=nil{t.Fatal(err)}
    if len(nativeSource)==0||len(nativeSource)>16384 {t.Fatal("unbounded original cognition source")}
    cmd:=exec.Command("git","show","HEAD:cognition/latent_bind.go")
    cmd.Dir=root
    committed,err:=cmd.Output()
    if err!=nil || !bytes.Equal(committed,nativeSource){
        t.Fatal("original native cognition source differs from committed source",err)
    }
    p:=tb.Policy{Workspace:root,AllowedTools:[]string{"read_file"},AllowedPaths:[]string{"cognition/latent_bind.go"},Timeout:time.Second,MaxOutputBytes:16384}
    observed,err:=tb.ObservePreaction(context.Background(),p,tb.Proposal{Tool:"read_file",Path:"cognition/latent_bind.go"})
    if err!=nil||observed.Resource!="regular"{t.Fatal("native preaction source not authenticated",err)}
    source,err:=NativePreactionFeatures(observed)
    if err!=nil{t.Fatal(err)}
    return source
}

func p6lExecuteOriginalMemoryAction(t *testing.T,change bool)([FeatureWidth]float32,bool) {
    t.Helper()
    var native LatentState
    before,err:=native.Snapshot()
    if err!=nil{t.Fatal(err)}
    src:=[FeatureWidth]float32{-0.6,-0.3,0,0}
    dst:=[FeatureWidth]float32{0.6,0.3,0,0}
    if change {
        err=native.Observe(src,dst,0,1)
    } else {
        _,err=native.Predict(src,dst,0)
    }
    if err!=nil{t.Fatal("original native cognition operation failed",err)}
    after,err:=native.Snapshot()
    if err!=nil{t.Fatal(err)}
    changed:=!bytes.Equal(before,after)
    if changed!=change {t.Fatal("original native state-change witness inconsistent")}
    // The effect is sourced exclusively from actual before/after bytes.
    // This is recorded AFTER executing native cognition, not a future label.
    observedEffect:=[FeatureWidth]float32{-1,0,1,0}
    if changed {observedEffect[0]=1}
    return observedEffect,changed
}

func TestNCI04P4P6LAuthenticConflictingNativeEffectsRequireAbstention(t *testing.T){
    p6fPinOriginalSource(t)
    pre:=p6lOriginalCognitionPreaction(t)
    var learned LatentState
    if _,known,err:=learned.PredictNativeTransition(pre,1);err!=nil||known{
        t.Fatal("untrained native effect was overconfident",err)
    }
    // Three genuinely executed state changes; then one real read-only call.
    // Both are original native project methods, on fresh identical states.
    var changes [4]bool
    for i,mutating:=range []bool{true,true,true,false}{
        after,actualChange:=p6lExecuteOriginalMemoryAction(t,mutating)
        changes[i]=actualChange
        if err:=learned.ObserveNativeTransition(pre,after,1);err!=nil{
            t.Fatal("native effect learning rejected admitted observation",err)
        }
    }
    if changes!=[4]bool{true,true,true,false}{t.Fatal("native effect sampling failed")}
    // The previously implemented P6H majority readout would claim effect
    // known (3 versus 1), but the actual native operations are still
    // observationally indistinguishable to this preaction feature pair.
    prediction,known,err:=learned.PredictNativeTransition(pre,1)
    if err!=nil||known||prediction!=([FeatureWidth]float32{}){
        t.Fatal("false confident native action-effect prediction without action-affordance information",err)
    }
    snap,err:=learned.Snapshot()
    if err!=nil{t.Fatal(err)}
    var reopened LatentState
    if err:=reopened.Restore(snap);err!=nil{t.Fatal(err)}
    again,known,err:=reopened.PredictNativeTransition(pre,1)
    if err!=nil||known||again!=prediction{t.Fatal("native conflict-abstention failed persisted-state restore",err)}
    if unsafe.Sizeof(learned)!=6472{t.Fatal("native learned capacity expanded")}
    report:=map[string]any{
        "gate":"NCI-04P6L",
        "true_original_native_actions":4,
        "native_mutating_observations":3,
        "native_readonly_observations":1,
        "source_authentication":"EXACT_GIT_BLOB_AND_AUTHORIZED_ROOTED_READ",
        "original_execution_families":[]string{"LatentState.Observe","LatentState.Predict"},
        "native_state_delta_measured_after_action":true,
        "preaction_perception_identical":true,
        "same_action_channel_supplied":true,
        "learner_retains_conflicting_effects":true,
        "learner_predictions":"ABSTAIN_BEFORE_AND_AFTER_CONFLICT",
        "fixed_state_bytes":unsafe.Sizeof(learned),
        "learned_snapshot_bytes":len(snap),
        "cold_os_process_restart":"NOT_TESTED; in-process independent-state restore only",
        "same_info_conventional_model":"ABSTAIN",
        "learned_transfer":"NOT_PROVEN",
        "all_in_resource_parity":"NOT_QUALIFIED",
        "scientific_first_use":false,
        "decision":"EXPERIENTIAL_CONFLICT_OBSERVED_BUT_ACTION_AFFORDANCE_NOT_IDENTIFIABLE",
    }
    b,err:=json.Marshal(report);if err!=nil{t.Fatal(err)}
    t.Log("P6L_NATIVE_EXPERIENTIAL_CONFLICT_RECEIPT="+string(b))
}

func TestNCI04P4P6LUnambiguousLearnedExperienceStillReadout(t *testing.T){
    p6fPinOriginalSource(t)
    pre:=p6lOriginalCognitionPreaction(t)
    var learned LatentState
    for i:=0;i<2;i++ {
        after,changed:=p6lExecuteOriginalMemoryAction(t,true)
        if !changed{t.Fatal("actual source consequence vanished")}
        if err:=learned.ObserveNativeTransition(pre,after,1);err!=nil{t.Fatal(err)}
    }
    _,known,err:=learned.PredictNativeTransition(pre,1)
    if err!=nil||!known {t.Fatal("false abstention on unambiguous genuine experiential action",err)}
}
