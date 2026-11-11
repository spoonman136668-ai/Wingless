$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up19-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up19-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-19 DIRECT ANONYMOUS GRAM READOUT ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP19_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP19_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP19_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP19_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-anonymous-gram-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP19_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP19_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-anonymous-gram-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP19_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-anonymous-gram-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP19_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP19_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.unitary-anonymous-gram.v1') { throw 'UP19_SCHEMA_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP19_NOT_SINGLE_RUNTIME_STATE' }
    if ([int]$Result.anonymous_channels -ne 6) { throw 'UP19_CHANNEL_COUNT_MISMATCH' }
    if ($Result.semantic_channel_locations_known) { throw 'UP19_SEMANTIC_LOCATIONS_EXPOSED' }
    if ($Result.known_mixer_exposed_to_learner) { throw 'UP19_MIXER_EXPOSED' }
    if ($Result.runtime_mixer_inverse_applied) { throw 'UP19_RUNTIME_INVERSE_APPLIED' }
    if ($Result.learned_demixer_used) { throw 'UP19_DEMIXER_USED' }
    if (-not $Result.direct_anonymous_readout) { throw 'UP19_DIRECT_READOUT_MISSING' }
    if ($Result.runtime_prototype_lookup) { throw 'UP19_PROTOTYPE_LOOKUP_ENABLED' }
    if ($Result.explicit_inverse_transport_readout) { throw 'UP19_TRANSPORT_INVERSE_ENABLED' }
    if ($Result.explicit_depth_provided) { throw 'UP19_DEPTH_PROVIDED' }
    if (-not $Result.memory_only_global_phase_nuisance) { throw 'UP19_GLOBAL_PHASE_NUISANCE_MISSING' }
    if (-not $Result.canonical_training_only) { throw 'UP19_TRAINING_SCOPE_EXPANDED' }
    if ([int]$Result.linear_feature_dimension -ne 36) { throw 'UP19_LINEAR_DIMENSION_MISMATCH' }
    if ([int]$Result.quadratic_feature_dimension -ne 702) { throw 'UP19_QUADRATIC_DIMENSION_MISMATCH' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Feature invariance pass:      $($Result.diagnosis.feature_invariance_pass)"
    Write-Host "Direct quadratic readout pass:$($Result.diagnosis.direct_quadratic_readout_pass)"
    Write-Host "Unseen-depth pass:            $($Result.diagnosis.unseen_depth_pass)"
    Write-Host "Mutable integration pass:     $($Result.diagnosis.mutable_integration_pass)"
    Write-Host "Linear held-out accuracy:     $($Result.diagnosis.linear_heldout_accuracy)"
    Write-Host "Quadratic train accuracy:     $($Result.diagnosis.quadratic_train_accuracy)"
    Write-Host "Quadratic held-out accuracy:  $($Result.diagnosis.quadratic_heldout_accuracy)"
    Write-Host "Matched control accuracy:     $($Result.diagnosis.matched_control_accuracy)"
    Write-Host "Unitary feature drift:        $($Result.diagnosis.max_unitary_feature_drift)"
    Write-Host "Non-unitary feature drift:    $($Result.diagnosis.max_nonunitary_feature_drift)"
    Write-Host 'WINGLESS_UP19_HARNESS_PASS'
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
