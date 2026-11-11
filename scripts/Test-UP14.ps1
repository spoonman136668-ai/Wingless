$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up14-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up14-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-14 TASK-LEARNED PILOT SUPPORT ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP14_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP14_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP14_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP14_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-learned-support-frame-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP14_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP14_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-learned-support-frame-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP14_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-learned-support-frame-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP14_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP14_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.unitary-learned-support-frame.v1') { throw 'UP14_SCHEMA_MISMATCH' }
    if ($Result.runtime_prototype_lookup) { throw 'UP14_PROTOTYPE_LOOKUP_ENABLED' }
    if ($Result.explicit_inverse_readout) { throw 'UP14_INVERSE_READOUT_ENABLED' }
    if ($Result.explicit_depth_provided) { throw 'UP14_DEPTH_PROVIDED' }
    if (-not $Result.global_phase_nuisance) { throw 'UP14_GLOBAL_PHASE_NUISANCE_MISSING' }
    if ($Result.phase_alphabet_learned) { throw 'UP14_PHASE_SCOPE_EXPANDED' }
    if (-not $Result.entity_support_learned) { throw 'UP14_SUPPORT_LEARNING_DISABLED' }
    if ($Result.anchor_learned) { throw 'UP14_ANCHOR_SCOPE_EXPANDED' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Support learning pass:      $($Result.diagnosis.support_learning_pass)"
    Write-Host "Unseen-depth pass:          $($Result.diagnosis.unseen_depth_pass)"
    Write-Host "Mutable integration pass:   $($Result.diagnosis.mutable_integration_pass)"
    Write-Host "Initial capacity accuracy:  $($Result.diagnosis.initial_capacity_accuracy)"
    Write-Host "Final held-out accuracy:    $($Result.diagnosis.final_heldout_accuracy)"
    Write-Host "Matched control accuracy:   $($Result.diagnosis.matched_control_accuracy)"
    Write-Host "Initial mean purity:        $($Result.diagnosis.initial_mean_target_purity)"
    Write-Host "Learned mean purity:        $($Result.diagnosis.learned_mean_target_purity)"
    Write-Host "Minimum learned purity:     $($Result.diagnosis.minimum_learned_target_purity)"
    Write-Host 'WINGLESS_UP14_HARNESS_PASS'
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
