$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up129c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up129c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP129C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP129C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP129C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP129C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up129c-consolidation-window -count=1;if($LASTEXITCODE-ne 0){throw 'UP129C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP129C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up129c-consolidation-window)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP129C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up129c-consolidation-window)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP129C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP129C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up129c-consolidation-window.v1'){throw 'UP129C_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 3 -or [int]$R.first_wave_qualified-ne 16 -or [int]$R.second_wave_qualified-ne 4){throw 'UP129C_DESIGN'}
 if([int]$R.pre_overflow_target_sightings-ne 2 -or [int]$R.pre_overflow_writes-ne 64){throw 'UP129C_TIMING_BUDGET'}
 if([int]$R.admission_memory_bytes-ne 128 -or [int]$R.exact_recall_cap-ne 16 -or [int]$R.history_entries-ne 32){throw 'UP129C_MEMORY'}
 $Expected=@('adjacent_same_generation','wide_same_generation','boundary_split');for($i=0;$i-lt 3;$i++){if($R.points[$i].arm-cne $Expected[$i]){throw 'UP129C_ORDER'}}
 if($R.semantic_priority_used -or $R.query_priority_used -or $R.future_oracle_used -or $R.memory_increased){throw 'UP129C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "arm=$($M.arm) boundary=$($M.boundary_between_sightings) valid_admit=$($M.valid_admission_rate) valid=$($M.valid_accuracy) exact=$($M.target16_exact_accuracy) evict=$($M.post_first_wave_evictions)"}
 Write-Host 'WINGLESS_UP129_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
