$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up37-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up37-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-37 ADAPTIVE CONTINUOUS SYMMETRY FUSION ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP37_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP37_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP37_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP37_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-adaptive-continuous-fusion-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP37_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP37_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-adaptive-continuous-fusion-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP37_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-adaptive-continuous-fusion-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP37_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP37_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.adaptive-continuous-fusion.v1') { throw 'UP37_SCHEMA_MISMATCH' }
    if ([int]$Result.latent_dimension -ne 96) { throw 'UP37_LATENT_DIMENSION_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP37_NOT_SINGLE_RUNTIME_STATE' }
    if (-not $Result.full_coordinate_mixing) { throw 'UP37_FULL_MIXING_MISSING' }
    if (-not $Result.starts_from_independent_offsets) { throw 'UP37_INDEPENDENT_START_MISSING' }
    if ($Result.finished_partition_menu_provided) { throw 'UP37_FINISHED_MENU_LEAK' }
    if ($Result.hard_merge_candidates_evaluated) { throw 'UP37_HARD_MERGE_SEARCH_PRESENT' }
    if (-not $Result.pair_probes_continuous_only) { throw 'UP37_CONTINUOUS_PROBE_BOUNDARY_MISSING' }
    if (-not $Result.affinity_recomputed_every_step) { throw 'UP37_ADAPTIVE_REESTIMATION_MISSING' }
    if (-not $Result.fusion_flow_simultaneous) { throw 'UP37_SIMULTANEOUS_FLOW_MISSING' }
    if (-not $Result.fusion_sticky_after_tolerance) { throw 'UP37_STICKY_FUSION_MISSING' }
    if ($Result.selection_uses_heldout_data) { throw 'UP37_HELDOUT_SELECTION_LEAK' }
    if ([double]$Result.capacity_price -ne 0.02) { throw 'UP37_CAPACITY_PRICE_MISMATCH' }
    if ([double]$Result.probe_fraction -ne 0.5) { throw 'UP37_PROBE_FRACTION_MISMATCH' }
    if ([double]$Result.flow_step_rate -ne 0.45) { throw 'UP37_FLOW_RATE_MISMATCH' }
    if ([double]$Result.fusion_tolerance -ne 0.0025) { throw 'UP37_FUSION_TOLERANCE_MISMATCH' }
    if ([int]$Result.maximum_flow_steps -ne 10) { throw 'UP37_FLOW_STEP_COUNT_MISMATCH' }
    if ([int]$Result.fit_tables -ne 96) { throw 'UP37_FIT_COUNT_MISMATCH' }
    if ([int]$Result.validation_tables -ne 32) { throw 'UP37_VALIDATION_COUNT_MISMATCH' }
    if ([int]$Result.true_heldout_tables -ne 128) { throw 'UP37_TRUE_HELDOUT_COUNT_MISMATCH' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Affinity re-estimations:       $($Result.diagnosis.affinity_reestimation_count)"
    Write-Host "Total positive affinities:     $($Result.diagnosis.total_positive_affinity_count)"
    Write-Host "Selected flow step:            $($Result.diagnosis.selected_step)"
    Write-Host "Selected capacity:             $($Result.diagnosis.selected_capacity)"
    Write-Host "Additional fusion beyond UP36: $($Result.diagnosis.additional_fusion_beyond_up36)"
    Write-Host "Selected held-out accuracy:    $($Result.diagnosis.selected_heldout_accuracy)"
    Write-Host "Selected commit accuracy:      $($Result.diagnosis.selected_commit_accuracy)"
    Write-Host "Full held-out accuracy:        $($Result.diagnosis.full_capacity_heldout_accuracy)"
    Write-Host "Full commit accuracy:          $($Result.diagnosis.full_capacity_commit_accuracy)"
    Write-Host "Held retention delta:          $($Result.diagnosis.heldout_retention_delta)"
    Write-Host "Commit retention delta:        $($Result.diagnosis.commit_retention_delta)"
    Write-Host "Adaptive fusion pass:          $($Result.diagnosis.adaptive_fusion_pass)"
    Write-Host 'WINGLESS_UP37_HARNESS_PASS'
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
