$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up6-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up6-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-6 LEARNED OBSERVATION / QUERY ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP6_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP6_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP6_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP6_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-readout-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP6_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP6_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-readout-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP6_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-readout-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP6_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP6_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.unitary-learned-readout-probe.v1') { throw 'UP6_SCHEMA_MISMATCH' }
    if ($Result.runtime_prototype_lookup) { throw 'UP6_PROTOTYPE_LOOKUP_STILL_ENABLED' }
    if ([double]$Result.value_decoder.train_accuracy -ne 1.0) { throw 'UP6_VALUE_DECODER_TRAIN_ACCURACY' }
    if ([double]$Result.relation_head.train_accuracy -ne 1.0) { throw 'UP6_RELATION_HEAD_TRAIN_ACCURACY' }
    if ([double]$Result.unitary.commit_decode_accuracy -ne 1.0) { throw 'UP6_UNITARY_COMMIT_ACCURACY' }
    if ([double]$Result.unitary.exact_final_table_accuracy -ne 1.0) { throw 'UP6_UNITARY_FINAL_TABLE_ACCURACY' }
    if ([double]$Result.unitary.relational_query_accuracy -ne 1.0) { throw 'UP6_UNITARY_RELATION_ACCURACY' }
    if ([double]$Result.unitary.max_round_trip_error -gt 1e-10) { throw 'UP6_UNITARY_ROUND_TRIP' }
    if ([double]$Result.unitary.max_forward_norm_drift -gt 1e-12) { throw 'UP6_UNITARY_NORM_DRIFT' }
    if ([double]$Result.unitary.min_value_margin -le 0.5) { throw 'UP6_UNITARY_VALUE_MARGIN' }
    if ([double]$Result.unitary.min_relation_margin -le 0.5) { throw 'UP6_UNITARY_RELATION_MARGIN' }

    Write-Host $Probe1
    Write-Host 'WINGLESS_UP6_FOCUSED_PASS'
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
