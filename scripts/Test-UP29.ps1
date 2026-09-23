$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up29-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up29-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-29 EXACT REAL ORTHOGONAL EQUIVALENCE ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP29_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP29_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP29_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP29_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-real-orthogonal-equivalence-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP29_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP29_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-real-orthogonal-equivalence-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP29_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-real-orthogonal-equivalence-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP29_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP29_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json

    if ($Result.schema -cne 'wingless.real-orthogonal-equivalence.v1') { throw 'UP29_SCHEMA_MISMATCH' }
    if ([int]$Result.complex_latent_dimension -ne 96) { throw 'UP29_COMPLEX_DIMENSION_MISMATCH' }
    if ([int]$Result.real_latent_dimension -ne 192) { throw 'UP29_REAL_DIMENSION_MISMATCH' }
    if (-not $Result.equal_real_scalar_degrees_of_freedom) { throw 'UP29_DOF_MISMATCH' }
    if (-not $Result.complex_transport_unitary) { throw 'UP29_COMPLEX_UNITARY_FLAG_MISSING' }
    if (-not $Result.real_transport_orthogonal) { throw 'UP29_REAL_ORTHOGONAL_FLAG_MISSING' }
    if (-not $Result.exact_realification_comparator) { throw 'UP29_REALIFICATION_FLAG_MISSING' }
    if (-not $Result.runtime_real_arithmetic_only) { throw 'UP29_REAL_RUNTIME_FLAG_MISSING' }
    if ($Result.independent_real_encoder_learned) { throw 'UP29_INDEPENDENT_ENCODER_FALSE_BOUNDARY_BROKEN' }
    if (-not $Result.encoder_realified_from_complex_construction) { throw 'UP29_ENCODER_BOUNDARY_MISSING' }
    if (-not $Result.transport_realified_from_complex_construction) { throw 'UP29_TRANSPORT_BOUNDARY_MISSING' }
    if (-not $Result.observable_bank_frozen_from_up28) { throw 'UP29_FROZEN_BANK_MISSING' }
    if (-not $Result.decoder_frozen_across_representations) { throw 'UP29_FROZEN_DECODER_MISSING' }
    if ([int]$Result.runtime_observable_count -ne 64) { throw 'UP29_OBSERVABLE_COUNT_MISMATCH' }
    if ($Result.runtime_prototype_lookup) { throw 'UP29_PROTOTYPE_LOOKUP_ENABLED' }
    if ($Result.explicit_depth_provided) { throw 'UP29_DEPTH_PROVIDED' }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Baseline reproduction pass:$($Result.diagnosis.baseline_reproduction_pass)"
    Write-Host "Orthogonality pass:        $($Result.diagnosis.orthogonality_pass)"
    Write-Host "State equivalence pass:    $($Result.diagnosis.state_equivalence_pass)"
    Write-Host "Feature equivalence pass:  $($Result.diagnosis.feature_equivalence_pass)"
    Write-Host "Decision equivalence pass: $($Result.diagnosis.decision_equivalence_pass)"
    Write-Host "Mutable integration pass:  $($Result.diagnosis.mutable_integration_pass)"
    Write-Host "Max orthogonality error:   $($Result.diagnosis.max_orthogonality_error)"
    Write-Host "Max state error:           $($Result.diagnosis.max_state_equivalence_error)"
    Write-Host "Max feature error:         $($Result.diagnosis.max_feature_equivalence_error)"
    Write-Host "Complex held accuracy:     $($Result.diagnosis.complex_heldout_accuracy)"
    Write-Host "Real held accuracy:        $($Result.diagnosis.real_heldout_accuracy)"
    Write-Host "Decision disagreements:    $($Result.diagnosis.decision_disagreements)"
    Write-Host 'WINGLESS_UP29_HARNESS_PASS'
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
