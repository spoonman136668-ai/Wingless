$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'

& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}

$raw=& go run ./cmd/wlm-lm-raw-rep-zeroshot-multistep-r2
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'){throw 'WINGLESS_RESULT_SCHEMA_INVALID'}
if([string]$probe.experiment-cne'WLM-LM-RAW-REP-ZEROSHOT-MULTISTEP-R2'){throw 'WINGLESS_RESULT_EXPERIMENT_INVALID'}

$required=@(
'bridge_record_count',
'selected_surface_representation_count',
'selected_true_surface_motif_match_count',
'induced_alias_class_count',
'correct_evaluator_alias_triplet_count',
'task_training_record_count',
'family2_task_training_record_count',
'class_transition_key_count',
'heldout_program_count',
'family2_operation_step_count',
'program_decode_failure_count',
'conjugated_right_program_count',
'induced_conjugated_right_accuracy',
'long_program_count',
'induced_long_program_accuracy',
'induced_overall_multistep_accuracy',
'raw_output_pair_exact_accuracy',
'exact_surface_overall_multistep_accuracy',
'exact_surface_programs_with_missing_transition',
'shuffled_bridge_overall_multistep_accuracy',
'tokenizer_use_count',
'external_model_call_count',
'capacity_growth_event_count',
'invalid_row_count',
'counter_overflow_count'
)
$metrics=[ordered]@{}
foreach($name in $required){
  $p=$probe.metrics.PSObject.Properties[$name]
  if($null-eq$p){throw "WINGLESS_METRIC_MISSING:$name"}
  $v=[double]$p.Value
  if([double]::IsNaN($v)-or[double]::IsInfinity($v)){throw "WINGLESS_METRIC_NONFINITE:$name"}
  $metrics[$name]=$v
}

if($metrics.bridge_record_count-ne 864-or$metrics.selected_surface_representation_count-ne 27){throw 'WINGLESS_REPRESENTATION_CARDINALITY_INVALID'}
if($metrics.task_training_record_count-ne 3200-or$metrics.family2_task_training_record_count-ne 0){throw 'WINGLESS_TASK_TRAINING_CARDINALITY_INVALID'}
if($metrics.heldout_program_count-ne 192-or$metrics.conjugated_right_program_count-ne 64-or$metrics.long_program_count-ne 128){throw 'WINGLESS_PROGRAM_CARDINALITY_INVALID'}
if($metrics.family2_operation_step_count-ne 960){throw 'WINGLESS_FAMILY2_STEP_CARDINALITY_INVALID'}
if($metrics.program_decode_failure_count-ne 0){throw 'WINGLESS_PROGRAM_DECODE_INVALID'}
if($metrics.tokenizer_use_count-ne 0-or$metrics.external_model_call_count-ne 0-or$metrics.capacity_growth_event_count-ne 0-or$metrics.invalid_row_count-ne 0-or$metrics.counter_overflow_count-ne 0){throw 'WINGLESS_CONTROL_INVALID'}
foreach($name in @('induced_conjugated_right_accuracy','induced_long_program_accuracy','induced_overall_multistep_accuracy','raw_output_pair_exact_accuracy','exact_surface_overall_multistep_accuracy','shuffled_bridge_overall_multistep_accuracy')){
  if($metrics[$name]-lt 0-or$metrics[$name]-gt 1){throw "WINGLESS_ACCURACY_RANGE_INVALID:$name"}
}

$result=[ordered]@{
 schema='wingless.research-scientific-result.v1'
 experiment='WLM-LM-RAW-REP-ZEROSHOT-MULTISTEP-R2'
 metrics=$metrics
}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Zero-shot raw family-2 multi-step reasoning completed; classify only from the frozen preregistration.'