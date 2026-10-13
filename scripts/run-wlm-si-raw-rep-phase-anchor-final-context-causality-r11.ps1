$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'
$prereg=Get-Content -Raw 'docs/experiments/wlm-si-raw-rep-phase-anchor-final-context-causality-r11.json'|ConvertFrom-Json
if([string]$prereg.exact_parent_sha-cne'6abdc1d7e457401a7e0f9319bd6b21eccfed88c6'){throw 'WINGLESS_PARENT_IDENTITY_INVALID'}
$north=(Get-FileHash -Algorithm SHA256 'research/bootstrap/wingless-lm-north-star.txt').Hash.ToLowerInvariant()
if($north-cne'e9a7b9b7371d6636d791e47fb1fb166c143687d8d41013a8fb5826ef3a6ac465'){throw 'WINGLESS_NORTH_STAR_IDENTITY_INVALID'}
& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}
$r10raw=& go run ./cmd/wlm-si-raw-rep-phase-anchor-budget-onset-exact-r10
if($LASTEXITCODE-ne 0){throw 'WINGLESS_R10_REPLICATION_FAILED'}
$r10=($r10raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$r10.schema-cne'wingless.research-scientific-result.v1'-or[string]$r10.experiment-cne'WLM-SI-RAW-REP-PHASE-ANCHOR-BUDGET-ONSET-EXACT-R10'){throw 'WINGLESS_R10_IDENTITY_INVALID'}
$r10m=$r10.metrics
if([double]$r10m.budget107_alias_class_count-ne10-or[double]$r10m.budget107_correct_alias_triplet_count-ne8-or[double]$r10m.budget107_split_entity_triplet_count-ne1-or[double]$r10m.budget107_cross_entity_collision_pair_count-ne0-or[double]$r10m.budget108_alias_class_count-ne9-or[double]$r10m.budget108_correct_alias_triplet_count-ne9-or[double]$r10m.budget108_split_entity_triplet_count-ne0-or[double]$r10m.budget108_cross_entity_collision_pair_count-ne0){throw 'WINGLESS_R10_ENDPOINT_CONTROL_DRIFT'}
$raw=& go run ./cmd/wlm-si-raw-rep-phase-anchor-final-context-causality-r11
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'){throw 'WINGLESS_RESULT_SCHEMA_INVALID'}
if([string]$probe.experiment-cne'WLM-SI-RAW-REP-PHASE-ANCHOR-FINAL-CONTEXT-CAUSALITY-R11'){throw 'WINGLESS_RESULT_EXPERIMENT_INVALID'}
$required=@('selected_surface_representation_count','selected_true_surface_motif_match_count','budget107_actual108_prefix_mismatch_count','sham_prefix_mismatch_count','sham_duplicate_mismatch_count','actual108_novel_support_addition_count','actual108_restoring_surface_count','sham108_novel_support_addition_count','evaluation_entity_pair_count_per_arm','task_training_record_count','tokenizer_use_count','external_model_call_count','capacity_growth_event_count','invalid_row_count','counter_overflow_count')
foreach($arm in @('budget107','actual108','sham108')){foreach($suffix in @('bridge_record_count','records_per_surface_min','records_per_surface_max','distinct_contexts_per_surface_min','distinct_contexts_per_surface_max','max_same_entity_support_symmetric_difference','alias_class_count','three_member_class_count','correct_alias_triplet_count','split_entity_triplet_count','cross_entity_collision_pair_count')){$required+="${arm}_$suffix"}}
$metrics=[ordered]@{}
foreach($name in $required){$p=$probe.metrics.PSObject.Properties[$name];if($null-eq$p){throw "WINGLESS_METRIC_MISSING:$name"};$v=[double]$p.Value;if([double]::IsNaN($v)-or[double]::IsInfinity($v)){throw "WINGLESS_METRIC_NONFINITE:$name"};$metrics[$name]=$v}
if($metrics.selected_surface_representation_count-ne27-or$metrics.selected_true_surface_motif_match_count-ne27-or$metrics.budget107_actual108_prefix_mismatch_count-ne0-or$metrics.sham_prefix_mismatch_count-ne0-or$metrics.sham_duplicate_mismatch_count-ne0){throw 'WINGLESS_CONSTRUCTION_INVALID'}
if($metrics.budget107_bridge_record_count-ne2889-or$metrics.budget107_alias_class_count-ne10-or$metrics.budget107_correct_alias_triplet_count-ne8-or$metrics.budget107_split_entity_triplet_count-ne1-or$metrics.budget107_cross_entity_collision_pair_count-ne0){throw 'WINGLESS_107_CONTROL_INVALID'}
if($metrics.actual108_bridge_record_count-ne2916-or$metrics.actual108_alias_class_count-ne9-or$metrics.actual108_correct_alias_triplet_count-ne9-or$metrics.actual108_split_entity_triplet_count-ne0-or$metrics.actual108_cross_entity_collision_pair_count-ne0){throw 'WINGLESS_ACTUAL108_CONTROL_INVALID'}
if($metrics.sham108_bridge_record_count-ne2916-or$metrics.sham108_records_per_surface_min-ne108-or$metrics.sham108_records_per_surface_max-ne108-or$metrics.sham108_novel_support_addition_count-ne0-or$metrics.sham108_alias_class_count-ne10-or$metrics.sham108_correct_alias_triplet_count-ne8-or$metrics.sham108_split_entity_triplet_count-ne1-or$metrics.sham108_cross_entity_collision_pair_count-ne0){throw 'WINGLESS_SHAM_CONTROL_INVALID'}
if($metrics.actual108_novel_support_addition_count-lt1-or$metrics.actual108_novel_support_addition_count-gt27-or$metrics.actual108_restoring_surface_count-lt1-or$metrics.actual108_restoring_surface_count-gt27){throw 'WINGLESS_NOVEL_SUPPORT_RANGE_INVALID'}
if($metrics.evaluation_entity_pair_count_per_arm-ne36-or$metrics.task_training_record_count-ne0-or$metrics.tokenizer_use_count-ne0-or$metrics.external_model_call_count-ne0-or$metrics.capacity_growth_event_count-ne0-or$metrics.invalid_row_count-ne0-or$metrics.counter_overflow_count-ne0){throw 'WINGLESS_CONTROL_INVALID'}
$result=[ordered]@{schema='wingless.research-scientific-result.v1';experiment='WLM-SI-RAW-REP-PHASE-ANCHOR-FINAL-CONTEXT-CAUSALITY-R11';metrics=$metrics}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Frozen final-context causality panel completed; classify only from the preregistration.'