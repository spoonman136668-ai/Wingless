package unitary

const UP94BEventStructuredRoutingSchema = "wingless.up94b-event-structured-routing.v1"

type UP94BPoint struct {
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

type UP94BEventStructuredRoutingResult struct {
	Schema              string       `json:"schema"`
	Experiment          string       `json:"experiment"`
	SourceUP93BSeal     string       `json:"source_up93b_seal"`
	SourceUPLM0CSeal    string       `json:"source_up_lm0c_seal"`
	SeedBases           []int        `json:"seed_bases"`
	ExactRecallCap      int          `json:"exact_recall_cap"`
	RecurrentStateBytes int          `json:"recurrent_state_bytes"`
	ExternalSalienceBit bool         `json:"external_salience_bit"`
	Points              []UP94BPoint `json:"points"`
}

func up94bRun(arm string,totalWrites,targetKeys int,seedBases []int) UP94BPoint {
	queryHits,queryTotal:=0,0
	exactHits,episodes:=0,0
	maxRecall:=0
	admissions,truePos,storeEvents,falsePos:=0,0,0,0

	for _,base:=range seedBases{
		for ep:=0;ep<64;ep++{
			seed:=sq0Seed(base,811+totalWrites*7+targetKeys,ep)
			rng:=newSQ0RNG(seed)
			m:=newSQ0Machine("transport_gated_correction",seed)
			truth:=make(map[int]int,targetKeys)

			write:=func(event byte,key,value int){
				preOld,_,_,present:=m.decodeCarrier(key,32)
				m.write(key,value,32)
				isStore:=event=='S'
				if isStore{storeEvents++}
				admit:=false
				switch arm{
				case "fifo_all_writes":
					admit=true
				case "structured_store_only":
					admit=isStore
				case "structure_plus_change":
					admit=isStore && present && preOld!=value
				}
				if admit{
					m.recallWrite(key,value);admissions++
					if isStore{truePos++}else{falsePos++}
				}
			}

			for k:=0;k<targetKeys;k++{
				v:=rng.intn(32);truth[k]=v;write('S',k,v)
			}
			distractors:=totalWrites-2*targetKeys
			pre:=distractors/2;post:=distractors-pre
			for d:=0;d<pre;d++{write('O',1000+ep*1000+d,rng.intn(32))}
			for k:=0;k<targetKeys;k++{
				old:=truth[k];v:=(old+1+rng.intn(31))%32;truth[k]=v;write('S',k,v)
			}
			for d:=0;d<post;d++{write('O',200000+ep*1000+d,rng.intn(32))}

			exact:=true
			for k:=0;k<targetKeys;k++{
				got:=up91bQuery(m,k,32);queryTotal++
				if got==truth[k]{queryHits++}else{exact=false}
			}
			episodes++;if exact{exactHits++}
			if m.recallUsed>maxRecall{maxRecall=m.recallUsed}
		}
	}

	precision:=1.0
	if admissions>0{precision=float64(truePos)/float64(admissions)}
	recall:=0.0
	if storeEvents>0{recall=float64(truePos)/float64(storeEvents)}
	return UP94BPoint{
		Arm:arm,TotalWrites:totalWrites,TargetKeys:targetKeys,
		QueryAccuracy:float64(queryHits)/float64(queryTotal),
		ExactTargetSetAccuracy:float64(exactHits)/float64(episodes),
		RecallEntriesUsed:maxRecall,AdmissionPrecision:precision,AdmissionRecall:recall,
		FalsePositiveAdmissions:falsePos,
	}
}

func RunUP94B()(UP94BEventStructuredRoutingResult,error){
	seedBases:=[]int{135000000,136000000}
	result:=UP94BEventStructuredRoutingResult{
		Schema:UP94BEventStructuredRoutingSchema,Experiment:"UP-94B-event-structured-routing",
		SourceUP93BSeal:"b667821eb97e31fc69db72d6113135aeebeaa918",
		SourceUPLM0CSeal:"6323443c61da2841f9ad6889d13abacd23590e61",
		SeedBases:append([]int(nil),seedBases...),ExactRecallCap:16,RecurrentStateBytes:512,ExternalSalienceBit:false,
	}
	for _,arm:=range []string{"fifo_all_writes","structured_store_only","structure_plus_change"}{
		for _,writes:=range []int{32,64,128,256}{
			for _,targets:=range []int{4,8,16}{
				result.Points=append(result.Points,up94bRun(arm,writes,targets,seedBases))
			}
		}
	}
	return result,nil
}
