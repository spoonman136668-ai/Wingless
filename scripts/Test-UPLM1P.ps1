$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm1p-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm1p-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM1P_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM1P_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM1P_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM1P_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm1p-fifth-family-scale -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1P_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1P_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm1p-fifth-family-scale)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1P_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm1p-fifth-family-scale)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1P_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM1P_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm1p-fifth-family-scale.v1'){throw 'UPLM1P_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64 -or [int]$R.exact_recall_cap-ne 16){throw 'UPLM1P_CAPACITY'}
 if([int]$R.fifth_router_epochs-ne 20 -or [int]$R.adaptation_epochs-ne 4){throw 'UPLM1P_EPOCHS'}
 if([double]$R.full_step_learning_rate-ne 0.08 -or [double]$R.half_step_learning_rate-ne 0.04){throw 'UPLM1P_LR'}
 if($R.recurrent_parameters_trained -or $R.recall_cap_changed -or $R.attention_used -or $R.future_oracle_used -or $R.fifth_correction_used){throw 'UPLM1P_BOUNDARY'}
 if([int]$R.byte_metrics.Count-ne 15){throw "UPLM1P_BYTE_COUNT $($R.byte_metrics.Count)"}
 if([int]$R.summaries.Count-ne 3){throw "UPLM1P_SUMMARY_COUNT $($R.summaries.Count)"}
 if([int]$R.router_metrics.Count-ne 7){throw "UPLM1P_ROUTER_COUNT $($R.router_metrics.Count)"}
 if([int]$R.integrated_metrics.Count-ne 30){throw "UPLM1P_INTEGRATED_COUNT $($R.integrated_metrics.Count)"}
 foreach($M in $R.integrated_metrics){if([int]$M.metric.max_recall_entries-gt 16){throw 'UPLM1P_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.summaries){Write-Host "arm=$($M.arm) min=$($M.min_heldout_accuracy) mean=$($M.mean_heldout_accuracy) spread=$($M.heldout_accuracy_spread)"}
 foreach($M in $R.router_metrics){Write-Host "router family=$($M.family) split=$($M.split) acc=$($M.accuracy) store=$($M.store_recall) observe=$($M.observe_recall) report=$($M.report_recall)"}
 Write-Host 'WINGLESS_UP178_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
