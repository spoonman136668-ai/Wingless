$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up124b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up124b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP124B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP124B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP124B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP124B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up124b-classwide-store-projector -count=1;if($LASTEXITCODE-ne 0){throw 'UP124B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP124B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up124b-classwide-store-projector)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP124B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up124b-classwide-store-projector)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP124B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP124B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up124b-classwide-store-projector.v1'){throw 'UP124B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UP124B_STATE_DIM'}
 if(-not $R.classwide_projector){throw 'UP124B_PROJECTOR_MISSING'}
 if($R.per_verb_recompute){throw 'UP124B_PER_VERB_RECOMPUTE'}
 if($R.adaptive_geometry){throw 'UP124B_ADAPTIVE_GEOMETRY'}
 if([int]$R.points.Count-ne 168){throw "UP124B_POINT_COUNT_$($R.points.Count)"}
 foreach($P in $R.points){if([int]$P.store_metrics.Count-ne 4){throw 'UP124B_STORE_METRIC_COUNT'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($A in @('local_corrections','classwide_store_projector')){
  $P=@($R.points|Where-Object{$_.arm-ceq $A})
  Write-Host "arm=$A min_base_orig=$(($P|Measure-Object base_original_accuracy -Minimum).Minimum) min_base_unseen=$(($P|Measure-Object base_unseen_accuracy -Minimum).Minimum) min_acquired=$(($P|Where-Object{[int]$_.stage-gt 0}|Measure-Object acquired_aggregate_accuracy -Minimum).Minimum)"
 }
 Write-Host 'WINGLESS_UP172_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
