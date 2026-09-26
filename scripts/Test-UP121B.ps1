$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up121b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up121b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP121B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP121B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP121B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP121B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up121b-orthogonal-generalization -count=1;if($LASTEXITCODE-ne 0){throw 'UP121B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP121B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up121b-orthogonal-generalization)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP121B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up121b-orthogonal-generalization)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP121B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP121B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up121b-orthogonal-generalization.v1'){throw 'UP121B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UP121B_STATE_DIM'}
 if($R.correction_recomputed){throw 'UP121B_CORRECTION_RECOMPUTED'}
 if($R.adaptive_geometry){throw 'UP121B_ADAPTIVE_GEOMETRY'}
 if([int]$R.original_subjects.Count-ne 6 -or [int]$R.unseen_subjects.Count-ne 6){throw 'UP121B_SUBJECT_COUNT'}
 if([int]$R.points.Count-ne 336){throw "UP121B_POINT_COUNT_$($R.points.Count)"}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 $Final=@($R.points|Where-Object{[int]$_.stage-eq 6})
 foreach($M in $Final){Write-Host "arm=$($M.representation_arm) lexical=$($M.lexical_set) order=$($M.class_order) policy=$($M.replay_policy) stores_orig=$($M.stores_original_accuracy) stores_unseen=$($M.stores_unseen_accuracy) base_orig=$($M.base_original_accuracy) base_unseen=$($M.base_unseen_accuracy) acquired=$($M.acquired_aggregate_accuracy) unseen_so=$($M.mean_store_observe_margin_unseen) unseen_sr=$($M.mean_store_report_margin_unseen)"}
 Write-Host 'WINGLESS_UP163_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
