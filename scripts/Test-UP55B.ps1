$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP55B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP55B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP55B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP55B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up55b-margin-replication -count=1;if($LASTEXITCODE-ne 0){throw 'UP55B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP55B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up55b-margin-replication)|Out-String).Trim()
 $P2=((go run ./cmd/unitary-up55b-margin-replication)|Out-String).Trim()
 if($P1-cne $P2){throw 'UP55B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up55b-margin-replication.v1'){throw 'UP55B_SCHEMA_MISMATCH'}
 if($R.training_changed -or $R.commit_behavior_changed){throw 'UP55B_BEHAVIOR_CHANGED'}
 if([int]$R.points.Count-ne 12){throw 'UP55B_POINT_COUNT'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "All schedules high-margin > low-margin: $($R.all_schedules_separated)"
 foreach($P in $R.points){Write-Host "schedule=$($P.schedule_base) noise=$($P.memory_noise) writes=$($P.writes_per_scenario) commit=$($P.commit_accuracy) low=$($P.low_margin_accuracy) high=$($P.high_margin_accuracy) correct_margin=$($P.mean_margin_correct) error_margin=$($P.mean_margin_incorrect)"}
 Write-Host 'WINGLESS_UP55_HARNESS_PASS'
}finally{Pop-Location}
