$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm1o-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm1o-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM1O_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM1O_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM1O_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM1O_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm1o-balanced-reintegration -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1O_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1O_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm1o-balanced-reintegration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1O_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm1o-balanced-reintegration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1O_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM1O_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm1o-balanced-reintegration.v1'){throw 'UPLM1O_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64 -or [int]$R.exact_recall_cap-ne 16){throw 'UPLM1O_CAPACITY'}
 if([int]$R.adaptation_epochs-ne 4 -or [double]$R.full_step_learning_rate-ne 0.08 -or [double]$R.half_step_learning_rate-ne 0.04){throw 'UPLM1O_SCHEDULE'}
 if($R.router_retraining_used -or $R.recurrent_parameters_trained -or $R.recall_cap_changed -or $R.attention_used -or $R.future_oracle_used){throw 'UPLM1O_BOUNDARY'}
 if([int]$R.byte_metrics.Count-ne 12){throw "UPLM1O_BYTE_COUNT $($R.byte_metrics.Count)"}
 if([int]$R.routing_metrics.Count-ne 48){throw "UPLM1O_ROUTE_COUNT $($R.routing_metrics.Count)"}
 foreach($M in $R.routing_metrics){if([int]$M.metric.max_recall_entries-gt 16){throw 'UPLM1O_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.byte_metrics){Write-Host "arm=$($M.arm) split=$($M.split) accuracy=$($M.top1_accuracy) perplexity=$($M.perplexity)"}
 foreach($M in $R.routing_metrics){if($M.metric.split-ceq 'fourth_block_stream1'){Write-Host "arm=$($M.arm) fourth=$($M.metric.top1_accuracy) exact=$($M.metric.query_set_exact_accuracy) event=$($M.metric.event_routing_accuracy) report=$($M.metric.report_routing_accuracy)"}}
 Write-Host 'WINGLESS_UP178_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
