$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up33-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up33-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-33 EQUAL COMMUTANT-CAPACITY PARTITION CONTROL ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP33_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP33_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP33_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP33_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-commutant-capacity-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP33_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP33_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-commutant-capacity-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP33_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-commutant-capacity-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP33_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP33_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.commutant-capacity-control.v1') { throw 'UP33_SCHEMA_MISMATCH' }
    if ([int]$Result.latent_dimension -ne 96) { throw 'UP33_LATENT_DIMENSION_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP33_NOT_SINGLE_RUNTIME_STATE' }
    if (-not $Result.full_coordinate_mixing) { throw 'UP33_FULL_MIXING_MISSING' }
    if (-not $Result.split_directions_anonymous) { throw 'UP33_SPLIT_DIRECTIONS_NOT_ANONYMOUS' }
    if ([double]$Result.fixed_nonzero_offset_rms -ne 0.05) { throw 'UP33_OFFSET_RMS_MISMATCH' }
    if (-not $Result.all_arms_norm_preserving) { throw 'UP33_NORM_PRESERVATION_MISSING' }
    if (-not $Result.all_arms_real_orthogonal_equivalent) { throw 'UP33_REAL_EQUIVALENCE_MISSING' }
    if (-not $Result.ablation_uses_known_partition) { throw 'UP33_ABLATION_BOUNDARY_MISSING' }
    if (-not $Result.observer_discovery_from_transport) { throw 'UP33_DISCOVERY_MISSING' }
    if ($Result.observer_uses_known_factorization) { throw 'UP33_OBSERVER_FACTOR_LEAK' }
    if (-not $Result.selector_uses_training_labels) { throw 'UP33_TRAINING_SELECTOR_MISSING' }
    if ($Result.selector_uses_heldout_data) { throw 'UP33_HELDOUT_LEAK' }
    if (-not $Result.phase_alphabet_supervision) { throw 'UP33_PHASE_SUPERVISION_MISSING' }
    if ([int]$Result.candidate_observable_count -ne 128) { throw 'UP33_CANDIDATE_COUNT_MISMATCH' }
    if ([int]$Result.runtime_observable_count -ne 64) { throw 'UP33_RUNTIME_COUNT_MISMATCH' }
    if ([int]$Result.permutation_names.Count -ne 3) { throw 'UP33_PERMUTATION_COUNT_MISMATCH' }
    if ([int]$Result.variants.Count -ne 12) { throw 'UP33_VARIANT_COUNT_MISMATCH' }
    if ([int]$Result.aggregates.Count -ne 4) { throw 'UP33_AGGREGATE_COUNT_MISMATCH' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Structural validity pass:      $($Result.diagnosis.structural_validity_pass)"
    Write-Host "Base control pass:             $($Result.diagnosis.base_control_pass)"
    Write-Host "Capacity-18 equivalence pass:  $($Result.diagnosis.capacity_18_equivalence_pass)"
    Write-Host "Capacity-12 equivalence pass:  $($Result.diagnosis.capacity_12_equivalence_pass)"
    Write-Host "Higher-capacity ordering pass: $($Result.diagnosis.higher_capacity_ordering_pass)"
    Write-Host "Capacity hypothesis supported: $($Result.diagnosis.capacity_hypothesis_supported)"
    Write-Host "C18 held difference:           $($Result.diagnosis.capacity_18_held_difference)"
    Write-Host "C18 commit difference:         $($Result.diagnosis.capacity_18_commit_difference)"
    Write-Host "C12 held difference:           $($Result.diagnosis.capacity_12_held_difference)"
    Write-Host "C12 commit difference:         $($Result.diagnosis.capacity_12_commit_difference)"
    Write-Host "Mean C18 held:                 $($Result.diagnosis.mean_capacity_18_held)"
    Write-Host "Mean C12 held:                 $($Result.diagnosis.mean_capacity_12_held)"
    Write-Host "Mean C18 commit:               $($Result.diagnosis.mean_capacity_18_commit)"
    Write-Host "Mean C12 commit:               $($Result.diagnosis.mean_capacity_12_commit)"
    Write-Host 'WINGLESS_UP33_HARNESS_PASS'
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
