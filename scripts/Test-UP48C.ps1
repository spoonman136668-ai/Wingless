$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up48c-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up48c-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-48C JOINT DIMENSION/DEPTH STRESS ==='
    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP48C_GO_ENV_FAILED' }
    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP48C_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP48C_ICE_BUILD_FAILED' }
    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP48C_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-up48c-scale-stress -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP48C_FOCUSED_TEST_FAILED' }
    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP48C_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-up48c-scale-stress) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP48C_PROBE1_FAILED' }
    $Probe2 = ((go run ./cmd/unitary-up48c-scale-stress) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP48C_PROBE2_FAILED' }
    if ($Probe1 -cne $Probe2) { throw 'UP48C_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.up48c-joint-scale-stress.v1') { throw 'UP48C_SCHEMA_MISMATCH' }
    if ([int]$Result.cases.Count -ne 4) { throw 'UP48C_CASE_COUNT_MISMATCH' }
    $ExpectedDims = @(16,32,64,128)
    $ExpectedDepths = @(0,128,512,2048)
    for ($i = 0; $i -lt 4; $i++) {
        if ([int]$Result.cases[$i].dimension -ne $ExpectedDims[$i]) { throw "UP48C_DIMENSION_MISMATCH_$i" }
        if ([int]$Result.cases[$i].depths.Count -ne 4) { throw "UP48C_DEPTH_COUNT_MISMATCH_$i" }
        for ($j = 0; $j -lt 4; $j++) {
            if ([int]$Result.cases[$i].depths[$j] -ne $ExpectedDepths[$j]) { throw "UP48C_DEPTH_MISMATCH_$($i)_$j" }
        }
    }

    Write-Host $Probe1
    Write-Host ''
    Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
    Write-Host "All scale gates: $($Result.all_scale_gates)"
    foreach ($Case in $Result.cases) {
        $Last = $Case.unitary.metrics[$Case.unitary.metrics.Count - 1]
        Write-Host "dim=$($Case.dimension) depth=$($Last.depth) acc=$($Last.accuracy) norm=$($Last.max_norm_drift) gram=$($Last.max_gram_error) perturb=$($Last.max_perturbation_gain) roundtrip=$($Case.unitary_max_round_trip_error) gate=$($Case.unitary_scale_gate)"
    }
    Write-Host 'WINGLESS_UP48_HARNESS_PASS'
}
finally {
    Pop-Location
    if ($null -eq $PriorGoCache) { Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue } else { $env:GOCACHE = $PriorGoCache }
    if ($null -eq $PriorGoTmp) { Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue } else { $env:GOTMPDIR = $PriorGoTmp }
    Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
    Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
