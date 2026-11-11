$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-1 TRAINED UNITARY PROPAGATOR ==='
    Write-Host "Repo: $Repo"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP1_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP1_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP1_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP1_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-train-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP1_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP1_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-train-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP1_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-train-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP1_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP1_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.unitary-training-probe.v1') { throw 'UP1_SCHEMA_MISMATCH' }
    if ([double]$Result.initial_accuracy -ne 0.5) { throw 'UP1_INITIAL_ACCURACY_MISMATCH' }
    if ([double]$Result.train_accuracy -ne 1.0) { throw 'UP1_TRAIN_ACCURACY_MISMATCH' }
    if ([double]$Result.held_out_accuracy -ne 1.0) { throw 'UP1_HELDOUT_ACCURACY_MISMATCH' }
    if ([double]$Result.final_loss -gt 1e-5) { throw 'UP1_FINAL_LOSS_MISMATCH' }
    if ([math]::Abs([double]$Result.learned_theta - ([math]::PI / 4)) -gt 1e-3) { throw 'UP1_THETA_MISMATCH' }
    if ([double]$Result.max_norm_drift -gt 1e-12) { throw 'UP1_NORM_DRIFT' }
    if ([double]$Result.max_round_trip_error -gt 1e-12) { throw 'UP1_ROUND_TRIP_ERROR' }

    Write-Host $Probe1
    Write-Host 'WINGLESS_UP1_FOCUSED_PASS'
}
finally {
    Pop-Location
}
