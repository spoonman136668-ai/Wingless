$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'

& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}

$raw=& go run ./cmd/wlm-si-raw-rep-alias-invariance-falsification-r1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'){throw 'WINGLESS_RESULT_SCHEMA_INVALID'}
if([string]$probe.experiment-cne'WLM-SI-RAW-REP-ALIAS-INVARIANCE-FALSIFICATION-R1'){throw 'WINGLESS_RESULT_EXPERIMENT_INVALID'}

$required=@(
'training_record_count',
'selected_representation_count',
'selected_true_motif_match_count',
'surface_value_representation_count',
'surface_operation_representation_count',
'training_decode_failure_count',
'same_family_evaluation_count',
'same_family_accuracy',
'cross_family_evaluation_count',
'cross_family_accuracy',
'cross_family_missing_transition_count',
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
if($metrics.same_family_evaluation_count-ne 160-or$metrics.cross_family_evaluation_count-ne 160){throw 'WINGLESS_EVALUATION_CARDINALITY_INVALID'}
if($metrics.tokenizer_use_count-ne 0-or$metrics.external_model_call_count-ne 0-or$metrics.capacity_growth_event_count-ne 0-or$metrics.counter_overflow_count-ne 0){throw 'WINGLESS_CONTROL_INVALID'}
foreach($name in @('same_family_accuracy','cross_family_accuracy')){
  if($metrics[$name]-lt 0-or$metrics[$name]-gt 1){throw "WINGLESS_ACCURACY_RANGE_INVALID:$name"}
}
$result=[ordered]@{
 schema='wingless.research-scientific-result.v1'
 experiment='WLM-SI-RAW-REP-ALIAS-INVARIANCE-FALSIFICATION-R1'
 metrics=$metrics
}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Raw representation alias-invariance falsification completed; classify only from the frozen preregistration.'
