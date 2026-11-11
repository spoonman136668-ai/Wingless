$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up12-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up12-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-12 FULL-RANK DECODER SATURATION / INTEGRATION ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP12_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP12_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP12_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP12_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-fullrank-integration-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP12_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP12_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-fullrank-integration-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP12_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-fullrank-integration-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP12_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP12_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.unitary-fullrank-integration.v1') { throw 'UP12_SCHEMA_MISMATCH' }
    if ($Result.runtime_prototype_lookup) { throw 'UP12_PROTOTYPE_LOOKUP_ENABLED' }
    if ($Result.explicit_inverse_readout) { throw 'UP12_INVERSE_READOUT_ENABLED' }
    if ($Result.explicit_depth_provided) { throw 'UP12_DEPTH_PROVIDED' }
    if (-not $Result.global_phase_nuisance) { throw 'UP12_GLOBAL_PHASE_NUISANCE_MISSING' }
    if ([int]$Result.code_rank -ne 4) { throw 'UP12_CODE_RANK_MISMATCH' }
    if ([int]$Result.full_train_pool_tables -ne 128 -or [int]$Result.full_held_out_pool_tables -ne 128) {
        throw 'UP12_FULL_POOL_MISMATCH'
    }
    if ([int]$Result.min_train_marginal_count -ne 32 -or [int]$Result.min_held_out_marginal_count -ne 32) {
        throw 'UP12_MARGINAL_COUNT_MISMATCH'
    }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Static saturation pass:      $($Result.diagnosis.static_saturation_pass)"
    Write-Host "Mutable integration pass:    $($Result.diagnosis.mutable_integration_pass)"
    Write-Host "UP-11 noisy rank-4:          $($Result.diagnosis.up11_rank4_noisy_accuracy)"
    Write-Host "Saturated static accuracy:   $($Result.diagnosis.saturated_static_accuracy)"
    Write-Host "Static accuracy gain:        $($Result.diagnosis.static_accuracy_gain)"
    Write-Host "Matched control accuracy:    $($Result.diagnosis.matched_control_accuracy)"
    Write-Host "Nonlinear readout indicated: $($Result.diagnosis.nonlinear_readout_indicated)"
    Write-Host 'WINGLESS_UP12_HARNESS_PASS'
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
