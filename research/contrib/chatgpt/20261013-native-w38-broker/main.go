// W-R38: experimental probe of native broker.Select authority, never production invocation.
package main

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "os"
    "time"

    "github.com/spoonman136668-ai/Wingless/broker"
    "github.com/spoonman136668-ai/Wingless/inference"
    "github.com/spoonman136668-ai/Wingless/resources"
)

type row struct {
    Index int `json:"index"`
    Flags string `json:"flags"`
    Oracle string `json:"oracle"`
    Native string `json:"native"`
    Naive string `json:"naive"`
}
type report struct {
    Schema string `json:"schema"`
    Classification string `json:"classification"`
    Parent string `json:"parent"`
    Cases int `json:"cases"`
    NativeVsOracleMismatches int `json:"native_vs_oracle_mismatches"`
    NaiveUnsafeSelections int `json:"naive_unsafe_selections"`
    NativeCounts map[string]int `json:"native_counts"`
    OracleCounts map[string]int `json:"oracle_counts"`
    CaseSha256 string `json:"case_sha256"`
    MismatchExamples []row `json:"mismatch_examples"`
    ModelCalls int `json:"model_calls"`
    Limit string `json:"claim_limit"`
}
func bit(n, shift int) bool {return n&(1<<shift)!=0}
func expected(b [10]bool) string {
    fh,dh,fc,dc,af,ad,hard,auth,fallback,resourceOK :=
        b[0],b[1],b[2],b[3],b[4],b[5],b[6],b[7],b[8],b[9]
    if !resourceOK || (hard && !auth) {return "DENY"}
    if hard {
        if ad && dh && dc {return "DEEP"}
        if fallback && af && fh && fc {return "FAST"}
        return "DENY"
    }
    if af && fh && fc {return "FAST"}
    return "DENY"
}
func naive(b [10]bool) string {
    if b[4] {return "FAST"}
    if b[5] {return "DEEP"}
    return "DENY"
}
func flagString(b [10]bool) string {
    s := ""
    for _, x := range b {if x {s+="1"} else {s+="0"}}
    return s
}
func main() {
    rep:=report{
        Schema:"wingless.native-broker-authority.w38.v1",
        Parent:"5d84c80d07fa8aee18ebe254577753136d350628",
        Cases:1024,
        NativeCounts:map[string]int{"DENY":0,"FAST":0,"DEEP":0},
        OracleCounts:map[string]int{"DENY":0,"FAST":0,"DEEP":0},
        ModelCalls:0,
        Limit:"Exhaustive fixed mock-backend authority routing, not inference capability or production authority.",
    }
    digest:=sha256.New()
    var sufficient uint64=2
    for index:=0;index<1024;index++ {
        var b [10]bool
        for i:=0;i<10;i++ {b[i]=bit(index,i)}
        flags:=flagString(b)
        reg:=&broker.Registry{}
        fastFeatures:=[]string{}
        deepFeatures:=[]string{}
        if b[2] {fastFeatures=[]string{"code"}}
        if b[3] {deepFeatures=[]string{"code"}}
        for _,e:=range []broker.Entry{
            {Backend:&inference.Mock{Name:"fast",Features:fastFeatures,Unhealthy:!b[0]},Class:"mock",Tier:"fast"},
            {Backend:&inference.Mock{Name:"deep",Features:deepFeatures,Unhealthy:!b[1]},Class:"mock",Tier:"deep"},
        } {
            if err:=reg.Register(e);err!=nil {panic(err)}
        }
        allowed:=[]string{}
        if b[4] {allowed=append(allowed,"fast")}
        if b[5] {allowed=append(allowed,"deep")}
        reason:="routine"
        if b[6] {reason="hard_failure"}
        metrics:=resources.Metrics{}
        if b[9] {metrics.RAMFree=&sufficient}
        req:=inference.Request{
            ID:"r",ParentWorkID:"p",Role:"code",Context:"x",
            MaxContextBytes:4096,MaxOutputTokens:32,
            Deadline:time.Now().Add(5*time.Minute),
            Workspace:"fixture",Capabilities:[]string{"code"},
            Resources:resources.Policy{MinRAM:1},
        }
        policy:=broker.Policy{AllowedBackends:allowed,Reason:reason,
            AllowDeep:b[7],AllowFallback:b[8],MaxRepairs:0}
        backend,decision,err:=broker.Select(context.Background(),reg,req,policy,metrics)
        native:="DENY"
        if err==nil {
            if backend==nil {panic("selected nil backend")}
            switch backend.ID() {
            case "fast":native="FAST"
            case "deep":native="DEEP"
            default:panic("unknown backend")
            }
            if decision.Backend!=backend.ID() {panic("mismatched decision")}
        } else if backend!=nil {panic("backend with error")}
        want:=expected(b)
        raw:=naive(b)
        rep.NativeCounts[native]++
        rep.OracleCounts[want]++
        if raw!="DENY" && raw!=want {rep.NaiveUnsafeSelections++}
        fmt.Fprintf(digest,"%04d:%s:%s:%s:%s\n",index,flags,native,want,raw)
        if native!=want {
            rep.NativeVsOracleMismatches++
            if len(rep.MismatchExamples)<16 {
                rep.MismatchExamples=append(rep.MismatchExamples,row{
                    Index:index,Flags:flags,Oracle:want,Native:native,Naive:raw,
                })
            }
        }
    }
    rep.CaseSha256=hex.EncodeToString(digest.Sum(nil))
    if rep.NativeVsOracleMismatches==0 && rep.NaiveUnsafeSelections>0 {
        rep.Classification="NATIVE_AUTHORITY_ENFORCEMENT"
    } else {rep.Classification="NEGATIVE_BROKER_AUTHORITY_FAILURE"}
    if err:=json.NewEncoder(os.Stdout).Encode(rep);err!=nil {panic(err)}
}
