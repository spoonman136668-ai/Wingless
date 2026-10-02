$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'

& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}

$raw=& go run ./cmd/wlm-si-raw-utf8-surface-invariance-falsification-r1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'){throw 'WINGLESS_RESULT_SCHEMA_INVALID'}
if([string]$probe.experiment-cne'WLM-SI-RAW-UTF8-SURFACE-INVARIANCE-FALSIFICATION-R1'){throw 'WINGLESS_RESULT_EXPERIMENT_INVALID'}

$required=@(
'evaluation_variant_count',
'evaluation_examples_per_variant',
'minimum_original_reasoning_accuracy',
'minimum_original_exact_raw_output_accuracy',
'minimum_punctuation_reasoning_accuracy',
'minimum_punctuation_exact_raw_output_accuracy',
'minimum_double_space_reasoning_accuracy',
'minimum_double_space_exact_raw_output_accuracy',
'maximum_tab_reasoning_accuracy',
'maximum_tab_exact_raw_output_accuracy',
'maximum_uppercase_reasoning_accuracy',
'maximum_uppercase_exact_raw_output_accuracy',
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
if($metrics.evaluation_variant_count-ne 5-or$metrics.evaluation_examples_per_variant-ne 128){throw 'WINGLESS_EVALUATION_CARDINALITY_INVALID'}
if($metrics.tokenizer_use_count-ne 0-or$metrics.external_model_call_count-ne 0-or$metrics.capacity_growth_event_count-ne 0-or$metrics.invalid_row_count-ne 0){throw 'WINGLESS_CONTROL_INVALID'}
foreach($name in @(
'minimum_original_reasoning_accuracy',
'minimum_original_exact_raw_output_accuracy',
'minimum_punctuation_reasoning_accuracy',
'minimum_punctuation_exact_raw_output_accuracy',
'minimum_double_space_reasoning_accuracy',
'minimum_double_space_exact_raw_output_accuracy',
'maximum_tab_reasoning_accuracy',
'maximum_tab_exact_raw_output_accuracy',
'maximum_uppercase_reasoning_accuracy',
'maximum_uppercase_exact_raw_output_accuracy'
)){
  if($metrics[$name]-lt 0-or$metrics[$name]-gt 1){throw "WINGLESS_ACCURACY_RANGE_INVALID:$name"}
}

$result=[ordered]@{
 schema='wingless.research-scientific-result.v1'
 experiment='WLM-SI-RAW-UTF8-SURFACE-INVARIANCE-FALSIFICATION-R1'
 metrics=$metrics
}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Raw UTF8 surface invariance falsification completed; classify only from the frozen preregistration.'
