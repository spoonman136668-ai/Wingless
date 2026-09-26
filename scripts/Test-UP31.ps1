$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up31-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up31-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-31 ISOSPECTRAL HIDDEN PER-COPY MISALIGNMENT ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP31_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP31_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP31_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP31_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-isospectral-misalignment-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP31_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP31_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-isospectral-misalignment-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP31_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-isospectral-misalignment-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP31_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP31_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.isospectral-misalignment.v1') { throw 'UP31_SCHEMA_MISMATCH' }
    if ([int]$Result.latent_dimension -ne 96) { throw 'UP31_LATENT_DIMENSION_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP31_NOT_SINGLE_RUNTIME_STATE' }
    if (-not $Result.full_coordinate_mixing) { throw 'UP31_FULL_MIXING_MISSING' }
    if (-not $Result.sixfold_spectrum_preserved) { throw 'UP31_SPECTRUM_NOT_PRESERVED' }
    if (-not $Result.hidden_per_copy_basis_change) { throw 'UP31_HIDDEN_BASIS_CHANGE_MISSING' }
    if (-not $Result.encoder_transport_conjugated_together) { throw 'UP31_CONSISTENT_CONJUGACY_MISSING' }
    if ($Result.observer_uses_hidden_basis) { throw 'UP31_HIDDEN_BASIS_LEAK' }
    if (-not $Result.observer_discovery_from_transport) { throw 'UP31_DISCOVERY_MISSING' }
    if (-not $Result.selector_uses_training_labels) { throw 'UP31_TRAINING_SELECTOR_MISSING' }
    if ($Result.selector_uses_heldout_data) { throw 'UP31_HELDOUT_LEAK' }
    if ($Result.selector_uses_explicit_depth) { throw 'UP31_DEPTH_LEAK' }
    if (-not $Result.phase_alphabet_supervision) { throw 'UP31_PHASE_SUPERVISION_MISSING' }
    if ([int]$Result.candidate_observable_count -ne 128) { throw 'UP31_CANDIDATE_COUNT_MISMATCH' }
    if ([int]$Result.runtime_observable_count -ne 64) { throw 'UP31_RUNTIME_COUNT_MISMATCH' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Base control pass:               $($Result.diagnosis.base_control_pass)"
    Write-Host "Misaligned orthogonality pass:   $($Result.diagnosis.misaligned_orthogonality_pass)"
    Write-Host "Isospectral equivalence pass:    $($Result.diagnosis.isospectral_equivalence_pass)"
    Write-Host "Naive alignment break pass:      $($Result.diagnosis.naive_alignment_break_pass)"
    Write-Host "Correct intertwiner pass:        $($Result.diagnosis.correct_intertwiner_pass)"
    Write-Host "Discovery commutator pass:       $($Result.diagnosis.discovery_commutator_pass)"
    Write-Host "Feature invariance pass:         $($Result.diagnosis.feature_invariance_pass)"
    Write-Host "Misaligned training pass:        $($Result.diagnosis.misaligned_training_pass)"
    Write-Host "Misaligned unseen-depth pass:    $($Result.diagnosis.misaligned_unseen_depth_pass)"
    Write-Host "Misaligned mutable pass:         $($Result.diagnosis.misaligned_mutable_pass)"
    Write-Host "Alignment not required:          $($Result.diagnosis.alignment_not_required_supported)"
    Write-Host "Naive shift commutator:          $($Result.diagnosis.naive_shift_commutator)"
    Write-Host "Correct shift commutator:        $($Result.diagnosis.correct_shift_commutator)"
    Write-Host "Misaligned held-out accuracy:    $($Result.misaligned_arm.held_out_accuracy)"
    Write-Host "Misaligned commit accuracy:      $($Result.misaligned_mutable_integration.commit_decode_accuracy)"
    Write-Host 'WINGLESS_UP31_HARNESS_PASS'
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
