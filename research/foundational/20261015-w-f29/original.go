// W-F29: independently generated multi-encounter, role-conditioned identity development.
// This is a DEVELOPMENT experiment, not evidence-born F3 or a production model.
package main

import (
 "encoding/json"
 "fmt"
 "math"
 "math/rand"
 "os"
 "time"
 "github.com/spoonman136668-ai/Wingless/research/foundational/20261015-w-f29/custody"
)
type row struct {
 Seed,Family,Depth,RoleShift,Condition int
 Truth bool
 Evidence []custody.Evidence
 Predictions map[string]bool
 HistoryBytes,Comparisons int
}
type tally struct {Correct,FP,FN,Total int; Accuracy float64}
type result struct {
 Schema,Stage,Parent,Verdict string
 LearnedWeights [10]float64
 TrainExamples,Questions,GraphStates,NativeBytes,ExternalHistoryBytes,TotalWitnessComparisons int
 TrainSeedRange,TestSeedRange [2]int
 TrainFamilies,TestFamilies []int
 Conditions []string
 Scores map[string]tally
 FamilyScores map[int]map[string]tally
 Rows []row
 SealedFirstUse,IndependentlyBornSource,AcceptedRefMutation bool
}
func bounded(v float64)float64{if v<0{return 0};if v>1{return 1};return v}
func sig(z float64)float64{if z>30{return 1};if z< -30{return 0};return 1/(1+math.Exp(-z))}
func generate(seed,family,depth,shift,cond int,truth bool)custody.PrivateCase{
 // True equivalence is determined by a private join in a variable-length
 // role path. Neither join identity nor the query answer is given to scorer.
 salt:=int64(seed)*991+int64(family)*10007+int64(depth)*131+
   int64(shift)*733+int64(cond)*1999
 if !truth{salt+=874319}
 rng:=rand.New(rand.NewSource(salt))
 var observations []custody.Evidence
 for encounter:=0;encounter<3;encounter++{
  e:=custody.Evidence{EncounterID:uint64(100+encounter),Role:uint8((encounter+shift+family+depth)%4)}
  for source:=0;source<3;source++{
   base:=rng.Float64()*1.6-.8
   // A relational witness is an episode-local measurement, not an ID hash.
   // The witness's role changes across independent encounters.
   e.Left[source]=base
   match:=true
   if truth{
    // A genuine match suffers an independently moved, corrupt observation.
    match = source != (encounter+cond+shift)%3
    if cond==2 && encounter==2 && source==((encounter+cond+shift+1)%3) {match=false}
   }else{
    // A coincidental impostor agrees in two channels in some encounters.
    // Its agreement rotates across encounters, rather than persisting.
    match = source==(encounter+shift+cond)%3
    if cond==1 || cond==3 { match=source!=(encounter+shift+cond)%3 }
    if cond==2 && encounter==0 {match=source!=((encounter+shift+cond)%3)}
   }
   if match {e.Right[source]=base+(rng.Float64()-.5)*.055
   }else{e.Right[source]=base+.16+(rng.Float64())*.18}
  }
  observations=append(observations,e)
 }
 return custody.PrivateCase{Input:custody.PublicInput{Evidence:observations},Truth:truth}
}
// All features use public purchased observation values only. 10 floating
// coefficients are trained; feature construction is acknowledged programmed.
func features(input custody.PublicInput)([10]float64,int){
 var hits [3][3]float64
 count:=0
 for i,e:=range input.Evidence{
  for j:=0;j<3;j++{
   if math.Abs(e.Left[j]-e.Right[j])<custody.ConsensusThreshold {hits[i][j]=1}
   count++
  }
 }
 var f [10]float64
 f[0]=1
 across:=0.
 for j:=0;j<3;j++{
  sum:=0.;persist:=1.
  for i:=0;i<3;i++{sum+=hits[i][j];persist*=hits[i][j]}
  f[1+j]=sum/3
  across+=persist
 }
 f[4]=across/3
 majority:=0.; minimum:=3.; maximum:=0.
 for i:=0;i<3;i++{
  n:=hits[i][0]+hits[i][1]+hits[i][2]
  if n>=2{majority++}
  if n<minimum{minimum=n}
  if n>maximum{maximum=n}
 }
 f[5]=majority/3
 f[6]=minimum/3
 f[7]=maximum/3
 f[8]=(maximum-minimum)/3
 // Role/encounter interaction cannot directly encode label or test entity.
 f[9]=float64((int(input.Evidence[0].Role)+int(input.Evidence[1].Role))%4)/3
 return f,count
}
func train()[10]float64{
 var x [][10]float64
 var labels []float64
 for seed:=151000;seed<151100;seed++{for family:=112;family<116;family++{
  for depth:=9;depth<=11;depth++{for cond:=0;cond<4;cond++{
   for k:=0;k<2;k++{
    input:=generate(seed,family,depth,seed%4,cond,k==0)
    feat,_:=features(input.Input)
    x=append(x,feat)
    if k==0{labels=append(labels,1)}else{labels=append(labels,0)}
   }
  }}
 }}
 var w [10]float64
 // Full-batch logistic fit on TRAIN only, deterministic seed-locked cohort.
 for epoch:=0;epoch<130;epoch++{
  var gradient [10]float64
  for i,v:=range x{
   z:=0.;for j:=range w{z+=w[j]*v[j]}
   delta:=sig(z)-labels[i]
   for j:=range w{gradient[j]+=delta*v[j]}
  }
  for j:=range w{w[j]-=.33*(gradient[j]/float64(len(x))+.0003*w[j])}
 }
 return w
}
func predict(input custody.PublicInput,w [10]float64,mode string)(bool,int){
 feat,comp:=features(input)
 switch mode {
 case "learned":
  z:=0.;for j:=range w{z+=w[j]*feat[j]};return sig(z)>=.5,comp
 case "multiview_fixed":
  // Two distinct channels must agree across ALL encounters.
  return feat[4]>=.6666666,comp
 case "static_first":
  return feat[5]>=.3333333,comp
 case "fixed_persistent":
  return feat[4]>=.3333333,comp
 case "consensus":
  m:=custody.Meter{}
  a,e:=custody.Consensus(input,&m);if e!=nil{panic(e)};return a,m.WitnessComparisons
 case "no_link":
  return false,0
 default:panic("unregistered model")
 }
}
func audit(rows []row,arm string)tally{
 var t tally
 for _,r:=range rows{
  pred:=r.Predictions[arm];t.Total++
  if pred==r.Truth{t.Correct++}
  if pred&&!r.Truth{t.FP++}
  if !pred&&r.Truth{t.FN++}
 }
 t.Accuracy=float64(t.Correct)/float64(t.Total)
 return t
}
func original()(result,error){
 start:=time.Now();_ = start
 trainManifest:=custody.Manifest{GraphFamilies:[]uint64{112,113,114,115}}
 testManifest:=custody.Manifest{GraphFamilies:[]uint64{116,117,118,119}}
 for i:=151000;i<151100;i++{trainManifest.Seeds=append(trainManifest.Seeds,uint64(i));trainManifest.EntityIDs=append(trainManifest.EntityIDs,uint64(500000+i))}
 for i:=152300;i<152312;i++{testManifest.Seeds=append(testManifest.Seeds,uint64(i));testManifest.EntityIDs=append(testManifest.EntityIDs,uint64(700000+i))}
 if e:=custody.ValidateSplit(trainManifest,testManifest);e!=nil{return result{},e}
 weights:=train()
 modes:=[]string{"learned","consensus","multiview_fixed","static_first","fixed_persistent","no_link","privileged_oracle"}
 out:=result{Schema:"wingless.f29.cross-encounter-identity.development.v1",Stage:"DEVELOPMENT_ONLY",
  Parent:"331faa14ffe039fffdaac8c4ee6640244d15d149",
  LearnedWeights:weights,TrainExamples:9600,NativeBytes:576,
  TrainSeedRange:[2]int{151000,151100},TestSeedRange:[2]int{152300,152312},
  TrainFamilies:[]int{112,113,114,115},TestFamilies:[]int{116,117,118,119},
  Conditions:[]string{"unrelated","correlated-impostor","witness-disagreement","counterfactual-context"},
  Scores:map[string]tally{}, FamilyScores:map[int]map[string]tally{}}
 for seed:=152300;seed<152312;seed++{for family:=116;family<120;family++{
  for depth:=9;depth<=11;depth++{for cond:=0;cond<4;cond++{
   for k:=0;k<2;k++{
    source:=generate(seed,family,depth,seed%4,cond,k==0)
    public,meter,err:=custody.Prepare(source,576);if err!=nil{return result{},err}
    predictions:=map[string]bool{"privileged_oracle":source.Truth}
    comps:=0
    for _,mode:=range modes[:len(modes)-1]{
     ans,c:=predict(public,weights,mode);predictions[mode]=ans;comps+=c
    }
    out.Rows=append(out.Rows,row{Seed:seed,Family:family,Depth:depth,RoleShift:seed%4,
      Condition:cond,Truth:source.Truth,Evidence:public.Evidence,Predictions:predictions,
      HistoryBytes:meter.HistoryBytes,Comparisons:comps})
    out.ExternalHistoryBytes+=meter.HistoryBytes;out.TotalWitnessComparisons+=comps
   }
  }}
 }}
 out.Questions=len(out.Rows);out.GraphStates=12*4*3
 if out.Questions!=1152{return result{},fmt.Errorf("unexpected case count %d",out.Questions)}
 for _,mode:=range modes{
  out.Scores[mode]=audit(out.Rows,mode)
  for family:=116;family<120;family++{
   subset:=[]row{};for _,r:=range out.Rows{if r.Family==family{subset=append(subset,r)}}
   if out.FamilyScores[family]==nil{out.FamilyScores[family]=map[string]tally{}}
   out.FamilyScores[family][mode]=audit(subset,mode)
  }
 }
 if out.Scores["privileged_oracle"].Correct!=out.Questions{return result{},fmt.Errorf("invalid source oracle")}
 fairBest:=0.;fairFP:=out.Questions
 for _,mode:=range modes[1:len(modes)-1]{
  s:=out.Scores[mode];if s.Accuracy>fairBest{fairBest=s.Accuracy;fairFP=s.FP}
 }
 learned:=out.Scores["learned"]
 advance:=learned.Accuracy-fairBest>=.03 && learned.FP<fairFP && learned.FN<=out.Scores["consensus"].FN+5
 for f:=116;f<120;f++{
  best:=0.;for _,mode:=range modes[1:len(modes)-1]{best=math.Max(best,out.FamilyScores[f][mode].Accuracy)}
  if out.FamilyScores[f]["learned"].Accuracy<best{advance=false}
 }
 out.Verdict="DEVELOPMENT_CROSS_ENCOUNTER_NO_QUALIFIED_GAIN"
 if advance{out.Verdict="DEVELOPMENT_CROSS_ENCOUNTER_LEARNED_GAIN"}
 return out,nil
}
func main(){
 r,e:=original();if e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)}
 enc:=json.NewEncoder(os.Stdout);enc.SetIndent("","  ");if e:=enc.Encode(r);e!=nil{panic(e)}
}
