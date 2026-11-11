$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up25-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up25-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-25 DYNAMICS-DISCOVERED COMMUTING OBSERVABLES ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP25_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP25_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP25_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP25_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-discovered-commutant-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP25_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP25_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-discovered-commutant-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP25_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-discovered-commutant-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP25_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP25_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.unitary-discovered-commutant.v1') { throw 'UP25_SCHEMA_MISMATCH' }
    if ([int]$Result.latent_dimension -ne 96) { throw 'UP25_LATENT_DIMENSION_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP25_NOT_SINGLE_RUNTIME_STATE' }
    if ($Result.visible_channel_blocks) { throw 'UP25_CHANNEL_BLOCKS_VISIBLE' }
    if (-not $Result.full_coordinate_mixing) { throw 'UP25_FULL_MIXING_MISSING' }
    if (-not $Result.observable_discovery_from_transport) { throw 'UP25_TRANSPORT_DISCOVERY_MISSING' }
    if ($Result.discovery_uses_hidden_multiplicity) { throw 'UP25_HIDDEN_MULTIPLICITY_USED' }
    if ($Result.discovery_uses_hidden_mixer) { throw 'UP25_HIDDEN_MIXER_USED' }
    if (-not $Result.transport_adjoint_used_in_discovery) { throw 'UP25_DISCOVERY_ADJOINT_FLAG_MISSING' }
    if ($Result.runtime_adjoint_applied) { throw 'UP25_RUNTIME_ADJOINT_ENABLED' }
    if ([int]$Result.discovery_seed_count -ne 32) { throw 'UP25_SEED_COUNT_MISMATCH' }
    if ([int]$Result.discovery_projection_rounds -ne 20) { throw 'UP25_ROUND_COUNT_MISMATCH' }
    if ([int]$Result.discovery_effective_orbit_size -ne 1048576) { throw 'UP25_ORBIT_SIZE_MISMATCH' }
    if ($Result.known_full_mixer_exposed_to_learner) { throw 'UP25_MIXER_EXPOSED' }
    if ($Result.runtime_unmix_applied) { throw 'UP25_RUNTIME_UNMIX_ENABLED' }
    if ($Result.learned_demixer_used) { throw 'UP25_DEMIXER_USED' }
    if (-not $Result.phase_alphabet_supervision) { throw 'UP25_PHASE_SUPERVISION_MISSING' }
    if ([int]$Result.phase_code_dimension -ne 2) { throw 'UP25_PHASE_DIMENSION_MISMATCH' }
    if ([int]$Result.raw_feature_dimension -ne 64) { throw 'UP25_RAW_FEATURE_DIMENSION_MISMATCH' }
    if ([int]$Result.quadratic_feature_dimension -ne 2144) { throw 'UP25_QUADRATIC_DIMENSION_MISMATCH' }
    if ($Result.runtime_prototype_lookup) { throw 'UP25_PROTOTYPE_LOOKUP_ENABLED' }
    if ($Result.explicit_inverse_transport_readout) { throw 'UP25_TRANSPORT_INVERSE_ENABLED' }
    if ($Result.explicit_depth_provided) { throw 'UP25_DEPTH_PROVIDED' }
    if (-not $Result.memory_only_global_phase_nuisance) { throw 'UP25_GLOBAL_PHASE_NUISANCE_MISSING' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Discovery commutator pass: $($Result.diagnosis.discovery_commutator_pass)"
    Write-Host "Feature invariance pass:   $($Result.diagnosis.feature_invariance_pass)"
    Write-Host "Phase-code learning pass:  $($Result.diagnosis.phase_code_learning_pass)"
    Write-Host "Unseen-depth pass:         $($Result.diagnosis.unseen_depth_pass)"
    Write-Host "Mutable integration pass:  $($Result.diagnosis.mutable_integration_pass)"
    Write-Host "Max commutator error:      $($Result.diagnosis.max_commutator_entry_error)"
    Write-Host "Max feature drift:         $($Result.diagnosis.max_feature_drift)"
    Write-Host "Min projected norm:        $($Result.diagnosis.min_pre_normalization_frobenius)"
    Write-Host "Train accuracy:            $($Result.diagnosis.train_accuracy)"
    Write-Host "Held-out accuracy:         $($Result.diagnosis.held_out_accuracy)"
    Write-Host "Matched control accuracy:  $($Result.diagnosis.matched_control_accuracy)"
    Write-Host "Mean held phase cosine:    $($Result.diagnosis.mean_held_phase_cosine)"
    Write-Host 'WINGLESS_UP25_HARNESS_PASS'
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
