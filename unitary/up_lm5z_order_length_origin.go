package unitary

import (
	"fmt"
	"sort"
	"strconv"
)

const UPLM5ZOriginSchema = "wingless.up-lm5z-order-length-origin.v1"

type UPLM5ZCounts struct {
	UsedThisDecision   int `json:"used_this_decision"`
	Empty              int `json:"empty"`
	NonOriginalVictim  int `json:"non_original_victim"`
	AlreadyReported    int `json:"already_reported"`
	Eligible           int `json:"eligible"`
}

type UPLM5ZCellSummary struct {
	DeadlineProfile    string                    `json:"deadline_profile"`
	BudgetReduction    int                       `json:"budget_reduction"`
	ThroughputReduction int                      `json:"throughput_reduction"`
	Rotation           int                       `json:"rotation"`
	Permutation        string                    `json:"permutation"`
	Conditions         int                       `json:"conditions"`
	DecisionPoints     int                       `json:"decision_points"`
	Counts             UPLM5ZCounts              `json:"counts"`
	RawLengthHistogram map[string]int             `json:"raw_length_histogram"`
	CategoryHistograms map[string]map[string]int  `json:"category_histograms"`
}

type UPLM5ZResult struct {
	Schema                       string                    `json:"schema"`
	Experiment                   string                    `json:"experiment"`
	SourceUPLM5YSeal             string                    `json:"source_up_lm5y_seal"`
	ParentDecisionPoints         int                       `json:"parent_decision_points"`
	ParentEligibleCandidates     int                       `json:"parent_eligible_candidates"`
	ParentObservedLengths        []int                     `json:"parent_observed_lengths"`
	ParentSubhazardCandidates    int                       `json:"parent_subhazard_candidates"`
	TotalConditions              int                       `json:"total_conditions"`
	TotalDecisionPoints          int                       `json:"total_decision_points"`
	Counts                       UPLM5ZCounts              `json:"counts"`
	RawSub16Candidates           int                       `json:"raw_sub16_candidates"`
	EligibleSub16Candidates      int                       `json:"eligible_sub16_candidates"`
	RawLengthHistogram           map[string]int            `json:"raw_length_histogram"`
	CategoryHistograms           map[string]map[string]int `json:"category_histograms"`
	CanonicalHazardSelectorUsed  bool                      `json:"canonical_hazard_selector_used"`
	PolicyChanged                bool                      `json:"policy_changed"`
	FutureScheduleUsed           bool                      `json:"future_schedule_used"`
	AdaptivePolicySelectionUsed  bool                      `json:"adaptive_policy_selection_used"`
	CounterfactualOnly           bool                      `json:"counterfactual_only"`
	LiveActivation               bool                      `json:"live_activation"`
	Summaries                    []UPLM5ZCellSummary       `json:"summaries"`
}

func uplm5zCategory(arms []uplm2xArm, used map[int]bool, rot int, i int) (string,int) {
	if used[i] {
		if len(arms[i].r.order)==0{return "USED_THIS_DECISION",0}
		return "USED_THIS_DECISION",len(arms[i].r.order)
	}
	if len(arms[i].r.order)==0{return "EMPTY",0}
	n:=len(arms[i].r.order)
	victim:=arms[i].r.order[0]
	original:=false
	for j:=0;j<12;j++{if victim==uplm2nName(j,rot){original=true;break}}
	if !original{return "NON_ORIGINAL_VICTIM",n}
	if arms[i].reported[victim]{return "ALREADY_REPORTED",n}
	return "ELIGIBLE",n
}

func uplm5zRun(rot int,perm,profile string,budget,start,tp int,w UPLM4WWindow,bred,tred int)(decisionPoints int,counts UPLM5ZCounts,raw map[int]int,bycat map[string]map[int]int){
	arms:=uplm2yPermute(uplm2xArms(rot),perm)
	uplm3tPrepressure(arms,rot,profile)
	original:=map[string]bool{};for i:=0;i<12;i++{original[uplm2nName(i,rot)]=true}
	raw=map[int]int{}
	bycat=map[string]map[int]int{"USED_THIS_DECISION":{},"EMPTY":{},"NON_ORIGINAL_VICTIM":{},"ALREADY_REPORTED":{},"ELIGIBLE":{}}
	actions:=0
	for round:=0;round<6;round++{
		if round>=start{
			cap:=uplm4wCap(budget,round,w.BudgetCut,w.BudgetRestore,bred)
			eff:=uplm4wTP(tp,round,w.ThroughputCut,w.ThroughputRestore,tred)
			used:=map[int]bool{}
			for k:=0;k<eff && actions<cap;k++{
				decisionPoints++
				for i:=range arms{
					cat,n:=uplm5zCategory(arms,used,rot,i)
					switch cat{
					case "USED_THIS_DECISION":counts.UsedThisDecision++
					case "EMPTY":counts.Empty++
					case "NON_ORIGINAL_VICTIM":counts.NonOriginalVictim++
					case "ALREADY_REPORTED":counts.AlreadyReported++
					case "ELIGIBLE":counts.Eligible++
					}
					if n>0{raw[n]++;bycat[cat][n]++}
				}
				i:=uplm5uThreatArm(arms,used,rot,false)
				if i<0{break}
				if len(arms[i].r.order)==0{break}
				n:=arms[i].r.order[0]
				if !original[n] || arms[i].reported[n]{
					if pn,_,ok:=uplm2xFirstPending(&arms[i]);ok{n=pn}else{break}
				}
				arms[i].reported[n]=true;actions++;used[i]=true
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("5z-%s-%s-%d-%d-%d-%d-%d-%d",profile,perm,bred,tred,budget,start,tp,round),"x")}
	}
	return
}

func uplm5zAddCounts(a *UPLM5ZCounts,b UPLM5ZCounts){
	a.UsedThisDecision+=b.UsedThisDecision;a.Empty+=b.Empty;a.NonOriginalVictim+=b.NonOriginalVictim;a.AlreadyReported+=b.AlreadyReported;a.Eligible+=b.Eligible
}
func uplm5zStringHist(src map[int]int)map[string]int{out:=map[string]int{};keys:=make([]int,0,len(src));for k:=range src{keys=append(keys,k)};sort.Ints(keys);for _,k:=range keys{out[strconv.Itoa(k)]=src[k]};return out}
func uplm5zStringCats(src map[string]map[int]int)map[string]map[string]int{out:=map[string]map[string]int{};for k,v:=range src{out[k]=uplm5zStringHist(v)};return out}

func RunUPLM5Z()(UPLM5ZResult,error){
	parent,err:=RunUPLM5Y();if err!=nil{return UPLM5ZResult{},err}
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	windows:=[]UPLM4WWindow{{1,3,2,4},{2,4,1,3},{1,4,2,5},{2,5,1,4},{1,3,3,5},{3,5,1,3}}
	breds:=[]int{1,2};treds:=[]int{1,2};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"};budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{2,3,4,5}
	res:=UPLM5ZResult{Schema:UPLM5ZOriginSchema,Experiment:"UP-LM5Z-order-length-origin",SourceUPLM5YSeal:"f026b4aa1488ab1bd5d186d13ef1a5ddf1e64467",ParentDecisionPoints:parent.TotalDecisionPoints,ParentEligibleCandidates:parent.TotalEligibleCandidates,ParentObservedLengths:append([]int{},parent.ObservedLengths...),ParentSubhazardCandidates:parent.TotalSubhazardCandidates,RawLengthHistogram:map[string]int{},CategoryHistograms:map[string]map[string]int{},CanonicalHazardSelectorUsed:true,PolicyChanged:false,FutureScheduleUsed:false,AdaptivePolicySelectionUsed:false,CounterfactualOnly:true,LiveActivation:false}
	globalRaw:=map[int]int{};globalCats:=map[string]map[int]int{"USED_THIS_DECISION":{},"EMPTY":{},"NON_ORIGINAL_VICTIM":{},"ALREADY_REPORTED":{},"ELIGIBLE":{}}
	for _,profile:=range profiles{for _,br:=range breds{for _,tr:=range treds{for _,rot:=range rots{for _,perm:=range perms{
		sm:=UPLM5ZCellSummary{DeadlineProfile:profile,BudgetReduction:br,ThroughputReduction:tr,Rotation:rot,Permutation:perm};cellRaw:=map[int]int{};cellCats:=map[string]map[int]int{"USED_THIS_DECISION":{},"EMPTY":{},"NON_ORIGINAL_VICTIM":{},"ALREADY_REPORTED":{},"ELIGIBLE":{}}
		for _,w:=range windows{for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
			d,c,h,cats:=uplm5zRun(rot,perm,profile,budget,start,tp,w,br,tr);sm.Conditions++;sm.DecisionPoints+=d;uplm5zAddCounts(&sm.Counts,c);res.TotalConditions++;res.TotalDecisionPoints+=d;uplm5zAddCounts(&res.Counts,c)
			for n,v:=range h{cellRaw[n]+=v;globalRaw[n]+=v}
			for cat,hist:=range cats{for n,v:=range hist{cellCats[cat][n]+=v;globalCats[cat][n]+=v}}
		}}}}
		sm.RawLengthHistogram=uplm5zStringHist(cellRaw);sm.CategoryHistograms=uplm5zStringCats(cellCats);res.Summaries=append(res.Summaries,sm)
	}}}}}
	for n,v:=range globalRaw{if n<16{res.RawSub16Candidates+=v}}
	for n,v:=range globalCats["ELIGIBLE"]{if n<16{res.EligibleSub16Candidates+=v}}
	res.RawLengthHistogram=uplm5zStringHist(globalRaw);res.CategoryHistograms=uplm5zStringCats(globalCats)
	return res,nil
}
