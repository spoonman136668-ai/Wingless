$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up5-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up5-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-5 RELATIONAL READ/WRITE MEMORY ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP5_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP5_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP5_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP5_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-memory-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP5_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP5_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-memory-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP5_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-memory-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP5_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP5_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.unitary-memory-probe.v1') { throw 'UP5_SCHEMA_MISMATCH' }
    if ([double]$Result.unitary.commit_decode_accuracy -ne 1.0) { throw 'UP5_UNITARY_COMMIT_ACCURACY' }
    if ([double]$Result.unitary.exact_final_table_accuracy -ne 1.0) { throw 'UP5_UNITARY_FINAL_TABLE_ACCURACY' }
    if ([double]$Result.unitary.relational_query_accuracy -ne 1.0) { throw 'UP5_UNITARY_RELATION_ACCURACY' }
    if ([double]$Result.unitary.min_decode_margin -le 0.1) { throw 'UP5_UNITARY_DECODE_MARGIN' }
    if ([double]$Result.unitary.max_norm_drift -gt 1e-12) { throw 'UP5_UNITARY_NORM_DRIFT' }

    Write-Host $Probe1
    Write-Host 'WINGLESS_UP5_FOCUSED_PASS'
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
