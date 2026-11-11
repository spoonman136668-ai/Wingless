$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up7-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up7-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-7 DIRECT TRANSPORTED-STATE OBSERVER ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP7_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP7_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP7_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP7_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-direct-observer-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP7_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP7_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-direct-observer-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP7_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-direct-observer-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP7_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP7_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.unitary-direct-observer-probe.v1') { throw 'UP7_SCHEMA_MISMATCH' }
    if ($Result.runtime_prototype_lookup) { throw 'UP7_PROTOTYPE_LOOKUP_ENABLED' }
    if ($Result.explicit_inverse_readout) { throw 'UP7_INVERSE_READOUT_ENABLED' }
    if ($Result.depth_feature_provided) { throw 'UP7_DEPTH_FEATURE_ENABLED' }
    if ([int]$Result.train_tables -ne 128 -or [int]$Result.held_out_tables -ne 128) { throw 'UP7_TABLE_SPLIT_MISMATCH' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC GATE (does not control harness acceptance) ==="
    Write-Host "Unitary hypothesis pass:    $($Result.unitary.hypothesis_pass)"
    Write-Host "Control hypothesis pass:    $($Result.non_unitary_matched.hypothesis_pass)"
    Write-Host 'WINGLESS_UP7_HARNESS_PASS'
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
