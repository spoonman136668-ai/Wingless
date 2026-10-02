$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up26-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up26-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-26 DISCOVERED OBSERVABLE BREADTH ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP26_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP26_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP26_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP26_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-discovery-breadth-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP26_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP26_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-discovery-breadth-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP26_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-discovery-breadth-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP26_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP26_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.unitary-discovery-breadth.v1') { throw 'UP26_SCHEMA_MISMATCH' }
    if ([int]$Result.latent_dimension -ne 96) { throw 'UP26_LATENT_DIMENSION_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP26_NOT_SINGLE_RUNTIME_STATE' }
    if ($Result.visible_channel_blocks) { throw 'UP26_CHANNEL_BLOCKS_VISIBLE' }
    if (-not $Result.full_coordinate_mixing) { throw 'UP26_FULL_MIXING_MISSING' }
    if (-not $Result.observable_discovery_from_transport) { throw 'UP26_TRANSPORT_DISCOVERY_MISSING' }
    if ($Result.discovery_uses_hidden_multiplicity) { throw 'UP26_HIDDEN_MULTIPLICITY_USED' }
    if ($Result.discovery_uses_hidden_mixer) { throw 'UP26_HIDDEN_MIXER_USED' }
    if (-not $Result.transport_adjoint_used_in_discovery) { throw 'UP26_DISCOVERY_ADJOINT_FLAG_MISSING' }
    if ($Result.runtime_adjoint_applied) { throw 'UP26_RUNTIME_ADJOINT_ENABLED' }
    if ([int]$Result.projection_rounds -ne 20) { throw 'UP26_ROUND_COUNT_MISMATCH' }
    if ([int]$Result.effective_orbit_size -ne 1048576) { throw 'UP26_ORBIT_SIZE_MISMATCH' }
    if ([int]$Result.baseline_32.observable_count -ne 32) { throw 'UP26_BASELINE_COUNT_MISMATCH' }
    if ([int]$Result.baseline_32.quadratic_feature_dimension -ne 2144) { throw 'UP26_BASELINE_DIMENSION_MISMATCH' }
    if ([int]$Result.primary_64.observable_count -ne 64) { throw 'UP26_PRIMARY_COUNT_MISMATCH' }
    if ([int]$Result.primary_64.quadratic_feature_dimension -ne 8384) { throw 'UP26_PRIMARY_DIMENSION_MISMATCH' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Baseline reproduction pass:$($Result.diagnosis.baseline_reproduction_pass)"
    Write-Host "Discovery commutator pass: $($Result.diagnosis.discovery_commutator_pass)"
    Write-Host "Feature invariance pass:   $($Result.diagnosis.feature_invariance_pass)"
    Write-Host "Phase-code learning pass:  $($Result.diagnosis.phase_code_learning_pass)"
    Write-Host "Unseen-depth pass:         $($Result.diagnosis.unseen_depth_pass)"
    Write-Host "Mutable integration pass:  $($Result.diagnosis.mutable_integration_pass)"
    Write-Host "Baseline held delta:       $($Result.diagnosis.baseline_heldout_delta)"
    Write-Host "Primary train accuracy:    $($Result.diagnosis.primary_train_accuracy)"
    Write-Host "Primary held accuracy:     $($Result.diagnosis.primary_heldout_accuracy)"
    Write-Host "Matched control accuracy:  $($Result.diagnosis.matched_control_accuracy)"
    Write-Host "Mean held phase cosine:    $($Result.diagnosis.mean_held_phase_cosine)"
    Write-Host 'WINGLESS_UP26_HARNESS_PASS'
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
