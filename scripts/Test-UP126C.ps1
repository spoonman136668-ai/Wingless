$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up126c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up126c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP126C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP126C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP126C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP126C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up126c-second-overflow-wave -count=1;if($LASTEXITCODE-ne 0){throw 'UP126C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP126C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up126c-second-overflow-wave)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP126C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up126c-second-overflow-wave)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP126C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP126C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up126c-second-overflow-wave.v1'){throw 'UP126C_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 4 -or [int]$R.first_wave_qualified-ne 16){throw 'UP126C_DESIGN'}
 if([int]$R.admission_memory_bytes-ne 128 -or [int]$R.exact_recall_cap-ne 16 -or [int]$R.history_entries-ne 32){throw 'UP126C_MEMORY'}
 if($R.semantic_priority_used -or $R.query_priority_used -or $R.future_oracle_used -or $R.memory_increased){throw 'UP126C_BOUNDARY_LEAK'}
 $Expected=@(0,1,4,8);for($i=0;$i-lt 4;$i++){if([int]$R.points[$i].second_wave_qualified-ne $Expected[$i]){throw 'UP126C_ORDER'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "q2=$($M.second_wave_qualified) first_evict=$($M.first_wave_evictions) second_evict=$($M.second_wave_evictions) valid_admit=$($M.valid_admission_rate) valid=$($M.valid_accuracy) exact=$($M.target16_exact_accuracy) age2=$($M.age2_evictions)"}
 Write-Host 'WINGLESS_UP126_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
