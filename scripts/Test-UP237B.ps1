$ErrorActionPreference = 'Stop'
$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$OldGOROOT = $env:GOROOT
$OldCache = $env:GOCACHE
$OldTmp = $env:GOTMPDIR
$GoExe = (Get-Command go -ErrorAction Stop).Source
$GoRoot = Split-Path (Split-Path $GoExe -Parent) -Parent
$Cache = Join-Path $env:TEMP ("up237b-cache-" + $PID)
$Tmp = Join-Path $env:TEMP ("up237b-tmp-" + $PID)
New-Item -ItemType Directory -Force -Path $Cache,$Tmp | Out-Null
$env:GOROOT = $GoRoot
$env:GOCACHE = $Cache
$env:GOTMPDIR = $Tmp
Push-Location $Repo
try {
    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP237B_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP237B_GO_CLEAN_FAILED' }

    go test ./unitary ./cmd/unitary-up237b-far-centered-nearzero-parity -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP237B_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP237B_FULL_REGRESSION_FAILED' }

    $P1 = ((go run ./cmd/unitary-up237b-far-centered-nearzero-parity) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP237B_PROBE1_FAILED' }

    $P2 = ((go run ./cmd/unitary-up237b-far-centered-nearzero-parity) | Out-String).Trim()
    if ($P1 -cne $P2) { throw 'UP237B_NONDETERMINISTIC_OUTPUT' }

    $R = $P1 | ConvertFrom-Json
    if ($R.schema -cne 'wingless.up237b-far-centered-nearzero-parity.v1' -or
        [int]$R.training_phases.Count -ne 15 -or
        [int]$R.far_evaluation_phases.Count -ne 32 -or
        [int]$R.summaries.Count -ne 4) {
        throw 'UP237B_DESIGN'
    }
    if ($R.parent_classification -cne 'CENTERING_REQUIRED' -or
        [double]$R.parent_centered_accuracy -ne 1.0) {
        throw 'UP237B_PARENT_ANCHOR_DRIFT'
    }
    if ($R.evaluation_label_fitting_used -or
        $R.phase_input_at_inference_used -or
        $R.adaptive_feature_selection_used -or
        $R.nonlinear_classifier_used -or
        $R.live_activation) {
        throw 'UP237B_BOUNDARY_LEAK'
    }

    Write-Host $P1
    Write-Host ''
    Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
    Write-Host "far_accuracy=$($R.far_accuracy) classification=$($R.classification)"
    foreach ($S in $R.summaries) {
        Write-Host "family=$($S.family) correct=$($S.correct)/$($S.states)"
    }
    Write-Host 'WINGLESS_UP384_HARNESS_PASS'
}
finally {
    Pop-Location
    if ($null -eq $OldGOROOT) { Remove-Item Env:GOROOT -ErrorAction SilentlyContinue } else { $env:GOROOT = $OldGOROOT }
    if ($null -eq $OldCache) { Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue } else { $env:GOCACHE = $OldCache }
    if ($null -eq $OldTmp) { Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue } else { $env:GOTMPDIR = $OldTmp }
    Remove-Item -Recurse -Force $Cache,$Tmp -ErrorAction SilentlyContinue
}
