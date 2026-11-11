$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up47-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up47-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-47 FROZEN DENSE CHECKPOINT REPLAY ==='
    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP47_GO_ENV_FAILED' }
    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP47_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP47_ICE_BUILD_FAILED' }
    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP47_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-dense-checkpoint-replay-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP47_FOCUSED_TEST_FAILED' }
    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP47_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-dense-checkpoint-replay-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP47_PROBE1_FAILED' }
    $Probe2 = ((go run ./cmd/unitary-dense-checkpoint-replay-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP47_PROBE2_FAILED' }
    if ($Probe1 -cne $Probe2) { throw 'UP47_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.dense-checkpoint-replay.v1') { throw 'UP47_SCHEMA_MISMATCH' }
    if (-not $Result.trajectory_frozen_from_up45) { throw 'UP47_TRAJECTORY_NOT_FROZEN' }
    if ($Result.optimizer_run) { throw 'UP47_OPTIMIZER_PRESENT' }
    if ($Result.lambda_changed) { throw 'UP47_LAMBDA_CHANGED' }
    if ($Result.selection_uses_heldout_data) { throw 'UP47_HELDOUT_SELECTION_LEAK' }
    if (-not $Result.reference_evaluator_used) { throw 'UP47_REFERENCE_EVALUATOR_MISSING' }
    if ([int]$Result.candidate_count -ne 49) { throw 'UP47_CANDIDATE_COUNT_MISMATCH' }
    if ([int]$Result.points.Count -ne 49) { throw 'UP47_POINT_COUNT_MISMATCH' }
    for ($i = 0; $i -lt 49; $i++) {
        if ([int]$Result.points[$i].step -ne $i) { throw "UP47_STEP_ORDER_MISMATCH_$i" }
    }
    if ([int]$Result.points[24].exact_capacity -ne 6) { throw 'UP47_STEP24_CAPACITY_MISMATCH' }
    if ([int]$Result.points[29].exact_capacity -ne 8) { throw 'UP47_STEP29_CAPACITY_MISMATCH' }
    if (-not $Result.points[29].exact_fusion) { throw 'UP47_STEP29_FUSION_MISSING' }

    Write-Host $Probe1
    Write-Host ''
    Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
    Write-Host "Dense selected step:          $($Result.diagnosis.dense_selected_step)"
    Write-Host "Dense selected objective:     $($Result.diagnosis.dense_selected_objective)"
    Write-Host "Dense selected capacity:      $($Result.diagnosis.dense_selected_exact_capacity)"
    Write-Host "Dense selected exact fusion:  $($Result.diagnosis.dense_selected_exact_fusion)"
    Write-Host "Dense-sparse objective delta: $($Result.diagnosis.dense_minus_sparse_objective)"
    Write-Host "Held-out accuracy:            $($Result.diagnosis.selected_heldout_accuracy)"
    Write-Host "Commit accuracy:              $($Result.diagnosis.selected_commit_accuracy)"
    Write-Host "Final table accuracy:         $($Result.diagnosis.selected_final_accuracy)"
    Write-Host "Relation accuracy:            $($Result.diagnosis.selected_relation_accuracy)"
    Write-Host "Capability gates:             $($Result.diagnosis.capability_gates)"
    Write-Host 'WINGLESS_UP47_HARNESS_PASS'
}
finally {
    Pop-Location
    if ($null -eq $PriorGoCache) { Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue } else { $env:GOCACHE = $PriorGoCache }
    if ($null -eq $PriorGoTmp) { Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue } else { $env:GOTMPDIR = $PriorGoTmp }
    Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
    Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
