$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up30-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up30-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-30 ORTHOGONAL TRANSPORT MULTIPLICITY ABLATION ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP30_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP30_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP30_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP30_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-multiplicity-ablation-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP30_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP30_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-multiplicity-ablation-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP30_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-multiplicity-ablation-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP30_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP30_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.multiplicity-ablation.v1') { throw 'UP30_SCHEMA_MISMATCH' }
    if ([int]$Result.latent_dimension -ne 96) { throw 'UP30_LATENT_DIMENSION_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP30_NOT_SINGLE_RUNTIME_STATE' }
    if (-not $Result.full_coordinate_mixing) { throw 'UP30_FULL_MIXING_MISSING' }
    if (-not $Result.both_transports_norm_preserving) { throw 'UP30_NORM_PRESERVATION_MISSING' }
    if (-not $Result.both_have_real_orthogonal_equivalent) { throw 'UP30_REAL_EQUIVALENCE_MISSING' }
    if (-not $Result.ablation_uses_known_factorization) { throw 'UP30_ABLATION_BOUNDARY_MISSING' }
    if ($Result.observer_uses_known_factorization) { throw 'UP30_OBSERVER_FACTOR_LEAK' }
    if (-not $Result.observer_discovery_from_transport) { throw 'UP30_DISCOVERY_MISSING' }
    if (-not $Result.phase_alphabet_supervision) { throw 'UP30_PHASE_SUPERVISION_MISSING' }
    if ([int]$Result.candidate_observable_count -ne 128) { throw 'UP30_CANDIDATE_COUNT_MISMATCH' }
    if ([int]$Result.runtime_observable_count -ne 64) { throw 'UP30_RUNTIME_COUNT_MISMATCH' }
    if ([int]$Result.split_offsets_radians_per_step.Count -ne 6) { throw 'UP30_SPLIT_OFFSET_COUNT_MISMATCH' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Base control pass:              $($Result.diagnosis.base_control_pass)"
    Write-Host "Base orthogonality pass:        $($Result.diagnosis.base_orthogonality_pass)"
    Write-Host "Split orthogonality pass:       $($Result.diagnosis.split_orthogonality_pass)"
    Write-Host "Degeneracy break pass:          $($Result.diagnosis.degeneracy_break_pass)"
    Write-Host "Split discovery commutator pass:$($Result.diagnosis.split_discovery_commutator_pass)"
    Write-Host "Split feature invariance pass:  $($Result.diagnosis.split_feature_invariance_pass)"
    Write-Host "Split training pass:            $($Result.diagnosis.split_training_pass)"
    Write-Host "Split unseen-depth pass:        $($Result.diagnosis.split_unseen_depth_pass)"
    Write-Host "Split mutable pass:             $($Result.diagnosis.split_mutable_pass)"
    Write-Host "Material ablation effect:       $($Result.diagnosis.material_ablation_effect)"
    Write-Host "Multiplicity dependence:        $($Result.diagnosis.multiplicity_dependence_supported)"
    Write-Host "Base shift commutator:          $($Result.diagnosis.base_cross_shift_commutator)"
    Write-Host "Split shift commutator:         $($Result.diagnosis.split_cross_shift_commutator)"
    Write-Host "Held-out accuracy drop:         $($Result.diagnosis.heldout_accuracy_drop)"
    Write-Host "Mutable commit drop:            $($Result.diagnosis.mutable_commit_drop)"
    Write-Host 'WINGLESS_UP30_HARNESS_PASS'
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
