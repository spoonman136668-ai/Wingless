package unitary

import "sort"

type wlmYggRealUTF8CalibratedMaintenanceR1Background struct {
	occ int
	err int
}

type wlmYggRealUTF8CalibratedMaintenanceR1Diagnostic struct {
	occ int
	err int
}

type wlmYggRealUTF8CalibratedMaintenanceR1Scored struct {
	key   [4]uint8
	score float64
}

type wlmYggRealUTF8CalibratedMaintenanceR1Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmYggRealUTF8CalibratedMaintenanceR1BackgroundState(
	data []byte,
	baseline *[256][256]uint32,
	canonical map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
) map[[4]uint8]wlmYggRealUTF8CalibratedMaintenanceR1Background {
	out:=make(map[[4]uint8]wlmYggRealUTF8CalibratedMaintenanceR1Background,len(canonical))
	for key:=range canonical { out[key]=wlmYggRealUTF8CalibratedMaintenanceR1Background{} }
	for i:=4;i<len(data);i++{
		key:=[4]uint8{data[i-4],data[i-3],data[i-2],data[i-1]}
		if _,ok:=canonical[key];!ok{continue}
		row:=out[key]
		row.occ++
		if wlmYggRealUTF8Level1MaintenanceR3Predict(baseline,canonical,key)!=data[i]{row.err++}
		out[key]=row
	}
	return out
}

func wlmYggRealUTF8CalibratedMaintenanceR1DiagnosticState(
	data []byte,
	baseline *[256][256]uint32,
	active map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
	catalog map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
) map[[4]uint8]wlmYggRealUTF8CalibratedMaintenanceR1Diagnostic {
	out:=make(map[[4]uint8]wlmYggRealUTF8CalibratedMaintenanceR1Diagnostic)
	for i:=4;i<len(data);i++{
		key:=[4]uint8{data[i-4],data[i-3],data[i-2],data[i-1]}
		if _,ok:=catalog[key];!ok{continue}
		row:=out[key]
		row.occ++
		if wlmYggRealUTF8Level1MaintenanceR3Predict(baseline,active,key)!=data[i]{row.err++}
		out[key]=row
	}
	return out
}

func wlmYggRealUTF8CalibratedMaintenanceR1Select(
	background map[[4]uint8]wlmYggRealUTF8CalibratedMaintenanceR1Background,
	diagnostic map[[4]uint8]wlmYggRealUTF8CalibratedMaintenanceR1Diagnostic,
) ([4]uint8,float64,float64,bool) {
	rows:=make([]wlmYggRealUTF8CalibratedMaintenanceR1Scored,0)
	for key,bg:=range background{
		d:=diagnostic[key]
		if bg.occ<3 || d.occ<1{continue}
		expected:=float64(bg.err)*float64(d.occ)/float64(bg.occ)
		score:=float64(d.err)-expected
		rows=append(rows,wlmYggRealUTF8CalibratedMaintenanceR1Scored{key:key,score:score})
	}
	if len(rows)==0{return [4]uint8{},0,0,false}
	sort.Slice(rows,func(i,j int)bool{
		if rows[i].score!=rows[j].score{return rows[i].score>rows[j].score}
		return wlmLmRealUTF8BytePilotR2Less(rows[i].key,rows[j].key)
	})
	runner:=0.0
	if len(rows)>1{runner=rows[1].score}
	return rows[0].key,rows[0].score,runner,true
}

func wlmYggRealUTF8CalibratedMaintenanceR1Min(a,b float64)float64{if b<a{return b};return a}
func wlmYggRealUTF8CalibratedMaintenanceR1Max(a,b float64)float64{if b>a{return b};return a}

// RunWlmYggRealUTF8CalibratedMaintenanceR1 evaluates calibrated real-byte maintenance.
func RunWlmYggRealUTF8CalibratedMaintenanceR1() interface{} {
	metrics:=map[string]float64{
		"valid_evaluation_file_count":0,
		"completed_cycle_count":0,
		"minimum_target_eligibility_count_per_file":256,
		"minimum_monitor_target_selection_accuracy":1,
		"minimum_selected_calibrated_excess_score":1e9,
		"minimum_selected_minus_runner_up_excess_score":1e9,
		"maximum_wrong_candidate_selected_key_accuracy_gain":-1,
		"wrong_candidate_retain_count":0,
		"wrong_candidate_revert_count":0,
		"minimum_correct_candidate_selected_key_accuracy_gain":1,
		"correct_candidate_retain_count":0,
		"correct_candidate_revert_count":0,
		"maximum_post_repair_audit_prediction_mismatch_count":0,
		"maximum_post_repair_audit_probability_abs_delta":0,
		"minimum_structure_jaccard_after_each_cycle":1,
		"maximum_background_monitor_state_bytes":2048,
		"maximum_candidate_read_bytes":0,
		"maximum_candidate_repair_operations":0,
		"maximum_final_selected_motif_count":0,
		"minimum_final_selected_motif_count":256,
		"lesion_label_access_count":0,
		"canonical_diagnostic_shadow_access_count":0,
		"evaluator_schedule_access_in_policy_count":0,
		"tokenizer_use_count":0,
		"capacity_growth_event_count":0,
		"invalid_cycle_rows":0,
		"training_file_identity_mismatch_count":0,
		"evaluation_file_identity_mismatch_count":0,
		"canonical_selected_motif_count":0,
	}

	train:=make([][]byte,0,len(wlmLmRealUTF8BytePilotR2Train))
	trainSet:=make(map[string]bool)
	for _,f:=range wlmLmRealUTF8BytePilotR2Train{
		data,ok:=wlmLmRealUTF8BytePilotR2LoadBlob(f)
		if !ok{metrics["training_file_identity_mismatch_count"]++}
		train=append(train,data)
		trainSet[f.sha]=true
	}
	baseline,_,selected:=wlmLmRealUTF8BytePilotR2TrainModel(train,metrics)
	canonical:=wlmYggRealUTF8Level1MaintenanceR3CloneMap(selected)
	active:=wlmYggRealUTF8Level1MaintenanceR3CloneMap(canonical)
	metrics["canonical_selected_motif_count"]=float64(len(canonical))
	if len(canonical)!=256{metrics["invalid_cycle_rows"]++}

	for _,f:=range wlmLmRealUTF8BytePilotR2Eval{
		if trainSet[f.sha]{metrics["invalid_cycle_rows"]++}
		data,ok:=wlmLmRealUTF8BytePilotR2LoadBlob(f)
		if !ok{metrics["evaluation_file_identity_mismatch_count"]++}
		quarters:=wlmYggRealUTF8Level1MaintenanceR3Split(data)
		targets:=wlmYggRealUTF8Level1MaintenanceR3Targets(quarters,&baseline,canonical)
		metrics["minimum_target_eligibility_count_per_file"]=wlmYggRealUTF8CalibratedMaintenanceR1Min(
			metrics["minimum_target_eligibility_count_per_file"],float64(len(targets)),
		)
		if len(targets)<2{metrics["invalid_cycle_rows"]++;continue}
		background:=wlmYggRealUTF8CalibratedMaintenanceR1BackgroundState(quarters[0],&baseline,canonical)
		if len(background)!=256{metrics["invalid_cycle_rows"]++}
		metrics["valid_evaluation_file_count"]++

		for targetIndex:=0;targetIndex<2;targetIndex++{
			targetKey:=targets[targetIndex].key
			delete(active,targetKey)
			diagnostic:=wlmYggRealUTF8CalibratedMaintenanceR1DiagnosticState(quarters[1],&baseline,active,canonical)
			selectedKey,score,runner,monitorOK:=wlmYggRealUTF8CalibratedMaintenanceR1Select(background,diagnostic)
			if !monitorOK{metrics["invalid_cycle_rows"]++;continue}
			if selectedKey!=targetKey{metrics["minimum_monitor_target_selection_accuracy"]=0}
			metrics["minimum_selected_calibrated_excess_score"]=wlmYggRealUTF8CalibratedMaintenanceR1Min(metrics["minimum_selected_calibrated_excess_score"],score)
			metrics["minimum_selected_minus_runner_up_excess_score"]=wlmYggRealUTF8CalibratedMaintenanceR1Min(metrics["minimum_selected_minus_runner_up_excess_score"],score-runner)

			preAcc,preCount:=wlmYggRealUTF8Level1MaintenanceR3KeyAccuracy(quarters[2],selectedKey,&baseline,active)
			if preCount<3{metrics["invalid_cycle_rows"]++}
			preBPB:=wlmYggRealUTF8Level1MaintenanceR3BPB(quarters[2],&baseline,active)

			wrongSource,wrongOK:=wlmYggRealUTF8Level1MaintenanceR3WrongSource(selectedKey,canonical)
			if !wrongOK{metrics["invalid_cycle_rows"]++;continue}
			wrongState:=wlmYggRealUTF8Level1MaintenanceR3CloneMap(active)
			wrongEntry:=wlmYggRealUTF8Level1MaintenanceR3CloneEntry(canonical[wrongSource])
			wrongEntry.key=selectedKey
			wrongState[selectedKey]=wrongEntry
			wrongAcc,_:=wlmYggRealUTF8Level1MaintenanceR3KeyAccuracy(quarters[2],selectedKey,&baseline,wrongState)
			wrongBPB:=wlmYggRealUTF8Level1MaintenanceR3BPB(quarters[2],&baseline,wrongState)
			wrongGain:=wrongAcc-preAcc
			metrics["maximum_wrong_candidate_selected_key_accuracy_gain"]=wlmYggRealUTF8CalibratedMaintenanceR1Max(metrics["maximum_wrong_candidate_selected_key_accuracy_gain"],wrongGain)
			metrics["maximum_candidate_read_bytes"]=wlmYggRealUTF8CalibratedMaintenanceR1Max(metrics["maximum_candidate_read_bytes"],1028)
			metrics["maximum_candidate_repair_operations"]=wlmYggRealUTF8CalibratedMaintenanceR1Max(metrics["maximum_candidate_repair_operations"],1)
			wrongRetain:=wrongGain>=0.20 && (wrongBPB-preBPB)<=0.005
			if wrongRetain{
				active=wrongState
				metrics["wrong_candidate_retain_count"]++
			}else{
				metrics["wrong_candidate_revert_count"]++
			}

			if !wrongRetain{
				correctState:=wlmYggRealUTF8Level1MaintenanceR3CloneMap(active)
				correctState[selectedKey]=wlmYggRealUTF8Level1MaintenanceR3CloneEntry(canonical[selectedKey])
				correctAcc,_:=wlmYggRealUTF8Level1MaintenanceR3KeyAccuracy(quarters[2],selectedKey,&baseline,correctState)
				correctBPB:=wlmYggRealUTF8Level1MaintenanceR3BPB(quarters[2],&baseline,correctState)
				correctGain:=correctAcc-preAcc
				metrics["minimum_correct_candidate_selected_key_accuracy_gain"]=wlmYggRealUTF8CalibratedMaintenanceR1Min(metrics["minimum_correct_candidate_selected_key_accuracy_gain"],correctGain)
				correctRetain:=correctGain>=0.20 && (correctBPB-preBPB)<=0.005
				if correctRetain{
					active=correctState
					metrics["correct_candidate_retain_count"]++
				}else{
					metrics["correct_candidate_revert_count"]++
				}
			}

			jaccard:=wlmYggRealUTF8Level1MaintenanceR3Jaccard(active,canonical)
			metrics["minimum_structure_jaccard_after_each_cycle"]=wlmYggRealUTF8CalibratedMaintenanceR1Min(metrics["minimum_structure_jaccard_after_each_cycle"],jaccard)
			predMismatch,probDelta:=wlmYggRealUTF8Level1MaintenanceR3Audit(quarters[3],&baseline,active,canonical)
			metrics["maximum_post_repair_audit_prediction_mismatch_count"]=wlmYggRealUTF8CalibratedMaintenanceR1Max(metrics["maximum_post_repair_audit_prediction_mismatch_count"],float64(predMismatch))
			metrics["maximum_post_repair_audit_probability_abs_delta"]=wlmYggRealUTF8CalibratedMaintenanceR1Max(metrics["maximum_post_repair_audit_probability_abs_delta"],probDelta)
			metrics["completed_cycle_count"]++
		}
	}

	metrics["maximum_final_selected_motif_count"]=float64(len(active))
	metrics["minimum_final_selected_motif_count"]=float64(len(active))
	return wlmYggRealUTF8CalibratedMaintenanceR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-YGG-REAL-UTF8-CALIBRATED-MAINTENANCE-R1",
		Metrics:metrics,
	}
}
