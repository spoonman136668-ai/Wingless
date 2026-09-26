$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm2h-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm2h-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM2H_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM2H_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM2H_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM2H_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm2h-combined-locality-stream-depth -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2H_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2H_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm2h-combined-locality-stream-depth)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2H_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm2h-combined-locality-stream-depth)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2H_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM2H_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm2h-combined-locality-stream-depth.v1'){throw 'UPLM2H_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 500 -or [int]$R.depth_summaries.Count-ne 100){throw 'UPLM2H_COUNT'}
 if([int]$R.allocations-ne 5 -or [int]$R.training_schedules-ne 2 -or [int]$R.variants-ne 2 -or [int]$R.families-ne 5){throw 'UPLM2H_DESIGN'}
 if([int]$R.exact_recall_cap-ne 16 -or $R.recall_cap_changed -or -not $R.store_observe_global_all_variants -or $R.evaluation_retraining_used){throw 'UPLM2H_CAPACITY'}
 if($R.recurrent_parameters_trained -or $R.router_modified -or $R.attention_used -or $R.future_oracle_used){throw 'UPLM2H_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.depth_summaries){Write-Host "allocation=$($S.allocation) schedule=$($S.training_schedule) variant=$($S.variant) stream=$($S.stream) fifth=$($S.fifth_integrated_mean) prior=$($S.prior_integrated_mean) min=$($S.min_family_integrated_mean) exact=$($S.structural_exact_all_families) recall=$($S.max_recall_entries)"}
 Write-Host 'WINGLESS_UP167_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
