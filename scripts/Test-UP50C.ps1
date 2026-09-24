$ErrorActionPreference = 'Stop'
$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp = $env:GOTMPDIR
$GoCache = Join-Path $env:TEMP ("wingless-up50c-gocache-" + $PID)
$GoTmp = Join-Path $env:TEMP ("wingless-up50c-gotmp-" + $PID)
New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp
Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-50C POWERED MILLION-DEPTH ==='
    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP50C_GO_ENV_FAILED' }
    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP50C_GO_CLEAN_FAILED' }
    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP50C_ICE_BUILD_FAILED' }
    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP50C_ICE_VALIDATE_FAILED' }
    go test ./unitary ./cmd/unitary-up50c-powered-depth -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP50C_FOCUSED_TEST_FAILED' }
    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP50C_FULL_REGRESSION_FAILED' }
    $Probe1 = ((go run ./cmd/unitary-up50c-powered-depth) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP50C_PROBE1_FAILED' }
    $Probe2 = ((go run ./cmd/unitary-up50c-powered-depth) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP50C_PROBE2_FAILED' }
    if ($Probe1 -cne $Probe2) { throw 'UP50C_NONDETERMINISTIC_OUTPUT' }
    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.up50c-powered-depth.v1') { throw 'UP50C_SCHEMA_MISMATCH' }
    if ([int]$Result.overlap_dimension -ne 64 -or [int]$Result.overlap_depth -ne 2048) { throw 'UP50C_OVERLAP_CONFIG' }
    if ([int]$Result.deep_depth -ne 1048576) { throw 'UP50C_DEEP_DEPTH' }
    if ($Result.overlap_equivalence_gate -and [int]$Result.deep_metrics.Count -ne 2) { throw 'UP50C_DEEP_COUNT' }

    Write-Host $Probe1
    Write-Host ''
    Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
    Write-Host "Overlap max state error: $($Result.overlap_max_state_error)"
    Write-Host "Overlap equivalence:      $($Result.overlap_equivalence_gate)"
    Write-Host "All million-depth gates:  $($Result.all_deep_gates)"
    foreach ($Metric in $Result.deep_metrics) {
        Write-Host "dim=$($Metric.dimension) depth=$($Metric.depth) acc=$($Metric.accuracy) norm=$($Metric.max_norm_drift) gram=$($Metric.max_gram_error) perturb=$($Metric.max_perturbation_gain) roundtrip=$($Metric.max_round_trip_error) gate=$($Metric.gate)"
    }
    Write-Host 'WINGLESS_UP50_HARNESS_PASS'
}
finally {
    Pop-Location
    if ($null -eq $PriorGoCache) { Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue } else { $env:GOCACHE = $PriorGoCache }
    if ($null -eq $PriorGoTmp) { Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue } else { $env:GOTMPDIR = $PriorGoTmp }
    Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
    Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
