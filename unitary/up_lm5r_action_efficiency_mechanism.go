package unitary

import "fmt"

const UPLM5RActionEfficiencySchema = "wingless.up-lm5r-action-efficiency-mechanism.v1"

type UPLM5RSummary struct {
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	MatchedConditions int `json:"matched_conditions"`
	EarliestFixedFailureDiffering int `json:"earliest_fixed_failure_differing"`
	EarliestFixedSameActionsDiffering int `json:"earliest_fixed_same_actions_differing"`
	EarliestLowerFixedSameActions int `json:"earliest_lower_fixed_same_actions"`
	FixedLowerEarliestSameActions int `json:"fixed_lower_earliest_same_actions"`
	EarliestLatestFailureDiffering int `json:"earliest_latest_failure_differing"`
	EarliestLatestSameActionsDiffering int `json:"earliest_latest_same_actions_differing"`
	EarliestLowerLatestSameActions int `json:"earliest_lower_latest_same_actions"`
	LatestLowerEarliestSameActions int `json:"latest_lower_earliest_same_actions"`
	EarliestActions int `json:"earliest_actions"`
	FixedActions int `json:"fixed_actions"`
	LatestActions int `json:"latest_actions"`
}

type UPLM5RResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM5QSeal string `json:"source_up_lm5q_seal"`
	Policies []string `json:"policies"`
	ConditionCells int `json:"condition_cells"`
	MatchedConditionsPerCell int `json:"matched_conditions_per_cell"`
	ResourceReductionsUsed bool `json:"resource_reductions_used"`
	PrepressureUsed bool `json:"prepressure_used"`
	ActionCountMeasured bool `json:"action_count_measured"`
	OnlyPolicyChanged bool `json:"only_policy_changed"`
	AdaptivePolicySelectionUsed bool `json:"adaptive_policy_selection_used"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UPLM5RSummary `json:"summaries"`
}

func uplm5rRun(rot int,perm,policy string,budget,start,tp int)(actions,failed int) {
	arms:=uplm2yPermute(uplm2xArms(rot),perm)
	for round:=0;round<6;round++ {
		if round>=start {
			used:=map[int]bool{}
			for k:=0;k<tp&&actions<budget;k++ {
				i:=-1
				if policy=="latest_deadline" {i=uplm5oChooseLatest(arms,used)} else {i=uplm3cChoose(arms,policy,used)}
				if i<0 {break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok {
					arms[i].reported[n]=true; actions++; used[i]=true
				} else {break}
			}
		}
		for i:=range arms {
			arms[i].r.write(fmt.Sprintf("5r-%s-%s-%d-%d-%d-%d-%d",policy,perm,budget,start,tp,round,i),"x")
		}
	}
	for i:=range arms {_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f}
	return
}

func RunUPLM5R()(UPLM5RResult,error) {
	rots:=[]int{5,13}
	perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7}
	starts:=[]int{2,3,4,5}
	tps:=[]int{2,3,4,5}
	res:=UPLM5RResult{
		Schema:UPLM5RActionEfficiencySchema,
		Experiment:"UP-LM5R-action-efficiency-mechanism",
		SourceUPLM5QSeal:"42c018a98747a50b4b2a3e02cf8a4774c19b5ae9",
		Policies:[]string{"earliest_deadline","fixed_order","latest_deadline"},
		ConditionCells:6,
		MatchedConditionsPerCell:64,
		ResourceReductionsUsed:false,
		PrepressureUsed:false,
		ActionCountMeasured:true,
		OnlyPolicyChanged:true,
		AdaptivePolicySelectionUsed:false,
		CounterfactualOnly:true,
		LiveActivation:false,
	}
	for _,rot:=range rots {
		for _,perm:=range perms {
			s:=UPLM5RSummary{Rotation:rot,Permutation:perm}
			for _,budget:=range budgets {
				for _,start:=range starts {
					for _,tp:=range tps {
						ea,ef:=uplm5rRun(rot,perm,"earliest_deadline",budget,start,tp)
						fa,ff:=uplm5rRun(rot,perm,"fixed_order",budget,start,tp)
						la,lf:=uplm5rRun(rot,perm,"latest_deadline",budget,start,tp)
						s.MatchedConditions++; s.EarliestActions+=ea; s.FixedActions+=fa; s.LatestActions+=la
						if ef!=ff {
							s.EarliestFixedFailureDiffering++
							if ea==fa {
								s.EarliestFixedSameActionsDiffering++
								if ef<ff {s.EarliestLowerFixedSameActions++} else {s.FixedLowerEarliestSameActions++}
							}
						}
						if ef!=lf {
							s.EarliestLatestFailureDiffering++
							if ea==la {
								s.EarliestLatestSameActionsDiffering++
								if ef<lf {s.EarliestLowerLatestSameActions++} else {s.LatestLowerEarliestSameActions++}
							}
						}
					}
				}
			}
			res.Summaries=append(res.Summaries,s)
		}
	}
	return res,nil
}
