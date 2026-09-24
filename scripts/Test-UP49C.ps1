$ErrorActionPreference = 'Stop'
$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp = $env:GOTMPDIR
$GoCache = Join-Path $env:TEMP ("wingless-up49c-gocache-" + $PID)
$GoTmp = Join-Path $env:TEMP ("wingless-up49c-gotmp-" + $PID)
New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp
Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-49C EXTENDED SCALE BOUNDARY SCREEN ==='
    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP49C_GO_ENV_FAILED' }
    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP49C_GO_CLEAN_FAILED' }
    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP49C_ICE_BUILD_FAILED' }
    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP49C_ICE_VALIDATE_FAILED' }
    go test ./unitary ./cmd/unitary-up49c-scale-boundary -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP49C_FOCUSED_TEST_FAILED' }
    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP49C_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-up49c-scale-boundary) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP49C_PROBE1_FAILED' }
    $Probe2 = ((go run ./cmd/unitary-up49c-scale-boundary) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP49C_PROBE2_FAILED' }
    if ($Probe1 -cne $Probe2) { throw 'UP49C_NONDETERMINISTIC_OUTPUT' }
    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.up49c-extended-scale-boundary.v1') { throw 'UP49C_SCHEMA_MISMATCH' }
    if (-not $Result.screening_only) { throw 'UP49C_NOT_SCREENING' }
    if ($Result.matched_control_run) { throw 'UP49C_UNPLANNED_CONTROL' }
    if ([int]$Result.trials_per_class -ne 2) { throw 'UP49C_TRIAL_COUNT' }
    if ([int]$Result.cases.Count -ne 4) { throw 'UP49C_CASE_COUNT' }

    Write-Host $Probe1
    Write-Host ''
    Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
    Write-Host "All gates:          $($Result.all_gates)"
    Write-Host "First failing case: $($Result.first_failing_case)"
    foreach ($Case in $Result.cases) {
        $Last = $Case.unitary.metrics[$Case.unitary.metrics.Count - 1]
        Write-Host "dim=$($Case.spec.dimension) depth=$($Case.spec.depth) acc=$($Last.accuracy) norm=$($Last.max_norm_drift) gram=$($Last.max_gram_error) perturb=$($Last.max_perturbation_gain) roundtrip=$($Case.unitary_max_round_trip_error) gate=$($Case.gate)"
    }
    Write-Host 'WINGLESS_UP49_HARNESS_PASS'
}
finally {
    Pop-Location
    if ($null -eq $PriorGoCache) { Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue } else { $env:GOCACHE = $PriorGoCache }
    if ($null -eq $PriorGoTmp) { Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue } else { $env:GOTMPDIR = $PriorGoTmp }
    Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
    Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
