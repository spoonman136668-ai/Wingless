package unitary

import "fmt"

const UPLM5UPrehazardSchema = "wingless.up-lm5u-prehazard-targeting.v1"

type UPLM5USummary struct {
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	Policy string `json:"policy"`
	Conditions int `json:"conditions"`
	Actions int `json:"actions"`
	Failures int `json:"failures"`
	ProtectedHazards int `json:"protected_hazards"`
	UnprotectedHazards int `json:"unprotected_hazards"`
}

type UPLM5UResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM5TSeal string `json:"source_up_lm5t_seal"`
	Policies []string `json:"policies"`
	ConditionCells int `json:"condition_cells"`
	MatchedConditionsPerPolicyCell int `json:"matched_conditions_per_policy_cell"`
	ResourceReductionsUsed bool `json:"resource_reductions_used"`
	PrepressureUsed bool `json:"prepressure_used"`
	PrehazardRuleFixed bool `json:"prehazard_rule_fixed"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	AdaptivePolicySelectionUsed bool `json:"adaptive_policy_selection_used"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UPLM5USummary `json:"summaries"`
}

func uplm5uThreatArm(arms []uplm2xArm,used map[int]bool,rot int,pre bool)int{
	original:=map[string]bool{}
	for i:=0;i<12;i++{original[uplm2nName(i,rot)]=true}
	for _,wantLen:=range []int{16,15}{
		if wantLen==15&&!pre{continue}
		for i:=range arms{
			if used[i]||len(arms[i].r.order)<wantLen{continue}
			if wantLen==15&&len(arms[i].r.order)!=15{continue}
			victim:=arms[i].r.order[0]
			if original[victim]&&!arms[i].reported[victim]{return i}
		}
	}
	return -1
}

func uplm5uRun(rot int,perm,policy string,budget,start,tp int)(actions,failed,protectedHazards,unprotectedHazards int){
	arms:=uplm2yPermute(uplm2xArms(rot),perm)
	original:=map[string]bool{}
	for i:=0;i<12;i++{original[uplm2nName(i,rot)]=true}
	for round:=0;round<6;round++{
		if round>=start{
			used:=map[int]bool{}
			for k:=0;k<tp&&actions<budget;k++{
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
				n:=arms[i].r.order[0]
				if !original[n]||arms[i].reported[n]{
					if pn,_,ok:=uplm2xFirstPending(&arms[i]);ok{n=pn}else{break}
				}
				arms[i].reported[n]=true;actions++;used[i]=true
			}
		}
		for i:=range arms{
			if len(arms[i].r.order)>=16{
				victim:=arms[i].r.order[0]
				if original[victim]{
					if arms[i].reported[victim]{protectedHazards++}else{unprotectedHazards++}
				}
			}
			arms[i].r.write(fmt.Sprintf("5u-%s-%s-%d-%d-%d-%d-%d",policy,perm,budget,start,tp,round,i),"x")
		}
	}
	for i:=range arms{_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f}
	return
}

func RunUPLM5U()(UPLM5UResult,error){
	rots:=[]int{5,13}
	perms:=[]string{"identity","reverse","rotate2"}
	policies:=[]string{"earliest_deadline","hazard_triggered","prehazard_triggered"}
	budgets:=[]int{4,5,6,7}
	starts:=[]int{2,3,4,5}
	tps:=[]int{2,3,4,5}
	res:=UPLM5UResult{
		Schema:UPLM5UPrehazardSchema,
		Experiment:"UP-LM5U-prehazard-targeting",
		SourceUPLM5TSeal:"0757294cbaa8c08c104fc6056d67ce236c320116",
		Policies:policies,
		ConditionCells:18,
		MatchedConditionsPerPolicyCell:64,
		ResourceReductionsUsed:false,
		PrepressureUsed:false,
		PrehazardRuleFixed:true,
		FutureScheduleUsed:false,
		AdaptivePolicySelectionUsed:false,
		CounterfactualOnly:true,
		LiveActivation:false,
	}
	for _,rot:=range rots{
		for _,perm:=range perms{
			for _,policy:=range policies{
				s:=UPLM5USummary{Rotation:rot,Permutation:perm,Policy:policy}
				for _,budget:=range budgets{
					for _,start:=range starts{
						for _,tp:=range tps{
							a,f,p,u:=uplm5uRun(rot,perm,policy,budget,start,tp)
							s.Conditions++;s.Actions+=a;s.Failures+=f;s.ProtectedHazards+=p;s.UnprotectedHazards+=u
						}
					}
				}
				res.Summaries=append(res.Summaries,s)
			}
		}
	}
	return res,nil
}
