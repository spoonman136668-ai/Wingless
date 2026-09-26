$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up144c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up144c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP144C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP144C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP144C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP144C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up144c-live-age-stack-replacement -count=1;if($LASTEXITCODE-ne 0){throw 'UP144C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP144C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up144c-live-age-stack-replacement)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP144C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up144c-live-age-stack-replacement)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP144C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP144C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up144c-live-age-stack-replacement.v1'){throw 'UP144C_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 2 -or [int]$R.history_entries-ne 32 -or [int]$R.max_history_age-ne 2){throw 'UP144C_DESIGN'}
 if($R.direct_state_fixture -or $R.memory_increased -or $R.future_oracle_used){throw 'UP144C_BOUNDARY'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "assignment=$($P.assignment) age2ev=$($P.age2_evictions) old1_admit=$($P.oldest1_admission_rate) old2_admit=$($P.oldest2_admission_rate) young1=$($P.younger1_admission_rate) young2=$($P.younger2_admission_rate) hot=$($P.hot_accuracy)"}
 Write-Host 'WINGLESS_UP144_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
