$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up101c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up101c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP101C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP101C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP101C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP101C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up101c-promotion-checkpoint -count=1;if($LASTEXITCODE-ne 0){throw 'UP101C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP101C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up101c-promotion-checkpoint)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP101C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up101c-promotion-checkpoint)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP101C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP101C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up101c-promotion-checkpoint.v1'){throw 'UP101C_SCHEMA_MISMATCH'}
 if([int]$R.hot_keys-ne 16 -or [int]$R.exact_recall_cap-ne 16){throw 'UP101C_CAP'}
 if($R.explicit_phase_label_used -or $R.future_oracle_used -or $R.query_labels_used_for_admission){throw 'UP101C_FORBIDDEN_SIGNAL'}
 if([int]$R.points.Count-ne 6){throw "UP101C_POINT_COUNT $($R.points.Count)"}
 foreach($M in $R.points){
  if([int]$M.checkpoints.Count-ne 3){throw 'UP101C_CHECKPOINT_COUNT'}
  if([int]$M.recall_entries_used-gt 16){throw 'UP101C_RECALL_CAP_EXCEEDED'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){foreach($C in $M.checkpoints){Write-Host "arm=$($M.arm) order=$($M.initial_order) checkpoint=$($C.name) resident_hot=$($C.mean_resident_hot_keys) hot_acc=$($C.hot_query_accuracy) exact=$($C.hot_set_exact_accuracy) resident_cold=$($C.mean_resident_cold_keys)"};Write-Host "repeat=$($M.mean_repeat_admissions) trigger_eps=$($M.trigger_fired_episodes) trigger_idx=$($M.mean_filtered_write_index_at_trigger) fp=$($M.false_positive_churn_admissions)"}
 Write-Host 'WINGLESS_UP160_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
