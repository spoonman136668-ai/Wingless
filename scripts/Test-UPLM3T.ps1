$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm3t-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm3t-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM3T_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM3T_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM3T_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM3T_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm3t-deadline-profile-transfer -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3T_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3T_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm3t-deadline-profile-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3T_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm3t-deadline-profile-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3T_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM3T_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm3t-deadline-profile-transfer.v1'){throw 'UPLM3T_SCHEMA_MISMATCH'}
 if([int]$R.pressure_level-ne 4 -or [int]$R.metrics.Count-ne 54 -or [int]$R.coordinate_summaries.Count-ne 6 -or [int]$R.profile_summaries.Count-ne 3){throw 'UPLM3T_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_profile_selection_used -or $R.adaptive_coordinate_search_used){throw 'UPLM3T_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.profile_summaries){Write-Host "profile=$($S.deadline_profile) global_min=$($S.min_earliest_failed) global_max=$($S.max_earliest_failed) range=$($S.outcome_range) distinct=$($S.distinct_outcomes)"}
 foreach($S in $R.coordinate_summaries){Write-Host "profile=$($S.deadline_profile) coordinate=$($S.coordinate) groups=$($S.groups) nonzero=$($S.nonzero_spread_groups) max_spread=$($S.max_failure_spread) mean_spread=$($S.mean_failure_spread)"}
 Write-Host 'WINGLESS_UP227_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
