$ErrorActionPreference = 'Stop'
$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp = $env:GOTMPDIR
$GoCache = Join-Path $env:TEMP ("wingless-up49a-gocache-" + $PID)
$GoTmp = Join-Path $env:TEMP ("wingless-up49a-gotmp-" + $PID)
New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp
Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-49A ROLE-BINDING COMPOSITION ==='
    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP49A_GO_ENV_FAILED' }
    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP49A_GO_CLEAN_FAILED' }
    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP49A_ICE_BUILD_FAILED' }
    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP49A_ICE_VALIDATE_FAILED' }
    go test ./unitary ./cmd/unitary-up49a-role-binding -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP49A_FOCUSED_TEST_FAILED' }
    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP49A_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-up49a-role-binding) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP49A_PROBE1_FAILED' }
    $Probe2 = ((go run ./cmd/unitary-up49a-role-binding) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP49A_PROBE2_FAILED' }
    if ($Probe1 -cne $Probe2) { throw 'UP49A_NONDETERMINISTIC_OUTPUT' }
    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.up49a-role-binding-composition.v1') { throw 'UP49A_SCHEMA_MISMATCH' }
    if (-not $Result.train_single_step_only) { throw 'UP49A_TRAIN_SEQUENCE_LEAK' }
    if ($Result.right_primitive_trained) { throw 'UP49A_RIGHT_PRIMITIVE_LEAK' }
    if ([int]$Result.values_per_role -ne 4 -or [int]$Result.joint_dimension -ne 16) { throw 'UP49A_DIMENSION_MISMATCH' }

    Write-Host $Probe1
    Write-Host ''
    Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
    Write-Host "Unitary train accuracy:      $($Result.diagnosis.unitary_train_accuracy)"
    Write-Host "Unitary held-out accuracy:   $($Result.diagnosis.unitary_heldout_accuracy)"
    Write-Host "Conjugated-right accuracy:   $($Result.diagnosis.conjugated_right_accuracy)"
    Write-Host "Long-program accuracy:       $($Result.diagnosis.long_program_accuracy)"
    Write-Host "Role-binding gate:           $($Result.diagnosis.role_binding_gate)"
    Write-Host "Matched nonunitary held-out: $($Result.diagnosis.nonunitary_heldout_accuracy)"
    Write-Host 'WINGLESS_UP49_HARNESS_PASS'
}
finally {
    Pop-Location
    if ($null -eq $PriorGoCache) { Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue } else { $env:GOCACHE = $PriorGoCache }
    if ($null -eq $PriorGoTmp) { Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue } else { $env:GOTMPDIR = $PriorGoTmp }
    Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
    Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
