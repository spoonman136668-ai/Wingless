$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up138c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up138c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP138C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP138C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP138C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP138C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up138c-counter-time-vs-event-time -count=1;if($LASTEXITCODE-ne 0){throw 'UP138C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP138C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up138c-counter-time-vs-event-time)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP138C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up138c-counter-time-vs-event-time)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP138C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP138C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up138c-counter-time-vs-event-time.v1'){throw 'UP138C_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 8 -or [int]$R.raw_operations_per_arm-ne 96){throw 'UP138C_DESIGN'}
 foreach($P in $R.points){
  if(([int]$P.candidate_writes+[int]$P.read_only_queries)-ne 96){throw 'UP138C_RAW_OPERATION_COUNT'}
  if([int]$P.candidate_boundaries-ne ([int]$P.candidate_writes/32)){throw 'UP138C_BOUNDARY_COUNT'}
 }
 if([int]$R.max_history_age-ne 2 -or [int]$R.exact_recall_cap-ne 16 -or [int]$R.history_entries-ne 32){throw 'UP138C_MEMORY'}
 if($R.query_operations_mutate_admission_state -or $R.semantic_priority_used -or $R.future_oracle_used -or $R.memory_increased){throw 'UP138C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "target=$($P.target) arm=$($P.arm) writes=$($P.candidate_writes) queries=$($P.read_only_queries) hist_present=$($P.history_present_before_pair_rate) hist_age=$($P.mean_history_age_before_pair) admit=$($P.admission_rate)"}
 Write-Host 'WINGLESS_UP138_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
