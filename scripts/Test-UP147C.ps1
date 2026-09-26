$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up147c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up147c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP147C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP147C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP147C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP147C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up147c-durable-recency-gradient -count=1;if($LASTEXITCODE-ne 0){throw 'UP147C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP147C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up147c-durable-recency-gradient)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP147C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up147c-durable-recency-gradient)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP147C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP147C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up147c-durable-recency-gradient.v1'){throw 'UP147C_SCHEMA'}
 if([int]$R.points.Count-ne 16 -or [int]$R.cohorts-ne 4 -or [int]$R.refresh_arms-ne 4){throw 'UP147C_DESIGN'}
 if([int]$R.exact_recall_cap-ne 16 -or [int]$R.admissions_per_arm-ne 4 -or [int]$R.refresh_attempts_per_arm-ne 4){throw 'UP147C_BUDGET'}
 foreach($P in $R.points){if([int]$P.recall_entries_used-ne 16){throw 'UP147C_RECALL_COUNT'}}
 if($R.memory_policy_changed -or $R.adaptive_query_selection -or $R.capacity_increased){throw 'UP147C_BOUNDARY'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "cohort=$($P.cohort) refresh_before=$($P.refresh_before_admission) subsequent=$($P.subsequent_admissions) present_pre=$($P.present_before_refresh) refreshed=$($P.successful_refresh_queries) final=$($P.target_final_accuracy) hand=$($P.final_hand)"}
 Write-Host 'WINGLESS_UP147_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
