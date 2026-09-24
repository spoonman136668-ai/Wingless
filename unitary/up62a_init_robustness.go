package unitary

import (
	"fmt"
	"math"
)

const UP62AInitSchema = "wingless.up62a-initialization-robustness.v1"

type UP62AInitPoint struct {
	Initialization      []float64 `json:"initialization"`
	Parameters          []float64 `json:"parameters"`
	DoubleAngleError    float64   `json:"double_angle_error"`
	TrainAccuracy       float64   `json:"train_accuracy"`
	Long128Accuracy     float64   `json:"long128_accuracy"`
	MaxNormDrift        float64   `json:"max_norm_drift"`
	MaxRoundTripError   float64   `json:"max_round_trip_error"`
	Gate                bool      `json:"gate"`
}
type UP62AInitResult struct {
	Schema          string           `json:"schema"`
	Experiment      string           `json:"experiment"`
	SourceUP61ASeal string           `json:"source_up61a_seal"`
	Steps           int              `json:"steps"`
	LearningRate    float64          `json:"learning_rate"`
	GradientEpsilon float64          `json:"gradient_epsilon"`
	Initializations [][]float64      `json:"initializations"`
	Points          []UP62AInitPoint `json:"points"`
}

func up62aTrain(train []up56aSample, init []float64)([]float64,error){
	const(steps=1440;lr=0.10;eps=1e-6)
	params:=append([]float64(nil),init...)
	for step:=0;step<steps;step++{
		g:=make([]float64,2)
		for j:=0;j<2;j++{
			plus:=append([]float64(nil),params...);minus:=append([]float64(nil),params...)
			plus[j]+=eps;minus[j]-=eps
			lp,e:=up56aLoss(plus,train,up56aApplyUnitary);if e!=nil{return nil,e}
			lm,e:=up56aLoss(minus,train,up56aApplyUnitary);if e!=nil{return nil,e}
			g[j]=(lp-lm)/(2*eps)
			if !finite(g[j]){return nil,fmt.Errorf("UP62A nonfinite gradient")}
		}
		for j:=range params{params[j]-=lr*g[j]}
	}
	return params,nil
}

func RunUP62A()(UP62AInitResult,error){
	const(steps=1440;lr=0.10;eps=1e-6)
	inits:=[][]float64{{0.05,0.05},{0.1,0.1},{0.25,0.25},{0.5,0.5},{-0.1,0.1}}
	train,err:=up56aTrainSamples();if err!=nil{return UP62AInitResult{},err}
	held,err:=up56aHeldSamples([]int{128});if err!=nil{return UP62AInitResult{},err}
	long:=up57aLongOnly(held)
	result:=UP62AInitResult{Schema:UP62AInitSchema,Experiment:"UP-62A-initialization-robustness",SourceUP61ASeal:"e3b319607b677dd7e9537d9651767c013e50d07e",Steps:steps,LearningRate:lr,GradientEpsilon:eps,Initializations:inits}
	for _,init:=range inits{
		params,err:=up62aTrain(train,init);if err!=nil{return UP62AInitResult{},err}
		_,ta,_,_,_,err:=up56aEvaluate(params,train,up56aApplyUnitary,true);if err!=nil{return UP62AInitResult{},err}
		_,ha,norm,rt,_,err:=up56aEvaluate(params,long,up56aApplyUnitary,true);if err!=nil{return UP62AInitResult{},err}
		result.Points=append(result.Points,UP62AInitPoint{
			Initialization:append([]float64(nil),init...),Parameters:append([]float64(nil),params...),
			DoubleAngleError:params[1]-math.Pi/2,TrainAccuracy:ta,Long128Accuracy:ha,
			MaxNormDrift:norm,MaxRoundTripError:rt,
			Gate:ha>=0.98&&norm<=1e-9&&rt<=1e-8,
		})
	}
	return result,nil
}
