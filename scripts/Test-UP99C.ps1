$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up99c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up99c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP99C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP99C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP99C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP99C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up99c-admission-trigger -count=1;if($LASTEXITCODE-ne 0){throw 'UP99C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP99C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up99c-admission-trigger)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP99C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up99c-admission-trigger)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP99C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP99C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up99c-admission-trigger.v1'){throw 'UP99C_SCHEMA_MISMATCH'}
 if([int]$R.filter_bits-ne 1024){throw 'UP99C_FILTER_BITS'}
 if([int]$R.generation_interval-ne 32){throw 'UP99C_INTERVAL'}
 if([int]$R.post_trigger_counter-ne 8){throw 'UP99C_POST_TRIGGER_COUNTER'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UP99C_RECALL_CAP'}
 if($R.explicit_phase_label_used){throw 'UP99C_PHASE_LABEL_PRESENT'}
 if($R.future_oracle_used){throw 'UP99C_ORACLE_PRESENT'}
 if($R.query_labels_used_for_admission){throw 'UP99C_QUERY_LABEL_LEAK'}
 if([int]$R.points.Count-ne 4){throw 'UP99C_POINT_COUNT'}
 foreach($M in $R.points){
  if([int]$M.recall_entries_used-gt 16){throw 'UP99C_RECALL_CAP_EXCEEDED'}
  if($M.arm-ceq 'continuous_control'){
   if([int]$M.policy_metadata_bytes-ne 134){throw 'UP99C_CONTROL_METADATA'}
  }else{
   if([int]$M.policy_metadata_bytes-ne 135){throw 'UP99C_TRIGGER_METADATA'}
  }
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "arm=$($M.arm) trigger=$($M.trigger_threshold) hot=$($M.hot_query_accuracy) exact=$($M.hot_set_exact_accuracy) fp=$($M.false_positive_churn_admissions) fired=$($M.trigger_fired_episodes) trigger_index=$($M.mean_filtered_write_index_at_trigger)"}
 Write-Host 'WINGLESS_UP165_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
