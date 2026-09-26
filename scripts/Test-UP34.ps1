$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up34-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up34-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-34 TASK-ALLOCATED COMMUTANT CAPACITY ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP34_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP34_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP34_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP34_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-task-allocated-symmetry-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP34_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP34_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-task-allocated-symmetry-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP34_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-task-allocated-symmetry-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP34_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP34_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.task-allocated-symmetry.v1') { throw 'UP34_SCHEMA_MISMATCH' }
    if ([int]$Result.latent_dimension -ne 96) { throw 'UP34_LATENT_DIMENSION_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP34_NOT_SINGLE_RUNTIME_STATE' }
    if (-not $Result.full_coordinate_mixing) { throw 'UP34_FULL_MIXING_MISSING' }
    if (-not $Result.candidate_menu_hand_specified) { throw 'UP34_MENU_BOUNDARY_MISSING' }
    if (-not $Result.allocation_uses_training_only) { throw 'UP34_TRAINING_ONLY_FLAG_MISSING' }
    if ($Result.selection_uses_heldout_data) { throw 'UP34_HELDOUT_SELECTION_LEAK' }
    if ([double]$Result.capacity_price -ne 0.02) { throw 'UP34_CAPACITY_PRICE_MISMATCH' }
    if ([int]$Result.true_heldout_tables -ne 128) { throw 'UP34_TRUE_HELDOUT_POOL_MISMATCH' }
    if ([int]$Result.candidates.Count -ne 22) { throw 'UP34_CANDIDATE_COUNT_MISMATCH' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Inner split balanced:       $($Result.diagnosis.inner_split_balanced_pass)"
    Write-Host "Selected capacity:          $($Result.diagnosis.selected_capacity)"
    Write-Host "Max available capacity:     $($Result.diagnosis.max_available_capacity)"
    Write-Host "Capacity reduction:         $($Result.diagnosis.capacity_reduction_fraction)"
    Write-Host "Selected validation score:  $($Result.diagnosis.selected_validation_score)"
    Write-Host "Selected held-out accuracy: $($Result.diagnosis.selected_heldout_accuracy)"
    Write-Host "Selected commit accuracy:   $($Result.diagnosis.selected_commit_accuracy)"
    Write-Host "Full held-out accuracy:     $($Result.diagnosis.full_capacity_heldout_accuracy)"
    Write-Host "Full commit accuracy:       $($Result.diagnosis.full_capacity_commit_accuracy)"
    Write-Host "Held retention delta:       $($Result.diagnosis.heldout_retention_delta)"
    Write-Host "Commit retention delta:     $($Result.diagnosis.commit_retention_delta)"
    Write-Host "Task allocation pass:       $($Result.diagnosis.task_allocation_pass)"
    Write-Host 'WINGLESS_UP34_HARNESS_PASS'
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
