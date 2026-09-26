$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up124c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up124c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP124C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP124C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP124C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP124C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up124c-history-overflow-boundary -count=1;if($LASTEXITCODE-ne 0){throw 'UP124C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP124C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up124c-history-overflow-boundary)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP124C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up124c-history-overflow-boundary)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP124C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP124C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up124c-history-overflow-boundary.v1'){throw 'UP124C_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 4){throw 'UP124C_POINT_COUNT'}
 if([int]$R.admission_memory_bytes-ne 128 -or [int]$R.exact_recall_cap-ne 16 -or [int]$R.history_entries-ne 32){throw 'UP124C_MEMORY'}
 if($R.query_evidence_used_for_admission -or $R.semantic_priority_used -or $R.future_oracle_used -or $R.replacement_rule_added){throw 'UP124C_BOUNDARY_LEAK'}
 $Expected=@(0,1,4,8)
 for($i=0;$i-lt 4;$i++){
  if([int]$R.points[$i].overflow_qualified_candidates-ne $Expected[$i]){throw 'UP124C_PRESSURE_ORDER'}
  if([int]$R.points[$i].max_history_table_entries-gt 32){throw 'UP124C_HISTORY_CAP_EXCEEDED'}
  if([int]$R.points[$i].recall_entries_used-gt 16){throw 'UP124C_RECALL_CAP_EXCEEDED'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "q=$($M.overflow_qualified_candidates) overflow=$($M.overflow_episode_rate) msg=$($M.overflow_message) evaluated=$($M.valid_evaluated_episodes) history_max=$($M.max_history_table_entries) valid=$($M.valid_accuracy) exact=$($M.target16_exact_accuracy)"}
 Write-Host 'WINGLESS_UP124_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
