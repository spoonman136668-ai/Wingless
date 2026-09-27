$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm3z-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm3z-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM3Z_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM3Z_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm3z-broad-profile-pressure-transfer -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3Z_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3Z_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm3z-broad-profile-pressure-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3Z_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm3z-broad-profile-pressure-transfer)|Out-String).Trim();if($P1-cne $P2){throw 'UPLM3Z_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm3z-broad-profile-pressure-transfer.v1' -or [int]$R.metrics.Count-ne 2160 -or [int]$R.coordinate_summaries.Count-ne 81 -or [int]$R.profile_summaries.Count-ne 27){throw 'UPLM3Z_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_selection_used){throw 'UPLM3Z_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.profile_summaries){Write-Host "profile=$($S.deadline_profile) pressure=$($S.pressure_level) range=$($S.outcome_range) distinct=$($S.distinct_outcomes)"}
 foreach($S in $R.coordinate_summaries){Write-Host "profile=$($S.deadline_profile) pressure=$($S.pressure_level) coordinate=$($S.coordinate) nonzero=$($S.nonzero_spread_groups) max_spread=$($S.max_failure_spread)"}
 Write-Host 'WINGLESS_UP243_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
