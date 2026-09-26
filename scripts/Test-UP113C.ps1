$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up113c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up113c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP113C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP113C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP113C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP113C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up113c-exact-probation-mixedset -count=1;if($LASTEXITCODE-ne 0){throw 'UP113C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP113C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up113c-exact-probation-mixedset)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP113C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up113c-exact-probation-mixedset)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP113C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP113C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up113c-exact-probation-mixedset.v1'){throw 'UP113C_SCHEMA_MISMATCH'}
 if([int]$R.hot_keys-ne 12 -or [int]$R.recurring_keys-ne 4 -or [int]$R.target_working_set-ne 16){throw 'UP113C_WORKING_SET'}
 if([int]$R.generation_interval-ne 32 -or [int]$R.admission_memory_bytes-ne 128){throw 'UP113C_ADMISSION'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UP113C_RECALL_CAP'}
 if([int]$R.episodes_per_seed-ne 32){throw 'UP113C_EPISODES'}
 if($R.query_evidence_used_for_admission -or $R.future_oracle_used -or $R.phase_label_used){throw 'UP113C_LABEL_LEAK'}
 if([int]$R.points.Count-ne 12){throw "UP113C_POINT_COUNT $($R.points.Count)"}
 foreach($M in $R.points){if([int]$M.recall_entries_used-gt 16){throw 'UP113C_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "gap=$($M.repeat_gap) churn=$($M.churn_writes) admit=$($M.recurring_admission_rate) recurring=$($M.recurring_value_accuracy) hot=$($M.hot_query_accuracy) exact=$($M.mixed_set_exact_accuracy) fp=$($M.one_shot_false_admissions)"}
 Write-Host 'WINGLESS_UP167_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
