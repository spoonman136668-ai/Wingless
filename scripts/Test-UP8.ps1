$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up8-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up8-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-8 CONFOUND-CONTROLLED OBSERVER ABLATION ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP8_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP8_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP8_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP8_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-observer-ablation-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP8_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP8_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-observer-ablation-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP8_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-observer-ablation-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP8_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP8_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.unitary-observer-ablation.v1') { throw 'UP8_SCHEMA_MISMATCH' }
    if (-not $Result.diagnosis.balanced_split_valid) { throw 'UP8_BALANCED_SPLIT_INVALID' }
    if ([int]$Result.min_train_marginal_count -le 0) { throw 'UP8_TRAIN_MARGINAL_MISSING' }
    if ([int]$Result.min_held_out_marginal_count -le 0) { throw 'UP8_HELD_MARGINAL_MISSING' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Measurement-loss supported:        $($Result.diagnosis.measurement_loss_supported)"
    Write-Host "Fixed coherence pass:              $($Result.diagnosis.fixed_coherence_pass)"
    Write-Host "Cross-depth coherence pass:        $($Result.diagnosis.cross_depth_coherence_pass)"
    Write-Host "Frame ambiguity supported:         $($Result.diagnosis.frame_ambiguity_supported)"
    Write-Host 'WINGLESS_UP8_HARNESS_PASS'
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
