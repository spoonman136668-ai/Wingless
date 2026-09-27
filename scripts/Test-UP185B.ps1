$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up185b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up185b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP185B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP185B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP185B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP185B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up185b-parity-diagnostic -count=1;if($LASTEXITCODE-ne 0){throw 'UP185B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP185B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up185b-parity-diagnostic)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP185B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up185b-parity-diagnostic)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP185B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP185B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up185b-parity-diagnostic.v1'){throw 'UP185B_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 6 -or [int]$R.summaries.Count-ne 2){throw 'UP185B_DESIGN'}
 if(-not $R.parity_is_diagnostic_only -or $R.native_deployment_feature -or $R.evaluation_fitting_used -or $R.maintenance_triggered){throw 'UP185B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "model=$($S.model) mean_abs_mass_error=$($S.mean_absolute_mass_ratio_error)"}
 Write-Host 'WINGLESS_UP185_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
