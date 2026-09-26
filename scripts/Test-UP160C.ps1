$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up160c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up160c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP160C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP160C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP160C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP160C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up160c-horizon16-geometry -count=1;if($LASTEXITCODE-ne 0){throw 'UP160C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP160C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up160c-horizon16-geometry)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP160C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up160c-horizon16-geometry)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP160C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP160C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up160c-horizon16-geometry.v1' -or [int]$R.primary_horizon-ne 16){throw 'UP160C_DESIGN'}
 if($R.future_schedule_used -or $R.semantic_class_used -or $R.intervention_policy_changed){throw 'UP160C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "cohort=$($P.cohort) hand=$($P.initial_hand) dist=$($P.cyclic_distance) eage=$($P.endangered_age) zeros=$($P.zero_age_slots) agesum=$($P.age_sum) durable=$($P.actual_durable)"}
 Write-Host 'WINGLESS_UP160_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
