$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up66b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up66b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP66B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP66B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP66B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP66B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up66b-confidence-carry-512 -count=1;if($LASTEXITCODE-ne 0){throw 'UP66B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP66B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up66b-confidence-carry-512)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP66B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up66b-confidence-carry-512)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP66B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP66B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up66b-confidence-carry-512.v1'){throw 'UP66B_SCHEMA_MISMATCH'}
 if($R.oracle_used){throw 'UP66B_ORACLE_PRESENT'}
 if($R.training_changed){throw 'UP66B_TRAINING_CHANGED'}
 if($R.threshold_selected){throw 'UP66B_THRESHOLD_SELECTION_PRESENT'}
 if([int]$R.metrics.Count-ne 64){throw "UP66B_METRIC_COUNT $($R.metrics.Count)"}
 $Depths=@($R.metrics|ForEach-Object{[int]$_.depth}|Sort-Object -Unique)
 if(($Depths -join ',') -cne '32,128,512,1024'){throw "UP66B_DEPTH_LEVELS $($Depths -join ',')"}
 foreach($M in $R.metrics){if([int]$M.metric.writes_per_scenario-ne 512){throw 'UP66B_WRITE_COUNT'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "depth=$($M.depth) schedule=$($M.metric.schedule_base) noise=$($M.metric.memory_noise) policy=$($M.metric.policy) threshold=$($M.metric.margin_threshold) commit=$($M.metric.commit_accuracy) final=$($M.metric.final_accuracy) relation=$($M.metric.relation_accuracy) low=$($M.metric.low_confidence_events)"}
 Write-Host 'WINGLESS_UP66_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
