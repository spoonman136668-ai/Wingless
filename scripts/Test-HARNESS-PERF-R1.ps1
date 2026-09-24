$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-hperf-r1-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-hperf-r1-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS HARNESS PERF EQUIVALENCE R1 ==='

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'HPERF_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'HPERF_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'HPERF_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'HPERF_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-harness-perf-equivalence-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'HPERF_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'HPERF_FULL_REGRESSION_FAILED' }

    $Probe = ((go run ./cmd/unitary-harness-perf-equivalence-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'HPERF_PROBE_FAILED' }

    $Result = $Probe | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.harness-perf-equivalence.v1') { throw 'HPERF_SCHEMA_MISMATCH' }
    if (-not $Result.reference_untouched) { throw 'HPERF_REFERENCE_NOT_UNTOUCHED' }
    if (-not $Result.only_invariant_reuse) { throw 'HPERF_SCOPE_VIOLATION' }
    if (-not $Result.exact_equivalence) { throw 'HPERF_SCIENTIFIC_MISMATCH' }
    if (-not $Result.optimized_deterministic) { throw 'HPERF_OPTIMIZED_NONDETERMINISTIC' }
    if ($Result.reference_sha256 -cne $Result.optimized_sha256) { throw 'HPERF_HASH_MISMATCH' }
    if ($Result.optimized_sha256 -cne $Result.second_optimized_sha256) { throw 'HPERF_REPLAY_HASH_MISMATCH' }
    if ([double]$Result.minimum_speedup -ne 1.05) { throw 'HPERF_MIN_SPEEDUP_MISMATCH' }
    if (-not $Result.speedup_gate_pass) { throw "HPERF_SPEEDUP_GATE_FAILED speedup=$($Result.speedup)" }
    if (-not $Result.accepted) { throw 'HPERF_NOT_ACCEPTED' }

    Write-Host $Probe
    Write-Host ""
    Write-Host "Reference ms: $($Result.reference_duration_millis)"
    Write-Host "Optimized ms: $($Result.optimized_duration_millis)"
    Write-Host "Optimized replay ms: $($Result.second_optimized_duration_millis)"
    Write-Host "Speedup: $($Result.speedup)x"
    Write-Host "Result SHA256: $($Result.reference_sha256)"
    Write-Host 'WINGLESS_HARNESS_PERF_EQUIVALENCE_R1_PASS'
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
