$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up126b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up126b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP126B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP126B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP126B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP126B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up126b-store-gate-topology -count=1;if($LASTEXITCODE-ne 0){throw 'UP126B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP126B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up126b-store-gate-topology)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP126B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up126b-store-gate-topology)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP126B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP126B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up126b-store-gate-topology.v1'){throw 'UP126B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UP126B_STATE_DIM'}
 if($R.gate_retrained -or $R.threshold_changed -or $R.projector_used){throw 'UP126B_DIAGNOSTIC_MUTATION'}
 if([int]$R.surface_metrics.Count-ne 45){throw "UP126B_SURFACE_COUNT $($R.surface_metrics.Count)"}
 if([int]$R.split_metrics.Count-ne 3){throw 'UP126B_SPLIT_COUNT'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.split_metrics){Write-Host "split=$($M.split) accuracy=$($M.accuracy) precision=$($M.precision) recall=$($M.recall)"}
 foreach($M in $R.surface_metrics){Write-Host "split=$($M.split) surface=$($M.surface) store=$($M.target_store) mean=$($M.mean_probability) min=$($M.min_probability) max=$($M.max_probability) positive=$($M.positive_rate)"}
 Write-Host 'WINGLESS_UP166_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
