$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up60b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up60b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP60B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP60B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP60B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP60B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up60b-confidence-noise-extension -count=1;if($LASTEXITCODE-ne 0){throw 'UP60B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP60B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up60b-confidence-noise-extension)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP60B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up60b-confidence-noise-extension)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP60B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP60B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up60b-confidence-noise-extension.v1'){throw 'UP60B_SCHEMA_MISMATCH'}
 if($R.oracle_used){throw 'UP60B_ORACLE_PRESENT'}
 if($R.training_changed){throw 'UP60B_TRAINING_CHANGED'}
 if($R.threshold_selected){throw 'UP60B_THRESHOLD_SELECTION_PRESENT'}
 if([int]$R.metrics.Count-ne 36){throw 'UP60B_METRIC_COUNT'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "schedule=$($M.schedule_base) noise=$($M.memory_noise) writes=$($M.writes_per_scenario) policy=$($M.policy) threshold=$($M.margin_threshold) commit=$($M.commit_accuracy) final=$($M.final_accuracy) relation=$($M.relation_accuracy)"}
 Write-Host 'WINGLESS_UP60_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
