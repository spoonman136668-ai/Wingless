$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE
$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up56b-gocache-"+$PID)
$GoTmp=Join-Path $env:TEMP ("wingless-up56b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache
$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP56B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP56B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP56B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP56B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up56b-selective-risk -count=1;if($LASTEXITCODE-ne 0){throw 'UP56B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP56B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up56b-selective-risk)|Out-String).Trim()
 $P2=((go run ./cmd/unitary-up56b-selective-risk)|Out-String).Trim()
 if($P1-cne $P2){throw 'UP56B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up56b-selective-risk.v1'){throw 'UP56B_SCHEMA_MISMATCH'}
 if($R.training_changed -or $R.commit_behavior_changed){throw 'UP56B_BEHAVIOR_CHANGED'}
 if([int]$R.points.Count-ne 12){throw 'UP56B_POINT_COUNT'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){foreach($M in $P.metrics){Write-Host "schedule=$($P.schedule_base) noise=$($P.memory_noise) writes=$($P.writes_per_scenario) threshold=$($M.threshold) coverage=$($M.coverage) selected_accuracy=$($M.selected_accuracy) baseline=$($M.baseline_accuracy)"}}
 Write-Host 'WINGLESS_UP56_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
