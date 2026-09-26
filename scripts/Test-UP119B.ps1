$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up119b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up119b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP119B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP119B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP119B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP119B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up119b-shifted-margin-topology -count=1;if($LASTEXITCODE-ne 0){throw 'UP119B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP119B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up119b-shifted-margin-topology)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP119B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up119b-shifted-margin-topology)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP119B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP119B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up119b-shifted-margin-topology.v1'){throw 'UP119B_SCHEMA_MISMATCH'}
 if($R.training_changed){throw 'UP119B_TRAINING_CHANGED'}
 if([int]$R.points.Count-ne 168){throw "UP119B_POINT_COUNT $($R.points.Count)"}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){if([double]$M.stores_accuracy-lt 1){Write-Host "arm=$($M.representation_arm) order=$($M.class_order) policy=$($M.replay_policy) stage=$($M.stage) acc=$($M.stores_accuracy) SO=$($M.mean_store_observe_margin) SR=$($M.mean_store_report_margin) wrongO=$($M.wrong_observe_count) wrongR=$($M.wrong_report_count)"}}
 Write-Host 'WINGLESS_UP161_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
