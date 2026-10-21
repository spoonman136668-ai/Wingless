package cognition

// P6G exercises the original Git-pinned Wingless native read_file executor.
// Captured facts precede execution; no original executor/model changes or science first use.
import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "os"
    "path/filepath"
    "testing"
    "time"

    tb "github.com/spoonman136668-ai/Wingless/toolboundary"
)

func TestNCI04P4P6GOriginalWorkspacePreaction(t *testing.T) {
    p6fPinOriginalSource(t)
    root:=t.TempDir()
    p:=tb.Policy{Workspace:root, AllowedTools:[]string{"read_file"}, AllowedPaths:[]string{"ok.txt"}, Timeout:time.Second, MaxOutputBytes:100}
    action:=tb.Proposal{Tool:"read_file", Path:"ok.txt"}

    missing,err:=tb.ObservePreaction(context.Background(),p,action)
    if err!=nil || missing.Resource!="missing" || missing.SHA256!="" {t.Fatalf("missing-file observation: %+v %v",missing,err)}
    _,err=tb.Execute(context.Background(),p,action)
    if err==nil {t.Fatal("original missing-file native execution succeeded")}
    if err:=os.WriteFile(filepath.Join(root,"ok.txt"),[]byte("alpha"),0600);err!=nil {t.Fatal(err)}
    present,err:=tb.ObservePreaction(context.Background(),p,action)
    if err!=nil || present.Resource!="regular" || present.Size!=5 {t.Fatalf("existing-file observation: %+v %v",present,err)}
    expected:=sha256.Sum256([]byte("alpha"))
    if present.SHA256!=hex.EncodeToString(expected[:]) {t.Fatal("preaction bytes were not source-authenticated")}
    observed,err:=tb.Execute(context.Background(),p,action)
    if err!=nil || observed.Stdout!="alpha" || observed.ExitCode!=0 {t.Fatalf("original native valid read failed: %+v %v",observed,err)}

    if err:=os.WriteFile(filepath.Join(root,"ok.txt"),[]byte("bravo"),0600);err!=nil {t.Fatal(err)}
    changed,err:=tb.ObservePreaction(context.Background(),p,action)
    if err!=nil || changed.Resource!="regular" || changed.Size!=present.Size || changed.SHA256==present.SHA256 {
        t.Fatalf("same-length content change was missed: %+v %v",changed,err)
    }
    observed,err=tb.Execute(context.Background(),p,action)
    if err!=nil || observed.Stdout!="bravo" {t.Fatal("changed original file native call not completed",err)}

    // Equal public information: a conventional rooted source inspector knows
    // the same existence signal. There is no learned advantage here.
    conventional:=[]bool{false,true,true}
    candidate:=[]bool{missing.Resource=="regular",present.Resource=="regular",changed.Resource=="regular"}
    for i:=range conventional {if conventional[i]!=candidate[i] {t.Fatal("candidate inconsistent with equally informed conventional source inspector")}}
    report:=map[string]any{
        "gate":"NCI-04P6G", "original_native_git_blob":p6fToolBlob,
        "three_original_native_calls":true,
        "candidate_observes_missing_vs_present":true,
        "candidate_observes_same_size_content_change":true,
        "conventional_source_inspection_equivalent":true,
        "candidate_superiority":false,
        "native_preaction_to_outcome_atomicity":"NOT_GUARANTEED_TOCTOU",
        "learner_mutation":false,"scientific_first_use":false,
        "classification":"AUTHENTIC_PREACTION_ENVIRONMENT_OBSERVATION_ONLY__NO_SCIENTIFIC_ADVANTAGE",
    }
    b,err:=json.Marshal(report);if err!=nil {t.Fatal(err)}
    t.Log("P6G_NATIVE_EVIDENCE="+string(b))
}

func TestNCI04P4P6GRejectsUnauthorizedAndUnobservable(t *testing.T) {
    p6fPinOriginalSource(t)
    root:=t.TempDir()
    p:=tb.Policy{Workspace:root,AllowedTools:[]string{"read_file"},AllowedPaths:[]string{"ok.txt","directory"},Timeout:time.Second,MaxOutputBytes:4}
    before:=tb.Proposal{Tool:"read_file",Path:"ok.txt"}
    for _,a:=range []tb.Proposal{
        {Tool:"shell",Path:"ok.txt"},
        {Tool:"read_file",Path:"../escape"},
        {Tool:"read_file",Path:"unknown"},
    } {
        if _,err:=tb.ObservePreaction(context.Background(),p,a);err==nil {t.Fatalf("unauthorized preaction observation: %+v",a)}
    }
    if err:=os.Mkdir(filepath.Join(root,"directory"),0700);err!=nil {t.Fatal(err)}
    if _,err:=tb.ObservePreaction(context.Background(),p,tb.Proposal{Tool:"read_file",Path:"directory"});err==nil {t.Fatal("nonregular source was admitted")}
    if err:=os.WriteFile(filepath.Join(root,"ok.txt"),[]byte("too long"),0600);err!=nil {t.Fatal(err)}
    if _,err:=tb.ObservePreaction(context.Background(),p,before);err==nil {t.Fatal("over-limit preaction read was admitted")}
    ctx,cancel:=context.WithCancel(context.Background());cancel()
    if _,err:=tb.ObservePreaction(ctx,p,before);err==nil {t.Fatal("cancelled preaction read was admitted")}
}
