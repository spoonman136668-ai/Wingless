$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up48a-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up48a-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-48A SIZE-TRANSFER COMPOSITION ==='
    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP48A_GO_ENV_FAILED' }
    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP48A_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP48A_ICE_BUILD_FAILED' }
    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP48A_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-up48a-size-transfer -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP48A_FOCUSED_TEST_FAILED' }
    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP48A_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-up48a-size-transfer) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP48A_PROBE1_FAILED' }
    $Probe2 = ((go run ./cmd/unitary-up48a-size-transfer) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP48A_PROBE2_FAILED' }
    if ($Probe1 -cne $Probe2) { throw 'UP48A_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.up48a-size-transfer-composition.v1') { throw 'UP48A_SCHEMA_MISMATCH' }
    if ($Result.train_programs_long) { throw 'UP48A_TRAINING_LEAKED_LONG_PROGRAMS' }
    if ([int]$Result.train_dimensions.Count -ne 2) { throw 'UP48A_TRAIN_DIM_COUNT' }
    if ([int]$Result.heldout_dimensions.Count -ne 2) { throw 'UP48A_HELD_DIM_COUNT' }
    if ([int]$Result.train_dimensions[0] -ne 4 -or [int]$Result.train_dimensions[1] -ne 6) { throw 'UP48A_TRAIN_DIM_MISMATCH' }
    if ([int]$Result.heldout_dimensions[0] -ne 8 -or [int]$Result.heldout_dimensions[1] -ne 12) { throw 'UP48A_HELD_DIM_MISMATCH' }
    if ([int]$Result.heldout_lengths.Count -ne 4) { throw 'UP48A_HELD_LENGTH_COUNT' }

    Write-Host $Probe1
    Write-Host ''
    Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
    Write-Host "Unitary train accuracy:       $($Result.diagnosis.unitary_train_accuracy)"
    Write-Host "Unitary held-out accuracy:    $($Result.diagnosis.unitary_heldout_accuracy)"
    Write-Host "Minimum dimension accuracy:   $($Result.diagnosis.minimum_unitary_dimension_accuracy)"
    Write-Host "Algorithmic transfer gate:    $($Result.diagnosis.unitary_algorithmic_transfer_gate)"
    Write-Host "Matched nonunitary held-out:  $($Result.diagnosis.nonunitary_heldout_accuracy)"
    Write-Host "Unitary held-out advantage:   $($Result.diagnosis.unitary_heldout_advantage)"
    Write-Host 'WINGLESS_UP48_HARNESS_PASS'
}
finally {
    Pop-Location
    if ($null -eq $PriorGoCache) { Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue } else { $env:GOCACHE = $PriorGoCache }
    if ($null -eq $PriorGoTmp) { Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue } else { $env:GOTMPDIR = $PriorGoTmp }
    Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
    Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
