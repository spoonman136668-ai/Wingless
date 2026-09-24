package unitary

import (
	"fmt"
	"math"
)

const UP61ATrainingBudgetSchema = "wingless.up61a-training-budget.v1"

type UP61ABudgetPoint struct {
	Steps              int       `json:"steps"`
	Parameters         []float64 `json:"parameters"`
	DoubleAngleError   float64   `json:"double_angle_error"`
	TrainAccuracy      float64   `json:"train_accuracy"`
	Long128Accuracy    float64   `json:"long128_accuracy"`
	MaxNormDrift       float64   `json:"max_norm_drift"`
	MaxRoundTripError  float64   `json:"max_round_trip_error"`
	Gate               bool      `json:"gate"`
}
type UP61ATrainingBudgetResult struct {
	Schema            string             `json:"schema"`
	Experiment        string             `json:"experiment"`
	SourceUP60ASeal   string             `json:"source_up60a_seal"`
	Budgets           []int              `json:"budgets"`
	LearningRate      float64            `json:"learning_rate"`
	GradientEpsilon   float64            `json:"gradient_epsilon"`
	Initialization    []float64          `json:"initialization"`
	OtherHyperparametersChanged bool     `json:"other_hyperparameters_changed"`
	Points            []UP61ABudgetPoint `json:"points"`
}

func up61aTrainBudget(train []up56aSample, steps int) ([]float64,error) {
	const(lr=0.10;eps=1e-6)
	params:=[]float64{0.1,0.1}
	for step:=0;step<steps;step++ {
		gradient:=make([]float64,2)
		for j:=0;j<2;j++ {
			plus:=append([]float64(nil),params...)
			minus:=append([]float64(nil),params...)
			plus[j]+=eps;minus[j]-=eps
			lp,e:=up56aLoss(plus,train,up56aApplyUnitary);if e!=nil{return nil,e}
			lm,e:=up56aLoss(minus,train,up56aApplyUnitary);if e!=nil{return nil,e}
			gradient[j]=(lp-lm)/(2*eps)
			if !finite(gradient[j]){return nil,fmt.Errorf("UP61A nonfinite gradient")}
		}
		for j:=range params{params[j]-=lr*gradient[j]}
	}
	return params,nil
}

func RunUP61A()(UP61ATrainingBudgetResult,error){
	const(lr=0.10;eps=1e-6)
	budgets:=[]int{90,180,360,720,1440}
	train,err:=up56aTrainSamples();if err!=nil{return UP61ATrainingBudgetResult{},err}
	held,err:=up56aHeldSamples([]int{128});if err!=nil{return UP61ATrainingBudgetResult{},err}
	long:=up57aLongOnly(held)
	result:=UP61ATrainingBudgetResult{
		Schema:UP61ATrainingBudgetSchema,Experiment:"UP-61A-training-budget",
		SourceUP60ASeal:"c6d8494bbbb2ddc1bf2895c4405790558de78051",
		Budgets:append([]int(nil),budgets...),LearningRate:lr,GradientEpsilon:eps,
		Initialization:[]float64{0.1,0.1},OtherHyperparametersChanged:false,
	}
	for _,steps:=range budgets {
		params,err:=up61aTrainBudget(train,steps);if err!=nil{return UP61ATrainingBudgetResult{},err}
		_,trainAcc,_,_,_,err:=up56aEvaluate(params,train,up56aApplyUnitary,true);if err!=nil{return UP61ATrainingBudgetResult{},err}
		_,longAcc,norm,rt,_,err:=up56aEvaluate(params,long,up56aApplyUnitary,true);if err!=nil{return UP61ATrainingBudgetResult{},err}
		errorAngle:=params[1]-math.Pi/2
		result.Points=append(result.Points,UP61ABudgetPoint{
			Steps:steps,Parameters:append([]float64(nil),params...),DoubleAngleError:errorAngle,
			TrainAccuracy:trainAcc,Long128Accuracy:longAcc,MaxNormDrift:norm,MaxRoundTripError:rt,
			Gate:longAcc>=0.98&&norm<=1e-9&&rt<=1e-8,
		})
	}
	return result,nil
}
