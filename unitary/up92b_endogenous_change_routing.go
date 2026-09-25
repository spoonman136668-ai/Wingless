package unitary

const UP92BEndogenousChangeRoutingSchema = "wingless.up92b-endogenous-change-routing.v1"

type UP92BEndogenousChangePoint struct {
	Arm                    string  `json:"arm"`
	TotalWrites            int     `json:"total_writes"`
	TargetKeys             int     `json:"target_keys"`
	QueryAccuracy          float64 `json:"query_accuracy"`
	ExactTargetSetAccuracy float64 `json:"exact_target_set_accuracy"`
	RecallEntriesUsed      int     `json:"recall_entries_used"`
	AdmissionPrecision     float64 `json:"admission_precision"`
	AdmissionRecall        float64 `json:"admission_recall"`
	FalsePositiveAdmissions int    `json:"false_positive_admissions"`
}

type UP92BEndogenousChangeRoutingResult struct {
	Schema              string                         `json:"schema"`
	Experiment          string                         `json:"experiment"`
	SourceUP91BSeal     string                         `json:"source_up91b_seal"`
	SeedBases           []int                          `json:"seed_bases"`
	ExactRecallCap      int                            `json:"exact_recall_cap"`
	RecurrentStateBytes int                            `json:"recurrent_state_bytes"`
	Points              []UP92BEndogenousChangePoint   `json:"points"`
}

func up92bMaybeAdmit(m *sq0Machine, arm string, key,value,vocab int, trueSalient bool, preOld int, prePresent bool) (admitted bool) {
	switch arm {
	case "fifo_all_writes":
		admitted=true
	case "explicit_target_rewrite":
		admitted=trueSalient
	case "endogenous_change":
		admitted=prePresent && preOld!=value
	}
	if admitted { m.recallWrite(key,value) }
	return
}

func up92bRun(arm string,totalWrites,targetKeys int,seedBases []int) UP92BEndogenousChangePoint {
	queryHits,queryTotal:=0,0
	exactHits,episodes:=0,0
	maxRecall:=0
	admissions,truePos,trueSalientTotal,falsePos:=0,0,0,0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,401+totalWrites*5+targetKeys,ep)
			rng:=newSQ0RNG(seed)
			m:=newSQ0Machine("transport_gated_correction",seed)
			truth:=make(map[int]int,targetKeys)

			writeEvent:=func(key,value int,trueSalient bool){
				preOld,_,_,prePresent:=m.decodeCarrier(key,32)
				m.write(key,value,32)
				admitted:=up92bMaybeAdmit(m,arm,key,value,32,trueSalient,preOld,prePresent)
				if trueSalient { trueSalientTotal++ }
				if admitted {
					admissions++
					if trueSalient { truePos++ } else { falsePos++ }
				}
			}

			for k:=0;k<targetKeys;k++ {
				v:=rng.intn(32)
				truth[k]=v
				writeEvent(k,v,false)
			}

			distractors:=totalWrites-2*targetKeys
			pre:=distractors/2
			post:=distractors-pre
			for d:=0;d<pre;d++ {
				key:=1000+ep*1000+d
				writeEvent(key,rng.intn(32),false)
			}

			for k:=0;k<targetKeys;k++ {
				old:=truth[k]
				v:=(old+1+rng.intn(31))%32
				truth[k]=v
				writeEvent(k,v,true)
			}

			for d:=0;d<post;d++ {
				key:=200000+ep*1000+d
				writeEvent(key,rng.intn(32),false)
			}

			exact:=true
			for k:=0;k<targetKeys;k++ {
				got:=up91bQuery(m,k,32)
				queryTotal++
				if got==truth[k] { queryHits++ } else { exact=false }
			}
			episodes++
			if exact { exactHits++ }
			if m.recallUsed>maxRecall { maxRecall=m.recallUsed }
		}
	}

	precision:=1.0
	if admissions>0 { precision=float64(truePos)/float64(admissions) }
	recall:=0.0
	if trueSalientTotal>0 { recall=float64(truePos)/float64(trueSalientTotal) }

	return UP92BEndogenousChangePoint{
		Arm:arm,TotalWrites:totalWrites,TargetKeys:targetKeys,
		QueryAccuracy:float64(queryHits)/float64(queryTotal),
		ExactTargetSetAccuracy:float64(exactHits)/float64(episodes),
		RecallEntriesUsed:maxRecall,
		AdmissionPrecision:precision,AdmissionRecall:recall,
		FalsePositiveAdmissions:falsePos,
	}
}

func RunUP92B()(UP92BEndogenousChangeRoutingResult,error){
	seedBases:=[]int{127000000,128000000}
	result:=UP92BEndogenousChangeRoutingResult{
		Schema:UP92BEndogenousChangeRoutingSchema,
		Experiment:"UP-92B-endogenous-change-routing",
		SourceUP91BSeal:"9c2b5c8195bba9589afb488236c5ce6ec9271939",
		SeedBases:append([]int(nil),seedBases...),
		ExactRecallCap:16,RecurrentStateBytes:512,
	}
	for _,arm:=range []string{"fifo_all_writes","explicit_target_rewrite","endogenous_change"} {
		for _,writes:=range []int{32,64,128,256} {
			for _,targets:=range []int{4,8,16} {
				result.Points=append(result.Points,up92bRun(arm,writes,targets,seedBases))
			}
		}
	}
	return result,nil
}
