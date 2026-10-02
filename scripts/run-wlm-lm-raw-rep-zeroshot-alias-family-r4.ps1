$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'

& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}

$raw=& go run ./cmd/wlm-lm-raw-rep-zeroshot-alias-family-r4
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'){throw 'WINGLESS_RESULT_SCHEMA_INVALID'}
if([string]$probe.experiment-cne'WLM-LM-RAW-REP-ZEROSHOT-ALIAS-FAMILY-R4'){throw 'WINGLESS_RESULT_EXPERIMENT_INVALID'}

$required=@(
'bridge_record_count',
'selected_surface_representation_count',
'selected_true_surface_motif_match_count',
'bridge_decode_failure_count',
'induced_alias_class_count',
'induced_three_member_class_count',
'correct_evaluator_alias_triplet_count',
'task_training_record_count',
'family2_task_training_record_count',
'task_training_decode_failure_count',
'class_transition_key_count',
'unseen_family_evaluation_count',
'unseen_family_accuracy',
'exact_surface_unseen_family_accuracy',
'exact_surface_unseen_missing_transition_count',
'shuffled_family2_alias_class_count',
'shuffled_family2_unseen_accuracy',
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
if($metrics.bridge_record_count-ne 864){throw 'WINGLESS_BRIDGE_CARDINALITY_INVALID'}
if($metrics.selected_surface_representation_count-ne 27){throw 'WINGLESS_SURFACE_CAPACITY_INVALID'}
if($metrics.task_training_record_count-ne 3200-or$metrics.family2_task_training_record_count-ne 0){throw 'WINGLESS_TASK_TRAINING_CARDINALITY_INVALID'}
if($metrics.unseen_family_evaluation_count-ne 80){throw 'WINGLESS_EVALUATION_CARDINALITY_INVALID'}
if($metrics.bridge_decode_failure_count-ne 0-or$metrics.task_training_decode_failure_count-ne 0){throw 'WINGLESS_DECODE_VALIDITY_INVALID'}
if($metrics.tokenizer_use_count-ne 0-or$metrics.external_model_call_count-ne 0-or$metrics.capacity_growth_event_count-ne 0-or$metrics.invalid_row_count-ne 0-or$metrics.counter_overflow_count-ne 0){throw 'WINGLESS_CONTROL_INVALID'}
foreach($name in @('unseen_family_accuracy','exact_surface_unseen_family_accuracy','shuffled_family2_unseen_accuracy')){
  if($metrics[$name]-lt 0-or$metrics[$name]-gt 1){throw "WINGLESS_ACCURACY_RANGE_INVALID:$name"}
}

$result=[ordered]@{
 schema='wingless.research-scientific-result.v1'
 experiment='WLM-LM-RAW-REP-ZEROSHOT-ALIAS-FAMILY-R4'
 metrics=$metrics
}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Zero-shot raw alias-family transfer R4 completed; classify only from the frozen preregistration.'
