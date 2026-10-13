$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'
$prereg=Get-Content -Raw 'docs/experiments/wlm-si-raw-rep-context-phase-anchor-identifiability-r5.json'|ConvertFrom-Json
if([string]$prereg.exact_parent_sha-cne'ef535451a788c8b9db328884b7d54fd4e77272be'){throw 'WINGLESS_PARENT_IDENTITY_INVALID'}
$north=(Get-FileHash -Algorithm SHA256 'research/bootstrap/wingless-lm-north-star.txt').Hash.ToLowerInvariant()
if($north-cne'e9a7b9b7371d6636d791e47fb1fb166c143687d8d41013a8fb5826ef3a6ac465'){throw 'WINGLESS_NORTH_STAR_IDENTITY_INVALID'}
& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}
$r4raw=& go run ./cmd/wlm-si-raw-rep-context-support-identifiability-r4
if($LASTEXITCODE-ne 0){throw 'WINGLESS_R4_REPLICATION_FAILED'}
$r4=($r4raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$r4.schema-cne'wingless.research-scientific-result.v1'-or[string]$r4.experiment-cne'WLM-SI-RAW-REP-CONTEXT-SUPPORT-IDENTIFIABILITY-R4'){throw 'WINGLESS_R4_IDENTITY_INVALID'}
$r4m=$r4.metrics
if([double]$r4m.exhaustive_bridge_record_count-ne 432-or[double]$r4m.exhaustive_records_per_surface_min-ne 16-or[double]$r4m.exhaustive_records_per_surface_max-ne 16-or[double]$r4m.exhaustive_distinct_contexts_per_surface_min-ne 16-or[double]$r4m.exhaustive_distinct_contexts_per_surface_max-ne 16-or[double]$r4m.exhaustive_alias_class_count-ne 8-or[double]$r4m.exhaustive_three_member_class_count-ne 7-or[double]$r4m.exhaustive_correct_alias_triplet_count-ne 9-or[double]$r4m.exhaustive_cross_entity_collision_pair_count-ne 1){throw 'WINGLESS_R4_CONTROL_DRIFT'}
$raw=& go run ./cmd/wlm-si-raw-rep-context-phase-anchor-identifiability-r5
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'){throw 'WINGLESS_RESULT_SCHEMA_INVALID'}
if([string]$probe.experiment-cne'WLM-SI-RAW-REP-CONTEXT-PHASE-ANCHOR-IDENTIFIABILITY-R5'){throw 'WINGLESS_RESULT_EXPERIMENT_INVALID'}
$required=@('anchored_bridge_record_count','anchored_records_per_surface_min','anchored_records_per_surface_max','anchored_distinct_contexts_per_surface_min','anchored_distinct_contexts_per_surface_max','anchored_distinct_phase_anchors_per_surface_min','anchored_distinct_phase_anchors_per_surface_max','unanchored_control_bridge_record_count','selected_surface_representation_count','selected_true_surface_motif_match_count','unanchored_control_alias_class_count','unanchored_control_three_member_class_count','unanchored_control_correct_alias_triplet_count','unanchored_control_cross_entity_collision_pair_count','anchored_alias_class_count','anchored_three_member_class_count','anchored_correct_alias_triplet_count','anchored_cross_entity_collision_pair_count','evaluation_entity_pair_count','task_training_record_count','tokenizer_use_count','external_model_call_count','capacity_growth_event_count','invalid_row_count','counter_overflow_count')
$metrics=[ordered]@{}
foreach($name in $required){$p=$probe.metrics.PSObject.Properties[$name];if($null-eq$p){throw "WINGLESS_METRIC_MISSING:$name"};$v=[double]$p.Value;if([double]::IsNaN($v)-or[double]::IsInfinity($v)){throw "WINGLESS_METRIC_NONFINITE:$name"};$metrics[$name]=$v}
if($metrics.anchored_bridge_record_count-ne 432-or$metrics.anchored_records_per_surface_min-ne 16-or$metrics.anchored_records_per_surface_max-ne 16-or$metrics.anchored_distinct_contexts_per_surface_min-ne 16-or$metrics.anchored_distinct_contexts_per_surface_max-ne 16-or$metrics.anchored_distinct_phase_anchors_per_surface_min-ne 16-or$metrics.anchored_distinct_phase_anchors_per_surface_max-ne 16){throw 'WINGLESS_ANCHORED_CARDINALITY_INVALID'}
if($metrics.unanchored_control_bridge_record_count-ne 432-or$metrics.selected_surface_representation_count-ne 27-or$metrics.selected_true_surface_motif_match_count-ne 27){throw 'WINGLESS_REPRESENTATION_CONTROL_INVALID'}
if($metrics.unanchored_control_alias_class_count-ne 8-or$metrics.unanchored_control_three_member_class_count-ne 7-or$metrics.unanchored_control_correct_alias_triplet_count-ne 9-or$metrics.unanchored_control_cross_entity_collision_pair_count-ne 1){throw 'WINGLESS_UNANCHORED_CONTROL_INVALID'}
if($metrics.evaluation_entity_pair_count-ne 36-or$metrics.task_training_record_count-ne 0-or$metrics.tokenizer_use_count-ne 0-or$metrics.external_model_call_count-ne 0-or$metrics.capacity_growth_event_count-ne 0-or$metrics.invalid_row_count-ne 0-or$metrics.counter_overflow_count-ne 0){throw 'WINGLESS_CONTROL_INVALID'}
foreach($name in @('anchored_alias_class_count','anchored_three_member_class_count','anchored_correct_alias_triplet_count','anchored_cross_entity_collision_pair_count')){if($metrics[$name]-lt 0-or$metrics[$name]-gt 36){throw "WINGLESS_METRIC_RANGE_INVALID:$name"}}
$result=[ordered]@{schema='wingless.research-scientific-result.v1';experiment='WLM-SI-RAW-REP-CONTEXT-PHASE-ANCHOR-IDENTIFIABILITY-R5';metrics=$metrics}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Frozen phase-anchor context-support identifiability diagnostic completed; classify only from the preregistration.'