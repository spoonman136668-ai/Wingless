package unitary

import "fmt"

const UPLM5THazardTriggerSchema = "wingless.up-lm5t-hazard-triggered-minimal-intervention.v1"

type UPLM5TSummary struct {
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	Policy string `json:"policy"`
	Conditions int `json:"conditions"`
	Actions int `json:"actions"`
	Failures int `json:"failures"`
	ProtectedHazards int `json:"protected_hazards"`
	UnprotectedHazards int `json:"unprotected_hazards"`
}

type UPLM5TResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM5SSeal string `json:"source_up_lm5s_seal"`
	Policies []string `json:"policies"`
	ConditionCells int `json:"condition_cells"`
	MatchedConditionsPerPolicyCell int `json:"matched_conditions_per_policy_cell"`
	ResourceReductionsUsed bool `json:"resource_reductions_used"`
	PrepressureUsed bool `json:"prepressure_used"`
	HazardTriggeredPolicyFixed bool `json:"hazard_triggered_policy_fixed"`
	AdaptivePolicySelectionUsed bool `json:"adaptive_policy_selection_used"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UPLM5TSummary `json:"summaries"`
}

func uplm5tHazardArm(arms []uplm2xArm,used map[int]bool,rot int)int{
	original:=map[string]bool{}
	for i:=0;i<12;i++{original[uplm2nName(i,rot)]=true}
	for i:=range arms{
		if used[i]||len(arms[i].r.order)<16{continue}
		victim:=arms[i].r.order[0]
		if original[victim]&&!arms[i].reported[victim]{return i}
	}
	return -1
}

func uplm5tRun(rot int,perm,policy string,budget,start,tp int)(actions,failed,protectedHazards,unprotectedHazards int){
	arms:=uplm2yPermute(uplm2xArms(rot),perm)
	original:=map[string]bool{}
	for i:=0;i<12;i++{original[uplm2nName(i,rot)]=true}
	for round:=0;round<6;round++{
		if round>=start{
			used:=map[int]bool{}
			for k:=0;k<tp&&actions<budget;k++{
				i:=-1
				if policy=="hazard_triggered"{i=uplm5tHazardArm(arms,used,rot)}else{i=uplm3cChoose(arms,policy,used)}
				if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{
					if policy=="hazard_triggered"&&len(arms[i].r.order)>=16{
						victim:=arms[i].r.order[0]
						if original[victim]&&!arms[i].reported[victim]{n=victim}
					}
					arms[i].reported[n]=true;actions++;used[i]=true
				}else{break}
			}
		}
		for i:=range arms{
			if len(arms[i].r.order)>=16{
				victim:=arms[i].r.order[0]
				if original[victim]{
					if arms[i].reported[victim]{protectedHazards++}else{unprotectedHazards++}
				}
			}
			arms[i].r.write(fmt.Sprintf("5t-%s-%s-%d-%d-%d-%d-%d",policy,perm,budget,start,tp,round,i),"x")
		}
	}
	for i:=range arms{_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f}
	return
}

func RunUPLM5T()(UPLM5TResult,error){
	rots:=[]int{5,13}
	perms:=[]string{"identity","reverse","rotate2"}
	policies:=[]string{"earliest_deadline","hazard_triggered","fixed_order"}
	budgets:=[]int{4,5,6,7}
	starts:=[]int{2,3,4,5}
	tps:=[]int{2,3,4,5}
	res:=UPLM5TResult{
		Schema:UPLM5THazardTriggerSchema,
		Experiment:"UP-LM5T-hazard-triggered-minimal-intervention",
		SourceUPLM5SSeal:"bbadd2b17bdb4b204eb43cf9e981dc8eaab97df2",
		Policies:policies,
		ConditionCells:18,
		MatchedConditionsPerPolicyCell:64,
		ResourceReductionsUsed:false,
		PrepressureUsed:false,
		HazardTriggeredPolicyFixed:true,
		AdaptivePolicySelectionUsed:false,
		CounterfactualOnly:true,
		LiveActivation:false,
	}
	for _,rot:=range rots{
		for _,perm:=range perms{
			for _,policy:=range policies{
				s:=UPLM5TSummary{Rotation:rot,Permutation:perm,Policy:policy}
				for _,budget:=range budgets{
					for _,start:=range starts{
						for _,tp:=range tps{
							a,f,p,u:=uplm5tRun(rot,perm,policy,budget,start,tp)
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
