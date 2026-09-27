$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up192b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up192b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP192B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP192B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP192B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP192B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up192b-native-geometry-residual -count=1;if($LASTEXITCODE-ne 0){throw 'UP192B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP192B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up192b-native-geometry-residual)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP192B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up192b-native-geometry-residual)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP192B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP192B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up192b-native-geometry-residual.v1'){throw 'UP192B_SCHEMA_MISMATCH'}
 if([int]$R.profiles.Count-ne 12 -or [int]$R.metrics.Count-ne 36 -or [int]$R.correlations.Count-ne 6){throw 'UP192B_DESIGN'}
 if($R.new_correction_fit -or $R.correction_applied -or $R.adaptive_feature_search_used -or $R.maintenance_triggered){throw 'UP192B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($C in $R.correlations){Write-Host "feature=$($C.feature) pearson_required_factor=$($C.pearson_required_factor)"}
 Write-Host 'WINGLESS_UP192B_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
