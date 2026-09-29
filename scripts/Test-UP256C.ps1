$ErrorActionPreference = 'Stop'
$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$OldGOROOT = $env:GOROOT
$OldCache = $env:GOCACHE
$OldTmp = $env:GOTMPDIR
$GoExe = (Get-Command go -ErrorAction Stop).Source
$GoRoot = Split-Path (Split-Path $GoExe -Parent) -Parent
$Cache = Join-Path $env:TEMP ("up256c-cache-" + $PID)
$Tmp = Join-Path $env:TEMP ("up256c-tmp-" + $PID)
New-Item -ItemType Directory -Force -Path $Cache,$Tmp | Out-Null
$env:GOROOT = $GoRoot
$env:GOCACHE = $Cache
$env:GOTMPDIR = $Tmp
Push-Location $Repo
try {
    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP256C_GO_ENV_FAILED' }
    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP256C_GO_CLEAN_FAILED' }

    go test ./unitary ./cmd/unitary-up256c-current-state-component-sufficiency -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP256C_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP256C_FULL_REGRESSION_FAILED' }

    $P1 = ((go run ./cmd/unitary-up256c-current-state-component-sufficiency) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP256C_PROBE1_FAILED' }
    $P2 = ((go run ./cmd/unitary-up256c-current-state-component-sufficiency) | Out-String).Trim()
    if ($P1 -cne $P2) { throw 'UP256C_NONDETERMINISTIC_OUTPUT' }

    $R = $P1 | ConvertFrom-Json
    if ($R.schema -cne 'wingless.up256c-current-state-component-sufficiency.v1' -or
        [int]$R.paired_initial_states -ne 64 -or
        [int]$R.partitions.Count -ne 4) {
        throw 'UP256C_DESIGN'
    }
    if ($R.parent_classification -cne 'CURRENT_STATE_RESOLVES_ALIASING' -or
        [int]$R.parent_current_mixed_groups -ne 0) {
        throw 'UP256C_PARENT_ANCHOR_DRIFT'
    }
    if (-not $R.diagnostic_only -or
        $R.intervention_changed -or
        $R.classifier_training_used -or
        $R.adaptive_feature_selection_used -or
        $R.live_activation) {
        throw 'UP256C_BOUNDARY_LEAK'
    }

    Write-Host $P1
    Write-Host ''
    Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
    Write-Host "classification=$($R.classification)"
    foreach ($P in $R.partitions) {
        Write-Host "partition=$($P.name) mixed=$($P.mixed_groups) groups=$($P.groups)"
    }
    Write-Host 'WINGLESS_UP387_HARNESS_PASS'
}
finally {
    Pop-Location
    if ($null -eq $OldGOROOT) { Remove-Item Env:GOROOT -ErrorAction SilentlyContinue } else { $env:GOROOT = $OldGOROOT }
    if ($null -eq $OldCache) { Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue } else { $env:GOCACHE = $OldCache }
    if ($null -eq $OldTmp) { Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue } else { $env:GOTMPDIR = $OldTmp }
    Remove-Item -Recurse -Force $Cache,$Tmp -ErrorAction SilentlyContinue
}
