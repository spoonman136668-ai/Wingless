package unitary

const UP153CHandAblationSchema="wingless.up153c-hand-state-ablation.v1"

type UP153CPoint struct{
 InitialHand int `json:"initial_hand"`
 Admission int `json:"admission"`
 ScheduledRefreshKey int `json:"scheduled_refresh_key"`
 RefreshHit bool `json:"refresh_hit"`
 HandBeforeWrite int `json:"hand_before_write"`
 FullStatePredictedSlot int `json:"full_state_predicted_slot"`
 AgeOnlyPredictedSlot int `json:"age_only_predicted_slot"`
 ActualSlot int `json:"actual_slot"`
 FullStateCorrect bool `json:"full_state_correct"`
 AgeOnlyCorrect bool `json:"age_only_correct"`
}
type UP153CResult struct{
 Schema string `json:"schema"`
 Experiment string `json:"experiment"`
 SourceUP152CSeal string `json:"source_up152c_seal"`
 ExactRecallCap int `json:"exact_recall_cap"`
 InitialHands []int `json:"initial_hands"`
 AdmissionsPerArm int `json:"admissions_per_arm"`
 SparseRefreshAdmissions []int `json:"sparse_refresh_admissions"`
 SparseRefreshKeys []int `json:"sparse_refresh_keys"`
 FullStateCorrect int `json:"full_state_correct"`
 AgeOnlyCorrect int `json:"age_only_correct"`
 TotalPredictions int `json:"total_predictions"`
 FullStateAccuracy float64 `json:"full_state_accuracy"`
 AgeOnlyAccuracy float64 `json:"age_only_accuracy"`
 FutureOracleUsed bool `json:"future_oracle_used"`
 MemoryPolicyChanged bool `json:"memory_policy_changed"`
 Points []UP153CPoint `json:"points"`
}
func up153cRefresh(adm int)(int,bool){switch adm{case 3:return 4,true;case 7:return 8,true;case 10:return 12,true;case 13:return 0,true;case 16:return 6,true};return -1,false}
func up153cAgeOnly(age [16]uint8)int{
 h:=0
 for{
  if age[h]==0{return h}
  age[h]--
  h=(h+1)%16
 }
}
func up153cRun(initialHand int)([]UP153CPoint,int,int){
 x:=&up81cAging{};for _,d:=range up145cDurable(){x.write(d.key,d.class)};x.hand=initialHand
 out:=[]UP153CPoint{};fullOK,ageOK:=0,0
 for j:=1;j<=16;j++{
  rk,scheduled:=up153cRefresh(j);hit:=false;if scheduled{_,hit=x.query(rk)}
  hb:=x.hand;full:=up151cPredict(x.hand,x.age);ageOnly:=up153cAgeOnly(x.age)
  newKey:=4000+j;x.write(newKey,j%3);actual:=x.find(newKey)
  fc,ac:=full==actual,ageOnly==actual;if fc{fullOK++};if ac{ageOK++}
  if !scheduled{rk=-1}
  out=append(out,UP153CPoint{InitialHand:initialHand,Admission:j,ScheduledRefreshKey:rk,RefreshHit:hit,HandBeforeWrite:hb,FullStatePredictedSlot:full,AgeOnlyPredictedSlot:ageOnly,ActualSlot:actual,FullStateCorrect:fc,AgeOnlyCorrect:ac})
 }
 return out,fullOK,ageOK
}
func RunUP153C()(UP153CResult,error){
 hands:=[]int{0,4,8,12};res:=UP153CResult{Schema:UP153CHandAblationSchema,Experiment:"UP-153C-hand-state-ablation",SourceUP152CSeal:"45772d81c0ea096e21fa0f1c33b7cfe43cc60051",ExactRecallCap:16,InitialHands:hands,AdmissionsPerArm:16,SparseRefreshAdmissions:[]int{3,7,10,13,16},SparseRefreshKeys:[]int{4,8,12,0,6},FutureOracleUsed:false,MemoryPolicyChanged:false}
 for _,h:=range hands{p,f,a:=up153cRun(h);res.Points=append(res.Points,p...);res.FullStateCorrect+=f;res.AgeOnlyCorrect+=a}
 res.TotalPredictions=len(res.Points);res.FullStateAccuracy=float64(res.FullStateCorrect)/float64(res.TotalPredictions);res.AgeOnlyAccuracy=float64(res.AgeOnlyCorrect)/float64(res.TotalPredictions)
 return res,nil
}
