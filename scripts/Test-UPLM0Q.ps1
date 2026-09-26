$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm0q-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm0q-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM0Q_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM0Q_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM0Q_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM0Q_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm0q-update-order -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0Q_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0Q_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm0q-update-order)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0Q_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm0q-update-order)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0Q_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM0Q_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm0q-update-order.v1'){throw 'UPLM0Q_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UPLM0Q_STATE_DIM'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UPLM0Q_RECALL_CAP'}
 if([int]$R.base_epochs-ne 20){throw 'UPLM0Q_BASE_EPOCHS'}
 if([int]$R.adaptation_epochs-ne 4){throw 'UPLM0Q_ADAPT_EPOCHS'}
 if([double]$R.learning_rate-ne 0.08){throw 'UPLM0Q_LR'}
 if([int]$R.paraphrase_passes_per_epoch-ne 1 -or [int]$R.base_passes_per_epoch-ne 1){throw 'UPLM0Q_PASS_COUNTS'}
 if($R.router_retraining_used){throw 'UPLM0Q_ROUTER_RETRAINING'}
 if($R.attention_used){throw 'UPLM0Q_ATTENTION_PRESENT'}
 if($R.future_oracle_used){throw 'UPLM0Q_ORACLE_PRESENT'}
 if([int]$R.byte_metrics.Count-ne 6){throw 'UPLM0Q_BYTE_COUNT'}
 if([int]$R.routing_metrics.Count-ne 15){throw 'UPLM0Q_ROUTE_COUNT'}
 foreach($A in @('para_then_base','base_then_para','example_interleaved')){
  if(@($R.byte_metrics|Where-Object{$_.schedule-ceq $A}).Count-ne 2){throw "UPLM0Q_BYTE_$A"}
  if(@($R.routing_metrics|Where-Object{$_.schedule-ceq $A}).Count-ne 5){throw "UPLM0Q_ROUTE_$A"}
 }
 foreach($M in $R.routing_metrics){if([int]$M.metric.max_recall_entries-gt 16){throw 'UPLM0Q_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.byte_metrics){Write-Host "schedule=$($M.schedule) split=$($M.split) accuracy=$($M.top1_accuracy) perplexity=$($M.perplexity)"}
 foreach($M in $R.routing_metrics){Write-Host "schedule=$($M.schedule) split=$($M.metric.split) accuracy=$($M.metric.top1_accuracy) exact=$($M.metric.query_set_exact_accuracy)"}
 Write-Host 'WINGLESS_UP142_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
