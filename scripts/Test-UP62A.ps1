$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up62a-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up62a-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP62A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP62A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP62A_ICE_BUILD_FAILED'};go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP62A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up62a-init-robustness -count=1;if($LASTEXITCODE-ne 0){throw 'UP62A_FOCUSED_TEST_FAILED'};go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP62A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up62a-init-robustness)|Out-String).Trim();$P2=((go run ./cmd/unitary-up62a-init-robustness)|Out-String).Trim();if($P1-cne $P2){throw 'UP62A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json;if($R.schema-cne 'wingless.up62a-initialization-robustness.v1'){throw 'UP62A_SCHEMA_MISMATCH'};if([int]$R.points.Count-ne 5){throw 'UP62A_POINT_COUNT'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ===';foreach($P in $R.points){Write-Host "init=$($P.initialization -join ',') double_error=$($P.double_angle_error) long128=$($P.long128_accuracy) gate=$($P.gate)"};Write-Host 'WINGLESS_UP62_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue}
