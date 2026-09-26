$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up42-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up42-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-42 ROLLING FULL-RANK GRADIENT HORIZON CALIBRATION ==='

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP42_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP42_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP42_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP42_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-rolling-gradient-horizon-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP42_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP42_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-rolling-gradient-horizon-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP42_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-rolling-gradient-horizon-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP42_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP42_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.rolling-gradient-horizon.v1') { throw 'UP42_SCHEMA_MISMATCH' }
    if (-not $Result.calibration_only) { throw 'UP42_NOT_CALIBRATION_ONLY' }
    if ($Result.task_labels_used) { throw 'UP42_TASK_LABEL_LEAK' }
    if ($Result.heldout_data_used) { throw 'UP42_HELDOUT_LEAK' }
    if ($Result.sticky_projection_used) { throw 'UP42_STICKY_PROJECTION_PRESENT' }
    if ($Result.architecture_selection_authorized) { throw 'UP42_ARCHITECTURE_SELECTION_AUTHORIZED' }
    if ($Result.target_source -cne 'UP-35-seal-09284aa70203adf8881be42f67084c49b338dabe') { throw 'UP42_TARGET_SOURCE_MISMATCH' }
    if ([int]$Result.target_capacity -ne 18) { throw 'UP42_TARGET_CAPACITY_MISMATCH' }
    if ([int]$Result.rolling_buffer_size -ne 7) { throw 'UP42_BUFFER_MISMATCH' }
    if ([double]$Result.perturbation -ne 0.002) { throw 'UP42_PERTURBATION_MISMATCH' }
    if ([double]$Result.learning_rate -ne 0.002) { throw 'UP42_LEARNING_RATE_MISMATCH' }
    if ([double]$Result.maximum_coordinate_update -ne 0.004) { throw 'UP42_MAX_UPDATE_MISMATCH' }
    if ([double]$Result.exact_fusion_tolerance -ne 0.0025) { throw 'UP42_FUSION_TOLERANCE_MISMATCH' }

    $ExpectedHorizons = @(12, 24, 48, 96)
    if ([int]$Result.horizons.Count -ne 4) { throw 'UP42_HORIZON_COUNT_MISMATCH' }
    for ($i = 0; $i -lt 4; $i++) {
        if ([int]$Result.horizons[$i] -ne $ExpectedHorizons[$i]) {
            throw "UP42_HORIZON_MISMATCH_$i"
        }
    }
    if ([int]$Result.runs.Count -ne 12) { throw 'UP42_RUN_COUNT_MISMATCH' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Sequential minimum tested horizon: $($Result.diagnosis.sequential_minimum_tested_horizon)"
    Write-Host "Rolling minimum tested horizon:    $($Result.diagnosis.rolling_minimum_tested_horizon)"
    Write-Host "Analytic minimum tested horizon:   $($Result.diagnosis.analytic_minimum_tested_horizon)"
    Write-Host "Rolling matches analytic horizon:  $($Result.diagnosis.rolling_matches_analytic_horizon)"
    Write-Host "Rolling beats sequential at max:   $($Result.diagnosis.rolling_beats_sequential_at_max)"
    Write-Host "Rolling efficient candidate:       $($Result.diagnosis.rolling_efficient_candidate)"
    Write-Host 'WINGLESS_UP42_HARNESS_PASS'
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
