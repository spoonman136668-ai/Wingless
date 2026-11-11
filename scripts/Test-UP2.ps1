$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-2 MULTILAYER PHASE ROUTING ==='
    Write-Host "Repo: $Repo"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP2_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP2_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP2_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP2_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-multilayer-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP2_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP2_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-multilayer-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP2_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-multilayer-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP2_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP2_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.unitary-multilayer-probe.v1') { throw 'UP2_SCHEMA_MISMATCH' }
    if ([double]$Result.unitary.initial_accuracy -ne 0.25) { throw 'UP2_UNITARY_INITIAL_ACCURACY_MISMATCH' }
    if ([double]$Result.unitary.train_accuracy -ne 1.0) { throw 'UP2_UNITARY_TRAIN_ACCURACY_MISMATCH' }
    if ([double]$Result.unitary.held_out_accuracy -ne 1.0) { throw 'UP2_UNITARY_HELDOUT_ACCURACY_MISMATCH' }
    if ([double]$Result.non_unitary.initial_accuracy -ne 0.25) { throw 'UP2_CONTROL_INITIAL_ACCURACY_MISMATCH' }
    if ([double]$Result.non_unitary.train_accuracy -ne 1.0) { throw 'UP2_CONTROL_TRAIN_ACCURACY_MISMATCH' }
    if ([double]$Result.non_unitary.held_out_accuracy -ne 1.0) { throw 'UP2_CONTROL_HELDOUT_ACCURACY_MISMATCH' }
    if ([double]$Result.unitary.final_loss -gt 1e-6) { throw 'UP2_UNITARY_FINAL_LOSS_MISMATCH' }
    if ([double]$Result.unitary.max_norm_drift -gt 1e-12) { throw 'UP2_UNITARY_NORM_DRIFT' }
    if ([double]$Result.unitary_max_round_trip_error -gt 1e-12) { throw 'UP2_UNITARY_ROUND_TRIP_ERROR' }

    Write-Host $Probe1
    Write-Host 'WINGLESS_UP2_FOCUSED_PASS'
}
finally {
    Pop-Location
}
