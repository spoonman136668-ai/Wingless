package unitary

const UP89AMultiQueryLatestValueSchema = "wingless.up89a-multiquery-latest-value.v1"

type UP89AMultiQueryLatestValuePoint struct {
	Representation        string  `json:"representation"`
	FeatureDimension      int     `json:"feature_dimension"`
	Length                int     `json:"length"`
	NoiseAmplitude        float64 `json:"noise_amplitude"`
	Scenarios             int     `json:"scenarios"`
	QueriesPerScenario    int     `json:"queries_per_scenario"`
	QueryAccuracy         float64 `json:"query_accuracy"`
	ExactQuerySetAccuracy float64 `json:"exact_query_set_accuracy"`
	Gate                  bool    `json:"gate"`
}

type UP89AMultiQueryLatestValueResult struct {
	Schema          string                            `json:"schema"`
	Experiment      string                            `json:"experiment"`
	SourceUP88ASeal string                            `json:"source_up88a_seal"`
	TrainingStates  int                               `json:"training_states"`
	TrainingSteps   int                               `json:"training_steps"`
	LearningRate    float64                           `json:"learning_rate"`
	Scenarios       int                               `json:"scenarios"`
	QueriesPerScenario int                            `json:"queries_per_scenario"`
	Lengths         []int                             `json:"lengths"`
	NoiseLevels     []float64                         `json:"noise_levels"`
	Retraining      bool                              `json:"retraining"`
	ExactRecallSideChannel bool                       `json:"exact_recall_side_channel"`
	Points          []UP89AMultiQueryLatestValuePoint `json:"points"`
}

func up89aQueries(scenario int)[4]int{
	return [4]int{
		scenario%6,
		(scenario*3+1)%6,
		(scenario*5+2)%6,
		(scenario*7+3)%6,
	}
}

func up89aRunRepresentation(name string,dim,width int,featureFn func(up81aRoles)[]float64,lengths []int,noises []float64)([]UP89AMultiQueryLatestValuePoint,error){
	all:=up81aAllRoles()
	var trainRoles []up81aRoles
	for _,r:=range all{
		if up81aPrimary(r)==0&&up81aSecondary(r)==0&&up81aTertiary(r)==0{
			trainRoles=append(trainRoles,r)
		}
	}
	const scenarios=64
	const queryCount=4
	var points []UP89AMultiQueryLatestValuePoint
	for _,length:=range lengths{
		for _,noise:=range noises{
			queryHits:=0
			exactError:=make([]bool,scenarios)
			for role:=0;role<6;role++{
				train:=make([]headSample,0,len(trainRoles))
				for _,r:=range trainRoles{
					train=append(train,headSample{features:featureFn(r),target:r[role]})
				}
				head,_,err:=trainLinearSoftmax(train,3,dim,1,1.0)
				if err!=nil{return nil,err}
				for scenario:=0;scenario<scenarios;scenario++{
					truth:=up87aScenarioRoles(scenario)
					state:=append([]float64(nil),featureFn(truth)...)
					for step:=0;step<length;step++{
						writeRole:=(scenario*5+step*7)%6
						writeValue:=(scenario+2*step+writeRole+1)%3
						truth[writeRole]=writeValue
						exact:=featureFn(truth)
						start:=writeRole*width
						copy(state[start:start+width],exact[start:start+width])
						up87aAddNoise(state,scenario,step,noise)
					}
					p,err:=head.probabilities(state)
					if err!=nil{return nil,err}
					got,_,err:=classAndMargin(p)
					if err!=nil{return nil,err}
					queries:=up89aQueries(scenario)
					for _,q:=range queries{
						if q!=role{continue}
						if got==truth[role]{queryHits++}else{exactError[scenario]=true}
					}
				}
			}
			exactHits:=0
			for _,bad:=range exactError{if !bad{exactHits++}}
			queryAcc:=float64(queryHits)/float64(scenarios*queryCount)
			exactAcc:=float64(exactHits)/scenarios
			points=append(points,UP89AMultiQueryLatestValuePoint{
				Representation:name,FeatureDimension:dim,Length:length,NoiseAmplitude:noise,
				Scenarios:scenarios,QueriesPerScenario:queryCount,
				QueryAccuracy:queryAcc,ExactQuerySetAccuracy:exactAcc,
				Gate:queryAcc>=0.99&&exactAcc>=0.95,
			})
		}
	}
	return points,nil
}

func RunUP89A()(UP89AMultiQueryLatestValueResult,error){
	lengths:=[]int{64,256,1024}
	noises:=[]float64{0,0.001,0.002}
	result:=UP89AMultiQueryLatestValueResult{
		Schema:UP89AMultiQueryLatestValueSchema,
		Experiment:"UP-89A-multiquery-latest-value",
		SourceUP88ASeal:"e97939a2ede84e22349c0ca2658be15fa48ceda7",
		TrainingStates:27,TrainingSteps:1,LearningRate:1.0,
		Scenarios:64,QueriesPerScenario:4,
		Lengths:append([]int(nil),lengths...),
		NoiseLevels:append([]float64(nil),noises...),
		Retraining:false,ExactRecallSideChannel:false,
	}
	onehot,err:=up89aRunRepresentation("raw_factorized_one_hot",18,3,up81aRawOneHot,lengths,noises)
	if err!=nil{return UP89AMultiQueryLatestValueResult{},err}
	result.Points=append(result.Points,onehot...)
	simplex,err:=up89aRunRepresentation("role_factorized_simplex",12,2,up81aSimplex,lengths,noises)
	if err!=nil{return UP89AMultiQueryLatestValueResult{},err}
	result.Points=append(result.Points,simplex...)
	return result,nil
}
