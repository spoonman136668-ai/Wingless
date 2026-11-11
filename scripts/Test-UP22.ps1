$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up22-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up22-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-22 FULL 96D LATENT MIXING ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP22_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP22_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP22_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP22_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-full-latent-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP22_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP22_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-full-latent-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP22_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-full-latent-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP22_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP22_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.unitary-full-latent-mix.v1') { throw 'UP22_SCHEMA_MISMATCH' }
    if ([int]$Result.latent_dimension -ne 96) { throw 'UP22_LATENT_DIMENSION_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP22_NOT_SINGLE_RUNTIME_STATE' }
    if ($Result.visible_channel_blocks) { throw 'UP22_CHANNEL_BLOCKS_VISIBLE' }
    if (-not $Result.full_coordinate_mixing) { throw 'UP22_FULL_MIXING_MISSING' }
    if (-not $Result.conjugated_transport) { throw 'UP22_CONJUGATED_TRANSPORT_MISSING' }
    if ($Result.known_full_mixer_exposed_to_learner) { throw 'UP22_MIXER_EXPOSED' }
    if ($Result.oracle_unmix_used_for_learning) { throw 'UP22_ORACLE_LEAKED' }
    if ($Result.runtime_unmix_applied) { throw 'UP22_RUNTIME_UNMIX_ENABLED' }
    if ($Result.learned_demixer_used) { throw 'UP22_DEMIXER_USED' }
    if (-not $Result.phase_alphabet_supervision) { throw 'UP22_PHASE_SUPERVISION_MISSING' }
    if ([int]$Result.phase_code_dimension -ne 2) { throw 'UP22_PHASE_DIMENSION_MISMATCH' }
    if ([int]$Result.full_hermitian_feature_dimension -ne 9216) { throw 'UP22_FEATURE_DIMENSION_MISMATCH' }
    if ($Result.runtime_prototype_lookup) { throw 'UP22_PROTOTYPE_LOOKUP_ENABLED' }
    if ($Result.explicit_inverse_transport_readout) { throw 'UP22_TRANSPORT_INVERSE_ENABLED' }
    if ($Result.explicit_depth_provided) { throw 'UP22_DEPTH_PROVIDED' }
    if (-not $Result.memory_only_global_phase_nuisance) { throw 'UP22_GLOBAL_PHASE_NUISANCE_MISSING' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Conjugation equivalence pass:$($Result.diagnosis.conjugation_equivalence_pass)"
    Write-Host "Full latent learning pass:   $($Result.diagnosis.full_latent_learning_pass)"
    Write-Host "Unseen-depth pass:           $($Result.diagnosis.unseen_depth_pass)"
    Write-Host "Mutable integration pass:    $($Result.diagnosis.mutable_integration_pass)"
    Write-Host "Oracle recovery error:       $($Result.diagnosis.max_oracle_recovery_error)"
    Write-Host "Train accuracy:              $($Result.diagnosis.train_accuracy)"
    Write-Host "Held-out accuracy:           $($Result.diagnosis.held_out_accuracy)"
    Write-Host "Matched control accuracy:    $($Result.diagnosis.matched_control_accuracy)"
    Write-Host "Mean held phase cosine:      $($Result.diagnosis.mean_held_phase_cosine)"
    Write-Host "Minimum mixer participation: $($Result.diagnosis.minimum_mixer_participation_ratio)"
    Write-Host 'WINGLESS_UP22_HARNESS_PASS'
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
