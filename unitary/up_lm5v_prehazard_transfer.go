package unitary

import "fmt"

const UPLM5VPrehazardTransferSchema = "wingless.up-lm5v-prehazard-transfer.v1"

type UPLM5VSummary struct {
	DeadlineProfile string `json:"deadline_profile"`
	BudgetReduction int `json:"budget_reduction"`
	ThroughputReduction int `json:"throughput_reduction"`
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	Policy string `json:"policy"`
	MatchedConditions int `json:"matched_conditions"`
	Actions int `json:"actions"`
	Failures int `json:"failures"`
}

type UPLM5VResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM5USeal string `json:"source_up_lm5u_seal"`
	Policies []string `json:"policies"`
	DeadlineProfiles []string `json:"deadline_profiles"`
	PooledConditionCells int `json:"pooled_condition_cells"`
	MatchedConditionsPerPolicyCell int `json:"matched_conditions_per_policy_cell"`
	ResourceReductionsUsed bool `json:"resource_reductions_used"`
	PrepressureUsed bool `json:"prepressure_used"`
	PrehazardRuleFixed bool `json:"prehazard_rule_fixed"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	AdaptivePolicySelectionUsed bool `json:"adaptive_policy_selection_used"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UPLM5VSummary `json:"summaries"`
}

func uplm5vRun(rot int,perm,profile,policy string,budget,start,tp int,w UPLM4WWindow,bred,tred int)(actions,failed int){
	arms:=uplm2yPermute(uplm2xArms(rot),perm)
	uplm3tPrepressure(arms,rot,profile)
	for round:=0;round<6;round++{
		if round>=start{
			cap:=uplm4wCap(budget,round,w.BudgetCut,w.BudgetRestore,bred)
			eff:=uplm4wTP(tp,round,w.ThroughputCut,w.ThroughputRestore,tred)
			used:=map[int]bool{}
			for k:=0;k<eff&&actions<cap;k++{
				i:=-1
				switch policy{
				case "earliest_deadline":
					i=uplm3cChoose(arms,policy,used)
				case "hazard_triggered":
					i=uplm5uThreatArm(arms,used,rot,false)
				case "prehazard_triggered":
					i=uplm5uThreatArm(arms,used,rot,true)
				}
				if i<0{break}
				if len(arms[i].r.order)==0{break}
				original:=map[string]bool{}
				for q:=0;q<12;q++{original[uplm2nName(q,rot)]=true}
				n:=arms[i].r.order[0]
				if !original[n]||arms[i].reported[n]{
					if pn,_,ok:=uplm2xFirstPending(&arms[i]);ok{n=pn}else{break}
				}
				arms[i].reported[n]=true;actions++;used[i]=true
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("5v-%s-%s-%s-%d-%d-%d-%d-%d-%d",profile,policy,perm,bred,tred,budget,start,tp,round),"x")}
	}
	for i:=range arms{_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f}
	return
}

func RunUPLM5V()(UPLM5VResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	policies:=[]string{"earliest_deadline","hazard_triggered","prehazard_triggered"}
	windows:=[]UPLM4WWindow{{1,3,2,4},{2,4,1,3},{1,4,2,5},{2,5,1,4},{1,3,3,5},{3,5,1,3}}
	breds:=[]int{1,2};treds:=[]int{1,2}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{2,3,4,5}
	res:=UPLM5VResult{
		Schema:UPLM5VPrehazardTransferSchema,
		Experiment:"UP-LM5V-prehazard-transfer",
		SourceUPLM5USeal:"c32c07b4337cee68163ae173ca8dc47f4003333a",
		Policies:policies,
		DeadlineProfiles:profiles,
		PooledConditionCells:72,
		MatchedConditionsPerPolicyCell:384,
		ResourceReductionsUsed:true,
		PrepressureUsed:true,
		PrehazardRuleFixed:true,
		FutureScheduleUsed:false,
		AdaptivePolicySelectionUsed:false,
		CounterfactualOnly:true,
		LiveActivation:false,
	}
	for _,profile:=range profiles{
		for _,br:=range breds{
			for _,tr:=range treds{
				for _,rot:=range rots{
					for _,perm:=range perms{
						for _,policy:=range policies{
							sm:=UPLM5VSummary{DeadlineProfile:profile,BudgetReduction:br,ThroughputReduction:tr,Rotation:rot,Permutation:perm,Policy:policy}
							for _,w:=range windows{
								for _,budget:=range budgets{
									for _,start:=range starts{
										for _,tp:=range tps{
											a,f:=uplm5vRun(rot,perm,profile,policy,budget,start,tp,w,br,tr)
											sm.MatchedConditions++;sm.Actions+=a;sm.Failures+=f
										}
									}
								}
							}
							res.Summaries=append(res.Summaries,sm)
						}
					}
				}
			}
		}
	}
	return res,nil
}
