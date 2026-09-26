$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm2w-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm2w-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM2W_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM2W_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM2W_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM2W_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm2w-budget-frontier -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2W_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2W_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm2w-budget-frontier)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2W_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm2w-budget-frontier)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2W_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM2W_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm2w-budget-frontier.v1'){throw 'UPLM2W_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 4 -or [int]$R.points.Count-ne 96 -or [int]$R.exact_recall_cap-ne 16){throw 'UPLM2W_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.future_schedule_used -or $R.capacity_changed -or $R.extra_training_used){throw 'UPLM2W_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "budget=$($M.budget) baseline=$($M.baseline_failed) guided=$($M.guided_failed) sham=$($M.sham_failed) guided_prevented=$($M.guided_prevented) guided_per_action=$($M.guided_prevented_per_action)"}
 Write-Host 'WINGLESS_UP184_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
