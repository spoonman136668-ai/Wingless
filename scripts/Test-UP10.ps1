$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up10-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up10-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-10 COMPRESSED / NOISY INTERNAL FRAME ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP10_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP10_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP10_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP10_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-compressed-frame-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP10_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP10_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-compressed-frame-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP10_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-compressed-frame-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP10_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP10_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.unitary-compressed-frame.v1') { throw 'UP10_SCHEMA_MISMATCH' }
    if ($Result.runtime_prototype_lookup) { throw 'UP10_PROTOTYPE_LOOKUP_ENABLED' }
    if ($Result.explicit_inverse_readout) { throw 'UP10_INVERSE_READOUT_ENABLED' }
    if ($Result.explicit_depth_provided) { throw 'UP10_DEPTH_PROVIDED' }
    if ([int]$Result.full_pilot_states -ne 5) { throw 'UP10_FULL_PILOT_COUNT' }
    if ([int]$Result.compressed_pilot_states -ne 3) { throw 'UP10_COMPRESSED_PILOT_COUNT' }
    if ([int]$Result.features_per_entity -ne 1) { throw 'UP10_FEATURE_BUDGET' }
    if ([int]$Result.min_train_marginal_count -le 0 -or [int]$Result.min_held_out_marginal_count -le 0) {
        throw 'UP10_BALANCED_SPLIT_INVALID'
    }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Compression pass:       $($Result.diagnosis.compression_pass)"
    Write-Host "Noisy frame pass:       $($Result.diagnosis.noisy_frame_pass)"
    Write-Host "Mutable program pass:   $($Result.diagnosis.mutable_program_pass)"
    Write-Host "Pilot-state reduction:  $($Result.diagnosis.pilot_state_reduction)"
    Write-Host "Drift 1e-4 accuracy:    $($Result.diagnosis.reference_drift_1e_4_accuracy)"
    Write-Host "Drift 1e-3 accuracy:    $($Result.diagnosis.reference_drift_1e_3_accuracy)"
    Write-Host 'WINGLESS_UP10_HARNESS_PASS'
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
