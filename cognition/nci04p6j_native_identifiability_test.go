package cognition

// P6J native identifiability audit. This deliberately changes no learner,
// action executor, protected science or capacity. Two different existing
// project-native operations receive equivalent public preaction observations.
// If their consequences differ, a source-blind predictor without any additional
// observed action-affordance information cannot identify both effects.
import (
    "bytes"
    "context"
    "encoding/json"
    "os/exec"
    "path/filepath"
    "reflect"
    "runtime"
    "strings"
    "testing"
    "time"
    "unsafe"

    "github.com/spoonman136668-ai/Wingless/ice"
    tb "github.com/spoonman136668-ai/Wingless/toolboundary"
)

func p6jExactDisposableClone(t *testing.T,source,head string)string{
    t.Helper()
    checkout:=filepath.Join(t.TempDir(),"original-native")
    command:=exec.Command("git","clone","--quiet","--no-hardlinks","--no-checkout",source,checkout)
    if b,err:=command.CombinedOutput();err!=nil {t.Fatalf("exact-source native clone: %v %s",err,b)}
    p6iGit(t,checkout,"checkout","--quiet","--detach",head)
    if p6iGit(t,checkout,"rev-parse","HEAD")!=head || p6iGit(t,checkout,"status","--porcelain")!=""{
        t.Fatal("native action clone was not clean and source-identical")
    }
    return checkout
}

func p6jOriginalPreaction(t *testing.T,clone string)(tb.PreactionObservation,[FeatureWidth]float32){
    t.Helper()
    policy:=tb.Policy{
        Workspace:clone,AllowedTools:[]string{"read_file"},
        AllowedPaths:[]string{"ice/index.go"},Timeout:time.Second,MaxOutputBytes:65536,
    }
    obs,err:=tb.ObservePreaction(context.Background(),policy,tb.Proposal{Tool:"read_file",Path:"ice/index.go"})
    if err!=nil || obs.Resource!="regular" || obs.SHA256==""{t.Fatal("missing authorized native source preaction information",err)}
    vector,err:=NativePreactionFeatures(obs)
    if err!=nil {t.Fatal(err)}
    return obs,vector
}

func TestNCI04P4P6JNativeActionEffectIdentifiability(t *testing.T){
    p6fPinOriginalSource(t)
    state,_,_,_:=p6hNativeTrain(t)
    frozen,err:=state.Snapshot()
    if err!=nil {t.Fatal(err)}
    if unsafe.Sizeof(state)!=6472{t.Fatal("fixed learner state drift")}
    root,err:=filepath.Abs("..")
    if err!=nil {t.Fatal(err)}
    head:=p6iGit(t,root,"rev-parse","HEAD")
    if len(head)!=40 {t.Fatal("source HEAD not pinned")}
    // Build and Save are genuine original ICE package operations with
    // different side effects, performed on two exact-source clean clones.
    buildClone:=p6jExactDisposableClone(t,root,head)
    saveClone:=p6jExactDisposableClone(t,root,head)
    sourceBuild,featuresBuild:=p6jOriginalPreaction(t,buildClone)
    sourceSave,featuresSave:=p6jOriginalPreaction(t,saveClone)
    if sourceBuild!=sourceSave || !reflect.DeepEqual(featuresBuild,featuresSave){
        t.Fatal("native preaction observations were not identical; uninformative pair")
    }
    // A single previously UNTRAINED action channel is used for both. No
    // operation symbol, hardcoded read/write kind, source name, or hidden
    // post-action outcome is delivered to the prediction function.
    const unlearnedAction=1
    buildPrior,buildKnown,err:=state.PredictNativeTransition(featuresBuild,unlearnedAction)
    if err!=nil {t.Fatal(err)}
    savePrior,saveKnown,err:=state.PredictNativeTransition(featuresSave,unlearnedAction)
    if err!=nil {t.Fatal(err)}
    if buildKnown||saveKnown||buildPrior!=savePrior {
        t.Fatal("native zero-information action channel produced unwarranted confident distinct effects")
    }
    if err:=ice.Validate(buildClone);err==nil {t.Fatal("native Build control unexpectedly already has a fresh ICE index")}
    if err:=ice.Validate(saveClone);err==nil {t.Fatal("native Save control unexpectedly already has a fresh ICE index")}

    // Original read-only build: genuine source inspection, no write.
    t0:=time.Now()
    built,err:=ice.Build(buildClone)
    buildWall:=time.Since(t0)
    if err!=nil || len(built.Files)==0 {t.Fatal("native read-only ice.Build failed",err)}
    if p6iGit(t,buildClone,"status","--porcelain")!="" {
        t.Fatal("original read-only native Build unexpectedly mutated source/index")
    }
    if err:=ice.Validate(buildClone);err==nil {t.Fatal("read-only native Build was mislabeled as state-changing")}

    // Original native Save: actually changes persistent indexed state.
    var mem0,mem1 runtime.MemStats
    runtime.ReadMemStats(&mem0)
    t0=time.Now()
    if err:=ice.Save(saveClone);err!=nil {t.Fatal("original native state-changing ice.Save failed",err)}
    saveWall:=time.Since(t0)
    runtime.ReadMemStats(&mem1)
    if err:=ice.Validate(saveClone);err!=nil {t.Fatal("native Save failed to repair original source index",err)}
    status:=p6iGit(t,saveClone,"status","--porcelain")
    if status=="" {t.Fatal("original native Save did not write any canonical artifacts")}
    for _,line:=range strings.Split(status,"\n"){
        fields:=strings.Fields(line)
        if len(fields)!=2 || !strings.HasPrefix(fields[1],".ice/"){
            t.Fatalf("native Save escaped expected artifacts: %q",line)
        }
    }
    if p6iGit(t,buildClone,"rev-parse","HEAD")!=head ||
       p6iGit(t,saveClone,"rev-parse","HEAD")!=head {
        t.Fatal("original source commit changed during native actions")
    }
    if after,err:=state.Snapshot();err!=nil||!bytes.Equal(after,frozen){
        t.Fatal("unseen native effects contaminated learned state",err)
    }
    // A conventional agent given only these identical preaction percepts
    // ALSO cannot distinguish them. A fully informed source inspector with
    // actual action/function identity could, but that is MORE information;
    // do not claim such a comparison is matched or award candidate credit.
    report:=map[string]any{
        "gate":"NCI-04P6J",
        "exact_source_commit":head,
        "learner_previously_trained_native_family":"toolboundary.Execute/read_file",
        "original_native_first_operation":"ice.Build",
        "original_native_second_operation":"ice.Save",
        "preaction_feature_identity":"IDENTICAL",
        "action_information_supplied_to_learner":"SAME_UNTRAINED_CHANNEL_1",
        "predictions_before_execution":"BOTH_ABSTAIN",
        "original_readonly_operation_changed_workspace":false,
        "original_native_save_changed_workspace":true,
        "original_native_save_repaired_index":true,
        "build_wall_ns":buildWall.Nanoseconds(),
        "save_wall_ns":saveWall.Nanoseconds(),
        "snapshot_bytes":len(frozen),
        "fixed_native_state_bytes":unsafe.Sizeof(state),
        "parent_heap_before_bytes":mem0.Alloc,
        "parent_heap_after_bytes":mem1.Alloc,
        "fully_informed_static_source_control":"NOT_FAIR_SAME_INFO_BECAUSE_ACTION_IDENTITY_MISSING",
        "equal_info_conventional_prediction":"ABSTAIN",
        "new_action_effect_representation_learned":false,
        "all_in_cpu_rss_io":"NOT_QUALIFIED",
        "scientific_first_use":false,
        "decision":"NO_GO_ACTION_EFFECT_INFORMATION_NOT_IDENTIFIABLE_FROM_CURRENT_PREACTION_VECTOR",
    }
    b,err:=json.Marshal(report)
    if err!=nil {t.Fatal(err)}
    t.Log("P6J_NATIVE_IDENTIFIABILITY_RECEIPT="+string(b))
}

func TestNCI04P4P6JNoOutcomeLeakageOnUnknownAction(t *testing.T){
    p6fPinOriginalSource(t)
    trained,_,_,_:=p6hNativeTrain(t)
    frozen,err:=trained.Snapshot()
    if err!=nil {t.Fatal(err)}
    before:=[FeatureWidth]float32{1,1,-1,0}
    prediction,known,err:=trained.PredictNativeTransition(before,1)
    if err!=nil || known || prediction!=([FeatureWidth]float32{}) {
        t.Fatal("unknown native action made confident pre-action prediction",err)
    }
    now,err:=trained.Snapshot()
    if err!=nil||!bytes.Equal(frozen,now) {t.Fatal("unknown future outcome mutated learning state",err)}
}
