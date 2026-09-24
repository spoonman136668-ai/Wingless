$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP55C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP55C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP55C_ICE_BUILD_FAILED'};go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP55C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up55c-tag-packing -count=1;if($LASTEXITCODE-ne 0){throw 'UP55C_FOCUSED_TEST_FAILED'};go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP55C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up55c-tag-packing)|Out-String).Trim();$P2=((go run ./cmd/unitary-up55c-tag-packing)|Out-String).Trim();if($P1-cne $P2){throw 'UP55C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json;if($R.schema-cne 'wingless.up55c-phase-tag-packing.v1'){throw 'UP55C_SCHEMA_MISMATCH'};if($R.selection_performed){throw 'UP55C_SELECTION_PRESENT'};if([int]$R.metrics.Count-ne 8){throw 'UP55C_METRIC_COUNT'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ===';foreach($M in $R.metrics){Write-Host "family=$($M.family) noise=$($M.noise) gate=$($M.gate) value=$($M.value_accuracy) exact=$($M.exact_scenario_accuracy) margin=$($M.minimum_margin)"};Write-Host 'WINGLESS_UP55_HARNESS_PASS'
}finally{Pop-Location}
