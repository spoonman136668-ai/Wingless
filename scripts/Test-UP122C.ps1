$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up122c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up122c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP122C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP122C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP122C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP122C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up122c-history-occupancy-pressure -count=1;if($LASTEXITCODE-ne 0){throw 'UP122C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP122C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up122c-history-occupancy-pressure)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP122C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up122c-history-occupancy-pressure)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP122C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP122C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up122c-history-occupancy-pressure.v1'){throw 'UP122C_SCHEMA_MISMATCH'}
 if([int]$R.hot_keys-ne 12 -or [int]$R.valid_targets-ne 4){throw 'UP122C_TARGET_SET'}
 if([int]$R.admission_memory_bytes-ne 128 -or [int]$R.exact_recall_cap-ne 16 -or [int]$R.history_entries-ne 32 -or [int]$R.key_bits-ne 14){throw 'UP122C_MEMORY'}
 if([int]$R.churn_writes-ne 12288 -or [int]$R.episodes_per_seed-ne 32){throw 'UP122C_WORKLOAD'}
 if($R.query_evidence_used_for_admission -or $R.semantic_priority_used -or $R.future_oracle_used -or $R.adaptive_horizon_used){throw 'UP122C_LABEL_LEAK'}
 if([int]$R.points.Count-ne 4){throw "UP122C_POINT_COUNT $($R.points.Count)"}
 $Expected=@(0,4,8,12)
 for($i=0;$i-lt 4;$i++){
  if([int]$R.points[$i].ephemeral_per_generation-ne $Expected[$i]){throw 'UP122C_PRESSURE_ORDER'}
  if([int]$R.points[$i].recall_entries_used-gt 16){throw 'UP122C_RECALL_CAP_EXCEEDED'}
  if([int]$R.points[$i].max_history_table_entries-gt 32){throw 'UP122C_HISTORY_CAP_EXCEEDED'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "pressure=$($M.ephemeral_per_generation) admit=$($M.valid_admission_rate) target=$($M.valid_accuracy) exact=$($M.target16_exact_accuracy) eph_admit=$($M.ephemeral_admission_rate) eph_keep=$($M.ephemeral_retention_rate) current_max=$($M.max_current_table_entries) history_max=$($M.max_history_table_entries) fp=$($M.one_shot_false_admissions)"}
 Write-Host 'WINGLESS_UP188_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
