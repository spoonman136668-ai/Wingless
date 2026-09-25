package unitary

import "math"

const UP93BConfidenceChangeRoutingSchema = "wingless.up93b-confidence-change-routing.v1"

type UP93BConfidenceChangePoint struct {
	Arm                    string  `json:"arm"`
	TotalWrites            int     `json:"total_writes"`
	TargetKeys             int     `json:"target_keys"`
	QueryAccuracy          float64 `json:"query_accuracy"`
	ExactTargetSetAccuracy float64 `json:"exact_target_set_accuracy"`
	RecallEntriesUsed      int     `json:"recall_entries_used"`
	AdmissionPrecision     float64 `json:"admission_precision"`
	AdmissionRecall        float64 `json:"admission_recall"`
	FalsePositiveAdmissions int    `json:"false_positive_admissions"`
	MarginThreshold        float64 `json:"margin_threshold"`
}

type UP93BConfidenceChangeRoutingResult struct {
	Schema              string                         `json:"schema"`
	Experiment          string                         `json:"experiment"`
	SourceUP92BSeal     string                         `json:"source_up92b_seal"`
	SeedBases           []int                          `json:"seed_bases"`
	ExactRecallCap      int                            `json:"exact_recall_cap"`
	RecurrentStateBytes int                            `json:"recurrent_state_bytes"`
	ThresholdSelection  bool                           `json:"threshold_selection"`
	Points              []UP93BConfidenceChangePoint   `json:"points"`
}

type up93bArm struct {
	name string
	explicit bool
	raw bool
	threshold float64
}

func up93bRun(cfg up93bArm,totalWrites,targetKeys int,seedBases []int) UP93BConfidenceChangePoint {
	queryHits,queryTotal:=0,0
	exactHits,episodes:=0,0
	maxRecall:=0
	admissions,truePos,trueSalientTotal,falsePos:=0,0,0,0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,601+totalWrites*5+targetKeys,ep)
			rng:=newSQ0RNG(seed)
			m:=newSQ0Machine("transport_gated_correction",seed)
			truth:=make(map[int]int,targetKeys)

			writeEvent:=func(key,value int,trueSalient bool){
				preOld,best,second,present:=m.decodeCarrier(key,32)
				den:=math.Max(1,math.Abs(best)+math.Abs(second))
				margin:=(best-second)/den
				m.write(key,value,32)

				admit:=false
				switch {
				case cfg.explicit:
					admit=trueSalient
				case cfg.raw:
					admit=present && preOld!=value
				default:
					admit=present && preOld!=value && margin>=cfg.threshold
				}
				if trueSalient { trueSalientTotal++ }
				if admit {
					m.recallWrite(key,value)
					admissions++
					if trueSalient { truePos++ } else { falsePos++ }
				}
			}

			for k:=0;k<targetKeys;k++ {
				v:=rng.intn(32);truth[k]=v;writeEvent(k,v,false)
			}
			distractors:=totalWrites-2*targetKeys
			pre:=distractors/2;post:=distractors-pre
			for d:=0;d<pre;d++ { writeEvent(1000+ep*1000+d,rng.intn(32),false) }
			for k:=0;k<targetKeys;k++ {
				old:=truth[k];v:=(old+1+rng.intn(31))%32;truth[k]=v;writeEvent(k,v,true)
			}
			for d:=0;d<post;d++ { writeEvent(200000+ep*1000+d,rng.intn(32),false) }

			exact:=true
			for k:=0;k<targetKeys;k++ {
				got:=up91bQuery(m,k,32);queryTotal++
				if got==truth[k] {queryHits++} else {exact=false}
			}
			episodes++;if exact{exactHits++}
			if m.recallUsed>maxRecall {maxRecall=m.recallUsed}
		}
	}

	precision:=1.0
	if admissions>0 { precision=float64(truePos)/float64(admissions) }
	recall:=0.0
	if trueSalientTotal>0 { recall=float64(truePos)/float64(trueSalientTotal) }

	return UP93BConfidenceChangePoint{
		Arm:cfg.name,TotalWrites:totalWrites,TargetKeys:targetKeys,
		QueryAccuracy:float64(queryHits)/float64(queryTotal),
		ExactTargetSetAccuracy:float64(exactHits)/float64(episodes),
		RecallEntriesUsed:maxRecall,AdmissionPrecision:precision,AdmissionRecall:recall,
		FalsePositiveAdmissions:falsePos,MarginThreshold:cfg.threshold,
	}
}

func RunUP93B()(UP93BConfidenceChangeRoutingResult,error){
	seedBases:=[]int{131000000,132000000}
	arms:=[]up93bArm{
		{name:"explicit_target_rewrite",explicit:true},
		{name:"endogenous_change_raw",raw:true},
		{name:"endogenous_change_margin_025",threshold:0.25},
		{name:"endogenous_change_margin_050",threshold:0.50},
		{name:"endogenous_change_margin_075",threshold:0.75},
	}
	result:=UP93BConfidenceChangeRoutingResult{
		Schema:UP93BConfidenceChangeRoutingSchema,
		Experiment:"UP-93B-confidence-change-routing",
		SourceUP92BSeal:"ece60e33e57c3bf830948f2a2cf4d1487b3af085",
		SeedBases:append([]int(nil),seedBases...),ExactRecallCap:16,RecurrentStateBytes:512,
		ThresholdSelection:false,
	}
	for _,arm:=range arms {
		for _,writes:=range []int{32,64,128,256} {
			for _,targets:=range []int{4,8,16} {
				result.Points=append(result.Points,up93bRun(arm,writes,targets,seedBases))
			}
		}
	}
	return result,nil
}
