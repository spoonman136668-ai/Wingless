$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP186C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP186C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up186c-warning-lead-calibration -count=1;if($LASTEXITCODE-ne 0){throw 'UP186C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP186C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up186c-warning-lead-calibration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP186C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up186c-warning-lead-calibration)|Out-String).Trim();if($P1-cne $P2){throw 'UP186C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up186c-warning-lead-calibration.v1' -or [int]$R.metrics.Count-ne 2){throw 'UP186C_DESIGN'}
 if($R.intervention_used -or $R.warning_threshold_changed -or $R.live_activation){throw 'UP186C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "cadence=$($M.cadence) arms=$($M.arms) warned=$($M.warned_before_loss) lost=$($M.lost_by_64) min_lead=$($M.min_lead_writes) max_lead=$($M.max_lead_writes) mean_lead=$($M.mean_lead_writes) within1cad=$($M.lead_within_cadence) within2cad=$($M.lead_within_two_cadences) after_loss=$($M.warnings_after_loss)"}
 Write-Host 'WINGLESS_UP229_HARNESS_PASS'
}finally{Pop-Location}
