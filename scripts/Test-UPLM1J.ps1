$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm1j-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm1j-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM1J_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM1J_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM1J_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM1J_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm1j-fourth-family-boundary -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1J_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1J_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm1j-fourth-family-boundary)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1J_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm1j-fourth-family-boundary)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1J_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM1J_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm1j-fourth-family-boundary.v1'){throw 'UPLM1J_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64 -or [int]$R.exact_recall_cap-ne 16){throw 'UPLM1J_STATE_BOUND'}
 if([int]$R.base_pretrain_epochs-ne 20 -or [int]$R.joint_interleaved_epochs-ne 20 -or [int]$R.three_family_adapt_epochs-ne 4 -or [int]$R.fourth_router_epochs-ne 20){throw 'UPLM1J_EPOCHS'}
 if($R.fourth_family_byte_training -or $R.recurrent_parameters_trained -or $R.attention_used -or $R.future_oracle_used){throw 'UPLM1J_BOUNDARY'}
 if([int]$R.router_metrics.Count-ne 1){throw 'UPLM1J_ROUTER_COUNT'}
 if([int]$R.metrics.Count-ne 16){throw "UPLM1J_METRIC_COUNT $($R.metrics.Count)"}
 foreach($M in $R.metrics){if([int]$M.max_recall_entries-gt 16){throw 'UPLM1J_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.router_metrics){Write-Host "router split=$($M.split) accuracy=$($M.accuracy) examples=$($M.examples)"}
 foreach($M in $R.metrics){Write-Host "split=$($M.split) accuracy=$($M.top1_accuracy) perplexity=$($M.perplexity) dependent=$($M.dependent_first_byte_accuracy) exact=$($M.query_set_exact_accuracy) event=$($M.event_routing_accuracy) report=$($M.report_routing_accuracy) entries=$($M.max_recall_entries)"}
 Write-Host 'WINGLESS_UP171_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
