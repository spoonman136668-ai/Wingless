$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm3q-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm3q-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM3Q_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM3Q_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM3Q_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM3Q_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm3q-rich-profile-pressure-sweep -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3Q_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3Q_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm3q-rich-profile-pressure-sweep)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3Q_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm3q-rich-profile-pressure-sweep)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3Q_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM3Q_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm3q-rich-profile-pressure-sweep.v1'){throw 'UPLM3Q_SCHEMA_MISMATCH'}
 if([int]$R.pressure_levels.Count-ne 5 -or [int]$R.metrics.Count-ne 90 -or [int]$R.summaries.Count-ne 10){throw 'UPLM3Q_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_pressure_levels_used -or $R.adaptive_coordinate_search_used){throw 'UPLM3Q_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "pressure=$($S.pressure_level) coordinate=$($S.coordinate) groups=$($S.groups) nonzero=$($S.nonzero_spread_groups) max_spread=$($S.max_failure_spread) mean_spread=$($S.mean_failure_spread)"}
 Write-Host 'WINGLESS_UP219_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
