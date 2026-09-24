$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP58A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP58A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP58A_ICE_BUILD_FAILED'};go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP58A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up58a-exact-angle -count=1;if($LASTEXITCODE-ne 0){throw 'UP58A_FOCUSED_TEST_FAILED'};go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP58A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up58a-exact-angle)|Out-String).Trim();$P2=((go run ./cmd/unitary-up58a-exact-angle)|Out-String).Trim();if($P1-cne $P2){throw 'UP58A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json;if($R.schema-cne 'wingless.up58a-exact-angle-control.v1'){throw 'UP58A_SCHEMA_MISMATCH'};if($R.training_changed){throw 'UP58A_TRAINING_CHANGED'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ===';Write-Host "Exact restores all: $($R.exact_restores_all)";foreach($P in $R.points){Write-Host "length=$($P.length) learned=$($P.learned_accuracy) exact=$($P.exact_accuracy)"};Write-Host 'WINGLESS_UP58_HARNESS_PASS'
}finally{Pop-Location}
