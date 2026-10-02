$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'

& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}

$raw=& go run ./cmd/wlm-lm-raw-rep-contextual-alias-multistep-r1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'){throw 'WINGLESS_RESULT_SCHEMA_INVALID'}
if([string]$probe.experiment-cne'WLM-LM-RAW-REP-CONTEXTUAL-ALIAS-MULTISTEP-R1'){throw 'WINGLESS_RESULT_EXPERIMENT_INVALID'}

$required=@(
'training_record_count',
'bridge_record_count',
'selected_surface_representation_count',
'induced_alias_class_count',
'correct_evaluator_alias_pair_count',
'heldout_program_count',
'mixed_alias_program_count',
'program_decode_failure_count',
'conjugated_right_program_count',
'induced_conjugated_right_accuracy',
'long_program_count',
'induced_long_program_accuracy',
'induced_overall_multistep_accuracy',
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
if($metrics.training_record_count-ne 3200){throw 'WINGLESS_TRAINING_CARDINALITY_INVALID'}
if($metrics.bridge_record_count-ne 576){throw 'WINGLESS_BRIDGE_CARDINALITY_INVALID'}
if($metrics.heldout_program_count-ne 192-or$metrics.conjugated_right_program_count-ne 64-or$metrics.long_program_count-ne 128){throw 'WINGLESS_PROGRAM_CARDINALITY_INVALID'}
if($metrics.mixed_alias_program_count-ne 192){throw 'WINGLESS_MIXED_ALIAS_CARDINALITY_INVALID'}
if($metrics.tokenizer_use_count-ne 0-or$metrics.external_model_call_count-ne 0-or$metrics.capacity_growth_event_count-ne 0-or$metrics.counter_overflow_count-ne 0){throw 'WINGLESS_CONTROL_INVALID'}
foreach($name in @('induced_conjugated_right_accuracy','induced_long_program_accuracy','induced_overall_multistep_accuracy','exact_surface_overall_multistep_accuracy','shuffled_bridge_overall_multistep_accuracy')){
  if($metrics[$name]-lt 0-or$metrics[$name]-gt 1){throw "WINGLESS_ACCURACY_RANGE_INVALID:$name"}
}
$result=[ordered]@{
 schema='wingless.research-scientific-result.v1'
 experiment='WLM-LM-RAW-REP-CONTEXTUAL-ALIAS-MULTISTEP-R1'
 metrics=$metrics
}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Context-induced alias multi-step reasoning completed; classify only from the frozen preregistration.'
