package unitary
import("fmt";"sort";"strconv")
const UPLM6ASchema="wingless.up-lm6a-prepressure-length-support.v1"
type UPLM6AArm struct{Prepressure bool `json:"prepressure"`;Conditions int `json:"conditions"`;DecisionPoints int `json:"decision_points"`;RawSub16 int `json:"raw_sub16"`;RawLengthHistogram map[string]int `json:"raw_length_histogram"`}
type UPLM6AResult struct{Schema string `json:"schema"`;Experiment string `json:"experiment"`;SourceUPLM5ZSeal string `json:"source_up_lm5z_seal"`;ParentDecisionPoints int `json:"parent_decision_points"`;ParentRawSub16 int `json:"parent_raw_sub16"`;ParentEligibleSub16 int `json:"parent_eligible_sub16"`;On UPLM6AArm `json:"prepressure_on"`;Off UPLM6AArm `json:"prepressure_off"`;CanonicalHazardSelectorUsed bool `json:"canonical_hazard_selector_used"`;PolicyChanged bool `json:"policy_changed"`;FutureScheduleUsed bool `json:"future_schedule_used"`;AdaptivePolicySelectionUsed bool `json:"adaptive_policy_selection_used"`;CounterfactualOnly bool `json:"counterfactual_only"`;LiveActivation bool `json:"live_activation"`;Classification string `json:"classification"`}
func uplm6aRun(rot int,perm,profile string,budget,start,tp int,w UPLM4WWindow,bred,tred int,pre bool)(dec int,hist map[int]int){
 arms:=uplm2yPermute(uplm2xArms(rot),perm);if pre{uplm3tPrepressure(arms,rot,profile)}
 original:=map[string]bool{};for i:=0;i<12;i++{original[uplm2nName(i,rot)]=true};hist=map[int]int{};actions:=0
 for round:=0;round<6;round++{
  if round>=start{cap:=uplm4wCap(budget,round,w.BudgetCut,w.BudgetRestore,bred);eff:=uplm4wTP(tp,round,w.ThroughputCut,w.ThroughputRestore,tred);used:=map[int]bool{}
   for k:=0;k<eff && actions<cap;k++{dec++;for i:=range arms{if len(arms[i].r.order)>0{hist[len(arms[i].r.order)]++}}
    i:=uplm5uThreatArm(arms,used,rot,false);if i<0{break};if len(arms[i].r.order)==0{break};n:=arms[i].r.order[0]
    if !original[n]||arms[i].reported[n]{if pn,_,ok:=uplm2xFirstPending(&arms[i]);ok{n=pn}else{break}}
    arms[i].reported[n]=true;actions++;used[i]=true
   }
  }
  for i:=range arms{arms[i].r.write(fmt.Sprintf("6a-%t-%s-%s-%d-%d-%d-%d-%d-%d",pre,profile,perm,bred,tred,budget,start,tp,round),"x")}
 }
 return
}
func uplm6aStringHist(src map[int]int)map[string]int{out:=map[string]int{};keys:=make([]int,0,len(src));for k:=range src{keys=append(keys,k)};sort.Ints(keys);for _,k:=range keys{out[strconv.Itoa(k)]=src[k]};return out}
func RunUPLM6A()(UPLM6AResult,error){
 p,err:=RunUPLM5Z();if err!=nil{return UPLM6AResult{},err}
 profiles:=[]string{"deferred_only","layout_only","hybrid_min"};windows:=[]UPLM4WWindow{{1,3,2,4},{2,4,1,3},{1,4,2,5},{2,5,1,4},{1,3,3,5},{3,5,1,3}};breds:=[]int{1,2};treds:=[]int{1,2};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"};budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{2,3,4,5}
 r:=UPLM6AResult{Schema:UPLM6ASchema,Experiment:"UP-LM6A-prepressure-length-support",SourceUPLM5ZSeal:"2a06f1b6d6a6eff416b7e231aebecbca94e1ef76",ParentDecisionPoints:p.TotalDecisionPoints,ParentRawSub16:p.RawSub16Candidates,ParentEligibleSub16:p.EligibleSub16Candidates,CanonicalHazardSelectorUsed:true,CounterfactualOnly:true}
 onHist:=map[int]int{};offHist:=map[int]int{}
 for _,profile:=range profiles{for _,br:=range breds{for _,tr:=range treds{for _,rot:=range rots{for _,perm:=range perms{for _,w:=range windows{for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
  d,h:=uplm6aRun(rot,perm,profile,budget,start,tp,w,br,tr,true);r.On.Conditions++;r.On.DecisionPoints+=d;for n,v:=range h{onHist[n]+=v;if n<16{r.On.RawSub16+=v}}
  d,h=uplm6aRun(rot,perm,profile,budget,start,tp,w,br,tr,false);r.Off.Conditions++;r.Off.DecisionPoints+=d;for n,v:=range h{offHist[n]+=v;if n<16{r.Off.RawSub16+=v}}
 }}}}}}}}}
 r.On.Prepressure=true;r.Off.Prepressure=false;r.On.RawLengthHistogram=uplm6aStringHist(onHist);r.Off.RawLengthHistogram=uplm6aStringHist(offHist)
 if r.ParentRawSub16!=0||r.ParentEligibleSub16!=0||r.On.RawSub16!=0||r.On.DecisionPoints!=r.ParentDecisionPoints{r.Classification="ANCHOR_NOT_REPRODUCED"}else if r.Off.RawSub16>0{r.Classification="PREPRESSURE_CREATES_LENGTH16_FLOOR"}else if r.Off.RawSub16==0{r.Classification="BASE_DYNAMICS_LENGTH16_FLOOR"}else{r.Classification="OTHER_VALID_PATTERN"}
 return r,nil
}
