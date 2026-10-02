$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-0 UNITARY INTERFERENCE ROUTING ==='
    Write-Host "Repo: $Repo"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP0_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP0_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP0_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP0_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP0_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP0_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP0_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP0_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP0_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.unitary-probe.v1') { throw 'UP0_SCHEMA_MISMATCH' }
    if ([double]$Result.baseline_accuracy -ne 0.5) { throw 'UP0_BASELINE_ACCURACY_MISMATCH' }
    if ([double]$Result.unitary_accuracy -ne 1.0) { throw 'UP0_UNITARY_ACCURACY_MISMATCH' }
    if ([double]$Result.max_norm_drift -gt 1e-12) { throw 'UP0_NORM_DRIFT' }
    if ([double]$Result.max_round_trip_error -gt 1e-12) { throw 'UP0_ROUND_TRIP_ERROR' }

    Write-Host $Probe1
    Write-Host 'WINGLESS_UP0_FOCUSED_PASS'
}
finally {
    Pop-Location
}
