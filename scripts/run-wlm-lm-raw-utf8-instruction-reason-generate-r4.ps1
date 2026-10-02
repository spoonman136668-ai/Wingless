$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'

& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}

$raw=& go run ./cmd/wlm-lm-raw-utf8-instruction-reason-generate-r4
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'){throw 'WINGLESS_RESULT_SCHEMA_INVALID'}
if([string]$probe.experiment-cne'WLM-LM-RAW-UTF8-INSTRUCTION-REASON-GENERATE-R4'){throw 'WINGLESS_RESULT_EXPERIMENT_INVALID'}

$required=@(
'training_example_count',
'heldout_example_count',
'training_context_count',
'heldout_context_count',
'heldout_exact_instruction_leak_count',
'selected_operation_anchor_count',
'selected_value_anchor_count',
'maximum_anchor_width_bytes',
'maximum_total_anchor_bytes',
'value_anchor_duplicate_count',
'value_operation_anchor_collision_count',
'minimum_training_operation_accuracy',
'minimum_training_left_accuracy',
'minimum_training_right_accuracy',
'maximum_training_operand_decode_failure_count',
'minimum_heldout_operation_accuracy',
'minimum_heldout_left_accuracy',
'minimum_heldout_right_accuracy',
'maximum_heldout_operand_decode_failure_count',
'minimum_heldout_full_representation_accuracy',
'reasoning_counter_count',
'minimum_heldout_reasoning_result_accuracy',
'minimum_output_generator_class_coverage',
'maximum_output_generator_state_count',
'minimum_heldout_exact_raw_output_accuracy',
'maximum_whole_instruction_lookup_exact_output_accuracy',
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
if($metrics.training_example_count-ne 768-or$metrics.heldout_example_count-ne 128){throw 'WINGLESS_EXAMPLE_CARDINALITY_INVALID'}
if($metrics.training_context_count-ne 192-or$metrics.heldout_context_count-ne 64){throw 'WINGLESS_CONTEXT_CARDINALITY_INVALID'}
if($metrics.heldout_exact_instruction_leak_count-ne 0){throw 'WINGLESS_INSTRUCTION_LEAK_INVALID'}
if($metrics.selected_operation_anchor_count-ne 4-or$metrics.selected_value_anchor_count-ne 8){throw 'WINGLESS_ANCHOR_CARDINALITY_INVALID'}
if($metrics.maximum_anchor_width_bytes-gt 7-or$metrics.maximum_total_anchor_bytes-gt 84){throw 'WINGLESS_ANCHOR_CAPACITY_INVALID'}
if($metrics.reasoning_counter_count-ne 96){throw 'WINGLESS_REASONING_CAPACITY_INVALID'}
if($metrics.minimum_output_generator_class_coverage-ne 8){throw 'WINGLESS_GENERATOR_CLASS_COVERAGE_INVALID'}
if($metrics.maximum_output_generator_state_count-gt 160){throw 'WINGLESS_GENERATOR_STATE_CAPACITY_INVALID'}
if($metrics.tokenizer_use_count-ne 0-or$metrics.external_model_call_count-ne 0-or$metrics.capacity_growth_event_count-ne 0-or$metrics.invalid_row_count-ne 0-or$metrics.counter_overflow_count-ne 0){throw 'WINGLESS_CONTROL_INVALID'}
foreach($name in @(
'minimum_training_operation_accuracy',
'minimum_training_left_accuracy',
'minimum_training_right_accuracy',
'minimum_heldout_operation_accuracy',
'minimum_heldout_left_accuracy',
'minimum_heldout_right_accuracy',
'minimum_heldout_full_representation_accuracy',
'minimum_heldout_reasoning_result_accuracy',
'minimum_heldout_exact_raw_output_accuracy',
'maximum_whole_instruction_lookup_exact_output_accuracy'
)){
  if($metrics[$name]-lt 0-or$metrics[$name]-gt 1){throw "WINGLESS_ACCURACY_RANGE_INVALID:$name"}
}
$result=[ordered]@{
 schema='wingless.research-scientific-result.v1'
 experiment='WLM-LM-RAW-UTF8-INSTRUCTION-REASON-GENERATE-R4'
 metrics=$metrics
 operation_anchors=@($probe.operation_anchors)
 value_anchors=@($probe.value_anchors)
}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Raw UTF8 precision-first value-anchor R4 completed; classify only from the frozen preregistration.'
