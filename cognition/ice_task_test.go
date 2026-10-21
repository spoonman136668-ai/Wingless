package cognition

import (
 "bytes"
 "os"
 "path/filepath"
 "testing"
 "github.com/spoonman136668-ai/Wingless/ice"
)
func TestICEStateTaskRoundtripAndStrongExactControl(t *testing.T){
 root:=t.TempDir()
 code:=[]byte("package sample\nfunc Old() { New() }\nfunc New() {}\nfunc Other() {}\n")
 if err:=os.WriteFile(filepath.Join(root,"sample.go"),code,0600);err!=nil{t.Fatal(err)}
 if err:=ice.Save(root);err!=nil{t.Fatal(err)}
 session,err:=OpenSourceReasoner(root);if err!=nil{t.Fatal(err)}
 before,err:=session.PredictCall("Old","New");if err!=nil{t.Fatal(err)}
 if !before.VerifiedSyntaxCall||before.LearnedProbability!=.5{t.Fatalf("baseline: %#v",before)}
 for i:=0;i<180;i++{if err:=session.ObserveExistingCall("Old","New");err!=nil{t.Fatal(err)}}
 after,err:=session.PredictCall("Old","New");if err!=nil||after.LearnedProbability<=.9{t.Fatalf("no learned improvement: %#v %v",after,err)}
 if err:=session.ObserveExistingCall("Old","Other");err==nil{t.Fatal("unverified call accepted")}
 memory:=filepath.Join(t.TempDir(),"state.json")
 if err:=session.Save(memory);err!=nil{t.Fatal(err)}
 restarted,err:=OpenSourceReasoner(root);if err!=nil{t.Fatal(err)}
 if err=restarted.Load(memory);err!=nil{t.Fatal(err)}
 restored,err:=restarted.PredictCall("Old","New");if err!=nil{t.Fatal(err)}
 if after.LearnedProbability!=restored.LearnedProbability{t.Fatal("restart lost knowledge")}
 prior,err:=restarted.State.Snapshot();if err!=nil{t.Fatal(err)}
 _,err=restarted.PredictCall("Old","Other");if err!=nil{t.Fatal(err)}
 following,_:=restarted.State.Snapshot();if !bytes.Equal(prior,following){t.Fatal("query changed persistent state")}
 // The source index is deliberately stale after a source edit. Fail closed.
 if err=os.WriteFile(filepath.Join(root,"sample.go"),append(code,[]byte("func Drift(){}\n")...),0600);err!=nil{t.Fatal(err)}
 if err=restarted.ObserveExistingCall("Old","New");err==nil{t.Fatal("stale source accepted")}
 if err=restarted.Save(memory);err==nil{t.Fatal("saved state against stale source")}
}
func TestICEStateAmbiguityAndNoInventedNegatives(t *testing.T){
 root:=t.TempDir()
 if err:=os.WriteFile(filepath.Join(root,"sample.go"),[]byte("package sample\nfunc Yes() {}\nfunc No() {}\n"),0600);err!=nil{t.Fatal(err)}
 if err:=ice.Save(root);err!=nil{t.Fatal(err)}
 s,err:=OpenSourceReasoner(root);if err!=nil{t.Fatal(err)}
 if err=s.ObserveExistingCall("Yes","No");err==nil{t.Fatal("absent source fact treated as negative supervision")}
 if _,err=s.PredictCall("Missing","Yes");err==nil{t.Fatal("missing source identity accepted")}
}
