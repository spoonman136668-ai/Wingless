package main

import (
 "os"
 "path/filepath"
 "testing"
)
func TestNCI04P2NativeCLIRejectsTargetAnswersAndUntrustedSnapshotPaths(t *testing.T) {
 for _,in:=range [][2]string{{"","-"},{"-1","-"},{"3","-"},{"1","1"},{"1","outcome=0"}} {
  if _,err:=checkedNativeAction(in[0],in[1]);err==nil{t.Fatalf("accepted untrusted native action/answer %q",in)}
 }
 for i:=0;i<3;i++ {
  s:=string(rune('0'+i))
  if n,err:=checkedNativeAction(s,"-");err!=nil||n!=i{t.Fatalf("valid bounded source action rejected: %d %v",n,err)}
 }
 root:=t.TempDir()
 dir:=filepath.Join(root,"nested")
 if err:=os.Mkdir(dir,0700);err!=nil{t.Fatal(err)}
 if err:=checkedNativeSnapshot(root,filepath.Join(dir,"state.json"));err==nil {
  t.Fatal("mutable native event receipt within indexed source root was accepted")
 }
 out:=t.TempDir()
 if err:=checkedNativeSnapshot(root,filepath.Join(out,"state.json"));err!=nil{
  t.Fatal("external durable state file unexpectedly rejected",err)
 }
}
