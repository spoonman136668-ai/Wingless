package unitary

import "fmt"

const UP63AProgramDistributionSchema="wingless.up63a-program-distribution.v1"

type UP63AProgramDistributionResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP61ASeal string `json:"source_up61a_seal"`
	TrainingSteps int `json:"training_steps"`
	LearningRate float64 `json:"learning_rate"`
	GradientEpsilon float64 `json:"gradient_epsilon"`
	Lengths []int `json:"lengths"`
	Parameters []float64 `json:"parameters"`
	AggregateAccuracy float64 `json:"aggregate_accuracy"`
	CategoryMetrics []UP56ACategoryMetric `json:"category_metrics"`
	MaxNormDrift float64 `json:"max_norm_drift"`
	MaxRoundTripError float64 `json:"max_round_trip_error"`
	Gate bool `json:"gate"`
}

func up63aInstruction(i,seed int,r up56aRoles,mode string) UP56AInstruction {
	switch mode {
	case "control_burst":
		target:=(i+seed)%3
		ca:=(target+1)%3
		cb:=(target+2)%3
		return UP56AInstruction{
			Kind:"double_controlled_value_swap",
			TargetRole:target,
			ControlRoleA:ca,ControlValueA:(i+seed+r[0])%3,
			ControlRoleB:cb,ControlValueB:(2*i+seed+r[1])%3,
			Start:(i*i+seed+r[2])%3,
		}
	case "swap_heavy":
		if i%4!=3{return UP56AInstruction{Kind:"role_swap",Position:(i+seed)%2}}
		target:=(i+seed)%3;ca:=(target+1)%3;cb:=(target+2)%3
		return UP56AInstruction{
			Kind:"double_controlled_value_swap",TargetRole:target,
			ControlRoleA:ca,ControlValueA:(seed+r[1])%3,
			ControlRoleB:cb,ControlValueB:(i+r[2])%3,
			Start:(i+seed+r[0])%3,
		}
	case "shifted_controls":
		target:=(i+seed)%3;ca:=(target+1)%3;cb:=(target+2)%3
		return UP56AInstruction{
			Kind:"double_controlled_value_swap",TargetRole:target,
			ControlRoleA:ca,ControlValueA:1+(i+seed)%2,
			ControlRoleB:cb,ControlValueB:1+(i+seed+1)%2,
			Start:(i+2*seed+r[0])%3,
		}
	default:
		return UP56AInstruction{Kind:"role_swap",Position:(i+seed)%2}
	}
}

func up63aHeldSamples(lengths []int)([]up56aSample,error){
	var out []up56aSample
	for _,length:=range lengths{
		for _,r:=range up56aAllRoles(){
			for seed:=0;seed<3;seed++{
				for _,mode:=range []string{"control_burst","swap_heavy","shifted_controls"}{
					p:=make([]UP56AInstruction,length)
					for i:=range p{p[i]=up63aInstruction(i,seed,r,mode)}
					t,err:=up56aTarget(r,p);if err!=nil{return nil,err}
					out=append(out,up56aSample{Roles:r,Program:p,Target:t,Category:fmt.Sprintf("%s_%d",mode,length)})
				}
				half:=length/2
				prefix:=make([]UP56AInstruction,half)
				for i:=range prefix{prefix[i]=up63aInstruction(i,seed,r,"control_burst")}
				p:=append([]UP56AInstruction(nil),prefix...)
				for i:=len(prefix)-1;i>=0;i--{p=append(p,prefix[i])}
				t,err:=up56aTarget(r,p);if err!=nil{return nil,err}
				out=append(out,up56aSample{Roles:r,Program:p,Target:t,Category:fmt.Sprintf("palindrome_%d",length)})
			}
		}
	}
	return out,nil
}

func RunUP63A()(UP63AProgramDistributionResult,error){
	const(steps=1440;lr=0.10;eps=1e-6)
	lengths:=[]int{32,64,128,256}
	train,err:=up56aTrainSamples();if err!=nil{return UP63AProgramDistributionResult{},err}
	params,err:=up61aTrainBudget(train,steps);if err!=nil{return UP63AProgramDistributionResult{},err}
	held,err:=up63aHeldSamples(lengths);if err!=nil{return UP63AProgramDistributionResult{},err}
	_,accuracy,norm,rt,cats,err:=up56aEvaluate(params,held,up56aApplyUnitary,true);if err!=nil{return UP63AProgramDistributionResult{},err}
	gate:=accuracy>=0.98&&norm<=1e-9&&rt<=1e-8
	for _,m:=range cats{if m.Accuracy<0.98{gate=false}}
	return UP63AProgramDistributionResult{
		Schema:UP63AProgramDistributionSchema,Experiment:"UP-63A-program-distribution",
		SourceUP61ASeal:"e3b319607b677dd7e9537d9651767c013e50d07e",
		TrainingSteps:steps,LearningRate:lr,GradientEpsilon:eps,Lengths:append([]int(nil),lengths...),
		Parameters:append([]float64(nil),params...),AggregateAccuracy:accuracy,CategoryMetrics:cats,
		MaxNormDrift:norm,MaxRoundTripError:rt,Gate:gate,
	},nil
}
