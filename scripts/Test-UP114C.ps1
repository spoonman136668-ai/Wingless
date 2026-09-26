$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up114c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up114c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP114C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP114C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP114C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP114C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up114c-overcapacity-recurrence -count=1;if($LASTEXITCODE-ne 0){throw 'UP114C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP114C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up114c-overcapacity-recurrence)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP114C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up114c-overcapacity-recurrence)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP114C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP114C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up114c-overcapacity-recurrence.v1'){throw 'UP114C_SCHEMA_MISMATCH'}
 if([int]$R.hot_keys-ne 12 -or [int]$R.exact_recall_cap-ne 16){throw 'UP114C_CAPACITY'}
 if([int]$R.repeat_gap-ne 16 -or [int]$R.generation_interval-ne 32){throw 'UP114C_WINDOW'}
 if([int]$R.admission_memory_bytes-ne 128 -or [int]$R.episodes_per_seed-ne 32){throw 'UP114C_CONFIG'}
 if($R.query_evidence_used_for_admission -or $R.priority_labels_used -or $R.future_oracle_used -or $R.phase_label_used){throw 'UP114C_LABEL_LEAK'}
 if([int]$R.points.Count-ne 3){throw "UP114C_POINT_COUNT $($R.points.Count)"}
 $Expected=@(4,5,8)
 for($i=0;$i-lt 3;$i++){
  if([int]$R.points[$i].recurring_demand-ne $Expected[$i]){throw 'UP114C_DEMAND_ORDER'}
  if([int]$R.points[$i].recall_entries_used-gt 16){throw 'UP114C_RECALL_CAP_EXCEEDED'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "demand=$($M.total_target_demand) admit=$($M.recurring_admission_rate) hot=$($M.hot_query_accuracy) recurring=$($M.recurring_accuracy) retained=$($M.mean_retained_recurring_count) target=$($M.target_set_accuracy) exact=$($M.whole_target_exact_accuracy) fp=$($M.one_shot_false_admissions)"}
 Write-Host 'WINGLESS_UP170_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
