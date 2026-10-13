// W-R39 native bounded verifier and repair: real broker.Runner, mock backends.
package main

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "time"

    "github.com/spoonman136668-ai/Wingless/broker"
    "github.com/spoonman136668-ai/Wingless/inference"
    "github.com/spoonman136668-ai/Wingless/resources"
)
type metricProvider struct{}
func (metricProvider) Snapshot() (resources.Metrics,error) {return resources.Metrics{},nil}
type fixtureVerifier struct{}
func (fixtureVerifier) Verify(_ context.Context, out inference.Result) error {
    if out.Text=="good" {return nil}
    return errors.New("fixture failed")
}
type row struct {
    Index int `json:"index"`
    Flags string `json:"flags"`
    NativeStatus string `json:"native_status"`
    ExpectedStatus string `json:"expected_status"`
    Attempts int `json:"attempts"`
    ExpectedAttempts int `json:"expected_attempts"`
    SecondBackend string `json:"second_backend"`
    ExpectedSecondBackend string `json:"expected_second_backend"`
    UnverifiedStatus string `json:"unverified_status"`
    Acceptance string `json:"acceptance"`
}
type output struct {
    Schema string `json:"schema"`
    Classification string `json:"classification"`
    Parent string `json:"parent"`
    TotalCases int `json:"total_cases"`
    OracleDisagreements int `json:"oracle_disagreements"`
    AcceptanceAuthorityViolations int `json:"acceptance_authority_violations"`
    VerifiedCandidates int `json:"verified_candidates"`
    OperatorBlocked int `json:"operator_blocked"`
    UnverifiedControlsResultReady int `json:"unverified_controls_result_ready"`
    CaseTraceSha256 string `json:"case_trace_sha256"`
    Mismatches []row `json:"mismatch_examples"`
    ModelCalls int `json:"model_calls"`
    ClaimLimit string `json:"claim_limit"`
}
func bit(n,k int)bool {return n&(1<<k)!=0}
func flags(n int)string {
    s:=""
    for k:=0;k<7;k++{if bit(n,k){s+="1"}else{s+="0"}}
    return s
}
func fixture(n int) (*broker.Registry,inference.Request,broker.Policy) {
    good:="bad"
    if bit(n,0) {good="good"}
    deepReply:="bad"
    if bit(n,4) {deepReply="good"}
    r:=&broker.Registry{}
    for _,e:=range []broker.Entry{
        {Backend:&inference.Mock{Name:"fast",Features:[]string{"code"},Reply:good},Class:"mock",Tier:"fast"},
        {Backend:&inference.Mock{Name:"deep",Features:[]string{"code"},Reply:deepReply,Unhealthy:!bit(n,3)},Class:"mock",Tier:"deep"},
    } {
        if err:=r.Register(e);err!=nil {panic(err)}
    }
    maxBytes:=1
    if bit(n,5){maxBytes=4096}
    req:=inference.Request{ID:"r",ParentWorkID:"p",Role:"code",Context:"x",
        MaxContextBytes:maxBytes,MaxOutputTokens:32,
        Deadline:time.Now().Add(5*time.Minute),Workspace:"fixture",
        Capabilities:[]string{"code"}}
    p:=broker.Policy{AllowedBackends:[]string{"fast","deep"},
        AllowDeep:bit(n,1),AllowFallback:bit(n,6)}
    if bit(n,2){p.MaxRepairs=1}
    return r,req,p
}
func oracle(n int)(string,int,string) {
    if bit(n,0) {return "verified_candidate",1,""}
    if !bit(n,1)||!bit(n,2)||!bit(n,5) {
        return "operator_blocked",1,""
    }
    if bit(n,3) {
        status:="operator_blocked"
        if bit(n,4){status="verified_candidate"}
        return status,2,"deep"
    }
    if bit(n,6) {return "operator_blocked",2,"fast"}
    return "operator_blocked",2,""
}
func main(){
    result:=output{
        Schema:"wingless.w-r39-native-verifier-boundary.v1",
        Parent:"8ac505090055ff2eefa4c177362e94d823cdd0eb",
        TotalCases:128,ModelCalls:0,
        ClaimLimit:"Real native Go broker.Runner verified using local fixed mock replies, not actual reasoning nor live authority.",
    }
    digest:=sha256.New()
    for n:=0;n<128;n++ {
        reg,req,p:=fixture(n)
        native:=(broker.Runner{Registry:reg,Metrics:metricProvider{},Verifier:fixtureVerifier{}}).
            Run(context.Background(),req,p)
        controlPolicy:=p
        controlPolicy.MaxRepairs=0
        control:=(broker.Runner{Registry:reg,Metrics:metricProvider{}}).
            Run(context.Background(),req,controlPolicy)
        status,count,second:=oracle(n)
        secondObserved:=""
        if len(native.Events)>1{secondObserved=native.Events[1].Route.Backend}
        mismatch:=native.Status!=status||len(native.Events)!=count||secondObserved!=second||
            control.Status!="result_ready"
        for _,x:=range []broker.Outcome{native,control} {
            if x.Acceptance!="external_required"{
                result.AcceptanceAuthorityViolations++
                mismatch=true
            }
        }
        if native.Status=="verified_candidate"{result.VerifiedCandidates++}
        if native.Status=="operator_blocked"{result.OperatorBlocked++}
        if control.Status=="result_ready"{result.UnverifiedControlsResultReady++}
        fmt.Fprintf(digest,"%03d:%s:%s:%d:%s:%s\n",n,flags(n),
            native.Status,len(native.Events),secondObserved,control.Status)
        if mismatch {
            result.OracleDisagreements++
            if len(result.Mismatches)<16 {
                result.Mismatches=append(result.Mismatches,row{
                    Index:n,Flags:flags(n),NativeStatus:native.Status,
                    ExpectedStatus:status,Attempts:len(native.Events),
                    ExpectedAttempts:count,SecondBackend:secondObserved,
                    ExpectedSecondBackend:second,UnverifiedStatus:control.Status,
                    Acceptance:native.Acceptance,
                })
            }
        }
    }
    result.CaseTraceSha256=hex.EncodeToString(digest.Sum(nil))
    if result.OracleDisagreements==0 && result.AcceptanceAuthorityViolations==0 &&
        result.VerifiedCandidates>0 && result.OperatorBlocked>0 {
        result.Classification="NATIVE_VERIFIER_BOUNDARY"
    } else{result.Classification="NEGATIVE_BOUNDARY_FAILURE"}
    if err:=json.NewEncoder(os.Stdout).Encode(result);err!=nil{panic(err)}
}
