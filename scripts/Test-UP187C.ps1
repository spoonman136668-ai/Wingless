$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP187C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP187C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up187c-native-risk-calibration -count=1;if($LASTEXITCODE-ne 0){throw 'UP187C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP187C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up187c-native-risk-calibration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP187C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up187c-native-risk-calibration)|Out-String).Trim();if($P1-cne $P2){throw 'UP187C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up187c-native-risk-calibration.v1' -or [int]$R.metrics.Count-ne 6 -or [int]$R.risk_buckets.Count-ne 3){throw 'UP187C_DESIGN'}
 if($R.intervention_used -or $R.adaptive_bucketing_used -or $R.warning_threshold_changed -or $R.live_activation){throw 'UP187C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "cadence=$($M.cadence) bucket=$($M.risk_bucket) checkpoints=$($M.checkpoints) one_fail=$($M.failures_within_one_cadence) one_rate=$($M.one_cadence_failure_rate) two_fail=$($M.failures_within_two_cadences) two_rate=$($M.two_cadence_failure_rate) mean_native_horizon=$($M.mean_adversarial_shield_horizon) mean_remaining=$($M.mean_writes_remaining_to_loss)"}
 Write-Host 'WINGLESS_UP232_HARNESS_PASS'
}finally{Pop-Location}
