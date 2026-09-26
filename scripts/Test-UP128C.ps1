$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up128c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up128c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP128C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP128C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP128C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP128C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up128c-recurrence-evidence-dose -count=1;if($LASTEXITCODE-ne 0){throw 'UP128C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP128C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up128c-recurrence-evidence-dose)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP128C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up128c-recurrence-evidence-dose)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP128C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP128C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up128c-recurrence-evidence-dose.v1'){throw 'UP128C_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 3 -or [int]$R.total_target_recurrence_sightings-ne 2){throw 'UP128C_DESIGN'}
 if([int]$R.admission_memory_bytes-ne 128 -or [int]$R.exact_recall_cap-ne 16 -or [int]$R.history_entries-ne 32){throw 'UP128C_MEMORY'}
 if($R.semantic_priority_used -or $R.query_priority_used -or $R.future_oracle_used -or $R.memory_increased -or $R.adaptive_timing_used){throw 'UP128C_BOUNDARY_LEAK'}
 for($i=0;$i-lt 3;$i++){
  if([int]$R.points[$i].pre_overflow_sightings-ne $i -or [int]$R.points[$i].post_overflow_sightings-ne (2-$i)){throw 'UP128C_TIMING_ORDER'}
  if([int]$R.points[$i].max_history_table_entries-gt 32 -or [int]$R.points[$i].recall_entries_used-gt 16){throw 'UP128C_CAP'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "arm=$($M.arm) second_evict=$($M.second_wave_evictions) valid_admit=$($M.valid_admission_rate) valid=$($M.valid_accuracy) exact=$($M.target16_exact_accuracy)"}
 Write-Host 'WINGLESS_UP128_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
