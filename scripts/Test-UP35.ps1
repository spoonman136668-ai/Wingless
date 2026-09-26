$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up35-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up35-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-35 GREEDY TASK-DRIVEN SYMMETRY INDUCTION ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP35_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP35_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP35_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP35_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-greedy-symmetry-induction-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP35_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP35_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-greedy-symmetry-induction-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP35_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-greedy-symmetry-induction-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP35_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP35_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.greedy-symmetry-induction.v1') { throw 'UP35_SCHEMA_MISMATCH' }
    if ([int]$Result.latent_dimension -ne 96) { throw 'UP35_LATENT_DIMENSION_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP35_NOT_SINGLE_RUNTIME_STATE' }
    if (-not $Result.full_coordinate_mixing) { throw 'UP35_FULL_MIXING_MISSING' }
    if (-not $Result.starts_from_minimum_symmetry) { throw 'UP35_MINIMUM_START_MISSING' }
    if ($Result.finished_partition_menu_provided) { throw 'UP35_FINISHED_MENU_LEAK' }
    if (-not $Result.pairwise_merge_only) { throw 'UP35_PAIRWISE_MERGE_BOUNDARY_MISSING' }
    if (-not $Result.merge_selection_training_only) { throw 'UP35_TRAINING_ONLY_SELECTION_MISSING' }
    if ($Result.selection_uses_heldout_data) { throw 'UP35_HELDOUT_SELECTION_LEAK' }
    if ([double]$Result.capacity_price -ne 0.02) { throw 'UP35_CAPACITY_PRICE_MISMATCH' }
    if ([double]$Result.minimum_accepted_gain -ne 0.005) { throw 'UP35_GAIN_RULE_MISMATCH' }
    if ([int]$Result.fit_tables -ne 96) { throw 'UP35_FIT_COUNT_MISMATCH' }
    if ([int]$Result.validation_tables -ne 32) { throw 'UP35_VALIDATION_COUNT_MISMATCH' }
    if ([int]$Result.true_heldout_tables -ne 128) { throw 'UP35_TRUE_HELDOUT_COUNT_MISMATCH' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Accepted merges:            $($Result.diagnosis.accepted_merge_count)"
    Write-Host "Selected capacity:          $($Result.diagnosis.selected_capacity)"
    Write-Host "Capacity created fraction:  $($Result.diagnosis.capacity_created_fraction)"
    Write-Host "Stopped by frozen gain:     $($Result.diagnosis.stopped_by_frozen_gain_rule)"
    Write-Host "Selected held-out accuracy: $($Result.diagnosis.selected_heldout_accuracy)"
    Write-Host "Selected commit accuracy:   $($Result.diagnosis.selected_commit_accuracy)"
    Write-Host "Full held-out accuracy:     $($Result.diagnosis.full_capacity_heldout_accuracy)"
    Write-Host "Full commit accuracy:       $($Result.diagnosis.full_capacity_commit_accuracy)"
    Write-Host "Held retention delta:       $($Result.diagnosis.heldout_retention_delta)"
    Write-Host "Commit retention delta:     $($Result.diagnosis.commit_retention_delta)"
    Write-Host "Greedy induction pass:      $($Result.diagnosis.greedy_induction_pass)"
    Write-Host 'WINGLESS_UP35_HARNESS_PASS'
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
