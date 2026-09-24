$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up46-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up46-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-46 OBJECTIVE DECOMPOSITION AUDIT ==='

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP46_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP46_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP46_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP46_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-objective-decomposition-audit -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP46_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP46_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-objective-decomposition-audit) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP46_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-objective-decomposition-audit) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP46_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP46_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.objective-decomposition-audit.v1') { throw 'UP46_SCHEMA_MISMATCH' }
    if ($Result.optimizer_run) { throw 'UP46_OPTIMIZER_PRESENT' }
    if ($Result.lambda_changed) { throw 'UP46_LAMBDA_CHANGED' }
    if ($Result.selection_uses_heldout_data) { throw 'UP46_HELDOUT_SELECTION_LEAK' }
    if (-not $Result.state_pair_frozen_from_up45) { throw 'UP46_STATE_PAIR_NOT_FROZEN' }
    if ([int]$Result.selected_step -ne 24) { throw 'UP46_SELECTED_STEP_MISMATCH' }
    if ([int]$Result.exact_fusion_step -ne 29) { throw 'UP46_FUSION_STEP_MISMATCH' }
    if ([int]$Result.selected_exact_capacity -ne 6) { throw 'UP46_SELECTED_CAPACITY_MISMATCH' }
    if ([int]$Result.fused_exact_capacity -ne 8) { throw 'UP46_FUSED_CAPACITY_MISMATCH' }

    Write-Host $Probe1
    Write-Host ''
    Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
    Write-Host "Objective delta fused-selected:  $($Result.smooth_delta_fused_minus_selected.objective)"
    Write-Host "Phase delta:                     $($Result.smooth_delta_fused_minus_selected.phase_score)"
    Write-Host "Value probability delta:         $($Result.smooth_delta_fused_minus_selected.value_probability)"
    Write-Host "Relation probability delta:      $($Result.smooth_delta_fused_minus_selected.relation_probability)"
    Write-Host "Harmonic task delta:             $($Result.smooth_delta_fused_minus_selected.harmonic_task_score)"
    Write-Host "Soft capacity delta:             $($Result.smooth_delta_fused_minus_selected.soft_capacity)"
    Write-Host "Resource penalty delta:          $($Result.smooth_delta_fused_minus_selected.resource_penalty)"
    Write-Host "Held-out accuracy delta:         $($Result.hard_delta_fused_minus_selected.heldout_accuracy)"
    Write-Host "Commit accuracy delta:           $($Result.hard_delta_fused_minus_selected.commit_accuracy)"
    Write-Host "Final table accuracy delta:      $($Result.hard_delta_fused_minus_selected.final_table_accuracy)"
    Write-Host "Relation accuracy delta:         $($Result.hard_delta_fused_minus_selected.relation_accuracy)"
    Write-Host 'WINGLESS_UP46_HARNESS_PASS'
}
finally {
    Pop-Location
    if ($null -eq $PriorGoCache) { Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue } else { $env:GOCACHE = $PriorGoCache }
    if ($null -eq $PriorGoTmp) { Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue } else { $env:GOTMPDIR = $PriorGoTmp }
    Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
    Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
