$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up43-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up43-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-43 ROLLING GRADIENT FREE-RUNNING TASK BRIDGE ==='

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP43_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP43_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP43_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP43_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-rolling-task-bridge-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP43_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP43_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-rolling-task-bridge-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP43_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-rolling-task-bridge-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP43_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP43_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.rolling-task-bridge.v1') { throw 'UP43_SCHEMA_MISMATCH' }
    if ($Result.known_target_used) { throw 'UP43_KNOWN_TARGET_LEAK' }
    if (-not $Result.rolling_full_rank_directional_memory) { throw 'UP43_ROLLING_MEMORY_MISSING' }
    if ([int]$Result.rolling_buffer_size -ne 7) { throw 'UP43_BUFFER_MISMATCH' }
    if ($Result.sticky_projection_used) { throw 'UP43_STICKY_PROJECTION_PRESENT' }
    if ($Result.teacher_forced_smooth_mutable) { throw 'UP43_TEACHER_FORCING_PRESENT' }
    if (-not $Result.free_running_probabilistic_mutable) { throw 'UP43_FREE_RUNNING_SURROGATE_MISSING' }
    if ($Result.selection_uses_heldout_data) { throw 'UP43_HELDOUT_SELECTION_LEAK' }
    if ([int]$Result.maximum_steps -ne 48) { throw 'UP43_STEP_COUNT_MISMATCH' }
    if ([double]$Result.perturbation -ne 0.002) { throw 'UP43_PERTURBATION_MISMATCH' }
    if ([double]$Result.learning_rate -ne 0.002) { throw 'UP43_LEARNING_RATE_MISMATCH' }
    if ([double]$Result.maximum_coordinate_update -ne 0.004) { throw 'UP43_MAX_UPDATE_MISMATCH' }
    if ([double]$Result.exact_fusion_tolerance -ne 0.0025) { throw 'UP43_FUSION_TOLERANCE_MISMATCH' }
    if ([double]$Result.resource_price -ne 0.02) { throw 'UP43_RESOURCE_PRICE_MISMATCH' }
    if ([double]$Result.soft_capacity_tau -ne 0.0025) { throw 'UP43_SOFT_CAPACITY_TAU_MISMATCH' }
    if ([int]$Result.fit_tables -ne 96) { throw 'UP43_FIT_COUNT_MISMATCH' }
    if ([int]$Result.validation_tables -ne 32) { throw 'UP43_VALIDATION_COUNT_MISMATCH' }
    if ([int]$Result.true_heldout_tables -ne 128) { throw 'UP43_HELDOUT_COUNT_MISMATCH' }
    if ([int]$Result.steps.Count -ne 48) { throw 'UP43_RESULT_STEP_COUNT_MISMATCH' }
    if ([int]$Result.checkpoints.Count -ne 3) { throw 'UP43_CHECKPOINT_COUNT_MISMATCH' }

    $ExpectedCheckpoints = @(12, 24, 48)
    for ($i = 0; $i -lt 3; $i++) {
        if ([int]$Result.checkpoint_steps[$i] -ne $ExpectedCheckpoints[$i]) {
            throw "UP43_CHECKPOINT_SCHEDULE_MISMATCH_$i"
        }
        if ([int]$Result.checkpoints[$i].step -ne $ExpectedCheckpoints[$i]) {
            throw "UP43_CHECKPOINT_RESULT_MISMATCH_$i"
        }
    }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Selected checkpoint:               $($Result.diagnosis.selected_checkpoint_step)"
    Write-Host "Objective gain:                    $($Result.diagnosis.selected_objective_gain)"
    Write-Host "First provisional fusion step:     $($Result.diagnosis.first_provisional_fusion_step)"
    Write-Host "Maximum provisional capacity:      $($Result.diagnosis.maximum_provisional_capacity)"
    Write-Host "Selected provisional capacity:     $($Result.diagnosis.selected_provisional_capacity)"
    Write-Host "Selected nearest offset gap:       $($Result.diagnosis.selected_nearest_offset_gap)"
    Write-Host "Task gradient entered fusion basin:$($Result.diagnosis.task_gradient_entered_fusion_basin)"
    Write-Host "Held-out accuracy:                 $($Result.diagnosis.selected_heldout_accuracy)"
    Write-Host "Commit accuracy:                   $($Result.diagnosis.selected_commit_accuracy)"
    Write-Host "Final table accuracy:              $($Result.diagnosis.selected_final_accuracy)"
    Write-Host "Relation accuracy:                 $($Result.diagnosis.selected_relation_accuracy)"
    Write-Host "Capability gates sans exact fusion:$($Result.diagnosis.capability_gates_without_exact_fusion)"
    Write-Host 'WINGLESS_UP43_HARNESS_PASS'
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
