package cognition
import("bytes";"fmt";"math";"testing";"unsafe")
func TestPartitionDynamicStateAndAba(t *testing.T){
 m:=NewPredictivePartition();a:=[4]float64{-.8,-.6,0,0};b:=[4]float64{.8,.7,0,0}
 for i:=0;i<30;i++{if e:=m.Learn(a,a,1,1);e!=nil{t.Fatal(e)}}
 p,ok,e:=m.Infer(a,1);if e!=nil||!ok||p<.9{t.Fatalf("A before distractor: %v %v %v",p,ok,e)}
 for i:=0;i<30;i++{if e:=m.Learn(b,b,1,0);e!=nil{t.Fatal(e)}}
 p,ok,e=m.Infer(a,1);if e!=nil||!ok||p<.85{t.Fatalf("A forgotten after B: %v %v %v",p,ok,e)}
 if m.used()<2||m.used()>PartitionCapacity{t.Fatalf("bounded learned state count %d",m.used())}
 snap,err:=m.Snapshot();if err!=nil{t.Fatal(err)}
 n:=NewPredictivePartition();if err:=n.Restore(snap);err!=nil{t.Fatal(err)}
 again,_:=n.Snapshot();if !bytes.Equal(snap,again){t.Fatal("snapshot drift")}
 for i:=0;i<10;i++{_,_,e:=n.Infer(b,1);if e!=nil{t.Fatal(e)}}
 again,_=n.Snapshot();if !bytes.Equal(snap,again){t.Fatal("query mutated state")}
 fmt.Printf("NCI04C_W_RETURN_PROB=%.6f LEARNED_STATES=%d BYTES=%d\n",p,m.used(),unsafe.Sizeof(*m))
}
func TestPartitionAmbiguousIdentifiabilityAndFailClosed(t *testing.T){
 m:=NewPredictivePartition();view:=[4]float64{}
 for i:=0;i<20;i++{if e:=m.Learn(view,view,0,1);e!=nil{t.Fatal(e)}}
 for i:=0;i<20;i++{if e:=m.Learn(view,view,0,0);e!=nil{t.Fatal(e)}}
 _,ok,e:=m.Infer(view,0);if e!=nil{t.Fatal(e)}
 if ok&&m.used()>1{t.Fatal("invented identity from indistinguishable view")}
 before,_:=m.Snapshot()
 if e:=m.Learn(view,[4]float64{math.NaN()},0,1);e==nil{t.Fatal("accepted NaN")}
 if e:=m.Learn(view,view,0,2);e==nil{t.Fatal("invalid outcome")}
 after,_:=m.Snapshot();if !bytes.Equal(before,after){t.Fatal("invalid input mutated state")}
 if e:=m.Restore([]byte("{\"schema\":2}"));e==nil{t.Fatal("accepted source drift")}
}
type trainedEMA struct{centres [8][4]float64;effect [8]float64;used int}
func(m *trainedEMA)Learn(x [4]float64,y float64){
 winner:=-1;dist:=math.Inf(1)
 for i:=0;i<m.used;i++{d:=0.0;for k:=0;k<4;k++{v:=x[k]-m.centres[i][k];d+=v*v};if d<dist{winner,dist=i,d}}
 if winner<0||dist>.04{if m.used<8{winner=m.used;m.used++}else{winner=m.used-1};m.centres[winner]=x;m.effect[winner]=.5}
 m.effect[winner]=.85*m.effect[winner]+.15*y
}
func(m *trainedEMA)Predict(x [4]float64)float64{
 winner:=-1;dist:=math.Inf(1)
 for i:=0;i<m.used;i++{d:=0.0;for k:=0;k<4;k++{v:=x[k]-m.centres[i][k];d+=v*v};if d<dist{winner,dist=i,d}}
 if winner<0{return .5};return m.effect[winner]
}
func TestPartitionTrainedPeerWithoutFalseAcceptance(t *testing.T){
 n:=NewPredictivePartition();var peer trainedEMA
 correct,baseline:=0,0
 for i:=0;i<96;i++{
  id:=i%4;x:=[4]float64{float64(id-2)*.4,.1,0,0};y:=float64(id%2)
  if i>=48{
   p,ok,e:=n.Infer(x,1);if e!=nil{t.Fatal(e)}
   if ok&&(p>=.5)==(y==1){correct++}
   q:=peer.Predict(x);if (q>=.5)==(y==1){baseline++}
  }
  if e:=n.Learn(x,x,1,y);e!=nil{t.Fatal(e)};peer.Learn(x,y)
 }
 fmt.Printf("NCI04C_W_DEV_NATIVE=%d/48 TRAINED_CONTROL=%d/48 NATIVE_BYTES=%d PEER_BYTES=%d\n",
 correct,baseline,unsafe.Sizeof(*n),unsafe.Sizeof(peer))
 if correct>48||baseline>48{t.Fatal("invalid measurement")}
 // This test MUST pass when the candidate loses: retain a negative.
}
