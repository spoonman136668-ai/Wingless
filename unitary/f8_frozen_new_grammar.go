package unitary

import("bytes";"compress/zlib";"crypto/sha256";"encoding/base64";"encoding/binary";"encoding/hex";_ "embed";"fmt";"io";"math";"time";"unsafe";"strings")

//go:embed f8_frozen_projection.b64
var wf8FrozenProjection string

// F7's 257 trained weights are frozen BYTES from archived F7 original run.
// No F7 training or first-use replay is performed. F7 coded upper was invalid
// (the F7 fixture did not fill Eq); F8 fixes that reference on a NEW cohort.
func wf8LoadFrozen()(wf7Encoder,error){
 var m wf7Encoder
 zipped,err:=base64.StdEncoding.DecodeString(strings.TrimSpace(wf8FrozenProjection))
 if err!=nil{return m,err}
 r,err:=zlib.NewReader(bytes.NewReader(zipped));if err!=nil{return m,err}
 raw,err:=io.ReadAll(io.LimitReader(r,257*8+1));r.Close()
 if err!=nil{return m,err}
 expected:="c7df241ca99627ea816091a84d321819828328a02a56a507cd7862d8afe67fe4"
 sum:=sha256.Sum256(raw)
 if len(raw)!=257*8||hex.EncodeToString(sum[:])!=expected{
  return m,fmt.Errorf("frozen F7 native parameters drift")
 }
 at:=0
 for k:=0;k<16;k++{for j:=0;j<16;j++{
  m.W[k][j]=math.Float64frombits(binary.LittleEndian.Uint64(raw[at:at+8]));at+=8
 }}
 m.Bias=math.Float64frombits(binary.LittleEndian.Uint64(raw[at:at+8]))
 return m,nil
}
type wf8Row struct{
 Seed int `json:"seed"`
 Family int `json:"family"`
 Truth [2]bool `json:"truth"`
 Questions [2][2]int `json:"questions"`
 Edges [12][2]int `json:"edges"`
 Predictions map[string][2]bool `json:"predictions"`
}
func wf8Topology(f int)([12][2]int,error){
 switch f{
 case 28:return [12][2]int{{0,1},{1,2},{2,3},{3,4},{0,6},{6,7},{7,8},{8,9},{10,11},{11,12},{12,13},{13,14}},nil
 case 29:return [12][2]int{{0,1},{1,2},{2,3},{3,4},{1,5},{5,6},{6,1},{2,7},{7,8},{8,2},{9,10},{10,11}},nil
 case 30:return [12][2]int{{0,1},{1,2},{2,3},{3,4},{0,5},{5,2},{1,6},{6,7},{7,3},{8,9},{9,10},{10,8}},nil
 case 31:return [12][2]int{{0,1},{1,2},{2,3},{3,4},{3,0},{0,5},{5,6},{6,2},{2,7},{7,8},{8,3},{9,10}},nil
 default:return [12][2]int{},fmt.Errorf("unknown F8 family %d",f)
 }
}
func wf8Fixture(seed,family int)(wf1Datum,error){
 motif,err:=wf8Topology(family)
 if err!=nil{return wf1Datum{},err}
 var d wf1Datum
 var edge [12][2]int
 for i,p:=range motif{edge[i]=[2]int{wf1Rename(p[0],seed,1),wf1Rename(p[1],seed,1)}}
 for i:=11;i>0;i--{
  j:=(seed*19+family*17+i*11)%(i+1)
  edge[i],edge[j]=edge[j],edge[i]
 }
 for i,p:=range edge{d.Edges[i]=nci02WEdge{From:p[0],To:p[1]}}
 a,b:=wf1Rename(0,seed,1),wf1Rename(4,seed,1)
 if !nci04WOracle(d.Edges[:],a,b,4){return d,fmt.Errorf("F8 positive not reachable")}
 used:=map[int]bool{}
 for _,e:=range d.Edges{used[e.From]=true;used[e.To]=true}
 negatives:=[]int{}
 for k:=0;k<16;k++{
  id:=wf1Rename(k,seed,1)
  if used[id]&&id!=b&&!nci04WOracle(d.Edges[:],a,id,4){negatives=append(negatives,id)}
 }
 if len(negatives)==0{return d,fmt.Errorf("F8 no present negative")}
 c:=negatives[(seed+family)%len(negatives)]
 if (seed+family)%2==0{
  d.Questions=[2][2]int{{a,b},{a,c}};d.Labels=[2]bool{true,false}
 }else{
  d.Questions=[2][2]int{{a,c},{a,b}};d.Labels=[2]bool{false,true}
 }
 for i,src:=range d.Edges{for j,dst:=range d.Edges{
  if src.To==dst.From{d.Eq[i][j]=1}
 }}
 for i,q:=range d.Questions{
  if nci04WOracle(d.Edges[:],q[0],q[1],4)!=d.Labels[i]{
   return d,fmt.Errorf("F8 ground truth changed")
  }
 }
 return d,nil
}
func RunWFoundationF8()(map[string]any,error){
 begun:=time.Now()
 frozen,err:=wf8LoadFrozen();if err!=nil{return nil,err}
 prior:=wf5Matcher()
 composed:=wf5Composer{W:[5]float64{6.2402561172942175,-.44436093640380137,
 -2.6062431834020403,8.506846585720483,-4.236991668757422}}
 fixed:=wf1Model{Form:"exact_binding",W:[3]float64{13,0,-6.5}}
 metrics:=map[string][2]int{}
 rows:=[]wf8Row{}
 for seed:=49800;seed<49848;seed++{for family:=28;family<=31;family++{
  d,err:=wf8Fixture(seed,family);if err!=nil{return nil,err}
  state:=wf7Memory(d,frozen)
  oldState:=wf4Memory(d,prior)
  exactState:=wf1Memory(d,fixed)
  row:=wf8Row{Seed:seed,Family:family,Truth:d.Labels,
   Questions:d.Questions,Predictions:map[string][2]bool{}}
  for i,e:=range d.Edges{row.Edges[i]=[2]int{e.From,e.To}}
  for _,name:=range []string{"frozen_F7_projection_and_learned_composer","frozen_F7_projection_coded_OR",
   "frozen_F5_matcher_and_learned_composer","corrected_fixed_exact_equality_upper"}{
   var predicted [2]bool
   for q,pair:=range d.Questions{
    switch name{
    case "frozen_F7_projection_and_learned_composer":predicted[q]=wf7Answer(state,pair,frozen,composed)>=.5
    case "frozen_F7_projection_coded_OR":predicted[q]=wf7CodedTraversal(state,pair,frozen)>=.5
    case "frozen_F5_matcher_and_learned_composer":predicted[q]=wf5Readout(oldState,pair,prior,composed)>=.5
    case "corrected_fixed_exact_equality_upper":predicted[q]=wf1Answer(exactState,pair)>=.5
    }
    count:=metrics[name]
    if predicted[q]==d.Labels[q]{count[0]++}
    if !d.Labels[q]&&predicted[q]{count[1]++}
    metrics[name]=count
   }
   row.Predictions[name]=predicted
  }
  rows=append(rows,row)
 }}
 if len(rows)!=192||unsafe.Sizeof(wf1State{})!=576{return nil,fmt.Errorf("F8 population/memory identity drift")}
 scores:=map[string]map[string]float64{}
 for name,v:=range metrics{
  scores[name]=map[string]float64{"correct":float64(v[0]),"accuracy":float64(v[0])/384,
   "false_positives":float64(v[1]),"negative_fp_rate":float64(v[1])/192}
 }
 verdict:="INVALID_CODED_REFERENCE_CONTROL"
 if metrics["corrected_fixed_exact_equality_upper"][0]==384{
  c:=scores["frozen_F7_projection_and_learned_composer"]
  old:=scores["frozen_F5_matcher_and_learned_composer"]
  if c["accuracy"]>=.85&&c["negative_fp_rate"]<=.15&&c["accuracy"]-old["accuracy"]>=.10{
   verdict="DEVELOPMENT_FROZEN_CONTRASTIVE_NEW_GRAMMAR_TRANSFER"
  }else{verdict="DEVELOPMENT_NO_FROZEN_CONTRASTIVE_NEW_GRAMMAR_TRANSFER"}
 }
 return map[string]any{"schema":"wingless.foundation.f8.frozen-contrastive-corrected-control-dev.v1",
  "stage":"DEVELOPMENT_ONLY","source_F7_sha":"33b45a16c3cbfd470d8ffe8b3e601f5218966985",
  "source_F7_raw_result_sha256":"799364227498074bb3ef45f9017a67d83a721b71e7e92cc5ad6b01fc8506b140",
  "source_F7_trained_weights_sha256":"c7df241ca99627ea816091a84d321819828328a02a56a507cd7862d8afe67fe4",
  "new_motif_families":[]int{28,29,30,31},
  "new_validation_seed_bounds":[]int{49800,49848},
  "question_count":384,"graph_states":192,"negative_questions":192,
  "resident_state_bytes":576,"learned_projection_parameters":257,
  "retraining":false,"original_F7_reference_invalid":true,
  "f8_exact_control_fixed":true,
  "rows":rows,"scores":scores,"registered_verdict":verdict,
  "source_independently_witnessed":false,"sealed_first_use":false,
  "accepted_ref_mutation":false,"external_llm_calls":0,
  "original_wall_ms":time.Since(begun).Milliseconds()},nil
}
