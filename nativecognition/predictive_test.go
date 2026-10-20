package nativecognition
import ("encoding/json";"fmt";"math";"os";"testing";"unsafe")
type nci04RNG struct{state uint32}
func (r *nci04RNG)next()uint32{x:=r.state;x^=x<<13;x^=x>>17;x^=x<<5;r.state=x;return x}
func (r *nci04RNG)value()float64{return float64(r.next()%10000)/10000}
func nci04Operators(seed uint32,stress bool)[2][4][4]float64{
 r:=nci04RNG{state:seed*747796405+2891336453}
 var m [2][4][4]float64
 for k:=0;k<2;k++{
  for i:=0;i<4;i++{m[k][i][i]=.65+.1*r.value()}
  for j:=0;j<4;j++{
   a:=int(r.next()%4);b:=int(r.next()%3);if b>=a{b++}
   sign:=1.0;if r.next()%2==0{sign=-1}
   m[k][a][b]+=sign*(.20+.09*r.value())
   if stress{m[k][b][a]+=sign*.08}
  }
 }
 return m
}
func nci04World(seed uint32,stress bool)[73][4]float64{
 m:=nci04Operators(seed,stress)
 r:=nci04RNG{state:seed^0x51943af}
 var x [73][4]float64
 x[0]=[4]float64{.8,-.4,.3,-.7}
 for t:=0;t<72;t++{
  context:=0;if t>=24&&t<48{context=1}
  for i:=0;i<4;i++{for j:=0;j<4;j++{x[t+1][i]+=m[context][i][j]*x[t][j]}}
  if t%9==0{j:=int(r.next()%4);x[t+1][j]+=.06*(2*r.value()-1)}
  for i:=0;i<4;i++{x[t+1][i]=math.Max(-1,math.Min(1,x[t+1][i]))}
 }
 return x
}
func nci04Seed(seed uint32,stress,erase bool)(map[string]float64,error){
 candidate,_:=NewPredictiveState(.24);control:=NewMatchedRecurrent(.24)
 world:=nci04World(seed,stress);a:=map[string]float64{}
 for t:=0;t<72;t++{
  if erase&&t==48{candidate.Erase()}
  x,y:=world[t],world[t+1]
  p,e:=candidate.Predict(x);if e!=nil{return nil,e}
  q:=control.Predict(x)
  phase:="A1";if t>=24&&t<48{phase="B"};if t>=48{phase="A2"}
  a[phase+"_candidate"]+=squareError(p,y)/24
  a[phase+"_recurrent"]+=squareError(q,y)/24
  if e:=candidate.Learn(x,y);e!=nil{return nil,e}
  control.Learn(x,y)
 }
 return a,nil
}
func TestNCI04PredictiveSourceContracts(t *testing.T){
 p,_:=NewPredictiveState(.2);x:=[4]float64{.1,-.2,.3,.4};before:=*p
 for i:=0;i<3;i++{if _,e:=p.Predict(x);e!=nil{t.Fatal(e)}}
 if *p!=before{t.Fatal("query mutated state")}
 if _,e:=p.Rollout(x,0);e==nil{t.Fatal("accepted invalid horizon")}
 if e:=p.Learn(x,[4]float64{math.NaN()});e==nil{t.Fatal("invalid observed transition")}
 if *p!=before{t.Fatal("invalid transition mutated state")}
}
func TestNCI04FrozenPredictiveABA(t *testing.T){
 result:=map[string]float64{}
 for _,c:=range []struct{name string;start uint32;stress bool}{{"train",41000,false},{"development",42000,false},{"stress",43000,true}}{
  for seed:=c.start;seed<c.start+32;seed++{
   value,err:=nci04Seed(seed,c.stress,false);if err!=nil{t.Fatal(err)}
   for k,v:=range value{result[c.name+"_"+k]+=v/32}
   if c.name=="stress"{ablated,err:=nci04Seed(seed,true,true);if err!=nil{t.Fatal(err)}
    result["stress_A2_erased_candidate"]+=ablated["A2_candidate"]/32}
  }
 }
 base:=result["stress_A2_recurrent"];candidate:=result["stress_A2_candidate"]
 gain:=0.0;if base>0{gain=(base-candidate)/base}
 memoryLoss:=result["stress_A2_erased_candidate"]-candidate
 passed:=base>0&&gain>=.05&&memoryLoss>=candidate*.10&&candidate<=result["stress_A1_candidate"]*1.10
 disposition:="NEGATIVE";if passed{disposition="PASS_DEVELOPMENT_SIGNAL_ONLY"}
 report:=map[string]any{"schema":"nci04.wingless.predictive_state.sourceonly.v1",
 "result":disposition,"metrics":result,"relative_stress_A2_gain":gain,"memory_erase_extra_error":memoryLoss,
 "candidate_parameters":32,"recurrent_control_parameters":32,
 "candidate_bytes":int(unsafe.Sizeof(PredictiveState{})),"control_bytes":int(unsafe.Sizeof(MatchedRecurrent{})),
 "task_identity_exposed":false,"hidden_causal_matrix_exposed":false,
 "intermediate_causal_labels":false,"symbolic_traversal":false,
 "entity_binding_proven":false,"cognitive_transfer_accepted":false,"first_use":0}
 raw,e:=json.Marshal(report);if e!=nil{t.Fatal(e)}
 fmt.Println("NCI04_W_RESULT="+string(raw))
 if path:=os.Getenv("NCI04_EVIDENCE");path!=""{if e:=os.WriteFile(path,append(raw,'\n'),0600);e!=nil{t.Fatal(e)}}
}
