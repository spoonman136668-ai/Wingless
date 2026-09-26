package unitary

const UP152CAblationSchema="wingless.up152c-age-state-ablation.v1"

type UP152CPoint struct{
 InitialHand int `json:"initial_hand"`
 Admission int `json:"admission"`
 ScheduledRefreshKey int `json:"scheduled_refresh_key"`
 RefreshHit bool `json:"refresh_hit"`
 HandBeforeWrite int `json:"hand_before_write"`
 FullStatePredictedSlot int `json:"full_state_predicted_slot"`
 HandOnlyPredictedSlot int `json:"hand_only_predicted_slot"`
 ActualSlot int `json:"actual_slot"`
 FullStateCorrect bool `json:"full_state_correct"`
 HandOnlyCorrect bool `json:"hand_only_correct"`
}
type UP152CResult struct{
 Schema string `json:"schema"`
 Experiment string `json:"experiment"`
 SourceUP151CSeal string `json:"source_up151c_seal"`
 ExactRecallCap int `json:"exact_recall_cap"`
 InitialHands []int `json:"initial_hands"`
 AdmissionsPerArm int `json:"admissions_per_arm"`
 SparseRefreshAdmissions []int `json:"sparse_refresh_admissions"`
 SparseRefreshKeys []int `json:"sparse_refresh_keys"`
 FullStateCorrect int `json:"full_state_correct"`
 HandOnlyCorrect int `json:"hand_only_correct"`
 TotalPredictions int `json:"total_predictions"`
 FullStateAccuracy float64 `json:"full_state_accuracy"`
 HandOnlyAccuracy float64 `json:"hand_only_accuracy"`
 FutureOracleUsed bool `json:"future_oracle_used"`
 MemoryPolicyChanged bool `json:"memory_policy_changed"`
 Points []UP152CPoint `json:"points"`
}
func up152cRefresh(adm int)(int,bool){switch adm{case 2:return 1,true;case 5:return 5,true;case 8:return 9,true;case 11:return 13,true;case 14:return 2,true};return -1,false}
func up152cRun(initialHand int)([]UP152CPoint,int,int){
 x:=&up81cAging{};for _,d:=range up145cDurable(){x.write(d.key,d.class)};x.hand=initialHand
 out:=[]UP152CPoint{};fullOK,handOK:=0,0
 for j:=1;j<=16;j++{
  rk,scheduled:=up152cRefresh(j);hit:=false;if scheduled{_,hit=x.query(rk)}
  hb:=x.hand;full:=up151cPredict(x.hand,x.age);hand:=x.hand
  newKey:=3000+j;x.write(newKey,j%3);actual:=x.find(newKey)
  fc,hc:=full==actual,hand==actual;if fc{fullOK++};if hc{handOK++}
  if !scheduled{rk=-1}
  out=append(out,UP152CPoint{InitialHand:initialHand,Admission:j,ScheduledRefreshKey:rk,RefreshHit:hit,HandBeforeWrite:hb,FullStatePredictedSlot:full,HandOnlyPredictedSlot:hand,ActualSlot:actual,FullStateCorrect:fc,HandOnlyCorrect:hc})
 }
 return out,fullOK,handOK
}
func RunUP152C()(UP152CResult,error){
 hands:=[]int{0,4,8,12};res:=UP152CResult{Schema:UP152CAblationSchema,Experiment:"UP-152C-age-state-ablation",SourceUP151CSeal:"8f59d4ee2494e5d641be5f1b2c76922b3d657033",ExactRecallCap:16,InitialHands:hands,AdmissionsPerArm:16,SparseRefreshAdmissions:[]int{2,5,8,11,14},SparseRefreshKeys:[]int{1,5,9,13,2},FutureOracleUsed:false,MemoryPolicyChanged:false}
 for _,h:=range hands{p,f,g:=up152cRun(h);res.Points=append(res.Points,p...);res.FullStateCorrect+=f;res.HandOnlyCorrect+=g}
 res.TotalPredictions=len(res.Points);res.FullStateAccuracy=float64(res.FullStateCorrect)/float64(res.TotalPredictions);res.HandOnlyAccuracy=float64(res.HandOnlyCorrect)/float64(res.TotalPredictions)
 return res,nil
}
