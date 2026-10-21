package cognition

// P6K is a pre-action affordance-information audit, not a new learned
// representation or a scientific result. It examines ORIGINAL native Go source
// AST before running either operation and never includes the operation name,
// post-action outcome, or source SHA in the 4-float observation vector.
import (
    "bytes"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "errors"
    "go/ast"
    "go/parser"
    "go/token"
    "os"
    "os/exec"
    "path/filepath"
    "reflect"
    "strings"
    "testing"
    "unsafe"
)

const p6kOriginalICEBlob = "b1f88d298a639244064c1f81a08db1c657d83c3a"

func p6kOriginalSource(t *testing.T) []byte {
    t.Helper()
    root,err:=filepath.Abs("..")
    if err!=nil{t.Fatal(err)}
    working,err:=os.ReadFile(filepath.Join(root,"ice","index.go"))
    if err!=nil{t.Fatal(err)}
    cmd:=exec.Command("git","show","HEAD:ice/index.go")
    cmd.Dir=root
    committed,err:=cmd.Output()
    if err!=nil || !bytes.Equal(working,committed){
        t.Fatal("original native action source differs from exact committed blob",err)
    }
    cmd=exec.Command("git","hash-object","--no-filters","ice/index.go")
    cmd.Dir=root
    blob,err:=cmd.Output()
    if err!=nil || strings.TrimSpace(string(blob))!=p6kOriginalICEBlob {
        t.Fatal("original ICE Go source blob differs from frozen action evidence",err)
    }
    return working
}

// Structural features are available before native execution, with no API
// name/destination/return code, yet cannot guarantee semantic effect parity.
// Selecting a declared entry point is provided by the caller's actual intended
// action, not discovered from a future test result.
func p6kASTPreactionStructure(original []byte, selectedEntrypoint string)([FeatureWidth]float32,error) {
    var features [FeatureWidth]float32
    if len(original)==0 || len(original)>131072 || selectedEntrypoint==""{
        return features,errors.New("invalid bounded original source")
    }
    f,err:=parser.ParseFile(token.NewFileSet(),"verified-original-native.go",original,0)
    if err!=nil{return features,err}
    var target *ast.FuncDecl
    for _,d:=range f.Decls{
        x,ok:=d.(*ast.FuncDecl)
        if !ok || x.Body==nil || x.Name.Name!=selectedEntrypoint {continue}
        if target!=nil{return features,errors.New("ambiguous native action entrypoint")}
        target=x
    }
    if target==nil{return features,errors.New("native action entrypoint unavailable")}
    var counts [FeatureWidth]uint32
    ast.Inspect(target.Body,func(n ast.Node)bool{
        switch n.(type){
        case *ast.CallExpr:counts[0]++
        case *ast.IfStmt:counts[1]++
        case *ast.RangeStmt:counts[2]++
        case *ast.ReturnStmt:counts[3]++
        }
        return true
    })
    for i,v:=range counts {
        if v>127{return features,errors.New("native structural feature overflow")}
        features[i]=float32(v)/127
    }
    return features,nil
}

func TestNCI04P4P6KPreactionActionStructureExists(t *testing.T){
    p6fPinOriginalSource(t)
    source:=p6kOriginalSource(t)
    build,err:=p6kASTPreactionStructure(source,"Build")
    if err!=nil{t.Fatal(err)}
    save,err:=p6kASTPreactionStructure(source,"Save")
    if err!=nil{t.Fatal(err)}
    if reflect.DeepEqual(build,save){
        t.Fatal("actual native preaction structure did not distinguish original Build and Save")
    }
    // Only the original native selected function name changes; features are
    // invariant to this change because names do not enter the descriptor.
    renamed:=bytes.Replace(source,[]byte("func Build("),[]byte("func UnnamedNativeEntrypoint("),1)
    if bytes.Equal(renamed,source){t.Fatal("native source rename control missing")}
    renamedFeatures,err:=p6kASTPreactionStructure(renamed,"UnnamedNativeEntrypoint")
    if err!=nil || !reflect.DeepEqual(renamedFeatures,build){
        t.Fatal("lexical function identity leaked into preaction structure",err)
    }
    if unsafe.Sizeof(LatentState{})!=6472{
        t.Fatal("learner state representation/capacity drift")
    }
    report:=map[string]any{
        "gate":"NCI-04P6K",
        "original_ice_go_git_blob":p6kOriginalICEBlob,
        "structure_available_preaction":true,
        "source_entrypoint_names_in_learner_vector":false,
        "build_feature":build,
        "save_feature":save,
        "structurally_distinguishable":true,
        "function_rename_invariance":true,
        "source_inspector_same_information":true,
        "new_learned_effect_representation":false,
        "effect_classification_or_semantic_guarantee":false,
        "source_independent_effect_transfer":"NOT_ESTABLISHED",
        "protected_scientific_first_use":false,
    }
    b,err:=json.Marshal(report)
    if err!=nil{t.Fatal(err)}
    t.Log("P6K_NATIVE_AFFORDANCE_RECEIPT="+string(b))
}

func TestNCI04P4P6KOriginalLiteralDecisionChangeCanBeInvisible(t *testing.T){
    source:=p6kOriginalSource(t)
    // Alter only a literal inside a genuine original Build decision.
    // Static AST node counts cannot see the changed condition semantics.
    replaced:=bytes.Replace(source,[]byte(`ext != ".go"`),[]byte(`ext != ".rs"`),1)
    if bytes.Equal(replaced,source){t.Fatal("original literal mutation was not present")}
    before,err:=p6kASTPreactionStructure(source,"Build")
    if err!=nil{t.Fatal(err)}
    after,err:=p6kASTPreactionStructure(replaced,"Build")
    if err!=nil{t.Fatal(err)}
    if !reflect.DeepEqual(before,after){
        t.Fatal("non-semantic preaction AST histogram changed on literal-only source edit")
    }
    h0,h1:=sha256.Sum256(source),sha256.Sum256(replaced)
    if h0==h1{t.Fatal("counterfactual source actually unchanged")}
    t.Log("P6K_PREACTION_SHAPE_LIMITATION="+
       `{"feature_equal":true,"source_bytes_differ":true,"semantic_effect_undetermined":true,"learned_transfer":false}`)
    t.Log("P6K_RAW_SOURCE_DIGEST="+hex.EncodeToString(h0[:]))
}
