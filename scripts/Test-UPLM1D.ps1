$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm1d-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm1d-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM1D_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM1D_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM1D_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM1D_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm1d-third-family-byte-adaptation -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1D_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1D_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm1d-third-family-byte-adaptation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1D_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm1d-third-family-byte-adaptation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1D_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM1D_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm1d-third-family-byte-adaptation.v1'){throw 'UPLM1D_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UPLM1D_STATE_DIM'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UPLM1D_RECALL_CAP'}
 if([int]$R.base_pretrain_epochs-ne 20 -or [int]$R.joint_interleaved_epochs-ne 20){throw 'UPLM1D_START_TRAINING'}
 if($R.router_retraining_used){throw 'UPLM1D_ROUTER_RETRAINED'}
 if($R.recurrent_parameters_trained){throw 'UPLM1D_RECURRENT_TRAINED'}
 if($R.attention_used){throw 'UPLM1D_ATTENTION'}
 if($R.future_oracle_used){throw 'UPLM1D_ORACLE'}
 if(-not $R.output_alphabet_extended){throw 'UPLM1D_ALPHABET_NOT_EXTENDED'}
 if([int]$R.metrics.Count-ne 48){throw "UPLM1D_METRIC_COUNT_$($R.metrics.Count)"}
 foreach($E in @(0,1,2,4)){
  $M=@($R.metrics|Where-Object{[int]$_.adaptation_epochs-eq $E})
  if($M.Count-ne 12){throw "UPLM1D_DEPTH_$E"}
 }
 foreach($M in $R.metrics){if([int]$M.metric.max_recall_entries-gt 16){throw 'UPLM1D_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "new_output_bytes=$($R.new_output_bytes -join ',')"
 foreach($M in $R.metrics){Write-Host "adapt=$($M.adaptation_epochs) split=$($M.metric.split) accuracy=$($M.metric.top1_accuracy) perplexity=$($M.metric.perplexity) dependent=$($M.metric.dependent_first_byte_accuracy) exact=$($M.metric.query_set_exact_accuracy) event=$($M.metric.event_routing_accuracy)"}
 Write-Host 'WINGLESS_UP165_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
