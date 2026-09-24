$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP60A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP60A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP60A_ICE_BUILD_FAILED'};go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP60A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up60a-angle-tolerance -count=1;if($LASTEXITCODE-ne 0){throw 'UP60A_FOCUSED_TEST_FAILED'};go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP60A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up60a-angle-tolerance)|Out-String).Trim();$P2=((go run ./cmd/unitary-up60a-angle-tolerance)|Out-String).Trim();if($P1-cne $P2){throw 'UP60A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json;if($R.schema-cne 'wingless.up60a-double-angle-tolerance.v1'){throw 'UP60A_SCHEMA_MISMATCH'};if($R.training_changed){throw 'UP60A_TRAINING_CHANGED'};if([int]$R.points.Count-ne 30){throw 'UP60A_POINT_COUNT'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ===';foreach($P in $R.points){Write-Host "delta=$($P.double_control_delta) length=$($P.length) accuracy=$($P.accuracy) gate=$($P.gate)"};Write-Host 'WINGLESS_UP60_HARNESS_PASS'
}finally{Pop-Location}
