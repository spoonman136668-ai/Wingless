package r1

import (
 "bytes"
 "encoding/base64"
 "math"
 "os"
 "os/exec"
 "reflect"
 "strconv"
 "testing"
)

func scene() PublicGrid {
 return PublicGrid{X:3,Y:3,GoalX:6,GoalY:5,Direction:0,FrontBlocked:false}
}

func TestR1AllThreeActionsCanBeLearned(t *testing.T) {
 p:=New()
 before:=scene()
 a,e:=p.Choose(before)
 if e!=nil||a!=0 {t.Fatalf("zero model %v %v",a,e)}
 if e:=p.Observe(before,1,before,1,true);e!=nil {t.Fatal(e)}
 a,e=p.Choose(before)
 if e!=nil||a!=1 {t.Fatalf("trained turn-right=%d error %v",a,e)}
 snapshot:=p.Fast[1]
 for i:=0;i<5;i++ {
  if e:=p.Observe(before,2,before,1,true);e!=nil {t.Fatal(e)}
 }
 a,e=p.Choose(before)
 if e!=nil||a!=2 {t.Fatalf("learned-forward=%d error %v",a,e)}
 if p.Fast[1]!=snapshot {t.Fatal("unselected arm update contaminated")}
 if p.Updates!=6{t.Fatalf("genuine public action count %d",p.Updates)}
}

func TestR1ValidPublicTransitionsAndAtomicFailures(t *testing.T) {
 p:=New()
 before:=scene()
 after:=before;after.X=4
 if err:=p.Observe(before,2,after,0,false);err!=nil{t.Fatal(err)}
 snapshot:=*p
 invalid:=[]struct{action int;reward float64}{
  {-1,0},{3,0},{2,math.NaN()},{2,math.Inf(1)},{2,101}}
 for _,bad:=range invalid {
  if err:=p.Observe(before,bad.action,after,bad.reward,false);err==nil {
   t.Fatalf("accepted illegal event %v",bad)
  }
  if !reflect.DeepEqual(*p,snapshot){t.Fatal("invalid event mutated state")}
 }
 malformed:=before;malformed.X=65
 if _,err:=p.Choose(malformed);err==nil {t.Fatal("invalid public coordinates")}
 if err:=p.Observe(before,1,malformed,1,true);err==nil{t.Fatal("invalid next source")}
 if !reflect.DeepEqual(*p,snapshot){t.Fatal("invalid next mutated state")}
}

func TestR1AuthenticatedCheckpointAndNoMutation(t *testing.T) {
 p:=New()
 b:=scene()
 if err:=p.Observe(b,2,b,1,true);err!=nil{t.Fatal(err)}
 raw,err:=p.Snapshot();if err!=nil{t.Fatal(err)}
 out,err:=Restore(raw);if err!=nil{t.Fatal(err)}
 if !reflect.DeepEqual(*p,*out){t.Fatal("checkpoint did not persist local state")}
 p0,err:=p.Scores(b);if err!=nil{t.Fatal(err)}
 p1,err:=out.Scores(b);if err!=nil||p0!=p1{t.Fatal("checkpoint changed predictions")}
 bad:=bytes.Replace(raw,[]byte("\"updates\":1"),[]byte("\"updates\":2"),1)
 if _,err:=Restore(bad);err==nil{t.Fatal("altered checkpoint accepted")}
 if _,err:=Restore(append(raw,[]byte("{}")...));err==nil{t.Fatal("trailing checkpoint accepted")}
 if len(raw)>16384{t.Fatal("oversize checkpoint")}
}

func TestR1ColdChild(t *testing.T) {
 data:=os.Getenv("ITD02_R1_CHILD")
 if data==""{return}
 raw,err:=base64.StdEncoding.DecodeString(data)
 if err!=nil{t.Fatal(err)}
 p,err:=Restore(raw)
 if err!=nil{t.Fatal(err)}
 b:=scene()
 got,err:=p.Choose(b)
 if err!=nil{t.Fatal(err)}
 wanted,err:=strconv.Atoi(os.Getenv("ITD02_R1_EXPECT"))
 if err!=nil||got!=wanted{t.Fatalf("fresh process restored action=%d expected=%d err=%v",got,wanted,err)}
}

func TestR1ActualColdOSProcessRestoration(t *testing.T) {
 if os.Getenv("ITD02_R1_CHILD")!=""{return}
 p:=New();b:=scene()
 for i:=0;i<4;i++{
  if e:=p.Observe(b,2,b,1,true);e!=nil{t.Fatal(e)}
 }
 expected,e:=p.Choose(b);if e!=nil{t.Fatal(e)}
 raw,e:=p.Snapshot();if e!=nil{t.Fatal(e)}
 cmd:=exec.Command(os.Args[0],"-test.run=^TestR1ColdChild$","-test.v")
 cmd.Env=append(os.Environ(),"ITD02_R1_CHILD="+base64.StdEncoding.EncodeToString(raw),
  "ITD02_R1_EXPECT="+strconv.Itoa(expected))
 out,e:=cmd.CombinedOutput()
 if e!=nil{t.Fatalf("cold native OS restart failed: %s; %v",string(out),e)}
}

func TestR1FutureSourceNotClaimed(t *testing.T) {
 if Schema!="wingless.itd02.r1.public-grid-action-state.v1"{t.Fatal("schema drift")}
 p:=New()
 if p.Updates!=0{t.Fatal("native source-only protocol consumed future experience")}
}
