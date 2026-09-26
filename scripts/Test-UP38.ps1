$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up38-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up38-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-38 DIRECT CONTINUOUS SPECTRAL-OFFSET OPTIMIZATION ==='
    Write-Host "Repo: $Repo"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP38_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP38_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP38_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP38_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-direct-offset-optimization-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP38_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP38_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-direct-offset-optimization-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP38_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-direct-offset-optimization-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP38_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP38_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.direct-offset-optimization.v1') { throw 'UP38_SCHEMA_MISMATCH' }
    if ([int]$Result.latent_dimension -ne 96) { throw 'UP38_LATENT_DIMENSION_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP38_NOT_SINGLE_RUNTIME_STATE' }
    if (-not $Result.full_coordinate_mixing) { throw 'UP38_FULL_MIXING_MISSING' }
    if (-not $Result.starts_from_independent_offsets) { throw 'UP38_INDEPENDENT_START_MISSING' }
    if ($Result.finished_partition_menu_provided) { throw 'UP38_FINISHED_MENU_LEAK' }
    if ($Result.hard_merge_candidates_evaluated) { throw 'UP38_HARD_MERGE_SEARCH_PRESENT' }
    if ($Result.pair_affinity_field_used) { throw 'UP38_PAIR_AFFINITY_FIELD_PRESENT' }
    if (-not $Result.direct_offsets_optimized) { throw 'UP38_DIRECT_OPTIMIZATION_MISSING' }
    if ($Result.optimizer -cne 'deterministic-projected-central-difference') { throw 'UP38_OPTIMIZER_MISMATCH' }
    if (-not $Result.smooth_task_objective) { throw 'UP38_SMOOTH_TASK_OBJECTIVE_MISSING' }
    if (-not $Result.smooth_resource_objective) { throw 'UP38_SMOOTH_RESOURCE_OBJECTIVE_MISSING' }
    if (-not $Result.sticky_exact_fusion_projection) { throw 'UP38_STICKY_FUSION_BOUNDARY_MISSING' }
    if ($Result.selection_uses_heldout_data) { throw 'UP38_HELDOUT_SELECTION_LEAK' }
    if ([double]$Result.resource_price -ne 0.02) { throw 'UP38_RESOURCE_PRICE_MISMATCH' }
    if ([double]$Result.soft_capacity_tau -ne 0.0025) { throw 'UP38_SOFT_CAPACITY_TAU_MISMATCH' }
    if ([double]$Result.perturbation -ne 0.002) { throw 'UP38_PERTURBATION_MISMATCH' }
    if ([double]$Result.learning_rate -ne 0.002) { throw 'UP38_LEARNING_RATE_MISMATCH' }
    if ([double]$Result.maximum_coordinate_update -ne 0.004) { throw 'UP38_MAX_UPDATE_MISMATCH' }
    if ([int]$Result.maximum_steps -ne 12) { throw 'UP38_STEP_COUNT_MISMATCH' }
    if ([int]$Result.fit_tables -ne 96) { throw 'UP38_FIT_COUNT_MISMATCH' }
    if ([int]$Result.validation_tables -ne 32) { throw 'UP38_VALIDATION_COUNT_MISMATCH' }
    if ([int]$Result.true_heldout_tables -ne 128) { throw 'UP38_TRUE_HELDOUT_COUNT_MISMATCH' }
    if ([int]$Result.steps.Count -ne 12) { throw 'UP38_RESULT_STEP_COUNT_MISMATCH' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Selected optimization step: $($Result.diagnosis.selected_step)"
    Write-Host "Selected capacity:          $($Result.diagnosis.selected_capacity)"
    Write-Host "Exact fusion occurred:      $($Result.diagnosis.exact_fusion_occurred)"
    Write-Host "Selected held-out accuracy: $($Result.diagnosis.selected_heldout_accuracy)"
    Write-Host "Selected commit accuracy:   $($Result.diagnosis.selected_commit_accuracy)"
    Write-Host "Full held-out accuracy:     $($Result.diagnosis.full_capacity_heldout_accuracy)"
    Write-Host "Full commit accuracy:       $($Result.diagnosis.full_capacity_commit_accuracy)"
    Write-Host "Held retention delta:       $($Result.diagnosis.heldout_retention_delta)"
    Write-Host "Commit retention delta:     $($Result.diagnosis.commit_retention_delta)"
    Write-Host "Direct optimization pass:   $($Result.diagnosis.direct_optimization_pass)"
    Write-Host 'WINGLESS_UP38_HARNESS_PASS'
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
