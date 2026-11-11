$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up16-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up16-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-16 SINGLE COMPOSITE RECURRENT STATE ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP16_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP16_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP16_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP16_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-composite-state-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP16_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP16_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-composite-state-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP16_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-composite-state-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP16_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP16_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.unitary-composite-state.v1') { throw 'UP16_SCHEMA_MISMATCH' }
    if ([int]$Result.runtime_state_objects -ne 1) { throw 'UP16_NOT_SINGLE_RUNTIME_STATE' }
    if ($Result.frame_side_vectors_materialized_after_pack) { throw 'UP16_SIDE_VECTORS_STILL_MATERIALIZED' }
    if (-not $Result.channel_subspaces_explicit) { throw 'UP16_SCOPE_EXPANDED_BEYOND_CHANNEL_FOLD' }
    if ($Result.runtime_prototype_lookup) { throw 'UP16_PROTOTYPE_LOOKUP_ENABLED' }
    if ($Result.explicit_inverse_readout) { throw 'UP16_INVERSE_READOUT_ENABLED' }
    if ($Result.explicit_depth_provided) { throw 'UP16_DEPTH_PROVIDED' }
    if (-not $Result.memory_only_global_phase_nuisance) { throw 'UP16_GLOBAL_PHASE_NUISANCE_MISSING' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Observable equivalence pass: $($Result.diagnosis.observable_equivalence_pass)"
    Write-Host "Static fold pass:             $($Result.diagnosis.static_fold_pass)"
    Write-Host "Mutable integration pass:     $($Result.diagnosis.mutable_integration_pass)"
    Write-Host "Max observable error:         $($Result.diagnosis.max_observable_error)"
    Write-Host "Unitary held-out accuracy:    $($Result.diagnosis.unitary_heldout_accuracy)"
    Write-Host "Matched control accuracy:     $($Result.diagnosis.matched_control_accuracy)"
    Write-Host "Unitary max norm drift:       $($Result.diagnosis.unitary_max_norm_drift)"
    Write-Host 'WINGLESS_UP16_HARNESS_PASS'
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
