package unitary

import "math"

type wlmSiContextCountJitterFalsificationR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmSiContextCountJitterFalsificationR1BridgeRecords() [][]uint8 {
	records:=make([][]uint8,0,576)
	for entity:=0;entity<9;entity++ {
		for family:=0;family<2;family++ {
			counts:=[4]int{8,8,8,8}
			if family==1 {
				counts=[4]int{9,7,8,8}
			}
			for contextIndex:=0;contextIndex<4;contextIndex++ {
				for repeat:=0;repeat<counts[contextIndex];repeat++ {
					records=append(records,wlmLmRawRepContextualAliasInductionR2BridgeRecord(family,entity,contextIndex,repeat,false))
				}
			}
		}
	}
	return records
}

func wlmSiContextCountJitterFalsificationR1ObservedCounts(
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	bridge [][]uint8,
)([]map[wlmLmRawRepContextualAliasInductionR2Context]uint32,int) {
	repIndex:=make(map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int,len(reps))
	for id,key:=range reps {
		repIndex[key]=id
	}
	out:=make([]map[wlmLmRawRepContextualAliasInductionR2Context]uint32,len(reps))
	for i:=range out {
		out[i]=make(map[wlmLmRawRepContextualAliasInductionR2Context]uint32)
	}
	failures:=0
	for _,data:=range bridge {
		matches:=0
		for i:=2;i+6<=len(data);i++ {
			key:=wlmSiRawRepAliasInvarianceFalsificationR1Key{data[i],data[i+1],data[i+2],data[i+3]}
			repID,ok:=repIndex[key]
			if !ok {
				continue
			}
			context:=wlmLmRawRepContextualAliasInductionR2Context{data[i-2],data[i-1],data[i+4],data[i+5]}
			out[repID][context]++
			matches++
			i+=3
		}
		if matches!=1 {
			failures++
		}
	}
	return out,failures
}

func wlmSiContextCountJitterFalsificationR1PairStats(
	a,b map[wlmLmRawRepContextualAliasInductionR2Context]uint32,
)(float64,float64,bool) {
	union:=make(map[wlmLmRawRepContextualAliasInductionR2Context]bool)
	for key:=range a { union[key]=true }
	for key:=range b { union[key]=true }
	if len(union)==0 {
		return 0,0,false
	}
	intersection:=0
	var totalA,totalB uint64
	for _,count:=range a { totalA+=uint64(count) }
	for _,count:=range b { totalB+=uint64(count) }
	if totalA==0||totalB==0 {
		return 0,0,false
	}
	tv:=0.0
	for key:=range union {
		ca,oka:=a[key]
		cb,okb:=b[key]
		if oka&&okb {
			intersection++
		}
		pa:=float64(ca)/float64(totalA)
		pb:=float64(cb)/float64(totalB)
		tv+=math.Abs(pa-pb)
	}
	jaccard:=float64(intersection)/float64(len(union))
	return jaccard,0.5*tv,true
}

func wlmSiContextCountJitterFalsificationR1TwoMemberCount(classes wlmLmRawRepContextualAliasInductionR2Classes) int {
	count:=0
	for _,members:=range classes.members {
		if len(members)==2 { count++ }
	}
	return count
}

// RunWlmSiContextCountJitterFalsificationR1 tests exact-signature robustness to tiny count drift.
func RunWlmSiContextCountJitterFalsificationR1() interface{} {
	metrics:=map[string]float64{
		"task_training_record_count":0,
		"selected_surface_representation_count":0,
		"selected_true_surface_motif_match_count":0,
		"normal_bridge_record_count":0,
		"normal_bridge_decode_failure_count":0,
		"normal_induced_alias_class_count":0,
		"normal_two_member_class_count":0,
		"normal_correct_evaluator_alias_pair_count":0,
		"jitter_bridge_record_count":0,
		"jitter_bridge_decode_failure_count":0,
		"minimum_alias_pair_context_support_jaccard":1,
		"maximum_alias_pair_total_variation_distance":0,
		"jitter_induced_alias_class_count":0,
		"jitter_two_member_class_count":0,
		"jitter_correct_evaluator_alias_pair_count":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"capacity_growth_event_count":0,
		"invalid_row_count":0,
	}

	taskRecords:=wlmLmRawRepContextualAliasInductionR2TaskRecords()
	metrics["task_training_record_count"]=float64(len(taskRecords))
	reps:=wlmSiRawRepAliasInvarianceFalsificationR1LearnRepresentation(taskRecords)
	metrics["selected_surface_representation_count"]=float64(len(reps))
	metrics["selected_true_surface_motif_match_count"]=float64(wlmSiRawRepAliasInvarianceFalsificationR1TrueMotifCount(reps))
	repIndex:=wlmLmRawRepContextualAliasInductionR2RepIndex(reps)

	normal:=wlmLmRawRepContextualAliasInductionR2BridgeRecords(false)
	metrics["normal_bridge_record_count"]=float64(len(normal))
	normalClasses:=wlmLmRawRepContextualAliasInductionR2Induce(reps,normal)
	metrics["normal_bridge_decode_failure_count"]=float64(normalClasses.decodeFailures)
	metrics["normal_induced_alias_class_count"]=float64(len(normalClasses.members))
	metrics["normal_two_member_class_count"]=float64(wlmSiContextCountJitterFalsificationR1TwoMemberCount(normalClasses))
	metrics["normal_correct_evaluator_alias_pair_count"]=float64(wlmLmRawRepContextualAliasInductionR2CorrectAliasPairs(reps,repIndex,normalClasses))

	jitter:=wlmSiContextCountJitterFalsificationR1BridgeRecords()
	metrics["jitter_bridge_record_count"]=float64(len(jitter))
	jitterClasses:=wlmLmRawRepContextualAliasInductionR2Induce(reps,jitter)
	metrics["jitter_bridge_decode_failure_count"]=float64(jitterClasses.decodeFailures)
	metrics["jitter_induced_alias_class_count"]=float64(len(jitterClasses.members))
	metrics["jitter_two_member_class_count"]=float64(wlmSiContextCountJitterFalsificationR1TwoMemberCount(jitterClasses))
	metrics["jitter_correct_evaluator_alias_pair_count"]=float64(wlmLmRawRepContextualAliasInductionR2CorrectAliasPairs(reps,repIndex,jitterClasses))

	observed,observedFailures:=wlmSiContextCountJitterFalsificationR1ObservedCounts(reps,jitter)
	if observedFailures!=0 {
		metrics["invalid_row_count"]+=float64(observedFailures)
	}
	for entity:=0;entity<9;entity++ {
		key0:=wlmLmRawRepContextualAliasInductionR2SurfaceKey(0,entity)
		key1:=wlmLmRawRepContextualAliasInductionR2SurfaceKey(1,entity)
		id0,ok0:=repIndex[key0]
		id1,ok1:=repIndex[key1]
		if !ok0||!ok1 {
			metrics["invalid_row_count"]++
			continue
		}
		jaccard,tv,ok:=wlmSiContextCountJitterFalsificationR1PairStats(observed[id0],observed[id1])
		if !ok {
			metrics["invalid_row_count"]++
			continue
		}
		if jaccard<metrics["minimum_alias_pair_context_support_jaccard"] {
			metrics["minimum_alias_pair_context_support_jaccard"]=jaccard
		}
		if tv>metrics["maximum_alias_pair_total_variation_distance"] {
			metrics["maximum_alias_pair_total_variation_distance"]=tv
		}
	}
	for _,value:=range metrics {
		if math.IsNaN(value)||math.IsInf(value,0) {
			metrics["invalid_row_count"]++
		}
	}
	return wlmSiContextCountJitterFalsificationR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-SI-CONTEXT-COUNT-JITTER-FALSIFICATION-R1",
		Metrics:metrics,
	}
}
