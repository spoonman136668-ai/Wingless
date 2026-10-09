// W-SYN01 is a separate synthetic-only falsification pilot.
// It is never an R53 external-data cohort or READY_RESEARCH authority.
package synw01

import (
 "encoding/json"
 "fmt"
 "math"
 "os"
 "testing"

 r53 "github.com/spoonman136668-ai/Wingless/research/r53prefreeze"
)

type generator struct{ x uint32 }
func newGenerator(seed uint32, regime int) generator {
 x:=seed ^ uint32(0x9e3779b9) ^ (uint32(regime+1)*uint32(0x85ebca6b))
 if x==0 {x=1}
 return generator{x:x}
}
func (g *generator) uniform() float64 {
 x:=g.x
 x^=x<<13
 x^=x>>17
 x^=x<<5
 g.x=x
 return float64(x)/4294967296.0
}
type regime struct { Name string; Coeff [6]float64 }
var regimes=[...]regime{
 {"early-only-tied",[6]float64{0.8,0.8,-0.9,0.4,1.0,-0.3}},
 {"late-only-tied",[6]float64{0.9,-0.6,0.4,0.4,0.4,0.4}},
 {"neither-tied",[6]float64{0.8,-0.6,-0.9,0.4,1.0,-0.3}},
 {"both-tied",[6]float64{0.8,0.8,0.4,0.4,0.4,0.4}},
}
var trainSeeds=[4]uint32{11,13,17,19}
var evalSeeds=[4]uint32{101,103,107,109}
func sample(g *generator, truth [6]float64)([6]float64,float64) {
 alpha:=r53.SelectedAlpha()
 var state [6]float64
 y:=0.2
 for i:=0;i<6;i++ {
  state[i]=(2*g.uniform()-1)*alpha[i]
  y+=state[i]*truth[i]
 }
 y+=0.20*(g.uniform()-0.5)
 return state,y
}
func fit(arm r53.ArmID, coeff [6]float64, seed uint32, ri int)([7]float64,error) {
 var theta [7]float64
 _, n, err:=r53.Groups(arm)
 if err!=nil {return theta,err}
 if n<2||n>7 {return theta,fmt.Errorf("invalid capacity")}
 var mat [7][8]float64
 gen:=newGenerator(seed,ri)
 for t:=0;t<436;t++ {
  x,y:=sample(&gen,coeff)
  row,e:=r53.DesignRow(arm,x)
  if e!=nil {return theta,e}
  for i:=0;i<n;i++ {
   mat[i][7]+=row[i]*y
   for j:=0;j<n;j++ {mat[i][j]+=row[i]*row[j]}
  }
 }
 for i:=1;i<n;i++ {mat[i][i]+=0.1}
 // Deterministic maximal absolute-value pivot; first row wins ties.
 for i:=0;i<n;i++ {
  pivot:=i
  for j:=i+1;j<n;j++ {
   if math.Abs(mat[j][i])>math.Abs(mat[pivot][i]) {pivot=j}
  }
  if math.Abs(mat[pivot][i])<1e-12 {return theta,fmt.Errorf("singular matrix")}
  if pivot!=i {mat[i],mat[pivot]=mat[pivot],mat[i]}
  den:=mat[i][i]
  for col:=i;col<=7;col++ {mat[i][col]/=den}
  for row:=0;row<n;row++ {
   if row==i {continue}
   mul:=mat[row][i]
   for col:=i;col<=7;col++ {mat[row][col]-=mul*mat[i][col]}
  }
 }
 for i:=0;i<n;i++ {
  theta[i]=mat[i][7]
  if math.IsInf(theta[i],0)||math.IsNaN(theta[i]) {return [7]float64{},fmt.Errorf("nonfinite theta")}
 }
 return theta,nil
}
func mse(arm r53.ArmID, theta [7]float64, truth [6]float64, seed uint32, ri int)(float64,error) {
 gen:=newGenerator(seed,ri)
 sum:=0.0
 for i:=0;i<256;i++ {
  x,y:=sample(&gen,truth)
  p,e:=r53.Predict(arm,x,theta)
  if e!=nil {return 0,e}
  d:=p-y
  sum+=d*d
 }
 return sum/256,nil
}
type cell struct {
 Regime string `json:"regime"`
 SeedPair int `json:"seed_pair"`
 MSE map[string]float64 `json:"mse"`
 EarlyGap float64 `json:"early_a2_minus_a1,omitempty"`
 LateGap float64 `json:"late_a1_minus_a2,omitempty"`
}
func evaluate() (map[string]any,error) {
 rows:=make([]cell,0,16)
 successPairs:=0
 partialEarly:=0
 partialLate:=0
 for si:=0;si<4;si++ {
  earlyOK,lateOK,earlyMatched,lateMatched:=false,false,false,false
  for ri,cfg:=range regimes {
   row:=cell{Regime:cfg.Name,SeedPair:si,MSE:map[string]float64{}}
   for _,arm:=range r53.Arms() {
    theta,e:=fit(arm,cfg.Coeff,trainSeeds[si],ri)
    if e!=nil {return nil,e}
    err,e:=mse(arm,theta,cfg.Coeff,evalSeeds[si],ri)
    if e!=nil||math.IsNaN(err)||math.IsInf(err,0) {return nil,fmt.Errorf("invalid result %s: %v",arm,e)}
    row.MSE[string(arm)]=err
   }
   if ri==0 {
    row.EarlyGap=row.MSE["A2"]-row.MSE["A1"]
    earlyOK=row.EarlyGap>=0.01
    earlyMatched=row.MSE["A1"]<=row.MSE["A0"]+0.01
   }
   if ri==1 {
    row.LateGap=row.MSE["A1"]-row.MSE["A2"]
    lateOK=row.LateGap>=0.01
    lateMatched=row.MSE["A2"]<=row.MSE["A0"]+0.01
   }
   rows=append(rows,row)
  }
  if earlyOK {partialEarly++}
  if lateOK {partialLate++}
  if earlyOK&&lateOK&&earlyMatched&&lateMatched {successPairs++}
 }
 classification:="NEGATIVE"
 if successPairs==4 {classification="SUPPORTED"
 } else if partialEarly>=2&&partialLate>=2||partialEarly==4||partialLate==4 {classification="MIXED"}
 return map[string]any{
  "schema":"ckb.synthetic-pilot-result.v1","experiment":"W-SYN01-READOUT-SLOPE-CONSTRAINT-FALSIFICATION",
  "data_class":"SYNTHETIC_ONLY","source_origin_external_claim":false,
  "official_ckb_acceptance":false,"classification":classification,
  "successful_seed_pairs":successPairs,"early_gap_seed_count":partialEarly,"late_gap_seed_count":partialLate,
  "training_rows":27904,"heldout_predictions":4096,
  "max_parameters":7,"state_dimension":6,"model_calls":0,
  "results":rows,
 },nil
}
func TestMechanicalAPIOnly(t *testing.T) {
 arms:=r53.Arms()
 if len(arms)!=4 {t.Fatal("arm drift")}
 for _,arm:=range arms {
  _,n,err:=r53.Groups(arm)
  if err!=nil||n<2||n>7 {t.Fatalf("capacity %s",arm)}
 }
 if len(regimes)!=4 {t.Fatal("regime drift")}
 // Only zero-outcome generator smoke from NON-study seed 1.
 g:=newGenerator(1,0)
 u:=g.uniform()
 if !(u>=0&&u<1) {t.Fatal("generator domain")}
}
func TestPrimarySyntheticResearch(t *testing.T) {
 // CKB must authenticate the receipt, implementation SHA and frozen
 // prereg independently BEFORE it sets this non-authoritative marker.
 if os.Getenv("CKB_SYN_W01_QUALIFIED_RUN")!="yes" {t.Skip("PRIMARY_NOT_ADMITTED_BY_CKB")}
 result,err:=evaluate()
 if err!=nil {t.Fatal(err)}
 b,err:=json.MarshalIndent(result,"","  ")
 if err!=nil {t.Fatal(err)}
 fmt.Fprintln(os.Stdout,string(b))
}
