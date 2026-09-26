$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm2t-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm2t-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM2T_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM2T_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM2T_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM2T_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm2t-sham-cost-control -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2T_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2T_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm2t-sham-cost-control)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2T_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm2t-sham-cost-control)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2T_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM2T_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm2t-sham-cost-control.v1'){throw 'UPLM2T_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 24){throw 'UPLM2T_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.capacity_changed -or $R.extra_training_used){throw 'UPLM2T_BOUNDARY_LEAK'}
 if([int]$R.guided_interventions-ne [int]$R.sham_interventions){throw 'UPLM2T_BUDGET_MISMATCH'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "baseline=$($R.baseline_harm) guided=$($R.guided_harm) sham=$($R.sham_harm) guided_prevented=$($R.guided_prevented) sham_prevented=$($R.sham_prevented)"
 Write-Host 'WINGLESS_UP181_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
