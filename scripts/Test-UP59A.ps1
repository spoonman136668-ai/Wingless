$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP59A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP59A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP59A_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP59A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up59a-angle-attribution -count=1;if($LASTEXITCODE-ne 0){throw 'UP59A_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP59A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up59a-angle-attribution)|Out-String).Trim()
 $P2=((go run ./cmd/unitary-up59a-angle-attribution)|Out-String).Trim()
 if($P1-cne $P2){throw 'UP59A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up59a-angle-component-attribution.v1'){throw 'UP59A_SCHEMA_MISMATCH'}
 if($R.training_changed){throw 'UP59A_TRAINING_CHANGED'}
 if([int]$R.metrics.Count-ne 16){throw 'UP59A_METRIC_COUNT'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "arm=$($M.arm) length=$($M.length) accuracy=$($M.accuracy) gate=$($M.gate)"}
 Write-Host 'WINGLESS_UP59_HARNESS_PASS'
}finally{Pop-Location}
