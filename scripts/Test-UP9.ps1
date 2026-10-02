$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up9-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up9-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-9 CO-EVOLVING INTERNAL FRAME ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP9_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP9_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP9_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP9_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-coevolving-frame-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP9_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP9_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-coevolving-frame-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP9_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-coevolving-frame-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP9_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP9_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.unitary-coevolving-frame.v1') { throw 'UP9_SCHEMA_MISMATCH' }
    if (-not $Result.diagnosis.balanced_split_valid) { throw 'UP9_BALANCED_SPLIT_INVALID' }
    if ($Result.runtime_prototype_lookup) { throw 'UP9_PROTOTYPE_LOOKUP_ENABLED' }
    if ($Result.explicit_inverse_readout) { throw 'UP9_INVERSE_READOUT_ENABLED' }
    if ($Result.explicit_depth_provided) { throw 'UP9_DEPTH_PROVIDED' }
    if (-not $Result.global_phase_nuisance) { throw 'UP9_GLOBAL_PHASE_NUISANCE_MISSING' }
    if ([int]$Result.features_per_entity -ne 2) { throw 'UP9_FEATURE_BUDGET_MISMATCH' }
    if ([int]$Result.pilot_states -ne 5) { throw 'UP9_PILOT_BUDGET_MISMATCH' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Unitary co-evolving pass:  $($Result.diagnosis.unitary_coevolving_unseen_depth_pass)"
    Write-Host "Unitary static pass:       $($Result.diagnosis.unitary_static_unseen_depth_pass)"
    Write-Host "Control co-evolving pass:  $($Result.diagnosis.control_coevolving_unseen_depth_pass)"
    Write-Host "Internal frame supported:  $($Result.diagnosis.internal_frame_supported)"
    Write-Host "Unitary frame gain:        $($Result.diagnosis.unitary_frame_gain)"
    Write-Host 'WINGLESS_UP9_HARNESS_PASS'
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
