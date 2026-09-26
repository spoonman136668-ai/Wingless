$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm1f-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm1f-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM1F_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM1F_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM1F_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM1F_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm1f-symmetric-triplet-gradient -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1F_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1F_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm1f-symmetric-triplet-gradient)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1F_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm1f-symmetric-triplet-gradient)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1F_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM1F_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm1f-symmetric-triplet-gradient.v1'){throw 'UPLM1F_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UPLM1F_STATE_DIM'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UPLM1F_RECALL_CAP'}
 if([int]$R.adaptation_epochs-ne 4){throw 'UPLM1F_ADAPT_EPOCHS'}
 if($R.router_retraining_used){throw 'UPLM1F_ROUTER_RETRAINED'}
 if($R.recurrent_parameters_trained){throw 'UPLM1F_RECURRENT_TRAINED'}
 if($R.attention_used){throw 'UPLM1F_ATTENTION'}
 if($R.future_oracle_used){throw 'UPLM1F_ORACLE'}
 if([int]$R.byte_metrics.Count-ne 9){throw "UPLM1F_BYTE_COUNT_$($R.byte_metrics.Count)"}
 if([int]$R.routing_metrics.Count-ne 30){throw "UPLM1F_ROUTE_COUNT_$($R.routing_metrics.Count)"}
 foreach($A in @('cyclic_by_index','symmetric_triplet_average','symmetric_triplet_sum')){
  if(@($R.byte_metrics|Where-Object{$_.arm-ceq $A}).Count-ne 3){throw "UPLM1F_BYTE_$A"}
  if(@($R.routing_metrics|Where-Object{$_.arm-ceq $A}).Count-ne 10){throw "UPLM1F_ROUTE_$A"}
 }
 foreach($M in $R.routing_metrics){if([int]$M.metric.max_recall_entries-gt 16){throw 'UPLM1F_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.byte_metrics){Write-Host "arm=$($M.arm) split=$($M.split) accuracy=$($M.top1_accuracy) perplexity=$($M.perplexity)"}
 foreach($M in $R.routing_metrics){Write-Host "arm=$($M.arm) split=$($M.metric.split) exact=$($M.metric.query_set_exact_accuracy) event=$($M.metric.event_routing_accuracy)"}
 Write-Host 'WINGLESS_UP171_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
