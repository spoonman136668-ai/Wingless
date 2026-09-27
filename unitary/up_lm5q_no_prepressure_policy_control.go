package unitary

import "fmt"

const UPLM5QNoPrepressureSchema = "wingless.up-lm5q-no-prepressure-policy-control.v1"

type UPLM5QSummary struct {
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	MatchedConditions int `json:"matched_conditions"`
	EarliestFixedDiffering int `json:"earliest_fixed_differing"`
	EarliestLatestDiffering int `json:"earliest_latest_differing"`
	EarliestLowerThanFixed int `json:"earliest_lower_than_fixed"`
	FixedLowerThanEarliest int `json:"fixed_lower_than_earliest"`
	EarliestFixedEqual int `json:"earliest_fixed_equal"`
	EarliestLowerThanLatest int `json:"earliest_lower_than_latest"`
	LatestLowerThanEarliest int `json:"latest_lower_than_earliest"`
	EarliestLatestEqual int `json:"earliest_latest_equal"`
	EarliestFailures int `json:"earliest_failures"`
	FixedFailures int `json:"fixed_failures"`
	LatestFailures int `json:"latest_failures"`
}

type UPLM5QResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM5PSeal string `json:"source_up_lm5p_seal"`
	Policies []string `json:"policies"`
	ConditionCells int `json:"condition_cells"`
	MatchedConditionsPerCell int `json:"matched_conditions_per_cell"`
	ResourceReductionsUsed bool `json:"resource_reductions_used"`
	PrepressureUsed bool `json:"prepressure_used"`
	OnlyPolicyChanged bool `json:"only_policy_changed"`
	AdaptivePolicySelectionUsed bool `json:"adaptive_policy_selection_used"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UPLM5QSummary `json:"summaries"`
}

func uplm5qRun(rot int,perm,policy string,budget,start,tp int) int {
	arms := uplm2yPermute(uplm2xArms(rot),perm)
	actions := 0
	for round:=0; round<6; round++ {
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
			arms[i].r.write(fmt.Sprintf("5q-%s-%s-%d-%d-%d-%d-%d",policy,perm,budget,start,tp,round,i),"x")
		}
	}
	failed:=0
	for i:=range arms {_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f}
	return failed
}

func RunUPLM5Q()(UPLM5QResult,error) {
	rots:=[]int{5,13}
	perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7}
	starts:=[]int{2,3,4,5}
	tps:=[]int{2,3,4,5}
	res:=UPLM5QResult{
		Schema:UPLM5QNoPrepressureSchema,
		Experiment:"UP-LM5Q-no-prepressure-policy-control",
		SourceUPLM5PSeal:"de9a735a34dd90578e9ebbfe42be77b24960c080",
		Policies:[]string{"earliest_deadline","fixed_order","latest_deadline"},
		ConditionCells:6,
		MatchedConditionsPerCell:64,
		ResourceReductionsUsed:false,
		PrepressureUsed:false,
		OnlyPolicyChanged:true,
		AdaptivePolicySelectionUsed:false,
		CounterfactualOnly:true,
		LiveActivation:false,
	}
	for _,rot:=range rots {
		for _,perm:=range perms {
			s:=UPLM5QSummary{Rotation:rot,Permutation:perm}
			for _,budget:=range budgets {
				for _,start:=range starts {
					for _,tp:=range tps {
						e:=uplm5qRun(rot,perm,"earliest_deadline",budget,start,tp)
						f:=uplm5qRun(rot,perm,"fixed_order",budget,start,tp)
						l:=uplm5qRun(rot,perm,"latest_deadline",budget,start,tp)
						s.MatchedConditions++; s.EarliestFailures+=e; s.FixedFailures+=f; s.LatestFailures+=l
						if e!=f {s.EarliestFixedDiffering++}
						if e!=l {s.EarliestLatestDiffering++}
						if e<f {s.EarliestLowerThanFixed++} else if f<e {s.FixedLowerThanEarliest++} else {s.EarliestFixedEqual++}
						if e<l {s.EarliestLowerThanLatest++} else if l<e {s.LatestLowerThanEarliest++} else {s.EarliestLatestEqual++}
					}
				}
			}
			res.Summaries=append(res.Summaries,s)
		}
	}
	return res,nil
}
