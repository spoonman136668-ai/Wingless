package cognition

// P6I native cross-family falsifier: no new learner, action executor or model.
// Source-changing ICE index repair occurs exclusively in a disposable exact-HEAD
// clone. Original P6H read_file training remains the only learner experience.
// This is a development-selected family, NOT independently held scientific first use.
import (
    "bytes"
    "context"
    "encoding/json"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "testing"
    "time"
    "unsafe"

    "github.com/spoonman136668-ai/Wingless/ice"
    tb "github.com/spoonman136668-ai/Wingless/toolboundary"
)

func p6iGit(t *testing.T, root string, args ...string)string {
    t.Helper()
    cmd:=exec.Command("git",args...)
    cmd.Dir=root
    out,err:=cmd.CombinedOutput()
    if err!=nil {t.Fatalf("exact Git command failed: %v %q %s",err,args,out)}
    return strings.TrimSpace(string(out))
}

func TestNCI04P4P6IAuthenticStateChangingNativeTransfer(t *testing.T) {
    p6fPinOriginalSource(t)
    // Freeze the P6H source-native learned state BEFORE a new family is used.
    trained,_,_,_:=p6hNativeTrain(t)
    frozen,err:=trained.Snapshot()
    if err!=nil {t.Fatal(err)}
    if unsafe.Sizeof(trained)!=6472 {t.Fatal("learner fixed-capacity byte drift")}
    sourceRoot,err:=filepath.Abs("..")
    if err!=nil {t.Fatal(err)}
    head:=p6iGit(t,sourceRoot,"rev-parse","HEAD")
    if len(head)!=40 {t.Fatal("unverifiable Git source identity")}
    clone:=filepath.Join(t.TempDir(),"isolated-original-native")
    cmd:=exec.Command("git","clone","--quiet","--no-hardlinks","--no-checkout",sourceRoot,clone)
    if out,err:=cmd.CombinedOutput();err!=nil {t.Fatalf("disposable local source clone failed: %v %s",err,out)}
    _=p6iGit(t,clone,"checkout","--quiet","--detach",head)
    if actual:=p6iGit(t,clone,"rev-parse","HEAD");actual!=head {t.Fatal("native action source identity drift")}
    if dirty:=p6iGit(t,clone,"status","--porcelain");dirty!="" {t.Fatalf("scratch source was not clean: %s",dirty)}

    // The new family is the existing native ice.Save(root), a genuine
    // state-changing materialization of canonical source-index artifacts.
    // The learner sees ONLY public, allowlisted PREACTION source facts.
    policy:=tb.Policy{
        Workspace:clone,AllowedTools:[]string{"read_file"},
        AllowedPaths:[]string{"ice/index.go"},Timeout:time.Second,
        MaxOutputBytes:65536,
    }
    proposal:=tb.Proposal{Tool:"read_file",Path:"ice/index.go"}
    permitted,err:=tb.ObservePreaction(context.Background(),policy,proposal)
    if err!=nil || permitted.Resource!="regular" {t.Fatal("untrusted native pre-action source:",err)}
    before,err:=NativePreactionFeatures(permitted)
    if err!=nil {t.Fatal(err)}
    // P6H was trained exclusively on action 0 (read_file). Action 1 is
    // purposefully a distinct native state-changing family, not a relabel.
    const repairAction=1
    predicted,known,predictErr:=trained.PredictNativeTransition(before,repairAction)
    if predictErr!=nil {t.Fatal(predictErr)}
    if known {t.Fatalf("unearned cross-family confident transfer: %v",predicted)}

    // Conventional code/source inspector receives the SAME genuine native
    // tree and can compute the exact canonical index pre-action, without
    // looking at the learner or reading the future Save outcome.
    existing,err:=os.ReadFile(filepath.Join(clone,".ice","manifest.json"))
    if err!=nil {t.Fatal(err)}
    canonical,err:=ice.Build(clone)
    if err!=nil {t.Fatal("native conventional source inspection failed:",err)}
    wouldWrite,err:=json.MarshalIndent(canonical,"","  ")
    if err!=nil {t.Fatal(err)}
    if bytes.Equal(existing,append(wouldWrite,'\n')) {
        t.Fatal("no real index repair needed; native state-change precondition absent")
    }
    if err:=ice.Validate(clone);err==nil {
        t.Fatal("source-bound indexed state unexpectedly fresh before native repair")
    }

    started:=time.Now()
    if err:=ice.Save(clone);err!=nil {t.Fatal("original native ICE.Save failed:",err)}
    repairWall:=time.Since(started)
    if err:=ice.Validate(clone);err!=nil {t.Fatal("native state-changing ICE.Save did not repair canonical source index:",err)}
    after,err:=os.ReadFile(filepath.Join(clone,".ice","manifest.json"))
    if err!=nil {t.Fatal(err)}
    if bytes.Equal(existing,after) || !bytes.Equal(after,append(wouldWrite,'\n')) {
        t.Fatal("state-changing action did not produce the source-predicted canonical artifact")
    }
    changed:=p6iGit(t,clone,"status","--porcelain")
    if changed=="" {t.Fatal("source-changing native action produced no artifacts")}
    // git status --porcelain uses XY + space + path. p6iGit trims the
    // left padding on status lines; normalize XY rather than indexing at 3.
    for _,line:=range strings.Split(changed,"\n"){
        fields:=strings.Fields(line)
        if len(fields)!=2 || !strings.HasPrefix(fields[1],".ice/") {
            t.Fatalf("original source escaped isolated index-only action: %q",line)
        }
    }
    if p6iGit(t,clone,"rev-parse","HEAD")!=head {t.Fatal("original source HEAD changed during native action")}

    afterLearner,err:=trained.Snapshot()
    if err!=nil || !bytes.Equal(afterLearner,frozen) {
        t.Fatal("held-out native action mutated the P6H learner")
    }
    report:=map[string]any{
        "gate":"NCI-04P6I",
        "source_commit":head,
        "native_training_family":"toolboundary.Execute / read_file",
        "different_native_action_family":"ice.Save / canonical source-index generation",
        "separate_original_go_package":true,
        "source_was_stale_before_action":true,
        "native_source_index_repaired":true,
        "native_original_index_artifact_only":true,
        "source_restore_scope":"disposable exact-HEAD Git clone removed after test",
        "learner_prior_on_new_action":"ABSTAIN",
        "learner_training_after_holdout":false,
        "conventional_source_inspection_predicts_artifact":true,
        "candidate_advantage":false,
        "learner_fixed_bytes":unsafe.Sizeof(trained),
        "learner_snapshot_bytes":len(frozen),
        "native_repair_wall_ns":repairWall.Nanoseconds(),
        "all_in_equal_cpu_rss_io":"NOT_QUALIFIED",
        "independent_prospective_first_use":false,
        "scientific_transfer":"NOT_PROVEN",
        "classification":"P6I_REAL_STATE_CHANGE_AND_CROSS_FAMILY_ABSTENTION__DEV_NEGATIVE",
    }
    output,err:=json.Marshal(report)
    if err!=nil {t.Fatal(err)}
    t.Log("P6I_NATIVE_CROSSFAMILY_RECEIPT="+string(output))
}

func TestNCI04P4P6IUnseenNativeActionAbstainsWithoutLearnedEvidence(t *testing.T){
    p6fPinOriginalSource(t)
    trained,_,_,_:=p6hNativeTrain(t)
    snapshot,err:=trained.Snapshot()
    if err!=nil {t.Fatal(err)}
    pre:=[FeatureWidth]float32{1,1,-1,0}
    after,known,err:=trained.PredictNativeTransition(pre,1)
    if err!=nil || known || after!=([FeatureWidth]float32{}) {
        t.Fatal("unseen source-native action predicted without genuine learned evidence",err)
    }
    current,err:=trained.Snapshot()
    if err!=nil || !bytes.Equal(snapshot,current) {t.Fatal("pre-action prediction modified learning state")}
}
