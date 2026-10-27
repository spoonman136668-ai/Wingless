package cognition
import("bytes";"math";"testing")
func TestOnlineLearningSnapshot(t *testing.T){
 var s State
 before,_:=s.Predict(1,0)
 for i:=0;i<200;i++{if e:=s.Observe(1,0,1,.08);e!=nil{t.Fatal(e)}}
 after,_:=s.Predict(1,0)
 if after<.9||after<=before{t.Fatalf("no online learning: %g to %g",before,after)}
 snap,_:=s.Snapshot()
 var copy State
 if e:=copy.Restore(snap);e!=nil{t.Fatal(e)}
 restored,_:=copy.Snapshot()
 if !bytes.Equal(snap,restored){t.Fatal("snapshot mismatch")}
 p,_:=s.ForecastTwo(0,1,0);if !(p>0&&p<1){t.Fatal("invalid multistep forecast")}
}
func TestFailClosedAndMarginal(t *testing.T){
 var s State;initial,_:=s.Snapshot()
 for _,bad:=range []func()error{
 func()error{return s.Observe(2,0,1,.1)},
 func()error{return s.Observe(1,0,2,.1)},
 func()error{return s.Observe(1,0,1,math.NaN())},
 func()error{return s.Restore([]byte(`{"weights":[1,2,3,4,5]}`))},
 }{if bad()==nil{t.Fatal("invalid input accepted")}}
 final,_:=s.Snapshot();if !bytes.Equal(initial,final){t.Fatal("state changed on invalid input")}
 s=State{Weights:[4]float64{.2,.4,-.8,.7}}
 p,_:=s.Predict(0,1);p0,_:=s.Predict(1,0);p1,_:=s.Predict(1,1)
 got,_:=s.ForecastTwo(1,0,1)
 if math.Abs(got-((1-p)*p0+p*p1))>1e-12{t.Fatal("oracle shortcut")}
}
