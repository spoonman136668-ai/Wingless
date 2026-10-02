package unitary

import "math"

type wlmLmRawRepZeroshotMultistepR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmRawRepZeroshotMultistepR1ProgramRecord(left,right int,ops []int,tag int) []uint8 {
	state:=uint32(0x6a09e667 ^ uint32((left+1)*131+(right+1)*977+(tag+1)*65537+len(ops)*8191))
	out:=make([]uint8,0,140)
	out=append(out,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,6+(tag%5))...)
	lk:=wlmLmRawRepZeroshotAliasFamilyR1Family2ValueMotif(left)
	out=append(out,lk[:]...)
	out=append(out,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,8+((left+tag)%5))...)
	rk:=wlmLmRawRepZeroshotAliasFamilyR1Family2ValueMotif(right)
	out=append(out,rk[:]...)
	for step,op:=range ops {
		out=append(out,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,7+((step+op+tag)%7))...)
		ok:=wlmLmRawRepZeroshotAliasFamilyR1Family2OpMotif(op)
		out=append(out,ok[:]...)
	}
	out=append(out,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,9+((tag+len(ops))%7))...)
	return out
}

func wlmLmRawRepZeroshotMultistepR1ExecuteClass(
	raw []uint8,
	ops []int,
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	classes wlmLmRawRepContextualAliasInductionR2Classes,
	model wlmLmRawRepContextualAliasInductionR2ClassModel,
	metrics map[string]float64,
)(int,int,bool) {
	surface:=wlmSiRawRepAliasInvarianceFalsificationR1Decode(raw,reps)
	if len(surface)!=2+len(ops) {
		return 0,0,false
	}
	if surface[0]<0||surface[0]>=len(classes.surfaceToClass)||surface[1]<0||surface[1]>=len(classes.surfaceToClass) {
		metrics["invalid_row_count"]++
		return 0,0,false
	}
	leftClass:=classes.surfaceToClass[surface[0]]
	rightClass:=classes.surfaceToClass[surface[1]]
	left,ok0:=model.valueMap[leftClass]
	right,ok1:=model.valueMap[rightClass]
	if !ok0||!ok1 {
		metrics["invalid_row_count"]++
		return 0,0,false
	}
	for step:=range ops {
		surfaceOp:=surface[2+step]
		if surfaceOp<0||surfaceOp>=len(classes.surfaceToClass) {
			metrics["invalid_row_count"]++
			return 0,0,false
		}
		classOp:=classes.surfaceToClass[surfaceOp]
		op,ok:=model.opMap[classOp]
		if !ok {
			metrics["invalid_row_count"]++
			return 0,0,false
		}
		left,right=model.transition.predict(op,left,right)
	}
	return left,right,true
}

func wlmLmRawRepZeroshotMultistepR1ExecuteSurface(
	raw []uint8,
	ops []int,
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	model wlmLmRawRepZeroshotAliasFamilyR1SurfaceModel,
	metrics map[string]float64,
)(int,int,bool,bool) {
	surface:=wlmSiRawRepAliasInvarianceFalsificationR1Decode(raw,reps)
	if len(surface)!=2+len(ops) {
		return 0,0,false,false
	}
	left,right:=surface[0],surface[1]
	anyMissing:=false
	for step:=range ops {
		op:=surface[2+step]
		pl,pr,found:=model.predict(op,left,right)
		if !found {
			anyMissing=true
			pl,pr=0,0
		}
		left,right=pl,pr
	}
	if left<0||left>=len(reps)||right<0||right>=len(reps) {
		metrics["invalid_row_count"]++
		return 0,0,false,anyMissing
	}
	return left,right,true,anyMissing
}

func wlmLmRawRepZeroshotMultistepR1Family2RawOutputExact(
	leftDense,rightDense,targetLeftSemantic,targetRightSemantic int,
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	repIndex map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int,
	classes wlmLmRawRepContextualAliasInductionR2Classes,
	model wlmLmRawRepContextualAliasInductionR2ClassModel,
) bool {
	if leftDense<0||leftDense>=len(model.valueIDs)||rightDense<0||rightDense>=len(model.valueIDs) {
		return false
	}
	leftClass:=model.valueIDs[leftDense]
	rightClass:=model.valueIDs[rightDense]
	leftFound:=false
	rightFound:=false
	var leftKey,rightKey wlmSiRawRepAliasInvarianceFalsificationR1Key
	for semantic:=0;semantic<4;semantic++ {
		key:=wlmLmRawRepZeroshotAliasFamilyR1Family2ValueMotif(semantic)
		surfaceID,ok:=repIndex[key]
		if !ok||surfaceID<0||surfaceID>=len(classes.surfaceToClass) {
			continue
		}
		classID:=classes.surfaceToClass[surfaceID]
		if classID==leftClass {
			leftKey=key
			leftFound=true
		}
		if classID==rightClass {
			rightKey=key
			rightFound=true
		}
	}
	if !leftFound||!rightFound {
		return false
	}
	return leftKey==wlmLmRawRepZeroshotAliasFamilyR1Family2ValueMotif(targetLeftSemantic)&&
		rightKey==wlmLmRawRepZeroshotAliasFamilyR1Family2ValueMotif(targetRightSemantic)
}

// RunWlmLmRawRepZeroshotMultistepR1 tests zero-shot multi-step execution through family-2 raw representations.
func RunWlmLmRawRepZeroshotMultistepR1() interface{} {
	metrics:=map[string]float64{
		"bridge_record_count":0,
		"selected_surface_representation_count":0,
		"selected_true_surface_motif_match_count":0,
		"induced_alias_class_count":0,
		"correct_evaluator_alias_triplet_count":0,
		"task_training_record_count":0,
		"family2_task_training_record_count":0,
		"class_transition_key_count":0,
		"heldout_program_count":0,
		"family2_operation_step_count":0,
		"program_decode_failure_count":0,
		"conjugated_right_program_count":0,
		"induced_conjugated_right_accuracy":0,
		"long_program_count":0,
		"induced_long_program_accuracy":0,
		"induced_overall_multistep_accuracy":0,
		"raw_output_pair_exact_accuracy":0,
		"exact_surface_overall_multistep_accuracy":0,
		"exact_surface_programs_with_missing_transition":0,
		"shuffled_bridge_overall_multistep_accuracy":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"capacity_growth_event_count":0,
		"invalid_row_count":0,
		"counter_overflow_count":0,
	}

	bridge:=wlmLmRawRepZeroshotAliasFamilyR5BridgeRecords(false)
	metrics["bridge_record_count"]=float64(len(bridge))
	reps:=wlmLmRawRepZeroshotAliasFamilyR1SelectSurface(bridge)
	metrics["selected_surface_representation_count"]=float64(len(reps))
	metrics["selected_true_surface_motif_match_count"]=float64(wlmLmRawRepZeroshotAliasFamilyR1TrueMotifCount(reps))
	repIndex:=wlmLmRawRepContextualAliasInductionR2RepIndex(reps)

	classes:=wlmLmRawRepContextualAliasInductionR2Induce(reps,bridge)
	metrics["induced_alias_class_count"]=float64(len(classes.members))
	metrics["correct_evaluator_alias_triplet_count"]=float64(wlmLmRawRepZeroshotAliasFamilyR1CorrectTriplets(repIndex,classes))
	if classes.decodeFailures!=0 {
		metrics["invalid_row_count"]+=float64(classes.decodeFailures)
	}

	taskRecords:=wlmLmRawRepContextualAliasInductionR2TaskRecords()
	metrics["task_training_record_count"]=float64(len(taskRecords))
	classModel:=wlmLmRawRepContextualAliasInductionR2BuildClassModel(taskRecords,reps,classes,metrics)
	if classModel.decodeFails!=0 {
		metrics["invalid_row_count"]+=float64(classModel.decodeFails)
	}
	if classModel.valid {
		metrics["class_transition_key_count"]=float64(80-classModel.transition.missingCount())
	}

	surfaceModel,surfaceDecodeFails:=wlmLmRawRepZeroshotAliasFamilyR1BuildSurfaceModel(taskRecords,reps,metrics)
	if surfaceDecodeFails!=0 {
		metrics["invalid_row_count"]+=float64(surfaceDecodeFails)
	}

	shuffledBridge:=wlmLmRawRepZeroshotAliasFamilyR5BridgeRecords(true)
	shuffledClasses:=wlmLmRawRepContextualAliasInductionR2Induce(reps,shuffledBridge)
	if shuffledClasses.decodeFailures!=0 {
		metrics["invalid_row_count"]+=float64(shuffledClasses.decodeFailures)
	}
	shuffledModel:=wlmLmRawRepContextualAliasInductionR2BuildClassModel(taskRecords,reps,shuffledClasses,metrics)
	if shuffledModel.decodeFails!=0 {
		metrics["invalid_row_count"]+=float64(shuffledModel.decodeFails)
	}

	conjHits:=0
	conjTotal:=0
	longHits:=0
	longTotal:=0
	overallHits:=0
	rawHits:=0
	exactHits:=0
	exactMissingPrograms:=0
	shuffledHits:=0
	overallTotal:=0

	evaluate:=func(left,right int,ops []int,tag int,conjugated bool){
		raw:=wlmLmRawRepZeroshotMultistepR1ProgramRecord(left,right,ops,tag)
		metrics["heldout_program_count"]++
		metrics["family2_operation_step_count"]+=float64(len(ops))
		if conjugated {
			metrics["conjugated_right_program_count"]++
		} else {
			metrics["long_program_count"]++
		}

		pl,pr,ok:=wlmLmRawRepZeroshotMultistepR1ExecuteClass(raw,ops,reps,classes,classModel,metrics)
		spl,spr,sok:=wlmLmRawRepZeroshotMultistepR1ExecuteClass(raw,ops,reps,shuffledClasses,shuffledModel,metrics)
		epl,epr,eok,missing:=wlmLmRawRepZeroshotMultistepR1ExecuteSurface(raw,ops,reps,surfaceModel,metrics)
		if !ok||!sok||!eok {
			metrics["program_decode_failure_count"]++
		}
		targetLeft,targetRight:=wlmLmRawRepRoleCompositionR1Target(left,right,ops)
		targetDenseLeft,okL:=wlmLmRawRepZeroshotAliasFamilyR1ExpectedDenseValue(targetLeft,reps,repIndex,classes,classModel)
		targetDenseRight,okR:=wlmLmRawRepZeroshotAliasFamilyR1ExpectedDenseValue(targetRight,reps,repIndex,classes,classModel)
		mainCorrect:=ok&&okL&&okR&&pl==targetDenseLeft&&pr==targetDenseRight
		if mainCorrect {
			overallHits++
			if conjugated {conjHits++} else {longHits++}
		}
		if ok&&wlmLmRawRepZeroshotMultistepR1Family2RawOutputExact(pl,pr,targetLeft,targetRight,reps,repIndex,classes,classModel) {
			rawHits++
		}
		if missing {
			exactMissingPrograms++
		}
		if eok {
			semanticLeft:=wlmLmRawRepZeroshotAliasFamilyR1SemanticOfValueKey(reps[epl])
			semanticRight:=wlmLmRawRepZeroshotAliasFamilyR1SemanticOfValueKey(reps[epr])
			if semanticLeft==targetLeft&&semanticRight==targetRight {
				exactHits++
			}
		}
		if sok {
			shTargetLeft,shOkL:=wlmLmRawRepZeroshotAliasFamilyR1ExpectedDenseValue(targetLeft,reps,repIndex,shuffledClasses,shuffledModel)
			shTargetRight,shOkR:=wlmLmRawRepZeroshotAliasFamilyR1ExpectedDenseValue(targetRight,reps,repIndex,shuffledClasses,shuffledModel)
			if shOkL&&shOkR&&spl==shTargetLeft&&spr==shTargetRight {
				shuffledHits++
			}
		}
		if conjugated {conjTotal++} else {longTotal++}
		overallTotal++
	}

	for left:=0;left<4;left++ {
		for right:=0;right<4;right++ {
			for start:=0;start<4;start++ {
				ops:=[]int{4,start,4}
				evaluate(left,right,ops,100+left*100+right*10+start,true)
			}
		}
	}
	for _,length:=range []int{5,7} {
		for left:=0;left<4;left++ {
			for right:=0;right<4;right++ {
				for seed:=0;seed<4;seed++ {
					ops:=wlmLmRawRepRoleCompositionR1LongProgram(left,right,seed,length)
					evaluate(left,right,ops,1000+length*100+left*20+right*4+seed,false)
				}
			}
		}
	}

	if conjTotal>0 {
		metrics["induced_conjugated_right_accuracy"]=float64(conjHits)/float64(conjTotal)
	}
	if longTotal>0 {
		metrics["induced_long_program_accuracy"]=float64(longHits)/float64(longTotal)
	}
	if overallTotal>0 {
		metrics["induced_overall_multistep_accuracy"]=float64(overallHits)/float64(overallTotal)
		metrics["raw_output_pair_exact_accuracy"]=float64(rawHits)/float64(overallTotal)
		metrics["exact_surface_overall_multistep_accuracy"]=float64(exactHits)/float64(overallTotal)
		metrics["shuffled_bridge_overall_multistep_accuracy"]=float64(shuffledHits)/float64(overallTotal)
	}
	metrics["exact_surface_programs_with_missing_transition"]=float64(exactMissingPrograms)

	for _,value:=range metrics {
		if math.IsNaN(value)||math.IsInf(value,0) {
			metrics["invalid_row_count"]++
		}
	}

	return wlmLmRawRepZeroshotMultistepR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-RAW-REP-ZEROSHOT-MULTISTEP-R1",
		Metrics:metrics,
	}
}
