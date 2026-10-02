$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'

& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}

$raw=& go run ./cmd/wlm-lm-raw-rep-distributional-alias-bridge-r1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'){throw 'WINGLESS_RESULT_SCHEMA_INVALID'}
if([string]$probe.experiment-cne'WLM-LM-RAW-REP-DISTRIBUTIONAL-ALIAS-BRIDGE-R1'){throw 'WINGLESS_RESULT_EXPERIMENT_INVALID'}

$required=@(
'transition_training_record_count',
'selected_representation_count',
'selected_true_motif_match_count',
'value_alias_pair_count',
'operation_alias_pair_count',
'ambiguous_signature_group_count',
'unpaired_selected_motif_count',
'training_decode_failure_count',
'same_family_evaluation_count',
'bridged_same_family_accuracy',
'unbridged_same_family_accuracy',
'cross_family_evaluation_count',
'bridged_cross_family_accuracy',
'unbridged_cross_family_accuracy',
'bridged_minus_unbridged_cross_family_accuracy',
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
if($metrics.transition_training_record_count-ne 3200){throw 'WINGLESS_TRAINING_CARDINALITY_INVALID'}
if($metrics.same_family_evaluation_count-ne 160-or$metrics.cross_family_evaluation_count-ne 160){throw 'WINGLESS_EVALUATION_CARDINALITY_INVALID'}
if($metrics.tokenizer_use_count-ne 0-or$metrics.external_model_call_count-ne 0-or$metrics.capacity_growth_event_count-ne 0-or$metrics.counter_overflow_count-ne 0){throw 'WINGLESS_CONTROL_INVALID'}
foreach($name in @('bridged_same_family_accuracy','unbridged_same_family_accuracy','bridged_cross_family_accuracy','unbridged_cross_family_accuracy')){
  if($metrics[$name]-lt 0-or$metrics[$name]-gt 1){throw "WINGLESS_ACCURACY_RANGE_INVALID:$name"}
}
$result=[ordered]@{
 schema='wingless.research-scientific-result.v1'
 experiment='WLM-LM-RAW-REP-DISTRIBUTIONAL-ALIAS-BRIDGE-R1'
 metrics=$metrics
}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Distributional alias bridge completed; classify only from the frozen preregistration.'
