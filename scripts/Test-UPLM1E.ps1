$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm1e-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm1e-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM1E_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM1E_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM1E_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM1E_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm1e-triplet-order -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1E_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1E_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm1e-triplet-order)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1E_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm1e-triplet-order)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1E_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM1E_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm1e-triplet-order.v1'){throw 'UPLM1E_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UPLM1E_STATE_DIM'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UPLM1E_RECALL_CAP'}
 if([int]$R.adaptation_epochs-ne 4){throw 'UPLM1E_ADAPT_EPOCHS'}
 if($R.router_retraining_used){throw 'UPLM1E_ROUTER_RETRAINED'}
 if($R.recurrent_parameters_trained){throw 'UPLM1E_RECURRENT_TRAINED'}
 if($R.attention_used){throw 'UPLM1E_ATTENTION'}
 if($R.future_oracle_used){throw 'UPLM1E_ORACLE'}
 if([int]$R.byte_metrics.Count-ne 12){throw "UPLM1E_BYTE_COUNT_$($R.byte_metrics.Count)"}
 if([int]$R.routing_metrics.Count-ne 40){throw "UPLM1E_ROUTE_COUNT_$($R.routing_metrics.Count)"}
 foreach($A in @('third_base_para','base_para_third','cyclic_by_index','cyclic_by_epoch')){
  if(@($R.byte_metrics|Where-Object{$_.schedule-ceq $A}).Count-ne 3){throw "UPLM1E_BYTE_$A"}
  if(@($R.routing_metrics|Where-Object{$_.schedule-ceq $A}).Count-ne 10){throw "UPLM1E_ROUTE_$A"}
 }
 foreach($M in $R.routing_metrics){if([int]$M.metric.max_recall_entries-gt 16){throw 'UPLM1E_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.byte_metrics){Write-Host "schedule=$($M.schedule) split=$($M.split) accuracy=$($M.top1_accuracy) perplexity=$($M.perplexity)"}
 foreach($M in $R.routing_metrics){Write-Host "schedule=$($M.schedule) split=$($M.metric.split) exact=$($M.metric.query_set_exact_accuracy) event=$($M.metric.event_routing_accuracy)"}
 Write-Host 'WINGLESS_UP168_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
