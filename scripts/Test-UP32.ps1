$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp = $env:GOTMPDIR
$GoCache = Join-Path $env:TEMP ('wingless-up32-gocache-' + $PID)
$GoTmp = Join-Path $env:TEMP ('wingless-up32-gotmp-' + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-32 MULTIPLICITY DOSE RESPONSE ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP32_GO_ENV_FAILED' }
    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP32_GO_CLEAN_FAILED' }
    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP32_ICE_BUILD_FAILED' }
    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP32_ICE_VALIDATE_FAILED' }
    go test ./unitary ./cmd/unitary-multiplicity-dose-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP32_FOCUSED_TEST_FAILED' }
    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP32_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-multiplicity-dose-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP32_PROBE1_FAILED' }
    $Probe2 = ((go run ./cmd/unitary-multiplicity-dose-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP32_PROBE2_FAILED' }
    if ($Probe1 -cne $Probe2) { throw 'UP32_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.multiplicity-dose.v1') { throw 'UP32_SCHEMA_MISMATCH' }
    if ([int]$Result.latent_dimension -ne 96) { throw 'UP32_LATENT_DIMENSION_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP32_NOT_SINGLE_RUNTIME_STATE' }
    if (-not $Result.full_coordinate_mixing) { throw 'UP32_FULL_MIXING_MISSING' }
    if (-not $Result.split_directions_anonymous) { throw 'UP32_SPLIT_DIRECTIONS_NOT_ANONYMOUS' }
    if ([double]$Result.fixed_nonzero_offset_rms -ne 0.05) { throw 'UP32_RMS_MISMATCH' }
    if (-not $Result.all_arms_norm_preserving) { throw 'UP32_NORM_PRESERVATION_MISSING' }
    if (-not $Result.all_arms_real_orthogonal_equivalent) { throw 'UP32_REAL_EQUIVALENCE_MISSING' }
    if (-not $Result.observer_discovery_from_transport) { throw 'UP32_DISCOVERY_MISSING' }
    if ($Result.observer_uses_known_factorization) { throw 'UP32_FACTORIZATION_LEAK' }
    if (-not $Result.selector_uses_training_labels) { throw 'UP32_TRAINING_SELECTOR_MISSING' }
    if ($Result.selector_uses_heldout_data) { throw 'UP32_HELDOUT_LEAK' }
    if (-not $Result.phase_alphabet_supervision) { throw 'UP32_PHASE_SUPERVISION_MISSING' }
    if ([int]$Result.candidate_observable_count -ne 128) { throw 'UP32_CANDIDATE_COUNT_MISMATCH' }
    if ([int]$Result.runtime_observable_count -ne 64) { throw 'UP32_RUNTIME_COUNT_MISMATCH' }
    if (@($Result.arms).Count -ne 5) { throw 'UP32_ARM_COUNT_MISMATCH' }

    $ExpectedMultiplicity = @(6, 4, 3, 2, 1)
    for ($i = 0; $i -lt 5; $i++) {
        if ([int]$Result.arms[$i].max_multiplicity -ne $ExpectedMultiplicity[$i]) {
            throw "UP32_MULTIPLICITY_MISMATCH index=$i"
        }
    }

    Write-Host $Probe1
    Write-Host ''
    Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
    Write-Host "Structural validity pass:   $($Result.diagnosis.structural_validity_pass)"
    Write-Host "Base control pass:          $($Result.diagnosis.base_control_pass)"
    Write-Host "Singleton material effect:  $($Result.diagnosis.singleton_material_effect_pass)"
    Write-Host "Held ordered steps:         $($Result.diagnosis.held_ordered_steps)"
    Write-Host "Commit ordered steps:       $($Result.diagnosis.commit_ordered_steps)"
    Write-Host "Multiplicity dose response: $($Result.diagnosis.multiplicity_dose_response_pass)"
    Write-Host "Held endpoint drop:         $($Result.diagnosis.held_endpoint_drop)"
    Write-Host "Commit endpoint drop:       $($Result.diagnosis.commit_endpoint_drop)"
    foreach ($Arm in $Result.arms) {
        Write-Host ('ARM m={0} held={1} commit={2} final={3} relation={4}' -f $Arm.max_multiplicity, $Arm.static.held_out_accuracy, $Arm.mutable_integration.commit_decode_accuracy, $Arm.mutable_integration.exact_final_table_accuracy, $Arm.mutable_integration.relational_query_accuracy)
    }
    Write-Host 'WINGLESS_UP32_HARNESS_PASS'
}
finally {
    Pop-Location
    if ($null -eq $PriorGoCache) { Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue } else { $env:GOCACHE = $PriorGoCache }
    if ($null -eq $PriorGoTmp) { Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue } else { $env:GOTMPDIR = $PriorGoTmp }
    Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
    Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
