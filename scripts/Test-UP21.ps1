$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up21-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up21-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-21 ANONYMOUS PHASE-CODE BOTTLENECK ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP21_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP21_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP21_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP21_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-phase-code-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP21_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP21_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-phase-code-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP21_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-phase-code-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP21_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP21_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.unitary-phase-code.v1') { throw 'UP21_SCHEMA_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP21_NOT_SINGLE_RUNTIME_STATE' }
    if ([int]$Result.anonymous_channels -ne 6) { throw 'UP21_CHANNEL_COUNT_MISMATCH' }
    if ($Result.semantic_channel_locations_known) { throw 'UP21_SEMANTIC_LOCATIONS_EXPOSED' }
    if ($Result.known_mixer_exposed_to_learner) { throw 'UP21_MIXER_EXPOSED' }
    if ($Result.oracle_used_for_learning) { throw 'UP21_ORACLE_LEAKED' }
    if ($Result.runtime_mixer_inverse_applied) { throw 'UP21_RUNTIME_INVERSE_APPLIED' }
    if ($Result.learned_demixer_used) { throw 'UP21_DEMIXER_USED' }
    if (-not $Result.direct_anonymous_readout) { throw 'UP21_DIRECT_READOUT_MISSING' }
    if (-not $Result.phase_alphabet_supervision) { throw 'UP21_PHASE_SUPERVISION_MISSING' }
    if ([int]$Result.phase_code_dimension -ne 2) { throw 'UP21_PHASE_DIMENSION_MISMATCH' }
    if ($Result.runtime_prototype_lookup) { throw 'UP21_PROTOTYPE_LOOKUP_ENABLED' }
    if ($Result.explicit_inverse_transport_readout) { throw 'UP21_TRANSPORT_INVERSE_ENABLED' }
    if ($Result.explicit_depth_provided) { throw 'UP21_DEPTH_PROVIDED' }
    if (-not $Result.memory_only_global_phase_nuisance) { throw 'UP21_GLOBAL_PHASE_NUISANCE_MISSING' }
    if ([int]$Result.quadratic_feature_dimension -ne 702) { throw 'UP21_FEATURE_DIMENSION_MISMATCH' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Phase-code learning pass:  $($Result.diagnosis.phase_code_learning_pass)"
    Write-Host "Unseen-depth pass:         $($Result.diagnosis.unseen_depth_pass)"
    Write-Host "Mutable integration pass:  $($Result.diagnosis.mutable_integration_pass)"
    Write-Host "Train accuracy:            $($Result.diagnosis.train_accuracy)"
    Write-Host "Held-out accuracy:         $($Result.diagnosis.held_out_accuracy)"
    Write-Host "Matched control accuracy:  $($Result.diagnosis.matched_control_accuracy)"
    Write-Host "Mean held phase cosine:    $($Result.diagnosis.mean_held_phase_cosine)"
    Write-Host 'WINGLESS_UP21_HARNESS_PASS'
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
