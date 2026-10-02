$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'

& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}

$raw=& go run ./cmd/wlm-si-context-count-jitter-falsification-r1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'){throw 'WINGLESS_RESULT_SCHEMA_INVALID'}
if([string]$probe.experiment-cne'WLM-SI-CONTEXT-COUNT-JITTER-FALSIFICATION-R1'){throw 'WINGLESS_RESULT_EXPERIMENT_INVALID'}

$required=@(
'task_training_record_count',
'selected_surface_representation_count',
'selected_true_surface_motif_match_count',
'normal_bridge_record_count',
'normal_bridge_decode_failure_count',
'normal_induced_alias_class_count',
'normal_two_member_class_count',
'normal_correct_evaluator_alias_pair_count',
'jitter_bridge_record_count',
'jitter_bridge_decode_failure_count',
'minimum_alias_pair_context_support_jaccard',
'maximum_alias_pair_total_variation_distance',
'jitter_induced_alias_class_count',
'jitter_two_member_class_count',
'jitter_correct_evaluator_alias_pair_count',
'tokenizer_use_count',
'external_model_call_count',
'capacity_growth_event_count',
'invalid_row_count'
)
$metrics=[ordered]@{}
foreach($name in $required){
  $p=$probe.metrics.PSObject.Properties[$name]
  if($null-eq$p){throw "WINGLESS_METRIC_MISSING:$name"}
  $v=[double]$p.Value
  if([double]::IsNaN($v)-or[double]::IsInfinity($v)){throw "WINGLESS_METRIC_NONFINITE:$name"}
  $metrics[$name]=$v
}
if($metrics.task_training_record_count-ne 3200){throw 'WINGLESS_TRAINING_CARDINALITY_INVALID'}
if($metrics.selected_surface_representation_count-ne 18){throw 'WINGLESS_SURFACE_CAPACITY_INVALID'}
if($metrics.normal_bridge_record_count-ne 576-or$metrics.jitter_bridge_record_count-ne 576){throw 'WINGLESS_BRIDGE_CARDINALITY_INVALID'}
if($metrics.normal_bridge_decode_failure_count-ne 0-or$metrics.jitter_bridge_decode_failure_count-ne 0){throw 'WINGLESS_BRIDGE_DECODE_INVALID'}
if($metrics.tokenizer_use_count-ne 0-or$metrics.external_model_call_count-ne 0-or$metrics.capacity_growth_event_count-ne 0-or$metrics.invalid_row_count-ne 0){throw 'WINGLESS_CONTROL_INVALID'}
if($metrics.minimum_alias_pair_context_support_jaccard-lt 0-or$metrics.minimum_alias_pair_context_support_jaccard-gt 1){throw 'WINGLESS_JACCARD_RANGE_INVALID'}
if($metrics.maximum_alias_pair_total_variation_distance-lt 0-or$metrics.maximum_alias_pair_total_variation_distance-gt 1){throw 'WINGLESS_TV_RANGE_INVALID'}

$result=[ordered]@{
 schema='wingless.research-scientific-result.v1'
 experiment='WLM-SI-CONTEXT-COUNT-JITTER-FALSIFICATION-R1'
 metrics=$metrics
}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Context-count jitter falsification completed; classify only from the frozen preregistration.'
