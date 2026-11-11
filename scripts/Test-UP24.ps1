$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up24-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up24-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-24 COMMUTING WEYL OBSERVABLE ALGEBRA ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP24_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP24_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP24_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP24_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-commuting-weyl-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP24_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP24_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-commuting-weyl-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP24_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-commuting-weyl-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP24_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP24_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.unitary-commuting-weyl.v1') { throw 'UP24_SCHEMA_MISMATCH' }
    if ([int]$Result.latent_dimension -ne 96) { throw 'UP24_LATENT_DIMENSION_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP24_NOT_SINGLE_RUNTIME_STATE' }
    if ($Result.visible_channel_blocks) { throw 'UP24_CHANNEL_BLOCKS_VISIBLE' }
    if (-not $Result.full_coordinate_mixing) { throw 'UP24_FULL_MIXING_MISSING' }
    if ([int]$Result.transport_multiplicity_dimension -ne 6) { throw 'UP24_MULTIPLICITY_MISMATCH' }
    if (-not $Result.commuting_observable_algebra) { throw 'UP24_COMMUTING_ALGEBRA_MISSING' }
    if (-not $Result.observable_algebra_hand_constructed) { throw 'UP24_HAND_CONSTRUCTED_FLAG_MISSING' }
    if ([int]$Result.weyl_observable_count -ne 36) { throw 'UP24_OBSERVABLE_COUNT_MISMATCH' }
    if ($Result.known_full_mixer_exposed_to_learner) { throw 'UP24_MIXER_EXPOSED' }
    if ($Result.runtime_unmix_applied) { throw 'UP24_RUNTIME_UNMIX_ENABLED' }
    if ($Result.learned_demixer_used) { throw 'UP24_DEMIXER_USED' }
    if (-not $Result.phase_alphabet_supervision) { throw 'UP24_PHASE_SUPERVISION_MISSING' }
    if ([int]$Result.phase_code_dimension -ne 2) { throw 'UP24_PHASE_DIMENSION_MISMATCH' }
    if ([int]$Result.raw_feature_dimension -ne 72) { throw 'UP24_RAW_FEATURE_DIMENSION_MISMATCH' }
    if ([int]$Result.quadratic_feature_dimension -ne 2700) { throw 'UP24_QUADRATIC_DIMENSION_MISMATCH' }
    if ($Result.runtime_prototype_lookup) { throw 'UP24_PROTOTYPE_LOOKUP_ENABLED' }
    if ($Result.explicit_inverse_transport_readout) { throw 'UP24_TRANSPORT_INVERSE_ENABLED' }
    if ($Result.explicit_depth_provided) { throw 'UP24_DEPTH_PROVIDED' }
    if (-not $Result.memory_only_global_phase_nuisance) { throw 'UP24_GLOBAL_PHASE_NUISANCE_MISSING' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Commutator pass:           $($Result.diagnosis.commutator_pass)"
    Write-Host "Feature invariance pass:   $($Result.diagnosis.feature_invariance_pass)"
    Write-Host "Phase-code learning pass:  $($Result.diagnosis.phase_code_learning_pass)"
    Write-Host "Unseen-depth pass:         $($Result.diagnosis.unseen_depth_pass)"
    Write-Host "Mutable integration pass:  $($Result.diagnosis.mutable_integration_pass)"
    Write-Host "Max commutator error:      $($Result.diagnosis.max_commutator_entry_error)"
    Write-Host "Max feature drift:         $($Result.diagnosis.max_feature_drift)"
    Write-Host "Train accuracy:            $($Result.diagnosis.train_accuracy)"
    Write-Host "Held-out accuracy:         $($Result.diagnosis.heldout_accuracy)"
    Write-Host "Matched control accuracy:  $($Result.diagnosis.matched_control_accuracy)"
    Write-Host "Mean held phase cosine:    $($Result.diagnosis.mean_held_phase_cosine)"
    Write-Host 'WINGLESS_UP24_HARNESS_PASS'
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
