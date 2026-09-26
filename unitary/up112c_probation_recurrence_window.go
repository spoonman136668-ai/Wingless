package unitary

const UP112CProbationWindowSchema = "wingless.up112c-probation-recurrence-window.v1"

type UP112CPoint struct {
	Arm                      string  `json:"arm"`
	RepeatAfterWrites        int     `json:"repeat_after_writes"`
	ExpectedRepeatAdmission  bool    `json:"expected_repeat_admission"`
	RepeatAdmissionRate      float64 `json:"repeat_admission_rate"`
	ExpectedDecisionAccuracy float64 `json:"expected_decision_accuracy"`
	OneShotFalseAdmissions   int     `json:"one_shot_false_admissions"`
	OneShotWrites            int     `json:"one_shot_writes"`
	OneShotFalseAdmissionRate float64 `json:"one_shot_false_admission_rate"`
	MeanGenerationClears     float64 `json:"mean_generation_clears"`
	AdmissionMemoryBytes     int     `json:"admission_memory_bytes"`
	CounterBytes             int     `json:"counter_bytes"`
}

type UP112CProbationWindowResult struct {
	Schema                 string        `json:"schema"`
	Experiment             string        `json:"experiment"`
	SourceUP111CSeal       string        `json:"source_up111c_seal"`
	GenerationInterval     int           `json:"generation_interval"`
	AdmissionMemoryBytes   int           `json:"admission_memory_bytes"`
	CounterBytes           int           `json:"counter_bytes"`
	EpisodesPerCell        int           `json:"episodes_per_cell"`
	ExactMemoryLayerUsed   bool          `json:"exact_memory_layer_used"`
	QueryEvidenceUsed      bool          `json:"query_evidence_used"`
	FutureOracleUsed       bool          `json:"future_oracle_used"`
	Points                 []UP112CPoint `json:"points"`
}

type up112cBloomSidecar struct {
	filter *up95cMemory
	counter int
	clears int
}

func newUP112CBloom() *up112cBloomSidecar {
	return &up112cBloomSidecar{filter:&up95cMemory{width:1024,seen:make([]uint64,16)}}
}

func (x *up112cBloomSidecar) process(key int) bool {
	admit:=x.filter.seenBefore(key)
	if !admit { x.filter.markSeen(key) }
	x.counter++
	if x.counter==32 {
		x.filter.clearSeen()
		x.counter=0
		x.clears++
	}
	return admit
}

type up112cExactSidecar struct {
	table [32]uint32
	counter int
	clears int
}

func (x *up112cExactSidecar) clear() {
	for i:=range x.table { x.table[i]=0 }
}

func (x *up112cExactSidecar) process(key int) bool {
	code:=uint32(key)+1
	admit:=false
	for _,v:=range x.table {
		if v==code { admit=true;break }
	}
	if !admit {
		inserted:=false
		for i,v:=range x.table {
			if v==0 { x.table[i]=code;inserted=true;break }
		}
		if !inserted { panic("UP112C exact probation table full before boundary") }
	}
	x.counter++
	if x.counter==32 {
		x.clear()
		x.counter=0
		x.clears++
	}
	return admit
}

func up112cRun(arm string,gap int) UP112CPoint {
	repeatAdmit,correct:=0,0
	falseOneShot,oneShot,totalClears:=0,0,0
	expected:=gap<=31

	for ep:=0;ep<64;ep++ {
		target:=1000000+gap*100000+ep*1000
		nextUnique:=0
		writes:=0
		var process func(int)bool
		var clears func()int

		if arm=="bloom1024" {
			x:=newUP112CBloom()
			process=x.process
			clears=func()int{return x.clears}
		}else{
			x:=&up112cExactSidecar{}
			process=x.process
			clears=func()int{return x.clears}
		}

		process(target)
		writes++

		for i:=0;i<gap-1;i++ {
			k:=target+1+nextUnique;nextUnique++
			if process(k) { falseOneShot++ }
			oneShot++;writes++
		}

		admitted:=process(target)
		writes++
		if admitted { repeatAdmit++ }
		if admitted==expected { correct++ }

		for writes<64 {
			k:=target+1+nextUnique;nextUnique++
			if process(k) { falseOneShot++ }
			oneShot++;writes++
		}
		totalClears+=clears()
	}

	rate:=0.0
	if oneShot>0 { rate=float64(falseOneShot)/float64(oneShot) }
	return UP112CPoint{
		Arm:arm,RepeatAfterWrites:gap,ExpectedRepeatAdmission:expected,
		RepeatAdmissionRate:float64(repeatAdmit)/64.0,
		ExpectedDecisionAccuracy:float64(correct)/64.0,
		OneShotFalseAdmissions:falseOneShot,OneShotWrites:oneShot,OneShotFalseAdmissionRate:rate,
		MeanGenerationClears:float64(totalClears)/64.0,
		AdmissionMemoryBytes:128,CounterBytes:1,
	}
}

func RunUP112C()(UP112CProbationWindowResult,error){
	result:=UP112CProbationWindowResult{
		Schema:UP112CProbationWindowSchema,Experiment:"UP-112C-probation-recurrence-window",
		SourceUP111CSeal:"bab1e8fdd41476f0d889e2e68dd7275c2b91d4d7",
		GenerationInterval:32,AdmissionMemoryBytes:128,CounterBytes:1,EpisodesPerCell:64,
		ExactMemoryLayerUsed:false,QueryEvidenceUsed:false,FutureOracleUsed:false,
	}
	for _,arm:=range []string{"bloom1024","exact32"} {
		for _,gap:=range []int{1,8,16,24,31,32,33,40} {
			result.Points=append(result.Points,up112cRun(arm,gap))
		}
	}
	return result,nil
}
