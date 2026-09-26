$ErrorActionPreference = 'Stop'

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache = $env:GOCACHE
$PriorGoTmp   = $env:GOTMPDIR
$GoCache      = Join-Path $env:TEMP ("wingless-up11-gocache-" + $PID)
$GoTmp        = Join-Path $env:TEMP ("wingless-up11-gotmp-" + $PID)

New-Item -ItemType Directory -Force -Path $GoCache | Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp | Out-Null
$env:GOCACHE = $GoCache
$env:GOTMPDIR = $GoTmp

Push-Location $Repo
try {
    Write-Host '=== WINGLESS UP-11 PILOT-RANK / MEMORY-NOISE SWEEP ==='
    Write-Host "Repo: $Repo"
    Write-Host "Isolated GOCACHE: $GoCache"
    Write-Host "Isolated GOTMPDIR: $GoTmp"

    go env -w GOFLAGS=-buildvcs=false
    if ($LASTEXITCODE -ne 0) { throw 'UP11_GO_ENV_FAILED' }

    go clean -cache -testcache
    if ($LASTEXITCODE -ne 0) { throw 'UP11_GO_CLEAN_FAILED' }

    go run ./cmd/ice build
    if ($LASTEXITCODE -ne 0) { throw 'UP11_ICE_BUILD_FAILED' }

    go run ./cmd/ice validate
    if ($LASTEXITCODE -ne 0) { throw 'UP11_ICE_VALIDATE_FAILED' }

    go test ./unitary ./cmd/unitary-pilot-rank-probe -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP11_FOCUSED_TEST_FAILED' }

    go test ./... -count=1
    if ($LASTEXITCODE -ne 0) { throw 'UP11_FULL_REGRESSION_FAILED' }

    $Probe1 = ((go run ./cmd/unitary-pilot-rank-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP11_PROBE1_FAILED' }

    $Probe2 = ((go run ./cmd/unitary-pilot-rank-probe) | Out-String).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'UP11_PROBE2_FAILED' }

    if ($Probe1 -cne $Probe2) { throw 'UP11_NONDETERMINISTIC_OUTPUT' }

    $Result = $Probe1 | ConvertFrom-Json
    if ($Result.schema -cne 'wingless.unitary-pilot-rank-probe.v1') { throw 'UP11_SCHEMA_MISMATCH' }
    if ($Result.runtime_prototype_lookup) { throw 'UP11_PROTOTYPE_LOOKUP_ENABLED' }
    if ($Result.explicit_inverse_readout) { throw 'UP11_INVERSE_READOUT_ENABLED' }
    if ($Result.explicit_depth_provided) { throw 'UP11_DEPTH_PROVIDED' }
    if (-not $Result.global_phase_nuisance) { throw 'UP11_GLOBAL_PHASE_NUISANCE_MISSING' }
    if ([int]$Result.min_train_marginal_count -le 0 -or [int]$Result.min_held_out_marginal_count -le 0) {
        throw 'UP11_BALANCED_SPLIT_INVALID'
    }

    Write-Host $Probe1
    Write-Host ""
    Write-Host "=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==="
    Write-Host "Minimal noiseless rank:        $($Result.diagnosis.minimal_noiseless_rank)"
    Write-Host "Minimal noise-0.05 rank:       $($Result.diagnosis.minimal_noise_0_05_rank)"
    Write-Host "Intrinsic rank loss supported: $($Result.diagnosis.intrinsic_rank_loss_supported)"
    Write-Host "Memory-noise sensitivity:      $($Result.diagnosis.memory_noise_sensitivity_supported)"
    Write-Host "Four-pilot candidate:          $($Result.diagnosis.four_pilot_candidate_supported)"
    Write-Host 'WINGLESS_UP11_HARNESS_PASS'
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
