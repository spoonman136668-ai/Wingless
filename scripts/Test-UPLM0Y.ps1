$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm0y-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm0y-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM0Y_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM0Y_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM0Y_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM0Y_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm0y-delta-balance -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0Y_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0Y_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm0y-delta-balance)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0Y_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm0y-delta-balance)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0Y_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM0Y_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm0y-delta-balance.v1'){throw 'UPLM0Y_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UPLM0Y_STATE_DIM'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UPLM0Y_RECALL_CAP'}
 if([int]$R.base_epochs-ne 20 -or [int]$R.adaptation_epochs-ne 4){throw 'UPLM0Y_EPOCHS'}
 if([double]$R.learning_rate-ne 0.08){throw 'UPLM0Y_LR'}
 if($R.router_retraining_used){throw 'UPLM0Y_ROUTER_RETRAINING'}
 if($R.attention_used){throw 'UPLM0Y_ATTENTION_PRESENT'}
 if($R.future_oracle_used){throw 'UPLM0Y_ORACLE_PRESENT'}
 if($R.anchoring_used){throw 'UPLM0Y_ANCHOR_PRESENT'}
 if([int]$R.byte_metrics.Count-ne 6){throw 'UPLM0Y_BYTE_COUNT'}
 if([int]$R.routing_metrics.Count-ne 15){throw 'UPLM0Y_ROUTE_COUNT'}
 if([int]$R.arm_diagnostics.Count-ne 3){throw 'UPLM0Y_DIAG_COUNT'}
 foreach($A in @('paired_delta_sum','norm_balanced_delta','norm_balanced_row_projected')){
  if(@($R.byte_metrics|Where-Object{$_.arm-ceq $A}).Count-ne 2){throw "UPLM0Y_BYTE_$A"}
  if(@($R.routing_metrics|Where-Object{$_.arm-ceq $A}).Count-ne 5){throw "UPLM0Y_ROUTE_$A"}
  if(@($R.arm_diagnostics|Where-Object{$_.arm-ceq $A}).Count-ne 1){throw "UPLM0Y_DIAG_$A"}
 }
 foreach($M in $R.routing_metrics){if([int]$M.metric.max_recall_entries-gt 16){throw 'UPLM0Y_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.arm_diagnostics){Write-Host "arm=$($M.arm) raw_ratio=$($M.mean_raw_para_to_base_norm_ratio) projections=$($M.pair_row_projections) pairs=$($M.pairs_processed)"}
 foreach($M in $R.byte_metrics){Write-Host "arm=$($M.arm) split=$($M.split) accuracy=$($M.top1_accuracy) perplexity=$($M.perplexity)"}
 Write-Host 'WINGLESS_UP166_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
