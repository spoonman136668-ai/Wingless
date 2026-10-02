package unitary

import (
	"math"
	"sort"
)

type wlmLmContextSimilarityAliasInductionR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

type wlmLmContextSimilarityAliasInductionR1DSU struct {
	parent []int
}

func wlmLmContextSimilarityAliasInductionR1NewDSU(n int) *wlmLmContextSimilarityAliasInductionR1DSU {
	p:=make([]int,n)
	for i:=range p { p[i]=i }
	return &wlmLmContextSimilarityAliasInductionR1DSU{parent:p}
}

func (d *wlmLmContextSimilarityAliasInductionR1DSU) find(x int) int {
	if d.parent[x]!=x {
		d.parent[x]=d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *wlmLmContextSimilarityAliasInductionR1DSU) union(a,b int) {
	ra:=d.find(a); rb:=d.find(b)
	if ra==rb { return }
	if ra<rb { d.parent[rb]=ra } else { d.parent[ra]=rb }
}

func wlmLmContextSimilarityAliasInductionR1BridgeRecords(shuffled bool) [][]uint8 {
	records:=make([][]uint8,0,576)
	for entity:=0;entity<9;entity++ {
		for family:=0;family<2;family++ {
			counts:=[4]int{8,8,8,8}
			if family==1 {
				counts=[4]int{9,7,8,8}
			}
			for contextIndex:=0;contextIndex<4;contextIndex++ {
				for repeat:=0;repeat<counts[contextIndex];repeat++ {
					records=append(records,wlmLmRawRepContextualAliasInductionR2BridgeRecord(family,entity,contextIndex,repeat,shuffled))
				}
			}
		}
	}
	return records
}

func wlmLmContextSimilarityAliasInductionR1Induce(
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	bridge [][]uint8,
	tvMax float64,
) wlmLmRawRepContextualAliasInductionR2Classes {
	observed,decodeFailures:=wlmSiContextCountJitterFalsificationR1ObservedCounts(reps,bridge)
	dsu:=wlmLmContextSimilarityAliasInductionR1NewDSU(len(reps))
	for i:=0;i<len(reps);i++ {
		for j:=i+1;j<len(reps);j++ {
			jaccard,tv,ok:=wlmSiContextCountJitterFalsificationR1PairStats(observed[i],observed[j])
			if ok && jaccard>=1.0 && tv<=tvMax {
				dsu.union(i,j)
			}
		}
	}
	groupMap:=make(map[int][]int)
	for i:=0;i<len(reps);i++ {
		root:=dsu.find(i)
		groupMap[root]=append(groupMap[root],i)
	}
	roots:=make([]int,0,len(groupMap))
	for root:=range groupMap { roots=append(roots,root) }
	sort.Ints(roots)
	surfaceToClass:=make([]int,len(reps))
	members:=make([][]int,0,len(roots))
	for classID,root:=range roots {
		group:=append([]int(nil),groupMap[root]...)
		sort.Ints(group)
		members=append(members,group)
		for _,surfaceID:=range group {
			surfaceToClass[surfaceID]=classID
		}
	}
	return wlmLmRawRepContextualAliasInductionR2Classes{
		surfaceToClass:surfaceToClass,
		members:members,
		decodeFailures:decodeFailures,
	}
}

// RunWlmLmContextSimilarityAliasInductionR1 tests bounded context-similarity alias induction.
func RunWlmLmContextSimilarityAliasInductionR1() interface{} {
	const tvThreshold=0.05
	metrics:=map[string]float64{
		"task_training_record_count":0,
		"selected_surface_representation_count":0,
		"selected_true_surface_motif_match_count":0,
		"jitter_bridge_record_count":0,
		"bridge_decode_failure_count":0,
		"similarity_tv_threshold":tvThreshold,
		"induced_alias_class_count":0,
		"induced_two_member_class_count":0,
		"correct_evaluator_alias_pair_count":0,
		"task_training_decode_failure_count":0,
		"class_transition_missing_count":0,
		"same_family_evaluation_count":0,
		"induced_same_family_accuracy":0,
		"cross_family_evaluation_count":0,
		"induced_cross_family_accuracy":0,
		"exact_surface_cross_family_accuracy":0,
		"exact_surface_cross_family_missing_transition_count":0,
		"shuffled_alias_class_count":0,
		"shuffled_correct_evaluator_alias_pair_count":0,
		"shuffled_cross_family_accuracy":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"capacity_growth_event_count":0,
		"invalid_row_count":0,
		"counter_overflow_count":0,
	}

	taskRecords:=wlmLmRawRepContextualAliasInductionR2TaskRecords()
	metrics["task_training_record_count"]=float64(len(taskRecords))
	reps:=wlmSiRawRepAliasInvarianceFalsificationR1LearnRepresentation(taskRecords)
	metrics["selected_surface_representation_count"]=float64(len(reps))
	metrics["selected_true_surface_motif_match_count"]=float64(wlmSiRawRepAliasInvarianceFalsificationR1TrueMotifCount(reps))
	repIndex:=wlmLmRawRepContextualAliasInductionR2RepIndex(reps)

	bridge:=wlmLmContextSimilarityAliasInductionR1BridgeRecords(false)
	metrics["jitter_bridge_record_count"]=float64(len(bridge))
	classes:=wlmLmContextSimilarityAliasInductionR1Induce(reps,bridge,tvThreshold)
	metrics["bridge_decode_failure_count"]=float64(classes.decodeFailures)
	metrics["induced_alias_class_count"]=float64(len(classes.members))
	metrics["induced_two_member_class_count"]=float64(wlmSiContextCountJitterFalsificationR1TwoMemberCount(classes))
	metrics["correct_evaluator_alias_pair_count"]=float64(wlmLmRawRepContextualAliasInductionR2CorrectAliasPairs(reps,repIndex,classes))

	classModel:=wlmLmRawRepContextualAliasInductionR2BuildClassModel(taskRecords,reps,classes,metrics)
	metrics["task_training_decode_failure_count"]=float64(classModel.decodeFails)
	if classModel.valid {
		metrics["class_transition_missing_count"]=float64(classModel.transition.missingCount())
	}

	surfaceModel:=wlmLmRawRepContextualAliasInductionR2BuildSurfaceModel(taskRecords,reps,metrics)

	shuffledBridge:=wlmLmContextSimilarityAliasInductionR1BridgeRecords(true)
	shuffledClasses:=wlmLmContextSimilarityAliasInductionR1Induce(reps,shuffledBridge,tvThreshold)
	metrics["shuffled_alias_class_count"]=float64(len(shuffledClasses.members))
	metrics["shuffled_correct_evaluator_alias_pair_count"]=float64(wlmLmRawRepContextualAliasInductionR2CorrectAliasPairs(reps,repIndex,shuffledClasses))
	if shuffledClasses.decodeFailures!=0 {
		metrics["invalid_row_count"]+=float64(shuffledClasses.decodeFailures)
	}
	shuffledModel:=wlmLmRawRepContextualAliasInductionR2BuildClassModel(taskRecords,reps,shuffledClasses,metrics)
	if shuffledModel.decodeFails!=0 {
		metrics["invalid_row_count"]+=float64(shuffledModel.decodeFails)
	}

	sameHits:=0
	sameTotal:=0
	for family:=0;family<2;family++ {
		for op:=0;op<5;op++ {
			for left:=0;left<4;left++ {
				for right:=0;right<4;right++ {
					if wlmLmRawRepContextualAliasInductionR2EvaluateClass(
						family,family,left,right,op,
						1000+family*100+op*20+left*4+right,
						reps,repIndex,classes,classModel,metrics,
					) {
						sameHits++
					}
					sameTotal++
				}
			}
		}
	}
	metrics["same_family_evaluation_count"]=float64(sameTotal)
	if sameTotal>0 {
		metrics["induced_same_family_accuracy"]=float64(sameHits)/float64(sameTotal)
	}

	crossHits:=0
	exactHits:=0
	exactMissing:=0
	shuffledHits:=0
	crossTotal:=0
	for valueFamily:=0;valueFamily<2;valueFamily++ {
		opFamily:=1-valueFamily
		for op:=0;op<5;op++ {
			for left:=0;left<4;left++ {
				for right:=0;right<4;right++ {
					tag:=2000+valueFamily*100+op*20+left*4+right
					if wlmLmRawRepContextualAliasInductionR2EvaluateClass(
						valueFamily,opFamily,left,right,op,tag,
						reps,repIndex,classes,classModel,metrics,
					) {
						crossHits++
					}
					if ok,missing:=wlmLmRawRepContextualAliasInductionR2EvaluateSurfaceCross(
						valueFamily,opFamily,left,right,op,tag,
						reps,repIndex,surfaceModel,metrics,
					); ok {
						exactHits++
						if missing { exactMissing++ }
					} else if missing {
						exactMissing++
					}
					if wlmLmRawRepContextualAliasInductionR2EvaluateClass(
						valueFamily,opFamily,left,right,op,tag,
						reps,repIndex,shuffledClasses,shuffledModel,metrics,
					) {
						shuffledHits++
					}
					crossTotal++
				}
			}
		}
	}
	metrics["cross_family_evaluation_count"]=float64(crossTotal)
	if crossTotal>0 {
		metrics["induced_cross_family_accuracy"]=float64(crossHits)/float64(crossTotal)
		metrics["exact_surface_cross_family_accuracy"]=float64(exactHits)/float64(crossTotal)
		metrics["shuffled_cross_family_accuracy"]=float64(shuffledHits)/float64(crossTotal)
	}
	metrics["exact_surface_cross_family_missing_transition_count"]=float64(exactMissing)

	for _,value:=range metrics {
		if math.IsNaN(value)||math.IsInf(value,0) {
			metrics["invalid_row_count"]++
		}
	}
	return wlmLmContextSimilarityAliasInductionR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-CONTEXT-SIMILARITY-ALIAS-INDUCTION-R1",
		Metrics:metrics,
	}
}
