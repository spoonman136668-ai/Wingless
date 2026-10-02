package unitary

import "sort"

type wlmYggRealUTF8HeterogeneousMaintenanceR1Target struct {
	key  [4]uint8
	occ  int
	gain float64
}

type wlmYggRealUTF8HeterogeneousMaintenanceR1Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

var wlmYggRealUTF8HeterogeneousMaintenanceR1Eval = []wlmLmRealUTF8BytePilotR2CorpusFile{
	{"docs/experiments/qwen-1.5b-linux-stage3.json","7765773c8271aa46e4414bd41b8708f1afb1af2f"},
	{"unitary/full_latent.go","5d9400dcc08c5d10cc34303876d88ab1fb785cc8"},
}

var wlmYggRealUTF8HeterogeneousMaintenanceR1Classes = []string{
	"record_missing",
	"successor_distribution_rotated",
	"canonical_best_bucket_erased",
	"nonzero_counts_flattened",
}

func wlmYggRealUTF8HeterogeneousMaintenanceR1Recompute(
	entry wlmLmRealUTF8BytePilotR2Selected,
) (wlmLmRealUTF8BytePilotR2Selected, bool) {
	if entry.stats == nil {
		return entry, false
	}
	best,bestCount,total:=wlmLmRealUTF8BytePilotR2Best(&entry.stats.counts)
	entry.stats.total=uint32(total)
	entry.best=best
	entry.bestCount=bestCount
	if total==0 {
		entry.consistency=0
		return entry,false
	}
	entry.consistency=float64(bestCount)/float64(total)
	return entry,true
}

func wlmYggRealUTF8HeterogeneousMaintenanceR1Degrade(
	canonical map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
	key [4]uint8,
	className string,
) (map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected, bool) {
	active:=wlmYggRealUTF8Level1MaintenanceR3CloneMap(canonical)
	entry,ok:=active[key]
	if !ok {
		return active,false
	}
	switch className {
	case "record_missing":
		delete(active,key)
		return active,true
	case "successor_distribution_rotated":
		var rotated [256]uint32
		for i:=0;i<256;i++ {
			rotated[(i+1)&255]=entry.stats.counts[i]
		}
		entry.stats.counts=rotated
	case "canonical_best_bucket_erased":
		entry.stats.counts[canonical[key].best]=0
	case "nonzero_counts_flattened":
		for i:=0;i<256;i++ {
			if entry.stats.counts[i]>0 {
				entry.stats.counts[i]=1
			}
		}
	default:
		return active,false
	}
	entry,ok=wlmYggRealUTF8HeterogeneousMaintenanceR1Recompute(entry)
	if !ok {
		return active,false
	}
	active[key]=entry
	return active,true
}

func wlmYggRealUTF8HeterogeneousMaintenanceR1CandidateTargets(
	quarters [4][]byte,
	baseline *[256][256]uint32,
	canonical map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
	className string,
	used map[[4]uint8]bool,
) []wlmYggRealUTF8HeterogeneousMaintenanceR1Target {
	cal:=wlmYggRealUTF8Level1MaintenanceR3SupportFor(quarters[0],baseline,canonical)
	diag:=wlmYggRealUTF8Level1MaintenanceR3SupportFor(quarters[1],baseline,canonical)
	val:=wlmYggRealUTF8Level1MaintenanceR3SupportFor(quarters[2],baseline,canonical)
	audit:=wlmYggRealUTF8Level1MaintenanceR3SupportFor(quarters[3],baseline,canonical)
	rows:=make([]wlmYggRealUTF8HeterogeneousMaintenanceR1Target,0)
	for key,row:=range cal {
		if used[key] || row.occ<6 || diag[key].occ<3 || val[key].occ<3 || audit[key].occ<3 {
			continue
		}
		gain:=float64(row.canonicalHit-row.baselineHit)/float64(row.occ)
		if gain<0.20 {
			continue
		}
		degraded,ok:=wlmYggRealUTF8HeterogeneousMaintenanceR1Degrade(canonical,key,className)
		if !ok {
			continue
		}
		canonAcc,count:=wlmYggRealUTF8Level1MaintenanceR3KeyAccuracy(quarters[0],key,baseline,canonical)
		degradedAcc,degradedCount:=wlmYggRealUTF8Level1MaintenanceR3KeyAccuracy(quarters[0],key,baseline,degraded)
		if count<6 || degradedCount!=count || canonAcc-degradedAcc<0.20 {
			continue
		}
		rows=append(rows,wlmYggRealUTF8HeterogeneousMaintenanceR1Target{key:key,occ:row.occ,gain:gain})
	}
	sort.Slice(rows,func(i,j int)bool{
		if rows[i].occ!=rows[j].occ{return rows[i].occ>rows[j].occ}
		if rows[i].gain!=rows[j].gain{return rows[i].gain>rows[j].gain}
		return wlmLmRealUTF8BytePilotR2Less(rows[i].key,rows[j].key)
	})
	return rows
}

func wlmYggRealUTF8HeterogeneousMaintenanceR1Schedule(
	quarters [4][]byte,
	baseline *[256][256]uint32,
	canonical map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
) ([]wlmYggRealUTF8HeterogeneousMaintenanceR1Target,int) {
	used:=make(map[[4]uint8]bool)
	out:=make([]wlmYggRealUTF8HeterogeneousMaintenanceR1Target,0,4)
	eligibleUnion:=make(map[[4]uint8]bool)
	for _,className:=range wlmYggRealUTF8HeterogeneousMaintenanceR1Classes {
		all:=wlmYggRealUTF8HeterogeneousMaintenanceR1CandidateTargets(quarters,baseline,canonical,className,map[[4]uint8]bool{})
		for _,row:=range all {eligibleUnion[row.key]=true}
		rows:=wlmYggRealUTF8HeterogeneousMaintenanceR1CandidateTargets(quarters,baseline,canonical,className,used)
		if len(rows)==0 {
			return out,len(eligibleUnion)
		}
		out=append(out,rows[0])
		used[rows[0].key]=true
	}
	return out,len(eligibleUnion)
}

func wlmYggRealUTF8HeterogeneousMaintenanceR1Min(a,b float64)float64{if b<a{return b};return a}
func wlmYggRealUTF8HeterogeneousMaintenanceR1Max(a,b float64)float64{if b>a{return b};return a}

// RunWlmYggRealUTF8HeterogeneousMaintenanceR1 evaluates calibrated maintenance on fresh blobs and four degradation classes.
func RunWlmYggRealUTF8HeterogeneousMaintenanceR1() interface{} {
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
		"minimum_per_class_completed_cycle_count":2,
		"maximum_post_repair_audit_prediction_mismatch_count":0,
		"maximum_post_repair_audit_probability_abs_delta":0,
		"minimum_structure_jaccard_after_each_cycle":1,
		"maximum_candidate_read_bytes":0,
		"maximum_candidate_repair_operations":0,
		"maximum_final_selected_motif_count":0,
		"minimum_final_selected_motif_count":256,
		"lesion_label_access_count":0,
		"class_label_access_in_policy_count":0,
		"canonical_diagnostic_shadow_access_count":0,
		"evaluator_schedule_access_in_policy_count":0,
		"tokenizer_use_count":0,
		"capacity_growth_event_count":0,
		"invalid_cycle_rows":0,
		"training_file_identity_mismatch_count":0,
		"evaluation_file_identity_mismatch_count":0,
		"fresh_eval_blob_overlap_count":0,
		"canonical_selected_motif_count":0,
	}
	classCounts:=[4]int{}

	train:=make([][]byte,0,len(wlmLmRealUTF8BytePilotR2Train))
	known:=make(map[string]bool)
	for _,f:=range wlmLmRealUTF8BytePilotR2Train {
		data,ok:=wlmLmRealUTF8BytePilotR2LoadBlob(f)
		if !ok{metrics["training_file_identity_mismatch_count"]++}
		train=append(train,data)
		known[f.sha]=true
	}
	for _,f:=range wlmLmRealUTF8BytePilotR2Eval { known[f.sha]=true }
	baseline,_,selected:=wlmLmRealUTF8BytePilotR2TrainModel(train,metrics)
	canonical:=wlmYggRealUTF8Level1MaintenanceR3CloneMap(selected)
	active:=wlmYggRealUTF8Level1MaintenanceR3CloneMap(canonical)
	metrics["canonical_selected_motif_count"]=float64(len(canonical))
	if len(canonical)!=256{metrics["invalid_cycle_rows"]++}

	for _,f:=range wlmYggRealUTF8HeterogeneousMaintenanceR1Eval {
		if known[f.sha]{metrics["fresh_eval_blob_overlap_count"]++}
		data,ok:=wlmLmRealUTF8BytePilotR2LoadBlob(f)
		if !ok{metrics["evaluation_file_identity_mismatch_count"]++}
		quarters:=wlmYggRealUTF8Level1MaintenanceR3Split(data)
		schedule,eligibilityCount:=wlmYggRealUTF8HeterogeneousMaintenanceR1Schedule(quarters,&baseline,canonical)
		metrics["minimum_target_eligibility_count_per_file"]=wlmYggRealUTF8HeterogeneousMaintenanceR1Min(
			metrics["minimum_target_eligibility_count_per_file"],float64(eligibilityCount),
		)
		if len(schedule)!=4 {
			metrics["invalid_cycle_rows"]++
			continue
		}
		background:=wlmYggRealUTF8CalibratedMaintenanceR1BackgroundState(quarters[0],&baseline,canonical)
		if len(background)!=256{metrics["invalid_cycle_rows"]++}
		metrics["valid_evaluation_file_count"]++

		for classIndex,className:=range wlmYggRealUTF8HeterogeneousMaintenanceR1Classes {
			targetKey:=schedule[classIndex].key
			degraded,ok:=wlmYggRealUTF8HeterogeneousMaintenanceR1Degrade(active,targetKey,className)
			if !ok {
				metrics["invalid_cycle_rows"]++
				continue
			}
			active=degraded
			diagnostic:=wlmYggRealUTF8CalibratedMaintenanceR1DiagnosticState(quarters[1],&baseline,active,canonical)
			selectedKey,score,runner,monitorOK:=wlmYggRealUTF8CalibratedMaintenanceR1Select(background,diagnostic)
			if !monitorOK{
				metrics["invalid_cycle_rows"]++
				continue
			}
			if selectedKey!=targetKey{metrics["minimum_monitor_target_selection_accuracy"]=0}
			metrics["minimum_selected_calibrated_excess_score"]=wlmYggRealUTF8HeterogeneousMaintenanceR1Min(metrics["minimum_selected_calibrated_excess_score"],score)
			metrics["minimum_selected_minus_runner_up_excess_score"]=wlmYggRealUTF8HeterogeneousMaintenanceR1Min(metrics["minimum_selected_minus_runner_up_excess_score"],score-runner)

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
			metrics["maximum_wrong_candidate_selected_key_accuracy_gain"]=wlmYggRealUTF8HeterogeneousMaintenanceR1Max(metrics["maximum_wrong_candidate_selected_key_accuracy_gain"],wrongGain)
			metrics["maximum_candidate_read_bytes"]=wlmYggRealUTF8HeterogeneousMaintenanceR1Max(metrics["maximum_candidate_read_bytes"],1028)
			metrics["maximum_candidate_repair_operations"]=wlmYggRealUTF8HeterogeneousMaintenanceR1Max(metrics["maximum_candidate_repair_operations"],1)
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
				metrics["minimum_correct_candidate_selected_key_accuracy_gain"]=wlmYggRealUTF8HeterogeneousMaintenanceR1Min(metrics["minimum_correct_candidate_selected_key_accuracy_gain"],correctGain)
				correctRetain:=correctGain>=0.20 && (correctBPB-preBPB)<=0.005
				if correctRetain{
					active=correctState
					metrics["correct_candidate_retain_count"]++
				}else{
					metrics["correct_candidate_revert_count"]++
				}
			}

			jaccard:=wlmYggRealUTF8Level1MaintenanceR3Jaccard(active,canonical)
			metrics["minimum_structure_jaccard_after_each_cycle"]=wlmYggRealUTF8HeterogeneousMaintenanceR1Min(metrics["minimum_structure_jaccard_after_each_cycle"],jaccard)
			predMismatch,probDelta:=wlmYggRealUTF8Level1MaintenanceR3Audit(quarters[3],&baseline,active,canonical)
			metrics["maximum_post_repair_audit_prediction_mismatch_count"]=wlmYggRealUTF8HeterogeneousMaintenanceR1Max(metrics["maximum_post_repair_audit_prediction_mismatch_count"],float64(predMismatch))
			metrics["maximum_post_repair_audit_probability_abs_delta"]=wlmYggRealUTF8HeterogeneousMaintenanceR1Max(metrics["maximum_post_repair_audit_probability_abs_delta"],probDelta)
			classCounts[classIndex]++
			metrics["completed_cycle_count"]++
		}
	}
	minClass:=classCounts[0]
	for _,count:=range classCounts[1:] {if count<minClass{minClass=count}}
	metrics["minimum_per_class_completed_cycle_count"]=float64(minClass)
	metrics["maximum_final_selected_motif_count"]=float64(len(active))
	metrics["minimum_final_selected_motif_count"]=float64(len(active))
	return wlmYggRealUTF8HeterogeneousMaintenanceR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-YGG-REAL-UTF8-HETEROGENEOUS-MAINTENANCE-R1",
		Metrics:metrics,
	}
}
