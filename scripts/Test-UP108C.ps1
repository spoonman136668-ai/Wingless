$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up108c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up108c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP108C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP108C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP108C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP108C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up108c-phasefree-reuse-checkpoint -count=1;if($LASTEXITCODE-ne 0){throw 'UP108C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP108C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up108c-phasefree-reuse-checkpoint)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP108C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up108c-phasefree-reuse-checkpoint)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP108C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP108C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up108c-phasefree-reuse-checkpoint.v1'){throw 'UP108C_SCHEMA_MISMATCH'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UP108C_RECALL_CAP'}
 if([int]$R.reuse_mask_bytes-ne 2){throw 'UP108C_MASK_BYTES'}
 if($R.phasefree_arm_uses_phase_label){throw 'UP108C_PHASE_LABEL'}
 if($R.query_labels_used_for_detector){throw 'UP108C_QUERY_LABEL'}
 if($R.future_oracle_used){throw 'UP108C_ORACLE'}
 if([int]$R.points.Count-ne 12){throw "UP108C_POINT_COUNT_$($R.points.Count)"}
 foreach($M in $R.points){
  if([int]$M.recall_entries_used-gt 16){throw 'UP108C_RECALL_CAP_EXCEEDED'}
  if($M.arm-ceq 'phasefree_reuse_checkpoint' -and [int]$M.policy_metadata_bytes-ne 136){throw 'UP108C_PHASEFREE_META'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "initial=$($M.initial_order) arm=$($M.arm) checkpoints=$($M.mean_checkpoints_per_episode) first=$($M.mean_first_checkpoint_filtered_write_index) resident=$($M.mean_resident_hot_before_churn) fp=$($M.false_positive_churn_admissions) final=$($M.final_hot_query_accuracy) exact=$($M.final_hot_set_exact_accuracy)"}
 Write-Host 'WINGLESS_UP173_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
