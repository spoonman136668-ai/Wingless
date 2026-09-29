$ErrorActionPreference = 'Stop'
$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$OldGOROOT = $env:GOROOT
$OldCache = $env:GOCACHE
$OldTmp = $env:GOTMPDIR
$GoExe = (Get-Command go -ErrorAction Stop).Source
$GoRoot = Split-Path (Split-Path $GoExe -Parent) -Parent
$Cache = Join-Path $env:TEMP ("up255c-cache-" + $PID)
$Tmp = Join-Path $env:TEMP ("up255c-tmp-" + $PID)
New-Item -ItemType Directory -Force -Path $Cache,$Tmp | Out-Null
$env:GOROOT = $GoRoot
$env:GOCACHE = $Cache
$env:GOTMPDIR = $Tmp
Push-Location $Repo
try {
    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP255C_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP255C_GO_CLEAN_FAILED' }

    go test ./unitary ./cmd/unitary-up255c-history-augmented-purity -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP255C_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP255C_FULL_REGRESSION_FAILED' }

    $P1 = ((go run ./cmd/unitary-up255c-history-augmented-purity) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP255C_PROBE1_FAILED' }

    $P2 = ((go run ./cmd/unitary-up255c-history-augmented-purity) | Out-String).Trim()
    if ($P1 -cne $P2) { throw 'UP255C_NONDETERMINISTIC_OUTPUT' }

    $R = $P1 | ConvertFrom-Json
    if ($R.schema -cne 'wingless.up255c-history-augmented-purity.v1' -or
        [int]$R.paired_initial_states -ne 64 -or
        [int]$R.partitions.Count -ne 3) {
        throw 'UP255C_DESIGN'
    }
    if ($R.parent_classification -cne 'FULL_NATIVE_STATE_NONSEPARABLE' -or
        [int]$R.partitions[0].mixed_groups -ne [int]$R.parent_mixed_groups) {
        throw 'UP255C_PARENT_ANCHOR_DRIFT'
    }
    if (-not $R.diagnostic_only -or
        $R.intervention_changed -or
        $R.classifier_training_used -or
        $R.adaptive_feature_selection_used -or
        $R.live_activation) {
        throw 'UP255C_BOUNDARY_LEAK'
    }

    Write-Host $P1
    Write-Host ''
    Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
    Write-Host "classification=$($R.classification)"
    Write-Host "static_mixed=$($R.partitions[0].mixed_groups)"
    Write-Host "current_mixed=$($R.partitions[1].mixed_groups)"
    Write-Host "trajectory_mixed=$($R.partitions[2].mixed_groups)"
    Write-Host 'WINGLESS_UP385_HARNESS_PASS'
}
finally {
    Pop-Location
    if ($null -eq $OldGOROOT) { Remove-Item Env:GOROOT -ErrorAction SilentlyContinue } else { $env:GOROOT = $OldGOROOT }
    if ($null -eq $OldCache) { Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue } else { $env:GOCACHE = $OldCache }
    if ($null -eq $OldTmp) { Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue } else { $env:GOTMPDIR = $OldTmp }
    Remove-Item -Recurse -Force $Cache,$Tmp -ErrorAction SilentlyContinue
}
