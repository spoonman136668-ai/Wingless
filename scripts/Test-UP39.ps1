$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up39-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up39-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-39 RICH-TASK DIRECT SPECTRAL-OFFSET OPTIMIZATION ==='

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP39_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP39_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP39_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP39_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-rich-direct-objective-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP39_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP39_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-rich-direct-objective-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP39_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-rich-direct-objective-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP39_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP39_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.rich-direct-objective.v1') { throw 'UP39_SCHEMA_MISMATCH' }
    if ([int]$Result.latent_dimension -ne 96) { throw 'UP39_LATENT_DIMENSION_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP39_NOT_SINGLE_RUNTIME_STATE' }
    if (-not $Result.full_coordinate_mixing) { throw 'UP39_FULL_MIXING_MISSING' }
    if (-not $Result.starts_from_independent_offsets) { throw 'UP39_INDEPENDENT_START_MISSING' }
    if ($Result.finished_partition_menu_provided) { throw 'UP39_FINISHED_MENU_LEAK' }
    if ($Result.hard_merge_candidates_evaluated) { throw 'UP39_HARD_MERGE_SEARCH_PRESENT' }
    if ($Result.pair_affinity_field_used) { throw 'UP39_PAIR_AFFINITY_FIELD_PRESENT' }
    if (-not $Result.direct_offsets_optimized) { throw 'UP39_DIRECT_OPTIMIZATION_MISSING' }
    if ($Result.optimizer -cne 'deterministic-projected-central-difference') { throw 'UP39_OPTIMIZER_MISMATCH' }
    if (-not $Result.teacher_forced_smooth_mutable) { throw 'UP39_SMOOTH_MUTABLE_MISSING' }
    if (-not $Result.smooth_relation_probability) { throw 'UP39_SMOOTH_RELATION_MISSING' }
    if (-not $Result.harmonic_bottleneck_objective) { throw 'UP39_HARMONIC_OBJECTIVE_MISSING' }
    if (-not $Result.sticky_exact_fusion_projection) { throw 'UP39_STICKY_FUSION_BOUNDARY_MISSING' }
    if ($Result.selection_uses_heldout_data) { throw 'UP39_HELDOUT_SELECTION_LEAK' }
    if ([double]$Result.resource_price -ne 0.02) { throw 'UP39_RESOURCE_PRICE_MISMATCH' }
    if ([double]$Result.soft_capacity_tau -ne 0.0025) { throw 'UP39_SOFT_CAPACITY_TAU_MISMATCH' }
    if ([double]$Result.perturbation -ne 0.002) { throw 'UP39_PERTURBATION_MISMATCH' }
    if ([double]$Result.learning_rate -ne 0.002) { throw 'UP39_LEARNING_RATE_MISMATCH' }
    if ([double]$Result.maximum_coordinate_update -ne 0.004) { throw 'UP39_MAX_UPDATE_MISMATCH' }
    if ([int]$Result.maximum_steps -ne 12) { throw 'UP39_STEP_COUNT_MISMATCH' }
    if ([int]$Result.fit_tables -ne 96) { throw 'UP39_FIT_COUNT_MISMATCH' }
    if ([int]$Result.validation_tables -ne 32) { throw 'UP39_VALIDATION_COUNT_MISMATCH' }
    if ([int]$Result.true_heldout_tables -ne 128) { throw 'UP39_TRUE_HELDOUT_COUNT_MISMATCH' }
    if ([int]$Result.steps.Count -ne 12) { throw 'UP39_RESULT_STEP_COUNT_MISMATCH' }

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
    Write-Host "Rich direct pass:           $($Result.diagnosis.rich_direct_pass)"
    Write-Host 'WINGLESS_UP39_HARNESS_PASS'
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
