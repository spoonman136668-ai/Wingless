$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm3l-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm3l-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM3L_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM3L_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM3L_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM3L_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm3l-temporal-resource-coordinate -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3L_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3L_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm3l-temporal-resource-coordinate)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3L_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm3l-temporal-resource-coordinate)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3L_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM3L_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm3l-temporal-resource-coordinate.v1'){throw 'UPLM3L_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 18 -or [int]$R.candidate_coordinates.Count-ne 4 -or [int]$R.coordinate_summaries.Count-ne 4){throw 'UPLM3L_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_grouping_used -or $R.adaptive_resources_used){throw 'UPLM3L_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.coordinate_summaries){Write-Host "coordinate=$($S.coordinate) groups=$($S.groups) nonzero=$($S.nonzero_spread_groups) max_spread=$($S.max_failure_spread) mean_spread=$($S.mean_failure_spread)"}
 Write-Host 'WINGLESS_UP204_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
