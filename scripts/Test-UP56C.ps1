$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP56C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP56C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP56C_ICE_BUILD_FAILED'};go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP56C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up56c-golden-confirmation -count=1;if($LASTEXITCODE-ne 0){throw 'UP56C_FOCUSED_TEST_FAILED'};go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP56C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up56c-golden-confirmation)|Out-String).Trim();$P2=((go run ./cmd/unitary-up56c-golden-confirmation)|Out-String).Trim();if($P1-cne $P2){throw 'UP56C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up56c-golden-confirmation.v1'){throw 'UP56C_SCHEMA_MISMATCH'}
 if($R.selection_performed){throw 'UP56C_SELECTION_PRESENT'}
 if([int]$R.points.Count-ne 18){throw 'UP56C_POINT_COUNT'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "schedule=$($P.schedule_base) noise=$($P.memory_noise) gate=$($P.gate) value=$($P.value_accuracy) exact=$($P.exact_scenario_accuracy) margin=$($P.minimum_margin)"}
 Write-Host 'WINGLESS_UP56_HARNESS_PASS'
}finally{Pop-Location}
