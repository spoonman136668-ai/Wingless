$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up41-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up41-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-41 DIRECT OPTIMIZER REACHABILITY CALIBRATION ==='

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP41_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP41_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP41_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP41_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-optimizer-reachability-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP41_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP41_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-optimizer-reachability-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP41_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-optimizer-reachability-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP41_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP41_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.optimizer-reachability-calibration.v1') { throw 'UP41_SCHEMA_MISMATCH' }
    if (-not $Result.calibration_only) { throw 'UP41_NOT_CALIBRATION_ONLY' }
    if ($Result.task_labels_used) { throw 'UP41_TASK_LABEL_LEAK' }
    if ($Result.heldout_data_used) { throw 'UP41_HELDOUT_LEAK' }
    if ($Result.architecture_selection_authorized) { throw 'UP41_ARCHITECTURE_SELECTION_AUTHORIZED' }
    if ($Result.target_source -cne 'UP-35-seal-09284aa70203adf8881be42f67084c49b338dabe') { throw 'UP41_TARGET_SOURCE_MISMATCH' }
    if ([int]$Result.target_capacity -ne 18) { throw 'UP41_TARGET_CAPACITY_MISMATCH' }
    if ([int]$Result.target_groups.Count -ne 3) { throw 'UP41_TARGET_GROUP_COUNT_MISMATCH' }
    if ([double]$Result.fixed_rms -ne 0.05) { throw 'UP41_RMS_MISMATCH' }
    if ([double]$Result.perturbation -ne 0.002) { throw 'UP41_PERTURBATION_MISMATCH' }
    if ([double]$Result.learning_rate -ne 0.002) { throw 'UP41_LEARNING_RATE_MISMATCH' }
    if ([double]$Result.maximum_coordinate_update -ne 0.004) { throw 'UP41_MAX_UPDATE_MISMATCH' }
    if ([int]$Result.optimization_steps -ne 12) { throw 'UP41_STEP_COUNT_MISMATCH' }
    if ([double]$Result.exact_fusion_tolerance -ne 0.0025) { throw 'UP41_FUSION_TOLERANCE_MISMATCH' }
    if ([int]$Result.arms.Count -ne 4) { throw 'UP41_ARM_COUNT_MISMATCH' }

    foreach ($Arm in $Result.arms) {
        if ([int]$Arm.steps.Count -ne 12) { throw "UP41_ARM_STEP_COUNT_MISMATCH_$($Arm.name)" }
    }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Current optimizer enters target basin:    $($Result.diagnosis.current_optimizer_enters_target_basin)"
    Write-Host "Current optimizer exact target grouping:  $($Result.diagnosis.current_optimizer_exact_target_grouping)"
    Write-Host "Analytic reference enters target basin:   $($Result.diagnosis.analytic_reference_enters_target_basin)"
    Write-Host "Sequential estimator trails reference:    $($Result.diagnosis.sequential_estimator_trails_analytic_reference)"
    Write-Host "Sticky cross-target fusion:               $($Result.diagnosis.sticky_projection_cross_target_fusion)"
    Write-Host "Twelve-step reachability limited:         $($Result.diagnosis.twelve_step_reachability_limited)"
    Write-Host "Estimator best-cosine delta:              $($Result.diagnosis.estimator_best_cosine_delta)"
    Write-Host "Optimizer mechanics implicated:           $($Result.diagnosis.optimizer_mechanics_implicated)"
    Write-Host 'WINGLESS_UP41_HARNESS_PASS'
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
