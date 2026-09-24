$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up53c-gocache-"+$PID)
$GoTmp=Join-Path $env:TEMP ("wingless-up53c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 Write-Host '=== WINGLESS UP-53C PHASE MULTIPLEX MEMORY ==='
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP53C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP53C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP53C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP53C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up53c-phase-multiplex -count=1;if($LASTEXITCODE-ne 0){throw 'UP53C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP53C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up53c-phase-multiplex)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP53C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up53c-phase-multiplex)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP53C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP53C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up53c-phase-multiplex-memory.v1'){throw 'UP53C_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 16){throw 'UP53C_STATE_DIMENSION'}
 if([int]$R.metrics.Count-ne 6){throw 'UP53C_LEVEL_COUNT'}
 if(-not $R.coherence_global_phase_invariant){throw 'UP53C_COHERENCE_FLAG'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "Maximum passing banks: $($R.maximum_passing_banks)"
 Write-Host "First failing banks:   $($R.first_failing_banks)"
 foreach($M in $R.metrics){
  Write-Host "banks=$($M.banks) gate=$($M.gate) coherence_value=$($M.coherence_value_accuracy) coherence_exact=$($M.coherence_exact_scenario_accuracy) magnitude_value=$($M.magnitude_value_accuracy) margin=$($M.coherence_minimum_margin) drift=$($M.max_norm_drift)"
 }
 Write-Host 'WINGLESS_UP53_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}