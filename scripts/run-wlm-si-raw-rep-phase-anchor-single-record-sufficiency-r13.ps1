$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'
$prereg=Get-Content -Raw 'docs/experiments/wlm-si-raw-rep-phase-anchor-single-record-sufficiency-r13.json'|ConvertFrom-Json
if([string]$prereg.exact_parent_sha-cne'c2779baabf197661511fda82d89007e44be30af0'){throw 'WINGLESS_PARENT_IDENTITY_INVALID'}
$north=(Get-FileHash -Algorithm SHA256 'research/bootstrap/wingless-lm-north-star.txt').Hash.ToLowerInvariant()
if($north-cne'e9a7b9b7371d6636d791e47fb1fb166c143687d8d41013a8fb5826ef3a6ac465'){throw 'WINGLESS_NORTH_STAR_IDENTITY_INVALID'}
& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}
$r12raw=& go run ./cmd/wlm-si-raw-rep-phase-anchor-final-record-locality-r12
if($LASTEXITCODE-ne 0){throw 'WINGLESS_R12_REPLICATION_FAILED'}
$r12=($r12raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$r12.schema-cne'wingless.research-scientific-result.v1'-or[string]$r12.experiment-cne'WLM-SI-RAW-REP-PHASE-ANCHOR-FINAL-RECORD-LOCALITY-R12'){throw 'WINGLESS_R12_IDENTITY_INVALID'}
$r12m=$r12.metrics
if([double]$r12m.actual108_novel_support_addition_count-ne1-or[double]$r12m.novel_removal_intervention_count-ne1-or[double]$r12m.destabilizing_intervention_count-ne1-or[double]$r12m.novel_removal_without_destabilization_count-ne0-or[double]$r12m.nonnovel_removal_with_destabilization_count-ne0-or[double]$r12m.invalid_row_count-ne0){throw 'WINGLESS_R12_CONTROL_DRIFT'}
$raw=& go run ./cmd/wlm-si-raw-rep-phase-anchor-single-record-sufficiency-r13
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'-or[string]$probe.experiment-cne'WLM-SI-RAW-REP-PHASE-ANCHOR-SINGLE-RECORD-SUFFICIENCY-R13'){throw 'WINGLESS_RESULT_IDENTITY_INVALID'}
$required=@('selected_surface_representation_count','selected_true_surface_motif_match_count','budget107_actual108_prefix_mismatch_count','actual108_novel_support_addition_count','target_surface_count','targeted_actual_final_record_count','targeted_duplicate_final_record_count','construction_mismatch_count','all_sham108_novel_support_addition_count','targeted_rescue108_novel_support_addition_count','targeted_actual108_class_membership_mismatch_count','evaluation_entity_pair_count_total','task_training_record_count','tokenizer_use_count','external_model_call_count','capacity_growth_event_count','invalid_row_count','counter_overflow_count')
foreach($arm in @('budget107','actual108','all_sham108','targeted_rescue108')){foreach($suffix in @('bridge_record_count','records_per_surface_min','records_per_surface_max','alias_class_count','three_member_class_count','correct_alias_triplet_count','split_entity_triplet_count','cross_entity_collision_pair_count')){$required+="${arm}_$suffix"}}
$metrics=[ordered]@{}
foreach($name in $required){$p=$probe.metrics.PSObject.Properties[$name];if($null-eq$p){throw "WINGLESS_METRIC_MISSING:$name"};$v=[double]$p.Value;if([double]::IsNaN($v)-or[double]::IsInfinity($v)){throw "WINGLESS_METRIC_NONFINITE:$name"};$metrics[$name]=$v}
if($metrics.selected_surface_representation_count-ne27-or$metrics.selected_true_surface_motif_match_count-ne27-or$metrics.budget107_actual108_prefix_mismatch_count-ne0-or$metrics.actual108_novel_support_addition_count-ne1-or$metrics.target_surface_count-ne1-or$metrics.targeted_actual_final_record_count-ne1-or$metrics.targeted_duplicate_final_record_count-ne26-or$metrics.construction_mismatch_count-ne0){throw 'WINGLESS_CONSTRUCTION_INVALID'}
foreach($arm in @('actual108','all_sham108','targeted_rescue108')){if($metrics["${arm}_bridge_record_count"]-ne2916-or$metrics["${arm}_records_per_surface_min"]-ne108-or$metrics["${arm}_records_per_surface_max"]-ne108){throw "WINGLESS_ARM_CARDINALITY_INVALID:$arm"}}
if($metrics.budget107_alias_class_count-ne10-or$metrics.budget107_correct_alias_triplet_count-ne8-or$metrics.budget107_split_entity_triplet_count-ne1-or$metrics.budget107_cross_entity_collision_pair_count-ne0-or$metrics.actual108_alias_class_count-ne9-or$metrics.actual108_correct_alias_triplet_count-ne9-or$metrics.actual108_split_entity_triplet_count-ne0-or$metrics.actual108_cross_entity_collision_pair_count-ne0){throw 'WINGLESS_FROZEN_CONTROL_INVALID'}
if($metrics.evaluation_entity_pair_count_total-ne144-or$metrics.task_training_record_count-ne0-or$metrics.tokenizer_use_count-ne0-or$metrics.external_model_call_count-ne0-or$metrics.capacity_growth_event_count-ne0-or$metrics.invalid_row_count-ne0-or$metrics.counter_overflow_count-ne0){throw 'WINGLESS_VALIDITY_CONTROL_INVALID'}
$supported=$metrics.all_sham108_novel_support_addition_count-eq0-and$metrics.all_sham108_alias_class_count-eq10-and$metrics.all_sham108_correct_alias_triplet_count-eq8-and$metrics.all_sham108_split_entity_triplet_count-eq1-and$metrics.all_sham108_cross_entity_collision_pair_count-eq0-and$metrics.targeted_rescue108_novel_support_addition_count-eq1-and$metrics.targeted_rescue108_alias_class_count-eq9-and$metrics.targeted_rescue108_correct_alias_triplet_count-eq9-and$metrics.targeted_rescue108_split_entity_triplet_count-eq0-and$metrics.targeted_rescue108_cross_entity_collision_pair_count-eq0-and$metrics.targeted_actual108_class_membership_mismatch_count-eq0
$recordCountOnly=$metrics.all_sham108_alias_class_count-eq9-and$metrics.all_sham108_correct_alias_triplet_count-eq9-and$metrics.all_sham108_split_entity_triplet_count-eq0
$classification=if($supported){'single_record_sufficiency_supported'}elseif($recordCountOnly){'record_count_only_supported'}elseif($metrics.all_sham108_alias_class_count-eq10-and$metrics.all_sham108_split_entity_triplet_count-eq1){'single_record_not_sufficient'}else{'mixed'}
$result=[ordered]@{schema='wingless.research-scientific-result.v1';experiment='WLM-SI-RAW-REP-PHASE-ANCHOR-SINGLE-RECORD-SUFFICIENCY-R13';classification=$classification;metrics=$metrics}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host "Frozen single-record sufficiency panel completed; classification=$classification"