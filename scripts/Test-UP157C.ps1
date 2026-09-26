$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up157c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up157c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP157C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP157C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP157C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP157C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up157c-loss-redistribution -count=1;if($LASTEXITCODE-ne 0){throw 'UP157C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP157C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up157c-loss-redistribution)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP157C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up157c-loss-redistribution)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP157C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP157C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up157c-loss-redistribution.v1'){throw 'UP157C_SCHEMA_MISMATCH'}
 if([int]$R.arms-ne 16 -or [int]$R.points.Count-ne 16){throw 'UP157C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation){throw 'UP157C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "permanent=$($R.permanent_saves) delayed=$($R.delayed_losses) protected_gain=$($R.guided_final_protected_gain) original_gain=$($R.guided_final_original_gain) noncohort_delta=$($R.guided_final_noncohort_original_delta)"
 Write-Host 'WINGLESS_UP157_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
