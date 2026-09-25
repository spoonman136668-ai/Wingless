$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path;$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up64b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up64b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null;$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP64B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP64B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP64B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP64B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up64b-depth-stratified-confidence-carry -count=1;if($LASTEXITCODE-ne 0){throw 'UP64B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP64B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up64b-depth-stratified-confidence-carry)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP64B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up64b-depth-stratified-confidence-carry)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP64B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP64B_NONDETERMINISTIC_OUTPUT'};$R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up64b-depth-stratified-confidence-carry.v1'){throw 'UP64B_SCHEMA_MISMATCH'}
 if($R.oracle_used){throw 'UP64B_ORACLE_PRESENT'};if($R.training_changed){throw 'UP64B_TRAINING_CHANGED'};if($R.threshold_selected){throw 'UP64B_THRESHOLD_SELECTED'}
 if([int]$R.metrics.Count-ne 64){throw "UP64B_METRIC_COUNT $($R.metrics.Count)"}
 $Depths=@($R.metrics|ForEach-Object{[int]$_.depth}|Sort-Object -Unique);if(($Depths -join ',') -cne '32,128,512,1024'){throw "UP64B_DEPTHS $($Depths -join ',')"}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($X in $R.metrics){$M=$X.metric;Write-Host "depth=$($X.depth) schedule=$($M.schedule_base) noise=$($M.memory_noise) policy=$($M.policy) threshold=$($M.margin_threshold) commit=$($M.commit_accuracy) final=$($M.final_accuracy) relation=$($M.relation_accuracy) low=$($M.low_confidence_events)"}
 Write-Host 'WINGLESS_UP64_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue}
