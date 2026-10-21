package cognition

import (
 "errors"
 "fmt"
 "strings"
 "testing"
)

func TestNCI04P2NativeProjectAssertionOwnership(t *testing.T){
 task:=nativeReasoningTest{Name:"TestSourceWins",Package:"./ice",Source:"ice/index_test.go"}
 pkg:="github.com/spoonman136668-ai/Wingless/ice"
 mk:=func(action,test,output string)string{
  fields:=fmt.Sprintf(`"Action":%q,"Package":%q`,action,pkg)
  if test!=""{fields+=fmt.Sprintf(`,"Test":%q`,test)}
  if output!=""{fields+=fmt.Sprintf(`,"Output":%q`,output)}
  return "{"+fields+"}
"
 }
 success:=mk("start","","")+mk("run",task.Name,"")+mk("pass",task.Name,"")+mk("pass","","")
 label,err:=nativeGoTestVerdict([]byte(success),task,nil)
 if err!=nil||label!=1{t.Fatalf("valid executable success denied: label=%d err=%v",label,err)}
 failure:=mk("start","","")+mk("run",task.Name,"")+
  mk("output",task.Name,"    source contract assertion failed
")+
  mk("fail",task.Name,"")+mk("fail","","")
 label,err=nativeGoTestVerdict([]byte(failure),task,errors.New("exit status 1"))
 if err!=nil||label!=0{t.Fatalf("real test assertion negative could not be represented: label=%d err=%v",label,err)}
 t.Log("NCI04P2_EXECUTION_RECEIPT_PARSER_PASS; simulated JSON is parser-unit evidence only, not novel native task execution")
}
func TestNCI04P2RejectsNoisyNonexecutionsAndMismatchedVerdicts(t *testing.T){
 task:=nativeReasoningTest{Name:"TestSourceWins",Package:"./ice",Source:"ice/index_test.go"}
 pkg:="github.com/spoonman136668-ai/Wingless/ice"
 ev:=func(a,p,test,output string)string {
  return fmt.Sprintf(`{"Action":%q,"Package":%q,"Test":%q,"Output":%q}`+"
",a,p,test,output)
 }
 run:=ev("run",pkg,task.Name,"")
 pass:=run+ev("pass",pkg,task.Name,"")+ev("pass",pkg,"","")
 fail:=run+ev("fail",pkg,task.Name,"")+ev("fail",pkg,"","")
 cases:=[]struct{name,data string;status error}{
  {"no_test_run",ev("fail",pkg,"","compilation errors"),errors.New("build failed")},
  {"plain_text","FAIL
",errors.New("exit status 1")},
  {"wrong_package",ev("run","evil/package",task.Name,"")+ev("pass","evil/package",task.Name,"")+ev("pass","evil/package","",""),nil},
  {"wrong_test",ev("run",pkg,"TestOther","")+ev("pass",pkg,"TestOther","")+ev("pass",pkg,"",""),nil},
  {"no_summary",run+ev("pass",pkg,task.Name,""),nil},
  {"no_complete_test",run+ev("pass",pkg,"",""),nil},
  {"test_skip",run+ev("skip",pkg,task.Name,"")+ev("pass",pkg,"",""),nil},
  {"conflict_process_success",pass,errors.New("exit status 1")},
  {"conflict_process_failure",fail,nil},
  {"panic_in_test",run+ev("output",pkg,task.Name,"panic: target exploded")+ev("fail",pkg,task.Name,"")+ev("fail",pkg,"",""),errors.New("exit status 1")},
  {"duplicate_test",pass+ev("pass",pkg,task.Name,""),nil},
  {"oversized",strings.Repeat("X",1<<20+1),nil},
 }
 for _,tc:=range cases{
  t.Run(tc.name,func(t *testing.T){
   if label,err:=nativeGoTestVerdict([]byte(tc.data),task,tc.status);err==nil{
    t.Fatalf("unverified execution was converted to learned outcome %d",label)
   }
  })
 }
}
