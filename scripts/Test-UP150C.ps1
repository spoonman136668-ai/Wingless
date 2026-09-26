$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up150c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up150c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP150C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP150C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP150C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP150C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up150c-eviction-horizon-map -count=1;if($LASTEXITCODE-ne 0){throw 'UP150C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP150C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up150c-eviction-horizon-map)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP150C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up150c-eviction-horizon-map)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP150C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP150C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up150c-eviction-horizon-map.v1'){throw 'UP150C_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 64 -or [int]$R.initial_hands.Count-ne 4 -or [int]$R.admissions_per_arm-ne 16){throw 'UP150C_DESIGN'}
 if([int]$R.durable_queries_after_initialization-ne 0 -or $R.memory_policy_changed -or $R.capacity_increased){throw 'UP150C_BOUNDARY_LEAK'}
 foreach($P in $R.points){if([int]$P.recall_entries_used-ne 16){throw 'UP150C_CAPACITY'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "hand=$($P.initial_hand) admission=$($P.admission) evicted=$($P.newly_evicted_original_key) remaining=$($P.original_facts_remaining) current_hand=$($P.current_hand)"}
 Write-Host 'WINGLESS_UP150_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
