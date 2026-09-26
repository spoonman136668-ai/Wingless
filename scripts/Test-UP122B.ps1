$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up122b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up122b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP122B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP122B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP122B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP122B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up122b-alternate-failure-topology -count=1;if($LASTEXITCODE-ne 0){throw 'UP122B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP122B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up122b-alternate-failure-topology)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP122B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up122b-alternate-failure-topology)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP122B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP122B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up122b-alternate-failure-topology.v1'){throw 'UP122B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UP122B_STATE_DIM'}
 if($R.training_changed){throw 'UP122B_TRAINING_CHANGED'}
 if($R.adaptive_geometry){throw 'UP122B_ADAPTIVE_GEOMETRY'}
 if([int]$R.controls.Count-ne 84){throw "UP122B_CONTROL_COUNT_$($R.controls.Count)"}
 if([int]$R.diagnostics.Count-ne 1008){throw "UP122B_DIAGNOSTIC_COUNT_$($R.diagnostics.Count)"}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 $D=@($R.diagnostics|Where-Object{$_.acquired -and [double]$_.accuracy-lt 1})
 foreach($M in $D){Write-Host "order=$($M.class_order) policy=$($M.replay_policy) stage=$($M.stage) verb=$($M.verb) family=$($M.subject_family) acc=$($M.accuracy) mean_margin=$($M.mean_correct_class_margin) min_margin=$($M.min_correct_class_margin) wrongS=$($M.wrong_store_count) wrongO=$($M.wrong_observe_count) wrongR=$($M.wrong_report_count)"}
 Write-Host 'WINGLESS_UP166_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
