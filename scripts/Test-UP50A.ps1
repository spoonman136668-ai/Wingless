$ErrorActionPreference = 'Stop'
$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp = $env:GOTMPDIR
$GoCache = Join-Path $env:TEMP ("wingless-up50a-gocache-" + $PID)
$GoTmp = Join-Path $env:TEMP ("wingless-up50a-gotmp-" + $PID)
New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp
Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-50A THREE-ROLE TRANSFER ==='
    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP50A_GO_ENV_FAILED' }
    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP50A_GO_CLEAN_FAILED' }
    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP50A_ICE_BUILD_FAILED' }
    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP50A_ICE_VALIDATE_FAILED' }
    go test ./unitary ./cmd/unitary-up50a-three-role -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP50A_FOCUSED_TEST_FAILED' }
    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP50A_FULL_REGRESSION_FAILED' }
    $Probe1 = ((go run ./cmd/unitary-up50a-three-role) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP50A_PROBE1_FAILED' }
    $Probe2 = ((go run ./cmd/unitary-up50a-three-role) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP50A_PROBE2_FAILED' }
    if ($Probe1 -cne $Probe2) { throw 'UP50A_NONDETERMINISTIC_OUTPUT' }
    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.up50a-three-role-transfer.v1') { throw 'UP50A_SCHEMA_MISMATCH' }
    if (-not $Result.train_single_step_only) { throw 'UP50A_SEQUENCE_TRAIN_LEAK' }
    if ($Result.swap12_trained -or $Result.role1_primitive_trained -or $Result.role2_primitive_trained) { throw 'UP50A_HELDOUT_PRIMITIVE_LEAK' }
    if ([int]$Result.joint_dimension -ne 27) { throw 'UP50A_DIMENSION_MISMATCH' }

    Write-Host $Probe1
    Write-Host ''
    Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
    Write-Host "Unitary train:          $($Result.diagnosis.unitary_train_accuracy)"
    Write-Host "Unitary held-out:       $($Result.diagnosis.unitary_heldout_accuracy)"
    Write-Host "Unseen swap12:          $($Result.diagnosis.unseen_swap12_accuracy)"
    Write-Host "Derived role1:          $($Result.diagnosis.derived_role1_accuracy)"
    Write-Host "Derived role2:          $($Result.diagnosis.derived_role2_accuracy)"
    Write-Host "Long programs:          $($Result.diagnosis.long_program_accuracy)"
    Write-Host "Three-role gate:        $($Result.diagnosis.three_role_transfer_gate)"
    Write-Host "Matched nonunitary:     $($Result.diagnosis.nonunitary_heldout_accuracy)"
    Write-Host 'WINGLESS_UP50_HARNESS_PASS'
}
finally {
    Pop-Location
    if ($null -eq $PriorGoCache) { Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue } else { $env:GOCACHE = $PriorGoCache }
    if ($null -eq $PriorGoTmp) { Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue } else { $env:GOTMPDIR = $PriorGoTmp }
    Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
    Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
