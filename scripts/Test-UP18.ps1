$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up18-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up18-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-18 ORTHOGONAL ANONYMOUS-CHANNEL DEMIXER ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP18_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP18_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP18_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP18_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-orthogonal-demixer-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP18_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP18_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-orthogonal-demixer-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP18_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-orthogonal-demixer-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP18_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP18_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.unitary-orthogonal-demixer.v1') { throw 'UP18_SCHEMA_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP18_NOT_SINGLE_RUNTIME_STATE' }
    if ($Result.semantic_channel_locations_known) { throw 'UP18_SEMANTIC_LOCATIONS_EXPOSED' }
    if ($Result.known_mixer_exposed_to_learner) { throw 'UP18_MIXER_EXPOSED_TO_LEARNER' }
    if ($Result.oracle_used_for_learning) { throw 'UP18_ORACLE_LEAKED_INTO_LEARNING' }
    if ($Result.runtime_mixer_inverse_applied) { throw 'UP18_RUNTIME_INVERSE_APPLIED' }
    if (-not $Result.learned_readout_projection) { throw 'UP18_READOUT_NOT_LEARNED' }
    if (-not $Result.readout_projection_orthogonal) { throw 'UP18_PROJECTION_NOT_ORTHOGONAL' }
    if ([int]$Result.trainable_demixer_parameters -ne 15) { throw 'UP18_PARAMETER_COUNT_MISMATCH' }
    if ($Result.runtime_prototype_lookup) { throw 'UP18_PROTOTYPE_LOOKUP_ENABLED' }
    if ($Result.explicit_inverse_transport_readout) { throw 'UP18_TRANSPORT_INVERSE_ENABLED' }
    if ($Result.explicit_depth_provided) { throw 'UP18_DEPTH_PROVIDED' }
    if (-not $Result.memory_only_global_phase_nuisance) { throw 'UP18_GLOBAL_PHASE_NUISANCE_MISSING' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Oracle recoverability pass: $($Result.diagnosis.oracle_recoverability_pass)"
    Write-Host "Structured learning pass:   $($Result.diagnosis.structured_learning_pass)"
    Write-Host "Unseen-depth pass:          $($Result.diagnosis.unseen_depth_pass)"
    Write-Host "Mutable integration pass:   $($Result.diagnosis.mutable_integration_pass)"
    Write-Host "Initial capacity accuracy:  $($Result.diagnosis.initial_capacity_accuracy)"
    Write-Host "Final held-out accuracy:    $($Result.diagnosis.final_heldout_accuracy)"
    Write-Host "Oracle held-out accuracy:   $($Result.diagnosis.oracle_heldout_accuracy)"
    Write-Host "Matched control accuracy:   $($Result.diagnosis.matched_control_accuracy)"
    Write-Host "Oracle feature error:       $($Result.diagnosis.oracle_feature_error)"
    Write-Host "Mean role alignment:        $($Result.diagnosis.mean_role_alignment)"
    Write-Host "Minimum role alignment:     $($Result.diagnosis.minimum_role_alignment)"
    Write-Host "Angle L2 shift:             $($Result.diagnosis.angle_l2_shift)"
    Write-Host "WINGLESS_UP18_HARNESS_PASS"
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
