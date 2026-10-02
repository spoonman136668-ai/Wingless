package unitary

import "math"

type wlmLmRawRepContextualAliasMultistepR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmRawRepContextualAliasMultistepR1ProgramRecord(left,right int,ops []int,valueFamily,tag int) []uint8 {
	state:=uint32(0x7f4a7c15 ^ uint32((left+1)*131+(right+1)*977+(tag+1)*65537+len(ops)*8191))
	out:=make([]uint8,0,140)
	out=append(out,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,6+(tag%5))...)
	lk:=wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(valueFamily,left)
	out=append(out,lk[:]...)
	out=append(out,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,8+((left+tag)%5))...)
	rk:=wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(valueFamily,right)
	out=append(out,rk[:]...)
	for step,op:=range ops {
		out=append(out,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,7+((step+op+tag)%7))...)
		opFamily:=step%2
		ok:=wlmSiRawRepAliasInvarianceFalsificationR1OpMotif(opFamily,op)
		out=append(out,ok[:]...)
	}
	out=append(out,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,9+((tag+len(ops))%7))...)
	return out
}

func wlmLmRawRepContextualAliasMultistepR1ExecuteClass(
	raw []uint8,
	leftSemantic,rightSemantic int,
	ops []int,
	valueFamily int,
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	repIndex map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int,
	classes wlmLmRawRepContextualAliasInductionR2Classes,
	model wlmLmRawRepContextualAliasInductionR2ClassModel,
	metrics map[string]float64,
)(bool,bool) {
	surface:=wlmSiRawRepAliasInvarianceFalsificationR1Decode(raw,reps)
	if len(surface)!=2+len(ops) {
		return false,false
	}
	if surface[0]<0||surface[0]>=len(classes.surfaceToClass)||surface[1]<0||surface[1]>=len(classes.surfaceToClass) {
		metrics["invalid_row_count"]++
		return false,true
	}
	leftClass:=classes.surfaceToClass[surface[0]]
	rightClass:=classes.surfaceToClass[surface[1]]
	left,ok0:=model.valueMap[leftClass]
	right,ok1:=model.valueMap[rightClass]
	if !ok0||!ok1 {
		metrics["invalid_row_count"]++
		return false,true
	}
	for step:=range ops {
		surfaceOp:=surface[2+step]
		if surfaceOp<0||surfaceOp>=len(classes.surfaceToClass) {
			metrics["invalid_row_count"]++
			return false,true
		}
		classOp:=classes.surfaceToClass[surfaceOp]
		op,ok:=model.opMap[classOp]
		if !ok {
			metrics["invalid_row_count"]++
			return false,true
		}
		left,right=model.transition.predict(op,left,right)
	}
	targetLeftSemantic,targetRightSemantic:=wlmLmRawRepRoleCompositionR1Target(leftSemantic,rightSemantic,ops)
	targetLeft,ok2:=wlmLmRawRepContextualAliasInductionR2ExpectedDenseValue(valueFamily,targetLeftSemantic,reps,repIndex,classes,model)
	targetRight,ok3:=wlmLmRawRepContextualAliasInductionR2ExpectedDenseValue(valueFamily,targetRightSemantic,reps,repIndex,classes,model)
	if !ok2||!ok3 {
		metrics["invalid_row_count"]++
		return false,true
	}
	return left==targetLeft&&right==targetRight,true
}

func wlmLmRawRepContextualAliasMultistepR1ExecuteSurface(
	raw []uint8,
	leftSemantic,rightSemantic int,
	ops []int,
	valueFamily int,
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	repIndex map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int,
	model wlmLmRawRepContextualAliasInductionR2SurfaceModel,
	metrics map[string]float64,
)(bool,bool,bool) {
	surface:=wlmSiRawRepAliasInvarianceFalsificationR1Decode(raw,reps)
	if len(surface)!=2+len(ops) {
		return false,false,false
	}
	left,ok0:=model.valueMap[surface[0]]
	right,ok1:=model.valueMap[surface[1]]
	if !ok0||!ok1 {
		metrics["invalid_row_count"]++
		return false,false,true
	}
	hadMissing:=false
	for step:=range ops {
		op,ok:=model.opMap[surface[2+step]]
		if !ok {
			metrics["invalid_row_count"]++
			return false,hadMissing,true
		}
		if !model.transition.has(op,left,right) {
			hadMissing=true
		}
		left,right=model.transition.predict(op,left,right)
	}
	targetLeftSemantic,targetRightSemantic:=wlmLmRawRepRoleCompositionR1Target(leftSemantic,rightSemantic,ops)
	targetLeftSurface,ok2:=repIndex[wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(valueFamily,targetLeftSemantic)]
	targetRightSurface,ok3:=repIndex[wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(valueFamily,targetRightSemantic)]
	if !ok2||!ok3 {
		metrics["invalid_row_count"]++
		return false,hadMissing,true
	}
	targetLeft,ok4:=model.valueMap[targetLeftSurface]
	targetRight,ok5:=model.valueMap[targetRightSurface]
	if !ok4||!ok5 {
		metrics["invalid_row_count"]++
		return false,hadMissing,true
	}
	return left==targetLeft&&right==targetRight,hadMissing,true
}

func wlmLmRawRepContextualAliasMultistepR1EvaluateProgram(
	left,right int,
	ops []int,
	tag int,
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	repIndex map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int,
	classes wlmLmRawRepContextualAliasInductionR2Classes,
	classModel wlmLmRawRepContextualAliasInductionR2ClassModel,
	surfaceModel wlmLmRawRepContextualAliasInductionR2SurfaceModel,
	shuffledClasses wlmLmRawRepContextualAliasInductionR2Classes,
	shuffledModel wlmLmRawRepContextualAliasInductionR2ClassModel,
	metrics map[string]float64,
)(bool,bool,bool,bool,bool) {
	valueFamily:=(left+right+tag)%2
	raw:=wlmLmRawRepContextualAliasMultistepR1ProgramRecord(left,right,ops,valueFamily,tag)
	induced,decoded:=wlmLmRawRepContextualAliasMultistepR1ExecuteClass(
		raw,left,right,ops,valueFamily,reps,repIndex,classes,classModel,metrics,
	)
	exact,missing,decodedSurface:=wlmLmRawRepContextualAliasMultistepR1ExecuteSurface(
		raw,left,right,ops,valueFamily,reps,repIndex,surfaceModel,metrics,
	)
	shuffled,decodedShuffled:=wlmLmRawRepContextualAliasMultistepR1ExecuteClass(
		raw,left,right,ops,valueFamily,reps,repIndex,shuffledClasses,shuffledModel,metrics,
	)
	return induced,exact,shuffled,missing,decoded&&decodedSurface&&decodedShuffled
}

// RunWlmLmRawRepContextualAliasMultistepR1 tests whether R2 alias classes survive unseen multi-step composition.
func RunWlmLmRawRepContextualAliasMultistepR1() interface{} {
	metrics:=map[string]float64{
		"training_record_count":0,
		"bridge_record_count":0,
		"selected_surface_representation_count":0,
		"induced_alias_class_count":0,
		"correct_evaluator_alias_pair_count":0,
		"heldout_program_count":0,
		"mixed_alias_program_count":0,
		"program_decode_failure_count":0,
		"conjugated_right_program_count":0,
		"induced_conjugated_right_accuracy":0,
		"long_program_count":0,
		"induced_long_program_accuracy":0,
		"induced_overall_multistep_accuracy":0,
		"exact_surface_overall_multistep_accuracy":0,
		"exact_surface_programs_with_missing_transition":0,
		"shuffled_bridge_overall_multistep_accuracy":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"capacity_growth_event_count":0,
		"invalid_row_count":0,
		"counter_overflow_count":0,
	}

	task:=wlmLmRawRepContextualAliasInductionR2TaskRecords()
	metrics["training_record_count"]=float64(len(task))
	reps:=wlmSiRawRepAliasInvarianceFalsificationR1LearnRepresentation(task)
	metrics["selected_surface_representation_count"]=float64(len(reps))
	repIndex:=wlmLmRawRepContextualAliasInductionR2RepIndex(reps)

	bridge:=wlmLmRawRepContextualAliasInductionR2BridgeRecords(false)
	metrics["bridge_record_count"]=float64(len(bridge))
	classes:=wlmLmRawRepContextualAliasInductionR2Induce(reps,bridge)
	metrics["induced_alias_class_count"]=float64(len(classes.members))
	metrics["correct_evaluator_alias_pair_count"]=float64(wlmLmRawRepContextualAliasInductionR2CorrectAliasPairs(reps,repIndex,classes))
	classModel:=wlmLmRawRepContextualAliasInductionR2BuildClassModel(task,reps,classes,metrics)
	surfaceModel:=wlmLmRawRepContextualAliasInductionR2BuildSurfaceModel(task,reps,metrics)

	shuffledBridge:=wlmLmRawRepContextualAliasInductionR2BridgeRecords(true)
	shuffledClasses:=wlmLmRawRepContextualAliasInductionR2Induce(reps,shuffledBridge)
	shuffledModel:=wlmLmRawRepContextualAliasInductionR2BuildClassModel(task,reps,shuffledClasses,metrics)

	if classes.decodeFailures!=0||classModel.decodeFails!=0||shuffledClasses.decodeFailures!=0||shuffledModel.decodeFails!=0 {
		metrics["invalid_row_count"]+=float64(classes.decodeFailures+classModel.decodeFails+shuffledClasses.decodeFailures+shuffledModel.decodeFails)
	}

	inducedHits:=0
	exactHits:=0
	shuffledHits:=0
	overall:=0
	conjHits:=0
	conjTotal:=0
	longHits:=0
	longTotal:=0

	for left:=0;left<4;left++ {
		for right:=0;right<4;right++ {
			for start:=0;start<4;start++ {
				ops:=[]int{4,start,4}
				tag:=100+left*100+right*10+start
				i,e,s,m,d:=wlmLmRawRepContextualAliasMultistepR1EvaluateProgram(
					left,right,ops,tag,reps,repIndex,classes,classModel,surfaceModel,shuffledClasses,shuffledModel,metrics,
				)
				metrics["heldout_program_count"]++
				metrics["conjugated_right_program_count"]++
				metrics["mixed_alias_program_count"]++
				if !d { metrics["program_decode_failure_count"]++ }
				if i { inducedHits++;conjHits++ }
				if e { exactHits++ }
				if s { shuffledHits++ }
				if m { metrics["exact_surface_programs_with_missing_transition"]++ }
				overall++;conjTotal++
			}
		}
	}

	for _,length:=range []int{5,7} {
		for left:=0;left<4;left++ {
			for right:=0;right<4;right++ {
				for seed:=0;seed<4;seed++ {
					ops:=wlmLmRawRepRoleCompositionR1LongProgram(left,right,seed,length)
					tag:=1000+length*100+left*20+right*4+seed
					i,e,s,m,d:=wlmLmRawRepContextualAliasMultistepR1EvaluateProgram(
						left,right,ops,tag,reps,repIndex,classes,classModel,surfaceModel,shuffledClasses,shuffledModel,metrics,
					)
					metrics["heldout_program_count"]++
					metrics["long_program_count"]++
					metrics["mixed_alias_program_count"]++
					if !d { metrics["program_decode_failure_count"]++ }
					if i { inducedHits++;longHits++ }
					if e { exactHits++ }
					if s { shuffledHits++ }
					if m { metrics["exact_surface_programs_with_missing_transition"]++ }
					overall++;longTotal++
				}
			}
		}
	}

	if conjTotal>0 { metrics["induced_conjugated_right_accuracy"]=float64(conjHits)/float64(conjTotal) }
	if longTotal>0 { metrics["induced_long_program_accuracy"]=float64(longHits)/float64(longTotal) }
	if overall>0 {
		metrics["induced_overall_multistep_accuracy"]=float64(inducedHits)/float64(overall)
		metrics["exact_surface_overall_multistep_accuracy"]=float64(exactHits)/float64(overall)
		metrics["shuffled_bridge_overall_multistep_accuracy"]=float64(shuffledHits)/float64(overall)
	}
	for _,v:=range metrics {
		if math.IsNaN(v)||math.IsInf(v,0) { metrics["invalid_row_count"]++ }
	}
	return wlmLmRawRepContextualAliasMultistepR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-RAW-REP-CONTEXTUAL-ALIAS-MULTISTEP-R1",
		Metrics:metrics,
	}
}
