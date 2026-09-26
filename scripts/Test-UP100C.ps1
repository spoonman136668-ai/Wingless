$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up100c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up100c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP100C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP100C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP100C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP100C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up100c-trigger8-generality -count=1;if($LASTEXITCODE-ne 0){throw 'UP100C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP100C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up100c-trigger8-generality)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP100C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up100c-trigger8-generality)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP100C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP100C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up100c-trigger8-generality.v1'){throw 'UP100C_SCHEMA_MISMATCH'}
 if([int]$R.trigger_threshold-ne 8){throw 'UP100C_TRIGGER'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UP100C_RECALL_CAP'}
 if($R.explicit_phase_label_used){throw 'UP100C_PHASE_LABEL'}
 if($R.future_oracle_used){throw 'UP100C_ORACLE'}
 if($R.query_labels_used_for_admission){throw 'UP100C_QUERY_LABEL'}
 if([int]$R.points.Count-ne 18){throw "UP100C_POINT_COUNT $($R.points.Count)"}
 foreach($M in $R.points){
  if([int]$M.recall_entries_used-gt 16){throw 'UP100C_RECALL_CAP_EXCEEDED'}
  if($M.arm-ceq 'trigger8' -and [int]$M.policy_metadata_bytes-ne 135){throw 'UP100C_TRIGGER_META'}
  if($M.arm-ceq 'continuous_control' -and [int]$M.policy_metadata_bytes-ne 134){throw 'UP100C_CONTROL_META'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "arm=$($M.arm) hot=$($M.hot_keys) order=$($M.initial_order) hot_acc=$($M.hot_query_accuracy) exact=$($M.hot_set_exact_accuracy) fp=$($M.false_positive_churn_admissions) trigger_eps=$($M.trigger_fired_episodes) trigger_idx=$($M.mean_filtered_write_index_at_trigger)"}
 Write-Host 'WINGLESS_UP157_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
