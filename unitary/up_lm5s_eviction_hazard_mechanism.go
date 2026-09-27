package unitary

import "fmt"

const UPLM5SEvictionHazardSchema = "wingless.up-lm5s-eviction-hazard-mechanism.v1"

type UPLM5SSummary struct {
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	Policy string `json:"policy"`
	Conditions int `json:"conditions"`
	Actions int `json:"actions"`
	Failures int `json:"failures"`
	ProtectedEvictionHazards int `json:"protected_eviction_hazards"`
	UnprotectedEvictionHazards int `json:"unprotected_eviction_hazards"`
}

type UPLM5SResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM5RSeal string `json:"source_up_lm5r_seal"`
	Policies []string `json:"policies"`
	ConditionCells int `json:"condition_cells"`
	MatchedConditionsPerPolicyCell int `json:"matched_conditions_per_policy_cell"`
	ResourceReductionsUsed bool `json:"resource_reductions_used"`
	PrepressureUsed bool `json:"prepressure_used"`
	EvictionHazardMeasured bool `json:"eviction_hazard_measured"`
	OnlyPolicyChanged bool `json:"only_policy_changed"`
	AdaptivePolicySelectionUsed bool `json:"adaptive_policy_selection_used"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UPLM5SSummary `json:"summaries"`
}

func uplm5sRun(rot int,perm,policy string,budget,start,tp int)(actions,failed,protectedHazards,unprotectedHazards int) {
	arms:=uplm2yPermute(uplm2xArms(rot),perm)
	original:=map[string]bool{}
	for i:=0;i<12;i++ {original[uplm2nName(i,rot)]=true}
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
			if len(arms[i].r.order)>=16 {
				victim:=arms[i].r.order[0]
				if original[victim] {
					if arms[i].reported[victim] {protectedHazards++} else {unprotectedHazards++}
				}
			}
			arms[i].r.write(fmt.Sprintf("5s-%s-%s-%d-%d-%d-%d-%d",policy,perm,budget,start,tp,round,i),"x")
		}
	}
	for i:=range arms {_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f}
	return
}

func RunUPLM5S()(UPLM5SResult,error) {
	rots:=[]int{5,13}
	perms:=[]string{"identity","reverse","rotate2"}
	policies:=[]string{"earliest_deadline","fixed_order","latest_deadline"}
	budgets:=[]int{4,5,6,7}
	starts:=[]int{2,3,4,5}
	tps:=[]int{2,3,4,5}
	res:=UPLM5SResult{
		Schema:UPLM5SEvictionHazardSchema,
		Experiment:"UP-LM5S-eviction-hazard-mechanism",
		SourceUPLM5RSeal:"62c0264596499891be9f1702b5b696acd5af5e18",
		Policies:policies,
		ConditionCells:18,
		MatchedConditionsPerPolicyCell:64,
		ResourceReductionsUsed:false,
		PrepressureUsed:false,
		EvictionHazardMeasured:true,
		OnlyPolicyChanged:true,
		AdaptivePolicySelectionUsed:false,
		CounterfactualOnly:true,
		LiveActivation:false,
	}
	for _,rot:=range rots {
		for _,perm:=range perms {
			for _,policy:=range policies {
				s:=UPLM5SSummary{Rotation:rot,Permutation:perm,Policy:policy}
				for _,budget:=range budgets {
					for _,start:=range starts {
						for _,tp:=range tps {
							a,f,p,u:=uplm5sRun(rot,perm,policy,budget,start,tp)
							s.Conditions++; s.Actions+=a; s.Failures+=f; s.ProtectedEvictionHazards+=p; s.UnprotectedEvictionHazards+=u
						}
					}
				}
				res.Summaries=append(res.Summaries,s)
			}
		}
	}
	return res,nil
}
