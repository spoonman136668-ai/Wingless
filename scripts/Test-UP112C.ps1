$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up112c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up112c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP112C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP112C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP112C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP112C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up112c-probation-recurrence-window -count=1;if($LASTEXITCODE-ne 0){throw 'UP112C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP112C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up112c-probation-recurrence-window)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP112C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up112c-probation-recurrence-window)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP112C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP112C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up112c-probation-recurrence-window.v1'){throw 'UP112C_SCHEMA_MISMATCH'}
 if([int]$R.generation_interval-ne 32 -or [int]$R.admission_memory_bytes-ne 128 -or [int]$R.counter_bytes-ne 1){throw 'UP112C_MEMORY_BOUND'}
 if([int]$R.episodes_per_cell-ne 64){throw 'UP112C_EPISODES'}
 if($R.exact_memory_layer_used -or $R.query_evidence_used -or $R.future_oracle_used){throw 'UP112C_BOUNDARY'}
 if([int]$R.points.Count-ne 16){throw "UP112C_POINT_COUNT $($R.points.Count)"}
 foreach($M in $R.points){
  if([int]$M.admission_memory_bytes-ne 128 -or [int]$M.counter_bytes-ne 1){throw 'UP112C_POINT_MEMORY'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "arm=$($M.arm) gap=$($M.repeat_after_writes) expected=$($M.expected_repeat_admission) repeat=$($M.repeat_admission_rate) correct=$($M.expected_decision_accuracy) one_shot_fp=$($M.one_shot_false_admissions) fp_rate=$($M.one_shot_false_admission_rate) clears=$($M.mean_generation_clears)"}
 Write-Host 'WINGLESS_UP173_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
