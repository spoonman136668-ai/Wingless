$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP207B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP207B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up207b-native-offset-proxy -count=1;if($LASTEXITCODE-ne 0){throw 'UP207B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP207B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up207b-native-offset-proxy)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP207B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up207b-native-offset-proxy)|Out-String).Trim();if($P1-cne $P2){throw 'UP207B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up207b-native-offset-proxy.v1' -or [int]$R.training_points-ne 60 -or [int]$R.metrics.Count-ne 168 -or [int]$R.summaries.Count-ne 4){throw 'UP207B_DESIGN'}
 if($R.heldout_fitting_used -or $R.adaptive_feature_selection_used -or $R.phase_input_used -or $R.parity_input_used -or $R.heldout_anchor_measurement_used -or $R.maintenance_triggered){throw 'UP207B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "composition=$($S.composition) static_mae=$($S.static_mean_absolute_error) native_mae=$($S.native_offset_mean_absolute_error) oracle_mae=$($S.oracle_anchor_mean_absolute_error) native_vs_static=$($S.native_better_than_static)/$($S.evaluation_points) static_max=$($S.static_max_absolute_error) native_max=$($S.native_offset_max_absolute_error) oracle_max=$($S.oracle_anchor_max_absolute_error)"}
 Write-Host 'WINGLESS_UP231_HARNESS_PASS'
}finally{Pop-Location}
