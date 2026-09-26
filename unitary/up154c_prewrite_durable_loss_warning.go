package unitary

const UP154CLossWarningSchema="wingless.up154c-prewrite-durable-loss-warning.v1"

type UP154CPoint struct{
 InitialHand int `json:"initial_hand"`
 WriteOrdinal int `json:"write_ordinal"`
 WriteKey int `json:"write_key"`
 Rewrite bool `json:"rewrite"`
 ScheduledRefreshKey int `json:"scheduled_refresh_key"`
 RefreshHit bool `json:"refresh_hit"`
 PredictedSlot int `json:"predicted_slot"`
 ActualSlot int `json:"actual_slot"`
 Warning bool `json:"warning"`
 ActualDurableLoss bool `json:"actual_durable_loss"`
 LostOriginalKey int `json:"lost_original_key"`
 Correct bool `json:"correct"`
}
type UP154CResult struct{
 Schema string `json:"schema"`
 Experiment string `json:"experiment"`
 SourceUP153CSeal string `json:"source_up153c_seal"`
 ExactRecallCap int `json:"exact_recall_cap"`
 InitialHands []int `json:"initial_hands"`
 WritesPerArm int `json:"writes_per_arm"`
 UniqueWrites int `json:"unique_writes"`
 RewriteWrites int `json:"rewrite_writes"`
 TruePositive int `json:"true_positive"`
 FalsePositive int `json:"false_positive"`
 FalseNegative int `json:"false_negative"`
 TrueNegative int `json:"true_negative"`
 Precision float64 `json:"precision"`
 Recall float64 `json:"recall"`
 FutureOracleUsed bool `json:"future_oracle_used"`
 MemoryPolicyChanged bool `json:"memory_policy_changed"`
 InterventionTriggered bool `json:"intervention_triggered"`
 Points []UP154CPoint `json:"points"`
}
func up154cRate(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
func up154cRefresh(ord int)(int,bool){switch ord{case 3:return 2,true;case 7:return 6,true;case 11:return 10,true;case 15:return 14,true;case 19:return 1,true;case 23:return 5,true};return -1,false}
func up154cRun(initialHand int)([]UP154CPoint,int,int,int,int,int,int){
 x:=&up81cAging{};for _,d:=range up145cDurable(){x.write(d.key,d.class)};x.hand=initialHand
 out:=[]UP154CPoint{};tp,fp,fn,tn,uniq,rewrite:=0,0,0,0,0,0
 lastNewKey:=-1
 for ord:=1;ord<=24;ord++{
  rk,scheduled:=up154cRefresh(ord);hit:=false;if scheduled{_,hit=x.query(rk)}
  isRewrite:=ord%4==0
  key:=0
  if isRewrite{key=lastNewKey;rewrite++}else{key=5000+ord;lastNewKey=key;uniq++}
  predictedSlot:=-1;warning:=false
  if x.find(key)<0{
   predictedSlot=up151cPredict(x.hand,x.age)
   if x.entries[predictedSlot].used&&x.entries[predictedSlot].key>=0&&x.entries[predictedSlot].key<16{warning=true}
  }
  actualLoss:=false;lostKey:=-1
  if x.find(key)<0{
   slot:=up151cPredict(x.hand,x.age)
   if x.entries[slot].used&&x.entries[slot].key>=0&&x.entries[slot].key<16{actualLoss=true;lostKey=x.entries[slot].key}
  }
  x.write(key,ord%3);actualSlot:=x.find(key)
  ok:=warning==actualLoss
  if warning&&actualLoss{tp++}else if warning&&!actualLoss{fp++}else if !warning&&actualLoss{fn++}else{tn++}
  if !scheduled{rk=-1}
  out=append(out,UP154CPoint{InitialHand:initialHand,WriteOrdinal:ord,WriteKey:key,Rewrite:isRewrite,ScheduledRefreshKey:rk,RefreshHit:hit,PredictedSlot:predictedSlot,ActualSlot:actualSlot,Warning:warning,ActualDurableLoss:actualLoss,LostOriginalKey:lostKey,Correct:ok})
 }
 return out,tp,fp,fn,tn,uniq,rewrite
}
func RunUP154C()(UP154CResult,error){
 hands:=[]int{0,4,8,12};res:=UP154CResult{Schema:UP154CLossWarningSchema,Experiment:"UP-154C-prewrite-durable-loss-warning",SourceUP153CSeal:"da9e256704ddafc9cb20c990f347e9bb59cf2181",ExactRecallCap:16,InitialHands:hands,WritesPerArm:24,FutureOracleUsed:false,MemoryPolicyChanged:false,InterventionTriggered:false}
 for _,h:=range hands{p,tp,fp,fn,tn,u,r:=up154cRun(h);res.Points=append(res.Points,p...);res.TruePositive+=tp;res.FalsePositive+=fp;res.FalseNegative+=fn;res.TrueNegative+=tn;res.UniqueWrites+=u;res.RewriteWrites+=r}
 res.Precision=up154cRate(res.TruePositive,res.TruePositive+res.FalsePositive);res.Recall=up154cRate(res.TruePositive,res.TruePositive+res.FalseNegative)
 return res,nil
}
