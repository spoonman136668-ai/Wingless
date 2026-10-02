$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'

& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}

$raw=& go run ./cmd/wlm-lm-raw-rep-pred-fresh-holdout-r2
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'){throw 'WINGLESS_RESULT_SCHEMA_INVALID'}
if([string]$probe.experiment-cne'WLM-LM-RAW-REP-PRED-FRESH-HOLDOUT-R2'){throw 'WINGLESS_RESULT_EXPERIMENT_INVALID'}

$expected=[ordered]@{
 training_file_identity_mismatch_count=0.0
 evaluation_file_identity_mismatch_count=0.0
 train_eval_blob_overlap_count=0.0
 valid_evaluation_file_count=4.0
 selected_motif_count=512.0
 minimum_eval_representation_coverage_fraction=0.36012885277106671
 minimum_covered_top1_accuracy_gain=0.38054777185211969
 minimum_effective_event_reduction_fraction=0.37991362272161633
 maximum_representation_minus_baseline_bits_per_byte=-0.38933721355583994
 files_with_positive_covered_top1_gain=4.0
 files_with_nonworse_bits_per_byte=4.0
 capacity_growth_event_count=0.0
 tokenizer_use_count=0.0
 external_model_call_count=0.0
 invalid_byte_rows=0.0
 counter_overflow_rows=0.0
}

$metrics=[ordered]@{}
foreach($name in $expected.Keys){
 $p=$probe.metrics.PSObject.Properties[$name]
 if($null-eq$p){throw "WINGLESS_METRIC_MISSING:$name"}
 $v=[double]$p.Value
 if([double]::IsNaN($v)-or[double]::IsInfinity($v)){throw "WINGLESS_METRIC_NONFINITE:$name"}
 if($v-ne[double]$expected[$name]){throw "WINGLESS_R1_REPLICATION_MISMATCH:$name expected=$($expected[$name]) actual=$v"}
 $metrics[$name]=$v
}

if($metrics.minimum_eval_representation_coverage_fraction-lt 0.04){throw 'WINGLESS_FROZEN_THRESHOLD_COVERAGE_FAILED'}
if($metrics.minimum_covered_top1_accuracy_gain-lt 0.05){throw 'WINGLESS_FROZEN_THRESHOLD_PREDICTION_FAILED'}
if($metrics.minimum_effective_event_reduction_fraction-lt 0.08){throw 'WINGLESS_FROZEN_THRESHOLD_EVENT_REDUCTION_FAILED'}
if($metrics.maximum_representation_minus_baseline_bits_per_byte-gt 0.0){throw 'WINGLESS_FROZEN_THRESHOLD_BPB_FAILED'}
if($metrics.files_with_positive_covered_top1_gain-ne 4){throw 'WINGLESS_FROZEN_THRESHOLD_POSITIVE_FILES_FAILED'}
if($metrics.files_with_nonworse_bits_per_byte-ne 4){throw 'WINGLESS_FROZEN_THRESHOLD_NONWORSE_FILES_FAILED'}

$result=[ordered]@{
 schema='wingless.research-scientific-result.v1'
 experiment='WLM-LM-RAW-REP-PRED-FRESH-HOLDOUT-R2'
 metrics=$metrics
}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Exact R1 scientific replication completed under the established harness marker; classify only from the unchanged frozen R1 thresholds.'
