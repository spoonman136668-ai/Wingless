package cognition

import (
 "bytes"
 "context"
 "os"
 "path/filepath"
 "testing"
 "unsafe"
)

// These are engineering-only observations of the exact generated, executed
// source fixture, not a new scientific result or learned-policy claim.
func TestNCI04G5NativeExecutableActionConsequence(t *testing.T) {
 root:=filepath.Join(t.TempDir(),"owned-probe")
 p,err:=NewNativeProbe(root);if err!=nil{t.Fatal(err)}
 original,err:=p.Snapshot();if err!=nil{t.Fatal(err)}
 if len(original)>512 || unsafe.Sizeof(*p)>128{t.Fatal("unbounded probe receipt/state")}
 obs0,err:=p.Intervene(context.Background(),false)
 if err!=nil{t.Fatal(err)}
 if obs0.Consequence!=0||obs0.Action!=0||obs0.Sequence!=1{t.Fatalf("incorrect observed native response %#v",obs0)}
 obs1,err:=p.Intervene(context.Background(),true)
 if err!=nil{t.Fatal(err)}
 if obs1.Consequence!=1||obs1.Action!=1||obs1.Sequence!=2{t.Fatalf("no distinguishing native response %#v",obs1)}
 if obs0.SourceSHA256==obs1.SourceSHA256 || obs0.CanonicalICESHA256==obs1.CanonicalICESHA256{
  t.Fatal("changed source did not have distinct authenticated source/index digests")
 }
 state,err:=p.Snapshot();if err!=nil{t.Fatal(err)}
 restarted,err:=OpenNativeProbe(root);if err!=nil{t.Fatal(err)}
 if err:=restarted.Restore(original);err==nil{t.Fatal("old source receipt accepted after native transition")}
 if err:=restarted.Restore(state);err!=nil{t.Fatal(err)}
 check,err:=restarted.Snapshot();if err!=nil||!bytes.Equal(state,check){t.Fatal("non-equivalent native restart",err)}
 obs2,err:=restarted.Intervene(context.Background(),false)
 if err!=nil||obs2.Sequence!=3||obs2.Consequence!=0{t.Fatal("failed third and final causal intervention",err)}
 if _,err=restarted.Intervene(context.Background(),true);err==nil{t.Fatal("hidden fourth intervention allowed")}
 // Fresh process without a receipt must still see the consumed durable budget.
 reopened,err:=OpenNativeProbe(root);if err!=nil{t.Fatal(err)}
 if _,err=reopened.Intervene(context.Background(),true);err==nil{
  t.Fatal("reopened sandbox reset durable action limit")
 }
 if err=reopened.Restore(state);err==nil{t.Fatal("stale receipt rewound durable sequence")}
 t.Logf("NCI04G5_WINGLESS_NATIVE_SEAM_PASS source_authenticated_interventions=3 receipt_bytes=%d resident_struct_bytes=%d",len(state),unsafe.Sizeof(*p))
}

func TestNCI04G5RejectUntrustedSourceOrExistingWorktree(t *testing.T){
 root:=filepath.Join(t.TempDir(),"probe")
 p,err:=NewNativeProbe(root);if err!=nil{t.Fatal(err)}
 if _,err:=NewNativeProbe(root);err==nil{t.Fatal("existing directory overwritten")}
 raw,err:=os.ReadFile(filepath.Join(root,"probe.go"));if err!=nil{t.Fatal(err)}
 if err=os.WriteFile(filepath.Join(root,"probe.go"),append(raw,[]byte("\n// altered source\n")...),0600);err!=nil{t.Fatal(err)}
 if _,err=p.Intervene(context.Background(),true);err==nil{t.Fatal("source identity drift accepted")}
 if _,err=OpenNativeProbe(root);err==nil{t.Fatal("unreviewed source accepted on restart")}
}

func TestNCI04G5CancelledExecutionCannotReplayReservedAction(t *testing.T) {
 root:=filepath.Join(t.TempDir(),"private")
 p,err:=NewNativeProbe(root);if err!=nil{t.Fatal(err)}
 ctx,cancel:=context.WithCancel(context.Background());cancel()
 if _,err:=p.Intervene(ctx,true);err==nil{t.Fatal("cancelled native execution accepted")}
 if _,err:=OpenNativeProbe(root);err==nil {
  t.Fatal("crash-window reopen accepted source/custody mismatch")
 }
 if _,err:=p.Intervene(context.Background(),false);err==nil {
  t.Fatal("attempt replayed reserved action after failed execution")
 }
}
