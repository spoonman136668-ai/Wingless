$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up13-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up13-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-13 TASK-LEARNED PHASE FRAME ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP13_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP13_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP13_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP13_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-learned-phase-frame-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP13_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP13_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-learned-phase-frame-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP13_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-learned-phase-frame-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP13_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP13_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.unitary-learned-phase-frame.v1') { throw 'UP13_SCHEMA_MISMATCH' }
    if ($Result.runtime_prototype_lookup) { throw 'UP13_PROTOTYPE_LOOKUP_ENABLED' }
    if ($Result.explicit_inverse_readout) { throw 'UP13_INVERSE_READOUT_ENABLED' }
    if ($Result.explicit_depth_provided) { throw 'UP13_DEPTH_PROVIDED' }
    if (-not $Result.global_phase_nuisance) { throw 'UP13_GLOBAL_PHASE_NUISANCE_MISSING' }
    if (-not $Result.phase_alphabet_learned) { throw 'UP13_PHASE_LEARNING_DISABLED' }
    if ($Result.entity_support_learned) { throw 'UP13_SUPPORT_SCOPE_EXPANDED' }
    if ($Result.anchor_learned) { throw 'UP13_ANCHOR_SCOPE_EXPANDED' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Phase learning pass:       $($Result.diagnosis.phase_learning_pass)"
    Write-Host "Unseen-depth pass:         $($Result.diagnosis.unseen_depth_pass)"
    Write-Host "Mutable integration pass:  $($Result.diagnosis.mutable_integration_pass)"
    Write-Host "Initial capacity accuracy: $($Result.diagnosis.initial_capacity_accuracy)"
    Write-Host "Final held-out accuracy:   $($Result.diagnosis.final_heldout_accuracy)"
    Write-Host "Matched control accuracy:  $($Result.diagnosis.matched_control_accuracy)"
    Write-Host "Initial min separation:    $($Result.diagnosis.initial_min_phase_separation)"
    Write-Host "Learned min separation:    $($Result.diagnosis.learned_min_phase_separation)"
    Write-Host 'WINGLESS_UP13_HARNESS_PASS'
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
