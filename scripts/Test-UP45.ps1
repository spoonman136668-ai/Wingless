$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up45-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up45-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-45 REVERSIBLE CONVEX EXACT FUSION ==='
    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP45_GO_ENV_FAILED' }
    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP45_GO_CLEAN_FAILED' }
    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP45_ICE_BUILD_FAILED' }
    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP45_ICE_VALIDATE_FAILED' }
    go test ./unitary ./cmd/unitary-proximal-task-bridge-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP45_FOCUSED_TEST_FAILED' }
    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP45_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-proximal-task-bridge-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP45_PROBE1_FAILED' }
    $Probe2 = ((go run ./cmd/unitary-proximal-task-bridge-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP45_PROBE2_FAILED' }
    if ($Probe1 -cne $Probe2) { throw 'UP45_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.proximal-task-bridge.v1') { throw 'UP45_SCHEMA_MISMATCH' }
    if ($Result.known_target_used) { throw 'UP45_KNOWN_TARGET_LEAK' }
    if ($Result.finished_partition_menu_provided) { throw 'UP45_PARTITION_MENU_PRESENT' }
    if ($Result.hard_merge_candidates_evaluated) { throw 'UP45_HARD_MERGE_SEARCH_PRESENT' }
    if ($Result.pair_affinity_field_used) { throw 'UP45_PAIR_AFFINITY_PRESENT' }
    if (-not $Result.rolling_full_rank_directional_memory) { throw 'UP45_ROLLING_MEMORY_MISSING' }
    if ([int]$Result.rolling_buffer_size -ne 7) { throw 'UP45_BUFFER_MISMATCH' }
    if ($Result.sticky_projection_used) { throw 'UP45_STICKY_PROJECTION_PRESENT' }
    if (-not $Result.convex_proximal_fusion_used) { throw 'UP45_PROXIMAL_FUSION_MISSING' }
    if (-not $Result.proximal_lambda_derived) { throw 'UP45_DERIVED_LAMBDA_MISSING' }
    if ([double]$Result.proximal_lambda -ne 0.00004) { throw 'UP45_LAMBDA_MISMATCH' }
    if (-not $Result.sqrt_probability_soft_memory) { throw 'UP45_SQRT_MEMORY_MISSING' }
    if ($Result.global_soft_memory_renormalization) { throw 'UP45_GLOBAL_RENORMALIZATION_PRESENT' }
    if ($Result.selection_uses_heldout_data) { throw 'UP45_HELDOUT_SELECTION_LEAK' }
    if ($Result.exact_fusion_used_for_selection) { throw 'UP45_FUSION_SELECTION_LEAK' }
    if ([int]$Result.maximum_steps -ne 48) { throw 'UP45_STEP_COUNT_MISMATCH' }
    if ([double]$Result.perturbation -ne 0.002) { throw 'UP45_PERTURBATION_MISMATCH' }
    if ([double]$Result.learning_rate -ne 0.002) { throw 'UP45_LEARNING_RATE_MISMATCH' }
    if ([double]$Result.maximum_coordinate_update -ne 0.004) { throw 'UP45_MAX_UPDATE_MISMATCH' }
    if ([double]$Result.resource_price -ne 0.02) { throw 'UP45_RESOURCE_PRICE_MISMATCH' }
    if ([double]$Result.soft_capacity_tau -ne 0.0025) { throw 'UP45_SOFT_CAPACITY_TAU_MISMATCH' }
    if ([int]$Result.fit_tables -ne 96) { throw 'UP45_FIT_COUNT_MISMATCH' }
    if ([int]$Result.validation_tables -ne 32) { throw 'UP45_VALIDATION_COUNT_MISMATCH' }
    if ([int]$Result.true_heldout_tables -ne 128) { throw 'UP45_HELDOUT_COUNT_MISMATCH' }
    if ([int]$Result.steps.Count -ne 48) { throw 'UP45_RESULT_STEP_COUNT_MISMATCH' }
    if ([int]$Result.checkpoints.Count -ne 3) { throw 'UP45_CHECKPOINT_COUNT_MISMATCH' }

    $ExpectedCheckpoints = @(12, 24, 48)
    for ($i = 0; $i -lt 3; $i++) {
        if ([int]$Result.checkpoint_steps[$i] -ne $ExpectedCheckpoints[$i]) { throw "UP45_CHECKPOINT_SCHEDULE_MISMATCH_$i" }
        if ([int]$Result.checkpoints[$i].step -ne $ExpectedCheckpoints[$i]) { throw "UP45_CHECKPOINT_RESULT_MISMATCH_$i" }
    }

    Write-Host $Probe1
    Write-Host ''
    Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
    Write-Host "Selected checkpoint:          $($Result.diagnosis.selected_checkpoint_step)"
    Write-Host "Objective gain:               $($Result.diagnosis.selected_objective_gain)"
    Write-Host "First exact fusion step:      $($Result.diagnosis.first_exact_fusion_step)"
    Write-Host "Maximum exact capacity:       $($Result.diagnosis.maximum_exact_capacity)"
    Write-Host "Selected exact capacity:      $($Result.diagnosis.selected_exact_capacity)"
    Write-Host "Selected exact fusion:        $($Result.diagnosis.selected_exact_fusion)"
    Write-Host "Fusion separated later:       $($Result.diagnosis.any_fusion_separated_later)"
    Write-Host "Held-out accuracy:            $($Result.diagnosis.selected_heldout_accuracy)"
    Write-Host "Commit accuracy:              $($Result.diagnosis.selected_commit_accuracy)"
    Write-Host "Final table accuracy:         $($Result.diagnosis.selected_final_accuracy)"
    Write-Host "Relation accuracy:            $($Result.diagnosis.selected_relation_accuracy)"
    Write-Host "Capability gates:             $($Result.diagnosis.capability_gates)"
    Write-Host 'WINGLESS_UP45_HARNESS_PASS'
}
finally {
    Pop-Location
    if ($null -eq $PriorGoCache) { Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue } else { $env:GOCACHE = $PriorGoCache }
    if ($null -eq $PriorGoTmp) { Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue } else { $env:GOTMPDIR = $PriorGoTmp }
    Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
    Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
