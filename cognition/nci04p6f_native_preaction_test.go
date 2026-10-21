package cognition

// P6F native read-only original-source falsifier.
// Candidate-only test: no learner updates, production changes or science first use.
import (
    "context"
    "crypto/sha1"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "reflect"
    "testing"
    "time"

    tb "github.com/spoonman136668-ai/Wingless/toolboundary"
)

const p6fToolBlob = "9775bc0ecd7b9946bc0e8fa7c3f67f9780d00a18"
const p6fToolSHA256 = "46e90c5a05df8887756dfd5f9cdedbb303ef59ce7b755cd5a43e546edabb6480"

func p6fPinOriginalSource(t *testing.T) {
    t.Helper()
    b, err := os.ReadFile("../toolboundary/tools.go")
    if err != nil { t.Fatal(err) }
    gitBlob := sha1.Sum(append([]byte(fmt.Sprintf("blob %d\x00", len(b))), b...))
    sha := sha256.Sum256(b)
    if hex.EncodeToString(gitBlob[:]) != p6fToolBlob || hex.EncodeToString(sha[:]) != p6fToolSHA256 {
        t.Fatal("original Wingless native toolboundary source identity drift")
    }
}

type p6fArguments struct {
    Tool string
    Path string
    Workspace string
    AllowedTools []string
    AllowedPaths []string
    Timeout time.Duration
    MaxOutputBytes int
}
func p6fPre(p tb.Policy, a tb.Proposal) p6fArguments {
    return p6fArguments{
        Tool:a.Tool, Path:a.Path, Workspace:p.Workspace,
        AllowedTools:append([]string(nil),p.AllowedTools...),
        AllowedPaths:append([]string(nil),p.AllowedPaths...),
        Timeout:p.Timeout,MaxOutputBytes:p.MaxOutputBytes,
    }
}
func p6fExec(ctx context.Context,p tb.Policy,a tb.Proposal) (p6fArguments,bool,error) {
    pre := p6fPre(p,a) // pre-action snapshot, before any native outcome exists
    observed,err := tb.Execute(ctx,p,a) // unmodified original native executor
    return pre,err==nil && observed.ExitCode==0,err
}
func p6fPermitted(p tb.Policy,a tb.Proposal) bool {
    tool,path := false,false
    for _,v := range p.AllowedTools { if v==a.Tool {tool=true} }
    for _,v := range p.AllowedPaths { if v==a.Path {path=true} }
    return a.Tool=="read_file" && tool && path
}
func TestNCI04P4P6FNativePreactionArgumentsAreInsufficient(t *testing.T) {
    p6fPinOriginalSource(t)
    root:=t.TempDir()
    p:=tb.Policy{Workspace:root,AllowedTools:[]string{"read_file"},AllowedPaths:[]string{"ok.txt"},Timeout:time.Second,MaxOutputBytes:100}
    a:=tb.Proposal{Tool:"read_file",Path:"ok.txt"}
    beforeA,passedA,errA:=p6fExec(context.Background(),p,a)
    if errA==nil || passedA {t.Fatal("original native missing-file must fail")}
    if err:=os.WriteFile(filepath.Join(root,"ok.txt"),[]byte("hello"),0600);err!=nil {t.Fatal(err)}
    beforeB,passedB,errB:=p6fExec(context.Background(),p,a)
    if errB!=nil || !passedB {t.Fatalf("original native present-file must pass: %v",errB)}
    if !reflect.DeepEqual(beforeA,beforeB) {t.Fatal("pre-action arguments unexpectedly differed")}
    if !p6fPermitted(p,a) {t.Fatal("original source guard did not admit calls")}
    result:=map[string]any{
        "gate":"NCI-04P6F",
        "original_native_git_blob":p6fToolBlob,
        "native_go_calls":2,
        "same_preaction_arguments":true,
        "missing_file_pass":passedA,
        "present_file_pass":passedB,
        "entry_guard_admits_both":true,
        "arguments_only_sufficient":false,
        "additional_required_preaction_signal":"native rooted filesystem existence",
        "learner_mutation":false,
        "scientific_first_use":false,
        "classification":"P6F_AUTHENTIC_NATIVE_ARGUMENT_ONLY_OBSERVABILITY_NO_GO",
    }
    b,err:=json.Marshal(result);if err!=nil {t.Fatal(err)}
    t.Log("P6F_NATIVE_EVIDENCE="+string(b))
}
func TestNCI04P4P6FNativeOriginalAuthorization(t *testing.T) {
    p6fPinOriginalSource(t)
    root:=t.TempDir()
    if err:=os.WriteFile(filepath.Join(root,"ok.txt"),[]byte("hello"),0600);err!=nil {t.Fatal(err)}
    p:=tb.Policy{Workspace:root,AllowedTools:[]string{"read_file"},AllowedPaths:[]string{"ok.txt"},Timeout:time.Second,MaxOutputBytes:100}
    for _,a:=range []tb.Proposal{
        {Tool:"read_file",Path:"ok.txt"},
        {Tool:"shell",Path:"ok.txt"},
        {Tool:"read_file",Path:"unknown"},
        {Tool:"read_file",Path:"../escape"},
    } {
        pre,passed,err:=p6fExec(context.Background(),p,a)
        if pre.Tool!=a.Tool || pre.Path!=a.Path {t.Fatal("preaction receipt drift")}
        if !p6fPermitted(p,a) && err==nil {t.Fatal("original unauthorized action was accepted")}
        if (err==nil)!=passed {t.Fatal("native outcome accounting contradiction")}
    }
}
