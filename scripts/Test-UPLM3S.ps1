$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm3s-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm3s-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM3S_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM3S_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM3S_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM3S_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm3s-pressure-saturation-diagnostic -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3S_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3S_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm3s-pressure-saturation-diagnostic)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3S_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm3s-pressure-saturation-diagnostic)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3S_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM3S_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm3s-pressure-saturation-diagnostic.v1'){throw 'UPLM3S_SCHEMA_MISMATCH'}
 if([int]$R.pressure_levels.Count-ne 17 -or [int]$R.metrics.Count-ne 17 -or [int]$R.resource_configurations_per_pressure-ne 18){throw 'UPLM3S_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_pressure_selection_used -or $R.resource_tuning_used){throw 'UPLM3S_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "pressure=$($M.pressure_level) min=$($M.min_earliest_failed) max=$($M.max_earliest_failed) range=$($M.outcome_range) distinct=$($M.distinct_outcomes) at_min=$($M.configurations_at_min) at_max=$($M.configurations_at_max)"}
 Write-Host 'WINGLESS_UP224_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
