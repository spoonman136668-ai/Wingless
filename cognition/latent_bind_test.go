package cognition
import("bytes";"testing";"math")
func TestNCI04ALearnedBindingAndComposition(t *testing.T){
 var s LatentState
 a:=[4]float32{1,0,0,0};b:=[4]float32{0,1,0,0};c:=[4]float32{0,0,1,0}
 prior,_:=s.Predict(a,b,0)
 for i:=0;i<180;i++ {if e:=s.Observe(a,b,0,1);e!=nil{t.Fatal(e)};if e:=s.Observe(b,c,1,1);e!=nil{t.Fatal(e)}}
 seen,_:=s.Predict(a,b,0);if seen<=.95||seen<=prior{t.Fatalf("learning failed %f",seen)}
 alias:=[4]float32{.98,.02,0,0}
 transferred,_:=s.Predict(alias,b,0);if math.Abs(seen-transferred)>1e-8{t.Fatalf("no alias transfer %.8f %.8f",seen,transferred)}
 path,_:=s.ForecastVia(alias,b,c,0,1);if path<=.90{t.Fatalf("no composition %f",path)}
 snap,_:=s.Snapshot();var restored LatentState
 if e:=restored.Restore(snap);e!=nil{t.Fatal(e)}
 copy,_:=restored.Snapshot();if !bytes.Equal(snap,copy){t.Fatal("non-deterministic restore")}
 if restored.Count!=3{t.Fatal("unexpected learned identity count",restored.Count)}
}
func TestNCI04AFailClosedAndNoQueryMutation(t *testing.T){
 var s LatentState
 a:=[4]float32{0,0,0,0};b:=[4]float32{1,0,0,0}
 x,_:=s.Snapshot();s.Predict(a,b,0);y,_:=s.Snapshot();if !bytes.Equal(x,y){t.Fatal("query modified memory")}
 if e:=s.Observe(a,b,2,5);e==nil{t.Fatal("accepted impossible consequence")}
 if e:=s.Observe([4]float32{100},b,0,1);e==nil{t.Fatal("unbounded sensory input")}
 z,_:=s.Snapshot();if !bytes.Equal(x,z){t.Fatal("invalid update changed state")}
 for i:=0;i<16;i++ {var f [4]float32;f[0]=float32(i)/16;f[1]=float32(i%2);f[2]=float32(i%3)/3;f[3]=float32(i%4)/4;if e:=s.Observe(f,f,0,1);e!=nil{t.Fatal(e)}}
 before,_:=s.Snapshot();newv:=[4]float32{-1,-1,-1,-1};if e:=s.Observe(newv,newv,0,1);e==nil{t.Fatal("capacity extension")};after,_:=s.Snapshot();if !bytes.Equal(before,after){t.Fatal("capacity failure partial commit")}
}
