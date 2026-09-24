$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up48b-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up48b-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-48B FROZEN GEOMETRY CAUSAL AUDIT ==='
    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP48B_GO_ENV_FAILED' }
    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP48B_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP48B_ICE_BUILD_FAILED' }
    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP48B_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-up48b-geometry-audit -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP48B_FOCUSED_TEST_FAILED' }
    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP48B_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-up48b-geometry-audit) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP48B_PROBE1_FAILED' }
    $Probe2 = ((go run ./cmd/unitary-up48b-geometry-audit) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP48B_PROBE2_FAILED' }
    if ($Probe1 -cne $Probe2) { throw 'UP48B_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.up48b-frozen-geometry-audit.v1') { throw 'UP48B_SCHEMA_MISMATCH' }
    if ($Result.optimizer_run) { throw 'UP48B_OPTIMIZER_PRESENT' }
    if ($Result.selection_uses_heldout_data) { throw 'UP48B_HELDOUT_SELECTION_LEAK' }
    if ([int]$Result.control_gap_source_step -ne 28) { throw 'UP48B_CONTROL_SOURCE_MISMATCH' }
    if ([int]$Result.arms.Count -ne 5) { throw 'UP48B_ARM_COUNT_MISMATCH' }
    if (-not $Result.arms[1].exact_fusion) { throw 'UP48B_STEP29_FUSION_MISSING' }
    if ($Result.arms[2].exact_fusion) { throw 'UP48B_MATCHED_SPLIT_STILL_FUSED' }
    if (-not $Result.arms[4].exact_fusion) { throw 'UP48B_STEP42_CONTROL_FUSION_MISSING' }

    Write-Host $Probe1
    Write-Host ''
    Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
    Write-Host "Step29 exact - matched split: $($Result.diagnosis.step29_minus_matched_split | ConvertTo-Json -Compress)"
    Write-Host "Step42 - step29:             $($Result.diagnosis.step42_minus_step29 | ConvertTo-Json -Compress)"
    Write-Host "Step42 fused - step42:       $($Result.diagnosis.step42_fused_minus_step42 | ConvertTo-Json -Compress)"
    Write-Host 'WINGLESS_UP48_HARNESS_PASS'
}
finally {
    Pop-Location
    if ($null -eq $PriorGoCache) { Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue } else { $env:GOCACHE = $PriorGoCache }
    if ($null -eq $PriorGoTmp) { Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue } else { $env:GOTMPDIR = $PriorGoTmp }
    Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
    Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
