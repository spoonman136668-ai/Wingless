$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM3U_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM3U_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm3u-profile-pressure-map -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3U_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3U_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm3u-profile-pressure-map)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3U_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm3u-profile-pressure-map)|Out-String).Trim();if($P1-cne $P2){throw 'UPLM3U_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm3u-profile-pressure-map.v1' -or [int]$R.metrics.Count-ne 486 -or [int]$R.coordinate_summaries.Count-ne 54 -or [int]$R.profile_summaries.Count-ne 27){throw 'UPLM3U_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_selection_used){throw 'UPLM3U_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.profile_summaries){Write-Host "profile=$($S.deadline_profile) pressure=$($S.pressure_level) range=$($S.outcome_range) distinct=$($S.distinct_outcomes)"}
 foreach($S in $R.coordinate_summaries){Write-Host "profile=$($S.deadline_profile) pressure=$($S.pressure_level) coordinate=$($S.coordinate) nonzero=$($S.nonzero_spread_groups) max_spread=$($S.max_failure_spread)"}
 Write-Host 'WINGLESS_UP230_HARNESS_PASS'
}finally{Pop-Location}
