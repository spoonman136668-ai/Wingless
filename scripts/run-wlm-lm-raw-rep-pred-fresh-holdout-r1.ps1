$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'
& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}
$raw=& go run ./cmd/wlm-lm-raw-rep-pred-fresh-holdout-r1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'){throw 'WINGLESS_RESULT_SCHEMA_INVALID'}
if([string]$probe.experiment-cne'WLM-LM-RAW-REP-PRED-FRESH-HOLDOUT-R1'){throw 'WINGLESS_RESULT_EXPERIMENT_INVALID'}
$required=@(
'training_file_identity_mismatch_count',
'evaluation_file_identity_mismatch_count',
'train_eval_blob_overlap_count',
'valid_evaluation_file_count',
'selected_motif_count',
'minimum_eval_representation_coverage_fraction',
'minimum_covered_top1_accuracy_gain',
'minimum_effective_event_reduction_fraction',
'maximum_representation_minus_baseline_bits_per_byte',
'files_with_positive_covered_top1_gain',
'files_with_nonworse_bits_per_byte',
'capacity_growth_event_count',
'tokenizer_use_count',
'external_model_call_count',
'invalid_byte_rows',
'counter_overflow_rows'
)
$metrics=[ordered]@{}
foreach($name in $required){
  $p=$probe.metrics.PSObject.Properties[$name]
  if($null-eq$p){throw "WINGLESS_METRIC_MISSING:$name"}
  $v=[double]$p.Value
  if([double]::IsNaN($v)-or[double]::IsInfinity($v)){throw "WINGLESS_METRIC_NONFINITE:$name"}
  $metrics[$name]=$v
}
if($metrics.training_file_identity_mismatch_count-ne 0-or$metrics.evaluation_file_identity_mismatch_count-ne 0-or$metrics.train_eval_blob_overlap_count-ne 0){throw 'WINGLESS_CORPUS_IDENTITY_INVALID'}
if($metrics.valid_evaluation_file_count-ne 4-or$metrics.selected_motif_count-ne 512){throw 'WINGLESS_CARDINALITY_INVALID'}
if($metrics.capacity_growth_event_count-ne 0-or$metrics.tokenizer_use_count-ne 0-or$metrics.external_model_call_count-ne 0-or$metrics.invalid_byte_rows-ne 0-or$metrics.counter_overflow_rows-ne 0){throw 'WINGLESS_CONTROL_INVALID'}
if($metrics.minimum_eval_representation_coverage_fraction-lt 0-or$metrics.minimum_eval_representation_coverage_fraction-gt 1){throw 'WINGLESS_COVERAGE_RANGE_INVALID'}
if($metrics.minimum_covered_top1_accuracy_gain-lt -1-or$metrics.minimum_covered_top1_accuracy_gain-gt 1){throw 'WINGLESS_ACCURACY_RANGE_INVALID'}
if($metrics.minimum_effective_event_reduction_fraction-lt 0-or$metrics.minimum_effective_event_reduction_fraction-gt 1){throw 'WINGLESS_REDUCTION_RANGE_INVALID'}
if($metrics.files_with_positive_covered_top1_gain-lt 0-or$metrics.files_with_positive_covered_top1_gain-gt 4){throw 'WINGLESS_POSITIVE_FILE_COUNT_INVALID'}
if($metrics.files_with_nonworse_bits_per_byte-lt 0-or$metrics.files_with_nonworse_bits_per_byte-gt 4){throw 'WINGLESS_BPB_FILE_COUNT_INVALID'}
$result=[ordered]@{schema='wingless.research-scientific-result.v1';experiment='WLM-LM-RAW-REP-PRED-FRESH-HOLDOUT-R1';metrics=$metrics}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC GATE'
Write-Host 'Raw-input representation prediction holdout completed; classify only from the frozen preregistration.'
