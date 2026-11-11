$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up40-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up40-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-40 FREE-RUNNING PROBABILISTIC DIRECT OBJECTIVE ==='

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP40_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP40_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP40_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP40_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-free-running-direct-objective-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP40_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP40_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-free-running-direct-objective-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP40_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-free-running-direct-objective-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP40_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP40_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.free-running-direct-objective.v1') { throw 'UP40_SCHEMA_MISMATCH' }
    if ([int]$Result.latent_dimension -ne 96) { throw 'UP40_LATENT_DIMENSION_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP40_NOT_SINGLE_RUNTIME_STATE' }
    if (-not $Result.full_coordinate_mixing) { throw 'UP40_FULL_MIXING_MISSING' }
    if (-not $Result.starts_from_independent_offsets) { throw 'UP40_INDEPENDENT_START_MISSING' }
    if ($Result.finished_partition_menu_provided) { throw 'UP40_FINISHED_MENU_LEAK' }
    if ($Result.hard_merge_candidates_evaluated) { throw 'UP40_HARD_MERGE_SEARCH_PRESENT' }
    if ($Result.pair_affinity_field_used) { throw 'UP40_PAIR_AFFINITY_FIELD_PRESENT' }
    if (-not $Result.direct_offsets_optimized) { throw 'UP40_DIRECT_OPTIMIZATION_MISSING' }
    if ($Result.optimizer -cne 'deterministic-projected-central-difference') { throw 'UP40_OPTIMIZER_MISMATCH' }
    if ($Result.teacher_forced_smooth_mutable) { throw 'UP40_TEACHER_FORCING_PRESENT' }
    if (-not $Result.free_running_probabilistic_mutable) { throw 'UP40_FREE_RUNNING_SURROGATE_MISSING' }
    if ($Result.argmax_feedback_used) { throw 'UP40_ARGMAX_FEEDBACK_PRESENT' }
    if ([int]$Result.soft_memory_dimension -ne 16) { throw 'UP40_SOFT_MEMORY_DIMENSION_MISMATCH' }
    if ($Result.soft_write_boundary -cne 'decode_probability_state_commanded_write_renormalize') { throw 'UP40_SOFT_WRITE_BOUNDARY_MISMATCH' }
    if (-not $Result.smooth_relation_probability) { throw 'UP40_SMOOTH_RELATION_MISSING' }
    if (-not $Result.harmonic_bottleneck_objective) { throw 'UP40_HARMONIC_OBJECTIVE_MISSING' }
    if (-not $Result.sticky_exact_fusion_projection) { throw 'UP40_STICKY_FUSION_BOUNDARY_MISSING' }
    if ($Result.selection_uses_heldout_data) { throw 'UP40_HELDOUT_SELECTION_LEAK' }
    if ([double]$Result.resource_price -ne 0.02) { throw 'UP40_RESOURCE_PRICE_MISMATCH' }
    if ([double]$Result.soft_capacity_tau -ne 0.0025) { throw 'UP40_SOFT_CAPACITY_TAU_MISMATCH' }
    if ([double]$Result.perturbation -ne 0.002) { throw 'UP40_PERTURBATION_MISMATCH' }
    if ([double]$Result.learning_rate -ne 0.002) { throw 'UP40_LEARNING_RATE_MISMATCH' }
    if ([double]$Result.maximum_coordinate_update -ne 0.004) { throw 'UP40_MAX_UPDATE_MISMATCH' }
    if ([int]$Result.maximum_steps -ne 12) { throw 'UP40_STEP_COUNT_MISMATCH' }
    if ([int]$Result.fit_tables -ne 96) { throw 'UP40_FIT_COUNT_MISMATCH' }
    if ([int]$Result.validation_tables -ne 32) { throw 'UP40_VALIDATION_COUNT_MISMATCH' }
    if ([int]$Result.true_heldout_tables -ne 128) { throw 'UP40_TRUE_HELDOUT_COUNT_MISMATCH' }
    if ([int]$Result.steps.Count -ne 12) { throw 'UP40_RESULT_STEP_COUNT_MISMATCH' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Selected optimization step: $($Result.diagnosis.selected_step)"
    Write-Host "Selected capacity:          $($Result.diagnosis.selected_capacity)"
    Write-Host "Exact fusion occurred:      $($Result.diagnosis.exact_fusion_occurred)"
    Write-Host "Selected held-out accuracy: $($Result.diagnosis.selected_heldout_accuracy)"
    Write-Host "Selected commit accuracy:   $($Result.diagnosis.selected_commit_accuracy)"
    Write-Host "Selected final accuracy:    $($Result.diagnosis.selected_final_accuracy)"
    Write-Host "Selected relation accuracy: $($Result.diagnosis.selected_relation_accuracy)"
    Write-Host "Held retention delta:       $($Result.diagnosis.heldout_retention_delta)"
    Write-Host "Commit retention delta:     $($Result.diagnosis.commit_retention_delta)"
    Write-Host "Free-running direct pass:   $($Result.diagnosis.free_running_direct_pass)"
    Write-Host 'WINGLESS_UP40_HARNESS_PASS'
}
finally {
    Pop-Location

    if ($null -eq $PriorGoCache) { Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue }
    else { $env:GOCACHE = $PriorGoCache }

    if ($null -eq $PriorGoTmp) { Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue }
    else { $env:GOTMPDIR = $PriorGoTmp }

    Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
    Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
