package cognition

// NCI-04G7A: a controlled *native source reasoning* counterfactual.
// The executor analyzes one exact committed Wingless Go function in a private
// temporary ICE index, then alters ONE admissible call site inside the copy.
// It changes the real SourceReasoner source-analysis task outcome without
// editing any original source. This is source-bound development evidence, not
// an intervention on a running environment or independent scientific first-use.
import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "fmt"
 "os"
 "os/exec"
 "path/filepath"
 "testing"
 "time"
 "github.com/spoonman136668-ai/Wingless/ice"
)

const causalSourcePath="cmd/wingless/ice_cognition.go"
const originalCall="session.ObserveExistingCall(from,to)"
const changedCall="session.Save(snapshot)"

func TestNCI04G7ARealSourceCausalIntervention(t *testing.T){
 repo,err:=filepath.Abs("..");if err!=nil{t.Fatal(err)}
 source,err:=os.ReadFile(filepath.Join(repo,causalSourcePath));if err!=nil{t.Fatal(err)}
 cmd:=exec.Command("git","show","HEAD:"+causalSourcePath);cmd.Dir=repo
 committed,err:=cmd.Output();if err!=nil{t.Fatal(err)}
 if !bytes.Equal(committed,source){t.Fatal("working source differs from verified Git original")}
 if bytes.Count(source,[]byte(originalCall))!=1{
  t.Fatal("causal intervention source identity no longer unique")
 }
 originalSHA:=sha256.Sum256(source)
 root:=t.TempDir()
 file:=filepath.Join(root,"observed.go")
 if err:=os.WriteFile(file,source,0600);err!=nil{t.Fatal(err)}
 sourceRuntime:=time.Now()
 if err:=ice.Save(root);err!=nil{t.Fatal(err)}
 baselineSymbols,err:=ice.Query(root,"runICESourceCognition");if err!=nil{t.Fatal(err)}
 baselineElapsed:=time.Since(sourceRuntime)
 if len(baselineSymbols)!=1{t.Fatal("native source analysis did not resolve unique task")}
 from:=baselineSymbols[0]
 observedSourceCall:=ice.Symbol{Name:"ObserveExistingCall",Kind:"method"}
 before:=lexicalCall(from,observedSourceCall)
 if !before{t.Fatal("authentic source callable not recorded by ICE")}
 // Intervention is a bounded, explicit source-code edit in a disposable file.
 // This intentionally changes the source analysis outcome; it cannot change
 // project source or become a silent negative label for an unrelated call.
 edited:=bytes.Replace(source,[]byte(originalCall),[]byte(changedCall),1)
 if len(edited)!=len(source)-len(originalCall)+len(changedCall) {t.Fatal("unexpected edit")}
 if err:=os.WriteFile(file,edited,0600);err!=nil{t.Fatal(err)}
 if _,err:=ice.Query(root,"runICESourceCognition");err==nil{
  t.Fatal("ICE should reject the now-stale source as authoritative")
 }
 sourceRuntime=time.Now()
 if err:=ice.Save(root);err!=nil{t.Fatal(err)}
 editedSymbols,err:=ice.Query(root,"runICESourceCognition");if err!=nil{t.Fatal(err)}
 editedElapsed:=time.Since(sourceRuntime)
 if len(editedSymbols)!=1{t.Fatal("causal task symbol disappeared")}
 after:=lexicalCall(editedSymbols[0],observedSourceCall)
 if after||before==after{t.Fatalf("native causal intervention did not change observed syntax call: %v -> %v",before,after)}
 // Strong conventional exact-string control sees the identical bytes; no
 // learned advantage is earned by source indexing or a trivial intervention.
 controlBefore:=bytes.Contains(source,[]byte(originalCall))
 controlAfter:=bytes.Contains(edited,[]byte(originalCall))
 if before!=controlBefore||after!=controlAfter{t.Fatal("strong exact-source control disagrees")}
 // Preserve the original state exactly: no source mutations in original repo.
 final,err:=os.ReadFile(filepath.Join(repo,causalSourcePath));if err!=nil||!bytes.Equal(final,source){t.Fatal("original source was changed")}
 editedSHA:=sha256.Sum256(edited)
 report:=map[string]any{
  "gate":"NCI-04G7A-W",
  "source_file":causalSourcePath,
  "original_git_verified":true,
  "original_sha256":hex.EncodeToString(originalSHA[:]),
  "intervention_private_copy":true,
  "edited_sha256":hex.EncodeToString(editedSHA[:]),
  "native_workflow":"ICE Save + Query on one exact native Wingless function",
  "pre_intervention_verified_syntax_call":before,
  "post_intervention_verified_syntax_call":after,
  "stale_index_fails_closed":true,
  "trained_exact_string_control_equivalent":true,
  "native_original_index_wall_ns":baselineElapsed.Nanoseconds(),
  "native_intervened_index_wall_ns":editedElapsed.Nanoseconds(),
  "source_bytes_before":len(source),"source_bytes_after":len(edited),
  "simulation_or_physical_world":"CONTROLLED_DISPOSABLE_NATIVE_SOURCE_COPY",
  "all_in_matched_cpu_rss":"NOT_QUALIFIED: no independent same-budget recurrent control and model memory cost; test time includes indexing",
  "novel_causal_transfer":"NOT_ESTABLISHED",
  "independent_scientific_first_use":false,
  "disposition":"NECESSARY_SOURCE_CAUSAL_EFFECT_PASS; SCIENTIFIC_NO_GO"}
 j,_:=json.Marshal(report);fmt.Println("NCI04G7A_W_SOURCE_CAUSAL_EVIDENCE="+string(j))
}
