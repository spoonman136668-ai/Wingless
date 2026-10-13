$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'
$prereg=Get-Content -Raw 'docs/experiments/wlm-si-raw-rep-phase-anchor-sample-budget-r7.json'|ConvertFrom-Json
if([string]$prereg.exact_parent_sha-cne'871311d16777ab340ad4ee08627ba843fc32d8bf'){throw 'WINGLESS_PARENT_IDENTITY_INVALID'}
$north=(Get-FileHash -Algorithm SHA256 'research/bootstrap/wingless-lm-north-star.txt').Hash.ToLowerInvariant()
if($north-cne'e9a7b9b7371d6636d791e47fb1fb166c143687d8d41013a8fb5826ef3a6ac465'){throw 'WINGLESS_NORTH_STAR_IDENTITY_INVALID'}
& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}
$r5raw=& go run ./cmd/wlm-si-raw-rep-context-phase-anchor-identifiability-r5
if($LASTEXITCODE-ne 0){throw 'WINGLESS_R5_REPLICATION_FAILED'}
$r5=($r5raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$r5.schema-cne'wingless.research-scientific-result.v1'-or[string]$r5.experiment-cne'WLM-SI-RAW-REP-CONTEXT-PHASE-ANCHOR-IDENTIFIABILITY-R5'){throw 'WINGLESS_R5_IDENTITY_INVALID'}
$r5m=$r5.metrics
if([double]$r5m.anchored_bridge_record_count-ne432-or[double]$r5m.anchored_alias_class_count-ne9-or[double]$r5m.anchored_three_member_class_count-ne9-or[double]$r5m.anchored_correct_alias_triplet_count-ne9-or[double]$r5m.anchored_cross_entity_collision_pair_count-ne0){throw 'WINGLESS_R5_CONTROL_DRIFT'}
$r6raw=& go run ./cmd/wlm-si-raw-rep-phase-anchor-sampling-r6
if($LASTEXITCODE-ne 0){throw 'WINGLESS_R6_REPLICATION_FAILED'}
$r6=($r6raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$r6.schema-cne'wingless.research-scientific-result.v1'-or[string]$r6.experiment-cne'WLM-SI-RAW-REP-PHASE-ANCHOR-SAMPLING-R6'){throw 'WINGLESS_R6_IDENTITY_INVALID'}
$r6m=$r6.metrics
if([double]$r6m.independent_bridge_record_count-ne1728-or[double]$r6m.independent_max_same_entity_support_symmetric_difference-ne2-or[double]$r6m.independent_alias_class_count-ne16-or[double]$r6m.independent_three_member_class_count-ne3-or[double]$r6m.independent_correct_alias_triplet_count-ne3-or[double]$r6m.independent_split_entity_triplet_count-ne6-or[double]$r6m.independent_cross_entity_collision_pair_count-ne0){throw 'WINGLESS_R6_CONTROL_DRIFT'}
$raw=& go run ./cmd/wlm-si-raw-rep-phase-anchor-sample-budget-r7
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'){throw 'WINGLESS_RESULT_SCHEMA_INVALID'}
if([string]$probe.experiment-cne'WLM-SI-RAW-REP-PHASE-ANCHOR-SAMPLE-BUDGET-R7'){throw 'WINGLESS_RESULT_EXPERIMENT_INVALID'}
$required=@('exhaustive_bridge_record_count','selected_surface_representation_count','selected_true_surface_motif_match_count','exhaustive_alias_class_count','exhaustive_three_member_class_count','exhaustive_correct_alias_triplet_count','exhaustive_split_entity_triplet_count','exhaustive_cross_entity_collision_pair_count','budget64_bridge_record_count','budget64_records_per_surface_min','budget64_records_per_surface_max','budget64_distinct_contexts_per_surface_min','budget64_distinct_contexts_per_surface_max','budget64_max_same_entity_support_symmetric_difference','budget64_alias_class_count','budget64_three_member_class_count','budget64_correct_alias_triplet_count','budget64_split_entity_triplet_count','budget64_cross_entity_collision_pair_count','budget128_bridge_record_count','budget128_records_per_surface_min','budget128_records_per_surface_max','budget128_distinct_contexts_per_surface_min','budget128_distinct_contexts_per_surface_max','budget128_max_same_entity_support_symmetric_difference','budget128_alias_class_count','budget128_three_member_class_count','budget128_correct_alias_triplet_count','budget128_split_entity_triplet_count','budget128_cross_entity_collision_pair_count','budget256_bridge_record_count','budget256_records_per_surface_min','budget256_records_per_surface_max','budget256_distinct_contexts_per_surface_min','budget256_distinct_contexts_per_surface_max','budget256_max_same_entity_support_symmetric_difference','budget256_alias_class_count','budget256_three_member_class_count','budget256_correct_alias_triplet_count','budget256_split_entity_triplet_count','budget256_cross_entity_collision_pair_count','nested_prefix_mismatch_count','evaluation_entity_pair_count_per_arm','task_training_record_count','tokenizer_use_count','external_model_call_count','capacity_growth_event_count','invalid_row_count','counter_overflow_count')
$metrics=[ordered]@{}
foreach($name in $required){$p=$probe.metrics.PSObject.Properties[$name];if($null-eq$p){throw "WINGLESS_METRIC_MISSING:$name"};$v=[double]$p.Value;if([double]::IsNaN($v)-or[double]::IsInfinity($v)){throw "WINGLESS_METRIC_NONFINITE:$name"};$metrics[$name]=$v}
if($metrics.exhaustive_bridge_record_count-ne432-or$metrics.selected_surface_representation_count-ne27-or$metrics.selected_true_surface_motif_match_count-ne27-or$metrics.exhaustive_alias_class_count-ne9-or$metrics.exhaustive_three_member_class_count-ne9-or$metrics.exhaustive_correct_alias_triplet_count-ne9-or$metrics.exhaustive_split_entity_triplet_count-ne0-or$metrics.exhaustive_cross_entity_collision_pair_count-ne0){throw 'WINGLESS_EXHAUSTIVE_CONTROL_INVALID'}
if($metrics.budget64_bridge_record_count-ne1728-or$metrics.budget64_records_per_surface_min-ne64-or$metrics.budget64_records_per_surface_max-ne64-or$metrics.budget64_max_same_entity_support_symmetric_difference-ne2-or$metrics.budget64_alias_class_count-ne16-or$metrics.budget64_three_member_class_count-ne3-or$metrics.budget64_correct_alias_triplet_count-ne3-or$metrics.budget64_split_entity_triplet_count-ne6-or$metrics.budget64_cross_entity_collision_pair_count-ne0){throw 'WINGLESS_R6_CONTROL_INVALID'}
if($metrics.budget128_bridge_record_count-ne3456-or$metrics.budget128_records_per_surface_min-ne128-or$metrics.budget128_records_per_surface_max-ne128-or$metrics.budget256_bridge_record_count-ne6912-or$metrics.budget256_records_per_surface_min-ne256-or$metrics.budget256_records_per_surface_max-ne256-or$metrics.nested_prefix_mismatch_count-ne0){throw 'WINGLESS_BUDGET_CONSTRUCTION_INVALID'}
foreach($prefix in @('budget128','budget256')){if($metrics[$prefix+'_alias_class_count']-lt9-or$metrics[$prefix+'_alias_class_count']-gt27-or$metrics[$prefix+'_correct_alias_triplet_count']-lt0-or$metrics[$prefix+'_correct_alias_triplet_count']-gt9-or$metrics[$prefix+'_split_entity_triplet_count']-lt0-or$metrics[$prefix+'_split_entity_triplet_count']-gt9-or$metrics[$prefix+'_cross_entity_collision_pair_count']-ne0){throw "WINGLESS_ARM_RANGE_INVALID:$prefix"}}
if($metrics.evaluation_entity_pair_count_per_arm-ne36-or$metrics.task_training_record_count-ne0-or$metrics.tokenizer_use_count-ne0-or$metrics.external_model_call_count-ne0-or$metrics.capacity_growth_event_count-ne0-or$metrics.invalid_row_count-ne0-or$metrics.counter_overflow_count-ne0){throw 'WINGLESS_CONTROL_INVALID'}
$result=[ordered]@{schema='wingless.research-scientific-result.v1';experiment='WLM-SI-RAW-REP-PHASE-ANCHOR-SAMPLE-BUDGET-R7';metrics=$metrics}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Frozen phase-anchor sample-budget response completed; classify only from the preregistration.'