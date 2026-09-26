$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up125c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up125c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP125C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP125C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP125C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP125C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up125c-age-eviction-overflow -count=1;if($LASTEXITCODE-ne 0){throw 'UP125C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP125C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up125c-age-eviction-overflow)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP125C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up125c-age-eviction-overflow)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP125C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP125C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up125c-age-eviction-overflow.v1'){throw 'UP125C_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 4){throw 'UP125C_POINT_COUNT'}
 if([int]$R.admission_memory_bytes-ne 128 -or [int]$R.exact_recall_cap-ne 16 -or [int]$R.history_entries-ne 32){throw 'UP125C_MEMORY'}
 if($R.semantic_priority_used -or $R.query_priority_used -or $R.future_oracle_used -or $R.memory_increased){throw 'UP125C_BOUNDARY_LEAK'}
 if($R.replacement_signal-cne 'maximum_history_age_only' -or $R.tie_break-cne 'lowest_table_index'){throw 'UP125C_RULE'}
 $Expected=@(0,1,8,16)
 for($i=0;$i-lt 4;$i++){
  if([int]$R.points[$i].overflow_qualified_candidates-ne $Expected[$i]){throw 'UP125C_PRESSURE_ORDER'}
  if([int]$R.points[$i].max_history_table_entries-gt 32 -or [int]$R.points[$i].recall_entries_used-gt 16){throw 'UP125C_CAP'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "q=$($M.overflow_qualified_candidates) panic=$($M.panic_rate) evict=$($M.total_evictions) age2=$($M.age2_evictions) age1=$($M.age1_evictions) valid=$($M.valid_accuracy) exact=$($M.target16_exact_accuracy) nonkeep=$($M.nonpersistent_retention_rate)"}
 Write-Host 'WINGLESS_UP125_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
