$ErrorActionPreference = 'Stop'
$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp = $env:GOTMPDIR
$GoCache = Join-Path $env:TEMP ("wingless-up49b-gocache-" + $PID)
$GoTmp = Join-Path $env:TEMP ("wingless-up49b-gotmp-" + $PID)
New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp
Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-49B SINGLE-PAIR FUSION SCREEN ==='
    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP49B_GO_ENV_FAILED' }
    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP49B_GO_CLEAN_FAILED' }
    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP49B_ICE_BUILD_FAILED' }
    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP49B_ICE_VALIDATE_FAILED' }
    go test ./unitary ./cmd/unitary-up49b-pair-fusion -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP49B_FOCUSED_TEST_FAILED' }
    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP49B_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-up49b-pair-fusion) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP49B_PROBE1_FAILED' }
    $Probe2 = ((go run ./cmd/unitary-up49b-pair-fusion) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP49B_PROBE2_FAILED' }
    if ($Probe1 -cne $Probe2) { throw 'UP49B_NONDETERMINISTIC_OUTPUT' }
    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.up49b-single-pair-fusion-screen.v1') { throw 'UP49B_SCHEMA_MISMATCH' }
    if ($Result.optimizer_run) { throw 'UP49B_OPTIMIZER_PRESENT' }
    if ($Result.selection_uses_heldout_data) { throw 'UP49B_HELDOUT_SELECTION_LEAK' }
    if ([int]$Result.candidate_count -ne 15) { throw 'UP49B_CANDIDATE_COUNT' }
    if ([int]$Result.candidates.Count -ne 15) { throw 'UP49B_CANDIDATE_ARRAY_COUNT' }

    Write-Host $Probe1
    Write-Host ''
    Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
    Write-Host "Selected pair:                 $($Result.diagnosis.selected_pair -join ',')"
    Write-Host "Selected objective gain:       $($Result.diagnosis.selected_objective_gain)"
    Write-Host "Selected hard mean:            $($Result.diagnosis.selected_hard_mean)"
    Write-Host "Selected hard gate:            $($Result.diagnosis.selected_hard_gate)"
    Write-Host "Posthoc best hard pair:        $($Result.diagnosis.posthoc_best_hard_pair -join ',')"
    Write-Host "Selected matches hard best:    $($Result.diagnosis.selected_matches_posthoc_best)"
    Write-Host "Smooth-hard Pearson:           $($Result.diagnosis.smooth_hard_pearson_correlation)"
    Write-Host 'WINGLESS_UP49_HARNESS_PASS'
}
finally {
    Pop-Location
    if ($null -eq $PriorGoCache) { Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue } else { $env:GOCACHE = $PriorGoCache }
    if ($null -eq $PriorGoTmp) { Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue } else { $env:GOTMPDIR = $PriorGoTmp }
    Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
    Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
