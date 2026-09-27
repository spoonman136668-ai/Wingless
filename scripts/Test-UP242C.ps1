$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false
 go clean -cache -testcache
 go test ./unitary ./cmd/unitary-up242c-active-cell-paired-transfer -count=1;if($LASTEXITCODE-ne 0){throw 'UP242C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP242C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up242c-active-cell-paired-transfer)|Out-String).Trim()
 $P2=((go run ./cmd/unitary-up242c-active-cell-paired-transfer)|Out-String).Trim()
 if($P1-cne $P2){throw 'UP242C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up242c-active-cell-paired-transfer.v1' -or [int]$R.paired_initial_states-ne 64){throw 'UP242C_DESIGN'}
 if($R.intervention_changed -or $R.threshold_fitting_used -or $R.classifier_training_used -or $R.adaptive_feature_selection_used -or $R.new_native_field_used -or $R.live_activation){throw 'UP242C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "pairs=$($R.paired_initial_states) both_fail=$($R.both_fail) a_only_fail=$($R.a_only_fail) b_only_fail=$($R.b_only_fail) both_survive=$($R.both_survive) both_event=$($R.both_event_present) current_match=$($R.exact_current_signature_matches) trajectory_match=$($R.exact_trajectory_signature_matches)"
 Write-Host 'WINGLESS_UP353_HARNESS_PASS'
}finally{Pop-Location}
