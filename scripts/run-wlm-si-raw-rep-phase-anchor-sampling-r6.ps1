$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'
$prereg=Get-Content -Raw 'docs/experiments/wlm-si-raw-rep-phase-anchor-sampling-r6.json'|ConvertFrom-Json
if([string]$prereg.exact_parent_sha-cne'9c14096453be4ae60a5fc4e1a67b413374ca812a'){throw 'WINGLESS_PARENT_IDENTITY_INVALID'}
$north=(Get-FileHash -Algorithm SHA256 'research/bootstrap/wingless-lm-north-star.txt').Hash.ToLowerInvariant()
if($north-cne'e9a7b9b7371d6636d791e47fb1fb166c143687d8d41013a8fb5826ef3a6ac465'){throw 'WINGLESS_NORTH_STAR_IDENTITY_INVALID'}
& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}
$r5raw=& go run ./cmd/wlm-si-raw-rep-context-phase-anchor-identifiability-r5
if($LASTEXITCODE-ne 0){throw 'WINGLESS_R5_REPLICATION_FAILED'}
$r5=($r5raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$r5.schema-cne'wingless.research-scientific-result.v1'-or[string]$r5.experiment-cne'WLM-SI-RAW-REP-CONTEXT-PHASE-ANCHOR-IDENTIFIABILITY-R5'){throw 'WINGLESS_R5_IDENTITY_INVALID'}
$r5m=$r5.metrics
if([double]$r5m.anchored_bridge_record_count-ne432-or[double]$r5m.anchored_distinct_contexts_per_surface_min-ne16-or[double]$r5m.anchored_distinct_contexts_per_surface_max-ne16-or[double]$r5m.anchored_distinct_phase_anchors_per_surface_min-ne16-or[double]$r5m.anchored_distinct_phase_anchors_per_surface_max-ne16-or[double]$r5m.anchored_alias_class_count-ne9-or[double]$r5m.anchored_three_member_class_count-ne9-or[double]$r5m.anchored_correct_alias_triplet_count-ne9-or[double]$r5m.anchored_cross_entity_collision_pair_count-ne0){throw 'WINGLESS_R5_CONTROL_DRIFT'}
$raw=& go run ./cmd/wlm-si-raw-rep-phase-anchor-sampling-r6
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'){throw 'WINGLESS_RESULT_SCHEMA_INVALID'}
if([string]$probe.experiment-cne'WLM-SI-RAW-REP-PHASE-ANCHOR-SAMPLING-R6'){throw 'WINGLESS_RESULT_EXPERIMENT_INVALID'}
$required=@('exhaustive_bridge_record_count','matched_bridge_record_count','independent_bridge_record_count','selected_surface_representation_count','selected_true_surface_motif_match_count','exhaustive_distinct_contexts_per_surface_min','exhaustive_distinct_contexts_per_surface_max','exhaustive_alias_class_count','exhaustive_three_member_class_count','exhaustive_correct_alias_triplet_count','exhaustive_split_entity_triplet_count','exhaustive_cross_entity_collision_pair_count','matched_records_per_surface_min','matched_records_per_surface_max','matched_max_same_entity_support_symmetric_difference','matched_alias_class_count','matched_three_member_class_count','matched_correct_alias_triplet_count','matched_split_entity_triplet_count','matched_cross_entity_collision_pair_count','independent_records_per_surface_min','independent_records_per_surface_max','independent_max_same_entity_support_symmetric_difference','independent_alias_class_count','independent_three_member_class_count','independent_correct_alias_triplet_count','independent_split_entity_triplet_count','independent_cross_entity_collision_pair_count','evaluation_entity_pair_count_per_arm','task_training_record_count','tokenizer_use_count','external_model_call_count','capacity_growth_event_count','invalid_row_count','counter_overflow_count')
$metrics=[ordered]@{}
foreach($name in $required){$p=$probe.metrics.PSObject.Properties[$name];if($null-eq$p){throw "WINGLESS_METRIC_MISSING:$name"};$v=[double]$p.Value;if([double]::IsNaN($v)-or[double]::IsInfinity($v)){throw "WINGLESS_METRIC_NONFINITE:$name"};$metrics[$name]=$v}
if($metrics.exhaustive_bridge_record_count-ne432-or$metrics.matched_bridge_record_count-ne1728-or$metrics.independent_bridge_record_count-ne1728-or$metrics.selected_surface_representation_count-ne27-or$metrics.selected_true_surface_motif_match_count-ne27){throw 'WINGLESS_CARDINALITY_INVALID'}
if($metrics.exhaustive_distinct_contexts_per_surface_min-ne16-or$metrics.exhaustive_distinct_contexts_per_surface_max-ne16-or$metrics.exhaustive_alias_class_count-ne9-or$metrics.exhaustive_three_member_class_count-ne9-or$metrics.exhaustive_correct_alias_triplet_count-ne9-or$metrics.exhaustive_split_entity_triplet_count-ne0-or$metrics.exhaustive_cross_entity_collision_pair_count-ne0){throw 'WINGLESS_EXHAUSTIVE_CONTROL_INVALID'}
if($metrics.matched_records_per_surface_min-ne64-or$metrics.matched_records_per_surface_max-ne64-or$metrics.matched_max_same_entity_support_symmetric_difference-ne0-or$metrics.matched_alias_class_count-ne9-or$metrics.matched_three_member_class_count-ne9-or$metrics.matched_correct_alias_triplet_count-ne9-or$metrics.matched_split_entity_triplet_count-ne0-or$metrics.matched_cross_entity_collision_pair_count-ne0){throw 'WINGLESS_MATCHED_CONTROL_INVALID'}
if($metrics.independent_records_per_surface_min-ne64-or$metrics.independent_records_per_surface_max-ne64-or$metrics.independent_max_same_entity_support_symmetric_difference-le0){throw 'WINGLESS_INDEPENDENT_SAMPLING_INVALID'}
if($metrics.evaluation_entity_pair_count_per_arm-ne36-or$metrics.task_training_record_count-ne0-or$metrics.tokenizer_use_count-ne0-or$metrics.external_model_call_count-ne0-or$metrics.capacity_growth_event_count-ne0-or$metrics.invalid_row_count-ne0-or$metrics.counter_overflow_count-ne0){throw 'WINGLESS_CONTROL_INVALID'}
foreach($name in @('independent_alias_class_count','independent_three_member_class_count','independent_correct_alias_triplet_count','independent_split_entity_triplet_count','independent_cross_entity_collision_pair_count')){if($metrics[$name]-lt0-or$metrics[$name]-gt36){throw "WINGLESS_METRIC_RANGE_INVALID:$name"}}
$result=[ordered]@{schema='wingless.research-scientific-result.v1';experiment='WLM-SI-RAW-REP-PHASE-ANCHOR-SAMPLING-R6';metrics=$metrics}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Frozen phase-anchor finite-sampling diagnostic completed; classify only from the preregistration.'