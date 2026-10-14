package unitary

import ("fmt";"math";"time";"unsafe")

// W-F7: learned 16D entity representation from self-supervised identity pairs.
// The endpoint/edge comparisons at inference are learned projection distances,
// never a coded ID equality signal. Four-hop composition is either frozen W-F5
// learned recurrence or a separately identified coded reference control.
// TRAIN and VALID graph grammars (12..15 / 16..19) are wholly disjoint from
// W-F1 through W-F6 motifs. All results are DEVELOPMENT ONLY, never F3 sealed.
type wf7Encoder struct{ W [16][16]float64; Bias float64 }
type wf7Row struct {
 Seed int `json:"seed"`
 Family int `json:"family"`
 Truth [2]bool `json:"truth"`
 Predictions map[string][2]bool `json:"predictions"`
}
func wf7Init()wf7Encoder{
 m:=wf7Encoder{}
 for k:=0;k<16;k++{for i:=0;i<16;i++{
  m.W[k][i]=.12*math.Sin(float64((k+2)*(i+5)*17))
 }}
 return m
}
func wf7Pair(m wf7Encoder,a,b int)(float64,[16]float64,[16]float64){
 var diff [16]float64
 mask:=(a^b)&65535
 for i:=0;i<16;i++{if mask&(1<<i)!=0{
  if a&(1<<i)!=0{diff[i]=1}else{diff[i]=-1}
 }}
 dist:=0.;var projected [16]float64
 for k:=0;k<16;k++{
  for i:=0;i<16;i++{projected[k]+=m.W[k][i]*diff[i]}
  dist+=projected[k]*projected[k]
 }
 return wf1Prob(m.Bias-dist),diff,projected
}
func wf7Train()(wf7Encoder,int64,float64){
 m:=wf7Init()
 var momentW,varianceW [16][16]float64
 var momentBias,varianceBias float64
 start:=time.Now()
 t:=uint32(946711)
 sample:=func()uint32{
  t^=t<<13;t^=t>>17;t^=t<<5
  return t
 }
 loss:=0.
 for step:=1;step<=135;step++{
  var gradW [16][16]float64
  gradBias:=0.;loss=0.
  for index:=0;index<256;index++{
   a:=int(sample()&65535)
   b:=a
   positive:=index%2==0
   if !positive{
    bit0:=int(sample()%16)
    b^=1<<bit0
    if index%6==1{b^=1<<int((bit0+1+int(sample()%15))%16)}
   }
   p,diff,projected:=wf7Pair(m,a,b)
   p=max(.000001,min(.999999,p))
   label:=0.
   if positive{label=1;loss-=math.Log(p)}else{loss-=math.Log(1-p)}
   delta:=p-label
   gradBias+=delta
   for k:=0;k<16;k++{for i:=0;i<16;i++{
    gradW[k][i]+=-2*delta*projected[k]*diff[i]
   }}
  }
  loss/=256;gradBias/=256
  momentBias=.9*momentBias+.1*gradBias
  varianceBias=.999*varianceBias+.001*gradBias*gradBias
  beta1:=1-math.Pow(.9,float64(step))
  beta2:=1-math.Pow(.999,float64(step))
  m.Bias-=.05*(momentBias/beta1)/(math.Sqrt(varianceBias/beta2)+1e-8)
  m.Bias=max(-20,min(20,m.Bias))
  for k:=0;k<16;k++{for i:=0;i<16;i++{
   g:=max(-10,min(10,gradW[k][i]/256))
   momentW[k][i]=.9*momentW[k][i]+.1*g
   varianceW[k][i]=.999*varianceW[k][i]+.001*g*g
   m.W[k][i]-=.05*(momentW[k][i]/beta1)/(math.Sqrt(varianceW[k][i]/beta2)+1e-8)
  }}
 }
 return m,time.Since(start).Milliseconds(),loss
}
func wf7Topology(f int)([12][2]int,error){
 switch f{
 case 12:return [12][2]int{{0,1},{1,2},{2,3},{3,4},{0,5},{5,6},{6,7},{7,8},{9,10},{10,11},{11,12},{12,13}},nil
 case 13:return [12][2]int{{0,1},{1,2},{2,3},{3,4},{1,5},{5,6},{6,3},{2,7},{7,8},{8,9},{10,11},{11,12}},nil
 case 14:return [12][2]int{{0,1},{1,2},{2,3},{3,4},{1,0},{2,1},{3,2},{4,3},{5,6},{6,7},{7,8},{8,9}},nil
 case 15:return [12][2]int{{0,1},{1,2},{2,3},{3,4},{0,5},{5,2},{2,6},{6,7},{7,3},{8,9},{9,10},{10,11}},nil
 case 16:return [12][2]int{{0,1},{1,2},{2,3},{3,4},{0,6},{6,2},{1,7},{7,8},{8,3},{9,10},{10,11},{11,12}},nil
 case 17:return [12][2]int{{0,1},{1,2},{2,3},{3,4},{1,5},{5,1},{2,6},{6,2},{7,8},{8,9},{9,7},{3,10}},nil
 case 18:return [12][2]int{{0,1},{1,2},{2,3},{3,4},{0,5},{5,6},{6,1},{2,7},{7,8},{8,4},{9,10},{10,9}},nil
 case 19:return [12][2]int{{0,1},{1,2},{2,3},{3,4},{1,6},{6,7},{7,8},{8,2},{3,9},{9,10},{10,11},{11,0}},nil
 default:return [12][2]int{},fmt.Errorf("not a preregistered topology family: %d",f)
 }
}
func wf7Fixture(seed,family int)(wf1Datum,error){
 motif,err:=wf7Topology(family);if err!=nil{return wf1Datum{},err}
 var d wf1Datum
 var edges [12][2]int
 for i,p:=range motif{
  edges[i]=[2]int{wf1Rename(p[0],seed,2),wf1Rename(p[1],seed,2)}
 }
 for i:=11;i>0;i--{j:=(seed*13+i*17+family*19)%(i+1);edges[i],edges[j]=edges[j],edges[i]}
 for i,p:=range edges{d.Edges[i]=nci02WEdge{From:p[0],To:p[1]}}
 src,end:=wf1Rename(0,seed,2),wf1Rename(4,seed,2)
 if !nci04WOracle(d.Edges[:],src,end,4){return d,fmt.Errorf("positive oracle changed")}
 used:=map[int]bool{}
 for _,e:=range d.Edges{used[e.From]=true;used[e.To]=true}
 negatives:=[]int{}
 for id:=0;id<=15;id++{
  candidate:=wf1Rename(id,seed,2)
  if candidate!=end&&used[candidate]&&!nci04WOracle(d.Edges[:],src,candidate,4){
   negatives=append(negatives,candidate)
  }
 }
 if len(negatives)==0{return d,fmt.Errorf("no observed hard negatives")}
 target:=negatives[(seed+family)%len(negatives)]
 if (seed+family)%2==0{
  d.Questions=[2][2]int{{src,end},{src,target}}
  d.Labels=[2]bool{true,false}
 }else{
  d.Questions=[2][2]int{{src,target},{src,end}}
  d.Labels=[2]bool{false,true}
 }
 return d,nil
}
func wf7Memory(d wf1Datum,m wf7Encoder)wf1State{
 var state wf1State
 for i,e:=range d.Edges{
  state[i]=math.Float32frombits(uint32(e.From)<<16|uint32(e.To))
 }
 off:=12
 for i,a:=range d.Edges{for j,b:=range d.Edges{
  if i==j{continue}
  prob,_,_:=wf7Pair(m,a.To,b.From)
  state[off]=float32(prob);off++
 }}
 return state
}
func wf7Answer(mem wf1State,q [2]int,m wf7Encoder,compose wf5Composer)float64{
 var src,dst [12]int
 for i:=0;i<12;i++{
  bits:=math.Float32bits(mem[i])
  src[i]=int(bits>>16);dst[i]=int(bits&65535)
 }
 var active [12]float64
 for i:=0;i<12;i++{active[i],_,_=wf7Pair(m,q[0],src[i])}
 for round:=0;round<3;round++{
  var next [12]float64
  for j:=0;j<12;j++{
   sum:=0.
   for i:=0;i<12;i++{
    if i==j{continue}
    off:=12+i*11+j;if j>i{off--}
    sum+=active[i]*float64(mem[off])
   }
   next[j]=wf1Prob(compose.W[0]*sum+compose.W[1]*active[j]+compose.W[2])
  }
  active=next
 }
 combined:=0.
 for i:=0;i<12;i++{
  score,_,_:=wf7Pair(m,dst[i],q[1])
  combined+=active[i]*score
 }
 return wf1Prob(compose.W[3]*combined+compose.W[4])
}
func wf7CodedTraversal(mem wf1State,q [2]int,m wf7Encoder)float64{
 var src,dst [12]int
 for i:=0;i<12;i++{bits:=math.Float32bits(mem[i]);src[i]=int(bits>>16);dst[i]=int(bits&65535)}
 var act [12]float64
 for i:=0;i<12;i++{act[i],_,_=wf7Pair(m,q[0],src[i])}
 for h:=0;h<3;h++{
  var next [12]float64
  for j:=0;j<12;j++{
   missed:=1.
   for i:=0;i<12;i++{
    if i==j{continue}
    off:=12+i*11+j;if j>i{off--}
    missed*=1-act[i]*float64(mem[off])
   }
   next[j]=1-missed
  }
  act=next
 }
 miss:=1.
 for i:=0;i<12;i++{score,_,_:=wf7Pair(m,dst[i],q[1]);miss*=1-act[i]*score}
 return 1-miss
}
func RunWFoundationF7()(map[string]any,error){
 start:=time.Now()
 encoder,trainMs,trainingLoss:=wf7Train()
 frozenMatcher:=wf5Matcher()
 frozenComposer:=wf5Composer{W:[5]float64{6.2402561172942175,-.44436093640380137,
  -2.6062431834020403,8.506846585720483,-4.236991668757422}}
 coded:=wf1Model{Form:"exact_binding",W:[3]float64{13,0,-6.5}}
 scores:=map[string][2]int{}
 rows:=[]wf7Row{}
 for seed:=47200;seed<47212;seed++{for family:=16;family<=19;family++{
  d,err:=wf7Fixture(seed,family);if err!=nil{return nil,err}
  state:=wf7Memory(d,encoder)
  frozenState:=wf4Memory(d,frozenMatcher)
  oracleState:=wf1Memory(d,coded)
  row:=wf7Row{Seed:seed,Family:family,Truth:d.Labels,Predictions:map[string][2]bool{}}
  for _,name:=range []string{"contrastive_projection_learned_recurrence","contrastive_projection_coded_traversal",
   "frozen_F5_learned_recurrence","fixed_coded_exact_upper"}{
   var predicted [2]bool
   for i,q:=range d.Questions{
    switch name{
    case "contrastive_projection_learned_recurrence":predicted[i]=wf7Answer(state,q,encoder,frozenComposer)>=.5
    case "contrastive_projection_coded_traversal":predicted[i]=wf7CodedTraversal(state,q,encoder)>=.5
    case "frozen_F5_learned_recurrence":predicted[i]=wf5Readout(frozenState,q,frozenMatcher,frozenComposer)>=.5
    case "fixed_coded_exact_upper":predicted[i]=wf1Answer(oracleState,q)>=.5
    }
    pair:=scores[name]
    if predicted[i]==d.Labels[i]{pair[0]++}
    if !d.Labels[i]&&predicted[i]{pair[1]++}
    scores[name]=pair
   }
   row.Predictions[name]=predicted
  }
  rows=append(rows,row)
 }}
 if len(rows)!=48||unsafe.Sizeof(wf1State{})!=576{
  return nil,fmt.Errorf("unfrozen graph population or resident memory")
 }
 aggregate:=map[string]map[string]float64{}
 for name,s:=range scores{
  aggregate[name]=map[string]float64{"accuracy":float64(s[0])/96,
   "correct":float64(s[0]),"false_positives":float64(s[1]),
   "negative_fp_rate":float64(s[1])/48}
 }
 candidate:=aggregate["contrastive_projection_learned_recurrence"]
 frozen:=aggregate["frozen_F5_learned_recurrence"]
 upper:=aggregate["fixed_coded_exact_upper"]
 pass:=candidate["accuracy"]>=.85&&candidate["negative_fp_rate"]<=.15&&
   candidate["accuracy"]-frozen["accuracy"]>=.10&&
   upper["accuracy"]-candidate["accuracy"]<=.15
 verdict:="DEVELOPMENT_NO_CONTRASTIVE_IDENTITY_STRUCTURAL_TRANSFER"
 if pass{verdict="DEVELOPMENT_CONTRASTIVE_IDENTITY_STRUCTURAL_TRANSFER_OBSERVED"}
 return map[string]any{"schema":"wingless.foundation.f7.contrastive-identity-structural-transfer-development.v1",
  "stage":"DEVELOPMENT_ONLY","source_F6_sha":"70f7ccd375a4b2d1bcbca4ae5af191e6de5ea437",
  "training_data":"deterministic 256 identity pairs per optimizer step, 135 steps, no graph question labels",
  "new_topology_train_only_families":[]int{12,13,14,15},
  "new_topology_validation_families":[]int{16,17,18,19},
  "validation_seed_bounds":[]int{47200,47212},"graph_states":len(rows),
  "question_count":96,"negative_questions":48,"resident_state_bytes":576,
  "learned_projection_dimensions":16,"learned_projection_parameter_count":257,
  "identity_pair_training_steps":135,"identity_training_loss":trainingLoss,
  "identity_training_ms":trainMs,"projection":encoder.W,"projection_bias":encoder.Bias,
  "frozen_composer":frozenComposer.W,"scores":aggregate,"rows":rows,
  "registered_verdict":verdict,"direct_coded_identity_in_candidate":false,
  "coded_or_traversal_in_candidate":false,"engineered_bit_inputs":true,
  "four_hop_computation_schedule_fixed":true,
  "development_source_by_candidate_team":true,"independently_witnessed_source_birth":false,
  "sealed_first_use":false,"accepted_ref_mutation":false,"external_llm_calls":0,
  "wall_ms":time.Since(start).Milliseconds()},nil
}
