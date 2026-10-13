$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'
$prereg=Get-Content -Raw 'docs/experiments/wlm-si-raw-rep-phase-anchor-budget-onset-exact-r10.json'|ConvertFrom-Json
if([string]$prereg.exact_parent_sha-cne'dee3c6ce659d4f51205da912447c5cb68ddb0e74'){throw 'WINGLESS_PARENT_IDENTITY_INVALID'}
$north=(Get-FileHash -Algorithm SHA256 'research/bootstrap/wingless-lm-north-star.txt').Hash.ToLowerInvariant()
if($north-cne'e9a7b9b7371d6636d791e47fb1fb166c143687d8d41013a8fb5826ef3a6ac465'){throw 'WINGLESS_NORTH_STAR_IDENTITY_INVALID'}
& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}
$r9raw=& go run ./cmd/wlm-si-raw-rep-phase-anchor-budget-onset-refinement-r9
if($LASTEXITCODE-ne 0){throw 'WINGLESS_R9_REPLICATION_FAILED'}
$r9=($r9raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$r9.schema-cne'wingless.research-scientific-result.v1'-or[string]$r9.experiment-cne'WLM-SI-RAW-REP-PHASE-ANCHOR-BUDGET-ONSET-REFINEMENT-R9'){throw 'WINGLESS_R9_IDENTITY_INVALID'}
$r9m=$r9.metrics
if([double]$r9m.budget104_alias_class_count-ne10-or[double]$r9m.budget104_correct_alias_triplet_count-ne8-or[double]$r9m.budget104_split_entity_triplet_count-ne1-or[double]$r9m.budget104_cross_entity_collision_pair_count-ne0-or[double]$r9m.budget108_alias_class_count-ne9-or[double]$r9m.budget108_correct_alias_triplet_count-ne9-or[double]$r9m.budget108_split_entity_triplet_count-ne0-or[double]$r9m.budget108_cross_entity_collision_pair_count-ne0){throw 'WINGLESS_R9_ENDPOINT_CONTROL_DRIFT'}
$raw=& go run ./cmd/wlm-si-raw-rep-phase-anchor-budget-onset-exact-r10
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'){throw 'WINGLESS_RESULT_SCHEMA_INVALID'}
if([string]$probe.experiment-cne'WLM-SI-RAW-REP-PHASE-ANCHOR-BUDGET-ONSET-EXACT-R10'){throw 'WINGLESS_RESULT_EXPERIMENT_INVALID'}
$required=@('selected_surface_representation_count','selected_true_surface_motif_match_count','nested_prefix_mismatch_count','evaluation_entity_pair_count_per_arm','first_restored_budget','restoration_reversal_count','class_count_increase_count','split_triplet_increase_count','task_training_record_count','tokenizer_use_count','external_model_call_count','capacity_growth_event_count','invalid_row_count','counter_overflow_count')
foreach($budget in @(104,105,106,107,108)){foreach($suffix in @('bridge_record_count','records_per_surface_min','records_per_surface_max','distinct_contexts_per_surface_min','distinct_contexts_per_surface_max','max_same_entity_support_symmetric_difference','alias_class_count','three_member_class_count','correct_alias_triplet_count','split_entity_triplet_count','cross_entity_collision_pair_count')){$required+="budget${budget}_$suffix"}}
$metrics=[ordered]@{}
foreach($name in $required){$p=$probe.metrics.PSObject.Properties[$name];if($null-eq$p){throw "WINGLESS_METRIC_MISSING:$name"};$v=[double]$p.Value;if([double]::IsNaN($v)-or[double]::IsInfinity($v)){throw "WINGLESS_METRIC_NONFINITE:$name"};$metrics[$name]=$v}
if($metrics.selected_surface_representation_count-ne27-or$metrics.selected_true_surface_motif_match_count-ne27-or$metrics.nested_prefix_mismatch_count-ne0){throw 'WINGLESS_REPRESENTATION_OR_PREFIX_INVALID'}
$expected=@{104=2808;105=2835;106=2862;107=2889;108=2916}
foreach($budget in @(104,105,106,107,108)){if($metrics["budget${budget}_bridge_record_count"]-ne$expected[$budget]-or$metrics["budget${budget}_records_per_surface_min"]-ne$budget-or$metrics["budget${budget}_records_per_surface_max"]-ne$budget){throw "WINGLESS_BUDGET_CONSTRUCTION_INVALID:$budget"};if($metrics["budget${budget}_alias_class_count"]-lt9-or$metrics["budget${budget}_alias_class_count"]-gt27-or$metrics["budget${budget}_correct_alias_triplet_count"]-lt0-or$metrics["budget${budget}_correct_alias_triplet_count"]-gt9-or$metrics["budget${budget}_split_entity_triplet_count"]-lt0-or$metrics["budget${budget}_split_entity_triplet_count"]-gt9-or$metrics["budget${budget}_cross_entity_collision_pair_count"]-ne0){throw "WINGLESS_ARM_RANGE_INVALID:$budget"}}
if($metrics.budget104_alias_class_count-ne10-or$metrics.budget104_correct_alias_triplet_count-ne8-or$metrics.budget104_split_entity_triplet_count-ne1-or$metrics.budget104_cross_entity_collision_pair_count-ne0){throw 'WINGLESS_104_CONTROL_INVALID'}
if($metrics.budget108_alias_class_count-ne9-or$metrics.budget108_correct_alias_triplet_count-ne9-or$metrics.budget108_split_entity_triplet_count-ne0-or$metrics.budget108_cross_entity_collision_pair_count-ne0){throw 'WINGLESS_108_CONTROL_INVALID'}
if(@(105,106,107,108)-notcontains[double]$metrics.first_restored_budget){throw 'WINGLESS_ONSET_VALUE_INVALID'}
if($metrics.evaluation_entity_pair_count_per_arm-ne36-or$metrics.task_training_record_count-ne0-or$metrics.tokenizer_use_count-ne0-or$metrics.external_model_call_count-ne0-or$metrics.capacity_growth_event_count-ne0-or$metrics.invalid_row_count-ne0-or$metrics.counter_overflow_count-ne0){throw 'WINGLESS_CONTROL_INVALID'}
$result=[ordered]@{schema='wingless.research-scientific-result.v1';experiment='WLM-SI-RAW-REP-PHASE-ANCHOR-BUDGET-ONSET-EXACT-R10';metrics=$metrics}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Frozen phase-anchor exact budget-onset panel completed; classify only from the preregistration.'