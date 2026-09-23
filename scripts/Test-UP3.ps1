$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-3 DEPTH RETENTION STRESS ==='
    Write-Host "Repo: $Repo"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP3_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP3_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP3_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP3_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-stress-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP3_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP3_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-stress-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP3_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-stress-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP3_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP3_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.unitary-stress-probe.v1') { throw 'UP3_SCHEMA_MISMATCH' }

    foreach ($Metric in $Result.unitary.metrics) {
        if ([double]$Metric.accuracy -ne 1.0) { throw "UP3_UNITARY_ACCURACY depth=$($Metric.depth)" }
        if ([double]$Metric.max_norm_drift -gt 1e-12) { throw "UP3_UNITARY_NORM_DRIFT depth=$($Metric.depth)" }
        if ([double]$Metric.max_gram_error -gt 1e-12) { throw "UP3_UNITARY_GRAM_ERROR depth=$($Metric.depth)" }
        if ([math]::Abs([double]$Metric.max_perturbation_gain - 1.0) -gt 1e-12) { throw "UP3_UNITARY_PERTURBATION_GAIN depth=$($Metric.depth)" }
    }

    if ([double]$Result.unitary_max_round_trip_error -gt 1e-11) {
        throw 'UP3_UNITARY_ROUND_TRIP_ERROR'
    }

    Write-Host $Probe1
    Write-Host 'WINGLESS_UP3_FOCUSED_PASS'
}
finally {
    Pop-Location
}
