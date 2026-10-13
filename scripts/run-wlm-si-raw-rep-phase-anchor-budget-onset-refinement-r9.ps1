$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'
$prereg=Get-Content -Raw 'docs/experiments/wlm-si-raw-rep-phase-anchor-budget-onset-refinement-r9.json'|ConvertFrom-Json
if([string]$prereg.exact_parent_sha-cne'f13d53d2e4f83bf76dec1147e75deee8116c2b36'){throw 'WINGLESS_PARENT_IDENTITY_INVALID'}
$north=(Get-FileHash -Algorithm SHA256 'research/bootstrap/wingless-lm-north-star.txt').Hash.ToLowerInvariant()
if($north-cne'e9a7b9b7371d6636d791e47fb1fb166c143687d8d41013a8fb5826ef3a6ac465'){throw 'WINGLESS_NORTH_STAR_IDENTITY_INVALID'}
& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}
$r8raw=& go run ./cmd/wlm-si-raw-rep-phase-anchor-budget-onset-r8
if($LASTEXITCODE-ne 0){throw 'WINGLESS_R8_REPLICATION_FAILED'}
$r8=($r8raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$r8.schema-cne'wingless.research-scientific-result.v1'-or[string]$r8.experiment-cne'WLM-SI-RAW-REP-PHASE-ANCHOR-BUDGET-ONSET-R8'){throw 'WINGLESS_R8_IDENTITY_INVALID'}
$r8m=$r8.metrics
if([double]$r8m.budget96_alias_class_count-ne12-or[double]$r8m.budget96_correct_alias_triplet_count-ne6-or[double]$r8m.budget96_split_entity_triplet_count-ne3-or[double]$r8m.budget96_cross_entity_collision_pair_count-ne0-or[double]$r8m.budget112_alias_class_count-ne9-or[double]$r8m.budget112_correct_alias_triplet_count-ne9-or[double]$r8m.budget112_split_entity_triplet_count-ne0-or[double]$r8m.budget112_cross_entity_collision_pair_count-ne0){throw 'WINGLESS_R8_ENDPOINT_CONTROL_DRIFT'}
$raw=& go run ./cmd/wlm-si-raw-rep-phase-anchor-budget-onset-refinement-r9
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'){throw 'WINGLESS_RESULT_SCHEMA_INVALID'}
if([string]$probe.experiment-cne'WLM-SI-RAW-REP-PHASE-ANCHOR-BUDGET-ONSET-REFINEMENT-R9'){throw 'WINGLESS_RESULT_EXPERIMENT_INVALID'}
$required=@('selected_surface_representation_count','selected_true_surface_motif_match_count','nested_prefix_mismatch_count','evaluation_entity_pair_count_per_arm','interior_first_restored_budget','restoration_reversal_count','class_count_increase_count','split_triplet_increase_count','task_training_record_count','tokenizer_use_count','external_model_call_count','capacity_growth_event_count','invalid_row_count','counter_overflow_count')
foreach($budget in @(96,100,104,108,112)){foreach($suffix in @('bridge_record_count','records_per_surface_min','records_per_surface_max','distinct_contexts_per_surface_min','distinct_contexts_per_surface_max','max_same_entity_support_symmetric_difference','alias_class_count','three_member_class_count','correct_alias_triplet_count','split_entity_triplet_count','cross_entity_collision_pair_count')){$required+="budget${budget}_$suffix"}}
$metrics=[ordered]@{}
foreach($name in $required){$p=$probe.metrics.PSObject.Properties[$name];if($null-eq$p){throw "WINGLESS_METRIC_MISSING:$name"};$v=[double]$p.Value;if([double]::IsNaN($v)-or[double]::IsInfinity($v)){throw "WINGLESS_METRIC_NONFINITE:$name"};$metrics[$name]=$v}
if($metrics.selected_surface_representation_count-ne27-or$metrics.selected_true_surface_motif_match_count-ne27-or$metrics.nested_prefix_mismatch_count-ne0){throw 'WINGLESS_REPRESENTATION_OR_PREFIX_INVALID'}
$expected=@{96=2592;100=2700;104=2808;108=2916;112=3024}
foreach($budget in @(96,100,104,108,112)){if($metrics["budget${budget}_bridge_record_count"]-ne$expected[$budget]-or$metrics["budget${budget}_records_per_surface_min"]-ne$budget-or$metrics["budget${budget}_records_per_surface_max"]-ne$budget){throw "WINGLESS_BUDGET_CONSTRUCTION_INVALID:$budget"};if($metrics["budget${budget}_alias_class_count"]-lt9-or$metrics["budget${budget}_alias_class_count"]-gt27-or$metrics["budget${budget}_correct_alias_triplet_count"]-lt0-or$metrics["budget${budget}_correct_alias_triplet_count"]-gt9-or$metrics["budget${budget}_split_entity_triplet_count"]-lt0-or$metrics["budget${budget}_split_entity_triplet_count"]-gt9-or$metrics["budget${budget}_cross_entity_collision_pair_count"]-ne0){throw "WINGLESS_ARM_RANGE_INVALID:$budget"}}
if($metrics.budget96_alias_class_count-ne12-or$metrics.budget96_correct_alias_triplet_count-ne6-or$metrics.budget96_split_entity_triplet_count-ne3-or$metrics.budget96_cross_entity_collision_pair_count-ne0){throw 'WINGLESS_96_CONTROL_INVALID'}
if($metrics.budget112_alias_class_count-ne9-or$metrics.budget112_correct_alias_triplet_count-ne9-or$metrics.budget112_split_entity_triplet_count-ne0-or$metrics.budget112_cross_entity_collision_pair_count-ne0){throw 'WINGLESS_112_CONTROL_INVALID'}
if(@(0,100,104,108)-notcontains[double]$metrics.interior_first_restored_budget){throw 'WINGLESS_ONSET_VALUE_INVALID'}
if($metrics.evaluation_entity_pair_count_per_arm-ne36-or$metrics.task_training_record_count-ne0-or$metrics.tokenizer_use_count-ne0-or$metrics.external_model_call_count-ne0-or$metrics.capacity_growth_event_count-ne0-or$metrics.invalid_row_count-ne0-or$metrics.counter_overflow_count-ne0){throw 'WINGLESS_CONTROL_INVALID'}
$result=[ordered]@{schema='wingless.research-scientific-result.v1';experiment='WLM-SI-RAW-REP-PHASE-ANCHOR-BUDGET-ONSET-REFINEMENT-R9';metrics=$metrics}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Frozen phase-anchor budget-onset refinement panel completed; classify only from the preregistration.'