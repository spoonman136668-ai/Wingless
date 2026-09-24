$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP57A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP57A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP57A_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP57A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up57a-depth-ladder -count=1;if($LASTEXITCODE-ne 0){throw 'UP57A_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP57A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up57a-depth-ladder)|Out-String).Trim();$P2=((go run ./cmd/unitary-up57a-depth-ladder)|Out-String).Trim()
 if($P1-cne $P2){throw 'UP57A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up57a-double-control-depth-ladder.v1'){throw 'UP57A_SCHEMA_MISMATCH'}
 if($R.training_changed){throw 'UP57A_TRAINING_CHANGED'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "Last passing length: $($R.last_passing_length)";Write-Host "First failing length: $($R.first_failing_length)"
 foreach($P in $R.points){Write-Host "length=$($P.length) gate=$($P.gate) unitary=$($P.unitary_accuracy) control=$($P.nonunitary_accuracy)"}
 Write-Host 'WINGLESS_UP57_HARNESS_PASS'
}finally{Pop-Location}
