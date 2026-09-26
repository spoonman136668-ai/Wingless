$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up169b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up169b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP169B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP169B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP169B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP169B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up169b-margin-rank-class-calibration -count=1;if($LASTEXITCODE-ne 0){throw 'UP169B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP169B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up169b-margin-rank-class-calibration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP169B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up169b-margin-rank-class-calibration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP169B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP169B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up169b-margin-rank-class-calibration.v1'){throw 'UP169B_SCHEMA_MISMATCH'}
 if([int]$R.class_metrics.Count-ne 3 -or [int]$R.calibration_phase-ne 16 -or [int]$R.parent_phase-ne 15){throw 'UP169B_DESIGN'}
 if(-not $R.bands_frozen_before_run -or $R.adaptive_band_used -or $R.maintenance_triggered){throw 'UP169B_BOUNDARY'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.class_metrics){Write-Host "class=$($M.class) band=$($M.preregistered_band_max_rank) crossings=$($M.total_crossings) maxrank=$($M.max_crossing_rank) coverage=$($M.class_band_coverage) universal4=$($M.universal4_coverage) density=$($M.class_band_crossing_density)"}
 Write-Host 'WINGLESS_UP169_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
