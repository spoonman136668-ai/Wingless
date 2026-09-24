$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up57b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up57b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP57B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP57B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP57B_ICE_BUILD_FAILED'};go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP57B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up57b-confidence-carry -count=1;if($LASTEXITCODE-ne 0){throw 'UP57B_FOCUSED_TEST_FAILED'};go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP57B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up57b-confidence-carry)|Out-String).Trim();$P2=((go run ./cmd/unitary-up57b-confidence-carry)|Out-String).Trim();if($P1-cne $P2){throw 'UP57B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json;if($R.schema-cne 'wingless.up57b-confidence-carry.v1'){throw 'UP57B_SCHEMA_MISMATCH'};if($R.oracle_used -or $R.training_changed){throw 'UP57B_BOUNDARY_VIOLATION'};if([int]$R.metrics.Count-ne 24){throw 'UP57B_METRIC_COUNT'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ===';foreach($M in $R.metrics){Write-Host "schedule=$($M.schedule_base) noise=$($M.memory_noise) policy=$($M.policy) threshold=$($M.margin_threshold) low=$($M.low_confidence_events) commit=$($M.commit_accuracy) final=$($M.final_accuracy) relation=$($M.relation_accuracy)"};Write-Host 'WINGLESS_UP57_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue}
