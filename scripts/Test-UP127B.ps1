$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up127b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up127b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP127B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP127B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP127B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP127B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up127b-balanced-store-gate -count=1;if($LASTEXITCODE-ne 0){throw 'UP127B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP127B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up127b-balanced-store-gate)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP127B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up127b-balanced-store-gate)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP127B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP127B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up127b-balanced-store-gate.v1'){throw 'UP127B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UP127B_STATE_DIM'}
 if([int]$R.gate_epochs-ne 20){throw 'UP127B_EPOCHS'}
 if([double]$R.gate_learning_rate-ne 0.08 -or [double]$R.gate_threshold-ne 0.5 -or [double]$R.positive_weight-ne 2){throw 'UP127B_PARAMS'}
 if($R.threshold_changed -or $R.encoder_changed){throw 'UP127B_BOUNDARY'}
 if([int]$R.split_metrics.Count-ne 6){throw "UP127B_SPLIT_COUNT $($R.split_metrics.Count)"}
 if([int]$R.surface_metrics.Count-ne 90){throw "UP127B_SURFACE_COUNT $($R.surface_metrics.Count)"}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.split_metrics){Write-Host "arm=$($M.arm) split=$($M.split) accuracy=$($M.accuracy) precision=$($M.precision) recall=$($M.recall)"}
 foreach($M in $R.surface_metrics){if($M.target_store){Write-Host "arm=$($M.arm) split=$($M.split) surface=$($M.surface) mean=$($M.mean_probability) positive=$($M.positive_rate)"}}
 Write-Host 'WINGLESS_UP169_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
