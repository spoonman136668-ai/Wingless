$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'
& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}
$raw=& go run ./cmd/wlm-si-raw-rep-context-sampling-falsification-r2
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'){throw 'WINGLESS_RESULT_SCHEMA_INVALID'}
if([string]$probe.experiment-cne'WLM-SI-RAW-REP-CONTEXT-SAMPLING-FALSIFICATION-R2'){throw 'WINGLESS_RESULT_EXPERIMENT_INVALID'}
$required=@('matched_bridge_record_count','independent_bridge_record_count','selected_surface_representation_count','selected_true_surface_motif_match_count','matched_max_context_count_l1_distance','independent_max_context_count_l1_distance','matched_alias_class_count','matched_three_member_class_count','matched_correct_alias_triplet_count','matched_unseen_family_accuracy','independent_alias_class_count','independent_three_member_class_count','independent_correct_alias_triplet_count','independent_missing_transition_count','independent_unseen_family_accuracy','exact_surface_unseen_family_accuracy','task_training_record_count','family2_task_training_record_count','unseen_family_evaluation_count','tokenizer_use_count','external_model_call_count','capacity_growth_event_count','invalid_row_count','counter_overflow_count')
$metrics=[ordered]@{}
foreach($name in $required){$p=$probe.metrics.PSObject.Properties[$name];if($null-eq$p){throw "WINGLESS_METRIC_MISSING:$name"};$v=[double]$p.Value;if([double]::IsNaN($v)-or[double]::IsInfinity($v)){throw "WINGLESS_METRIC_NONFINITE:$name"};$metrics[$name]=$v}
if($metrics.matched_bridge_record_count-ne 1728-or$metrics.independent_bridge_record_count-ne 1728){throw 'WINGLESS_BRIDGE_CARDINALITY_INVALID'}
if($metrics.selected_surface_representation_count-ne 27-or$metrics.selected_true_surface_motif_match_count-ne 27){throw 'WINGLESS_REPRESENTATION_CARDINALITY_INVALID'}
if($metrics.matched_max_context_count_l1_distance-ne 0-or$metrics.independent_max_context_count_l1_distance-le 0){throw 'WINGLESS_SAMPLING_CONTROL_INVALID'}
if($metrics.task_training_record_count-ne 3200-or$metrics.family2_task_training_record_count-ne 0-or$metrics.unseen_family_evaluation_count-ne 80){throw 'WINGLESS_TASK_CARDINALITY_INVALID'}
if($metrics.tokenizer_use_count-ne 0-or$metrics.external_model_call_count-ne 0-or$metrics.capacity_growth_event_count-ne 0-or$metrics.invalid_row_count-ne 0-or$metrics.counter_overflow_count-ne 0){throw 'WINGLESS_CONTROL_INVALID'}
foreach($name in @('matched_unseen_family_accuracy','independent_unseen_family_accuracy','exact_surface_unseen_family_accuracy')){if($metrics[$name]-lt 0-or$metrics[$name]-gt 1){throw "WINGLESS_ACCURACY_RANGE_INVALID:$name"}}
$result=[ordered]@{schema='wingless.research-scientific-result.v1';experiment='WLM-SI-RAW-REP-CONTEXT-SAMPLING-FALSIFICATION-R2';metrics=$metrics}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Finite context-sampling falsification completed; classify only from the frozen preregistration.'