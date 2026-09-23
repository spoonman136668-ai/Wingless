$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-4 SEQUENTIAL COMPOSITION ==='
    Write-Host "Repo: $Repo"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP4_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP4_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP4_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP4_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-sequential-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP4_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP4_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-sequential-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP4_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-sequential-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP4_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP4_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.unitary-sequential-probe.v1') { throw 'UP4_SCHEMA_MISMATCH' }
    if ([double]$Result.unitary.initial_train_accuracy -ne 0.5) { throw 'UP4_UNITARY_INITIAL_ACCURACY' }
    if ([double]$Result.unitary.train_accuracy -ne 1.0) { throw 'UP4_UNITARY_TRAIN_ACCURACY' }
    if ([double]$Result.unitary.held_out_accuracy -ne 1.0) { throw 'UP4_UNITARY_HELDOUT_ACCURACY' }
    if ([double]$Result.unitary.final_loss -gt 1e-4) { throw 'UP4_UNITARY_FINAL_LOSS' }
    if ([double]$Result.unitary.max_norm_drift -gt 1e-12) { throw 'UP4_UNITARY_NORM_DRIFT' }
    if ([double]$Result.unitary.max_round_trip_error -gt 1e-11) { throw 'UP4_UNITARY_ROUND_TRIP' }
    if ($Result.order_sensitive_ab_target -eq $Result.order_sensitive_ba_target) { throw 'UP4_NOT_ORDER_SENSITIVE' }

    Write-Host $Probe1
    Write-Host 'WINGLESS_UP4_FOCUSED_PASS'
}
finally {
    Pop-Location
}
