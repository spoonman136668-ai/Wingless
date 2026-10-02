$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up20-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up20-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-20 CLOSED-FORM ANONYMOUS QUADRATIC READOUT ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP20_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP20_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP20_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP20_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-ridge-anonymous-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP20_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP20_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-ridge-anonymous-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP20_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-ridge-anonymous-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP20_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP20_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.unitary-ridge-anonymous.v1') { throw 'UP20_SCHEMA_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP20_NOT_SINGLE_RUNTIME_STATE' }
    if ([int]$Result.anonymous_channels -ne 6) { throw 'UP20_CHANNEL_COUNT_MISMATCH' }
    if ($Result.semantic_channel_locations_known) { throw 'UP20_SEMANTIC_LOCATIONS_EXPOSED' }
    if ($Result.known_mixer_exposed_to_learner) { throw 'UP20_MIXER_EXPOSED' }
    if ($Result.oracle_witness_used_for_learning) { throw 'UP20_ORACLE_LEAKED_INTO_LEARNING' }
    if ($Result.runtime_mixer_inverse_applied) { throw 'UP20_RUNTIME_INVERSE_APPLIED' }
    if ($Result.learned_demixer_used) { throw 'UP20_DEMIXER_USED' }
    if (-not $Result.direct_anonymous_readout) { throw 'UP20_DIRECT_READOUT_MISSING' }
    if ($Result.runtime_prototype_lookup) { throw 'UP20_PROTOTYPE_LOOKUP_ENABLED' }
    if ($Result.explicit_inverse_transport_readout) { throw 'UP20_TRANSPORT_INVERSE_ENABLED' }
    if ($Result.explicit_depth_provided) { throw 'UP20_DEPTH_PROVIDED' }
    if (-not $Result.memory_only_global_phase_nuisance) { throw 'UP20_GLOBAL_PHASE_NUISANCE_MISSING' }
    if ([int]$Result.quadratic_feature_dimension -ne 702) { throw 'UP20_FEATURE_DIMENSION_MISMATCH' }
    if ([double]$Result.ridge_lambda -ne 0.000001) { throw 'UP20_RIDGE_LAMBDA_MISMATCH' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Oracle witness pass:         $($Result.diagnosis.oracle_witness_pass)"
    Write-Host "Closed-form learning pass:   $($Result.diagnosis.closed_form_learning_pass)"
    Write-Host "Unseen-depth pass:           $($Result.diagnosis.unseen_depth_pass)"
    Write-Host "Mutable integration pass:    $($Result.diagnosis.mutable_integration_pass)"
    Write-Host "Oracle witness max error:    $($Result.diagnosis.oracle_witness_max_error)"
    Write-Host "Oracle witness held accuracy:$($Result.diagnosis.oracle_witness_held_accuracy)"
    Write-Host "Ridge train accuracy:        $($Result.diagnosis.ridge_train_accuracy)"
    Write-Host "Ridge held accuracy:         $($Result.diagnosis.ridge_held_accuracy)"
    Write-Host "Matched control accuracy:    $($Result.diagnosis.matched_control_accuracy)"
    Write-Host 'WINGLESS_UP20_HARNESS_PASS'
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
