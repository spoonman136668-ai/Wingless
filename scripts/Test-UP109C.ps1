$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up109c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up109c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP109C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP109C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP109C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP109C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up109c-slot-reuse-confirmation -count=1;if($LASTEXITCODE-ne 0){throw 'UP109C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP109C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up109c-slot-reuse-confirmation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP109C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up109c-slot-reuse-confirmation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP109C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP109C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up109c-slot-reuse-confirmation.v1'){throw 'UP109C_SCHEMA_MISMATCH'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UP109C_RECALL_CAP'}
 if($R.phasefree_arm_uses_phase_label){throw 'UP109C_PHASE_LABEL'}
 if($R.query_labels_used_for_detector){throw 'UP109C_QUERY_LABEL'}
 if($R.future_oracle_used){throw 'UP109C_ORACLE'}
 if([int]$R.points.Count-ne 12){throw "UP109C_POINT_COUNT $($R.points.Count)"}
 foreach($M in $R.points){if([int]$M.recall_entries_used-gt 16){throw 'UP109C_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "arm=$($M.arm) order=$($M.initial_order) checkpoints=$($M.mean_checkpoints_per_episode) idx=$($M.mean_first_checkpoint_filtered_write_index) fp=$($M.false_positive_churn_admissions) exact=$($M.final_hot_set_exact_accuracy)"}
 Write-Host 'WINGLESS_UP163_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
