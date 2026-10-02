$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up28-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up28-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-28 INTERACTION-AWARE COMMUTANT SELECTION ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP28_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP28_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP28_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP28_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-interaction-selected-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP28_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP28_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-interaction-selected-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP28_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-interaction-selected-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP28_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP28_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.unitary-interaction-selected.v1') { throw 'UP28_SCHEMA_MISMATCH' }
    if ([int]$Result.latent_dimension -ne 96) { throw 'UP28_LATENT_DIMENSION_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP28_NOT_SINGLE_RUNTIME_STATE' }
    if ($Result.visible_channel_blocks) { throw 'UP28_CHANNEL_BLOCKS_VISIBLE' }
    if (-not $Result.full_coordinate_mixing) { throw 'UP28_FULL_MIXING_MISSING' }
    if (-not $Result.observable_discovery_from_transport) { throw 'UP28_TRANSPORT_DISCOVERY_MISSING' }
    if ($Result.discovery_uses_hidden_multiplicity) { throw 'UP28_HIDDEN_MULTIPLICITY_USED' }
    if ($Result.discovery_uses_hidden_mixer) { throw 'UP28_HIDDEN_MIXER_USED' }
    if (-not $Result.selector_uses_training_labels) { throw 'UP28_TRAINING_SELECTOR_MISSING' }
    if ($Result.selector_uses_heldout_data) { throw 'UP28_HELDOUT_SELECTOR_LEAK' }
    if ($Result.selector_uses_explicit_depth) { throw 'UP28_DEPTH_SELECTOR_LEAK' }
    if (-not $Result.selector_uses_quadratic_interactions) { throw 'UP28_INTERACTION_SELECTOR_MISSING' }
    if ([int]$Result.candidate_observable_count -ne 128) { throw 'UP28_CANDIDATE_COUNT_MISMATCH' }
    if ([int]$Result.runtime_observable_count -ne 64) { throw 'UP28_RUNTIME_COUNT_MISMATCH' }
    if ([int]$Result.projection_rounds -ne 20) { throw 'UP28_ROUND_COUNT_MISMATCH' }
    if ([int]$Result.effective_orbit_size -ne 1048576) { throw 'UP28_ORBIT_SIZE_MISMATCH' }
    if ([int]$Result.selected_indices.Count -ne 64) { throw 'UP28_SELECTED_INDEX_COUNT_MISMATCH' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Prefix reproduction pass: $($Result.diagnosis.prefix64_reproduction_pass)"
    Write-Host "Training-only selector:   $($Result.diagnosis.selector_training_only_pass)"
    Write-Host "Discovery commutator pass:$($Result.diagnosis.discovery_commutator_pass)"
    Write-Host "Feature invariance pass:  $($Result.diagnosis.feature_invariance_pass)"
    Write-Host "Phase-code learning pass: $($Result.diagnosis.phase_code_learning_pass)"
    Write-Host "Unseen-depth pass:        $($Result.diagnosis.unseen_depth_pass)"
    Write-Host "Mutable integration pass: $($Result.diagnosis.mutable_integration_pass)"
    Write-Host "Prefix held delta:        $($Result.diagnosis.prefix64_heldout_delta)"
    Write-Host "Overlap prefix64:         $($Result.diagnosis.selected_overlap_with_prefix64)"
    Write-Host "Overlap UP27:             $($Result.diagnosis.selected_overlap_with_up27)"
    Write-Host "Selected train accuracy:  $($Result.diagnosis.selected_train_accuracy)"
    Write-Host "Selected held accuracy:   $($Result.diagnosis.selected_heldout_accuracy)"
    Write-Host "Matched control accuracy: $($Result.diagnosis.matched_control_accuracy)"
    Write-Host "Mean held phase cosine:   $($Result.diagnosis.mean_held_phase_cosine)"
    Write-Host 'WINGLESS_UP28_HARNESS_PASS'
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
