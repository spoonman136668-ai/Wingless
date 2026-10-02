$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'

& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}

$raw=& go run ./cmd/wlm-lm-raw-rep-compositional-reasoning-r1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'){throw 'WINGLESS_RESULT_SCHEMA_INVALID'}
if([string]$probe.experiment-cne'WLM-LM-RAW-REP-COMPOSITIONAL-REASONING-R1'){throw 'WINGLESS_RESULT_EXPERIMENT_INVALID'}

$required=@(
'training_record_count',
'heldout_pair_count',
'selected_representation_count',
'selected_true_motif_match_count',
'training_representation_decode_failure_count',
'heldout_representation_decode_failure_count',
'heldout_context_leak_count',
'factorized_bit_coverage_missing_count',
'factorized_heldout_accuracy',
'lookup_heldout_accuracy',
'raw_output_exact_accuracy',
'factorized_counter_count',
'lookup_counter_count',
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
if($metrics.training_record_count-ne 1536-or$metrics.heldout_pair_count-ne 16){throw 'WINGLESS_REASONING_CARDINALITY_INVALID'}
if($metrics.factorized_counter_count-ne 24-or$metrics.lookup_counter_count-ne 512){throw 'WINGLESS_REASONING_STATE_CARDINALITY_INVALID'}
if($metrics.tokenizer_use_count-ne 0-or$metrics.external_model_call_count-ne 0-or$metrics.capacity_growth_event_count-ne 0-or$metrics.counter_overflow_count-ne 0){throw 'WINGLESS_REASONING_CONTROL_INVALID'}
if($metrics.factorized_heldout_accuracy-lt 0-or$metrics.factorized_heldout_accuracy-gt 1-or$metrics.lookup_heldout_accuracy-lt 0-or$metrics.lookup_heldout_accuracy-gt 1-or$metrics.raw_output_exact_accuracy-lt 0-or$metrics.raw_output_exact_accuracy-gt 1){throw 'WINGLESS_REASONING_METRIC_RANGE_INVALID'}

$result=[ordered]@{
 schema='wingless.research-scientific-result.v1'
 experiment='WLM-LM-RAW-REP-COMPOSITIONAL-REASONING-R1'
 metrics=$metrics
}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Raw representation compositional reasoning gate completed; classify only from the frozen preregistration.'
