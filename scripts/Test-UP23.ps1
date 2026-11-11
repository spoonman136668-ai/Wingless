$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up23-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up23-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-23 FULL-LATENT TRANSPORT SPECTRAL INVARIANTS ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP23_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP23_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP23_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP23_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-spectral-invariants-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP23_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP23_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-spectral-invariants-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP23_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-spectral-invariants-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP23_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP23_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.unitary-spectral-invariants.v1') { throw 'UP23_SCHEMA_MISMATCH' }
    if ([int]$Result.latent_dimension -ne 96) { throw 'UP23_LATENT_DIMENSION_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP23_NOT_SINGLE_RUNTIME_STATE' }
    if ($Result.visible_channel_blocks) { throw 'UP23_CHANNEL_BLOCKS_VISIBLE' }
    if (-not $Result.full_coordinate_mixing) { throw 'UP23_FULL_MIXING_MISSING' }
    if (-not $Result.transport_operator_used_as_reference) { throw 'UP23_TRANSPORT_REFERENCE_MISSING' }
    if ($Result.known_full_mixer_exposed_to_learner) { throw 'UP23_MIXER_EXPOSED' }
    if ($Result.runtime_unmix_applied) { throw 'UP23_RUNTIME_UNMIX_ENABLED' }
    if ($Result.learned_demixer_used) { throw 'UP23_DEMIXER_USED' }
    if (-not $Result.phase_alphabet_supervision) { throw 'UP23_PHASE_SUPERVISION_MISSING' }
    if ([int]$Result.phase_code_dimension -ne 2) { throw 'UP23_PHASE_DIMENSION_MISMATCH' }
    if ([int]$Result.spectral_moment_count -ne 16) { throw 'UP23_MOMENT_COUNT_MISMATCH' }
    if ([int]$Result.spectral_feature_dimension -ne 32) { throw 'UP23_FEATURE_DIMENSION_MISMATCH' }
    if ($Result.runtime_prototype_lookup) { throw 'UP23_PROTOTYPE_LOOKUP_ENABLED' }
    if ($Result.explicit_inverse_transport_readout) { throw 'UP23_TRANSPORT_INVERSE_ENABLED' }
    if ($Result.explicit_depth_provided) { throw 'UP23_DEPTH_PROVIDED' }
    if (-not $Result.memory_only_global_phase_nuisance) { throw 'UP23_GLOBAL_PHASE_NUISANCE_MISSING' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Moment invariance pass:    $($Result.diagnosis.moment_invariance_pass)"
    Write-Host "Spectral learning pass:    $($Result.diagnosis.spectral_learning_pass)"
    Write-Host "Unseen-depth pass:         $($Result.diagnosis.unseen_depth_pass)"
    Write-Host "Mutable integration pass:  $($Result.diagnosis.mutable_integration_pass)"
    Write-Host "Max moment drift:          $($Result.diagnosis.max_moment_drift)"
    Write-Host "Train accuracy:            $($Result.diagnosis.train_accuracy)"
    Write-Host "Held-out accuracy:         $($Result.diagnosis.held_out_accuracy)"
    Write-Host "Matched control accuracy:  $($Result.diagnosis.matched_control_accuracy)"
    Write-Host "Mean held phase cosine:    $($Result.diagnosis.mean_held_phase_cosine)"
    Write-Host 'WINGLESS_UP23_HARNESS_PASS'
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
